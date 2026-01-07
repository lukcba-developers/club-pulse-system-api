'use client';

import React from 'react';
import { CreditCard, AlertCircle, CheckCircle } from 'lucide-react';
import { useAuth } from '@/hooks/use-auth';

// Mock data, eventually this should come from a billing service/API

export function BillingSection() {
    const { user } = useAuth();
    // In a real app, use SWR or useEffect to fetch /api/v1/billing/balance

    // Simulation for demo:
    const balance = 15000; // Example debt
    const hasDebt = balance > 0;

    const [loadingPayment, setLoadingPayment] = React.useState(false);

    const handlePay = async () => {
        if (!user) return;
        setLoadingPayment(true);
        try {
            // Call Backend to get Preference URL
            const response = await import('@/lib/axios').then(mod => mod.default.post('/payments/checkout', {
                amount: balance,
                description: `Pago de Saldo Pendiente - ${user.name}`,
                payer_email: user.email,
                reference_id: user.id, // In real app, this might be specific invoice ID
                reference_type: 'MEMBERSHIP'
            }));

            if (response.data.url) {
                window.location.href = response.data.url;
            } else {
                alert('Error al iniciar el pago: URL no recibida');
            }
        } catch (error) {
            console.error('Payment Error:', error);
            alert('Hubo un error al procesar la solicitud de pago.');
        } finally {
            setLoadingPayment(false);
        }
    };

    return (
        <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden mb-6">
            <div className="px-4 py-5 sm:px-6 border-b border-gray-200 dark:border-gray-700 flex items-center gap-2">
                <CreditCard className="h-5 w-5 text-brand-500" />
                <h3 className="text-lg font-medium leading-6 text-gray-900 dark:text-white">
                    Estado de Cuenta
                </h3>
            </div>
            <div className="p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="flex flex-col">
                        <span className="text-sm text-gray-500 dark:text-gray-400">Balance Pendiente</span>
                        <div className="flex items-baseline mt-1">
                            <span className="text-3xl font-bold text-gray-900 dark:text-white">
                                ${balance.toLocaleString()}
                            </span>
                            <span className="ml-2 text-sm text-gray-500">ARS</span>
                        </div>
                        {hasDebt ? (
                            <div className="mt-2 flex items-center text-sm text-red-600 dark:text-red-400">
                                <AlertCircle className="h-4 w-4 mr-1.5" />
                                Pago requerido para habilitar acceso
                            </div>
                        ) : (
                            <div className="mt-2 flex items-center text-sm text-green-600 dark:text-green-400">
                                <CheckCircle className="h-4 w-4 mr-1.5" />
                                Estás al día
                            </div>
                        )}
                    </div>

                    <div className="flex flex-col justify-center items-start md:items-end">
                        <div className="w-full md:w-auto">
                            <button
                                onClick={handlePay}
                                disabled={!hasDebt || loadingPayment}
                                className={`w-full md:w-auto flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white 
                                    ${hasDebt && !loadingPayment
                                        ? 'bg-brand-600 hover:bg-brand-700 focus:ring-brand-500'
                                        : 'bg-gray-400 cursor-not-allowed'} 
                                    focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors`}
                            >
                                <CreditCard className="h-4 w-4 mr-2" />
                                {loadingPayment ? 'Procesando...' : (hasDebt ? 'Pagar Ahora' : 'Sin Deuda')}
                            </button>
                            <p className="mt-2 text-xs text-gray-500 text-center md:text-right">
                                Próximo vencimiento: 05/02/2026
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
