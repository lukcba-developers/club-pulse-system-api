import api from '../lib/axios';

export interface AccessResult {
    status: 'GRANTED' | 'DENIED';
    reason?: string;
    timestamp?: string;
}

export const accessService = {
    simulateEntry: async (userId: string) => {
        const response = await api.post<AccessResult>('/access/entry', { user_id: userId, direction: 'IN' });
        return response.data;
    }
}
