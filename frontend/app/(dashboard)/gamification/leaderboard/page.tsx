"use client"

import { useEffect, useState } from "react"
import { gamificationService, Leaderboard } from "@/services/gamification-service"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Trophy, Medal, Crown } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

export default function LeaderboardPage() {
    const [leaderboard, setLeaderboard] = useState<Leaderboard | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchLeaderboard = async () => {
            try {
                const data = await gamificationService.getGlobalLeaderboard('MONTHLY', 50)
                setLeaderboard(data)
            } catch (error) {
                console.error("Failed to fetch leaderboard", error)
            } finally {
                setLoading(false)
            }
        }
        fetchLeaderboard()
    }, [])

    if (loading) {
        return <div className="p-8 text-center">Cargando ranking...</div>
    }

    if (!leaderboard) {
        return <div className="p-8 text-center">No se pudo cargar el ranking.</div>
    }

    const getRankIcon = (rank: number) => {
        switch (rank) {
            case 1: return <Crown className="w-6 h-6 text-yellow-500 fill-yellow-500" />
            case 2: return <Medal className="w-6 h-6 text-gray-400 fill-gray-400" />
            case 3: return <Medal className="w-6 h-6 text-amber-700 fill-amber-700" />
            default: return <span className="font-bold text-gray-500 w-6 text-center">{rank}</span>
        }
    }

    return (
        <div className="container mx-auto py-6">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
                <div>
                    <h1 className="text-3xl font-bold flex items-center gap-2">
                        <Trophy className="w-8 h-8 text-primary" />
                        Ranking Global
                    </h1>
                    <p className="text-muted-foreground mt-1">Los miembros más activos del club este mes.</p>
                </div>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Top 50 del Mes</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        {leaderboard.entries.map((entry) => (
                            <div key={entry.user_id} className="flex items-center justify-between p-4 bg-gray-50/50 hover:bg-gray-100 rounded-lg transition-colors border">
                                <div className="flex items-center gap-4">
                                    <div className="flex items-center justify-center w-8">
                                        {getRankIcon(entry.rank)}
                                    </div>
                                    <Avatar>
                                        <AvatarImage src={entry.avatar_url} />
                                        <AvatarFallback>{entry.user_name.substring(0, 2).toUpperCase()}</AvatarFallback>
                                    </Avatar>
                                    <div>
                                        <div className="font-semibold text-lg">{entry.user_name}</div>
                                        <div className="text-xs text-muted-foreground uppercase">Miembro</div>
                                    </div>
                                </div>
                                <div className="text-right">
                                    <div className="text-2xl font-bold text-primary">{entry.score.toLocaleString()}</div>
                                    <div className="text-xs text-muted-foreground font-medium">Puntos de XP</div>
                                </div>
                            </div>
                        ))}
                        {leaderboard.entries.length === 0 && (
                            <div className="text-center py-8 text-gray-500">
                                Aún no hay datos para este ranking.
                            </div>
                        )}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
