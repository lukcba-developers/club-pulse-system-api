'use client';

import React, { Suspense } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { CheckCircle, XCircle, AlertCircle } from 'lucide-react';

function PaymentResultContent() {
    const searchParams = useSearchParams();
    const router = useRouter();
    const status = searchParams.get('status'); // success, failure, pending

    let icon = <AlertCircle className="h-16 w-16 text-yellow-500" />;
    let title = 'Pago Pendiente';
    let message = 'Tu pago está siendo procesado. Te notificaremos cuando se complete.';

    if (status === 'success') {
        icon = <CheckCircle className="h-16 w-16 text-green-500" />;
        title = '¡Pago Exitoso!';
        message = 'Tu transacción ha sido completada correctamente. Tu saldo ha sido actualizado.';
    } else if (status === 'failure') {
        icon = <XCircle className="h-16 w-16 text-red-500" />;
        title = 'Pago Fallido';
        message = 'Hubo un problema con tu pago. Por favor intenta nuevamente.';
    }

    return (
        <div className="flex flex-col items-center justify-center min-h-[50vh] p-4 text-center">
            <div className="mb-4">{icon}</div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">{title}</h1>
            <p className="text-gray-600 dark:text-gray-300 mb-8 max-w-md">{message}</p>

            <button
                onClick={() => router.push('/profile')}
                className="px-6 py-2 bg-brand-600 text-white rounded-md hover:bg-brand-700 transition-colors"
            >
                Volver al Perfil
            </button>
        </div>
    );
}

export default function PaymentResultPage() {
    return (
        <Suspense fallback={<div>Cargando...</div>}>
            <PaymentResultContent />
        </Suspense>
    );
}
