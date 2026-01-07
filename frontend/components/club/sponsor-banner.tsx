"use client";

import { useEffect, useState } from "react";
import { clubService, AdPlacement } from "@/services/club-service";

export function SponsorBanner() {
    const [ads, setAds] = useState<AdPlacement[]>([]);

    useEffect(() => {
        // Fetch ads quietly
        clubService.getActiveAds().then(data => {
            if (data) setAds(data);
        }).catch(err => console.error("Failed to load ads", err));
    }, []);

    if (ads.length === 0) return null;

    // Simple rotation or list. For now, just show the first one or a grid.
    // Let's make a simple consistent banner at the bottom or top.

    return (
        <div className="w-full bg-slate-50 border-t p-4 mt-8">
            <div className="text-xs text-center text-muted-foreground mb-2 uppercase tracking-widest">Nuestros Sponsors</div>
            <div className="flex flex-wrap justify-center gap-8 items-center opacity-80 grayscale hover:grayscale-0 transition-all duration-500">
                {ads.map((ad, idx) => (
                    <div key={ad.id || idx} className="text-sm font-semibold text-slate-600">
                        {/* If we had logos, we would put <img> here. For now using text/detail */}
                        {ad.location_detail || "Sponsor"}
                    </div>
                ))}
            </div>
        </div>
    );
}
