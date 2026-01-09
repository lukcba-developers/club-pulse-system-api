package domain

// XPRewardType represents the type of action that grants XP.
type XPRewardType string

const (
	XPMatchPlayed         XPRewardType = "MATCH_PLAYED"
	XPMatchWon            XPRewardType = "MATCH_WON"
	XPBookingComplete     XPRewardType = "BOOKING_COMPLETE"
	XPBookingFirstOfMonth XPRewardType = "BOOKING_FIRST_OF_MONTH"
	XPAttendance          XPRewardType = "ATTENDANCE"
	XPAttendanceStreak    XPRewardType = "ATTENDANCE_STREAK_BONUS"
	XPTournamentEnd       XPRewardType = "TOURNAMENT_END"
	XPTournamentTop3      XPRewardType = "TOURNAMENT_TOP3"
	XPProfileComplete     XPRewardType = "PROFILE_COMPLETE"
	XPReferral            XPRewardType = "REFERRAL"
	XPDailyMission        XPRewardType = "DAILY_MISSION"
	XPWeeklyMission       XPRewardType = "WEEKLY_MISSION"
)

// XPRewards contains the base XP values for each action type.
var XPRewards = map[XPRewardType]int{
	XPMatchPlayed:         100,
	XPMatchWon:            50, // Bonus adicional sobre Match Played
	XPBookingComplete:     25,
	XPBookingFirstOfMonth: 10, // Bonus adicional
	XPAttendance:          50,
	XPAttendanceStreak:    25, // Bonus si racha activa
	XPTournamentEnd:       200,
	XPTournamentTop3:      100, // Bonus adicional por top 3
	XPProfileComplete:     500,
	XPReferral:            300,
	XPDailyMission:        50,
	XPWeeklyMission:       200,
}

// GetXPForAction returns the base XP for a given action type.
func GetXPForAction(actionType XPRewardType) int {
	if xp, ok := XPRewards[actionType]; ok {
		return xp
	}
	return 0
}

// StreakMultiplier returns the XP multiplier based on current streak days.
// Returns a value between 1.0 and 2.0.
func StreakMultiplier(streakDays int) float64 {
	switch {
	case streakDays >= 30:
		return 2.0 // +100% XP
	case streakDays >= 14:
		return 1.5 // +50% XP
	case streakDays >= 7:
		return 1.25 // +25% XP
	case streakDays >= 3:
		return 1.1 // +10% XP
	default:
		return 1.0
	}
}

// CalculateXPWithStreak applies the streak multiplier to base XP.
func CalculateXPWithStreak(baseXP, streakDays int) int {
	multiplier := StreakMultiplier(streakDays)
	return int(float64(baseXP) * multiplier)
}
