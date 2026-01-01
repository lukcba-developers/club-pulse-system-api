package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Enums for Status and Cycle
type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "ACTIVE"
	MembershipStatusInactive  MembershipStatus = "INACTIVE"
	MembershipStatusPending   MembershipStatus = "PENDING"
	MembershipStatusCancelled MembershipStatus = "CANCELLED"
	MembershipStatusExpired   MembershipStatus = "EXPIRED"
)

type BillingCycle string

const (
	BillingCycleMonthly    BillingCycle = "MONTHLY"
	BillingCycleQuarterly  BillingCycle = "QUARTERLY"
	BillingCycleSemiAnnual BillingCycle = "SEMI_ANNUAL"
	BillingCycleAnnual     BillingCycle = "ANNUAL"
)

// MembershipTier defines the types of memberships available (Gold, Silver, etc.)
type MembershipTier struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string          `json:"name" gorm:"not null;size:255"`
	Description string          `json:"description" gorm:"type:text"`
	MonthlyFee  decimal.Decimal `json:"monthly_fee" gorm:"type:decimal(10,2);not null"`
	Colors      string          `json:"colors" gorm:"size:50"` // e.g. "bg-amber-100 text-amber-800" for frontend
	Benefits    pq.StringArray  `json:"benefits" gorm:"type:text[]"`
	IsActive    bool            `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Membership represents a user's subscription
type Membership struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	MembershipTierID uuid.UUID      `json:"membership_tier_id" gorm:"type:uuid;not null"`
	MembershipTier   MembershipTier `json:"membership_tier" gorm:"foreignKey:MembershipTierID"`

	Status       MembershipStatus `json:"status" gorm:"not null;default:'PENDING'"`
	BillingCycle BillingCycle     `json:"billing_cycle" gorm:"not null;default:'MONTHLY'"`

	StartDate          time.Time       `json:"start_date" gorm:"not null"`
	EndDate            *time.Time      `json:"end_date,omitempty"`
	NextBillingDate    time.Time       `json:"next_billing_date" gorm:"not null"`
	OutstandingBalance decimal.Decimal `json:"outstanding_balance" gorm:"type:decimal(10,2);default:0"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// CalculateLateFee determines if a late fee should be applied based on days past due
// Returns the fee amount (e.g., 10% of monthly fee)
func (m *Membership) CalculateLateFee() decimal.Decimal {
	now := time.Now()
	if now.Before(m.NextBillingDate) {
		return decimal.Zero
	}

	// Example Logic: 10% penalty if overdue
	// In real app, this might be configurable per Tier
	return m.MembershipTier.MonthlyFee.Mul(decimal.NewFromFloat(0.10))
}

// Repository Interface
type MembershipRepository interface {
	Create(ctx context.Context, membership *Membership) error
	GetByID(ctx context.Context, id uuid.UUID) (*Membership, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Membership, error)
	ListTiers(ctx context.Context) ([]MembershipTier, error)
	GetTierByID(ctx context.Context, id uuid.UUID) (*MembershipTier, error)
}
