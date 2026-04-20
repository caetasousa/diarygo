<script lang="ts">
	interface Props {
		value: string; // bindable, formato YYYY-MM-DDTHH:mm
		minAdiantamentoHoras?: number;
	}

	let { value = $bindable(''), minAdiantamentoHoras = 24 }: Props = $props();

	let datePart = $state('');
	let timePart = $state('');
	let erroAntecedencia = $state(false);

	// Sincroniza datePart/timePart com value externo ao montar
	$effect(() => {
		if (value && value.includes('T')) {
			const [d, t] = value.split('T');
			datePart = d;
			timePart = t.slice(0, 5);
		}
	});

	// Mínimo = amanhã
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
</script>

<div class="calendario">
	<div class="campos">
		<div class="campo">
			<label for="agend-data">Data</label>
			<input
				id="agend-data"
				type="date"
				min={minDate}
				bind:value={datePart}
				onchange={combinar}
				onblur={combinar}
				aria-invalid={erroAntecedencia}
				maxlength="10"
				autocomplete="off"
			/>
		</div>
		<div class="campo">
			<label for="agend-hora">Horário</label>
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
	{#if erroAntecedencia}
		<p class="erro" role="alert">Mínimo {minAdiantamentoHoras}h de antecedência.</p>
	{/if}
</div>

<style>
	.calendario { display: flex; flex-direction: column; gap: 8px; }

	.campos {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 12px;
	}

	.campo { display: flex; flex-direction: column; gap: 4px; }

	.campo label {
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
	}
	.campo input:focus { border-color: var(--color-orange-10); }
	.campo input[aria-invalid='true'] { border-color: #ef4444; }

	.erro {
		font-size: 0.75rem;
		color: #ef4444;
		margin: 0;
	}
</style>
