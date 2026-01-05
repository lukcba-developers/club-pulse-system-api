package e2e

import (
	"encoding/json"
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

func TestMultiTenantIsolation(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clean state
	_ = db.Migrator().DropTable(&facilitiesRepo.FacilityModel{}, &domain.Booking{}, &domain.RecurringRule{})
	_ = db.AutoMigrate(&facilitiesRepo.FacilityModel{}, &domain.Booking{}, &domain.RecurringRule{})

	// Clear PostgreSQL cached prepared statements after schema change
	db.Exec("DISCARD ALL")

	// Repos
	facRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	bookRepo := bookingRepo.NewPostgresBookingRepository(db)
	recRepo := bookingRepo.NewPostgresRecurringRepository(db)
	usRepo := userRepo.NewPostgresUserRepository(db)

	// Logic
	// Using Booking Module as proxy for isolation test (most data heavy)
	bookUC := bookingApp.NewBookingUseCases(bookRepo, recRepo, facRepo, usRepo, nil) // Notifier nil
	bookH := bookingHttp.NewBookingHandler(bookUC)

	r := gin.New()

	// Middleware factory to simulate different clubs
	createMiddleware := func(clubID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("clubID", clubID)
			c.Set("userID", "admin-user")
			c.Set("userRole", "ADMIN") // Use Admin role to access ListAll
			c.Next()
		}
	}

	// Route Groups
	clubA := r.Group("/clubA/api/v1")
	clubA.Use(createMiddleware("club-A"))
	bookingHttp.RegisterRoutes(clubA, bookH, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	clubB := r.Group("/clubB/api/v1")
	clubB.Use(createMiddleware("club-B"))
	bookingHttp.RegisterRoutes(clubB, bookH, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Create Data for Club A
	// Create Facility A
	facAID := uuid.New()
	db.Create(&facilitiesRepo.FacilityModel{
		ID:     facAID.String(),
		ClubID: "club-A",
		Name:   "Field A",
		Type:   "football",
		Status: "active",
	})

	userAID := uuid.New()
	userBID := uuid.New()

	// Create Users in DB to satisfy FKs (assuming Booking has FK to User)
	db.Create(&userRepo.UserModel{ID: userAID.String(), ClubID: "club-A", Name: "User A", Email: "usera@test.com", Role: "MEMBER"})
	// 3. Create Data for Club B
	db.Create(&userRepo.UserModel{ID: userBID.String(), ClubID: "club-B", Name: "User B", Email: "userb@test.com", Role: "MEMBER"})

	// Create Booking A
	startA := time.Now().Add(24 * time.Hour)
	bookA := &domain.Booking{
		ID:           uuid.New(),
		ClubID:       "club-A",
		FacilityID:   facAID,
		UserID:       userAID,
		StartTime:    startA,
		EndTime:      startA.Add(1 * time.Hour),
		Status:       domain.BookingStatusConfirmed,
		GuestDetails: []domain.GuestDetail{},
	}
	err := bookRepo.Create(bookA)
	require.NoError(t, err)

	// Create Facility B
	facBID := uuid.New()
	db.Create(&facilitiesRepo.FacilityModel{
		ID:     facBID.String(),
		ClubID: "club-B", // Different Club
		Name:   "Field B",
		Type:   "padel",
		Status: "active",
	})

	// Create Booking B
	startB := time.Now().Add(48 * time.Hour)
	bookB := &domain.Booking{
		ID:           uuid.New(),
		ClubID:       "club-B",
		FacilityID:   facBID, // Valid for this club
		UserID:       userBID,
		StartTime:    startB,
		EndTime:      startB.Add(1 * time.Hour),
		Status:       domain.BookingStatusConfirmed,
		GuestDetails: []domain.GuestDetail{},
	}
	err = bookRepo.Create(bookB)
	require.NoError(t, err)

	// 4. Verify Isolation: Club A should NOT see Booking B
	t.Run("Club A Isolation", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, _ = http.NewRequest("GET", "/clubA/api/v1/bookings", nil) // Keeping original endpoint but with admin role?
		// Wait, Handler List() does NOT check Admin. ListAll() does.
		// So we MUST use /bookings/all
		req, _ := http.NewRequest("GET", "/clubA/api/v1/bookings/all", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		// Handle potential nil or type mismatch if needed, but assuming valid response
		bookingsData, ok := resp["data"].([]interface{})
		if !ok {
			// fallback or fail
			bookingsData = []interface{}{}
		}

		foundA := false
		foundB := false
		for _, b := range bookingsData {
			bMap := b.(map[string]interface{})
			idStr, ok := bMap["id"].(string)
			if !ok {
				continue
			}
			if idStr == bookA.ID.String() {
				foundA = true
			}
			if idStr == bookB.ID.String() {
				foundB = true
			}
		}

		assert.True(t, foundA, "Club A should see Booking A")
		assert.False(t, foundB, "Club A should NOT see Booking B")
	})

	// 5. Verify Isolation: Club B should NOT see Booking A
	t.Run("Club B Isolation", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/clubB/api/v1/bookings/all", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		bookingsData, ok := resp["data"].([]interface{})
		if !ok {
			bookingsData = []interface{}{}
		}

		foundA := false
		foundB := false
		for _, b := range bookingsData {
			bMap := b.(map[string]interface{})
			idStr, ok := bMap["id"].(string)
			if !ok {
				continue
			}
			if idStr == bookA.ID.String() {
				foundA = true
			}
			if idStr == bookB.ID.String() {
				foundB = true
			}
		}

		assert.False(t, foundA, "Club B should NOT see Booking A")
		assert.True(t, foundB, "Club B should see Booking B")
	})
}
