import api from "@/lib/axios";

export interface Club {
    id: string;
    name: string;
    slug: string;
    domain?: string;
    logo_url?: string;
    theme_config?: string; // or parsed object if strict
    settings?: string;
    status: 'ACTIVE' | 'INACTIVE';
    created_at: string;
    updated_at: string;
}

export interface AdPlacement {
    id: string;
    sponsor_id: string;
    location_type: string;
    location_detail: string;
    image_url?: string;
}

export const clubService = {
    getActiveAds: async () => {
        const response = await api.get("/club/ads");
        return response.data.data;
    },

    listClubs: async (): Promise<Club[]> => {
        const response = await api.get("/admin/clubs");
        // Backend returns { data: clubs }
        return response.data.data;
    },

    createClub: async (data: Partial<Club>): Promise<Club> => {
        const response = await api.post("/admin/clubs", data);
        return response.data;
    }
};
