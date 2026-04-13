import { test, expect } from '@playwright/test';
import { emailUnico, limparSessao } from './helpers';

test.describe('Registro de Cliente', () => {
	test.beforeEach(async ({ page }) => {
		await limparSessao(page);
	});

	test('cadastro com sucesso redireciona para login', async ({ page }) => {
		await page.goto('/registro');
		await page.waitForLoadState('networkidle');

		const email = emailUnico('cliente');
		const senha = 'Senha@123';

		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);

		await page.click('button[type="submit"]');

		// Aguarda redirecionamento para /login
		await page.waitForURL('**/login', { timeout: 8000 });
		expect(page.url()).toContain('/login');
	});

	test('cadastro com sucesso redireciona para login com mensagem', async ({ page }) => {
		await page.goto('/registro');
		await page.waitForLoadState('networkidle');

		const email = emailUnico('toast-cliente');
		const senha = 'Senha@123';

		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);

		await page.locator('button[type="submit"]').click();
		await page.waitForURL('**/login', { timeout: 12000 });

		expect(page.url()).toContain('/login');
	});

	test('email duplicado exibe erro 409', async ({ page }) => {
		const email = emailUnico('dup');
		const senha = 'Senha@123';

		// Primeira vez
		await page.goto('/registro');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/login', { timeout: 8000 });

		// Segunda vez — mesmo email
		await page.goto('/registro');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);
		await page.click('button[type="submit"]');

		const toast = page.locator('.toast-error').first();
		await toast.waitFor({ state: 'visible', timeout: 8000 });
		const text = await toast.innerText();
		expect(text.toLowerCase()).toContain('email');
	});

	test('senha fraca exibe erro de validação inline', async ({ page }) => {
		await page.goto('/registro');
		await page.fill('input[type="email"]', emailUnico('fraco'));
		await page.fill('input[id="senha"]', '123');
		await page.fill('input[id="confirmar"]', '123');
		await page.click('button[type="submit"]');

		// Erro inline — sem toast, sem redirect
		const erro = page.locator('.form-error').first();
		await expect(erro).toBeVisible();
		expect(page.url()).not.toContain('/login');
	});

	test('senhas diferentes bloqueiam envio', async ({ page }) => {
		await page.goto('/registro');
		await page.fill('input[type="email"]', emailUnico('diff'));
		await page.fill('input[id="senha"]', 'Senha@123');
		await page.fill('input[id="confirmar"]', 'Outra@456');
		await page.click('button[type="submit"]');

		const erro = page.locator('.form-error').last();
		await expect(erro).toBeVisible();
		expect(page.url()).not.toContain('/login');
	});

	test('campo email vazio exibe erro', async ({ page }) => {
		await page.goto('/registro');
		await page.waitForLoadState('networkidle');
		await page.fill('input[id="senha"]', 'Senha@123');
		await page.fill('input[id="confirmar"]', 'Senha@123');
		await page.locator('button[type="submit"]').click();

		const erro = page.locator('.form-error').first();
		await expect(erro).toBeVisible({ timeout: 5000 });
	});
});
