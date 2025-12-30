package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"gorm.io/gorm"
)

type PostgresMembershipRepository struct {
	db *gorm.DB
}

func NewPostgresMembershipRepository(db *gorm.DB) *PostgresMembershipRepository {
	// AutoMigrate tables for MVP
	db.AutoMigrate(&domain.MembershipTier{}, &domain.Membership{})
	return &PostgresMembershipRepository{db: db}
}

func (r *PostgresMembershipRepository) Create(ctx context.Context, membership *domain.Membership) error {
	return r.db.WithContext(ctx).Create(membership).Error
}

func (r *PostgresMembershipRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Membership, error) {
	var membership domain.Membership
	if err := r.db.WithContext(ctx).Preload("MembershipTier").First(&membership, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("membership not found")
		}
		return nil, err
	}
	return &membership, nil
}

func (r *PostgresMembershipRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Membership, error) {
	var memberships []domain.Membership
	if err := r.db.WithContext(ctx).Preload("MembershipTier").Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *PostgresMembershipRepository) ListTiers(ctx context.Context) ([]domain.MembershipTier, error) {
	var tiers []domain.MembershipTier
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("monthly_fee asc").Find(&tiers).Error; err != nil {
		return nil, err
	}
	return tiers, nil
}

func (r *PostgresMembershipRepository) GetTierByID(ctx context.Context, id uuid.UUID) (*domain.MembershipTier, error) {
	var tier domain.MembershipTier
	if err := r.db.WithContext(ctx).First(&tier, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("membership tier not found")
		}
		return nil, err
	}
	return &tier, nil
}
