package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MatchEvent struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TrainingGroupID uuid.UUID  `json:"training_group_id" gorm:"type:uuid;not null;index"`
	OpponentName    string     `json:"opponent_name,omitempty"`
	Location        string     `json:"location,omitempty"` // Home/Away desc
	IsHomeGame      bool       `json:"is_home_game" gorm:"default:true"`
	MeetupTime      time.Time  `json:"meetup_time" gorm:"not null"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type PlayerAvailabilityStatus string

const (
	AvailabilityConfirmed PlayerAvailabilityStatus = "CONFIRMED"
	AvailabilityDeclined  PlayerAvailabilityStatus = "DECLINED"
	AvailabilityMaybe     PlayerAvailabilityStatus = "MAYBE"
)

type PlayerAvailability struct {
	MatchEventID uuid.UUID                `json:"match_event_id" gorm:"type:uuid;primary_key"`
	UserID       string                   `json:"user_id" gorm:"primary_key"`
	Status       PlayerAvailabilityStatus `json:"status" gorm:"not null"`
	Reason       string                   `json:"reason,omitempty"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

type TeamRepository interface {
	CreateMatchEvent(ctx context.Context, event *MatchEvent) error
	GetMatchEvent(ctx context.Context, clubID, id string) (*MatchEvent, error)
	SetPlayerAvailability(ctx context.Context, availability *PlayerAvailability) error
	GetEventAvailabilities(ctx context.Context, clubID, eventID string) ([]PlayerAvailability, error)
}
