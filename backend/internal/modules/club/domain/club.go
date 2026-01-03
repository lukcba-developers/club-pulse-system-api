package domain

import (
	"time"
)

type ClubStatus string

const (
	ClubStatusActive   ClubStatus = "ACTIVE"
	ClubStatusInactive ClubStatus = "INACTIVE"
)

type Club struct {
	ID        string     `json:"id" gorm:"primaryKey"` // Custom ID (slug) or UUID
	Name      string     `json:"name" gorm:"not null"`
	Domain    string     `json:"domain,omitempty"`
	Status    ClubStatus `json:"status" gorm:"default:'ACTIVE'"`
	Settings  string     `json:"settings" gorm:"type:jsonb;serializer:json"` // JSON settings
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ClubRepository interface {
	Create(club *Club) error
	GetByID(id string) (*Club, error)
	List(limit, offset int) ([]Club, error)
	Update(club *Club) error
	Delete(id string) error
}
