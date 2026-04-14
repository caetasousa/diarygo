<script lang="ts">
	interface Props {
		value: string;
		error?: string;
		disabled?: boolean;
	}

	let { value = $bindable(''), error = '', disabled = false }: Props = $props();

	function aplicarMascara(v: string): string {
		const d = v.replace(/\D/g, '').slice(0, 11);
		if (d.length <= 3) return d;
		if (d.length <= 6) return `${d.slice(0, 3)}.${d.slice(3)}`;
		if (d.length <= 9) return `${d.slice(0, 3)}.${d.slice(3, 6)}.${d.slice(6)}`;
		return `${d.slice(0, 3)}.${d.slice(3, 6)}.${d.slice(6, 9)}-${d.slice(9)}`;
	}

	function validarCPF(cpf: string): boolean {
		const d = cpf.replace(/\D/g, '');
		if (d.length !== 11) return false;
		if (/^(\d)\1{10}$/.test(d)) return false;
		let soma = 0;
		for (let i = 0; i < 9; i++) soma += parseInt(d[i]) * (10 - i);
		let r = soma % 11;
		const d1 = r < 2 ? 0 : 11 - r;
		if (parseInt(d[9]) !== d1) return false;
		soma = 0;
		for (let i = 0; i < 10; i++) soma += parseInt(d[i]) * (11 - i);
		r = soma % 11;
		const d2 = r < 2 ? 0 : 11 - r;
		return parseInt(d[10]) === d2;
	}

	let cpfValido = $derived(value.replace(/\D/g, '').length === 11 ? validarCPF(value) : true);

	function handleInput(e: Event) {
		const input = e.target as HTMLInputElement;
		value = aplicarMascara(input.value);
	}
</script>

<div class="form-group">
	<label for="cpf-input" class="form-label">CPF</label>
	<input
		id="cpf-input"
		type="text"
		inputmode="numeric"
		autocomplete="off"
		maxlength="14"
		placeholder="000.000.000-00"
		class="form-input"
		class:input-error={!cpfValido || !!error}
		{value}
		{disabled}
		oninput={handleInput}
	/>
	{#if !cpfValido && value.replace(/\D/g, '').length === 11}
		<p class="form-error">CPF inválido</p>
	{:else if error}
		<p class="form-error">{error}</p>
	{/if}
</div>
