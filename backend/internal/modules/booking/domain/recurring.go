package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecurrenceType string

const (
	RecurrenceTypeClass       RecurrenceType = "CLASS"
	RecurrenceTypeMaintenance RecurrenceType = "MAINTENANCE"
	RecurrenceTypeFixed       RecurrenceType = "FIXED"
)

type RecurringRule struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID     string         `json:"club_id" gorm:"index;not null"`
	FacilityID uuid.UUID      `json:"facility_id" gorm:"type:uuid;not null;index"`
	Type       RecurrenceType `json:"type" gorm:"type:varchar(20);not null"`

	// DayOfWeek: 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	DayOfWeek int       `json:"day_of_week" gorm:"not null"`
	StartTime time.Time `json:"start_time" gorm:"type:time;not null"`
	EndTime   time.Time `json:"end_time" gorm:"type:time;not null"`

	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`

	// Optional: Who owns this recurrence? (e.g., Coach for Class, User for Fixed)
	OwnerID *uuid.UUID `json:"owner_id,omitempty" gorm:"type:uuid"`
	GroupID *uuid.UUID `json:"group_id,omitempty" gorm:"type:uuid"` // For Schools/Classes (TrainingGroup)

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type RecurringRepository interface {
	Create(ctx context.Context, rule *RecurringRule) error
	GetByFacility(ctx context.Context, clubID string, facilityID uuid.UUID) ([]RecurringRule, error)
	GetAllActive(ctx context.Context, clubID string) ([]RecurringRule, error)
	Delete(ctx context.Context, clubID string, id uuid.UUID) error
}
