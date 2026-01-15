export type MembershipStatus = 'ACTIVE' | 'INACTIVE' | 'PENDING' | 'CANCELLED' | 'EXPIRED';
export type BillingCycle = 'MONTHLY' | 'QUARTERLY' | 'SEMI_ANNUAL' | 'ANNUAL';

export interface MembershipTier {
    id: string;
    club_id: string; // Required - NOT NULL in backend
    name: string;
    description: string;
    monthly_fee: number; // Backend sends decimal.Decimal serialized as number
    duration_days?: number | null; // Optional, overrides standard billing cycle
    colors: string;
    benefits: string[];
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface Membership {
    id: string;
    club_id: string;
    user_id: string;
    membership_tier_id: string;
    membership_tier: MembershipTier;
    status: MembershipStatus;
    billing_cycle: BillingCycle;
    auto_renew: boolean;
    start_date: string;
    end_date?: string;
    next_billing_date: string;
    outstanding_balance: number;
    created_at: string;
    updated_at: string;
}
