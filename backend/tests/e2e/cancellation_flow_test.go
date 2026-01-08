package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	bookingHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingCancellationFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clean state
	_ = db.Migrator().DropTable(&domain.Booking{}, &facilitiesRepo.FacilityModel{})
	_ = db.AutoMigrate(&domain.Booking{}, &facilitiesRepo.FacilityModel{})

	// Deps
	facRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	bookRepo := bookingRepo.NewPostgresBookingRepository(db)
	recRepo := bookingRepo.NewPostgresRecurringRepository(db)
	uRepo := userRepo.NewPostgresUserRepository(db)

	// Mock Notifier & Refund
	sharedMock := &SharedMockNotifier{}
	bookingUC := bookingApp.NewBookingUseCases(bookRepo, recRepo, facRepo, uRepo, sharedMock, sharedMock)
	h := bookingHttp.NewBookingHandler(bookingUC)

	r := gin.New()
	authMw := func(userID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("userID", userID)
			c.Set("clubID", "test-club-cancel")
			c.Next()
		}
	}

	facID := uuid.New()
	db.Create(&facilitiesRepo.FacilityModel{
		ID:     facID.String(),
		ClubID: "test-club-cancel",
		Name:   "Court Cancel",
		Status: "active",
	})

	userID := uuid.New()
	db.Create(&userRepo.UserModel{ID: userID.String(), ClubID: "test-club-cancel", Name: "U Cancel", Email: "cancel@test.com", MedicalCertStatus: "VALID"})

	bookingID := uuid.New()
	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(time.Hour)

	db.Create(&domain.Booking{
		ID:         bookingID,
		UserID:     userID,
		FacilityID: facID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     domain.BookingStatusConfirmed,
		ClubID:     "test-club-cancel",
	})

	t.Run("User Cancel Booking", func(t *testing.T) {
		// Route
		group := r.Group("/api/v1")
		// Register cancellation route
		bookingHttp.RegisterRoutes(group, h, authMw(userID.String()), func(c *gin.Context) { c.Next() })

		w := httptest.NewRecorder()
		// Endpoint: DELETE /bookings/:id
		req, _ := http.NewRequest("DELETE", "/api/v1/bookings/"+bookingID.String(), nil)
		r.ServeHTTP(w, req)

		// Assert
		require.True(t, w.Code == http.StatusOK || w.Code == http.StatusNoContent, "Expected 200 or 204, got %d. Body: %s", w.Code, w.Body.String())

		var b domain.Booking
		err := db.First(&b, "id = ?", bookingID).Error
		assert.NoError(t, err)
		assert.Equal(t, domain.BookingStatusCancelled, b.Status)
	})
}
