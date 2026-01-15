-- Revert Roster Snapshot
DROP TABLE IF EXISTS tournament_team_members;

-- Revert Settings
ALTER TABLE championships DROP COLUMN IF EXISTS settings;

-- Revert Match Scores (Warning: Data Loss of decimals)
ALTER TABLE tournament_matches 
ALTER COLUMN home_score TYPE INT,
ALTER COLUMN away_score TYPE INT;
