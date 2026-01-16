-- Migration: Add club_id to subscriptions table for multi-tenancy isolation
-- SECURITY FIX (VUL-002): Adds tenant isolation column to subscriptions table

-- Add club_id column
ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS club_id VARCHAR(36);

-- Create index for efficient tenant-scoped queries
CREATE INDEX IF NOT EXISTS idx_subscriptions_club_id ON subscriptions(club_id);

-- Backfill club_id from related memberships
-- Subscriptions are linked to memberships via membership_id, which has club_id
UPDATE subscriptions s
SET club_id = m.club_id
FROM memberships m
WHERE s.membership_id = m.id
AND s.club_id IS NULL;

-- Make club_id NOT NULL after backfill (only if all records have been populated)
-- Run separately after verifying backfill:
-- ALTER TABLE subscriptions ALTER COLUMN club_id SET NOT NULL;
