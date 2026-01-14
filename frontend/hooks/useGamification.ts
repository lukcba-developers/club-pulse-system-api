"use client";

import { useState, useEffect, useCallback } from "react";

interface UserStats {
    level: number;
    experience: number;
    totalXp: number;
    currentStreak: number;
    longestStreak: number;
    matchesPlayed: number;
    matchesWon: number;
    rankingPoints: number;
    nextLevelXP: number;
}

interface GamificationState {
    stats: UserStats | null;
    previousLevel: number;
    showLevelUpModal: boolean;
    newLevel: number;
    loading: boolean;
    error: string | null;
}

interface UseGamificationReturn extends GamificationState {
    fetchStats: () => Promise<void>;
    closeLevelUpModal: () => void;
    calculateProgress: () => number;
    getNextLevelXP: () => number;
}




export function useGamification(userId?: string): UseGamificationReturn {
    const [state, setState] = useState<GamificationState>({
        stats: null,
        previousLevel: 0,
        showLevelUpModal: false,
        newLevel: 0,
        loading: false,
        error: null,
    });

    const fetchStats = useCallback(async () => {
        if (!userId) return;

        setState((prev) => ({ ...prev, loading: true, error: null }));

        try {
            const response = await fetch(`/api/v1/users/${userId}/stats`, {
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to fetch stats");
            }

            const data = await response.json();

            setState((prev) => {
                const newStats: UserStats = {
                    level: data.level || 1,
                    experience: data.experience || 0,
                    totalXp: data.total_xp || 0,
                    currentStreak: data.current_streak || 0,
                    longestStreak: data.longest_streak || 0,
                    matchesPlayed: data.matches_played || 0,
                    matchesWon: data.matches_won || 0,
                    rankingPoints: data.ranking_points || 0,
                    nextLevelXP: data.next_level_xp || 0,
                };

                // Check for level up
                const previousLevel = prev.stats?.level || prev.previousLevel;
                const didLevelUp = previousLevel > 0 && newStats.level > previousLevel;

                return {
                    ...prev,
                    stats: newStats,
                    previousLevel: newStats.level,
                    showLevelUpModal: didLevelUp,
                    newLevel: didLevelUp ? newStats.level : prev.newLevel,
                    loading: false,
                };
            });
        } catch (error) {
            setState((prev) => ({
                ...prev,
                loading: false,
                error: error instanceof Error ? error.message : "Unknown error",
            }));
        }
    }, [userId]);

    const closeLevelUpModal = useCallback(() => {
        setState((prev) => ({ ...prev, showLevelUpModal: false }));
    }, []);

    const calculateProgress = useCallback((): number => {
        if (!state.stats) return 0;
        const required = state.stats.nextLevelXP;
        return Math.min((state.stats.experience / required) * 100, 100);
    }, [state.stats]);

    const getNextLevelXP = useCallback((): number => {
        if (!state.stats) return 575; // Default level 1 requirement
        return state.stats.nextLevelXP;
    }, [state.stats]);

    // Initial fetch
    useEffect(() => {
        fetchStats();
    }, [fetchStats]);

    // Periodic refresh (every 30 seconds)
    useEffect(() => {
        const interval = setInterval(fetchStats, 30000);
        return () => clearInterval(interval);
    }, [fetchStats]);

    return {
        ...state,
        fetchStats,
        closeLevelUpModal,
        calculateProgress,
        getNextLevelXP,
    };
}

// Utility hook for displaying streak info
export function useStreakStatus(currentStreak: number): {
    label: string;
    emoji: string;
    multiplier: number;
    color: string;
} {
    if (currentStreak >= 30) {
        return {
            label: "¡Racha Legendaria!",
            emoji: "🔥🔥🔥",
            multiplier: 2.0,
            color: "text-orange-500",
        };
    }
    if (currentStreak >= 14) {
        return {
            label: "¡Racha Épica!",
            emoji: "🔥🔥",
            multiplier: 1.5,
            color: "text-purple-500",
        };
    }
    if (currentStreak >= 7) {
        return {
            label: "¡Buena Racha!",
            emoji: "🔥",
            multiplier: 1.25,
            color: "text-blue-500",
        };
    }
    if (currentStreak >= 3) {
        return {
            label: "Racha Activa",
            emoji: "✨",
            multiplier: 1.1,
            color: "text-green-500",
        };
    }
    return {
        label: "Sin Racha",
        emoji: "💤",
        multiplier: 1.0,
        color: "text-gray-500",
    };
}
