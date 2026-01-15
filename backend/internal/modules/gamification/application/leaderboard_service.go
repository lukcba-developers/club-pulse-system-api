package application

import (
	"context"
	"sort"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/gamification/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// LeaderboardServiceImpl implements the LeaderboardService interface.
type LeaderboardServiceImpl struct {
	userRepo userDomain.UserRepository
}

// NewLeaderboardService creates a new LeaderboardService.
func NewLeaderboardService(userRepo userDomain.UserRepository) *LeaderboardServiceImpl {
	return &LeaderboardServiceImpl{
		userRepo: userRepo,
	}
}

// GetGlobalLeaderboard returns the global XP leaderboard for a club.
func (s *LeaderboardServiceImpl) GetGlobalLeaderboard(ctx context.Context, clubID string, period domain.LeaderboardPeriod, limit, offset int) (*domain.Leaderboard, error) {
	// Get all users with stats
	users, err := s.userRepo.List(ctx, clubID, 1000, 0, nil) // TODO: Add pagination/filtering in repo
	if err != nil {
		return nil, err
	}

	// Build entries from users
	entries := s.buildEntriesFromUsers(users, period)

	// Sort by score (TotalXP or Experience based on period)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Level > entries[j].Level
	})

	// Apply ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Apply pagination
	total := len(entries)
	if offset >= len(entries) {
		entries = []domain.LeaderboardEntry{}
	} else {
		end := offset + limit
		if end > len(entries) {
			end = len(entries)
		}
		entries = entries[offset:end]
	}

	return &domain.Leaderboard{
		Type:      domain.LeaderboardTypeGlobal,
		Period:    period,
		ClubID:    clubID,
		Entries:   entries,
		UpdatedAt: time.Now(),
		Total:     total,
	}, nil
}

// GetUserRank returns a user's rank in the specified leaderboard.
func (s *LeaderboardServiceImpl) GetUserRank(ctx context.Context, clubID, userID string, lbType domain.LeaderboardType, period domain.LeaderboardPeriod) (int, error) {
	leaderboard, err := s.GetGlobalLeaderboard(ctx, clubID, period, 1000, 0)
	if err != nil {
		return 0, err
	}

	for _, entry := range leaderboard.Entries {
		if entry.UserID == userID {
			return entry.Rank, nil
		}
	}

	return 0, nil // User not found in leaderboard
}

// GetUserContext returns the user's position with surrounding entries.
func (s *LeaderboardServiceImpl) GetUserContext(ctx context.Context, clubID, userID string, period domain.LeaderboardPeriod) (*domain.LeaderboardContext, error) {
	leaderboard, err := s.GetGlobalLeaderboard(ctx, clubID, period, 1000, 0)
	if err != nil {
		return nil, err
	}

	var userIndex = -1
	for i, entry := range leaderboard.Entries {
		if entry.UserID == userID {
			userIndex = i
			break
		}
	}

	if userIndex == -1 {
		return nil, nil // User not in leaderboard
	}

	context := &domain.LeaderboardContext{
		UserEntry: leaderboard.Entries[userIndex],
	}
	context.UserEntry.IsCurrentUser = true

	// Get 2 users above
	start := userIndex - 2
	if start < 0 {
		start = 0
	}
	for i := start; i < userIndex; i++ {
		context.Above = append(context.Above, leaderboard.Entries[i])
	}

	// Get 2 users below
	end := userIndex + 3
	if end > len(leaderboard.Entries) {
		end = len(leaderboard.Entries)
	}
	for i := userIndex + 1; i < end; i++ {
		context.Below = append(context.Below, leaderboard.Entries[i])
	}

	return context, nil
}

// GetBookingsLeaderboard returns leaderboard based on total bookings.
func (s *LeaderboardServiceImpl) GetBookingsLeaderboard(ctx context.Context, clubID string, period domain.LeaderboardPeriod, limit int) (*domain.Leaderboard, error) {
	users, err := s.userRepo.List(ctx, clubID, 1000, 0, nil)
	if err != nil {
		return nil, err
	}

	entries := s.buildBookingEntriesFromUsers(users)

	// Sort by booking count
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	// Apply ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Limit entries
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}

	return &domain.Leaderboard{
		Type:      domain.LeaderboardTypeBookings,
		Period:    period,
		ClubID:    clubID,
		Entries:   entries,
		UpdatedAt: time.Now(),
		Total:     len(entries),
	}, nil
}

// buildEntriesFromUsers converts user data to leaderboard entries.
func (s *LeaderboardServiceImpl) buildEntriesFromUsers(users []userDomain.User, period domain.LeaderboardPeriod) []domain.LeaderboardEntry {
	var entries []domain.LeaderboardEntry

	for _, user := range users {
		if user.Stats == nil {
			continue
		}

		score := user.Stats.TotalXP
		if period == domain.LeaderboardPeriodMonthly || period == domain.LeaderboardPeriodWeekly {
			// For periodic leaderboards, we'd need to track XP gained in period
			// For now, use TotalXP as fallback
			score = user.Stats.TotalXP
		}

		entries = append(entries, domain.LeaderboardEntry{
			UserID:   user.ID,
			UserName: user.Name,
			Score:    float64(score),
			Level:    user.Stats.Level,
			Change:   0, // Would require historical tracking
		})
	}

	return entries
}

// buildBookingEntriesFromUsers creates entries based on booking count.
func (s *LeaderboardServiceImpl) buildBookingEntriesFromUsers(users []userDomain.User) []domain.LeaderboardEntry {
	var entries []domain.LeaderboardEntry

	for _, user := range users {
		bookings := 0
		if user.Stats != nil {
			bookings = user.Stats.MatchesPlayed // Using MatchesPlayed as proxy; would use TotalBookings from stats
		}

		entries = append(entries, domain.LeaderboardEntry{
			UserID:   user.ID,
			UserName: user.Name,
			Score:    float64(bookings),
			Level:    1,
		})
	}

	return entries
}
