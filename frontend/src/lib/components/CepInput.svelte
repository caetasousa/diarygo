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
	let avisoForaGoiania = $state(false);

	function aplicarMascara(v: string): string {
		const d = v.replace(/\D/g, '').slice(0, 8);
		if (d.length <= 5) return d;
		return `${d.slice(0, 5)}-${d.slice(5)}`;
	}

	function handleInput(e: Event) {
		const input = e.target as HTMLInputElement;
		value = aplicarMascara(input.value);
		erroCEP = '';
		avisoForaGoiania = false;
		const digits = value.replace(/\D/g, '');
		if (digits.length === 8) {
			avisoForaGoiania = !digits.startsWith('74');
			buscarCEP(digits);
		}
	}

	async function buscarViaCEP(cep: string): Promise<DadosViaCEP | null> {
		const controller = new AbortController();
		const timer = setTimeout(() => controller.abort(), 4000);
		try {
			const res = await fetch(`https://viacep.com.br/ws/${cep}/json/`, {
				signal: controller.signal
			});
			if (!res.ok) return null;
			const data = await res.json();
			if (data.erro) return null;
			return {
				logradouro: data.logradouro ?? '',
				bairro: data.bairro ?? '',
				localidade: data.localidade ?? '',
				uf: data.uf ?? ''
			};
		} catch {
			return null;
		} finally {
			clearTimeout(timer);
		}
	}

	async function buscarBrasilAPI(cep: string): Promise<DadosViaCEP | null> {
		const controller = new AbortController();
		const timer = setTimeout(() => controller.abort(), 4000);
		try {
			const res = await fetch(`https://brasilapi.com.br/api/cep/v2/${cep}`, {
				signal: controller.signal
			});
			if (!res.ok) return null;
			const data = await res.json();
			return {
				logradouro: data.street ?? '',
				bairro: data.neighborhood ?? '',
				localidade: data.city ?? '',
				uf: data.state ?? ''
			};
		} catch {
			return null;
		} finally {
			clearTimeout(timer);
		}
	}

	async function buscarCEP(cep: string) {
		buscando = true;
		erroCEP = '';
		try {
			const dados = (await buscarViaCEP(cep)) ?? (await buscarBrasilAPI(cep));
			if (!dados) {
				erroCEP = 'CEP não encontrado';
				return;
			}
			onPreenchido?.(dados);
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
	{:else if avisoForaGoiania}
		<p class="form-warning">
			Este CEP não parece ser de Goiânia. O cadastro só é aceito para endereços em Goiânia/GO.
		</p>
	{/if}
</div>

<style>
	.form-warning {
		margin-top: 0.35rem;
		font-size: 0.8rem;
		color: var(--color-orange-10, #ff801f);
	}
</style>
