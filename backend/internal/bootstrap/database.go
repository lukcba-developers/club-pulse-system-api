package bootstrap

import (
	"fmt"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/audit"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/logger"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"

	// Domains for migration (only if strictly necessary to keep AutoMigrate for now, generally we'd move this)

	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	paymentDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"

	"gorm.io/gorm"
)

type Infrastructure struct {
	DB         *gorm.DB
	Redis      *platformRedis.RedisClient
	AuditQueue *audit.AuditQueue
}

func InitInfrastructure() (*Infrastructure, error) {
	// 1. Initialize Database
	database.InitDB()
	db := database.GetDB()

	// 2. Initialize Redis
	platformRedis.InitRedis()
	redisClient := platformRedis.GetClient()
	logger.Info("Redis client initialized")

	// 3. Initialize Audit Queue
	// Using a 5-minute flush interval as per previous configuration
	auditQueue := audit.NewAuditQueue(db, 5*time.Minute)
	auditQueue.StartFlushWorker()
	logger.Info("Audit queue worker started (flush every 5 min)")

	// 4. Run basic AutoMigrate (Ideally this should be separate, but keeping parity for now)
	// We only migrate Payment here because it was explicitly in main.go
	if err := db.AutoMigrate(&paymentDomain.Payment{}); err != nil {
		logger.Error(fmt.Sprintf("Failed to migrate payment table: %v", err))
	}
	if err := db.AutoMigrate(&bookingDomain.RecurringRule{}); err != nil {
		logger.Error(fmt.Sprintf("Failed to migrate recurring_rule table: %v", err))
	}

	return &Infrastructure{
		DB:         db,
		Redis:      redisClient,
		AuditQueue: auditQueue,
	}, nil
}

func (i *Infrastructure) Shutdown() {
	if i.AuditQueue != nil {
		i.AuditQueue.Stop()
	}
	// Close DB/Redis connections if necessary (GORM manages pool, Redis client has Close)
	if i.Redis != nil {
		i.Redis.Close()
	}
}
