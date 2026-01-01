package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalculateLateFee_Overdue(t *testing.T) {
	// Setup
	tier := domain.MembershipTier{
		MonthlyFee: decimal.NewFromFloat(100.00),
	}

	// Due date was yesterday
	dueDate := time.Now().Add(-24 * time.Hour)

	membership := domain.Membership{
		ID:              uuid.New(),
		NextBillingDate: dueDate,
		MembershipTier:  tier,
	}

	// Execution
	fee := membership.CalculateLateFee()

	// Validation
	// Expect 10% of 100.00 = 10.00
	expectedFee := decimal.NewFromFloat(10.00)
	assert.True(t, expectedFee.Equal(fee), "Late fee should be 10% of monthly fee")
}

func TestCalculateLateFee_OnTime(t *testing.T) {
	// Setup
	tier := domain.MembershipTier{
		MonthlyFee: decimal.NewFromFloat(100.00),
	}

	// Due date is tomorrow
	dueDate := time.Now().Add(24 * time.Hour)

	membership := domain.Membership{
		ID:              uuid.New(),
		NextBillingDate: dueDate,
		MembershipTier:  tier,
	}

	// Execution
	fee := membership.CalculateLateFee()

	// Validation
	assert.True(t, decimal.Zero.Equal(fee), "Late fee should be zero if not overdue")
}
