'use client';

import { useState } from 'react';
import { useAuth } from '@/hooks/use-auth';
import { TournamentWizard } from '@/components/championship/tournament-wizard';
import { AdminStandingsTable as StandingsTable } from '@/components/championship/AdminStandingsTable';
import { AdminFixtureList as FixtureList } from '@/components/championship/AdminFixtureList';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Trophy, LayoutList, PlusCircle } from 'lucide-react';
import { useToast } from '@/components/ui/use-toast';
import { championshipService, Tournament, Standing } from '@/services/championship-service';

export default function ChampionshipPage() {
    const { user } = useAuth();
    const { toast } = useToast();
    const [view, setView] = useState<'LIST' | 'CREATE'>('LIST');

    // State for data
    const [tournaments, setTournaments] = useState<Tournament[]>([]);
    const [selectedTournament, setSelectedTournament] = useState<Tournament | null>(null);
    const [standings, setStandings] = useState<Standing[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    // Initial Fetch
    const fetchTournaments = async () => {
        if (!user || user.role === 'MEMBER') return; // Or handle members differently
        setLoading(true);
        try {
            const clubId = user.club_id || "default-club-id";
            const data = await championshipService.listTournaments(clubId);
            setTournaments(data);
            if (data.length > 0) {
                // Auto-select first active or first element
                const active = data.find(t => t.status === 'ACTIVE') || data[0];
                setSelectedTournament(active);
                // Try to selected first group of first stage
                if (active.stages && active.stages.length > 0 && active.stages[0].groups && active.stages[0].groups.length > 0) {
                    const firstGroupId = active.stages[0].groups[0].id;
                    setSelectedGroupId(firstGroupId);
                    fetchStandings(firstGroupId);
                }
            }
        } catch (error) {
            console.error(error);
            toast({ title: "Error", description: "No se pudieron cargar los torneos.", variant: "destructive" });
        } finally {
            setLoading(false);
        }
    };

    const fetchStandings = async (groupId: string) => {
        try {
            const data = await championshipService.getStandings(groupId);
            setStandings(data);
        } catch (error) {
            console.error(error);
        }
    };

    // Run effect

    useState(() => {
        fetchTournaments();
    }); // Simple mount effect equivalent logic via useState initializer or separate useEffect. 
    // Ideally use useEffect.
    const [mounted, setMounted] = useState(false);
    if (!mounted && typeof window !== 'undefined') {
        setMounted(true);
        fetchTournaments();
    }

    if (!user) return null;

    return (
        <div className="space-y-6 max-w-7xl mx-auto p-4 md:p-8">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-gray-900 dark:text-white flex items-center gap-2">
                        <Trophy className="w-8 h-8 text-yellow-500" />
                        Campeonatos y Torneos
                    </h1>
                    <p className="text-gray-500 dark:text-gray-400 mt-1">
                        Gestiona ligas, fixtures y tablas de posiciones.
                    </p>
                </div>
                {view === 'LIST' && (user.role === 'ADMIN' || user.role === 'SUPER_ADMIN') && (
                    <button
                        onClick={() => setView('CREATE')}
                        className="bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors shadow-sm"
                    >
                        <PlusCircle className="w-4 h-4" />
                        Nuevo Torneo
                    </button>
                )}
                {view === 'CREATE' && (
                    <button
                        onClick={() => {
                            setView('LIST');
                            fetchTournaments(); // Refresh list on return
                        }}
                        className="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors"
                    >
                        <LayoutList className="w-4 h-4" />
                        Volver al Listado
                    </button>
                )}
            </div>

            {view === 'CREATE' ? (
                <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
                    <TournamentWizard clubId={user.club_id || "default-club-id"} />
                </div>
            ) : (
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-in fade-in duration-500">
                    <div className="lg:col-span-2 space-y-6">
                        {/* Featured Tournament Standings */}
                        {selectedTournament ? (
                            <Card className="border-brand-100 dark:border-zinc-800 shadow-sm">
                                <CardHeader>
                                    <CardTitle className="flex items-center justify-between">
                                        <span>{selectedTournament.name}</span>
                                        <span className={`text-xs font-normal px-2 py-1 rounded-full ${selectedTournament.status === 'ACTIVE' ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'}`}>
                                            {selectedTournament.status === 'ACTIVE' ? 'En Curso' : selectedTournament.status}
                                        </span>
                                    </CardTitle>
                                </CardHeader>
                                <CardContent>
                                    {standings.length > 0 ? (
                                        <StandingsTable standings={standings} />
                                    ) : (
                                        <div className="text-center py-8 text-gray-500">
                                            No hay tablas de posiciones disponibles o no se ha seleccionado un grupo.
                                        </div>
                                    )}
                                </CardContent>
                            </Card>
                        ) : (
                            <Card>
                                <CardContent className="py-8 text-center text-gray-500">
                                    {loading ? 'Cargando torneos...' : 'No hay torneos activos. Crea uno nuevo para comenzar.'}
                                </CardContent>
                            </Card>
                        )}
                    </div>

                    <div className="space-y-6">
                        {/* Tournament Selector (if multiple) */}
                        {tournaments.length > 1 && (
                            <Card>
                                <CardHeader><CardTitle className="text-sm">Mis Torneos</CardTitle></CardHeader>
                                <CardContent className="space-y-2">
                                    {tournaments.map(t => (
                                        <div
                                            key={t.id}
                                            onClick={() => {
                                                setSelectedTournament(t);
                                                if (t.stages?.[0]?.groups?.[0]) {
                                                    const gid = t.stages[0].groups[0].id;
                                                    setSelectedGroupId(gid);
                                                    fetchStandings(gid);
                                                } else {
                                                    setStandings([]);
                                                    setSelectedGroupId(null);
                                                }
                                            }}
                                            className={`cursor-pointer p-2 rounded text-sm ${selectedTournament?.id === t.id ? 'bg-brand-50 text-brand-700' : 'hover:bg-gray-50'}`}
                                        >
                                            {t.name}
                                        </div>
                                    ))}
                                </CardContent>
                            </Card>
                        )}

                        {/* Fixture List */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-lg">Fixture y Resultados</CardTitle>
                            </CardHeader>
                            <CardContent>
                                {selectedGroupId ? (
                                    <FixtureList
                                        groupId={selectedGroupId}
                                        clubId={user.club_id || "default-club-id"}
                                    />
                                ) : (
                                    <div className="text-center py-4 text-gray-500 text-sm">
                                        Selecciona un torneo/grupo.
                                    </div>
                                )}
                            </CardContent>
                        </Card>
                    </div>
                </div>
            )}
        </div>
    );
}
