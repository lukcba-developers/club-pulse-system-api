package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Scholarship struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	Percentage decimal.Decimal `json:"percentage"` // e.g., 0.50
	Reason     string          `json:"reason"`
	GrantorID  string          `json:"grantor_id"`
	ValidUntil *time.Time      `json:"valid_until"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type ScholarshipRepository interface {
	Create(scholarship *Scholarship) error
	GetByUserID(userID string) ([]*Scholarship, error)
	GetActiveByUserID(userID string) (*Scholarship, error) // Assuming one active per user
	ListActiveByUserIDs(userIDs []string) (map[string]*Scholarship, error)
}

// ApplyDiscount calculates the discounted amount.
// If the scholarship is invalid or not active, returns original amount.
func (s *Scholarship) ApplyDiscount(amount decimal.Decimal) decimal.Decimal {
	if !s.IsActive {
		return amount
	}
	if s.ValidUntil != nil && s.ValidUntil.Before(time.Now()) {
		return amount
	}

	discount := amount.Mul(s.Percentage)
	return amount.Sub(discount)
}
