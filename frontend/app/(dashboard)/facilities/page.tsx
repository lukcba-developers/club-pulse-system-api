'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Plus, MapPin, Users, DollarSign } from 'lucide-react';
import Link from 'next/link';
import api from '@/lib/axios';
import { BookingModal } from '@/components/booking-modal';

interface Facility {
    id: string;
    name: string;
    type: string;
    status: string;
    capacity: number;
    hourly_rate: number;
    opening_time: string;
    closing_time: string;
    location: {
        address?: string;
    };
}

export default function FacilitiesPage() {
    const [facilities, setFacilities] = useState<Facility[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedFacility, setSelectedFacility] = useState<Facility | null>(null);
    const [isBookingModalOpen, setIsBookingModalOpen] = useState(false);

    const handleBooking = (facility: Facility) => {
        setSelectedFacility(facility);
        setIsBookingModalOpen(true);
    };

    useEffect(() => {
        const fetchFacilities = async () => {
            try {
                const res = await api.get('/facilities');
                const facilitiesData = Array.isArray(res.data) ? res.data : (res.data.data || []);
                setFacilities(facilitiesData);
            } catch (error) {
                console.error("Failed to fetch facilities", error);
            } finally {
                setLoading(false);
            }
        };
        fetchFacilities();
    }, []);

    if (loading) return <div>Cargando instalaciones...</div>;

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Instalaciones</h1>
                    <p className="text-muted-foreground">Gestiona las canchas y espacios del club.</p>
                </div>
                <Link href="/facilities/create">
                    <Button>
                        <Plus className="mr-2 h-4 w-4" />
                        Nueva Instalación
                    </Button>
                </Link>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {facilities.map((fac) => (
                    <Card key={fac.id}>
                        <CardHeader>
                            <div className="flex justify-between items-start">
                                <div>
                                    <CardTitle className="text-xl">{fac.name}</CardTitle>
                                    <CardDescription>{fac.type}</CardDescription>
                                </div>
                                <Badge variant={fac.status === 'active' ? 'default' : 'secondary'}>
                                    {fac.status}
                                </Badge>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                <div className="space-y-2 text-sm">
                                    <div className="flex items-center text-muted-foreground">
                                        <Users className="mr-2 h-4 w-4" />
                                        Capacidad: {fac.capacity} personas
                                    </div>
                                    <div className="flex items-center text-muted-foreground">
                                        <DollarSign className="mr-2 h-4 w-4" />
                                        Precio: ${fac.hourly_rate}/hora
                                    </div>
                                    <div className="flex items-center text-muted-foreground">
                                        <span className="mr-2">🕐</span>
                                        Horario: {fac.opening_time} - {fac.closing_time}
                                    </div>
                                    {fac.location?.address && (
                                        <div className="flex items-center text-muted-foreground">
                                            <MapPin className="mr-2 h-4 w-4" />
                                            {fac.location.address}
                                        </div>
                                    )}
                                </div>
                                <Button
                                    className="w-full"
                                    onClick={() => handleBooking(fac)}
                                    disabled={fac.status !== 'active'}
                                >
                                    Reservar
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
            {facilities.length === 0 && (
                <div className="text-center p-8 text-muted-foreground">No hay instalaciones registradas</div>
            )}
            {selectedFacility && (
                <BookingModal
                    isOpen={isBookingModalOpen}
                    onClose={() => setIsBookingModalOpen(false)}
                    facilityId={selectedFacility.id}
                    facilityName={selectedFacility.name}
                />
            )}
        </div>
    );
}
