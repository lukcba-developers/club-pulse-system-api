'use client';

import { Suspense, useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import api from '@/lib/axios';
import { useAuth } from '@/hooks/use-auth';
import { Loader2 } from 'lucide-react';

function GoogleCallbackContent() {
    const searchParams = useSearchParams();
    const router = useRouter();
    const { login } = useAuth();
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const code = searchParams.get('code');
        const process = async () => {
            if (code) {
                try {
                    await api.post('/auth/google', { code });
                    await login();
                    router.push('/');
                } catch (err) {
                    console.error('Google login failed:', err);
                    setError('Error al autenticar con Google. Por favor, intente de nuevo.');
                }
            } else {
                setError('No se recibió el código de autorización.');
            }
        };
        process();
    }, [searchParams, login, router]);

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center min-h-screen p-4">
                <div className="bg-red-50 text-red-800 p-4 rounded-lg shadow-sm max-w-md w-full text-center">
                    <p>{error}</p>
                    <button
                        onClick={() => router.push('/login')}
                        className="mt-4 text-brand-600 font-medium hover:underline"
                    >
                        Volver al inicio de sesión
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="flex flex-col items-center justify-center min-h-screen">
            <Loader2 className="h-8 w-8 animate-spin text-brand-600" />
            <p className="mt-4 text-gray-600">Autenticando con Google...</p>
        </div>
    );
}

export default function GoogleCallbackPage() {
    return (
        <Suspense fallback={<div>Cargando...</div>}>
            <GoogleCallbackContent />
        </Suspense>
    );
}
