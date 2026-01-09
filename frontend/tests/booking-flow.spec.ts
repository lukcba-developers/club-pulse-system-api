import { test, expect } from '@playwright/test';
import { BookingPage } from './pages/BookingPage';

test.describe('Booking Flow', () => {
    test('should open booking modal and select a slot', async ({ page }) => {
        page.on('console', msg => console.log(`[Browser]: ${msg.text()}`));
        const bookingPage = new BookingPage(page);

        // 1. Setup Mocks
        await bookingPage.mockApis();

        // 2. Navigate
        await bookingPage.goto();

        // 3. Perform Booking
        // Use tomorrow and afternoon slot to avoid past date/time issues in CI
        await bookingPage.bookCourtForTomorrow('Cancha 1', '15:00');

        // 4. Verify
        await bookingPage.expectSuccessMessage();
    });

    test('should see booked slots as disabled', async ({ page }) => {
        const bookingPage = new BookingPage(page);
        await bookingPage.mockApis();
        await bookingPage.goto();

        // Open modal manually or reuse a helper that just opens it without booking?
        // Let's reuse the flow but stop before booking, or add a specific method.
        // For now, let's just inspect the modal state after opening it.

        await expect(page.getByText('Cancha 1')).toBeVisible();
        await page.getByRole('button', { name: 'Reservar' }).first().click();

        // 16:00 should be disabled (booked)
        await bookingPage.expectSlotBooked('16:00');
    });
});
