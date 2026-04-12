<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import { toasts } from '$lib/stores/toasts';

	let email = '';
	let senha = '';
	let confirmar = '';
	let loading = false;
	let errors: { email?: string; senha?: string; confirmar?: string } = {};

	function validate() {
		errors = {};
		if (!email) errors.email = 'Email obrigatório';
		else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) errors.email = 'Email inválido';
		if (!senha) errors.senha = 'Senha obrigatória';
		else if (senha.length < 8) errors.senha = 'Mínimo 8 caracteres';
		else if (senha.length > 72) errors.senha = 'Máximo 72 caracteres';
		if (senha && confirmar !== senha) errors.confirmar = 'Senhas não conferem';
		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!validate()) return;

		loading = true;
		try {
			await api.registrarCliente({ email, senha });
			toasts.success('Conta criada! Faça login para continuar.');
			goto('/login');
		} catch (err) {
			toasts.error(err instanceof Error ? err.message : 'Erro ao criar conta');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Criar conta — DiaryGo</title>
</svelte:head>

<div class="auth-page section">
	<div class="container-sm">
		<div class="auth-header">
			<a href="/" class="auth-brand">
				<span style="color: var(--color-orange-10)">✦</span> DiaryGo
			</a>
			<div class="row gap-3" style="align-items: center; flex-wrap: wrap;">
				<h1 class="text-heading auth-title">Criar conta</h1>
				<span class="badge badge-green">Cliente</span>
			</div>
			<p class="text-body">Crie sua conta e agende seu primeiro serviço.</p>
		</div>

		<form class="auth-form card" on:submit={handleSubmit} novalidate>
			<div class="form-group">
				<label for="email" class="label">Email</label>
				<input
					id="email"
					type="email"
					class="input"
					placeholder="seu@email.com"
					bind:value={email}
					on:input={() => { errors = { ...errors, email: undefined }; }}
					disabled={loading}
					autocomplete="email"
				/>
				{#if errors.email}<span class="form-error">{errors.email}</span>{/if}
			</div>

			<div class="form-group">
				<label for="senha" class="label">Senha</label>
				<input
					id="senha"
					type="password"
					class="input"
					placeholder="Mínimo 8 caracteres"
					bind:value={senha}
					on:input={() => { errors = { ...errors, senha: undefined }; }}
					disabled={loading}
					autocomplete="new-password"
				/>
				{#if errors.senha}<span class="form-error">{errors.senha}</span>{/if}
			</div>

			<div class="form-group">
				<label for="confirmar" class="label">Confirmar senha</label>
				<input
					id="confirmar"
					type="password"
					class="input"
					placeholder="••••••••"
					bind:value={confirmar}
					on:input={() => { errors = { ...errors, confirmar: undefined }; }}
					disabled={loading}
					autocomplete="new-password"
				/>
				{#if errors.confirmar}<span class="form-error">{errors.confirmar}</span>{/if}
			</div>

			<div class="password-strength">
				<div class="strength-bar">
					{#each Array(4) as _, i}
						<div
							class="strength-segment"
							class:active={senha.length > i * 18}
							class:strong={senha.length >= 12}
						></div>
					{/each}
				</div>
				<span class="text-small" style="color: var(--color-text-tertiary);">
					{#if senha.length === 0}Digite uma senha
					{:else if senha.length < 8}Muito curta
					{:else if senha.length < 12}Razoável
					{:else if senha.length < 20}Boa
					{:else}Excelente{/if}
				</span>
			</div>

			<button type="submit" class="btn btn-white btn-lg" style="width: 100%;" disabled={loading}>
				{#if loading}
					<span class="spinner"></span> Criando conta...
				{:else}
					Criar conta de cliente
				{/if}
			</button>

			<p class="text-caption" style="text-align: center; color: var(--color-text-tertiary);">
				Ao criar conta, você concorda com nossos termos de uso.
			</p>
		</form>

		<div class="divider" style="margin: var(--space-6) 0;"></div>

		<p class="text-body" style="text-align: center;">
			Já tem conta? <a href="/login">Entrar</a>
			· É diarista? <a href="/registro/profissional">Cadastre-se aqui</a>
		</p>
	</div>
</div>

<style>
	.auth-page {
		display: flex;
		align-items: flex-start;
		justify-content: center;
		padding-top: var(--space-16);
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

	.password-strength {
		display: flex;
		align-items: center;
		gap: var(--space-3);
	}

	.strength-bar {
		display: flex;
		gap: 4px;
		flex: 1;
	}

	.strength-segment {
		flex: 1;
		height: 3px;
		background: var(--border-frost);
		border-radius: var(--radius-pill);
		transition: background 0.3s ease;
	}

	.strength-segment.active {
		background: var(--color-orange-10);
	}

	.strength-segment.active.strong {
		background: var(--color-green-10);
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
