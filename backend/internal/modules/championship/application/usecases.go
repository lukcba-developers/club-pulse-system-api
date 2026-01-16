package application

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
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
	if err := uc.repo.RegisterTeam(ctx, standing); err != nil {
		return nil, err
	}
	return standing, nil
}

type CreateTeamInput struct {
	ClubID  string `json:"club_id"`
	Name    string `json:"name" binding:"required"`
	LogoURL string `json:"logo_url"`
	Contact string `json:"contact"`
}

func (uc *ChampionshipUseCases) CreateTeam(ctx context.Context, input CreateTeamInput) (*domain.Team, error) {
	// TODO: Validate ClubID if Team belongs to Club (current domain.Team doesn't have ClubID field, but likely should.
	// For now, teams are global or shared? Actually Championship domain Team doesn't have ClubID.
	// Assuming they are shared or we rely on repo logic if we add ClubID later.
	// Ideally we should add ClubID to Team struct, but keeping minimal changes as requested.)
	// Wait, if Name is unique per Club?

	team := &domain.Team{
		ID:        uuid.New(),
		Name:      input.Name,
		LogoURL:   input.LogoURL,
		Contact:   input.Contact,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.CreateTeam(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (uc *ChampionshipUseCases) AddMember(ctx context.Context, clubID, teamID, userID string) error {
	// 1. Verify Team belongs to Club?
	// Current Team struct has no ClubID, but we assume it's global or implicitly trustworthy for Admin.
	// But ideally we should check if User belongs to Club (ClubID check).
	// We lack user repo here to check club membership easily without importing user module.
	// Proceeding with adding member.
	return uc.repo.AddMember(ctx, teamID, userID)
}

func (uc *ChampionshipUseCases) GenerateGroupFixture(ctx context.Context, clubID, groupID string) ([]domain.TournamentMatch, error) {
	// 1. Get Teams in Group (via Standings)
	standings, err := uc.repo.GetStandings(ctx, clubID, groupID)
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

	group, err := uc.repo.GetGroup(ctx, clubID, groupID)
	if err != nil {
		return nil, err
	}

	stage, err := uc.repo.GetStage(ctx, clubID, group.StageID.String())
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
	if err := uc.repo.CreateMatchesBatch(ctx, matches); err != nil {
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
	if err := uc.repo.UpdateMatchResult(ctx, input.ClubID, input.MatchID, input.HomeScore, input.AwayScore); err != nil {
		return err
	}

	// Trigger Recalculate Standings
	// 1. Get Match to find GroupID
	match, err := uc.repo.GetMatch(ctx, input.ClubID, input.MatchID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.New("match not found after update")
	}

	// Update User Stats
	if uc.userService != nil {
		tournament, err := uc.repo.GetTournament(ctx, input.ClubID, match.TournamentID.String())
		if err == nil {
			clubID := tournament.ClubID.String()
			homeWon := input.HomeScore > input.AwayScore
			awayWon := input.AwayScore > input.HomeScore

			// XP Logic
			xp := 100 // base XP

			// Home Players
			homePlayers, _ := uc.repo.GetTeamMembers(ctx, match.HomeTeamID.String())
			for _, userID := range homePlayers {
				_ = uc.userService.UpdateMatchStats(ctx, clubID, userID, homeWon, xp)
			}

			// Away Players
			awayPlayers, _ := uc.repo.GetTeamMembers(ctx, match.AwayTeamID.String())
			for _, userID := range awayPlayers {
				_ = uc.userService.UpdateMatchStats(ctx, clubID, userID, awayWon, xp)
			}
		}
	}

	if match.GroupID == nil {
		return nil // Not a group match
	}

	// Safe dereference as we checked for nil
	groupID := (*match.GroupID).String()
	return uc.recalculateStandings(ctx, input.ClubID, groupID)
}

func (uc *ChampionshipUseCases) GetMyMatches(ctx context.Context, clubID, userID string) ([]domain.TournamentMatch, error) {
	return uc.repo.GetMatchesByUserID(ctx, clubID, userID)
}

// HeadToHeadResult represents the summary and history of matches between two teams
type HeadToHeadResult struct {
	TeamAID    string                   `json:"team_a_id"`
	TeamBID    string                   `json:"team_b_id"`
	TeamAWins  int                      `json:"team_a_wins"`
	TeamBWins  int                      `json:"team_b_wins"`
	Draws      int                      `json:"draws"`
	TeamAGoals int                      `json:"team_a_goals"`
	TeamBGoals int                      `json:"team_b_goals"`
	Matches    []domain.TournamentMatch `json:"matches"`
}

func (uc *ChampionshipUseCases) GetHeadToHeadHistory(ctx context.Context, clubID, groupID, teamAID, teamBID string) (*HeadToHeadResult, error) {
	matches, err := uc.repo.GetMatchesByGroup(ctx, clubID, groupID)
	if err != nil {
		return nil, err
	}

	teamA := uuid.MustParse(teamAID)
	teamB := uuid.MustParse(teamBID)

	result := &HeadToHeadResult{
		TeamAID: teamAID,
		TeamBID: teamBID,
	}

	for _, m := range matches {
		if m.Status != domain.MatchCompleted || m.HomeScore == nil || m.AwayScore == nil {
			continue
		}

		isH2H := (m.HomeTeamID == teamA && m.AwayTeamID == teamB) ||
			(m.HomeTeamID == teamB && m.AwayTeamID == teamA)
		if !isH2H {
			continue
		}

		result.Matches = append(result.Matches, m)
		homeScore := int(*m.HomeScore)
		awayScore := int(*m.AwayScore)

		// Normalize: teamA goals vs teamB goals
		var aGoals, bGoals int
		if m.HomeTeamID == teamA {
			aGoals, bGoals = homeScore, awayScore
		} else {
			aGoals, bGoals = awayScore, homeScore
		}

		result.TeamAGoals += aGoals
		result.TeamBGoals += bGoals

		if aGoals > bGoals {
			result.TeamAWins++
		} else if bGoals > aGoals {
			result.TeamBWins++
		} else {
			result.Draws++
		}
	}

	return result, nil
}

// compareHeadToHead returns:
//   - positive if teamA has advantage over teamB in direct matches
//   - negative if teamB has advantage
//   - 0 if tied or no direct matches
func compareHeadToHead(teamA, teamB uuid.UUID, matches []domain.TournamentMatch) int {
	var aPoints, bPoints int
	for _, m := range matches {
		if m.Status != domain.MatchCompleted || m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		homeScore := *m.HomeScore
		awayScore := *m.AwayScore

		// Match between teamA (home) vs teamB (away)
		if m.HomeTeamID == teamA && m.AwayTeamID == teamB {
			if homeScore > awayScore {
				aPoints += 3
			} else if awayScore > homeScore {
				bPoints += 3
			} else {
				aPoints++
				bPoints++
			}
		}
		// Match between teamB (home) vs teamA (away)
		if m.HomeTeamID == teamB && m.AwayTeamID == teamA {
			if homeScore > awayScore {
				bPoints += 3
			} else if awayScore > homeScore {
				aPoints += 3
			} else {
				aPoints++
				bPoints++
			}
		}
	}
	return aPoints - bPoints
}

func (uc *ChampionshipUseCases) recalculateStandings(ctx context.Context, clubID, groupID string) error {
	standings, err := uc.repo.GetStandings(ctx, clubID, groupID)
	if err != nil {
		return err
	}

	matches, err := uc.repo.GetMatchesByGroup(ctx, clubID, groupID)
	if err != nil {
		return err
	}

	// 1. Get Group for StageID
	group, err := uc.repo.GetGroup(ctx, clubID, groupID)
	if err != nil {
		return err
	}

	// 2. Get Stage for TournamentID
	stage, err := uc.repo.GetStage(ctx, clubID, group.StageID.String())
	if err != nil {
		return err
	}

	// 3. Get Tournament for Settings
	tournament, err := uc.repo.GetTournament(ctx, clubID, stage.TournamentID.String())
	if err != nil {
		return err
	}

	// Default Points
	pointsWin := 3.0
	pointsDraw := 1.0

	// Parse custom points if available
	if tournament.Settings != nil {
		var settings struct {
			PointsWin  float64 `json:"points_win"`
			PointsDraw float64 `json:"points_draw"`
		}

		if err := json.Unmarshal(tournament.Settings, &settings); err == nil {
			if settings.PointsWin > 0 {
				pointsWin = settings.PointsWin
			}
			if settings.PointsDraw > 0 {
				pointsDraw = settings.PointsDraw
			}
		}
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
				home.Points += pointsWin
				away.Lost++
			} else if awayScore > homeScore {
				away.Won++
				away.Points += pointsWin
				home.Lost++
			} else {
				home.Drawn++
				home.Points += pointsDraw
				away.Drawn++
				away.Points += pointsDraw
			}
		}
	}

	// Convert map to slice
	var standingsToUpdate []domain.Standing
	for _, s := range stats {
		standingsToUpdate = append(standingsToUpdate, *s)
	}

	// Parse tiebreaker criteria from settings
	tiebreakerCriteria := []string{"GOAL_DIFF", "GOALS_FOR"} // Default order
	if tournament.Settings != nil {
		var tSettings struct {
			TiebreakerCriteria []string `json:"tiebreaker_criteria"`
		}
		if err := json.Unmarshal(tournament.Settings, &tSettings); err == nil && len(tSettings.TiebreakerCriteria) > 0 {
			tiebreakerCriteria = tSettings.TiebreakerCriteria
		}
	}

	// Sort standings by Points (desc), then by tiebreaker criteria
	sort.SliceStable(standingsToUpdate, func(i, j int) bool {
		a, b := standingsToUpdate[i], standingsToUpdate[j]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		for _, criterion := range tiebreakerCriteria {
			switch criterion {
			case "GOAL_DIFF":
				if a.GoalDifference != b.GoalDifference {
					return a.GoalDifference > b.GoalDifference
				}
			case "GOALS_FOR":
				if a.GoalsFor != b.GoalsFor {
					return a.GoalsFor > b.GoalsFor
				}
			case "HEAD_TO_HEAD":
				// Compare head-to-head results between team A and team B
				h2hResult := compareHeadToHead(a.TeamID, b.TeamID, matches)
				if h2hResult != 0 {
					return h2hResult > 0 // Positive means A won more
				}
			}
		}
		return false // Maintain original order if all criteria are equal
	})

	// Assign position based on sorted order
	for i := range standingsToUpdate {
		standingsToUpdate[i].Position = i + 1
	}

	if len(standingsToUpdate) > 0 {
		if err := uc.repo.UpdateStandingsBatch(ctx, standingsToUpdate); err != nil {
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

func (uc *ChampionshipUseCases) ScheduleMatch(ctx context.Context, input ScheduleMatchInput) error {
	// 1. Create Booking via Service
	notes := "Partido de Torneo: " + input.MatchID
	bookingID, err := uc.bookingService.CreateSystemBooking(input.ClubID, input.CourtID, input.StartTime, input.EndTime, notes)
	if err != nil {
		return errors.New("failed to book court: " + err.Error())
	}

	// 2. Update Match with BookingID and Date
	return uc.repo.UpdateMatchScheduling(ctx, input.ClubID, input.MatchID, input.StartTime, *bookingID)
}

// GenerateKnockoutBracketInput defines the input for generating a knockout bracket.
type GenerateKnockoutBracketInput struct {
	ClubID    string   `json:"club_id"`
	StageID   string   `json:"stage_id"`
	SeedOrder []string `json:"seed_order"` // Team IDs in seed order (1st vs 8th, 2nd vs 7th, etc.)
}

// GenerateKnockoutBracket generates elimination bracket matches for a stage.
// It pairs teams based on seeding: #1 vs #N, #2 vs #(N-1), etc.
// Supports 2, 4, 8, 16, 32 team brackets (must be power of 2).
func (uc *ChampionshipUseCases) GenerateKnockoutBracket(ctx context.Context, input GenerateKnockoutBracketInput) ([]domain.TournamentMatch, error) {
	// 1. Validate stage exists and is KNOCKOUT type
	stage, err := uc.repo.GetStage(ctx, input.ClubID, input.StageID)
	if err != nil {
		return nil, err
	}
	if stage == nil {
		return nil, errors.New("stage not found")
	}
	if stage.Type != domain.StageKnockout {
		return nil, errors.New("stage must be of type KNOCKOUT to generate bracket")
	}

	// 2. Validate seed order count (must be power of 2 and >= 2)
	numTeams := len(input.SeedOrder)
	if numTeams < 2 {
		return nil, errors.New("need at least 2 teams to generate knockout bracket")
	}
	if (numTeams & (numTeams - 1)) != 0 {
		return nil, errors.New("number of teams must be a power of 2 (e.g., 2, 4, 8, 16, 32)")
	}

	// 3. Parse team UUIDs
	teamIDs := make([]uuid.UUID, numTeams)
	for i, idStr := range input.SeedOrder {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, errors.New("invalid team ID at position " + string(rune(i+1)))
		}
		teamIDs[i] = id
	}

	// 4. Generate bracket pairings (seed-based: 1 vs N, 2 vs N-1, ...)
	var matches []domain.TournamentMatch
	numMatches := numTeams / 2

	for i := 0; i < numMatches; i++ {
		homeTeam := teamIDs[i]            // Seed 1, 2, 3...
		awayTeam := teamIDs[numTeams-1-i] // Seed N, N-1, N-2...

		match := domain.TournamentMatch{
			ID:           uuid.New(),
			TournamentID: stage.TournamentID,
			StageID:      stage.ID,
			GroupID:      nil, // Knockout has no group
			HomeTeamID:   homeTeam,
			AwayTeamID:   awayTeam,
			Status:       domain.MatchScheduled,
			Date:         time.Now().Add(time.Hour * 24 * 7), // Scheduled 1 week out by default
		}
		matches = append(matches, match)
	}

	// 5. Create all matches atomically
	if err := uc.repo.CreateMatchesBatch(ctx, matches); err != nil {
		return nil, err
	}

	return matches, nil
}
