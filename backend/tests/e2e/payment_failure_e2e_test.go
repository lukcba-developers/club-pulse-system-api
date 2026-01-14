package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	paymentApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/application"
	paymentDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	paymentHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/http"
	paymentRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/repository"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// FailureMockPaymentGateway that allows us to simulate different statuses
type FailureMockPaymentGateway struct {
	StatusToReturn string
}

func (m *FailureMockPaymentGateway) CreatePreference(ctx context.Context, p *paymentDomain.Payment, email, desc string) (string, error) {
	return "http://mock-mp.com/pay/" + p.ID.String(), nil
}

func (m *FailureMockPaymentGateway) ProcessWebhook(ctx context.Context, payload interface{}) (*paymentDomain.Payment, error) {
	// Payload is ID
	idStr := fmt.Sprintf("%v", payload)

	// We assume external ref is the UUID in our DB
	// For testing, we mock the return
	status := paymentDomain.PaymentStatusPending
	if m.StatusToReturn == "rejected" {
		status = paymentDomain.PaymentStatusFailed
	} else if m.StatusToReturn == "approved" {
		status = paymentDomain.PaymentStatusCompleted
	}

	return &paymentDomain.Payment{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"), // This needs to match the test record
		Status:     status,
		ExternalID: idStr,
	}, nil
}

func (m *FailureMockPaymentGateway) ValidateWebhook(req *http.Request) error {
	return nil // Skip validation for test simplicity
}

func (m *FailureMockPaymentGateway) Refund(ctx context.Context, externalID string) error {
	return nil
}

func TestPaymentFailureFlow(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clean tables
	_ = db.Migrator().DropTable(&paymentDomain.Payment{}, &bookingDomain.Booking{})
	_ = db.AutoMigrate(&paymentDomain.Payment{}, &bookingDomain.Booking{})

	// Wire Payment
	repo := paymentRepo.NewPostgresPaymentRepository(db)
	mockGateway := &FailureMockPaymentGateway{StatusToReturn: "rejected"}
	uc := paymentApp.NewPaymentUseCases(repo, mockGateway)

	// Wire Booking Responder
	bRepo := bookingRepo.NewPostgresBookingRepository(db)
	rRepo := bookingRepo.NewPostgresRecurringRepository(db)
	fRepo := facilitiesRepo.NewPostgresFacilityRepository(db)
	uRepo := userRepo.NewPostgresUserRepository(db)
	// Mock Notifier & Refund
	sharedMock := &SharedMockNotifier{}
	bookingUC := bookingApp.NewBookingUseCases(bRepo, rRepo, fRepo, uRepo, sharedMock, sharedMock)
	uc.RegisterResponder("BOOKING", bookingUC)

	h := paymentHttp.NewPaymentHandler(uc)

	r := gin.New()
	clubID := "test-club-payment"

	// Data Setup
	paymentID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	bookingID := uuid.New()

	// Create Booking
	db.Create(&bookingDomain.Booking{
		ID:     bookingID,
		ClubID: clubID,
		Status: bookingDomain.BookingStatusPendingPayment,
	})

	// Create Payment record
	db.Create(&paymentDomain.Payment{
		ID:            paymentID,
		ClubID:        clubID,
		ReferenceID:   bookingID,
		ReferenceType: "BOOKING",
		Status:        paymentDomain.PaymentStatusPending,
		Amount:        decimal.NewFromFloat(100.0),
		PayerID:       uuid.New(),
		Method:        paymentDomain.PaymentMethodMercadoPago,
		ExternalID:    "12345",
	})

	t.Run("Webhook Rejected Update Payment Status", func(t *testing.T) {
		group := r.Group("/api/v1")
		paymentHttp.RegisterRoutes(group, h, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

		w := httptest.NewRecorder()
		// MP Webhook format: POST /webhook?type=payment&data.id=12345
		req, _ := http.NewRequest("POST", "/api/v1/payments/webhook?type=payment&data.id=12345", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify Payment in DB
		var p paymentDomain.Payment
		db.First(&p, "id = ?", paymentID)
		assert.Equal(t, paymentDomain.PaymentStatusFailed, p.Status)

		// Verify Booking in DB (should NOT be confirmed)
		var b bookingDomain.Booking
		db.First(&b, "id = ?", bookingID)
		assert.Equal(t, bookingDomain.BookingStatusPendingPayment, b.Status)
	})

	t.Run("Webhook Approved Updates Booking to Confirmed", func(t *testing.T) {
		mockGateway.StatusToReturn = "approved"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/payments/webhook?type=payment&data.id=12345", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify Booking in DB
		var b bookingDomain.Booking
		db.First(&b, "id = ?", bookingID)
		assert.Equal(t, bookingDomain.BookingStatusConfirmed, b.Status)
		assert.Nil(t, b.PaymentExpiry)
	})
}
