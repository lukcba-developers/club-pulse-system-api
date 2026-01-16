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
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	bookingHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	clubDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	clubRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/repository"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
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
	_ = db.Migrator().DropTable(&bookingDomain.Booking{}, &facilitiesRepo.FacilityModel{}, &userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{}, &clubDomain.Club{})
	_ = db.AutoMigrate(&userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{}, &facilitiesRepo.FacilityModel{}, &bookingDomain.Booking{}, &clubDomain.Club{})

	// Deps
	fRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	bRepo := bookingRepo.NewPostgresBookingRepository(db)
	rRepo := bookingRepo.NewPostgresRecurringRepository(db)
	uRepo := userRepo.NewPostgresUserRepository(db)
	cRepo := clubRepo.NewPostgresClubRepository(db)
	mockNotifier := &SharedMockNotifier{}

	useCases := bookingApp.NewBookingUseCases(bRepo, rRepo, fRepo, cRepo, uRepo, mockNotifier, mockNotifier)
	h := bookingHttp.NewBookingHandler(useCases)

	r := gin.New()
	authMw := func(userID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("userID", userID)
			c.Set("clubID", "test-club-cancel")
			c.Next()
		}
	}

	clubID := "test-club-cancel"
	facID := uuid.New()
	db.Create(&facilitiesRepo.FacilityModel{
		ID:     facID.String(),
		ClubID: clubID,
		Name:   "Court Cancel",
		Status: "active",
	})

	userID := uuid.New()
	// 0. Create Club
	db.Create(&clubDomain.Club{
		ID:       clubID,
		Name:     "Cancellation Club",
		Timezone: "UTC",
	})

	// 1. Create User
	db.Create(&userRepo.UserModel{ID: userID.String(), ClubID: clubID, Name: "U Cancel", Email: "cancel@test.com", MedicalCertStatus: "VALID"})

	bookingID := uuid.New()
	startTime := time.Now().Add(48 * time.Hour).Truncate(24 * time.Hour).Add(10 * time.Hour)
	endTime := startTime.Add(time.Hour)

	db.Create(&bookingDomain.Booking{
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
