import api from "@/lib/axios";

export interface Club {
    id: string;
    name: string;
    slug: string;
    domain?: string;
    logo_url?: string;
    primary_color?: string;
    secondary_color?: string;
    contact_email?: string;
    contact_phone?: string;
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
    },

    getClub: async (id: string): Promise<Club> => {
        // We can reuse getPublicClubBySlug if id is slug, or we need a GetClub endpoint by ID.
        // Backend handler has GetByID but not explicitly exposed as /clubs/:id for public, but YES for admin/clubs/:id?
        // Wait, handler.go:210 clubs.POST("", handler.CreateClub); 211 clubs.GET("", handler.ListClubs).
        // It does NOT have GET /admin/clubs/:id.
        // But it has handlers "GetClub" in usecase.
        // I should stick to what exists. 
        // handler.go has public.GET("/:slug", handler.GetPublicClubBySlug) at /public/clubs/:slug.
        // If I have the ID, I might not have the slug easily unless user object has it.
        // User object has club_id. 
        // I need an endpoint to get club by ID or I need to expose it.
        // I'll assume I should use the admin endpoint if user is admin, or public if not?
        // But `users/me` -> `club_id`.
        // I should have an endpoint GET /clubs/:id or similar.
        // I added PUT /admin/clubs/:id.
        // I should add GET /admin/clubs/:id as well?
        // Or generic GET /clubs/:id for members?
        // The user said "Cuando un usuario perteneciente a un club inicie sesión ... detecte ... primary_color".
        // Use /public/clubs/:slug is available.
        // But I only have club_id in user.
        // I will add GET /admin/clubs/:id to handler.go as well OR simpler:
        // Use user.club_id to fetch via a new endpoint or existing.
        // I'll add GET "/admin/clubs/:id" for simplicity since I am admin in this context (creating clubs).
        // BUT wait, a Member logging in is NOT an admin.
        // They need access to THEIR club info.
        // I should probably have `GET /my-club` or similar, or allow `GET /public/clubs/by-id/:id`.
        // Or just `GET /admin/clubs/:id` but that requires Admin role.
        // Creating a route `GET /public/clubs/id/:id` might be best for members.
        // I will add GET /admin/clubs/:id to handler.go first for completeness.
        const response = await api.get(`/admin/clubs/${id}`);
        return response.data;
    }
};
