<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import type { SolicitacaoResponse, StatusSolicitacao } from '$lib/types';
	import { Plus } from 'lucide-svelte';

	type Filtro = StatusSolicitacao | 'TODAS';

	const FILTROS: { label: string; value: Filtro }[] = [
		{ label: 'Todas', value: 'TODAS' },
		{ label: 'Aguardando', value: 'AGUARDANDO' },
		{ label: 'Atribuída', value: 'ATRIBUIDA' },
		{ label: 'Confirmada', value: 'CONFIRMADA' },
		{ label: 'Em andamento', value: 'EM_ANDAMENTO' },
		{ label: 'Concluída', value: 'CONCLUIDA' },
		{ label: 'Cancelada', value: 'CANCELADA' }
	];

	let solicitacoes = $state<SolicitacaoResponse[]>([]);
	let carregando = $state(true);
	let filtroAtivo = $state<Filtro>('TODAS');

	onMount(carregar);

	async function carregar() {
		carregando = true;
		try {
			solicitacoes = await api.listarSolicitacoes();
		} catch {
			toasts.error('Não foi possível carregar suas solicitações.');
		} finally {
			carregando = false;
		}
	}

	let lista = $derived(
		filtroAtivo === 'TODAS'
			? solicitacoes
			: solicitacoes.filter((s) => s.status === filtroAtivo)
	);

	function formatBRL(v: number) {
		return v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
	}

	function formatData(iso: string) {
		return new Date(iso).toLocaleString('pt-BR', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="page">
	<div class="page-header">
		<div>
			<h1>Meus serviços</h1>
			<p>Acompanhe suas solicitações</p>
		</div>
		<a href="/solicitacoes/novo" class="btn-novo">
			<Plus size={14} />
			Nova solicitação
		</a>
	</div>

	<div class="filtros">
		{#each FILTROS as f}
			<button
				type="button"
				class="pill"
				class:active={filtroAtivo === f.value}
				onclick={() => (filtroAtivo = f.value)}
			>{f.label}</button>
		{/each}
	</div>

	{#if carregando}
		<p class="estado">Carregando…</p>
	{:else if lista.length === 0}
		<div class="vazio">
			<p>
				{filtroAtivo === 'TODAS'
					? 'Você ainda não fez solicitações.'
					: `Nenhuma solicitação com status "${FILTROS.find((f) => f.value === filtroAtivo)?.label}".`}
			</p>
			{#if filtroAtivo === 'TODAS'}
				<a href="/solicitacoes/novo" class="btn-novo-vazio">Criar primeira solicitação</a>
			{/if}
		</div>
	{:else}
		<ul class="lista">
			{#each lista as sol (sol.id)}
				<li>
					<a
						class="card"
						href="/dashboard/solicitacoes/{sol.id}"
						aria-label="Ver detalhes da solicitação de {formatData(sol.criada_em)}"
					>
						<div class="card-top">
							<span class="card-data">{formatData(sol.data_servico)}</span>
							<StatusBadge status={sol.status} />
						</div>
						<div class="card-valor">{formatBRL(sol.valor_total)}</div>
						<div class="card-rodape">
							<span>Criada em {formatData(sol.criada_em)}</span>
						</div>
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.page {
		padding: 2rem 2rem 4rem;
		max-width: 720px;
	}

	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 16px;
		margin-bottom: 24px;
	}

	.page-header h1 {
		margin: 0 0 2px;
		font-size: 1.375rem;
		font-weight: 600;
		color: var(--color-text-primary);
		letter-spacing: -0.5px;
	}

	.page-header p {
		margin: 0;
		font-size: 0.875rem;
		color: var(--color-text-tertiary);
	}

	.btn-novo {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 9px 16px;
		background: var(--color-orange-10);
		color: #fff;
		border-radius: var(--radius-pill);
		font-size: 0.8125rem;
		font-weight: 600;
		text-decoration: none;
		white-space: nowrap;
		transition: background 0.12s;
	}
	.btn-novo:hover { background: var(--color-orange-11); }

	.filtros {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
		margin-bottom: 20px;
	}

	.pill {
		padding: 6px 14px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		background: transparent;
		color: var(--color-text-secondary);
		font-size: 0.8125rem;
		cursor: pointer;
		transition: background 0.12s, border-color 0.12s, color 0.12s;
	}
	.pill:hover { background: var(--bg-hover-subtle); color: var(--color-text-primary); }
	.pill.active {
		background: rgba(255, 128, 31, 0.08);
		border-color: var(--color-orange-10);
		color: var(--color-orange-10);
	}

	.estado {
		color: var(--color-text-tertiary);
		font-size: 0.875rem;
	}

	.vazio {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		padding: 48px 24px;
		text-align: center;
		color: var(--color-text-tertiary);
		border: 1px dashed var(--border-frost);
		border-radius: var(--radius-large);
	}

	.btn-novo-vazio {
		padding: 9px 18px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		font-size: 0.875rem;
		color: var(--color-text-primary);
		text-decoration: none;
		transition: background 0.12s;
	}
	.btn-novo-vazio:hover { background: var(--bg-hover-subtle); }

	.lista {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.card {
		padding: 16px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		background: var(--bg-card);
		cursor: pointer;
		transition: border-color 0.12s, background 0.12s;
		display: flex;
		flex-direction: column;
		gap: 8px;
		text-decoration: none;
	}
	.card:hover,
	.card:focus-visible { border-color: var(--border-frost-hover); background: var(--bg-hover-subtle); }

	li { list-style: none; }

	.card-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 8px;
	}

	.card-data {
		font-size: 0.875rem;
		color: var(--color-text-primary);
		font-weight: 500;
	}

	.card-valor {
		font-size: 1.125rem;
		font-weight: 700;
		color: var(--color-orange-10);
		letter-spacing: -0.5px;
	}

	.card-rodape {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
	}
</style>
