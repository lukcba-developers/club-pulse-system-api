'use client';

import { MembershipTier } from '@/types/membership';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Check } from 'lucide-react';
import { cn } from '@/lib/utils';
import { formatARS } from '@/lib/currency';

interface PricingCardsProps {
    tiers: MembershipTier[];
    onSelectTier: (tierId: string) => void;
    loadingId?: string | null;
}

export function PricingCards({ tiers, onSelectTier, loadingId }: PricingCardsProps) {
    if (!tiers || tiers.length === 0) {
        return <div className="text-center p-8 text-gray-500">No hay planes de membresía disponibles por el momento.</div>;
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {tiers.map((tier) => (
                <Card
                    key={tier.id}
                    className={cn(
                        "relative flex flex-col hover:border-brand-500 transition-colors duration-300 dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800",
                        tier.name === 'Gold' ? 'border-brand-500 shadow-lg scale-105 z-10' : ''
                    )}
                >
                    {tier.name === 'Gold' && (
                        <div className="absolute -top-4 left-0 right-0 flex justify-center">
                            <span className="bg-brand-600 text-white text-xs font-bold px-3 py-1 rounded-full uppercase tracking-wider">
                                Más Popular
                            </span>
                        </div>
                    )}
                    <CardHeader>
                        <CardTitle className="text-2xl font-bold">{tier.name}</CardTitle>
                        <CardDescription>{tier.description}</CardDescription>
                    </CardHeader>
                    <CardContent className="flex-grow">
                        <div className="mb-6">
                            <span className="text-4xl font-extrabold">{formatARS(tier.monthly_fee)}</span>
                            <span className="text-gray-500 ml-2">/mes</span>
                        </div>
                        <ul className="space-y-3">
                            {tier.benefits && tier.benefits.map((benefit, index) => (
                                <li key={index} className="flex items-start">
                                    <Check className="h-5 w-5 text-green-500 mr-2 flex-shrink-0" />
                                    <span className="text-sm text-gray-600 dark:text-gray-300">{benefit}</span>
                                </li>
                            ))}
                        </ul>
                    </CardContent>
                    <CardFooter>
                        <Button
                            className="w-full"
                            variant={tier.name === 'Gold' ? 'default' : 'outline'}
                            onClick={() => onSelectTier(tier.id)}
                            disabled={!!loadingId}
                        >
                            {loadingId === tier.id ? 'Procesando...' : 'Suscribirse Ahora'}
                        </Button>
                    </CardFooter>
                </Card>
            ))}
        </div>
    );
}
