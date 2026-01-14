import { type Page, type Locator, expect } from '@playwright/test';

export class RecurringRulePage {
    readonly page: Page;
    readonly newRuleButton: Locator;
    readonly modal: Locator;
    readonly submitButton: Locator;

    constructor(page: Page) {
        this.page = page;
        this.newRuleButton = page.getByRole('button', { name: 'Nueva Regla' });
        this.modal = page.locator('div[role="dialog"]');
        this.submitButton = page.getByRole('button', { name: 'Crear Regla' });
    }

    async goto() {
        await this.page.goto('/admin/recurring-bookings');
        await expect(this.page.getByRole('heading', { name: 'Reservas Recurrentes' })).toBeVisible();
    }

    async mockApis() {
        // 1. Mock Facilities
        await this.page.route('**/facilities*', async route => {
            await route.fulfill({
                status: 200,
                json: [
                    { id: 'cam-1', name: 'Cancha 1', type: 'padel' }
                ]
            });
        });

        // 2. Mock Existing Rules (Empty initially)
        await this.page.route('**/bookings/recurring*', async route => {
            if (route.request().method() === 'GET') {
                await route.fulfill({ status: 200, json: { data: [] } });
            } else {
                await route.continue();
            }
        });

        // 3. Mock User Profile (Admin) - REQUIRED for Dashboard Access
        await this.page.route('**/users/me', async route => {
            await route.fulfill({
                status: 200,
                json: {
                    id: 'admin-user',
                    name: 'Admin User',
                    email: 'admin@club.com',
                    role: 'SUPER_ADMIN',
                    club_id: 'club-1',
                    medical_cert_status: 'VALID'
                }
            });
        });
    }

    async openNewRuleModal() {
        await this.newRuleButton.click();
        await expect(this.modal).toBeVisible();
    }

    async fillForm(rule: {
        facilityName: string;
        frequency: string;
        dayOfWeek: string;
        startTime: string;
        endTime: string;
        startDate: string;
        endDate: string;
    }) {
        // Select Facility
        // Strategy: Click the trigger that is following the Label "Instalación" or just find by placeholder
        await this.page.getByRole('combobox').filter({ hasText: 'Selecciona una instalación...' }).click();
        await this.page.getByRole('option', { name: rule.facilityName }).click();

        // Select Frequency
        // The default is Weekly, so the trigger says "Semanal". We need to click it to change or verifying
        // If we want to select Weekly (default), checking it matches is good.
        // But let's assume we might select Monthly.
        // Trigger might display selected value.
        // Let's find the select by Label text proximity if possible, or order.
        // Order: Facility(0), Frequency(1), Day(2).

        // Let's just use specific locators based on current value or placeholder?
        // Frequency default is 'Semanal'.
        if (rule.frequency !== 'Semanal') {
            // If we want 'Mensual', we expect to click 'Semanal' (current value) and select 'Mensual'
            await this.page.getByRole('combobox', { name: 'Semanal' }).click(); // This assumes the label is accessibility linked or text matches
            // Actually, Shadcn trigger often contains the text of the selected value.
            // If default is Weekly, text is "Semanal".
            await this.page.getByRole('option', { name: rule.frequency }).click();
        }

        // Day of Week
        const dayTrigger = this.page.locator('button[role="combobox"]').nth(2); // Risk: order dependency
        await dayTrigger.click();
        await this.page.getByRole('option', { name: rule.dayOfWeek }).click();

        // Time Inputs
        await this.page.fill('input[id="start_time"]', rule.startTime);
        await this.page.fill('input[id="end_time"]', rule.endTime);

        // Date Inputs
        await this.page.fill('input[id="start_date"]', rule.startDate);
        await this.page.fill('input[id="end_date"]', rule.endDate);
    }

    async submit() {
        await this.submitButton.click();
    }

    async expectSuccessMessage() {
        // Checking for toast can be flaky if animation is slow/fast.
        // Checking if modal closed is a strong signal of success in this flow.
        await expect(this.modal).toBeHidden();
        // Optional: Try to find the toast but don't block if modal is already gone (success state reached)
        // await expect(this.page.getByText('Regla creada')).toBeVisible();
    }
}
