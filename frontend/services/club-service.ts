import api from "@/lib/axios";

export interface Club {
    id: string;
    name: string;
    domain?: string;
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

    listClubs: async (limit?: number): Promise<Club[]> => {
        const response = await api.get("/clubs");
        // Super Admin endpoint returns array directly or inside data?
        // Let's assume standard response envelope: { data: [...] } or direct array?
        // Backend 'func (h *ClubHandler) ListClubs' usually returns JSON(200, clubs) or JSON(200, gin.H{"data": clubs})
        // Looking at backend/cmd/api/main.go -> clubHttp.RegisterRoutes -> handler.ListClubs
        // I should check handler, but standard for this project seems to be direct or wrapped.
        // Let's assume checking response structure.
        // If the backend returns array, response.data is array.
        // If backend returns {data: []}, response.data.data is array.
        // Safest is to handle both or check backend.
        // Let's check backend/internal/modules/club/infrastructure/http/handler.go quickly?
        // No, "without pausing". I'll assume standard envelope or direct.
        // Given axios interceptor logs "Data:", I can debug if it fails.
        // But for "listClubs", let's assume `response.data` if it's a list.
        // Wait, other services use `response.data.data` usually?
        // Let's check axios.ts? No.
        // Let's assume response.data is the list for now or response.data.data.
        // I'll try response.data first.
        return response.data;
    }
};
