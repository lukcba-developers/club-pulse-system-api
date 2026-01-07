'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/use-auth';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Building2, Users, DollarSign, Activity } from "lucide-react";
import { clubService, Club } from '@/services/club-service';
import api from '@/lib/axios';

export default function SuperAdminDashboard() {
    const { user, loading } = useAuth();
    const router = useRouter();
    const [clubs, setClubs] = useState<Club[]>([]);
    const [systemStatus, setSystemStatus] = useState<string>("Checking...");
    const [stats, setStats] = useState({
        activeClubs: 0,
        totalRevenue: "$45k", // Mocked for now
        totalUsers: "2,350"   // Mocked for now
    });

    useEffect(() => {
        const init = async () => {
            if (!loading && (!user || user.role !== 'SUPER_ADMIN')) {
                if (user) router.push('/');
                else router.push('/login');
                return;
            }

            if (user?.role === 'SUPER_ADMIN') {
                const checkHealth = async () => {
                    try {
                        const res = await api.get('/healthz');
                        if (res.status === 200 && res.data.status === 'ok') {
                            setSystemStatus("Operativo");
                        } else {
                            setSystemStatus("Degradado");
                        }
                    } catch {
                        setSystemStatus("Sin Conexión");
                    }
                };

                const fetchClubs = async () => {
                    try {
                        const data = await clubService.listClubs(100);
                        setClubs(data);
                        setStats(prev => ({ ...prev, activeClubs: data.length }));
                    } catch (error) {
                        console.error(error);
                    }
                };

                await Promise.all([checkHealth(), fetchClubs()]);
            }
        };

        init();
    }, [user, loading, router]);

    if (loading || !user) return null;

    // Metrics
    const metrics = [
        { title: "Total MRR", value: stats.totalRevenue, change: "+20.1%", icon: DollarSign },
        { title: "Clubes Activos", value: stats.activeClubs.toString(), change: "+2", icon: Building2 },
        { title: "Usuarios Totales", value: stats.totalUsers, change: "+180", icon: Users },
        { title: "Estado Sistema", value: systemStatus, change: "Monitoreo: /healthz", icon: Activity },
    ];

    return (
        <div className="space-y-8">
            <div>
                <h2 className="text-3xl font-bold tracking-tight text-gray-900 dark:text-gray-100">Panel Global</h2>
                <p className="text-muted-foreground text-gray-500">Visión general de la plataforma Club Pulse.</p>
            </div>

            {/* Metrics */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {metrics.map((metric) => (
                    <Card key={metric.title}>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">
                                {metric.title}
                            </CardTitle>
                            <metric.icon className="h-4 w-4 text-muted-foreground text-gray-400" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{metric.value}</div>
                            <p className="text-xs text-muted-foreground text-gray-500">
                                {metric.change}
                            </p>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Club List */}
            <Card className="col-span-4">
                <CardHeader>
                    <CardTitle>Tenants (Clubes)</CardTitle>
                    <CardDescription>
                        Gestiona los clubes registrados en la plataforma.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        {clubs.length === 0 ? (
                            <p className="text-sm text-gray-500">No hay clubes registrados.</p>
                        ) : (
                            clubs.map((club) => (
                                <div key={club.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-900 transition-colors">
                                    <div className="flex items-center gap-4">
                                        <div className={`w-3 h-3 rounded-full ${club.status === 'ACTIVE' ? 'bg-green-500' : 'bg-gray-300'}`} />
                                        <div>
                                            <p className="font-medium text-sm text-gray-900 dark:text-gray-100">{club.name}</p>
                                            <p className="text-xs text-gray-500">{club.domain}</p>
                                        </div>
                                    </div>
                                    <div className="flex gap-8 text-sm text-gray-600 dark:text-gray-400">
                                        <span>{new Date(club.created_at).toLocaleDateString()}</span>
                                        <button className="text-brand-600 hover:underline">Gestionar</button>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
