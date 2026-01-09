-- Migration: Add gamification streak fields to user_stats
-- Part of Gamification Phase 1

-- Ensure table exists (safeguard for CI/CD if 001 didn't run fully)
CREATE TABLE IF NOT EXISTS user_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    matches_played INT DEFAULT 0,
    matches_won INT DEFAULT 0,
    ranking_points INT DEFAULT 0,
    level INT DEFAULT 1,
    experience INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

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
