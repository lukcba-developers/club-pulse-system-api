'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Calendar, Clock, MapPin, Users } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { facilityService, Facility } from '@/services/facility-service';
import { bookingService } from '@/services/booking-service';


interface TimeSlot {
    time: string;
    available: boolean;
}

export default function NewBookingPage() {
    const router = useRouter();
    const [facilities, setFacilities] = useState<Facility[]>([]);
    const [selectedFacility, setSelectedFacility] = useState<Facility | null>(null);
    const [selectedDate, setSelectedDate] = useState<string>('');
    const [selectedTime, setSelectedTime] = useState<string>('');
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Generate time slots from 8:00 to 22:00
    const timeSlots: TimeSlot[] = [];
    for (let hour = 8; hour <= 22; hour++) {
        const time = `${hour.toString().padStart(2, '0')}:00`;
        timeSlots.push({ time, available: true });
    }

    useEffect(() => {
        loadFacilities();
    }, []);

    const loadFacilities = async () => {
        try {
            const data = await facilityService.list();
            setFacilities(data);
        } catch (err) {
            console.error('Error loading facilities:', err);
            setError('Error al cargar instalaciones');
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = async () => {
        if (!selectedFacility || !selectedDate || !selectedTime) {
            setError('Por favor complete todos los campos');
            return;
        }

        setSubmitting(true);
        setError(null);

        try {
            const startTime = new Date(`${selectedDate}T${selectedTime}:00`);
            const endTime = new Date(startTime);
            endTime.setHours(endTime.getHours() + 1);

            await bookingService.createBooking({
                facility_id: selectedFacility.id,
                start_time: startTime.toISOString(),
                end_time: endTime.toISOString(),
            });

            router.push('/bookings');
        } catch (err) {
            console.error('Error creating booking:', err);
            setError('Error al crear la reserva. Por favor intente nuevamente.');
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-600"></div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                    Nueva Reserva
                </h1>
                <p className="text-gray-500 text-sm mt-1">
                    Selecciona una instalación, fecha y horario para tu reserva.
                </p>
            </div>

            {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
                    {error}
                </div>
            )}

            {/* Step 1: Select Facility */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <MapPin className="h-5 w-5 text-brand-600" />
                        1. Selecciona Instalación
                    </CardTitle>
                    <CardDescription>
                        Elige la cancha o espacio que deseas reservar.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                        {facilities.map((facility) => (
                            <div
                                key={facility.id}
                                onClick={() => setSelectedFacility(facility)}
                                className={`p-4 border-2 rounded-lg cursor-pointer transition-all ${selectedFacility?.id === facility.id
                                    ? 'border-brand-600 bg-brand-50 dark:bg-brand-900/20'
                                    : 'border-gray-200 hover:border-gray-300'
                                    }`}
                            >
                                <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                                    {facility.name}
                                </h3>
                                <p className="text-sm text-gray-500 capitalize">{facility.type}</p>
                                <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
                                    <span className="flex items-center gap-1">
                                        <Users className="h-3 w-3" />
                                        {facility.capacity || 4}
                                    </span>
                                    <span className="font-medium text-brand-600">
                                        ${facility.hourly_rate}/hora
                                    </span>
                                </div>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>

            {/* Step 2: Select Date */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Calendar className="h-5 w-5 text-brand-600" />
                        2. Selecciona Fecha
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <input
                        type="date"
                        value={selectedDate}
                        onChange={(e) => setSelectedDate(e.target.value)}
                        min={new Date().toISOString().split('T')[0]}
                        className="w-full md:w-auto px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
                    />
                </CardContent>
            </Card>

            {/* Step 3: Select Time */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Clock className="h-5 w-5 text-brand-600" />
                        3. Selecciona Horario
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-2">
                        {timeSlots.map((slot) => (
                            <button
                                key={slot.time}
                                onClick={() => setSelectedTime(slot.time)}
                                disabled={!slot.available}
                                className={`px-3 py-2 text-sm rounded-lg transition-colors ${selectedTime === slot.time
                                    ? 'bg-brand-600 text-white'
                                    : slot.available
                                        ? 'bg-gray-100 hover:bg-gray-200 text-gray-700'
                                        : 'bg-gray-100 text-gray-400 cursor-not-allowed'
                                    }`}
                            >
                                {slot.time}
                            </button>
                        ))}
                    </div>
                </CardContent>
            </Card>

            {/* Summary and Submit */}
            {selectedFacility && selectedDate && selectedTime && (
                <Card className="bg-brand-50 dark:bg-brand-900/20 border-brand-200">
                    <CardContent className="pt-6">
                        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                            <div>
                                <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                                    Resumen de Reserva
                                </h3>
                                <p className="text-sm text-gray-600 dark:text-gray-400">
                                    {selectedFacility.name} - {selectedDate} a las {selectedTime}
                                </p>
                                <p className="text-lg font-bold text-brand-600 mt-1">
                                    Total: ${selectedFacility.hourly_rate}
                                </p>
                            </div>
                            <Button
                                onClick={handleSubmit}
                                disabled={submitting}
                                className="bg-brand-600 hover:bg-brand-700 text-white px-8"
                            >
                                {submitting ? 'Procesando...' : 'Confirmar Reserva'}
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
