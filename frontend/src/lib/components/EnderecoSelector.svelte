<script lang="ts">
	import { MapPin, Plus } from 'lucide-svelte';
	import type { EnderecoResponse } from '$lib/types';

	interface Props {
		enderecos: EnderecoResponse[];
		selectedId: string;
		onNovo: () => void;
	}

	let { enderecos, selectedId = $bindable(), onNovo }: Props = $props();
</script>

<div class="selector">
	{#if enderecos.length === 0}
		<div class="empty">
			<MapPin size={24} />
			<p>Você ainda não tem endereços cadastrados.</p>
			<button type="button" class="btn-add" onclick={onNovo}>
				<Plus size={14} /> Adicionar endereço
			</button>
		</div>
	{:else}
		<div class="cards">
			{#each enderecos as end (end.id)}
				<label class="card" class:selected={selectedId === end.id}>
					<input
						type="radio"
						name="endereco"
						value={end.id}
						bind:group={selectedId}
						class="sr-only"
					/>
					<div class="card-radio" aria-hidden="true">
						{#if selectedId === end.id}
							<div class="dot"></div>
						{/if}
					</div>
					<div class="card-body">
						<span class="card-rua">{end.logradouro}, {end.numero}</span>
						{#if end.complemento}
							<span class="card-extra">{end.complemento}</span>
						{/if}
						<span class="card-extra">{end.bairro} — {end.cidade}/{end.estado} · {end.cep.replace(/(\d{5})(\d{3})/, '$1-$2')}</span>
					</div>
					{#if end.principal}
						<span class="badge-principal">Principal</span>
					{/if}
				</label>
			{/each}

			<button type="button" class="card card-add" onclick={onNovo}>
				<Plus size={16} />
				<span>Adicionar novo endereço</span>
			</button>
		</div>
	{/if}
</div>

<style>
	.selector { display: flex; flex-direction: column; gap: 12px; }

	.empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12px;
		padding: 32px 16px;
		text-align: center;
		color: var(--color-text-tertiary);
		border: 1px dashed var(--border-frost);
		border-radius: var(--radius-card);
	}

	.btn-add {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 8px 16px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		background: var(--bg-card);
		color: var(--color-text-primary);
		font-size: 0.875rem;
		cursor: pointer;
		transition: border-color 0.12s, background 0.12s;
	}
	.btn-add:hover { border-color: var(--border-frost-hover); background: var(--bg-hover-subtle); }

	.cards { display: flex; flex-direction: column; gap: 8px; }

	.card {
		display: flex;
		align-items: flex-start;
		gap: 12px;
		padding: 14px 16px;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		background: var(--bg-card);
		cursor: pointer;
		transition: border-color 0.12s, background 0.12s;
		position: relative;
	}
	.card:hover { border-color: var(--border-frost-hover); background: var(--bg-hover-subtle); }
	.card.selected { border-color: var(--color-orange-10); background: rgba(255, 128, 31, 0.04); }

	.card-radio {
		width: 18px;
		height: 18px;
		border-radius: 50%;
		border: 2px solid var(--border-frost);
		flex-shrink: 0;
		margin-top: 2px;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: border-color 0.12s;
	}
	.card.selected .card-radio { border-color: var(--color-orange-10); }
	.dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: var(--color-orange-10);
	}

	.card-body { display: flex; flex-direction: column; gap: 2px; flex: 1; }
	.card-rua { font-size: 0.875rem; color: var(--color-text-primary); font-weight: 500; }
	.card-extra { font-size: 0.75rem; color: var(--color-text-tertiary); }

	.badge-principal {
		font-size: 0.625rem;
		font-weight: 600;
		letter-spacing: 0.5px;
		text-transform: uppercase;
		padding: 3px 8px;
		border-radius: var(--radius-pill);
		border: 1px solid var(--color-green-4);
		color: var(--color-green-10);
		background: transparent;
	}

	.card-add {
		justify-content: center;
		gap: 8px;
		color: var(--color-text-tertiary);
		font-size: 0.875rem;
		border-style: dashed;
	}
	.card-add:hover { color: var(--color-text-primary); }

	.sr-only {
		position: absolute;
		width: 1px; height: 1px;
		padding: 0; margin: -1px;
		overflow: hidden;
		clip: rect(0,0,0,0);
		white-space: nowrap;
		border-width: 0;
	}
</style>
