'use client';

import { useAuthContext } from '@/context/auth-context';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
    const { user, loading } = useAuthContext();
    const router = useRouter();

    useEffect(() => {
        if (!loading) {
            if (!user || user.role !== 'SUPER_ADMIN') {
                router.push('/dashboard'); // or 404
            }
        }
    }, [user, loading, router]);

    if (loading) return <div className="p-8">Verifying Admin Privileges...</div>;

    if (!user || user.role !== 'SUPER_ADMIN') return null;

    return (
        <div className="admin-layout">
            <div className="bg-slate-900 text-white p-4 mb-6 rounded-lg">
                <h2 className="text-xl font-bold">🔧 Super Admin Area</h2>
            </div>
            {children}
        </div>
    );
}
