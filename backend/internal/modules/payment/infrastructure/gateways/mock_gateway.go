package gateways

import (
	"context"
	"fmt"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
)

// MockGateway is a mock implementation for testing
type MockGateway struct{}

func NewMockGateway() *MockGateway {
	return &MockGateway{}
}

func (p *MockGateway) CreatePreference(ctx context.Context, payment *domain.Payment, payerEmail string, description string) (string, error) {
	// For mock, simply return a fake checkout URL with the Payment ID
	mockURL := fmt.Sprintf("https://www.mercadopago.com.ar/checkout/v1/redirect?pref_id=mock-%s", payment.ID.String())
	return mockURL, nil
}

func (p *MockGateway) ProcessWebhook(ctx context.Context, payload interface{}) (*domain.Payment, error) {
	// Simulate parsing a webhook payload and returning the updated payment status
	now := time.Now()
	// In a real Mock, we might inspect payload to decide success/fail
	return &domain.Payment{
		// ID would be matched from ExternalReference in payload ideally,
		// but here we just return a stub to be updated.
		// Logic in handler handles retrieval by ID usually.
		Status: domain.PaymentStatusCompleted,
		Method: domain.PaymentMethodMercadoPago,
		PaidAt: &now,
	}, nil
}
