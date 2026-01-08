package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
)

// VolunteerService maneja la lógica de negocio para voluntarios
type VolunteerService struct {
	volunteerRepo domain.VolunteerRepository
}

// NewVolunteerService crea una nueva instancia del servicio
func NewVolunteerService(volunteerRepo domain.VolunteerRepository) *VolunteerService {
	return &VolunteerService{
		volunteerRepo: volunteerRepo,
	}
}

// AssignVolunteer asigna un voluntario a un partido
func (s *VolunteerService) AssignVolunteer(ctx context.Context, clubID string, matchID uuid.UUID, userID string, role domain.VolunteerRole, assignedBy string, notes string) error {
	assignment := &domain.VolunteerAssignment{
		ClubID:     clubID,
		MatchID:    matchID,
		UserID:     userID,
		Role:       role,
		Notes:      notes,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
	}

	return s.volunteerRepo.Create(assignment)
}

// GetMatchVolunteers obtiene todos los voluntarios de un partido
func (s *VolunteerService) GetMatchVolunteers(ctx context.Context, clubID string, matchID uuid.UUID) ([]domain.VolunteerAssignment, error) {
	return s.volunteerRepo.GetByMatchID(clubID, matchID)
}

// GetVolunteerSummary obtiene un resumen de voluntarios para un partido
func (s *VolunteerService) GetVolunteerSummary(ctx context.Context, clubID string, matchID uuid.UUID) (*domain.VolunteerSummary, error) {
	volunteers, err := s.volunteerRepo.GetByMatchID(clubID, matchID)
	if err != nil {
		return nil, err
	}

	summary := &domain.VolunteerSummary{
		MatchID:     matchID,
		TotalSlots:  10, // Configurable
		FilledSlots: len(volunteers),
		ByRole:      make(map[domain.VolunteerRole]int),
		Volunteers:  volunteers,
	}

	for _, v := range volunteers {
		summary.ByRole[v.Role]++
	}

	return summary, nil
}

// RemoveVolunteer remueve un voluntario de un partido
func (s *VolunteerService) RemoveVolunteer(ctx context.Context, clubID string, assignmentID uuid.UUID) error {
	return s.volunteerRepo.Delete(clubID, assignmentID)
}

// GetUserAssignments obtiene todas las asignaciones de un usuario
func (s *VolunteerService) GetUserAssignments(ctx context.Context, clubID, userID string) ([]domain.VolunteerAssignment, error) {
	return s.volunteerRepo.GetByUserID(clubID, userID)
}

// ValidateAssignment verifica que un rol específico tenga espacio
func (s *VolunteerService) ValidateAssignment(ctx context.Context, clubID string, matchID uuid.UUID, role domain.VolunteerRole, maxPerRole int) error {
	existing, err := s.volunteerRepo.GetByRoleAndMatch(clubID, matchID, role)
	if err != nil {
		return err
	}

	if len(existing) >= maxPerRole {
		return fmt.Errorf("el rol %s ya tiene el máximo de voluntarios asignados", role)
	}

	return nil
}
