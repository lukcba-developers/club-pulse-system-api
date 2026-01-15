'use client';

import { useState, useEffect, useCallback } from 'react';
import { membershipService, Membership } from '@/services/membership-service';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, Users, AlertCircle, CheckCircle } from 'lucide-react';
import { formatARS } from '@/lib/currency';

const statusColors: Record<string, string> = {
    ACTIVE: 'bg-green-100 text-green-800',
    PENDING: 'bg-yellow-100 text-yellow-800',
    INACTIVE: 'bg-gray-100 text-gray-800',
    CANCELLED: 'bg-red-100 text-red-800',
    EXPIRED: 'bg-orange-100 text-orange-800',
};

const statusLabels: Record<string, string> = {
    ACTIVE: 'Activo',
    PENDING: 'Pendiente',
    INACTIVE: 'Inactivo',
    CANCELLED: 'Cancelado',
    EXPIRED: 'Expirado',
};

export default function MembershipAdminPage() {
    const [memberships, setMemberships] = useState<Membership[]>([]);
    const [loading, setLoading] = useState(true);

    const loadMemberships = useCallback(async () => {
        setLoading(true);
        try {
            const data = await membershipService.listAllMemberships();
            setMemberships(data || []);
        } catch (error) {
            console.error('Failed to load memberships', error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadMemberships();
    }, [loadMemberships]);

    // Summary calculations
    const activeCount = memberships.filter(m => m.status === 'ACTIVE').length;
    const pendingCount = memberships.filter(m => m.status === 'PENDING').length;
    const overdueCount = memberships.filter(m => Number(m.outstanding_balance || 0) > 0).length;

    return (
        <div className="space-y-6 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Administración de Membresías</h1>
                <p className="text-muted-foreground">Vista general de todos los socios y sus estados de pago.</p>
            </div>

            {/* Summary Cards */}
            <div className="grid gap-4 md:grid-cols-3">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Socios Activos</CardTitle>
                        <Users className="h-4 w-4 text-green-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{activeCount}</div>
                        <p className="text-xs text-muted-foreground">Membresías en estado activo</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Pendientes</CardTitle>
                        <AlertCircle className="h-4 w-4 text-yellow-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{pendingCount}</div>
                        <p className="text-xs text-muted-foreground">Esperando confirmación</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Con Saldo Pendiente</CardTitle>
                        <AlertCircle className="h-4 w-4 text-red-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{overdueCount}</div>
                        <p className="text-xs text-muted-foreground">Socios con deuda</p>
                    </CardContent>
                </Card>
            </div>

            {/* Memberships Table */}
            <Card>
                <CardHeader>
                    <CardTitle>Lista de Membresías</CardTitle>
                    <CardDescription>Todos los socios registrados en el club.</CardDescription>
                </CardHeader>
                <CardContent>
                    {loading ? (
                        <div className="flex justify-center py-10">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    ) : memberships.length === 0 ? (
                        <div className="text-center py-10 text-muted-foreground">
                            No hay membresías registradas.
                        </div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="border-b">
                                        <th className="text-left p-2">ID Usuario</th>
                                        <th className="text-left p-2">Estado</th>
                                        <th className="text-left p-2">Inicio</th>
                                        <th className="text-left p-2">Saldo Pendiente</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {memberships.map((membership) => (
                                        <tr key={membership.id} className="border-b hover:bg-muted/50">
                                            <td className="p-2 font-mono text-xs">{membership.user_id}</td>
                                            <td className="p-2">
                                                <Badge className={statusColors[membership.status] || 'bg-gray-100'}>
                                                    {statusLabels[membership.status] || membership.status}
                                                </Badge>
                                            </td>
                                            <td className="p-2">{new Date(membership.start_date).toLocaleDateString()}</td>
                                            <td className="p-2">
                                                {Number(membership.outstanding_balance || 0) > 0 ? (
                                                    <span className="text-red-600 font-medium">
                                                        {formatARS(membership.outstanding_balance)}
                                                    </span>
                                                ) : (
                                                    <span className="flex items-center text-green-600">
                                                        <CheckCircle className="h-4 w-4 mr-1" />
                                                        Al día
                                                    </span>
                                                )}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
