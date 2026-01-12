package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Sponsor struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID      string     `json:"club_id" gorm:"not null;index"`
	Name        string     `json:"name" gorm:"not null"`
	ContactInfo string     `json:"contact_info,omitempty"`
	LogoURL     string     `json:"logo_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type LocationType string

const (
	LocationWebsiteBanner  LocationType = "WEBSITE_BANNER"
	LocationPhysicalBanner LocationType = "PHYSICAL_BANNER"
	LocationJersey         LocationType = "JERSEY"
)

type AdPlacement struct {
	ID             uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SponsorID      uuid.UUID    `json:"sponsor_id" gorm:"type:uuid;not null;index"`
	Sponsor        *Sponsor     `json:"sponsor,omitempty" gorm:"foreignKey:SponsorID"`
	LocationType   LocationType `json:"location_type"`
	LocationDetail string       `json:"location_detail,omitempty"`
	ContractStart  *time.Time   `json:"contract_start,omitempty"`
	ContractEnd    time.Time    `json:"contract_end" gorm:"not null"`
	AmountPaid     float64      `json:"amount_paid" gorm:"type:decimal(10,2)"`
	IsActive       bool         `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type SponsorRepository interface {
	CreateSponsor(ctx context.Context, sponsor *Sponsor) error
	CreateAdPlacement(ctx context.Context, ad *AdPlacement) error
	GetActiveAds(ctx context.Context, clubID string) ([]AdPlacement, error)
}
