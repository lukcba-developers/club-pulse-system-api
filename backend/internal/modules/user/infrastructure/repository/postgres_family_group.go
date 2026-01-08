package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"gorm.io/gorm"
)

type PostgresFamilyGroupRepository struct {
	db *gorm.DB
}

func NewPostgresFamilyGroupRepository(db *gorm.DB) *PostgresFamilyGroupRepository {
	// AutoMigrate FamilyGroup
	_ = db.AutoMigrate(&domain.FamilyGroup{})
	return &PostgresFamilyGroupRepository{db: db}
}

func (r *PostgresFamilyGroupRepository) Create(group *domain.FamilyGroup) error {
	if group.ID == uuid.Nil {
		group.ID = uuid.New()
	}
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	return r.db.Create(group).Error
}

func (r *PostgresFamilyGroupRepository) GetByID(clubID string, id uuid.UUID) (*domain.FamilyGroup, error) {
	var group domain.FamilyGroup
	result := r.db.Preload("Members").Where("id = ? AND club_id = ?", id, clubID).First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &group, nil
}

func (r *PostgresFamilyGroupRepository) GetByHeadUserID(clubID, headUserID string) (*domain.FamilyGroup, error) {
	var group domain.FamilyGroup
	result := r.db.Preload("Members").Where("head_user_id = ? AND club_id = ?", headUserID, clubID).First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &group, nil
}

func (r *PostgresFamilyGroupRepository) GetByMemberID(clubID, userID string) (*domain.FamilyGroup, error) {
	// Find user first to get their family_group_id
	var user struct {
		FamilyGroupID *uuid.UUID
	}
	result := r.db.Table("users").Select("family_group_id").Where("id = ? AND club_id = ?", userID, clubID).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.FamilyGroupID == nil {
		return nil, nil
	}
	return r.GetByID(clubID, *user.FamilyGroupID)
}

func (r *PostgresFamilyGroupRepository) AddMember(clubID string, groupID uuid.UUID, userID string) error {
	// Verify group exists and belongs to club
	group, err := r.GetByID(clubID, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return errors.New("family group not found")
	}

	// Update user's family_group_id
	return r.db.Table("users").
		Where("id = ? AND club_id = ?", userID, clubID).
		Update("family_group_id", groupID).Error
}

func (r *PostgresFamilyGroupRepository) RemoveMember(clubID string, groupID uuid.UUID, userID string) error {
	// Verify group exists
	group, err := r.GetByID(clubID, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return errors.New("family group not found")
	}

	// Cannot remove head user
	if group.HeadUserID == userID {
		return errors.New("cannot remove head user from family group")
	}

	// Set user's family_group_id to NULL
	return r.db.Table("users").
		Where("id = ? AND club_id = ? AND family_group_id = ?", userID, clubID, groupID).
		Update("family_group_id", nil).Error
}
