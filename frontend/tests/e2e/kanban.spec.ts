import { expect, test } from '@playwright/test';

test.describe('Kanban Board', () => {
	test('redirects to login when not authenticated', async ({ page }) => {
		await page.goto('/projects/project-1/boards/board-1');
		await expect(page).toHaveURL(/\/login/);
	});

	test('displays login page with correct elements', async ({ page }) => {
		await page.goto('/login');
		await expect(page.locator('h2')).toContainText('Sign in');
		await expect(page.locator('button[type="submit"]')).toContainText('Sign in');
	});
});
