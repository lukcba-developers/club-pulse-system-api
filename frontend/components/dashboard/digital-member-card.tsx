import React, { useState } from 'react';
import { User } from '@/services/user-service';
import { accessService } from '@/services/access-service';
import { QrCode, Scan } from 'lucide-react';

interface DigitalMemberCardProps {
    user: User;
}

export function DigitalMemberCard({ user }: DigitalMemberCardProps) {
    const [accessStatus, setAccessStatus] = useState<'IDLE' | 'GRANTED' | 'DENIED'>('IDLE');
    const [message, setMessage] = useState('');
    const [loading, setLoading] = useState(false);

    const simulateAccess = async () => {
        setLoading(true);
        try {
            const result = await accessService.simulateEntry(user.id);
            setAccessStatus(result.status);
            if (result.status === 'DENIED') {
                setMessage(result.reason || 'Acceso Denegado');
            } else {
                setMessage('¡Bienvenido/a!');
            }
        } catch (error) {
            console.error(error);
            setAccessStatus('DENIED');
            setMessage('Error de conexión');
        } finally {
            setLoading(false);
            // Reset after 3 seconds
            setTimeout(() => {
                setAccessStatus('IDLE');
                setMessage('');
            }, 3000);
        }
    };

    return (
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-6 flex flex-col items-center border border-gray-100 dark:border-gray-700 h-full">
            <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">Credencial Digital</h3>

            <div className={`p-4 rounded-xl bg-white mb-6 border-4 transition-colors duration-500 ${accessStatus === 'GRANTED' ? 'border-green-500' :
                    accessStatus === 'DENIED' ? 'border-red-500' : 'border-gray-100'
                }`}>
                <QrCode className="h-40 w-40 text-gray-900" />
            </div>

            <div className="text-center mb-6">
                <p className="font-bold text-2xl text-gray-900 dark:text-white mb-1">{user.name}</p>
                <div className="inline-flex px-3 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-sm font-medium text-gray-600 dark:text-gray-300">
                    {user.role}
                </div>
            </div>

            <div className="min-h-[2rem] mb-4 text-center">
                {accessStatus !== 'IDLE' && (
                    <div className={`px-4 py-2 rounded-lg text-sm font-bold animate-pulse ${accessStatus === 'GRANTED' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                        }`}>
                        {message}
                    </div>
                )}
            </div>

            <button
                onClick={simulateAccess}
                disabled={loading}
                className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-gray-900 hover:bg-gray-800 text-white rounded-lg transition-colors font-medium disabled:opacity-50"
            >
                <Scan className="h-5 w-5" />
                {loading ? 'Validando...' : 'Simular Acceso (QR)'}
            </button>
        </div>
    );
}
