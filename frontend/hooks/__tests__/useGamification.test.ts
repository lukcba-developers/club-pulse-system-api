import { renderHook, waitFor } from '@testing-library/react';
import { useGamification } from '../useGamification';

// Mock global fetch
global.fetch = jest.fn();

describe('useGamification', () => {
    beforeEach(() => {
        (global.fetch as jest.Mock).mockClear();
    });

    it('should fetch stats and update state correctly', async () => {
        const mockStats = {
            level: 2,
            experience: 600,
            total_xp: 1100,
            current_streak: 5,
            longest_streak: 10,
            matches_played: 20,
            matches_won: 15,
            ranking_points: 100,
            next_level_xp: 661 // 500 * 1.15^2
        };

        (global.fetch as jest.Mock).mockResolvedValueOnce({
            ok: true,
            json: async () => mockStats,
        });

        const { result } = renderHook(() => useGamification('user-1'));

        // Initial state
        expect(result.current.loading).toBe(true);

        // Wait for update
        await waitFor(() => {
            expect(result.current.stats).not.toBeNull();
        });

        // Verify stats mapped correctly
        expect(result.current.stats).toEqual({
            level: 2,
            experience: 600,
            totalXp: 1100,
            currentStreak: 5,
            longestStreak: 10,
            matchesPlayed: 20,
            matchesWon: 15,
            rankingPoints: 100,
            nextLevelXP: 661 // Verified mapping from API
        });

        // Verify loaded state
        expect(result.current.loading).toBe(false);
    });

    it('should calculate progress using nextLevelXP from API', async () => {
        const mockStats = {
            level: 1,
            experience: 250,
            total_xp: 250,
            next_level_xp: 500 // 500 needed for level 2
        };

        (global.fetch as jest.Mock).mockResolvedValueOnce({
            ok: true,
            json: async () => mockStats,
        });

        const { result } = renderHook(() => useGamification('user-1'));

        await waitFor(() => {
            expect(result.current.stats).not.toBeNull();
        });

        // 250 / 500 = 50%
        const progress = result.current.calculateProgress();
        expect(progress).toBe(50);
    });

    it('should return nextLevelXP from getter', async () => {
        const mockStats = {
            level: 1,
            next_level_xp: 1000 // Custom req
        };

        (global.fetch as jest.Mock).mockResolvedValueOnce({
            ok: true,
            json: async () => mockStats,
        });

        const { result } = renderHook(() => useGamification('user-1'));

        await waitFor(() => {
            expect(result.current.stats).not.toBeNull();
        });

        expect(result.current.getNextLevelXP()).toBe(1000);
    });

    it('should handle error', async () => {
        (global.fetch as jest.Mock).mockRejectedValueOnce(new Error('Network Error'));

        const { result } = renderHook(() => useGamification('user-1'));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Network Error');
    });
});
