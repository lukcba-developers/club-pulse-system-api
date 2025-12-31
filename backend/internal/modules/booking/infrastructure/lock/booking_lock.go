package lock

import (
	"context"
	"fmt"
	"time"

	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

// BookingLock provides distributed locking for booking slots
// Prevents race conditions when multiple users try to book the same slot
type BookingLock struct {
	redis *platformRedis.RedisClient
}

// NewBookingLock creates a new booking lock service
func NewBookingLock() *BookingLock {
	return &BookingLock{
		redis: platformRedis.GetClient(),
	}
}

// lockKey generates a unique key for a facility/time slot combination
func lockKey(facilityID string, start, end time.Time) string {
	return fmt.Sprintf("booking_lock:%s:%d:%d",
		facilityID,
		start.Unix(),
		end.Unix(),
	)
}

// AcquireLock attempts to acquire a lock for a booking slot
// Returns true if lock acquired, false if slot is already locked by another user
// TTL is the maximum time the lock will be held (e.g., 5 minutes for checkout flow)
func (l *BookingLock) AcquireLock(ctx context.Context, facilityID string, start, end time.Time, userID string, ttl time.Duration) (bool, error) {
	key := lockKey(facilityID, start, end)

	// Try to set the lock (SetNX = SET if Not eXists)
	acquired, err := l.redis.SetNX(ctx, key, userID, ttl)
	if err != nil {
		return false, fmt.Errorf("failed to acquire booking lock: %w", err)
	}

	return acquired, nil
}

// ReleaseLock releases a booking lock
func (l *BookingLock) ReleaseLock(ctx context.Context, facilityID string, start, end time.Time) error {
	key := lockKey(facilityID, start, end)
	return l.redis.Del(ctx, key)
}

// IsLocked checks if a slot is currently locked
func (l *BookingLock) IsLocked(ctx context.Context, facilityID string, start, end time.Time) (bool, error) {
	key := lockKey(facilityID, start, end)
	return l.redis.Exists(ctx, key)
}

// GetLockHolder returns the user ID holding the lock, or empty string if unlocked
func (l *BookingLock) GetLockHolder(ctx context.Context, facilityID string, start, end time.Time) (string, error) {
	key := lockKey(facilityID, start, end)
	userID, err := l.redis.Get(ctx, key)
	if err != nil {
		// redis.Nil error means key doesn't exist
		return "", nil
	}
	return userID, nil
}

// ExtendLock extends the TTL of an existing lock (for long checkout processes)
func (l *BookingLock) ExtendLock(ctx context.Context, facilityID string, start, end time.Time, additionalTTL time.Duration) error {
	key := lockKey(facilityID, start, end)
	return l.redis.Expire(ctx, key, additionalTTL)
}
