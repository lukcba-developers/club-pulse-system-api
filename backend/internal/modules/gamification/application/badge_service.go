package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/gamification/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// BadgeService handles badge-related operations.
type BadgeService struct {
	badgeRepo domain.BadgeRepository
	userRepo  userDomain.UserRepository
}

// NewBadgeService creates a new BadgeService instance.
func NewBadgeService(badgeRepo domain.BadgeRepository, userRepo userDomain.UserRepository) *BadgeService {
	return &BadgeService{
		badgeRepo: badgeRepo,
		userRepo:  userRepo,
	}
}

// CheckAndAwardProgressionBadges checks if user qualifies for level-based badges.
func (s *BadgeService) CheckAndAwardProgressionBadges(ctx context.Context, clubID, userID string, level int) error {
	levelBadges := map[int]string{
		5:   "LEVEL_5",
		10:  "LEVEL_10",
		25:  "LEVEL_25",
		50:  "LEVEL_50",
		100: "LEVEL_100",
	}

	for requiredLevel, badgeCode := range levelBadges {
		if level >= requiredLevel {
			if err := s.awardBadgeIfNotOwned(ctx, clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckAndAwardStreakBadges checks if user qualifies for streak-based badges.
func (s *BadgeService) CheckAndAwardStreakBadges(ctx context.Context, clubID, userID string, streak int) error {
	streakBadges := map[int]string{
		7:   "STREAK_7",
		30:  "STREAK_30",
		100: "STREAK_100",
	}

	for requiredStreak, badgeCode := range streakBadges {
		if streak >= requiredStreak {
			if err := s.awardBadgeIfNotOwned(ctx, clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckAndAwardBookingBadges checks if user qualifies for booking-based badges.
func (s *BadgeService) CheckAndAwardBookingBadges(ctx context.Context, clubID, userID string, totalBookings int) error {
	bookingBadges := map[int]string{
		10:  "BOOKING_10",
		50:  "BOOKING_50",
		100: "BOOKING_100",
	}

	for requiredBookings, badgeCode := range bookingBadges {
		if totalBookings >= requiredBookings {
			if err := s.awardBadgeIfNotOwned(ctx, clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// AwardTournamentBadge awards a tournament-related badge.
func (s *BadgeService) AwardTournamentBadge(ctx context.Context, clubID, userID string, position int) error {
	var badgeCode string
	switch position {
	case 1:
		badgeCode = "TOURNAMENT_CHAMPION"
	case 2:
		badgeCode = "TOURNAMENT_RUNNER_UP"
	default:
		badgeCode = "TOURNAMENT_PARTICIPANT"
	}

	return s.awardBadgeIfNotOwned(ctx, clubID, userID, badgeCode)
}

// AwardReferralBadge checks and awards referral badges based on count.
func (s *BadgeService) AwardReferralBadge(ctx context.Context, clubID, userID string, referralCount int) error {
	referralBadges := map[int]string{
		1:  "REFERRAL_1",
		5:  "REFERRAL_5",
		10: "REFERRAL_10",
	}

	for requiredReferrals, badgeCode := range referralBadges {
		if referralCount >= requiredReferrals {
			if err := s.awardBadgeIfNotOwned(ctx, clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserBadges returns all badges earned by a user.
func (s *BadgeService) GetUserBadges(ctx context.Context, clubID, userID string) ([]domain.UserBadge, error) {
	return s.badgeRepo.GetUserBadges(ctx, clubID, userID)
}

// GetFeaturedBadges returns the user's featured badges (max 3).
func (s *BadgeService) GetFeaturedBadges(ctx context.Context, clubID, userID string) ([]domain.UserBadge, error) {
	return s.badgeRepo.GetFeaturedBadges(ctx, clubID, userID)
}

// GetAllBadges returns all badges available in a club.
func (s *BadgeService) GetAllBadges(ctx context.Context, clubID string) ([]domain.Badge, error) {
	return s.badgeRepo.List(ctx, clubID)
}

// SetFeaturedBadge toggles whether a badge is featured on the user's profile.
func (s *BadgeService) SetFeaturedBadge(ctx context.Context, clubID, userID string, badgeID uuid.UUID, featured bool) error {
	if featured {
		// Check if user already has 3 featured badges
		current, err := s.badgeRepo.GetFeaturedBadges(ctx, clubID, userID)
		if err != nil {
			return err
		}
		if len(current) >= 3 {
			return ErrMaxFeaturedBadges
		}
	}

	return s.badgeRepo.SetFeatured(ctx, clubID, userID, badgeID, featured)
}

// awardBadgeIfNotOwned awards a badge to a user if they don't already have it.
func (s *BadgeService) awardBadgeIfNotOwned(ctx context.Context, clubID, userID, badgeCode string) error {
	// Check if user already has this badge
	hasBadge, err := s.badgeRepo.HasBadge(ctx, clubID, userID, badgeCode)
	if err != nil {
		return err
	}
	if hasBadge {
		return nil // Already has badge
	}

	// Get the badge
	badge, err := s.badgeRepo.GetByCode(ctx, clubID, badgeCode)
	if err != nil {
		return err
	}
	if badge == nil {
		return nil // Badge doesn't exist in this club
	}

	// Award the badge
	userBadge := &domain.UserBadge{
		ID:        uuid.New(),
		UserID:    userID,
		BadgeID:   badge.ID,
		AwardedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.badgeRepo.AwardBadge(ctx, userBadge); err != nil {
		return err
	}

	// Grant XP reward if applicable
	if badge.XPReward > 0 {
		s.grantBadgeXP(ctx, clubID, userID, badge.XPReward)
	}

	return nil
}

// grantBadgeXP adds XP to user for earning a badge.
func (s *BadgeService) grantBadgeXP(ctx context.Context, clubID, userID string, xp int) {
	user, err := s.userRepo.GetByID(context.Background(), clubID, userID)
	if err != nil || user == nil || user.Stats == nil {
		return
	}

	user.Stats.Experience += xp
	user.Stats.TotalXP += xp
	user.Stats.UpdatedAt = time.Now()

	_ = s.userRepo.Update(ctx, user)
}

// SeedBadges for a club.
func (s *BadgeService) SeedBadgesForClub(ctx context.Context, clubID string) error {
	for _, badge := range domain.PredefinedBadges {
		badge.ID = uuid.New()
		badge.ClubID = clubID
		badge.CreatedAt = time.Now()
		badge.UpdatedAt = time.Now()
		badge.IsActive = true

		if err := s.badgeRepo.Create(ctx, &badge); err != nil {
			// Ignore duplicate errors
			continue
		}
	}
	return nil
}

// Custom errors
type BadgeError string

func (e BadgeError) Error() string { return string(e) }

const ErrMaxFeaturedBadges = BadgeError("maximum of 3 featured badges allowed")
