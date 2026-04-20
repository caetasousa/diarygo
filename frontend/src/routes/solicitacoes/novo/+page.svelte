<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type {
		EnderecoResponse,
		CategoriaResponse,
		OpcionalResponse,
		RegiaoResponse,
		FrequenciaServico
	} from '$lib/types';
	import StepIndicator from '$lib/components/StepIndicator.svelte';
	import EnderecoSelector from '$lib/components/EnderecoSelector.svelte';
	import CalendarioAgendamento from '$lib/components/CalendarioAgendamento.svelte';
	import ResumoSolicitacao from '$lib/components/ResumoSolicitacao.svelte';
	import ModalConfirmacao from '$lib/components/ModalConfirmacao.svelte';
	import { Plus, Minus } from 'lucide-svelte';

	const STEPS = ['Endereço', 'Serviço', 'Agenda', 'Resumo'];

	let step = $state(0);

	// Step 1 — Endereço
	let enderecos = $state<EnderecoResponse[]>([]);
	let enderecoID = $state('');
	let carregandoEnderecos = $state(true);

	// Step 2 — Categoria / cômodos / opcionais
	let categorias = $state<CategoriaResponse[]>([]);
	let categoriaID = $state('');
	let opcionais = $state<OpcionalResponse[]>([]);
	let opcionaisSelecionados = $state<Set<string>>(new Set());
	let numQuartos = $state(1);
	let numBanheiros = $state(1);
	let numSalas = $state(1);
	let numCozinhas = $state(1);

	// Step 3 — Data / frequência / região
	let dataServico = $state('');
	let frequencia = $state<FrequenciaServico>('UNICA');
	let regioes = $state<RegiaoResponse[]>([]);
	let regiaoID = $state('');

	// Step 4 — Resumo / confirmação
	import type { CalculoPrecoResponse } from '$lib/types';
	let calculo = $state<CalculoPrecoResponse | null>(null);
	let calculando = $state(false);
	let observacao = $state('');
	let modalAberto = $state(false);
	let enviando = $state(false);

	// Dados derivados para ResumoSolicitacao
	let enderecoSelecionado = $derived(enderecos.find((e) => e.id === enderecoID) ?? null);
	let categoriaSelecionada = $derived(categorias.find((c) => c.id === categoriaID) ?? null);
	let regiaoSelecionada = $derived(regioes.find((r) => r.id === regiaoID) ?? null);

	let horasAteServico = $derived(() => {
		if (!dataServico) return Infinity;
		return (new Date(dataServico).getTime() - Date.now()) / 3600000;
	});

	onMount(async () => {
		try {
			const [ends, cats, regs] = await Promise.all([
				api.listarEnderecos(),
				api.listarCategorias(),
				api.listarRegioes()
			]);
			enderecos = ends;
			categorias = cats;
			regioes = regs;
			if (cats.length > 0) categoriaID = cats[0].id;
			// Pré-seleciona endereço principal
			const principal = ends.find((e) => e.principal);
			if (principal) enderecoID = principal.id;
		} catch {
			toasts.error('Erro ao carregar dados. Tente novamente.');
		} finally {
			carregandoEnderecos = false;
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

	// Auto-seleciona região com base no endereço (heurística: primeiro disponível)
	$effect(() => {
		if (regioes.length > 0 && !regiaoID) {
			regiaoID = regioes[0].id;
		}
	});

	function toggleOpcional(id: string) {
		const novo = new Set(opcionaisSelecionados);
		if (novo.has(id)) novo.delete(id);
		else novo.add(id);
		opcionaisSelecionados = novo;
	}

	function inc(
		campo: 'numQuartos' | 'numBanheiros' | 'numSalas' | 'numCozinhas',
		delta: number
	) {
		if (campo === 'numQuartos') numQuartos = Math.max(1, Math.min(20, numQuartos + delta));
		if (campo === 'numBanheiros') numBanheiros = Math.max(0, Math.min(20, numBanheiros + delta));
		if (campo === 'numSalas') numSalas = Math.max(0, Math.min(20, numSalas + delta));
		if (campo === 'numCozinhas') numCozinhas = Math.max(0, Math.min(20, numCozinhas + delta));
	}

	async function irParaResumo() {
		calculando = true;
		try {
			calculo = await api.calcularPreco({
				categoria_id: categoriaID,
				regiao_id: regiaoID,
				num_quartos: numQuartos,
				num_banheiros: numBanheiros,
				num_salas: numSalas,
				num_cozinhas: numCozinhas,
				opcionais_ids: Array.from(opcionaisSelecionados),
				frequencia,
				data_servico: dataServico || undefined
			});
			step = 3;
		} catch (e: unknown) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao calcular orçamento.');
		} finally {
			calculando = false;
		}
	}

	async function confirmarSolicitacao() {
		modalAberto = false;
		enviando = true;
		try {
			const resp = await api.criarSolicitacao({
				endereco_id: enderecoID,
				categoria_id: categoriaID,
				regiao_id: regiaoID,
				num_quartos: numQuartos,
				num_banheiros: numBanheiros,
				num_salas: numSalas,
				num_cozinhas: numCozinhas,
				opcionais_ids: Array.from(opcionaisSelecionados),
				frequencia,
				data_servico: dataServico,
				observacao: observacao || undefined
			});
			window.location.href = `/dashboard/solicitacoes/${resp.id}`;
		} catch (e: unknown) {
			toasts.error(e instanceof Error ? e.message : 'Não foi possível criar a solicitação.');
			enviando = false;
		}
	}

	const cmodRows = $derived([
		{ label: 'Quartos', campo: 'numQuartos' as const, valor: numQuartos },
		{ label: 'Banheiros', campo: 'numBanheiros' as const, valor: numBanheiros },
		{ label: 'Salas', campo: 'numSalas' as const, valor: numSalas },
		{ label: 'Cozinhas', campo: 'numCozinhas' as const, valor: numCozinhas }
	]);
</script>

<div class="container">
	<div class="header">
		<h1>Solicitar serviço</h1>
		<StepIndicator steps={STEPS} current={step} />
	</div>

	<div class="card">
		{#if step === 0}
			<!-- ══════════════ Step 1 — Endereço ══════════════ -->
			<h2>Onde será o serviço?</h2>
			{#if carregandoEnderecos}
				<p class="loading">Carregando endereços…</p>
			{:else}
				<EnderecoSelector
					{enderecos}
					bind:selectedId={enderecoID}
					onNovo={() => (window.location.href = '/dashboard/enderecos')}
				/>
			{/if}
			<div class="acoes">
				<button
					class="btn-primary"
					disabled={!enderecoID}
					onclick={() => (step = 1)}
				>Continuar</button>
			</div>

		{:else if step === 1}
			<!-- ══════════════ Step 2 — Serviço ══════════════ -->
			<h2>Qual serviço você precisa?</h2>

			<div class="field">
				<label for="categoria">Tipo de serviço</label>
				<select id="categoria" bind:value={categoriaID}>
					{#each categorias as cat (cat.id)}
						<option value={cat.id}>{cat.nome}</option>
					{/each}
				</select>
			</div>

			<fieldset class="comodos">
				<legend>Cômodos</legend>
				{#each cmodRows as row}
					<div class="comodo-row">
						<span class="comodo-label">{row.label}</span>
						<div class="stepper">
							<button
								type="button"
								class="step-btn"
								onclick={() => inc(row.campo, -1)}
								aria-label="Diminuir {row.label}"
							><Minus size={14} /></button>
							<span class="step-val">{row.valor}</span>
							<button
								type="button"
								class="step-btn"
								onclick={() => inc(row.campo, 1)}
								aria-label="Aumentar {row.label}"
							><Plus size={14} /></button>
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
								<small>+{op.valor_extra.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })} · +{op.tempo_extra_min} min</small>
							</span>
						</label>
					{/each}
				</fieldset>
			{/if}

			<div class="acoes">
				<button class="btn-secondary" onclick={() => (step = 0)}>Voltar</button>
				<button class="btn-primary" disabled={!categoriaID} onclick={() => (step = 2)}>Continuar</button>
			</div>

		{:else if step === 2}
			<!-- ══════════════ Step 3 — Agenda ══════════════ -->
			<h2>Quando você precisa?</h2>

			<div class="field">
				<span class="field-label">Data e horário</span>
				<CalendarioAgendamento bind:value={dataServico} />
			</div>

			<div class="field">
				<label for="freq">Frequência</label>
				<select id="freq" bind:value={frequencia}>
					<option value="UNICA">Única (sem desconto)</option>
					<option value="SEMANAL">Semanal (−10%)</option>
					<option value="QUINZENAL">Quinzenal (−5%)</option>
					<option value="DUAS_POR_SEMANA">2× por semana (−15%)</option>
				</select>
			</div>

			{#if regioes.length > 1}
				<div class="field">
					<label for="regiao">Região de atendimento</label>
					<select id="regiao" bind:value={regiaoID}>
						{#each regioes as r (r.id)}
							<option value={r.id}>{r.nome} — {r.cidade}/{r.estado}</option>
						{/each}
					</select>
				</div>
			{/if}

			<div class="acoes">
				<button class="btn-secondary" onclick={() => (step = 1)}>Voltar</button>
				<button
					class="btn-primary"
					disabled={!dataServico || calculando}
					onclick={irParaResumo}
				>
					{calculando ? 'Calculando…' : 'Ver resumo'}
				</button>
			</div>

		{:else if step === 3}
			<!-- ══════════════ Step 4 — Resumo ══════════════ -->
			<h2>Confirme sua solicitação</h2>

			{#if calculo && enderecoSelecionado && categoriaSelecionada && regiaoSelecionada}
				<ResumoSolicitacao
					categoria={categoriaSelecionada}
					regiao={regiaoSelecionada}
					endereco={enderecoSelecionado}
					breakdown={calculo.itens}
					valorTotal={calculo.valor_total}
					duracaoMin={calculo.duracao_min}
					{dataServico}
				/>
			{/if}

			<div class="field">
				<label for="obs">Observações <small>(opcional, máx. 500 caracteres)</small></label>
				<textarea
					id="obs"
					bind:value={observacao}
					maxlength="500"
					rows="3"
					autocomplete="off"
					placeholder="Informações adicionais para a profissional (porteiro, animais, etc.)"
				></textarea>
			</div>

			<div class="acoes">
				<button class="btn-secondary" onclick={() => (step = 2)}>Voltar</button>
				<button
					class="btn-primary"
					disabled={enviando}
					onclick={() => (modalAberto = true)}
				>
					{enviando ? 'Enviando…' : 'Confirmar solicitação'}
				</button>
			</div>
		{/if}
	</div>
</div>

<ModalConfirmacao
	bind:open={modalAberto}
	titulo="Confirmar solicitação"
	mensagem={horasAteServico() < 48
		? 'Atenção: cancelamentos com menos de 24h de antecedência penalizam seu score em 5 pontos. Deseja confirmar?'
		: 'Confirmar a solicitação de serviço com os dados informados?'}
	textoConfirmar="Confirmar"
	onConfirm={confirmarSolicitacao}
	onCancel={() => (modalAberto = false)}
/>

<style>
	.container {
		max-width: 640px;
		margin: 0 auto;
		padding: 2rem 1.5rem 4rem;
	}

	.header {
		display: flex;
		flex-direction: column;
		gap: 20px;
		margin-bottom: 24px;
	}

	.header h1 {
		font-size: 1.375rem;
		font-weight: 600;
		color: var(--color-text-primary);
		letter-spacing: -0.5px;
		margin: 0;
	}

	.card {
		background: rgba(255, 255, 255, 0.018);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-large);
		padding: 24px;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.card h2 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.loading {
		color: var(--color-text-tertiary);
		font-size: 0.875rem;
		text-align: center;
		padding: 24px 0;
	}

	.acoes {
		display: flex;
		gap: 10px;
		justify-content: flex-end;
		padding-top: 4px;
	}

	.btn-primary,
	.btn-secondary {
		padding: 10px 20px;
		border-radius: var(--radius-pill);
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.12s, border-color 0.12s, opacity 0.12s;
	}

	.btn-primary {
		background: var(--color-orange-10);
		border: 1px solid var(--color-orange-10);
		color: #fff;
	}
	.btn-primary:hover:not(:disabled) { background: var(--color-orange-11); }
	.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

	.btn-secondary {
		background: transparent;
		border: 1px solid var(--border-frost);
		color: var(--color-text-primary);
	}
	.btn-secondary:hover { background: var(--bg-hover-subtle); border-color: var(--border-frost-hover); }

	.field {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.field label,
	.field-label {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-secondary);
	}

	.field select,
	.field textarea {
		padding: 10px 12px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		background: var(--bg-card);
		color: var(--color-text-primary);
		font-size: 0.875rem;
		outline: none;
		transition: border-color 0.12s;
		width: 100%;
		box-sizing: border-box;
	}
	.field select:focus,
	.field textarea:focus { border-color: var(--color-orange-10); }

	.field textarea { resize: vertical; font-family: inherit; line-height: 1.5; }

	/* Cômodos */
	.comodos {
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		padding: 16px;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
	.comodos legend {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		padding: 0 4px;
	}

	.comodo-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.comodo-label {
		font-size: 0.875rem;
		color: var(--color-text-primary);
	}

	.stepper {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.step-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		border: 1px solid var(--border-frost);
		border-radius: 50%;
		background: var(--bg-card);
		color: var(--color-text-primary);
		cursor: pointer;
		transition: background 0.1s, border-color 0.1s;
	}
	.step-btn:hover { background: var(--bg-hover-subtle); border-color: var(--border-frost-hover); }

	.step-val {
		font-size: 0.9375rem;
		font-weight: 600;
		min-width: 20px;
		text-align: center;
		color: var(--color-text-primary);
		font-variant-numeric: tabular-nums;
	}

	/* Opcionais */
	.opcionais {
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		padding: 16px;
		display: flex;
		flex-direction: column;
		gap: 10px;
	}
	.opcionais legend {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		padding: 0 4px;
	}

	.opcional {
		display: flex;
		align-items: flex-start;
		gap: 10px;
		cursor: pointer;
	}

	.opcional input { margin-top: 2px; accent-color: var(--color-orange-10); }

	.opcional-info {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}
	.opcional-info strong { font-size: 0.875rem; color: var(--color-text-primary); font-weight: 500; }
	.opcional-info small { font-size: 0.75rem; color: var(--color-text-tertiary); }
</style>
