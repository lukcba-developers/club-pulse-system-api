package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/gateways"
	paymentHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/http"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPaymentWebhookFlow(t *testing.T) {
	// 1. Setup Environment
	gin.SetMode(gin.TestMode)

	database.InitDB()
	db := database.GetDB()

	// Clean state
	_ = db.Migrator().DropTable(&domain.Payment{})
	_ = db.AutoMigrate(&domain.Payment{})

	// Clear PostgreSQL cached prepared statements after schema change
	db.Exec("DISCARD ALL")

	// Repositories & Services
	repo := repository.NewPostgresPaymentRepository(db)
	gateway := gateways.NewMockGateway()
	useCases := application.NewPaymentUseCases(repo, gateway)
	handler := paymentHttp.NewPaymentHandler(useCases)

	// Router
	r := gin.Default()
	paymentHttp.RegisterRoutes(r.Group("/api/v1"), handler, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Prepare Data: Create a pending payment manually in DB
	amount := decimal.NewFromFloat(5000.00)

	paymentID := uuid.New()
	payment := &domain.Payment{
		ID:            paymentID,
		Amount:        amount,
		Currency:      "ARS",
		Status:        domain.PaymentStatusPending,
		Method:        domain.PaymentMethodMercadoPago,
		PayerID:       uuid.New(),
		ReferenceID:   uuid.New(),
		ReferenceType: "MEMBERSHIP",
	}

	err := repo.Create(context.Background(), payment)
	assert.NoError(t, err)

	// 3. Simulate Webhook Call (Success)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/payments/webhook?type=payment", nil)
	r.ServeHTTP(w, req)

	// 4. Verification
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPaymentCheckout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()
	_ = db.AutoMigrate(&domain.Payment{})

	repo := repository.NewPostgresPaymentRepository(db)
	gateway := gateways.NewMockGateway()
	useCases := application.NewPaymentUseCases(repo, gateway)
	handler := paymentHttp.NewPaymentHandler(useCases)

	r := gin.New()
	// Mock Auth and Tenant Middleware
	authMw := func(c *gin.Context) {
		c.Set("userID", uuid.New().String())
		c.Next()
	}
	tenantMw := func(c *gin.Context) {
		c.Set("clubID", "test-club-id")
		c.Next()
	}

	paymentHttp.RegisterRoutes(r.Group("/api/v1"), handler, authMw, tenantMw)

	t.Run("Checkout Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"amount": 100.00, "description": "Test", "payer_email": "test@test.com", "reference_id": "` + uuid.New().String() + `", "reference_type": "MEMBERSHIP"}`
		req, _ := http.NewRequest("POST", "/api/v1/payments/checkout", strings.NewReader(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Checkout Missing ClubID", func(t *testing.T) {
		// Router without Tenant Middleware
		r2 := gin.New()
		authMw := func(c *gin.Context) {
			c.Set("userID", uuid.New().String())
			c.Next()
		}
		// Empty Tenant Middleware (Simulate failure)
		tenantMwFail := func(c *gin.Context) {
			// ClubID not set
			c.Next()
		}
		paymentHttp.RegisterRoutes(r2.Group("/api/v1"), handler, authMw, tenantMwFail)

		w := httptest.NewRecorder()
		body := `{"amount": 100.00, "description": "Test", "payer_email": "test@test.com", "reference_id": "` + uuid.New().String() + `", "reference_type": "MEMBERSHIP"}`
		req, _ := http.NewRequest("POST", "/api/v1/payments/checkout", strings.NewReader(body))
		r2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
