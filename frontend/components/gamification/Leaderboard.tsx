"use client";

import { useState, useEffect, useCallback } from "react";
import { cn } from "@/lib/utils";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuthContext as useAuth } from "@/context/auth-context";

interface LeaderboardEntry {
    rank: number;
    user_id: string;
    first_name: string;
    last_name: string;
    avatar_url?: string;
    points: number;
    change?: number;
    isCurrentUser?: boolean;
}

interface ApiLeaderboardEntry {
    rank: number;
    user_id: string;
    user_name: string;
    avatar_url?: string;
    score: number;
    level: number;
    change?: number;
}

export function Leaderboard() {
    const { user } = useAuth();
    const [period, setPeriod] = useState<"WEEKLY" | "MONTHLY" | "ALL_TIME">("MONTHLY");
    const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
    const [loading, setLoading] = useState(true);
    // Removed unused setter, kept variable as it is used in the JSX
    const [userContext] = useState<{ rank: number } | null>(null);

    const fetchLeaderboard = useCallback(async () => {
        setLoading(true);
        try {
            const currentUserId = user?.id; // Define currentUserId here
            const res = await fetch(`/api/v1/gamification/leaderboard?period=${period}&limit=20`, {
                credentials: "include",
            });
            if (res.ok) {
                const data = await res.json();
                setEntries(
                    data.entries?.map((e: ApiLeaderboardEntry) => ({
                        rank: e.rank,
                        user_id: e.user_id,
                        first_name: e.user_name.split(' ')[0] || '', // Fallback splitting name
                        last_name: e.user_name.split(' ').slice(1).join(' ') || '',
                        avatar_url: e.avatar_url,
                        points: e.score,
                        change: e.change,
                        isCurrentUser: e.user_id === currentUserId,
                    })) || []
                );
            }
        } catch (error) {
            console.error("Failed to fetch leaderboard:", error);
        }
        setLoading(false);
    }, [period, user]);

    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        fetchLeaderboard();
    }, [fetchLeaderboard]);

    const getRankStyle = (rank: number) => {
        switch (rank) {
            case 1:
                return "bg-gradient-to-r from-yellow-400 to-amber-500 text-white shadow-lg";
            case 2:
                return "bg-gradient-to-r from-gray-300 to-gray-400 text-gray-800";
            case 3:
                return "bg-gradient-to-r from-amber-600 to-amber-700 text-white";
            default:
                return "bg-muted";
        }
    };

    const getRankEmoji = (rank: number) => {
        switch (rank) {
            case 1:
                return "🥇";
            case 2:
                return "🥈";
            case 3:
                return "🥉";
            default:
                return null;
        }
    };

    return (
        <Card className="w-full">
            <CardHeader className="pb-2">
                <CardTitle className="flex items-center gap-2">
                    <span>🏆</span> Tabla de Posiciones
                </CardTitle>
            </CardHeader>

            <CardContent>
                {/* Period Tabs */}
                <Tabs
                    value={period}
                    onValueChange={(v) => setPeriod(v as typeof period)}
                    className="mb-4"
                >
                    <TabsList className="grid w-full grid-cols-3">
                        <TabsTrigger value="WEEKLY">Semanal</TabsTrigger>
                        <TabsTrigger value="MONTHLY">Mensual</TabsTrigger>
                        <TabsTrigger value="ALL_TIME">Total</TabsTrigger>
                    </TabsList>
                </Tabs>

                {/* Leaderboard Entries */}
                <div className="space-y-2">
                    {loading ? (
                        <div className="text-center py-8 text-muted-foreground">
                            Cargando...
                        </div>
                    ) : entries.length === 0 ? (
                        <div className="text-center py-8 text-muted-foreground">
                            No hay datos aún
                        </div>
                    ) : (
                        entries.map((entry) => (
                            <LeaderboardRow
                                key={entry.user_id}
                                entry={entry}
                                rankStyle={getRankStyle(entry.rank)}
                                rankEmoji={getRankEmoji(entry.rank)}
                            />
                        ))
                    )}
                </div>

                {/* User's Position (if not in top 20) */}
                {user?.id && userContext && !entries.some((e) => e.isCurrentUser) && (
                    <div className="mt-4 pt-4 border-t">
                        <p className="text-xs text-muted-foreground mb-2">Tu posición:</p>
                        <LeaderboardRow
                            entry={{
                                rank: userContext.rank,
                                user_id: user.id,
                                first_name: "Tú",
                                last_name: "",
                                points: 0,
                                isCurrentUser: true,
                            }}
                            rankStyle="bg-primary/10 border-2 border-primary"
                            rankEmoji={null}
                        />
                    </div>
                )}
            </CardContent>
        </Card>
    );
}

function LeaderboardRow({
    entry,
    rankStyle,
    rankEmoji,
}: {
    entry: LeaderboardEntry;
    rankStyle: string;
    rankEmoji: string | null;
}) {
    return (
        <div
            className={cn(
                "flex items-center gap-3 p-2 rounded-lg transition-all",
                entry.isCurrentUser && "ring-2 ring-primary ring-offset-2",
                "hover:bg-muted/50"
            )}
        >
            {/* Rank */}
            <div
                className={cn(
                    "w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold",
                    rankStyle
                )}
            >
                {rankEmoji || entry.rank}
            </div>

            {/* Avatar & Name */}
            <Avatar className="w-8 h-8">
                <AvatarImage src={entry.avatar_url} />
                <AvatarFallback>{entry.first_name.charAt(0).toUpperCase()}</AvatarFallback>
            </Avatar>

            <div className="flex-1 min-w-0">
                <p className={cn("font-medium truncate", entry.isCurrentUser && "text-primary")}>
                    {entry.first_name} {entry.last_name}
                </p>
                {/* Removed missing 'level' property as it is not in the interface */}
            </div>

            {/* Score */}
            <div className="text-right">
                <p className="font-bold">{entry.points.toLocaleString()}</p>
                <p className="text-xs text-muted-foreground">XP</p>
            </div>

            {/* Change Indicator */}
            {entry.change !== undefined && entry.change !== 0 && (
                <div
                    className={cn(
                        "text-xs font-medium",
                        entry.change > 0 ? "text-green-500" : "text-red-500"
                    )}
                >
                    {entry.change > 0 ? `↑${entry.change}` : `↓${Math.abs(entry.change)}`}
                </div>
            )}
        </div>
    );
}
