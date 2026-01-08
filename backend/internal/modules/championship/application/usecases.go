package application

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
)

// BookingService defines the dependency on the Booking module
type BookingService interface {
	CreateSystemBooking(clubID, courtID string, startTime, endTime time.Time, notes string) (*uuid.UUID, error)
}

type UserService interface {
	UpdateMatchStats(clubID, userID string, won bool, xpGained int) error
}

type ChampionshipUseCases struct {
	repo           domain.ChampionshipRepository
	bookingService BookingService
	userService    UserService
}

func NewChampionshipUseCases(repo domain.ChampionshipRepository, bookingService BookingService, userService UserService) *ChampionshipUseCases {
	return &ChampionshipUseCases{
		repo:           repo,
		bookingService: bookingService,
		userService:    userService,
	}
}

type CreateTournamentInput struct {
	ClubID    string    `json:"club_id"`
	Name      string    `json:"name"`
	Sport     string    `json:"sport"`
	Category  string    `json:"category"`
	StartDate time.Time `json:"start_date"`
}

func (uc *ChampionshipUseCases) CreateTournament(input CreateTournamentInput) (*domain.Tournament, error) {
	tournament := &domain.Tournament{
		ID:        uuid.New(),
		ClubID:    uuid.MustParse(input.ClubID),
		Name:      input.Name,
		Sport:     input.Sport,
		Category:  input.Category,
		Status:    domain.TournamentDraft, // Initial status
		StartDate: input.StartDate,
	}

	if err := uc.repo.CreateTournament(tournament); err != nil {
		return nil, err
	}

	return tournament, nil
}

func (uc *ChampionshipUseCases) ListTournaments(clubID string) ([]domain.Tournament, error) {
	return uc.repo.ListTournaments(clubID)
}

func (uc *ChampionshipUseCases) GetTournament(clubID, id string) (*domain.Tournament, error) {
	return uc.repo.GetTournament(clubID, id)
}

func (uc *ChampionshipUseCases) GetStandings(clubID, groupID string) ([]domain.Standing, error) {
	return uc.repo.GetStandings(clubID, groupID)
}

func (uc *ChampionshipUseCases) GetMatchesByGroup(clubID, groupID string) ([]domain.TournamentMatch, error) {
	return uc.repo.GetMatchesByGroup(clubID, groupID)
}

type AddStageInput struct {
	TournamentID string `json:"tournament_id"` // Path param usually, but can be in DTO
	ClubID       string `json:"club_id"`
	Name         string `json:"name"`
	Type         string `json:"type"` // "GROUP" or "KNOCKOUT"
	Order        int    `json:"order"`
}

func (uc *ChampionshipUseCases) AddStage(tournamentID string, input AddStageInput) (*domain.TournamentStage, error) {
	// Verify tournament belongs to club
	if _, err := uc.repo.GetTournament(input.ClubID, tournamentID); err != nil {
		return nil, errors.New("tournament not found or access denied")
	}

	stage := &domain.TournamentStage{
		ID:           uuid.New(),
		TournamentID: uuid.MustParse(tournamentID),
		Name:         input.Name,
		Type:         domain.StageType(input.Type),
		Order:        input.Order,
		Status:       domain.StagePending,
	}

	if err := uc.repo.CreateStage(stage); err != nil {
		return nil, err
	}
	return stage, nil
}

type AddGroupInput struct {
	StageID string `json:"stage_id"`
	ClubID  string `json:"club_id"`
	Name    string `json:"name"`
}

func (uc *ChampionshipUseCases) AddGroup(stageID string, input AddGroupInput) (*domain.Group, error) {
	// Verify stage belongs to club
	if _, err := uc.repo.GetStage(input.ClubID, stageID); err != nil {
		return nil, errors.New("stage not found or access denied")
	}

	group := &domain.Group{
		ID:      uuid.New(),
		StageID: uuid.MustParse(stageID),
		Name:    input.Name,
	}
	if err := uc.repo.CreateGroup(group); err != nil {
		return nil, err
	}
	return group, nil
}

type RegisterTeamInput struct {
	TeamID string `json:"team_id"`
}

func (uc *ChampionshipUseCases) RegisterTeam(clubID, groupID string, input RegisterTeamInput) (*domain.Standing, error) {
	// Verify group belongs to club
	if _, err := uc.repo.GetGroup(clubID, groupID); err != nil {
		return nil, errors.New("group not found or access denied")
	}

	// 1. Create Standing entry (which implicitly registers the team in the group)
	standing := &domain.Standing{
		ID:      uuid.New(),
		GroupID: uuid.MustParse(groupID),
		TeamID:  uuid.MustParse(input.TeamID),
	}

	// 2. Add to repo
	if err := uc.repo.RegisterTeam(standing); err != nil {
		return nil, err
	}
	return standing, nil
}

func (uc *ChampionshipUseCases) GenerateGroupFixture(clubID, groupID string) ([]domain.TournamentMatch, error) {
	// 1. Get Teams in Group (via Standings)
	standings, err := uc.repo.GetStandings(clubID, groupID)
	if err != nil {
		return nil, err
	}

	var teamIDs []uuid.UUID
	for _, s := range standings {
		teamIDs = append(teamIDs, s.TeamID)
	}

	if len(teamIDs) < 2 {
		return nil, errors.New("need at least 2 teams to generate fixture")
	}

	matches := []domain.TournamentMatch{}

	group, err := uc.repo.GetGroup(clubID, groupID)
	if err != nil {
		return nil, err
	}

	stage, err := uc.repo.GetStage(clubID, group.StageID.String())
	if err != nil {
		return nil, err
	}

	n := len(teamIDs)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			gID := group.ID
			match := domain.TournamentMatch{
				ID:           uuid.New(),
				TournamentID: stage.TournamentID,
				StageID:      stage.ID,
				GroupID:      &gID,
				HomeTeamID:   teamIDs[i],
				AwayTeamID:   teamIDs[j],
				Status:       domain.MatchScheduled,
				Date:         time.Now().Add(time.Hour * 24), // TBD
			}
			matches = append(matches, match)
		}
	}

	// Create all matches atomically in a single transaction
	if err := uc.repo.CreateMatchesBatch(matches); err != nil {
		return nil, err
	}

	return matches, nil
}

type UpdateMatchResultInput struct {
	ClubID    string `json:"club_id"`
	MatchID   string `json:"match_id"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

func (uc *ChampionshipUseCases) UpdateMatchResult(input UpdateMatchResultInput) error {
	if err := uc.repo.UpdateMatchResult(input.ClubID, input.MatchID, input.HomeScore, input.AwayScore); err != nil {
		return err
	}

	// Trigger Recalculate Standings
	// 1. Get Match to find GroupID
	match, err := uc.repo.GetMatch(input.ClubID, input.MatchID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.New("match not found after update")
	}

	// Update User Stats
	if uc.userService != nil {
		tournament, err := uc.repo.GetTournament(input.ClubID, match.TournamentID.String())
		if err == nil {
			clubID := tournament.ClubID.String()
			homeWon := input.HomeScore > input.AwayScore
			awayWon := input.AwayScore > input.HomeScore

			// XP Logic
			xp := 100 // base XP

			// Home Players
			homePlayers, _ := uc.repo.GetTeamMembers(match.HomeTeamID.String())
			for _, userID := range homePlayers {
				_ = uc.userService.UpdateMatchStats(clubID, userID, homeWon, xp)
			}

			// Away Players
			awayPlayers, _ := uc.repo.GetTeamMembers(match.AwayTeamID.String())
			for _, userID := range awayPlayers {
				_ = uc.userService.UpdateMatchStats(clubID, userID, awayWon, xp)
			}
		}
	}

	if match.GroupID == nil {
		return nil // Not a group match
	}

	return uc.recalculateStandings(input.ClubID, match.GroupID.String())
}

func (uc *ChampionshipUseCases) recalculateStandings(clubID, groupID string) error {
	standings, err := uc.repo.GetStandings(clubID, groupID)
	if err != nil {
		return err
	}

	matches, err := uc.repo.GetMatchesByGroup(clubID, groupID)
	if err != nil {
		return err
	}

	stats := make(map[uuid.UUID]*domain.Standing)
	for i := range standings {
		s := &standings[i]
		s.Played = 0
		s.Won = 0
		s.Drawn = 0
		s.Lost = 0
		s.GoalsFor = 0
		s.GoalsAgainst = 0
		s.GoalDifference = 0
		s.Points = 0
		stats[s.TeamID] = s
	}

	for _, m := range matches {
		if m.Status != domain.MatchCompleted {
			continue
		}
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}

		home, okH := stats[m.HomeTeamID]
		away, okA := stats[m.AwayTeamID]

		if okH && okA {
			homeScore := *m.HomeScore
			awayScore := *m.AwayScore

			home.Played++
			away.Played++
			home.GoalsFor += homeScore
			home.GoalsAgainst += awayScore
			home.GoalDifference = home.GoalsFor - home.GoalsAgainst

			away.GoalsFor += awayScore
			away.GoalsAgainst += homeScore
			away.GoalDifference = away.GoalsFor - away.GoalsAgainst

			if homeScore > awayScore {
				home.Won++
				home.Points += 3
				away.Lost++
			} else if awayScore > homeScore {
				away.Won++
				away.Points += 3
				home.Lost++
			} else {
				home.Drawn++
				home.Points += 1
				away.Drawn++
				away.Points += 1
			}
		}
	}

	for _, s := range stats {
		if err := uc.repo.UpdateStanding(s); err != nil {
			return err
		}
	}

	return nil
}

type ScheduleMatchInput struct {
	ClubID    string    `json:"club_id"`
	MatchID   string    `json:"match_id"`
	CourtID   string    `json:"court_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (uc *ChampionshipUseCases) ScheduleMatch(input ScheduleMatchInput) error {
	// 1. Create Booking via Service
	notes := "Partido de Torneo: " + input.MatchID
	bookingID, err := uc.bookingService.CreateSystemBooking(input.ClubID, input.CourtID, input.StartTime, input.EndTime, notes)
	if err != nil {
		return errors.New("failed to book court: " + err.Error())
	}

	// 2. Update Match with BookingID and Date
	return uc.repo.UpdateMatchScheduling(input.ClubID, input.MatchID, input.StartTime, *bookingID)
}
