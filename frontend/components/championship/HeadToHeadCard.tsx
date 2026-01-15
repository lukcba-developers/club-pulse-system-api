'use client';

import { useState, useEffect } from 'react';
import { championshipService, HeadToHeadResult, Match } from '@/services/championship-service';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Loader2, Trophy, Minus, Target } from 'lucide-react';

interface HeadToHeadCardProps {
    groupId: string;
    teamAId: string;
    teamBId: string;
    teamAName?: string;
    teamBName?: string;
    onClose?: () => void;
}

export function HeadToHeadCard({
    groupId,
    teamAId,
    teamBId,
    teamAName = 'Equipo A',
    teamBName = 'Equipo B',
    onClose
}: HeadToHeadCardProps) {
    const [data, setData] = useState<HeadToHeadResult | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                const result = await championshipService.getHeadToHead(groupId, teamAId, teamBId);
                setData(result);
            } catch (err) {
                console.error('Error fetching head-to-head:', err);
                setError('No se pudo cargar el historial de enfrentamientos');
            } finally {
                setLoading(false);
            }
        };

        if (groupId && teamAId && teamBId) {
            fetchData();
        }
    }, [groupId, teamAId, teamBId]);

    if (loading) {
        return (
            <Card className="w-full max-w-lg">
                <CardContent className="flex justify-center items-center py-12">
                    <Loader2 className="h-8 w-8 animate-spin text-brand-500" />
                </CardContent>
            </Card>
        );
    }

    if (error || !data) {
        return (
            <Card className="w-full max-w-lg">
                <CardContent className="py-8 text-center text-red-500">
                    {error || 'Sin datos disponibles'}
                </CardContent>
            </Card>
        );
    }

    const totalMatches = data.team_a_wins + data.team_b_wins + data.draws;

    return (
        <Card className="w-full max-w-lg overflow-hidden">
            <CardHeader className="bg-gradient-to-r from-brand-600 to-brand-700 text-white">
                <div className="flex justify-between items-center">
                    <CardTitle className="text-lg">⚔️ Enfrentamientos Directos</CardTitle>
                    {onClose && (
                        <button
                            onClick={onClose}
                            className="text-white/70 hover:text-white transition-colors"
                        >
                            ✕
                        </button>
                    )}
                </div>
                <CardDescription className="text-white/80">
                    {teamAName} vs {teamBName}
                </CardDescription>
            </CardHeader>

            <CardContent className="p-6">
                {/* Stats Summary */}
                <div className="grid grid-cols-3 gap-4 mb-6">
                    {/* Team A Wins */}
                    <div className="text-center p-4 bg-green-50 dark:bg-green-900/20 rounded-xl">
                        <div className="text-3xl font-bold text-green-600 dark:text-green-400">
                            {data.team_a_wins}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400 mt-1 truncate" title={teamAName}>
                            {teamAName}
                        </div>
                    </div>

                    {/* Draws */}
                    <div className="text-center p-4 bg-gray-50 dark:bg-gray-800/50 rounded-xl">
                        <div className="text-3xl font-bold text-gray-600 dark:text-gray-400">
                            {data.draws}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                            Empates
                        </div>
                    </div>

                    {/* Team B Wins */}
                    <div className="text-center p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl">
                        <div className="text-3xl font-bold text-blue-600 dark:text-blue-400">
                            {data.team_b_wins}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400 mt-1 truncate" title={teamBName}>
                            {teamBName}
                        </div>
                    </div>
                </div>

                {/* Goals */}
                <div className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800/30 rounded-xl mb-6">
                    <div className="flex items-center gap-2">
                        <Target className="h-5 w-5 text-gray-400" />
                        <span className="text-sm text-gray-600 dark:text-gray-400">Goles</span>
                    </div>
                    <div className="flex items-center gap-3">
                        <span className="text-xl font-bold text-green-600">{data.team_a_goals}</span>
                        <Minus className="h-4 w-4 text-gray-400" />
                        <span className="text-xl font-bold text-blue-600">{data.team_b_goals}</span>
                    </div>
                </div>

                {/* Match History */}
                {data.matches.length > 0 ? (
                    <div>
                        <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
                            Historial ({totalMatches} partidos)
                        </h4>
                        <div className="space-y-2 max-h-48 overflow-y-auto">
                            {data.matches.map((match: Match) => (
                                <div
                                    key={match.id}
                                    className="flex items-center justify-between p-3 bg-white dark:bg-gray-800 rounded-lg border border-gray-100 dark:border-gray-700"
                                >
                                    <div className="text-xs text-gray-500">
                                        {new Date(match.date).toLocaleDateString('es-AR', {
                                            day: '2-digit',
                                            month: 'short',
                                            year: 'numeric'
                                        })}
                                    </div>
                                    <div className="flex items-center gap-2 font-mono text-sm">
                                        <span className={match.home_score !== null && match.away_score !== null && match.home_score > match.away_score ? 'font-bold text-green-600' : ''}>
                                            {match.home_score ?? '-'}
                                        </span>
                                        <span className="text-gray-400">-</span>
                                        <span className={match.home_score !== null && match.away_score !== null && match.away_score > match.home_score ? 'font-bold text-blue-600' : ''}>
                                            {match.away_score ?? '-'}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                ) : (
                    <div className="text-center py-4 text-gray-500 dark:text-gray-400">
                        <Trophy className="h-8 w-8 mx-auto mb-2 opacity-30" />
                        <p className="text-sm">No hay partidos jugados entre estos equipos</p>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}
