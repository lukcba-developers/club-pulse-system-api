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
        // This makes the test readable: "Go there, book Court 1 at 10:00, expect success"
        await bookingPage.bookCourt('Cancha 1', '10:00');

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

        // 11:00 should be disabled (booked)
        await bookingPage.expectSlotBooked('11:00');
    });
});
