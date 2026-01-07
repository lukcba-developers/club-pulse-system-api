import { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { facilityService, Facility } from '@/services/facility-service';
import { championshipService, Match } from '@/services/championship-service';
// import { format } from 'date-fns';

interface MatchSchedulerModalProps {
    isOpen: boolean;
    onClose: () => void;
    match: Match;
    clubId: string;
    onSuccess: () => void;
}

export function MatchSchedulerModal({ isOpen, onClose, match, clubId, onSuccess }: MatchSchedulerModalProps) {
    const { toast } = useToast();
    const [loading, setLoading] = useState(false);
    const [facilities, setFacilities] = useState<Facility[]>([]);

    // Form state
    const [courtId, setCourtId] = useState<string>('');
    const [date, setDate] = useState<string>('');
    const [startTime, setStartTime] = useState<string>('');
    const [duration, setDuration] = useState<string>('60'); // Minutes

    useEffect(() => {
        if (isOpen) {
            fetchFacilities();
        }
    }, [isOpen]);

    const fetchFacilities = async () => {
        try {
            const data = await facilityService.list(100); // Fetch enough
            setFacilities(data);
        } catch (error) {
            console.error(error);
            toast({
                title: "Error",
                description: "No se pudieron cargar las canchas.",
                variant: "destructive",
            });
        }
    };

    const handleSchedule = async () => {
        if (!courtId || !date || !startTime) {
            toast({
                title: "Faltan datos",
                description: "Por favor completa todos los campos requeridos.",
                variant: "destructive",
            });
            return;
        }

        setLoading(true);
        try {
            // Construct start and end times
            // Date is YYYY-MM-DD, StartTime is HH:MM
            const startDateTime = new Date(`${date}T${startTime}`);
            const endDateTime = new Date(startDateTime.getTime() + parseInt(duration) * 60000);

            await championshipService.scheduleMatch({
                club_id: clubId,
                match_id: match.id,
                court_id: courtId,
                start_time: startDateTime.toISOString(),
                end_time: endDateTime.toISOString(),
            });

            toast({
                title: "Partido Programado",
                description: "La reserva de cancha se ha creado exitosamente.",
            });
            onSuccess();
            onClose();
        } catch (error) {
            console.error(error);
            toast({
                title: "Error",
                description: "No se pudo programar el partido. Verifica la disponibilidad.",
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
                    <DialogTitle>Programar Partido</DialogTitle>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="court" className="text-right">
                            Cancha
                        </Label>
                        <Select onValueChange={setCourtId} value={courtId}>
                            <SelectTrigger className="w-[280px]">
                                <SelectValue placeholder="Seleccionar cancha" />
                            </SelectTrigger>
                            <SelectContent>
                                {facilities.map((facility) => (
                                    <SelectItem key={facility.id} value={facility.id}>
                                        {facility.name} ({facility.location.name})
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="date" className="text-right">
                            Fecha
                        </Label>
                        <Input
                            id="date"
                            type="date"
                            className="col-span-3"
                            value={date}
                            onChange={(e) => setDate(e.target.value)}
                        />
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="time" className="text-right">
                            Hora
                        </Label>
                        <Input
                            id="time"
                            type="time"
                            className="col-span-3"
                            value={startTime}
                            onChange={(e) => setStartTime(e.target.value)}
                        />
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="duration" className="text-right">
                            Duración
                        </Label>
                        <Select onValueChange={setDuration} value={duration}>
                            <SelectTrigger className="w-[280px]">
                                <SelectValue placeholder="Duración" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="60">60 min</SelectItem>
                                <SelectItem value="90">90 min</SelectItem>
                                <SelectItem value="120">120 min</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={onClose}>Cancelar</Button>
                    <Button onClick={handleSchedule} disabled={loading}>
                        {loading ? "Programando..." : "Confirmar"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
