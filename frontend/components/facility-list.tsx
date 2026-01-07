'use client';

import { useEffect, useState } from 'react';
import api from '@/lib/axios';
import { Loader2, Activity } from 'lucide-react';
import { FacilityCard } from '@/components/ui/facility-card';

interface Facility {
    id: string;
    name: string;
    type: string;
    location: {
        name: string;
        description?: string;
    };
    capacity: number;
    status: string;
    hourly_rate: number;
}

export function FacilityList() {
    const [facilities, setFacilities] = useState<Facility[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchFacilities = async () => {
            try {
                const response = await api.get('/facilities');

                // Handle different response structures:
                // 1. { data: [...] } (Standard API wrapper)
                // 2. [...] (Direct array)
                const facilitiesData = response.data.data || response.data || [];

                // Ensure it's an array before setting
                if (Array.isArray(facilitiesData)) {
                    setFacilities(facilitiesData);
                } else {
                    console.error('Facilities data is not an array:', facilitiesData);
                    setFacilities([]);
                }
            } catch (err) {
                console.error("Failed to fetch facilities", err);
                setError('No se pudieron cargar las instalaciones. Por favor, inténtelo de nuevo.');
            } finally {
                setLoading(false);
            }
        };

        fetchFacilities();
    }, []);

    if (loading) {
        return (
            <div className="flex justify-center items-center py-20">
                <Loader2 className="h-10 w-10 text-brand-500 animate-spin" />
            </div>
        );
    }

    if (error) {
        return (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl relative text-sm" role="alert">
                <span className="block sm:inline">{error}</span>
            </div>
        );
    }

    if (facilities.length === 0) {
        return (
            <div className="text-center py-20 bg-white dark:bg-zinc-900 rounded-2xl border border-dashed border-gray-300 dark:border-zinc-700">
                <div className="mx-auto h-12 w-12 text-gray-300 dark:text-gray-600">
                    <Activity className="h-full w-full" />
                </div>
                <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-gray-200">Sin instalaciones</h3>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Comienza creando una nueva instalación.</p>
            </div>
        );
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
            {facilities.map((facility) => (
                <FacilityCard key={facility.id} facility={facility} />
            ))}
        </div>
    );
}
