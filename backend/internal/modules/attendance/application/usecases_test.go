package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/domain"
	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockAttendanceRepo struct {
	mock.Mock
}

func (m *MockAttendanceRepo) CreateList(ctx context.Context, list *domain.AttendanceList) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockAttendanceRepo) GetListByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.AttendanceList, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AttendanceList), args.Error(1)
}

func (m *MockAttendanceRepo) GetListByGroupAndDate(ctx context.Context, clubID string, group string, date time.Time) (*domain.AttendanceList, error) {
	args := m.Called(ctx, clubID, group, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AttendanceList), args.Error(1)
}

func (m *MockAttendanceRepo) GetListByTrainingGroupAndDate(ctx context.Context, clubID string, groupID uuid.UUID, date time.Time) (*domain.AttendanceList, error) {
	args := m.Called(ctx, clubID, groupID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AttendanceList), args.Error(1)
}

func (m *MockAttendanceRepo) UpdateRecord(ctx context.Context, record *domain.AttendanceRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockAttendanceRepo) UpsertRecord(ctx context.Context, record *domain.AttendanceRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockAttendanceRepo) GetAttendanceStats(ctx context.Context, clubID, userID string, from, to time.Time) (present, total int, err error) {
	args := m.Called(ctx, clubID, userID, from, to)
	return args.Int(0), args.Int(1), args.Error(2)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) List(ctx context.Context, clubID string, limit, offset int, filters map[string]interface{}) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, limit, offset, filters)
	return args.Get(0).([]userDomain.User), args.Error(1)
}

func (m *MockUserRepo) ListByIDs(ctx context.Context, clubID string, ids []string) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, ids)
	return args.Get(0).([]userDomain.User), args.Error(1)
}

// Minimal implementation of other methods if needed (usually just used as interface)
func (m *MockUserRepo) GetByID(ctx context.Context, clubID, id string) (*userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Update(ctx context.Context, user *userDomain.User) error { return nil }
func (m *MockUserRepo) Delete(ctx context.Context, clubID, id string) error     { return nil }
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
func (m *MockUserRepo) AnonymizeForGDPR(ctx context.Context, clubID, id string) error { return nil }

type MockMembershipRepo struct {
	mock.Mock
}

func (m *MockMembershipRepo) GetByUserIDs(ctx context.Context, clubID string, userIDs []uuid.UUID) ([]membershipDomain.Membership, error) {
	args := m.Called(ctx, clubID, userIDs)
	return args.Get(0).([]membershipDomain.Membership), args.Error(1)
}

// Minimal implementation of other methods
func (m *MockMembershipRepo) Create(ctx context.Context, membership *membershipDomain.Membership) error {
	return nil
}
func (m *MockMembershipRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*membershipDomain.Membership, error) {
	return nil, nil
}
func (m *MockMembershipRepo) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]membershipDomain.Membership, error) {
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
func (m *MockMembershipRepo) UpdateBalance(ctx context.Context, clubID string, membershipID uuid.UUID, newBalance decimal.Decimal, nextBilling time.Time) error {
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

func TestAttendanceUseCases_GetOrCreateList(t *testing.T) {
	repo := new(MockAttendanceRepo)
	userRepo := new(MockUserRepo)
	membershipRepo := new(MockMembershipRepo)
	uc := application.NewAttendanceUseCases(repo, userRepo, membershipRepo)

	ctx := context.TODO()
	clubID := "club-1"
	group := "2012"
	date := time.Now().Truncate(24 * time.Hour)
	coachID := "coach-1"

	t.Run("Returns existing list", func(t *testing.T) {
		existing := &domain.AttendanceList{ID: uuid.New(), Group: group, Date: date}
		repo.On("GetListByGroupAndDate", ctx, clubID, group, date).Return(existing, nil).Once()

		res, err := uc.GetOrCreateList(ctx, clubID, group, date, coachID)
		assert.NoError(t, err)
		assert.Equal(t, existing.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Creates and auto-populates list", func(t *testing.T) {
		repo.On("GetListByGroupAndDate", ctx, clubID, group, date).Return(nil, nil).Once()

		userID := uuid.New().String()
		users := []userDomain.User{{ID: userID}}
		userRepo.On("List", ctx, clubID, 100, 0, map[string]interface{}{"category": group}).Return(users, nil).Once()

		repo.On("CreateList", ctx, mock.Anything).Return(nil).Once()
		repo.On("UpsertRecord", ctx, mock.Anything).Return(nil).Once()

		res, err := uc.GetOrCreateList(ctx, clubID, group, date, coachID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Records, 1)
		assert.Equal(t, userID, res.Records[0].UserID)
		assert.Equal(t, domain.StatusAbsent, res.Records[0].Status)
		repo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})
}

func TestAttendanceUseCases_MarkAttendance(t *testing.T) {
	repo := new(MockAttendanceRepo)
	uc := application.NewAttendanceUseCases(repo, nil, nil)

	ctx := context.TODO()
	clubID := "club-1"
	listID := uuid.New()
	userID := "user-123"

	t.Run("Update existing record", func(t *testing.T) {
		recordID := uuid.New()
		list := &domain.AttendanceList{
			ID: listID,
			Records: []domain.AttendanceRecord{
				{ID: recordID, UserID: userID, Status: domain.StatusAbsent},
			},
		}

		repo.On("GetListByID", ctx, clubID, listID).Return(list, nil).Once()
		repo.On("UpsertRecord", ctx, mock.MatchedBy(func(r *domain.AttendanceRecord) bool {
			return r.ID == recordID && r.UserID == userID && r.Status == domain.StatusPresent
		})).Return(nil).Once()

		err := uc.MarkAttendance(ctx, clubID, listID, application.MarkAttendanceDTO{
			UserID: userID,
			Status: domain.StatusPresent,
		})
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Insert new record into existing list", func(t *testing.T) {
		list := &domain.AttendanceList{ID: listID, Records: []domain.AttendanceRecord{}}
		repo.On("GetListByID", ctx, clubID, listID).Return(list, nil).Once()
		repo.On("UpsertRecord", ctx, mock.MatchedBy(func(r *domain.AttendanceRecord) bool {
			return r.UserID == userID && r.Status == domain.StatusLate
		})).Return(nil).Once()

		err := uc.MarkAttendance(ctx, clubID, listID, application.MarkAttendanceDTO{
			UserID: userID,
			Status: domain.StatusLate,
		})
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
