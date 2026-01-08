import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
    test('should login successfully with HttpOnly cookies', async ({ page, context }) => {
        // 1. Navigate to Login (required for page context)
        await page.goto('/login');

        // 0. Setup: Register via Browser Context (to bypass Node network issues)
        await page.evaluate(async () => {
            const CLUB_ID = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';
            try {
                // Using 127.0.0.1 to be safe, though localhost usually works in browser
                const response = await fetch('http://127.0.0.1:8080/api/v1/auth/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-Club-ID': CLUB_ID
                    },
                    body: JSON.stringify({
                        name: 'System Admin',
                        email: 'admin@clubpulse.com',
                        password: 'admin123',
                    })
                });
                console.log('Registration status:', response.status);
            } catch (e) {
                console.log('Registration fetch warning:', e);
            }
        });

        // 2. Fill Credentials
        await page.fill('input[name="email"]', 'admin@clubpulse.com');
        await page.fill('input[name="password"]', 'admin123');

        // 3. Submit
        await page.click('button[type="submit"]');

        // 4. Verify Redirect to Dashboard
        await expect(page).toHaveURL('/');

        // 5. Verify Dashboard Content
        // Assuming "Bienvenido" or user name is shown.
        await expect(page.getByText('Bienvenido, System Admin')).toBeVisible();

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
