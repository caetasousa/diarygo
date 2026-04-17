<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import TelefoneInput from '$lib/components/TelefoneInput.svelte';
	import type { ReferenciaResponse } from '$lib/types';

	let referencias = $state<ReferenciaResponse[]>([]);
	let carregando = $state(true);
	let adicionando = $state(false);

	let nomeContato = $state('');
	let telefoneContato = $state('');
	let erros = $state<Record<string, string>>({});

	const statusLabel: Record<string, string> = {
		PENDENTE: 'Pendente', CONFIRMADA: 'Confirmada', NAO_CONFIRMADA: 'Não confirmada'
	};
	const statusBadge: Record<string, string> = {
		PENDENTE: 'badge-orange', CONFIRMADA: 'badge-green', NAO_CONFIRMADA: 'badge-red'
	};

	onMount(carregarReferencias);

	async function carregarReferencias() {
		carregando = true;
		try {
			referencias = await api.listarReferencias();
		} catch {
			toasts.error('Erro ao carregar referências');
		} finally {
			carregando = false;
		}
	}

	function validar(): boolean {
		erros = {};
		if (!nomeContato.trim()) erros.nome = 'Nome é obrigatório';
		if (telefoneContato.replace(/\D/g, '').length < 10) erros.telefone = 'Telefone inválido';
		return Object.keys(erros).length === 0;
	}

	async function adicionar() {
		if (!validar()) return;
		adicionando = true;
		try {
			await api.adicionarReferencia({
				nome_contato: nomeContato.trim(),
				telefone_contato: telefoneContato.replace(/\D/g, '')
			});
			nomeContato = '';
			telefoneContato = '';
			erros = {};
			toasts.success('Referência adicionada!');
			await carregarReferencias();
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao adicionar referência');
		} finally {
			adicionando = false;
		}
	}
</script>

<svelte:head>
	<title>Referências — DiaryGo</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<div>
			<h1 class="page-title">Referências</h1>
			<p class="page-sub">Contatos que atestam sua experiência profissional.</p>
		</div>
		{#if referencias.length > 0}
			<span class="badge badge-{referencias.length >= 2 ? 'green' : 'orange'}">
				{referencias.length} / 2 mínimo
			</span>
		{/if}
	</div>

	<div class="layout">
		<div class="form-card">
			<h2 class="card-title">Adicionar referência</h2>
			<form novalidate onsubmit={(e) => { e.preventDefault(); adicionar(); }}>
				<div class="form-group">
					<label for="nome-contato" class="form-label">Nome do contato</label>
					<input
						id="nome-contato" type="text" maxlength="100"
						placeholder="Nome completo" class="form-input"
						class:input-error={!!erros.nome} bind:value={nomeContato}
						autocomplete="off"
					/>
					{#if erros.nome}<p class="form-error">{erros.nome}</p>{/if}
				</div>

				<TelefoneInput bind:value={telefoneContato} error={erros.telefone} />

				<button type="submit" class="btn btn-white btn-full" disabled={adicionando}>
					{adicionando ? 'Adicionando...' : 'Adicionar referência'}
				</button>
			</form>
		</div>

		<div class="list-col">
			{#if carregando}
				<div class="loading-state">
					<div class="loading-spinner"></div>
					<span>Carregando...</span>
				</div>
			{:else if referencias.length === 0}
				<div class="empty-state">
					<p class="empty-title">Nenhuma referência adicionada</p>
					<p class="empty-sub">Adicione ao menos 2 referências profissionais para concluir o cadastro.</p>
				</div>
			{:else}
				<div class="ref-list">
					{#each referencias as ref (ref.id)}
						<div class="ref-row">
							<div class="ref-avatar">
								{ref.nome_contato[0].toUpperCase()}
							</div>
							<div class="ref-info">
								<span class="ref-nome">{ref.nome_contato}</span>
								<span class="ref-tel">{ref.telefone_contato}</span>
							</div>
							<span class="badge {statusBadge[ref.status] ?? 'badge-blue'}">
								{statusLabel[ref.status] ?? ref.status}
							</span>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.page { padding: 40px; max-width: 900px; margin: 0 auto; width: 100%; }

	.page-header {
		display: flex; align-items: flex-start;
		justify-content: space-between; gap: 16px;
		margin-bottom: 32px; flex-wrap: wrap;
	}

	.page-title {
		font-size: 1.75rem; font-weight: 400;
		letter-spacing: -1px; color: var(--color-text-primary);
		line-height: 1; margin-bottom: 6px;
	}

	.page-sub { font-size: 0.875rem; color: var(--color-text-tertiary); }

	.layout {
		display: grid;
		grid-template-columns: 300px 1fr;
		gap: 20px;
		align-items: start;
	}

	.form-card {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 12px;
		padding: 20px;
		box-shadow: var(--shadow-ring);
	}

	.card-title {
		font-size: 0.8125rem; font-weight: 500;
		color: var(--color-text-tertiary);
		text-transform: uppercase; letter-spacing: 0.5px;
		margin-bottom: 16px; padding-bottom: 12px;
		border-bottom: 1px solid var(--border-frost);
	}

	.list-col { min-width: 0; }

	.loading-state {
		display: flex; align-items: center; gap: 10px;
		color: var(--color-text-tertiary); font-size: 0.875rem; padding: 32px 0;
	}

	.loading-spinner {
		width: 16px; height: 16px;
		border: 2px solid var(--border-frost);
		border-top-color: var(--color-text-secondary);
		border-radius: 50%;
		animation: spin 0.7s linear infinite;
	}

	@keyframes spin { to { transform: rotate(360deg); } }

	.empty-state {
		padding: 40px 24px; text-align: center;
		background: var(--bg-card); border: 1px solid var(--border-frost);
		border-radius: 12px; box-shadow: var(--shadow-ring);
	}
	.empty-title { font-size: 0.9375rem; font-weight: 500; color: var(--color-text-primary); margin-bottom: 6px; }
	.empty-sub { font-size: 0.8125rem; color: var(--color-text-tertiary); }

	.ref-list {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 12px;
		overflow: hidden;
		box-shadow: var(--shadow-ring);
	}

	.ref-row {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 14px 18px;
		border-bottom: 1px solid var(--border-frost);
	}

	.ref-row:last-child { border-bottom: none; }

	.ref-avatar {
		width: 34px; height: 34px; border-radius: 50%;
		background: var(--color-blue-4);
		border: 1px solid var(--color-blue-5);
		display: flex; align-items: center; justify-content: center;
		font-size: 0.8125rem; font-weight: 600;
		color: var(--color-blue-10); flex-shrink: 0;
	}

	.ref-info {
		display: flex; flex-direction: column; gap: 2px; flex: 1; min-width: 0;
	}

	.ref-nome {
		font-size: 0.875rem; font-weight: 500; color: var(--color-text-primary);
	}

	.ref-tel {
		font-size: 0.75rem; color: var(--color-text-tertiary);
	}

	@media (max-width: 768px) {
		.page { padding: 24px 16px; }
		.layout { grid-template-columns: 1fr; }
	}
</style>
