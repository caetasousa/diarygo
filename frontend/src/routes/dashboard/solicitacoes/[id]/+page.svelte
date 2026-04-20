<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import ModalConfirmacao from '$lib/components/ModalConfirmacao.svelte';
	import type { SolicitacaoResponse, StatusSolicitacao } from '$lib/types';
	import { ArrowLeft, Clock, MapPin, Calendar } from 'lucide-svelte';

	const STATUS_TIMELINE: StatusSolicitacao[] = [
		'AGUARDANDO',
		'ATRIBUIDA',
		'CONFIRMADA',
		'EM_ANDAMENTO',
		'CONCLUIDA'
	];

	const STATUS_LABELS: Record<StatusSolicitacao, string> = {
		AGUARDANDO: 'Aguardando',
		ATRIBUIDA: 'Atribuída',
		CONFIRMADA: 'Confirmada',
		EM_ANDAMENTO: 'Em andamento',
		CONCLUIDA: 'Concluída',
		CANCELADA: 'Cancelada'
	};

	const CANCELAVEIS: StatusSolicitacao[] = ['AGUARDANDO', 'ATRIBUIDA', 'CONFIRMADA'];

	let sol = $state<SolicitacaoResponse | null>(null);
	let carregando = $state(true);
	let modalAberto = $state(false);
	let cancelando = $state(false);

	let id = $derived(page.params.id ?? '');

	let podeCancelar = $derived(sol ? CANCELAVEIS.includes(sol.status) : false);

	let horasAteServico = $derived(() => {
		if (!sol) return Infinity;
		return (new Date(sol.data_servico).getTime() - Date.now()) / 3600000;
	});

	let statusIdx = $derived(sol ? STATUS_TIMELINE.indexOf(sol.status) : -1);

	onMount(async () => {
		if (!id) { goto('/dashboard/solicitacoes'); return; }
		try {
			sol = await api.buscarSolicitacao(id);
		} catch {
			toasts.error('Solicitação não encontrada.');
			goto('/dashboard/solicitacoes');
		} finally {
			carregando = false;
		}
	});

	async function cancelar() {
		modalAberto = false;
		cancelando = true;
		try {
			await api.cancelarSolicitacao(id);
			sol = await api.buscarSolicitacao(id);
			toasts.error('Solicitação cancelada.');
		} catch (e: unknown) {
			toasts.error(e instanceof Error ? e.message : 'Não foi possível cancelar.');
		} finally {
			cancelando = false;
		}
	}

	function formatBRL(v: number) {
		return v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
	}

	function formatData(iso: string) {
		return new Date(iso).toLocaleString('pt-BR', {
			weekday: 'long',
			day: '2-digit',
			month: 'long',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatDuracao(min: number) {
		const h = Math.floor(min / 60);
		const m = min % 60;
		if (h === 0) return `${m} min`;
		return m > 0 ? `${h}h ${m}min` : `${h}h`;
	}

	function tipoLabel(tipo: string) {
		const m: Record<string, string> = {
			BASE: 'Base', COMODO: 'Cômodo', OPCIONAL: 'Opcional',
			ACRESCIMO: 'Acréscimo', DESCONTO: 'Desconto', TOTAL: 'Total'
		};
		return m[tipo] ?? tipo;
	}
</script>

<div class="page">
	<div class="page-header">
		<a href="/dashboard/solicitacoes" class="btn-voltar">
			<ArrowLeft size={14} />
			Voltar
		</a>
		{#if sol}
			<StatusBadge status={sol.status} />
		{/if}
	</div>

	{#if carregando}
		<p class="estado">Carregando…</p>
	{:else if sol}
		<!-- ══════ Resumo do valor ══════ -->
		<div class="valor-card">
			<span class="valor-label">Valor combinado</span>
			<strong class="valor-total">{formatBRL(sol.valor_total)}</strong>
			<span class="valor-duracao"><Clock size={12} /> {formatDuracao(sol.duracao_min)}</span>
		</div>

		<!-- ══════ Info principal ══════ -->
		<div class="info-grid">
			<div class="info-item">
				<span class="info-icon"><Calendar size={14} /></span>
				<div>
					<span class="info-titulo">Data do serviço</span>
					<span class="info-detalhe">{formatData(sol.data_servico)}</span>
				</div>
			</div>
			<div class="info-item">
				<span class="info-icon"><MapPin size={14} /></span>
				<div>
					<span class="info-titulo">Endereço</span>
					<span class="info-detalhe">ID: {sol.endereco_id}</span>
				</div>
			</div>
		</div>

		<!-- ══════ Timeline de status ══════ -->
		{#if sol.status !== 'CANCELADA'}
			<div class="timeline">
				<h3>Progresso</h3>
				<ol class="tl-list">
					{#each STATUS_TIMELINE as s, i}
						<li
							class="tl-item"
							class:done={i < statusIdx}
							class:active={i === statusIdx}
						>
							<div class="tl-dot"></div>
							{#if i < STATUS_TIMELINE.length - 1}
								<div class="tl-line" class:done={i < statusIdx}></div>
							{/if}
							<span class="tl-label">{STATUS_LABELS[s]}</span>
						</li>
					{/each}
				</ol>
			</div>
		{:else}
			<div class="cancelada-banner">
				Solicitação cancelada
				{#if sol.cancelada_em}
					em {formatData(sol.cancelada_em)}
				{/if}
			</div>
		{/if}

		<!-- ══════ Breakdown ══════ -->
		{#if sol.breakdown?.length > 0}
			<div class="breakdown">
				<h3>Como chegamos nesse valor</h3>
				<ul>
					{#each sol.breakdown as item, i (i)}
						<li
							class="bk-item"
							class:total={item.tipo === 'TOTAL'}
							class:desconto={item.tipo === 'DESCONTO'}
							class:acrescimo={item.tipo === 'ACRESCIMO'}
						>
							<span class="bk-tipo">{tipoLabel(item.tipo)}</span>
							<span class="bk-label">{item.label}</span>
							<span class="bk-valor">{formatBRL(item.valor)}</span>
						</li>
					{/each}
				</ul>
			</div>
		{/if}

		<!-- ══════ Observação ══════ -->
		{#if sol.observacao}
			<div class="obs-card">
				<span class="obs-label">Observações</span>
				<p class="obs-texto">{sol.observacao}</p>
			</div>
		{/if}

		<!-- ══════ Cancelar ══════ -->
		{#if podeCancelar}
			<div class="acoes">
				<button
					class="btn-cancelar"
					disabled={cancelando}
					onclick={() => (modalAberto = true)}
				>
					{cancelando ? 'Cancelando…' : 'Cancelar solicitação'}
				</button>
			</div>
		{/if}
	{/if}
</div>

<ModalConfirmacao
	bind:open={modalAberto}
	titulo="Cancelar solicitação"
	tipo="danger"
	mensagem={horasAteServico() < 24
		? 'Atenção: este cancelamento penalizará seu score de confiabilidade em 5 pontos pois faltam menos de 24h para o serviço.'
		: 'Tem certeza que deseja cancelar esta solicitação?'}
	textoConfirmar="Confirmar cancelamento"
	onConfirm={cancelar}
	onCancel={() => (modalAberto = false)}
/>

<style>
	.page {
		padding: 2rem 2rem 4rem;
		max-width: 640px;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.page-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
	}

	.btn-voltar {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		font-size: 0.875rem;
		color: var(--color-text-secondary);
		text-decoration: none;
		transition: color 0.12s;
	}
	.btn-voltar:hover { color: var(--color-text-primary); }

	.estado {
		color: var(--color-text-tertiary);
		font-size: 0.875rem;
	}

	/* Valor */
	.valor-card {
		display: flex;
		flex-direction: column;
		gap: 4px;
		padding: 16px;
		background: rgba(255, 128, 31, 0.06);
		border: 1px solid rgba(255, 128, 31, 0.2);
		border-radius: var(--radius-card);
	}

	.valor-label {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
	}

	.valor-total {
		font-size: 1.75rem;
		font-weight: 700;
		color: var(--color-orange-10);
		letter-spacing: -1px;
	}

	.valor-duracao {
		display: flex;
		align-items: center;
		gap: 4px;
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
	}

	/* Info grid */
	.info-grid { display: flex; flex-direction: column; gap: 12px; }

	.info-item { display: flex; gap: 10px; align-items: flex-start; }

	.info-icon { color: var(--color-text-tertiary); flex-shrink: 0; margin-top: 2px; }

	.info-item > div { display: flex; flex-direction: column; gap: 2px; }

	.info-titulo {
		font-size: 0.6875rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		color: var(--color-text-tertiary);
	}

	.info-detalhe { font-size: 0.875rem; color: var(--color-text-primary); }

	/* Timeline */
	.timeline h3,
	.breakdown h3 {
		margin: 0 0 12px;
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		color: var(--color-text-tertiary);
	}

	.tl-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: row;
		align-items: flex-start;
		gap: 0;
	}

	.tl-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		flex: 1;
		position: relative;
	}

	.tl-dot {
		width: 12px;
		height: 12px;
		border-radius: 50%;
		border: 2px solid var(--border-frost);
		background: var(--bg-card);
		z-index: 1;
		transition: background 0.15s, border-color 0.15s;
	}

	.tl-item.done .tl-dot {
		background: var(--color-green-10);
		border-color: var(--color-green-10);
	}

	.tl-item.active .tl-dot {
		background: var(--color-orange-10);
		border-color: var(--color-orange-10);
	}

	.tl-line {
		position: absolute;
		top: 5px;
		left: 50%;
		width: 100%;
		height: 2px;
		background: var(--border-frost);
		transition: background 0.15s;
	}

	.tl-line.done { background: var(--color-green-4); }

	.tl-label {
		margin-top: 6px;
		font-size: 0.625rem;
		text-align: center;
		color: var(--color-text-tertiary);
		white-space: nowrap;
	}

	.tl-item.active .tl-label { color: var(--color-text-primary); font-weight: 500; }
	.tl-item.done .tl-label { color: var(--color-text-secondary); }

	/* Cancelada */
	.cancelada-banner {
		padding: 12px 16px;
		background: rgba(239, 68, 68, 0.08);
		border: 1px solid rgba(239, 68, 68, 0.25);
		border-radius: var(--radius-card);
		font-size: 0.875rem;
		color: #ef4444;
		font-weight: 500;
	}

	/* Breakdown */
	.breakdown ul {
		list-style: none;
		padding: 0;
		margin: 0;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		overflow: hidden;
	}

	.bk-item {
		display: grid;
		grid-template-columns: 70px 1fr auto;
		gap: 8px;
		align-items: center;
		padding: 10px 14px;
		border-bottom: 1px solid var(--border-frost);
		font-size: 0.8125rem;
	}
	.bk-item:last-child { border-bottom: none; }
	.bk-item.total { background: rgba(255, 255, 255, 0.03); font-weight: 600; }
	.bk-item.desconto .bk-valor { color: #22c55e; }
	.bk-item.acrescimo .bk-valor { color: var(--color-orange-10); }

	.bk-tipo {
		font-size: 0.625rem;
		text-transform: uppercase;
		letter-spacing: 0.4px;
		font-weight: 600;
		color: var(--color-text-tertiary);
	}

	.bk-label { color: var(--color-text-primary); }
	.bk-valor { color: var(--color-text-primary); text-align: right; font-variant-numeric: tabular-nums; }

	/* Observação */
	.obs-card {
		padding: 14px 16px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.obs-label {
		font-size: 0.6875rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		color: var(--color-text-tertiary);
	}

	.obs-texto {
		margin: 0;
		font-size: 0.875rem;
		color: var(--color-text-primary);
		line-height: 1.5;
	}

	/* Ações */
	.acoes { display: flex; justify-content: flex-end; padding-top: 4px; }

	.btn-cancelar {
		padding: 9px 18px;
		border: 1px solid rgba(239, 68, 68, 0.4);
		border-radius: var(--radius-pill);
		background: transparent;
		color: #ef4444;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.12s, border-color 0.12s;
	}
	.btn-cancelar:hover:not(:disabled) {
		background: rgba(239, 68, 68, 0.08);
		border-color: #ef4444;
	}
	.btn-cancelar:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
