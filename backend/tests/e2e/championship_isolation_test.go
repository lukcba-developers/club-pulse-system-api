package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	championshipApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	championshipHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/http"
	championshipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/repository"
	championshipSvc "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/service"
	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	clubRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock mocks for notification
type MockEmailProvider2 struct{}

func (m *MockEmailProvider2) SendEmail(ctx context.Context, to string, subject, body string) (*service.DeliveryResult, error) {
	return &service.DeliveryResult{}, nil
}

type MockSMSProvider2 struct{}

func (m *MockSMSProvider2) SendSMS(ctx context.Context, to string, body string) (*service.DeliveryResult, error) {
	return &service.DeliveryResult{}, nil
}

func TestChampionshipIsolation(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clean/Migrate
	_ = db.Migrator().DropTable(&domain.Tournament{}, &domain.TournamentStage{}, &domain.Group{}, &domain.TournamentMatch{}, &domain.Standing{})
	_ = db.AutoMigrate(&domain.Tournament{}, &domain.TournamentStage{}, &domain.Group{}, &domain.TournamentMatch{}, &domain.Standing{})

	db.Exec("DISCARD ALL")

	// Dependencies
	champRepo := championshipRepo.NewPostgresChampionshipRepository(db)
	cRepo := clubRepo.NewPostgresClubRepository(db)

	// Create Clubs in DB
	clubAId := uuid.New()
	clubBId := uuid.New()

	notifier := service.NewNotificationService(&MockEmailProvider2{}, &MockSMSProvider2{})
	clubUC := clubApp.NewClubUseCases(cRepo, cRepo, cRepo, notifier)

	bookingAdapter := championshipSvc.NewChampionshipBookingAdapter(nil)
	champUC := championshipApp.NewChampionshipUseCases(champRepo, bookingAdapter, nil)

	// Handler
	handler := championshipHttp.NewChampionshipHandler(champUC, nil, clubUC)

	handleMiddleware := func(clubID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("clubID", clubID)
			c.Next()
		}
	}

	r := gin.New()

	// Admin Routes (using Context)
	adminA := r.Group("/admin/clubA")
	adminA.Use(handleMiddleware(clubAId.String()))
	handler.RegisterRoutes(adminA, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	adminB := r.Group("/admin/clubB")
	adminB.Use(handleMiddleware(clubBId.String()))
	handler.RegisterRoutes(adminB, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Data Creation
	t1 := &domain.Tournament{
		ID:     uuid.New(),
		ClubID: clubAId,
		Name:   "Tournament A",
		Sport:  "FUTBOL",
		Status: domain.TournamentActive,
	}
	require.NoError(t, champRepo.CreateTournament(t1))

	t2 := &domain.Tournament{
		ID:     uuid.New(),
		ClubID: clubBId,
		Name:   "Tournament B",
		Sport:  "TENNIS",
		Status: domain.TournamentActive,
	}
	require.NoError(t, champRepo.CreateTournament(t2))

	// 3. Verify Isolation
	t.Run("Club A Admin sees only T1", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Corrected URL with trailing slash (handler registers "/")
		req, _ := http.NewRequest("GET", "/admin/clubA/championships/", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)

		// Unmarshal into slice, not map["data"]
		var data []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &data)

		foundA := false
		foundB := false
		for _, m := range data {
			if m["id"] == t1.ID.String() {
				foundA = true
			}
			if m["id"] == t2.ID.String() {
				foundB = true
			}
		}
		assert.True(t, foundA, "Club A should see T1")
		assert.False(t, foundB, "Club A should NOT see T2")
	})

	t.Run("Club B Admin sees only T2", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/clubB/championships/", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)

		var data []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &data)

		foundA := false
		foundB := false
		for _, m := range data {
			if m["id"] == t1.ID.String() {
				foundA = true
			}
			if m["id"] == t2.ID.String() {
				foundB = true
			}
		}
		assert.False(t, foundA, "Club B should NOT see T1")
		assert.True(t, foundB, "Club B should see T2")
	})
}
