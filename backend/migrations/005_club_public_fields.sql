-- Add public fields to clubs table
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS slug VARCHAR(255);
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS logo_url TEXT;
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS theme_config JSONB;

-- Create unique index for slug
CREATE UNIQUE INDEX IF NOT EXISTS idx_clubs_slug ON clubs(slug);

-- Update existing clubs with a slug based on name (to avoid nulls in unique column)
-- This is a best-effort update for existing data only.
UPDATE clubs 
SET slug = LOWER(REPLACE(name, ' ', '-')) 
WHERE slug IS NULL;

-- Now make slug NOT NULL
ALTER TABLE clubs ALTER COLUMN slug SET NOT NULL;
