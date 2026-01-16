package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
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
	ClubID          string    `gorm:"index;not null"` // SECURITY FIX (VUL-002): Tenant isolation
	UserID          uuid.UUID
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

func (r *PostgresSubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	model := r.toModel(subscription)
	return r.db.WithContext(ctx).Create(model).Error
}

// SECURITY FIX (VUL-002): Added clubID parameter for tenant isolation
func (r *PostgresSubscriptionRepository) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Subscription, error) {
	var model SubscriptionModel
	if err := r.db.WithContext(ctx).Scopes(database.TenantScope(clubID)).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

// SECURITY FIX (VUL-002): Added clubID parameter for tenant isolation
func (r *PostgresSubscriptionRepository) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Subscription, error) {
	var models []SubscriptionModel
	if err := r.db.WithContext(ctx).Scopes(database.TenantScope(clubID)).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	var subscriptions []domain.Subscription
	for _, m := range models {
		subscriptions = append(subscriptions, *r.toDomain(&m))
	}
	return subscriptions, nil
}

// SECURITY FIX (VUL-002): Update validates club_id before updating
func (r *PostgresSubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	model := r.toModel(subscription)
	result := r.db.WithContext(ctx).
		Model(&SubscriptionModel{}).
		Where("id = ? AND club_id = ?", model.ID, model.ClubID).
		Updates(map[string]interface{}{
			"status":            model.Status,
			"amount":            model.Amount,
			"next_billing_date": model.NextBillingDate,
			"last_payment_date": model.LastPaymentDate,
			"fail_count":        model.FailCount,
			"updated_at":        time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PostgresSubscriptionRepository) toModel(d *domain.Subscription) *SubscriptionModel {
	return &SubscriptionModel{
		ID:              d.ID,
		ClubID:          d.ClubID, // SECURITY FIX (VUL-002)
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
		ClubID:          m.ClubID, // SECURITY FIX (VUL-002)
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
