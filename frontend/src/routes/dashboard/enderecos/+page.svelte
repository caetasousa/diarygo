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
	let mostarFormulario = $state(false);
	let editandoID = $state<string | null>(null);
	let salvando = $state(false);

	// Campos do formulário
	let cep = $state('');
	let logradouro = $state('');
	let numero = $state('');
	let complemento = $state('');
	let bairro = $state('');
	let cidade = $state('');
	let estado = $state('');
	let numQuartos = $state(1);
	let numBanheiros = $state(1);
	let numSalas = $state(1);
	let numCozinhas = $state(1);
	let areaM2 = $state(0);
	let principal = $state(false);
	let erros = $state<Record<string, string>>({});

	onMount(async () => {
		await carregarEnderecos();
	});

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
		resetarFormulario();
		editandoID = null;
		mostarFormulario = true;
	}

	function abrirEditar(id: string) {
		const e = enderecos.find((en) => en.id === id);
		if (!e) return;
		cep = e.cep;
		logradouro = e.logradouro;
		numero = e.numero;
		complemento = e.complemento;
		bairro = e.bairro;
		cidade = e.cidade;
		estado = e.estado;
		numQuartos = e.num_quartos;
		numBanheiros = e.num_banheiros;
		numSalas = e.num_salas;
		numCozinhas = e.num_cozinhas;
		areaM2 = e.area_m2;
		principal = e.principal;
		editandoID = id;
		mostarFormulario = true;
		erros = {};
	}

	function resetarFormulario() {
		cep = ''; logradouro = ''; numero = ''; complemento = '';
		bairro = ''; cidade = ''; estado = '';
		numQuartos = 1; numBanheiros = 1; numSalas = 1; numCozinhas = 1;
		areaM2 = 0; principal = false; erros = {};
	}

	function cancelar() {
		mostarFormulario = false;
		resetarFormulario();
	}

	function validar(): boolean {
		erros = {};
		const cepDigits = cep.replace(/\D/g, '');
		if (cepDigits.length !== 8) erros.cep = 'CEP inválido';
		if (!logradouro.trim()) erros.logradouro = 'Logradouro é obrigatório';
		if (!cidade.trim()) erros.cidade = 'Cidade é obrigatória';
		if (estado.length !== 2) erros.estado = 'Estado inválido (UF)';
		if (numQuartos < 1) erros.num_quartos = 'Mínimo 1 quarto';
		return Object.keys(erros).length === 0;
	}

	async function salvar() {
		if (!validar()) return;
		salvando = true;
		try {
			const req = {
				logradouro: logradouro.trim(),
				numero: numero.trim(),
				complemento: complemento.trim(),
				bairro: bairro.trim(),
				cidade: cidade.trim(),
				estado: estado.toUpperCase().trim(),
				cep: cep.replace(/\D/g, ''),
				num_quartos: numQuartos,
				num_banheiros: numBanheiros,
				num_salas: numSalas,
				num_cozinhas: numCozinhas,
				area_m2: areaM2,
				principal
			};

			if (editandoID) {
				await api.atualizarEndereco(editandoID, req);
				toasts.success('Endereço atualizado!');
			} else {
				await api.criarEndereco(req);
				toasts.success('Endereço criado!');
			}

			mostarFormulario = false;
			resetarFormulario();
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

<div class="container" style="max-width: 720px; padding: 2rem 1rem;">
	<div class="flex" style="justify-content: space-between; align-items: center; margin-bottom: 2rem; flex-wrap: wrap; gap: 1rem;">
		<div>
			<h1 class="text-2xl font-bold">Meus Endereços</h1>
			<p class="text-muted">Gerencie os locais onde o serviço será realizado.</p>
		</div>
		{#if !mostarFormulario}
			<button class="btn btn-primary" onclick={abrirNovo}>+ Novo endereço</button>
		{/if}
	</div>

	{#if mostarFormulario}
		<div class="card" style="margin-bottom: 2rem;">
			<h2 class="font-bold" style="margin-bottom: 1.5rem;">
				{editandoID ? 'Editar endereço' : 'Novo endereço'}
			</h2>

			<form novalidate onsubmit={(e) => { e.preventDefault(); salvar(); }}>
				<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
					<div style="grid-column: 1 / -1;">
						<CepInput bind:value={cep} error={erros.cep} onPreenchido={preencherViaCEP} />
					</div>

					<div class="form-group" style="grid-column: 1 / -1;">
						<label for="logradouro" class="form-label">Logradouro</label>
						<input id="logradouro" type="text" maxlength="200" class="form-input"
							class:input-error={!!erros.logradouro} bind:value={logradouro}
							placeholder="Rua, Avenida, etc." autocomplete="street-address" />
						{#if erros.logradouro}<p class="form-error">{erros.logradouro}</p>{/if}
					</div>

					<div class="form-group">
						<label for="numero" class="form-label">Número</label>
						<input id="numero" type="text" maxlength="20" class="form-input"
							bind:value={numero} placeholder="123" autocomplete="off" />
					</div>

					<div class="form-group">
						<label for="complemento" class="form-label">Complemento</label>
						<input id="complemento" type="text" maxlength="100" class="form-input"
							bind:value={complemento} placeholder="Apto, Bloco..." autocomplete="off" />
					</div>

					<div class="form-group">
						<label for="bairro" class="form-label">Bairro</label>
						<input id="bairro" type="text" maxlength="100" class="form-input"
							bind:value={bairro} placeholder="Bairro" autocomplete="off" />
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
							{#each ufs as uf}
								<option value={uf}>{uf}</option>
							{/each}
						</select>
						{#if erros.estado}<p class="form-error">{erros.estado}</p>{/if}
					</div>
				</div>

				<h3 class="font-bold" style="margin: 1.5rem 0 1rem;">Detalhes do imóvel</h3>
				<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem;">
					<div class="form-group">
						<label for="quartos" class="form-label">Quartos</label>
						<input id="quartos" type="number" min="1" max="20" class="form-input"
							class:input-error={!!erros.num_quartos} bind:value={numQuartos} />
						{#if erros.num_quartos}<p class="form-error">{erros.num_quartos}</p>{/if}
					</div>
					<div class="form-group">
						<label for="banheiros" class="form-label">Banheiros</label>
						<input id="banheiros" type="number" min="0" max="20" class="form-input" bind:value={numBanheiros} />
					</div>
					<div class="form-group">
						<label for="salas" class="form-label">Salas</label>
						<input id="salas" type="number" min="0" max="20" class="form-input" bind:value={numSalas} />
					</div>
					<div class="form-group">
						<label for="cozinhas" class="form-label">Cozinhas</label>
						<input id="cozinhas" type="number" min="0" max="10" class="form-input" bind:value={numCozinhas} />
					</div>
					<div class="form-group">
						<label for="area" class="form-label">Área (m²)</label>
						<input id="area" type="number" min="0" step="0.1" class="form-input" bind:value={areaM2} />
					</div>
					<div class="form-group" style="display: flex; align-items: flex-end; padding-bottom: 0.5rem;">
						<label class="flex gap-sm" style="align-items: center; cursor: pointer;">
							<input type="checkbox" bind:checked={principal} style="width: 1rem; height: 1rem;" />
							<span>Endereço principal</span>
						</label>
					</div>
				</div>

				<div class="flex gap-sm" style="margin-top: 1.5rem;">
					<button type="submit" class="btn btn-primary" disabled={salvando}>
						{salvando ? 'Salvando...' : 'Salvar'}
					</button>
					<button type="button" class="btn btn-secondary" onclick={cancelar}>Cancelar</button>
				</div>
			</form>
		</div>
	{/if}

	{#if carregando}
		<p class="text-muted">Carregando endereços...</p>
	{:else if enderecos.length === 0 && !mostarFormulario}
		<div class="card" style="text-align: center; padding: 3rem; color: var(--color-text-muted);">
			<p style="margin-bottom: 1rem;">Nenhum endereço cadastrado.</p>
			<button class="btn btn-primary" onclick={abrirNovo}>Adicionar primeiro endereço</button>
		</div>
	{:else}
		<div style="display: flex; flex-direction: column; gap: 1rem;">
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
