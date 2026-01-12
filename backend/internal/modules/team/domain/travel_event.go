package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// EventType define el tipo de evento
type EventType string

const (
	EventTypeTravel     EventType = "TRAVEL"     // Viaje
	EventTypeMatch      EventType = "MATCH"      // Partido
	EventTypeTournament EventType = "TOURNAMENT" // Torneo
	EventTypeTraining   EventType = "TRAINING"   // Entrenamiento especial
)

// RSVPStatus define el estado de confirmación de asistencia
type RSVPStatus string

const (
	RSVPStatusPending   RSVPStatus = "PENDING"   // Sin confirmar
	RSVPStatusConfirmed RSVPStatus = "CONFIRMED" // Confirmado
	RSVPStatusDeclined  RSVPStatus = "DECLINED"  // Rechazado
)

// TravelEvent representa un evento de viaje del equipo
type TravelEvent struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID      string    `json:"club_id" gorm:"index;not null"`
	TeamID      uuid.UUID `json:"team_id" gorm:"type:uuid;not null;index"`
	Type        EventType `json:"type" gorm:"not null;default:'TRAVEL'"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`

	// Detalles del viaje
	Destination   string     `json:"destination" gorm:"not null"`
	DepartureDate time.Time  `json:"departure_date" gorm:"not null"`
	ReturnDate    *time.Time `json:"return_date,omitempty"`
	MeetingPoint  string     `json:"meeting_point"`
	MeetingTime   time.Time  `json:"meeting_time" gorm:"not null"`

	// Costos
	EstimatedCost decimal.Decimal `json:"estimated_cost" gorm:"type:decimal(10,2);default:0"`
	ActualCost    decimal.Decimal `json:"actual_cost" gorm:"type:decimal(10,2);default:0"`
	CostPerPerson decimal.Decimal `json:"cost_per_person" gorm:"type:decimal(10,2);default:0"` // Calculado

	// Metadata
	MaxParticipants *int      `json:"max_participants,omitempty"`
	CreatedBy       string    `json:"created_by" gorm:"not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones (cargadas bajo demanda)
	RSVPs []EventRSVP `json:"rsvps,omitempty" gorm:"foreignKey:EventID"`
}

// EventRSVP representa la confirmación de asistencia de un usuario a un evento
type EventRSVP struct {
	ID      uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID uuid.UUID  `json:"event_id" gorm:"type:uuid;not null;index"`
	UserID  string     `json:"user_id" gorm:"not null;index"`
	Status  RSVPStatus `json:"status" gorm:"not null;default:'PENDING'"`
	Notes   string     `json:"notes"`

	// Metadata
	RespondedAt *time.Time `json:"responded_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// EventSummary representa un resumen de un evento con estadísticas
type EventSummary struct {
	Event          *TravelEvent    `json:"event"`
	TotalInvited   int             `json:"total_invited"`
	TotalConfirmed int             `json:"total_confirmed"`
	TotalDeclined  int             `json:"total_declined"`
	TotalPending   int             `json:"total_pending"`
	CostPerPerson  decimal.Decimal `json:"cost_per_person"`
	ConfirmedUsers []string        `json:"confirmed_users,omitempty"`
}

// CalculateCostPerPerson calcula el costo por persona basado en confirmados
func (e *TravelEvent) CalculateCostPerPerson(confirmedCount int) decimal.Decimal {
	if confirmedCount == 0 {
		return decimal.Zero
	}

	cost := e.ActualCost
	if cost.IsZero() {
		cost = e.EstimatedCost
	}

	return cost.Div(decimal.NewFromInt(int64(confirmedCount)))
}

// IsOpen verifica si el evento aún acepta confirmaciones
func (e *TravelEvent) IsOpen() bool {
	// El evento está abierto si la fecha de salida es futura
	return time.Now().Before(e.DepartureDate)
}

// IsFull verifica si el evento alcanzó el máximo de participantes
func (e *TravelEvent) IsFull(currentCount int) bool {
	if e.MaxParticipants == nil {
		return false
	}
	return currentCount >= *e.MaxParticipants
}

// TravelEventRepository define las operaciones de persistencia para eventos
type TravelEventRepository interface {
	Create(ctx context.Context, event *TravelEvent) error
	GetByID(ctx context.Context, clubID string, id uuid.UUID) (*TravelEvent, error)
	GetByTeamID(ctx context.Context, clubID string, teamID uuid.UUID) ([]TravelEvent, error)
	GetUpcoming(ctx context.Context, clubID string, teamID uuid.UUID) ([]TravelEvent, error)
	Update(ctx context.Context, event *TravelEvent) error
	Delete(ctx context.Context, clubID string, id uuid.UUID) error

	// RSVP operations
	CreateRSVP(ctx context.Context, rsvp *EventRSVP) error
	GetRSVPsByEventID(ctx context.Context, eventID uuid.UUID) ([]EventRSVP, error)
	GetRSVPByUserAndEvent(ctx context.Context, eventID uuid.UUID, userID string) (*EventRSVP, error)
	UpdateRSVP(ctx context.Context, rsvp *EventRSVP) error
	DeleteRSVP(ctx context.Context, id uuid.UUID) error
}
