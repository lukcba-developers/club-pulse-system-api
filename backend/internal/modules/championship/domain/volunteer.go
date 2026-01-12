package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// VolunteerRole define el rol de un voluntario
type VolunteerRole string

const (
	VolunteerRoleBuffet    VolunteerRole = "BUFFET"    // Atención de buffet
	VolunteerRoleSecurity  VolunteerRole = "SECURITY"  // Seguridad
	VolunteerRoleTransport VolunteerRole = "TRANSPORT" // Transporte
	VolunteerRoleFirstAid  VolunteerRole = "FIRST_AID" // Primeros auxilios
	VolunteerRoleCoach     VolunteerRole = "COACH"     // Asistente técnico
)

// VolunteerAssignment representa la asignación de un voluntario a un partido
type VolunteerAssignment struct {
	ID      uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID  string        `json:"club_id" gorm:"index;not null"`
	MatchID uuid.UUID     `json:"match_id" gorm:"type:uuid;not null;index"`
	UserID  string        `json:"user_id" gorm:"not null;index"`
	Role    VolunteerRole `json:"role" gorm:"not null"`
	Notes   string        `json:"notes"`

	// Metadata
	AssignedBy string    `json:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// VolunteerSummary representa un resumen de voluntarios para un partido
type VolunteerSummary struct {
	MatchID     uuid.UUID             `json:"match_id"`
	TotalSlots  int                   `json:"total_slots"`
	FilledSlots int                   `json:"filled_slots"`
	ByRole      map[VolunteerRole]int `json:"by_role"`
	Volunteers  []VolunteerAssignment `json:"volunteers"`
}

// VolunteerRepository define las operaciones de persistencia para voluntarios
type VolunteerRepository interface {
	Create(ctx context.Context, assignment *VolunteerAssignment) error
	GetByMatchID(ctx context.Context, clubID string, matchID uuid.UUID) ([]VolunteerAssignment, error)
	GetByUserID(ctx context.Context, clubID, userID string) ([]VolunteerAssignment, error)
	GetByRoleAndMatch(ctx context.Context, clubID string, matchID uuid.UUID, role VolunteerRole) ([]VolunteerAssignment, error)
	Update(ctx context.Context, assignment *VolunteerAssignment) error
	Delete(ctx context.Context, clubID string, id uuid.UUID) error
}
