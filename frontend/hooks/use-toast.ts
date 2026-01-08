"use client"

import { useState, useCallback } from "react"

interface Toast {
    id: string
    title?: string
    description?: string
    variant?: "default" | "destructive"
}

interface ToastState {
    toasts: Toast[]
}

// Simple toast hook without provider (stores state locally per component tree)
// For a full implementation, you'd want a global state manager or context

let globalToasts: Toast[] = []
let listeners: Array<(toasts: Toast[]) => void> = []

function emitChange() {
    for (const listener of listeners) {
        listener(globalToasts)
    }
}

export function useToast() {
    const [toasts, setToasts] = useState<Toast[]>(globalToasts)

    // Subscribe to changes
    useState(() => {
        const listener = (newToasts: Toast[]) => setToasts([...newToasts])
        listeners.push(listener)
        return () => {
            listeners = listeners.filter(l => l !== listener)
        }
    })

    const toast = useCallback((props: Omit<Toast, "id">) => {
        const id = Math.random().toString(36).substring(2, 9)
        const newToast = { ...props, id }
        globalToasts = [...globalToasts, newToast]
        emitChange()

        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            globalToasts = globalToasts.filter((t) => t.id !== id)
            emitChange()
        }, 5000)
    }, [])

    const dismiss = useCallback((id: string) => {
        globalToasts = globalToasts.filter((t) => t.id !== id)
        emitChange()
    }, [])

    return {
        toast,
        toasts,
        dismiss,
    }
}
