import { MapPin, Users, DollarSign, CalendarPlus } from 'lucide-react';
import Image from 'next/image';
import { cn } from '@/lib/utils';
import { useState } from 'react';
import { BookingModal } from '@/components/booking-modal';


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

const getFacilityImage = (type: string) => {
    // Using high quality Unsplash placeholders
    switch (type.toLowerCase()) {
        case 'court': return 'https://images.unsplash.com/photo-1622279457486-62dcc4a431d6?q=80&w=800&auto=format&fit=crop';
        case 'swimming': return 'https://images.unsplash.com/photo-1576458088443-04a19bb13da6?q=80&w=800&auto=format&fit=crop';
        case 'gym': return 'https://images.unsplash.com/photo-1534438327276-14e5300c3a48?q=80&w=800&auto=format&fit=crop';
        case 'padel': return 'https://images.unsplash.com/photo-1626248316347-06daec5a0695?q=80&w=800&auto=format&fit=crop'; // Padel-like
        default: return 'https://images.unsplash.com/photo-1599058945522-28d584b6f0ff?q=80&w=800&auto=format&fit=crop'; // Generic Sports
    }
}

export function FacilityCard({ facility }: { facility: Facility }) {
    const isAvailable = facility.status === 'active';
    const imageUrl = getFacilityImage(facility.type);
    const [isBookingModalOpen, setBookingModalOpen] = useState(false);


    return (
        <div className="group bg-white dark:bg-zinc-900 rounded-2xl border border-gray-200 dark:border-zinc-800 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300 overflow-hidden flex flex-col h-full">
            {/* Card Header / Image */}
            <div className="h-48 relative overflow-hidden bg-gray-100 dark:bg-zinc-800">
                <Image
                    src={imageUrl}
                    alt={facility.name}
                    fill
                    className="object-cover transition-transform duration-700 group-hover:scale-110"
                    sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                />

                {/* Gradient Overlay */}
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent z-10" />

                <div className="absolute bottom-4 left-4 z-20 text-white">
                    <h3 className="text-xl font-bold truncate leading-tight tracking-tight">{facility.name}</h3>
                    <p className="text-xs text-gray-200 opacity-90 capitalize font-medium flex items-center gap-1">
                        {facility.type}
                    </p>
                </div>

                {/* Status Badge */}
                <div className="absolute top-4 right-4 z-20 bg-white/10 backdrop-blur-md border border-white/20 px-3 py-1 rounded-full flex items-center gap-2 shadow-lg">
                    <span className={cn(
                        "h-2 w-2 rounded-full shadow-[0_0_8px_rgba(0,0,0,0.5)]",
                        isAvailable ? "bg-emerald-400 shadow-emerald-400/50" : "bg-red-400 shadow-red-400/50"
                    )} />
                    <span className="text-xs font-bold text-white uppercase tracking-wider">
                        {isAvailable ? 'Disponible' : 'Ocupado'}
                    </span>
                </div>
            </div>

            {/* Card Body */}
            <div className="p-5 flex-1 flex flex-col gap-4">
                <div className="grid grid-cols-2 gap-4">
                    <div className="flex items-center gap-2 text-gray-600 dark:text-gray-400">
                        <div className="p-2 bg-brand-50 dark:bg-brand-900/20 rounded-lg text-brand-600 dark:text-brand-400">
                            <Users className="h-4 w-4" />
                        </div>
                        <div className="flex flex-col">
                            <span className="text-[10px] uppercase text-gray-400 font-bold tracking-wider">Capacidad</span>
                            <span className="text-sm font-semibold text-gray-900 dark:text-gray-200">{facility.capacity} Jugadores</span>
                        </div>
                    </div>
                    <div className="flex items-center gap-2 text-gray-600 dark:text-gray-400">
                        <div className="p-2 bg-brand-50 dark:bg-brand-900/20 rounded-lg text-brand-600 dark:text-brand-400">
                            <DollarSign className="h-4 w-4" />
                        </div>
                        <div className="flex flex-col">
                            <span className="text-[10px] uppercase text-gray-400 font-bold tracking-wider">Tarifa</span>
                            <span className="text-sm font-semibold text-gray-900 dark:text-gray-200">${facility.hourly_rate}/h</span>
                        </div>
                    </div>
                </div>

                <div className="flex items-start gap-2 text-gray-500 dark:text-gray-500 mt-2 py-2 border-t border-gray-100 dark:border-zinc-800/50">
                    <MapPin className="h-4 w-4 mt-0.5 shrink-0 text-gray-400" />
                    <span className="text-xs leading-relaxed line-clamp-2 font-medium">
                        {facility.location.name} {facility.location.description ? `- ${facility.location.description}` : ''}
                    </span>
                </div>
            </div>

            {/* Card Footer */}
            <div className="px-5 py-4 bg-gray-50/50 dark:bg-zinc-900/50 flex items-center justify-between border-t border-gray-100 dark:border-zinc-800 gap-2">
                <button
                    onClick={() => setBookingModalOpen(true)}
                    className="text-sm font-semibold text-brand-600 dark:text-brand-400 hover:text-brand-700 dark:hover:text-brand-300 transition-colors"
                >
                    Ver Disponibilidad
                </button>
                <div className="flex-1">
                    <button
                        onClick={() => setBookingModalOpen(true)}
                        className="flex items-center gap-1.5 px-3 py-1.5 bg-brand-600 text-white text-xs font-semibold rounded-lg hover:bg-brand-700 transition w-full justify-center shadow-sm"
                    >
                        <CalendarPlus className="w-3.5 h-3.5" />
                        Reservar
                    </button>
                    <BookingModal
                        isOpen={isBookingModalOpen}
                        onClose={() => setBookingModalOpen(false)}
                        facilityId={facility.id}
                        facilityName={facility.name}
                    />
                </div>
            </div>
        </div>
    )
}
