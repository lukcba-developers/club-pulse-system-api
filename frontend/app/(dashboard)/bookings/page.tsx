'use client';

import { useEffect, useState } from 'react';
import api from '@/lib/axios';
import { useAuth } from '@/hooks/use-auth';
import { Loader2, CalendarX, Calendar } from 'lucide-react';
import { format } from 'date-fns';
import { es } from 'date-fns/locale';
import { Button } from '@/components/ui/button';
import { BookingExpiryTimer } from '@/components/booking-expiry-timer';
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

import { Booking } from '@/types/booking';

interface Facility {
    id: string;
    name: string;
}

const getStatusConfig = (status: string) => {
    const s = status.toUpperCase();
    switch (s) {
        case 'CONFIRMED': return { label: 'Confirmada', color: 'text-green-600 bg-green-50 dark:bg-green-900/20' };
        case 'PENDING_PAYMENT': return { label: 'Pendiente de Pago', color: 'text-amber-600 bg-amber-50 dark:bg-amber-900/20' };
        case 'CANCELLED': return { label: 'Cancelada', color: 'text-red-600 bg-red-50 dark:bg-red-900/20' };
        case 'EXPIRED': return { label: 'Expirada', color: 'text-gray-600 bg-gray-50 dark:bg-gray-800' };
        case 'COMPLETED': return { label: 'Completada', color: 'text-blue-600 bg-blue-50 dark:bg-blue-900/20' };
        case 'NO_SHOW': return { label: 'Ausente', color: 'text-purple-600 bg-purple-50 dark:bg-purple-900/20' };
        default: return { label: status, color: 'text-gray-600 bg-gray-50' };
    }
};

export default function BookingsPage() {
    const { user, loading: authLoading } = useAuth();
    const [bookings, setBookings] = useState<Booking[]>([]);
    const [facilities, setFacilities] = useState<Record<string, Facility>>({});
    const [loading, setLoading] = useState(true);
    const [cancellingId, setCancellingId] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            if (!user) return;
            try {
                // Parallel fetch
                const [bookingsRes, facilitiesRes] = await Promise.all([
                    api.get('/bookings'),
                    api.get('/facilities')
                ]);

                const facilitiesMap: Record<string, Facility> = {};
                (facilitiesRes.data.data || []).forEach((f: Facility) => {
                    facilitiesMap[f.id] = f;
                });
                setFacilities(facilitiesMap);

                setBookings(bookingsRes.data.data || []);
            } catch (err) {
                console.error("Failed to fetch data", err);
            } finally {
                setLoading(false);
            }
        };

        if (!authLoading) {
            fetchData();
        }
    }, [user, authLoading]);

    const handleCancel = async (id: string) => {
        setCancellingId(id);
        try {
            await api.delete(`/bookings/${id}`);
            setBookings(prev => prev.filter(b => b.id !== id));
            // Show success message (would be better with a toast context)
            // For now we just remove from list silently - the UI update indicates success
        } catch (error: unknown) {
            console.error("Failed to cancel booking", error);
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const err = error as any;
            // Mostrar mensaje de error del backend si existe (ej. "payment provider error")
            const msg = err.response?.data?.error || "No se pudo cancelar la reserva";
            alert(msg);
        } finally {
            setCancellingId(null);
        }
    };

    const isCancellable = (startTime: string) => {
        const start = new Date(startTime);
        const now = new Date();
        const hoursDiff = (start.getTime() - now.getTime()) / (1000 * 60 * 60);
        return hoursDiff >= 24;
    };

    if (authLoading || loading) {
        return (
            <div className="h-full flex items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-brand-600" />
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex flex-col gap-2">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Mis Reservas</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">Gestiona tus próximas reservas de instalaciones.</p>
            </div>

            {bookings.length === 0 ? (
                <div className="text-center py-12 bg-white dark:bg-zinc-900 rounded-xl border border-dashed border-gray-300 dark:border-zinc-700">
                    <Calendar className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 dark:text-white">Aún no tienes reservas</h3>
                    <p className="text-gray-500 text-sm mt-1">Comienza reservando una instalación desde el Panel de Control.</p>
                </div>
            ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {bookings.map((booking) => {
                        const facility = facilities[booking.facility_id];
                        const startDate = new Date(booking.start_time);
                        const endDate = new Date(booking.end_time);
                        const statusConfig = getStatusConfig(booking.status);

                        return (
                            <div key={booking.id} className="bg-white dark:bg-zinc-900 rounded-xl border border-gray-200 dark:border-zinc-800 p-5 shadow-sm hover:shadow-md transition-shadow relative overflow-hidden">
                                {booking.status === 'PENDING_PAYMENT' && (
                                    <div className="absolute top-0 right-0 p-2">
                                        <span className="relative flex h-3 w-3">
                                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-400 opacity-75"></span>
                                            <span className="relative inline-flex rounded-full h-3 w-3 bg-amber-500"></span>
                                        </span>
                                    </div>
                                )}
                                <div className="flex justify-between items-start mb-4">
                                    <div>
                                        <h3 className="font-semibold text-gray-900 dark:text-white text-lg">
                                            {facility ? facility.name : 'Instalación Desconocida'}
                                        </h3>
                                        <div className={`flex items-center gap-1.5 text-xs font-medium px-2 py-0.5 rounded mt-1 w-fit capitalize ${statusConfig.color}`}>
                                            {statusConfig.label}
                                        </div>
                                    </div>
                                    <div className="bg-gray-50 dark:bg-zinc-800 p-2 rounded-lg">
                                        <Calendar className="h-5 w-5 text-gray-500" />
                                    </div>
                                </div>

                                <div className="space-y-3 text-sm text-gray-600 dark:text-gray-400">
                                    <div className="flex items-center gap-2">
                                        <div className="flex flex-col">
                                            <span className="capitalize">
                                                {format(startDate, 'PPP', { locale: es })}
                                            </span>
                                            <span className="text-gray-500 text-xs">
                                                {format(startDate, 'p', { locale: es })} - {format(endDate, 'p', { locale: es })}
                                            </span>
                                        </div>
                                    </div>

                                    {(Number(booking.total_price) > 0 || (booking.guest_details && booking.guest_details.length > 0)) && (
                                        <div className="pt-2 border-t border-gray-100 dark:border-zinc-800 flex justify-between items-center text-xs">
                                            {Number(booking.total_price) > 0 && (
                                                <span className="font-medium text-gray-900 dark:text-gray-100">
                                                    Precio: ${booking.total_price}
                                                </span>
                                            )}
                                            {booking.guest_details && booking.guest_details.length > 0 && (
                                                <span className="text-muted-foreground">
                                                    {booking.guest_details.length} Invitado(s)
                                                </span>
                                            )}
                                        </div>
                                    )}

                                    {booking.status === 'PENDING_PAYMENT' && booking.payment_expiry && (
                                        <div className="mt-3 bg-amber-50 dark:bg-amber-900/10 p-2 rounded border border-amber-100 dark:border-amber-900/20">
                                            <BookingExpiryTimer
                                                expiry={booking.payment_expiry}
                                                onExpire={() => {
                                                    // Refresh bookings to show EXPIRED status
                                                    // For now user can refresh page
                                                }}
                                            />
                                        </div>
                                    )}
                                </div>

                                <div className="mt-6 pt-4 border-t border-gray-100 dark:border-zinc-800 flex justify-end">
                                    <AlertDialog>
                                        <TooltipProvider>
                                            <Tooltip>
                                                <TooltipTrigger asChild>
                                                    <span tabIndex={0} className="w-full sm:w-auto block"> {/* Wrapper for disabled button tooltip */}
                                                        <AlertDialogTrigger asChild>
                                                            <Button
                                                                variant="destructive"
                                                                size="sm"
                                                                disabled={cancellingId === booking.id || !isCancellable(booking.start_time)}
                                                                className="w-full sm:w-auto"
                                                            >
                                                                {cancellingId === booking.id ? (
                                                                    <Loader2 className="h-4 w-4 animate-spin mr-2" />
                                                                ) : (
                                                                    <CalendarX className="h-4 w-4 mr-2" />
                                                                )}
                                                                Cancelar Reserva
                                                            </Button>
                                                        </AlertDialogTrigger>
                                                    </span>
                                                </TooltipTrigger>
                                                {!isCancellable(booking.start_time) && (
                                                    <TooltipContent>
                                                        <p>Solo se puede cancelar con 24hs de anticipación</p>
                                                    </TooltipContent>
                                                )}
                                            </Tooltip>
                                        </TooltipProvider>

                                        <AlertDialogContent>
                                            <AlertDialogHeader>
                                                <AlertDialogTitle>¿Está seguro?</AlertDialogTitle>
                                                <AlertDialogDescription className="space-y-2">
                                                    <span>Esta acción no se puede deshacer. Esto cancelará permanentemente su reserva para {facility?.name}.</span>
                                                    <span className="block text-sm text-green-600 dark:text-green-400 font-medium">
                                                        💰 Si realizaste un pago, el reembolso se procesará automáticamente.
                                                    </span>
                                                </AlertDialogDescription>
                                            </AlertDialogHeader>
                                            <AlertDialogFooter>
                                                <AlertDialogCancel>Cancelar</AlertDialogCancel>
                                                <AlertDialogAction onClick={() => handleCancel(booking.id)} className="bg-red-600 hover:bg-red-700">
                                                    Sí, cancelar reserva
                                                </AlertDialogAction>
                                            </AlertDialogFooter>
                                        </AlertDialogContent>
                                    </AlertDialog>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
}
