package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SubscriptionStatus string

const (
	SubscriptionActive    SubscriptionStatus = "ACTIVE"
	SubscriptionPaused    SubscriptionStatus = "PAUSED"
	SubscriptionCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionPastDue   SubscriptionStatus = "PAST_DUE" // Payment failed, retrying
)

// Subscription represents a recurring payment agreement for a membership.
type Subscription struct {
	ID              uuid.UUID          `json:"id"`
	UserID          string             `json:"user_id"`
	MembershipID    uuid.UUID          `json:"membership_id"`
	Amount          decimal.Decimal    `json:"amount"`
	Currency        string             `json:"currency"`
	Status          SubscriptionStatus `json:"status"`
	PaymentMethodID string             `json:"payment_method_id"` // Token or ID from Payment Provider (e.g., MP Card ID)
	NextBillingDate time.Time          `json:"next_billing_date"`
	LastPaymentDate *time.Time         `json:"last_payment_date,omitempty"`
	FailCount       int                `json:"fail_count"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

func NewSubscription(userID string, membershipID uuid.UUID, amount decimal.Decimal, paymentMethodID string) *Subscription {
	return &Subscription{
		ID:              uuid.New(),
		UserID:          userID,
		MembershipID:    membershipID,
		Amount:          amount,
		Currency:        "ARS",
		Status:          SubscriptionActive,
		PaymentMethodID: paymentMethodID,
		NextBillingDate: time.Now().AddDate(0, 1, 0), // Default next month? detailed logic needed in service
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

type SubscriptionRepository interface {
	Create(subscription *Subscription) error
	GetByID(id uuid.UUID) (*Subscription, error)
	GetByUserID(userID string) ([]Subscription, error)
	Update(subscription *Subscription) error
}
