package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

// Session represents a user's active session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	DeviceID  string    `json:"device_id,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	IP        string    `json:"ip,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Store defines the interface for session management
type Store interface {
	Create(ctx context.Context, userID, deviceID, userAgent, ip string) (*Session, error)
	Validate(ctx context.Context, sessionID string) (*Session, error)
	Revoke(ctx context.Context, sessionID string) error
	RevokeAllForUser(ctx context.Context, userID string) error
	ListUserSessions(ctx context.Context, userID string) ([]Session, error)
}

// RedisSessionStore implements Store using Redis
type RedisSessionStore struct {
	redis *platformRedis.RedisClient
	ttl   time.Duration
}

// NewRedisSessionStore creates a new session store
func NewRedisSessionStore(ttl time.Duration) *RedisSessionStore {
	return &RedisSessionStore{
		redis: platformRedis.GetClient(),
		ttl:   ttl,
	}
}

// sessionKey generates the Redis key for a session
func sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

// userSessionsKey generates the Redis key for user's session list
func userSessionsKey(userID string) string {
	return fmt.Sprintf("user_sessions:%s", userID)
}

// Create creates a new session and stores it in Redis
func (s *RedisSessionStore) Create(ctx context.Context, userID, deviceID, userAgent, ip string) (*Session, error) {
	session := &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		DeviceID:  deviceID,
		UserAgent: userAgent,
		IP:        ip,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.ttl),
	}

	data, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	// Store session data
	if err := s.redis.Set(ctx, sessionKey(session.ID), string(data), s.ttl); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	// Add session ID to user's session list (for RevokeAll)
	if err := s.redis.LPush(ctx, userSessionsKey(userID), session.ID); err != nil {
		return nil, fmt.Errorf("failed to add session to user list: %w", err)
	}

	return session, nil
}

// Validate checks if a session is valid and returns it
func (s *RedisSessionStore) Validate(ctx context.Context, sessionID string) (*Session, error) {
	data, err := s.redis.Get(ctx, sessionKey(sessionID))
	if err != nil {
		return nil, fmt.Errorf("session not found or expired")
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.Revoke(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// Revoke invalidates a session by deleting it from Redis
func (s *RedisSessionStore) Revoke(ctx context.Context, sessionID string) error {
	return s.redis.Del(ctx, sessionKey(sessionID))
}

// RevokeAllForUser invalidates all sessions for a user (Global Logout)
func (s *RedisSessionStore) RevokeAllForUser(ctx context.Context, userID string) error {
	// Get all session IDs for the user
	sessionIDs, err := s.redis.LRange(ctx, userSessionsKey(userID), 0, -1)
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Delete each session
	for _, sid := range sessionIDs {
		_ = s.redis.Del(ctx, sessionKey(sid))
	}

	// Clear the user's session list
	return s.redis.Del(ctx, userSessionsKey(userID))
}

// ListUserSessions returns all active sessions for a user
func (s *RedisSessionStore) ListUserSessions(ctx context.Context, userID string) ([]Session, error) {
	sessionIDs, err := s.redis.LRange(ctx, userSessionsKey(userID), 0, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	var sessions []Session
	for _, sid := range sessionIDs {
		session, err := s.Validate(ctx, sid)
		if err != nil {
			// Session expired or invalid, skip it
			continue
		}
		sessions = append(sessions, *session)
	}

	return sessions, nil
}
