package e2e_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authRepository "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	clubDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	clubRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/repository"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
)

func TestMultiTenantIsolation(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clean
	db.Exec("TRUNCATE TABLE bookings CASCADE")
	// Ensure UserStats and Wallet exist (Booking flow might access them)
	_ = db.AutoMigrate(&userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{}, &clubDomain.Club{})

	// Dependencies
	bookRepo := bookingRepo.NewPostgresBookingRepository(db)
	recRepo := bookingRepo.NewPostgresRecurringRepository(db)
	facRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	usRepo := userRepo.NewPostgresUserRepository(db)
	clRepo := clubRepo.NewPostgresClubRepository(db)

	// Add proper Auth migration (Fixes missing google_id column)
	authRepo := authRepository.NewPostgresAuthRepository(db)
	_ = authRepo.Migrate()

	sharedMock := &SharedMockNotifier{}
	bookUC := bookingApp.NewBookingUseCases(bookRepo, recRepo, facRepo, clRepo, usRepo, sharedMock, sharedMock)
	bookH := bookingHttp.NewBookingHandler(bookUC)

	r := gin.New()
	authMw := func(c *gin.Context) {
		c.Set("userID", uuid.New().String())
		c.Next()
	}
	tenantMw := func(clubID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("clubID", clubID)
			c.Next()
		}
	}

	// Create Club A
	db.Create(&clubDomain.Club{
		ID:       "club-a",
		Name:     "Club A",
		Timezone: "UTC",
	})
	// Create Club B
	db.Create(&clubDomain.Club{
		ID:       "club-b",
		Name:     "Club B",
		Timezone: "UTC",
	})

	// Route for Club A
	userIDClubA := uuid.New().String()
	r.POST("/a/bookings", tenantMw("club-a"), func(c *gin.Context) {
		c.Set("userID", userIDClubA)
		c.Next()
	}, bookH.Create)
	// Route for Club B
	r.GET("/b/bookings", tenantMw("club-b"), authMw, bookH.List)

	// Create test user with valid medical certificate for Club A
	validStatus := "VALID"
	futureExpiry := time.Now().Add(365 * 24 * time.Hour)
	db.Exec("DELETE FROM users WHERE email = ?", "test-multi@test.com")
	testUser := &userRepo.UserModel{
		ID:                   userIDClubA,
		Name:                 "Test User A",
		Email:                "test-multi@test.com",
		Role:                 "MEMBER",
		ClubID:               "club-a",
		Password:             "$2a$10$placeholder",
		MedicalCertStatus:    validStatus,
		MedicalCertExpiry:    &futureExpiry,
		TermsAcceptedAt:      &futureExpiry, // Just needs to be non-nil
		PrivacyPolicyVersion: "2026-01",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	db.Create(testUser)

	// Create default Stats and Wallet
	db.Create(&userDomain.UserStats{UserID: userIDClubA, Level: 1})
	db.Create(&userDomain.Wallet{UserID: userIDClubA, Balance: 0})

	// Create test facility for Club A
	facID := uuid.New()
	db.Exec("DELETE FROM facilities WHERE id = ?", facID)
	// Create test facility using struct to ensure defaults and types are correct
	testFacility := &facilitiesRepo.FacilityModel{
		ID:          facID.String(),
		ClubID:      "club-a",
		Name:        "Test Court A",
		Type:        "CANCHA",
		Status:      "active",
		HourlyRate:  0.0,
		GuestFee:    0.0,
		OpeningTime: "08:00",
		ClosingTime: "22:00",
		Capacity:    10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(testFacility)

	// 2. Scenario: Create in Club A, should NOT see in Club B
	startTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	endTime := time.Now().Add(25 * time.Hour).Format(time.RFC3339)
	body := fmt.Sprintf(`{"facility_id": "%s", "start_time": "%s", "end_time": "%s"}`, facID.String(), startTime, endTime)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/a/bookings", strings.NewReader(body))
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/b/bookings", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Response body should have 0 bookings for Club B
	assert.Equal(t, `{"data":[]}`, w2.Body.String())
}
