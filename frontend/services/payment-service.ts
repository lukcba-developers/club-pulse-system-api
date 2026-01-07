import api from '../lib/axios';

export interface PaymentInitiation {
    checkout_url: string; // The MercadoPago redirect URL
    payment_id: string;
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
    }
}
