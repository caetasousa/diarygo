import { test, expect } from '@playwright/test';

test('debug goto após login', async ({ page }) => {
	const BACKEND = 'http://localhost:8080/api/v1';
	const email = `debug_${Date.now()}@teste.com`;
	const senha = 'Senha@123';

	await fetch(`${BACKEND}/auth/registro/cliente`, {
		method: 'POST', headers: {'Content-Type':'application/json'},
		body: JSON.stringify({email, senha})
	});

	await page.goto('/login');
	await page.fill('input[type="email"]', email);
	await page.fill('input[id="senha"]', senha);
	await page.click('button[type="submit"]');
	await page.waitForTimeout(3000);

	console.log('URL:', page.url());
	const heading = await page.locator('h1').first().textContent();
	console.log('H1:', heading);
});
