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
	ID        uuid.UUID          `json:"id"`
	Date      time.Time          `json:"date"`
	Group     string             `json:"group"` // e.g. "2012", "Pre-Novena"
	CoachID   string             `json:"coach_id"`
	Records   []AttendanceRecord `json:"records,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// AttendanceRecord represents the status of a single user in an AttendanceList.
type AttendanceRecord struct {
	ID               uuid.UUID        `json:"id"`
	AttendanceListID uuid.UUID        `json:"attendance_list_id"`
	UserID           string           `json:"user_id"`
	Status           AttendanceStatus `json:"status"`
	Notes            string           `json:"notes,omitempty"`
	// Populated for response convenience if needed, essentially a DTO field, or we fetch separately.
	User *userDomain.User `json:"user,omitempty" gorm:"-"`
}

type AttendanceRepository interface {
	CreateList(list *AttendanceList) error
	GetListByID(id uuid.UUID) (*AttendanceList, error)
	GetListByGroupAndDate(group string, date time.Time) (*AttendanceList, error)
	UpdateRecord(record *AttendanceRecord) error
	// UpsertRecord updates or creates a record if it doesn't exist within a list
	UpsertRecord(record *AttendanceRecord) error
}
