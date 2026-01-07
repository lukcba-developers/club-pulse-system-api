import api from '../lib/axios';

export interface Facility {
    id: string;
    name: string;
    type: string;
    status: string;
    capacity: number;
    hourly_rate: number;
    specifications: {
        surface_type?: string;
        lighting: boolean;
        covered: boolean;
        equipment?: string[];
    };
    location: {
        name: string;
        description?: string;
    };
    created_at: string;
    updated_at: string;
}

export interface CreateFacilityRequest {
    name: string;
    type: string;
    description?: string;
    hourly_rate: number;
    capacity: number;
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
            type: data.type, // "Tennis Court", "Swimming Pool", etc.
            description: data.description,
            hourly_rate: data.hourly_rate,
            capacity: data.capacity,
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
    }
};
