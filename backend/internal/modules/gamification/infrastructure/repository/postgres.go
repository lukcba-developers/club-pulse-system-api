package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/gamification/domain"
	"gorm.io/gorm"
)

type PostgresBadgeRepository struct {
	db *gorm.DB
}

func NewPostgresBadgeRepository(db *gorm.DB) *PostgresBadgeRepository {
	return &PostgresBadgeRepository{db: db}
}

// Badge CRUD

func (r *PostgresBadgeRepository) Create(ctx context.Context, badge *domain.Badge) error {
	return r.db.WithContext(ctx).Create(badge).Error
}

func (r *PostgresBadgeRepository) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Badge, error) {
	var badge domain.Badge
	err := r.db.WithContext(ctx).Where("id = ? AND club_id = ?", id, clubID).First(&badge).Error
	return &badge, err
}

func (r *PostgresBadgeRepository) GetByCode(ctx context.Context, clubID, code string) (*domain.Badge, error) {
	var badge domain.Badge
	err := r.db.WithContext(ctx).Where("code = ? AND club_id = ?", code, clubID).First(&badge).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if explicitly not found
		}
		return nil, err
	}
	return &badge, nil
}

func (r *PostgresBadgeRepository) List(ctx context.Context, clubID string) ([]domain.Badge, error) {
	var badges []domain.Badge
	err := r.db.WithContext(ctx).Where("club_id = ?", clubID).Find(&badges).Error
	return badges, err
}

func (r *PostgresBadgeRepository) Update(ctx context.Context, badge *domain.Badge) error {
	return r.db.WithContext(ctx).Where("id = ? AND club_id = ?", badge.ID, badge.ClubID).Updates(badge).Error
}

// UserBadge operations with strict tenant isolation

func (r *PostgresBadgeRepository) AwardBadge(ctx context.Context, userBadge *domain.UserBadge) error {
	// Note: Caller should verify badge ownership/existence before calling this if needed.
	// We trust the service layer has fetched the badge ID correctly scoped to the club.
	return r.db.WithContext(ctx).Create(userBadge).Error
}

func (r *PostgresBadgeRepository) GetUserBadges(ctx context.Context, clubID, userID string) ([]domain.UserBadge, error) {
	var badges []domain.UserBadge
	// Join with Badges table to enforce ClubID
	err := r.db.WithContext(ctx).Table("user_badges").
		Joins("JOIN badges ON badges.id = user_badges.badge_id").
		Where("user_badges.user_id = ? AND badges.club_id = ?", userID, clubID).
		Find(&badges).Error
	return badges, err
}

func (r *PostgresBadgeRepository) HasBadge(ctx context.Context, clubID, userID string, badgeCode string) (bool, error) {
	var count int64
	// Check existence joining badge table
	err := r.db.WithContext(ctx).Table("user_badges").
		Joins("JOIN badges ON badges.id = user_badges.badge_id").
		Where("user_badges.user_id = ? AND badges.code = ? AND badges.club_id = ?", userID, badgeCode, clubID).
		Count(&count).Error
	return count > 0, err
}

func (r *PostgresBadgeRepository) SetFeatured(ctx context.Context, clubID, userID string, badgeID uuid.UUID, featured bool) error {
	// Only allow setting featured if the badge belongs to the club (security check)
	// We do this by verifying existence via JOIN or Subquery in the update condition.
	// But GORM Updates with Join is tricky.

	// Option 1: Verify first
	var count int64
	if err := r.db.WithContext(ctx).Table("user_badges").
		Joins("JOIN badges ON badges.id = user_badges.badge_id").
		Where("user_badges.user_id = ? AND user_badges.badge_id = ? AND badges.club_id = ?", userID, badgeID, clubID).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound // Or access denied
	}

	return r.db.WithContext(ctx).Model(&domain.UserBadge{}).
		Where("user_id = ? AND badge_id = ?", userID, badgeID).
		Update("featured", featured).Error
}

func (r *PostgresBadgeRepository) GetFeaturedBadges(ctx context.Context, clubID, userID string) ([]domain.UserBadge, error) {
	var badges []domain.UserBadge
	err := r.db.WithContext(ctx).Table("user_badges").
		Joins("JOIN badges ON badges.id = user_badges.badge_id").
		Where("user_badges.user_id = ? AND user_badges.featured = ? AND badges.club_id = ?", userID, true, clubID).
		Find(&badges).Error
	return badges, err
}
