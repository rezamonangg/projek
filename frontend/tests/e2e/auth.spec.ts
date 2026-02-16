import { expect, test } from '@playwright/test';

test.describe('Auth Flow', () => {
	test('redirects to login when not authenticated', async ({ page }) => {
		await page.goto('/dashboard');
		await expect(page).toHaveURL(/\/login/);
	});

	test('displays login form', async ({ page }) => {
		await page.goto('/login');
		await expect(page.locator('h2')).toContainText('Sign in');
		await expect(page.locator('input[type="email"]')).toBeVisible();
		await expect(page.locator('input[type="password"]')).toBeVisible();
		await expect(page.locator('button[type="submit"]')).toBeVisible();
	});

	test('can navigate to register page', async ({ page }) => {
		await page.goto('/login');
		await page.click('a[href="/register"]');
		await expect(page).toHaveURL(/\/register/);
		await expect(page.locator('h2')).toContainText('Create your community');
	});

	test('displays register form', async ({ page }) => {
		await page.goto('/register');
		await expect(page.locator('input[name="communityName"]')).toBeVisible();
		await expect(page.locator('input[name="adminEmail"]')).toBeVisible();
		await expect(page.locator('input[name="adminPassword"]')).toBeVisible();
		await expect(page.locator('input[name="confirmPassword"]')).toBeVisible();
	});
});
