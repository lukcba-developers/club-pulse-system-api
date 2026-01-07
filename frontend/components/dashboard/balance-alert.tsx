import React from 'react';
import { AlertCircle } from 'lucide-react';
import { Membership } from '@/services/membership-service';
import { paymentService } from '@/services/payment-service';

interface BalanceAlertProps {
    membership: Membership;
    onPaymentSuccess: () => void;
}

export function BalanceAlert({ membership, onPaymentSuccess }: BalanceAlertProps) {
    if (!membership || (membership.outstanding_balance ?? 0) <= 0) return null;

    const handlePay = async () => {
        const confirmed = window.confirm(`¿Pagar $${membership.outstanding_balance ?? 0} ahora con MercadoPago (Simulado)?`);
        if (confirmed) {
            await paymentService.initiatePayment(membership.outstanding_balance ?? 0, membership.currency ?? 'ARS');
            alert('¡Pago procesado correctamente!');
            onPaymentSuccess();
        }
    };

    return (
        <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-8 rounded-r shadow-sm">
            <div className="flex justify-between items-center flex-wrap gap-4">
                <div className="flex">
                    <div className="flex-shrink-0">
                        <AlertCircle className="h-5 w-5 text-red-500" aria-hidden="true" />
                    </div>
                    <div className="ml-3">
                        <p className="text-sm text-red-800 font-bold">
                            Saldo Pendiente
                        </p>
                        <p className="text-sm text-red-700">
                            Tienes una deuda de <span className="font-bold">${membership.outstanding_balance} {membership.currency}</span>. El acceso al club puede estar restringido.
                        </p>
                    </div>
                </div>
                <button
                    onClick={handlePay}
                    className="px-4 py-2 bg-red-600 text-white text-sm font-bold rounded hover:bg-red-700 transition shadow-sm"
                >
                    Pagar Deuda
                </button>
            </div>
        </div>
    );
}
