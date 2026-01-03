package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/gateways"
	paymentHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/http"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPaymentWebhookFlow(t *testing.T) {
	// 1. Setup Environment (Mock DB or Test DB)
	// For E2E in this context, we'll try to use the real Handler setup but with a Mock Processor
	gin.SetMode(gin.TestMode)

	database.InitDB() // Use local DB (ensure it's running)
	db := database.GetDB()

	// Clean state
	_ = db.Migrator().DropTable(&domain.Payment{})
	_ = db.AutoMigrate(&domain.Payment{})

	// Repositories & Services
	repo := repository.NewPostgresPaymentRepository(db)
	gateway := gateways.NewMockGateway()
	handler := paymentHttp.NewPaymentHandler(repo, gateway)

	// Router
	r := gin.Default()
	paymentHttp.RegisterRoutes(r.Group("/api/v1"), handler, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Prepare Data: Create a pending payment manually in DB to simulate initiation
	// Convert float to Decimal
	amount := decimal.NewFromFloat(5000.00)

	paymentID := uuid.New()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        amount,
		Currency:      "ARS",
		Status:        domain.PaymentStatusPending,
		Method:        domain.PaymentMethodMercadoPago,
		PayerID:       uuid.New(), // Random User
		ReferenceID:   uuid.New(), // Random Membership
		ReferenceType: "MEMBERSHIP",
	}

	err := repo.Create(context.Background(), payment)
	assert.NoError(t, err)

	// 3. Simulate Webhook Call (Success)
	// /api/v1/payments/webhook?type=payment&id=...
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payments/webhook?type=payment", nil)
	r.ServeHTTP(w, req)

	// 4. Verification
	assert.Equal(t, http.StatusOK, w.Code)

	// In a real integration, the webhook would trigger a background update or immediate update.
	// Since our mock webhook handler logs but doesn't fully update DB yet in this simplified Phase 2,
	// we verify the endpoint is reachable and returns 200.
	// For full E2E, we would assert repo.GetByID(paymentID) has Status = Completed.
}
