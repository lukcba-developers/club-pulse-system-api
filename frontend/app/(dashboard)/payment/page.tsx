'use client';

import { useState, useEffect, useCallback } from 'react';
import { paymentService, Payment, OfflinePaymentRequest } from '@/services/payment-service';
import { userService, User } from '@/services/user-service';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { Plus, DollarSign, Clock, CheckCircle, XCircle, Loader2, Search } from 'lucide-react';

const statusColors: Record<string, string> = {
    PENDING: 'bg-yellow-100 text-yellow-800',
    COMPLETED: 'bg-green-100 text-green-800',
    FAILED: 'bg-red-100 text-red-800',
    REFUNDED: 'bg-gray-100 text-gray-800',
};

const methodLabels: Record<string, string> = {
    CASH: 'Efectivo',
    MERCADOPAGO: 'Mercado Pago',
    TRANSFER: 'Transferencia',
    LABOR_EXCHANGE: 'Canje',
};

export default function PaymentDashboardPage() {
    const [payments, setPayments] = useState<Payment[]>([]);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [users, setUsers] = useState<User[]>([]);
    const [userSearch, setUserSearch] = useState('');
    const [submitting, setSubmitting] = useState(false);

    // Form state
    const [formData, setFormData] = useState<OfflinePaymentRequest>({
        amount: 0,
        method: 'CASH',
        payer_id: '',
        notes: '',
    });

    const loadPayments = useCallback(async () => {
        setLoading(true);
        try {
            const { data, total: totalCount } = await paymentService.getPayments();
            setPayments(data || []);
            setTotal(totalCount);
        } catch (error) {
            console.error('Failed to load payments', error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadPayments();
    }, [loadPayments]);

    const handleUserSearch = async () => {
        if (userSearch.length < 2) return;
        try {
            const results = await userService.searchUsers(userSearch);
            setUsers(results || []);
        } catch (error) {
            console.error('Failed to search users', error);
        }
    };

    const handleSubmitOfflinePayment = async () => {
        if (!formData.payer_id || formData.amount <= 0) return;
        setSubmitting(true);
        try {
            await paymentService.createOfflinePayment(formData);
            setIsModalOpen(false);
            setFormData({ amount: 0, method: 'CASH', payer_id: '', notes: '' });
            loadPayments();
        } catch (error) {
            console.error('Failed to create offline payment', error);
            alert('Error al registrar el pago. Intenta de nuevo.');
        } finally {
            setSubmitting(false);
        }
    };

    // Summary calculations
    const totalRevenue = payments
        .filter(p => p.status === 'COMPLETED')
        .reduce((sum, p) => sum + parseFloat(p.amount), 0);
    const pendingCount = payments.filter(p => p.status === 'PENDING').length;
    const completedCount = payments.filter(p => p.status === 'COMPLETED').length;

    return (
        <div className="space-y-6 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Pagos</h1>
                    <p className="text-muted-foreground">Gestión de transacciones y pagos offline.</p>
                </div>
                <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
                    <DialogTrigger asChild>
                        <Button>
                            <Plus className="mr-2 h-4 w-4" />
                            Registrar Pago
                        </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-[480px]">
                        <DialogHeader>
                            <DialogTitle>Registrar Pago Offline</DialogTitle>
                            <DialogDescription>
                                Registra un pago realizado en efectivo, transferencia o canje de trabajo.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="grid gap-4 py-4">
                            {/* User Search */}
                            <div className="space-y-2">
                                <Label>Buscar Socio</Label>
                                <div className="flex gap-2">
                                    <Input
                                        placeholder="Nombre o email..."
                                        value={userSearch}
                                        onChange={(e) => setUserSearch(e.target.value)}
                                        onKeyDown={(e) => e.key === 'Enter' && handleUserSearch()}
                                    />
                                    <Button variant="outline" size="icon" onClick={handleUserSearch}>
                                        <Search className="h-4 w-4" />
                                    </Button>
                                </div>
                                {users.length > 0 && (
                                    <div className="border rounded-md max-h-32 overflow-y-auto">
                                        {users.map((user) => (
                                            <div
                                                key={user.id}
                                                onClick={() => {
                                                    setFormData({ ...formData, payer_id: user.id });
                                                    setUserSearch(user.name);
                                                    setUsers([]);
                                                }}
                                                className="p-2 hover:bg-muted cursor-pointer text-sm"
                                            >
                                                {user.name} <span className="text-muted-foreground">({user.email})</span>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>
                            {/* Amount */}
                            <div className="space-y-2">
                                <Label htmlFor="amount">Monto ($)</Label>
                                <Input
                                    id="amount"
                                    type="number"
                                    value={formData.amount || ''}
                                    onChange={(e) => setFormData({ ...formData, amount: parseFloat(e.target.value) || 0 })}
                                />
                            </div>
                            {/* Method */}
                            <div className="space-y-2">
                                <Label>Método de Pago</Label>
                                <Select
                                    value={formData.method}
                                    onValueChange={(value) => setFormData({ ...formData, method: value as OfflinePaymentRequest['method'] })}
                                >
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="CASH">Efectivo</SelectItem>
                                        <SelectItem value="TRANSFER">Transferencia</SelectItem>
                                        <SelectItem value="LABOR_EXCHANGE">Canje de Trabajo</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                            {/* Notes */}
                            <div className="space-y-2">
                                <Label htmlFor="notes">Notas (opcional)</Label>
                                <Textarea
                                    id="notes"
                                    value={formData.notes || ''}
                                    onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                                    placeholder="Detalles del canje, concepto, etc."
                                />
                            </div>
                        </div>
                        <DialogFooter>
                            <Button variant="outline" onClick={() => setIsModalOpen(false)}>Cancelar</Button>
                            <Button onClick={handleSubmitOfflinePayment} disabled={submitting || !formData.payer_id}>
                                {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                Registrar
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            {/* Summary Cards */}
            <div className="grid gap-4 md:grid-cols-3">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Ingresos Totales</CardTitle>
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">${totalRevenue.toFixed(2)}</div>
                        <p className="text-xs text-muted-foreground">{completedCount} pagos completados</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Pendientes</CardTitle>
                        <Clock className="h-4 w-4 text-yellow-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{pendingCount}</div>
                        <p className="text-xs text-muted-foreground">Pagos por confirmar</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Transacciones</CardTitle>
                        <CheckCircle className="h-4 w-4 text-green-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{total}</div>
                        <p className="text-xs text-muted-foreground">Histórico completo</p>
                    </CardContent>
                </Card>
            </div>

            {/* Payments Table */}
            <Card>
                <CardHeader>
                    <CardTitle>Historial de Transacciones</CardTitle>
                    <CardDescription>Lista de todos los pagos registrados.</CardDescription>
                </CardHeader>
                <CardContent>
                    {loading ? (
                        <div className="flex justify-center py-10">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    ) : payments.length === 0 ? (
                        <div className="text-center py-10 text-muted-foreground">
                            No hay transacciones registradas.
                        </div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="border-b">
                                        <th className="text-left p-2">Fecha</th>
                                        <th className="text-left p-2">Monto</th>
                                        <th className="text-left p-2">Método</th>
                                        <th className="text-left p-2">Estado</th>
                                        <th className="text-left p-2">Notas</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {payments.map((payment) => (
                                        <tr key={payment.id} className="border-b hover:bg-muted/50">
                                            <td className="p-2">{new Date(payment.created_at).toLocaleDateString()}</td>
                                            <td className="p-2 font-medium">${parseFloat(payment.amount).toFixed(2)}</td>
                                            <td className="p-2">{methodLabels[payment.method] || payment.method}</td>
                                            <td className="p-2">
                                                <Badge className={statusColors[payment.status]}>
                                                    {payment.status}
                                                </Badge>
                                            </td>
                                            <td className="p-2 max-w-xs truncate text-muted-foreground">{payment.notes || '-'}</td>
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
