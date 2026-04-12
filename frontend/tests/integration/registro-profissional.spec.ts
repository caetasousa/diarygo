import { test, expect } from '@playwright/test';
import { emailUnico, limparSessao } from './helpers';

test.describe('Registro de Profissional', () => {
	test.beforeEach(async ({ page }) => {
		await limparSessao(page);
	});

	test('página exibe card de requisitos de documentação', async ({ page }) => {
		await page.goto('/registro/profissional');
		await page.waitForLoadState('networkidle');

		const infoCard = page.locator('.info-card');
		await expect(infoCard).toBeVisible();
		await expect(infoCard).toContainText('RG');
		await expect(infoCard).toContainText('Comprovante');
	});

	test('cadastro de profissional com sucesso redireciona para login', async ({ page }) => {
		await page.goto('/registro/profissional');
		await page.waitForLoadState('networkidle');

		const email = emailUnico('prof');
		const senha = 'Senha@123';

		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);

		await page.click('button[type="submit"]');

		await page.waitForURL('**/login', { timeout: 8000 });
		expect(page.url()).toContain('/login');
	});

	test('exibe toast com mensagem de análise de cadastro', async ({ page }) => {
		await page.goto('/registro/profissional');
		await page.waitForLoadState('networkidle');

		const email = emailUnico('prof');
		const senha = 'Senha@123';

		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);

		await page.click('button[type="submit"]');

		const toast = page.locator('.toast-success').first();
		await toast.waitFor({ state: 'visible', timeout: 8000 });
		const text = await toast.innerText();
		// Mensagem deve mencionar aprovação/análise
		expect(text.toLowerCase()).toMatch(/aprovação|análise|aguarde/);
	});

	test('profissional logado vê dashboard com status PENDENTE', async ({ page }) => {
		const email = emailUnico('pendente');
		const senha = 'Senha@123';

		// Registra como profissional
		await page.goto('/registro/profissional');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/login', { timeout: 8000 });

		// Loga
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard', { timeout: 8000 });

		// Dashboard deve mostrar status PENDENTE
		await expect(page.locator('.badge-orange')).toBeVisible();
		const dashContent = await page.content();
		expect(dashContent).toMatch(/PENDENTE|análise|aprovação/i);
	});

	test('email duplicado exibe erro', async ({ page }) => {
		const email = emailUnico('dupprof');
		const senha = 'Senha@123';

		await page.goto('/registro/profissional');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/login', { timeout: 8000 });

		await page.goto('/registro/profissional');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);
		await page.click('button[type="submit"]');

		const toast = page.locator('.toast-error').first();
		await toast.waitFor({ state: 'visible', timeout: 8000 });
		await expect(toast).toBeVisible();
	});
});
