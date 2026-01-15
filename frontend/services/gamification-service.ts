import api from '@/lib/axios';

export type LeaderboardPeriod = 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'ALL_TIME';

export interface LeaderboardEntry {
    rank: number;
    user_id: string;
    user_name: string;
    avatar_url?: string;
    score: number;
    level: number;
    change: number; // +2, -1, 0
    is_current_user?: boolean;
}

export interface Leaderboard {
    type: string;
    period: LeaderboardPeriod;
    club_id: string;
    entries: LeaderboardEntry[];
    updated_at: string;
    total: number;
}

export interface LeaderboardContext {
    user_entry: LeaderboardEntry;
    above: LeaderboardEntry[];
    below: LeaderboardEntry[];
}

export interface Badge {
    id: string;
    code: string;
    name: string;
    description: string;
    image_url: string;
    category: string;
}

export interface UserBadge {
    badge_id: string;
    user_id: string;
    earned_at: string;
    is_featured: boolean;
    badge: Badge; // Enriched info
}

export const gamificationService = {
    getGlobalLeaderboard: async (period: LeaderboardPeriod = 'MONTHLY', limit: number = 20) => {
        const response = await api.get<Leaderboard>('/gamification/leaderboard', {
            params: { period, limit }
        });
        return response.data;
    },

    getMyLeaderboardContext: async (period: LeaderboardPeriod = 'MONTHLY') => {
        const response = await api.get<LeaderboardContext>('/gamification/leaderboard/context', {
            params: { period }
        });
        return response.data;
    },

    getMyRank: async (period: LeaderboardPeriod = 'MONTHLY') => {
        const response = await api.get<{ rank: number }>('/gamification/leaderboard/rank', {
            params: { period }
        });
        return response.data;
    },

    listBadges: async () => {
        const response = await api.get<Badge[]>('/gamification/badges');
        return response.data;
    },

    getMyBadges: async () => {
        const response = await api.get<UserBadge[]>('/gamification/badges/my');
        return response.data;
    },

    getFeaturedBadges: async (userId: string) => {
        const response = await api.get<UserBadge[]>(`/gamification/badges/featured/${userId}`);
        return response.data;
    },

    setFeaturedBadge: async (badgeId: string, featured: boolean) => {
        const response = await api.put<{ message: string }>(`/gamification/badges/${badgeId}/feature`, { featured });
        return response.data;
    }
};
