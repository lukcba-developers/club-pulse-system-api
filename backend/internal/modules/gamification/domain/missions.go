package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// MissionType defines whether a mission is daily or weekly.
type MissionType string

const (
	MissionTypeDaily  MissionType = "DAILY"
	MissionTypeWeekly MissionType = "WEEKLY"
)

// MissionStatus tracks the state of a user's mission progress.
type MissionStatus string

const (
	MissionStatusActive    MissionStatus = "ACTIVE"
	MissionStatusCompleted MissionStatus = "COMPLETED"
	MissionStatusClaimed   MissionStatus = "CLAIMED"
	MissionStatusExpired   MissionStatus = "EXPIRED"
)

// Mission represents a challenge that users can complete for rewards.
type Mission struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID      string      `json:"club_id" gorm:"index;not null"`
	Code        string      `json:"code" gorm:"not null"` // "EARLY_BIRD", "EXPLORER"
	Type        MissionType `json:"type" gorm:"not null"` // DAILY, WEEKLY
	Name        string      `json:"name" gorm:"not null"`
	Description string      `json:"description"`
	XPReward    int         `json:"xp_reward" gorm:"default:50"`
	BadgeID     *uuid.UUID  `json:"badge_id,omitempty" gorm:"type:uuid"` // Optional badge reward
	TargetValue int         `json:"target_value" gorm:"default:1"`       // e.g., "3 bookings" -> 3
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (Mission) TableName() string {
	return "missions"
}

// UserMission tracks a user's progress on a specific mission.
type UserMission struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      string        `json:"user_id" gorm:"not null;index:idx_user_missions_user"`
	MissionID   uuid.UUID     `json:"mission_id" gorm:"type:uuid;not null"`
	Status      MissionStatus `json:"status" gorm:"default:'ACTIVE'"`
	Progress    int           `json:"progress" gorm:"default:0"`
	AssignedAt  time.Time     `json:"assigned_at"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	ClaimedAt   *time.Time    `json:"claimed_at,omitempty"`
	ExpiresAt   time.Time     `json:"expires_at"` // Daily missions expire at midnight, weekly on Sunday
}

func (UserMission) TableName() string {
	return "user_missions"
}

// MissionRepository defines the interface for mission persistence.
type MissionRepository interface {
	// Mission CRUD
	Create(ctx context.Context, mission *Mission) error
	GetByID(ctx context.Context, clubID string, id uuid.UUID) (*Mission, error)
	GetByCode(ctx context.Context, clubID, code string) (*Mission, error)
	ListActive(ctx context.Context, clubID string, missionType MissionType) ([]Mission, error)

	// UserMission operations
	AssignMission(ctx context.Context, userMission *UserMission) error
	GetUserMissions(ctx context.Context, userID string, missionType MissionType) ([]UserMission, error)
	GetActiveUserMissions(ctx context.Context, userID string) ([]UserMission, error)
	UpdateProgress(ctx context.Context, userMissionID uuid.UUID, progress int) error
	CompleteMission(ctx context.Context, userMissionID uuid.UUID) error
	ClaimReward(ctx context.Context, userMissionID uuid.UUID) error
	ExpireOldMissions(ctx context.Context, before time.Time) error
}

// PredefinedMissions contains the default missions for the system.
var PredefinedMissions = []Mission{
	// Daily Missions
	{Code: "EARLY_BIRD", Type: MissionTypeDaily, Name: "Madrugador", Description: "Haz una reserva antes de las 10am", XPReward: 50, TargetValue: 1},
	{Code: "TEAM_PLAYER", Type: MissionTypeDaily, Name: "Compañero", Description: "Completa una reserva con 2+ jugadores", XPReward: 75, TargetValue: 1},
	{Code: "EXPLORER", Type: MissionTypeDaily, Name: "Explorador", Description: "Reserva en una instalación nueva", XPReward: 100, TargetValue: 1},
	{Code: "DOUBLE_PLAY", Type: MissionTypeDaily, Name: "Doble Jornada", Description: "Completa 2 reservas en un día", XPReward: 100, TargetValue: 2},

	// Weekly Missions
	{Code: "TRIATHLETE", Type: MissionTypeWeekly, Name: "Triatleta", Description: "3 reservas en 3 deportes diferentes", XPReward: 500, TargetValue: 3},
	{Code: "IRON_PLAYER", Type: MissionTypeWeekly, Name: "Iron Player", Description: "5 partidos en la semana", XPReward: 400, TargetValue: 5},
	{Code: "SOCIAL_BUTTERFLY", Type: MissionTypeWeekly, Name: "Social", Description: "Invita a 1 nuevo miembro", XPReward: 300, TargetValue: 1},
	{Code: "CONSISTENCY", Type: MissionTypeWeekly, Name: "Constancia", Description: "Actividad 5 días de la semana", XPReward: 350, TargetValue: 5},
}
