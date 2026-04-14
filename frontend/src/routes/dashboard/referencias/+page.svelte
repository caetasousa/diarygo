<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import TelefoneInput from '$lib/components/TelefoneInput.svelte';
	import type { ReferenciaResponse } from '$lib/types';

	let referencias = $state<ReferenciaResponse[]>([]);
	let carregando = $state(true);
	let adicionando = $state(false);

	let nomeContato = $state('');
	let telefoneContato = $state('');
	let erros = $state<Record<string, string>>({});

	const statusLabel: Record<string, string> = {
		PENDENTE: 'Pendente',
		CONFIRMADA: 'Confirmada',
		NAO_CONFIRMADA: 'Não confirmada'
	};
	const statusBadge: Record<string, string> = {
		PENDENTE: 'badge-warning',
		CONFIRMADA: 'badge-success',
		NAO_CONFIRMADA: 'badge-danger'
	};

	onMount(async () => {
		await carregarReferencias();
	});

	async function carregarReferencias() {
		carregando = true;
		try {
			referencias = await api.listarReferencias();
		} catch {
			toasts.error('Erro ao carregar referências');
		} finally {
			carregando = false;
		}
	}

	function validar(): boolean {
		erros = {};
		if (!nomeContato.trim()) erros.nome = 'Nome é obrigatório';
		const digits = telefoneContato.replace(/\D/g, '');
		if (digits.length < 10) erros.telefone = 'Telefone inválido';
		return Object.keys(erros).length === 0;
	}

	async function adicionar() {
		if (!validar()) return;
		adicionando = true;
		try {
			await api.adicionarReferencia({
				nome_contato: nomeContato.trim(),
				telefone_contato: telefoneContato.replace(/\D/g, '')
			});
			nomeContato = '';
			telefoneContato = '';
			erros = {};
			toasts.success('Referência adicionada!');
			await carregarReferencias();
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao adicionar referência');
		} finally {
			adicionando = false;
		}
	}
</script>

<div class="container" style="max-width: 640px; padding: 2rem 1rem;">
	<h1 class="text-2xl font-bold" style="margin-bottom: 0.5rem;">Referências Profissionais</h1>
	<p class="text-muted" style="margin-bottom: 2rem;">
		Adicione contatos que possam atestar sua experiência como diarista.
	</p>

	<div class="card" style="margin-bottom: 2rem;">
		<h2 class="font-bold" style="margin-bottom: 1rem;">Adicionar referência</h2>
		<form novalidate onsubmit={(e) => { e.preventDefault(); adicionar(); }}>
			<div class="form-group">
				<label for="nome-contato" class="form-label">Nome do contato</label>
				<input
					id="nome-contato"
					type="text"
					maxlength="100"
					placeholder="Nome completo"
					class="form-input"
					class:input-error={!!erros.nome}
					bind:value={nomeContato}
					autocomplete="off"
				/>
				{#if erros.nome}<p class="form-error">{erros.nome}</p>{/if}
			</div>

			<TelefoneInput bind:value={telefoneContato} error={erros.telefone} />

			<button type="submit" class="btn btn-primary" disabled={adicionando}>
				{adicionando ? 'Adicionando...' : 'Adicionar referência'}
			</button>
		</form>
	</div>

	{#if carregando}
		<p class="text-muted">Carregando referências...</p>
	{:else if referencias.length === 0}
		<p class="text-muted">Nenhuma referência adicionada ainda.</p>
	{:else}
		<h2 class="font-bold" style="margin-bottom: 1rem;">Referências cadastradas</h2>
		<div style="display: flex; flex-direction: column; gap: 0.75rem;">
			{#each referencias as ref (ref.id)}
				<div class="card" style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.75rem;">
					<div>
						<strong>{ref.nome_contato}</strong>
						<br />
						<span class="text-muted" style="font-size: 0.85rem;">{ref.telefone_contato}</span>
					</div>
					<span class="badge {statusBadge[ref.status] ?? 'badge-info'}">
						{statusLabel[ref.status] ?? ref.status}
					</span>
				</div>
			{/each}
		</div>
	{/if}
</div>
