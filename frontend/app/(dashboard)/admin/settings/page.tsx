'use client';

import { useState } from 'react';
import { Settings, Server, Database, Shield, Mail, Globe } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export default function AdminSettingsPage() {
    const [activeTab, setActiveTab] = useState('system');

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                    <Settings className="h-6 w-6 text-brand-600" />
                    Configuración Global
                </h1>
                <p className="text-gray-500 text-sm mt-1">
                    Configuración a nivel de plataforma (solo Super Admin).
                </p>
            </div>

            {/* Tabs */}
            <div className="flex gap-2 border-b">
                {[
                    { id: 'system', label: 'Sistema', icon: <Server className="h-4 w-4" /> },
                    { id: 'database', label: 'Base de Datos', icon: <Database className="h-4 w-4" /> },
                    { id: 'security', label: 'Seguridad', icon: <Shield className="h-4 w-4" /> },
                    { id: 'email', label: 'Email', icon: <Mail className="h-4 w-4" /> },
                    { id: 'domains', label: 'Dominios', icon: <Globe className="h-4 w-4" /> },
                ].map((tab) => (
                    <button
                        key={tab.id}
                        onClick={() => setActiveTab(tab.id)}
                        className={`flex items-center gap-2 px-4 py-2 text-sm transition-colors border-b-2 -mb-[2px] ${activeTab === tab.id
                                ? 'border-brand-600 text-brand-600'
                                : 'border-transparent text-gray-600 hover:text-gray-900'
                            }`}
                    >
                        {tab.icon}
                        {tab.label}
                    </button>
                ))}
            </div>

            {/* System Settings */}
            {activeTab === 'system' && (
                <div className="grid gap-6 md:grid-cols-2">
                    <Card>
                        <CardHeader>
                            <CardTitle>Estado del Sistema</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-3">
                            <div className="flex justify-between">
                                <span className="text-sm text-gray-600">API Status</span>
                                <span className="text-sm text-green-600 font-medium">● Operativo</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-sm text-gray-600">Database</span>
                                <span className="text-sm text-green-600 font-medium">● Conectado</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-sm text-gray-600">Redis Cache</span>
                                <span className="text-sm text-green-600 font-medium">● Activo</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-sm text-gray-600">Versión</span>
                                <span className="text-sm font-mono">v1.0.0</span>
                            </div>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle>Modo Mantenimiento</CardTitle>
                            <CardDescription>
                                Activar para realizar actualizaciones.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center justify-between">
                                <span className="text-sm">Estado actual: Desactivado</span>
                                <Button variant="outline" size="sm">
                                    Activar
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            )}

            {/* Database Settings */}
            {activeTab === 'database' && (
                <Card>
                    <CardHeader>
                        <CardTitle>Configuración de Base de Datos</CardTitle>
                        <CardDescription>
                            Gestión de conexiones y backups.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="p-4 bg-gray-50 rounded-lg">
                            <h3 className="font-medium text-sm mb-2">Conexión Actual</h3>
                            <p className="text-xs text-gray-600 font-mono">PostgreSQL @ localhost:5432</p>
                        </div>
                        <div className="flex gap-2">
                            <Button variant="outline">Ejecutar Backup</Button>
                            <Button variant="outline">Ver Logs</Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Security Settings */}
            {activeTab === 'security' && (
                <Card>
                    <CardHeader>
                        <CardTitle>Configuración de Seguridad</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="flex items-center justify-between py-2">
                            <div>
                                <p className="text-sm font-medium">Rate Limiting</p>
                                <p className="text-xs text-gray-500">Limitar peticiones por IP</p>
                            </div>
                            <label className="relative inline-flex items-center cursor-pointer">
                                <input type="checkbox" defaultChecked className="sr-only peer" />
                                <div className="w-11 h-6 bg-gray-200 rounded-full peer peer-checked:bg-brand-600 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:after:translate-x-full"></div>
                            </label>
                        </div>
                        <div className="flex items-center justify-between py-2">
                            <div>
                                <p className="text-sm font-medium">2FA Obligatorio</p>
                                <p className="text-xs text-gray-500">Para todos los admins</p>
                            </div>
                            <label className="relative inline-flex items-center cursor-pointer">
                                <input type="checkbox" className="sr-only peer" />
                                <div className="w-11 h-6 bg-gray-200 rounded-full peer peer-checked:bg-brand-600 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:after:translate-x-full"></div>
                            </label>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Email Settings */}
            {activeTab === 'email' && (
                <Card>
                    <CardHeader>
                        <CardTitle>Configuración de Email</CardTitle>
                        <CardDescription>
                            SMTP y plantillas de correo.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium mb-1">SMTP Host</label>
                            <input
                                type="text"
                                defaultValue="smtp.example.com"
                                className="w-full px-3 py-2 border rounded-lg"
                            />
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">Puerto</label>
                                <input
                                    type="number"
                                    defaultValue="587"
                                    className="w-full px-3 py-2 border rounded-lg"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium mb-1">Seguridad</label>
                                <select className="w-full px-3 py-2 border rounded-lg">
                                    <option>TLS</option>
                                    <option>SSL</option>
                                    <option>Ninguna</option>
                                </select>
                            </div>
                        </div>
                        <Button className="bg-brand-600 hover:bg-brand-700 text-white">
                            Guardar Configuración
                        </Button>
                    </CardContent>
                </Card>
            )}

            {/* Domains Settings */}
            {activeTab === 'domains' && (
                <Card>
                    <CardHeader>
                        <CardTitle>Dominios Permitidos</CardTitle>
                        <CardDescription>
                            Gestiona los dominios autorizados para la plataforma.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2 mb-4">
                            {['clubpulse.com', 'app.clubpulse.com', 'localhost:3000'].map((domain) => (
                                <div key={domain} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                                    <span className="text-sm font-mono">{domain}</span>
                                    <Button variant="ghost" size="sm" className="text-red-600">
                                        Eliminar
                                    </Button>
                                </div>
                            ))}
                        </div>
                        <Button variant="outline">+ Agregar Dominio</Button>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
