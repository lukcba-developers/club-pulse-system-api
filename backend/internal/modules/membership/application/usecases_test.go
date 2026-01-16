package application_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---
type MockMembershipRepo struct {
	mock.Mock
}

func (m *MockMembershipRepo) Create(ctx context.Context, membership *domain.Membership) error {
	args := m.Called(ctx, membership)
	return args.Error(0)
}
func (m *MockMembershipRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Membership, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) ListTiers(ctx context.Context, clubID string) ([]domain.MembershipTier, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.MembershipTier), args.Error(1)
}
func (m *MockMembershipRepo) GetTierByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.MembershipTier, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MembershipTier), args.Error(1)
}
func (m *MockMembershipRepo) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) ListBillable(ctx context.Context, clubID string, date time.Time) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, date)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) Update(ctx context.Context, membership *domain.Membership) error {
	args := m.Called(ctx, membership)
	return args.Error(0)
}
func (m *MockMembershipRepo) UpdateBalancesBatch(ctx context.Context, updates map[uuid.UUID]struct {
	Balance     decimal.Decimal
	NextBilling time.Time
}) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}
func (m *MockMembershipRepo) ListAll(ctx context.Context, clubID string) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) GetByUserIDs(ctx context.Context, clubID string, userIDs []uuid.UUID) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, userIDs)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) UpdateBalance(ctx context.Context, clubID string, membershipID uuid.UUID, newBalance decimal.Decimal, nextBilling time.Time) error {
	args := m.Called(ctx, clubID, membershipID, newBalance, nextBilling)
	return args.Error(0)
}

type MockScholarshipRepo struct {
	mock.Mock
}

func (m *MockScholarshipRepo) ListActiveByUserIDs(ctx context.Context, userIDs []string) (map[string]*domain.Scholarship, error) {
	args := m.Called(ctx, userIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*domain.Scholarship), args.Error(1)
}
func (m *MockScholarshipRepo) Create(ctx context.Context, scholarship *domain.Scholarship) error {
	args := m.Called(ctx, scholarship)
	return args.Error(0)
}
func (m *MockScholarshipRepo) GetByUserID(ctx context.Context, userID string) ([]*domain.Scholarship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Scholarship), args.Error(1)
}
func (m *MockScholarshipRepo) GetActiveByUserID(ctx context.Context, userID string) (*domain.Scholarship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Scholarship), args.Error(1)
}

type MockSubscriptionRepo struct {
	mock.Mock
}

func (m *MockSubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	return nil
}
func (m *MockSubscriptionRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Subscription, error) {
	return nil, nil
}
func (m *MockSubscriptionRepo) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Subscription, error) {
	return nil, nil
}
func (m *MockSubscriptionRepo) Update(ctx context.Context, sub *domain.Subscription) error {
	return nil
}

func TestProcessMonthlyBilling(t *testing.T) {
	repo := new(MockMembershipRepo)
	sRepo := new(MockScholarshipRepo)
	subRepo := new(MockSubscriptionRepo)
	uc := application.NewMembershipUseCases(repo, sRepo, subRepo)

	ctx := context.TODO()
	clubID := "club-1"
	userID := uuid.New()
	memID := uuid.New()

	t.Run("Calculates Fees Correctly", func(t *testing.T) {
		// Mock ListBillable
		memberships := []domain.Membership{
			{
				ID: memID, UserID: userID, AutoRenew: true,
				OutstandingBalance: decimal.Zero,
				NextBillingDate:    time.Now().AddDate(0, -1, 0), // Last month
				MembershipTier:     domain.MembershipTier{MonthlyFee: decimal.NewFromFloat(100)},
			},
		}
		repo.On("ListBillable", ctx, clubID, mock.Anything).Return(memberships, nil).Once()

		// Use Explicit Make to ensure type correctness
		scholarships := make(map[string]*domain.Scholarship)
		scholarships[userID.String()] = &domain.Scholarship{
			Percentage: decimal.NewFromFloat(0.5),
			IsActive:   true,
		}
		fmt.Printf("TEST DEBUG: Type of scholarships map: %T\n", scholarships)

		sRepo.On("ListActiveByUserIDs", ctx, []string{userID.String()}).Return(scholarships, nil).Once()

		// Mock UpdateBatch
		repo.On("UpdateBalancesBatch", ctx, mock.MatchedBy(func(updates map[uuid.UUID]struct {
			Balance     decimal.Decimal
			NextBilling time.Time
		}) bool {
			update, ok := updates[memID]
			if !ok {
				return false
			}
			return update.Balance.Equal(decimal.NewFromFloat(50))
		})).Return(nil).Once()

		count, err := uc.ProcessMonthlyBilling(ctx, clubID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		repo.AssertExpectations(t)
	})
}
