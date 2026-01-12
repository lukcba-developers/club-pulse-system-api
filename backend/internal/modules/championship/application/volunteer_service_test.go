package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVolunteerRepo struct {
	mock.Mock
}

func (m *MockVolunteerRepo) Create(ctx context.Context, v *domain.VolunteerAssignment) error {
	return m.Called(ctx, v).Error(0)
}

func (m *MockVolunteerRepo) GetByMatchID(ctx context.Context, clubID string, matchID uuid.UUID) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, matchID)
	var res []domain.VolunteerAssignment
	if args.Get(0) != nil {
		res = args.Get(0).([]domain.VolunteerAssignment)
	}
	return res, args.Error(1)
}

func (m *MockVolunteerRepo) GetByUserID(ctx context.Context, clubID, userID string) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, userID)
	var res []domain.VolunteerAssignment
	if args.Get(0) != nil {
		res = args.Get(0).([]domain.VolunteerAssignment)
	}
	return res, args.Error(1)
}

func (m *MockVolunteerRepo) GetByRoleAndMatch(ctx context.Context, clubID string, matchID uuid.UUID, role domain.VolunteerRole) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, matchID, role)
	var res []domain.VolunteerAssignment
	if args.Get(0) != nil {
		res = args.Get(0).([]domain.VolunteerAssignment)
	}
	return res, args.Error(1)
}

func (m *MockVolunteerRepo) Update(ctx context.Context, v *domain.VolunteerAssignment) error {
	return m.Called(ctx, v).Error(0)
}

func (m *MockVolunteerRepo) Delete(ctx context.Context, clubID string, id uuid.UUID) error {
	return m.Called(ctx, clubID, id).Error(0)
}

func TestVolunteerService(t *testing.T) {
	repo := new(MockVolunteerRepo)
	svc := application.NewVolunteerService(repo)
	ctx := context.TODO()
	cID := "club-1"
	mID := uuid.New()

	t.Run("AssignVolunteer Success", func(t *testing.T) {
		repo.On("Create", ctx, mock.Anything).Return(nil).Once()
		err := svc.AssignVolunteer(ctx, cID, mID, "u1", "REFREE", "admin", "notes")
		assert.NoError(t, err)
	})

	t.Run("GetMatchVolunteers Success", func(t *testing.T) {
		repo.On("GetByMatchID", ctx, cID, mID).Return([]domain.VolunteerAssignment{{UserID: "u1"}}, nil).Once()
		res, err := svc.GetMatchVolunteers(ctx, cID, mID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GetVolunteerSummary Success", func(t *testing.T) {
		repo.On("GetByMatchID", ctx, cID, mID).Return([]domain.VolunteerAssignment{
			{UserID: "u1", Role: "REFREE"},
			{UserID: "u2", Role: "BALL_BOY"},
		}, nil).Once()
		res, err := svc.GetVolunteerSummary(ctx, cID, mID)
		assert.NoError(t, err)
		assert.Equal(t, 2, res.FilledSlots)
		assert.Equal(t, 1, res.ByRole["REFREE"])
	})

	t.Run("GetVolunteerSummary Error", func(t *testing.T) {
		repo.On("GetByMatchID", ctx, cID, mID).Return(nil, errors.New("db error")).Once()
		_, err := svc.GetVolunteerSummary(ctx, cID, mID)
		assert.Error(t, err)
	})

	t.Run("RemoveVolunteer Success", func(t *testing.T) {
		aID := uuid.New()
		repo.On("Delete", ctx, cID, aID).Return(nil).Once()
		err := svc.RemoveVolunteer(ctx, cID, aID)
		assert.NoError(t, err)
	})

	t.Run("GetUserAssignments Success", func(t *testing.T) {
		repo.On("GetByUserID", ctx, cID, "u1").Return([]domain.VolunteerAssignment{{MatchID: mID}}, nil).Once()
		res, err := svc.GetUserAssignments(ctx, cID, "u1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("ValidateAssignment Success", func(t *testing.T) {
		repo.On("GetByRoleAndMatch", ctx, cID, mID, domain.VolunteerRole("REFREE")).Return([]domain.VolunteerAssignment{}, nil).Once()
		err := svc.ValidateAssignment(ctx, cID, mID, "REFREE", 2)
		assert.NoError(t, err)
	})

	t.Run("ValidateAssignment LimitReached", func(t *testing.T) {
		repo.On("GetByRoleAndMatch", ctx, cID, mID, domain.VolunteerRole("REFREE")).Return([]domain.VolunteerAssignment{{ID: uuid.New()}}, nil).Once()
		err := svc.ValidateAssignment(ctx, cID, mID, "REFREE", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "máximo")
	})

	t.Run("ValidateAssignment RepoError", func(t *testing.T) {
		repo.On("GetByRoleAndMatch", ctx, cID, mID, domain.VolunteerRole("REFREE")).Return(nil, errors.New("db error")).Once()
		err := svc.ValidateAssignment(ctx, cID, mID, "REFREE", 2)
		assert.Error(t, err)
	})
}
