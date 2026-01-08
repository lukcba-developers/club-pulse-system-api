package gateways

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/mercadopago/sdk-go/pkg/config"
	mp_payment "github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type MercadoPagoGateway struct {
	accessToken   string
	webhookSecret string
}

func NewMercadoPagoGateway() *MercadoPagoGateway {
	token := os.Getenv("MP_ACCESS_TOKEN")
	if token == "" {
		// Fallback for dev/test if not set, though SDK might complain
		token = "TEST-ACCESS-TOKEN-PLACEHOLDER"
	}
	secret := os.Getenv("MP_WEBHOOK_SECRET")
	return &MercadoPagoGateway{
		accessToken:   token,
		webhookSecret: secret,
	}
}

func (g *MercadoPagoGateway) CreatePreference(ctx context.Context, payment *domain.Payment, payerEmail string, description string) (string, error) {
	cfg, err := config.New(g.accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to create payment config: %w", err)
	}

	client := preference.NewClient(cfg)

	// Amount conversion (Decimal to Float64)
	amount, _ := payment.Amount.Float64()

	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:       description,
				Quantity:    1,
				UnitPrice:   amount,
				CurrencyID:  payment.Currency,
				Description: description,
			},
		},
		Payer: &preference.PayerRequest{
			Email: payerEmail,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:3000/payment/result?status=success",
			Failure: "http://localhost:3000/payment/result?status=failure",
			Pending: "http://localhost:3000/payment/result?status=pending",
		},
		AutoReturn:        "approved",
		ExternalReference: payment.ID.String(),
		NotificationURL:   "https://your-domain.ngrok.io/api/v1/payments/webhook", // Placeholder for local dev
	}

	resp, err := client.Create(ctx, request)
	if err != nil {
		return "", fmt.Errorf("error creating preference: %w", err)
	}

	return resp.InitPoint, nil
}

func (g *MercadoPagoGateway) ProcessWebhook(ctx context.Context, payload interface{}) (*domain.Payment, error) {
	// Payload is expected to be the Payment ID (string or int)
	var paymentID int
	switch v := payload.(type) {
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid payment id format: %w", err)
		}
		paymentID = id
	case int:
		paymentID = v
	case int64:
		paymentID = int(v)
	case float64:
		paymentID = int(v)
	default:
		return nil, fmt.Errorf("unsupported payload type for webhook")
	}

	cfg, err := config.New(g.accessToken)
	if err != nil {
		return nil, err
	}

	client := mp_payment.NewClient(cfg)
	mpdata, err := client.Get(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment from MP: %w", err)
	}

	// Map Status
	status := domain.PaymentStatusPending
	if mpdata.Status == "approved" {
		status = domain.PaymentStatusCompleted
	} else if mpdata.Status == "rejected" || mpdata.Status == "cancelled" {
		status = domain.PaymentStatusFailed
	} else if mpdata.Status == "refunded" {
		status = domain.PaymentStatusRefunded
	}

	// Extract External Reference (Our Payment UUID)
	paymentUUID, err := uuid.Parse(mpdata.ExternalReference)
	if err != nil {
		// Log warning?
		return nil, fmt.Errorf("invalid external reference in MP payment: %s", mpdata.ExternalReference)
	}

	now := time.Now()

	return &domain.Payment{
		ID:         paymentUUID,
		Status:     status,
		ExternalID: strconv.Itoa(paymentID),
		Method:     domain.PaymentMethodMercadoPago,
		PaidAt:     &now, // Approximate
	}, nil
}
func (g *MercadoPagoGateway) ValidateWebhook(req *http.Request) error {
	if g.webhookSecret == "" {
		// SECURITY FIX (VUL-007): Log critical security warning
		fmt.Println("[SECURITY] CRITICAL: MP_WEBHOOK_SECRET not configured - webhook signature validation disabled")
		return fmt.Errorf("SECURITY: webhook signature validation disabled - MP_WEBHOOK_SECRET not configured")
	}

	// 1. Get Headers
	xSignature := req.Header.Get("x-signature")
	xRequestID := req.Header.Get("x-request-id")

	if xSignature == "" || xRequestID == "" {
		return fmt.Errorf("missing signature headers")
	}

	// 2. Parse TS and V1
	// Format: ts=...;v1=...
	parts := strings.Split(xSignature, ";")
	var ts, v1 string
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			if kv[0] == "ts" {
				ts = kv[1]
			} else if kv[0] == "v1" {
				v1 = kv[1]
			}
		}
	}

	if ts == "" || v1 == "" {
		return fmt.Errorf("invalid signature format")
	}

	// 3. Reconstruct Manifest
	// Template: "id:[data.id_url_param];request-id:[x-request-id];ts:[ts];"
	// Wait, MP documentation says:
	// "id:[data.id];request-id:[x-request-id];ts:[ts];"
	// Where data.id is from the info.
	// But we don't have the body/query parsed yet passed to this function easily unless we parse it here.
	// But `req` is passed.
	// We need 'data.id'.
	// In the handler, we extracted it from Query.
	// Let's assume URL query "data.id".
	dataID := req.URL.Query().Get("data.id")
	if dataID == "" {
		dataID = req.URL.Query().Get("id")
	}
	if dataID == "" {
		// If we can't get ID, we can't validate signature which depends on it.
		return fmt.Errorf("missing data.id for signature validation")
	}

	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", dataID, xRequestID, ts)

	// 4. Compute HMAC
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write([]byte(manifest))
	computedHash := hex.EncodeToString(mac.Sum(nil))

	if computedHash != v1 {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
