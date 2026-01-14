"use client"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { cn } from "@/lib/utils"

interface Standing {
    id: string
    team_name: string
    points: number
    played: number
    won: number
    drawn: number
    lost: number
    goals_for: number
    goals_against: number
    goal_difference: number
}

interface StandingsTableProps {
    standings: Standing[]
    loading?: boolean
}

export function PublicStandingsTable({ standings, loading }: StandingsTableProps) {
    if (loading) {
        return <div className="p-8 text-center text-muted-foreground">Cargando posiciones...</div>
    }

    if (standings.length === 0) {
        return <div className="p-8 text-center text-muted-foreground">Aún no hay datos de posiciones.</div>
    }

    return (
        <div className="rounded-md border">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[10px]">#</TableHead>
                        <TableHead>Equipo</TableHead>
                        <TableHead className="text-center">PTS</TableHead>
                        <TableHead className="text-center">PJ</TableHead>
                        <TableHead className="text-center hidden sm:table-cell">PG</TableHead>
                        <TableHead className="text-center hidden sm:table-cell">PE</TableHead>
                        <TableHead className="text-center hidden sm:table-cell">PP</TableHead>
                        <TableHead className="text-center hidden md:table-cell">GF</TableHead>
                        <TableHead className="text-center hidden md:table-cell">GC</TableHead>
                        <TableHead className="text-center hidden md:table-cell">DG</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {standings.map((team, index) => (
                        <TableRow key={team.id}>
                            <TableCell className="font-medium">{index + 1}</TableCell>
                            <TableCell className="font-semibold">{team.team_name}</TableCell>
                            <TableCell className="text-center font-bold bg-muted/30">{team.points}</TableCell>
                            <TableCell className="text-center">{team.played}</TableCell>
                            <TableCell className="text-center hidden sm:table-cell">{team.won}</TableCell>
                            <TableCell className="text-center hidden sm:table-cell">{team.drawn}</TableCell>
                            <TableCell className="text-center hidden sm:table-cell">{team.lost}</TableCell>
                            <TableCell className="text-center hidden md:table-cell">{team.goals_for}</TableCell>
                            <TableCell className="text-center hidden md:table-cell">{team.goals_against}</TableCell>
                            <TableCell className={cn("text-center hidden md:table-cell", team.goal_difference > 0 ? "text-green-600" : team.goal_difference < 0 ? "text-red-600" : "")}>
                                {team.goal_difference > 0 ? "+" : ""}{team.goal_difference}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    )
}
