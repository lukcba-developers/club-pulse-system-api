package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/shopspring/decimal"
)

type PaymentHandler struct {
	repo    domain.PaymentRepository
	gateway domain.PaymentGateway
}

func NewPaymentHandler(repo domain.PaymentRepository, gateway domain.PaymentGateway) *PaymentHandler {
	return &PaymentHandler{
		repo:    repo,
		gateway: gateway,
	}
}

type CheckoutRequest struct {
	Amount        float64 `json:"amount" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PayerEmail    string  `json:"payer_email" binding:"required,email"`
	ReferenceID   string  `json:"reference_id" binding:"required"`
	ReferenceType string  `json:"reference_type" binding:"required"` // MEMBERSHIP, BOOKING
}

// Checkout creates a payment Intent and returns the MP Preference URL
func (h *PaymentHandler) Checkout(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Create Payment Record (Pending)
	// We assume PayerID comes from Context (Auth Middleware)
	userIDStr := c.GetString("userID")
	var userID uuid.UUID
	var err error

	if userIDStr != "" {
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			log.Printf("Invalid User ID in context: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user session"})
			return
		}
	} else {
		// If not authenticated (should be prevented by middleware), fail or prompt login
		// For now, we return 401
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	refID, err := uuid.Parse(req.ReferenceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reference_id format"})
		return
	}

	clubID := c.GetString("clubID")
	if clubID == "" {
		// Defensive: Middleware should ensure this, but handle legacy/error cases
		// For payments, maybe we want to allow global context? No, strict tenancy.
		c.JSON(http.StatusBadRequest, gin.H{"error": "club context required for payment"})
		return
	}

	payment := &domain.Payment{
		ID:            uuid.New(),
		Amount:        decimal.NewFromFloat(req.Amount),
		Currency:      "ARS",
		Status:        domain.PaymentStatusPending,
		Method:        domain.PaymentMethodMercadoPago, // Default for this endpoint
		PayerID:       userID,
		ClubID:        clubID,
		ReferenceID:   refID,
		ReferenceType: req.ReferenceType,
	}

	if err := h.repo.Create(c.Request.Context(), payment); err != nil {
		log.Printf("Failed to create payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment record", "details": err.Error()})
		return
	}

	// 2. Call Gateway
	url, err := h.gateway.CreatePreference(c.Request.Context(), payment, req.PayerEmail, req.Description)
	if err != nil {
		log.Printf("Gateway Error: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact payment gateway"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// HandleWebhook receives notifications from payment providers
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// 0. Validate Signature
	if err := h.gateway.ValidateWebhook(c.Request); err != nil {
		log.Printf("Invalid Webhook Signature: %v", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// MP query param type=payment
	webhookType := c.Query("type")

	// Sometimes MP sends topic=payment
	if webhookType == "" {
		webhookType = c.Query("topic")
	}

	if webhookType == "payment" {
		dataID := c.Query("data.id")
		if dataID == "" {
			dataID = c.Query("id")
		}

		if dataID != "" {
			updatedPayment, err := h.gateway.ProcessWebhook(c.Request.Context(), dataID)
			if err != nil {
				log.Printf("Webhook processing failed (gateway): %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "processing failed"})
				return
			}

			if updatedPayment != nil {
				existing, err := h.repo.GetByID(c.Request.Context(), updatedPayment.ID)
				if err != nil {
					log.Printf("Payment not found for update: %s", updatedPayment.ID)
					// If payment not found, maybe 404? But usually we want to retry if it's a race?
					// Or if it simply doesn't exist, maybe 200 to stop retry loops?
					// Let's return 200 to stop retries if we really can't find it.
					c.Status(http.StatusOK)
					return
				}

				existing.Status = updatedPayment.Status
				existing.PaidAt = updatedPayment.PaidAt
				existing.ExternalID = updatedPayment.ExternalID

				if err := h.repo.Update(c.Request.Context(), existing); err != nil {
					log.Printf("Failed to update payment status (db): %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"})
					return
				}
				log.Printf("Payment %s updated to %s", existing.ID, existing.Status)
			}
		}
	}

	c.Status(http.StatusOK)
}

func RegisterRoutes(r *gin.RouterGroup, handler *PaymentHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	payments := r.Group("/payments")
	{
		// Protected
		payments.POST("/checkout", authMiddleware, tenantMiddleware, handler.Checkout)

		// Public (Webhook)
		payments.POST("/webhook", handler.HandleWebhook)
	}
}
