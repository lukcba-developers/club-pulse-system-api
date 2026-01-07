import api from '../lib/axios';

// Matches Backend 'Membership' Domain/DTO
export interface Membership {
    id: string;
    tier_id: string;
    user_id: string;
    status: string;
    start_date: string;
    end_date: string;
    price: number;
    currency: string;
    outstanding_balance?: number;
}

export interface Tier {
    id: string;
    name: string;
    price: number;
    currency: string;
    description: string;
    features: string[];
}

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

    // Not yet implemented in Backend
    // cancelMembership: async (id: string) => {
    //     const response = await api.delete<{ message: string }>(`/memberships/${id}`);
    //     return response.data;
    // }
}
