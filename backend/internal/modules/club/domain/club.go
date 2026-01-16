package domain

import (
	"context"
	"time"
)

type ClubStatus string

const (
	ClubStatusActive   ClubStatus = "ACTIVE"
	ClubStatusInactive ClubStatus = "INACTIVE"
)

type ClubSettings struct {
	Timezone     string `json:"timezone"`
	Currency     string `json:"currency"`
	Language     string `json:"language"`
	SupportEmail string `json:"support_email"`
}

type Club struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name" gorm:"not null"`
	Slug           string     `json:"slug" gorm:"uniqueIndex;not null"`
	LogoURL        string     `json:"logo_url,omitempty"`
	PrimaryColor   string     `json:"primary_color,omitempty"`
	SecondaryColor string     `json:"secondary_color,omitempty"`
	ContactEmail   string     `json:"contact_email,omitempty"`
	ContactPhone   string     `json:"contact_phone,omitempty"`
	SocialLinks    string     `json:"social_links" gorm:"type:jsonb;serializer:json"` // JSON with social links
	Timezone       string     `json:"timezone" gorm:"default:'UTC'"`
	ThemeConfig    string     `json:"theme_config" gorm:"type:jsonb;serializer:json"` // JSON with colors, fonts
	Domain         string     `json:"domain,omitempty"`
	Status         ClubStatus `json:"status" gorm:"default:'ACTIVE'"`
	Settings       string     `json:"settings" gorm:"type:jsonb;serializer:json"` // JSON settings
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ClubRepository interface {
	Create(ctx context.Context, club *Club) error
	GetByID(ctx context.Context, id string) (*Club, error)
	GetBySlug(ctx context.Context, slug string) (*Club, error)
	GetMemberEmails(ctx context.Context, clubID string) ([]string, error)
	List(ctx context.Context, limit, offset int) ([]Club, error)
	Update(ctx context.Context, club *Club) error
	Delete(ctx context.Context, id string) error
}
