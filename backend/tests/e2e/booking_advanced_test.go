package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	bookingHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockNotifier struct{}

func (m *MockNotifier) Send(ctx context.Context, n service.Notification) error {
	return nil
}

func TestBookingAdvancedFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	_ = db.Migrator().DropTable(&domain.Booking{}, &domain.Waitlist{}, &facilitiesRepo.FacilityModel{})
	_ = db.AutoMigrate(&domain.Booking{}, &domain.Waitlist{}, &facilitiesRepo.FacilityModel{}, &userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{})

	facRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	bookRepo := bookingRepo.NewPostgresBookingRepository(db)
	recRepo := bookingRepo.NewPostgresRecurringRepository(db)

	usRepo := userRepo.NewPostgresUserRepository(db)

	sharedMock := &SharedMockNotifier{}
	uc := bookingApp.NewBookingUseCases(bookRepo, recRepo, facRepo, usRepo, sharedMock, sharedMock)
	h := bookingHttp.NewBookingHandler(uc)

	// Middleware Helper (Removed as unused)

	// 2. Data Setup
	// Facility
	facID := uuid.New().String()
	db.Create(&facilitiesRepo.FacilityModel{
		ID:     facID,
		ClubID: "test-club-adv-booking",
		Name:   "Padel Court 1",
		Status: "active",
	})

	// Users
	user1ID := uuid.New().String()
	user2ID := uuid.New().String()
	db.Create(&userRepo.UserModel{ID: user1ID, ClubID: "test-club-adv-booking", Name: "U1", Email: "u1@test.com", MedicalCertStatus: "VALID"})
	db.Create(&userRepo.UserModel{ID: user2ID, ClubID: "test-club-adv-booking", Name: "U2", Email: "u2@test.com", MedicalCertStatus: "VALID"})

	// 3. Test: User 1 Books
	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(time.Hour)

	// Duplicate test removed

	// Better Router Setup
	r := gin.New()
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

		if w.Code != http.StatusCreated {
			t.Logf("Booking Failed Body: %s", w.Body.String())
		}
		require.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("User 2 Conflict", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":     user2ID,
			"facility_id": facID,
			"start_time":  startTime, // Same time
			"end_time":    endTime,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		req.Header.Set("X-User-ID", user2ID)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Logf("Conflict Failed Body: %s", w.Body.String())
		}
		require.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("User 2 Join Waitlist", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"resource_id": facID,
			"target_date": startTime,
			"user_id":     user2ID, // Required by binding, though overridden/verified by token
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/bookings/waitlist", bytes.NewBuffer(body))
		req.Header.Set("X-User-ID", user2ID)
		r.ServeHTTP(w, req)

		// Assuming /waitlist endpoint exists as per docs
		require.Equal(t, http.StatusCreated, w.Code)

		// Verify DB
		var entry domain.Waitlist
		result := db.Where("user_id = ? AND resource_id = ?", user2ID, facID).First(&entry)
		assert.NoError(t, result.Error)
	})
}
