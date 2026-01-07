'use client';

import { useEffect } from 'react';

export default function Error({
    error,
    reset,
}: {
    error: Error & { digest?: string };
    reset: () => void;
}) {
    useEffect(() => {
        // Log the error to an error reporting service
        console.error('[Global Error Boundary]:', error);
    }, [error]);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-neutral-50 p-4 text-center">
            <h2 className="text-2xl font-bold text-neutral-800 mb-4">
                Something went wrong!
            </h2>
            <p className="text-neutral-600 mb-6 max-w-md">
                We encountered an unexpected error. Check the console for more details.
            </p>
            <div className="flex gap-4">
                <button
                    onClick={
                        // Attempt to recover by trying to re-render the segment
                        () => reset()
                    }
                    className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition"
                >
                    Try again
                </button>
                <button
                    onClick={() => window.location.reload()}
                    className="px-4 py-2 border border-neutral-300 bg-white text-neutral-700 rounded-md hover:bg-neutral-50 transition"
                >
                    Reload Page
                </button>
            </div>
            {process.env.NODE_ENV === 'development' && (
                <div className="mt-8 p-4 bg-red-50 border border-red-200 rounded-md text-left w-full max-w-2xl overflow-auto">
                    <p className="font-mono text-sm text-red-800 whitespace-pre-wrap">{error.message}</p>
                    <p className="font-mono text-xs text-red-600 mt-2 whitespace-pre-wrap">{error.stack}</p>
                </div>
            )}
        </div>
    );
}
