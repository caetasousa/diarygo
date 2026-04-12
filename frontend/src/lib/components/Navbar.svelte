<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth, isAuthenticated, currentUser } from '$lib/stores/auth';

	let menuOpen = false;

	function handleLogout() {
		auth.logout();
		goto('/');
	}
</script>

<header class="navbar">
	<div class="container navbar-inner">
		<a href="/" class="navbar-brand">
			<span class="brand-icon">✦</span>
			<span class="brand-name">DiaryGo</span>
		</a>

		<nav class="navbar-links" class:open={menuOpen}>
			<a href="/" class="nav-link text-nav">Início</a>
			{#if $isAuthenticated}
				{#if $currentUser?.tipo === 'CLIENTE'}
					<a href="/dashboard" class="nav-link text-nav">Dashboard</a>
					<a href="/solicitacoes" class="nav-link text-nav">Serviços</a>
				{:else if $currentUser?.tipo === 'PROFISSIONAL'}
					<a href="/dashboard" class="nav-link text-nav">Dashboard</a>
					<a href="/agenda" class="nav-link text-nav">Agenda</a>
				{:else if $currentUser?.tipo === 'ADMIN'}
					<a href="/admin" class="nav-link text-nav">Admin</a>
				{/if}
			{:else}
				<a href="/como-funciona" class="nav-link text-nav">Como funciona</a>
				<a href="/profissionais" class="nav-link text-nav">Para diaristas</a>
			{/if}
		</nav>

		<div class="navbar-actions">
			{#if $isAuthenticated}
				<span class="badge badge-blue text-small">{$currentUser?.tipo}</span>
				<button class="btn btn-primary btn-sm" on:click={handleLogout}>
					Sair
				</button>
			{:else}
				<a href="/login" class="btn btn-primary btn-sm">Entrar</a>
				<a href="/registro" class="btn btn-white btn-sm">Começar</a>
			{/if}

			<button
				class="menu-toggle"
				aria-label="Menu"
				on:click={() => (menuOpen = !menuOpen)}
			>
				<span></span>
				<span></span>
				<span></span>
			</button>
		</div>
	</div>
</header>

<style>
	.navbar {
		position: sticky;
		top: 0;
		z-index: 100;
		width: 100%;
		background: rgba(0, 0, 0, 0.85);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1px solid var(--border-frost);
	}

	.navbar-inner {
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 60px;
		gap: var(--space-6);
	}

	.navbar-brand {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		text-decoration: none;
		color: var(--color-text-primary);
		flex-shrink: 0;
	}

	.brand-icon {
		color: var(--color-orange-10);
		font-size: 1.25rem;
	}

	.brand-name {
		font-size: 1rem;
		font-weight: 600;
		letter-spacing: -0.3px;
	}

	.navbar-links {
		display: flex;
		align-items: center;
		gap: var(--space-6);
		flex: 1;
	}

	.nav-link {
		color: var(--color-text-secondary);
		text-decoration: none;
		transition: color 0.15s ease;
	}

	.nav-link:hover {
		color: var(--color-text-primary);
		opacity: 1;
	}

	.navbar-actions {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		flex-shrink: 0;
	}

	.menu-toggle {
		display: none;
		flex-direction: column;
		justify-content: center;
		gap: 5px;
		width: 24px;
		height: 24px;
		background: none;
		border: none;
		cursor: pointer;
		padding: 0;
	}

	.menu-toggle span {
		display: block;
		width: 100%;
		height: 1px;
		background: var(--color-text-primary);
		transition: transform 0.2s ease, opacity 0.2s ease;
	}

	@media (max-width: 600px) {
		.navbar-links {
			display: none;
			position: absolute;
			top: 60px;
			left: 0;
			right: 0;
			background: var(--color-black);
			border-bottom: 1px solid var(--border-frost);
			flex-direction: column;
			padding: var(--space-4) var(--space-6);
			gap: var(--space-4);
			align-items: flex-start;
		}

		.navbar-links.open {
			display: flex;
		}

		.menu-toggle {
			display: flex;
		}
	}
</style>
