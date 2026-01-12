package repository

import (
	"context"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	"gorm.io/gorm"
)

type PostgresTeamRepository struct {
	db *gorm.DB
}

func NewPostgresTeamRepository(db *gorm.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (r *PostgresTeamRepository) CreateMatchEvent(ctx context.Context, event *domain.MatchEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *PostgresTeamRepository) GetMatchEvent(ctx context.Context, clubID, id string) (*domain.MatchEvent, error) {
	var event domain.MatchEvent
	// Join with TrainingGroup to check club_id (assuming training_groups has club_id)
	err := r.db.WithContext(ctx).Joins("JOIN training_groups ON training_groups.id = match_events.training_group_id").
		Where("match_events.id = ? AND training_groups.club_id = ?", id, clubID).
		First(&event).Error
	return &event, err
}

func (r *PostgresTeamRepository) SetPlayerAvailability(ctx context.Context, availability *domain.PlayerAvailability) error {
	// Upsert: On conflict update status and reason
	return r.db.WithContext(ctx).Save(availability).Error
}

func (r *PostgresTeamRepository) GetEventAvailabilities(ctx context.Context, clubID, eventID string) ([]domain.PlayerAvailability, error) {
	var availabilities []domain.PlayerAvailability
	// Join MatchEvent -> TrainingGroup to check club_id
	err := r.db.WithContext(ctx).Table("player_availabilities").
		Joins("JOIN match_events ON match_events.id = player_availabilities.match_event_id").
		Joins("JOIN training_groups ON training_groups.id = match_events.training_group_id").
		Where("player_availabilities.match_event_id = ? AND training_groups.club_id = ?", eventID, clubID).
		Find(&availabilities).Error
	return availabilities, err
}
