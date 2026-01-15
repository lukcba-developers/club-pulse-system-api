'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import { useAuthContext } from '@/context/auth-context';
import { clubService, Club } from '@/services/club-service';

interface BrandContextType {
    club: Club | null;
}

const BrandContext = createContext<BrandContextType>({ club: null });

export function useBrand() {
    return useContext(BrandContext);
}

export function BrandProvider({ children }: { children: React.ReactNode }) {
    const { user } = useAuthContext();
    const [club, setClub] = useState<Club | null>(null);

    useEffect(() => {
        const applyBranding = async () => {
            if (user?.club_id) {
                try {
                    const fetchedClub = await clubService.getClub(user.club_id);
                    setClub(fetchedClub);
                    const root = document.documentElement;

                    if (fetchedClub.primary_color) {
                        root.style.setProperty('--brand-600', fetchedClub.primary_color);
                        root.style.setProperty('--brand-500', fetchedClub.primary_color);
                        // Also set a reliable default for text if needed, but brand-600 is widely used.
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
            }
        };

        applyBranding();
    }, [user?.club_id]);

    return (
        <BrandContext.Provider value={{ club }}>
            {children}
        </BrandContext.Provider>
    );
}
