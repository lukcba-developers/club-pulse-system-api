import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
    test('should login successfully with HttpOnly cookies', async ({ page, context }) => {
        // 0. Setup: Ensure user exists using default ClubID from axios.ts
        const CLUB_ID = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';
        try {
            await page.request.post('http://127.0.0.1:8080/api/v1/auth/register', {
                data: {
                    name: 'System Admin',
                    email: 'admin@clubpulse.com',
                    password: 'admin123',
                },
                headers: {
                    'X-Club-ID': CLUB_ID
                }
            });
            console.log('Test user registered or already exists');
        } catch (e) {
            console.log('Registration warning (might already exist):', e);
        }

        // 1. Navigate to Login
        await page.goto('/login');

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
