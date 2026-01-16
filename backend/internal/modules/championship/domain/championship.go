package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TournamentStatus string

const (
	TournamentDraft     TournamentStatus = "DRAFT"
	TournamentActive    TournamentStatus = "ACTIVE"
	TournamentCompleted TournamentStatus = "COMPLETED"
)

type Tournament struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID      uuid.UUID        `json:"club_id" gorm:"type:uuid;not null;index"` // Organizing Club
	Name        string           `json:"name" gorm:"not null"`
	Description string           `json:"description,omitempty"`
	Sport       string           `json:"sport" gorm:"not null"` // e.g., "FUTBOL", "PADEL"
	Category    string           `json:"category,omitempty"`    // e.g., "Libre", "Veteranos"
	Status      TournamentStatus `json:"status" gorm:"default:'DRAFT'"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	LogoURL     string           `json:"logo_url,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`

	Settings datatypes.JSON `json:"settings" gorm:"default:'{}'"` // Dynamic configuration (points per win, durations, etc)

	Stages []TournamentStage `json:"stages,omitempty" gorm:"foreignKey:TournamentID"`
}

func (Tournament) TableName() string {
	return "championships"
}

type StageType string

const (
	StageGroup    StageType = "GROUP"
	StageKnockout StageType = "KNOCKOUT"
)

type TournamentStage struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TournamentID uuid.UUID   `json:"tournament_id" gorm:"type:uuid;not null;index"`
	Order        int         `json:"order" gorm:"not null"` // 1, 2, 3...
	Name         string      `json:"name" gorm:"not null"`  // "Fase de Grupos", "Playoffs"
	Type         StageType   `json:"type" gorm:"not null"`
	Status       StageStatus `json:"status" gorm:"default:'PENDING'"`
	Groups       []Group     `json:"groups,omitempty" gorm:"foreignKey:StageID"`
}

type StageStatus string

const (
	StagePending   StageStatus = "PENDING"
	StageActive    StageStatus = "ACTIVE"
	StageCompleted StageStatus = "COMPLETED"
)

type Group struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	StageID uuid.UUID `json:"stage_id" gorm:"type:uuid;not null;index"`
	Name    string    `json:"name" gorm:"not null"` // "Grupo A", "Grupo B"

	Standings []Standing `json:"standings,omitempty" gorm:"foreignKey:GroupID"`
}

type Team struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"not null"`
	LogoURL   string    `json:"logo_url,omitempty"`
	Contact   string    `json:"contact,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Standing struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	GroupID        uuid.UUID `json:"group_id" gorm:"type:uuid;not null;index"`
	TeamID         uuid.UUID `json:"team_id" gorm:"type:uuid;not null;index"`
	Points         float64   `json:"points"`
	Played         int       `json:"played"`
	Won            int       `json:"won"`
	Drawn          int       `json:"drawn"`
	Lost           int       `json:"lost"`
	GoalsFor       float64   `json:"goals_for"`
	GoalsAgainst   float64   `json:"goals_against"`
	GoalDifference float64   `json:"goal_difference"`
	Position       int       `json:"position"` // Calculated ranking position
	UpdatedAt      time.Time `json:"updated_at"`

	// Enriched Fields
	TeamName string `json:"team_name,omitempty" gorm:"->"`
}

type MatchStatus string

const (
	MatchScheduled MatchStatus = "SCHEDULED"
	MatchCompleted MatchStatus = "COMPLETED"
	MatchCancelled MatchStatus = "CANCELLED"
)

type TournamentMatch struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TournamentID uuid.UUID  `json:"tournament_id" gorm:"type:uuid;not null;index"`
	StageID      uuid.UUID  `json:"stage_id" gorm:"type:uuid;not null;index"`
	GroupID      *uuid.UUID `json:"group_id,omitempty" gorm:"type:uuid;index"` // Nullable if knockout

	HomeTeamID uuid.UUID `json:"home_team_id" gorm:"type:uuid;not null;index"`
	AwayTeamID uuid.UUID `json:"away_team_id" gorm:"type:uuid;not null;index"`

	HomeScore *float64 `json:"home_score,omitempty"`
	AwayScore *float64 `json:"away_score,omitempty"`

	BookingID *uuid.UUID  `json:"booking_id,omitempty" gorm:"type:uuid;index"` // Link to Booking system
	Status    MatchStatus `json:"status" gorm:"default:'SCHEDULED'"`
	Date      time.Time   `json:"date"`

	// Enriched Fields (Filled via Joins)
	HomeTeamName string `json:"home_team_name,omitempty" gorm:"-"`
	AwayTeamName string `json:"away_team_name,omitempty" gorm:"-"`
}

func (TournamentMatch) TableName() string {
	return "tournament_matches"
}

type ChampionshipRepository interface {
	CreateTournament(ctx context.Context, tournament *Tournament) error
	GetTournament(ctx context.Context, clubID, id string) (*Tournament, error)
	ListTournaments(ctx context.Context, clubID string) ([]Tournament, error)
	CreateStage(ctx context.Context, stage *TournamentStage) error
	GetStage(ctx context.Context, clubID, id string) (*TournamentStage, error)
	CreateGroup(ctx context.Context, group *Group) error
	GetGroup(ctx context.Context, clubID, id string) (*Group, error)
	CreateMatch(ctx context.Context, clubID string, match *TournamentMatch) error
	CreateMatchesBatch(ctx context.Context, clubID string, matches []TournamentMatch) error // Atomic batch creation
	GetMatch(ctx context.Context, clubID, id string) (*TournamentMatch, error)
	GetMatchesByGroup(ctx context.Context, clubID, groupID string) ([]TournamentMatch, error)
	UpdateMatchResult(ctx context.Context, clubID, matchID string, homeScore, awayScore float64) error
	UpdateMatchScheduling(ctx context.Context, clubID, matchID string, date time.Time, bookingID uuid.UUID) error
	GetStandings(ctx context.Context, clubID, groupID string) ([]Standing, error)
	RegisterTeam(ctx context.Context, clubID string, standing *Standing) error
	UpdateStanding(ctx context.Context, standing *Standing) error
	UpdateStandingsBatch(ctx context.Context, clubID string, standings []Standing) error
	GetTeamMembers(ctx context.Context, teamID string) ([]string, error)
	CreateTeam(ctx context.Context, team *Team) error
	AddMember(ctx context.Context, teamID, userID string) error
	GetMatchesByUserID(ctx context.Context, clubID, userID string) ([]TournamentMatch, error)
	GetUpcomingMatches(ctx context.Context, clubID string, from, to time.Time) ([]TournamentMatch, error)
}
