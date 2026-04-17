<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import CpfInput from '$lib/components/CpfInput.svelte';
	import TelefoneInput from '$lib/components/TelefoneInput.svelte';
	import type { ClienteResponse, ProfissionalResponse } from '$lib/types';

	let tipo = $derived($auth.payload?.tipo ?? null);
	let carregando = $state(true);
	let salvando = $state(false);

	let perfilCliente = $state<ClienteResponse | null>(null);
	let nomeCliente = $state('');
	let cpfCliente = $state('');
	let telefoneCliente = $state('');
	let erros = $state<Record<string, string>>({});

	let perfilProf = $state<ProfissionalResponse | null>(null);
	let nomeProf = $state('');
	let cpfProf = $state('');
	let rgProf = $state('');
	let telefoneProf = $state('');
	let fotoUrlProf = $state('');
	let meiProf = $state(false);

	onMount(async () => {
		try {
			if (tipo === 'CLIENTE') {
				perfilCliente = await api.buscarPerfilCliente().catch(() => null);
				if (perfilCliente) {
					nomeCliente = perfilCliente.nome;
					cpfCliente = perfilCliente.cpf;
					telefoneCliente = perfilCliente.telefone;
				}
			} else if (tipo === 'PROFISSIONAL') {
				perfilProf = await api.buscarPerfilProfissional().catch(() => null);
				if (perfilProf) {
					nomeProf = perfilProf.nome;
					cpfProf = perfilProf.cpf;
					rgProf = perfilProf.rg;
					telefoneProf = perfilProf.telefone;
					fotoUrlProf = perfilProf.foto_url;
					meiProf = perfilProf.mei;
				}
			}
		} finally {
			carregando = false;
		}
	});

	function validarCliente(): boolean {
		erros = {};
		if (!nomeCliente.trim()) erros.nome = 'Nome é obrigatório';
		if (cpfCliente.replace(/\D/g, '').length !== 11) erros.cpf = 'CPF deve ter 11 dígitos';
		if (telefoneCliente.replace(/\D/g, '').length < 10) erros.telefone = 'Telefone inválido';
		return Object.keys(erros).length === 0;
	}

	function validarProfissional(): boolean {
		erros = {};
		if (!nomeProf.trim()) erros.nome = 'Nome é obrigatório';
		if (cpfProf.replace(/\D/g, '').length !== 11) erros.cpf = 'CPF deve ter 11 dígitos';
		if (telefoneProf.replace(/\D/g, '').length < 10) erros.telefone = 'Telefone inválido';
		return Object.keys(erros).length === 0;
	}

	async function salvarCliente() {
		if (!validarCliente()) return;
		salvando = true;
		try {
			const req = {
				nome: nomeCliente.trim(),
				cpf: cpfCliente.replace(/\D/g, ''),
				telefone: telefoneCliente.replace(/\D/g, '')
			};
			perfilCliente = perfilCliente
				? await api.atualizarPerfilCliente(req)
				: await api.criarPerfilCliente(req);
			toasts.success('Perfil salvo com sucesso!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar perfil');
		} finally {
			salvando = false;
		}
	}

	async function salvarProfissional() {
		if (!validarProfissional()) return;
		salvando = true;
		try {
			const req = {
				nome: nomeProf.trim(),
				cpf: cpfProf.replace(/\D/g, ''),
				rg: rgProf.trim(),
				telefone: telefoneProf.replace(/\D/g, ''),
				foto_url: fotoUrlProf.trim(),
				mei: meiProf
			};
			perfilProf = perfilProf
				? await api.atualizarPerfilProfissional(req)
				: await api.criarPerfilProfissional(req);
			toasts.success('Perfil salvo com sucesso!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar perfil');
		} finally {
			salvando = false;
		}
	}

	const statusLabel: Record<string, string> = {
		PENDENTE: 'Pendente de aprovação',
		APROVADA: 'Aprovada',
		REPROVADA: 'Reprovada',
		SUSPENSA: 'Suspensa',
		DESCREDENCIADA: 'Descredenciada'
	};
	const statusBadge: Record<string, string> = {
		PENDENTE: 'badge-orange',
		APROVADA: 'badge-green',
		REPROVADA: 'badge-red',
		SUSPENSA: 'badge-red',
		DESCREDENCIADA: 'badge-red'
	};
</script>

<svelte:head>
	<title>Meu Perfil — DiaryGo</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<div>
			<h1 class="page-title">Meu perfil</h1>
			<p class="page-sub">Informações pessoais e dados de conta.</p>
		</div>
		{#if tipo === 'PROFISSIONAL' && perfilProf}
			<span class="badge {statusBadge[perfilProf.status] ?? 'badge-blue'}">
				{statusLabel[perfilProf.status] ?? perfilProf.status}
			</span>
		{/if}
	</div>

	{#if carregando}
		<div class="loading-state">
			<div class="loading-spinner"></div>
			<span>Carregando...</span>
		</div>
	{:else if tipo === 'CLIENTE'}

		<div class="form-layout">
			<div class="form-section">
				<h2 class="form-section-title">Dados pessoais</h2>
				<form novalidate onsubmit={(e) => { e.preventDefault(); salvarCliente(); }}>
					<div class="form-group">
						<label for="nome-cliente" class="form-label">Nome completo</label>
						<input
							id="nome-cliente"
							type="text"
							autocomplete="name"
							maxlength="100"
							placeholder="Seu nome completo"
							class="form-input"
							class:input-error={!!erros.nome}
							bind:value={nomeCliente}
						/>
						{#if erros.nome}<p class="form-error">{erros.nome}</p>{/if}
					</div>

					<CpfInput bind:value={cpfCliente} error={erros.cpf} />
					<TelefoneInput bind:value={telefoneCliente} error={erros.telefone} />

					<div class="form-actions">
						<button type="submit" class="btn btn-white" disabled={salvando}>
							{salvando ? 'Salvando...' : (perfilCliente ? 'Salvar alterações' : 'Criar perfil')}
						</button>
					</div>
				</form>
			</div>

			{#if perfilCliente}
				<div class="info-panel">
					<h2 class="form-section-title">Score de confiança</h2>
					<div class="score-display">
						<span class="score-value">{perfilCliente.score}</span>
						<span class="score-max">/100</span>
					</div>
					<div class="score-bar">
						<div class="score-fill" style="width: {perfilCliente.score}%"></div>
					</div>
					<p class="score-note">O score aumenta com serviços concluídos e avaliações positivas.</p>
				</div>
			{/if}
		</div>

	{:else if tipo === 'PROFISSIONAL'}

		<div class="form-layout">
			<div class="form-section">
				<h2 class="form-section-title">Dados pessoais</h2>
				<form novalidate onsubmit={(e) => { e.preventDefault(); salvarProfissional(); }}>
					<div class="form-group">
						<label for="nome-prof" class="form-label">Nome completo</label>
						<input
							id="nome-prof"
							type="text"
							autocomplete="name"
							maxlength="100"
							placeholder="Seu nome completo"
							class="form-input"
							class:input-error={!!erros.nome}
							bind:value={nomeProf}
						/>
						{#if erros.nome}<p class="form-error">{erros.nome}</p>{/if}
					</div>

					<CpfInput bind:value={cpfProf} error={erros.cpf} />

					<div class="form-group">
						<label for="rg-prof" class="form-label">RG</label>
						<input
							id="rg-prof"
							type="text"
							autocomplete="off"
							maxlength="20"
							placeholder="Número do RG"
							class="form-input"
							bind:value={rgProf}
						/>
					</div>

					<TelefoneInput bind:value={telefoneProf} error={erros.telefone} />

					<div class="form-group">
						<label for="foto-url" class="form-label">URL da foto de perfil</label>
						<input
							id="foto-url"
							type="url"
							autocomplete="off"
							maxlength="500"
							placeholder="https://..."
							class="form-input"
							bind:value={fotoUrlProf}
						/>
					</div>

					<div class="checkbox-row">
						<input type="checkbox" id="mei" bind:checked={meiProf} />
						<label for="mei" class="checkbox-label">
							Possuo MEI (Microempreendedor Individual)
						</label>
					</div>

					<div class="form-actions">
						<button type="submit" class="btn btn-white" disabled={salvando}>
							{salvando ? 'Salvando...' : (perfilProf ? 'Salvar alterações' : 'Criar perfil')}
						</button>
					</div>
				</form>
			</div>
		</div>

	{:else}
		<p class="form-label">Tipo de usuário não suportado.</p>
	{/if}
</div>

<style>
	.page {
		padding: 40px;
		max-width: 760px;
		margin: 0 auto;
		width: 100%;
	}

	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 16px;
		margin-bottom: 36px;
		flex-wrap: wrap;
	}

	.page-title {
		font-size: 1.75rem;
		font-weight: 400;
		letter-spacing: -1px;
		color: var(--color-text-primary);
		line-height: 1;
		margin-bottom: 6px;
	}

	.page-sub {
		font-size: 0.875rem;
		color: var(--color-text-tertiary);
	}

	.loading-state {
		display: flex;
		align-items: center;
		gap: 10px;
		color: var(--color-text-tertiary);
		font-size: 0.875rem;
	}

	.loading-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid var(--border-frost);
		border-top-color: var(--color-text-secondary);
		border-radius: 50%;
		animation: spin 0.7s linear infinite;
	}

	@keyframes spin { to { transform: rotate(360deg); } }

	/* Form layout */
	.form-layout {
		display: grid;
		grid-template-columns: 1fr 280px;
		gap: 24px;
		align-items: start;
	}

	.form-section {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 12px;
		padding: 24px;
		box-shadow: var(--shadow-ring);
	}

	.form-section-title {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		margin-bottom: 20px;
		padding-bottom: 12px;
		border-bottom: 1px solid var(--border-frost);
	}

	.form-actions {
		margin-top: 8px;
		padding-top: 20px;
		border-top: 1px solid var(--border-frost);
	}

	/* Checkbox */
	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 12px 0;
	}

	.checkbox-row input[type="checkbox"] {
		width: 15px;
		height: 15px;
		flex-shrink: 0;
		cursor: pointer;
	}

	.checkbox-label {
		font-size: 0.875rem;
		color: var(--color-text-secondary);
		cursor: pointer;
	}

	/* Info panel */
	.info-panel {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 12px;
		padding: 24px;
		box-shadow: var(--shadow-ring);
	}

	.score-display {
		display: flex;
		align-items: baseline;
		gap: 4px;
		margin: 12px 0;
	}

	.score-value {
		font-size: 3rem;
		font-weight: 400;
		letter-spacing: -2px;
		color: var(--color-green-10);
		line-height: 1;
	}

	.score-max {
		font-size: 1rem;
		color: var(--color-text-tertiary);
	}

	.score-bar {
		width: 100%;
		height: 3px;
		background: rgba(255, 255, 255, 0.08);
		border-radius: 9999px;
		margin-bottom: 12px;
	}

	.score-fill {
		height: 100%;
		border-radius: 9999px;
		background: var(--color-green-10);
		box-shadow: 0 0 6px rgba(34, 255, 153, 0.4);
		transition: width 0.6s ease;
	}

	.score-note {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		line-height: 1.5;
	}

	@media (max-width: 768px) {
		.page { padding: 24px 16px; }
		.form-layout { grid-template-columns: 1fr; }
	}
</style>
