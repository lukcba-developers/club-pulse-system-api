import axios from 'axios';

// Create generic axios instance
const api = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true, // Enable sending cookies
});

// Helper to read a cookie by name
function getCookie(name: string): string | undefined {
    if (typeof document === 'undefined') return undefined;
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(';').shift();
    return undefined;
}

// Add a request interceptor to Add Log, Club ID, and CSRF Token
api.interceptors.request.use(
    (config) => {
        // Client-side only
        if (typeof window !== 'undefined') {
            // Inject Club ID (Priority: LocalStorage > Env > Default)
            const clubID = localStorage.getItem('clubID') || process.env.NEXT_PUBLIC_DEFAULT_CLUB_ID || 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';
            config.headers['X-Club-ID'] = clubID;

            // Inject CSRF Token for state-changing methods
            const method = config.method?.toUpperCase();
            if (method === 'POST' || method === 'PUT' || method === 'DELETE' || method === 'PATCH') {
                const csrfToken = getCookie('csrf_token');
                if (csrfToken) {
                    config.headers['X-CSRF-Token'] = csrfToken;
                }
            }
        }

        // Enhanced Logging
        console.groupCollapsed(`[API Request] ${config.method?.toUpperCase()} ${config.url}`);
        console.log('Headers:', config.headers);
        if (config.data) console.log('Payload:', config.data);
        console.groupEnd();

        return config;
    },
    (error) => {
        console.error('[API Request Error]:', error);
        return Promise.reject(error);
    }
);

// Response Interceptor
api.interceptors.response.use(
    (response) => {
        console.groupCollapsed(`[API Response] ${response.status} ${response.config.url}`);
        console.log('Data:', response.data);
        console.groupEnd();
        return response;
    },
    (error: unknown) => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const axiosError = error as any;
        const status = axiosError.response?.status;
        const url = axiosError.config?.url;
        const responseData = axiosError.response?.data;

        // Extract error type from backend response for humanization
        // Backend sends: { type: "booking_conflict", error: "..." }
        if (responseData?.type) {
            axiosError.errorType = responseData.type;
        }

        // Skip logging for 401/403 as they are handled by AuthContext (session expiry)
        if (status === 401 || status === 403) {
            console.groupCollapsed(`[API Auth] ${status} ${url}`);
            console.log('Session expired or invalid token. Redirecting/Clearing session.');
            console.groupEnd();
        } else {
            console.group(`[API Error] ${status || 'Net'} ${url}`);
            console.error('Type:', responseData?.type || 'unknown');
            console.error('Message:', axiosError.message);
            if (axiosError.response) {
                console.error('Response Data:', responseData);
                console.error('Status:', status);
            }
            console.groupEnd();
        }

        return Promise.reject(error);
    }
);


export default api;
