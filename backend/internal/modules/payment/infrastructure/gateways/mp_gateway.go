package gateways

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/mercadopago/sdk-go/pkg/config"
	mp_payment "github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type MercadoPagoGateway struct {
	accessToken string
}

func NewMercadoPagoGateway() *MercadoPagoGateway {
	token := os.Getenv("MP_ACCESS_TOKEN")
	if token == "" {
		// Fallback for dev/test if not set, though SDK might complain
		token = "TEST-ACCESS-TOKEN-PLACEHOLDER"
	}
	return &MercadoPagoGateway{
		accessToken: token,
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
