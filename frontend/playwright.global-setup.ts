import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

/**
 * Global setup for Playwright E2E tests
 * Ensures backend services are running before tests start
 */
async function globalSetup() {
    console.log('🔧 [E2E Setup] Checking backend services...');

    try {
        // Check if backend is already running on port 8081
        const isBackendRunning = await checkPort(8081);

        if (isBackendRunning) {
            console.log('✅ [E2E Setup] Backend already running on port 8081');
            return;
        }

        console.log('🚀 [E2E Setup] Starting backend services via docker-compose...');

        // Start postgres, redis, and api services
        await execAsync('docker-compose up -d postgres redis api', {
            cwd: process.cwd().replace('/frontend', ''),
        });

        // Wait for services to be healthy
        console.log('⏳ [E2E Setup] Waiting for services to be healthy...');
        await waitForHealthy('http://localhost:8081/health', 30000);

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
        const { stdout } = await execAsync(`lsof -ti :${port}`);
        return stdout.trim().length > 0;
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
            const response = await fetch(url);
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
