-- Migration: Fix Membership Lifecycle Gaps
-- Description: Adds duration_days to membership_tiers and auto_renew to memberships to support fixed-duration plans.

-- 1. Add duration_days to membership_tiers
-- Supports "Day Pass" (1 day), "Week Pass" (7 days) etc.
-- If null, it follows standard BillingCycle (Monthly, Yearly, etc.)
ALTER TABLE membership_tiers 
ADD COLUMN IF NOT EXISTS duration_days INT;

COMMENT ON COLUMN membership_tiers.duration_days IS 'Fixed duration in days. If set, overrides standard billing cycle (e.g., 1 for Day Pass). NULL means standard recurring.';

-- 2. Add auto_renew to memberships
-- Explicit control over renewal. Default to TRUE for backward compatibility with existing recurring memberships.
ALTER TABLE memberships 
ADD COLUMN IF NOT EXISTS auto_renew BOOLEAN NOT NULL DEFAULT TRUE;

COMMENT ON COLUMN memberships.auto_renew IS 'If true, membership renews automatically at end of cycle. Forced to false for fixed-duration tiers.';
