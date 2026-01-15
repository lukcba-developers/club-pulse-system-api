'use client';

import { useState } from 'react';
import { Bell, X, Calendar, Clock } from 'lucide-react';
import { useWebSocket, WS_EVENTS, WebSocketMessage } from '@/hooks/use-websocket';

interface Notification {
    id: string;
    type: string;
    title: string;
    message: string;
    facilityId?: string;
    actionUrl?: string;
    timestamp: Date;
    read: boolean;
}

let notificationCounter = 0;

export function NotificationToast() {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [showPanel, setShowPanel] = useState(false);
    const { isConnected } = useWebSocket({
        onMessage: (msg) => handleNewMessage(msg)
    });

    const handleNewMessage = (message: WebSocketMessage) => {
        const notification: Notification = {
            id: `${message.timestamp}-${notificationCounter++}`,
            type: message.payload.type,
            title: getNotificationTitle(message.payload.type),
            message: message.payload.message || getDefaultMessage(message.payload.type),
            facilityId: message.payload.facility_id,
            actionUrl: message.payload.action_url,
            timestamp: new Date(message.timestamp),
            read: false
        };

        setNotifications((prev) => [notification, ...prev.slice(0, 9)]); // Keep last 10

        // Show toast briefly
        showToast(notification);
    };

    const showToast = (notification: Notification) => {
        // Create toast element
        const toast = document.createElement('div');
        toast.className = 'fixed bottom-4 right-4 z-50 animate-slide-up';
        toast.innerHTML = `
            <div class="bg-white dark:bg-zinc-900 rounded-xl shadow-2xl border border-gray-200 dark:border-zinc-700 p-4 max-w-sm cursor-pointer" onclick="${notification.actionUrl ? `window.location.href='${notification.actionUrl}'` : ''}">
                <div class="flex items-start gap-3">
                    <div class="flex-shrink-0 w-10 h-10 rounded-full bg-gradient-to-r from-brand-500 to-purple-500 flex items-center justify-center">
                        <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                    </div>
                    <div class="flex-1">
                        <h4 class="font-semibold text-gray-900 dark:text-gray-100">${notification.title}</h4>
                        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">${notification.message}</p>
                        ${notification.actionUrl ? `<p class="text-xs text-indigo-500 mt-2 font-medium">Click to view</p>` : ''}
                    </div>
                </div>
            </div>
        `;

        document.body.appendChild(toast);

        // Remove after 5 seconds
        setTimeout(() => {
            toast.classList.add('animate-fade-out');
            setTimeout(() => toast.remove(), 300);
        }, 5000);
    };

    const getNotificationTitle = (type: string): string => {
        switch (type) {
            case WS_EVENTS.BOOKING_CANCELLED:
                return '¡Pista Disponible!';
            case WS_EVENTS.SLOT_AVAILABLE:
                return 'Nueva Disponibilidad';
            case WS_EVENTS.MAINTENANCE_START:
                return 'Mantenimiento Iniciado';
            case WS_EVENTS.MAINTENANCE_END:
                return 'Mantenimiento Finalizado';
            default:
                return 'Notificación';
        }
    };

    const getDefaultMessage = (type: string): string => {
        switch (type) {
            case WS_EVENTS.BOOKING_CANCELLED:
                return 'Se ha liberado una pista. ¡Reserva ahora!';
            case WS_EVENTS.SLOT_AVAILABLE:
                return 'Hay nuevos horarios disponibles';
            default:
                return 'Hay actualizaciones en el sistema';
        }
    };

    const markAsRead = (id: string) => {
        setNotifications((prev) =>
            prev.map((n) => (n.id === id ? { ...n, read: true } : n))
        );
    };

    const clearAll = () => {
        setNotifications([]);
    };

    const unreadCount = notifications.filter((n) => !n.read).length;

    return (
        <>
            {/* Notification Bell Button */}
            <button
                onClick={() => setShowPanel(!showPanel)}
                aria-label={`Notificaciones${unreadCount > 0 ? `, ${unreadCount} sin leer` : ''}`}
                aria-expanded={showPanel}
                aria-haspopup="true"
                className="relative p-2 rounded-lg text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-zinc-800 transition-colors"
            >
                <Bell className="h-5 w-5" aria-hidden="true" />
                {unreadCount > 0 && (
                    <span
                        className="absolute -top-1 -right-1 h-5 w-5 rounded-full bg-red-500 text-white text-xs flex items-center justify-center font-medium"
                        aria-hidden="true"
                    >
                        {unreadCount}
                    </span>
                )}
                {/* Connection indicator */}
                <span
                    className={`absolute bottom-0 right-0 h-2 w-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-gray-400'}`}
                    aria-hidden="true"
                />
            </button>

            {/* Notification Panel */}
            {showPanel && (
                <div className="absolute right-0 top-full mt-2 w-80 bg-white dark:bg-zinc-900 rounded-xl shadow-xl border border-gray-200 dark:border-zinc-700 overflow-hidden z-50">
                    {/* Header */}
                    <div className="flex items-center justify-between px-4 py-3 border-b border-gray-100 dark:border-zinc-800">
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100">Notificaciones</h3>
                        <div className="flex items-center gap-2">
                            {notifications.length > 0 && (
                                <button
                                    onClick={clearAll}
                                    className="text-xs text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
                                >
                                    Limpiar
                                </button>
                            )}
                            <button
                                onClick={() => setShowPanel(false)}
                                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                            >
                                <X className="h-4 w-4" />
                            </button>
                        </div>
                    </div>

                    {/* Notifications List */}
                    <div className="max-h-96 overflow-y-auto">
                        {notifications.length === 0 ? (
                            <div className="p-8 text-center text-gray-500">
                                <Bell className="h-8 w-8 mx-auto mb-2 opacity-50" />
                                <p className="text-sm">No hay notificaciones</p>
                            </div>
                        ) : (
                            <ul className="divide-y divide-gray-100 dark:divide-zinc-800">
                                {notifications.map((notification) => (
                                    <li
                                        key={notification.id}
                                        onClick={() => markAsRead(notification.id)}
                                        className={`p-4 hover:bg-gray-50 dark:hover:bg-zinc-800 cursor-pointer transition-colors ${!notification.read ? 'bg-blue-50/50 dark:bg-blue-900/10' : ''
                                            }`}
                                    >
                                        <div className="flex items-start gap-3">
                                            <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${notification.type === WS_EVENTS.BOOKING_CANCELLED
                                                ? 'bg-green-100 dark:bg-green-900/30 text-green-600'
                                                : 'bg-blue-100 dark:bg-blue-900/30 text-blue-600'
                                                }`}>
                                                <Calendar className="h-4 w-4" />
                                            </div>
                                            <div className="flex-1 min-w-0">
                                                <p className="font-medium text-sm text-gray-900 dark:text-gray-100">
                                                    {notification.title}
                                                </p>
                                                <p className="text-xs text-gray-500 mt-1 truncate">
                                                    {notification.message}
                                                </p>
                                                <div className="flex items-center gap-1 mt-1 text-xs text-gray-400">
                                                    <Clock className="h-3 w-3" />
                                                    {notification.timestamp.toLocaleTimeString()}
                                                </div>
                                            </div>
                                            {!notification.read && (
                                                <div className="w-2 h-2 rounded-full bg-blue-500" />
                                            )}
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>

                    {/* Footer */}
                    <div className="px-4 py-2 border-t border-gray-100 dark:border-zinc-800 bg-gray-50 dark:bg-zinc-800">
                        <div className="flex items-center justify-center gap-2 text-xs text-gray-500">
                            <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
                            {isConnected ? 'Conectado en tiempo real' : 'Reconectando...'}
                        </div>
                    </div>
                </div>
            )}

            {/* Add animation styles */}
            <style jsx global>{`
                @keyframes slide-up {
                    from {
                        opacity: 0;
                        transform: translateY(20px);
                    }
                    to {
                        opacity: 1;
                        transform: translateY(0);
                    }
                }
                @keyframes fade-out {
                    to {
                        opacity: 0;
                        transform: translateY(20px);
                    }
                }
                .animate-slide-up {
                    animation: slide-up 0.3s ease-out;
                }
                .animate-fade-out {
                    animation: fade-out 0.3s ease-out forwards;
                }
            `}</style>
        </>
    );
}
