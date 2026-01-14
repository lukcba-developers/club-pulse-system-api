import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

// In CI, the backend runs on port 8080 (started by the workflow)
// Locally, it runs on port 8081 via docker-compose
const BACKEND_PORT = process.env.CI ? 8080 : 8081;
const HEALTH_URL = `http://localhost:${BACKEND_PORT}/health`;

/**
 * Global setup for Playwright E2E tests
 * Ensures backend services are running before tests start
 */
async function globalSetup() {
    console.log(`🔧 [E2E Setup] Checking backend services on port ${BACKEND_PORT}...`);

    try {
        // Check if backend is already running
        const isBackendRunning = await checkPort(BACKEND_PORT);

        if (isBackendRunning) {
            console.log(`✅ [E2E Setup] Backend already running on port ${BACKEND_PORT}`);
            // Still wait for health check to pass
            await waitForHealthy(HEALTH_URL, 30000);
            return;
        }

        // In CI, the backend should already be started by the workflow
        // If it's not running, just wait for it (the workflow starts it in background)
        if (process.env.CI) {
            console.log('⏳ [E2E Setup] CI environment detected. Waiting for backend to be ready...');
            await waitForHealthy(HEALTH_URL, 60000); // Longer timeout for CI
            console.log('✅ [E2E Setup] Backend ready in CI!');
            return;
        }

        // Local development: start services via docker-compose
        console.log('🚀 [E2E Setup] Starting backend services via docker-compose...');

        // Try docker compose (v2 syntax) first, fallback to docker-compose (v1)
        try {
            await execAsync('docker compose up -d postgres redis api', {
                cwd: process.cwd().replace('/frontend', ''),
            });
        } catch {
            await execAsync('docker-compose up -d postgres redis api', {
                cwd: process.cwd().replace('/frontend', ''),
            });
        }

        // Wait for services to be healthy
        console.log('⏳ [E2E Setup] Waiting for services to be healthy...');
        await waitForHealthy(HEALTH_URL, 30000);

        console.log('✅ [E2E Setup] All services ready!');
    } catch (error) {
        console.error('❌ [E2E Setup] Failed to start backend services:', error);
        throw error;
    }
}

/**
 * Check if a port is in use
 */
async function checkPort(port: number): Promise<boolean> {
    try {
        // Use fetch to check if the port is responding
        const response = await fetch(`http://localhost:${port}/health`, {
            signal: AbortSignal.timeout(2000),
        });
        return response.ok;
    } catch {
        return false;
    }
}

/**
 * Wait for a health endpoint to respond with 200
 */
async function waitForHealthy(url: string, timeout: number): Promise<void> {
    const startTime = Date.now();

    while (Date.now() - startTime < timeout) {
        try {
            const response = await fetch(url, {
                signal: AbortSignal.timeout(5000),
            });
            if (response.ok) {
                return;
            }
        } catch {
            // Endpoint not ready yet, continue waiting
        }

        await new Promise(resolve => setTimeout(resolve, 1000));
    }

    throw new Error(`Health check timeout for ${url}`);
}

export default globalSetup;
