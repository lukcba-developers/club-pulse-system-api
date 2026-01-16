import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';

interface CustomAxiosError extends AxiosError {
    errorType?: string;
    config: InternalAxiosRequestConfig & { _retry?: boolean };
}
import { v4 as uuidv4 } from 'uuid';

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
            // Inject Auth Header from LocalStorage (if available)
            const token = localStorage.getItem('token');
            if (token) {
                config.headers['Authorization'] = `Bearer ${token}`;
            }
            // Inject Club ID (Priority: LocalStorage > Env > Default)
            const clubID = localStorage.getItem('clubID') || process.env.NEXT_PUBLIC_DEFAULT_CLUB_ID || 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';
            config.headers['X-Club-ID'] = clubID;

            // Inject Trace ID / Correlation ID
            const requestID = uuidv4();
            config.headers['X-Request-ID'] = requestID;
            // Also useful for simple logs
            config.headers['X-Correlation-ID'] = requestID;

            // W3C Trace Context - traceparent
            // Version (2 hex) - TraceID (32 hex) - SpanID (16 hex) - TraceFlags (2 hex)
            const traceId = requestID.replace(/-/g, ''); // UUID to 32 hex
            const spanId = uuidv4().replace(/-/g, '').substring(0, 16); // Random 16 hex
            const traceFlags = '01'; // Sampled
            const traceparent = `00-${traceId}-${spanId}-${traceFlags}`;
            config.headers['traceparent'] = traceparent;

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
    // Success handler - just pass through successful responses
    (response) => response,
    // Error handler
    (error: unknown) => {
        const axiosError = error as CustomAxiosError;
        const status = axiosError.response?.status;
        const url = axiosError.config?.url;
        const responseData = axiosError.response?.data as { type?: string } | undefined;

        // Extract error type from backend response for humanization
        if (responseData?.type) {
            axiosError.errorType = responseData.type;
        }

        // Automatic Token Refresh Logic
        if (status === 401 && responseData?.type === 'TOKEN_EXPIRED' && !axiosError.config._retry) {
            axiosError.config._retry = true;
            console.log('[API Auth] Token expired. Attempting refresh...');

            return api.post('/auth/refresh', { refresh_token: 'cookie' })
                .then((res) => {
                    if (res.status === 200) {
                        console.log('[API Auth] Refresh successful. Retrying original request.');
                        return api(axiosError.config);
                    }
                    return Promise.reject(error);
                })
                .catch((refreshError) => {
                    console.error('[API Auth] Refresh failed. Redirecting to login.');
                    // Optional: Clear storage/redirect if needed
                    // window.location.href = '/login';
                    return Promise.reject(refreshError);
                });
        }

        // Skip logging for 401/403 as they are handled by AuthContext or refresh logic
        if (status === 401 || status === 403) {
            console.groupCollapsed(`[API Auth] ${status} ${url}`);
            console.log('Session expired or invalid token.');
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
