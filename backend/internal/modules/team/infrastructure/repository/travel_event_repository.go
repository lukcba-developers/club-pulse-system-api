package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	"gorm.io/gorm"
)

// PostgresTravelEventRepository implementa el repositorio de eventos usando PostgreSQL
type PostgresTravelEventRepository struct {
	db *gorm.DB
}

// NewPostgresTravelEventRepository crea una nueva instancia del repositorio
func NewPostgresTravelEventRepository(db *gorm.DB) *PostgresTravelEventRepository {
	return &PostgresTravelEventRepository{db: db}
}

// Create crea un nuevo evento
func (r *PostgresTravelEventRepository) Create(event *domain.TravelEvent) error {
	return r.db.Create(event).Error
}

// GetByID obtiene un evento por su ID
func (r *PostgresTravelEventRepository) GetByID(clubID string, id uuid.UUID) (*domain.TravelEvent, error) {
	var event domain.TravelEvent
	err := r.db.Where("club_id = ? AND id = ?", clubID, id).
		Preload("RSVPs").
		First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetByTeamID obtiene todos los eventos de un equipo
func (r *PostgresTravelEventRepository) GetByTeamID(clubID string, teamID uuid.UUID) ([]domain.TravelEvent, error) {
	var events []domain.TravelEvent
	err := r.db.Where("club_id = ? AND team_id = ?", clubID, teamID).
		Order("departure_date DESC").
		Preload("RSVPs").
		Find(&events).Error
	return events, err
}

// GetUpcoming obtiene eventos futuros de un equipo
func (r *PostgresTravelEventRepository) GetUpcoming(clubID string, teamID uuid.UUID) ([]domain.TravelEvent, error) {
	var events []domain.TravelEvent
	err := r.db.Where("club_id = ? AND team_id = ? AND departure_date > ?", clubID, teamID, time.Now()).
		Order("departure_date ASC").
		Preload("RSVPs").
		Find(&events).Error
	return events, err
}

// Update actualiza un evento existente
func (r *PostgresTravelEventRepository) Update(event *domain.TravelEvent) error {
	return r.db.Save(event).Error
}

// Delete elimina un evento
func (r *PostgresTravelEventRepository) Delete(clubID string, id uuid.UUID) error {
	return r.db.Where("club_id = ? AND id = ?", clubID, id).
		Delete(&domain.TravelEvent{}).Error
}

// CreateRSVP crea una nueva confirmación de asistencia
func (r *PostgresTravelEventRepository) CreateRSVP(rsvp *domain.EventRSVP) error {
	return r.db.Create(rsvp).Error
}

// GetRSVPsByEventID obtiene todas las confirmaciones de un evento
func (r *PostgresTravelEventRepository) GetRSVPsByEventID(eventID uuid.UUID) ([]domain.EventRSVP, error) {
	var rsvps []domain.EventRSVP
	err := r.db.Where("event_id = ?", eventID).
		Order("created_at ASC").
		Find(&rsvps).Error
	return rsvps, err
}

// GetRSVPByUserAndEvent obtiene la confirmación de un usuario para un evento específico
func (r *PostgresTravelEventRepository) GetRSVPByUserAndEvent(eventID uuid.UUID, userID string) (*domain.EventRSVP, error) {
	var rsvp domain.EventRSVP
	err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).
		First(&rsvp).Error
	if err != nil {
		return nil, err
	}
	return &rsvp, nil
}

// UpdateRSVP actualiza una confirmación existente
func (r *PostgresTravelEventRepository) UpdateRSVP(rsvp *domain.EventRSVP) error {
	return r.db.Save(rsvp).Error
}

// DeleteRSVP elimina una confirmación
func (r *PostgresTravelEventRepository) DeleteRSVP(id uuid.UUID) error {
	return r.db.Delete(&domain.EventRSVP{}, id).Error
}
