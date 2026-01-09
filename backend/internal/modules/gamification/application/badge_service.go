package application

import (
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
func (s *BadgeService) CheckAndAwardProgressionBadges(clubID, userID string, level int) error {
	levelBadges := map[int]string{
		5:   "LEVEL_5",
		10:  "LEVEL_10",
		25:  "LEVEL_25",
		50:  "LEVEL_50",
		100: "LEVEL_100",
	}

	for requiredLevel, badgeCode := range levelBadges {
		if level >= requiredLevel {
			if err := s.awardBadgeIfNotOwned(clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckAndAwardStreakBadges checks if user qualifies for streak-based badges.
func (s *BadgeService) CheckAndAwardStreakBadges(clubID, userID string, streak int) error {
	streakBadges := map[int]string{
		7:   "STREAK_7",
		30:  "STREAK_30",
		100: "STREAK_100",
	}

	for requiredStreak, badgeCode := range streakBadges {
		if streak >= requiredStreak {
			if err := s.awardBadgeIfNotOwned(clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckAndAwardBookingBadges checks if user qualifies for booking-based badges.
func (s *BadgeService) CheckAndAwardBookingBadges(clubID, userID string, totalBookings int) error {
	bookingBadges := map[int]string{
		10:  "BOOKING_10",
		50:  "BOOKING_50",
		100: "BOOKING_100",
	}

	for requiredBookings, badgeCode := range bookingBadges {
		if totalBookings >= requiredBookings {
			if err := s.awardBadgeIfNotOwned(clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// AwardTournamentBadge awards a tournament-related badge.
func (s *BadgeService) AwardTournamentBadge(clubID, userID string, position int) error {
	var badgeCode string
	switch position {
	case 1:
		badgeCode = "TOURNAMENT_CHAMPION"
	case 2:
		badgeCode = "TOURNAMENT_RUNNER_UP"
	default:
		badgeCode = "TOURNAMENT_PARTICIPANT"
	}

	return s.awardBadgeIfNotOwned(clubID, userID, badgeCode)
}

// AwardReferralBadge checks and awards referral badges based on count.
func (s *BadgeService) AwardReferralBadge(clubID, userID string, referralCount int) error {
	referralBadges := map[int]string{
		1:  "REFERRAL_1",
		5:  "REFERRAL_5",
		10: "REFERRAL_10",
	}

	for requiredReferrals, badgeCode := range referralBadges {
		if referralCount >= requiredReferrals {
			if err := s.awardBadgeIfNotOwned(clubID, userID, badgeCode); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserBadges returns all badges earned by a user.
func (s *BadgeService) GetUserBadges(userID string) ([]domain.UserBadge, error) {
	return s.badgeRepo.GetUserBadges(userID)
}

// GetFeaturedBadges returns the user's featured badges (max 3).
func (s *BadgeService) GetFeaturedBadges(userID string) ([]domain.UserBadge, error) {
	return s.badgeRepo.GetFeaturedBadges(userID)
}

// GetAllBadges returns all badges available in a club.
func (s *BadgeService) GetAllBadges(clubID string) ([]domain.Badge, error) {
	return s.badgeRepo.List(clubID)
}

// SetFeaturedBadge toggles whether a badge is featured on the user's profile.
func (s *BadgeService) SetFeaturedBadge(userID string, badgeID uuid.UUID, featured bool) error {
	if featured {
		// Check if user already has 3 featured badges
		current, err := s.badgeRepo.GetFeaturedBadges(userID)
		if err != nil {
			return err
		}
		if len(current) >= 3 {
			return ErrMaxFeaturedBadges
		}
	}

	return s.badgeRepo.SetFeatured(userID, badgeID, featured)
}

// awardBadgeIfNotOwned awards a badge to a user if they don't already have it.
func (s *BadgeService) awardBadgeIfNotOwned(clubID, userID, badgeCode string) error {
	// Check if user already has this badge
	hasBadge, err := s.badgeRepo.HasBadge(userID, badgeCode)
	if err != nil {
		return err
	}
	if hasBadge {
		return nil // Already has badge
	}

	// Get the badge
	badge, err := s.badgeRepo.GetByCode(clubID, badgeCode)
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

	if err := s.badgeRepo.AwardBadge(userBadge); err != nil {
		return err
	}

	// Grant XP reward if applicable
	if badge.XPReward > 0 {
		s.grantBadgeXP(clubID, userID, badge.XPReward)
	}

	return nil
}

// grantBadgeXP adds XP to user for earning a badge.
func (s *BadgeService) grantBadgeXP(clubID, userID string, xp int) {
	user, err := s.userRepo.GetByID(clubID, userID)
	if err != nil || user == nil || user.Stats == nil {
		return
	}

	user.Stats.Experience += xp
	user.Stats.TotalXP += xp
	user.Stats.UpdatedAt = time.Now()

	_ = s.userRepo.Update(user)
}

// SeedBadgesForClub creates default badges for a new club.
func (s *BadgeService) SeedBadgesForClub(clubID string) error {
	for _, badge := range domain.PredefinedBadges {
		badge.ID = uuid.New()
		badge.ClubID = clubID
		badge.CreatedAt = time.Now()
		badge.UpdatedAt = time.Now()
		badge.IsActive = true

		if err := s.badgeRepo.Create(&badge); err != nil {
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
