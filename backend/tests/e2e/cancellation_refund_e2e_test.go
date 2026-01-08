package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	bookingHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	paymentApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/application"
	paymentDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	paymentRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/repository"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCancellationRefundFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	db := SetupTestDB(t)

	// Ensure UserStats and Wallet tables exist
	_ = db.AutoMigrate(&userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{})

	// Repos
	bRepo := bookingRepo.NewPostgresBookingRepository(db)
	rRepo := bookingRepo.NewPostgresRecurringRepository(db)
	fRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	uRepo := userRepo.NewPostgresUserRepository(db)
	pRepo := paymentRepo.NewPostgresPaymentRepository(db)

	// Mocks
	notifier := &SharedMockNotifier{}
	recordingGw := &RecordingMockPaymentGateway{}

	// UseCases
	payUC := paymentApp.NewPaymentUseCases(pRepo, recordingGw)
	bookUC := bookingApp.NewBookingUseCases(bRepo, rRepo, fRepo, uRepo, notifier, payUC)

	// Register responder
	payUC.RegisterResponder("BOOKING", bookUC)

	h := bookingHttp.NewBookingHandler(bookUC)

	r := gin.New()
	userID := uuid.New().String()
	clubID := "test-club-refund"

	authMw := func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("userRole", "PLAYER")
		c.Next()
	}
	tenantMw := func(c *gin.Context) {
		c.Set("clubID", clubID)
		c.Next()
	}

	r.POST("/bookings", tenantMw, authMw, h.Create)
	r.DELETE("/bookings/:id", tenantMw, authMw, h.Cancel)

	// Create test user with valid medical certificate
	validStatus := "VALID"
	futureExpiry := time.Now().Add(365 * 24 * time.Hour)
	testUser := &userRepo.UserModel{
		ID:                userID,
		Name:              "Test User",
		Email:             "test-refund@test.com",
		Role:              "MEMBER",
		ClubID:            clubID,
		Password:          "$2a$10$placeholder",
		MedicalCertStatus: validStatus,
		MedicalCertExpiry: &futureExpiry,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	db.Create(testUser)

	// Create default Stats and Wallet
	db.Create(&userDomain.UserStats{UserID: userID, Level: 1})
	db.Create(&userDomain.Wallet{UserID: userID, Balance: 0})

	// Create test facility
	facID := uuid.New()
	// Create test facility using struct to ensure defaults and types are correct
	testFacility := &facilitiesRepo.FacilityModel{
		ID:          facID.String(),
		ClubID:      clubID,
		Name:        "Test Court",
		Type:        "CANCHA",
		Status:      "active",
		HourlyRate:  100.0,
		GuestFee:    0.0,
		OpeningHour: 8,
		ClosingHour: 22,
		Capacity:    10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(testFacility)

	// 2. Scenario: Create booking
	startTime := time.Now().Add(48 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(1 * time.Hour)

	body := `{"user_id": "` + userID + `", "facility_id": "` + facID.String() + `", "start_time": "` + startTime.Format(time.RFC3339) + `", "end_time": "` + endTime.Format(time.RFC3339) + `"}`
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/bookings", strings.NewReader(body))
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusCreated {
		t.Logf("Booking Failed Body: %s", w1.Body.String())
	}
	require.Equal(t, http.StatusCreated, w1.Code)

	var bookingCreated bookingDomain.Booking
	err := json.Unmarshal(w1.Body.Bytes(), &bookingCreated)
	require.NoError(t, err)
	bookingID := bookingCreated.ID

	// 3. Simulate Payment (Manual DB insert or use usecase)
	payment := &paymentDomain.Payment{
		ID:            uuid.New(),
		ClubID:        clubID,
		ReferenceID:   bookingID,
		ReferenceType: "BOOKING",
		Amount:        decimal.NewFromFloat(100),
		Status:        paymentDomain.PaymentStatusCompleted,
		ExternalID:    "mp-ext-12345",
	}
	err = pRepo.Create(context.Background(), payment)
	require.NoError(t, err)

	// Update booking to CONFIRMED (simulating what ProcessWebhook does)
	bookingCreated.Status = bookingDomain.BookingStatusConfirmed
	err = bRepo.Update(&bookingCreated)
	require.NoError(t, err)

	// 4. Cancel booking
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/bookings/"+bookingID.String(), nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// 5. Verify Refund was called
	assert.Contains(t, recordingGw.RefundCalledWith, "mp-ext-12345")

	// 6. Verify Payment status in DB
	updatedPayments, _, err := pRepo.List(context.Background(), clubID, paymentDomain.PaymentFilter{})
	require.NoError(t, err)

	found := false
	for _, p := range updatedPayments {
		if p.ReferenceID == bookingID {
			assert.Equal(t, paymentDomain.PaymentStatusRefunded, p.Status)
			found = true
		}
	}
	assert.True(t, found, "Payment for booking should be found and refunded")
}
