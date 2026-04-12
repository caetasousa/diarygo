/**
 * Helpers compartilhados para testes de integração frontend ↔ backend.
 * O backend usa repositório in-memory — cada restart reseta o estado.
 * Por isso cada teste usa emails únicos via timestamp.
 */

import type { Page } from '@playwright/test';

const BACKEND_URL = 'http://localhost:8080/api/v1';

/** Gera email único para evitar conflito entre testes */
export function emailUnico(prefixo = 'teste'): string {
	return `${prefixo}_${Date.now()}_${Math.random().toString(36).slice(2, 7)}@teste.com`;
}

/** Registra cliente diretamente via API e retorna token JWT */
export async function registrarELogarCliente(email: string, senha: string): Promise<string> {
	// Registra
	const regRes = await fetch(`${BACKEND_URL}/auth/registro/cliente`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, senha })
	});
	if (!regRes.ok) {
		const body = await regRes.text();
		throw new Error(`Falha ao registrar cliente: ${regRes.status} ${body}`);
	}

	// Loga
	const loginRes = await fetch(`${BACKEND_URL}/auth/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, senha })
	});
	if (!loginRes.ok) {
		const body = await loginRes.text();
		throw new Error(`Falha ao logar: ${loginRes.status} ${body}`);
	}
	const { access_token } = await loginRes.json();
	return access_token;
}

/** Injeta token JWT no localStorage e navega para a URL */
export async function logarNaInterface(page: Page, token: string, url = '/dashboard') {
	await page.goto('/');
	await page.evaluate((t) => {
		localStorage.setItem('diarygo_token', t);
	}, token);
	await page.goto(url);
}

/** Remove o token do localStorage (simula logout).
 *  Navega para a home primeiro para garantir contexto de página válido.
 */
export async function limparSessao(page: Page) {
	try {
		// Se a página ainda não navegou para nenhuma URL do app, leva para a home primeiro
		const url = page.url();
		if (!url.startsWith('http://localhost:5173')) {
			await page.goto('/');
		}
		await page.evaluate(() => localStorage.removeItem('diarygo_token'));
	} catch {
		// Se ainda falhar, navega para home e tenta novamente
		await page.goto('/');
		await page.evaluate(() => localStorage.removeItem('diarygo_token'));
	}
}

/** Aguarda e retorna o texto do toast mais recente */
export async function aguardarToast(page: Page): Promise<string> {
	const toast = page.locator('.toast').last();
	await toast.waitFor({ state: 'visible', timeout: 8000 });
	return toast.innerText();
}
