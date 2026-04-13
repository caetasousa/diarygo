<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import { toasts } from '$lib/stores/toasts';

	let email = '';
	let senha = '';
	let confirmar = '';
	let loading = false;
	let errors: { email?: string; senha?: string; confirmar?: string } = {};

	const requirements = [
		{ icon: '📄', label: 'RG (frente e verso) — enviado após aprovação inicial' },
		{ icon: '🏠', label: 'Comprovante de residência recente' },
		{ icon: '📸', label: 'Foto do rosto (selfie nítida)' },
		{ icon: '📞', label: 'Mínimo 2 referências profissionais com telefone' },
		{ icon: '📍', label: 'CEP das regiões onde deseja atender' }
	];

	function validate() {
		errors = {};
		if (!email) errors.email = 'Email obrigatório';
		else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) errors.email = 'Email inválido';
		if (!senha) errors.senha = 'Senha obrigatória';
		else if (senha.length < 8) errors.senha = 'Mínimo 8 caracteres';
		if (senha && confirmar !== senha) errors.confirmar = 'Senhas não conferem';
		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!validate()) return;

		loading = true;
		try {
			await api.registrarProfissional({ email, senha });
			toasts.success('Cadastro enviado! Aguarde aprovação da equipe.');
			goto('/login');
		} catch (err) {
			toasts.error(err instanceof Error ? err.message : 'Erro ao criar conta');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Cadastro de Diarista — DiaryGo</title>
</svelte:head>

<div class="auth-page section">
	<div class="container-sm">
		<div class="auth-header">
			<a href="/" class="auth-brand">
				<span style="color: var(--color-orange-10)">✦</span> DiaryGo
			</a>
			<div class="row gap-3" style="align-items: center; flex-wrap: wrap;">
				<h1 class="text-heading auth-title">Seja diarista</h1>
				<span class="badge badge-orange">Profissional</span>
			</div>
			<p class="text-body">
				Cadastre-se e comece a receber solicitações na sua região.
				Seu perfil será analisado em até 48h úteis.
			</p>
		</div>

		<div class="info-card card" style="margin-bottom: var(--space-6);">
			<h3 class="text-subheading" style="margin-bottom: var(--space-4);">O que você precisará</h3>
			<ul class="requirements-list">
				{#each requirements as req}
					<li class="requirement-item row gap-3">
						<span class="req-icon">{req.icon}</span>
						<span class="text-caption">{req.label}</span>
					</li>
				{/each}
			</ul>
		</div>

		<form class="auth-form card" on:submit={handleSubmit} novalidate>
			<div class="form-group">
				<label for="email" class="label">Email profissional</label>
				<input
					id="email"
					type="email"
					class="input"
					placeholder="seu@email.com"
					bind:value={email}
					on:input={() => { errors = { ...errors, email: undefined }; }}
					disabled={loading}
					autocomplete="email"
					maxlength="254"
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
					maxlength="72"
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
					maxlength="72"
				/>
				{#if errors.confirmar}<span class="form-error">{errors.confirmar}</span>{/if}
			</div>

			<button type="submit" class="btn btn-white btn-lg" style="width: 100%;" disabled={loading}>
				{#if loading}
					<span class="spinner"></span> Enviando cadastro...
				{:else}
					Enviar cadastro
				{/if}
			</button>

			<p class="text-caption" style="text-align: center; color: var(--color-text-tertiary);">
				Após o envio, sua documentação será analisada em até 48h úteis.
			</p>
		</form>

		<div class="divider" style="margin: var(--space-6) 0;"></div>
		<p class="text-body" style="text-align: center;">
			Já tem conta? <a href="/login">Entrar</a>
			· Precisa de serviço? <a href="/registro">Sou cliente</a>
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

	.info-card {
		border-color: var(--color-orange-4);
	}

	.requirements-list {
		list-style: none;
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.req-icon {
		font-size: 1rem;
		flex-shrink: 0;
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
