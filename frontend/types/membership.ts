export type MembershipStatus = 'ACTIVE' | 'INACTIVE' | 'PENDING' | 'CANCELLED' | 'EXPIRED';
export type BillingCycle = 'MONTHLY' | 'QUARTERLY' | 'SEMI_ANNUAL' | 'ANNUAL';
export type SubscriptionStatus = 'ACTIVE' | 'PAUSED' | 'CANCELLED' | 'PAST_DUE';

export interface MembershipTier {
    id: string;
    club_id: string;
    name: string;
    description: string;
    monthly_fee: string; // Changed from number to string for decimal precision (ARS)
    duration_days?: number | null;
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
    outstanding_balance: string; // Changed from number to string for decimal precision (ARS)
    created_at: string;
    updated_at: string;
}

export interface Subscription {
    id: string;
    user_id: string;
    membership_id: string;
    amount: string; // Decimal as string for precision
    currency: string;
    status: SubscriptionStatus;
    payment_method_id: string;
    next_billing_date: string;
    last_payment_date?: string;
    fail_count: number;
    created_at: string;
    updated_at: string;
}

