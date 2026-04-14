<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type { DisponibilidadeResponse } from '$lib/types';

	const diasSemana = [
		{ value: 0, label: 'Domingo' },
		{ value: 1, label: 'Segunda-feira' },
		{ value: 2, label: 'Terça-feira' },
		{ value: 3, label: 'Quarta-feira' },
		{ value: 4, label: 'Quinta-feira' },
		{ value: 5, label: 'Sexta-feira' },
		{ value: 6, label: 'Sábado' }
	];

	type SlotLocal = {
		ativo: boolean;
		horaInicio: string;
		horaFim: string;
	};

	let carregando = $state(true);
	let salvando = $state(false);

	// Estado local: um slot por dia da semana (0-6)
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

<div class="container" style="max-width: 640px; padding: 2rem 1rem;">
	<h1 class="text-2xl font-bold" style="margin-bottom: 0.5rem;">Disponibilidade Semanal</h1>
	<p class="text-muted" style="margin-bottom: 2rem;">
		Defina os dias e horários em que você está disponível para atender.
	</p>

	{#if carregando}
		<p class="text-muted">Carregando...</p>
	{:else}
		<div style="display: flex; flex-direction: column; gap: 0.75rem; margin-bottom: 2rem;">
			{#each diasSemana as dia, i}
				{@const slot = slots[i]}
				{@const erro = erroSlot(slot)}
				<div class="card" style="padding: 1rem;">
					<label style="display: flex; align-items: center; gap: 0.75rem; cursor: pointer; margin-bottom: {slot.ativo ? '1rem' : '0'};">
						<input
							type="checkbox"
							style="width: 1.1rem; height: 1.1rem;"
							bind:checked={slots[i].ativo}
						/>
						<strong>{dia.label}</strong>
					</label>

					{#if slot.ativo}
						<div style="display: flex; gap: 1rem; align-items: flex-end; flex-wrap: wrap;">
							<div class="form-group" style="margin: 0;">
								<label for="inicio-{i}" class="form-label" style="font-size: 0.85rem;">Início</label>
								<input
									id="inicio-{i}"
									type="time"
									class="form-input"
									style="width: 130px;"
									bind:value={slots[i].horaInicio}
								/>
							</div>
							<div class="form-group" style="margin: 0;">
								<label for="fim-{i}" class="form-label" style="font-size: 0.85rem;">Fim</label>
								<input
									id="fim-{i}"
									type="time"
									class="form-input"
									style="width: 130px;"
									bind:value={slots[i].horaFim}
								/>
							</div>
						</div>
						{#if erro}
							<p class="form-error" style="margin-top: 0.5rem;">{erro}</p>
						{/if}
					{/if}
				</div>
			{/each}
		</div>

		<button
			class="btn btn-primary btn-full"
			onclick={salvar}
			disabled={salvando || temErro}
		>
			{salvando ? 'Salvando...' : 'Salvar disponibilidade'}
		</button>
	{/if}
</div>
