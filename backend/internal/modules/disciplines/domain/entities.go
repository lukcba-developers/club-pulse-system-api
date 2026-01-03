package domain

import (
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
	CoachID      string         `json:"coach_id"`                         // User ID of the coach
	Schedule     string         `json:"schedule"`                         // e.g. "Mon/Wed 18:00"
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type DisciplineRepository interface {
	CreateDiscipline(discipline *Discipline) error
	ListDisciplines(clubID string) ([]Discipline, error)
	GetDisciplineByID(clubID string, id uuid.UUID) (*Discipline, error)

	CreateGroup(group *TrainingGroup) error
	ListGroups(clubID string, filter map[string]interface{}) ([]TrainingGroup, error)
	GetGroupByID(clubID string, id uuid.UUID) (*TrainingGroup, error)
}

type TournamentRepository interface {
	CreateTournament(tournament *Tournament) error
	GetTournamentByID(clubID string, id uuid.UUID) (*Tournament, error)
	ListTournaments(clubID string) ([]Tournament, error)
	UpdateTournament(tournament *Tournament) error

	CreateTeam(team *Team) error
	GetTeamByID(clubID string, id uuid.UUID) (*Team, error)
	ListTeams(clubID string, tournamentID uuid.UUID) ([]Team, error)

	CreateMatch(match *Match) error
	UpdateMatch(match *Match) error
	GetMatchByID(clubID string, id uuid.UUID) (*Match, error)
	ListMatches(clubID string, tournamentID uuid.UUID) ([]Match, error)

	GetStandings(clubID string, tournamentID uuid.UUID) ([]Standing, error)
}
