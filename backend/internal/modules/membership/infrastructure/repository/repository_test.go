package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- SQLite Compatible Models ---

type TestMembershipTier struct {
	ID         uuid.UUID       `gorm:"type:uuid;primary_key"`
	ClubID     string          `gorm:"index;not null"`
	Name       string          `gorm:"not null"`
	MonthlyFee decimal.Decimal `gorm:"type:decimal(10,2)"`
	IsActive   bool            `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (TestMembershipTier) TableName() string { return "membership_tiers" }

type TestMembership struct {
	ID                 uuid.UUID       `gorm:"type:uuid;primary_key"`
	ClubID             string          `gorm:"index;not null"`
	UserID             uuid.UUID       `gorm:"type:uuid;not null"`
	MembershipTierID   uuid.UUID       `gorm:"type:uuid;not null"`
	Status             string          `gorm:"not null"`
	BillingCycle       string          `gorm:"not null"`
	StartDate          time.Time       `gorm:"not null"`
	NextBillingDate    time.Time       `gorm:"not null"`
	OutstandingBalance decimal.Decimal `gorm:"type:decimal(10,2)"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (TestMembership) TableName() string { return "memberships" }

type TestScholarship struct {
	ID         string          `gorm:"primaryKey"`
	UserID     string          `gorm:"not null;index"`
	Percentage decimal.Decimal `gorm:"type:decimal(5,2);not null"`
	Reason     string
	GrantorID  string
	ValidUntil *time.Time
	IsActive   bool `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (TestScholarship) TableName() string { return "scholarships" }

type TestSubscription struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
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

func (TestSubscription) TableName() string { return "subscriptions" }

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	db.AutoMigrate(&TestMembershipTier{}, &TestMembership{}, &TestScholarship{}, &TestSubscription{})
	return db
}

func TestPostgresMembershipRepositories(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewPostgresMembershipRepository(db)
	scholarRepo := repository.NewPostgresScholarshipRepository(db)
	subRepo := repository.NewPostgresSubscriptionRepository(db)
	clubID := "club-r-final-2"

	t.Run("Full Membership Lifecycle", func(t *testing.T) {
		tierID := uuid.New()
		db.Create(&TestMembershipTier{ID: tierID, ClubID: clubID, Name: "Platinum", MonthlyFee: decimal.NewFromInt(200)})

		m := &domain.Membership{
			ID:               uuid.New(),
			ClubID:           clubID,
			UserID:           uuid.New(),
			MembershipTierID: tierID,
			Status:           domain.MembershipStatusActive,
			NextBillingDate:  time.Now().AddDate(0, 0, -1),
		}
		_ = repo.Create(context.Background(), m)

		billable, _ := repo.ListBillable(context.Background(), clubID, time.Now())
		assert.NotEmpty(t, billable)
		_ = repo.UpdateBalance(context.Background(), clubID, m.ID, decimal.NewFromInt(200), time.Now().AddDate(0, 1, 0))
		saved, _ := repo.GetByID(context.Background(), clubID, m.ID)
		assert.Equal(t, int64(200), saved.OutstandingBalance.IntPart())
	})

	t.Run("Scholarship Lifecycle", func(t *testing.T) {
		uID := uuid.New().String()
		// 1. Create Active
		s := &domain.Scholarship{
			ID:         uuid.New().String(),
			UserID:     uID,
			Percentage: decimal.NewFromFloat(0.2),
			IsActive:   true,
			Reason:     "Merit",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := scholarRepo.Create(context.Background(), s)
		assert.NoError(t, err)

		// 2. Get Active
		active, err := scholarRepo.GetActiveByUserID(context.Background(), uID)
		assert.NoError(t, err)
		assert.NotNil(t, active)
		assert.Equal(t, s.ID, active.ID)

		// 3. Create Expired
		expiredTime := time.Now().Add(-24 * time.Hour)
		sExpired := &domain.Scholarship{
			ID:         uuid.New().String(),
			UserID:     uID,
			Percentage: decimal.NewFromFloat(0.5),
			IsActive:   true,
			ValidUntil: &expiredTime,
		}
		_ = scholarRepo.Create(context.Background(), sExpired)

		// 4. Get Active (Should still return the first valid one if logic holds, or just verify list)
		// The GetActiveByUserID returns the *first* matching.
		// Let's verify ListActiveByUserIDs
		activeMap, err := scholarRepo.ListActiveByUserIDs(context.Background(), []string{uID})
		assert.NoError(t, err)
		assert.NotNil(t, activeMap[uID])
		assert.Equal(t, s.ID, activeMap[uID].ID) // Should be the valid one
	})

	t.Run("Subscription Lifecycle", func(t *testing.T) {
		uID := uuid.New().String()
		mID := uuid.New()

		sub := &domain.Subscription{
			ID:              uuid.New(),
			UserID:          uID,
			MembershipID:    mID,
			Amount:          decimal.NewFromInt(50),
			Currency:        "USD",
			Status:          domain.SubscriptionActive,
			NextBillingDate: time.Now().AddDate(0, 1, 0),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		err := subRepo.Create(context.Background(), sub)
		assert.NoError(t, err)

		saved, err := subRepo.GetByID(context.Background(), sub.ID)
		assert.NoError(t, err)
		assert.Equal(t, sub.Amount.String(), saved.Amount.String())

		// Update
		sub.Status = domain.SubscriptionCancelled
		err = subRepo.Update(context.Background(), sub)
		assert.NoError(t, err)

		updated, err := subRepo.GetByID(context.Background(), sub.ID)
		assert.NoError(t, err)
		assert.Equal(t, domain.SubscriptionCancelled, updated.Status)

		// Get By User
		list, err := subRepo.GetByUserID(context.Background(), uID)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("Queries and Errors", func(t *testing.T) {
		// 1. Not Found
		m, err := repo.GetByID(context.Background(), clubID, uuid.New())
		assert.Error(t, err)
		assert.Nil(t, m)

		tErr, err := repo.GetTierByID(context.Background(), clubID, uuid.New())
		assert.Error(t, err)
		assert.Nil(t, tErr)

		// 2. List Tiers
		tiers, err := repo.ListTiers(context.Background(), clubID)
		assert.NoError(t, err)
		// Assuming setupDB doesn't seed, but Full Lifecycle might have added one
		// We can add one here to be sure
		db.Create(&TestMembershipTier{ID: uuid.New(), ClubID: clubID, Name: "Gold", MonthlyFee: decimal.NewFromInt(100), IsActive: true})
		tiers, _ = repo.ListTiers(context.Background(), clubID)
		assert.NotEmpty(t, tiers)

		// 3. List All
		all, err := repo.ListAll(context.Background(), clubID)
		assert.NoError(t, err)
		assert.NotEmpty(t, all)

		// 4. Get By UserID
		uID := uuid.New()
		_ = repo.Create(context.Background(), &domain.Membership{
			ID:               uuid.New(),
			ClubID:           clubID,
			UserID:           uID,
			MembershipTierID: uuid.New(), // dummy
			Status:           domain.MembershipStatusActive,
			NextBillingDate:  time.Now(),
		})
		byUser, err := repo.GetByUserID(context.Background(), clubID, uID)
		assert.NoError(t, err)
		assert.Len(t, byUser, 1)

		// 5. Get By UserIDs batch
		byUsers, err := repo.GetByUserIDs(context.Background(), clubID, []uuid.UUID{uID})
		assert.NoError(t, err)
		assert.Len(t, byUsers, 1)
	})

	t.Run("Batch Updates", func(t *testing.T) {
		m1 := &domain.Membership{ID: uuid.New(), ClubID: clubID, UserID: uuid.New(), Status: domain.MembershipStatusActive}
		m2 := &domain.Membership{ID: uuid.New(), ClubID: clubID, UserID: uuid.New(), Status: domain.MembershipStatusActive}
		_ = repo.Create(context.Background(), m1)
		_ = repo.Create(context.Background(), m2)

		updates := map[uuid.UUID]struct {
			Balance     decimal.Decimal
			NextBilling time.Time
		}{
			m1.ID: {Balance: decimal.NewFromInt(50), NextBilling: time.Now().AddDate(0, 1, 0)},
			m2.ID: {Balance: decimal.NewFromInt(75), NextBilling: time.Now().AddDate(0, 1, 0)},
		}

		// This exercises the fallback logic for SQLite in postgres.go
		err := repo.UpdateBalancesBatch(context.Background(), updates)
		assert.NoError(t, err)

		saved1, _ := repo.GetByID(context.Background(), clubID, m1.ID)
		assert.Equal(t, "50", saved1.OutstandingBalance.String())

		saved2, _ := repo.GetByID(context.Background(), clubID, m2.ID)
		assert.Equal(t, "75", saved2.OutstandingBalance.String())
	})

	t.Run("Expanded Coverage", func(t *testing.T) {
		// 1. GetTierByID Success
		tierID := uuid.New()
		db.Create(&TestMembershipTier{ID: tierID, ClubID: clubID, Name: "Diamond", MonthlyFee: decimal.NewFromInt(300), IsActive: true})
		tier, err := repo.GetTierByID(context.Background(), clubID, tierID)
		assert.NoError(t, err)
		assert.Equal(t, "Diamond", tier.Name)

		// 2. ListBillable Edge Cases
		// Create a member with billing date in future
		futureDate := time.Now().AddDate(0, 1, 0)
		mFuture := &domain.Membership{ID: uuid.New(), ClubID: clubID, UserID: uuid.New(), Status: domain.MembershipStatusActive, NextBillingDate: futureDate, MembershipTierID: tierID}
		_ = repo.Create(context.Background(), mFuture)

		billable, err := repo.ListBillable(context.Background(), clubID, time.Now())
		assert.NoError(t, err)
		// Should NOT contain mFuture
		for _, b := range billable {
			assert.NotEqual(t, mFuture.ID, b.ID)
		}

		// 3. Scholarship Expiry Logic Check
		// Create active scholarship with future expiry
		uID := uuid.New().String()
		validUntil := time.Now().Add(24 * time.Hour)
		sValid := &domain.Scholarship{ID: uuid.New().String(), UserID: uID, IsActive: true, ValidUntil: &validUntil, Percentage: decimal.NewFromFloat(0.1)}
		_ = scholarRepo.Create(context.Background(), sValid)

		active, err := scholarRepo.GetActiveByUserID(context.Background(), uID)
		assert.NoError(t, err)
		assert.NotNil(t, active)
		assert.Equal(t, sValid.ID, active.ID)
	})
}
