'use client';

import { useEffect, useRef, useState, useCallback } from 'react';

export interface WebSocketMessage {
    type: string;
    payload: {
        type: string;
        facility_id: string;
        start_time: string;
        end_time: string;
        user_id?: string;
        message?: string;
        action_url?: string;
        timestamp: string;
    };
    timestamp: string;
}

interface UseWebSocketOptions {
    url?: string;
    onMessage?: (message: WebSocketMessage) => void;
    onConnect?: () => void;
    onDisconnect?: () => void;
    reconnectAttempts?: number;
    reconnectInterval?: number;
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
    const {
        url = 'ws://localhost:8081/ws/notifications',
        onMessage,
        onConnect,
        onDisconnect,
        reconnectAttempts = 5,
        reconnectInterval = 3000
    } = options;

    const [isConnected, setIsConnected] = useState(false);
    const [messages, setMessages] = useState<WebSocketMessage[]>([]);
    const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);
    const wsRef = useRef<WebSocket | null>(null);
    const reconnectCountRef = useRef(0);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    const connect = useCallback(() => {
        // Clean up existing connection
        if (wsRef.current?.readyState === WebSocket.OPEN || wsRef.current?.readyState === WebSocket.CONNECTING) {
            return;
        }

        try {
            const ws = new WebSocket(url);

            ws.onopen = () => {
                setIsConnected(true);
                reconnectCountRef.current = 0;
                onConnect?.();
                console.log('WebSocket connected');
            };

            ws.onmessage = (event) => {
                try {
                    const message: WebSocketMessage = JSON.parse(event.data);
                    setLastMessage(message);
                    setMessages((prev) => [...prev.slice(-49), message]); // Keep last 50 messages
                    onMessage?.(message);
                } catch (error) {
                    console.error('Failed to parse WebSocket message:', error);
                }
            };

            ws.onclose = () => {
                setIsConnected(false);
                onDisconnect?.();
                console.log('WebSocket disconnected');

                // Reconnection logic is now handled in useEffect to avoid hoisting issues
            };

            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };

            wsRef.current = ws;
        } catch (error) {
            console.error('Failed to create WebSocket:', error);
        }
    }, [url, onMessage, onConnect, onDisconnect]);

    // Handle Reconnection
    useEffect(() => {
        if (!isConnected && reconnectCountRef.current < reconnectAttempts) {
            const timer = setTimeout(() => {
                reconnectCountRef.current++;
                console.log(`Reconnecting... attempt ${reconnectCountRef.current}`);
                connect();
            }, reconnectInterval);
            return () => clearTimeout(timer);
        }
    }, [isConnected, connect, reconnectAttempts, reconnectInterval]);

    const disconnect = useCallback(() => {
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
        }
        reconnectCountRef.current = reconnectAttempts; // Prevent reconnection
        wsRef.current?.close();
    }, [reconnectAttempts]);

    const subscribe = useCallback((channels: string[]) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                action: 'subscribe',
                targets: channels
            }));
        }
    }, []);

    const unsubscribe = useCallback((channels: string[]) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                action: 'unsubscribe',
                targets: channels
            }));
        }
    }, []);

    const clearMessages = useCallback(() => {
        setMessages([]);
        setLastMessage(null);
    }, []);

    // Auto-connect on mount
    useEffect(() => {
        // WebSocket temporarily disabled
        // connect();
        return () => {
            disconnect();
        };
    }, [connect, disconnect]);

    return {
        isConnected,
        messages,
        lastMessage,
        connect,
        disconnect,
        subscribe,
        unsubscribe,
        clearMessages
    };
}

// Event type constants
export const WS_EVENTS = {
    BOOKING_CANCELLED: 'booking.cancelled',
    SLOT_AVAILABLE: 'slot.available',
    MAINTENANCE_START: 'maintenance.start',
    MAINTENANCE_END: 'maintenance.end'
} as const;
