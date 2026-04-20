<script lang="ts">
	import type { CategoriaResponse, RegiaoResponse, EnderecoResponse, ItemCalculo } from '$lib/types';
	import { Clock, MapPin, Calendar } from 'lucide-svelte';

	interface Props {
		categoria: CategoriaResponse;
		regiao: RegiaoResponse;
		endereco: EnderecoResponse;
		breakdown: ItemCalculo[];
		valorTotal: number;
		duracaoMin: number;
		dataServico: string;
	}

	let { categoria, regiao, endereco, breakdown, valorTotal, duracaoMin, dataServico }: Props =
		$props();

	function formatBRL(v: number) {
		return v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
	}

	function formatDuracao(min: number) {
		if (min < 60) return `${min} min`;
		const h = Math.floor(min / 60);
		const m = min % 60;
		return m > 0 ? `${h}h ${m}min` : `${h}h`;
	}

	function formatData(iso: string) {
		const d = new Date(iso);
		return d.toLocaleString('pt-BR', {
			weekday: 'long',
			day: '2-digit',
			month: 'long',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function tipoLabel(tipo: string) {
		const m: Record<string, string> = {
			BASE: 'Base',
			COMODO: 'Cômodo',
			OPCIONAL: 'Opcional',
			ACRESCIMO: 'Acréscimo',
			DESCONTO: 'Desconto',
			TOTAL: 'Total'
		};
		return m[tipo] ?? tipo;
	}
</script>

<div class="resumo">
	<div class="resumo-header">
		<span class="resumo-label">Valor combinado</span>
		<strong class="resumo-total">{formatBRL(valorTotal)}</strong>
		<span class="resumo-duracao"><Clock size={12} /> {formatDuracao(duracaoMin)}</span>
	</div>

	<div class="info-grid">
		<div class="info-item">
			<span class="info-icon"><MapPin size={14} /></span>
			<div>
				<span class="info-titulo">Endereço</span>
				<span class="info-detalhe">{endereco.logradouro}, {endereco.numero}{endereco.complemento ? ` — ${endereco.complemento}` : ''}</span>
				<span class="info-detalhe">{endereco.bairro} · {endereco.cidade}/{endereco.estado}</span>
			</div>
		</div>

		<div class="info-item">
			<span class="info-icon"><Calendar size={14} /></span>
			<div>
				<span class="info-titulo">Data e horário</span>
				<span class="info-detalhe">{formatData(dataServico)}</span>
			</div>
		</div>
	</div>

	<div class="breakdown">
		<h4>Como chegamos nesse valor</h4>
		<ul>
			{#each breakdown as item, i (i)}
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

	<p class="disclaimer">
		Categoria: <strong>{categoria.nome}</strong> · Região: <strong>{regiao.nome}</strong>
	</p>
</div>

<style>
	.resumo {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.resumo-header {
		display: flex;
		flex-direction: column;
		gap: 4px;
		padding: 16px;
		background: rgba(255, 128, 31, 0.06);
		border: 1px solid rgba(255, 128, 31, 0.2);
		border-radius: var(--radius-card);
	}

	.resumo-label {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
	}

	.resumo-total {
		font-size: 1.75rem;
		font-weight: 700;
		color: var(--color-orange-10);
		letter-spacing: -1px;
	}

	.resumo-duracao {
		display: flex;
		align-items: center;
		gap: 4px;
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
	}

	.info-grid {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.info-item {
		display: flex;
		gap: 10px;
		align-items: flex-start;
	}

	.info-icon {
		color: var(--color-text-tertiary);
		flex-shrink: 0;
		margin-top: 2px;
	}

	.info-item > div {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.info-titulo {
		font-size: 0.6875rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		color: var(--color-text-tertiary);
	}

	.info-detalhe {
		font-size: 0.875rem;
		color: var(--color-text-primary);
	}

	.breakdown h4 {
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		color: var(--color-text-tertiary);
		margin: 0 0 8px;
	}

	.breakdown ul {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0;
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

	.disclaimer {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		margin: 0;
	}
</style>
