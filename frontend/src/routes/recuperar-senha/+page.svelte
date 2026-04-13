<script lang="ts">
	import { api } from '$lib/api/client';
	import { toasts } from '$lib/stores/toasts';

	let email = '';
	let loading = false;
	let sent = false;
	let devToken = '';

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!email) return;

		loading = true;
		try {
			const resp = await api.solicitarRecuperacao({ email });
			sent = true;
			if (resp.token) devToken = resp.token; // apenas em dev
			toasts.info(resp.mensagem);
		} catch (err) {
			toasts.error(err instanceof Error ? err.message : 'Erro ao solicitar recuperação');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Recuperar senha — DiaryGo</title>
</svelte:head>

<div class="auth-page section">
	<div class="container-sm">
		<div class="auth-header">
			<a href="/" class="auth-brand">
				<span style="color: var(--color-orange-10)">✦</span> DiaryGo
			</a>
			<h1 class="text-heading auth-title">Recuperar senha</h1>
			<p class="text-body">
				Informe seu email e enviaremos instruções para redefinir sua senha.
			</p>
		</div>

		{#if sent}
			<div class="success-state card animate-in">
				<div class="success-icon">✓</div>
				<h2 class="text-subheading">Email enviado</h2>
				<p class="text-body">
					Se o email estiver cadastrado, você receberá as instruções em breve.
				</p>
				{#if devToken}
					<div class="dev-token">
						<p class="text-small" style="color: var(--color-text-tertiary); margin-bottom: 8px;">
							🛠 Ambiente de desenvolvimento — token de recuperação:
						</p>
						<code class="text-mono" style="color: var(--color-blue-10); word-break: break-all;">
							{devToken}
						</code>
						<a
							href="/redefinir-senha?token={devToken}"
							class="btn btn-primary"
							style="margin-top: var(--space-4); width: 100%; justify-content: center;"
						>
							Redefinir senha agora
						</a>
					</div>
				{/if}
				<a href="/login" class="btn btn-ghost" style="margin-top: var(--space-2);">
					← Voltar ao login
				</a>
			</div>
		{:else}
			<form class="auth-form card" on:submit={handleSubmit} novalidate>
				<div class="form-group">
					<label for="email" class="label">Email</label>
					<input
						id="email"
						type="email"
						class="input"
						placeholder="seu@email.com"
						bind:value={email}
						disabled={loading}
						autocomplete="email"
						maxlength="254"
					/>
				</div>

				<button type="submit" class="btn btn-white btn-lg" style="width: 100%;" disabled={loading || !email}>
					{#if loading}
						<span class="spinner"></span> Enviando...
					{:else}
						Enviar instruções
					{/if}
				</button>
			</form>

			<p class="text-body" style="text-align: center; margin-top: var(--space-6);">
				Lembrou a senha? <a href="/login">Entrar</a>
			</p>
		{/if}
	</div>
</div>

<style>
	.auth-page {
		display: flex;
		align-items: flex-start;
		justify-content: center;
		padding-top: var(--space-20);
	}

	.auth-header {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
		margin-bottom: var(--space-6);
	}

	.auth-brand {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		text-decoration: none;
	}

	.auth-title {
		font-size: 2rem;
		letter-spacing: -1px;
	}

	.auth-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-5);
	}

	.success-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: var(--space-4);
	}

	.success-icon {
		width: 48px;
		height: 48px;
		border-radius: 50%;
		background: var(--color-green-3);
		border: 1px solid rgba(34, 255, 153, 0.3);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.25rem;
		color: var(--color-green-10);
	}

	.dev-token {
		width: 100%;
		background: rgba(59, 158, 255, 0.05);
		border: 1px solid var(--color-blue-4);
		border-radius: var(--radius-standard);
		padding: var(--space-4);
		text-align: left;
	}

	.spinner {
		display: inline-block;
		width: 14px;
		height: 14px;
		border: 2px solid rgba(0, 0, 0, 0.3);
		border-top-color: #000;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
