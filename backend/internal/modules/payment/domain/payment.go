package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

type PaymentMethod string

const (
	PaymentMethodCash          PaymentMethod = "CASH"
	PaymentMethodMercadoPago   PaymentMethod = "MERCADOPAGO"
	PaymentMethodStripe        PaymentMethod = "STRIPE"
	PaymentMethodTransfer      PaymentMethod = "TRANSFER"
	PaymentMethodLaborExchange PaymentMethod = "LABOR_EXCHANGE"
)

type Payment struct {
	ID            uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Amount        decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency      string          `json:"currency" gorm:"not null;default:'ARS'"`
	Status        PaymentStatus   `json:"status" gorm:"not null;default:'PENDING'"`
	Method        PaymentMethod   `json:"method" gorm:"not null"`
	ExternalID    string          `json:"external_id"`          // ID from Payment Provider
	ClubID        string          `json:"club_id" gorm:"index"` // Tenant Isolation
	PayerID       uuid.UUID       `json:"payer_id" gorm:"type:uuid;not null;index"`
	ReferenceID   uuid.UUID       `json:"reference_id" gorm:"type:uuid;index"` // Could be Membership ID or Booking ID
	ReferenceType string          `json:"reference_type"`                      // "MEMBERSHIP", "BOOKING"
	Notes         string          `json:"notes" gorm:"type:text"`              // Details for Offline/Labor payments

	PaidAt    *time.Time     `json:"paid_at,omitempty"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type PaymentFilter struct {
	PayerID   uuid.UUID
	Status    PaymentStatus
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	Update(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	GetByExternalID(ctx context.Context, externalID string) (*Payment, error)
	List(ctx context.Context, clubID string, filter PaymentFilter) ([]*Payment, int64, error)
}

// PaymentStatusResponder allows other modules to react to payment status changes.
type PaymentStatusResponder interface {
	OnPaymentStatusChanged(ctx context.Context, clubID string, referenceID uuid.UUID, status PaymentStatus) error
}
