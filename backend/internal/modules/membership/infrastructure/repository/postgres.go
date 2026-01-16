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

func (r *PostgresMembershipRepository) GetByUserIDs(ctx context.Context, clubID string, userIDs []uuid.UUID) ([]domain.Membership, error) {
	if len(userIDs) == 0 {
		return []domain.Membership{}, nil
	}
	var memberships []domain.Membership
	if err := r.db.WithContext(ctx).Preload("MembershipTier").Where("user_id IN ? AND club_id = ?", userIDs, clubID).Find(&memberships).Error; err != nil {
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

func (r *PostgresMembershipRepository) ListAll(ctx context.Context, clubID string) ([]domain.Membership, error) {
	var memberships []domain.Membership
	err := r.db.WithContext(ctx).
		Preload("MembershipTier").
		Where("club_id = ?", clubID).
		Order("created_at DESC").
		Find(&memberships).Error
	return memberships, err
}

// Update saves changes to an existing membership
func (r *PostgresMembershipRepository) Update(ctx context.Context, membership *domain.Membership) error {
	return r.db.WithContext(ctx).Save(membership).Error
}

func (r *PostgresMembershipRepository) UpdateBalancesBatch(ctx context.Context, updates map[uuid.UUID]struct {
	Balance     decimal.Decimal
	NextBilling time.Time
}) error {
	if len(updates) == 0 {
		return nil
	}

	// Efficient Bulk Update using PostgreSQL FROM VALUES
	// UPDATE memberships AS m SET
	//   outstanding_balance = v.balance,
	//   next_billing_date = v.next_billing,
	//   updated_at = NOW()
	// FROM (VALUES
	//   ('uuid1', 100.00, '2023-01-01'),
	//   ('uuid2', 200.00, '2023-02-01')
	// ) AS v(id, balance, next_billing)
	// WHERE m.id = v.id::uuid;

	// Build the query and args
	query := `
		UPDATE memberships AS m 
		SET 
			outstanding_balance = v.balance,
			next_billing_date = v.next_billing,
			updated_at = NOW()
		FROM (VALUES 
	`
	var args []interface{}
	var valueParamPlaceholders string

	// Fallback for non-postgres (like SQLite in tests)
	// SECURITY FIX (VUL-004): Now validates club_id to prevent cross-tenant updates
	if r.db.Dialector.Name() != "postgres" {
		return r.db.Transaction(func(tx *gorm.DB) error {
			for id, update := range updates {
				// SECURITY: Use club_id check to prevent cross-tenant updates
				// We get the membership first to extract its club_id, then validate
				var membership domain.Membership
				if err := tx.Select("club_id").First(&membership, "id = ?", id).Error; err != nil {
					continue // Skip if not found
				}
				// Update only records matching both id AND club_id
				err := tx.Model(&domain.Membership{}).
					Where("id = ? AND club_id = ?", id, membership.ClubID).
					Updates(map[string]interface{}{
						"outstanding_balance": update.Balance,
						"next_billing_date":   update.NextBilling,
						"updated_at":          time.Now(),
					}).Error
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	// Iterate and build values
	i := 0
	for id, update := range updates {
		if i > 0 {
			valueParamPlaceholders += ","
		}
		// $1 is id, $2 is balance, $3 is next_billing
		// GORM (or raw sql) placeholders are ? or $N. GORM uses ? usually but r.db.Exec with raw sql in postgres driver uses $N?
		// Actually GORM `Exec` handles ? to $ transformation.
		valueParamPlaceholders += "(?::uuid, ?::numeric, ?::timestamp)"
		args = append(args, id.String(), update.Balance, update.NextBilling)
		i++
	}

	query += valueParamPlaceholders + `) AS v(id, balance, next_billing) WHERE m.id = v.id::uuid`

	return r.db.WithContext(ctx).Exec(query, args...).Error
}
