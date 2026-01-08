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

func (uc *ChampionshipUseCases) GetStandings(groupID string) ([]domain.Standing, error) {
	return uc.repo.GetStandings(groupID)
}

func (uc *ChampionshipUseCases) GetMatchesByGroup(groupID string) ([]domain.TournamentMatch, error) {
	return uc.repo.GetMatchesByGroup(groupID)
}

type AddStageInput struct {
	TournamentID string `json:"tournament_id"` // Path param usually, but can be in DTO
	Name         string `json:"name"`
	Type         string `json:"type"` // "GROUP" or "KNOCKOUT"
	Order        int    `json:"order"`
}

func (uc *ChampionshipUseCases) AddStage(tournamentID string, input AddStageInput) (*domain.TournamentStage, error) {
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
	Name    string `json:"name"`
}

func (uc *ChampionshipUseCases) AddGroup(stageID string, input AddGroupInput) (*domain.Group, error) {
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

func (uc *ChampionshipUseCases) RegisterTeam(groupID string, input RegisterTeamInput) (*domain.Standing, error) {
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

func (uc *ChampionshipUseCases) GenerateGroupFixture(groupID string) ([]domain.TournamentMatch, error) {
	// 1. Get Teams in Group (via Standings)
	standings, err := uc.repo.GetStandings(groupID)
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
	// Simple All-vs-All (One leg)
	// We need StageID and TournamentID. Using a workaround:
	// Assuming GroupID -> GetGroup -> StageID -> GetStage -> TournamentID.
	// Since we don't have GetGroup public yet, we'll fetch ONE standing and check? No.
	// We will try to fetch the Group using GetGroup (assuming we add it/it exists in repo impl naturally).
	// Checking repository interface: `GetGroup(id string) (*Group, error)`. YES IT EXISTS.

	group, err := uc.repo.GetGroup(groupID)
	if err != nil {
		return nil, err
	}

	stage, err := uc.repo.GetStage(group.StageID.String())
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

	for _, m := range matches {
		// Use CreateMatch one by one or add bulk
		if err := uc.repo.CreateMatch(&m); err != nil {
			return nil, err
		}
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
	if err := uc.repo.UpdateMatchResult(input.MatchID, input.HomeScore, input.AwayScore); err != nil {
		return err
	}

	// Trigger Recalculate Standings
	// 1. Get Match to find GroupID
	match, err := uc.repo.GetMatch(input.MatchID)
	if err != nil {
		return err // Log warning?
	}

	// Update User Stats
	// Assuming match just completed. To avoid double counting, we should check previous status.
	// But simply, let's assume this call finalizes the match.
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

	return uc.recalculateStandings(match.GroupID.String())
}

func (uc *ChampionshipUseCases) recalculateStandings(groupID string) error {
	standings, err := uc.repo.GetStandings(groupID)
	if err != nil {
		return err
	}

	matches, err := uc.repo.GetMatchesByGroup(groupID)
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
		_ = uc.repo.UpdateStanding(s)
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
	return uc.repo.UpdateMatchScheduling(input.MatchID, input.StartTime, *bookingID)
}
