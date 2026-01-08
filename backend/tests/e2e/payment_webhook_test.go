package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	paymentHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/http"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/require"
)

// Mock Repo to avoid complex DB setup for just handler testing
type mockPaymentRepo struct {
	domain.PaymentRepository
}

func (m *mockPaymentRepo) Create(ctx context.Context, p *domain.Payment) error { return nil }
func (m *mockPaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	return nil, nil
}
func (m *mockPaymentRepo) Update(ctx context.Context, p *domain.Payment) error { return nil }
func (m *mockPaymentRepo) GetByExternalID(ctx context.Context, ext string) (*domain.Payment, error) {
	return nil, nil
}
func (m *mockPaymentRepo) List(ctx context.Context, clubID string, filter domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	return nil, 0, nil
}

type mockGatewayStrict struct{}

func (m *mockGatewayStrict) CreatePreference(ctx context.Context, p *domain.Payment, e, d string) (string, error) {
	return "", nil
}
func (m *mockGatewayStrict) ProcessWebhook(ctx context.Context, pl interface{}) (*domain.Payment, error) {
	return nil, nil
}
func (m *mockGatewayStrict) ValidateWebhook(req *http.Request) error {
	if req.Header.Get("x-signature") == "" {
		return fmt.Errorf("missing signature")
	}
	return nil
}
func (m *mockGatewayStrict) Refund(ctx context.Context, externalID string) error {
	return nil
}

func TestPaymentWebhookSecurity(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB() // Init DB just in case dependencies need it

	mockG := &mockGatewayStrict{}
	mockR := &mockPaymentRepo{}
	useCases := application.NewPaymentUseCases(mockR, mockG)
	h := paymentHttp.NewPaymentHandler(useCases)

	r := gin.New()
	r.POST("/webhook", h.HandleWebhook)

	t.Run("Missing Signature Headers -> 403 Forbidden", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/webhook?type=payment&data.id=123", nil)
		// No headers
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Valid Signature -> 200 OK (logic wise)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/webhook?type=payment&data.id=123", nil)
		req.Header.Set("x-signature", "valid")
		req.Header.Set("x-request-id", "req-id")
		r.ServeHTTP(w, req)

		// It should pass validation (mock) and try to process
		require.Equal(t, http.StatusOK, w.Code)
	})
}
