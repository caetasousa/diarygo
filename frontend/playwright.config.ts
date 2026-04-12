import { defineConfig, devices } from '@playwright/test';

/**
 * Testes de integração frontend ↔ backend.
 * O backend deve estar rodando em localhost:8080.
 * O frontend é iniciado automaticamente pelo Playwright via webServer.
 */
export default defineConfig({
	testDir: './tests/integration',
	timeout: 30_000,
	retries: 0,
	reporter: [['list'], ['html', { open: 'never', outputFolder: 'playwright-report' }]],

	use: {
		baseURL: 'http://localhost:5173',
		trace: 'on-first-retry',
		screenshot: 'only-on-failure',
		headless: true
	},

	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] }
		}
	],

	// Inicia o frontend em modo dev; assume backend já rodando na :8080
	webServer: {
		command: 'npm run dev',
		url: 'http://localhost:5173',
		reuseExistingServer: true,
		timeout: 30_000
	}
});
