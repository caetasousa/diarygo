<script lang="ts">
	import { currentUser } from '$lib/stores/auth';

	const adminActions = [
		{ icon: '👔', title: 'Profissionais', desc: 'Aprovar e gerenciar cadastros', href: '/admin/profissionais' },
		{ icon: '👥', title: 'Clientes', desc: 'Ver e gerenciar clientes', href: '/admin/clientes' },
		{ icon: '📋', title: 'Serviços', desc: 'Acompanhar todos os serviços', href: '/admin/servicos' },
		{ icon: '💰', title: 'Tabela de preços', desc: 'Configurar preços por região', href: '/admin/precos' }
	];

	const statsByTipo = {
		CLIENTE: [
			{ label: 'Serviços realizados', value: '0', icon: '✓', color: 'green' },
			{ label: 'Agendados', value: '0', icon: '📅', color: 'blue' },
			{ label: 'Score', value: '100', icon: '⭐', color: 'yellow' }
		],
		PROFISSIONAL: [
			{ label: 'Serviços realizados', value: '0', icon: '✓', color: 'green' },
			{ label: 'Nota média', value: '—', icon: '⭐', color: 'yellow' },
			{ label: 'Status', value: 'PENDENTE', icon: '⏳', color: 'orange' }
		],
		ADMIN: [
			{ label: 'Usuários', value: '—', icon: '👥', color: 'blue' },
			{ label: 'Serviços hoje', value: '—', icon: '📋', color: 'green' },
			{ label: 'Pendentes', value: '—', icon: '⚠️', color: 'orange' }
		]
	};

	$: stats = $currentUser ? statsByTipo[$currentUser.tipo] ?? [] : [];
</script>

<svelte:head>
	<title>Dashboard — DiaryGo</title>
</svelte:head>

<div class="section">
	<div class="container">
		<!-- Header -->
		<div class="dashboard-header">
			<div>
				<div class="row gap-3" style="align-items: center; margin-bottom: var(--space-2);">
					<h1 class="text-heading dash-title">Dashboard</h1>
					{#if $currentUser}
						<span class="badge badge-{$currentUser.tipo === 'CLIENTE' ? 'green' : $currentUser.tipo === 'PROFISSIONAL' ? 'orange' : 'blue'}">
							{$currentUser.tipo}
						</span>
					{/if}
				</div>
				<p class="text-body">
					{#if $currentUser}
						Olá, <strong style="color: var(--color-text-primary)">{$currentUser.email}</strong>
					{/if}
				</p>
			</div>

			{#if $currentUser?.tipo === 'CLIENTE'}
				<a href="/solicitacoes/novo" class="btn btn-white">
					+ Solicitar serviço
				</a>
			{:else if $currentUser?.tipo === 'ADMIN'}
				<a href="/admin/profissionais" class="btn btn-white">
					Gerenciar profissionais
				</a>
			{/if}
		</div>

		<!-- Stats -->
		<div class="stats-grid">
			{#each stats as stat}
				<div class="stat-card card">
					<div class="row gap-3" style="justify-content: space-between; align-items: flex-start;">
						<span class="stat-card-label text-caption">{stat.label}</span>
						<span class="badge badge-{stat.color}">{stat.icon}</span>
					</div>
					<span class="stat-card-value">{stat.value}</span>
				</div>
			{/each}
		</div>

		<!-- Content por tipo -->
		{#if $currentUser?.tipo === 'CLIENTE'}
			{@render ClienteDashboard()}
		{:else if $currentUser?.tipo === 'PROFISSIONAL'}
			{@render ProfissionalDashboard()}
		{:else if $currentUser?.tipo === 'ADMIN'}
			{@render AdminDashboard()}
		{/if}
	</div>
</div>

<!-- Sub-componentes inline -->
{#snippet ClienteDashboard()}
	<div class="dash-section">
		<div class="section-title-row">
			<h2 class="text-feature-title">Próximos serviços</h2>
			<a href="/solicitacoes" class="btn btn-ghost btn-sm">Ver todos</a>
		</div>
		<div class="empty-state card">
			<span class="empty-icon">🧹</span>
			<p class="text-subheading">Nenhum serviço agendado</p>
			<p class="text-body">Solicite seu primeiro serviço e mantenha sua casa impecável.</p>
			<a href="/solicitacoes/novo" class="btn btn-white">Solicitar agora</a>
		</div>
	</div>
{/snippet}

{#snippet ProfissionalDashboard()}
	<div class="dash-section">
		<div class="pending-approval card" style="border-color: var(--color-orange-4);">
			<div class="row gap-4" style="align-items: flex-start;">
				<span style="font-size: 1.5rem;">⏳</span>
				<div class="stack gap-2">
					<h3 class="text-subheading">Cadastro em análise</h3>
					<p class="text-body">
						Sua documentação está sendo analisada pela equipe DiaryGo.
						Você receberá uma notificação em até 48h úteis.
					</p>
					<span class="badge badge-orange">PENDENTE</span>
				</div>
			</div>
		</div>

		<div class="section-title-row" style="margin-top: var(--space-8);">
			<h2 class="text-feature-title">Próximos serviços</h2>
		</div>
		<div class="empty-state card">
			<span class="empty-icon">📅</span>
			<p class="text-subheading">Nenhum serviço agendado</p>
			<p class="text-body">Após aprovação, você começará a receber solicitações na sua região.</p>
		</div>
	</div>
{/snippet}

{#snippet AdminDashboard()}
	<div class="dash-section">
		<div class="admin-grid">
			{#each adminActions as action}
				<a href={action.href} class="admin-action card row gap-4">
					<span class="admin-action-icon">{action.icon}</span>
					<div class="stack gap-1">
						<h3 class="text-subheading">{action.title}</h3>
						<p class="text-caption">{action.desc}</p>
					</div>
				</a>
			{/each}
		</div>
	</div>
{/snippet}


<style>
	.dashboard-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-4);
		margin-bottom: var(--space-10);
		flex-wrap: wrap;
	}

	.dash-title {
		font-size: 2.5rem;
		letter-spacing: -1.5px;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: var(--space-4);
		margin-bottom: var(--space-10);
	}

	.stat-card {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.stat-card-label {
		color: var(--color-text-secondary);
	}

	.stat-card-value {
		font-size: 2rem;
		font-weight: 600;
		color: var(--color-text-primary);
		letter-spacing: -1px;
	}

	.dash-section {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}

	.section-title-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-4);
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: var(--space-4);
		padding: var(--space-12) var(--space-6);
	}

	.empty-icon {
		font-size: 2.5rem;
	}

	.admin-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: var(--space-4);
	}

	.admin-action {
		text-decoration: none;
		align-items: flex-start;
		transition: border-color 0.15s ease;
	}

	.admin-action:hover {
		border-color: rgba(214, 235, 253, 0.35);
		opacity: 1;
	}

	.admin-action-icon {
		font-size: 1.5rem;
		flex-shrink: 0;
	}

	.pending-approval {
		padding: var(--space-6);
	}
</style>
