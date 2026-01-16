export type BookingStatus = 'PENDING_PAYMENT' | 'CONFIRMED' | 'CANCELLED' | 'EXPIRED' | 'COMPLETED' | 'NO_SHOW';

export interface GuestDetail {
    name: string;
    dni: string;
    fee_amount: number;
}

export interface Booking {
    id: string;
    club_id: string;
    user_id: string;
    facility_id: string;
    start_time: string; // ISO 8601
    end_time: string;   // ISO 8601
    total_price: string; // Decimal string from backend
    status: BookingStatus;
    guest_details?: GuestDetail[];
    payment_expiry?: string;
    created_at: string;
    updated_at: string;
}
