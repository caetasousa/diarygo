<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { toasts } from '$lib/stores/toasts';
	import { onMount } from 'svelte';

	let email = '';
	let senha = '';
	let loading = false;
	let errors: { email?: string; senha?: string } = {};

	onMount(() => {
		if ($isAuthenticated) goto('/dashboard');
	});

	function validate() {
		errors = {};
		if (!email) errors.email = 'Email obrigatório';
		else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) errors.email = 'Email inválido';
		if (!senha) errors.senha = 'Senha obrigatória';
		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!validate()) return;

		loading = true;
		try {
			const resp = await api.login({ email, senha });
			auth.login(resp.access_token);
			toasts.success('Login realizado com sucesso!');
			goto('/dashboard');
		} catch (err) {
			toasts.error(err instanceof Error ? err.message : 'Erro ao fazer login');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Entrar — DiaryGo</title>
</svelte:head>

<div class="auth-page section">
	<div class="container-sm">
		<div class="auth-header">
			<a href="/" class="auth-brand">
				<span style="color: var(--color-orange-10)">✦</span> DiaryGo
			</a>
			<h1 class="text-heading auth-title">Entrar</h1>
			<p class="text-body">Acesse sua conta para continuar.</p>
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
					disabled={loading}
					autocomplete="email"
				/>
				{#if errors.email}
					<span class="form-error">{errors.email}</span>
				{/if}
			</div>

			<div class="form-group">
				<div class="row" style="justify-content: space-between; align-items: center;">
					<label for="senha" class="label" style="margin: 0;">Senha</label>
					<a href="/recuperar-senha" class="text-caption" style="color: var(--color-blue-10);">
						Esqueceu a senha?
					</a>
				</div>
				<input
					id="senha"
					type="password"
					class="input"
					placeholder="••••••••"
					bind:value={senha}
					disabled={loading}
					autocomplete="current-password"
				/>
				{#if errors.senha}
					<span class="form-error">{errors.senha}</span>
				{/if}
			</div>

			<button type="submit" class="btn btn-white btn-lg" style="width: 100%;" disabled={loading}>
				{#if loading}
					<span class="spinner"></span> Entrando...
				{:else}
					Entrar
				{/if}
			</button>
		</form>

		<div class="divider" style="margin: var(--space-6) 0;"></div>

		<p class="text-body" style="text-align: center;">
			Não tem conta?
			<a href="/registro">Criar conta de cliente</a>
			·
			<a href="/registro/profissional">Sou diarista</a>
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
