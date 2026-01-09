'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import api from '@/lib/axios';
import { useAuth } from '@/hooks/use-auth';
import { Loader2, Lock, Mail, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';

export function LoginForm() {
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const { login } = useAuth();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');

        try {
            await api.post('/auth/login', { email, password });
            await login();
            router.push('/'); // SPA navigation for better performance
        } catch (err: unknown) {
            const errorMessage = (err as { response?: { data?: { error?: string } } })?.response?.data?.error;
            setError(errorMessage || 'Credenciales inválidas. Por favor, inténtelo de nuevo.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="w-full max-w-sm space-y-8">
            <div className="text-center">
                <h2 className="mt-6 text-3xl font-extrabold text-gray-900 tracking-tight">
                    Bienvenido de nuevo
                </h2>
                <p className="mt-2 text-sm text-gray-600">
                    Inicia sesión para acceder al panel de tu club
                </p>
            </div>

            <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
                <div className="space-y-4">
                    {/* Email Field */}
                    <div className="space-y-1.5">
                        <label
                            htmlFor="email-address"
                            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
                        >
                            Correo electrónico
                        </label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Mail className="h-5 w-5 text-gray-400" aria-hidden="true" />
                            </div>
                            <input
                                id="email-address"
                                name="email"
                                type="email"
                                autoComplete="email"
                                required
                                aria-describedby={error ? 'login-error' : undefined}
                                className="appearance-none relative block w-full px-3 py-3 pl-10 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-lg focus:outline-none focus:ring-brand-500 focus:border-brand-500 focus:z-10 sm:text-sm transition-all dark:bg-zinc-800 dark:border-zinc-700 dark:text-white dark:placeholder-gray-400"
                                placeholder="tu@email.com"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                        </div>
                    </div>

                    {/* Password Field */}
                    <div className="space-y-1.5">
                        <label
                            htmlFor="password"
                            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
                        >
                            Contraseña
                        </label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Lock className="h-5 w-5 text-gray-400" aria-hidden="true" />
                            </div>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                autoComplete="current-password"
                                required
                                aria-describedby={error ? 'login-error' : undefined}
                                className="appearance-none relative block w-full px-3 py-3 pl-10 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-lg focus:outline-none focus:ring-brand-500 focus:border-brand-500 focus:z-10 sm:text-sm transition-all dark:bg-zinc-800 dark:border-zinc-700 dark:text-white dark:placeholder-gray-400"
                                placeholder="••••••••"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                        </div>
                    </div>
                </div>

                {error && (
                    <div
                        id="login-error"
                        className="rounded-lg bg-red-50 dark:bg-red-900/20 p-4 border border-red-200 dark:border-red-800"
                        role="alert"
                        aria-live="polite"
                    >
                        <div className="flex items-start gap-3">
                            <AlertCircle className="h-5 w-5 text-red-500 mt-0.5 flex-shrink-0" aria-hidden="true" />
                            <div>
                                <h3 className="text-sm font-medium text-red-800 dark:text-red-200">
                                    No pudimos iniciar sesión
                                </h3>
                                <p className="text-sm text-red-700 dark:text-red-300 mt-1">
                                    {error}
                                </p>
                            </div>
                        </div>
                    </div>
                )}

                <div>
                    <button
                        type="submit"
                        disabled={isLoading}
                        className={cn(
                            "group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-brand-600 hover:bg-brand-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500 transition-all shadow-lg hover:shadow-xl hover:-translate-y-0.5",
                            isLoading && "opacity-75 cursor-not-allowed"
                        )}
                    >
                        {isLoading ? (
                            <Loader2 className="h-5 w-5 animate-spin" />
                        ) : (
                            "Iniciar Sesión"
                        )}
                    </button>
                </div>

                <div className="relative mt-6">
                    <div className="absolute inset-0 flex items-center">
                        <div className="w-full border-t border-gray-300"></div>
                    </div>
                    <div className="relative flex justify-center text-sm">
                        <span className="px-2 bg-white text-gray-500">O continuar con</span>
                    </div>
                </div>

                <div className="mt-6">
                    <button
                        type="button"
                        onClick={() => {
                            const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
                            const redirectUri = process.env.NEXT_PUBLIC_GOOGLE_REDIRECT_URL;
                            const scope = 'https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile';
                            window.location.href = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&redirect_uri=${redirectUri}&response_type=code&scope=${scope}&access_type=offline&prompt=consent`;
                        }}
                        className="w-full flex justify-center py-3 px-4 border border-gray-300 rounded-lg shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500 transition-all"
                    >
                        <svg className="h-5 w-5 mr-2" viewBox="0 0 24 24">
                            <path
                                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                                fill="#4285F4"
                            />
                            <path
                                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                                fill="#34A853"
                            />
                            <path
                                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                                fill="#FBBC05"
                            />
                            <path
                                d="M12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 12-4.53z"
                                fill="#EA4335"
                            />
                        </svg>
                        Google
                    </button>
                </div>
            </form>
        </div>
    );
}
