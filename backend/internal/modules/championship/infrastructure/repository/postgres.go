package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"

	"gorm.io/gorm"
)

type PostgresChampionshipRepository struct {
	db *gorm.DB
}

func NewPostgresChampionshipRepository(db *gorm.DB) *PostgresChampionshipRepository {
	return &PostgresChampionshipRepository{db: db}
}

func (r *PostgresChampionshipRepository) CreateTournament(tournament *domain.Tournament) error {
	return r.db.Create(tournament).Error
}

func (r *PostgresChampionshipRepository) GetTournament(id string) (*domain.Tournament, error) {
	var tournament domain.Tournament
	if err := r.db.Preload("Stages.Groups").First(&tournament, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *PostgresChampionshipRepository) ListTournaments(clubID string) ([]domain.Tournament, error) {
	var tournaments []domain.Tournament
	if err := r.db.Preload("Stages.Groups").Where("club_id = ?", clubID).Find(&tournaments).Error; err != nil {
		return nil, err
	}
	return tournaments, nil
}

func (r *PostgresChampionshipRepository) CreateStage(stage *domain.TournamentStage) error {
	return r.db.Create(stage).Error
}

func (r *PostgresChampionshipRepository) GetStage(id string) (*domain.TournamentStage, error) {
	var stage domain.TournamentStage
	err := r.db.First(&stage, "id = ?", id).Error
	return &stage, err
}

func (r *PostgresChampionshipRepository) CreateGroup(group *domain.Group) error {
	return r.db.Create(group).Error
}

func (r *PostgresChampionshipRepository) GetGroup(id string) (*domain.Group, error) {
	var group domain.Group
	// Preload Stage to get TournamentID eventually
	err := r.db.Preload("Standings").First(&group, "id = ?", id).Error
	return &group, err
}

func (r *PostgresChampionshipRepository) CreateMatch(match *domain.TournamentMatch) error {
	return r.db.Create(match).Error
}

func (r *PostgresChampionshipRepository) GetMatch(id string) (*domain.TournamentMatch, error) {
	var match domain.TournamentMatch
	err := r.db.First(&match, "id = ?", id).Error
	return &match, err
}

func (r *PostgresChampionshipRepository) GetMatchesByGroup(groupID string) ([]domain.TournamentMatch, error) {
	var matches []domain.TournamentMatch
	// Join with teams (twice, for home and away)
	// Assuming table name is "teams"
	err := r.db.Table("tournament_matches").
		Select("tournament_matches.*, h.name as home_team_name, a.name as away_team_name").
		Joins("LEFT JOIN teams h ON h.id = tournament_matches.home_team_id").
		Joins("LEFT JOIN teams a ON a.id = tournament_matches.away_team_id").
		Where("tournament_matches.group_id = ?", groupID).
		Scan(&matches).Error
	return matches, err
}

func (r *PostgresChampionshipRepository) UpdateMatchResult(matchID string, homeScore, awayScore int) error {
	return r.db.Model(&domain.TournamentMatch{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"home_score": homeScore,
		"away_score": awayScore,
		"status":     domain.MatchCompleted,
	}).Error
}

func (r *PostgresChampionshipRepository) UpdateMatchScheduling(matchID string, date time.Time, bookingID uuid.UUID) error {
	return r.db.Model(&domain.TournamentMatch{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"date":       date,
		"booking_id": bookingID,
		"status":     domain.MatchScheduled,
	}).Error
}

func (r *PostgresChampionshipRepository) GetStandings(groupID string) ([]domain.Standing, error) {
	var standings []domain.Standing
	// Order by Points DESC, Goal Difference DESC, Goals For DESC
	// And Join with teams
	err := r.db.Table("standings").
		Select("standings.*, teams.name as team_name").
		Joins("LEFT JOIN teams ON teams.id = standings.team_id").
		Where("standings.group_id = ?", groupID).
		Order("standings.points DESC, standings.goal_difference DESC, standings.goals_for DESC").
		Scan(&standings).Error
	return standings, err
}

func (r *PostgresChampionshipRepository) RegisterTeam(standing *domain.Standing) error {
	return r.db.Create(standing).Error
}

func (r *PostgresChampionshipRepository) UpdateStanding(standing *domain.Standing) error {
	return r.db.Save(standing).Error
}

func (r *PostgresChampionshipRepository) GetTeamMembers(teamID string) ([]string, error) {
	var userIDs []string
	// Assuming table 'team_members' with user_id column
	err := r.db.Table("team_members").
		Select("user_id").
		Where("team_id = ?", teamID).
		Scan(&userIDs).Error
	return userIDs, err
}
