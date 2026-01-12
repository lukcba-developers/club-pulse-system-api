'use client';

import { useState, useEffect } from 'react';
import { Calendar, ChevronLeft, ChevronRight, Clock } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { facilityService, Facility } from '@/services/facility-service';
import { bookingService, Booking } from '@/services/booking-service';


export default function BookingCalendarPage() {
    const [currentDate, setCurrentDate] = useState(new Date());
    const [facilities, setFacilities] = useState<Facility[]>([]);
    const [bookings, setBookings] = useState<Booking[]>([]);
    const [selectedFacility, setSelectedFacility] = useState<string>('all');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            try {
                const [facilitiesData, bookingsData] = await Promise.all([
                    facilityService.list(),
                    bookingService.getClubBookings(),
                ]);
                setFacilities(facilitiesData);
                setBookings(bookingsData);
            } catch (err) {
                console.error('Error loading data:', err);
            } finally {
                setLoading(false);
            }
        };
        loadData();
    }, []);

    const getDaysInMonth = (date: Date) => {
        const year = date.getFullYear();
        const month = date.getMonth();
        const firstDay = new Date(year, month, 1);
        const lastDay = new Date(year, month + 1, 0);
        const days: Date[] = [];

        // Add empty days for padding
        for (let i = 0; i < firstDay.getDay(); i++) {
            days.push(new Date(year, month, -firstDay.getDay() + i + 1));
        }

        // Add all days in month
        for (let i = 1; i <= lastDay.getDate(); i++) {
            days.push(new Date(year, month, i));
        }

        return days;
    };

    const getBookingsForDate = (date: Date) => {
        const dateStr = date.toISOString().split('T')[0];
        return bookings.filter((booking) => {
            const bookingDate = new Date(booking.start_time).toISOString().split('T')[0];
            const facilityMatch = selectedFacility === 'all' || booking.facility_id === selectedFacility;
            return bookingDate === dateStr && facilityMatch;
        });
    };

    const prevMonth = () => {
        setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1));
    };

    const nextMonth = () => {
        setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1));
    };

    const isCurrentMonth = (date: Date) => {
        return date.getMonth() === currentDate.getMonth();
    };

    const isToday = (date: Date) => {
        const today = new Date();
        return date.toDateString() === today.toDateString();
    };

    const days = getDaysInMonth(currentDate);
    const monthName = currentDate.toLocaleDateString('es-ES', { month: 'long', year: 'numeric' });

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-600"></div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                        <Calendar className="h-6 w-6 text-brand-600" />
                        Calendario de Reservas
                    </h1>
                    <p className="text-gray-500 text-sm mt-1">
                        Vista global de la ocupación de instalaciones.
                    </p>
                </div>

                <div className="flex items-center gap-4">
                    <select
                        value={selectedFacility}
                        onChange={(e) => setSelectedFacility(e.target.value)}
                        className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                    >
                        <option value="all">Todas las instalaciones</option>
                        {facilities.map((facility) => (
                            <option key={facility.id} value={facility.id}>
                                {facility.name}
                            </option>
                        ))}
                    </select>
                </div>
            </div>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between">
                    <CardTitle className="capitalize">{monthName}</CardTitle>
                    <div className="flex items-center gap-2">
                        <Button variant="outline" size="icon" onClick={prevMonth}>
                            <ChevronLeft className="h-4 w-4" />
                        </Button>
                        <Button variant="outline" size="icon" onClick={nextMonth}>
                            <ChevronRight className="h-4 w-4" />
                        </Button>
                    </div>
                </CardHeader>
                <CardContent>
                    {/* Day Headers */}
                    <div className="grid grid-cols-7 gap-1 mb-2">
                        {['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb'].map((day) => (
                            <div
                                key={day}
                                className="text-center text-xs font-medium text-gray-500 py-2"
                            >
                                {day}
                            </div>
                        ))}
                    </div>

                    {/* Calendar Grid */}
                    <div className="grid grid-cols-7 gap-1">
                        {days.map((date, index) => {
                            const dayBookings = getBookingsForDate(date);
                            const hasBookings = dayBookings.length > 0;
                            const inCurrentMonth = isCurrentMonth(date);
                            const today = isToday(date);

                            return (
                                <div
                                    key={index}
                                    className={`min-h-[80px] p-2 border rounded-lg transition-colors ${!inCurrentMonth
                                        ? 'bg-gray-50 opacity-50'
                                        : today
                                            ? 'bg-brand-50 border-brand-300'
                                            : 'bg-white hover:bg-gray-50'
                                        }`}
                                >
                                    <div className={`text-sm font-medium ${today ? 'text-brand-600' : 'text-gray-700'
                                        }`}>
                                        {date.getDate()}
                                    </div>
                                    {hasBookings && (
                                        <div className="mt-1 space-y-1">
                                            {dayBookings.slice(0, 2).map((booking) => (
                                                <div
                                                    key={booking.id}
                                                    className="text-xs bg-brand-100 text-brand-700 px-1 py-0.5 rounded truncate"
                                                    title={`${new Date(booking.start_time).toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' })}`}
                                                >
                                                    <Clock className="inline h-3 w-3 mr-1" />
                                                    {new Date(booking.start_time).toLocaleTimeString('es-ES', {
                                                        hour: '2-digit',
                                                        minute: '2-digit',
                                                    })}
                                                </div>
                                            ))}
                                            {dayBookings.length > 2 && (
                                                <div className="text-xs text-gray-500">
                                                    +{dayBookings.length - 2} más
                                                </div>
                                            )}
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </CardContent>
            </Card>

            {/* Legend */}
            <Card>
                <CardContent className="pt-4">
                    <div className="flex flex-wrap items-center gap-4 text-sm">
                        <div className="flex items-center gap-2">
                            <div className="w-4 h-4 bg-brand-50 border border-brand-300 rounded"></div>
                            <span>Hoy</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <div className="w-4 h-4 bg-brand-100 rounded"></div>
                            <span>Con reservas</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <div className="w-4 h-4 bg-gray-50 rounded border"></div>
                            <span>Otro mes</span>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
