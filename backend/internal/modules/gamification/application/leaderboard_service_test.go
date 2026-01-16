package application_test

import (
	"context"
	"testing"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/gamification/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/gamification/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of userDomain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, clubID, id string) (*userDomain.User, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *userDomain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, clubID, id string) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, clubID string, limit, offset int, filters map[string]interface{}) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userDomain.User), args.Error(1)
}

func (m *MockUserRepository) ListByIDs(ctx context.Context, clubID string, ids []string) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userDomain.User), args.Error(1)
}

func (m *MockUserRepository) FindChildren(ctx context.Context, clubID, parentID string) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userDomain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *userDomain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateIncident(ctx context.Context, incident *userDomain.IncidentLog) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, clubID, email string) (*userDomain.User, error) {
	args := m.Called(ctx, clubID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) AnonymizeForGDPR(ctx context.Context, clubID, id string) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

func TestLeaderboardService_GetGlobalLeaderboard(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewLeaderboardService(mockRepo)
	ctx := context.TODO()
	clubID := "club-1"

	users := []userDomain.User{
		{ID: "u1", Name: "Alice", Stats: &userDomain.UserStats{TotalXP: 100, Level: 2}},
		{ID: "u2", Name: "Bob", Stats: &userDomain.UserStats{TotalXP: 200, Level: 3}},
		{ID: "u3", Name: "Charlie", Stats: &userDomain.UserStats{TotalXP: 50, Level: 1}},
	}

	mockRepo.On("List", ctx, clubID, 1000, 0, mock.Anything).Return(users, nil)

	leaderboard, err := service.GetGlobalLeaderboard(ctx, clubID, domain.LeaderboardPeriodAllTime, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, leaderboard)
	assert.Equal(t, domain.LeaderboardTypeGlobal, leaderboard.Type)
	assert.Equal(t, 3, len(leaderboard.Entries))

	// Verify sorting (Desc by TotalXP)
	assert.Equal(t, "Bob", leaderboard.Entries[0].UserName)
	assert.Equal(t, float64(200), leaderboard.Entries[0].Score)
	assert.Equal(t, 1, leaderboard.Entries[0].Rank)

	assert.Equal(t, "Alice", leaderboard.Entries[1].UserName)
	assert.Equal(t, float64(100), leaderboard.Entries[1].Score)
	assert.Equal(t, 2, leaderboard.Entries[1].Rank)

	assert.Equal(t, "Charlie", leaderboard.Entries[2].UserName)
	assert.Equal(t, float64(50), leaderboard.Entries[2].Score)
	assert.Equal(t, 3, leaderboard.Entries[2].Rank)
}

func TestLeaderboardService_GetBookingsLeaderboard(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewLeaderboardService(mockRepo)
	ctx := context.TODO()
	clubID := "club-1"

	users := []userDomain.User{
		{ID: "u1", Name: "Alice", Stats: &userDomain.UserStats{MatchesPlayed: 5}},
		{ID: "u2", Name: "Bob", Stats: &userDomain.UserStats{MatchesPlayed: 10}},
		{ID: "u3", Name: "Charlie", Stats: &userDomain.UserStats{MatchesPlayed: 2}},
	}

	mockRepo.On("List", ctx, clubID, 1000, 0, mock.Anything).Return(users, nil)

	leaderboard, err := service.GetBookingsLeaderboard(ctx, clubID, domain.LeaderboardPeriodAllTime, 10)

	assert.NoError(t, err)
	assert.NotNil(t, leaderboard)
	assert.Equal(t, domain.LeaderboardTypeBookings, leaderboard.Type)
	assert.Equal(t, 3, len(leaderboard.Entries))

	// Verify sorting (Desc by MatchesPlayed)
	assert.Equal(t, "Bob", leaderboard.Entries[0].UserName)
	assert.Equal(t, float64(10), leaderboard.Entries[0].Score)

	assert.Equal(t, "Alice", leaderboard.Entries[1].UserName)
	assert.Equal(t, float64(5), leaderboard.Entries[1].Score)

	assert.Equal(t, "Charlie", leaderboard.Entries[2].UserName)
	assert.Equal(t, float64(2), leaderboard.Entries[2].Score)
}

func TestLeaderboardService_GetUserRank(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewLeaderboardService(mockRepo)
	ctx := context.TODO()
	clubID := "club-1"

	users := []userDomain.User{
		{ID: "u1", Name: "Alice", Stats: &userDomain.UserStats{TotalXP: 100}},
		{ID: "u2", Name: "Bob", Stats: &userDomain.UserStats{TotalXP: 200}},
		{ID: "u3", Name: "Charlie", Stats: &userDomain.UserStats{TotalXP: 50}},
	}

	mockRepo.On("List", ctx, clubID, 1000, 0, mock.Anything).Return(users, nil)

	rank, err := service.GetUserRank(ctx, clubID, "u1", domain.LeaderboardTypeGlobal, domain.LeaderboardPeriodAllTime)
	assert.NoError(t, err)
	assert.Equal(t, 2, rank) // Alice is 2nd (200 > 100 > 50)
}
