'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/use-auth';
import { AdminDashboardView } from '@/components/dashboard/admin-view';
import { MemberDashboardView } from '@/components/dashboard/member-view';

export default function DashboardPage() {
    const { user, loading } = useAuth();
    const router = useRouter();

    useEffect(() => {
        if (!loading && !user) {
            router.push('/login');
        } else if (!loading && user && user.role === 'SUPER_ADMIN') {
            // Super Admin has their own special dashboard, redirect there
            router.push('/admin/platform');
        }
    }, [user, loading, router]);

    if (loading || !user) {
        return (
            <div className="h-full flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-600"></div>
            </div>
        );
    }

    if (user.role === 'SUPER_ADMIN') {
        // Should have redirected, but show loader just in case
        return (
            <div className="h-full flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-600"></div>
            </div>
        );
    }

    if (user.role === 'ADMIN') {
        return <AdminDashboardView user={user} />;
    }

    return <MemberDashboardView user={user} />;
}
