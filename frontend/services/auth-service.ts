import api from '../lib/axios';

export interface Session {
    id: string;
    user_id: string;
    device_id: string;
    token: string;
    expires_at: string;
    is_revoked: boolean;
    created_at: string;
    updated_at: string;
}

export const authService = {
    listSessions: async () => {
        const response = await api.get<Session[]>('/auth/sessions');
        return response.data;
    },
    revokeSession: async (id: string) => {
        await api.delete(`/auth/sessions/${id}`);
    }
}
