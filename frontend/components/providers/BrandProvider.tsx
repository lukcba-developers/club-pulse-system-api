'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import { useAuthContext } from '@/context/auth-context';
import { clubService, Club } from '@/services/club-service';

interface BrandContextType {
    club: Club | null;
    isLoading: boolean;
}

const BrandContext = createContext<BrandContextType>({ club: null, isLoading: true });

export function useBrand() {
    return useContext(BrandContext);
}

// Helper to calculate contrast color (black or white) based on hex background
function getContrastColor(hexColor: string): string {
    // Remove # if present
    const hex = hexColor.replace('#', '');
    const r = parseInt(hex.substr(0, 2), 16);
    const g = parseInt(hex.substr(2, 2), 16);
    const b = parseInt(hex.substr(4, 2), 16);

    // Calculate Luminance (standard formula)
    const yiq = ((r * 299) + (g * 587) + (b * 114)) / 1000;

    // Returns black for bright colors, white for dark colors
    return (yiq >= 128) ? '#000000' : '#ffffff';
}

export function BrandProvider({ children }: { children: React.ReactNode }) {
    const { user } = useAuthContext();
    const [club, setClub] = useState<Club | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const applyBranding = async () => {
            setIsLoading(true);
            if (user?.club_id) {
                try {
                    const fetchedClub = await clubService.getClub(user.club_id);
                    setClub(fetchedClub);
                    const root = document.documentElement;

                    if (fetchedClub.primary_color) {
                        root.style.setProperty('--brand-600', fetchedClub.primary_color);
                        root.style.setProperty('--brand-500', fetchedClub.primary_color);

                        // Calculate and set foreground color for contrast
                        const contrastColor = getContrastColor(fetchedClub.primary_color);
                        root.style.setProperty('--brand-foreground', contrastColor);
                    }
                    if (fetchedClub.secondary_color) {
                        root.style.setProperty('--brand-900', fetchedClub.secondary_color);
                    }
                } catch (error) {
                    console.error("Failed to load club branding", error);
                }
            } else {
                setClub(null);
                const root = document.documentElement;
                root.style.removeProperty('--brand-600');
                root.style.removeProperty('--brand-500');
                root.style.removeProperty('--brand-900');
                root.style.removeProperty('--brand-foreground');
            }
            setIsLoading(false);
        };

        applyBranding();
    }, [user?.club_id]);

    return (
        <BrandContext.Provider value={{ club, isLoading }}>
            {children}
        </BrandContext.Provider>
    );
}
