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
		PENDENTE: 'Pendente', APROVADO: 'Aprovado', REPROVADO: 'Reprovado'
	};
	const statusBadge: Record<string, string> = {
		PENDENTE: 'badge-orange', APROVADO: 'badge-green', REPROVADO: 'badge-red'
	};

	onMount(carregarDocumentos);

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

	let pendentes = $derived(documentos.filter(d => d.status === 'PENDENTE').length);
	let aprovados = $derived(documentos.filter(d => d.status === 'APROVADO').length);
</script>

<svelte:head>
	<title>Documentos — DiaryGo</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<div>
			<h1 class="page-title">Documentos</h1>
			<p class="page-sub">Envie os documentos obrigatórios para análise.</p>
		</div>
		{#if documentos.length > 0}
			<div class="doc-stats">
				<span class="doc-stat"><span class="doc-stat-num">{aprovados}</span> aprovados</span>
				<span class="doc-stat-sep">·</span>
				<span class="doc-stat"><span class="doc-stat-num">{pendentes}</span> pendentes</span>
			</div>
		{/if}
	</div>

	<div class="layout">
		<!-- Formulário -->
		<div class="form-card">
			<h2 class="card-title">Enviar documento</h2>
			<form novalidate onsubmit={(e) => { e.preventDefault(); enviar(); }}>
				<div class="form-group">
					<label for="tipo-doc" class="form-label">Tipo</label>
					<select id="tipo-doc" class="form-input" bind:value={tipo}>
						{#each tiposDocumento as t}
							<option value={t.value}>{t.label}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label for="url-doc" class="form-label">URL do documento</label>
					<input
						id="url-doc" type="url" maxlength="500"
						placeholder="https://drive.google.com/..."
						class="form-input" class:input-error={!!erroURL}
						bind:value={url} autocomplete="off"
					/>
					{#if erroURL}<p class="form-error">{erroURL}</p>{/if}
					<p class="form-hint">Google Drive, Dropbox ou outro link público.</p>
				</div>
				<button type="submit" class="btn btn-white btn-full" disabled={enviando}>
					{enviando ? 'Enviando...' : 'Enviar para análise'}
				</button>
			</form>
		</div>

		<!-- Lista -->
		<div class="list-col">
			{#if carregando}
				<div class="loading-state">
					<div class="loading-spinner"></div>
					<span>Carregando...</span>
				</div>
			{:else if documentos.length === 0}
				<div class="empty-state">
					<p class="empty-title">Nenhum documento enviado</p>
					<p class="empty-sub">Use o formulário ao lado para enviar sua documentação.</p>
				</div>
			{:else}
				<div class="doc-list">
					{#each documentos as doc (doc.id)}
						<div class="doc-row">
							<div class="doc-info">
								<span class="doc-tipo">{tiposDocumento.find(t => t.value === doc.tipo)?.label ?? doc.tipo}</span>
								<a href={doc.url} target="_blank" rel="noopener noreferrer" class="doc-url">
									{doc.url.length > 50 ? doc.url.slice(0, 50) + '...' : doc.url}
								</a>
							</div>
							<span class="badge {statusBadge[doc.status] ?? 'badge-blue'}">
								{statusLabel[doc.status] ?? doc.status}
							</span>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.page { padding: 40px; max-width: 900px; margin: 0 auto; width: 100%; }

	.page-header {
		display: flex; align-items: flex-start;
		justify-content: space-between; gap: 16px;
		margin-bottom: 32px; flex-wrap: wrap;
	}

	.page-title {
		font-size: 1.75rem; font-weight: 400;
		letter-spacing: -1px; color: var(--color-text-primary);
		line-height: 1; margin-bottom: 6px;
	}

	.page-sub { font-size: 0.875rem; color: var(--color-text-tertiary); }

	.doc-stats {
		display: flex; align-items: center; gap: 8px;
		font-size: 0.8125rem; color: var(--color-text-tertiary);
	}
	.doc-stat-num { color: var(--color-text-primary); font-weight: 500; }
	.doc-stat-sep { color: var(--color-text-tertiary); }

	.layout {
		display: grid;
		grid-template-columns: 320px 1fr;
		gap: 20px;
		align-items: start;
	}

	.form-card {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 12px;
		padding: 20px;
		box-shadow: var(--shadow-ring);
	}

	.card-title {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		margin-bottom: 16px;
		padding-bottom: 12px;
		border-bottom: 1px solid var(--border-frost);
	}

	.form-hint {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		margin-top: 4px;
		line-height: 1.4;
	}

	.list-col { min-width: 0; }

	.loading-state {
		display: flex; align-items: center; gap: 10px;
		color: var(--color-text-tertiary); font-size: 0.875rem; padding: 32px 0;
	}

	.loading-spinner {
		width: 16px; height: 16px;
		border: 2px solid var(--border-frost);
		border-top-color: var(--color-text-secondary);
		border-radius: 50%;
		animation: spin 0.7s linear infinite;
	}

	@keyframes spin { to { transform: rotate(360deg); } }

	.empty-state {
		padding: 40px 24px; text-align: center;
		background: var(--bg-card); border: 1px solid var(--border-frost);
		border-radius: 12px; box-shadow: var(--shadow-ring);
	}
	.empty-title { font-size: 0.9375rem; font-weight: 500; color: var(--color-text-primary); margin-bottom: 6px; }
	.empty-sub { font-size: 0.8125rem; color: var(--color-text-tertiary); }

	.doc-list {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: 12px;
		overflow: hidden;
		box-shadow: var(--shadow-ring);
	}

	.doc-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		padding: 14px 18px;
		border-bottom: 1px solid var(--border-frost);
		flex-wrap: wrap;
	}

	.doc-row:last-child { border-bottom: none; }

	.doc-info {
		display: flex; flex-direction: column; gap: 3px; min-width: 0;
	}

	.doc-tipo {
		font-size: 0.875rem; font-weight: 500; color: var(--color-text-primary);
	}

	.doc-url {
		font-size: 0.75rem; color: var(--color-blue-10);
		white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
		max-width: 300px;
	}

	@media (max-width: 768px) {
		.page { padding: 24px 16px; }
		.layout { grid-template-columns: 1fr; }
	}
</style>
