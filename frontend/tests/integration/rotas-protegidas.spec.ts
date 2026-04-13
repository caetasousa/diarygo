import { test, expect } from '@playwright/test';
import { emailUnico, registrarELogarCliente, logarNaInterface, limparSessao } from './helpers';

test.describe('Rotas Protegidas', () => {
	test.beforeEach(async ({ page }) => {
		await limparSessao(page);
	});

	test('/dashboard sem autenticação redireciona para /login', async ({ page }) => {
		await page.goto('/dashboard');
		await page.waitForURL('**/login', { timeout: 8000 });
		expect(page.url()).toContain('/login');
	});

	test('/dashboard com token válido carrega corretamente', async ({ page }) => {
		const token = await registrarELogarCliente(emailUnico('prot'), 'Senha@123');
		await logarNaInterface(page, token, '/dashboard');

		await page.waitForLoadState('networkidle');
		await expect(page.locator('h1', { hasText: 'Dashboard' })).toBeVisible();
	});

	test('token expirado no localStorage redireciona para login', async ({ page }) => {
		// Cria um JWT expirado manualmente (exp no passado)
		// Header.Payload.Signature — assinatura inválida mas o store checa exp localmente
		const expiredPayload = btoa(JSON.stringify({
			sub: 'fake-uuid',
			email: 'fake@teste.com',
			tipo: 'CLIENTE',
			exp: Math.floor(Date.now() / 1000) - 3600, // 1h no passado
			iat: Math.floor(Date.now() / 1000) - 7200
		}));
		const fakeToken = `eyJhbGciOiJIUzI1NiJ9.${expiredPayload}.fakesig`;

		await page.goto('/');
		await page.evaluate((t) => localStorage.setItem('diarygo_token', t), fakeToken);
		await page.goto('/dashboard');

		await page.waitForURL('**/login', { timeout: 8000 });
		expect(page.url()).toContain('/login');
	});

	test('navbar mostra botão "Sair" quando autenticado', async ({ page }) => {
		const token = await registrarELogarCliente(emailUnico('navbar'), 'Senha@123');
		await logarNaInterface(page, token, '/dashboard');

		await expect(page.locator('button', { hasText: 'Sair' })).toBeVisible();
		await expect(page.locator('a', { hasText: 'Entrar' })).not.toBeVisible();
	});

	test('navbar mostra botão "Entrar" quando não autenticado', async ({ page }) => {
		await page.goto('/');

		await expect(page.locator('a', { hasText: 'Entrar' })).toBeVisible();
		await expect(page.locator('button', { hasText: 'Sair' })).not.toBeVisible();
	});

	test('rota /api/v1/me com token válido retorna dados do usuário', async ({ page }) => {
		const token = await registrarELogarCliente(emailUnico('meapi'), 'Senha@123');

		// Faz a chamada diretamente via fetch dentro do contexto da página
		const result = await page.evaluate(async (t) => {
			const res = await fetch('/api/v1/me', {
				headers: { 'Authorization': `Bearer ${t}` }
			});
			return { status: res.status, data: await res.json() };
		}, token);

		expect(result.status).toBe(200);
		expect(result.data).toHaveProperty('tipo', 'CLIENTE');
		expect(result.data).toHaveProperty('email');
	});

	test('rota /api/v1/me sem token retorna 401', async ({ page }) => {
		await page.goto('/');

		const result = await page.evaluate(async () => {
			const res = await fetch('/api/v1/me');
			return { status: res.status };
		});

		expect(result.status).toBe(401);
	});

	test('rota /api/v1/me com token inválido retorna 401', async ({ page }) => {
		await page.goto('/');

		const result = await page.evaluate(async () => {
			const res = await fetch('/api/v1/me', {
				headers: { 'Authorization': 'Bearer token.invalido.aqui' }
			});
			return { status: res.status };
		});

		expect(result.status).toBe(401);
	});
});

test.describe('Integração API — Fluxo Completo', () => {
	test('ciclo completo: registrar → login → me → logout → login bloqueado', async ({ page }) => {
		const email = emailUnico('ciclo');
		const senha = 'Senha@123';

		// 1. Registra
		await page.goto('/registro');
		await page.waitForLoadState('networkidle');
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.fill('input[id="confirmar"]', senha);
		await page.locator('button[type="submit"]').click();
		await page.waitForURL('**/login', { timeout: 12000 });

		// 2. Loga
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senha);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard', { timeout: 8000 });

		// 3. Confirma autenticação via /me
		const token = await page.evaluate(() => localStorage.getItem('diarygo_token'));
		const meResult = await page.evaluate(async (t) => {
			const res = await fetch('/api/v1/me', {
				headers: { 'Authorization': `Bearer ${t}` }
			});
			return res.status;
		}, token);
		expect(meResult).toBe(200);

		// 4. Faz logout
		await page.locator('button', { hasText: 'Sair' }).click();
		await page.waitForURL('**/', { timeout: 8000 });

		// 5. Token removido do localStorage
		const tokenPosLogout = await page.evaluate(() => localStorage.getItem('diarygo_token'));
		expect(tokenPosLogout).toBeNull();

		// 6. /dashboard redireciona para login
		await page.goto('/dashboard');
		await page.waitForURL('**/login', { timeout: 8000 });
		expect(page.url()).toContain('/login');
	});

	test('ciclo completo de recuperação de senha: solicitar → redefinir → login com nova senha', async ({ page }) => {
		const email = emailUnico('full');
		const senhaOriginal = 'Senha@123';
		const senhaNova = 'NovaSenha@999';

		// Registra
		await fetch('http://localhost:8080/api/v1/auth/registro/cliente', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email, senha: senhaOriginal })
		});

		// Solicita recuperação
		await page.goto('/recuperar-senha');
		await page.fill('input[type="email"]', email);
		await page.click('button[type="submit"]');
		await page.locator('.success-state').waitFor({ state: 'visible', timeout: 8000 });

		// Clica no link "Redefinir senha agora" (dev token exibido)
		await page.locator('a', { hasText: 'Redefinir senha agora' }).click();
		await page.waitForLoadState('networkidle');

		// Redefine a senha
		await page.fill('input[id="novaSenha"]', senhaNova);
		await page.fill('input[id="confirmar"]', senhaNova);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/login', { timeout: 8000 });

		// Tenta login com senha antiga — deve falhar
		await page.fill('input[type="email"]', email);
		await page.fill('input[id="senha"]', senhaOriginal);
		await page.click('button[type="submit"]');

		const toastErro = page.locator('.toast-error').first();
		await toastErro.waitFor({ state: 'visible', timeout: 8000 });
		expect(page.url()).not.toContain('/dashboard');

		// Faz login com nova senha — deve funcionar
		await page.fill('input[id="senha"]', senhaNova);
		await page.click('button[type="submit"]');
		await page.waitForURL('**/dashboard', { timeout: 8000 });
		expect(page.url()).toContain('/dashboard');
	});
});
