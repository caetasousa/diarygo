<script lang="ts">
	interface Props {
		value: string;
		error?: string;
		disabled?: boolean;
	}

	let { value = $bindable(''), error = '', disabled = false }: Props = $props();

	function aplicarMascara(v: string): string {
		const d = v.replace(/\D/g, '').slice(0, 11);
		if (d.length <= 2) return d.length ? `(${d}` : '';
		if (d.length <= 6) return `(${d.slice(0, 2)}) ${d.slice(2)}`;
		if (d.length <= 10) return `(${d.slice(0, 2)}) ${d.slice(2, 6)}-${d.slice(6)}`;
		return `(${d.slice(0, 2)}) ${d.slice(2, 7)}-${d.slice(7)}`;
	}

	function handleInput(e: Event) {
		const input = e.target as HTMLInputElement;
		value = aplicarMascara(input.value);
	}
</script>

<div class="form-group">
	<label for="telefone-input" class="form-label">Telefone</label>
	<input
		id="telefone-input"
		type="tel"
		inputmode="numeric"
		autocomplete="tel"
		maxlength="15"
		placeholder="(11) 91234-5678"
		class="form-input"
		class:input-error={!!error}
		{value}
		{disabled}
		oninput={handleInput}
	/>
	{#if error}
		<p class="form-error">{error}</p>
	{/if}
</div>
