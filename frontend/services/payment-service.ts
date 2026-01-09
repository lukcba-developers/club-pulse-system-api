import api from '../lib/axios';

export interface PaymentInitiation {
    checkout_url: string; // The MercadoPago redirect URL
    payment_id: string;
}

export interface Payment {
    id: string;
    amount: string;
    currency: string;
    status: 'PENDING' | 'COMPLETED' | 'FAILED' | 'REFUNDED';
    method: 'CASH' | 'MERCADOPAGO' | 'STRIPE' | 'TRANSFER' | 'LABOR_EXCHANGE';
    payer_id: string;
    reference_id?: string;
    reference_type?: string;
    notes?: string;
    paid_at?: string;
    created_at: string;
}

export interface PaymentListResponse {
    data: Payment[];
    total: number;
}

export interface OfflinePaymentRequest {
    amount: number;
    method: 'CASH' | 'TRANSFER' | 'LABOR_EXCHANGE';
    payer_id: string;
    reference_id?: string;
    reference_type?: string;
    notes?: string;
}

export const paymentService = {
    // Initiate payment via Backend -> MercadoPago
    initiatePayment: async (amount: number, currency: string, description: string = "Club Payment") => {
        const response = await api.post<PaymentInitiation>('/payments/checkout', {
            amount,
            currency,
            description
        });
        return response.data;
    },

    // Admin: List all payments
    getPayments: async (filters?: { payer_id?: string; status?: string }) => {
        const params = new URLSearchParams();
        if (filters?.payer_id) params.append('payer_id', filters.payer_id);
        if (filters?.status) params.append('status', filters.status);

        const response = await api.get<PaymentListResponse>(`/payments?${params.toString()}`);
        return response.data;
    },

    // Admin: Create offline payment
    createOfflinePayment: async (data: OfflinePaymentRequest) => {
        const response = await api.post<{ data: Payment }>('/payments/offline', data);
        return response.data.data;
    },

    // Admin: Process refund for a completed payment
    refundPayment: async (paymentId: string) => {
        const response = await api.post<{ message: string }>(`/payments/${paymentId}/refund`);
        return response.data;
    }
};

