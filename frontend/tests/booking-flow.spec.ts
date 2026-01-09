import { test, expect } from '@playwright/test';
import { BookingPage } from './pages/BookingPage';

test.describe('Booking Flow', () => {
    // Freeze clock at 8:00 AM to ensure all time slots are in the future
    test.beforeEach(async ({ page }) => {
        // Set the browser clock to 8:00 AM today
        const today = new Date();
        today.setHours(8, 0, 0, 0);
        await page.clock.install({ time: today });
    });

    test('should open booking modal and select a slot', async ({ page }) => {
        page.on('console', msg => console.log(`[Browser]: ${msg.text()}`));
        const bookingPage = new BookingPage(page);

        // 1. Setup Mocks
        await bookingPage.mockApis();

        // 2. Navigate
        await bookingPage.goto();

        // 3. Perform Booking - 10:00 is in the future since clock is frozen at 8:00 AM
        await bookingPage.bookCourt('Cancha 1', '10:00');

        // 4. Verify
        await bookingPage.expectSuccessMessage();
    });

    test('should see booked slots as disabled', async ({ page }) => {
        const bookingPage = new BookingPage(page);
        await bookingPage.mockApis();
        await bookingPage.goto();

        await expect(page.getByText('Cancha 1')).toBeVisible();
        await page.getByRole('button', { name: 'Reservar' }).first().click();

        // 11:00 should be disabled (booked)
        await bookingPage.expectSlotBooked('11:00');
    });
});
