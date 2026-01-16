package domain_test

import (
	"testing"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestConsentRecord_IsActive(t *testing.T) {
	t.Run("Active Consent", func(t *testing.T) {
		consent := &domain.ConsentRecord{
			Accepted:  true,
			RevokedAt: nil,
		}
		assert.True(t, consent.IsActive())
	})

	t.Run("Not Accepted", func(t *testing.T) {
		consent := &domain.ConsentRecord{
			Accepted:  false,
			RevokedAt: nil,
		}
		assert.False(t, consent.IsActive())
	})

	t.Run("Revoked Consent", func(t *testing.T) {
		now := time.Now()
		consent := &domain.ConsentRecord{
			Accepted:  true,
			RevokedAt: &now,
		}
		assert.False(t, consent.IsActive())
	})

	t.Run("Revoked and Not Accepted (Edge Case)", func(t *testing.T) {
		now := time.Now()
		consent := &domain.ConsentRecord{
			Accepted:  false,
			RevokedAt: &now,
		}
		assert.False(t, consent.IsActive())
	})
}
