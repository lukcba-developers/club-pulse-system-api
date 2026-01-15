'use client';

import { useEffect, useState } from 'react';
import { membershipService, Membership, MembershipTier, Subscription } from '@/services/membership-service';
import { PricingCards } from '@/components/pricing-cards';
import { Loader2, CreditCard, ShieldCheck } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuth } from '@/hooks/use-auth';
import { formatARS } from '@/lib/currency';
import { cn } from '@/lib/utils';

export default function MembershipPage() {
    const { user } = useAuth();
    const [tiers, setTiers] = useState<MembershipTier[]>([]);
    const [currentMembership, setCurrentMembership] = useState<Membership | null>(null);
    const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
    const [loading, setLoading] = useState(true);
    const [subscribing, setSubscribing] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // Fetch all data in parallel
                const [tiersData, membershipsData, subscriptionsData] = await Promise.all([
                    membershipService.listTiers(),
                    membershipService.listMyMemberships(),
                    membershipService.listSubscriptions()
                ]);

                setTiers(tiersData || []);
                setSubscriptions(subscriptionsData || []);

                // Find active membership
                const active = membershipsData?.find((m: Membership) => m.status === 'ACTIVE' || m.status === 'PENDING');
                setCurrentMembership(active || null);

            } catch (err) {
                console.error("Failed to load membership data", err);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    const handleSubscribe = async (tierId: string) => {
        if (!user) return;
        setSubscribing(tierId);
        try {
            await membershipService.createMembership({
                membership_tier_id: tierId,
                billing_cycle: 'MONTHLY' // Default for MVP
            });
            window.location.reload(); // Simple refresh to update state
        } catch (err) {
            console.error("Failed to subscribe", err);
            alert("Failed to subscribe. Please try again.");
        } finally {
            setSubscribing(null);
        }
    };

    if (loading) {
        return (
            <div className="flex h-full w-full justify-center items-center">
                <Loader2 className="h-8 w-8 animate-spin text-brand-500" />
            </div>
        );
    }

    return (
        <div className="space-y-8 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="space-y-2">
                <h1 className="text-3xl font-bold tracking-tight text-gray-900 dark:text-gray-100">Gestión de Membresías</h1>
                <p className="text-gray-500 dark:text-gray-400">Gestiona tu suscripción, consulta beneficios y mejora tu plan.</p>
            </div>

            {currentMembership ? (
                <>
                    <div className="grid gap-6 md:grid-cols-2">
                        <Card className="border-brand-200 dark:border-brand-800 bg-brand-50/50 dark:bg-brand-900/10">
                            <CardHeader className="pb-2">
                                <CardTitle className="flex items-center text-xl">
                                    <ShieldCheck className="h-5 w-5 mr-2 text-brand-600 dark:text-brand-400" />
                                    Plan Actual
                                </CardTitle>
                                <CardDescription>Detalles de tu suscripción activa</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="mt-4 space-y-4">
                                    <div>
                                        <p className="text-sm font-medium text-gray-500">Plan</p>
                                        <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                                            {currentMembership.membership_tier.name}
                                        </p>
                                    </div>
                                    <div className="grid grid-cols-2 gap-4">
                                        <div>
                                            <p className="text-sm font-medium text-gray-500">Estado</p>
                                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300 capitalize">
                                                {currentMembership.status.toLowerCase() === 'active' ? 'activo' : currentMembership.status.toLowerCase()}
                                            </span>
                                        </div>
                                        <div>
                                            <p className="text-sm font-medium text-gray-500">Ciclo de Facturación</p>
                                            <p className="text-sm font-medium capitalize">{currentMembership.billing_cycle.toLowerCase().replace('_', ' ') === 'monthly' ? 'mensual' : 'anual'}</p>
                                        </div>
                                        <div>
                                            <p className="text-sm font-medium text-gray-500">Fecha de Inicio</p>
                                            <p className="text-sm font-medium">
                                                {new Date(currentMembership.start_date).toLocaleDateString()}
                                            </p>
                                        </div>
                                        <div>
                                            <p className="text-sm font-medium text-gray-500">Próxima Facturación</p>
                                            <p className="text-sm font-medium">
                                                {new Date(currentMembership.next_billing_date).toLocaleDateString()}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader className="pb-2">
                                <CardTitle className="flex items-center text-xl">
                                    <CreditCard className="h-5 w-5 mr-2 text-gray-500" />
                                    Método de Pago
                                </CardTitle>
                                <CardDescription>Gestiona tu información de facturación</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="flex flex-col items-center justify-center h-40 text-center space-y-3">
                                    <p className="text-sm text-gray-500">La facturación se gestiona automáticamente contra tu saldo.</p>
                                    <Button variant="outline" size="sm" disabled>Gestionar Tarjetas (Próximamente)</Button>
                                </div>
                            </CardContent>
                        </Card>
                    </div>

                    {/* Historial de Suscripciones */}
                    <Card>
                        <CardHeader>
                            <CardTitle>Mis Suscripciones</CardTitle>
                            <CardDescription>Historial de tus suscripciones activas y pasadas.</CardDescription>
                        </CardHeader>
                        <CardContent>
                            {subscriptions.length === 0 ? (
                                <p className="text-sm text-gray-500 text-center py-4">No hay suscripciones registradas.</p>
                            ) : (
                                <div className="overflow-x-auto">
                                    <table className="w-full text-sm text-left">
                                        <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-800 dark:text-gray-400">
                                            <tr>
                                                <th className="px-4 py-3 rounded-l-lg">Fecha</th>
                                                <th className="px-4 py-3">Concepto</th>
                                                <th className="px-4 py-3">Monto</th>
                                                <th className="px-4 py-3 rounded-r-lg">Estado</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {subscriptions.map((sub, idx) => (
                                                <tr key={sub.id || idx} className="border-b dark:border-gray-700">
                                                    <td className="px-4 py-3">
                                                        {new Date(sub.created_at).toLocaleDateString()}
                                                    </td>
                                                    <td className="px-4 py-3">
                                                        Renovación Membresía
                                                    </td>
                                                    <td className="px-4 py-3 font-medium">
                                                        {formatARS(sub.amount)}
                                                    </td>
                                                    <td className="px-4 py-3">
                                                        <span className={cn(
                                                            "px-2 py-1 rounded-full text-xs font-medium",
                                                            sub.status === 'ACTIVE' ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300" :
                                                                sub.status === 'PAST_DUE' ? "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300" :
                                                                    "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300"
                                                        )}>
                                                            {sub.status}
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
                </>
            ) : (
                <div className="space-y-6">
                    <div className="flex flex-col items-center justify-center py-10 bg-gradient-to-br from-indigo-50 to-white dark:from-zinc-900 dark:to-zinc-950 rounded-2xl border border-indigo-100 dark:border-zinc-800">
                        <h2 className="text-2xl font-bold text-center mb-2">Mejora tu Experiencia</h2>
                        <p className="text-gray-500 text-center max-w-md mb-8">Elige el plan que mejor se adapte a tus necesidades y desbloquea instalaciones y funciones premium.</p>
                        <PricingCards tiers={tiers} onSelectTier={handleSubscribe} loadingId={subscribing} />
                    </div>
                </div>
            )}
        </div>
    );
}
