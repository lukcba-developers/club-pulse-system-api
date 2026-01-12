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

func (m *MockDisciplineRepo) CreateDiscipline(ctx context.Context, discipline *domain.Discipline) error {
	args := m.Called(ctx, discipline)
	return args.Error(0)
}

func (m *MockDisciplineRepo) ListDisciplines(ctx context.Context, clubID string) ([]domain.Discipline, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.Discipline), args.Error(1)
}

func (m *MockDisciplineRepo) GetDisciplineByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Discipline, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Discipline), args.Error(1)
}

func (m *MockDisciplineRepo) CreateGroup(ctx context.Context, group *domain.TrainingGroup) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockDisciplineRepo) ListGroups(ctx context.Context, clubID string, filter map[string]interface{}) ([]domain.TrainingGroup, error) {
	args := m.Called(ctx, clubID, filter)
	return args.Get(0).([]domain.TrainingGroup), args.Error(1)
}

func (m *MockDisciplineRepo) GetGroupByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.TrainingGroup, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TrainingGroup), args.Error(1)
}

type MockTournamentRepo struct {
	mock.Mock
}

func (m *MockTournamentRepo) CreateTournament(ctx context.Context, t *domain.Tournament) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTournamentRepo) GetTournamentByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Tournament, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tournament), args.Error(1)
}

func (m *MockTournamentRepo) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.Tournament), args.Error(1)
}

func (m *MockTournamentRepo) UpdateTournament(ctx context.Context, t *domain.Tournament) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTournamentRepo) CreateTeam(ctx context.Context, team *domain.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTournamentRepo) GetTeamByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Team, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockTournamentRepo) ListTeams(ctx context.Context, clubID string, tID uuid.UUID) ([]domain.Team, error) {
	args := m.Called(ctx, clubID, tID)
	return args.Get(0).([]domain.Team), args.Error(1)
}

func (m *MockTournamentRepo) CreateMatch(ctx context.Context, match *domain.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockTournamentRepo) UpdateMatch(ctx context.Context, match *domain.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockTournamentRepo) GetMatchByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Match, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Match), args.Error(1)
}

func (m *MockTournamentRepo) ListMatches(ctx context.Context, clubID string, tID uuid.UUID) ([]domain.Match, error) {
	args := m.Called(ctx, clubID, tID)
	return args.Get(0).([]domain.Match), args.Error(1)
}

func (m *MockTournamentRepo) GetStandings(ctx context.Context, clubID string, tID uuid.UUID) ([]domain.Standing, error) {
	args := m.Called(ctx, clubID, tID)
	return args.Get(0).([]domain.Standing), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *userDomain.User) error { return nil }
func (m *MockUserRepo) GetByID(ctx context.Context, clubID, id string) (*userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Update(ctx context.Context, user *userDomain.User) error { return nil }
func (m *MockUserRepo) List(ctx context.Context, clubID string, limit, offset int, filters map[string]interface{}) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, limit, offset, filters)
	return args.Get(0).([]userDomain.User), args.Error(1)
}
func (m *MockUserRepo) ListByIDs(ctx context.Context, clubID string, ids []string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FindChildren(ctx context.Context, clubID, parentID string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Delete(ctx context.Context, clubID, id string) error           { return nil }
func (m *MockUserRepo) AnonymizeForGDPR(ctx context.Context, clubID, id string) error { return nil }
func (m *MockUserRepo) CreateIncident(ctx context.Context, incident *userDomain.IncidentLog) error {
	return nil
}

// --- Tests ---

func TestDisciplineUseCases_Disciplines(t *testing.T) {
	mockRepo := new(MockDisciplineRepo)
	uc := application.NewDisciplineUseCases(mockRepo, nil, nil)
	clubID := "club-1"

	t.Run("Create Discipline", func(t *testing.T) {
		mockRepo.On("CreateDiscipline", mock.Anything, mock.Anything).Return(nil).Once()
		d, err := uc.CreateDiscipline(context.TODO(), clubID, "Fútbol", "Escuela de fútbol")
		assert.NoError(t, err)
		assert.Equal(t, "Fútbol", d.Name)
	})

	t.Run("Create Group", func(t *testing.T) {
		dID := uuid.New()
		mockRepo.On("CreateGroup", mock.Anything, mock.Anything).Return(nil).Once()
		g, err := uc.CreateGroup(context.TODO(), clubID, "Sub-12 A", dID, "2012", "coach-1", "Lun-Mie 18:00")
		assert.NoError(t, err)
		assert.Equal(t, "Sub-12 A", g.Name)
	})

	t.Run("List Disciplines", func(t *testing.T) {
		mockRepo.On("ListDisciplines", mock.Anything, clubID).Return([]domain.Discipline{{Name: "D1"}}, nil).Once()
		list, err := uc.ListDisciplines(context.TODO(), clubID)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("List Groups with Filters", func(t *testing.T) {
		dID := uuid.New()
		mockRepo.On("ListGroups", mock.Anything, clubID, mock.MatchedBy(func(f map[string]interface{}) bool {
			return f["discipline_id"] == dID && f["category"] == "2010"
		})).Return([]domain.TrainingGroup{{Name: "G1"}}, nil).Once()

		list, err := uc.ListGroups(context.TODO(), clubID, dID.String(), "2010")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestDisciplineUseCases_Students(t *testing.T) {
	mockRepo := new(MockDisciplineRepo)
	mockUserRepo := new(MockUserRepo)
	uc := application.NewDisciplineUseCases(mockRepo, nil, mockUserRepo)
	clubID := "club-1"
	gID := uuid.New()

	t.Run("List Students In Group", func(t *testing.T) {
		mockRepo.On("GetGroupByID", mock.Anything, clubID, gID).Return(&domain.TrainingGroup{Category: "2015"}, nil).Once()
		mockUserRepo.On("List", mock.Anything, clubID, 100, 0, map[string]interface{}{"category": "2015"}).
			Return([]userDomain.User{{Email: "student@test.com"}}, nil).Once()

		students, err := uc.ListStudentsInGroup(context.TODO(), clubID, gID)
		assert.NoError(t, err)
		assert.Len(t, students, 1)
	})
}

func TestDisciplineUseCases_Tournaments(t *testing.T) {
	mockTourneyRepo := new(MockTournamentRepo)
	uc := application.NewDisciplineUseCases(nil, mockTourneyRepo, nil)
	clubID := "club-1"
	dID := uuid.New()

	t.Run("Create Tournament", func(t *testing.T) {
		mockTourneyRepo.On("CreateTournament", mock.Anything, mock.Anything).Return(nil).Once()
		tourney, err := uc.CreateTournament(context.TODO(), clubID, "Torneo Invierno", dID.String(), time.Now(), time.Now().AddDate(0, 1, 0), "LEAGUE")
		assert.NoError(t, err)
		assert.Equal(t, "Torneo Invierno", tourney.Name)
	})

	t.Run("Register Team", func(t *testing.T) {
		tID := uuid.New()
		mockTourneyRepo.On("CreateTeam", mock.Anything, mock.Anything).Return(nil).Once()
		team, err := uc.RegisterTeam(context.TODO(), clubID, tID.String(), "Los Cracks", "user-captain", []string{"u1", "u2"})
		assert.NoError(t, err)
		assert.Equal(t, "Los Cracks", team.Name)
	})

	t.Run("Schedule and Result Match", func(t *testing.T) {
		tID := uuid.New()
		hID := uuid.New()
		aID := uuid.New()
		mID := uuid.New()

		mockTourneyRepo.On("CreateMatch", mock.Anything, mock.MatchedBy(func(m *domain.Match) bool {
			return m.HomeTeamID == hID && m.AwayTeamID == aID
		})).Return(nil).Once()

		match, err := uc.ScheduleMatch(context.TODO(), clubID, tID.String(), hID.String(), aID.String(), time.Now(), "Cancha 1", "R1")
		assert.NoError(t, err)
		assert.NotNil(t, match)

		mockTourneyRepo.On("GetMatchByID", mock.Anything, clubID, mID).Return(&domain.Match{ID: mID}, nil).Once()
		mockTourneyRepo.On("UpdateMatch", mock.Anything, mock.Anything).Return(nil).Once()

		updated, err := uc.UpdateMatchResult(context.TODO(), clubID, mID.String(), 2, 1)
		assert.NoError(t, err)
		assert.Equal(t, 2, updated.ScoreHome)
		assert.Equal(t, domain.MatchStatusPlayed, updated.Status)
	})
}
