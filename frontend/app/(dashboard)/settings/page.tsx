'use client';

import { useState } from 'react';
import { Settings, Building2, Bell, Palette, Shield, CreditCard, Users } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

interface SettingSection {
    id: string;
    title: string;
    description: string;
    icon: React.ReactNode;
}

export default function SettingsPage() {
    const [activeSection, setActiveSection] = useState('general');

    const sections: SettingSection[] = [
        {
            id: 'general',
            title: 'General',
            description: 'Información básica del club',
            icon: <Building2 className="h-5 w-5" />,
        },
        {
            id: 'notifications',
            title: 'Notificaciones',
            description: 'Preferencias de alertas',
            icon: <Bell className="h-5 w-5" />,
        },
        {
            id: 'appearance',
            title: 'Apariencia',
            description: 'Tema y personalización',
            icon: <Palette className="h-5 w-5" />,
        },
        {
            id: 'security',
            title: 'Seguridad',
            description: 'Contraseña y accesos',
            icon: <Shield className="h-5 w-5" />,
        },
        {
            id: 'billing',
            title: 'Facturación',
            description: 'Plan y pagos',
            icon: <CreditCard className="h-5 w-5" />,
        },
        {
            id: 'team',
            title: 'Equipo',
            description: 'Gestionar administradores',
            icon: <Users className="h-5 w-5" />,
        },
    ];

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                    <Settings className="h-6 w-6 text-brand-600" />
                    Configuración
                </h1>
                <p className="text-gray-500 text-sm mt-1">
                    Administra la configuración de tu club.
                </p>
            </div>

            <div className="grid gap-6 md:grid-cols-[250px_1fr]">
                {/* Sidebar Navigation */}
                <Card className="h-fit">
                    <CardContent className="p-2">
                        <nav className="space-y-1">
                            {sections.map((section) => (
                                <button
                                    key={section.id}
                                    onClick={() => setActiveSection(section.id)}
                                    className={`w-full flex items-center gap-3 px-3 py-2 text-sm rounded-lg transition-colors ${activeSection === section.id
                                            ? 'bg-brand-50 text-brand-700 dark:bg-brand-900/20'
                                            : 'text-gray-600 hover:bg-gray-100 dark:hover:bg-gray-800'
                                        }`}
                                >
                                    {section.icon}
                                    <span>{section.title}</span>
                                </button>
                            ))}
                        </nav>
                    </CardContent>
                </Card>

                {/* Content Area */}
                <div className="space-y-6">
                    {activeSection === 'general' && (
                        <Card>
                            <CardHeader>
                                <CardTitle>Información del Club</CardTitle>
                                <CardDescription>
                                    Actualiza los datos básicos de tu club.
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Nombre del Club
                                    </label>
                                    <input
                                        type="text"
                                        defaultValue="Club Pulse Demo"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Dirección
                                    </label>
                                    <input
                                        type="text"
                                        defaultValue="Av. Principal 123"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Teléfono de Contacto
                                    </label>
                                    <input
                                        type="tel"
                                        defaultValue="+54 11 1234-5678"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Email de Contacto
                                    </label>
                                    <input
                                        type="email"
                                        defaultValue="contacto@clubpulse.com"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                                    />
                                </div>
                                <div className="pt-4">
                                    <Button className="bg-brand-600 hover:bg-brand-700 text-white">
                                        Guardar Cambios
                                    </Button>
                                </div>
                            </CardContent>
                        </Card>
                    )}

                    {activeSection === 'notifications' && (
                        <Card>
                            <CardHeader>
                                <CardTitle>Preferencias de Notificaciones</CardTitle>
                                <CardDescription>
                                    Configura qué notificaciones deseas recibir.
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                {[
                                    { label: 'Nuevas reservas', defaultChecked: true },
                                    { label: 'Cancelaciones', defaultChecked: true },
                                    { label: 'Nuevos miembros', defaultChecked: true },
                                    { label: 'Recordatorios de pago', defaultChecked: false },
                                    { label: 'Actualizaciones del sistema', defaultChecked: false },
                                ].map((item) => (
                                    <div key={item.label} className="flex items-center justify-between py-2">
                                        <span className="text-sm text-gray-700">{item.label}</span>
                                        <label className="relative inline-flex items-center cursor-pointer">
                                            <input
                                                type="checkbox"
                                                defaultChecked={item.defaultChecked}
                                                className="sr-only peer"
                                            />
                                            <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-brand-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-brand-600"></div>
                                        </label>
                                    </div>
                                ))}
                            </CardContent>
                        </Card>
                    )}

                    {activeSection === 'appearance' && (
                        <Card>
                            <CardHeader>
                                <CardTitle>Tema y Apariencia</CardTitle>
                                <CardDescription>
                                    Personaliza la apariencia de tu panel.
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-6">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-3">
                                        Tema
                                    </label>
                                    <div className="flex gap-4">
                                        {['Claro', 'Oscuro', 'Sistema'].map((theme) => (
                                            <button
                                                key={theme}
                                                className={`px-4 py-2 border rounded-lg text-sm transition-colors ${theme === 'Claro'
                                                        ? 'border-brand-600 bg-brand-50 text-brand-700'
                                                        : 'border-gray-300 hover:border-gray-400'
                                                    }`}
                                            >
                                                {theme}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-3">
                                        Color Principal
                                    </label>
                                    <div className="flex gap-3">
                                        {['#3B82F6', '#10B981', '#8B5CF6', '#F59E0B', '#EF4444'].map((color) => (
                                            <button
                                                key={color}
                                                className="w-8 h-8 rounded-full ring-2 ring-offset-2 ring-transparent hover:ring-gray-300 transition-all"
                                                style={{ backgroundColor: color }}
                                            />
                                        ))}
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    )}

                    {activeSection === 'security' && (
                        <Card>
                            <CardHeader>
                                <CardTitle>Seguridad</CardTitle>
                                <CardDescription>
                                    Gestiona la seguridad de tu cuenta.
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Contraseña Actual
                                    </label>
                                    <input
                                        type="password"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Nueva Contraseña
                                    </label>
                                    <input
                                        type="password"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Confirmar Nueva Contraseña
                                    </label>
                                    <input
                                        type="password"
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-500"
                                    />
                                </div>
                                <div className="pt-4">
                                    <Button className="bg-brand-600 hover:bg-brand-700 text-white">
                                        Actualizar Contraseña
                                    </Button>
                                </div>
                            </CardContent>
                        </Card>
                    )}

                    {activeSection === 'billing' && (
                        <Card>
                            <CardHeader>
                                <CardTitle>Facturación y Plan</CardTitle>
                                <CardDescription>
                                    Gestiona tu suscripción y métodos de pago.
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="bg-brand-50 dark:bg-brand-900/20 p-4 rounded-lg mb-6">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <h3 className="font-semibold text-brand-700">Plan Pro</h3>
                                            <p className="text-sm text-gray-600">Próxima facturación: 1 Feb 2026</p>
                                        </div>
                                        <span className="text-2xl font-bold text-brand-600">$99/mes</span>
                                    </div>
                                </div>
                                <Button variant="outline" className="w-full">
                                    Cambiar Plan
                                </Button>
                            </CardContent>
                        </Card>
                    )}

                    {activeSection === 'team' && (
                        <Card>
                            <CardHeader>
                                <CardTitle>Gestión de Equipo</CardTitle>
                                <CardDescription>
                                    Administra los usuarios con acceso al panel.
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-brand-100 rounded-full flex items-center justify-center">
                                                <span className="text-brand-700 font-medium">SA</span>
                                            </div>
                                            <div>
                                                <p className="font-medium text-sm">System Admin</p>
                                                <p className="text-xs text-gray-500">admin@clubpulse.com</p>
                                            </div>
                                        </div>
                                        <span className="text-xs bg-brand-100 text-brand-700 px-2 py-1 rounded">
                                            Propietario
                                        </span>
                                    </div>
                                </div>
                                <Button className="mt-4 w-full" variant="outline">
                                    + Invitar Administrador
                                </Button>
                            </CardContent>
                        </Card>
                    )}
                </div>
            </div>
        </div>
    );
}
