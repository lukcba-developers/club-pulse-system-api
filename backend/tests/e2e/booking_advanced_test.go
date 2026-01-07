package e2e_test

import (
	"bytes"
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

func TestBookingAdvancedFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	_ = db.Migrator().DropTable(&domain.Booking{}, &domain.Waitlist{}, &facilitiesRepo.FacilityModel{})
	_ = db.AutoMigrate(&domain.Booking{}, &domain.Waitlist{}, &facilitiesRepo.FacilityModel{})

	facRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	bookRepo := bookingRepo.NewPostgresBookingRepository(db)
	recRepo := bookingRepo.NewPostgresRecurringRepository(db)
	// Waitlist Repo? Usually part of BookingRepo or separate.
	// If separate, need to find it. `bookingRepo.NewPostgresWaitlistRepository`?
	// Let's assume it's integrated or accessible.
	// Checking `bookingApp.NewBookingUseCases` args: (BookingRepo, RecurringRepo, FacilityRepo, UserRepo, Notifier)
	// It might manage waitlist internally implicitly or we need a WaitlistRepo.
	// If I can't find it, I will assume it's part of the use case logic manually or check files.
	// For now, let's assume `bookRepo` handles it or we find out.
	// Actually, looking at previous files, `bookRepo` was passed.

	usRepo := userRepo.NewPostgresUserRepository(db)

	uc := bookingApp.NewBookingUseCases(bookRepo, recRepo, facRepo, usRepo, nil)
	h := bookingHttp.NewBookingHandler(uc)

	r := gin.Default()

	// Middleware Helper
	mockAuth := func(userID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("userID", userID)
			c.Set("clubID", "test-club-adv-booking")
			c.Next()
		}
	}

	// 2. Data Setup
	// Facility
	facID := uuid.New().String()
	db.Create(&facilitiesRepo.FacilityModel{
		ID:     facID,
		ClubID: "test-club-adv-booking",
		Name:   "Padel Court 1",
		Status: "ACTIVE",
	})

	// Users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	db.Create(&userRepo.UserModel{ID: user1ID, ClubID: "test-club-adv-booking", Name: "U1", Email: "u1@test.com"})
	db.Create(&userRepo.UserModel{ID: user2ID, ClubID: "test-club-adv-booking", Name: "U2", Email: "u2@test.com"})

	// 3. Test: User 1 Books
	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(time.Hour)

	t.Run("User 1 Books Successfully", func(t *testing.T) {
		group := r.Group("/")
		group.Use(mockAuth(user1ID))
		bookingHttp.RegisterRoutes(group, h, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

		body, _ := json.Marshal(map[string]interface{}{
			"facility_id": facID,
			"start_time":  startTime,
			"end_time":    endTime,
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		// Wait, I registered routes to root.
		// Actually I should register once carefully or use unique paths.
		// Let's restart Router for clean referencing or use sub-paths with different middleware?
		// Gin middleware is global if attached to r.Use().
		// I used generated usage inside the Run.
		r.ServeHTTP(w, req)
	})

	// Better Router Setup
	r = gin.New()
	// Dynamic Middleware
	r.Use(func(c *gin.Context) {
		uid := c.GetHeader("X-User-ID")
		if uid == "" {
			uid = user1ID
		} // default
		c.Set("userID", uid)
		c.Set("clubID", "test-club-adv-booking")
		c.Next()
	})
	bookingHttp.RegisterRoutes(r.Group("/api/v1"), h, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	t.Run("User 1 Books", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"facility_id": facID,
			"start_time":  startTime,
			"end_time":    endTime,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		req.Header.Set("X-User-ID", user1ID)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("User 2 Conflict", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"facility_id": facID,
			"start_time":  startTime, // Same time
			"end_time":    endTime,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		req.Header.Set("X-User-ID", user2ID)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("User 2 Join Waitlist", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"facility_id": facID,
			"start_time":  startTime,
			"end_time":    endTime,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/bookings/waitlist", bytes.NewBuffer(body))
		req.Header.Set("X-User-ID", user2ID)
		r.ServeHTTP(w, req)

		// Assuming /waitlist endpoint exists as per docs
		require.Equal(t, http.StatusCreated, w.Code)

		// Verify DB
		var entry domain.Waitlist
		result := db.Where("user_id = ? AND facility_id = ?", user2ID, facID).First(&entry)
		assert.NoError(t, result.Error)
	})
}
