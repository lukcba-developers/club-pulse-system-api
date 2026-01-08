package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	facilityDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	facilityRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
)

func TestBookingPricing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clear DB
	db.Exec("TRUNCATE TABLE bookings CASCADE")
	db.Exec("TRUNCATE TABLE facilities CASCADE")
	db.Exec("TRUNCATE TABLE users CASCADE")

	// Dependencies
	bRepo := bookingRepo.NewPostgresBookingRepository(db)
	rRepo := bookingRepo.NewPostgresRecurringRepository(db)
	fRepo := facilityRepo.NewPostgresFacilityRepository(db)
	uRepo := userRepo.NewPostgresUserRepository(db)
	mockNotifier := &SharedMockNotifier{}

	useCases := bookingApp.NewBookingUseCases(bRepo, rRepo, fRepo, uRepo, mockNotifier, mockNotifier)
	handler := bookingHttp.NewBookingHandler(useCases)

	r := gin.New()
	authMw := func(c *gin.Context) {
		c.Set("userID", "11111111-1111-1111-1111-111111111111")
		c.Next()
	}
	tenantMw := func(c *gin.Context) {
		c.Set("clubID", "test-club-1")
		c.Next()
	}
	bookingHttp.RegisterRoutes(r.Group("/api/v1"), handler, authMw, tenantMw)

	// Setup Data
	clubID := "test-club-1"
	userID := "11111111-1111-1111-1111-111111111111"

	// 1. Create User
	validCert := userDomain.MedicalCertStatusValid // Create variable to take address
	now := time.Now()
	expiry := now.AddDate(1, 0, 0) // Valid for 1 year
	user := &userDomain.User{
		ID:                userID,
		ClubID:            clubID,
		Name:              "Test User",
		Email:             "test@test.com",
		MedicalCertStatus: &validCert,
		MedicalCertExpiry: &expiry,
	}
	err := uRepo.Create(user)
	assert.NoError(t, err)

	// 2. Create Facility with Rate and Guest Fee
	facID := uuid.New().String()
	facility := &facilityDomain.Facility{
		ID:          facID,
		ClubID:      clubID,
		Name:        "Tennis Court 1",
		Status:      facilityDomain.FacilityStatusActive,
		HourlyRate:  100.00,
		GuestFee:    50.00,
		OpeningHour: 8,
		ClosingHour: 22,
	}
	err = fRepo.Create(facility)
	assert.NoError(t, err)

	t.Run("Create Booking with Guests and Verify Price", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(2 * time.Hour) // 2 Hours

		// Guest Details: 2 Guests
		body := `{
			"facility_id": "` + facID + `",
			"user_id": "` + userID + `",
			"start_time": "` + startTime.Format(time.RFC3339) + `",
			"end_time": "` + endTime.Format(time.RFC3339) + `",
			"guest_details": [
				{"name": "Guest 1", "dni": "123"},
				{"name": "Guest 2", "dni": "456"}
			]
		}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/bookings", strings.NewReader(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify calculation in DB
		// Base: 2 hours * 100 = 200
		// Guests: 2 guests * 50 = 100
		// Total: 300

		// Need to fetch from DB to verify TotalPrice as response might not return it if struct tags/json didn't have it (we added it)
		// But we did add json tag "total_price"
		assert.Contains(t, w.Body.String(), `"total_price":300`)
	})
}
