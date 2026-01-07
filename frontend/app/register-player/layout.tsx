
import React from "react";

export default function RegisterPlayerLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <div className="min-h-screen bg-slate-50 flex flex-col items-center py-10 px-4">
            <div className="w-full max-w-lg">
                {children}
            </div>
        </div>
    );
}
