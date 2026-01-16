package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockDisciplineRepo struct {
	mock.Mock
}

func (m *MockDisciplineRepo) CreateDiscipline(ctx context.Context, d *domain.Discipline) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}
func (m *MockDisciplineRepo) ListDisciplines(ctx context.Context, clubID string) ([]domain.Discipline, error) {
	return nil, nil
}
func (m *MockDisciplineRepo) GetDisciplineByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Discipline, error) {
	return nil, nil
}
func (m *MockDisciplineRepo) CreateGroup(ctx context.Context, g *domain.TrainingGroup) error {
	return nil
}
func (m *MockDisciplineRepo) ListGroups(ctx context.Context, clubID string, filter map[string]interface{}) ([]domain.TrainingGroup, error) {
	return nil, nil
}
func (m *MockDisciplineRepo) GetGroupByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.TrainingGroup, error) {
	return nil, nil
}

type MockTournamentRepo struct {
	mock.Mock
}

func (m *MockTournamentRepo) CreateTournament(ctx context.Context, t *domain.Tournament) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *MockTournamentRepo) GetTournamentByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Tournament, error) {
	return nil, nil
}
func (m *MockTournamentRepo) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	return nil, nil
}
func (m *MockTournamentRepo) UpdateTournament(ctx context.Context, t *domain.Tournament) error {
	return nil
}
func (m *MockTournamentRepo) CreateTeam(ctx context.Context, t *domain.Team) error {
	return nil
}
func (m *MockTournamentRepo) GetTeamByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Team, error) {
	return nil, nil
}
func (m *MockTournamentRepo) ListTeams(ctx context.Context, clubID string, tournamentID uuid.UUID) ([]domain.Team, error) {
	return nil, nil
}
func (m *MockTournamentRepo) CreateMatch(ctx context.Context, match *domain.Match) error {
	return nil
}
func (m *MockTournamentRepo) UpdateMatch(ctx context.Context, match *domain.Match) error {
	return nil
}
func (m *MockTournamentRepo) GetMatchByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Match, error) {
	return nil, nil
}
func (m *MockTournamentRepo) ListMatches(ctx context.Context, clubID string, tournamentID uuid.UUID) ([]domain.Match, error) {
	return nil, nil
}
func (m *MockTournamentRepo) GetStandings(ctx context.Context, clubID string, tournamentID uuid.UUID) ([]domain.Standing, error) {
	return nil, nil
}

type MockUserRepo struct{ mock.Mock }

// Implement stubs if needed for ListStudentsInGroup
func (m *MockUserRepo) List(ctx context.Context, clubID string, limit, offset int, filters map[string]interface{}) ([]userDomain.User, error) {
	return nil, nil
}

// Add other user stubs if interface is larger, but usecase likely uses limited set.
// Checking userRepo usage in usecases.go -> List.
func (m *MockUserRepo) Create(ctx context.Context, user *userDomain.User) error { return nil }
func (m *MockUserRepo) GetByID(ctx context.Context, clubID, id string) (*userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, clubID, email string) (*userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Update(ctx context.Context, user *userDomain.User) error { return nil }
func (m *MockUserRepo) Delete(ctx context.Context, clubID, id string) error     { return nil }
func (m *MockUserRepo) ListByIDs(ctx context.Context, clubID string, ids []string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FindChildren(ctx context.Context, clubID, parentID string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) CreateIncident(ctx context.Context, incident *userDomain.IncidentLog) error {
	return nil
}
func (m *MockUserRepo) AnonymizeForGDPR(ctx context.Context, clubID, id string) error {
	return nil
}

// --- Tests ---

func TestDisciplineUseCases_CreateDiscipline(t *testing.T) {
	dRepo := new(MockDisciplineRepo)
	tRepo := new(MockTournamentRepo)
	uRepo := new(MockUserRepo)

	uc := application.NewDisciplineUseCases(dRepo, tRepo, uRepo)
	clubID := "c1"

	t.Run("Success", func(t *testing.T) {
		dRepo.On("CreateDiscipline", mock.Anything, mock.MatchedBy(func(d *domain.Discipline) bool {
			return d.Name == "Tennis" && d.ClubID == clubID
		})).Return(nil).Once()

		res, err := uc.CreateDiscipline(context.TODO(), clubID, "Tennis", "Desc")
		assert.NoError(t, err)
		assert.Equal(t, "Tennis", res.Name)
		dRepo.AssertExpectations(t)
	})
}

func TestDisciplineUseCases_CreateTournament(t *testing.T) {
	dRepo := new(MockDisciplineRepo)
	tRepo := new(MockTournamentRepo)
	uRepo := new(MockUserRepo)

	uc := application.NewDisciplineUseCases(dRepo, tRepo, uRepo)
	clubID := "c1"
	dID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		tRepo.On("CreateTournament", mock.Anything, mock.MatchedBy(func(tr *domain.Tournament) bool {
			return tr.Name == "Summer Cup" && tr.DisciplineID == dID
		})).Return(nil).Once()

		res, err := uc.CreateTournament(context.TODO(), clubID, "Summer Cup", dID.String(), time.Now(), time.Now().Add(24*time.Hour), "LEAGUE")
		assert.NoError(t, err)
		assert.Equal(t, "Summer Cup", res.Name)
		tRepo.AssertExpectations(t)
	})

	t.Run("Fail: Invalid UUID", func(t *testing.T) {
		_, err := uc.CreateTournament(context.TODO(), clubID, "Fail", "invalid-uuid", time.Now(), time.Now(), "LEAGUE")
		assert.Error(t, err)
	})
}
