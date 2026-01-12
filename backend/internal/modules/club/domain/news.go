package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID    string    `json:"club_id" gorm:"index;not null"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	ImageURL  string    `json:"image_url,omitempty"`
	Published bool      `json:"published" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewsRepository interface {
	CreateNews(ctx context.Context, news *News) error
	GetNewsByClub(ctx context.Context, clubID string, limit, offset int) ([]News, error)
	GetPublicNewsByClub(ctx context.Context, clubID string, limit, offset int) ([]News, error)
	GetNewsByID(ctx context.Context, id uuid.UUID) (*News, error)
	UpdateNews(ctx context.Context, news *News) error
	DeleteNews(ctx context.Context, id uuid.UUID) error
}
