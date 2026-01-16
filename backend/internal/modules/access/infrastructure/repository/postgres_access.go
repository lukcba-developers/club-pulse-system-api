package repository

import (
	"context"
	"errors"

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

func (r *PostgresAccessRepository) GetByEventID(ctx context.Context, clubID string, eventID string) (*domain.AccessLog, error) {
	var log domain.AccessLog
	if err := r.db.WithContext(ctx).Where("event_id = ? AND club_id = ?", eventID, clubID).First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}
