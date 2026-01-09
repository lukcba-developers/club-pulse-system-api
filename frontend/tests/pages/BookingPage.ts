import { type Page, type Locator, expect } from '@playwright/test';

export class BookingPage {
    readonly page: Page;
    readonly reserveButton: Locator;
    readonly modal: Locator;

    constructor(page: Page) {
        this.page = page;
        this.reserveButton = page.getByRole('button', { name: 'Reservar' });
        this.modal = page.locator('div[role="dialog"]');
    }

    async goto() {
        await this.page.goto('/');
    }

    /**
     * Mocks the necessary APIs for the booking flow
     */
    async mockApis() {
        // 1. Mock Availability API
        await this.page.route('**/bookings/availability*', async route => {
            const json = {
                data: [
                    { start_time: '10:00', end_time: '11:00', status: 'available' },
                    { start_time: '11:00', end_time: '12:00', status: 'booked' }
                ]
            };
            await route.fulfill({ status: 200, json });
        });

        // 2. Mock Facilities API
        await this.page.route('**/facilities*', async route => {
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

        // 3. Mock User Profile
        await this.page.route('**/users/me', async route => {
            const json = {
                id: 'user-1',
                name: 'Test User',
                email: 'test@example.com',
                role: 'ADMIN',
                club_id: 'club-1',
                medical_cert_status: 'VALID'
            };
            await route.fulfill({ status: 200, json });
        });

        // 4. Mock Booking Creation
        await this.page.route('**/bookings', async route => {
            await route.fulfill({ status: 201, json: { message: "Created" } });
        });
    }

    async login() {
        await this.page.context().addCookies([{
            name: 'access_token',
            value: 'dummy-token',
            domain: 'localhost',
            path: '/',
            httpOnly: true
        }]);
    }

    async bookCourt(courtName: string, timeSlot: string) {
        // Wait for the specific court card to be visible
        await expect(this.page.getByText(courtName)).toBeVisible({ timeout: 10000 });

        // Click "Reservar". 
        // Improvement: Scope to the specific card if possible, for now keeping it simple but explicit
        // In a real app we'd chain: page.locator('.court-card', { hasText: courtName }).getByRole(...)
        await this.page.getByRole('button', { name: 'Reservar' }).first().click();

        await expect(this.modal).toBeVisible();
        await expect(this.modal).toContainText(`Reservar ${courtName}`);

        // Select slot
        const slot = this.modal.getByRole('button', { name: timeSlot });
        await expect(slot).toBeEnabled();
        // Check it is indeed available (white bg)
        await expect(slot).toHaveClass(/bg-white/);

        // Use JS click to ensure interaction works with the custom calendar component
        await slot.evaluate((node) => (node as HTMLElement).click());

        // Check verification text (UI feedback)
        await expect(this.modal).toContainText('Seleccionado:');

        // Confirm
        await this.modal.getByRole('button', { name: 'Confirmar Reserva' }).click();
    }

    async expectSlotBooked(timeSlot: string) {
        // Used to verify disabled/booked slots
        await expect(this.modal).toBeVisible();
        const bookedSlot = this.modal.getByRole('button', { name: timeSlot });
        await expect(bookedSlot).toBeDisabled();
    }

    async expectSuccessMessage() {
        await expect(this.modal).toContainText('¡Reserva Exitosa!');
    }
}
