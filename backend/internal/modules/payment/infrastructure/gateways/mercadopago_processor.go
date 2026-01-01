package gateways

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
)

// PaymentProcessor defines the contract for payment gateways
type PaymentProcessor interface {
	CreatePreference(ctx context.Context, payment *domain.Payment) (string, error)
	ProcessWebhook(ctx context.Context, payload interface{}) (*domain.Payment, error)
}

// MercadoPagoMockProcessor is a mock implementation for local dev
type MercadoPagoMockProcessor struct{}

func NewMercadoPagoMockProcessor() *MercadoPagoMockProcessor {
	return &MercadoPagoMockProcessor{}
}

func (p *MercadoPagoMockProcessor) CreatePreference(ctx context.Context, payment *domain.Payment) (string, error) {
	// In a real implementation, we would call MP API here
	// For mock, we simply return a fake checkout URL
	mockURL := fmt.Sprintf("https://www.mercadopago.com.ar/checkout/v1/redirect?pref_id=mock-%s", payment.ID.String())
	return mockURL, nil
}

func (p *MercadoPagoMockProcessor) ProcessWebhook(ctx context.Context, payload interface{}) (*domain.Payment, error) {
	// Simulate parsing a webhook payload and returning the updated payment status
	// This is a simplified mock

	now := time.Now()
	return &domain.Payment{
		ID:        uuid.New(), // In real scenario extracting ID from payload
		Status:    domain.PaymentStatusCompleted,
		Method:    domain.PaymentMethodMercadoPago,
		PaidAt:    &now,
		UpdatedAt: now,
	}, nil
}
