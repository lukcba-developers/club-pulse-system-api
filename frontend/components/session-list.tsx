'use client';

import { useState, useEffect } from 'react';
import { authService, Session } from '@/services/auth-service';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Trash2, Monitor, Loader2 } from 'lucide-react';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog';

export function SessionList() {
    const [sessions, setSessions] = useState<Session[]>([]);
    const [loading, setLoading] = useState(true);
    const [revokingId, setRevokingId] = useState<string | null>(null);

    const fetchSessions = async () => {
        try {
            const data = await authService.listSessions();
            setSessions(data || []);
        } catch (error) {
            console.error('Error fetching sessions:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleRevoke = async (id: string) => {
        setRevokingId(id);
        try {
            await authService.revokeSession(id);
            setSessions(prev => prev.filter(s => s.id !== id));
        } catch (error) {
            console.error('Error revoking session:', error);
        } finally {
            setRevokingId(null);
        }
    };

    useEffect(() => {
        fetchSessions();
    }, []);

    if (loading) return <div>Cargando sesiones...</div>;

    return (
        <Card>
            <CardHeader>
                <CardTitle>Sesiones Activas</CardTitle>
                <CardDescription>Gestiona tus dispositivos y sesiones activas.</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {sessions.length === 0 ? (
                        <p className="text-muted-foreground text-sm">No se encontraron sesiones activas.</p>
                    ) : (
                        sessions.map((session) => (
                            <div key={session.id} className="flex items-center justify-between border-b pb-4 last:border-0 last:pb-0">
                                <div className="flex items-center space-x-4">
                                    <div className="bg-secondary p-2 rounded-full">
                                        <Monitor className="h-5 w-5 text-primary" />
                                    </div>
                                    <div>
                                        <p className="text-sm font-medium">
                                            {session.device_id || 'Dispositivo Desconocido'}
                                        </p>
                                        <p className="text-xs text-muted-foreground">
                                            Creado: {new Date(session.created_at).toLocaleDateString()}
                                        </p>
                                    </div>
                                </div>
                                <AlertDialog>
                                    <AlertDialogTrigger asChild>
                                        <Button variant="ghost" size="sm" disabled={revokingId === session.id}>
                                            {revokingId === session.id ? <Loader2 className="h-4 w-4 animate-spin text-destructive" /> : <Trash2 className="h-4 w-4 text-destructive" />}
                                        </Button>
                                    </AlertDialogTrigger>
                                    <AlertDialogContent>
                                        <AlertDialogHeader>
                                            <AlertDialogTitle>¿Revocar Sesión?</AlertDialogTitle>
                                            <AlertDialogDescription>
                                                Esta acción cerrará la sesión en el dispositivo seleccionado. ¿Estás seguro?
                                            </AlertDialogDescription>
                                        </AlertDialogHeader>
                                        <AlertDialogFooter>
                                            <AlertDialogCancel>Cancelar</AlertDialogCancel>
                                            <AlertDialogAction onClick={() => handleRevoke(session.id)} className="bg-red-600 hover:bg-red-700">
                                                Sí, revocar
                                            </AlertDialogAction>
                                        </AlertDialogFooter>
                                    </AlertDialogContent>
                                </AlertDialog>
                            </div>
                        ))
                    )}
                </div>
            </CardContent>
        </Card>
    );
}
