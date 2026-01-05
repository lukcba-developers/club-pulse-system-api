package repository

import (
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	"gorm.io/gorm"
)

type PostgresTeamRepository struct {
	db *gorm.DB
}

func NewPostgresTeamRepository(db *gorm.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (r *PostgresTeamRepository) CreateMatchEvent(event *domain.MatchEvent) error {
	return r.db.Create(event).Error
}

func (r *PostgresTeamRepository) GetMatchEvent(id string) (*domain.MatchEvent, error) {
	var event domain.MatchEvent
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *PostgresTeamRepository) SetPlayerAvailability(availability *domain.PlayerAvailability) error {
	// Upsert: On conflict update status and reason
	return r.db.Save(availability).Error
}

func (r *PostgresTeamRepository) GetEventAvailabilities(eventID string) ([]domain.PlayerAvailability, error) {
	var availabilities []domain.PlayerAvailability
	if err := r.db.Where("match_event_id = ?", eventID).Find(&availabilities).Error; err != nil {
		return nil, err
	}
	return availabilities, nil
}
