package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetXPForAction(t *testing.T) {
	tests := []struct {
		action   XPRewardType
		expected int
	}{
		{XPMatchPlayed, 100},
		{XPMatchWon, 50},
		{XPBookingComplete, 25},
		{XPBookingFirstOfMonth, 10},
		{XPAttendance, 50},
		{XPAttendanceStreak, 25},
		{XPTournamentEnd, 200},
		{XPTournamentTop3, 100},
		{XPProfileComplete, 500},
		{XPReferral, 300},
		{XPDailyMission, 50},
		{XPWeeklyMission, 200},
	}

	for _, tt := range tests {
		t.Run(string(tt.action), func(t *testing.T) {
			result := GetXPForAction(tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetXPForAction_Unknown(t *testing.T) {
	result := GetXPForAction(XPRewardType("UNKNOWN"))
	assert.Equal(t, 0, result)
}

func TestStreakMultiplier(t *testing.T) {
	tests := []struct {
		streakDays int
		expected   float64
	}{
		{0, 1.0},
		{1, 1.0},
		{2, 1.0},
		{3, 1.1},
		{6, 1.1},
		{7, 1.25},
		{13, 1.25},
		{14, 1.5},
		{29, 1.5},
		{30, 2.0},
		{100, 2.0},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.streakDays)), func(t *testing.T) {
			result := StreakMultiplier(tt.streakDays)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateXPWithStreak(t *testing.T) {
	tests := []struct {
		name       string
		baseXP     int
		streakDays int
		expected   int
	}{
		{"No streak", 100, 0, 100},
		{"3 day streak +10%", 100, 3, 110},
		{"7 day streak +25%", 100, 7, 125},
		{"14 day streak +50%", 100, 14, 150},
		{"30 day streak +100%", 100, 30, 200},
		{"With different base XP", 50, 7, 62}, // 50 * 1.25 = 62.5 -> 62
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateXPWithStreak(tt.baseXP, tt.streakDays)
			assert.Equal(t, tt.expected, result)
		})
	}
}
