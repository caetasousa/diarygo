<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import CpfInput from '$lib/components/CpfInput.svelte';
	import TelefoneInput from '$lib/components/TelefoneInput.svelte';
	import type { ClienteResponse, ProfissionalResponse } from '$lib/types';

	let tipo = $derived($auth.payload?.tipo ?? null);
	let carregando = $state(true);
	let salvando = $state(false);

	// Cliente
	let perfilCliente = $state<ClienteResponse | null>(null);
	let nomeCliente = $state('');
	let cpfCliente = $state('');
	let telefoneCliente = $state('');
	let erros = $state<Record<string, string>>({});

	// Profissional
	let perfilProf = $state<ProfissionalResponse | null>(null);
	let nomeProf = $state('');
	let cpfProf = $state('');
	let rgProf = $state('');
	let telefoneProf = $state('');
	let fotoUrlProf = $state('');
	let meiProf = $state(false);

	onMount(async () => {
		try {
			if (tipo === 'CLIENTE') {
				perfilCliente = await api.buscarPerfilCliente().catch(() => null);
				if (perfilCliente) {
					nomeCliente = perfilCliente.nome;
					cpfCliente = perfilCliente.cpf;
					telefoneCliente = perfilCliente.telefone;
				}
			} else if (tipo === 'PROFISSIONAL') {
				perfilProf = await api.buscarPerfilProfissional().catch(() => null);
				if (perfilProf) {
					nomeProf = perfilProf.nome;
					cpfProf = perfilProf.cpf;
					rgProf = perfilProf.rg;
					telefoneProf = perfilProf.telefone;
					fotoUrlProf = perfilProf.foto_url;
					meiProf = perfilProf.mei;
				}
			}
		} finally {
			carregando = false;
		}
	});

	function validarCliente(): boolean {
		erros = {};
		if (!nomeCliente.trim()) erros.nome = 'Nome é obrigatório';
		const cpfDigits = cpfCliente.replace(/\D/g, '');
		if (cpfDigits.length !== 11) erros.cpf = 'CPF deve ter 11 dígitos';
		const telDigits = telefoneCliente.replace(/\D/g, '');
		if (telDigits.length < 10) erros.telefone = 'Telefone inválido';
		return Object.keys(erros).length === 0;
	}

	function validarProfissional(): boolean {
		erros = {};
		if (!nomeProf.trim()) erros.nome = 'Nome é obrigatório';
		const cpfDigits = cpfProf.replace(/\D/g, '');
		if (cpfDigits.length !== 11) erros.cpf = 'CPF deve ter 11 dígitos';
		const telDigits = telefoneProf.replace(/\D/g, '');
		if (telDigits.length < 10) erros.telefone = 'Telefone inválido';
		return Object.keys(erros).length === 0;
	}

	async function salvarCliente() {
		if (!validarCliente()) return;
		salvando = true;
		try {
			const req = {
				nome: nomeCliente.trim(),
				cpf: cpfCliente.replace(/\D/g, ''),
				telefone: telefoneCliente.replace(/\D/g, '')
			};
			if (perfilCliente) {
				perfilCliente = await api.atualizarPerfilCliente(req);
			} else {
				perfilCliente = await api.criarPerfilCliente(req);
			}
			toasts.success('Perfil salvo com sucesso!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar perfil');
		} finally {
			salvando = false;
		}
	}

	async function salvarProfissional() {
		if (!validarProfissional()) return;
		salvando = true;
		try {
			const req = {
				nome: nomeProf.trim(),
				cpf: cpfProf.replace(/\D/g, ''),
				rg: rgProf.trim(),
				telefone: telefoneProf.replace(/\D/g, ''),
				foto_url: fotoUrlProf.trim(),
				mei: meiProf
			};
			if (perfilProf) {
				perfilProf = await api.atualizarPerfilProfissional(req);
			} else {
				perfilProf = await api.criarPerfilProfissional(req);
			}
			toasts.success('Perfil salvo com sucesso!');
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar perfil');
		} finally {
			salvando = false;
		}
	}

	const statusLabel: Record<string, string> = {
		PENDENTE: 'Pendente de aprovação',
		APROVADA: 'Aprovada',
		REPROVADA: 'Reprovada',
		SUSPENSA: 'Suspensa',
		DESCREDENCIADA: 'Descredenciada'
	};
	const statusBadge: Record<string, string> = {
		PENDENTE: 'badge-warning',
		APROVADA: 'badge-success',
		REPROVADA: 'badge-danger',
		SUSPENSA: 'badge-danger',
		DESCREDENCIADA: 'badge-danger'
	};
</script>

<div class="container" style="max-width: 640px; padding: 2rem 1rem;">
	<h1 class="text-2xl font-bold" style="margin-bottom: 0.5rem;">Meu Perfil</h1>
	<p class="text-muted" style="margin-bottom: 2rem;">Complete seus dados para usar a plataforma.</p>

	{#if carregando}
		<p class="text-muted">Carregando...</p>
	{:else if tipo === 'CLIENTE'}
		<form novalidate onsubmit={(e) => { e.preventDefault(); salvarCliente(); }}>
			<div class="form-group">
				<label for="nome-cliente" class="form-label">Nome completo</label>
				<input
					id="nome-cliente"
					type="text"
					autocomplete="name"
					maxlength="100"
					placeholder="Seu nome completo"
					class="form-input"
					class:input-error={!!erros.nome}
					bind:value={nomeCliente}
				/>
				{#if erros.nome}<p class="form-error">{erros.nome}</p>{/if}
			</div>

			<CpfInput bind:value={cpfCliente} error={erros.cpf} />
			<TelefoneInput bind:value={telefoneCliente} error={erros.telefone} />

			<button type="submit" class="btn btn-primary btn-full" disabled={salvando}>
				{salvando ? 'Salvando...' : (perfilCliente ? 'Atualizar perfil' : 'Criar perfil')}
			</button>
		</form>

		{#if perfilCliente}
			<div class="card" style="margin-top: 1.5rem; background: var(--color-bg-secondary);">
				<p style="font-size: 0.85rem; color: var(--color-text-muted);">Score de confiabilidade</p>
				<p style="font-size: 2rem; font-weight: 700; color: var(--color-primary);">{perfilCliente.score}/100</p>
			</div>
		{/if}

	{:else if tipo === 'PROFISSIONAL'}
		{#if perfilProf}
			<div style="margin-bottom: 1.5rem;">
				<span class="badge {statusBadge[perfilProf.status] ?? 'badge-info'}">
					{statusLabel[perfilProf.status] ?? perfilProf.status}
				</span>
				{#if perfilProf.status === 'PENDENTE'}
					<p class="text-muted" style="font-size: 0.85rem; margin-top: 0.5rem;">
						Sua conta está em análise. O prazo é de até 48h úteis.
					</p>
				{/if}
			</div>
		{/if}

		<form novalidate onsubmit={(e) => { e.preventDefault(); salvarProfissional(); }}>
			<div class="form-group">
				<label for="nome-prof" class="form-label">Nome completo</label>
				<input
					id="nome-prof"
					type="text"
					autocomplete="name"
					maxlength="100"
					placeholder="Seu nome completo"
					class="form-input"
					class:input-error={!!erros.nome}
					bind:value={nomeProf}
				/>
				{#if erros.nome}<p class="form-error">{erros.nome}</p>{/if}
			</div>

			<CpfInput bind:value={cpfProf} error={erros.cpf} />

			<div class="form-group">
				<label for="rg-prof" class="form-label">RG</label>
				<input
					id="rg-prof"
					type="text"
					autocomplete="off"
					maxlength="20"
					placeholder="Número do RG"
					class="form-input"
					bind:value={rgProf}
				/>
			</div>

			<TelefoneInput bind:value={telefoneProf} error={erros.telefone} />

			<div class="form-group">
				<label for="foto-url" class="form-label">URL da foto de perfil</label>
				<input
					id="foto-url"
					type="url"
					autocomplete="off"
					maxlength="500"
					placeholder="https://..."
					class="form-input"
					bind:value={fotoUrlProf}
				/>
			</div>

			<div class="form-group">
				<label class="flex gap-sm" style="align-items: center; cursor: pointer;">
					<input type="checkbox" bind:checked={meiProf} style="width: 1rem; height: 1rem;" />
					<span>Possuo MEI (Microempreendedor Individual)</span>
				</label>
			</div>

			<button type="submit" class="btn btn-primary btn-full" disabled={salvando}>
				{salvando ? 'Salvando...' : (perfilProf ? 'Atualizar perfil' : 'Criar perfil')}
			</button>
		</form>
	{:else}
		<p class="text-muted">Tipo de usuário não suportado nesta página.</p>
	{/if}
</div>
