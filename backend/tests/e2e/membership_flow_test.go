package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	authRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	authToken "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/token"
	membershipApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	membershipHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/http"
	membershipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMembershipFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	_ = db.Migrator().DropTable(&membershipDomain.Membership{}, &membershipDomain.MembershipTier{}, &userDomain.UserStats{}, &userDomain.Wallet{}, &userRepo.UserModel{})
	_ = db.AutoMigrate(&membershipDomain.Membership{}, &membershipDomain.MembershipTier{}, &userRepo.UserModel{})

	// Repos
	memRepo := membershipRepo.NewPostgresMembershipRepository(db)
	scholarshipRepo := membershipRepo.NewPostgresScholarshipRepository(db)
	subRepo := membershipRepo.NewPostgresSubscriptionRepository(db)

	// We need User UseCase / Service to create user or just use Auth helper
	authR := authRepo.NewPostgresAuthRepository(db)
	_ = authR.Migrate()
	tokenService := authToken.NewJWTService("secret")
	authUC := authApp.NewAuthUseCases(authR, tokenService, nil)

	// Membership Logic
	memUC := membershipApp.NewMembershipUseCases(memRepo, scholarshipRepo, subRepo) // Correct repo
	memH := membershipHttp.NewMembershipHandler(memUC)

	r := gin.Default()

	// Middleware factory
	r.Use(func(c *gin.Context) {
		// Default to a user context if not overridden
		if _, exists := c.Get("userID"); !exists {
			c.Set("userID", "ignored") // overridden in sub-tests
		}
		c.Set("clubID", "test-club-membership")
		c.Next()
	})

	membershipHttp.RegisterRoutes(r.Group("/api/v1"), memH, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Create User
	email := "mem_test_" + uuid.New().String() + "@example.com"
	_, err := authUC.Register(context.Background(), authApp.RegisterDTO{
		Name:                 "Mem User",
		Email:                email,
		Password:             "Password123!",
		AcceptTerms:          true,
		PrivacyPolicyVersion: "2026-01",
	}, "test-club-membership")
	require.NoError(t, err)
	user, _ := authR.FindUserByEmail(context.Background(), email, "test-club-membership")
	userID := user.ID

	// 3. Create Tier (Directly via DB or Handler if exists)
	// Usually admin creates tier. Let's do DB for speed.
	tierID := uuid.New()
	tier := &membershipDomain.MembershipTier{
		ID:          tierID,
		ClubID:      "test-club-membership",
		Name:        "Gold",
		MonthlyFee:  decimal.NewFromFloat(100.0),
		Description: "Gold Tier",
	}
	db.Create(tier)

	// 4. Test: Purchase Membership
	// 4. Test: Purchase (Skipped incomplete setup, moving to re-setup below)

	// Re-setup with adjustable context
	r = gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("clubID", "test-club-membership")
		c.Set("userRole", "MEMBER")
		c.Next()
	})
	membershipHttp.RegisterRoutes(r.Group("/api/v1"), memH, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	t.Run("Purchase Success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":            userID,
			"membership_tier_id": tierID,
			"billing_cycle":      "MONTHLY",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/memberships", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
	})

	// 5. Test: Get My Memberships
	t.Run("Get My Memberships", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/memberships", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Gold")
	})

	// 6. Test: Process Billing (Admin)
	// Needs new router/middleware for Admin
	rAdmin := gin.New()
	rAdmin.Use(func(c *gin.Context) {
		c.Set("userID", "admin-id")
		c.Set("clubID", "test-club-membership")
		c.Set("userRole", "ADMIN")
		c.Next()
	})
	membershipHttp.RegisterRoutes(rAdmin.Group("/api/v1"), memH, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 6. Test: Assign Scholarship (Admin)
	t.Run("Assign Scholarship", func(t *testing.T) {
		// Grant 50% scholarship
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":    userID,
			"percentage": 0.5,
			"reason":     "Talented athlete",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/memberships/scholarship", bytes.NewBuffer(body))
		rAdmin.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
	})

	// 7. Test: Process Billing (Admin) - Should apply scholarship
	t.Run("Process Billing with Scholarship", func(t *testing.T) {
		// Force billing due by moving NextBillingDate to past (Yesterday)
		// This ensures that when +1 month is added, the new date is ~1 month from now (definitely in future)
		uid, _ := uuid.Parse(userID)

		err := db.Model(&membershipDomain.Membership{}).
			Where("user_id = ?", uid).
			Update("next_billing_date", time.Now().AddDate(0, 0, -1)).Error
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/memberships/process-billing", nil)
		rAdmin.ServeHTTP(w, req)

		// This might fail if the endpoint logic requires complex setup (dates, etc)
		// But we expect at least 200 OK
		require.Equal(t, http.StatusOK, w.Code)

		// Check if balance updated?
		// Check Membership NextBillingDate updated?
		var mem membershipDomain.Membership
		db.Where("user_id = ?", userID).First(&mem)
		// Assuming billing moved NextBillingDate
		assert.True(t, mem.NextBillingDate.After(time.Now()))

		// Balance check: Original Fee 100. Scholarship 50%. Expected Charge 50.
		// NOTE: Assuming this runs after Purchase (bal=0) + Scholarship + Billing.
		// If balance logic is cumulative, we expect 50.
		expectedBalance := decimal.NewFromFloat(50.0)
		assert.True(t, mem.OutstandingBalance.Equal(expectedBalance), "Balance should be 50 (100 - 50%)")
	})
}
