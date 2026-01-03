package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	"gorm.io/gorm"
)

type PostgresTournamentRepository struct {
	db *gorm.DB
}

func NewPostgresTournamentRepository(db *gorm.DB) *PostgresTournamentRepository {
	_ = db.AutoMigrate(&domain.Tournament{}, &domain.Team{}, &domain.Match{})
	return &PostgresTournamentRepository{db: db}
}

// --- Tournament ---

func (r *PostgresTournamentRepository) CreateTournament(tournament *domain.Tournament) error {
	return r.db.Create(tournament).Error
}

func (r *PostgresTournamentRepository) GetTournamentByID(clubID string, id uuid.UUID) (*domain.Tournament, error) {
	var tournament domain.Tournament
	if err := r.db.Preload("Teams").Preload("Matches").First(&tournament, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tournament, nil
}

func (r *PostgresTournamentRepository) ListTournaments(clubID string) ([]domain.Tournament, error) {
	var tournaments []domain.Tournament
	err := r.db.Where("club_id = ?", clubID).Find(&tournaments).Error
	return tournaments, err
}

func (r *PostgresTournamentRepository) UpdateTournament(tournament *domain.Tournament) error {
	return r.db.Save(tournament).Error
}

// --- Team ---

func (r *PostgresTournamentRepository) CreateTeam(team *domain.Team) error {
	return r.db.Create(team).Error
}

func (r *PostgresTournamentRepository) GetTeamByID(clubID string, id uuid.UUID) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.First(&team, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (r *PostgresTournamentRepository) ListTeams(clubID string, tournamentID uuid.UUID) ([]domain.Team, error) {
	var teams []domain.Team
	err := r.db.Where("tournament_id = ? AND club_id = ?", tournamentID, clubID).Find(&teams).Error
	return teams, err
}

// --- Match ---

func (r *PostgresTournamentRepository) CreateMatch(match *domain.Match) error {
	return r.db.Create(match).Error
}

func (r *PostgresTournamentRepository) UpdateMatch(match *domain.Match) error {
	return r.db.Save(match).Error
}

func (r *PostgresTournamentRepository) GetMatchByID(clubID string, id uuid.UUID) (*domain.Match, error) {
	var match domain.Match
	if err := r.db.First(&match, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *PostgresTournamentRepository) ListMatches(clubID string, tournamentID uuid.UUID) ([]domain.Match, error) {
	var matches []domain.Match
	err := r.db.Where("tournament_id = ? AND club_id = ?", tournamentID, clubID).Order("start_time asc").Find(&matches).Error
	return matches, err
}

// --- Standings ---

func (r *PostgresTournamentRepository) GetStandings(clubID string, tournamentID uuid.UUID) ([]domain.Standing, error) {
	// 1. Get all matches for tournament
	matches, err := r.ListMatches(clubID, tournamentID)
	if err != nil {
		return nil, err
	}

	// 2. Get all teams to initialize map
	teams, err := r.ListTeams(clubID, tournamentID)
	if err != nil {
		return nil, err
	}

	standingMap := make(map[uuid.UUID]*domain.Standing)
	for _, team := range teams {
		standingMap[team.ID] = &domain.Standing{
			TournamentID: tournamentID,
			TeamID:       team.ID,
			TeamName:     team.Name,
		}
	}

	// 3. Iterate matches
	for _, m := range matches {
		if m.Status != domain.MatchStatusPlayed {
			continue
		}

		home := standingMap[m.HomeTeamID]
		away := standingMap[m.AwayTeamID]

		// Safety check if team was deleted but match exists
		if home == nil || away == nil {
			continue
		}

		home.Played++
		away.Played++
		home.GoalsFor += m.ScoreHome
		home.GoalsAgainst += m.ScoreAway
		away.GoalsFor += m.ScoreAway
		away.GoalsAgainst += m.ScoreHome

		if m.ScoreHome > m.ScoreAway {
			home.Won++
			home.Points += 3
			away.Lost++
		} else if m.ScoreAway > m.ScoreHome {
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

	// 4. Convert map to slice
	var standings []domain.Standing
	for _, s := range standingMap {
		standings = append(standings, *s)
	}

	// Sort by points (desc), then goal diff? For MVP just return list. UseCase can sort.
	return standings, nil
}
