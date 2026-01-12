package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresChampionshipRepository struct {
	db *gorm.DB
}

func NewPostgresChampionshipRepository(db *gorm.DB) *PostgresChampionshipRepository {
	return &PostgresChampionshipRepository{db: db}
}

func (r *PostgresChampionshipRepository) CreateTournament(ctx context.Context, tournament *domain.Tournament) error {
	return r.db.WithContext(ctx).Create(tournament).Error
}

func (r *PostgresChampionshipRepository) GetTournament(ctx context.Context, clubID, id string) (*domain.Tournament, error) {
	var tournament domain.Tournament
	if err := r.db.WithContext(ctx).Preload("Stages.Groups").First(&tournament, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *PostgresChampionshipRepository) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	var tournaments []domain.Tournament
	if err := r.db.WithContext(ctx).Preload("Stages.Groups").Where("club_id = ?", clubID).Find(&tournaments).Error; err != nil {
		return nil, err
	}
	return tournaments, nil
}

func (r *PostgresChampionshipRepository) CreateStage(ctx context.Context, stage *domain.TournamentStage) error {
	return r.db.WithContext(ctx).Create(stage).Error
}

func (r *PostgresChampionshipRepository) GetStage(ctx context.Context, clubID, id string) (*domain.TournamentStage, error) {
	var stage domain.TournamentStage
	// Join with Tournament to check club_id
	err := r.db.WithContext(ctx).Joins("JOIN championships ON championships.id = tournament_stages.tournament_id").
		Where("tournament_stages.id = ? AND championships.club_id = ?", id, clubID).
		First(&stage).Error
	return &stage, err
}

func (r *PostgresChampionshipRepository) CreateGroup(ctx context.Context, group *domain.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *PostgresChampionshipRepository) GetGroup(ctx context.Context, clubID, id string) (*domain.Group, error) {
	var group domain.Group
	// Join Group -> Stage -> Tournament to check club_id
	err := r.db.WithContext(ctx).Preload("Standings").
		Joins("JOIN tournament_stages ON tournament_stages.id = groups.stage_id").
		Joins("JOIN championships ON championships.id = tournament_stages.tournament_id").
		Where("groups.id = ? AND championships.club_id = ?", id, clubID).
		First(&group).Error
	return &group, err
}

func (r *PostgresChampionshipRepository) CreateMatch(ctx context.Context, match *domain.TournamentMatch) error {
	return r.db.WithContext(ctx).Create(match).Error
}

// CreateMatchesBatch creates multiple matches atomically using a database transaction.
// If any match fails to create, the entire batch is rolled back.
func (r *PostgresChampionshipRepository) CreateMatchesBatch(ctx context.Context, matches []domain.TournamentMatch) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range matches {
			if err := tx.Create(&matches[i]).Error; err != nil {
				return err // Transaction will be rolled back
			}
		}
		return nil
	})
}

func (r *PostgresChampionshipRepository) GetMatch(ctx context.Context, clubID, id string) (*domain.TournamentMatch, error) {
	var match domain.TournamentMatch
	// Join Tournament to check club_id
	err := r.db.WithContext(ctx).Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Where("tournament_matches.id = ? AND championships.club_id = ?", id, clubID).
		First(&match).Error
	return &match, err
}

func (r *PostgresChampionshipRepository) GetMatchesByGroup(ctx context.Context, clubID, groupID string) ([]domain.TournamentMatch, error) {
	var matches []domain.TournamentMatch
	// Validate membership to club via join
	err := r.db.WithContext(ctx).Table("tournament_matches").
		Select("tournament_matches.*, h.name as home_team_name, a.name as away_team_name").
		Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Joins("LEFT JOIN teams h ON h.id = tournament_matches.home_team_id").
		Joins("LEFT JOIN teams a ON a.id = tournament_matches.away_team_id").
		Where("tournament_matches.group_id = ? AND championships.club_id = ?", groupID, clubID).
		Scan(&matches).Error
	return matches, err
}

func (r *PostgresChampionshipRepository) UpdateMatchResult(ctx context.Context, clubID, matchID string, homeScore, awayScore int) error {
	// Verify club ownership before update
	var count int64
	r.db.WithContext(ctx).Table("tournament_matches").
		Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Where("tournament_matches.id = ? AND championships.club_id = ?", matchID, clubID).
		Count(&count)

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.WithContext(ctx).Model(&domain.TournamentMatch{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"home_score": homeScore,
		"away_score": awayScore,
		"status":     domain.MatchCompleted,
	}).Error
}

func (r *PostgresChampionshipRepository) UpdateMatchScheduling(ctx context.Context, clubID, matchID string, date time.Time, bookingID uuid.UUID) error {
	// Verify club ownership before update
	var count int64
	r.db.WithContext(ctx).Table("tournament_matches").
		Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Where("tournament_matches.id = ? AND championships.club_id = ?", matchID, clubID).
		Count(&count)

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.WithContext(ctx).Model(&domain.TournamentMatch{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"date":       date,
		"booking_id": bookingID,
		"status":     domain.MatchScheduled,
	}).Error
}

func (r *PostgresChampionshipRepository) GetStandings(ctx context.Context, clubID, groupID string) ([]domain.Standing, error) {
	var standings []domain.Standing
	// Join with Group -> Stage -> Tournament to check club_id
	err := r.db.WithContext(ctx).Table("standings").
		Select("standings.*, teams.name as team_name").
		Joins("JOIN groups ON groups.id = standings.group_id").
		Joins("JOIN tournament_stages ON tournament_stages.id = groups.stage_id").
		Joins("JOIN championships ON championships.id = tournament_stages.tournament_id").
		Joins("LEFT JOIN teams ON teams.id = standings.team_id").
		Where("standings.group_id = ? AND championships.club_id = ?", groupID, clubID).
		Order("standings.points DESC, standings.goal_difference DESC, standings.goals_for DESC").
		Scan(&standings).Error
	return standings, err
}

func (r *PostgresChampionshipRepository) RegisterTeam(ctx context.Context, standing *domain.Standing) error {
	return r.db.WithContext(ctx).Create(standing).Error
}

func (r *PostgresChampionshipRepository) UpdateStanding(ctx context.Context, standing *domain.Standing) error {
	return r.db.WithContext(ctx).Save(standing).Error
}

func (r *PostgresChampionshipRepository) UpdateStandingsBatch(ctx context.Context, standings []domain.Standing) error {
	// Use GORM's Clauses to perform a bulk upsert (INSERT ... ON CONFLICT DO UPDATE)
	// This generates a single SQL statement instead of N updates.
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // Conflict on Primary Key
		UpdateAll: true,                          // Update all columns if conflict
	}).Save(&standings).Error
}

func (r *PostgresChampionshipRepository) GetTeamMembers(ctx context.Context, teamID string) ([]string, error) {
	var userIDs []string
	// Assuming table 'team_members' with user_id column
	err := r.db.WithContext(ctx).Table("team_members").
		Select("user_id").
		Where("team_id = ?", teamID).
		Scan(&userIDs).Error
	return userIDs, err
}
