export type MembershipStatus = 'ACTIVE' | 'INACTIVE' | 'PENDING' | 'CANCELLED' | 'EXPIRED';
export type BillingCycle = 'MONTHLY' | 'QUARTERLY' | 'SEMI_ANNUAL' | 'ANNUAL';

export interface MembershipTier {
    id: string;
    name: string;
    description: string;
    monthly_fee: number; // Decimal in standard JSON often comes as number, handle string if needed
    colors: string;
    benefits: string[];
}

export interface Membership {
    id: string;
    user_id: string;
    membership_tier_id: string;
    membership_tier: MembershipTier;
    status: MembershipStatus;
    billing_cycle: BillingCycle;
    start_date: string;
    end_date?: string;
    next_billing_date: string;
}
