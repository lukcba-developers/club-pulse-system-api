-- Add frequency column to recurring_rules table
ALTER TABLE recurring_rules ADD COLUMN frequency VARCHAR(20) NOT NULL DEFAULT 'WEEKLY';

-- Update existing records if any (optional, assuming WEEKLY as default)
UPDATE recurring_rules SET frequency = 'WEEKLY' WHERE frequency IS NULL;
