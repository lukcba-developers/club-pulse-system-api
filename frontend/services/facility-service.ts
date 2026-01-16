import api from '../lib/axios';

import { Facility, FacilityType, FacilityStatus } from '@/types/facility';

export { FacilityType, FacilityStatus };
export type { Facility };

export interface CreateFacilityRequest {
    name: string;
    description?: string;
    type: FacilityType;
    hourly_rate: number;
    capacity: number;
    opening_time?: string; // HH:MM
    closing_time?: string; // HH:MM
    guest_fee?: number;
    location_name: string;
    location_description?: string;
    surface_type?: string;
    lighting?: boolean;
    covered?: boolean;
}

export interface SearchResult {
    facility: Facility;
    similarity: number;
}

export interface SearchResponse {
    query: string;
    count: number;
    results: SearchResult[];
}

export interface UpdateFacilityRequest {
    name?: string;
    description?: string;
    status?: FacilityStatus;
    opening_time?: string; // HH:MM
    closing_time?: string; // HH:MM
    hourly_rate?: number;
    guest_fee?: number;
    specifications?: {
        surface_type?: string;
        lighting?: boolean;
        covered?: boolean;
    };
}

export const facilityService = {
    /**
     * List all facilities with pagination
     */
    list: async (limit = 10, offset = 0): Promise<Facility[]> => {
        const response = await api.get<Facility[]>('/facilities', {
            params: { limit, offset }
        });
        return response.data;
    },

    /**
     * Get a single facility by ID
     */
    getById: async (id: string): Promise<Facility> => {
        const response = await api.get<Facility>(`/facilities/${id}`);
        return response.data;
    },

    /**
     * Create a new facility
     */
    create: async (data: CreateFacilityRequest): Promise<Facility> => {
        const response = await api.post<{ data: Facility }>('/facilities', {
            name: data.name,
            description: data.description,
            type: data.type,
            hourly_rate: data.hourly_rate,
            capacity: data.capacity,
            opening_time: data.opening_time ?? "08:00",
            closing_time: data.closing_time ?? "22:00",
            guest_fee: data.guest_fee ?? 0,
            location: {
                name: data.location_name,
                description: data.location_description
            },
            specifications: {
                surface_type: data.surface_type,
                lighting: data.lighting || false,
                covered: data.covered || false
            }
        });
        return response.data.data;
    },

    /**
     * Semantic search for facilities using natural language
     * Examples: "canchas techadas", "piscina nocturna", "tenis con luz"
     */
    search: async (query: string, limit = 10): Promise<SearchResponse> => {
        const response = await api.get<SearchResponse>('/facilities/search', {
            params: { q: query, limit }
        });
        return response.data;
    },

    /**
     * Generate embeddings for all facilities (admin operation)
     */
    generateEmbeddings: async (): Promise<{ message: string; processed: number }> => {
        const response = await api.post<{ message: string; processed: number }>('/facilities/embeddings/generate');
        return response.data;
    },

    /**
     * Update facility details including schedule (opening/closing hours)
     */
    update: async (id: string, data: UpdateFacilityRequest): Promise<Facility> => {
        const response = await api.put<Facility>(`/facilities/${id}`, data);
        return response.data;
    },

    /**
     * Update facility schedule (convenience method)
     */
    updateSchedule: async (id: string, openingTime: string, closingTime: string): Promise<Facility> => {
        return facilityService.update(id, {
            opening_time: openingTime,
            closing_time: closingTime
        });
    }
};
