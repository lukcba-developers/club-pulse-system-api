package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PostgresMembershipRepository struct {
	db *gorm.DB
}

func NewPostgresMembershipRepository(db *gorm.DB) *PostgresMembershipRepository {
	// AutoMigrate tables for MVP
	if err := db.AutoMigrate(&domain.MembershipTier{}, &domain.Membership{}); err != nil {
		// In a real app, we might panic or log fatal here, ensuring DB is consistent
		panic("failed to migrate membership tables: " + err.Error())
	}
	return &PostgresMembershipRepository{db: db}
}

func (r *PostgresMembershipRepository) Create(ctx context.Context, membership *domain.Membership) error {
	return r.db.WithContext(ctx).Create(membership).Error
}

func (r *PostgresMembershipRepository) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Membership, error) {
	var membership domain.Membership
	if err := r.db.WithContext(ctx).Preload("MembershipTier").First(&membership, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("membership not found")
		}
		return nil, err
	}
	return &membership, nil
}

func (r *PostgresMembershipRepository) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Membership, error) {
	var memberships []domain.Membership
	if err := r.db.WithContext(ctx).Preload("MembershipTier").Where("user_id = ? AND club_id = ?", userID, clubID).Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *PostgresMembershipRepository) ListTiers(ctx context.Context, clubID string) ([]domain.MembershipTier, error) {
	var tiers []domain.MembershipTier
	if err := r.db.WithContext(ctx).Where("is_active = ? AND club_id = ?", true, clubID).Order("monthly_fee asc").Find(&tiers).Error; err != nil {
		return nil, err
	}
	return tiers, nil
}

func (r *PostgresMembershipRepository) GetTierByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.MembershipTier, error) {
	var tier domain.MembershipTier
	if err := r.db.WithContext(ctx).First(&tier, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("membership tier not found")
		}
		return nil, err
	}
	return &tier, nil
}

func (r *PostgresMembershipRepository) ListBillable(ctx context.Context, clubID string, date time.Time) ([]domain.Membership, error) {
	var memberships []domain.Membership
	// Status Active AND NextBillingDate <= today
	if err := r.db.WithContext(ctx).
		Preload("MembershipTier").
		Where("status = ? AND next_billing_date <= ? AND club_id = ?", domain.MembershipStatusActive, date, clubID).
		Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *PostgresMembershipRepository) UpdateBalance(ctx context.Context, clubID string, membershipID uuid.UUID, newBalance decimal.Decimal, nextBilling time.Time) error {
	updates := map[string]interface{}{
		"outstanding_balance": newBalance,
		"next_billing_date":   nextBilling,
		"updated_at":          time.Now(),
	}
	return r.db.WithContext(ctx).Model(&domain.Membership{}).Where("id = ? AND club_id = ?", membershipID, clubID).Updates(updates).Error
}
