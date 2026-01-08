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
	ID          string     `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;not null"`
	LogoURL     string     `json:"logo_url,omitempty"`
	ThemeConfig string     `json:"theme_config" gorm:"type:jsonb;serializer:json"` // JSON with colors, fonts
	Domain      string     `json:"domain,omitempty"`
	Status      ClubStatus `json:"status" gorm:"default:'ACTIVE'"`
	Settings    string     `json:"settings" gorm:"type:jsonb;serializer:json"` // JSON settings
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ClubRepository interface {
	Create(club *Club) error
	GetByID(id string) (*Club, error)
	GetBySlug(slug string) (*Club, error)
	GetMemberEmails(clubID string) ([]string, error)
	List(limit, offset int) ([]Club, error)
	Update(club *Club) error
	Delete(id string) error
}
