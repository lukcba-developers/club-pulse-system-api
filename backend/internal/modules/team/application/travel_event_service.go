package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	"github.com/shopspring/decimal"
)

// TravelEventService maneja la lógica de negocio para eventos de viaje
type TravelEventService struct {
	eventRepo domain.TravelEventRepository
}

// NewTravelEventService crea una nueva instancia del servicio
func NewTravelEventService(eventRepo domain.TravelEventRepository) *TravelEventService {
	return &TravelEventService{
		eventRepo: eventRepo,
	}
}

// CreateEvent crea un nuevo evento de viaje
func (s *TravelEventService) CreateEvent(ctx context.Context, clubID string, teamID uuid.UUID, createdBy string, event *domain.TravelEvent) error {
	event.ClubID = clubID
	event.TeamID = teamID
	event.CreatedBy = createdBy

	return s.eventRepo.Create(ctx, event)
}

// GetEvent obtiene un evento por su ID
func (s *TravelEventService) GetEvent(ctx context.Context, clubID string, eventID uuid.UUID) (*domain.TravelEvent, error) {
	return s.eventRepo.GetByID(ctx, clubID, eventID)
}

// GetTeamEvents obtiene todos los eventos de un equipo
func (s *TravelEventService) GetTeamEvents(ctx context.Context, clubID string, teamID uuid.UUID) ([]domain.TravelEvent, error) {
	return s.eventRepo.GetByTeamID(ctx, clubID, teamID)
}

// GetUpcomingEvents obtiene eventos futuros de un equipo
func (s *TravelEventService) GetUpcomingEvents(ctx context.Context, clubID string, teamID uuid.UUID) ([]domain.TravelEvent, error) {
	return s.eventRepo.GetUpcoming(ctx, clubID, teamID)
}

// RespondToEvent registra la respuesta de un usuario a un evento
func (s *TravelEventService) RespondToEvent(ctx context.Context, eventID uuid.UUID, userID string, status domain.RSVPStatus, notes string) error {
	// Verificar si ya existe una respuesta
	existingRSVP, err := s.eventRepo.GetRSVPByUserAndEvent(ctx, eventID, userID)

	now := time.Now()

	if err == nil && existingRSVP != nil {
		// Actualizar respuesta existente
		existingRSVP.Status = status
		existingRSVP.Notes = notes
		existingRSVP.RespondedAt = &now
		return s.eventRepo.UpdateRSVP(ctx, existingRSVP)
	}

	// Crear nueva respuesta
	rsvp := &domain.EventRSVP{
		EventID:     eventID,
		UserID:      userID,
		Status:      status,
		Notes:       notes,
		RespondedAt: &now,
	}

	return s.eventRepo.CreateRSVP(ctx, rsvp)
}

// GetEventSummary obtiene un resumen del evento con estadísticas
func (s *TravelEventService) GetEventSummary(ctx context.Context, clubID string, eventID uuid.UUID) (*domain.EventSummary, error) {
	event, err := s.eventRepo.GetByID(ctx, clubID, eventID)
	if err != nil {
		return nil, err
	}

	rsvps, err := s.eventRepo.GetRSVPsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	summary := &domain.EventSummary{
		Event:          event,
		TotalInvited:   len(rsvps),
		ConfirmedUsers: []string{},
	}

	for _, rsvp := range rsvps {
		switch rsvp.Status {
		case domain.RSVPStatusConfirmed:
			summary.TotalConfirmed++
			summary.ConfirmedUsers = append(summary.ConfirmedUsers, rsvp.UserID)
		case domain.RSVPStatusDeclined:
			summary.TotalDeclined++
		case domain.RSVPStatusPending:
			summary.TotalPending++
		}
	}

	// Calcular costo per cápita
	summary.CostPerPerson = event.CalculateCostPerPerson(summary.TotalConfirmed)

	// Actualizar el evento con el costo calculado
	event.CostPerPerson = summary.CostPerPerson
	_ = s.eventRepo.Update(ctx, event)

	return summary, nil
}

// UpdateEventCost actualiza el costo real de un evento
func (s *TravelEventService) UpdateEventCost(ctx context.Context, clubID string, eventID uuid.UUID, actualCost decimal.Decimal) error {
	event, err := s.eventRepo.GetByID(ctx, clubID, eventID)
	if err != nil {
		return err
	}

	event.ActualCost = actualCost

	// Recalcular costo per cápita
	rsvps, err := s.eventRepo.GetRSVPsByEventID(ctx, eventID)
	if err != nil {
		return err
	}

	confirmedCount := 0
	for _, rsvp := range rsvps {
		if rsvp.Status == domain.RSVPStatusConfirmed {
			confirmedCount++
		}
	}

	event.CostPerPerson = event.CalculateCostPerPerson(confirmedCount)

	return s.eventRepo.Update(ctx, event)
}

// CancelEvent cancela un evento
func (s *TravelEventService) CancelEvent(ctx context.Context, clubID string, eventID uuid.UUID) error {
	// TODO: Enviar notificaciones a todos los confirmados
	return s.eventRepo.Delete(ctx, clubID, eventID)
}

// GetUserRSVP obtiene la respuesta de un usuario para un evento
func (s *TravelEventService) GetUserRSVP(ctx context.Context, eventID uuid.UUID, userID string) (*domain.EventRSVP, error) {
	return s.eventRepo.GetRSVPByUserAndEvent(ctx, eventID, userID)
}

// ValidateEventCapacity verifica si el evento puede aceptar más confirmaciones
func (s *TravelEventService) ValidateEventCapacity(ctx context.Context, clubID string, eventID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, clubID, eventID)
	if err != nil {
		return err
	}

	if !event.IsOpen() {
		return fmt.Errorf("el evento ya pasó")
	}

	rsvps, err := s.eventRepo.GetRSVPsByEventID(ctx, eventID)
	if err != nil {
		return err
	}

	confirmedCount := 0
	for _, rsvp := range rsvps {
		if rsvp.Status == domain.RSVPStatusConfirmed {
			confirmedCount++
		}
	}

	if event.IsFull(confirmedCount) {
		return fmt.Errorf("el evento alcanzó el máximo de participantes")
	}

	return nil
}
