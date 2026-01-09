import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { CalendarClock, CheckCircle, Trophy } from 'lucide-react';
import { championshipService, Match } from '@/services/championship-service';
import { MatchSchedulerModal } from './match-scheduler-modal';
import { MatchResultModal } from './match-result-modal';
import { format } from 'date-fns';
import { es } from 'date-fns/locale';

interface FixtureListProps {
    groupId: string;
    clubId: string;
}

export function FixtureList({ groupId, clubId }: FixtureListProps) {
    const [matches, setMatches] = useState<Match[]>([]);
    const [loading, setLoading] = useState(false);
    const [selectedMatch, setSelectedMatch] = useState<Match | null>(null);
    const [selectedResultMatch, setSelectedResultMatch] = useState<Match | null>(null);

    const fetchFixture = async () => {
        setLoading(true);
        try {
            const data = await championshipService.getMatches(groupId);
            setMatches(data);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleGenerateFixture = async () => {
        setLoading(true);
        try {
            await championshipService.generateFixture(groupId);
            await fetchFixture();
        } catch (error) {
            console.error("Error generating fixture:", error);
            // Optionally show toast here if we had access or pass prop
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const fetchFixtureData = async () => {
            setLoading(true);
            try {
                const data = await championshipService.getMatches(groupId);
                setMatches(data);
            } catch (error) {
                console.error(error);
            } finally {
                setLoading(false);
            }
        };

        if (groupId) {
            fetchFixtureData();
        }
    }, [groupId]);

    if (loading && matches.length === 0) {
        return <div className="text-center py-4">Cargando partidos...</div>;
    }

    if (matches.length === 0) {
        return (
            <div className="text-center py-8 text-gray-500 space-y-4">
                <p>Aún no hay partidos generados para este grupo.</p>
                <Button onClick={handleGenerateFixture} disabled={loading} variant="outline">
                    {loading ? "Generando..." : "Generar Fixture Automático"}
                </Button>
            </div>
        );
    }

    // Sort matches: Scheduled first, then TBD? Or Date?
    // Let's sort by Date
    const sortedMatches = [...matches].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

    return (
        <>
            <div className="space-y-4">
                {sortedMatches.map((match) => {
                    // Since we don't have team names in Match object (just IDs), we might need to map them locally
                    // But ideally the API should return names (joined) or we map them from Standings/Teams list.
                    // For now, let's just show IDs or "Equipo Local" placeholders if we can't map.
                    // IMPORTANT: The backend GenerateFixture response only has IDs.
                    // Ideally we should update the backend to Preload Teams or return a DTO with names.
                    // Or, the `standings` we fetched earlier have team names? No, `Standing` has `team_id` but maybe we can fetch team details?
                    // Wait, `Standing` interface in `championship-service.ts` has `TeamID`. It doesn't show names.
                    // BUT `mockStandings` had names.
                    // Real `Standing` struct in backend has `TeamID`. `Team` entity has Name.
                    // I need to fetch team names. Or better, update `GetMatchesByGroup` to return Team objects or names.

                    // To avoid backend refactor overhead right now, I will display "Equipo 1 vs Equipo 2" or assume IDs are readable? No IDs are UUIDs.
                    // I MUST Fetch team names or Update backend.
                    // Updating backend `GetMatchesByGroup` to Preload HomeTeam and AwayTeam is quick.

                    return (
                        <div key={match.id} className="flex items-center justify-between p-4 bg-white dark:bg-zinc-900 border border-gray-100 dark:border-zinc-800 rounded-lg shadow-sm">
                            <div className="flex flex-col gap-1">
                                <div className="flex items-center gap-2 font-medium">
                                    <span className="text-brand-600">
                                        {match.home_team_name || match.home_team_id.substring(0, 8)}
                                        {/* (Local) */}
                                    </span>
                                    <span className="text-xs text-gray-400">vs</span>
                                    <span className="text-red-600">
                                        {match.away_team_name || match.away_team_id.substring(0, 8)}
                                        {/* (Visita) */}
                                    </span>
                                </div>
                                <div className="text-xs text-gray-500 flex items-center gap-1">
                                    <CalendarClock className="w-3 h-3" />
                                    {match.booking_id ? (
                                        <span className="text-green-600 font-medium">
                                            {format(new Date(match.date), "dd MMM yyyy - HH:mm", { locale: es })} (Cancha Reservada)
                                        </span>
                                    ) : (
                                        <span>Sin programar</span>
                                    )}
                                </div>
                            </div>

                            <div className="flex items-center gap-2">
                                {match.home_score !== undefined && match.away_score !== undefined ? (
                                    <div className="flex items-center gap-2 mr-4">
                                        <span className="text-2xl font-bold">{match.home_score}</span>
                                        <span className="text-gray-400">-</span>
                                        <span className="text-2xl font-bold">{match.away_score}</span>
                                        <Button size="sm" variant="ghost" onClick={() => setSelectedResultMatch(match)}>
                                            <Trophy className="w-4 h-4 text-yellow-500" />
                                        </Button>
                                    </div>
                                ) : (
                                    match.booking_id && (
                                        <Button size="sm" variant="outline" onClick={() => setSelectedResultMatch(match)}>
                                            Ingresar Resultado
                                        </Button>
                                    )
                                )}

                                {!match.booking_id && (
                                    <Button size="sm" variant="secondary" onClick={() => setSelectedMatch(match)}>
                                        Programar
                                    </Button>
                                )}
                                {match.booking_id && !match.home_score && (
                                    <div className="text-green-600 flex items-center text-xs">
                                        <CheckCircle className="w-4 h-4 mr-1" /> Programado
                                    </div>
                                )}
                            </div>
                        </div>
                    );
                })}
            </div>

            {selectedMatch && (
                <MatchSchedulerModal
                    isOpen={!!selectedMatch}
                    onClose={() => setSelectedMatch(null)}
                    match={selectedMatch}
                    clubId={clubId}
                    onSuccess={fetchFixture}
                />
            )}

            {selectedResultMatch && (
                <MatchResultModal
                    isOpen={!!selectedResultMatch}
                    onClose={() => setSelectedResultMatch(null)}
                    match={selectedResultMatch}
                    onSuccess={fetchFixture}
                />
            )}
        </>
    );
}
