package application_test

import (
	"context"
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
	return m.Called(ctx, membership).Error(0)
}
func (m *MockMembershipRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Membership, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) GetByUserIDs(ctx context.Context, clubID string, userIDs []uuid.UUID) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, userIDs)
	return args.Get(0).([]domain.Membership), args.Error(1)
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
func (m *MockMembershipRepo) ListBillable(ctx context.Context, clubID string, date time.Time) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, date)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) Update(ctx context.Context, membership *domain.Membership) error {
	return m.Called(ctx, membership).Error(0)
}
func (m *MockMembershipRepo) UpdateBalance(ctx context.Context, clubID string, membershipID uuid.UUID, newBalance decimal.Decimal, nextBilling time.Time) error {
	return m.Called(ctx, clubID, membershipID, newBalance, nextBilling).Error(0)
}
func (m *MockMembershipRepo) UpdateBalancesBatch(ctx context.Context, updates map[uuid.UUID]struct {
	Balance     decimal.Decimal
	NextBilling time.Time
}) error {
	return m.Called(ctx, updates).Error(0)
}
func (m *MockMembershipRepo) ListAll(ctx context.Context, clubID string) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.Membership), args.Error(1)
}

type MockScholarshipRepo struct {
	mock.Mock
}

func (m *MockScholarshipRepo) Create(ctx context.Context, s *domain.Scholarship) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockScholarshipRepo) GetByUserID(ctx context.Context, userID string) ([]*domain.Scholarship, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Scholarship), args.Error(1)
}
func (m *MockScholarshipRepo) GetActiveByUserID(ctx context.Context, userID string) (*domain.Scholarship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Scholarship), args.Error(1)
}
func (m *MockScholarshipRepo) ListActiveByUserIDs(ctx context.Context, userIDs []string) (map[string]*domain.Scholarship, error) {
	args := m.Called(ctx, userIDs)
	return args.Get(0).(map[string]*domain.Scholarship), args.Error(1)
}

// --- Tests ---

func TestMembershipUseCases_Creation(t *testing.T) {
	mockRepo := new(MockMembershipRepo)
	uc := application.NewMembershipUseCases(mockRepo, nil)
	clubID := "club-1"
	tierID := uuid.New()
	userID := uuid.New()

	t.Run("Create Membership Successful", func(t *testing.T) {
		tier := &domain.MembershipTier{ID: tierID, MonthlyFee: decimal.NewFromInt(100)}
		mockRepo.On("GetTierByID", mock.Anything, clubID, tierID).Return(tier, nil).Once()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		req := application.CreateMembershipRequest{
			UserID:           userID,
			MembershipTierID: tierID,
			BillingCycle:     domain.BillingCycleMonthly,
		}

		m, err := uc.CreateMembership(context.Background(), clubID, req)
		assert.NoError(t, err)
		assert.Equal(t, domain.MembershipStatusActive, m.Status)
	})
}

func TestMembershipUseCases_Cancellation(t *testing.T) {
	mockRepo := new(MockMembershipRepo)
	uc := application.NewMembershipUseCases(mockRepo, nil)
	clubID := "club-1"
	mID := uuid.New()
	uID := uuid.New()

	t.Run("Cancel Success", func(t *testing.T) {
		m := &domain.Membership{ID: mID, UserID: uID, Status: domain.MembershipStatusActive}
		mockRepo.On("GetByID", mock.Anything, clubID, mID).Return(m, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		cancelled, err := uc.CancelMembership(context.Background(), clubID, mID, uID.String())
		assert.NoError(t, err)
		assert.Equal(t, domain.MembershipStatusCancelled, cancelled.Status)
	})
}

func TestMembershipUseCases_Billing(t *testing.T) {
	mockRepo := new(MockMembershipRepo)
	mockScholarRepo := new(MockScholarshipRepo)
	uc := application.NewMembershipUseCases(mockRepo, mockScholarRepo)
	clubID := "club-1"
	uID := uuid.New()
	mID := uuid.New()

	t.Run("Process Monthly Billing with scholarship", func(t *testing.T) {
		tier := domain.MembershipTier{MonthlyFee: decimal.NewFromInt(100)}
		billable := []domain.Membership{
			{
				ID:              mID,
				UserID:          uID,
				MembershipTier:  tier,
				NextBillingDate: time.Now().AddDate(0, 0, -1),
			},
		}

		mockRepo.On("ListBillable", mock.Anything, clubID, mock.Anything).Return(billable, nil).Once()

		scholarships := map[string]*domain.Scholarship{
			uID.String(): {
				Percentage: decimal.NewFromFloat(0.5),
				IsActive:   true,
			},
		}
		mockScholarRepo.On("ListActiveByUserIDs", mock.Anything, []string{uID.String()}).Return(scholarships, nil).Once()

		mockRepo.On("UpdateBalancesBatch", mock.Anything, mock.MatchedBy(func(updates map[uuid.UUID]struct {
			Balance     decimal.Decimal
			NextBilling time.Time
		}) bool {
			// Fee was 100, 50% discount -> 50 Balance update
			return updates[mID].Balance.Equal(decimal.NewFromInt(50))
		})).Return(nil).Once()

		count, err := uc.ProcessMonthlyBilling(context.Background(), clubID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestMembershipUseCases_ScholarshipAssignment(t *testing.T) {
	mockScholarRepo := new(MockScholarshipRepo)
	uc := application.NewMembershipUseCases(nil, mockScholarRepo)
	clubID := "club-1"

	t.Run("Assign scholarship", func(t *testing.T) {
		mockScholarRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		req := application.AssignScholarshipRequest{
			UserID:     uuid.New().String(),
			Percentage: decimal.NewFromFloat(0.2),
			Reason:     "Good student",
		}
		s, err := uc.AssignScholarship(context.Background(), clubID, req, "admin-1")
		assert.NoError(t, err)
		assert.Equal(t, "Good student", s.Reason)
		assert.True(t, s.IsActive)
	})
}
