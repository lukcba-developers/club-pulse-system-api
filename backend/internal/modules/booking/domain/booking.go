package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Waitlist struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID     string    `json:"club_id" gorm:"not null"`
	ResourceID uuid.UUID `json:"resource_id" gorm:"type:uuid;not null"`
	TargetDate time.Time `json:"target_date" gorm:"not null"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Status     string    `json:"status" gorm:"default:'PENDING'"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WaitlistEntry = Waitlist

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

type GuestDetail struct {
	Name      string  `json:"name"`
	DNI       string  `json:"dni"`
	FeeAmount float64 `json:"fee_amount"`
}

type GuestDetails []GuestDetail

func (g GuestDetails) Value() (driver.Value, error) {
	if g == nil {
		return "[]", nil
	}
	b, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (g *GuestDetails) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, g)
}

type Booking struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	ClubID       string        `json:"club_id" gorm:"index;not null"`
	UserID       uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	FacilityID   uuid.UUID     `json:"facility_id" gorm:"type:uuid;not null"`
	StartTime    time.Time     `json:"start_time" gorm:"not null"`
	EndTime      time.Time     `json:"end_time" gorm:"not null"`
	TotalPrice   float64       `json:"total_price" gorm:"type:decimal(10,2);default:0"`
	Status       BookingStatus `json:"status" gorm:"type:varchar(20);default:'CONFIRMED'"`
	GuestDetails GuestDetails  `json:"guest_details" gorm:"type:jsonb"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(clubID string, id uuid.UUID) (*Booking, error)
	List(clubID string, filter map[string]interface{}) ([]Booking, error)
	Update(booking *Booking) error
	HasTimeConflict(clubID string, facilityID uuid.UUID, start, end time.Time) (bool, error)
	ListByFacilityAndDate(clubID string, facilityID uuid.UUID, date time.Time) ([]Booking, error)
	ListAll(clubID string, filter map[string]interface{}, from, to *time.Time) ([]Booking, error)
	AddToWaitlist(ctx context.Context, entry *Waitlist) error
	GetNextInLine(ctx context.Context, clubID string, resourceID uuid.UUID, date time.Time) (*Waitlist, error)
}
