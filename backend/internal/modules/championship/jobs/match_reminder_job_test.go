package jobs_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/jobs"
	notificationSvc "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateTournament(ctx context.Context, tournament *domain.Tournament) error {
	return nil
}
func (m *MockRepo) GetTournament(ctx context.Context, clubID, id string) (*domain.Tournament, error) {
	return nil, nil
}
func (m *MockRepo) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	return nil, nil
}
func (m *MockRepo) CreateStage(ctx context.Context, stage *domain.TournamentStage) error { return nil }
func (m *MockRepo) GetStage(ctx context.Context, clubID, id string) (*domain.TournamentStage, error) {
	return nil, nil
}
func (m *MockRepo) CreateGroup(ctx context.Context, group *domain.Group) error { return nil }
func (m *MockRepo) GetGroup(ctx context.Context, clubID, id string) (*domain.Group, error) {
	return nil, nil
}
func (m *MockRepo) CreateMatch(ctx context.Context, clubID string, match *domain.TournamentMatch) error {
	return nil
}
func (m *MockRepo) CreateMatchesBatch(ctx context.Context, clubID string, matches []domain.TournamentMatch) error {
	return nil
}
func (m *MockRepo) GetMatch(ctx context.Context, clubID, id string) (*domain.TournamentMatch, error) {
	return nil, nil
}
func (m *MockRepo) GetMatchesByGroup(ctx context.Context, clubID, groupID string) ([]domain.TournamentMatch, error) {
	return nil, nil
}
func (m *MockRepo) UpdateMatchResult(ctx context.Context, clubID, matchID string, homeScore, awayScore float64) error {
	return nil
}
func (m *MockRepo) UpdateMatchScheduling(ctx context.Context, clubID, matchID string, date time.Time, bookingID uuid.UUID) error {
	return nil
}
func (m *MockRepo) GetStandings(ctx context.Context, clubID, groupID string) ([]domain.Standing, error) {
	return nil, nil
}
func (m *MockRepo) RegisterTeam(ctx context.Context, clubID string, standing *domain.Standing) error {
	return nil
}
func (m *MockRepo) UpdateStanding(ctx context.Context, standing *domain.Standing) error { return nil }
func (m *MockRepo) UpdateStandingsBatch(ctx context.Context, clubID string, standings []domain.Standing) error {
	return nil
}
func (m *MockRepo) CreateTeam(ctx context.Context, team *domain.Team) error    { return nil }
func (m *MockRepo) AddMember(ctx context.Context, teamID, userID string) error { return nil }
func (m *MockRepo) GetMatchesByUserID(ctx context.Context, clubID, userID string) ([]domain.TournamentMatch, error) {
	return nil, nil
}

func (m *MockRepo) GetUpcomingMatches(ctx context.Context, clubID string, from, to time.Time) ([]domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, from, to)
	return args.Get(0).([]domain.TournamentMatch), args.Error(1)
}

func (m *MockRepo) GetTeamMembers(ctx context.Context, teamID string) ([]string, error) {
	args := m.Called(ctx, teamID)
	return args.Get(0).([]string), args.Error(1)
}

// Mock Notification Provider
type MockEmailProvider struct{}

func (m *MockEmailProvider) SendEmail(ctx context.Context, to, subject, body string) (*notificationSvc.DeliveryResult, error) {
	return &notificationSvc.DeliveryResult{Success: true, ProviderID: "email-id"}, nil
}
func (m *MockEmailProvider) SendEmailWithTemplate(ctx context.Context, to, templateID string, data map[string]interface{}) (*notificationSvc.DeliveryResult, error) {
	return &notificationSvc.DeliveryResult{Success: true, ProviderID: "email-tpl-id"}, nil
}

type MockSMSProvider struct{}

func (m *MockSMSProvider) SendSMS(ctx context.Context, to, body string) (*notificationSvc.DeliveryResult, error) {
	return &notificationSvc.DeliveryResult{Success: true, ProviderID: "sms-id"}, nil
}

func TestMatchReminderJob(t *testing.T) {
	mockRepo := new(MockRepo)
	// Create service using real constructor but with mock providers (if interface allowed)
	// Since Update of NotificationService logic wasn't fully refactored to interface-based sender injection for Push,
	// we assume standard behavior. But checking log output is hard.
	// Ideally NotificationService should be interface. But struct is used.
	// For now we trust "Print" logic or check if we can mock providers.

	// We'll create a minimal notification service
	notifService := notificationSvc.NewNotificationService(&MockEmailProvider{}, &MockSMSProvider{})

	job := jobs.NewMatchReminderJob(mockRepo, notifService, 24)

	clubID := "test-club"

	// Setup Expectations
	matches := []domain.TournamentMatch{
		{
			ID:           uuid.New(),
			HomeTeamID:   uuid.New(),
			AwayTeamID:   uuid.New(),
			AwayTeamName: "AwayTeam",
			Date:         time.Now().Add(10 * time.Hour),
		},
	}

	mockRepo.On("GetUpcomingMatches", mock.Anything, clubID, mock.Anything, mock.Anything).Return(matches, nil)
	mockRepo.On("GetTeamMembers", mock.Anything, matches[0].HomeTeamID.String()).Return([]string{"user1"}, nil)
	mockRepo.On("GetTeamMembers", mock.Anything, matches[0].AwayTeamID.String()).Return([]string{"user2"}, nil)

	// Run
	err := job.Run(context.Background(), clubID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
