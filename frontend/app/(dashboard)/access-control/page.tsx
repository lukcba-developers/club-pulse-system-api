'use client';

import { useState } from 'react';
import { accessService, AccessResult } from '@/services/access-service';
import { Scan, CheckCircle, XCircle, Search } from 'lucide-react';

export default function AccessControlPage() {
    const [userId, setUserId] = useState('');
    const [result, setResult] = useState<AccessResult | null>(null);
    const [loading, setLoading] = useState(false);
    const [lastScan, setLastScan] = useState<string>('');

    const handleScan = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!userId) return;

        setLoading(true);
        setLastScan(userId);
        try {
            const data = await accessService.simulateEntry(userId);
            setResult(data);

            // Auto-clear after 3 seconds for the next person
            setTimeout(() => {
                if (userId === lastScan) { // Only clear if we haven't started a new one
                    setResult(null);
                    setUserId('');
                }
            }, 3000);

        } catch (error) {
            console.error('Scan failed', error);
            setResult({ status: 'DENIED', reason: 'Error de Sistema' });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-2xl mx-auto py-12">
            <div className="text-center mb-10">
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white flex items-center justify-center gap-3">
                    <Scan className="h-8 w-8" />
                    Punto de Acceso
                </h1>
                <p className="mt-2 text-gray-500 dark:text-gray-400">
                    Control de Ingreso - Kiosco
                </p>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl overflow-hidden border border-gray-100 dark:border-gray-700">
                <div className="p-8">
                    <form onSubmit={handleScan} className="flex gap-4 mb-8">
                        <div className="relative flex-1">
                            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                            <input
                                type="text"
                                value={userId}
                                onChange={(e) => setUserId(e.target.value)}
                                placeholder="Escanear QR / Ingresar UUID"
                                className="w-full pl-10 pr-4 py-3 border border-gray-300 dark:border-gray-600 rounded-xl bg-gray-50 dark:bg-gray-900 text-lg focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all font-mono"
                                autoFocus
                            />
                        </div>
                        <button
                            type="submit"
                            disabled={loading || !userId}
                            className="px-8 py-3 bg-brand-600 hover:bg-brand-700 text-white rounded-xl font-bold text-lg shadow-lg transition-transform active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {loading ? 'Validando...' : 'Escanear'}
                        </button>
                    </form>

                    {result ? (
                        <div className={`text-center p-8 rounded-xl animate-in fade-in zoom-in duration-300 ${result.status === 'GRANTED'
                                ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300 border-2 border-green-200 dark:border-green-800'
                                : 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 border-2 border-red-200 dark:border-red-800'
                            }`}>
                            {result.status === 'GRANTED' ? (
                                <CheckCircle className="h-24 w-24 mx-auto mb-4 text-green-500" />
                            ) : (
                                <XCircle className="h-24 w-24 mx-auto mb-4 text-red-500" />
                            )}

                            <h2 className="text-3xl font-extrabold mb-2">
                                {result.status === 'GRANTED' ? 'ACCESO PERMITIDO' : 'ACCESO DENEGADO'}
                            </h2>
                            {result.reason && (
                                <p className="text-xl opacity-80 font-medium">
                                    {result.reason}
                                </p>
                            )}
                            <p className="mt-4 text-sm opacity-60 font-mono">
                                ID: {lastScan}
                            </p>
                        </div>
                    ) : (
                        <div className="text-center p-12 text-gray-400 dark:text-gray-500 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-xl">
                            <Scan className="h-16 w-16 mx-auto mb-4 opacity-50" />
                            <p className="text-lg">Esperando escaneo...</p>
                        </div>
                    )}
                </div>
                <div className="bg-gray-50 dark:bg-gray-900/50 p-4 text-center text-sm text-gray-500">
                    Sistema de Control de Acceso v1.0
                </div>
            </div>
        </div>
    );
}
