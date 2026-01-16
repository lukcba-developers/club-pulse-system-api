package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// BadgeRarity represents the rarity tier of a badge.
type BadgeRarity string

const (
	BadgeRarityCommon    BadgeRarity = "COMMON"
	BadgeRarityRare      BadgeRarity = "RARE"
	BadgeRarityEpic      BadgeRarity = "EPIC"
	BadgeRarityLegendary BadgeRarity = "LEGENDARY"
)

// BadgeCategory groups badges by type.
type BadgeCategory string

const (
	BadgeCategoryProgression BadgeCategory = "PROGRESSION"
	BadgeCategoryStreak      BadgeCategory = "STREAK"
	BadgeCategorySocial      BadgeCategory = "SOCIAL"
	BadgeCategoryTournament  BadgeCategory = "TOURNAMENT"
	BadgeCategoryBooking     BadgeCategory = "BOOKING"
	BadgeCategorySpecial     BadgeCategory = "SPECIAL"
)

// Badge represents an achievement that can be earned by users.
type Badge struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID      string        `json:"club_id" gorm:"index;not null"`
	Code        string        `json:"code" gorm:"uniqueIndex:idx_badges_club_code;not null"` // "LEVEL_10", "WEEK_STREAK"
	Name        string        `json:"name" gorm:"not null"`
	Description string        `json:"description"`
	IconURL     string        `json:"icon_url"`
	Rarity      BadgeRarity   `json:"rarity" gorm:"default:'COMMON'"`
	Category    BadgeCategory `json:"category" gorm:"not null"`
	XPReward    int           `json:"xp_reward" gorm:"default:0"` // XP granted when badge is earned
	IsActive    bool          `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (Badge) TableName() string {
	return "badges"
}

// UserBadge represents a badge earned by a specific user.
type UserBadge struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"not null;index:idx_user_badges_user"`
	BadgeID   uuid.UUID `json:"badge_id" gorm:"type:uuid;not null;index:idx_user_badges_badge"`
	AwardedAt time.Time `json:"awarded_at"`
	Featured  bool      `json:"featured" gorm:"default:false"` // User can feature up to 3 badges
	CreatedAt time.Time `json:"created_at"`
}

func (UserBadge) TableName() string {
	return "user_badges"
}

// BadgeRepository defines the interface for badge persistence.
type BadgeRepository interface {
	// Badge CRUD
	Create(ctx context.Context, badge *Badge) error
	GetByID(ctx context.Context, clubID string, id uuid.UUID) (*Badge, error)
	GetByCode(ctx context.Context, clubID, code string) (*Badge, error)
	List(ctx context.Context, clubID string) ([]Badge, error)
	Update(ctx context.Context, badge *Badge) error

	// UserBadge operations
	AwardBadge(ctx context.Context, userBadge *UserBadge) error
	GetUserBadges(ctx context.Context, clubID, userID string) ([]UserBadge, error)
	HasBadge(ctx context.Context, clubID, userID string, badgeCode string) (bool, error)
	SetFeatured(ctx context.Context, clubID, userID string, badgeID uuid.UUID, featured bool) error
	GetFeaturedBadges(ctx context.Context, clubID, userID string) ([]UserBadge, error)
}

// PredefinedBadges contains the default badges for the system.
var PredefinedBadges = []Badge{
	// Progression Badges
	{Code: "LEVEL_5", Name: "Novato", Description: "Alcanzaste el nivel 5", Rarity: BadgeRarityCommon, Category: BadgeCategoryProgression, XPReward: 50},
	{Code: "LEVEL_10", Name: "Aprendiz", Description: "Alcanzaste el nivel 10", Rarity: BadgeRarityCommon, Category: BadgeCategoryProgression, XPReward: 100},
	{Code: "LEVEL_25", Name: "Veterano", Description: "Alcanzaste el nivel 25", Rarity: BadgeRarityRare, Category: BadgeCategoryProgression, XPReward: 250},
	{Code: "LEVEL_50", Name: "Experto", Description: "Alcanzaste el nivel 50", Rarity: BadgeRarityEpic, Category: BadgeCategoryProgression, XPReward: 500},
	{Code: "LEVEL_100", Name: "Leyenda", Description: "Alcanzaste el nivel 100", Rarity: BadgeRarityLegendary, Category: BadgeCategoryProgression, XPReward: 1000},

	// Streak Badges
	{Code: "STREAK_7", Name: "Semana Perfecta", Description: "7 días consecutivos de actividad", Rarity: BadgeRarityRare, Category: BadgeCategoryStreak, XPReward: 100},
	{Code: "STREAK_30", Name: "Mes Imbatible", Description: "30 días consecutivos de actividad", Rarity: BadgeRarityEpic, Category: BadgeCategoryStreak, XPReward: 300},
	{Code: "STREAK_100", Name: "Imparable", Description: "100 días consecutivos de actividad", Rarity: BadgeRarityLegendary, Category: BadgeCategoryStreak, XPReward: 1000},

	// Booking Badges
	{Code: "BOOKING_10", Name: "Habitual", Description: "10 reservas completadas", Rarity: BadgeRarityCommon, Category: BadgeCategoryBooking, XPReward: 50},
	{Code: "BOOKING_50", Name: "Cliente VIP", Description: "50 reservas completadas", Rarity: BadgeRarityRare, Category: BadgeCategoryBooking, XPReward: 200},
	{Code: "BOOKING_100", Name: "Miembro Platino", Description: "100 reservas completadas", Rarity: BadgeRarityEpic, Category: BadgeCategoryBooking, XPReward: 500},

	// Tournament Badges
	{Code: "TOURNAMENT_CHAMPION", Name: "Campeón", Description: "Ganaste un torneo", Rarity: BadgeRarityEpic, Category: BadgeCategoryTournament, XPReward: 500},
	{Code: "TOURNAMENT_RUNNER_UP", Name: "Subcampeón", Description: "Segundo lugar en un torneo", Rarity: BadgeRarityRare, Category: BadgeCategoryTournament, XPReward: 250},
	{Code: "TOURNAMENT_PARTICIPANT", Name: "Competidor", Description: "Participaste en un torneo", Rarity: BadgeRarityCommon, Category: BadgeCategoryTournament, XPReward: 50},

	// Social Badges
	{Code: "REFERRAL_1", Name: "Embajador", Description: "Invitaste a 1 amigo", Rarity: BadgeRarityCommon, Category: BadgeCategorySocial, XPReward: 100},
	{Code: "REFERRAL_5", Name: "Reclutador", Description: "Invitaste a 5 amigos", Rarity: BadgeRarityRare, Category: BadgeCategorySocial, XPReward: 300},
	{Code: "REFERRAL_10", Name: "Influencer", Description: "Invitaste a 10 amigos", Rarity: BadgeRarityEpic, Category: BadgeCategorySocial, XPReward: 500},
}
