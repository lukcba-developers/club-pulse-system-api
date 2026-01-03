package application

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type DisciplineUseCases struct {
	repo           domain.DisciplineRepository
	tournamentRepo domain.TournamentRepository
	userRepo       userDomain.UserRepository
}

func NewDisciplineUseCases(repo domain.DisciplineRepository, tournamentRepo domain.TournamentRepository, userRepo userDomain.UserRepository) *DisciplineUseCases {
	return &DisciplineUseCases{
		repo:           repo,
		tournamentRepo: tournamentRepo,
		userRepo:       userRepo,
	}
}

func (uc *DisciplineUseCases) CreateDiscipline(clubID string, name, description string) (*domain.Discipline, error) {
	d := &domain.Discipline{
		ID:          uuid.New(),
		ClubID:      clubID,
		Name:        name,
		Description: description,
		IsActive:    true,
	}
	if err := uc.repo.CreateDiscipline(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (uc *DisciplineUseCases) CreateGroup(clubID string, name string, dID uuid.UUID, category, coachID, schedule string) (*domain.TrainingGroup, error) {
	g := &domain.TrainingGroup{
		ID:           uuid.New(),
		ClubID:       clubID,
		Name:         name,
		DisciplineID: dID,
		Category:     category,
		CoachID:      coachID,
		Schedule:     schedule,
	}
	if err := uc.repo.CreateGroup(g); err != nil {
		return nil, err
	}
	return g, nil
}

func (uc *DisciplineUseCases) ListDisciplines(clubID string) ([]domain.Discipline, error) {
	return uc.repo.ListDisciplines(clubID)
}

func (uc *DisciplineUseCases) ListGroups(clubID string, disciplineID string, category string) ([]domain.TrainingGroup, error) {
	filter := make(map[string]interface{})
	if disciplineID != "" {
		id, err := uuid.Parse(disciplineID)
		if err == nil {
			filter["discipline_id"] = id
		}
	}
	if category != "" {
		filter["category"] = category
	}
	return uc.repo.ListGroups(clubID, filter)
}

func (uc *DisciplineUseCases) ListStudentsInGroup(clubID string, groupID uuid.UUID) ([]userDomain.User, error) {
	group, err := uc.repo.GetGroupByID(clubID, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	// Filter users by the group's category (year of birth)
	// We use the User module's List method with the "category" filter.
	return uc.userRepo.List(clubID, 100, 0, map[string]interface{}{"category": group.Category})
}

// --- Championships ---

func (uc *DisciplineUseCases) CreateTournament(clubID string, name, disciplineID string, startDate, endDate time.Time, format string) (*domain.Tournament, error) {
	dID, err := uuid.Parse(disciplineID)
	if err != nil {
		return nil, err
	}
	t := &domain.Tournament{
		ID:           uuid.New(),
		ClubID:       clubID,
		Name:         name,
		DisciplineID: dID,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       domain.TournamentStatusOpen,
		Format:       format,
	}
	if err := uc.tournamentRepo.CreateTournament(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (uc *DisciplineUseCases) ListTournaments(clubID string) ([]domain.Tournament, error) {
	return uc.tournamentRepo.ListTournaments(clubID)
}

func (uc *DisciplineUseCases) RegisterTeam(clubID, tournamentID string, name string, captainID string, memberIDs []string) (*domain.Team, error) {
	tID, err := uuid.Parse(tournamentID)
	if err != nil {
		return nil, err
	}

	team := &domain.Team{
		ID:           uuid.New(),
		ClubID:       clubID,
		TournamentID: tID,
		Name:         name,
		Members:      memberIDs,
	}
	if captainID != "" {
		team.CaptainID = &captainID
	}

	if err := uc.tournamentRepo.CreateTeam(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (uc *DisciplineUseCases) ScheduleMatch(clubID, tournamentID, homeTeamID, awayTeamID string, startTime time.Time, location, round string) (*domain.Match, error) {
	tID, err := uuid.Parse(tournamentID)
	if err != nil {
		return nil, err
	}
	hID, err := uuid.Parse(homeTeamID)
	if err != nil {
		return nil, err
	}
	aID, err := uuid.Parse(awayTeamID)
	if err != nil {
		return nil, err
	}

	match := &domain.Match{
		ID:           uuid.New(),
		ClubID:       clubID,
		TournamentID: tID,
		HomeTeamID:   hID,
		AwayTeamID:   aID,
		StartTime:    startTime,
		Location:     location,
		Round:        round,
		Status:       domain.MatchStatusScheduled,
	}

	if err := uc.tournamentRepo.CreateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

func (uc *DisciplineUseCases) UpdateMatchResult(clubID, matchID string, scoreHome, scoreAway int) (*domain.Match, error) {
	mID, err := uuid.Parse(matchID)
	if err != nil {
		return nil, err
	}

	match, err := uc.tournamentRepo.GetMatchByID(clubID, mID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, errors.New("match not found")
	}

	match.ScoreHome = scoreHome
	match.ScoreAway = scoreAway
	match.Status = domain.MatchStatusPlayed
	match.UpdatedAt = time.Now()

	if err := uc.tournamentRepo.UpdateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

func (uc *DisciplineUseCases) GetStandings(clubID, tournamentID string) ([]domain.Standing, error) {
	tID, err := uuid.Parse(tournamentID)
	if err != nil {
		return nil, err
	}
	return uc.tournamentRepo.GetStandings(clubID, tID)
}

func (uc *DisciplineUseCases) ListMatches(clubID, tournamentID string) ([]domain.Match, error) {
	tID, err := uuid.Parse(tournamentID)
	if err != nil {
		return nil, err
	}
	return uc.tournamentRepo.ListMatches(clubID, tID)
}
