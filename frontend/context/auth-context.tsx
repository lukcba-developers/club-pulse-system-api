'use client';

import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import api from '@/lib/axios';

export interface User {
    id: string;
    club_id?: string;
    name: string;
    email: string;
    role: string;
    medical_cert_status?: 'VALID' | 'EXPIRED' | 'PENDING';
    medical_cert_expiry?: string;
    family_group_id?: string;
    emergency_contact_name?: string;
    emergency_contact_phone?: string;
    insurance_provider?: string;
    insurance_number?: string;
}

interface AuthContextType {
    user: User | null;
    loading: boolean;
    login: () => Promise<void>;
    logout: () => void;
    checkAuth: () => Promise<void>;
    refreshUser: () => Promise<void>; // Alias for checkAuth
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    const checkAuth = useCallback(async () => {
        try {
            // Attempt to fetch profile. If cookie is present and valid, this succeeds.
            const response = await api.get('/users/me');
            setUser(response.data);
        } catch (error: unknown) {
            // Silently handle 401/403 (Invalid Token)
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            if ((error as any).response && ((error as any).response.status === 401 || (error as any).response.status === 403 || (error as any).response.status === 404)) {
                setUser(null);
            } else {
                console.error('Failed to fetch user', error);
            }
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        checkAuth();
    }, [checkAuth]);

    const login = useCallback(async () => {
        // Just refresh the auth state, cookie should be set by the login API call
        await checkAuth();
    }, [checkAuth]);

    const logout = useCallback(async () => {
        try {
            // Call logout endpoint to clear cookie
            // We need to pass something? Handler expects Refresh Token in body often. 
            // In cookie-based auth, we should just need to hit /logout (POST).
            // But our current handler expects `refresh_token` in body.
            // We need to update backend Logout as well to be cookie aware or just ignore body if cookie.
            // For now, let's just clear Client state. Backend change for full logout is needed too.
            // We'll call API anyway.
            await api.post('/auth/logout', { refresh_token: "cookie" }); // Dummy for now
        } catch (e) {
            console.error("Logout failed", e);
        } finally {
            setUser(null);
            // Optionally force reload to clear all states
            // window.location.href = '/login';
        }
    }, []);

    const value = React.useMemo(() => ({
        user, loading, login, logout, checkAuth, refreshUser: checkAuth
    }), [user, loading, login, logout, checkAuth]);

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuthContext() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuthContext must be used within an AuthProvider');
    }
    return context;
}
