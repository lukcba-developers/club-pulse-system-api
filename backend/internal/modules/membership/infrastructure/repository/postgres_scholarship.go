package repository

import (
	"errors"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PostgresScholarshipRepository struct {
	db *gorm.DB
}

func NewPostgresScholarshipRepository(db *gorm.DB) *PostgresScholarshipRepository {
	_ = db.AutoMigrate(&ScholarshipModel{})
	return &PostgresScholarshipRepository{db: db}
}

type ScholarshipModel struct {
	ID         string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID     string          `gorm:"not null;index"`
	Percentage decimal.Decimal `gorm:"type:decimal(5,2);not null"`
	Reason     string
	GrantorID  string
	ValidUntil *time.Time
	IsActive   bool `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (ScholarshipModel) TableName() string {
	return "scholarships"
}

func (r *PostgresScholarshipRepository) Create(scholarship *domain.Scholarship) error {
	model := ScholarshipModel{
		ID:         scholarship.ID,
		UserID:     scholarship.UserID,
		Percentage: scholarship.Percentage,
		Reason:     scholarship.Reason,
		GrantorID:  scholarship.GrantorID,
		ValidUntil: scholarship.ValidUntil,
		IsActive:   scholarship.IsActive,
		CreatedAt:  scholarship.CreatedAt,
		UpdatedAt:  scholarship.UpdatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *PostgresScholarshipRepository) GetByUserID(userID string) ([]*domain.Scholarship, error) {
	var models []ScholarshipModel
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	scholarships := make([]*domain.Scholarship, len(models))
	for i, m := range models {
		scholarships[i] = r.toDomain(m)
	}
	return scholarships, nil
}

func (r *PostgresScholarshipRepository) GetActiveByUserID(userID string) (*domain.Scholarship, error) {
	var model ScholarshipModel
	// Find first active scholarship that is either not expired or has no expiry date
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Where("valid_until IS NULL OR valid_until > ?", time.Now()).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(model), nil
}

func (r *PostgresScholarshipRepository) toDomain(m ScholarshipModel) *domain.Scholarship {
	return &domain.Scholarship{
		ID:         m.ID,
		UserID:     m.UserID,
		Percentage: m.Percentage,
		Reason:     m.Reason,
		GrantorID:  m.GrantorID,
		ValidUntil: m.ValidUntil,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
