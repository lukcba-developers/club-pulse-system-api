package domain

import (
	"context"
	"net/http"
)

type PaymentGateway interface {
	// CreatePreference creates a payment preference and returns the checkout URL (init_point)
	CreatePreference(ctx context.Context, payment *Payment, payerEmail string, description string) (string, error)

	// ProcessWebhook handles the notification from the provider
	ProcessWebhook(ctx context.Context, payload interface{}) (*Payment, error)

	// ValidateWebhook verifies the authenticity of the webhook request
	ValidateWebhook(req *http.Request) error

	// Refund reverses a payment
	Refund(ctx context.Context, externalID string) error
}
