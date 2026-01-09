import api from '../lib/axios';
import type { Membership, MembershipTier } from '../types/membership';

// Re-export for convenience
export type { Membership, MembershipTier };

// Alias para compatibilidad con código existente
export type Tier = MembershipTier;

export interface CreateMembershipRequest {
    tier_id: string;
    start_date: string; // YYYY-MM-DD
}


export const membershipService = {
    listMyMemberships: async () => {
        // Backend: GET /memberships -> returns { data: Membership[] }
        const response = await api.get<{ data: Membership[] }>('/memberships');
        return response.data.data;
    },

    listTiers: async () => {
        // Backend: GET /memberships/tiers -> returns { data: Tier[] }
        const response = await api.get<{ data: Tier[] }>('/memberships/tiers');
        return response.data.data;
    },

    createMembership: async (data: CreateMembershipRequest) => {
        // Backend: POST /memberships -> returns { data: Membership }
        const response = await api.post<{ data: Membership }>('/memberships', data);
        return response.data.data;
    },

    // Admin: List all memberships in the club
    listAllMemberships: async () => {
        const response = await api.get<{ data: Membership[] }>('/memberships/admin');
        return response.data.data;
    },

    // Cancel a membership
    cancelMembership: async (id: string) => {
        const response = await api.delete<{ message: string; data: Membership }>(`/memberships/${id}`);
        return response.data;
    }
}
