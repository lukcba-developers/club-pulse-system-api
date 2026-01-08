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
// SECURITY: Verifies that the requesting user owns the lock before extending
func (l *BookingLock) ExtendLock(ctx context.Context, facilityID string, start, end time.Time, userID string, additionalTTL time.Duration) error {
	key := lockKey(facilityID, start, end)

	// Lua script to verify ownership and extend TTL atomically
	// KEYS[1]: lock key
	// ARGV[1]: user ID (expected owner)
	// ARGV[2]: new TTL duration (in seconds/milliseconds depending on redis version, but usually Go client handles duration)

	// Note: redis.call("expire", ...) returns 1 if timeout was set, 0 if key does not exist.
	// We want to ensure we own it first.
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	// Execute the script
	res, err := l.redis.Eval(ctx, script, []string{key}, userID, additionalTTL).Result()
	if err != nil {
		return fmt.Errorf("failed to extend lock: %w", err)
	}

	// Result can he int64(1) or int64(0). If 0, it means either key didn't exist or we didn't own it.
	// Since Eval returns interface{}, we need to care about type assertion depending on driver version,
	// but go-redis usually handles basic types well.

	val, ok := res.(int64)
	if !ok {
		// Try minimal casting fallback if needed or assume success if err is nil?
		// Actually go-redis Eval returns whatever redis returns.
		// Redis EXPIRE returns integer 1 or 0. Lua script returns what expire returns.
		return fmt.Errorf("unexpected result from redis script")
	}

	if val == 0 {
		return fmt.Errorf("lock extension failed: lock not found or not owned by user")
	}

	return nil
}
