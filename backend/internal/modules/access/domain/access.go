package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessDirection string
type AccessStatus string

const (
	AccessDirectionIn  AccessDirection = "IN"
	AccessDirectionOut AccessDirection = "OUT"

	AccessStatusGranted AccessStatus = "GRANTED"
	AccessStatusDenied  AccessStatus = "DENIED"
)

type AccessLog struct {
	ID         uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
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
	GetByUserID(ctx context.Context, userID string, limit int) ([]AccessLog, error)
}
