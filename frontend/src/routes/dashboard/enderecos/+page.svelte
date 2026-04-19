<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/stores/toasts';
	import { api } from '$lib/api/client';
	import EnderecoCard from '$lib/components/EnderecoCard.svelte';
	import CepInput from '$lib/components/CepInput.svelte';
	import type { EnderecoResponse } from '$lib/types';
	import type { DadosViaCEP } from '$lib/components/CepInput.svelte';

	let enderecos = $state<EnderecoResponse[]>([]);
	let carregando = $state(true);
	let mostrarFormulario = $state(false);
	let editandoID = $state<string | null>(null);
	let salvando = $state(false);

	let cep = $state('');
	let logradouro = $state('');
	let numero = $state('');
	let complemento = $state('');
	let bairro = $state('');
	let cidade = $state('');
	let estado = $state('');
	let numQuartos = $state<number | null>(null);
	let numBanheiros = $state<number | null>(null);
	let numSalas = $state<number | null>(null);
	let numCozinhas = $state<number | null>(null);
	let areaM2 = $state<number | null>(null);
	let principal = $state(false);
	let erros = $state<Record<string, string>>({});

	onMount(carregarEnderecos);

	async function carregarEnderecos() {
		carregando = true;
		try {
			enderecos = await api.listarEnderecos();
		} catch {
			toasts.error('Erro ao carregar endereços');
		} finally {
			carregando = false;
		}
	}

	function preencherViaCEP(dados: DadosViaCEP) {
		logradouro = dados.logradouro;
		bairro = dados.bairro;
		cidade = dados.localidade;
		estado = dados.uf;
	}

	function abrirNovo() {
		resetar();
		editandoID = null;
		mostrarFormulario = true;
	}

	function abrirEditar(id: string) {
		const e = enderecos.find((en) => en.id === id);
		if (!e) return;
		cep = e.cep; logradouro = e.logradouro; numero = e.numero;
		complemento = e.complemento; bairro = e.bairro;
		cidade = e.cidade; estado = e.estado;
		numQuartos = e.num_quartos; numBanheiros = e.num_banheiros;
		numSalas = e.num_salas; numCozinhas = e.num_cozinhas;
		areaM2 = e.area_m2; principal = e.principal;
		editandoID = id; erros = {};
		mostrarFormulario = true;
	}

	function resetar() {
		cep = ''; logradouro = ''; numero = ''; complemento = '';
		bairro = ''; cidade = 'Goiânia'; estado = 'GO';
		numQuartos = null; numBanheiros = null; numSalas = null; numCozinhas = null;
		areaM2 = null; principal = false; erros = {};
	}

	function cancelar() {
		mostrarFormulario = false;
		resetar();
	}

	function validar(): boolean {
		erros = {};
		if (cep.replace(/\D/g, '').length !== 8) erros.cep = 'CEP inválido';
		if (!logradouro.trim()) erros.logradouro = 'Logradouro é obrigatório';
		if (!numero.trim()) erros.numero = 'Número é obrigatório';
		if (!bairro.trim()) erros.bairro = 'Bairro é obrigatório';
		if (!cidade.trim()) erros.cidade = 'Cidade é obrigatória';
		if (estado.length !== 2) erros.estado = 'Estado inválido';
		if (numQuartos == null || numQuartos < 1) erros.num_quartos = 'Mínimo 1 quarto';
		if (numBanheiros != null && numBanheiros < 0) erros.num_banheiros = 'Inválido';
		if (numSalas != null && numSalas < 0) erros.num_salas = 'Inválido';
		if (numCozinhas != null && numCozinhas < 0) erros.num_cozinhas = 'Inválido';
		if (areaM2 != null && areaM2 < 0) erros.area_m2 = 'Inválido';
		return Object.keys(erros).length === 0;
	}

	async function salvar() {
		if (!validar()) return;
		salvando = true;
		try {
			const req = {
				logradouro: logradouro.trim(), numero: numero.trim(),
				complemento: complemento.trim(), bairro: bairro.trim(),
				cidade: cidade.trim(), estado: estado.toUpperCase().trim(),
				cep: cep.replace(/\D/g, ''),
				num_quartos: numQuartos ?? 0,
				num_banheiros: numBanheiros ?? 0,
				num_salas: numSalas ?? 0,
				num_cozinhas: numCozinhas ?? 0,
				area_m2: areaM2 ?? 0,
				principal
			};
			if (editandoID) {
				await api.atualizarEndereco(editandoID, req);
				toasts.success('Endereço atualizado!');
			} else {
				await api.criarEndereco(req);
				toasts.success('Endereço adicionado!');
			}
			mostrarFormulario = false;
			resetar();
			await carregarEnderecos();
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao salvar endereço');
		} finally {
			salvando = false;
		}
	}

	async function remover(id: string) {
		if (!confirm('Remover este endereço?')) return;
		try {
			await api.removerEndereco(id);
			toasts.success('Endereço removido');
			await carregarEnderecos();
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao remover endereço');
		}
	}

	async function definirPrincipal(id: string) {
		try {
			await api.definirEnderecoPrincipal(id);
			toasts.success('Endereço principal definido');
			await carregarEnderecos();
		} catch (e) {
			toasts.error(e instanceof Error ? e.message : 'Erro ao definir principal');
		}
	}

	const ufs = ['AC','AL','AP','AM','BA','CE','DF','ES','GO','MA','MT','MS','MG','PA','PB','PR','PE','PI','RJ','RN','RS','RO','RR','SC','SP','SE','TO'];
</script>

<svelte:head>
	<title>Endereços — DiaryGo</title>
</svelte:head>

<div class="page">
	<div class="page-header">
		<div>
			<h1 class="page-title">Endereços</h1>
			<p class="page-sub">Locais onde o serviço será realizado. Atendemos apenas em Goiânia/GO.</p>
		</div>
		{#if !mostrarFormulario}
			<button class="btn btn-white" onclick={abrirNovo}>
				<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
				Novo endereço
			</button>
		{/if}
	</div>

	{#if mostrarFormulario}
		<div class="form-card">
			<div class="form-card-head">
				<h2 class="form-card-title">{editandoID ? 'Editar endereço' : 'Novo endereço'}</h2>
				<button class="close-btn" onclick={cancelar} aria-label="Fechar">
					<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
				</button>
			</div>

			<form novalidate onsubmit={(e) => { e.preventDefault(); salvar(); }}>
				<div class="form-cols">
					<div style="grid-column: 1 / -1;">
						<CepInput bind:value={cep} error={erros.cep} onPreenchido={preencherViaCEP} />
					</div>

					<div class="form-group" style="grid-column: 1 / -1;">
						<label for="logradouro" class="form-label">Logradouro</label>
						<input id="logradouro" type="text" maxlength="200" class="form-input"
							class:input-error={!!erros.logradouro} bind:value={logradouro}
							placeholder="Rua, Avenida..." autocomplete="street-address" />
						{#if erros.logradouro}<p class="form-error">{erros.logradouro}</p>{/if}
					</div>

					<div class="form-group">
						<label for="numero" class="form-label">Número</label>
						<input id="numero" type="text" maxlength="20" class="form-input"
							class:input-error={!!erros.numero} bind:value={numero} placeholder="123" />
						{#if erros.numero}<p class="form-error">{erros.numero}</p>{/if}
					</div>

					<div class="form-group">
						<label for="complemento" class="form-label">Complemento</label>
						<input id="complemento" type="text" maxlength="100" class="form-input"
							bind:value={complemento} placeholder="Apto, Bloco..." />
					</div>

					<div class="form-group">
						<label for="bairro" class="form-label">Bairro</label>
						<input id="bairro" type="text" maxlength="100" class="form-input"
							class:input-error={!!erros.bairro} bind:value={bairro} placeholder="Bairro" />
						{#if erros.bairro}<p class="form-error">{erros.bairro}</p>{/if}
					</div>

					<div class="form-group">
						<label for="cidade" class="form-label">Cidade</label>
						<input id="cidade" type="text" maxlength="100" class="form-input"
							class:input-error={!!erros.cidade} bind:value={cidade}
							placeholder="Cidade" autocomplete="address-level2" />
						{#if erros.cidade}<p class="form-error">{erros.cidade}</p>{/if}
					</div>

					<div class="form-group">
						<label for="estado" class="form-label">Estado (UF)</label>
						<select id="estado" class="form-input" class:input-error={!!erros.estado} bind:value={estado}>
							<option value="">Selecione</option>
							{#each ufs as uf}<option value={uf}>{uf}</option>{/each}
						</select>
						{#if erros.estado}<p class="form-error">{erros.estado}</p>{/if}
					</div>
				</div>

				<div class="form-divider">
					<span class="form-divider-label">Detalhes do imóvel</span>
				</div>

				<div class="form-cols">
					<div class="form-group">
						<label for="quartos" class="form-label">Quartos</label>
						<input id="quartos" type="number" min="1" max="20" class="form-input"
							class:input-error={!!erros.num_quartos} bind:value={numQuartos} />
						{#if erros.num_quartos}<p class="form-error">{erros.num_quartos}</p>{/if}
					</div>
					<div class="form-group">
						<label for="banheiros" class="form-label">Banheiros</label>
						<input id="banheiros" type="number" min="0" max="20" class="form-input"
							class:input-error={!!erros.num_banheiros} bind:value={numBanheiros} />
						{#if erros.num_banheiros}<p class="form-error">{erros.num_banheiros}</p>{/if}
					</div>
					<div class="form-group">
						<label for="salas" class="form-label">Salas</label>
						<input id="salas" type="number" min="0" max="20" class="form-input"
							class:input-error={!!erros.num_salas} bind:value={numSalas} />
						{#if erros.num_salas}<p class="form-error">{erros.num_salas}</p>{/if}
					</div>
					<div class="form-group">
						<label for="cozinhas" class="form-label">Cozinhas</label>
						<input id="cozinhas" type="number" min="0" max="10" class="form-input"
							class:input-error={!!erros.num_cozinhas} bind:value={numCozinhas} />
						{#if erros.num_cozinhas}<p class="form-error">{erros.num_cozinhas}</p>{/if}
					</div>
					<div class="form-group">
						<label for="area" class="form-label">Área (m²)</label>
						<input id="area" type="number" min="0" step="0.1" class="form-input"
							class:input-error={!!erros.area_m2} bind:value={areaM2} />
						{#if erros.area_m2}<p class="form-error">{erros.area_m2}</p>{/if}
					</div>
					<div class="form-group" style="display:flex; align-items:flex-end;">
						<div class="checkbox-row">
							<input type="checkbox" id="principal" bind:checked={principal} />
							<label for="principal" class="checkbox-label">Endereço principal</label>
						</div>
					</div>
				</div>

				<div class="form-actions">
					<button type="submit" class="btn btn-white" disabled={salvando}>
						{salvando ? 'Salvando...' : 'Salvar endereço'}
					</button>
					<button type="button" class="btn btn-primary" onclick={cancelar}>Cancelar</button>
				</div>
			</form>
		</div>
	{/if}

	{#if carregando}
		<div class="loading-state">
			<div class="loading-spinner"></div>
			<span>Carregando endereços...</span>
		</div>
	{:else if enderecos.length === 0 && !mostrarFormulario}
		<div class="empty-state">
			<div class="empty-icon">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
			</div>
			<p class="empty-title">Nenhum endereço cadastrado</p>
			<p class="empty-sub">Adicione o endereço onde os serviços serão realizados.</p>
			<button class="btn btn-white" onclick={abrirNovo}>Adicionar endereço</button>
		</div>
	{:else}
		<div class="list">
			{#each enderecos as endereco (endereco.id)}
				<EnderecoCard
					{endereco}
					onDefinirPrincipal={definirPrincipal}
					onEditar={abrirEditar}
					onRemover={remover}
				/>
			{/each}
		</div>
	{/if}
</div>

<style>
	.page { padding: 40px; max-width: 760px; margin: 0 auto; width: 100%; }

	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 16px;
		margin-bottom: 32px;
		flex-wrap: wrap;
	}

	.page-title {
		font-size: 1.75rem;
		font-weight: 400;
		letter-spacing: -1px;
		color: var(--color-text-primary);
		line-height: 1;
		margin-bottom: 6px;
	}

	.page-sub { font-size: 0.875rem; color: var(--color-text-tertiary); }

	.form-card {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		padding: 24px;
		margin-bottom: 24px;
		box-shadow: var(--shadow-ring);
	}

	.form-card-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 20px;
		padding-bottom: 16px;
		border-bottom: 1px solid var(--border-frost);
	}

	.form-card-title {
		font-size: 0.9375rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.close-btn {
		background: none; border: none; cursor: pointer;
		color: var(--color-text-tertiary); padding: 4px;
		border-radius: var(--radius-subtle); display: flex;
		transition: color 0.12s, background 0.12s;
	}
	.close-btn:hover { color: var(--color-text-primary); background: var(--bg-hover-white); }

	.form-cols {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0 16px;
	}

	.form-divider {
		display: flex;
		align-items: center;
		gap: 10px;
		margin: 20px 0 16px;
	}

	.form-divider-label {
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		white-space: nowrap;
	}

	.form-divider::before,
	.form-divider::after {
		content: ''; flex: 1; height: 1px; background: var(--border-frost);
	}

	.form-actions {
		display: flex;
		gap: 8px;
		margin-top: 20px;
		padding-top: 20px;
		border-top: 1px solid var(--border-frost);
	}

	.checkbox-row {
		display: flex; align-items: center; gap: 8px;
	}
	.checkbox-row input[type="checkbox"] { width: 15px; height: 15px; cursor: pointer; }
	.checkbox-label { font-size: 0.875rem; color: var(--color-text-secondary); cursor: pointer; }

	.loading-state {
		display: flex; align-items: center; gap: 10px;
		color: var(--color-text-tertiary); font-size: 0.875rem; padding: 40px 0;
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
		display: flex; flex-direction: column; align-items: center;
		text-align: center; padding: 60px 24px; gap: 10px;
		background: var(--bg-card); border: 1px solid var(--border-frost);
		border-radius: var(--radius-card); box-shadow: var(--shadow-ring);
	}

	.empty-icon {
		width: 48px; height: 48px; border-radius: var(--radius-card);
		background: var(--bg-hover-subtle); border: 1px solid var(--border-frost);
		display: flex; align-items: center; justify-content: center;
		color: var(--color-text-tertiary); margin-bottom: 4px;
	}

	.empty-title { font-size: 0.9375rem; font-weight: 500; color: var(--color-text-primary); }
	.empty-sub { font-size: 0.8125rem; color: var(--color-text-tertiary); margin-bottom: 8px; }

	.list { display: flex; flex-direction: column; gap: 10px; }

	@media (max-width: 640px) {
		.page { padding: 24px 16px; }
		.form-cols { grid-template-columns: 1fr; }
	}
</style>
