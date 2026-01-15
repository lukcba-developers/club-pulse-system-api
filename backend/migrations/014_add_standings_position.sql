-- Add position column to standings table for ranking display
ALTER TABLE standings ADD COLUMN IF NOT EXISTS position INTEGER DEFAULT 0;
