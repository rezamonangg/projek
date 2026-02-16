import { expect, test } from '@playwright/test';

test.describe('Projects', () => {
	test('displays projects list page after login', async ({ page }) => {
		await page.goto('/projects');
		await expect(page).toHaveURL(/\/login/);
		await expect(page.locator('h2')).toContainText('Sign in');
	});

	test('displays dashboard layout elements', async ({ page }) => {
		await page.goto('/login');
		await expect(page.locator('h1')).toContainText('Projek');
		await expect(page.locator('text=Demo credentials')).toBeVisible();
	});
});
