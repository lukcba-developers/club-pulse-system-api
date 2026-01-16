import api from '../lib/axios';
import { v4 as uuidv4 } from 'uuid';

export interface AccessResult {
    status: 'GRANTED' | 'DENIED';
    reason?: string;
    timestamp?: string;
}

export const accessService = {

    simulateEntry: async (userId: string, deviceId = "WEB_DASHBOARD") => {
        const response = await api.post<AccessResult>('/access/entry', {
            user_id: userId,
            direction: 'IN',
            device_id: deviceId,
            event_id: uuidv4(), // Generate unique event ID for idempotency
            timestamp: new Date().toISOString() // Send client timestamp for offline sync accuracy
        });
        return response.data;
    },

    /**
     * Generate a QR code for virtual access (stub)
     * In a real implementation, this would fetch a signed token from the backend
     */
    generateQR: async (): Promise<string> => {
        // Return a dummy value or fetch from backend
        // await api.get('/access/credentials/virtual');
        return "mock-qr-token-" + Date.now();
    }
}
