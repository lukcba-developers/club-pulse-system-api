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
	repo       domain.PaymentRepository
	gateway    domain.PaymentGateway
	responders map[string]domain.PaymentStatusResponder
}

// NewPaymentUseCases creates a new PaymentUseCases instance.
func NewPaymentUseCases(repo domain.PaymentRepository, gateway domain.PaymentGateway) *PaymentUseCases {
	return &PaymentUseCases{
		repo:       repo,
		gateway:    gateway,
		responders: make(map[string]domain.PaymentStatusResponder),
	}
}

// RegisterResponder registers a module to handle payment status changes for a specific reference type.
func (uc *PaymentUseCases) RegisterResponder(refType string, responder domain.PaymentStatusResponder) {
	uc.responders[refType] = responder
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

	// 4. Notify Responder if any
	if responder, ok := uc.responders[existing.ReferenceType]; ok {
		if err := responder.OnPaymentStatusChanged(ctx, existing.ClubID, existing.ReferenceID, existing.Status); err != nil {
			log.Printf("Responder failed for %s: %v", existing.ReferenceType, err)
			// We don't fail the webhook processing itself if responder fails,
			// though in a mission-critical app we might want to retry or use a queue.
		}
	}

	result.Processed = true
	result.PaymentID = existing.ID
	result.NewStatus = existing.Status

	return result, nil
}

// Refund finds the payment for a given reference and initiates a refund.
func (uc *PaymentUseCases) Refund(ctx context.Context, clubID string, referenceID uuid.UUID, referenceType string) error {
	// 1. Find the payment (Ideally we use a filter or a specific method)
	// For simplicity, let's assume we can list by reference.
	// But our Repo doesn't have ListByReference yet.
	// Let's use List with filter if possible.
	filter := domain.PaymentFilter{
		Status: domain.PaymentStatusCompleted,
	}
	payments, _, err := uc.repo.List(ctx, clubID, filter)
	if err != nil {
		return err
	}

	var target *domain.Payment
	for _, p := range payments {
		if p.ReferenceID == referenceID && p.ReferenceType == referenceType {
			target = p
			break
		}
	}

	if target == nil {
		return nil // No payment to refund or already refunded
	}

	if target.ExternalID == "" {
		// Not a gateway payment (maybe cash), mark as refunded in DB
		target.Status = domain.PaymentStatusRefunded
		return uc.repo.Update(ctx, target)
	}

	// 2. Call Gateway
	if err := uc.gateway.Refund(ctx, target.ExternalID); err != nil {
		return err
	}

	// 3. Update Status
	target.Status = domain.PaymentStatusRefunded
	return uc.repo.Update(ctx, target)
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

	// Notify Responder if any
	if responder, ok := uc.responders[payment.ReferenceType]; ok {
		if err := responder.OnPaymentStatusChanged(ctx, payment.ClubID, payment.ReferenceID, payment.Status); err != nil {
			log.Printf("Responder failed for %s (offline): %v", payment.ReferenceType, err)
		}
	}

	return payment, nil
}

// ListPayments retrieves filtered payments for a club.
func (uc *PaymentUseCases) ListPayments(ctx context.Context, clubID string, filter domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	return uc.repo.List(ctx, clubID, filter)
}
