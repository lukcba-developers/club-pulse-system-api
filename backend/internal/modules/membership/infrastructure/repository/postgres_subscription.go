package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PostgresSubscriptionRepository struct {
	db *gorm.DB
}

func NewPostgresSubscriptionRepository(db *gorm.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

type SubscriptionModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          string
	MembershipID    uuid.UUID
	Amount          decimal.Decimal `gorm:"type:decimal(10,2)"`
	Currency        string
	Status          string
	PaymentMethodID string
	NextBillingDate time.Time
	LastPaymentDate *time.Time
	FailCount       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}

func (r *PostgresSubscriptionRepository) Create(subscription *domain.Subscription) error {
	model := r.toModel(subscription)
	return r.db.Create(model).Error
}

func (r *PostgresSubscriptionRepository) GetByID(id uuid.UUID) (*domain.Subscription, error) {
	var model SubscriptionModel
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *PostgresSubscriptionRepository) GetByUserID(userID string) ([]domain.Subscription, error) {
	var models []SubscriptionModel
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	var subscriptions []domain.Subscription
	for _, m := range models {
		subscriptions = append(subscriptions, *r.toDomain(&m))
	}
	return subscriptions, nil
}

func (r *PostgresSubscriptionRepository) Update(subscription *domain.Subscription) error {
	model := r.toModel(subscription)
	return r.db.Save(model).Error
}

func (r *PostgresSubscriptionRepository) toModel(d *domain.Subscription) *SubscriptionModel {
	return &SubscriptionModel{
		ID:              d.ID,
		UserID:          d.UserID,
		MembershipID:    d.MembershipID,
		Amount:          d.Amount,
		Currency:        d.Currency,
		Status:          string(d.Status),
		PaymentMethodID: d.PaymentMethodID,
		NextBillingDate: d.NextBillingDate,
		LastPaymentDate: d.LastPaymentDate,
		FailCount:       d.FailCount,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

func (r *PostgresSubscriptionRepository) toDomain(m *SubscriptionModel) *domain.Subscription {
	return &domain.Subscription{
		ID:              m.ID,
		UserID:          m.UserID,
		MembershipID:    m.MembershipID,
		Amount:          m.Amount,
		Currency:        m.Currency,
		Status:          domain.SubscriptionStatus(m.Status),
		PaymentMethodID: m.PaymentMethodID,
		NextBillingDate: m.NextBillingDate,
		LastPaymentDate: m.LastPaymentDate,
		FailCount:       m.FailCount,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
