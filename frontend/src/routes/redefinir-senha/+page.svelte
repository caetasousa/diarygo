<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import { toasts } from '$lib/stores/toasts';

	let token = $page.url.searchParams.get('token') ?? '';
	let novaSenha = '';
	let confirmar = '';
	let loading = false;
	let errors: { token?: string; novaSenha?: string; confirmar?: string } = {};

	function validate() {
		errors = {};
		if (!token) errors.token = 'Token obrigatório';
		if (!novaSenha) errors.novaSenha = 'Nova senha obrigatória';
		else if (novaSenha.length < 8) errors.novaSenha = 'Mínimo 8 caracteres';
		if (novaSenha && confirmar !== novaSenha) errors.confirmar = 'Senhas não conferem';
		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!validate()) return;

		loading = true;
		try {
			await api.redefinirSenha({ token, nova_senha: novaSenha });
			toasts.success('Senha redefinida com sucesso!');
			goto('/login');
		} catch (err) {
			toasts.error(err instanceof Error ? err.message : 'Erro ao redefinir senha');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Redefinir senha — DiaryGo</title>
</svelte:head>

<div class="auth-page section">
	<div class="container-sm">
		<div class="auth-header">
			<a href="/" class="auth-brand">
				<span style="color: var(--color-orange-10)">✦</span> DiaryGo
			</a>
			<h1 class="text-heading auth-title">Nova senha</h1>
			<p class="text-body">Escolha uma senha forte para sua conta.</p>
		</div>

		<form class="auth-form card" on:submit={handleSubmit} novalidate>
			{#if !$page.url.searchParams.get('token')}
				<div class="form-group">
					<label for="token" class="label">Token de recuperação</label>
					<input
						id="token"
						type="text"
						class="input"
						placeholder="Cole o token recebido"
						bind:value={token}
						on:input={() => { errors = { ...errors, token: undefined }; }}
						disabled={loading}
					/>
					{#if errors.token}<span class="form-error">{errors.token}</span>{/if}
				</div>
			{/if}

			<div class="form-group">
				<label for="novaSenha" class="label">Nova senha</label>
				<input
					id="novaSenha"
					type="password"
					class="input"
					placeholder="Mínimo 8 caracteres"
					bind:value={novaSenha}
					on:input={() => { errors = { ...errors, novaSenha: undefined }; }}
					disabled={loading}
					autocomplete="new-password"
				/>
				{#if errors.novaSenha}<span class="form-error">{errors.novaSenha}</span>{/if}
			</div>

			<div class="form-group">
				<label for="confirmar" class="label">Confirmar nova senha</label>
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

			<button type="submit" class="btn btn-white btn-lg" style="width: 100%;" disabled={loading}>
				{#if loading}
					<span class="spinner"></span> Redefinindo...
				{:else}
					Redefinir senha
				{/if}
			</button>
		</form>

		<p class="text-body" style="text-align: center; margin-top: var(--space-6);">
			<a href="/login">← Voltar ao login</a>
		</p>
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
