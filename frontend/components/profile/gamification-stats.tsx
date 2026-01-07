'use client';

import { useEffect, useState } from 'react';
import { userService, UserStats, Wallet } from '@/services/user-service';
import { useAuthContext } from '@/context/auth-context';

export function GamificationStats() {
    const { user } = useAuthContext();
    const [stats, setStats] = useState<UserStats | null>(null);
    const [wallet, setWallet] = useState<Wallet | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            try {
                const [statsData, walletData] = await Promise.all([
                    userService.getStats(user!.id),
                    userService.getWallet(user!.id)
                ]);
                setStats(statsData);
                setWallet(walletData);
            } catch (error) {
                console.error("Failed to load gamification data", error);
            } finally {
                setLoading(false);
            }
        };

        if (user?.id) {
            loadData();
        }
    }, [user]);

    if (loading) return <div className="animate-pulse h-32 bg-gray-100 rounded-lg"></div>;

    // Use default values if no data (e.g. freshly created user before first periodic sync)
    // Though endpoint usually returns defaults.
    const points = stats?.ranking_points ?? 1000;
    const level = stats?.level ?? 1;
    const currentStreak = stats?.current_streak ?? 0;
    const balance = wallet?.balance ?? 0;

    return (
        <div className="bg-white rounded-lg shadow p-6 border border-gray-100">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Player Stats</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="text-center p-3 bg-indigo-50 rounded-lg">
                    <span className="block text-xl font-bold text-indigo-700">{level}</span>
                    <span className="text-xs text-indigo-500 uppercase font-semibold tracking-wide">Level</span>
                </div>
                <div className="text-center p-3 bg-purple-50 rounded-lg">
                    <span className="block text-xl font-bold text-purple-700">{points}</span>
                    <span className="text-xs text-purple-500 uppercase font-semibold tracking-wide">Points</span>
                </div>
                <div className="text-center p-3 bg-orange-50 rounded-lg">
                    <span className="block text-xl font-bold text-orange-700">{currentStreak} 🔥</span>
                    <span className="text-xs text-orange-500 uppercase font-semibold tracking-wide">Streak</span>
                </div>
                <div className="text-center p-3 bg-green-50 rounded-lg">
                    <span className="block text-xl font-bold text-green-700">${balance.toFixed(2)}</span>
                    <span className="text-xs text-green-500 uppercase font-semibold tracking-wide">Wallet</span>
                </div>
            </div>
        </div>
    );
}
