import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
    test('should login successfully with HttpOnly cookies', async ({ page, context }) => {
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
