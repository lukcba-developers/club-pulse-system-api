'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { championshipService } from '@/services/championship-service';
import { useToast } from '@/components/ui/use-toast';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
// import { Label } from '@/components/ui/label'; // Ensure this exists or use standard label
// import { Button } from '@/components/ui/button'; // Standard HTML button for now to avoid dependency issues if missing
import { Trophy, ArrowRight, Save, Plus } from 'lucide-react';
// import { clubService } from "@/services/club-service"; // Assuming we have this to get ClubID? Or user/auth context.

export function TournamentWizard({ clubId }: { clubId: string }) {
    const router = useRouter();
    const { toast } = useToast();
    const [step, setStep] = useState(1);
    const [isLoading, setIsLoading] = useState(false);

    // Tournament Data
    const [tournamentData, setTournamentData] = useState({
        name: '',
        description: '',
        sport: 'FUTBOL',
        category: 'Libre',
        start_date: '',
    });

    const [createdTournamentId, setCreatedTournamentId] = useState<string | null>(null);

    const handleCreateTournament = async () => {
        setIsLoading(true);
        try {
            const result = await championshipService.createTournament({
                ...tournamentData,
                club_id: clubId,
                status: 'DRAFT',
                start_date: new Date(tournamentData.start_date).toISOString(),
            });
            setCreatedTournamentId(result.id);
            toast({
                title: "Torneo Creado",
                description: "El torneo se ha guardado como borrador.",
            });
            setStep(2);
        } catch {
            toast({
                variant: "destructive",
                title: "Error",
                description: "No se pudo crear el torneo.",
            });
        } finally {
            setIsLoading(false);
        }
    };

    const handleConfigureStages = async () => {
        if (!createdTournamentId) return;
        setIsLoading(true);
        try {
            // Default configuration: 1 Group Stage
            const stage = await championshipService.addStage(createdTournamentId, {
                name: "Fase de Grupos",
                type: "GROUP",
                order: 1
            });

            // Add default Group A
            await championshipService.addGroup(stage.id, {
                name: "Grupo A"
            });

            toast({
                title: "Fase de Grupos Configurada",
                description: "Se ha creado la Fase de Grupos y el Grupo A por defecto.",
            });

            // Finished for MVP Wizard
            // router.push(`/championships/${createdTournamentId}`); // Future route
            router.refresh();
        } catch {
            toast({
                variant: "destructive",
                title: "Error",
                description: "No se pudo configurar las fases.",
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Card className="w-full max-w-2xl mx-auto border-brand-100 shadow-lg">
            <CardHeader className="bg-gradient-to-r from-brand-50 to-white dark:from-brand-900/20 dark:to-zinc-900 border-b">
                <CardTitle className="flex items-center gap-2 text-brand-700 dark:text-brand-400">
                    <Trophy className="w-6 h-6" />
                    Nuevo Torneo - Paso {step}
                </CardTitle>
            </CardHeader>
            <CardContent className="pt-6 space-y-4">
                {step === 1 && (
                    <div className="space-y-4">
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Nombre del Torneo</label>
                            <Input
                                placeholder="Ej: Torneo Verano 2026"
                                value={tournamentData.name}
                                onChange={(e) => setTournamentData({ ...tournamentData, name: e.target.value })}
                            />
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <label className="text-sm font-medium">Deporte</label>
                                <Input
                                    value={tournamentData.sport}
                                    onChange={(e) => setTournamentData({ ...tournamentData, sport: e.target.value })}
                                />
                            </div>
                            <div className="space-y-2">
                                <label className="text-sm font-medium">Categoría</label>
                                <Input
                                    placeholder="Ej: Veteranos"
                                    value={tournamentData.category}
                                    onChange={(e) => setTournamentData({ ...tournamentData, category: e.target.value })}
                                />
                            </div>
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Fecha de Inicio</label>
                            <Input
                                type="date"
                                value={tournamentData.start_date}
                                onChange={(e) => setTournamentData({ ...tournamentData, start_date: e.target.value })}
                            />
                        </div>
                    </div>
                )}

                {step === 2 && (
                    <div className="text-center py-8">
                        <div className="mx-auto w-12 h-12 bg-green-100 text-green-600 rounded-full flex items-center justify-center mb-4">
                            <Save className="w-6 h-6" />
                        </div>
                        <h3 className="text-lg font-medium text-gray-900 dark:text-white">¡Torneo Inicializado!</h3>
                        <p className="text-gray-500 mb-6">
                            Ahora vamos a configurar la estructura básica.
                            El sistema creará una &quot;Fase de Grupos&quot; automática para empezar.
                        </p>
                    </div>
                )}
            </CardContent>
            <CardFooter className="flex justify-between border-t pt-6 bg-gray-50 dark:bg-zinc-800/50">
                {step === 1 ? (
                    <button
                        onClick={handleCreateTournament}
                        disabled={isLoading || !tournamentData.name}
                        className="ml-auto bg-brand-600 text-white px-4 py-2 rounded-md hover:bg-brand-700 disabled:opacity-50 flex items-center gap-2"
                    >
                        {isLoading ? 'Guardando...' : 'Siguiente'} <ArrowRight className="w-4 h-4" />
                    </button>
                ) : (
                    <button
                        onClick={handleConfigureStages}
                        disabled={isLoading}
                        className="ml-auto bg-brand-600 text-white px-4 py-2 rounded-md hover:bg-brand-700 disabled:opacity-50 flex items-center gap-2"
                    >
                        {isLoading ? 'Configurando...' : 'Crear Fase de Grupos'} <Plus className="w-4 h-4" />
                    </button>
                )}
            </CardFooter>
        </Card>
    );
}
