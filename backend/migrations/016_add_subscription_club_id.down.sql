-- Rollback: Remove club_id from subscriptions table

-- Drop the index first
DROP INDEX IF EXISTS idx_subscriptions_club_id;

-- Remove the column
ALTER TABLE subscriptions DROP COLUMN IF EXISTS club_id;
