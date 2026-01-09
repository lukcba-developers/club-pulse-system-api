'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { FacilityList } from '@/components/facility-list';
import { SemanticSearch } from '@/components/semantic-search';
import { facilityService } from '@/services/facility-service';
import { bookingService } from '@/services/booking-service';
import { Plus, LayoutGrid, Calendar, TrendingUp } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { User } from '@/context/auth-context';

import { Booking } from '@/services/booking-service';
import { Facility } from '@/services/facility-service';
import { BookingCalendar } from './booking-calendar';

export function AdminDashboardView({ user }: { user: User }) {
    const router = useRouter();
    const [activeTab, setActiveTab] = useState<'overview' | 'calendar'>('overview');
    const [data, setData] = useState<{
        facilities: Facility[];
        bookings: Booking[];
        occupancy: number;
    }>({
        facilities: [],
        bookings: [],
        occupancy: 0
    });

    useEffect(() => {
        const fetchData = async () => {
            try {
                // 1. Fetch facilities
                const facilities = await facilityService.list(100);

                // 2. Fetch bookings for TODAY
                const today = new Date();
                const yyyy = today.getFullYear();
                const mm = String(today.getMonth() + 1).padStart(2, '0');
                const dd = String(today.getDate()).padStart(2, '0');
                const dateStr = `${yyyy}-${mm}-${dd}`;

                // Admin endpoint returns all bookings for the club
                const bookings = await bookingService.getClubBookings({
                    from: dateStr,
                    to: dateStr
                });

                // 3. Calculate Metrics
                const totalFacilities = facilities.length;
                const bookingsCount = bookings ? bookings.length : 0;

                // Estimation: 15 operating hours per facility per day
                const totalSlots = totalFacilities * 15;
                const occupancyRate = totalSlots > 0 ? Math.round((bookingsCount / totalSlots) * 100) : 0;

                setData({
                    facilities,
                    bookings: bookings || [],
                    occupancy: occupancyRate
                });

            } catch (error) {
                console.error('Failed to fetch admin metrics', error);
            }
        };

        fetchData();
    }, []);

    const handleFacilitySelect = useCallback((facilityId: string) => {
        console.log('Selected facility:', facilityId);
    }, []);

    const stats = [
        { title: "Instalaciones Activas", value: data.facilities.length.toString(), icon: LayoutGrid },
        { title: "Reservas Hoy", value: data.bookings.length.toString(), icon: Calendar },
        { title: "Ocupación Hoy", value: `${data.occupancy || 0}%`, icon: TrendingUp },
    ];

    return (
        <div className="space-y-8">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Panel de Administración</h1>
                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                        Gestiona tu club: {user.club_id || "Sede Principal"}
                    </p>
                </div>
                <button
                    onClick={() => router.push('/facilities/create')}
                    className="flex items-center gap-2 px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition-colors shadow-sm font-medium"
                >
                    <Plus className="h-4 w-4" />
                    Nueva Instalación
                </button>
            </div>

            {/* Admin Stats */}
            <div className="grid gap-4 md:grid-cols-3">
                {stats.map((stat) => (
                    <Card key={stat.title}>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
                                {stat.title}
                            </CardTitle>
                            <stat.icon className="h-4 w-4 text-brand-600 dark:text-brand-400" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold text-gray-900 dark:text-white">{stat.value}</div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Tabs */}
            <div className="border-b border-gray-200 dark:border-zinc-800">
                <nav className="-mb-px flex space-x-8">
                    <button
                        onClick={() => setActiveTab('overview')}
                        className={`whitespace-nowrap pb-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'overview'
                            ? 'border-brand-500 text-brand-600 dark:text-brand-400'
                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
                            }`}
                    >
                        Vista General
                    </button>
                    <button
                        onClick={() => setActiveTab('calendar')}
                        className={`whitespace-nowrap pb-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'calendar'
                            ? 'border-brand-500 text-brand-600 dark:text-brand-400'
                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
                            }`}
                    >
                        Calendario (Beta)
                    </button>
                </nav>
            </div>

            {activeTab === 'overview' ? (
                <div className="space-y-6">
                    <div className="flex items-center justify-between">
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Gestión de Instalaciones</h2>
                        <div className="w-full max-w-md">
                            <SemanticSearch onResultSelect={handleFacilitySelect} />
                        </div>
                    </div>
                    <FacilityList />
                </div>
            ) : (
                <div className="space-y-6">
                    <div className="flex items-center justify-between">
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Agenda del Día</h2>
                    </div>
                    {/* Calendar Component */}
                    <BookingCalendar
                        bookings={data.bookings}
                        facilities={data.facilities}
                        date={new Date()}
                    />
                </div>
            )}
        </div>
    );
}
