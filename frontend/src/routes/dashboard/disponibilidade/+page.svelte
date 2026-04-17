<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type { DisponibilidadeResponse } from '$lib/types';

	const diasSemana = [
		{ value: 0, label: 'Domingo',        short: 'Dom' },
		{ value: 1, label: 'Segunda-feira',   short: 'Seg' },
		{ value: 2, label: 'Terça-feira',     short: 'Ter' },
		{ value: 3, label: 'Quarta-feira',    short: 'Qua' },
		{ value: 4, label: 'Quinta-feira',    short: 'Qui' },
		{ value: 5, label: 'Sexta-feira',     short: 'Sex' },
		{ value: 6, label: 'Sábado',          short: 'Sáb' },
	];

	type SlotLocal = { ativo: boolean; horaInicio: string; horaFim: string };

	let carregando = $state(true);
	let salvando = $state(false);
	let slots = $state<SlotLocal[]>(
		diasSemana.map(() => ({ ativo: false, horaInicio: '08:00', horaFim: '17:00' }))
	);

	onMount(async () => {
		try {
			const disp = await api.listarDisponibilidades().catch(() => [] as DisponibilidadeResponse[]);
			for (const d of disp) {
				if (d.dia_semana >= 0 && d.dia_semana <= 6) {
					slots[d.dia_semana] = { ativo: true, horaInicio: d.hora_inicio, horaFim: d.hora_fim };
				}
			}
		} finally {
			carregando = false;
		}
	});

	function erroSlot(slot: SlotLocal): string | null {
		if (!slot.ativo) return null;
		if (slot.horaFim <= slot.horaInicio) return 'Fim deve ser após o início';
		return null;
	}

	let temErro = $derived(slots.some((s) => erroSlot(s) !== null));
	let ativosCount = $derived(slots.filter(s => s.ativo).length);

	async function salvar() {
		if (temErro) return;
		salvando = true;
		try {
			const slotsAtivos = slots
				.map((s, i) => ({ ...s, dia: i }))
				.filter((s) => s.ativo)
				.map((s) => ({ dia_semana: s.dia, hora_inicio: s.horaInicio, hora_fim: s.horaFim }));
			await api.definirDisponibilidades({ slots: slotsAtivos });
			toasts.success('Disponibilidade salva!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar disponibilidade');
		} finally {
			salvando = false;
		}
	}
</script>

<svelte:head>
	<title>Disponibilidade — DiaryGo</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<div>
			<h1 class="page-title">Disponibilidade</h1>
			<p class="page-sub">Dias e horários em que você está disponível para atender.</p>
		</div>
		{#if ativosCount > 0}
			<span class="badge badge-green">{ativosCount} dia{ativosCount !== 1 ? 's' : ''} ativo{ativosCount !== 1 ? 's' : ''}</span>
		{/if}
	</div>

	{#if carregando}
		<div class="loading-state">
			<div class="loading-spinner"></div>
			<span>Carregando...</span>
		</div>
	{:else}
		<div class="slots-list">
			{#each diasSemana as dia, i}
				{@const slot = slots[i]}
				{@const erro = erroSlot(slot)}
				<div class="slot-row" class:slot-active={slot.ativo}>
					<button
						class="day-toggle"
						class:on={slot.ativo}
						onclick={() => slots[i].ativo = !slots[i].ativo}
						type="button"
						aria-label="Ativar {dia.label}"
					>
						<span class="day-short">{dia.short}</span>
						<span class="day-full">{dia.label}</span>
					</button>

					{#if slot.ativo}
						<div class="time-inputs">
							<div class="time-group">
								<label for="inicio-{i}" class="time-label">Início</label>
								<input
									id="inicio-{i}" type="time" class="time-input"
									class:input-error={!!erro}
									bind:value={slots[i].horaInicio}
								/>
							</div>
							<span class="time-sep">→</span>
							<div class="time-group">
								<label for="fim-{i}" class="time-label">Fim</label>
								<input
									id="fim-{i}" type="time" class="time-input"
									class:input-error={!!erro}
									bind:value={slots[i].horaFim}
								/>
							</div>
							{#if erro}
								<span class="time-error">{erro}</span>
							{/if}
						</div>
					{:else}
						<span class="slot-off-label">Indisponível</span>
					{/if}
				</div>
			{/each}
		</div>

		<div class="actions">
			<button
				class="btn btn-white"
				onclick={salvar}
				disabled={salvando || temErro}
			>
				{salvando ? 'Salvando...' : 'Salvar disponibilidade'}
			</button>
			{#if temErro}
				<span class="error-hint">Corrija os horários em vermelho antes de salvar.</span>
			{/if}
		</div>
	{/if}
</div>

<style>
	.page { padding: 40px; max-width: 640px; margin: 0 auto; width: 100%; }

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

	/* Slots */
	.slots-list {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 12px;
		overflow: hidden;
		box-shadow: var(--shadow-ring);
		margin-bottom: 20px;
	}

	.slot-row {
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 14px 18px;
		border-bottom: 1px solid var(--border-frost);
		transition: background 0.12s;
	}

	.slot-row:last-child { border-bottom: none; }
	.slot-row.slot-active { background: rgba(255, 255, 255, 0.02); }

	.day-toggle {
		min-width: 130px;
		display: flex;
		align-items: center;
		gap: 10px;
		background: none;
		border: none;
		cursor: pointer;
		padding: 0;
		text-align: left;
	}

	.day-short {
		display: none;
		width: 32px; height: 32px;
		border-radius: 8px;
		background: rgba(255,255,255,0.05);
		border: 1px solid var(--border-frost);
		font-size: 0.75rem; font-weight: 600;
		color: var(--color-text-tertiary);
		align-items: center; justify-content: center;
		transition: background 0.12s, border-color 0.12s, color 0.12s;
	}

	.day-full {
		font-size: 0.9375rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		transition: color 0.12s;
	}

	.day-toggle.on .day-full { color: var(--color-text-primary); }

	.time-inputs {
		display: flex;
		align-items: center;
		gap: 10px;
		flex: 1;
		flex-wrap: wrap;
	}

	.time-group {
		display: flex;
		flex-direction: column;
		gap: 3px;
	}

	.time-label {
		font-size: 0.6875rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.4px;
	}

	.time-input {
		background: rgba(255,255,255,0.04);
		border: 1px solid var(--border-frost);
		border-radius: 7px;
		padding: 7px 10px;
		font-family: var(--font-mono);
		font-size: 0.9375rem;
		color: var(--color-text-primary);
		outline: none;
		width: 110px;
		transition: border-color 0.12s;
		color-scheme: dark;
	}

	.time-input:focus {
		border-color: var(--color-blue-10);
	}

	.time-input.input-error {
		border-color: var(--color-red-10);
	}

	.time-sep {
		font-size: 0.875rem;
		color: var(--color-text-tertiary);
		margin-top: 16px;
	}

	.time-error {
		font-size: 0.75rem;
		color: var(--color-red-10);
		margin-top: 16px;
	}

	.slot-off-label {
		font-size: 0.8125rem;
		color: var(--color-text-tertiary);
	}

	.actions {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.error-hint {
		font-size: 0.8125rem;
		color: var(--color-red-10);
	}

	@media (max-width: 640px) {
		.page { padding: 24px 16px; }

		.day-full { display: none; }
		.day-short { display: flex; }

		.slot-row { padding: 10px 14px; }
	}
</style>
