package application

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/shopspring/decimal"
)

// PaymentUseCases contains all payment-related business logic.
// Following Clean Architecture: handlers only handle HTTP, use cases contain logic.
type PaymentUseCases struct {
	repo    domain.PaymentRepository
	gateway domain.PaymentGateway
}

// NewPaymentUseCases creates a new PaymentUseCases instance.
func NewPaymentUseCases(repo domain.PaymentRepository, gateway domain.PaymentGateway) *PaymentUseCases {
	return &PaymentUseCases{
		repo:    repo,
		gateway: gateway,
	}
}

// CheckoutRequest represents the input for creating a checkout session.
type CheckoutRequest struct {
	Amount        float64
	Description   string
	PayerEmail    string
	ReferenceID   uuid.UUID
	ReferenceType string
	UserID        uuid.UUID
	ClubID        string
}

// Checkout creates a payment record and returns a payment gateway URL.
func (uc *PaymentUseCases) Checkout(ctx context.Context, req CheckoutRequest) (string, error) {
	payment := &domain.Payment{
		ID:            uuid.New(),
		Amount:        decimal.NewFromFloat(req.Amount),
		Currency:      "ARS",
		Status:        domain.PaymentStatusPending,
		Method:        domain.PaymentMethodMercadoPago,
		PayerID:       req.UserID,
		ClubID:        req.ClubID,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
	}

	if err := uc.repo.Create(ctx, payment); err != nil {
		log.Printf("Failed to create payment: %v", err)
		return "", errors.New("failed to create payment record")
	}

	url, err := uc.gateway.CreatePreference(ctx, payment, req.PayerEmail, req.Description)
	if err != nil {
		log.Printf("Gateway Error: %v", err)
		return "", errors.New("failed to contact payment gateway")
	}

	return url, nil
}

// ProcessWebhookRequest contains parsed webhook data.
type ProcessWebhookRequest struct {
	Type   string // "payment" for payment notifications
	DataID string // External payment ID from provider
}

// WebhookResult represents the outcome of webhook processing.
type WebhookResult struct {
	Processed bool
	PaymentID uuid.UUID
	NewStatus domain.PaymentStatus
	NotFound  bool
}

// ValidateWebhook validates the webhook signature.
func (uc *PaymentUseCases) ValidateWebhook(req *http.Request) error {
	return uc.gateway.ValidateWebhook(req)
}

// ProcessWebhook handles webhook notifications from payment providers.
// This contains the business logic previously in the handler.
func (uc *PaymentUseCases) ProcessWebhook(ctx context.Context, webhookReq ProcessWebhookRequest) (*WebhookResult, error) {
	result := &WebhookResult{}

	if webhookReq.Type != "payment" {
		// Not a payment webhook, acknowledge but don't process
		return result, nil
	}

	if webhookReq.DataID == "" {
		return result, nil
	}

	// 1. Get payment info from gateway
	updatedPayment, err := uc.gateway.ProcessWebhook(ctx, webhookReq.DataID)
	if err != nil {
		log.Printf("Webhook processing failed (gateway): %v", err)
		return nil, errors.New("gateway processing failed")
	}

	if updatedPayment == nil {
		return result, nil
	}

	// 2. Find existing payment in our DB
	existing, err := uc.repo.GetByID(ctx, updatedPayment.ID)
	if err != nil {
		log.Printf("Payment not found for update: %s", updatedPayment.ID)
		result.NotFound = true
		return result, nil // Return nil error to stop retries
	}

	// 3. Update payment status
	existing.Status = updatedPayment.Status
	existing.PaidAt = updatedPayment.PaidAt
	existing.ExternalID = updatedPayment.ExternalID

	if err := uc.repo.Update(ctx, existing); err != nil {
		log.Printf("Failed to update payment status (db): %v", err)
		return nil, errors.New("database update failed")
	}

	log.Printf("Payment %s updated to %s", existing.ID, existing.Status)

	result.Processed = true
	result.PaymentID = existing.ID
	result.NewStatus = existing.Status

	return result, nil
}

// CreateOfflinePaymentRequest represents input for offline payment registration.
type CreateOfflinePaymentRequest struct {
	Amount        float64
	Method        domain.PaymentMethod
	PayerID       uuid.UUID
	ReferenceID   uuid.UUID
	ReferenceType string
	Notes         string
	ClubID        string
}

// CreateOfflinePayment registers a payment made outside the system.
func (uc *PaymentUseCases) CreateOfflinePayment(ctx context.Context, req CreateOfflinePaymentRequest) (*domain.Payment, error) {
	now := time.Now()
	payment := &domain.Payment{
		ID:            uuid.New(),
		Amount:        decimal.NewFromFloat(req.Amount),
		Currency:      "ARS",
		Status:        domain.PaymentStatusCompleted, // Offline payments are recorded when completed
		Method:        req.Method,
		PayerID:       req.PayerID,
		ClubID:        req.ClubID,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
		Notes:         req.Notes,
		PaidAt:        &now,
	}

	if err := uc.repo.Create(ctx, payment); err != nil {
		log.Printf("Failed to create offline payment: %v", err)
		return nil, errors.New("failed to record payment")
	}

	return payment, nil
}

// ListPayments retrieves filtered payments for a club.
func (uc *PaymentUseCases) ListPayments(ctx context.Context, clubID string, filter domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	return uc.repo.List(ctx, clubID, filter)
}
