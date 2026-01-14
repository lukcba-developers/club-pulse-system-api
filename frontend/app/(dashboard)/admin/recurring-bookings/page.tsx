'use client';

import { useState, useEffect, useCallback } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
    bookingService,
    RecurringRule,
    CreateRecurringRuleDTO
} from '@/services/booking-service';
import { facilityService, Facility } from '@/services/facility-service';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue
} from '@/components/ui/select';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from '@/components/ui/dialog';
import {
    Plus,
    Calendar,
    Repeat,
    Loader2,
    Play,
    Trash2
} from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

const DAYS_OF_WEEK = [
    { value: 0, label: 'Domingo' },
    { value: 1, label: 'Lunes' },
    { value: 2, label: 'Martes' },
    { value: 3, label: 'Miércoles' },
    { value: 4, label: 'Jueves' },
    { value: 5, label: 'Viernes' },
    { value: 6, label: 'Sábado' },
];

// Form schema
const ruleSchema = z.object({
    facility_id: z.string().min(1, 'Selecciona una instalación'),
    frequency: z.enum(['WEEKLY', 'MONTHLY']),
    day_of_week: z.coerce.number().min(0).max(6),
    start_time: z.string().regex(/^\d{2}:\d{2}$/, 'Formato HH:MM'),
    end_time: z.string().regex(/^\d{2}:\d{2}$/, 'Formato HH:MM'),
    start_date: z.string().min(1, 'Fecha requerida'),
    end_date: z.string().min(1, 'Fecha requerida'),
});

type RuleFormValues = z.infer<typeof ruleSchema>;

export default function RecurringBookingsPage() {
    const [rules, setRules] = useState<RecurringRule[]>([]);
    const [facilities, setFacilities] = useState<Facility[]>([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [generating, setGenerating] = useState(false);
    const { toast } = useToast();

    const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<RuleFormValues>({
        resolver: zodResolver(ruleSchema),
        defaultValues: {
            frequency: 'WEEKLY',
            day_of_week: 1,
            start_time: '09:00',
            end_time: '10:00',
        }
    });

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const [facilitiesData] = await Promise.all([
                facilityService.list(),
                // Rules endpoint might not be implemented yet
                // bookingService.listRecurringRules().catch(() => []),
            ]);
            setFacilities(facilitiesData || []);
            // setRules(rulesData || []);
        } catch (error) {
            console.error('Failed to load data', error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const onSubmit = async (data: RuleFormValues) => {
        setSubmitting(true);
        try {
            const dto: CreateRecurringRuleDTO = {
                ...data,
                type: 'FIXED', // Default purpose
                start_time: data.start_time + ':00', // Add seconds
                end_time: data.end_time + ':00',
            };
            const newRule = await bookingService.createRecurringRule(dto);
            setRules(prev => [...prev, newRule]);
            setIsModalOpen(false);
            reset();
            toast({
                title: 'Regla creada',
                description: 'La regla de reserva recurrente fue creada exitosamente.',
            });
        } catch (error) {
            console.error('Failed to create rule', error);
            toast({
                title: 'Error',
                description: 'No se pudo crear la regla. Intenta de nuevo.',
                variant: 'destructive',
            });
        } finally {
            setSubmitting(false);
        }
    };

    const handleGenerate = async () => {
        setGenerating(true);
        try {
            await bookingService.generateFromRules(4);
            toast({
                title: 'Reservas generadas',
                description: 'Se generaron las reservas para las próximas 4 semanas.',
            });
        } catch (error) {
            console.error('Failed to generate bookings', error);
            toast({
                title: 'Error',
                description: 'No se pudieron generar las reservas.',
                variant: 'destructive',
            });
        } finally {
            setGenerating(false);
        }
    };

    const getFacilityName = (id: string) => {
        return facilities.find(f => f.id === id)?.name || 'Desconocida';
    };

    return (
        <div className="space-y-6 max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="flex justify-between items-center flex-wrap gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Reservas Recurrentes</h1>
                    <p className="text-muted-foreground">Gestiona reglas de reservas automáticas para clases, entrenamientos o eventos.</p>
                </div>
                <div className="flex gap-2">
                    <Button
                        variant="outline"
                        onClick={handleGenerate}
                        disabled={generating || rules.length === 0}
                    >
                        {generating ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        ) : (
                            <Play className="mr-2 h-4 w-4" />
                        )}
                        Generar Reservas
                    </Button>
                    <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
                        <DialogTrigger asChild>
                            <Button>
                                <Plus className="mr-2 h-4 w-4" />
                                Nueva Regla
                            </Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-[520px]">
                            <DialogHeader>
                                <DialogTitle>Crear Regla Recurrente</DialogTitle>
                                <DialogDescription>
                                    Define una regla para crear reservas automáticamente en un horario fijo.
                                </DialogDescription>
                            </DialogHeader>
                            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 py-4">
                                {/* Facility */}
                                <div className="space-y-2">
                                    <Label>Instalación</Label>
                                    <Select
                                        onValueChange={(v) => setValue('facility_id', v)}
                                        defaultValue=""
                                    >
                                        <SelectTrigger>
                                            <SelectValue placeholder="Selecciona una instalación..." />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {facilities.map(f => (
                                                <SelectItem key={f.id} value={f.id}>{f.name}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                    {errors.facility_id && <p className="text-xs text-red-500">{errors.facility_id.message}</p>}
                                </div>

                                {/* Frequency */}
                                <div className="space-y-2">
                                    <Label>Frecuencia</Label>
                                    <Select
                                        onValueChange={(v) => setValue('frequency', v as 'WEEKLY' | 'MONTHLY')}
                                        defaultValue="WEEKLY"
                                    >
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="WEEKLY">Semanal</SelectItem>
                                            <SelectItem value="MONTHLY">Mensual</SelectItem>
                                        </SelectContent>
                                    </Select>
                                    {errors.frequency && <p className="text-xs text-red-500">{errors.frequency.message}</p>}
                                </div>

                                {/* Day of Week */}
                                <div className="space-y-2">
                                    <Label>Día de la Semana</Label>
                                    <Select
                                        onValueChange={(v) => setValue('day_of_week', parseInt(v))}
                                        defaultValue="1"
                                    >
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {DAYS_OF_WEEK.map(d => (
                                                <SelectItem key={d.value} value={d.value.toString()}>{d.label}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>

                                {/* Time Range */}
                                <div className="grid grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <Label htmlFor="start_time">Hora Inicio</Label>
                                        <Input
                                            id="start_time"
                                            type="time"
                                            {...register('start_time')}
                                        />
                                        {errors.start_time && <p className="text-xs text-red-500">{errors.start_time.message}</p>}
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="end_time">Hora Fin</Label>
                                        <Input
                                            id="end_time"
                                            type="time"
                                            {...register('end_time')}
                                        />
                                        {errors.end_time && <p className="text-xs text-red-500">{errors.end_time.message}</p>}
                                    </div>
                                </div>

                                {/* Date Range */}
                                <div className="grid grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <Label htmlFor="start_date">Fecha Inicio</Label>
                                        <Input
                                            id="start_date"
                                            type="date"
                                            {...register('start_date')}
                                        />
                                        {errors.start_date && <p className="text-xs text-red-500">{errors.start_date.message}</p>}
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="end_date">Fecha Fin</Label>
                                        <Input
                                            id="end_date"
                                            type="date"
                                            {...register('end_date')}
                                        />
                                        {errors.end_date && <p className="text-xs text-red-500">{errors.end_date.message}</p>}
                                    </div>
                                </div>

                                <DialogFooter className="pt-4">
                                    <Button type="button" variant="outline" onClick={() => setIsModalOpen(false)}>
                                        Cancelar
                                    </Button>
                                    <Button type="submit" disabled={submitting}>
                                        {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                        Crear Regla
                                    </Button>
                                </DialogFooter>
                            </form>
                        </DialogContent>
                    </Dialog>
                </div>
            </div>

            {/* Rules List */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Repeat className="h-5 w-5" />
                        Reglas Activas
                    </CardTitle>
                    <CardDescription>Lista de reglas de reservas recurrentes configuradas.</CardDescription>
                </CardHeader>
                <CardContent>
                    {loading ? (
                        <div className="flex justify-center py-10">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    ) : rules.length === 0 ? (
                        <div className="text-center py-12 border border-dashed rounded-lg">
                            <Calendar className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                            <h3 className="text-lg font-medium text-gray-900 dark:text-white">Sin reglas configuradas</h3>
                            <p className="text-gray-500 text-sm mt-1">Crea una regla para automatizar reservas repetitivas.</p>
                            <Button className="mt-4" onClick={() => setIsModalOpen(true)}>
                                <Plus className="mr-2 h-4 w-4" />
                                Crear Primera Regla
                            </Button>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {rules.map(rule => (
                                <div
                                    key={rule.id}
                                    className="flex items-center justify-between p-4 bg-muted/50 rounded-lg border"
                                >
                                    <div className="flex items-center gap-4">
                                        <div className="bg-brand-100 dark:bg-brand-900/20 p-2 rounded-lg">
                                            <Repeat className="h-5 w-5 text-brand-600" />
                                        </div>
                                        <div>
                                            <p className="font-medium">{getFacilityName(rule.facility_id)}</p>
                                            <p className="text-sm text-muted-foreground">
                                                {DAYS_OF_WEEK.find(d => d.value === rule.day_of_week)?.label} •
                                                {rule.start_time} - {rule.end_time}
                                            </p>
                                            <p className="text-xs text-muted-foreground">
                                                {rule.start_date} → {rule.end_date}
                                            </p>
                                        </div>
                                    </div>
                                    <Button variant="ghost" size="icon" className="text-red-500 hover:text-red-600 hover:bg-red-50">
                                        <Trash2 className="h-4 w-4" />
                                    </Button>
                                </div>
                            ))}
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
