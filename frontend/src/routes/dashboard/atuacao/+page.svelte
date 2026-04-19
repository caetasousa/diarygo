<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type { RegiaoResponse, DisponibilidadeResponse } from '$lib/types';

	// Regioes
	let todasRegioes = $state<RegiaoResponse[]>([]);
	let regioesAtuacao = $state<RegiaoResponse[]>([]);
	let selecionadas = $state<Set<string>>(new Set());
	let regioesBaseIds = $state<string[]>([]);

	// Disponibilidade
	const diasSemana = [
		{ value: 0, label: 'Domingo',       short: 'Dom' },
		{ value: 1, label: 'Segunda-feira', short: 'Seg' },
		{ value: 2, label: 'Terça-feira',   short: 'Ter' },
		{ value: 3, label: 'Quarta-feira',  short: 'Qua' },
		{ value: 4, label: 'Quinta-feira',  short: 'Qui' },
		{ value: 5, label: 'Sexta-feira',   short: 'Sex' },
		{ value: 6, label: 'Sábado',        short: 'Sáb' }
	];

	type SlotLocal = { ativo: boolean; horaInicio: string; horaFim: string };
	let slots = $state<SlotLocal[]>(
		diasSemana.map(() => ({ ativo: false, horaInicio: '08:00', horaFim: '17:00' }))
	);
	let slotsBase = $state<string>('');

	let carregando = $state(true);
	let salvando = $state(false);

	onMount(async () => {
		try {
			const [todas, atuacao, disp] = await Promise.all([
				api.listarRegioes(),
				api.listarRegioesAtuacao().catch(() => [] as RegiaoResponse[]),
				api.listarDisponibilidades().catch(() => [] as DisponibilidadeResponse[])
			]);
			todasRegioes = todas;
			regioesAtuacao = atuacao;
			selecionadas = new Set(atuacao.map((r) => r.id));
			regioesBaseIds = [...selecionadas].sort();

			for (const d of disp) {
				if (d.dia_semana >= 0 && d.dia_semana <= 6) {
					slots[d.dia_semana] = {
						ativo: true,
						horaInicio: d.hora_inicio,
						horaFim: d.hora_fim
					};
				}
			}
			slotsBase = JSON.stringify(slots);
		} catch {
			toasts.error('Erro ao carregar dados de atuação');
		} finally {
			carregando = false;
		}
	});

	function toggleRegiao(id: string) {
		const novas = new Set(selecionadas);
		if (novas.has(id)) novas.delete(id);
		else novas.add(id);
		selecionadas = novas;
	}

	function erroSlot(slot: SlotLocal): string | null {
		if (!slot.ativo) return null;
		if (slot.horaFim <= slot.horaInicio) return 'Fim deve ser após o início';
		return null;
	}

	let temErroSlot = $derived(slots.some((s) => erroSlot(s) !== null));

	let regioesMudaram = $derived(
		JSON.stringify([...selecionadas].sort()) !== JSON.stringify(regioesBaseIds)
	);
	let slotsMudaram = $derived(JSON.stringify(slots) !== slotsBase);
	let temMudancas = $derived(regioesMudaram || slotsMudaram);

	let regioesSelecionadasCount = $derived(selecionadas.size);
	let diasAtivos = $derived(slots.filter((s) => s.ativo).length);
	let horasSemana = $derived(
		slots.reduce((acc, s) => {
			if (!s.ativo) return acc;
			const ini = parseInt(s.horaInicio.slice(0, 2)) * 60 + parseInt(s.horaInicio.slice(3));
			const fim = parseInt(s.horaFim.slice(0, 2)) * 60 + parseInt(s.horaFim.slice(3));
			return acc + Math.max(0, (fim - ini) / 60);
		}, 0)
	);

	async function salvarTudo() {
		if (temErroSlot) return;
		salvando = true;
		try {
			const promises: Promise<unknown>[] = [];
			if (regioesMudaram) {
				promises.push(api.definirRegioesAtuacao({ regiao_ids: [...selecionadas] }));
			}
			if (slotsMudaram) {
				const slotsAtivos = slots
					.map((s, i) => ({ ...s, dia: i }))
					.filter((s) => s.ativo)
					.map((s) => ({ dia_semana: s.dia, hora_inicio: s.horaInicio, hora_fim: s.horaFim }));
				promises.push(api.definirDisponibilidades({ slots: slotsAtivos }));
			}
			await Promise.all(promises);

			regioesAtuacao = todasRegioes.filter((r) => selecionadas.has(r.id));
			regioesBaseIds = [...selecionadas].sort();
			slotsBase = JSON.stringify(slots);

			toasts.success('Atuação atualizada!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar atuação');
		} finally {
			salvando = false;
		}
	}
</script>

<svelte:head>
	<title>Atuação — DiaryGo</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<div>
			<h1 class="page-title">Atuação</h1>
			<p class="page-sub">Onde e quando você está disponível para atender em Goiânia.</p>
		</div>
		<div class="progress-stats">
			<div class="stat">
				<span class="stat-num">{regioesSelecionadasCount}</span>
				<span class="stat-label">regiões</span>
			</div>
			<div class="stat">
				<span class="stat-num">{diasAtivos}</span>
				<span class="stat-label">dia{diasAtivos !== 1 ? 's' : ''}</span>
			</div>
			<div class="stat">
				<span class="stat-num">{horasSemana.toFixed(0)}h</span>
				<span class="stat-label">por semana</span>
			</div>
		</div>
	</div>

	{#if carregando}
		<div class="loading-state">
			<div class="loading-spinner"></div>
			<span>Carregando...</span>
		</div>
	{:else}
		<div class="layout">
			<section class="col">
				<header class="col-head">
					<h2 class="col-title">Onde atendo</h2>
					<span class="col-sub">
						{regioesSelecionadasCount} de {todasRegioes.length} regiões selecionadas
					</span>
				</header>
				<div class="regions-list">
					{#each todasRegioes as regiao (regiao.id)}
						<button
							type="button"
							class="region-card"
							class:selected={selecionadas.has(regiao.id)}
							onclick={() => toggleRegiao(regiao.id)}
						>
							<div class="region-check">
								{#if selecionadas.has(regiao.id)}
									<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
								{/if}
							</div>
							<div class="region-info">
								<span class="region-nome">{regiao.nome}</span>
								<span class="region-cep">
									CEP {regiao.cep_inicio.slice(0, 5)}-{regiao.cep_inicio.slice(5)} a
									{regiao.cep_fim.slice(0, 5)}-{regiao.cep_fim.slice(5)}
								</span>
							</div>
						</button>
					{/each}
				</div>
			</section>

			<section class="col">
				<header class="col-head">
					<h2 class="col-title">Quando atendo</h2>
					<span class="col-sub">
						{diasAtivos} dia{diasAtivos !== 1 ? 's' : ''} · {horasSemana.toFixed(0)}h/semana
					</span>
				</header>
				<div class="slots-list">
					{#each diasSemana as dia, i}
						{@const slot = slots[i]}
						{@const erro = erroSlot(slot)}
						<div class="slot-row" class:slot-active={slot.ativo}>
							<button
								class="day-toggle"
								class:on={slot.ativo}
								type="button"
								aria-label="Ativar {dia.label}"
								onclick={() => (slots[i].ativo = !slots[i].ativo)}
							>
								<span class="day-short">{dia.short}</span>
								<span class="day-full">{dia.label}</span>
							</button>

							{#if slot.ativo}
								<div class="time-inputs">
									<input
										type="time"
										class="time-input"
										class:input-error={!!erro}
										aria-label="Início"
										bind:value={slots[i].horaInicio}
									/>
									<span class="time-sep">→</span>
									<input
										type="time"
										class="time-input"
										class:input-error={!!erro}
										aria-label="Fim"
										bind:value={slots[i].horaFim}
									/>
								</div>
								{#if erro}<span class="time-error">{erro}</span>{/if}
							{:else}
								<span class="slot-off">Indisponível</span>
							{/if}
						</div>
					{/each}
				</div>
			</section>
		</div>

		<div class="save-bar" class:visible={temMudancas}>
			<span class="save-hint">
				{#if regioesMudaram && slotsMudaram}
					Alterações em regiões e horários
				{:else if regioesMudaram}
					Alterações em regiões
				{:else}
					Alterações em horários
				{/if}
			</span>
			<button
				type="button"
				class="btn btn-white"
				onclick={salvarTudo}
				disabled={salvando || temErroSlot}
			>
				{salvando ? 'Salvando...' : 'Salvar alterações'}
			</button>
		</div>
	{/if}
</div>

<style>
	.page { padding: 40px; max-width: 1200px; margin: 0 auto; width: 100%; }

	.page-header {
		display: flex; align-items: flex-start;
		justify-content: space-between; gap: 16px;
		margin-bottom: 24px; flex-wrap: wrap;
	}

	.page-title {
		font-size: 1.75rem; font-weight: 400;
		letter-spacing: -1px; color: var(--color-text-primary);
		line-height: 1; margin-bottom: 6px;
	}
	.page-sub { font-size: 0.875rem; color: var(--color-text-tertiary); }

	.progress-stats { display: flex; gap: 24px; flex-wrap: wrap; }
	.stat { display: flex; flex-direction: column; gap: 2px; }
	.stat-num {
		font-size: 1.25rem; font-weight: 500; color: var(--color-text-primary);
		font-variant-numeric: tabular-nums; line-height: 1;
	}
	.stat-label { font-size: 0.75rem; color: var(--color-text-tertiary); }

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

	.layout {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 20px;
		align-items: start;
		margin-bottom: 20px;
	}

	.col {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		box-shadow: var(--shadow-ring);
		overflow: hidden;
	}

	.col-head {
		padding: 14px 18px;
		border-bottom: 1px solid var(--border-frost);
		display: flex; justify-content: space-between; align-items: baseline;
		flex-wrap: wrap; gap: 6px;
	}
	.col-title {
		font-size: 0.875rem; font-weight: 500; color: var(--color-text-primary);
	}
	.col-sub {
		font-size: 0.75rem; color: var(--color-text-tertiary);
	}

	.regions-list {
		display: flex; flex-direction: column; gap: 4px;
		padding: 10px;
		max-height: 520px; overflow-y: auto;
	}
	.region-card {
		display: flex; gap: 10px; align-items: flex-start;
		padding: 10px 12px;
		background: transparent;
		border: 1px solid transparent;
		border-radius: var(--radius-standard);
		cursor: pointer; text-align: left;
		font-family: inherit;
		transition: background 0.12s, border-color 0.12s;
	}
	.region-card:hover {
		background: var(--bg-hover-subtle);
		border-color: var(--border-frost);
	}
	.region-card.selected {
		background: var(--color-blue-wash);
		border-color: var(--color-blue-5);
	}
	.region-check {
		width: 18px; height: 18px; border-radius: var(--radius-sharp);
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
	.region-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
	.region-nome {
		font-size: 0.875rem; font-weight: 500; color: var(--color-text-primary);
	}
	.region-cep {
		font-size: 0.7rem; color: var(--color-text-tertiary);
		font-family: var(--font-mono);
	}

	.slots-list {
		display: flex; flex-direction: column;
	}
	.slot-row {
		display: flex; align-items: center; gap: 12px; flex-wrap: wrap;
		padding: 12px 18px;
		border-bottom: 1px solid var(--border-frost);
		transition: background 0.12s;
	}
	.slot-row:last-child { border-bottom: none; }
	.slot-row.slot-active { background: var(--bg-card); }

	.day-toggle {
		min-width: 130px;
		display: flex; align-items: center; gap: 10px;
		background: none; border: none; cursor: pointer;
		padding: 0; text-align: left; font-family: inherit;
	}
	.day-short {
		display: none;
		width: 32px; height: 32px; border-radius: var(--radius-standard);
		background: var(--bg-hover-subtle);
		border: 1px solid var(--border-frost);
		font-size: 0.75rem; font-weight: 600;
		color: var(--color-text-tertiary);
		align-items: center; justify-content: center;
	}
	.day-full {
		font-size: 0.875rem; font-weight: 500;
		color: var(--color-text-tertiary);
	}
	.day-toggle.on .day-full { color: var(--color-text-primary); }

	.time-inputs {
		display: flex; align-items: center; gap: 8px; flex: 1; flex-wrap: wrap;
	}
	.time-input {
		background: var(--bg-hover-subtle);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-standard);
		padding: 6px 10px;
		font-family: var(--font-mono);
		font-size: 0.875rem;
		color: var(--color-text-primary);
		outline: none;
		width: 100px;
		color-scheme: dark;
		transition: border-color 0.12s;
	}
	.time-input:focus { border-color: var(--color-blue-10); }
	.time-input.input-error { border-color: var(--color-red-10); }
	.time-sep { font-size: 0.875rem; color: var(--color-text-tertiary); }
	.time-error { font-size: 0.75rem; color: var(--color-red-10); }
	.slot-off { font-size: 0.8125rem; color: var(--color-text-tertiary); }

	.save-bar {
		position: sticky;
		bottom: 20px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		padding: 12px 16px;
		background: rgba(10,10,10,0.92); /* save-bar backdrop opaco */
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		backdrop-filter: blur(8px);
		opacity: 0;
		transform: translateY(8px);
		pointer-events: none;
		transition: opacity 0.2s, transform 0.2s;
	}
	.save-bar.visible {
		opacity: 1; transform: translateY(0); pointer-events: auto;
	}
	.save-hint { font-size: 0.875rem; color: var(--color-text-secondary); }

	@media (max-width: 900px) {
		.layout { grid-template-columns: 1fr; }
	}
	@media (max-width: 640px) {
		.page { padding: 24px 16px; }
		.day-full { display: none; }
		.day-short { display: flex; }
	}
</style>
