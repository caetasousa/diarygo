<script lang="ts">
	import type { EnderecoResponse } from '$lib/types';

	interface Props {
		endereco: EnderecoResponse;
		onDefinirPrincipal?: (id: string) => void;
		onEditar?: (id: string) => void;
		onRemover?: (id: string) => void;
	}

	let { endereco, onDefinirPrincipal, onEditar, onRemover }: Props = $props();

	function formatarCEP(cep: string): string {
		return cep.length === 8 ? `${cep.slice(0, 5)}-${cep.slice(5)}` : cep;
	}
</script>

<div class="card" style="position: relative;">
	{#if endereco.principal}
		<span class="badge badge-success" style="position: absolute; top: 1rem; right: 1rem;">
			Principal
		</span>
	{/if}

	<div style="margin-bottom: 0.75rem;">
		<strong>{endereco.logradouro}, {endereco.numero}</strong>
		{#if endereco.complemento}
			<span style="color: var(--color-text-muted);">– {endereco.complemento}</span>
		{/if}
	</div>
	<div style="color: var(--color-text-muted); font-size: 0.9rem; margin-bottom: 0.5rem;">
		{endereco.bairro} – {endereco.cidade}/{endereco.estado} &middot; CEP {formatarCEP(endereco.cep)}
	</div>
	<div style="font-size: 0.85rem; color: var(--color-text-muted);">
		{endereco.num_quartos} quarto(s) &middot; {endereco.num_banheiros} banheiro(s) &middot;
		{endereco.num_salas} sala(s) &middot; {endereco.num_cozinhas} cozinha(s)
		{#if endereco.area_m2 > 0}
			&middot; {endereco.area_m2} m²
		{/if}
	</div>

	<div class="flex gap-sm" style="margin-top: 1rem; flex-wrap: wrap;">
		{#if !endereco.principal && onDefinirPrincipal}
			<button
				class="btn btn-secondary btn-sm"
				onclick={() => onDefinirPrincipal!(endereco.id)}
			>
				Definir como principal
			</button>
		{/if}
		{#if onEditar}
			<button class="btn btn-secondary btn-sm" onclick={() => onEditar!(endereco.id)}>
				Editar
			</button>
		{/if}
		{#if onRemover}
			<button class="btn btn-danger btn-sm" onclick={() => onRemover!(endereco.id)}>
				Remover
			</button>
		{/if}
	</div>
</div>
