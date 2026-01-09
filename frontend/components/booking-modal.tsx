'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Loader2 } from 'lucide-react';
import api from '@/lib/axios';
import { humanizeError } from '@/lib/error-messages';
import { useAuth } from '@/hooks/use-auth';
import { AvailabilityCalendar } from './availability-calendar';

interface BookingModalProps {
    isOpen: boolean;
    onClose: () => void;
    facilityId: string;
    facilityName: string;
}

export function BookingModal({ isOpen, onClose, facilityId, facilityName }: BookingModalProps) {
    const { user } = useAuth();
    const [date, setDate] = useState('');
    const [startTime, setStartTime] = useState('');
    const [addGuest, setAddGuest] = useState(false);
    const [guestName, setGuestName] = useState('');
    const [guestDNI, setGuestDNI] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        if (!user) {
            setError("Debes iniciar sesión para reservar.");
            setLoading(false);
            return;
        }

        // Validar certificado médico
        if (user.medical_cert_status !== 'VALID') {
            setError(humanizeError('medical_certificate_invalid'));
            setLoading(false);
            return;
        }

        try {
            // Simple validation
            if (!date || !startTime) {
                setError("Por favor selecciona fecha y hora.");
                setLoading(false);
                return;
            }

            // Validar que la fecha no sea pasada
            const startDateTime = new Date(`${date}T${startTime}:00`);
            if (startDateTime < new Date()) {
                setError("No podés reservar para fechas u horarios pasados.");
                setLoading(false);
                return;
            }
            const endDateTime = new Date(startDateTime.getTime() + 60 * 60 * 1000); // +1 hour

            await api.post('/bookings', {
                user_id: user.id,
                facility_id: facilityId,
                start_time: startDateTime.toISOString(),
                end_time: endDateTime.toISOString(),
                guest_details: addGuest && guestName ? [{
                    name: guestName,
                    dni: guestDNI,
                    fee_amount: 1500 // Hardcoded fee for guest for MVP
                }] : undefined
            });

            // Reusing error state for success message to be simple or alert, 
            // but user hates alerts. 
            // Let's close and show a nice toast? We don't have a global toast context easily accessible here yet without refactoring.
            // Let's just change the modal content to success state? 
            // I'll assume we can use the `onClose` and just alert for now? No, user said NO ALERT.
            // I will implement a success step in the modal.

            // Actually, I can replace the form with a success message.
            setSuccess(true);
            setTimeout(() => {
                onClose();
                setSuccess(false);
                // Reset form
                setDate('');
                setStartTime('');
            }, 2000);

        } catch (err: unknown) {
            console.error("Booking failed", err);
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const axiosError = err as any;
            if (axiosError.response?.status === 409) {
                setError(humanizeError('booking_conflict'));
            } else {
                const serverMessage = axiosError.response?.data?.error || '';
                setError(humanizeError(serverMessage));
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[425px] bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800">
                <DialogHeader>
                    <DialogTitle>Reservar {facilityName}</DialogTitle>
                    <DialogDescription>
                        Selecciona una fecha y hora de inicio. Todas las reservas son de 1 hora.
                    </DialogDescription>
                </DialogHeader>
                {success ? (
                    <div className="flex flex-col items-center justify-center py-8 space-y-4">
                        <div className="h-12 w-12 rounded-full bg-green-100 text-green-600 flex items-center justify-center">
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                        </div>
                        <h3 className="text-lg font-medium text-center">¡Reserva Exitosa!</h3>
                        <p className="text-sm text-gray-500 text-center">Tu pista ha sido reservada correctamente.</p>
                    </div>
                ) : (
                    <form onSubmit={handleSubmit} className="grid gap-4 py-4">
                        {error && (
                            <div className="p-2 text-sm text-red-600 bg-red-50 dark:bg-red-900/20 rounded">
                                {error}
                            </div>
                        )}

                        <AvailabilityCalendar
                            facilityId={facilityId}
                            onSlotSelect={(d, t) => {
                                setDate(d);
                                setStartTime(t);
                            }}
                        />

                        {/* Guest Section */}
                        <fieldset className="space-y-3 border-t pt-4">
                            <legend className="sr-only">Agregar invitado a la reserva</legend>

                            <div className="flex items-center justify-between">
                                <div className="flex items-center space-x-2">
                                    <input
                                        type="checkbox"
                                        id="addGuest"
                                        checked={addGuest}
                                        onChange={(e) => setAddGuest(e.target.checked)}
                                        className="h-4 w-4 rounded border-gray-300 text-brand-600 focus:ring-brand-500"
                                    />
                                    <label htmlFor="addGuest" className="text-sm font-medium text-gray-900 dark:text-gray-100">
                                        ¿Viene alguien más contigo?
                                    </label>
                                </div>
                                <span className="text-xs text-gray-500 bg-gray-100 dark:bg-zinc-800 px-2 py-1 rounded-full">
                                    +$1.500 por invitado
                                </span>
                            </div>

                            {addGuest && (
                                <div className="grid grid-cols-2 gap-3 mt-3 p-3 bg-gray-50 dark:bg-zinc-800 rounded-lg">
                                    <div>
                                        <label htmlFor="guestName" className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
                                            Nombre del invitado
                                        </label>
                                        <input
                                            id="guestName"
                                            type="text"
                                            placeholder="Juan Pérez"
                                            value={guestName}
                                            onChange={(e) => setGuestName(e.target.value)}
                                            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                                            required={addGuest}
                                        />
                                    </div>
                                    <div>
                                        <label htmlFor="guestDNI" className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
                                            DNI
                                        </label>
                                        <input
                                            id="guestDNI"
                                            type="text"
                                            placeholder="12.345.678"
                                            value={guestDNI}
                                            onChange={(e) => setGuestDNI(e.target.value)}
                                            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                                            required={addGuest}
                                        />
                                    </div>
                                </div>
                            )}
                        </fieldset>

                        <div className="text-xs text-center text-gray-500">
                            {date && startTime ? `Seleccionado: ${date} a las ${startTime}` : "Por favor selecciona un horario"}
                        </div>

                        <DialogFooter>
                            <Button type="button" variant="outline" onClick={onClose} disabled={loading}>
                                Cancelar
                            </Button>
                            <Button type="submit" className="bg-brand-600 hover:bg-brand-700 text-white" disabled={loading || !date || !startTime}>
                                {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                Confirmar Reserva
                            </Button>
                        </DialogFooter>
                    </form>
                )}
            </DialogContent>
        </Dialog>
    );
}
