'use client';

import { useEffect, useState } from 'react';
import api from '@/lib/axios';
import { useAuth } from '@/hooks/use-auth';
import { Loader2, CalendarX, Clock, Calendar } from 'lucide-react';
import { format } from 'date-fns';
import { es } from 'date-fns/locale';
import { Button } from '@/components/ui/button';
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

interface Booking {
    id: string;
    facility_id: string;
    start_time: string;
    end_time: string;
    status: string;
}

interface Facility {
    id: string;
    name: string;
}

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
        } catch (err) {
            console.error("Failed to cancel booking", err);
            alert("No se pudo cancelar la reserva");
        } finally {
            setCancellingId(null);
        }
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

                        return (
                            <div key={booking.id} className="bg-white dark:bg-zinc-900 rounded-xl border border-gray-200 dark:border-zinc-800 p-5 shadow-sm hover:shadow-md transition-shadow">
                                <div className="flex justify-between items-start mb-4">
                                    <div>
                                        <h3 className="font-semibold text-gray-900 dark:text-white text-lg">
                                            {facility ? facility.name : 'Instalación Desconocida'}
                                        </h3>
                                        <div className="flex items-center gap-1.5 text-xs font-medium text-brand-600 bg-brand-50 dark:bg-brand-900/20 px-2 py-0.5 rounded mt-1 w-fit capitalize">
                                            {booking.status === 'confirmed' ? 'confirmada' : booking.status}
                                        </div>
                                    </div>
                                    <div className="bg-gray-50 dark:bg-zinc-800 p-2 rounded-lg">
                                        <Calendar className="h-5 w-5 text-gray-500" />
                                    </div>
                                </div>

                                <div className="space-y-3 text-sm text-gray-600 dark:text-gray-400">
                                    <div className="flex items-center gap-2">
                                        <Clock className="h-4 w-4 text-gray-400" />
                                        <span className="capitalize">
                                            {format(startDate, 'PPP', { locale: es })} <br />
                                            <span className="text-gray-500 text-xs">
                                                {format(startDate, 'p', { locale: es })} - {format(endDate, 'p', { locale: es })}
                                            </span>
                                        </span>
                                    </div>
                                </div>

                                <div className="mt-6 pt-4 border-t border-gray-100 dark:border-zinc-800 flex justify-end">
                                    <AlertDialog>
                                        <AlertDialogTrigger asChild>
                                            <Button
                                                variant="destructive"
                                                size="sm"
                                                disabled={cancellingId === booking.id}
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
