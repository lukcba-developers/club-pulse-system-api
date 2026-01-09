-- Migration: Gamification Phase 2 - Badges, Missions, Leaderboards
-- Creates tables for the complete gamification system

-- ============================================
-- BADGES SYSTEM
-- ============================================

CREATE TABLE IF NOT EXISTS badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_url TEXT,
    rarity VARCHAR(20) NOT NULL DEFAULT 'COMMON', -- COMMON, RARE, EPIC, LEGENDARY
    category VARCHAR(50) NOT NULL, -- PROGRESSION, STREAK, SOCIAL, TOURNAMENT, BOOKING, SPECIAL
    xp_reward INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(club_id, code)
);

CREATE INDEX IF NOT EXISTS idx_badges_club ON badges(club_id);
CREATE INDEX IF NOT EXISTS idx_badges_category ON badges(category);
CREATE INDEX IF NOT EXISTS idx_badges_rarity ON badges(rarity);

COMMENT ON TABLE badges IS 'Achievement badges that users can earn';
COMMENT ON COLUMN badges.rarity IS 'Badge rarity tier: COMMON, RARE, EPIC, LEGENDARY';
COMMENT ON COLUMN badges.category IS 'Badge category: PROGRESSION, STREAK, SOCIAL, TOURNAMENT, BOOKING, SPECIAL';

CREATE TABLE IF NOT EXISTS user_badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, badge_id)
);

CREATE INDEX IF NOT EXISTS idx_user_badges_user ON user_badges(user_id);
CREATE INDEX IF NOT EXISTS idx_user_badges_badge ON user_badges(badge_id);
CREATE INDEX IF NOT EXISTS idx_user_badges_featured ON user_badges(user_id, featured) WHERE featured = TRUE;

COMMENT ON TABLE user_badges IS 'Badges earned by users';
COMMENT ON COLUMN user_badges.featured IS 'Users can feature up to 3 badges on their profile';

-- ============================================
-- MISSIONS SYSTEM
-- ============================================

CREATE TABLE IF NOT EXISTS missions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL, -- DAILY, WEEKLY
    name VARCHAR(255) NOT NULL,
    description TEXT,
    xp_reward INT DEFAULT 50,
    badge_id UUID REFERENCES badges(id) ON DELETE SET NULL,
    target_value INT DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(club_id, code)
);

CREATE INDEX IF NOT EXISTS idx_missions_club ON missions(club_id);
CREATE INDEX IF NOT EXISTS idx_missions_type ON missions(type);
CREATE INDEX IF NOT EXISTS idx_missions_active ON missions(is_active) WHERE is_active = TRUE;

COMMENT ON TABLE missions IS 'Daily and weekly challenges for users';
COMMENT ON COLUMN missions.type IS 'Mission type: DAILY or WEEKLY';
COMMENT ON COLUMN missions.target_value IS 'Number required to complete (e.g., 3 for "3 bookings")';

CREATE TABLE IF NOT EXISTS user_missions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mission_id UUID NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, COMPLETED, CLAIMED, EXPIRED
    progress INT DEFAULT 0,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    claimed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    UNIQUE(user_id, mission_id, assigned_at)
);

CREATE INDEX IF NOT EXISTS idx_user_missions_user ON user_missions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_missions_status ON user_missions(status);
CREATE INDEX IF NOT EXISTS idx_user_missions_expires ON user_missions(expires_at) WHERE status = 'ACTIVE';

COMMENT ON TABLE user_missions IS 'User progress on assigned missions';
COMMENT ON COLUMN user_missions.status IS 'Mission status: ACTIVE, COMPLETED, CLAIMED, EXPIRED';

-- ============================================
-- LEADERBOARD SUPPORT (Materialized View Approach)
-- ============================================

-- Add total_xp index for efficient leaderboard queries
CREATE INDEX IF NOT EXISTS idx_user_stats_leaderboard 
ON user_stats (level DESC, total_xp DESC);

-- Add bookings count to user_stats for booking leaderboards
ALTER TABLE user_stats ADD COLUMN IF NOT EXISTS total_bookings INT DEFAULT 0;

COMMENT ON COLUMN user_stats.total_bookings IS 'Total confirmed bookings for booking leaderboard';

-- ============================================
-- SEED DEFAULT BADGES (for new clubs)
-- ============================================

-- This would typically be done via application code during club setup,
-- but we provide the structure here for reference
