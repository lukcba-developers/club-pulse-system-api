-- Add club_id to payments table
ALTER TABLE payments ADD COLUMN club_id TEXT;
CREATE INDEX idx_payments_club_id ON payments(club_id);
-- For existing rows, we leave it NULL or set to a default if known. 
-- Since this is Multi-Tenant from start mostly, we might want to enforce NOT NULL later.
