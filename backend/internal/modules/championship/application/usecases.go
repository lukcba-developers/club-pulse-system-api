package application

import (
	"context"
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
	UpdateMatchStats(ctx context.Context, clubID, userID string, won bool, xpGained int) error
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
	Name      string    `json:"name" binding:"required"`
	Sport     string    `json:"sport" binding:"required"`
	Category  string    `json:"category"`
	StartDate time.Time `json:"start_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
}

func (uc *ChampionshipUseCases) CreateTournament(ctx context.Context, input CreateTournamentInput) (*domain.Tournament, error) {
	tournament := &domain.Tournament{
		ID:        uuid.New(),
		ClubID:    uuid.MustParse(input.ClubID),
		Name:      input.Name,
		Sport:     input.Sport,
		Category:  input.Category,
		Status:    domain.TournamentDraft, // Initial status
		StartDate: input.StartDate,
	}

	if err := uc.repo.CreateTournament(ctx, tournament); err != nil {
		return nil, err
	}

	return tournament, nil
}

func (uc *ChampionshipUseCases) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	return uc.repo.ListTournaments(ctx, clubID)
}

func (uc *ChampionshipUseCases) GetTournament(ctx context.Context, clubID, id string) (*domain.Tournament, error) {
	return uc.repo.GetTournament(ctx, clubID, id)
}

func (uc *ChampionshipUseCases) GetStandings(ctx context.Context, clubID, groupID string) ([]domain.Standing, error) {
	return uc.repo.GetStandings(ctx, clubID, groupID)
}

func (uc *ChampionshipUseCases) GetMatchesByGroup(ctx context.Context, clubID, groupID string) ([]domain.TournamentMatch, error) {
	return uc.repo.GetMatchesByGroup(ctx, clubID, groupID)
}

type AddStageInput struct {
	TournamentID string `json:"tournament_id"` // Path param usually, but can be in DTO
	ClubID       string `json:"club_id"`
	Name         string `json:"name"`
	Type         string `json:"type"` // "GROUP" or "KNOCKOUT"
	Order        int    `json:"order"`
}

func (uc *ChampionshipUseCases) AddStage(ctx context.Context, tournamentID string, input AddStageInput) (*domain.TournamentStage, error) {
	// Verify tournament belongs to club
	if _, err := uc.repo.GetTournament(ctx, input.ClubID, tournamentID); err != nil {
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

	if err := uc.repo.CreateStage(ctx, stage); err != nil {
		return nil, err
	}
	return stage, nil
}

type AddGroupInput struct {
	StageID string `json:"stage_id"`
	ClubID  string `json:"club_id"`
	Name    string `json:"name"`
}

func (uc *ChampionshipUseCases) AddGroup(ctx context.Context, stageID string, input AddGroupInput) (*domain.Group, error) {
	// Verify stage belongs to club
	if _, err := uc.repo.GetStage(ctx, input.ClubID, stageID); err != nil {
		return nil, errors.New("stage not found or access denied")
	}

	group := &domain.Group{
		ID:      uuid.New(),
		StageID: uuid.MustParse(stageID),
		Name:    input.Name,
	}
	if err := uc.repo.CreateGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

type RegisterTeamInput struct {
	TeamID string `json:"team_id"`
}

func (uc *ChampionshipUseCases) RegisterTeam(ctx context.Context, clubID, groupID string, input RegisterTeamInput) (*domain.Standing, error) {
	// Verify group belongs to club
	if _, err := uc.repo.GetGroup(ctx, clubID, groupID); err != nil {
		return nil, errors.New("group not found or access denied")
	}

	// 1. Create Standing entry (which implicitly registers the team in the group)
	standing := &domain.Standing{
		ID:      uuid.New(),
		GroupID: uuid.MustParse(groupID),
		TeamID:  uuid.MustParse(input.TeamID),
	}

	// 2. Add to repo
	if err := uc.repo.RegisterTeam(ctx, clubID, standing); err != nil {
		return nil, err
	}
	return standing, nil
}

// ... (CreateTeam, AddMember methods are fine or out of scope for now)

func (uc *ChampionshipUseCases) GenerateGroupFixture(ctx context.Context, clubID, groupID string) ([]domain.TournamentMatch, error) {
	// 1. Get Teams in Group (via Standings)
	standings, err := uc.repo.GetStandings(ctx, clubID, groupID)
	if err != nil {
		return nil, err
	}
	if len(standings) < 2 {
		return nil, errors.New("at least 2 teams are required to generate fixture")
	}

	// 2. Generate Matches (Round Robin)
	// Simple algorithm:
	// If odd number of teams, add a dummy team.
	// Rotate teams fixes the first one.
	// Here we implement a simple all-vs-all for single round.

	// Ensure we have the Group info for Tournament/Stage IDs
	group, err := uc.repo.GetGroup(ctx, clubID, groupID)
	if err != nil {
		return nil, err
	}
	stage, err := uc.repo.GetStage(ctx, clubID, group.StageID.String())
	if err != nil {
		return nil, err
	}

	var matches []domain.TournamentMatch
	n := len(standings)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			match := domain.TournamentMatch{
				ID:           uuid.New(),
				TournamentID: stage.TournamentID,
				StageID:      stage.ID,
				GroupID:      &group.ID,
				HomeTeamID:   standings[i].TeamID,
				AwayTeamID:   standings[j].TeamID,
				Status:       domain.MatchScheduled,
				Date:         time.Now(), // Default to now or TBD
			}
			matches = append(matches, match)
		}
	}

	// Create all matches atomically in a single transaction
	if err := uc.repo.CreateMatchesBatch(ctx, clubID, matches); err != nil {
		return nil, err
	}

	return matches, nil
}

// ...

func (uc *ChampionshipUseCases) recalculateStandings(ctx context.Context, clubID, groupID string) error {
	matches, err := uc.repo.GetMatchesByGroup(ctx, clubID, groupID)
	if err != nil {
		return err
	}
	standings, err := uc.repo.GetStandings(ctx, clubID, groupID)
	if err != nil {
		return err
	}

	standingMap := make(map[uuid.UUID]*domain.Standing)
	for i := range standings {
		s := &standings[i]
		// Reset stats
		s.Points = 0
		s.Played = 0
		s.Won = 0
		s.Drawn = 0
		s.Lost = 0
		s.GoalsFor = 0
		s.GoalsAgainst = 0
		s.GoalDifference = 0
		standingMap[s.TeamID] = s
	}

	for _, m := range matches {
		if m.Status != domain.MatchCompleted || m.HomeScore == nil || m.AwayScore == nil {
			continue
		}

		home, okH := standingMap[m.HomeTeamID]
		away, okA := standingMap[m.AwayTeamID]
		if !okH || !okA {
			continue // Should not happen if referential integrity holds
		}

		home.Played++
		away.Played++
		home.GoalsFor += *m.HomeScore
		home.GoalsAgainst += *m.AwayScore
		away.GoalsFor += *m.AwayScore
		away.GoalsAgainst += *m.HomeScore

		home.GoalDifference = home.GoalsFor - home.GoalsAgainst
		away.GoalDifference = away.GoalsFor - away.GoalsAgainst

		if *m.HomeScore > *m.AwayScore {
			home.Points += 3
			home.Won++
			away.Lost++
		} else if *m.AwayScore > *m.HomeScore {
			away.Points += 3
			away.Won++
			home.Lost++
		} else {
			home.Points += 1
			away.Points += 1
			home.Drawn++
			away.Drawn++
		}
	}

	// Identify changed standings
	var standingsToUpdate []domain.Standing
	standingsToUpdate = append(standingsToUpdate, standings...)

	if len(standingsToUpdate) > 0 {
		if err := uc.repo.UpdateStandingsBatch(ctx, clubID, standingsToUpdate); err != nil {
			return err
		}
	}

	return nil
}

// ...

type GenerateKnockoutBracketInput struct {
	ClubID    string   `json:"club_id"`
	StageID   string   `json:"stage_id"`
	SeedOrder []string `json:"seed_order"` // Team IDs in order of seeding
}

// GenerateKnockoutBracket generates elimination bracket matches for a stage.
func (uc *ChampionshipUseCases) GenerateKnockoutBracket(ctx context.Context, input GenerateKnockoutBracketInput) ([]domain.TournamentMatch, error) {
	if len(input.SeedOrder) < 2 {
		return nil, errors.New("at least 2 teams are required")
	}
	// Check if power of 2
	n := len(input.SeedOrder)
	if n&(n-1) != 0 {
		return nil, errors.New("number of teams must be a power of 2")
	}

	stage, err := uc.repo.GetStage(ctx, input.ClubID, input.StageID)
	if err != nil {
		if err.Error() == "record not found" {
			// Handle implementation detail of repo return
			return nil, errors.New("stage not found")
		}
		if stage == nil {
			return nil, errors.New("stage not found")
		}
		return nil, err
	}
	if stage.Type != domain.StageKnockout {
		return nil, errors.New("stage is not a knockout stage")
	}

	var matches []domain.TournamentMatch
	// Simple pairing: 1 vs N, 2 vs N-1, etc.
	// This assumes single elimination round 1.
	for i := 0; i < n/2; i++ {
		homeID, err := uuid.Parse(input.SeedOrder[i])
		if err != nil {
			return nil, errors.New("invalid team ID")
		}
		awayID, err := uuid.Parse(input.SeedOrder[n-1-i])
		if err != nil {
			return nil, errors.New("invalid team ID")
		}

		match := domain.TournamentMatch{
			ID:           uuid.New(),
			TournamentID: stage.TournamentID,
			StageID:      stage.ID,
			HomeTeamID:   homeID,
			AwayTeamID:   awayID,
			Status:       domain.MatchScheduled,
			Date:         time.Now(),
		}
		matches = append(matches, match)
	}

	// 5. Create all matches atomically
	if err := uc.repo.CreateMatchesBatch(ctx, input.ClubID, matches); err != nil {
		return nil, err
	}

	return matches, nil
}

type UpdateMatchResultInput struct {
	ClubID    string  `json:"club_id"`
	MatchID   string  `json:"match_id"`
	HomeScore float64 `json:"home_score"`
	AwayScore float64 `json:"away_score"`
}

func (uc *ChampionshipUseCases) UpdateMatchResult(ctx context.Context, input UpdateMatchResultInput) error {
	// 1. Update Match Score
	if err := uc.repo.UpdateMatchResult(ctx, input.ClubID, input.MatchID, input.HomeScore, input.AwayScore); err != nil {
		return err
	}

	// 2. Trigger async updates (XP, Standings)
	// We do this synchronously here for simplicity, but ideally async.

	match, err := uc.repo.GetMatch(ctx, input.ClubID, input.MatchID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.New("match not found after update")
	}

	if match.GroupID != nil {
		if err := uc.recalculateStandings(ctx, input.ClubID, match.GroupID.String()); err != nil {
			return err
		}
	}

	// XP Update
	homeMembers, _ := uc.repo.GetTeamMembers(ctx, match.HomeTeamID.String())
	awayMembers, _ := uc.repo.GetTeamMembers(ctx, match.AwayTeamID.String())

	homeWon := input.HomeScore > input.AwayScore
	awayWon := input.AwayScore > input.HomeScore

	for _, uid := range homeMembers {
		_ = uc.userService.UpdateMatchStats(ctx, input.ClubID, uid, homeWon, 100)
	}
	for _, uid := range awayMembers {
		_ = uc.userService.UpdateMatchStats(ctx, input.ClubID, uid, awayWon, 100)
	}

	return nil
}

type ScheduleMatchInput struct {
	ClubID  string    `json:"club_id"`
	MatchID string    `json:"match_id"`
	CourtID string    `json:"court_id"`
	Date    time.Time `json:"date"`
}

func (uc *ChampionshipUseCases) ScheduleMatch(ctx context.Context, input ScheduleMatchInput) error {
	// Default match duration 90m
	duration := 90 * time.Minute
	endTime := input.Date.Add(duration)

	bookingID, err := uc.bookingService.CreateSystemBooking(input.ClubID, input.CourtID, input.Date, endTime, "Championship Match")
	if err != nil {
		return errors.New("failed to book court: " + err.Error())
	}

	if err := uc.repo.UpdateMatchScheduling(ctx, input.ClubID, input.MatchID, input.Date, *bookingID); err != nil {
		return err
	}

	return nil
}

type CreateTeamInput struct {
	ClubID string `json:"club_id"`
	Name   string `json:"name"`
}

func (uc *ChampionshipUseCases) CreateTeam(ctx context.Context, input CreateTeamInput) (*domain.Team, error) {
	team := &domain.Team{
		ID:   uuid.New(),
		Name: input.Name,
	}
	if err := uc.repo.CreateTeam(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (uc *ChampionshipUseCases) AddMember(ctx context.Context, clubID, teamID, userID string) error {
	return uc.repo.AddMember(ctx, teamID, userID)
}

func (uc *ChampionshipUseCases) GetMyMatches(ctx context.Context, clubID, userID string) ([]domain.TournamentMatch, error) {
	return uc.repo.GetMatchesByUserID(ctx, clubID, userID)
}

type HeadToHeadResult struct {
	Matches    []domain.TournamentMatch `json:"matches"`
	TeamAWins  int                      `json:"team_a_wins"`
	TeamBWins  int                      `json:"team_b_wins"`
	Draws      int                      `json:"draws"`
	TeamAGoals float64                  `json:"team_a_goals"`
	TeamBGoals float64                  `json:"team_b_goals"`
}

func (uc *ChampionshipUseCases) GetHeadToHeadHistory(ctx context.Context, clubID, groupID, teamAID, teamBID string) (*HeadToHeadResult, error) {
	matches, err := uc.repo.GetMatchesByGroup(ctx, clubID, groupID)
	if err != nil {
		return nil, err
	}

	tA, err := uuid.Parse(teamAID)
	if err != nil {
		return nil, errors.New("invalid team A ID")
	}
	tB, err := uuid.Parse(teamBID)
	if err != nil {
		return nil, errors.New("invalid team B ID")
	}

	res := &HeadToHeadResult{
		Matches: []domain.TournamentMatch{},
	}

	for _, m := range matches {
		relevant := (m.HomeTeamID == tA && m.AwayTeamID == tB) || (m.HomeTeamID == tB && m.AwayTeamID == tA)
		if !relevant {
			continue
		}
		if m.Status != domain.MatchCompleted || m.HomeScore == nil || m.AwayScore == nil {
			continue
		}

		res.Matches = append(res.Matches, m)

		scoreA := 0.0
		scoreB := 0.0
		if m.HomeTeamID == tA {
			scoreA = *m.HomeScore
			scoreB = *m.AwayScore
		} else {
			scoreA = *m.AwayScore
			scoreB = *m.HomeScore
		}

		res.TeamAGoals += scoreA
		res.TeamBGoals += scoreB

		if scoreA > scoreB {
			res.TeamAWins++
		} else if scoreB > scoreA {
			res.TeamBWins++
		} else {
			res.Draws++
		}
	}

	return res, nil
}
