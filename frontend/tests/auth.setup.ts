import { test as setup } from '@playwright/test';

const authFile = 'playwright/.auth/user.json';

setup('authenticate', async ({ page }) => {
    // In our system, "login" is just adding the dummy access_token cookie
    // because we are mocking the backend APIs.
    await page.context().addCookies([{
        name: 'access_token',
        value: 'dummy-token',
        domain: 'localhost',
        path: '/',
        httpOnly: true,
        sameSite: 'Lax'
    }]);

    // Save storage state to be used by other tests
    await page.context().storageState({ path: authFile });
});
