'use client';

import { useState, useCallback, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { BalanceAlert } from '@/components/dashboard/balance-alert';
import { DigitalMemberCard } from '@/components/dashboard/digital-member-card';
import { membershipService, Membership } from '@/services/membership-service';
import { bookingService, Booking } from '@/services/booking-service';
import { facilityService, Facility } from '@/services/facility-service';
import { MatchScheduler } from '@/components/team/match-scheduler';
import { IncidentReportModal } from '@/components/user/incident-report-modal';
import { SponsorBanner } from '@/components/club/sponsor-banner';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, MapPin, Clock, CreditCard } from 'lucide-react';

import { User } from '@/context/auth-context';

export function MemberDashboardView({ user }: { user: User }) {
    const router = useRouter();
    const [membership, setMembership] = useState<Membership | null>(null);
    const [nextBooking, setNextBooking] = useState<Booking | null>(null);
    const [facilities, setFacilities] = useState<Record<string, Facility>>({});

    const fetchData = useCallback(async () => {
        try {
            const [memberships, bookingsData, facilitiesList] = await Promise.all([
                membershipService.listMyMemberships().catch(() => [] as Membership[]),
                bookingService.getMyBookings(),
                facilityService.list(100)
            ]);

            setMembership(memberships && memberships.length > 0 ? memberships[0] : null);

            // Map facilities for easy lookup
            const facilitiesMap: Record<string, Facility> = {};
            facilitiesList.forEach(f => facilitiesMap[f.id] = f);
            setFacilities(facilitiesMap);

            // Find next booking
            if (bookingsData && bookingsData.length > 0) {
                const now = new Date();
                const futureBookings = bookingsData
                    .filter(b => new Date(b.start_time) > now && b.status !== 'CANCELLED')
                    .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime());

                if (futureBookings.length > 0) {
                    setNextBooking(futureBookings[0]);
                }
            }
        } catch (error) {
            console.error('Failed to fetch dashboard data', error);
        }
    }, []);

    useEffect(() => {
        // eslint-disable-next-line
        fetchData();
    }, [fetchData]);

    const formatTime = (isoString: string) => {
        return new Date(isoString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    };

    const formatDate = (isoString: string) => {
        return new Date(isoString).toLocaleDateString([], { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });
    };

    return (
        <div className="max-w-7xl mx-auto">
            <div className="mb-8">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Hola, {user.name} 👋</h1>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    Bienvenido a tu panel de socio.
                </p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Main Column */}
                <div className="lg:col-span-2 space-y-6">
                    {/* Alerts */}
                    {membership && (
                        <BalanceAlert
                            membership={membership}
                            onPaymentSuccess={fetchData}
                        />
                    )}

                    {/* Next Booking Hero Card */}
                    {nextBooking ? (
                        <Card className="bg-gradient-to-br from-brand-50 to-white dark:from-brand-900/10 dark:to-zinc-900 border-brand-100 dark:border-brand-900/20">
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2 text-brand-700 dark:text-brand-400">
                                    <Calendar className="w-5 h-5" /> Tu Próximo Partido
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="flex flex-col sm:flex-row justify-between gap-4">
                                    <div className="space-y-2">
                                        <div className="flex items-center gap-2 text-gray-700 dark:text-gray-300">
                                            <Clock className="w-4 h-4 text-gray-400" />
                                            <span className="font-semibold text-lg">
                                                {formatDate(nextBooking.start_time)}, {formatTime(nextBooking.start_time)} - {formatTime(nextBooking.end_time)}
                                            </span>
                                        </div>
                                        <div className="flex items-center gap-2 text-gray-600 dark:text-gray-400">
                                            <MapPin className="w-4 h-4 text-gray-400" />
                                            <span>
                                                {facilities[nextBooking.facility_id]?.name || "Instalación Desconocida"}
                                                {facilities[nextBooking.facility_id]?.specifications?.surface_type && ` (${facilities[nextBooking.facility_id].specifications.surface_type})`}
                                            </span>
                                        </div>
                                    </div>
                                    <div className="flex items-center">
                                        <button className="px-4 py-2 bg-white dark:bg-zinc-800 text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-zinc-700 rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-700 text-sm font-medium shadow-sm transition-colors">
                                            Ver Detalles
                                        </button>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    ) : (
                        <Card className="bg-white dark:bg-zinc-900 border-dashed border-2 border-gray-200 dark:border-zinc-800">
                            <CardContent className="flex flex-col items-center justify-center py-8 text-center">
                                <Calendar className="w-10 h-10 text-gray-300 mb-2" />
                                <p className="text-gray-500 font-medium">No tienes reservas próximas</p>
                                <button
                                    onClick={() => router.push('/facilities')}
                                    className="mt-4 text-brand-600 hover:text-brand-700 font-medium text-sm"
                                >
                                    Reservar una cancha
                                </button>
                            </CardContent>
                        </Card>
                    )}

                    {/* Quick Actions */}
                    <div className="grid grid-cols-2 gap-4">
                        <button
                            onClick={() => router.push('/facilities')}
                            className="p-6 bg-white dark:bg-zinc-900 border border-gray-200 dark:border-zinc-800 rounded-xl hover:border-brand-300 dark:hover:border-brand-700 hover:shadow-md transition-all text-left group"
                        >
                            <div className="w-10 h-10 rounded-full bg-brand-100 dark:bg-brand-900/30 text-brand-600 dark:text-brand-400 flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
                                <Calendar className="w-5 h-5" />
                            </div>
                            <h3 className="font-semibold text-gray-900 dark:text-white">Nueva Reserva</h3>
                            <p className="text-xs text-gray-500 mt-1">Busca cancha para hoy</p>
                        </button>

                        <button
                            onClick={() => router.push('/membership')}
                            className="p-6 bg-white dark:bg-zinc-900 border border-gray-200 dark:border-zinc-800 rounded-xl hover:border-brand-300 dark:hover:border-brand-700 hover:shadow-md transition-all text-left group"
                        >
                            <div className="w-10 h-10 rounded-full bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
                                <CreditCard className="w-5 h-5" />
                            </div>
                            <h3 className="font-semibold text-gray-900 dark:text-white">Mi Membresía</h3>
                            <p className="text-xs text-gray-500 mt-1">Gestionar pagos y plan</p>
                        </button>

                        {/* Operational Features Quick Access */}
                        <div className="col-span-2 grid grid-cols-2 gap-4 border-t pt-4">
                            <MatchScheduler />
                            <IncidentReportModal />
                        </div>
                    </div>
                </div>

                <div className="space-y-6">
                    <DigitalMemberCard user={user} />
                </div>
            </div>

            <SponsorBanner />
        </div >
    );
}
