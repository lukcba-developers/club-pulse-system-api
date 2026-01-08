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

func (r *PostgresChampionshipRepository) GetTournament(clubID, id string) (*domain.Tournament, error) {
	var tournament domain.Tournament
	if err := r.db.Preload("Stages.Groups").First(&tournament, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
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

func (r *PostgresChampionshipRepository) GetStage(clubID, id string) (*domain.TournamentStage, error) {
	var stage domain.TournamentStage
	// Join with Tournament to check club_id
	err := r.db.Joins("JOIN championships ON championships.id = tournament_stages.tournament_id").
		Where("tournament_stages.id = ? AND championships.club_id = ?", id, clubID).
		First(&stage).Error
	return &stage, err
}

func (r *PostgresChampionshipRepository) CreateGroup(group *domain.Group) error {
	return r.db.Create(group).Error
}

func (r *PostgresChampionshipRepository) GetGroup(clubID, id string) (*domain.Group, error) {
	var group domain.Group
	// Join Group -> Stage -> Tournament to check club_id
	// Assuming table names: groups, tournament_stages, championships
	// Since GORM might infer singular/plural, need to be careful. The model says "TournamentStage" so it might be "tournament_stages".
	// The models use implicit naming. Let's assume standard snake_case plural.
	err := r.db.Preload("Standings").
		Joins("JOIN tournament_stages ON tournament_stages.id = groups.stage_id").
		Joins("JOIN championships ON championships.id = tournament_stages.tournament_id").
		Where("groups.id = ? AND championships.club_id = ?", id, clubID).
		First(&group).Error
	return &group, err
}

func (r *PostgresChampionshipRepository) CreateMatch(match *domain.TournamentMatch) error {
	return r.db.Create(match).Error
}

func (r *PostgresChampionshipRepository) GetMatch(clubID, id string) (*domain.TournamentMatch, error) {
	var match domain.TournamentMatch
	// Join Tournament to check club_id
	err := r.db.Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Where("tournament_matches.id = ? AND championships.club_id = ?", id, clubID).
		First(&match).Error
	return &match, err
}

func (r *PostgresChampionshipRepository) GetMatchesByGroup(clubID, groupID string) ([]domain.TournamentMatch, error) {
	var matches []domain.TournamentMatch
	// Validate membership to club via join
	err := r.db.Table("tournament_matches").
		Select("tournament_matches.*, h.name as home_team_name, a.name as away_team_name").
		Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Joins("LEFT JOIN teams h ON h.id = tournament_matches.home_team_id").
		Joins("LEFT JOIN teams a ON a.id = tournament_matches.away_team_id").
		Where("tournament_matches.group_id = ? AND championships.club_id = ?", groupID, clubID).
		Scan(&matches).Error
	return matches, err
}

func (r *PostgresChampionshipRepository) UpdateMatchResult(clubID, matchID string, homeScore, awayScore int) error {
	// Verify club ownership before update
	var count int64
	r.db.Table("tournament_matches").
		Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Where("tournament_matches.id = ? AND championships.club_id = ?", matchID, clubID).
		Count(&count)

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.Model(&domain.TournamentMatch{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"home_score": homeScore,
		"away_score": awayScore,
		"status":     domain.MatchCompleted,
	}).Error
}

func (r *PostgresChampionshipRepository) UpdateMatchScheduling(clubID, matchID string, date time.Time, bookingID uuid.UUID) error {
	// Verify club ownership before update
	var count int64
	r.db.Table("tournament_matches").
		Joins("JOIN championships ON championships.id = tournament_matches.tournament_id").
		Where("tournament_matches.id = ? AND championships.club_id = ?", matchID, clubID).
		Count(&count)

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.Model(&domain.TournamentMatch{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"date":       date,
		"booking_id": bookingID,
		"status":     domain.MatchScheduled,
	}).Error
}

func (r *PostgresChampionshipRepository) GetStandings(clubID, groupID string) ([]domain.Standing, error) {
	var standings []domain.Standing
	// Join with Group -> Stage -> Tournament to check club_id
	err := r.db.Table("standings").
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
