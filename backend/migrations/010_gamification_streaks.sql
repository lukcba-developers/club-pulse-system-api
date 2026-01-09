-- Migration: Add gamification streak fields to user_stats
-- Part of Gamification Phase 1

-- Add streak tracking columns
ALTER TABLE user_stats ADD COLUMN IF NOT EXISTS current_streak INT DEFAULT 0;
ALTER TABLE user_stats ADD COLUMN IF NOT EXISTS longest_streak INT DEFAULT 0;
ALTER TABLE user_stats ADD COLUMN IF NOT EXISTS last_activity_date DATE;
ALTER TABLE user_stats ADD COLUMN IF NOT EXISTS total_xp INT DEFAULT 0;

-- Initialize total_xp from existing experience (one-time migration)
UPDATE user_stats SET total_xp = experience WHERE total_xp = 0 AND experience > 0;

COMMENT ON COLUMN user_stats.current_streak IS 'Consecutive days with activity';
COMMENT ON COLUMN user_stats.longest_streak IS 'Best streak ever achieved';
COMMENT ON COLUMN user_stats.last_activity_date IS 'Date of last qualifying activity for streak';
COMMENT ON COLUMN user_stats.total_xp IS 'Lifetime XP earned (never decreases, unlike experience which resets on level up)';
