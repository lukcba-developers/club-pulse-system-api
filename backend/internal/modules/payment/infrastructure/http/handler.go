package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
)

// PaymentHandler handles HTTP requests for payments.
// Following Clean Architecture: handlers only parse HTTP, delegate logic to use cases.
type PaymentHandler struct {
	useCases *application.PaymentUseCases
}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler(useCases *application.PaymentUseCases) *PaymentHandler {
	return &PaymentHandler{useCases: useCases}
}

// CheckoutRequest is the HTTP request body for checkout.
type CheckoutRequest struct {
	Amount        float64 `json:"amount" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PayerEmail    string  `json:"payer_email" binding:"required,email"`
	ReferenceID   string  `json:"reference_id" binding:"required"`
	ReferenceType string  `json:"reference_type" binding:"required"` // MEMBERSHIP, BOOKING
}

// OfflinePaymentRequest is the HTTP request body for offline payments.
type OfflinePaymentRequest struct {
	Amount        float64 `json:"amount" binding:"required"`
	Method        string  `json:"method" binding:"required,oneof=CASH LABOR_EXCHANGE TRANSFER"`
	PayerID       string  `json:"payer_id" binding:"required"`
	ReferenceID   string  `json:"reference_id"`
	ReferenceType string  `json:"reference_type"`
	Notes         string  `json:"notes"`
}

// Checkout creates a payment intent and returns the MP Preference URL.
func (h *PaymentHandler) Checkout(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user session"})
		return
	}

	refID, err := uuid.Parse(req.ReferenceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reference_id format"})
		return
	}

	clubID := c.GetString("clubID")
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club context required for payment"})
		return
	}

	payment, url, err := h.useCases.Checkout(c.Request.Context(), application.CheckoutRequest{
		Amount:        req.Amount,
		Description:   req.Description,
		PayerEmail:    req.PayerEmail,
		ReferenceID:   refID,
		ReferenceType: req.ReferenceType,
		UserID:        userID,
		ClubID:        clubID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": url,
		"payment_id":   payment.ID,
	})
}

// HandleWebhook receives notifications from payment providers.
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// 1. Validate Signature (delegates to use case which uses gateway)
	if err := h.useCases.ValidateWebhook(c.Request); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// 2. Parse webhook type from query params
	webhookType := c.Query("type")
	if webhookType == "" {
		webhookType = c.Query("topic") // MP sometimes sends "topic" instead
	}

	dataID := c.Query("data.id")
	if dataID == "" {
		dataID = c.Query("id")
	}

	// 3. Delegate to use case
	result, err := h.useCases.ProcessWebhook(c.Request.Context(), application.ProcessWebhookRequest{
		Type:   webhookType,
		DataID: dataID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return 200 to acknowledge receipt (even if not found, to stop retries)
	c.Status(http.StatusOK)
	_ = result // Result logged internally, no need to expose
}

// ListPayments returns filtered payments for the dashboard.
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	clubID := c.GetString("clubID")

	var filter domain.PaymentFilter

	if payerID := c.Query("payer_id"); payerID != "" {
		if id, err := uuid.Parse(payerID); err == nil {
			filter.PayerID = id
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = domain.PaymentStatus(status)
	}

	filter.Limit = 20 // Default

	payments, total, err := h.useCases.ListPayments(c.Request.Context(), clubID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  payments,
		"total": total,
	})
}

// CreateOfflinePayment registers a payment made outside the system.
// SECURITY: Only ADMIN and STAFF can register offline payments.
func (h *PaymentHandler) CreateOfflinePayment(c *gin.Context) {
	// RBAC Check
	role := c.GetString("userRole")
	if role != "ADMIN" && role != "STAFF" && role != "SUPER_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to create offline payments"})
		return
	}

	var req OfflinePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")

	payerUUID, err := uuid.Parse(req.PayerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payer_id"})
		return
	}

	var refID uuid.UUID
	if req.ReferenceID != "" {
		refID, _ = uuid.Parse(req.ReferenceID)
	}

	payment, err := h.useCases.CreateOfflinePayment(c.Request.Context(), application.CreateOfflinePaymentRequest{
		Amount:        req.Amount,
		Method:        domain.PaymentMethod(req.Method),
		PayerID:       payerUUID,
		ReferenceID:   refID,
		ReferenceType: req.ReferenceType,
		Notes:         req.Notes,
		ClubID:        clubID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": payment})
}

// RefundPayment processes a manual refund for a completed payment.
// SECURITY: Only ADMIN can process refunds.
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	// RBAC Check
	role := c.GetString("userRole")
	if role != "ADMIN" && role != "SUPER_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to process refunds"})
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment ID required"})
		return
	}

	pID, err := uuid.Parse(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID format"})
		return
	}

	clubID := c.GetString("clubID")

	// Get the payment first to extract reference info
	payments, _, err := h.useCases.ListPayments(c.Request.Context(), clubID, domain.PaymentFilter{
		Limit: 1000, // Get all to find by ID
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find payment"})
		return
	}

	var targetPayment *domain.Payment
	for _, p := range payments {
		if p.ID == pID {
			targetPayment = p
			break
		}
	}

	if targetPayment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	if targetPayment.Status != domain.PaymentStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only completed payments can be refunded"})
		return
	}

	// Process refund
	if err := h.useCases.Refund(c.Request.Context(), clubID, targetPayment.ReferenceID, targetPayment.ReferenceType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "refund processed successfully"})
}

// RegisterRoutes registers payment HTTP routes.
func RegisterRoutes(r *gin.RouterGroup, handler *PaymentHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	payments := r.Group("/payments")
	{
		// Protected endpoints
		payments.POST("/checkout", authMiddleware, tenantMiddleware, handler.Checkout)
		payments.POST("/offline", authMiddleware, tenantMiddleware, handler.CreateOfflinePayment)
		payments.POST("/:id/refund", authMiddleware, tenantMiddleware, handler.RefundPayment)
		payments.GET("", authMiddleware, tenantMiddleware, handler.ListPayments)

		// Public endpoint (Webhook)
		payments.POST("/webhook", handler.HandleWebhook)
	}
}
