package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/gateways"
)

type PaymentHandler struct {
	repo      domain.PaymentRepository
	processor gateways.PaymentProcessor
}

func NewPaymentHandler(repo domain.PaymentRepository, processor gateways.PaymentProcessor) *PaymentHandler {
	return &PaymentHandler{
		repo:      repo,
		processor: processor,
	}
}

// HandleWebhook receives notifications from payment providers
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// In production, verify signature and determine provider
	// For Phase 1 Mock, we assume it's a success signal for a payment ID passed in Query

	// Real implementation would parse body
	// payload, _ := c.GetRawData()

	externalID := c.Query("type")
	if externalID == "payment" {
		// Mock logic: Update status
		log.Println("Received Payment Webhook")

		// In a real app we would:
		// 1. Parse ID from payload
		// 2. Fetch Payment
		// 3. Update Status
		// 4. Update Mebership/Booking if needed
	}

	c.Status(http.StatusOK)
}

func RegisterRoutes(r *gin.RouterGroup, handler *PaymentHandler) {
	payments := r.Group("/payments")
	{
		payments.POST("/webhook", handler.HandleWebhook)
	}
}
