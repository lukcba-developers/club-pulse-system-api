package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

type Booking struct {
	ID         uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	ClubID     string        `json:"club_id" gorm:"index;not null"`
	UserID     uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	FacilityID uuid.UUID     `json:"facility_id" gorm:"type:uuid;not null"`
	StartTime  time.Time     `json:"start_time" gorm:"not null"`
	EndTime    time.Time     `json:"end_time" gorm:"not null"`
	Status     BookingStatus `json:"status" gorm:"type:varchar(20);default:'CONFIRMED'"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(clubID string, id uuid.UUID) (*Booking, error)
	List(clubID string, filter map[string]interface{}) ([]Booking, error)
	Update(booking *Booking) error
	HasTimeConflict(clubID string, facilityID uuid.UUID, start, end time.Time) (bool, error)
	ListByFacilityAndDate(clubID string, facilityID uuid.UUID, date time.Time) ([]Booking, error)
	ListAll(clubID string, filter map[string]interface{}, from, to *time.Time) ([]Booking, error)
}
