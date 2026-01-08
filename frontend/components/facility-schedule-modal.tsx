'use client';

import { useState } from 'react';
import { Clock, Save, Loader2 } from 'lucide-react';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '@/components/ui/dialog';
import { facilityService } from '@/services/facility-service';

interface FacilityScheduleModalProps {
    isOpen: boolean;
    onClose: () => void;
    facilityId: string;
    facilityName: string;
    currentOpeningHour: number;
    currentClosingHour: number;
    onSuccess?: () => void;
}

// Generate hours array (0-23)
const hours = Array.from({ length: 24 }, (_, i) => i);

// Format hour for display (e.g., "08:00", "14:00")
const formatHour = (hour: number) => {
    return `${hour.toString().padStart(2, '0')}:00`;
};

export function FacilityScheduleModal({
    isOpen,
    onClose,
    facilityId,
    facilityName,
    currentOpeningHour,
    currentClosingHour,
    onSuccess,
}: FacilityScheduleModalProps) {
    const [openingHour, setOpeningHour] = useState(currentOpeningHour);
    const [closingHour, setClosingHour] = useState(currentClosingHour);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSave = async () => {
        // Validation
        if (openingHour >= closingHour) {
            setError('La hora de apertura debe ser anterior a la hora de cierre');
            return;
        }

        setLoading(true);
        setError('');

        try {
            await facilityService.updateSchedule(facilityId, openingHour, closingHour);
            onSuccess?.();
            onClose();
        } catch (err) {
            console.error('Error updating schedule:', err);
            setError('Error al actualizar horarios. Por favor, intente de nuevo.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-md bg-white dark:bg-zinc-900">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2 text-gray-900 dark:text-gray-100">
                        <Clock className="h-5 w-5 text-brand-500" />
                        Configurar Horarios
                    </DialogTitle>
                    <DialogDescription className="text-gray-500 dark:text-gray-400">
                        Configure las horas de apertura y cierre para <span className="font-medium">{facilityName}</span>
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-6 py-4">
                    {/* Opening Hour */}
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                            Hora de Apertura
                        </label>
                        <select
                            value={openingHour}
                            onChange={(e) => setOpeningHour(Number(e.target.value))}
                            className="w-full px-3 py-2 bg-white dark:bg-zinc-800 border border-gray-300 dark:border-zinc-700 rounded-lg text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-brand-500 focus:border-transparent transition-all"
                        >
                            {hours.map((hour) => (
                                <option key={`open-${hour}`} value={hour}>
                                    {formatHour(hour)}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Closing Hour */}
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                            Hora de Cierre
                        </label>
                        <select
                            value={closingHour}
                            onChange={(e) => setClosingHour(Number(e.target.value))}
                            className="w-full px-3 py-2 bg-white dark:bg-zinc-800 border border-gray-300 dark:border-zinc-700 rounded-lg text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-brand-500 focus:border-transparent transition-all"
                        >
                            {hours.map((hour) => (
                                <option key={`close-${hour}`} value={hour}>
                                    {formatHour(hour)}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Preview */}
                    <div className="p-4 bg-brand-50 dark:bg-brand-900/20 rounded-lg border border-brand-100 dark:border-brand-800">
                        <p className="text-sm text-brand-700 dark:text-brand-300">
                            <span className="font-medium">Horario:</span>{' '}
                            {formatHour(openingHour)} - {formatHour(closingHour)}{' '}
                            <span className="text-brand-500">({closingHour - openingHour} horas)</span>
                        </p>
                    </div>

                    {/* Error */}
                    {error && (
                        <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-sm text-red-600 dark:text-red-400">
                            {error}
                        </div>
                    )}
                </div>

                <DialogFooter className="gap-2 sm:gap-0">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-zinc-800 hover:bg-gray-200 dark:hover:bg-zinc-700 rounded-lg transition-colors"
                    >
                        Cancelar
                    </button>
                    <button
                        onClick={handleSave}
                        disabled={loading}
                        className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-brand-600 hover:bg-brand-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {loading ? (
                            <>
                                <Loader2 className="h-4 w-4 animate-spin" />
                                Guardando...
                            </>
                        ) : (
                            <>
                                <Save className="h-4 w-4" />
                                Guardar
                            </>
                        )}
                    </button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
