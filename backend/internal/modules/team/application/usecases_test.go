package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockTeamRepo struct {
	mock.Mock
}

func (m *MockTeamRepo) CreateMatchEvent(ctx context.Context, event *domain.MatchEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockTeamRepo) GetMatchEvent(ctx context.Context, clubID, id string) (*domain.MatchEvent, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MatchEvent), args.Error(1)
}

func (m *MockTeamRepo) SetPlayerAvailability(ctx context.Context, clubID string, availability *domain.PlayerAvailability) error {
	args := m.Called(ctx, clubID, availability)
	return args.Error(0)
}

func (m *MockTeamRepo) GetEventAvailabilities(ctx context.Context, clubID, eventID string) ([]domain.PlayerAvailability, error) {
	args := m.Called(ctx, clubID, eventID)
	return args.Get(0).([]domain.PlayerAvailability), args.Error(1)
}

// --- Tests ---

func TestTeamUseCases_ScheduleMatch(t *testing.T) {
	repo := new(MockTeamRepo)
	uc := application.NewTeamUseCases(repo)
	clubID := "c1"

	t.Run("Success", func(t *testing.T) {
		repo.On("CreateMatchEvent", mock.Anything, mock.MatchedBy(func(e *domain.MatchEvent) bool {
			return e.OpponentName == "Rivals" && e.IsHomeGame == true
		})).Return(nil).Once()

		groupID := uuid.New()
		event, err := uc.ScheduleMatch(context.Background(), clubID, groupID, "Rivals", true, time.Now(), "Stadium")
		assert.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, "Rivals", event.OpponentName)
		repo.AssertExpectations(t)
	})
}

func TestTeamUseCases_RespondAvailability(t *testing.T) {
	repo := new(MockTeamRepo)
	uc := application.NewTeamUseCases(repo)
	clubID := "c1"
	eventID := uuid.New()
	userID := "u1"

	t.Run("Success", func(t *testing.T) {
		repo.On("GetMatchEvent", mock.Anything, clubID, eventID.String()).Return(&domain.MatchEvent{ID: eventID}, nil).Once()
		repo.On("SetPlayerAvailability", mock.Anything, clubID, mock.MatchedBy(func(a *domain.PlayerAvailability) bool {
			return a.Status == domain.AvailabilityConfirmed && a.UserID == userID
		})).Return(nil).Once()

		err := uc.RespondAvailability(context.Background(), clubID, eventID.String(), userID, domain.AvailabilityConfirmed, "")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Fail: Event Not Found", func(t *testing.T) {
		repo.On("GetMatchEvent", mock.Anything, clubID, eventID.String()).Return(nil, errors.New("not found")).Once()

		err := uc.RespondAvailability(context.Background(), clubID, eventID.String(), userID, domain.AvailabilityConfirmed, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
