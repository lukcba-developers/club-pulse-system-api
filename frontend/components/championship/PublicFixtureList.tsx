"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Calendar, Clock } from "lucide-react"

interface Match {
    id: string
    home_team_name: string
    away_team_name: string
    home_score?: number
    away_score?: number
    date: string
    status: string
    home_team_id: string
    away_team_id: string
}

interface FixtureListProps {
    matches: Match[]
    loading?: boolean
}

export function PublicFixtureList({ matches, loading }: FixtureListProps) {
    if (loading) {
        return <div className="p-8 text-center text-muted-foreground">Cargando partidos...</div>
    }

    if (matches.length === 0) {
        return <div className="p-8 text-center text-muted-foreground">No hay partidos programados.</div>
    }

    // Sort by Date
    const sortedMatches = [...matches].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())

    return (
        <div className="space-y-4">
            {sortedMatches.map((match) => {
                const date = new Date(match.date)
                const isCompleted = match.status === 'COMPLETED'

                return (
                    <Card key={match.id} className="overflow-hidden">
                        <CardContent className="p-4 sm:p-6">
                            <div className="flex flex-col md:flex-row justify-between items-center gap-4">
                                {/* Date & Status */}
                                <div className="flex md:flex-col items-center md:items-start gap-2 min-w-[120px]">
                                    <div className="flex items-center text-sm text-muted-foreground">
                                        <Calendar className="mr-2 h-4 w-4" />
                                        {date.toLocaleDateString()}
                                    </div>
                                    <div className="flex items-center text-sm text-muted-foreground">
                                        <Clock className="mr-2 h-4 w-4" />
                                        {date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                    </div>
                                    <Badge variant={isCompleted ? "secondary" : "outline"} className="mt-1">
                                        {isCompleted ? "Finalizado" : "Programado"}
                                    </Badge>
                                </div>

                                {/* Teams & Score */}
                                <div className="flex-grow flex items-center justify-center gap-4 w-full md:w-auto">
                                    <div className="flex-1 text-right font-semibold text-lg truncate min-w-[100px]">
                                        {match.home_team_name || "Local"}
                                    </div>

                                    <div className="flex items-center gap-3 bg-muted px-4 py-2 rounded-lg font-mono text-xl font-bold">
                                        <span>{isCompleted ? match.home_score : "-"}</span>
                                        <span className="text-muted-foreground text-sm">:</span>
                                        <span>{isCompleted ? match.away_score : "-"}</span>
                                    </div>

                                    <div className="flex-1 text-left font-semibold text-lg truncate min-w-[100px]">
                                        {match.away_team_name || "Visitante"}
                                    </div>
                                </div>

                                {/* Action / Info placeholder */}
                                <div className="min-w-[50px] text-right hidden md:block">
                                    {/* Link to detail if needed */}
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                )
            })}
        </div>
    )
}
