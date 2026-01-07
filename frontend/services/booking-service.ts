import api from '../lib/axios';

export interface Booking {
    id: string;
    user_id: string;
    facility_id: string;
    start_time: string;
    end_time: string;
    status: 'CONFIRMED' | 'CANCELLED';
    guest_details?: GuestDetail[];
    created_at: string;
}

export interface GuestDetail {
    name: string;
    dni: string;
    fee_amount: number;
}

export interface CreateBookingDTO {
    facility_id: string;
    start_time: string;
    end_time: string;
    guest_details?: GuestDetail[];
}

export interface WaitlistEntry {
    id: string;
    resource_id: string;
    target_date: string;
    status: string;
}

export interface AvailabilitySlot {
    start_time: string;
    end_time: string;
    available: boolean;
}

export const bookingService = {
    // List user bookings
    getMyBookings: async () => {
        const response = await api.get<{ data: Booking[] }>('/bookings');
        return response.data.data;
    },

    // Create a new booking
    createBooking: async (data: CreateBookingDTO) => {
        const response = await api.post<Booking>('/bookings', data);
        return response.data;
    },

    // Cancel a booking
    cancelBooking: async (id: string) => {
        const response = await api.delete<{ message: string }>(`/bookings/${id}`);
        return response.data;
    },

    // Check availability for a facility on a specific date (YYYY-MM-DD)
    getAvailability: async (facilityId: string, date: string) => {
        const response = await api.get<{ data: AvailabilitySlot[] }>(`/bookings/availability?facility_id=${facilityId}&date=${date}`);
        return response.data.data;
    },

    // Admin: Get all bookings for the club
    getClubBookings: async (filters?: { facility_id?: string; from?: string; to?: string }) => {
        const params = new URLSearchParams();
        if (filters?.facility_id) params.append('facility_id', filters.facility_id);
        if (filters?.from) params.append('from', filters.from);
        if (filters?.to) params.append('to', filters.to);

        const response = await api.get<{ data: Booking[] }>(`/bookings/all?${params.toString()}`);
        return response.data.data;
    },

    addToWaitlist: async (data: { resource_id: string; target_date: string }) => {
        const response = await api.post<WaitlistEntry>('/bookings/waitlist', data);
        return response.data;
    }
};
