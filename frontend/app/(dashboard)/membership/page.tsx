'use client';

import { useEffect, useState } from 'react';
import api from '@/lib/axios';
import { Membership, MembershipTier } from '@/types/membership';
import { PricingCards } from '@/components/pricing-cards';
import { Loader2, CreditCard, ShieldCheck } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuth } from '@/hooks/use-auth';

export default function MembershipPage() {
    const { user } = useAuth();
    const [tiers, setTiers] = useState<MembershipTier[]>([]);
    const [currentMembership, setCurrentMembership] = useState<Membership | null>(null);
    const [loading, setLoading] = useState(true);
    const [subscribing, setSubscribing] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // Fetch tiers
                const tiersRes = await api.get('/memberships/tiers');
                setTiers(tiersRes.data.data || []);

                // Fetch current user memberships
                // For MVP, we presume user has one or none active subscription.
                const memRes = await api.get('/memberships');
                const memberships = memRes.data.data || [];

                // Find active one
                const active = memberships.find((m: Membership) => m.status === 'ACTIVE' || m.status === 'PENDING');
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
            await api.post('/memberships', {
                user_id: user.id, // In real app, this comes from token context
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
                                <p className="text-sm text-gray-500">La facturación se gestiona actualmente mediante factura.</p>
                                <Button variant="outline" size="sm">Ver Facturas</Button>
                            </div>
                        </CardContent>
                    </Card>
                </div>
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
