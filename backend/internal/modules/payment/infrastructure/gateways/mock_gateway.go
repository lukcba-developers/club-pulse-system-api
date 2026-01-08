package gateways

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
)

type MockPaymentGateway struct {
	ShouldFail bool
}

func (m *MockPaymentGateway) CreatePreference(ctx context.Context, payment *domain.Payment, payerEmail string, description string) (string, error) {
	if m.ShouldFail {
		return "", fmt.Errorf("gateway failure")
	}
	return "https://sandbox.mercadopago.com.ar/checkout/v1/redirect/mock-pref-id", nil
}

func (m *MockPaymentGateway) ProcessWebhook(ctx context.Context, payload interface{}) (*domain.Payment, error) {
	// Mock returns mock payment
	return &domain.Payment{
		ExternalID: "mock-ext-id",
		Status:     domain.PaymentStatusCompleted,
	}, nil
}

func (m *MockPaymentGateway) ValidateWebhook(req *http.Request) error {
	// Always valid for mock
	return nil
}
