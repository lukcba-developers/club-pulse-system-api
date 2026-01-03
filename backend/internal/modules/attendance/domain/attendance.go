package domain

import (
	"time"

	"github.com/google/uuid"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type AttendanceStatus string

const (
	StatusPresent AttendanceStatus = "PRESENT"
	StatusAbsent  AttendanceStatus = "ABSENT"
	StatusLate    AttendanceStatus = "LATE"
)

// AttendanceList represents a roll call session for a specific group/category on a specific date.
type AttendanceList struct {
	ID              uuid.UUID          `json:"id"`
	ClubID          string             `json:"club_id"`
	Date            time.Time          `json:"date"`
	Group           string             `json:"group"` // Display name, e.g. "Fútbol 2012"
	TrainingGroupID *uuid.UUID         `json:"training_group_id"`
	CoachID         string             `json:"coach_id"`
	Records         []AttendanceRecord `json:"records,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

// AttendanceRecord represents the status of a single user in an AttendanceList.
type AttendanceRecord struct {
	ID               uuid.UUID        `json:"id"`
	AttendanceListID uuid.UUID        `json:"attendance_list_id"`
	UserID           string           `json:"user_id"`
	Status           AttendanceStatus `json:"status"`
	Notes            string           `json:"notes,omitempty"`
	HasDebt          bool             `json:"has_debt"` // Computed field for UI
	// Populated for response convenience
	User *userDomain.User `json:"user,omitempty" gorm:"-"`
}

type AttendanceRepository interface {
	CreateList(list *AttendanceList) error
	GetListByID(clubID string, id uuid.UUID) (*AttendanceList, error)
	GetListByGroupAndDate(clubID string, group string, date time.Time) (*AttendanceList, error)
	GetListByTrainingGroupAndDate(clubID string, groupID uuid.UUID, date time.Time) (*AttendanceList, error)
	UpdateRecord(record *AttendanceRecord) error
	// UpsertRecord updates or creates a record if it doesn't exist within a list
	UpsertRecord(record *AttendanceRecord) error
}
