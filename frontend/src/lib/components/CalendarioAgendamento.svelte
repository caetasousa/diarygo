<script lang="ts">
	import { Calendar, Clock } from 'lucide-svelte';

	interface Props {
		value: string; // bindable, formato YYYY-MM-DDTHH:mm
		minAdiantamentoHoras?: number;
	}

	let { value = $bindable(''), minAdiantamentoHoras = 24 }: Props = $props();

	// Inicializa datePart/timePart a partir de value apenas na montagem.
	// Após isso, datePart/timePart são a fonte da verdade e propagam para value via combinar().
	const [initDate, initTime] = (() => {
		if (value && value.includes('T')) {
			const [d, t] = value.split('T');
			return [d, t.slice(0, 5)];
		}
		return ['', ''];
	})();

	let datePart = $state(initDate);
	let timePart = $state(initTime);
	let erroAntecedencia = $state(false);

	const minDate = (() => {
		const d = new Date();
		d.setDate(d.getDate() + 1);
		return d.toISOString().split('T')[0];
	})();

	function combinar() {
		if (!datePart || !timePart) {
			value = '';
			erroAntecedencia = false;
			return;
		}
		const iso = `${datePart}T${timePart}`;
		const ms = new Date(iso).getTime() - Date.now();
		erroAntecedencia = ms < minAdiantamentoHoras * 3600 * 1000;
		value = iso;
	}

	const atalhos = $derived.by(() => {
		const opts: { label: string; date: string }[] = [];
		const base = new Date();
		base.setHours(0, 0, 0, 0);
		for (let i = 1; i <= 3; i++) {
			const d = new Date(base);
			d.setDate(d.getDate() + i);
			const iso = d.toISOString().split('T')[0];
			let label = '';
			if (i === 1) label = 'Amanhã';
			else if (i === 2) label = 'Depois de amanhã';
			else
				label = d.toLocaleDateString('pt-BR', { weekday: 'short', day: '2-digit', month: '2-digit' });
			opts.push({ label, date: iso });
		}
		return opts;
	});

	const horarios = ['08:00', '10:00', '13:00', '15:00'];

	function selecionarData(iso: string) {
		datePart = iso;
		if (!timePart) timePart = '09:00';
		combinar();
	}

	function selecionarHora(h: string) {
		timePart = h;
		if (!datePart) datePart = minDate;
		combinar();
	}
</script>

<div class="calendario">
	<div class="campos">
		<div class="campo">
			<label for="agend-data"><Calendar size={12} /> Data</label>
			<input
				id="agend-data"
				type="date"
				min={minDate}
				bind:value={datePart}
				onchange={combinar}
				onblur={combinar}
				aria-invalid={erroAntecedencia}
				autocomplete="off"
			/>
		</div>
		<div class="campo">
			<label for="agend-hora"><Clock size={12} /> Horário</label>
			<input
				id="agend-hora"
				type="time"
				bind:value={timePart}
				onchange={combinar}
				onblur={combinar}
				aria-invalid={erroAntecedencia}
				autocomplete="off"
			/>
		</div>
	</div>

	<div class="atalhos">
		<span class="atalhos-titulo">Sugestões</span>
		<div class="chips">
			{#each atalhos as a}
				<button
					type="button"
					class="chip"
					class:ativo={datePart === a.date}
					onclick={() => selecionarData(a.date)}
				>{a.label}</button>
			{/each}
		</div>
		<div class="chips">
			{#each horarios as h}
				<button
					type="button"
					class="chip"
					class:ativo={timePart === h}
					onclick={() => selecionarHora(h)}
				>{h}</button>
			{/each}
		</div>
	</div>

	{#if erroAntecedencia}
		<p class="erro" role="alert">Mínimo {minAdiantamentoHoras}h de antecedência.</p>
	{/if}
</div>

<style>
	.calendario { display: flex; flex-direction: column; gap: 14px; }

	.campos {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 10px;
	}

	.campo { display: flex; flex-direction: column; gap: 6px; }

	.campo label {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		font-size: 0.75rem;
		color: var(--color-text-secondary);
		font-weight: 500;
	}

	.campo input {
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
		/* Dark mode para inputs date/time nativos (Chromium) */
		color-scheme: dark;
		font-family: inherit;
		line-height: 1.2;
		min-height: 40px;
	}
	.campo input:focus { border-color: var(--color-orange-10); }
	.campo input[aria-invalid='true'] { border-color: #ef4444; }

	/* Ícone do calendar/clock do input nativo (Chromium) */
	.campo input::-webkit-calendar-picker-indicator {
		filter: invert(0.7);
		cursor: pointer;
		opacity: 0.7;
		transition: opacity 0.12s;
	}
	.campo input::-webkit-calendar-picker-indicator:hover { opacity: 1; }

	/* Atalhos */
	.atalhos {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.atalhos-titulo {
		font-size: 0.6875rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
		color: var(--color-text-tertiary);
	}

	.chips {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
	}

	.chip {
		padding: 6px 12px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		background: transparent;
		color: var(--color-text-secondary);
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.12s, border-color 0.12s, color 0.12s;
	}
	.chip:hover { background: var(--bg-hover-subtle); color: var(--color-text-primary); }
	.chip.ativo {
		background: rgba(255, 128, 31, 0.08);
		border-color: var(--color-orange-10);
		color: var(--color-orange-10);
	}

	.erro {
		font-size: 0.75rem;
		color: #ef4444;
		margin: 0;
	}
</style>
