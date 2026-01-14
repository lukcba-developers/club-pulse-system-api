package domain

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        string    `json:"user_id" gorm:"type:varchar(100);not null;uniqueIndex"`
	MatchesPlayed int       `json:"matches_played" gorm:"default:0"`
	MatchesWon    int       `json:"matches_won" gorm:"default:0"`
	RankingPoints int       `json:"ranking_points" gorm:"default:0"`
	NextLevelXP   int       `json:"next_level_xp" gorm:"-"` // Computed field
	Level         int       `json:"level" gorm:"default:1"`
	Experience    int       `json:"experience" gorm:"default:0"`

	// Gamification: Streak System
	CurrentStreak    int        `json:"current_streak" gorm:"default:0"`
	LongestStreak    int        `json:"longest_streak" gorm:"default:0"`
	LastActivityDate *time.Time `json:"last_activity_date,omitempty"`

	// Gamification: Total XP (historical, never decreases)
	TotalXP int `json:"total_xp" gorm:"default:0"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserStats) TableName() string {
	return "user_stats"
}

// CalculateLevel determines level based on XP containing geometric progression (factor 1.15)
// Base XP for Level 2 is 500
func (s *UserStats) CalculateLevel() int {
	xp := float64(s.TotalXP)
	if xp < 500 {
		return 1
	}
	// Formula: XP = 500 * (1.15 ^ (level - 1))
	// log(XP/500) / log(1.15) = level - 1
	// level = 1 + int(math.Log(xp/500.0)/math.Log(1.15))
	level := 1 + int(math.Log(xp/500.0)/math.Log(1.15))
	return level
}

// CalculateNextLevelXP calculates the XP required to reach the next level
func (s *UserStats) CalculateNextLevelXP() int {
	return int(500 * math.Pow(1.15, float64(s.Level)))
}
