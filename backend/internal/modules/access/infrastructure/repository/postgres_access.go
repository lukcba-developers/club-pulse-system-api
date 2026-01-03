package repository

import (
	"context"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/domain"
	"gorm.io/gorm"
)

type PostgresAccessRepository struct {
	db *gorm.DB
}

func NewPostgresAccessRepository(db *gorm.DB) *PostgresAccessRepository {
	return &PostgresAccessRepository{db: db}
}

func (r *PostgresAccessRepository) Create(ctx context.Context, log *domain.AccessLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *PostgresAccessRepository) GetByUserID(ctx context.Context, clubID string, userID string, limit int) ([]domain.AccessLog, error) {
	var logs []domain.AccessLog
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND club_id = ?", userID, clubID).
		Order("timestamp desc").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
