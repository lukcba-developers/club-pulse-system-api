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
import { Calendar, MapPin, Clock, CreditCard, ArrowRight, Trophy, Sparkles } from 'lucide-react';
import { XPProgressBar, LevelUpModal } from '@/components/gamification';
import { useGamification } from '@/hooks/useGamification';

import { User } from '@/context/auth-context';

export function MemberDashboardView({ user }: { user: User }) {
    const router = useRouter();
    const [membership, setMembership] = useState<Membership | null>(null);
    const [nextBooking, setNextBooking] = useState<Booking | null>(null);
    const [facilities, setFacilities] = useState<Record<string, Facility>>({});

    // Gamification Hook
    const {
        stats,
        showLevelUpModal,
        newLevel,
        closeLevelUpModal,
        getNextLevelXP
    } = useGamification(user.id);

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
                        <Card className="overflow-hidden border-0 shadow-lg ring-1 ring-black/5 dark:ring-white/10">
                            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-brand-400 via-brand-500 to-brand-600" />
                            <CardHeader className="bg-gradient-to-br from-brand-50 via-white to-brand-50/30 dark:from-brand-900/40 dark:via-zinc-900 dark:to-zinc-900 pb-2">
                                <div className="flex justify-between items-start">
                                    <div className="space-y-1">
                                        <div className="flex items-center gap-2 text-brand-600 dark:text-brand-400 font-medium text-sm uppercase tracking-wide">
                                            <Sparkles className="w-4 h-4" />
                                            <span>¡Prepárate!</span>
                                        </div>
                                        <CardTitle className="text-xl sm:text-2xl text-gray-900 dark:text-white">
                                            Tu próximo juego se acerca 🎾
                                        </CardTitle>
                                    </div>
                                    <div className="bg-white dark:bg-zinc-800 shadow-sm border border-brand-100 dark:border-brand-900/30 px-3 py-1 rounded-full text-xs font-bold text-brand-700 dark:text-brand-300">
                                        Confirmado
                                    </div>
                                </div>
                            </CardHeader>
                            <CardContent className="pt-6 bg-white dark:bg-zinc-900">
                                <div className="flex flex-col sm:flex-row justify-between gap-6">
                                    <div className="space-y-4">
                                        <div className="flex items-start gap-3">
                                            <div className="p-2 bg-brand-50 dark:bg-brand-900/20 rounded-lg text-brand-600 dark:text-brand-400">
                                                <Clock className="w-5 h-5" />
                                            </div>
                                            <div>
                                                <p className="text-sm text-gray-500 dark:text-gray-400 font-medium">Fecha y Hora</p>
                                                <p className="text-lg font-bold text-gray-900 dark:text-white">
                                                    {formatDate(nextBooking.start_time)}
                                                </p>
                                                <p className="text-gray-700 dark:text-gray-300">
                                                    {formatTime(nextBooking.start_time)} - {formatTime(nextBooking.end_time)}
                                                </p>
                                            </div>
                                        </div>

                                        <div className="flex items-start gap-3">
                                            <div className="p-2 bg-gray-50 dark:bg-zinc-800 rounded-lg text-gray-600 dark:text-gray-400">
                                                <MapPin className="w-5 h-5" />
                                            </div>
                                            <div>
                                                <p className="text-sm text-gray-500 dark:text-gray-400 font-medium">Ubicación</p>
                                                <p className="font-semibold text-gray-900 dark:text-white">
                                                    {facilities[nextBooking.facility_id]?.name || "Instalación Desconocida"}
                                                </p>
                                                {facilities[nextBooking.facility_id]?.specifications?.surface_type && (
                                                    <p className="text-sm text-gray-600 dark:text-gray-400 capitalize">
                                                        {facilities[nextBooking.facility_id].specifications.surface_type}
                                                    </p>
                                                )}
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex items-end justify-end sm:justify-start">
                                        <button className="w-full sm:w-auto px-6 py-3 bg-brand-600 hover:bg-brand-700 text-white rounded-xl font-semibold shadow-brand-sm hover:shadow-brand-md transition-all flex items-center justify-center gap-2 group">
                                            Gestionar Reserva
                                            <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
                                        </button>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    ) : (
                        <Card className="bg-gray-50 dark:bg-zinc-800/50 border border-gray-100 dark:border-zinc-800 rounded-2xl overflow-hidden">
                            <CardContent className="flex flex-col items-center justify-center py-10 text-center px-4">
                                <div className="w-16 h-16 bg-white dark:bg-zinc-800 rounded-full flex items-center justify-center shadow-sm mb-4">
                                    <Trophy className="w-8 h-8 text-brand-500 dark:text-brand-400" />
                                </div>
                                <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
                                    ¿Sale partido esta semana? 🏃‍♂️
                                </h3>
                                <p className="text-gray-500 dark:text-gray-400 max-w-xs mb-6">
                                    No tienes juegos programados. ¡Es buen momento para reservar cancha y juntar al equipo!
                                </p>
                                <button
                                    onClick={() => router.push('/facilities')}
                                    className="px-6 py-2.5 bg-white dark:bg-zinc-800 text-brand-600 dark:text-brand-400 border border-brand-200 dark:border-brand-900/50 rounded-xl hover:bg-brand-50 dark:hover:bg-brand-900/20 font-semibold transition-all shadow-sm hover:shadow"
                                >
                                    Buscar Cancha Ahora
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
                    {/* Gamification Progress */}
                    {stats && (
                        <Card className="overflow-hidden">
                            <CardHeader className="bg-gradient-to-r from-purple-500/10 to-pink-500/10 pb-2">
                                <CardTitle className="text-lg flex items-center gap-2">
                                    <Sparkles className="w-5 h-5 text-purple-500" />
                                    Tu Progreso
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="pt-4">
                                <XPProgressBar
                                    level={stats.level}
                                    currentXP={stats.experience}
                                    requiredXP={getNextLevelXP()}
                                    totalXP={stats.totalXp}
                                    currentStreak={stats.currentStreak}
                                />
                            </CardContent>
                        </Card>
                    )}

                    <DigitalMemberCard user={user} />
                </div>
            </div>

            {/* Level Up Modal */}
            <LevelUpModal
                isOpen={showLevelUpModal}
                onClose={closeLevelUpModal}
                newLevel={newLevel}
            />

            <SponsorBanner />
        </div >
    );
}
