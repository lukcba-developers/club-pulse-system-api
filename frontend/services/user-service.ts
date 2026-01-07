import api from '../lib/axios';

export interface User {
    id: string;
    name: string;
    email: string;
    role: string;
    medical_cert_status?: 'VALID' | 'EXPIRED' | 'PENDING';
    medical_cert_expiry?: string;
    family_group_id?: string;
}

export interface UserStats {
    matches_played: number;
    matches_won: number;
    ranking_points: number;
    level: number;
    current_streak: number;
}

export interface Wallet {
    id: string;
    balance: number;
    points: number;
}

export const userService = {
    searchUsers: async (query: string) => {
        const response = await api.get<{ data: User[] }>(`/users?search=${encodeURIComponent(query)}`);
        return response.data.data;
    },
    getChildren: async () => {
        const response = await api.get<{ data: User[] }>(`/users/me/children`);
        return response.data.data;
    },
    registerChild: async (name: string, dateOfBirth: string, email?: string) => {
        const response = await api.post<{ data: User }>('/users/me/children', {
            name,
            date_of_birth: dateOfBirth ? new Date(dateOfBirth).toISOString() : null,
            email
        });
        return response.data.data;
    },

    // --- Gamification ---

    getStats: async (userId: string) => {
        const response = await api.get<UserStats>(`/users/${userId}/stats`);
        return response.data;
    },

    getWallet: async (userId: string) => {
        const response = await api.get<Wallet>(`/users/${userId}/wallet`);
        return response.data;
    },

    createManualDebt: async (userID: string, amount: number, concept: string) => {
        const response = await api.post(`/users/${userID}/debts/manual`, { amount, concept });
        return response.data;
    },

    updateEmergencyInfo: async (data: { contact_name: string; contact_phone: string; insurance_provider: string; insurance_number: string }) => {
        const response = await api.put("/users/me/emergency", data);
        return response.data;
    },

    logIncident: async (data: { injured_user_id?: string; description: string; witnesses?: string; action_taken?: string }) => {
        const response = await api.post("/users/incidents", data);
        return response.data;
    }
}
