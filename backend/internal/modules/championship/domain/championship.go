package domain

import (
	"time"

	"github.com/google/uuid"
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

	Stages []TournamentStage `json:"stages,omitempty" gorm:"foreignKey:TournamentID"`
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
	Points         int       `json:"points"`
	Played         int       `json:"played"`
	Won            int       `json:"won"`
	Drawn          int       `json:"drawn"`
	Lost           int       `json:"lost"`
	GoalsFor       int       `json:"goals_for"`
	GoalsAgainst   int       `json:"goals_against"`
	GoalDifference int       `json:"goal_difference"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Enriched Fields
	TeamName string `json:"team_name,omitempty" gorm:"-"`
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

	HomeScore *int `json:"home_score,omitempty"`
	AwayScore *int `json:"away_score,omitempty"`

	BookingID *uuid.UUID  `json:"booking_id,omitempty" gorm:"type:uuid;index"` // Link to Booking system
	Status    MatchStatus `json:"status" gorm:"default:'SCHEDULED'"`
	Date      time.Time   `json:"date"`

	// Enriched Fields (Filled via Joins)
	HomeTeamName string `json:"home_team_name,omitempty" gorm:"-"`
	AwayTeamName string `json:"away_team_name,omitempty" gorm:"-"`
}

type ChampionshipRepository interface {
	CreateTournament(tournament *Tournament) error
	GetTournament(id string) (*Tournament, error)
	ListTournaments(clubID string) ([]Tournament, error)
	CreateStage(stage *TournamentStage) error
	GetStage(id string) (*TournamentStage, error)
	CreateGroup(group *Group) error
	GetGroup(id string) (*Group, error)
	CreateMatch(match *TournamentMatch) error
	GetMatch(id string) (*TournamentMatch, error)
	GetMatchesByGroup(groupID string) ([]TournamentMatch, error)
	UpdateMatchResult(matchID string, homeScore, awayScore int) error
	UpdateMatchScheduling(matchID string, date time.Time, bookingID uuid.UUID) error
	GetStandings(groupID string) ([]Standing, error)
	RegisterTeam(standing *Standing) error
	UpdateStanding(standing *Standing) error
	GetTeamMembers(teamID string) ([]string, error)
}
