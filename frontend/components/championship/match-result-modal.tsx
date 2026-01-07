import { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { championshipService, Match } from '@/services/championship-service';

interface MatchResultModalProps {
    isOpen: boolean;
    onClose: () => void;
    match: Match;
    onSuccess: () => void;
}

export function MatchResultModal({ isOpen, onClose, match, onSuccess }: MatchResultModalProps) {
    const { toast } = useToast();
    const [loading, setLoading] = useState(false);
    const [homeScore, setHomeScore] = useState<string>(match.home_score?.toString() || '');
    const [awayScore, setAwayScore] = useState<string>(match.away_score?.toString() || '');

    const handleSave = async () => {
        if (homeScore === '' || awayScore === '') {
            toast({
                title: "Error",
                description: "Por favor ingresa ambos marcadores.",
                variant: "destructive",
            });
            return;
        }

        setLoading(true);
        try {
            await championshipService.updateMatchResult(
                match.id,
                parseInt(homeScore),
                parseInt(awayScore)
            );

            toast({
                title: "Resultado Actualizado",
                description: "El marcador ha sido guardado exitosamente.",
            });
            onSuccess();
            onClose();
        } catch (error) {
            console.error(error);
            toast({
                title: "Error",
                description: "No se pudo actualizar el resultado.",
                variant: "destructive",
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Ingresar Resultado</DialogTitle>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                    <div className="flex items-center justify-between gap-4">
                        <div className="flex flex-col items-center gap-2">
                            <Label>Local</Label>
                            <div className="text-sm font-bold text-gray-500">{match.home_team_id.substring(0, 6)}...</div>
                            <Input
                                type="number"
                                className="w-20 text-center text-lg font-bold"
                                value={homeScore}
                                onChange={(e) => setHomeScore(e.target.value)}
                            />
                        </div>
                        <span className="text-2xl font-bold text-gray-300">-</span>
                        <div className="flex flex-col items-center gap-2">
                            <Label>Visita</Label>
                            <div className="text-sm font-bold text-gray-500">{match.away_team_id.substring(0, 6)}...</div>
                            <Input
                                type="number"
                                className="w-20 text-center text-lg font-bold"
                                value={awayScore}
                                onChange={(e) => setAwayScore(e.target.value)}
                            />
                        </div>
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={onClose}>Cancelar</Button>
                    <Button onClick={handleSave} disabled={loading}>
                        {loading ? "Guardando..." : "Guardar Resultado"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
