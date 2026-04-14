<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type { RegiaoResponse } from '$lib/types';

	let todasRegioes = $state<RegiaoResponse[]>([]);
	let regioesAtuacao = $state<RegiaoResponse[]>([]);
	let carregando = $state(true);
	let salvando = $state(false);

	let selecionadas = $state<Set<string>>(new Set());

	onMount(async () => {
		try {
			const [todas, atuacao] = await Promise.all([
				api.listarRegioes(),
				api.listarRegioesAtuacao().catch(() => [])
			]);
			todasRegioes = todas;
			regioesAtuacao = atuacao;
			selecionadas = new Set(atuacao.map((r) => r.id));
		} catch {
			toasts.error('Erro ao carregar regiões');
		} finally {
			carregando = false;
		}
	});

	function toggle(id: string) {
		const novas = new Set(selecionadas);
		if (novas.has(id)) novas.delete(id);
		else novas.add(id);
		selecionadas = novas;
	}

	async function salvar() {
		salvando = true;
		try {
			regioesAtuacao = await api.definirRegioesAtuacao({ regiao_ids: [...selecionadas] });
			selecionadas = new Set(regioesAtuacao.map((r) => r.id));
			toasts.success('Regiões de atuação salvas!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar regiões');
		} finally {
			salvando = false;
		}
	}
</script>

<div class="container" style="max-width: 640px; padding: 2rem 1rem;">
	<h1 class="text-2xl font-bold" style="margin-bottom: 0.5rem;">Regiões de Atuação</h1>
	<p class="text-muted" style="margin-bottom: 2rem;">
		Selecione as regiões onde você está disponível para trabalhar.
	</p>

	{#if carregando}
		<p class="text-muted">Carregando regiões...</p>
	{:else if todasRegioes.length === 0}
		<p class="text-muted">Nenhuma região disponível no momento.</p>
	{:else}
		<div style="display: flex; flex-direction: column; gap: 0.75rem; margin-bottom: 2rem;">
			{#each todasRegioes as regiao (regiao.id)}
				<label class="card" style="cursor: pointer; display: flex; align-items: center; gap: 1rem; padding: 1rem;">
					<input
						type="checkbox"
						style="width: 1.1rem; height: 1.1rem; flex-shrink: 0;"
						checked={selecionadas.has(regiao.id)}
						onchange={() => toggle(regiao.id)}
					/>
					<div>
						<strong>{regiao.nome}</strong>
						<br />
						<span class="text-muted" style="font-size: 0.85rem;">
							{regiao.cidade}/{regiao.estado} &middot; CEP {regiao.cep_inicio} – {regiao.cep_fim}
						</span>
					</div>
				</label>
			{/each}
		</div>

		<button class="btn btn-primary btn-full" onclick={salvar} disabled={salvando}>
			{salvando ? 'Salvando...' : 'Salvar regiões selecionadas'}
		</button>
	{/if}
</div>
