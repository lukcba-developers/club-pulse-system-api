-- Migration to support Guest Checkout in Orders
-- Make user_id nullable
ALTER TABLE orders ALTER COLUMN user_id DROP NOT NULL;

-- Add guest fields
ALTER TABLE orders ADD COLUMN guest_name VARCHAR(255);
ALTER TABLE orders ADD COLUMN guest_email VARCHAR(255);
