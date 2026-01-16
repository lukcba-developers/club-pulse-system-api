package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"gorm.io/gorm"
)

type PostgresPaymentRepository struct {
	db *gorm.DB
}

func NewPostgresPaymentRepository(db *gorm.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

func (r *PostgresPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// Update updates a payment record.
// SECURITY FIX (VUL-003): Now validates that the payment belongs to the club before updating.
func (r *PostgresPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	// Validate that the payment exists and belongs to the club
	result := r.db.WithContext(ctx).
		Model(&domain.Payment{}).
		Where("id = ? AND club_id = ?", payment.ID, payment.ClubID).
		Updates(map[string]interface{}{
			"status":      payment.Status,
			"paid_at":     payment.PaidAt,
			"external_id": payment.ExternalID,
			"updated_at":  payment.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("payment not found or does not belong to this club")
	}

	return nil
}

// SECURITY FIX (VUL-003): Added clubID parameter for tenant isolation
func (r *PostgresPaymentRepository) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Scopes(database.TenantScope(clubID)).First(&payment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// SECURITY FIX (VUL-003): Added clubID parameter for tenant isolation
func (r *PostgresPaymentRepository) GetByExternalID(ctx context.Context, clubID string, externalID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Scopes(database.TenantScope(clubID)).First(&payment, "external_id = ?", externalID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// GetByExternalIDForWebhook is ONLY for webhook processing where clubID is unknown.
// SECURITY: This is safe because webhook signature is validated before this call,
// and the clubID from the returned payment is used for subsequent tenant-scoped operations.
func (r *PostgresPaymentRepository) GetByExternalIDForWebhook(ctx context.Context, externalID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).First(&payment, "external_id = ?", externalID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PostgresPaymentRepository) List(ctx context.Context, clubID string, filter domain.PaymentFilter) ([]*domain.Payment, int64, error) {
	var payments []*domain.Payment
	var total int64

	// Create base query with mandatory club_id
	query := r.db.WithContext(ctx).Model(&domain.Payment{}).Scopes(database.TenantScope(clubID))

	// Apply filters if present
	if filter.PayerID != uuid.Nil {
		query = query.Where("payer_id = ?", filter.PayerID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// Count total matching records before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Execute query ordered by newest first
	if err := query.Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
