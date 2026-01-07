'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { ArrowLeft, Loader2, Save } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { facilityService } from '@/services/facility-service';

// Form Validation Schema
const facilitySchema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    type: z.string().min(1, 'Type is required'),
    description: z.string().optional(),
    hourly_rate: z.coerce.number().min(0, 'Rate must be a positive number'),
    capacity: z.coerce.number().int().min(1, 'Capacity must be at least 1'),
    location_name: z.string().min(2, 'Location name is required'),
    location_description: z.string().optional(),
    surface_type: z.string().optional(),
    lighting: z.boolean().default(false),
    covered: z.boolean().default(false),
});

type FacilityFormValues = z.infer<typeof facilitySchema>;

export default function CreateFacilityPage() {
    const router = useRouter();
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const { register, handleSubmit, formState: { errors } } = useForm<FacilityFormValues>({
        resolver: zodResolver(facilitySchema),
        defaultValues: {
            hourly_rate: 0,
            capacity: 4,
            lighting: true,
            covered: false
        }
    });

    const onSubmit = async (data: FacilityFormValues) => {
        setIsSubmitting(true);
        setError(null);
        try {
            await facilityService.create(data);
            router.push('/');
            router.refresh();
        } catch (err: unknown) {
            console.error(err);
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            setError((err as any).response?.data?.error || 'Error al crear la instalación. Por favor, inténtelo de nuevo.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="max-w-3xl mx-auto py-8">
            <button
                onClick={() => router.back()}
                className="flex items-center text-sm text-gray-500 hover:text-gray-700 mb-6 transition-colors"
            >
                <ArrowLeft className="w-4 h-4 mr-1" />
                Volver al Panel
            </button>

            <div className="flex items-center justify-between mb-8">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Crear Nueva Instalación</h1>
                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Añade una nueva instalación deportiva a tu club.</p>
                </div>
            </div>

            <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
                {/* General Info Section */}
                <div className="bg-white dark:bg-zinc-900 p-6 rounded-xl border border-gray-200 dark:border-zinc-800 shadow-sm space-y-6">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-200 border-b border-gray-100 dark:border-zinc-800 pb-2">Información General</h2>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-2">
                            <Label htmlFor="name">Nombre de la Instalación</Label>
                            <Input id="name" placeholder="ej., Pista Central 1" {...register('name')} />
                            {errors.name && <p className="text-xs text-red-500 font-medium">{errors.name.message}</p>}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="type">Tipo</Label>
                            <div className="relative">
                                <select
                                    id="type"
                                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    {...register('type')}
                                >
                                    <option value="">Selecciona un tipo...</option>
                                    <option value="Tennis Court">Pista de Tenis</option>
                                    <option value="Padel Court">Pista de Pádel</option>
                                    <option value="Swimming Pool">Piscina</option>
                                    <option value="Gym">Gimnasio</option>
                                    <option value="Football Field">Campo de Fútbol</option>
                                    <option value="Golf Simulator">Simulador de Golf</option>
                                </select>
                            </div>
                            {errors.type && <p className="text-xs text-red-500 font-medium">{errors.type.message}</p>}
                        </div>

                        <div className="space-y-2 md:col-span-2">
                            <Label htmlFor="description">Descripción (Opcional)</Label>
                            <Input id="description" placeholder="Breve descripción de la instalación..." {...register('description')} />
                        </div>
                    </div>
                </div>

                {/* Capacity & Pricing Section */}
                <div className="bg-white dark:bg-zinc-900 p-6 rounded-xl border border-gray-200 dark:border-zinc-800 shadow-sm space-y-6">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-200 border-b border-gray-100 dark:border-zinc-800 pb-2">Capacidad y Precios</h2>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-2">
                            <Label htmlFor="capacity">Capacidad (Jugadores)</Label>
                            <Input id="capacity" type="number" min="1" {...register('capacity')} />
                            {errors.capacity && <p className="text-xs text-red-500 font-medium">{errors.capacity.message}</p>}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="hourly_rate">Tarifa por Hora ($)</Label>
                            <Input id="hourly_rate" type="number" min="0" step="0.01" {...register('hourly_rate')} />
                            {errors.hourly_rate && <p className="text-xs text-red-500 font-medium">{errors.hourly_rate.message}</p>}
                        </div>
                    </div>
                </div>

                {/* Location & Specs Section */}
                <div className="bg-white dark:bg-zinc-900 p-6 rounded-xl border border-gray-200 dark:border-zinc-800 shadow-sm space-y-6">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-200 border-b border-gray-100 dark:border-zinc-800 pb-2">Ubicación y Especificaciones</h2>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-2">
                            <Label htmlFor="location_name">Nombre de Ubicación / Dirección</Label>
                            <Input id="location_name" placeholder="ej., Edificio Principal, Ala Norte" {...register('location_name')} />
                            {errors.location_name && <p className="text-xs text-red-500 font-medium">{errors.location_name.message}</p>}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="surface_type">Tipo de Superficie (Opcional)</Label>
                            <Input id="surface_type" placeholder="ej., Tierra batida, Césped, Dura" {...register('surface_type')} />
                        </div>

                        <div className="flex items-center space-x-2 pt-8">
                            <input
                                type="checkbox"
                                id="lighting"
                                className="h-4 w-4 rounded border-gray-300 text-brand-600 focus:ring-brand-600"
                                {...register('lighting')}
                            />
                            <Label htmlFor="lighting" className="font-normal cursor-pointer">¿Tiene Iluminación?</Label>
                        </div>

                        <div className="flex items-center space-x-2 pt-8">
                            <input
                                type="checkbox"
                                id="covered"
                                className="h-4 w-4 rounded border-gray-300 text-brand-600 focus:ring-brand-600"
                                {...register('covered')}
                            />
                            <Label htmlFor="covered" className="font-normal cursor-pointer">¿Es Cubierta/Interior?</Label>
                        </div>
                    </div>
                </div>

                {/* Error & Submit */}
                {error && (
                    <div className="p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg flex items-center gap-2">
                        <span className="font-medium">Error:</span> {error}
                    </div>
                )}

                <div className="flex justify-end gap-4">
                    <Button type="button" variant="outline" onClick={() => router.back()}>Cancelar</Button>
                    <Button type="submit" disabled={isSubmitting} className="min-w-[120px] bg-brand-600 hover:bg-brand-700">
                        {isSubmitting ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                Guardando...
                            </>
                        ) : (
                            <>
                                <Save className="mr-2 h-4 w-4" />
                                Crear Instalación
                            </>
                        )}
                    </Button>
                </div>
            </form>
        </div>
    );
}
