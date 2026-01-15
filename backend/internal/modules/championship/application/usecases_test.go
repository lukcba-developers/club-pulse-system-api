package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockChampionshipRepo struct {
	mock.Mock
}

func (m *MockChampionshipRepo) CreateTournament(ctx context.Context, t *domain.Tournament) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetTournament(ctx context.Context, clubID, id string) (*domain.Tournament, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tournament), args.Error(1)
}

func (m *MockChampionshipRepo) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	args := m.Called(ctx, clubID)
	var res []domain.Tournament
	if args.Get(0) != nil {
		res = args.Get(0).([]domain.Tournament)
	}
	return res, args.Error(1)
}

func (m *MockChampionshipRepo) CreateStage(ctx context.Context, s *domain.TournamentStage) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetStage(ctx context.Context, clubID, id string) (*domain.TournamentStage, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TournamentStage), args.Error(1)
}

func (m *MockChampionshipRepo) CreateGroup(ctx context.Context, g *domain.Group) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetGroup(ctx context.Context, clubID, id string) (*domain.Group, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Group), args.Error(1)
}

func (m *MockChampionshipRepo) CreateMatch(ctx context.Context, match *domain.TournamentMatch) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockChampionshipRepo) CreateMatchesBatch(ctx context.Context, matches []domain.TournamentMatch) error {
	args := m.Called(ctx, matches)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetMatch(ctx context.Context, clubID, id string) (*domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TournamentMatch), args.Error(1)
}

func (m *MockChampionshipRepo) GetMatchesByGroup(ctx context.Context, clubID, groupID string) ([]domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, groupID)
	var res []domain.TournamentMatch
	if args.Get(0) != nil {
		res = args.Get(0).([]domain.TournamentMatch)
	}
	return res, args.Error(1)
}

func (m *MockChampionshipRepo) UpdateMatchResult(ctx context.Context, clubID, matchID string, homeScore, awayScore float64) error {
	args := m.Called(ctx, clubID, matchID, homeScore, awayScore)
	return args.Error(0)
}

func (m *MockChampionshipRepo) UpdateMatchScheduling(ctx context.Context, clubID, matchID string, date time.Time, bookingID uuid.UUID) error {
	args := m.Called(ctx, clubID, matchID, date, bookingID)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetStandings(ctx context.Context, clubID, groupID string) ([]domain.Standing, error) {
	args := m.Called(ctx, clubID, groupID)
	var res []domain.Standing
	if args.Get(0) != nil {
		res = args.Get(0).([]domain.Standing)
	}
	return res, args.Error(1)
}

func (m *MockChampionshipRepo) RegisterTeam(ctx context.Context, s *domain.Standing) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) UpdateStanding(ctx context.Context, s *domain.Standing) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) UpdateStandingsBatch(ctx context.Context, s []domain.Standing) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetTeamMembers(ctx context.Context, id string) ([]string, error) {
	args := m.Called(ctx, id)
	var res []string
	if args.Get(0) != nil {
		res = args.Get(0).([]string)
	}
	return res, args.Error(1)
}

func (m *MockChampionshipRepo) CreateTeam(ctx context.Context, team *domain.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockChampionshipRepo) AddMember(ctx context.Context, teamID, userID string) error {
	args := m.Called(ctx, teamID, userID)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetMatchesByUserID(ctx context.Context, clubID, userID string) ([]domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]domain.TournamentMatch), args.Error(1)
}

func (m *MockChampionshipRepo) GetUpcomingMatches(ctx context.Context, clubID string, from, to time.Time) ([]domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, from, to)
	return args.Get(0).([]domain.TournamentMatch), args.Error(1)
}

type MockBookingService struct {
	mock.Mock
}

func (m *MockBookingService) CreateSystemBooking(clubID, courtID string, start, end time.Time, notes string) (*uuid.UUID, error) {
	args := m.Called(clubID, courtID, start, end, notes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	id := args.Get(0).(uuid.UUID)
	return &id, args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) UpdateMatchStats(ctx context.Context, clubID, userID string, won bool, xp int) error {
	args := m.Called(ctx, clubID, userID, won, xp)
	return args.Error(0)
}

// --- Tests ---

func TestChampionshipUseCases_Tournament(t *testing.T) {
	repo := new(MockChampionshipRepo)
	uc := application.NewChampionshipUseCases(repo, nil, nil)
	clubID := uuid.New().String()

	t.Run("CreateTournament Success", func(t *testing.T) {
		repo.On("CreateTournament", mock.Anything, mock.Anything).Return(nil).Once()
		input := application.CreateTournamentInput{
			ClubID: clubID,
			Name:   "Copa America",
			Sport:  "FUTBOL",
		}
		res, err := uc.CreateTournament(context.TODO(), input)
		assert.NoError(t, err)
		assert.Equal(t, "Copa America", res.Name)
	})

	t.Run("CreateTournament Error", func(t *testing.T) {
		repo.On("CreateTournament", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
		_, err := uc.CreateTournament(context.TODO(), application.CreateTournamentInput{ClubID: clubID})
		assert.Error(t, err)
	})

	t.Run("List and Get", func(t *testing.T) {
		repo.On("ListTournaments", mock.Anything, clubID).Return([]domain.Tournament{{Name: "T1"}}, nil).Once()
		repo.On("GetTournament", mock.Anything, clubID, "id1").Return(&domain.Tournament{Name: "T1"}, nil).Once()

		list, _ := uc.ListTournaments(context.TODO(), clubID)
		get, _ := uc.GetTournament(context.TODO(), clubID, "id1")

		assert.Len(t, list, 1)
		assert.Equal(t, "T1", get.Name)
	})
}

func TestChampionshipUseCases_StagesAndGroups(t *testing.T) {
	repo := new(MockChampionshipRepo)
	uc := application.NewChampionshipUseCases(repo, nil, nil)
	cID := "club-1"
	tID := uuid.New().String()

	t.Run("AddStage Success", func(t *testing.T) {
		repo.On("GetTournament", mock.Anything, cID, tID).Return(&domain.Tournament{ID: uuid.MustParse(tID)}, nil).Once()
		repo.On("CreateStage", mock.Anything, mock.Anything).Return(nil).Once()

		res, err := uc.AddStage(context.TODO(), tID, application.AddStageInput{ClubID: cID, Name: "Phase 1", Type: "GROUP"})
		assert.NoError(t, err)
		assert.Equal(t, "Phase 1", res.Name)
	})

	t.Run("AddStage TournamentError", func(t *testing.T) {
		repo.On("GetTournament", mock.Anything, cID, tID).Return(nil, errors.New("not found")).Once()
		_, err := uc.AddStage(context.TODO(), tID, application.AddStageInput{ClubID: cID})
		assert.Error(t, err)
	})

	t.Run("AddStage CreateError", func(t *testing.T) {
		repo.On("GetTournament", mock.Anything, cID, tID).Return(&domain.Tournament{}, nil).Once()
		repo.On("CreateStage", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
		_, err := uc.AddStage(context.TODO(), tID, application.AddStageInput{ClubID: cID})
		assert.Error(t, err)
	})

	t.Run("AddGroup Success", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(&domain.TournamentStage{ID: uuid.MustParse(sID)}, nil).Once()
		repo.On("CreateGroup", mock.Anything, mock.Anything).Return(nil).Once()

		res, err := uc.AddGroup(context.TODO(), sID, application.AddGroupInput{ClubID: cID, Name: "Group A"})
		assert.NoError(t, err)
		assert.Equal(t, "Group A", res.Name)
	})

	t.Run("AddGroup StageError", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(nil, errors.New("not found")).Once()
		_, err := uc.AddGroup(context.TODO(), sID, application.AddGroupInput{ClubID: cID})
		assert.Error(t, err)
	})

	t.Run("AddGroup CreateError", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(&domain.TournamentStage{}, nil).Once()
		repo.On("CreateGroup", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
		_, err := uc.AddGroup(context.TODO(), sID, application.AddGroupInput{ClubID: cID})
		assert.Error(t, err)
	})
}

func TestChampionshipUseCases_Teams(t *testing.T) {
	repo := new(MockChampionshipRepo)
	uc := application.NewChampionshipUseCases(repo, nil, nil)
	cID := "club-1"
	gID := uuid.New().String()

	t.Run("RegisterTeam Success", func(t *testing.T) {
		repo.On("GetGroup", mock.Anything, cID, gID).Return(&domain.Group{ID: uuid.MustParse(gID)}, nil).Once()
		repo.On("RegisterTeam", mock.Anything, mock.Anything).Return(nil).Once()

		res, err := uc.RegisterTeam(context.TODO(), cID, gID, application.RegisterTeamInput{TeamID: uuid.New().String()})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("RegisterTeam GroupError", func(t *testing.T) {
		repo.On("GetGroup", mock.Anything, cID, gID).Return(nil, errors.New("not found")).Once()
		_, err := uc.RegisterTeam(context.TODO(), cID, gID, application.RegisterTeamInput{TeamID: uuid.New().String()})
		assert.Error(t, err)
	})

	t.Run("RegisterTeam CreateError", func(t *testing.T) {
		repo.On("GetGroup", mock.Anything, cID, gID).Return(&domain.Group{}, nil).Once()
		repo.On("RegisterTeam", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
		_, err := uc.RegisterTeam(context.TODO(), cID, gID, application.RegisterTeamInput{TeamID: uuid.New().String()})
		assert.Error(t, err)
	})
}

func TestChampionshipUseCases_FixtureLogic(t *testing.T) {
	repo := new(MockChampionshipRepo)
	uc := application.NewChampionshipUseCases(repo, nil, nil)
	cID := "club-1"
	gID := uuid.New().String()

	t.Run("GenerateGroupFixture (Round-robin)", func(t *testing.T) {
		teams := []domain.Standing{
			{TeamID: uuid.New()}, {TeamID: uuid.New()}, {TeamID: uuid.New()},
		}
		repo.On("GetStandings", mock.Anything, cID, gID).Return(teams, nil).Once()
		repo.On("GetGroup", mock.Anything, cID, gID).Return(&domain.Group{ID: uuid.MustParse(gID), StageID: uuid.New()}, nil).Once()
		repo.On("GetStage", mock.Anything, cID, mock.Anything).Return(&domain.TournamentStage{ID: uuid.New()}, nil).Once()
		repo.On("CreateMatchesBatch", mock.Anything, mock.Anything).Return(nil).Once()

		matches, err := uc.GenerateGroupFixture(context.TODO(), cID, gID)
		assert.NoError(t, err)
		assert.Len(t, matches, 3)
	})

	t.Run("GenerateGroupFixture MinTeamsError", func(t *testing.T) {
		repo.On("GetStandings", mock.Anything, cID, gID).Return([]domain.Standing{{TeamID: uuid.New()}}, nil).Once()
		_, err := uc.GenerateGroupFixture(context.TODO(), cID, gID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 2 teams")
	})

	t.Run("GenerateGroupFixture RepoError", func(t *testing.T) {
		repo.On("GetStandings", mock.Anything, cID, gID).Return(nil, errors.New("db error")).Once()
		_, err := uc.GenerateGroupFixture(context.TODO(), cID, gID)
		assert.Error(t, err)
	})

	t.Run("GenerateKnockoutBracket Success", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(&domain.TournamentStage{ID: uuid.MustParse(sID), Type: domain.StageKnockout}, nil).Once()
		repo.On("CreateMatchesBatch", mock.Anything, mock.Anything).Return(nil).Once()

		seeds := []string{uuid.New().String(), uuid.New().String(), uuid.New().String(), uuid.New().String()}
		matches, err := uc.GenerateKnockoutBracket(context.TODO(), application.GenerateKnockoutBracketInput{
			ClubID: cID, StageID: sID, SeedOrder: seeds,
		})
		assert.NoError(t, err)
		assert.Len(t, matches, 2)
	})

	t.Run("GenerateKnockoutBracket StageNotFound", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(nil, nil).Once()
		_, err := uc.GenerateKnockoutBracket(context.TODO(), application.GenerateKnockoutBracketInput{
			ClubID: cID, StageID: sID, SeedOrder: []string{uuid.New().String(), uuid.New().String()},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stage not found")
	})

	t.Run("GenerateKnockoutBracket WrongStageType", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(&domain.TournamentStage{Type: domain.StageGroup}, nil).Once()
		_, err := uc.GenerateKnockoutBracket(context.TODO(), application.GenerateKnockoutBracketInput{
			ClubID: cID, StageID: sID, SeedOrder: []string{uuid.New().String(), uuid.New().String()},
		})
		assert.Error(t, err)
	})

	t.Run("GenerateKnockoutBracket Invalid UUID", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(&domain.TournamentStage{ID: uuid.MustParse(sID), Type: domain.StageKnockout}, nil).Once()

		seeds := []string{"invalid-uuid", uuid.New().String()}
		_, err := uc.GenerateKnockoutBracket(context.TODO(), application.GenerateKnockoutBracketInput{
			ClubID: cID, StageID: sID, SeedOrder: seeds,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid team ID")
	})

	t.Run("GenerateKnockoutBracket MinTeamsError", func(t *testing.T) {
		sID := uuid.New().String()
		repo.On("GetStage", mock.Anything, cID, sID).Return(&domain.TournamentStage{ID: uuid.MustParse(sID), Type: domain.StageKnockout}, nil).Once()

		_, err := uc.GenerateKnockoutBracket(context.TODO(), application.GenerateKnockoutBracketInput{
			ClubID: cID, StageID: sID, SeedOrder: []string{uuid.New().String()},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 2 teams")
	})
}

func TestChampionshipUseCases_Results(t *testing.T) {
	repo := new(MockChampionshipRepo)
	userSvc := new(MockUserService)
	uc := application.NewChampionshipUseCases(repo, nil, userSvc)
	cID := "club-1"
	mID := uuid.New().String()

	t.Run("UpdateMatchResult and Recalculate", func(t *testing.T) {
		match := &domain.TournamentMatch{
			ID:           uuid.MustParse(mID),
			TournamentID: uuid.New(),
			GroupID:      func() *uuid.UUID { id := uuid.New(); return &id }(),
			HomeTeamID:   uuid.New(),
			AwayTeamID:   uuid.New(),
		}
		repo.On("UpdateMatchResult", mock.Anything, cID, mID, 2.0, 1.0).Return(nil).Once()
		repo.On("GetMatch", mock.Anything, cID, mID).Return(match, nil).Once()
		repo.On("GetGroup", mock.Anything, cID, mock.Anything).Return(&domain.Group{StageID: uuid.New()}, nil).Once()
		repo.On("GetStage", mock.Anything, cID, mock.Anything).Return(&domain.TournamentStage{TournamentID: uuid.New()}, nil).Once()
		repo.On("GetTournament", mock.Anything, cID, mock.Anything).Return(&domain.Tournament{ClubID: uuid.New()}, nil).Twice()
		repo.On("GetTeamMembers", mock.Anything, mock.Anything).Return([]string{"u1"}, nil).Twice()
		userSvc.On("UpdateMatchStats", mock.Anything, mock.Anything, "u1", mock.Anything, 100).Return(nil).Twice()

		// Recalculate logic
		repo.On("GetStandings", mock.Anything, cID, mock.Anything).Return([]domain.Standing{
			{TeamID: match.HomeTeamID}, {TeamID: match.AwayTeamID},
		}, nil).Once()
		hScore := 2.0
		aScore := 1.0
		repo.On("GetMatchesByGroup", mock.Anything, cID, mock.Anything).Return([]domain.TournamentMatch{
			{
				HomeTeamID: match.HomeTeamID, AwayTeamID: match.AwayTeamID,
				HomeScore: &hScore, AwayScore: &aScore, Status: domain.MatchCompleted,
			},
		}, nil).Once()
		repo.On("UpdateStandingsBatch", mock.Anything, mock.Anything).Return(nil).Once()

		err := uc.UpdateMatchResult(context.TODO(), application.UpdateMatchResultInput{
			ClubID: cID, MatchID: mID, HomeScore: 2.0, AwayScore: 1.0,
		})
		assert.NoError(t, err)
	})

	t.Run("UpdateMatchResult MatchNotFoundError", func(t *testing.T) {
		repo.On("UpdateMatchResult", mock.Anything, cID, mID, 2.0, 1.0).Return(nil).Once()
		repo.On("GetMatch", mock.Anything, cID, mID).Return(nil, nil).Once()
		err := uc.UpdateMatchResult(context.TODO(), application.UpdateMatchResultInput{
			ClubID: cID, MatchID: mID, HomeScore: 2.0, AwayScore: 1.0,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "match not found")
	})

	t.Run("UpdateMatchResult RepoError", func(t *testing.T) {
		repo.On("UpdateMatchResult", mock.Anything, cID, mID, 2.0, 1.0).Return(errors.New("db error")).Once()
		err := uc.UpdateMatchResult(context.TODO(), application.UpdateMatchResultInput{
			ClubID: cID, MatchID: mID, HomeScore: 2.0, AwayScore: 1.0,
		})
		assert.Error(t, err)
	})
}

func TestChampionshipUseCases_Scheduling(t *testing.T) {
	repo := new(MockChampionshipRepo)
	bookingSvc := new(MockBookingService)
	uc := application.NewChampionshipUseCases(repo, bookingSvc, nil)
	cID := "club-1"

	t.Run("ScheduleMatch Success", func(t *testing.T) {
		bID := uuid.New()
		bookingSvc.On("CreateSystemBooking", cID, "court-1", mock.Anything, mock.Anything, mock.Anything).Return(bID, nil).Once()
		repo.On("UpdateMatchScheduling", mock.Anything, cID, "match-1", mock.Anything, bID).Return(nil).Once()

		err := uc.ScheduleMatch(context.TODO(), application.ScheduleMatchInput{
			ClubID: cID, MatchID: "match-1", CourtID: "court-1",
		})
		assert.NoError(t, err)
	})

	t.Run("ScheduleMatch BookingError", func(t *testing.T) {
		bookingSvc.On("CreateSystemBooking", cID, "court-1", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("court busy")).Once()
		err := uc.ScheduleMatch(context.TODO(), application.ScheduleMatchInput{
			ClubID: cID, MatchID: "match-1", CourtID: "court-1",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to book court")
	})
}
