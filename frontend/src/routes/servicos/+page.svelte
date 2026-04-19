<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type {
		CategoriaResponse,
		OpcionalResponse,
		RegiaoResponse,
		CalculoPrecoResponse,
		FrequenciaServico
	} from '$lib/types';
	import {
		Sparkles,
		MapPin,
		Calendar,
		RefreshCw,
		Plus,
		Minus,
		Info,
		Loader2
	} from 'lucide-svelte';

	let categorias = $state<CategoriaResponse[]>([]);
	let regioes = $state<RegiaoResponse[]>([]);
	let opcionais = $state<OpcionalResponse[]>([]);
	let carregando = $state(true);
	let calculando = $state(false);
	let erroCalculo = $state<string | null>(null);

	let categoriaID = $state('');
	let regiaoID = $state('');
	let numQuartos = $state(1);
	let numBanheiros = $state(1);
	let numSalas = $state(1);
	let numCozinhas = $state(1);
	let opcionaisSelecionados = $state<Set<string>>(new Set());
	let frequencia = $state<FrequenciaServico>('UNICA');
	let dataServico = $state('');

	let resultado = $state<CalculoPrecoResponse | null>(null);

	onMount(async () => {
		try {
			const [cats, regs] = await Promise.all([api.listarCategorias(), api.listarRegioes()]);
			categorias = cats;
			regioes = regs;
			if (cats.length > 0) categoriaID = cats[0].id;
			if (regs.length > 0) regiaoID = regs[0].id;
		} catch {
			toasts.error('Não foi possível carregar o catálogo. Tente novamente em instantes.');
		} finally {
			carregando = false;
		}
	});

	$effect(() => {
		if (!categoriaID) return;
		api
			.listarOpcionais(categoriaID)
			.then((ops) => {
				opcionais = ops;
				opcionaisSelecionados = new Set();
			})
			.catch(() => {
				opcionais = [];
			});
	});

	function toggleOpcional(id: string) {
		const novos = new Set(opcionaisSelecionados);
		if (novos.has(id)) novos.delete(id);
		else novos.add(id);
		opcionaisSelecionados = novos;
	}

	function inc(campo: 'num_quartos' | 'num_banheiros' | 'num_salas' | 'num_cozinhas', delta: number) {
		if (campo === 'num_quartos') numQuartos = Math.max(1, Math.min(20, numQuartos + delta));
		if (campo === 'num_banheiros') numBanheiros = Math.max(0, Math.min(20, numBanheiros + delta));
		if (campo === 'num_salas') numSalas = Math.max(0, Math.min(20, numSalas + delta));
		if (campo === 'num_cozinhas') numCozinhas = Math.max(0, Math.min(20, numCozinhas + delta));
	}

	async function calcular(e: Event) {
		e.preventDefault();
		if (!categoriaID || !regiaoID) {
			erroCalculo = 'Selecione categoria e região.';
			return;
		}
		calculando = true;
		erroCalculo = null;
		try {
			resultado = await api.calcularPreco({
				categoria_id: categoriaID,
				regiao_id: regiaoID,
				num_quartos: numQuartos,
				num_banheiros: numBanheiros,
				num_salas: numSalas,
				num_cozinhas: numCozinhas,
				opcionais_ids: Array.from(opcionaisSelecionados),
				frequencia,
				data_servico: dataServico ? new Date(dataServico).toISOString() : undefined
			});
		} catch (err) {
			const msg = err instanceof Error ? err.message : 'Erro ao calcular orçamento';
			erroCalculo = msg;
			resultado = null;
		} finally {
			calculando = false;
		}
	}

	function formatBRL(v: number) {
		return v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
	}

	function formatDuracao(min: number) {
		const h = Math.floor(min / 60);
		const m = min % 60;
		if (h === 0) return `${m} min`;
		if (m === 0) return `${h}h`;
		return `${h}h ${m}min`;
	}

	function tipoLabel(tipo: string) {
		switch (tipo) {
			case 'BASE':
				return 'Base';
			case 'COMODO':
				return 'Cômodos';
			case 'OPCIONAL':
				return 'Opcional';
			case 'ACRESCIMO':
				return 'Acréscimo';
			case 'DESCONTO':
				return 'Desconto';
			case 'TOTAL':
				return 'Total';
			default:
				return tipo;
		}
	}

	// Data mínima = amanhã (agendamento min 24h antecedência)
	let minDate = $derived.by(() => {
		const d = new Date();
		d.setDate(d.getDate() + 1);
		return d.toISOString().slice(0, 10);
	});
</script>

<svelte:head>
	<title>Calcule seu orçamento — DiaryGo</title>
</svelte:head>

<section class="hero">
	<div class="container">
		<div class="badge"><Sparkles size={14} /> Orçamento transparente</div>
		<h1>Veja exatamente quanto custa sua faxina</h1>
		<p>Calcule em segundos, sem cadastro. Cada centavo é explicado no detalhamento.</p>
	</div>
</section>

<section class="calc-wrap">
	<div class="container grid">
		<form class="panel form" onsubmit={calcular} novalidate>
			<h2>Parâmetros</h2>

			{#if carregando}
				<div class="loading"><Loader2 size={20} class="spin" /> Carregando catálogo…</div>
			{:else}
				<div class="field">
					<label for="categoria"><Sparkles size={14} /> Categoria</label>
					<select id="categoria" bind:value={categoriaID} required>
						{#each categorias as c (c.id)}
							<option value={c.id}>{c.nome}</option>
						{/each}
					</select>
				</div>

				<div class="field">
					<label for="regiao"><MapPin size={14} /> Região</label>
					<select id="regiao" bind:value={regiaoID} required>
						{#each regioes as r (r.id)}
							<option value={r.id}>{r.nome}</option>
						{/each}
					</select>
				</div>

				<fieldset class="comodos">
					<legend>Cômodos</legend>
					{#each [{ label: 'Quartos', campo: 'num_quartos', valor: numQuartos, min: 1 }, { label: 'Banheiros', campo: 'num_banheiros', valor: numBanheiros, min: 0 }, { label: 'Salas', campo: 'num_salas', valor: numSalas, min: 0 }, { label: 'Cozinhas', campo: 'num_cozinhas', valor: numCozinhas, min: 0 }] as row (row.campo)}
						<div class="stepper">
							<span class="stepper-label">{row.label}</span>
							<div class="stepper-btns">
								<button
									type="button"
									class="step-btn"
									disabled={row.valor <= row.min}
									onclick={() => inc(row.campo as 'num_quartos' | 'num_banheiros' | 'num_salas' | 'num_cozinhas', -1)}
									aria-label="Diminuir {row.label}"
								>
									<Minus size={14} />
								</button>
								<span class="step-valor">{row.valor}</span>
								<button
									type="button"
									class="step-btn"
									onclick={() => inc(row.campo as 'num_quartos' | 'num_banheiros' | 'num_salas' | 'num_cozinhas', 1)}
									aria-label="Aumentar {row.label}"
								>
									<Plus size={14} />
								</button>
							</div>
						</div>
					{/each}
				</fieldset>

				{#if opcionais.length > 0}
					<fieldset class="opcionais">
						<legend>Opcionais</legend>
						{#each opcionais as op (op.id)}
							<label class="opcional">
								<input
									type="checkbox"
									checked={opcionaisSelecionados.has(op.id)}
									onchange={() => toggleOpcional(op.id)}
								/>
								<span class="opcional-info">
									<strong>{op.nome}</strong>
									<small>+{formatBRL(op.valor_extra)} · +{op.tempo_extra_min} min</small>
								</span>
							</label>
						{/each}
					</fieldset>
				{/if}

				<div class="field">
					<label for="freq"><RefreshCw size={14} /> Frequência</label>
					<select id="freq" bind:value={frequencia}>
						<option value="UNICA">Única (sem desconto)</option>
						<option value="SEMANAL">Semanal (−10%)</option>
						<option value="QUINZENAL">Quinzenal (−5%)</option>
						<option value="DUAS_POR_SEMANA">2× por semana (−15%)</option>
					</select>
				</div>

				<div class="field">
					<label for="data"><Calendar size={14} /> Data do serviço <small>(opcional)</small></label>
					<input id="data" type="date" bind:value={dataServico} min={minDate} maxlength="10" />
				</div>

				<button type="submit" class="btn-primary" disabled={calculando}>
					{#if calculando}
						<Loader2 size={16} class="spin" /> Calculando…
					{:else}
						Calcular orçamento
					{/if}
				</button>

				{#if erroCalculo}
					<p class="erro" role="alert">{erroCalculo}</p>
				{/if}
			{/if}
		</form>

		<aside class="panel resultado" aria-live="polite">
			{#if !resultado}
				<div class="resultado-vazio">
					<Info size={28} />
					<p>Preencha os parâmetros ao lado e veja o orçamento completo aqui.</p>
				</div>
			{:else}
				<header class="resultado-header">
					<span class="resultado-label">Orçamento estimado</span>
					<strong class="resultado-total">{formatBRL(resultado.valor_total)}</strong>
					<span class="resultado-duracao">Duração: {formatDuracao(resultado.duracao_min)}</span>
				</header>

				<div class="breakdown">
					<h3>Como chegamos a esse valor</h3>
					<ul>
						{#each resultado.itens as item, i (i)}
							<li class="breakdown-item" class:total={item.tipo === 'TOTAL'} class:desconto={item.tipo === 'DESCONTO'} class:acrescimo={item.tipo === 'ACRESCIMO'}>
								<span class="bk-tipo">{tipoLabel(item.tipo)}</span>
								<span class="bk-label">{item.label}</span>
								<span class="bk-valor">{formatBRL(item.valor)}</span>
							</li>
						{/each}
					</ul>
				</div>

				<p class="disclaimer">
					<Info size={12} /> Valor de referência. O preço final pode ser ajustado conforme
					avaliação da profissional no dia.
				</p>
			{/if}
		</aside>
	</div>
</section>

<style>
	.hero {
		position: relative;
		padding: 5rem 0 3rem;
		background:
			radial-gradient(80% 60% at 50% 0%, rgba(255, 128, 31, 0.08), transparent 70%),
			radial-gradient(60% 50% at 50% 0%, rgba(255, 255, 255, 0.015), transparent 70%),
			var(--color-black);
		border-bottom: 1px solid var(--border-frost);
		overflow: hidden;
	}

	.hero::before {
		content: '';
		position: absolute;
		inset: 0;
		background-image:
			linear-gradient(rgba(214, 235, 253, 0.04) 1px, transparent 1px),
			linear-gradient(90deg, rgba(214, 235, 253, 0.04) 1px, transparent 1px);
		background-size: 48px 48px;
		mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
		-webkit-mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
		pointer-events: none;
	}

	.container {
		max-width: 1100px;
		margin: 0 auto;
		padding: 0 1.5rem;
		position: relative;
	}

	.badge {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 5px 11px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		background: rgba(255, 128, 31, 0.06);
		color: var(--color-orange-11);
		font-size: 0.75rem;
		font-weight: 500;
		letter-spacing: -0.1px;
		margin-bottom: 1.25rem;
	}

	.hero h1 {
		font-size: clamp(2rem, 4.5vw, 3rem);
		line-height: 1.08;
		letter-spacing: -1px;
		margin: 0 0 0.75rem;
		color: var(--color-text-primary);
		font-weight: 600;
	}

	.hero p {
		color: var(--color-text-secondary);
		margin: 0;
		max-width: 560px;
		font-size: 1.0625rem;
		line-height: 1.55;
		letter-spacing: -0.1px;
	}

	.calc-wrap {
		padding: 3rem 0 5rem;
		background: var(--color-black);
	}

	.grid {
		display: grid;
		gap: 1.25rem;
		grid-template-columns: 1fr;
	}

	@media (min-width: 900px) {
		.grid {
			grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
			align-items: start;
		}
	}

	.panel {
		background: rgba(255, 255, 255, 0.018);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-large);
		padding: 1.75rem;
	}

	.panel h2 {
		margin: 0 0 1.25rem;
		font-size: 0.9375rem;
		font-weight: 600;
		color: var(--color-text-primary);
		letter-spacing: -0.2px;
	}

	.loading {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: var(--color-text-secondary);
		padding: 1rem 0;
		font-size: 0.875rem;
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
		margin-bottom: 1.1rem;
	}

	.field label {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.78rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		letter-spacing: -0.1px;
	}

	.field small {
		color: var(--color-text-tertiary);
		font-weight: 400;
	}

	.field select,
	.field input {
		padding: 0.65rem 0.85rem;
		border-radius: var(--radius-standard);
		border: 1px solid var(--border-frost);
		background: rgba(255, 255, 255, 0.02);
		color: var(--color-text-primary);
		font-family: var(--font-body);
		font-size: 0.9375rem;
		transition: border-color 0.15s, background 0.15s;
	}

	.field select:hover,
	.field input:hover {
		border-color: var(--border-frost-alt);
		background: rgba(255, 255, 255, 0.035);
	}

	.field select:focus,
	.field input:focus {
		outline: none;
		border-color: var(--color-orange-10);
		background: rgba(255, 128, 31, 0.04);
	}

	.field select option {
		background: #0a0a0a;
		color: var(--color-text-primary);
	}

	.comodos,
	.opcionais {
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		padding: 0.85rem 1rem 1rem;
		margin-bottom: 1.1rem;
		background: rgba(255, 255, 255, 0.012);
	}

	.comodos legend,
	.opcionais legend {
		font-size: 0.72rem;
		font-weight: 600;
		color: var(--color-text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.08em;
		padding: 0 0.5rem;
	}

	.stepper {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.45rem 0;
	}

	.stepper-label {
		font-size: 0.9rem;
		color: var(--color-text-primary);
		letter-spacing: -0.1px;
	}

	.stepper-btns {
		display: inline-flex;
		align-items: center;
		gap: 0.6rem;
	}

	.step-btn {
		width: 28px;
		height: 28px;
		border-radius: var(--radius-standard);
		border: 1px solid var(--border-frost);
		background: rgba(255, 255, 255, 0.02);
		display: inline-flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		color: var(--color-text-secondary);
		transition: border-color 0.15s, color 0.15s, background 0.15s;
	}

	.step-btn:hover:not(:disabled) {
		border-color: var(--color-orange-10);
		color: var(--color-orange-10);
		background: rgba(255, 128, 31, 0.08);
	}

	.step-btn:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.step-valor {
		min-width: 1.5rem;
		text-align: center;
		font-weight: 600;
		font-variant-numeric: tabular-nums;
		color: var(--color-text-primary);
	}

	.opcional {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		padding: 0.5rem 0;
		cursor: pointer;
	}

	.opcional input {
		width: 16px;
		height: 16px;
		accent-color: var(--color-orange-10);
	}

	.opcional-info {
		display: flex;
		flex-direction: column;
	}

	.opcional-info strong {
		font-size: 0.875rem;
		color: var(--color-text-primary);
		font-weight: 500;
	}

	.opcional-info small {
		color: var(--color-text-tertiary);
		font-size: 0.75rem;
		margin-top: 2px;
	}

	.btn-primary {
		width: 100%;
		padding: 0.8rem 1rem;
		border-radius: var(--radius-standard);
		border: 0;
		background: var(--color-white);
		color: var(--color-black);
		font-family: var(--font-body);
		font-weight: 600;
		font-size: 0.9375rem;
		letter-spacing: -0.1px;
		cursor: pointer;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		transition: opacity 0.15s, transform 0.15s;
	}

	.btn-primary:hover:not(:disabled) {
		opacity: 0.92;
		transform: translateY(-1px);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.erro {
		color: var(--color-red-10);
		background: rgba(255, 32, 71, 0.08);
		border: 1px solid var(--color-red-5);
		border-radius: var(--radius-standard);
		padding: 0.65rem 0.85rem;
		font-size: 0.8125rem;
		margin: 0.85rem 0 0;
		letter-spacing: -0.1px;
	}

	.resultado {
		position: sticky;
		top: 1.5rem;
	}

	.resultado-vazio {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.85rem;
		color: var(--color-text-tertiary);
		text-align: center;
		padding: 2.5rem 1rem;
	}

	.resultado-vazio p {
		font-size: 0.875rem;
		max-width: 280px;
		margin: 0;
		line-height: 1.5;
	}

	.resultado-header {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		padding-bottom: 1.25rem;
		border-bottom: 1px solid var(--border-frost);
		margin-bottom: 1.25rem;
	}

	.resultado-label {
		font-size: 0.7rem;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		font-weight: 500;
	}

	.resultado-total {
		font-size: 2.25rem;
		background: linear-gradient(135deg, var(--color-orange-11) 0%, var(--color-orange-10) 50%, #ff5900 100%);
		-webkit-background-clip: text;
		background-clip: text;
		color: transparent;
		font-weight: 700;
		letter-spacing: -1.2px;
		line-height: 1;
	}

	.resultado-duracao {
		font-size: 0.8125rem;
		color: var(--color-text-secondary);
		letter-spacing: -0.1px;
	}

	.breakdown h3 {
		font-size: 0.72rem;
		margin: 0 0 0.85rem;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		font-weight: 600;
	}

	.breakdown ul {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	.breakdown-item {
		display: grid;
		grid-template-columns: 88px 1fr auto;
		align-items: center;
		gap: 0.85rem;
		padding: 0.65rem 0.85rem;
		border-radius: var(--radius-standard);
		background: rgba(255, 255, 255, 0.018);
		border: 1px solid var(--border-frost);
		font-size: 0.8125rem;
	}

	.bk-tipo {
		font-size: 0.68rem;
		font-weight: 600;
		text-transform: uppercase;
		color: var(--color-text-tertiary);
		letter-spacing: 0.06em;
	}

	.bk-label {
		color: var(--color-text-secondary);
		letter-spacing: -0.1px;
	}

	.bk-valor {
		font-variant-numeric: tabular-nums;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.breakdown-item.desconto {
		background: rgba(34, 255, 153, 0.04);
		border-color: var(--color-green-4);
	}

	.breakdown-item.desconto .bk-tipo,
	.breakdown-item.desconto .bk-valor {
		color: var(--color-green-10);
	}

	.breakdown-item.acrescimo {
		background: rgba(255, 197, 61, 0.04);
		border-color: rgba(255, 197, 61, 0.25);
	}

	.breakdown-item.acrescimo .bk-tipo,
	.breakdown-item.acrescimo .bk-valor {
		color: var(--color-yellow-9);
	}

	.breakdown-item.total {
		background: linear-gradient(135deg, rgba(255, 128, 31, 0.1), rgba(255, 89, 0, 0.06));
		border-color: var(--color-orange-4);
	}

	.breakdown-item.total .bk-tipo {
		color: var(--color-orange-11);
	}

	.breakdown-item.total .bk-label {
		color: var(--color-text-primary);
		font-weight: 500;
	}

	.breakdown-item.total .bk-valor {
		color: var(--color-orange-11);
		font-size: 0.95rem;
	}

	.disclaimer {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.72rem;
		color: var(--color-text-tertiary);
		margin: 1.25rem 0 0;
		letter-spacing: -0.1px;
	}

	:global(.spin) {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
