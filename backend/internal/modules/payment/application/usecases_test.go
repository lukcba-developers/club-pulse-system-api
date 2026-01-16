package application_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepo) Update(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Payment, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByExternalID(ctx context.Context, clubID string, externalID string) (*domain.Payment, error) {
	args := m.Called(ctx, clubID, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByExternalIDForWebhook(ctx context.Context, externalID string) (*domain.Payment, error) {
	args := m.Called(ctx, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) List(ctx context.Context, clubID string, filter domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	args := m.Called(ctx, clubID, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Payment), int64(0), args.Error(2)
}

type MockPaymentGateway struct {
	mock.Mock
}

func (m *MockPaymentGateway) CreatePreference(ctx context.Context, payment *domain.Payment, payerEmail string, description string) (string, error) {
	args := m.Called(ctx, payment, payerEmail, description)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentGateway) ProcessWebhook(ctx context.Context, payload interface{}) (*domain.Payment, error) {
	args := m.Called(ctx, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentGateway) ValidateWebhook(req *http.Request) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockPaymentGateway) Refund(ctx context.Context, externalID string) error {
	args := m.Called(ctx, externalID)
	return args.Error(0)
}

type MockPaymentResponder struct {
	mock.Mock
}

func (m *MockPaymentResponder) OnPaymentStatusChanged(ctx context.Context, clubID string, referenceID uuid.UUID, status domain.PaymentStatus) error {
	args := m.Called(ctx, clubID, referenceID, status)
	return args.Error(0)
}

func TestPaymentUseCases_Checkout(t *testing.T) {
	repo := new(MockPaymentRepo)
	gateway := new(MockPaymentGateway)
	uc := application.NewPaymentUseCases(repo, gateway)
	ctx := context.TODO()

	req := application.CheckoutRequest{
		Amount:        "100.50",
		Description:   "Booking 123",
		PayerEmail:    "test@user.com",
		ReferenceID:   uuid.New(),
		ReferenceType: "BOOKING",
		UserID:        uuid.New(),
		ClubID:        "club-1",
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("Create", ctx, mock.Anything).Return(nil).Once()

		gateway.On("CreatePreference", ctx, mock.Anything, "test@user.com", "Booking 123").Return("http://checkout.url", nil).Once()

		payment, url, err := uc.Checkout(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "http://checkout.url", url)
		assert.NotNil(t, payment)
	})

	t.Run("Invalid Amount", func(t *testing.T) {
		reqInvalid := req
		reqInvalid.Amount = "invalid"
		_, _, err := uc.Checkout(ctx, reqInvalid)
		assert.Error(t, err)
	})

	t.Run("Repo Error", func(t *testing.T) {
		repo.On("Create", ctx, mock.Anything).Return(errors.New("db error")).Once()
		_, _, err := uc.Checkout(ctx, req)
		assert.Error(t, err)
	})
}

func TestPaymentUseCases_ProcessWebhook(t *testing.T) {
	repo := new(MockPaymentRepo)
	gateway := new(MockPaymentGateway)
	uc := application.NewPaymentUseCases(repo, gateway)
	ctx := context.TODO()

	// Register a responder
	responder := new(MockPaymentResponder)
	uc.RegisterResponder("BOOKING", responder)

	t.Run("Success - Payment Approved", func(t *testing.T) {
		payload := application.ProcessWebhookRequest{Type: "payment", DataID: "ext-123"}

		updatedTime := time.Now()
		gatewayPayment := &domain.Payment{
			ExternalID: "ext-123",
			Status:     domain.PaymentStatusCompleted,
			PaidAt:     &updatedTime,
			ID:         uuid.New(), // Match existing ID to simulate valid match
		}

		existingPayment := &domain.Payment{
			ID:            gatewayPayment.ID,
			ClubID:        "club-1",
			Status:        domain.PaymentStatusPending,
			ExternalID:    "ext-123",
			ReferenceID:   uuid.New(),
			ReferenceType: "BOOKING",
		}

		gateway.On("ProcessWebhook", ctx, "ext-123").Return(gatewayPayment, nil).Once()
		repo.On("GetByExternalIDForWebhook", ctx, "ext-123").Return(existingPayment, nil).Once()
		repo.On("Update", ctx, mock.MatchedBy(func(p *domain.Payment) bool {
			return p.Status == domain.PaymentStatusCompleted
		})).Return(nil).Once()
		responder.On("OnPaymentStatusChanged", ctx, "club-1", existingPayment.ReferenceID, domain.PaymentStatusCompleted).Return(nil).Once()

		res, err := uc.ProcessWebhook(ctx, payload)
		assert.NoError(t, err)
		assert.True(t, res.Processed)
		assert.Equal(t, domain.PaymentStatusCompleted, res.NewStatus)
	})

	t.Run("Security - Club ID Missing in Existing", func(t *testing.T) {
		payload := application.ProcessWebhookRequest{Type: "payment", DataID: "ext-bad"}
		gatewayPayment := &domain.Payment{ExternalID: "ext-bad", Status: domain.PaymentStatusCompleted, ID: uuid.New()}
		existingPayment := &domain.Payment{
			ID:         gatewayPayment.ID,
			ClubID:     "", // Vulnerability
			ExternalID: "ext-bad",
		}

		gateway.On("ProcessWebhook", ctx, "ext-bad").Return(gatewayPayment, nil).Once()
		repo.On("GetByExternalIDForWebhook", ctx, "ext-bad").Return(existingPayment, nil).Once()

		_, err := uc.ProcessWebhook(ctx, payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment missing club_id")
	})

	t.Run("Security - Payment ID Mismatch", func(t *testing.T) {
		payload := application.ProcessWebhookRequest{Type: "payment", DataID: "ext-mismatch"}
		// Gateway says this update is for payment ID X
		gatewayPayment := &domain.Payment{ExternalID: "ext-mismatch", Status: domain.PaymentStatusCompleted, ID: uuid.New()}
		// DB says external ID maps to payment ID Y (attacker trying to mix up?)
		existingPayment := &domain.Payment{
			ID:         uuid.New(), // Different ID
			ClubID:     "club-1",
			ExternalID: "ext-mismatch",
		}

		gateway.On("ProcessWebhook", ctx, "ext-mismatch").Return(gatewayPayment, nil).Once()
		repo.On("GetByExternalIDForWebhook", ctx, "ext-mismatch").Return(existingPayment, nil).Once()

		_, err := uc.ProcessWebhook(ctx, payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment ID mismatch")
	})
}

func TestPaymentUseCases_Refund(t *testing.T) {
	repo := new(MockPaymentRepo)
	gateway := new(MockPaymentGateway)
	uc := application.NewPaymentUseCases(repo, gateway)
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		refID := uuid.New()
		payment := &domain.Payment{
			ID:            uuid.New(),
			ReferenceID:   refID,
			ReferenceType: "BOOKING",
			ExternalID:    "ext-123",
			Status:        domain.PaymentStatusCompleted,
		}

		repo.On("List", ctx, "club-1", mock.MatchedBy(func(f domain.PaymentFilter) bool {
			return f.Status == domain.PaymentStatusCompleted
		})).Return([]*domain.Payment{payment}, int64(1), nil).Once()

		gateway.On("Refund", ctx, "ext-123").Return(nil).Once()
		repo.On("Update", ctx, mock.MatchedBy(func(p *domain.Payment) bool {
			return p.Status == domain.PaymentStatusRefunded
		})).Return(nil).Once()

		err := uc.Refund(ctx, "club-1", refID, "BOOKING")
		assert.NoError(t, err)
	})
}
