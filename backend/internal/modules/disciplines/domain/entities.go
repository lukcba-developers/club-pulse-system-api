package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Discipline struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID      string         `json:"club_id" gorm:"index;not null"`
	Name        string         `json:"name" gorm:"not null;unique;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type TrainingGroup struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID       string         `json:"club_id" gorm:"index;not null"`
	Name         string         `json:"name" gorm:"not null;size:100"`
	DisciplineID uuid.UUID      `json:"discipline_id" gorm:"type:uuid;not null"`
	Discipline   Discipline     `json:"discipline" gorm:"foreignKey:DisciplineID"`
	Category     string         `json:"category" gorm:"not null;size:20"` // e.g. "2012"
	CategoryYear int            `json:"category_year"`                    // Normalized year (e.g. 2010)
	CoachID      string         `json:"coach_id"`                         // User ID of the coach
	Schedule     string         `json:"schedule"`                         // e.g. "Mon/Wed 18:00"
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type DisciplineRepository interface {
	CreateDiscipline(ctx context.Context, discipline *Discipline) error
	ListDisciplines(ctx context.Context, clubID string) ([]Discipline, error)
	GetDisciplineByID(ctx context.Context, clubID string, id uuid.UUID) (*Discipline, error)

	CreateGroup(ctx context.Context, group *TrainingGroup) error
	ListGroups(ctx context.Context, clubID string, filter map[string]interface{}) ([]TrainingGroup, error)
	GetGroupByID(ctx context.Context, clubID string, id uuid.UUID) (*TrainingGroup, error)
}

type TournamentRepository interface {
	CreateTournament(ctx context.Context, tournament *Tournament) error
	GetTournamentByID(ctx context.Context, clubID string, id uuid.UUID) (*Tournament, error)
	ListTournaments(ctx context.Context, clubID string) ([]Tournament, error)
	UpdateTournament(ctx context.Context, tournament *Tournament) error

	CreateTeam(ctx context.Context, team *Team) error
	GetTeamByID(ctx context.Context, clubID string, id uuid.UUID) (*Team, error)
	ListTeams(ctx context.Context, clubID string, tournamentID uuid.UUID) ([]Team, error)

	CreateMatch(ctx context.Context, match *Match) error
	UpdateMatch(ctx context.Context, match *Match) error
	GetMatchByID(ctx context.Context, clubID string, id uuid.UUID) (*Match, error)
	ListMatches(ctx context.Context, clubID string, tournamentID uuid.UUID) ([]Match, error)

	GetStandings(ctx context.Context, clubID string, tournamentID uuid.UUID) ([]Standing, error)
}
