package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        string    `json:"user_id" gorm:"type:varchar(100);not null;uniqueIndex"`
	MatchesPlayed int       `json:"matches_played" gorm:"default:0"`
	MatchesWon    int       `json:"matches_won" gorm:"default:0"`
	RankingPoints int       `json:"ranking_points" gorm:"default:0"`
	Level         int       `json:"level" gorm:"default:1"`
	Experience    int       `json:"experience" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
