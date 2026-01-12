package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	authRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	clubDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	clubHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/http"
	clubRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/repository"
	notifService "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock Providers for testing
type mockEmailProvider struct{}

func (m *mockEmailProvider) SendEmail(ctx context.Context, to, subject, body string) (*notifService.DeliveryResult, error) {
	return &notifService.DeliveryResult{Success: true}, nil
}

type mockSMSProvider struct{}

func (m *mockSMSProvider) SendSMS(ctx context.Context, to, message string) (*notifService.DeliveryResult, error) {
	return &notifService.DeliveryResult{Success: true}, nil
}

func TestSuperAdminAccess(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Ensure clean state
	_ = db.Migrator().DropTable(&userDomain.UserStats{}, &userDomain.Wallet{}, &userRepo.UserModel{}, &clubDomain.Club{})
	_ = db.AutoMigrate(&clubDomain.Club{}, &userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{})

	// Clear PostgreSQL cached prepared statements after schema change
	// 2. Setup repos

	// Repos & Services
	cRepo := clubRepo.NewPostgresClubRepository(db)
	notifier := notifService.NewNotificationService(&mockEmailProvider{}, &mockSMSProvider{})
	// NewClubUseCases now requires SponsorRepo, ClubRepo, NewsRepo, NotificationService
	cUC := clubApp.NewClubUseCases(cRepo, cRepo, cRepo, notifier)
	cHandler := clubHttp.NewClubHandler(cUC)

	aRepo := authRepo.NewPostgresAuthRepository(db)

	// Create Super Admin User
	saEmail := "super@admin.com"
	saName := "Super Admin"
	// Ensure cleanup
	existingSA, _ := aRepo.FindUserByEmail(context.Background(), saEmail, "system")
	if existingSA != nil {
		db.Unscoped().Delete(existingSA)
	}

	saUser := &userRepo.UserModel{
		ID:        "sa-id-123",
		Name:      saName,
		Email:     saEmail,
		Role:      string(userDomain.RoleSuperAdmin),
		ClubID:    "system", // Special system club
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Direct DB creation to bypass any normal registration strictness
	db.Create(saUser)

	// Create Normal User
	normEmail := "normal@user.com"
	existingNorm, _ := aRepo.FindUserByEmail(context.Background(), normEmail, "club-A")
	if existingNorm != nil {
		db.Unscoped().Delete(existingNorm)
	}
	normUser := &userRepo.UserModel{
		ID:        "norm-id-456",
		Name:      "Normal User",
		Email:     normEmail,
		Role:      string(userDomain.RoleMember),
		ClubID:    "club-A",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(normUser)

	// Setup Router
	r := gin.New()

	// Mock Auth Middleware based on "Authorization" header
	// Mock Auth Middleware based on "Authorization" header
	mockAuth := func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		switch tokenString {
		case "Bearer super-token":
			c.Set("userID", saUser.ID)
			c.Set("userRole", saUser.Role)
			c.Set("userClubID", saUser.ClubID)
		case "Bearer normal-token":
			c.Set("userID", normUser.ID)
			c.Set("userRole", normUser.Role)
			c.Set("userClubID", normUser.ClubID)
		}
		c.Next()
	}

	// Route Group (Protected by TenantMiddleware)
	// Even though /clubs is global, we apply TenantMiddleware to verify Bypass logic.
	// But logically, /clubs shouldn't enforce X-Club-ID strictly if we are super admin?
	// Our TenantMiddleware logic explicitly allows bypass if RoleSuperAdmin.

	api := r.Group("/api/v1")
	api.Use(mockAuth)
	api.Use(middleware.TenantMiddleware(cRepo))

	// Strict Auth Middleware to enforce Super Admin role
	strictAuth := func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != string(userDomain.RoleSuperAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}

	clubHttp.RegisterRoutes(api, cHandler, strictAuth, func(c *gin.Context) { c.Next() }) // Use strictAuth

	// 2. Test: Normal User Cannot Access /clubs (Forbidden)
	t.Run("Normal User Forbidden", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/admin/clubs", nil)
		req.Header.Set("Authorization", "Bearer normal-token")
		req.Header.Set("X-Club-ID", "club-A") // Even with valid club ID

		// Create mock club-A so middleware passes strict check
		_ = cRepo.Create(context.Background(), &clubDomain.Club{ID: "club-A", Name: "Club A", Status: "ACTIVE"}) // Error if exists, ignore

		r.ServeHTTP(w, req)

		// Should be 403 Forbidden by *Handler* check, not Middleware
		require.Equal(t, http.StatusForbidden, w.Code)
	})

	// 3. Test: Super Admin Can Create Club (Success)
	t.Run("Super Admin Create Club", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"name": "New Club", "slug": "new-club", "domain": "new.club.com"}`
		req, _ := http.NewRequest("POST", "/api/v1/admin/clubs", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer super-token")
		// No X-Club-ID header needed for Super Admin if Bypass works!

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var createdClub map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &createdClub)
		createdID := createdClub["id"].(string)

		// Verify DB
		club, _ := cRepo.GetByID(context.Background(), createdID)
		assert.NotNil(t, club)
		assert.Equal(t, "New Club", club.Name)
	})

	// 4. Test: Super Admin Can List Clubs
	t.Run("Super Admin List Clubs", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/admin/clubs", nil)
		req.Header.Set("Authorization", "Bearer super-token")

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		clubs := resp["data"].([]interface{})

		found := false
		for _, c := range clubs {
			clubMap := c.(map[string]interface{})
			if clubMap["name"] == "New Club" {
				found = true
				break
			}
		}
		assert.True(t, found)
	})
}
