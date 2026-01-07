'use client';

import { useEffect } from 'react';
import { AlertTriangle, RefreshCcw } from 'lucide-react';

export default function Error({
    error,
    reset,
}: {
    error: Error & { digest?: string };
    reset: () => void;
}) {
    useEffect(() => {
        // Log the error to an error reporting service
        console.error(error);
    }, [error]);

    return (
        <div className="h-full flex flex-col items-center justify-center min-h-[50vh] space-y-4 px-4 text-center">
            <div className="p-4 bg-red-50 rounded-full dark:bg-red-900/20">
                <AlertTriangle className="h-10 w-10 text-red-600 dark:text-red-400" />
            </div>
            <div className="space-y-2">
                <h2 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-gray-100">
                    Algo salió mal
                </h2>
                <p className="text-gray-500 dark:text-gray-400 max-w-[400px]">
                    Ha ocurrido un error inesperado al cargar el panel de control. Por favor, inténtelo de nuevo.
                </p>
            </div>
            <button
                onClick={reset}
                className="flex items-center gap-2 px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition-colors shadow-sm font-medium"
            >
                <RefreshCcw className="h-4 w-4" />
                Intentar de nuevo
            </button>
        </div>
    );
}
