'use client';

import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAuth } from '@/hooks/use-auth';
import { ShieldAlert } from 'lucide-react';
import { useState } from 'react';
import { userService } from '@/services/user-service';
import { useToast } from '@/components/ui/use-toast';

export function HealthSection() {
    const { user } = useAuth();
    const { toast } = useToast();
    const [loading, setLoading] = useState(false);

    const [emergencyForm, setEmergencyForm] = useState({
        contact_name: user?.emergency_contact_name || '',
        contact_phone: user?.emergency_contact_phone || '',
        insurance_provider: user?.insurance_provider || '',
        insurance_number: user?.insurance_number || ''
    });

    if (!user) return null;

    const handleUpdateEmergency = async () => {
        try {
            setLoading(true);
            await userService.updateEmergencyInfo(emergencyForm);
            toast({
                title: "Información Actualizada",
                description: "Los datos de emergencia han sido guardados.",
            });
        } catch (err) {
            console.error(err);
            toast({
                title: "Error",
                description: "No se pudo actualizar la información.",
                variant: 'destructive'
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <div className="space-y-1">
                    <CardTitle className="text-2xl font-bold">Salud y Seguridad</CardTitle>
                    <CardDescription>Estado médico y contactos de emergencia</CardDescription>
                </div>
                <ShieldAlert className="h-8 w-8 text-muted-foreground" />
            </CardHeader>
            <CardContent>
                {status !== 'VALID' && (
                    <div className="rounded-md bg-blue-50 dark:bg-blue-900/20 p-4">
                        <div className="flex">
                            <div className="ml-3">
                                <h3 className="text-sm font-medium text-blue-800 dark:text-blue-300">
                                    Acción Requerida
                                </h3>
                                <div className="mt-2 text-sm text-blue-700 dark:text-blue-400">
                                    <p>
                                        Por favor cerca tu certificado médico actualizado a la administración del club para habilitar reservas.
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
                {/* Add Emergency Contact Form Here */}
                <div className="grid gap-4 py-4">
                    <div className="space-y-2">
                        <Label>Nombre Contacto Emergencia</Label>
                        <Input
                            value={emergencyForm.contact_name}
                            onChange={(e) => setEmergencyForm({ ...emergencyForm, contact_name: e.target.value })}
                        />
                    </div>
                    <div className="space-y-2">
                        <Label>Teléfono Contacto</Label>
                        <Input
                            value={emergencyForm.contact_phone}
                            onChange={(e) => setEmergencyForm({ ...emergencyForm, contact_phone: e.target.value })}
                        />
                    </div>
                    <div className="space-y-2">
                        <Label>Obra Social / Seguro</Label>
                        <Input
                            value={emergencyForm.insurance_provider}
                            onChange={(e) => setEmergencyForm({ ...emergencyForm, insurance_provider: e.target.value })}
                        />
                    </div>
                    <div className="space-y-2">
                        <Label>Nro. Afiliado / Póliza</Label>
                        <Input
                            value={emergencyForm.insurance_number}
                            onChange={(e) => setEmergencyForm({ ...emergencyForm, insurance_number: e.target.value })}
                        />
                    </div>
                    <Button onClick={handleUpdateEmergency} disabled={loading}>
                        {loading ? 'Guardando...' : 'Guardar Información'}
                    </Button>
                </div>
            </CardContent>
        </Card>
    );
}
