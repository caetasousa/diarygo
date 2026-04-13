import { test, expect } from '@playwright/test';
import { emailUnico, registrarELogarCliente, limparSessao } from './helpers';

const BACKEND = 'http://localhost:8080/api/v1';

/** Registra um cliente via API e retorna o email e senha */
async function criarCliente(): Promise<{ email: string; senha: string }> {
	const email = emailUnico('login');
	const senha = 'Senha@123';
	const res = await fetch(`${BACKEND}/auth/registro/cliente`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, senha })
	});
	if (!res.ok) throw new Error(`Falha criar cliente: ${await res.text()}`);
	return { email, senha };
}

test.describe('Login', () => {
	test.beforeEach(async ({ page }) => {
		await limparSessao(page);
	});

	test('login com sucesso redireciona para dashboard', async ({ page }) => {
		const { email, senha } = await criarCliente();

		await page.goto('/login');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.click('button[type="submit"]');

		await page.waitForURL('**/dashboard', { timeout: 8000 });
		expect(page.url()).toContain('/dashboard');
	});

	test('token JWT salvo no localStorage após login', async ({ page }) => {
		const { email, senha } = await criarCliente();

		await page.goto('/login');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard', { timeout: 8000 });

		const token = await page.evaluate(() => localStorage.getItem('diarygo_token'));
		expect(token).toBeTruthy();
		// JWT tem 3 partes separadas por ponto
		expect(token!.split('.').length).toBe(3);
	});

	test('dashboard exibe email do usuário logado', async ({ page }) => {
		const { email, senha } = await criarCliente();

		await page.goto('/login');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.click('button[type="submit"]');
		// Aguarda a URL do dashboard e o conteúdo estar estável
		await page.waitForURL('**/dashboard', { timeout: 8000 });
		await page.waitForLoadState('networkidle');
		// Garante que não houve redirect de volta para /login
		expect(page.url()).toContain('/dashboard');

		await expect(page.locator('text=' + email)).toBeVisible();
	});

	test('dashboard exibe badge CLIENTE', async ({ page }) => {
		const { email, senha } = await criarCliente();

		await page.goto('/login');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard', { timeout: 8000 });
		await page.waitForLoadState('networkidle');
		expect(page.url()).toContain('/dashboard');

		await expect(page.locator('.badge-green', { hasText: 'CLIENTE' })).toBeVisible();
	});

	test('credenciais inválidas exibem toast de erro', async ({ page }) => {
		await page.goto('/login');
		await page.fill('input[type="email"]', emailUnico('errado'));
		await page.fill('input[id="senha"]', 'SenhaErrada1');
		await page.click('button[type="submit"]');

		const toast = page.locator('.toast-error').first();
		await toast.waitFor({ state: 'visible', timeout: 8000 });
		await expect(toast).toBeVisible();
		expect(page.url()).not.toContain('/dashboard');
	});

	test('campos vazios exibem erros de validação inline', async ({ page }) => {
		await page.goto('/login');
		await page.waitForLoadState('networkidle');
		await page.locator('button[type="submit"]').click();

		const erros = page.locator('.form-error');
		await expect(erros.first()).toBeVisible({ timeout: 5000 });
		expect(page.url()).not.toContain('/dashboard');
	});

	test('usuário já autenticado é redirecionado do login para dashboard', async ({ page }) => {
		const token = await registrarELogarCliente(emailUnico('auto'), 'Senha@123');

		// Injeta token e vai para /login — deve redirecionar para /dashboard
		await page.goto('/');
		await page.evaluate((t) => localStorage.setItem('diarygo_token', t), token);
		await page.goto('/login');

		await page.waitForURL('**/dashboard', { timeout: 8000 });
		expect(page.url()).toContain('/dashboard');
	});

	test('logout remove token e redireciona para home', async ({ page }) => {
		const { email, senha } = await criarCliente();

		await page.goto('/login');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard', { timeout: 8000 });

		// Clica no botão "Sair" da navbar
		await page.locator('button', { hasText: 'Sair' }).click();

		await page.waitForURL('**/', { timeout: 8000 });

		const token = await page.evaluate(() => localStorage.getItem('diarygo_token'));
		expect(token).toBeNull();
	});
});
