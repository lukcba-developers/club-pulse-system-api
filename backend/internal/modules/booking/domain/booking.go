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
	GetByID(id uuid.UUID) (*Booking, error)
	List(filter map[string]interface{}) ([]Booking, error)
	Update(booking *Booking) error
	HasTimeConflict(facilityID uuid.UUID, start, end time.Time) (bool, error)
}
