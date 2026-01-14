import { test, expect } from '@playwright/test';
import { RecurringRulePage } from './pages/RecurringRulePage';

test.describe('Recurring Rules', () => {
    test('should create a recurring rule with correct payload', async ({ page }) => {
        const rulesPage = new RecurringRulePage(page);
        await rulesPage.mockApis();

        // Unroute existing GET mock to avoid conflict or just override carefully
        await page.route('**/bookings/recurring*', async route => {
            const method = route.request().method();
            if (method === 'POST') {
                await route.fulfill({
                    status: 201,
                    json: {
                        data: {
                            id: 'rule-1',
                            frequency: 'WEEKLY',
                            type: 'FIXED',
                            facility_id: 'cam-1',
                            day_of_week: 1,
                            start_time: '10:00:00',
                            end_time: '11:00:00',
                            start_date: '2024-01-01',
                            end_date: '2024-12-31'
                        }
                    }
                });
            } else if (method === 'GET') {
                await route.fulfill({ status: 200, json: { data: [] } });
            } else {
                await route.continue();
            }
        });

        await rulesPage.goto();
        await rulesPage.openNewRuleModal();

        await rulesPage.fillForm({
            facilityName: 'Cancha 1',
            frequency: 'Semanal',
            dayOfWeek: 'Lunes',
            startTime: '10:00',
            endTime: '11:00',
            startDate: '2024-01-01',
            endDate: '2024-12-31'
        });

        const requestPromise = page.waitForRequest(request =>
            request.url().includes('/bookings/recurring') && request.method() === 'POST'
        );

        await rulesPage.submit();
        const request = await requestPromise;
        const capturedPayload = request.postDataJSON();

        // Critical Verification: check payload
        expect(capturedPayload).toBeTruthy(); // Ensure request was made

        await rulesPage.expectSuccessMessage();
        expect(capturedPayload!.frequency).toBe('WEEKLY');
        expect(capturedPayload!.type).toBe('FIXED');
        expect(capturedPayload!.facility_id).toBe('cam-1');
    });
});
