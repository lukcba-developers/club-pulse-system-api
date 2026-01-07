"use client"

// Simplified toast hook for immediate usage
import { useState, useEffect } from "react"

type ToastProps = {
    title?: string
    description?: string
    variant?: "default" | "destructive"
}


export function useToast() {
    const toast = ({ title, description, variant }: ToastProps) => {
        const event = new CustomEvent("toast-trigger", {
            detail: { title, description, variant }
        });
        window.dispatchEvent(event);
    }
    return { toast }
}

// Minimal Toaster component to be placed in Layout

export function Toaster() {
    const [toasts, setToasts] = useState<ToastProps[]>([])

    useEffect(() => {
        const handleToast = (event: Event) => {
            const customEvent = event as CustomEvent<ToastProps>;
            setToasts(prev => [...prev, customEvent.detail]);

            // Auto remove after 3s
            setTimeout(() => {
                setToasts(prev => prev.slice(1));
            }, 3000);
        };

        window.addEventListener("toast-trigger", handleToast);
        return () => window.removeEventListener("toast-trigger", handleToast);
    }, []);

    if (toasts.length === 0) return null;

    return (
        <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
            {toasts.map((t, i) => (
                <div
                    key={i}
                    className={`p-4 rounded-md shadow-lg border w-80 bg-white dark:bg-zinc-900 ${t.variant === 'destructive' ? 'border-red-500 text-red-500' : 'border-gray-200 dark:border-zinc-800'}`}
                >
                    {t.title && <h4 className="font-semibold mb-1">{t.title}</h4>}
                    {t.description && <p className="text-sm text-gray-500">{t.description}</p>}
                </div>
            ))}
        </div>
    )
}
