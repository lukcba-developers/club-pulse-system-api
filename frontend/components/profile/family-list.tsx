'use client';

import React, { useEffect, useState } from 'react';
import { userService, User } from '@/services/user-service';
import { Baby, ChevronRight } from 'lucide-react';

export function FamilyList() {
    const [children, setChildren] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [newChildName, setNewChildName] = useState('');
    const [newChildDob, setNewChildDob] = useState('');
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        loadChildren();
    }, []);

    const loadChildren = () => {
        setLoading(true);
        userService.getChildren()
            .then(setChildren)
            .catch(console.error)
            .finally(() => setLoading(false));
    };

    const handleAddChild = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        try {
            await userService.registerChild(newChildName, newChildDob);
            setShowModal(false);
            setNewChildName('');
            setNewChildDob('');
            loadChildren();
        } catch (error) {
            console.error("Failed to add child", error);
            alert("Error al agregar hijo. Intente nuevamente.");
        } finally {
            setSubmitting(false);
        }
    }

    if (loading && children.length === 0) return <div className="text-gray-500">Cargando grupo familiar...</div>;

    return (
        <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
            <div className="px-4 py-5 sm:px-6 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
                <div className="flex items-center gap-2">
                    <Baby className="h-5 w-5 text-brand-500" />
                    <h3 className="text-lg font-medium leading-6 text-gray-900 dark:text-white">
                        Grupo Familiar
                    </h3>
                </div>
                <button
                    onClick={() => setShowModal(true)}
                    className="text-sm bg-brand-600 text-white px-3 py-1.5 rounded-md hover:bg-brand-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500"
                >
                    Agregar Hijo
                </button>
            </div>

            {children.length > 0 ? (
                <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                    {children.map((child) => (
                        <li key={child.id} className="px-4 py-4 sm:px-6 hover:bg-gray-50 dark:hover:bg-gray-750 transition flex items-center justify-between">
                            <div className="flex items-center">
                                <div className="h-10 w-10 rounded-full bg-brand-100 flex items-center justify-center text-brand-700 font-bold">
                                    {child.name.charAt(0)}
                                </div>
                                <div className="ml-4">
                                    <p className="text-sm font-medium text-gray-900 dark:text-white">{child.name}</p>
                                    <p className="text-xs text-gray-500 dark:text-gray-400">{child.email || 'Sin email'}</p>
                                </div>
                            </div>
                            <ChevronRight className="h-5 w-5 text-gray-400" />
                        </li>
                    ))}
                </ul>
            ) : (
                <div className="p-6 text-center text-gray-500">No tienes hijos registrados.</div>
            )}

            {showModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6">
                        <h2 className="text-xl font-bold mb-4 text-gray-900 dark:text-white">Agregar Hijo</h2>
                        <form onSubmit={handleAddChild} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Nombre Completo</label>
                                <input
                                    type="text"
                                    required
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-brand-500 focus:ring-brand-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white p-2 border"
                                    value={newChildName}
                                    onChange={e => setNewChildName(e.target.value)}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Fecha de Nacimiento</label>
                                <input
                                    type="date"
                                    required
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-brand-500 focus:ring-brand-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white p-2 border"
                                    value={newChildDob}
                                    onChange={e => setNewChildDob(e.target.value)}
                                />
                            </div>
                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    type="button"
                                    onClick={() => setShowModal(false)}
                                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500"
                                >
                                    Cancelar
                                </button>
                                <button
                                    type="submit"
                                    disabled={submitting}
                                    className="px-4 py-2 text-sm font-medium text-white bg-brand-600 border border-transparent rounded-md hover:bg-brand-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500 disabled:opacity-50"
                                >
                                    {submitting ? 'Guardando...' : 'Guardar'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
