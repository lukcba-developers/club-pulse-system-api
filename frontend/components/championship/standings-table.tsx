'use client';

import { Standing } from '@/services/championship-service';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Badge } from '@/components/ui/badge';

interface StandingsTableProps {
    standings: Standing[];
}

export function StandingsTable({ standings }: StandingsTableProps) {
    if (!standings || standings.length === 0) {
        return <p className="text-gray-500 text-sm text-center py-4">No hay datos disponibles en esta tabla.</p>;
    }

    // Sort by points desc, goal diff desc
    const sortedStandings = [...standings].sort((a, b) => {
        if (b.points !== a.points) return b.points - a.points;
        return b.goal_difference - a.goal_difference;
    });

    return (
        <div className="rounded-md border">
            <Table>
                <TableHeader>
                    <TableRow className="bg-gray-50 dark:bg-zinc-800">
                        <TableHead className="w-[100px] font-bold">Pos</TableHead>
                        <TableHead>Equipo</TableHead>
                        <TableHead className="text-center">PTS</TableHead>
                        <TableHead className="text-center hidden md:table-cell">PJ</TableHead>
                        <TableHead className="text-center hidden md:table-cell">PG</TableHead>
                        <TableHead className="text-center hidden md:table-cell">PE</TableHead>
                        <TableHead className="text-center hidden md:table-cell">PP</TableHead>
                        <TableHead className="text-center hidden lg:table-cell">GF</TableHead>
                        <TableHead className="text-center hidden lg:table-cell">GC</TableHead>
                        <TableHead className="text-center">DG</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {sortedStandings.map((team, index) => (
                        <TableRow key={team.team_id}>
                            <TableCell className="font-medium">
                                {index + 1}
                                {index === 0 && <span className="ml-2 text-yellow-500">🏆</span>}
                            </TableCell>
                            <TableCell className="font-semibold text-gray-900 dark:text-gray-200">
                                {team.team_name || <span className="font-mono text-xs text-gray-500">{team.team_id.substring(0, 8)}...</span>}
                            </TableCell>
                            <TableCell className="text-center font-bold text-lg">{team.points}</TableCell>
                            <TableCell className="text-center hidden md:table-cell">{team.played}</TableCell>
                            <TableCell className="text-center hidden md:table-cell text-green-600">{team.won}</TableCell>
                            <TableCell className="text-center hidden md:table-cell text-gray-500">{team.drawn}</TableCell>
                            <TableCell className="text-center hidden md:table-cell text-red-500">{team.lost}</TableCell>
                            <TableCell className="text-center hidden lg:table-cell">{team.goals_for}</TableCell>
                            <TableCell className="text-center hidden lg:table-cell">{team.goals_against}</TableCell>
                            <TableCell className="text-center font-semibold">
                                <Badge variant={team.goal_difference > 0 ? "default" : "secondary"}>
                                    {team.goal_difference > 0 ? `+${team.goal_difference}` : team.goal_difference}
                                </Badge>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
}
