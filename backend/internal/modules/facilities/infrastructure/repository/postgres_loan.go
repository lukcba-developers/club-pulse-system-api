package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"gorm.io/gorm"
)

type PostgresLoanRepository struct {
	db *gorm.DB
}

func NewPostgresLoanRepository(db *gorm.DB) *PostgresLoanRepository {
	// AutoMigrate is handled in FacilityRepository constructor usually,
	// but we can add it here safely or ensure it's called centrally.
	// For now, let's assume central migration or ensure struct existence.
	_ = db.AutoMigrate(&EquipmentLoanModel{})
	return &PostgresLoanRepository{db: db}
}

type EquipmentLoanModel struct {
	ID                string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	EquipmentID       string    `gorm:"not null;type:uuid;index"`
	UserID            string    `gorm:"not null;index"`
	LoanedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ExpectedReturnAt  time.Time `gorm:"not null"`
	ReturnedAt        *time.Time
	Status            string `gorm:"default:'ACTIVE'"`
	ConditionOnReturn string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (EquipmentLoanModel) TableName() string {
	return "equipment_loans"
}

func (r *PostgresLoanRepository) Create(ctx context.Context, loan *domain.EquipmentLoan) error {
	model := EquipmentLoanModel{
		ID:               loan.ID,
		EquipmentID:      loan.EquipmentID,
		UserID:           loan.UserID,
		LoanedAt:         loan.LoanedAt,
		ExpectedReturnAt: loan.ExpectedReturnAt,
		Status:           string(loan.Status),
		CreatedAt:        loan.CreatedAt,
		UpdatedAt:        loan.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *PostgresLoanRepository) GetByID(ctx context.Context, id string) (*domain.EquipmentLoan, error) {
	var model EquipmentLoanModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(model), nil
}

func (r *PostgresLoanRepository) ListByUser(ctx context.Context, userID string) ([]*domain.EquipmentLoan, error) {
	var models []EquipmentLoanModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	loans := make([]*domain.EquipmentLoan, len(models))
	for i, m := range models {
		loans[i] = r.toDomain(m)
	}
	return loans, nil
}

func (r *PostgresLoanRepository) ListByStatus(ctx context.Context, status domain.LoanStatus) ([]*domain.EquipmentLoan, error) {
	var models []EquipmentLoanModel
	if err := r.db.WithContext(ctx).Where("status = ?", string(status)).Find(&models).Error; err != nil {
		return nil, err
	}
	loans := make([]*domain.EquipmentLoan, len(models))
	for i, m := range models {
		loans[i] = r.toDomain(m)
	}
	return loans, nil
}

func (r *PostgresLoanRepository) Update(ctx context.Context, loan *domain.EquipmentLoan) error {
	model := EquipmentLoanModel{
		ID:                loan.ID,
		EquipmentID:       loan.EquipmentID,
		UserID:            loan.UserID,
		LoanedAt:          loan.LoanedAt,
		ExpectedReturnAt:  loan.ExpectedReturnAt,
		ReturnedAt:        loan.ReturnedAt,
		Status:            string(loan.Status),
		ConditionOnReturn: loan.ConditionOnReturn,
		CreatedAt:         loan.CreatedAt,
		UpdatedAt:         time.Now(),
	}
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *PostgresLoanRepository) toDomain(m EquipmentLoanModel) *domain.EquipmentLoan {
	return &domain.EquipmentLoan{
		ID:                m.ID,
		EquipmentID:       m.EquipmentID,
		UserID:            m.UserID,
		LoanedAt:          m.LoanedAt,
		ExpectedReturnAt:  m.ExpectedReturnAt,
		ReturnedAt:        m.ReturnedAt,
		Status:            domain.LoanStatus(m.Status),
		ConditionOnReturn: m.ConditionOnReturn,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
