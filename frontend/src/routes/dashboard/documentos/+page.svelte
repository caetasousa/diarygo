<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import type { DocumentoResponse, TipoDocumento } from '$lib/types';

	let documentos = $state<DocumentoResponse[]>([]);
	let carregando = $state(true);
	let enviando = $state(false);

	let tipo = $state<TipoDocumento>('RG_FRENTE');
	let url = $state('');
	let erroURL = $state('');

	const tiposDocumento: { value: TipoDocumento; label: string }[] = [
		{ value: 'RG_FRENTE', label: 'RG – Frente' },
		{ value: 'RG_VERSO', label: 'RG – Verso' },
		{ value: 'CPF', label: 'CPF' },
		{ value: 'COMPROVANTE', label: 'Comprovante de Residência' },
		{ value: 'FOTO', label: 'Foto de Perfil' },
		{ value: 'OUTRO', label: 'Outro' }
	];

	const statusLabel: Record<string, string> = {
		PENDENTE: 'Pendente',
		APROVADO: 'Aprovado',
		REPROVADO: 'Reprovado'
	};
	const statusBadge: Record<string, string> = {
		PENDENTE: 'badge-warning',
		APROVADO: 'badge-success',
		REPROVADO: 'badge-danger'
	};

	onMount(async () => {
		await carregarDocumentos();
	});

	async function carregarDocumentos() {
		carregando = true;
		try {
			documentos = await api.listarDocumentos();
		} catch {
			toasts.error('Erro ao carregar documentos');
		} finally {
			carregando = false;
		}
	}

	async function enviar() {
		erroURL = '';
		if (!url.trim()) { erroURL = 'URL é obrigatória'; return; }
		enviando = true;
		try {
			await api.enviarDocumento({ tipo, url: url.trim() });
			url = '';
			toasts.success('Documento enviado para análise!');
			await carregarDocumentos();
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao enviar documento');
		} finally {
			enviando = false;
		}
	}
</script>

<div class="container" style="max-width: 680px; padding: 2rem 1rem;">
	<h1 class="text-2xl font-bold" style="margin-bottom: 0.5rem;">Documentos</h1>
	<p class="text-muted" style="margin-bottom: 2rem;">
		Envie os documentos obrigatórios para análise. Prazo: até 48h úteis.
	</p>

	<div class="card" style="margin-bottom: 2rem;">
		<h2 class="font-bold" style="margin-bottom: 1rem;">Enviar novo documento</h2>
		<form novalidate onsubmit={(e) => { e.preventDefault(); enviar(); }}>
			<div class="form-group">
				<label for="tipo-doc" class="form-label">Tipo de documento</label>
				<select id="tipo-doc" class="form-input" bind:value={tipo}>
					{#each tiposDocumento as t}
						<option value={t.value}>{t.label}</option>
					{/each}
				</select>
			</div>
			<div class="form-group">
				<label for="url-doc" class="form-label">URL do documento</label>
				<input
					id="url-doc"
					type="url"
					maxlength="500"
					placeholder="https://..."
					class="form-input"
					class:input-error={!!erroURL}
					bind:value={url}
					autocomplete="off"
				/>
				{#if erroURL}<p class="form-error">{erroURL}</p>{/if}
				<p class="text-muted" style="font-size: 0.8rem; margin-top: 0.25rem;">
					No MVP, informe a URL pública do documento (ex: Google Drive, Dropbox).
				</p>
			</div>
			<button type="submit" class="btn btn-primary" disabled={enviando}>
				{enviando ? 'Enviando...' : 'Enviar documento'}
			</button>
		</form>
	</div>

	{#if carregando}
		<p class="text-muted">Carregando documentos...</p>
	{:else if documentos.length === 0}
		<p class="text-muted">Nenhum documento enviado ainda.</p>
	{:else}
		<h2 class="font-bold" style="margin-bottom: 1rem;">Documentos enviados</h2>
		<div style="display: flex; flex-direction: column; gap: 0.75rem;">
			{#each documentos as doc (doc.id)}
				<div class="card" style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.75rem;">
					<div>
						<strong>{tiposDocumento.find(t => t.value === doc.tipo)?.label ?? doc.tipo}</strong>
						<br />
						<a href={doc.url} target="_blank" rel="noopener noreferrer"
							style="font-size: 0.8rem; color: var(--color-primary); word-break: break-all;">
							{doc.url}
						</a>
					</div>
					<span class="badge {statusBadge[doc.status] ?? 'badge-info'}">
						{statusLabel[doc.status] ?? doc.status}
					</span>
				</div>
			{/each}
		</div>
	{/if}
</div>
