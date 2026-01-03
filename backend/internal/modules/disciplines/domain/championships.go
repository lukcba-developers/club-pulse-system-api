package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TournamentStatus string
type MatchStatus string

const (
	TournamentStatusOpen      TournamentStatus = "OPEN"
	TournamentStatusActive    TournamentStatus = "ACTIVE"
	TournamentStatusCompleted TournamentStatus = "COMPLETED"

	MatchStatusScheduled MatchStatus = "SCHEDULED"
	MatchStatusPlayed    MatchStatus = "PLAYED"
	MatchStatusCancelled MatchStatus = "CANCELLED"
)

type Tournament struct {
	ID           uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID       string           `json:"club_id" gorm:"index;not null"`
	Name         string           `json:"name" gorm:"not null"`
	DisciplineID uuid.UUID        `json:"discipline_id" gorm:"type:uuid;not null"`
	StartDate    time.Time        `json:"start_date"`
	EndDate      time.Time        `json:"end_date"`
	Status       TournamentStatus `json:"status" gorm:"default:'OPEN'"`
	Format       string           `json:"format"` // "LEAGUE", "BRACKET"
	Teams        []Team           `json:"teams,omitempty" gorm:"foreignKey:TournamentID"`
	Matches      []Match          `json:"matches,omitempty" gorm:"foreignKey:TournamentID"`
	CreatedAt    time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index"`
}

type Team struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID       string         `json:"club_id" gorm:"index;not null"`
	Name         string         `json:"name" gorm:"not null"`
	TournamentID uuid.UUID      `json:"tournament_id" gorm:"type:uuid;not null"`
	CaptainID    *string        `json:"captain_id" gorm:"type:varchar(100)"` // Optional link to User
	Members      []string       `json:"members,omitempty" gorm:"type:jsonb"` // List of UserIDs or strings
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Match struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID       string         `json:"club_id" gorm:"index;not null"`
	TournamentID uuid.UUID      `json:"tournament_id" gorm:"type:uuid;not null"`
	HomeTeamID   uuid.UUID      `json:"home_team_id" gorm:"type:uuid;not null"`
	AwayTeamID   uuid.UUID      `json:"away_team_id" gorm:"type:uuid;not null"`
	ScoreHome    int            `json:"score_home" gorm:"default:0"`
	ScoreAway    int            `json:"score_away" gorm:"default:0"`
	StartTime    time.Time      `json:"start_time"`
	Status       MatchStatus    `json:"status" gorm:"default:'SCHEDULED'"`
	Round        string         `json:"round"` // "1", "QuarterFinal", etc.
	Location     string         `json:"location"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Standing struct {
	TournamentID uuid.UUID `json:"tournament_id"`
	TeamID       uuid.UUID `json:"team_id"`
	TeamName     string    `json:"team_name"`
	Played       int       `json:"played"`
	Won          int       `json:"won"`
	Drawn        int       `json:"drawn"`
	Lost         int       `json:"lost"`
	Points       int       `json:"points"`
	GoalsFor     int       `json:"goals_for"`
	GoalsAgainst int       `json:"goals_against"`
}
