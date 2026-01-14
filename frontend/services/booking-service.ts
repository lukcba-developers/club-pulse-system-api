import api from '../lib/axios';

export interface Booking {
    id: string;
    club_id: string;
    user_id: string;
    facility_id: string;
    start_time: string;
    end_time: string;
    total_price: string;
    status: 'PENDING_PAYMENT' | 'CONFIRMED' | 'CANCELLED' | 'EXPIRED' | 'COMPLETED' | 'NO_SHOW';
    guest_details?: GuestDetail[];
    payment_expiry?: string;
    created_at: string;
    updated_at: string;
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

/**
 * Represents an availability slot for a facility.
 * Used by the booking calendar component.
 */
export interface AvailabilitySlot {
    /** Format: "HH:MM" (24h), e.g. "14:00" */
    start_time: string;
    /** Format: "HH:MM" (24h), e.g. "15:00" */
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

    /**
     * Add user to waitlist for a specific resource.
     * @param data - Waitlist entry data
     * @param data.target_date - ISO 8601 format: "YYYY-MM-DDTHH:mm:ssZ"
     */
    addToWaitlist: async (data: { resource_id: string; target_date: string }) => {
        const response = await api.post<WaitlistEntry>('/bookings/waitlist', data);
        return response.data;
    },

    // Admin: Create recurring booking rule
    createRecurringRule: async (data: CreateRecurringRuleDTO) => {
        const response = await api.post<{ data: RecurringRule }>('/bookings/recurring', data);
        return response.data.data;
    },

    // Admin: Generate bookings from rules
    generateFromRules: async (weeks: number = 4) => {
        const response = await api.post<{ message: string }>('/bookings/generate', { weeks });
        return response.data;
    },

    // Admin: List recurring rules (future implementation)
    listRecurringRules: async () => {
        const response = await api.get<{ data: RecurringRule[] }>('/bookings/recurring');
        return response.data.data;
    },

    // Update guest fee (Example usage, implementing as standard practice)
    formatGuestFee: (amount: number): string => {
        return amount.toFixed(2);
    }
};

// --- Recurring Booking Types ---

export type RecurrenceFrequency = 'WEEKLY' | 'MONTHLY';
export type RecurrenceType = 'CLASS' | 'MAINTENANCE' | 'FIXED';

export interface RecurringRule {
    id: string;
    facility_id: string;
    type: RecurrenceType;
    frequency: RecurrenceFrequency;
    day_of_week: number; // 0=Sunday, 6=Saturday
    start_time: string;
    end_time: string;
    start_date: string;
    end_date: string;
    created_at?: string;
}

export interface CreateRecurringRuleDTO {
    facility_id: string;
    type?: RecurrenceType;
    frequency: RecurrenceFrequency;
    day_of_week: number;
    start_time: string;
    end_time: string;
    start_date: string; // YYYY-MM-DD
    end_date: string;   // YYYY-MM-DD
}
