import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
    test('should login successfully with HttpOnly cookies', async ({ page, context }) => {
        // Debug: Log browser console messages to debug CI failures


        // 1. Navigate to Login (required for page context)
        await page.goto('/login');

        // 0. Setup: Register via Browser Context (to bypass Node network issues)
        // Use apiUrl passed from Node context to ensure it matches environment config (and avoids CSP issues)
        const apiUrl = process.env.TEST_API_URL || 'http://localhost:8080/api/v1';
        const uniqueEmail = `admin-${Date.now()}@clubpulse.com`;

        await page.evaluate(async ({ url, email }) => {
            const CLUB_ID = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';

            // Force ClubID in localStorage to ensure axios uses it matches our registration
            localStorage.setItem('clubID', CLUB_ID);

            const response = await fetch(`${url}/auth/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Club-ID': CLUB_ID
                },
                body: JSON.stringify({
                    name: 'System Admin',
                    email: email,
                    password: 'admin123',
                    accept_terms: true,
                    privacy_policy_version: '2026-01'
                })
            });

            if (!response.ok) {
                throw new Error(`Registration failed with status ${response.status}`);
            }
        }, { url: apiUrl, email: uniqueEmail });

        // 2. Fill Credentials
        await page.fill('input[name="email"]', uniqueEmail);
        await page.fill('input[name="password"]', 'admin123');

        // 3. Submit
        await page.click('button[type="submit"]');

        // 4. Verify Redirect to Dashboard
        await expect(page).toHaveURL('/');

        // 5. Verify Dashboard Content
        await expect(page.getByText('Hola, System Admin')).toBeVisible();

        // 6. Verify LocalStorage is Clean (Security Requirement)
        const token = await page.evaluate(() => localStorage.getItem('token'));
        expect(token).toBeNull();

        // 7. Verify HttpOnly Cookie presence
        const cookies = await context.cookies();
        const accessTokenCookie = cookies.find(c => c.name === 'access_token');

        expect(accessTokenCookie).toBeDefined();
        expect(accessTokenCookie?.httpOnly).toBe(true);
        expect(accessTokenCookie?.secure).toBe(false); // In Dev, Secure is false. Adjust if Env changes.
    });
});
