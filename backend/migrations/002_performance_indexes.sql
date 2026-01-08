-- Migration to add missing indexes for performance optimization
-- 002_performance_indexes.sql

-- 1. Booking Availability Indexes
-- Speeds up HasTimeConflict and ListByFacilityAndDate (critical for booking flow)
CREATE INDEX IF NOT EXISTS idx_bookings_availability 
ON bookings (club_id, facility_id, status, start_time, end_time);

-- Speeds up "My Bookings" queries
CREATE INDEX IF NOT EXISTS idx_bookings_user 
ON bookings (user_id, start_time DESC);

-- 2. Championship Module Indexes
-- Fixes slow joins in matches retrieval and standings calculation

-- Indexes for 'tournament_matches' (using name from repository code)
-- Supports GetMatchesByGroup (filtering by group and sorting/joining)
CREATE INDEX IF NOT EXISTS idx_tournament_matches_container 
ON tournament_matches (tournament_id, home_team_id, away_team_id);

CREATE INDEX IF NOT EXISTS idx_tournament_matches_group 
ON tournament_matches (group_id);

-- Indexes for 'teams'
-- Supports joins from matches to teams (finding team names)
CREATE INDEX IF NOT EXISTS idx_teams_championship 
ON teams (championship_id); 
-- Note: 'teams' table is shared, referencing 'championship_id' per schema 001.

-- Indexes for 'standings'
-- Supports GetStandings (filtering by group and complex sorting)
CREATE INDEX IF NOT EXISTS idx_standings_ranking 
ON standings (group_id, points DESC, goal_difference DESC, goals_for DESC);

-- Indexes for 'team_members'
-- Supports finding players in team (e.g. for user stats updates)
CREATE INDEX IF NOT EXISTS idx_team_members_user 
ON team_members (user_id);
