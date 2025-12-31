package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
	"gorm.io/gorm"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"index"`
	Action    string    `json:"action" gorm:"index"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details" gorm:"type:text"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName specifies the table name for GORM
func (AuditLog) TableName() string {
	return "audit_logs"
}

const (
	auditQueueKey = "audit:queue"
	batchSize     = 100
)

// AuditQueue provides asynchronous audit logging via Redis
type AuditQueue struct {
	redis      *platformRedis.RedisClient
	db         *gorm.DB
	flushEvery time.Duration
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// NewAuditQueue creates a new audit queue
func NewAuditQueue(db *gorm.DB, flushEvery time.Duration) *AuditQueue {
	return &AuditQueue{
		redis:      platformRedis.GetClient(),
		db:         db,
		flushEvery: flushEvery,
		stopChan:   make(chan struct{}),
	}
}

// Push adds an audit log entry to the queue (non-blocking, ~0ms impact on request)
func (q *AuditQueue) Push(ctx context.Context, log *AuditLog) error {
	if log.ID == "" {
		log.ID = fmt.Sprintf("%d-%s", time.Now().UnixNano(), log.UserID[:8])
	}
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal audit log: %w", err)
	}

	return q.redis.LPush(ctx, auditQueueKey, string(data))
}

// StartFlushWorker starts the background worker that flushes logs to PostgreSQL
func (q *AuditQueue) StartFlushWorker() {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		ticker := time.NewTicker(q.flushEvery)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := q.Flush(context.Background()); err != nil {
					log.Printf("Audit queue flush error: %v", err)
				}
			case <-q.stopChan:
				// Final flush before shutdown
				_ = q.Flush(context.Background())
				return
			}
		}
	}()
}

// Stop gracefully stops the flush worker
func (q *AuditQueue) Stop() {
	close(q.stopChan)
	q.wg.Wait()
}

// Flush moves logs from Redis to PostgreSQL
func (q *AuditQueue) Flush(ctx context.Context) error {
	// Get batch of logs from Redis
	entries, err := q.redis.LRange(ctx, auditQueueKey, 0, batchSize-1)
	if err != nil {
		return fmt.Errorf("failed to read audit queue: %w", err)
	}

	if len(entries) == 0 {
		return nil
	}

	// Parse and collect logs
	var logs []AuditLog
	for _, entry := range entries {
		var auditLog AuditLog
		if err := json.Unmarshal([]byte(entry), &auditLog); err != nil {
			log.Printf("Failed to unmarshal audit log: %v", err)
			continue
		}
		logs = append(logs, auditLog)
	}

	// Bulk insert to PostgreSQL
	if len(logs) > 0 {
		if err := q.db.Create(&logs).Error; err != nil {
			return fmt.Errorf("failed to insert audit logs: %w", err)
		}
	}

	// Remove processed entries from Redis
	if err := q.redis.LTrim(ctx, auditQueueKey, int64(len(entries)), -1); err != nil {
		return fmt.Errorf("failed to trim audit queue: %w", err)
	}

	log.Printf("Flushed %d audit logs to PostgreSQL", len(logs))
	return nil
}

// LogAction is a convenience method for logging common actions
func (q *AuditQueue) LogAction(ctx context.Context, userID, action, resource string, details map[string]interface{}, ip, userAgent string) {
	detailsJSON, _ := json.Marshal(details)
	_ = q.Push(ctx, &AuditLog{
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Details:   string(detailsJSON),
		IP:        ip,
		UserAgent: userAgent,
	})
}
