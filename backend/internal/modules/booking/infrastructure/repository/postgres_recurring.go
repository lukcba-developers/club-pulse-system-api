package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	"gorm.io/gorm"
)

type PostgresRecurringRepository struct {
	db *gorm.DB
}

func NewPostgresRecurringRepository(db *gorm.DB) *PostgresRecurringRepository {
	return &PostgresRecurringRepository{db: db}
}

func (r *PostgresRecurringRepository) Create(ctx context.Context, rule *domain.RecurringRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *PostgresRecurringRepository) GetByFacility(ctx context.Context, clubID string, facilityID uuid.UUID) ([]domain.RecurringRule, error) {
	var rules []domain.RecurringRule
	err := r.db.WithContext(ctx).
		Where("club_id = ? AND facility_id = ?", clubID, facilityID).
		Find(&rules).Error
	return rules, err
}

// GetAllActive returns rules that are currently valid (EndDate >= Today)
func (r *PostgresRecurringRepository) GetAllActive(ctx context.Context, clubID string) ([]domain.RecurringRule, error) {
	var rules []domain.RecurringRule
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Where("club_id = ? AND end_date >= ?", clubID, today).
		Find(&rules).Error
	return rules, err
}

func (r *PostgresRecurringRepository) Delete(ctx context.Context, clubID string, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.RecurringRule{}, "id = ? AND club_id = ?", id, clubID).Error
}
