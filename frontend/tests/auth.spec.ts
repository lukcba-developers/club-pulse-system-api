import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
    // Ensure clean state by overriding global storageState (which might have dummy auth)
    test.use({ storageState: { cookies: [], origins: [] } });

    test('should login successfully with HttpOnly cookies', async ({ page, context }) => {
        // Enable console logging for debugging
        page.on('console', msg => {
            if (msg.type() === 'error') {
                console.log(`[Browser Error]: ${msg.text()}`);
            }
        });

        // 1. Navigate to Login and ensure clean state
        await page.goto('/login');
        await page.evaluate(() => {
            localStorage.clear();
            sessionStorage.clear();
        });
        await context.clearCookies();
        await page.reload();

        // Wait for login form to be ready
        await page.waitForSelector('input[name="email"]', { timeout: 10000 });

        // 2. Setup: Register via Browser Context (to bypass Node network issues)
        const apiUrl = process.env.TEST_API_URL || 'http://localhost:8081/api/v1';
        const uniqueEmail = `admin-${Date.now()}@clubpulse.com`;

        const registrationResult = await page.evaluate(async ({ url, email }) => {
            const CLUB_ID = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'; // Deterministic ID from Seeder

            // Force ClubID in localStorage to ensure axios uses it
            localStorage.setItem('clubID', CLUB_ID);

            try {
                const response = await fetch(`${url}/auth/register`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-Club-ID': CLUB_ID
                    },
                    body: JSON.stringify({
                        name: 'System Admin',
                        email: email,
                        password: 'Admin123',
                        accept_terms: true,
                        privacy_policy_version: '2026-01'
                    })
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    return { success: false, status: response.status, error: errorText };
                }
                return { success: true, status: response.status };
            } catch (e) {
                return { success: false, error: String(e) };
            }
        }, { url: apiUrl, email: uniqueEmail });

        // Log registration result for debugging
        console.log('Registration result:', registrationResult);
        expect(registrationResult.success).toBe(true);

        // 3. Fill Credentials
        await page.fill('input[name="email"]', uniqueEmail);
        await page.fill('input[name="password"]', 'Admin123');

        // 4. Submit
        await page.click('button[type="submit"]');

        // 5. Wait for navigation to complete (either success or failure)
        await page.waitForURL(url => !url.pathname.includes('/login'), { timeout: 15000 });

        // 6. Verify we're NOT on login page anymore
        const currentUrl = page.url();
        expect(currentUrl).not.toContain('/login');

        // 7. Verify HttpOnly Cookie presence - This is the main goal of the test
        const cookies = await context.cookies();
        const accessTokenCookie = cookies.find(c => c.name === 'access_token');

        expect(accessTokenCookie).toBeDefined();
        expect(accessTokenCookie?.httpOnly).toBe(true);
        // In Dev mode, Secure is typically false
        expect(accessTokenCookie?.secure).toBe(false);

        // 8. Verify LocalStorage is Clean (Security Requirement - no tokens in localStorage)
        const token = await page.evaluate(() => localStorage.getItem('token'));
        expect(token).toBeNull();
    });
});
