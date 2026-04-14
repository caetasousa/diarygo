<script lang="ts">
	interface Props {
		value: string;
		error?: string;
		disabled?: boolean;
		onPreenchido?: (dados: DadosViaCEP) => void;
	}

	export interface DadosViaCEP {
		logradouro: string;
		bairro: string;
		localidade: string;
		uf: string;
	}

	let { value = $bindable(''), error = '', disabled = false, onPreenchido }: Props = $props();

	let buscando = $state(false);
	let erroCEP = $state('');

	function aplicarMascara(v: string): string {
		const d = v.replace(/\D/g, '').slice(0, 8);
		if (d.length <= 5) return d;
		return `${d.slice(0, 5)}-${d.slice(5)}`;
	}

	function handleInput(e: Event) {
		const input = e.target as HTMLInputElement;
		value = aplicarMascara(input.value);
		erroCEP = '';
		const digits = value.replace(/\D/g, '');
		if (digits.length === 8) {
			buscarCEP(digits);
		}
	}

	async function buscarCEP(cep: string) {
		buscando = true;
		erroCEP = '';
		try {
			const res = await fetch(`https://viacep.com.br/ws/${cep}/json/`);
			if (!res.ok) throw new Error('CEP não encontrado');
			const data = await res.json();
			if (data.erro) {
				erroCEP = 'CEP não encontrado';
				return;
			}
			if (onPreenchido) {
				onPreenchido({
					logradouro: data.logradouro ?? '',
					bairro: data.bairro ?? '',
					localidade: data.localidade ?? '',
					uf: data.uf ?? ''
				});
			}
		} catch {
			erroCEP = 'Não foi possível buscar o CEP';
		} finally {
			buscando = false;
		}
	}

	let erroExibido = $derived(erroCEP || error);
</script>

<div class="form-group">
	<label for="cep-input" class="form-label">
		CEP
		{#if buscando}
			<span class="badge badge-info" style="font-size:0.7rem; margin-left:0.4rem;">Buscando...</span>
		{/if}
	</label>
	<input
		id="cep-input"
		type="text"
		inputmode="numeric"
		autocomplete="postal-code"
		maxlength="9"
		placeholder="00000-000"
		class="form-input"
		class:input-error={!!erroExibido}
		{value}
		{disabled}
		oninput={handleInput}
	/>
	{#if erroExibido}
		<p class="form-error">{erroExibido}</p>
	{/if}
</div>
