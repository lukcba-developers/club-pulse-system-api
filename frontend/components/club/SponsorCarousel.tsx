"use client"

import { useEffect, useState } from "react"
// import { Ad } from '@/services/club-service'; // Unused
// import { Card, CardContent } from "@/components/ui/card" // Unused

// Defining a type that matches what the backend SHOULD return for display
interface AdDisplay {
    id: string
    sponsor_id: string
    sponsor_name?: string
    sponsor_logo?: string
    location_type: string
}

export function SponsorCarousel({ clubSlug }: { clubSlug: string }) {
    const [ads, setAds] = useState<AdDisplay[]>([])

    useEffect(() => {
        // Fetch ads
        // Backend endpoint: /public/clubs/:slug/ads
        const fetchAds = async () => {
            try {
                const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/public/clubs/${clubSlug}/ads`)
                if (res.ok) {
                    const data = await res.json()
                    // Mapear data si es necesario.
                    // PROBLEMA: El backend actual retorna AdPlacement struct. No tiene Logo ni Nombre de Sponsor.
                    // Necesito arreglar el backend para hacer Preload("Sponsor").
                    setAds(data.data || [])
                }
            } catch (e) {
                console.error(e)
            }
        }
        fetchAds()
    }, [clubSlug])

    if (!ads.length) return null

    return (
        <div className="w-full bg-slate-50 py-8 border-t">
            <div className="container mx-auto px-4">
                <h3 className="text-center text-sm font-semibold text-muted-foreground uppercase tracking-widest mb-6">Nuestros Sponsors</h3>
                <div className="flex overflow-x-auto gap-8 justify-center items-center pb-4 scrollbar-hide">
                    {ads.map((ad) => (
                        <div key={ad.id} className="flex-shrink-0 grayscale hover:grayscale-0 transition-all duration-300">
                            {/* Placeholder visual hasta que arregle backend */}
                            <div className="h-12 w-32 bg-slate-200 rounded flex items-center justify-center text-xs text-slate-400">
                                {ad.sponsor_name || "Sponsor"}
                            </div>
                            {/* Cuando backend tenga logo:
                         <img src={ad.sponsor_logo} alt="Sponsor" className="h-12 w-auto object-contain" />
                         */}
                        </div>
                    ))}
                    {/* Mock Sponsors si no hay data real para demo */}
                    {ads.length === 0 && Array.from({ length: 5 }).map((_, i) => (
                        <div key={i} className="h-12 w-32 bg-slate-200 rounded animate-pulse" />
                    ))}
                </div>
            </div>
        </div>
    )
}
