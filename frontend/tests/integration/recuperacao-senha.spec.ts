import { test, expect } from '@playwright/test';
import { emailUnico, limparSessao } from './helpers';

const BACKEND = 'http://localhost:8080/api/v1';

/** Cria usuário e solicita recuperação de senha via API, retornando o token de recuperação */
async function criarEsolicitarRecuperacao(email: string, senha: string): Promise<string> {
	// Registra
	await fetch(`${BACKEND}/auth/registro/cliente`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, senha })
	});

	// Solicita recuperação
	const res = await fetch(`${BACKEND}/auth/solicitar-recuperacao-senha`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email })
	});
	const data = await res.json();
	// Em ENV=development o backend retorna o token na resposta
	if (!data.token) throw new Error('Backend não retornou token (verifique ENV=development)');
	return data.token as string;
}

test.describe('Recuperação de Senha', () => {
	test.beforeEach(async ({ page }) => {
		await limparSessao(page);
	});

	test('formulário de recuperação exibe estado de sucesso após envio', async ({ page }) => {
		const email = emailUnico('rec');

		// Registra o usuário via API primeiro
		await fetch(`${BACKEND}/auth/registro/cliente`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email, senha: 'Senha@123' })
		});

		await page.goto('/recuperar-senha');
		await page.fill('input[type="email"]', email);
		await page.click('button[type="submit"]');

		// Deve mostrar o estado de sucesso
		const successState = page.locator('.success-state');
		await successState.waitFor({ state: 'visible', timeout: 8000 });
		await expect(successState).toBeVisible();
		await expect(successState).toContainText('Email enviado');
	});

	test('email inexistente também exibe sucesso (OWASP A07 — sem enumeração)', async ({ page }) => {
		await page.goto('/recuperar-senha');
		await page.fill('input[type="email"]', emailUnico('naoexiste'));
		await page.click('button[type="submit"]');

		// Deve mostrar sucesso mesmo para email inexistente
		const successState = page.locator('.success-state');
		await successState.waitFor({ state: 'visible', timeout: 8000 });
		await expect(successState).toBeVisible();
	});

	test('em desenvolvimento, exibe token de recuperação e link direto', async ({ page }) => {
		const email = emailUnico('devtoken');
		await fetch(`${BACKEND}/auth/registro/cliente`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email, senha: 'Senha@123' })
		});

		await page.goto('/recuperar-senha');
		await page.fill('input[type="email"]', email);
		await page.click('button[type="submit"]');

		await page.locator('.success-state').waitFor({ state: 'visible', timeout: 8000 });

		// O token de dev e o link direto devem aparecer
		const devToken = page.locator('.dev-token');
		await expect(devToken).toBeVisible();

		const link = devToken.locator('a', { hasText: 'Redefinir senha agora' });
		await expect(link).toBeVisible();
	});
});

test.describe('Redefinição de Senha', () => {
	test.beforeEach(async ({ page }) => {
		await limparSessao(page);
	});

	test('redefinição com token válido redireciona para login', async ({ page }) => {
		const email = emailUnico('redef');
		const senha = 'Senha@123';
		const token = await criarEsolicitarRecuperacao(email, senha);

		await page.goto(`/redefinir-senha?token=${token}`);
		await page.fill('input[id="novaSenha"]', 'NovaSenha@456');
		await page.fill('input[id="confirmar"]', 'NovaSenha@456');
		await page.click('button[type="submit"]');

		await page.waitForURL('**/login', { timeout: 8000 });
		expect(page.url()).toContain('/login');
	});

	test('após redefinição, login com nova senha funciona', async ({ page }) => {
		const email = emailUnico('newpass');
		const senhaNova = 'NovaSenha@789';
		const token = await criarEsolicitarRecuperacao(email, 'Senha@123');

		// Redefine via interface
		await page.goto(`/redefinir-senha?token=${token}`);
		await page.fill('input[id="novaSenha"]', senhaNova);
		await page.fill('input[id="confirmar"]', senhaNova);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/login', { timeout: 8000 });

		// Loga com nova senha
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senhaNova);
		await page.click('button[type="submit"]');

		await page.waitForURL('**/dashboard', { timeout: 8000 });
		expect(page.url()).toContain('/dashboard');
	});

	test('token inválido exibe toast de erro', async ({ page }) => {
		await page.goto('/redefinir-senha?token=token-invalido-abc123');
		await page.fill('input[id="novaSenha"]', 'NovaSenha@456');
		await page.fill('input[id="confirmar"]', 'NovaSenha@456');
		await page.click('button[type="submit"]');

		const toast = page.locator('.toast-error').first();
		await toast.waitFor({ state: 'visible', timeout: 8000 });
		await expect(toast).toBeVisible();
		expect(page.url()).not.toContain('/login');
	});

	test('senhas não conferem bloqueiam envio', async ({ page }) => {
		const email = emailUnico('mismatch');
		const token = await criarEsolicitarRecuperacao(email, 'Senha@123');

		await page.goto(`/redefinir-senha?token=${token}`);
		await page.fill('input[id="novaSenha"]', 'NovaSenha@456');
		await page.fill('input[id="confirmar"]', 'Diferente@789');
		await page.click('button[type="submit"]');

		const erro = page.locator('.form-error').last();
		await expect(erro).toBeVisible();
		expect(page.url()).not.toContain('/login');
	});

	test('senha fraca exibe erro de validação', async ({ page }) => {
		const email = emailUnico('senhafraca');
		const token = await criarEsolicitarRecuperacao(email, 'Senha@123');

		await page.goto(`/redefinir-senha?token=${token}`);
		await page.fill('input[id="novaSenha"]', '123');
		await page.fill('input[id="confirmar"]', '123');
		await page.click('button[type="submit"]');

		const erro = page.locator('.form-error').first();
		await expect(erro).toBeVisible();
	});

	test('sem token na URL exibe campo de token manual', async ({ page }) => {
		await page.goto('/redefinir-senha');

		// O campo de token manual deve aparecer quando não há token na URL
		const tokenInput = page.locator('input[id="token"]');
		await expect(tokenInput).toBeVisible();
	});

	test('redefinição via campo de token manual funciona', async ({ page }) => {
		const email = emailUnico('manual');
		const token = await criarEsolicitarRecuperacao(email, 'Senha@123');

		await page.goto('/redefinir-senha');

		// Preenche token manualmente
		await page.fill('input[id="token"]', token);
		await page.fill('input[id="novaSenha"]', 'NovaSenha@456');
		await page.fill('input[id="confirmar"]', 'NovaSenha@456');
		await page.click('button[type="submit"]');

		await page.waitForURL('**/login', { timeout: 8000 });
		expect(page.url()).toContain('/login');
	});
});
