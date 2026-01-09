'use client';

import { useEffect } from 'react';
import { AlertTriangle, RefreshCcw, Home } from 'lucide-react';

export default function Error({
    error,
    reset,
}: {
    error: Error & { digest?: string };
    reset: () => void;
}) {
    useEffect(() => {
        console.error('[Error Global]:', error);
    }, [error]);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 dark:from-zinc-900 dark:to-zinc-950 p-4 text-center">
            {/* Icono con animación sutil */}
            <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-full mb-6">
                <AlertTriangle className="h-12 w-12 text-red-500 dark:text-red-400" />
            </div>

            {/* Título empático */}
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
                ¡Ups! Algo no salió como esperábamos
            </h2>

            {/* Descripción orientada a acción */}
            <p className="text-gray-600 dark:text-gray-400 mb-8 max-w-md">
                Estamos trabajando para solucionarlo.
                Mientras tanto, puedes intentar recargar la página.
            </p>

            {/* CTAs con jerarquía clara */}
            <div className="flex gap-4">
                <button
                    onClick={reset}
                    className="flex items-center gap-2 px-6 py-3 bg-brand-600 text-white rounded-xl hover:bg-brand-700 transition-all shadow-lg hover:shadow-xl font-medium"
                >
                    <RefreshCcw className="h-4 w-4" />
                    Reintentar
                </button>
                <button
                    onClick={() => window.location.href = '/'}
                    className="flex items-center gap-2 px-6 py-3 bg-white dark:bg-zinc-800 text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-zinc-700 rounded-xl hover:bg-gray-50 dark:hover:bg-zinc-700 transition-all font-medium"
                >
                    <Home className="h-4 w-4" />
                    Ir al inicio
                </button>
            </div>

            {/* Debug info only in development */}
            {process.env.NODE_ENV === 'development' && (
                <div className="mt-8 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-left w-full max-w-2xl overflow-auto">
                    <p className="font-mono text-sm text-red-800 dark:text-red-200 whitespace-pre-wrap">{error.message}</p>
                    <p className="font-mono text-xs text-red-600 dark:text-red-400 mt-2 whitespace-pre-wrap">{error.stack}</p>
                </div>
            )}
        </div>
    );
}
