package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VirtualCredential struct {
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
	Token     string `json:"token"`
}

type AccessDirection string
type AccessStatus string

const (
	AccessDirectionIn  AccessDirection = "IN"
	AccessDirectionOut AccessDirection = "OUT"

	AccessStatusGranted AccessStatus = "GRANTED"
	AccessStatusDenied  AccessStatus = "DENIED"
)

type AccessLog struct {
	ID         uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ClubID     string          `gorm:"type:varchar(255);index;not null" json:"club_id"`
	UserID     string          `gorm:"type:text;not null" json:"user_id"`
	FacilityID *uuid.UUID      `gorm:"type:uuid" json:"facility_id,omitempty"`
	Direction  AccessDirection `gorm:"type:varchar(10);not null" json:"direction"`
	Status     AccessStatus    `gorm:"type:varchar(50);not null" json:"status"`
	Reason     string          `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Timestamp  time.Time       `gorm:"not null" json:"timestamp"`
	CreatedAt  time.Time       `json:"created_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

type AccessRepository interface {
	Create(ctx context.Context, log *AccessLog) error
	GetByUserID(ctx context.Context, clubID string, userID string, limit int) ([]AccessLog, error)
}
