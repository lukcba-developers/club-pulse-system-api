package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *RedisClient
	once   sync.Once
)

// RedisClient wraps the go-redis client with helper methods
type RedisClient struct {
	rdb *redis.Client
}

// InitRedis initializes the Redis connection
func InitRedis() {
	once.Do(func() {
		host := getEnv("REDIS_HOST", "localhost")
		port := getEnv("REDIS_PORT", "6379")
		password := getEnv("REDIS_PASSWORD", "")

		rdb := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Password:     password,
			DB:           0,
			PoolSize:     100,
			MinIdleConns: 10,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		// Test connection with retries
		ctx := context.Background()
		for i := 0; i < 5; i++ {
			_, err := rdb.Ping(ctx).Result()
			if err == nil {
				break
			}
			log.Printf("Failed to connect to Redis, retrying in 2 seconds... (%d/5)", i+1)
			time.Sleep(2 * time.Second)
		}

		// Final ping (will panic if still not connected)
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			log.Printf("Warning: Redis not available, some features disabled: %v", err)
		} else {
			log.Println("Redis connection established successfully")
		}

		client = &RedisClient{rdb: rdb}
	})
}

// GetClient returns the Redis client instance
func GetClient() *RedisClient {
	if client == nil {
		InitRedis()
	}
	return client
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.rdb.Close()
}

// --- Basic Operations ---

// Set stores a key-value with optional TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

// Del deletes one or more keys
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	return r.rdb.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.rdb.Exists(ctx, key).Result()
	return result > 0, err
}

// --- Rate Limiting Operations ---

// Incr increments a key and returns the new value
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.rdb.Incr(ctx, key).Result()
}

// Expire sets a TTL on a key
func (r *RedisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.rdb.Expire(ctx, key, ttl).Err()
}

// --- Session Operations ---

// SetNX sets a key only if it doesn't exist (for locks)
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return r.rdb.SetNX(ctx, key, value, ttl).Result()
}

// Scan iterates over keys matching a pattern
func (r *RedisClient) Scan(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	iter := r.rdb.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	return keys, iter.Err()
}

// --- Pub/Sub Operations ---

// Publish sends a message to a channel
func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.rdb.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to a channel and returns a PubSub
func (r *RedisClient) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.rdb.Subscribe(ctx, channel)
}

// Eval executes a Lua script
func (r *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return r.rdb.Eval(ctx, script, keys, args...)
}

// --- List Operations (for Audit Queue) ---

// LPush pushes a value to the head of a list
func (r *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.rdb.LPush(ctx, key, values...).Err()
}

// LRange gets a range of values from a list
func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.rdb.LRange(ctx, key, start, stop).Result()
}

// LTrim trims a list to the specified range
func (r *RedisClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	return r.rdb.LTrim(ctx, key, start, stop).Err()
}

// --- Health Check ---

// Ping checks if Redis is available
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
