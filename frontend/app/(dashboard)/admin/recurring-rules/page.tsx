'use client';

import { useState, useEffect } from 'react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Loader2, Calendar, RefreshCw, Plus, List } from 'lucide-react';
import api from '@/lib/axios';
import { bookingService, RecurringRule } from '@/services/booking-service';

interface RecurringRuleDTO {
    facility_id: string;
    type: 'WEEKLY' | 'MONTHLY';
    day_of_week: number;
    start_time: string;
    end_time: string;
    start_date: string;
    end_date: string;
}

const dayLabels = ['Domingo', 'Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado'];

export default function RecurringRulesPage() {
    const [formData, setFormData] = useState<RecurringRuleDTO>({
        facility_id: '',
        type: 'WEEKLY',
        day_of_week: 1,
        start_time: '09:00',
        end_time: '10:00',
        start_date: '',
        end_date: '',
    });
    const [loading, setLoading] = useState(false);
    const [generating, setGenerating] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [rules, setRules] = useState<RecurringRule[]>([]);
    const [loadingRules, setLoadingRules] = useState(true);

    // Fetch existing rules on mount
    useEffect(() => {
        const fetchRules = async () => {
            try {
                const data = await bookingService.listRecurringRules();
                setRules(data || []);
            } catch (err) {
                console.error('Failed to fetch rules', err);
            } finally {
                setLoadingRules(false);
            }
        };
        fetchRules();
    }, []);

    const handleCreateRule = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setMessage('');

        try {
            // Convert time strings to full ISO datetime
            const startTimeISO = new Date(`${formData.start_date}T${formData.start_time}:00`).toISOString();
            const endTimeISO = new Date(`${formData.start_date}T${formData.end_time}:00`).toISOString();

            await api.post('/bookings/recurring', {
                facility_id: formData.facility_id,
                type: formData.type,
                day_of_week: formData.day_of_week,
                start_time: startTimeISO,
                end_time: endTimeISO,
                start_date: formData.start_date,
                end_date: formData.end_date,
            });

            setMessage('Regla recurrente creada exitosamente.');
            // Refresh rules list
            const updatedRules = await bookingService.listRecurringRules();
            setRules(updatedRules || []);
            setFormData({
                facility_id: '',
                type: 'WEEKLY',
                day_of_week: 1,
                start_time: '09:00',
                end_time: '10:00',
                start_date: '',
                end_date: '',
            });
        } catch (err) {
            console.error('Failed to create rule', err);
            setError('Error al crear la regla. Verificá los datos.');
        } finally {
            setLoading(false);
        }
    };

    const handleGenerateBookings = async () => {
        setGenerating(true);
        setError('');
        setMessage('');

        try {
            await api.post('/bookings/generate');
            setMessage('Reservas generadas exitosamente a partir de las reglas activas.');
        } catch (err) {
            console.error('Failed to generate bookings', err);
            setError('Error al generar reservas.');
        } finally {
            setGenerating(false);
        }
    };

    return (
        <div className="space-y-6 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Reglas Recurrentes</h1>
                    <p className="text-muted-foreground">
                        Configura reservas automáticas semanales o mensuales.
                    </p>
                </div>
                <Button onClick={handleGenerateBookings} disabled={generating} variant="outline">
                    {generating ? (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    ) : (
                        <RefreshCw className="mr-2 h-4 w-4" />
                    )}
                    Generar Reservas
                </Button>
            </div>

            {message && (
                <div className="p-3 rounded-lg bg-green-50 text-green-800 border border-green-200">
                    {message}
                </div>
            )}

            {error && (
                <div className="p-3 rounded-lg bg-red-50 text-red-800 border border-red-200">
                    {error}
                </div>
            )}

            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Calendar className="h-5 w-5" />
                        Crear Nueva Regla
                    </CardTitle>
                    <CardDescription>
                        Define un patrón para crear reservas automáticamente.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleCreateRule} className="grid gap-4 md:grid-cols-2">
                        <div className="space-y-2">
                            <Label htmlFor="facility_id">ID de Instalación</Label>
                            <Input
                                id="facility_id"
                                placeholder="UUID de la cancha/sala"
                                value={formData.facility_id}
                                onChange={(e) => setFormData({ ...formData, facility_id: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label>Tipo de Recurrencia</Label>
                            <Select
                                value={formData.type}
                                onValueChange={(v) => setFormData({ ...formData, type: v as 'WEEKLY' | 'MONTHLY' })}
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="WEEKLY">Semanal</SelectItem>
                                    <SelectItem value="MONTHLY">Mensual</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <Label>Día de la Semana</Label>
                            <Select
                                value={formData.day_of_week.toString()}
                                onValueChange={(v) => setFormData({ ...formData, day_of_week: parseInt(v) })}
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {dayLabels.map((day, i) => (
                                        <SelectItem key={i} value={i.toString()}>{day}</SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="start_time">Hora de Inicio</Label>
                            <Input
                                id="start_time"
                                type="time"
                                value={formData.start_time}
                                onChange={(e) => setFormData({ ...formData, start_time: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="end_time">Hora de Fin</Label>
                            <Input
                                id="end_time"
                                type="time"
                                value={formData.end_time}
                                onChange={(e) => setFormData({ ...formData, end_time: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="start_date">Fecha de Inicio</Label>
                            <Input
                                id="start_date"
                                type="date"
                                value={formData.start_date}
                                onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="end_date">Fecha de Fin</Label>
                            <Input
                                id="end_date"
                                type="date"
                                value={formData.end_date}
                                onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                                required
                            />
                        </div>

                        <div className="md:col-span-2 flex justify-end pt-4">
                            <Button type="submit" disabled={loading}>
                                {loading ? (
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                ) : (
                                    <Plus className="mr-2 h-4 w-4" />
                                )}
                                Crear Regla
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>

            {/* Existing Rules List */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <List className="h-5 w-5" />
                        Reglas Activas
                    </CardTitle>
                    <CardDescription>
                        Reglas recurrentes configuradas actualmente.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {loadingRules ? (
                        <div className="flex items-center justify-center py-8">
                            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                        </div>
                    ) : rules.length === 0 ? (
                        <div className="text-center py-8 text-muted-foreground">
                            No hay reglas recurrentes configuradas.
                        </div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="border-b">
                                        <th className="text-left py-2 px-2 font-medium">Instalación</th>
                                        <th className="text-left py-2 px-2 font-medium">Día</th>
                                        <th className="text-left py-2 px-2 font-medium">Horario</th>
                                        <th className="text-left py-2 px-2 font-medium">Período</th>
                                        <th className="text-left py-2 px-2 font-medium">Tipo</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {rules.map((rule) => (
                                        <tr key={rule.id} className="border-b last:border-0 hover:bg-muted/50">
                                            <td className="py-2 px-2 font-mono text-xs truncate max-w-[120px]" title={rule.facility_id}>
                                                {rule.facility_id.slice(0, 8)}...
                                            </td>
                                            <td className="py-2 px-2">{dayLabels[rule.day_of_week]}</td>
                                            <td className="py-2 px-2">
                                                {rule.start_time.slice(0, 5)} - {rule.end_time.slice(0, 5)}
                                            </td>
                                            <td className="py-2 px-2 text-xs text-muted-foreground">
                                                {rule.start_date} → {rule.end_date}
                                            </td>
                                            <td className="py-2 px-2">
                                                <span className="px-2 py-0.5 rounded-full text-xs bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300">
                                                    {rule.type}
                                                </span>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
