import { test } from '@playwright/test';

test('debug registro de cliente', async ({ page }) => {
	// Captura erros de console
	page.on('console', msg => {
		if (msg.type() === 'error') console.log('CONSOLE ERROR:', msg.text());
	});
	page.on('pageerror', err => console.log('PAGE ERROR:', err.message));

	// Captura respostas de rede
	page.on('response', resp => {
		if (resp.url().includes('/auth/')) {
			console.log('RESPONSE:', resp.status(), resp.url());
		}
	});

	await page.goto('/registro');
	await page.waitForLoadState('networkidle');

	const email = `debug_reg_${Date.now()}@teste.com`;
	await page.fill('input[type="email"]', email);
	await page.fill('input[id="senha"]', 'Senha@123');
	await page.fill('input[id="confirmar"]', 'Senha@123');

	console.log('Clicando submit...');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(3000);

	console.log('URL final:', page.url());
	// Pega todos os toasts visíveis
	const toasts = await page.locator('.toast').allTextContents();
	console.log('Toasts:', toasts);
});
