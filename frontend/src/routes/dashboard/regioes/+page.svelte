<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type { RegiaoResponse } from '$lib/types';

	let todasRegioes = $state<RegiaoResponse[]>([]);
	let regioesAtuacao = $state<RegiaoResponse[]>([]);
	let carregando = $state(true);
	let salvando = $state(false);
	let selecionadas = $state<Set<string>>(new Set());

	let modificado = $derived(
		JSON.stringify([...selecionadas].sort()) !==
		JSON.stringify(regioesAtuacao.map(r => r.id).sort())
	);

	onMount(async () => {
		try {
			const [todas, atuacao] = await Promise.all([
				api.listarRegioes(),
				api.listarRegioesAtuacao().catch(() => [])
			]);
			todasRegioes = todas;
			regioesAtuacao = atuacao;
			selecionadas = new Set(atuacao.map((r) => r.id));
		} catch {
			toasts.error('Erro ao carregar regiões');
		} finally {
			carregando = false;
		}
	});

	function toggle(id: string) {
		const novas = new Set(selecionadas);
		if (novas.has(id)) novas.delete(id);
		else novas.add(id);
		selecionadas = novas;
	}

	async function salvar() {
		salvando = true;
		try {
			regioesAtuacao = await api.definirRegioesAtuacao({ regiao_ids: [...selecionadas] });
			selecionadas = new Set(regioesAtuacao.map((r) => r.id));
			toasts.success('Regiões salvas!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar regiões');
		} finally {
			salvando = false;
		}
	}
</script>

<svelte:head>
	<title>Regiões de Atuação — DiaryGo</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<div>
			<h1 class="page-title">Regiões de atuação</h1>
			<p class="page-sub">Selecione as regiões onde você está disponível para atender.</p>
		</div>
		{#if selecionadas.size > 0}
			<span class="badge badge-blue">{selecionadas.size} selecionada{selecionadas.size !== 1 ? 's' : ''}</span>
		{/if}
	</div>

	{#if carregando}
		<div class="loading-state">
			<div class="loading-spinner"></div>
			<span>Carregando regiões...</span>
		</div>
	{:else if todasRegioes.length === 0}
		<div class="empty-state">
			<p class="empty-title">Nenhuma região disponível</p>
			<p class="empty-sub">Nenhuma região está cadastrada no sistema no momento.</p>
		</div>
	{:else}
		<div class="regions-grid">
			{#each todasRegioes as regiao (regiao.id)}
				<button
					class="region-card"
					class:selected={selecionadas.has(regiao.id)}
					onclick={() => toggle(regiao.id)}
					type="button"
				>
					<div class="region-check">
						{#if selecionadas.has(regiao.id)}
							<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
						{/if}
					</div>
					<div class="region-info">
						<span class="region-nome">{regiao.nome}</span>
						<span class="region-sub">{regiao.cidade} / {regiao.estado}</span>
						<span class="region-cep">CEP {regiao.cep_inicio} – {regiao.cep_fim}</span>
					</div>
				</button>
			{/each}
		</div>

		<div class="save-bar" class:visible={modificado}>
			<span class="save-hint">Você tem alterações não salvas.</span>
			<button class="btn btn-white" onclick={salvar} disabled={salvando}>
				{salvando ? 'Salvando...' : 'Salvar regiões'}
			</button>
		</div>
	{/if}
</div>

<style>
	.page { padding: 40px; max-width: 800px; margin: 0 auto; width: 100%; }

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

	.loading-state {
		display: flex; align-items: center; gap: 10px;
		color: var(--color-text-tertiary); font-size: 0.875rem; padding: 40px 0;
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

	.regions-grid {
		display: flex;
		flex-direction: column;
		gap: 6px;
		margin-bottom: 24px;
	}

	.region-card {
		display: flex;
		align-items: flex-start;
		gap: 12px;
		padding: 14px 16px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 10px;
		cursor: pointer;
		text-align: left;
		transition: background 0.12s, border-color 0.12s;
		box-shadow: var(--shadow-ring);
	}

	.region-card:hover {
		background: rgba(255, 255, 255, 0.04);
	}

	.region-card.selected {
		border-color: var(--color-blue-5);
		background: rgba(0, 117, 255, 0.06);
	}

	.region-check {
		width: 18px; height: 18px; border-radius: 4px;
		border: 1.5px solid var(--border-frost);
		display: flex; align-items: center; justify-content: center;
		flex-shrink: 0; margin-top: 1px;
		transition: background 0.12s, border-color 0.12s;
	}

	.region-card.selected .region-check {
		background: var(--color-blue-10);
		border-color: var(--color-blue-10);
		color: #000;
	}

	.region-info {
		display: flex; flex-direction: column; gap: 3px;
	}

	.region-nome {
		font-size: 0.9375rem; font-weight: 500; color: var(--color-text-primary);
	}

	.region-sub {
		font-size: 0.8125rem; color: var(--color-text-secondary);
	}

	.region-cep {
		font-size: 0.75rem; color: var(--color-text-tertiary);
		font-family: var(--font-mono);
	}

	/* Barra de salvar fixa */
	.save-bar {
		position: sticky;
		bottom: 24px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		padding: 12px 16px;
		background: rgba(10, 10, 10, 0.9);
		border: 1px solid var(--border-frost);
		border-radius: 10px;
		backdrop-filter: blur(8px);
		opacity: 0;
		transform: translateY(8px);
		pointer-events: none;
		transition: opacity 0.2s, transform 0.2s;
	}

	.save-bar.visible {
		opacity: 1;
		transform: translateY(0);
		pointer-events: auto;
	}

	.save-hint {
		font-size: 0.875rem;
		color: var(--color-text-secondary);
	}

	@media (max-width: 640px) {
		.page { padding: 24px 16px; }
	}
</style>
