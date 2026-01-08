package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
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

	// Dependencies
	bookRepo := bookingRepo.NewPostgresBookingRepository(db)
	recRepo := bookingRepo.NewPostgresRecurringRepository(db)
	facRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	usRepo := userRepo.NewPostgresUserRepository(db)

	sharedMock := &SharedMockNotifier{}
	bookUC := bookingApp.NewBookingUseCases(bookRepo, recRepo, facRepo, usRepo, sharedMock, sharedMock)
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

	// Route for Club A
	r.POST("/a/bookings", tenantMw("club-a"), authMw, bookH.Create)
	// Route for Club B
	r.GET("/b/bookings", tenantMw("club-b"), authMw, bookH.List)

	// 2. Scenario: Create in Club A, should NOT see in Club B
	facID := uuid.New().String()
	body := `{"facility_id": "` + facID + `", "start_time": "2025-01-01T10:00:00Z", "end_time": "2025-01-01T11:00:00Z"}`

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/a/bookings", strings.NewReader(body))
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/b/bookings", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Response body should have 0 bookings for Club B
	assert.Equal(t, "[]", w2.Body.String())
}
