package domain

import (
	"context"
	"time"
)

// LeaderboardType defines the type of ranking.
type LeaderboardType string

const (
	LeaderboardTypeGlobal     LeaderboardType = "GLOBAL"
	LeaderboardTypeDiscipline LeaderboardType = "DISCIPLINE"
	LeaderboardTypeCategory   LeaderboardType = "CATEGORY"
	LeaderboardTypeFriends    LeaderboardType = "FRIENDS"
	LeaderboardTypeBookings   LeaderboardType = "BOOKINGS"
)

// LeaderboardPeriod defines the time range for the leaderboard.
type LeaderboardPeriod string

const (
	LeaderboardPeriodDaily   LeaderboardPeriod = "DAILY"
	LeaderboardPeriodWeekly  LeaderboardPeriod = "WEEKLY"
	LeaderboardPeriodMonthly LeaderboardPeriod = "MONTHLY"
	LeaderboardPeriodAllTime LeaderboardPeriod = "ALL_TIME"
)

// LeaderboardEntry represents a single entry in the leaderboard.
type LeaderboardEntry struct {
	Rank          int    `json:"rank"`
	UserID        string `json:"user_id"`
	UserName      string `json:"user_name"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	Score         int    `json:"score"`
	Level         int    `json:"level"`
	Change        int    `json:"change"` // Position change since last period (+2, -1, 0)
	IsCurrentUser bool   `json:"is_current_user,omitempty"`
}

// Leaderboard represents a complete ranking table.
type Leaderboard struct {
	Type      LeaderboardType    `json:"type"`
	Period    LeaderboardPeriod  `json:"period"`
	ClubID    string             `json:"club_id"`
	FilterID  string             `json:"filter_id,omitempty"` // Discipline ID, Category, etc.
	Entries   []LeaderboardEntry `json:"entries"`
	UpdatedAt time.Time          `json:"updated_at"`
	Total     int                `json:"total"` // Total users in ranking
}

// LeaderboardService defines the interface for leaderboard operations.
type LeaderboardService interface {
	// Get leaderboards
	GetGlobalLeaderboard(ctx context.Context, clubID string, period LeaderboardPeriod, limit, offset int) (*Leaderboard, error)
	GetDisciplineLeaderboard(ctx context.Context, clubID, disciplineID string, period LeaderboardPeriod, limit int) (*Leaderboard, error)
	GetCategoryLeaderboard(ctx context.Context, clubID, category string, period LeaderboardPeriod, limit int) (*Leaderboard, error)
	GetFriendsLeaderboard(ctx context.Context, clubID, userID string, period LeaderboardPeriod) (*Leaderboard, error)
	GetBookingsLeaderboard(ctx context.Context, clubID string, period LeaderboardPeriod, limit int) (*Leaderboard, error)

	// User-specific
	GetUserRank(ctx context.Context, clubID, userID string, lbType LeaderboardType, period LeaderboardPeriod) (int, error)
	GetUserContext(ctx context.Context, clubID, userID string, period LeaderboardPeriod) (*LeaderboardContext, error)
}

// LeaderboardContext provides surrounding context for a user's position.
type LeaderboardContext struct {
	UserEntry LeaderboardEntry   `json:"user_entry"`
	Above     []LeaderboardEntry `json:"above"` // 2 users above
	Below     []LeaderboardEntry `json:"below"` // 2 users below
}

// LeaderboardReward defines rewards for top positions.
type LeaderboardReward struct {
	Position    int    `json:"position"` // 1, 2, 3
	XPReward    int    `json:"xp_reward"`
	BadgeCode   string `json:"badge_code,omitempty"`
	DiscountPct int    `json:"discount_pct,omitempty"` // % discount on next booking
}

// DefaultLeaderboardRewards for monthly top 3
var DefaultLeaderboardRewards = []LeaderboardReward{
	{Position: 1, XPReward: 1000, BadgeCode: "MONTHLY_CHAMPION", DiscountPct: 20},
	{Position: 2, XPReward: 500, BadgeCode: "MONTHLY_RUNNER_UP", DiscountPct: 10},
	{Position: 3, XPReward: 250, BadgeCode: "MONTHLY_PODIUM", DiscountPct: 5},
}
