'use client';

import { SessionList } from '@/components/session-list';
import { FamilyList } from '@/components/profile/family-list';
import { BillingSection } from '@/components/profile/billing-section';
import { GamificationStats } from '@/components/profile/gamification-stats';
import { HealthSection } from '@/components/profile/health-section';

export default function ProfilePage() {
    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-bold tracking-tight">Perfil & Estadísticas</h1>
            </div>

            <div className="grid gap-6">
                <GamificationStats />
                <HealthSection />
                <BillingSection />
                <FamilyList />
                <SessionList />
            </div>
        </div>
    );
}
