package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/domain"
	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockAccessRepo struct {
	mock.Mock
}

func (m *MockAccessRepo) Create(ctx context.Context, log *domain.AccessLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAccessRepo) GetByEventID(ctx context.Context, clubID, eventID string) (*domain.AccessLog, error) {
	args := m.Called(ctx, clubID, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AccessLog), args.Error(1)
}

// Unused but required
func (m *MockAccessRepo) List(ctx context.Context, clubID string, filter map[string]interface{}, limit, offset int) ([]domain.AccessLog, error) {
	return nil, nil
}

func (m *MockAccessRepo) GetByUserID(ctx context.Context, clubID, userID string, limit int) ([]domain.AccessLog, error) {
	return nil, nil
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, clubID, id string) (*userDomain.User, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

// Stub other methods
func (m *MockUserRepo) Update(ctx context.Context, user *userDomain.User) error { return nil }
func (m *MockUserRepo) Delete(ctx context.Context, clubID, id string) error     { return nil }
func (m *MockUserRepo) List(ctx context.Context, clubID string, limit, offset int, filters map[string]interface{}) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FindChildren(ctx context.Context, clubID, parentID string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Create(ctx context.Context, user *userDomain.User) error { return nil }
func (m *MockUserRepo) CreateIncident(ctx context.Context, incident *userDomain.IncidentLog) error {
	return nil
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, clubID, email string) (*userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) ListByIDs(ctx context.Context, clubID string, ids []string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) AnonymizeForGDPR(ctx context.Context, clubID, id string) error { return nil }

type MockMembershipRepo struct {
	mock.Mock
}

func (m *MockMembershipRepo) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]membershipDomain.Membership, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]membershipDomain.Membership), args.Error(1)
}

// Stub others
func (m *MockMembershipRepo) Create(ctx context.Context, membership *membershipDomain.Membership) error {
	return nil
}
func (m *MockMembershipRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*membershipDomain.Membership, error) {
	return nil, nil
}
func (m *MockMembershipRepo) ListTiers(ctx context.Context, clubID string) ([]membershipDomain.MembershipTier, error) {
	return nil, nil
}
func (m *MockMembershipRepo) GetTierByID(ctx context.Context, clubID string, id uuid.UUID) (*membershipDomain.MembershipTier, error) {
	return nil, nil
}
func (m *MockMembershipRepo) ListBillable(ctx context.Context, clubID string, date time.Time) ([]membershipDomain.Membership, error) {
	return nil, nil
}
func (m *MockMembershipRepo) Update(ctx context.Context, membership *membershipDomain.Membership) error {
	return nil
}
func (m *MockMembershipRepo) UpdateBalancesBatch(ctx context.Context, updates map[uuid.UUID]struct {
	Balance     decimal.Decimal
	NextBilling time.Time
}) error {
	return nil
}
func (m *MockMembershipRepo) ListAll(ctx context.Context, clubID string) ([]membershipDomain.Membership, error) {
	return nil, nil
}
func (m *MockMembershipRepo) GetByUserIDs(ctx context.Context, clubID string, userIDs []uuid.UUID) ([]membershipDomain.Membership, error) {
	return nil, nil
}
func (m *MockMembershipRepo) UpdateBalance(ctx context.Context, clubID string, membershipID uuid.UUID, newBalance decimal.Decimal, nextBilling time.Time) error {
	return nil
}

func TestRequestEntry(t *testing.T) {
	ar := new(MockAccessRepo)
	ur := new(MockUserRepo)
	mr := new(MockMembershipRepo)
	uc := application.NewAccessUseCases(ar, ur, mr)

	clubID := "c1"
	userID := uuid.New()
	eventID := "evt-1"

	t.Run("Success: Access Granted", func(t *testing.T) {
		// Idempotency check
		ar.On("GetByEventID", mock.Anything, clubID, eventID).Return(nil, nil).Once()
		// User check
		ur.On("GetByID", mock.Anything, clubID, userID.String()).Return(&userDomain.User{ID: userID.String()}, nil).Once()
		// Membership check
		mr.On("GetByUserID", mock.Anything, clubID, userID).Return([]membershipDomain.Membership{
			{Status: membershipDomain.MembershipStatusActive, OutstandingBalance: decimal.Zero},
		}, nil).Once()
		// Log creation
		ar.On("Create", mock.Anything, mock.MatchedBy(func(l *domain.AccessLog) bool {
			return l.Status == domain.AccessStatusGranted && l.UserID == userID.String()
		})).Return(nil).Once()

		req := application.EntryRequest{UserID: userID.String(), EventID: eventID, Direction: "IN"}
		log, err := uc.RequestEntry(context.TODO(), clubID, req)
		assert.NoError(t, err)
		assert.Equal(t, domain.AccessStatusGranted, log.Status)
	})

	t.Run("Denied: Outstanding Debt", func(t *testing.T) {
		eID := "evt-2"
		ar.On("GetByEventID", mock.Anything, clubID, eID).Return(nil, nil).Once()
		ur.On("GetByID", mock.Anything, clubID, userID.String()).Return(&userDomain.User{ID: userID.String()}, nil).Once()
		mr.On("GetByUserID", mock.Anything, clubID, userID).Return([]membershipDomain.Membership{
			{Status: membershipDomain.MembershipStatusActive, OutstandingBalance: decimal.NewFromFloat(100)},
		}, nil).Once()

		ar.On("Create", mock.Anything, mock.MatchedBy(func(l *domain.AccessLog) bool {
			return l.Status == domain.AccessStatusDenied && l.Reason == "Outstanding debt"
		})).Return(nil).Once()

		req := application.EntryRequest{UserID: userID.String(), EventID: eID, Direction: "IN"}
		log, err := uc.RequestEntry(context.TODO(), clubID, req)
		assert.NoError(t, err)
		assert.Equal(t, domain.AccessStatusDenied, log.Status)
	})

	t.Run("Denied: User Not Found", func(t *testing.T) {
		eID := "evt-3"
		unknownID := uuid.New().String()
		ar.On("GetByEventID", mock.Anything, clubID, eID).Return(nil, nil).Once()
		ur.On("GetByID", mock.Anything, clubID, unknownID).Return(nil, nil).Once()

		ar.On("Create", mock.Anything, mock.MatchedBy(func(l *domain.AccessLog) bool {
			return l.Status == domain.AccessStatusDenied && l.Reason == "User not found"
		})).Return(nil).Once()

		req := application.EntryRequest{UserID: unknownID, EventID: eID, Direction: "IN"}
		log, err := uc.RequestEntry(context.TODO(), clubID, req)
		assert.NoError(t, err)
		assert.Equal(t, domain.AccessStatusDenied, log.Status)
	})

	t.Run("Idempotency: Return Existing Log", func(t *testing.T) {
		eID := "evt-4"
		existing := &domain.AccessLog{ID: uuid.New(), Status: domain.AccessStatusGranted, EventID: eID}
		ar.On("GetByEventID", mock.Anything, clubID, eID).Return(existing, nil).Once()

		req := application.EntryRequest{UserID: userID.String(), EventID: eID, Direction: "IN"}
		log, err := uc.RequestEntry(context.TODO(), clubID, req)
		assert.NoError(t, err)
		assert.Equal(t, existing.ID, log.ID)
	})
}
