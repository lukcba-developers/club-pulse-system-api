package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

// ReservationLock provides temporary locking for booking slots
// This prevents "slot stealing" when a user is completing the booking process
type ReservationLock struct {
	redis *platformRedis.RedisClient
	ttl   time.Duration
}

// LockInfo contains information about a lock
type LockInfo struct {
	LockID     string    `json:"lock_id"`
	FacilityID string    `json:"facility_id"`
	UserID     string    `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// NewReservationLock creates a new reservation lock manager
func NewReservationLock(ttl time.Duration) *ReservationLock {
	return &ReservationLock{
		redis: platformRedis.GetClient(),
		ttl:   ttl,
	}
}

// lockKey generates the Redis key for a slot lock
func lockKey(facilityID string, start, end time.Time) string {
	return fmt.Sprintf("lock:slot:%s:%d:%d", facilityID, start.Unix(), end.Unix())
}

// Acquire attempts to acquire a lock for a specific slot
// Returns the lock ID if successful, empty string if slot is already locked
func (l *ReservationLock) Acquire(ctx context.Context, facilityID, userID string, start, end time.Time) (string, error) {
	lockID := uuid.New().String()
	key := lockKey(facilityID, start, end)

	// Use SETNX to ensure only one lock can be created
	acquired, err := l.redis.SetNX(ctx, key, lockID, l.ttl)
	if err != nil {
		return "", fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !acquired {
		return "", fmt.Errorf("slot already locked by another user")
	}

	// Store additional lock info (optional, for debugging/admin)
	infoKey := fmt.Sprintf("lock:info:%s", lockID)
	info := fmt.Sprintf("%s|%s|%s|%d|%d",
		facilityID, userID, lockID, start.Unix(), end.Unix())
	_ = l.redis.Set(ctx, infoKey, info, l.ttl)

	return lockID, nil
}

// Release removes a lock for a specific slot
func (l *ReservationLock) Release(ctx context.Context, facilityID string, start, end time.Time) error {
	key := lockKey(facilityID, start, end)
	return l.redis.Del(ctx, key)
}

// ReleaseByID releases a lock by its ID (verifies ownership)
func (l *ReservationLock) ReleaseByID(ctx context.Context, lockID, facilityID string, start, end time.Time) error {
	key := lockKey(facilityID, start, end)

	// Get the current lock value
	currentLock, err := l.redis.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("lock not found or expired")
	}

	// Verify ownership
	if currentLock != lockID {
		return fmt.Errorf("lock owned by different user")
	}

	// Delete the lock
	if err := l.redis.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	// Clean up info key
	infoKey := fmt.Sprintf("lock:info:%s", lockID)
	_ = l.redis.Del(ctx, infoKey)

	return nil
}

// IsLocked checks if a specific slot is currently locked
func (l *ReservationLock) IsLocked(ctx context.Context, facilityID string, start, end time.Time) (bool, error) {
	key := lockKey(facilityID, start, end)
	exists, err := l.redis.Exists(ctx, key)
	if err != nil {
		// If Redis fails, assume not locked (fail open for availability)
		return false, nil
	}
	return exists, nil
}

// ExtendLock extends the TTL of an existing lock
func (l *ReservationLock) ExtendLock(ctx context.Context, lockID, facilityID string, start, end time.Time) error {
	key := lockKey(facilityID, start, end)

	// Get the current lock value
	currentLock, err := l.redis.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("lock not found or expired")
	}

	// Verify ownership
	if currentLock != lockID {
		return fmt.Errorf("lock owned by different user")
	}

	// Extend the TTL
	return l.redis.Expire(ctx, key, l.ttl)
}
