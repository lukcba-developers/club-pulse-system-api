import { test, expect } from '@playwright/test';

test.describe('Booking Flow', () => {
    test('should open booking modal and select a slot', async ({ page }) => {
        // 1. Mock Availability API
        await page.route('**/bookings/availability*', async route => {
            const json = {
                data: [
                    { start_time: '10:00', end_time: '11:00', status: 'available' },
                    { start_time: '11:00', end_time: '12:00', status: 'booked' }
                ]
            };
            await route.fulfill({ status: 200, json });
        });

        // 2. Mock Facilities API (CRITICAL: Must be exact path match or glob if query params)
        await page.route('**/facilities*', async route => {
            const json = {
                data: [
                    {
                        id: 'cam-1',
                        name: 'Cancha 1',
                        type: 'padel',
                        price_per_hour: 2000,
                        status: 'active',
                        capacity: 4,
                        location: { name: 'Sede Central' }
                    }
                ]
            };
            await route.fulfill({ status: 200, json });
        });

        // 3. Mock Login
        await page.context().addCookies([{
            name: 'access_token',
            value: 'dummy-token',
            domain: 'localhost',
            path: '/',
            httpOnly: true
        }]);

        // Mock User Profile (Auth Context needs this)
        await page.route('**/users/me', async route => {
            const json = {
                id: 'user-1',
                name: 'Test User',
                email: 'test@example.com',
                role: 'admin',
                club_id: 'club-1'
            };
            await route.fulfill({ status: 200, json });
        });

        // 4. Navigate to Dashboard
        await page.goto('/');

        // 5. Verify Facility Loaded
        await expect(page.getByText('Cancha 1')).toBeVisible({ timeout: 10000 });

        // Click "Reservar" - Ensure we click the one that opens modal directly or via card footer
        // The card has TWO buttons that open modal.
        // We target the explicit one.
        await page.getByRole('button', { name: 'Reservar' }).first().click();

        // 6. Verify Modal Open
        const modal = page.locator('div[role="dialog"]');
        await expect(modal).toBeVisible();
        await expect(modal).toContainText('Reservar Cancha 1');

        // 7. Verify Slots
        // 10:00 should be enabled (available)
        const availableSlot = modal.getByRole('button', { name: '10:00' });
        await expect(availableSlot).toBeEnabled();

        // 11:00 should be disabled (booked)
        const bookedSlot = modal.getByRole('button', { name: '11:00' });
        await expect(bookedSlot).toBeDisabled();

        // 8. Select Slot
        await availableSlot.click();
        await expect(modal).toContainText('Seleccionado:');

        // 9. Submit
        await page.route('**/bookings', async route => {
            await route.fulfill({ status: 201, json: { message: "Created" } });
        });

        await modal.getByRole('button', { name: 'Confirmar Reserva' }).click();

        // 10. Verify Success Message
        await expect(modal).toContainText('¡Reserva Exitosa!');
    });
});
