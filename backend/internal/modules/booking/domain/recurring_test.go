package domain_test

import (
	"testing"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	"github.com/stretchr/testify/assert"
)

func TestRecurringRule_Defaults(t *testing.T) {
	// Test that Type should ideally be treated as FIXED by default if logic enforces it,
	// but here we just check struct instantiation.
	rule := domain.RecurringRule{
		Frequency: "WEEKLY",
		Type:      domain.RecurrenceTypeFixed,
	}

	assert.Equal(t, "WEEKLY", rule.Frequency)
	assert.Equal(t, domain.RecurrenceType("FIXED"), rule.Type)
}

func TestRecurrenceType_Values(t *testing.T) {
	// Verify enums match expectations
	assert.Equal(t, domain.RecurrenceType("CLASS"), domain.RecurrenceTypeClass)
	assert.Equal(t, domain.RecurrenceType("MAINTENANCE"), domain.RecurrenceTypeMaintenance)
	assert.Equal(t, domain.RecurrenceType("FIXED"), domain.RecurrenceTypeFixed)
}

func TestRecurringRule_ActiveStatus(t *testing.T) {
	now := time.Now()
	rule := domain.RecurringRule{
		StartDate: now.Add(-24 * time.Hour),
		EndDate:   now.Add(24 * time.Hour), // Must be time.Time, not pointer
	}
	// Just asserting fields exist and are set
	assert.False(t, rule.EndDate.IsZero())
}
