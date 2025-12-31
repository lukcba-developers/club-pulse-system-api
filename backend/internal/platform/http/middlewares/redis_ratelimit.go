package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

// RedisRateLimiter implements rate limiting using Redis sliding window
type RedisRateLimiter struct {
	redis  *platformRedis.RedisClient
	limit  int
	window time.Duration
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		redis:  platformRedis.GetClient(),
		limit:  limit,
		window: window,
	}
}

// Allow checks if a request is allowed based on rate limits
func (r *RedisRateLimiter) Allow(ctx context.Context, ip string) (allowed bool, remaining int) {
	// Key based on IP and current window
	windowStart := time.Now().Unix() / int64(r.window.Seconds())
	key := fmt.Sprintf("ratelimit:%s:%d", ip, windowStart)

	count, err := r.redis.Incr(ctx, key)
	if err != nil {
		// If Redis is unavailable, allow the request (fail open)
		return true, r.limit
	}

	// Set TTL on first request in this window
	if count == 1 {
		_ = r.redis.Expire(ctx, key, r.window)
	}

	remaining = r.limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return count <= int64(r.limit), remaining
}

// Global Redis rate limiter instance
var globalRedisLimiter *RedisRateLimiter

// InitRedisRateLimiter initializes the global Redis rate limiter
func InitRedisRateLimiter(limit int, window time.Duration) {
	globalRedisLimiter = NewRedisRateLimiter(limit, window)
}

// RedisRateLimitMiddleware creates a Gin middleware for Redis-based rate limiting
func RedisRateLimitMiddleware() gin.HandlerFunc {
	// Initialize with defaults if not already initialized
	if globalRedisLimiter == nil {
		InitRedisRateLimiter(100, time.Minute) // 100 requests per minute
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		ctx := context.Background()

		allowed, remaining := globalRedisLimiter.Allow(ctx, ip)

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", globalRedisLimiter.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(globalRedisLimiter.window).Unix()))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"code":    http.StatusTooManyRequests,
				"message": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
