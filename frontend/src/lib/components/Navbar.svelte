<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth, isAuthenticated, currentUser } from '$lib/stores/auth';
	import { Sparkles, Menu, X, ArrowRight, LogOut } from 'lucide-svelte';

	let menuOpen = $state(false);

	function handleLogout() {
		auth.logout();
		goto('/');
	}
</script>

<header class="navbar">
	<div class="navbar-inner">
		<a href="/" class="navbar-brand">
			<span class="brand-mark">
				<Sparkles size={13} strokeWidth={2.2} />
			</span>
			<span class="brand-name">DiaryGo</span>
		</a>

		<nav class="navbar-links" class:open={menuOpen}>
			{#if $isAuthenticated}
				<a href="/dashboard" class="nav-link">Dashboard</a>
				{#if $currentUser?.tipo === 'CLIENTE'}
					<a href="/solicitacoes" class="nav-link">Meus serviços</a>
				{:else if $currentUser?.tipo === 'PROFISSIONAL'}
					<a href="/agenda" class="nav-link">Minha agenda</a>
				{/if}
			{:else}
				<a href="/#servicos" class="nav-link">Serviços</a>
				<a href="/#como-funciona" class="nav-link">Como funciona</a>
				<a href="/#para-diaristas" class="nav-link">Para diaristas</a>
				<a href="/#faq" class="nav-link">Perguntas</a>
			{/if}
		</nav>

		<div class="navbar-actions">
			{#if $isAuthenticated}
				<button class="link-btn" onclick={handleLogout} aria-label="Sair">
					<LogOut size={14} strokeWidth={2} />
					Sair
				</button>
			{:else}
				<a href="/login" class="link-btn desktop-only">Entrar</a>
				<a href="/registro" class="navbar-cta">
					Começar agora
					<ArrowRight size={13} strokeWidth={2.2} />
				</a>
			{/if}

			<button
				class="menu-toggle"
				aria-label={menuOpen ? 'Fechar menu' : 'Abrir menu'}
				onclick={() => (menuOpen = !menuOpen)}
			>
				{#if menuOpen}
					<X size={18} strokeWidth={2} />
				{:else}
					<Menu size={18} strokeWidth={2} />
				{/if}
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
		background: rgba(0, 0, 0, 0.72);
		backdrop-filter: blur(14px) saturate(140%);
		-webkit-backdrop-filter: blur(14px) saturate(140%);
		border-bottom: 1px solid var(--border-frost);
	}

	.navbar-inner {
		max-width: 1240px;
		margin: 0 auto;
		padding: 0 24px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 32px;
		height: 62px;
	}

	.navbar-brand {
		display: flex;
		align-items: center;
		gap: 9px;
		text-decoration: none;
		color: var(--color-text-primary);
		flex-shrink: 0;
	}

	.brand-mark {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 22px;
		height: 22px;
		border-radius: 6px;
		background: linear-gradient(135deg, var(--color-orange-10), #ff6b1a);
		color: #1a0a00;
	}

	.brand-name {
		font-size: 0.9375rem;
		font-weight: 600;
		letter-spacing: -0.3px;
	}

	.navbar-links {
		display: flex;
		align-items: center;
		gap: 4px;
		flex: 1;
		justify-content: center;
	}

	.nav-link {
		padding: 7px 12px;
		border-radius: 7px;
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		text-decoration: none;
		letter-spacing: -0.1px;
		transition: color 0.12s, background 0.12s;
	}

	.nav-link:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.04);
		opacity: 1;
	}

	.navbar-actions {
		display: flex;
		align-items: center;
		gap: 6px;
		flex-shrink: 0;
	}

	.link-btn {
		display: inline-flex;
		align-items: center;
		gap: 5px;
		padding: 7px 12px;
		background: transparent;
		border: none;
		border-radius: 7px;
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		font-family: var(--font-body);
		cursor: pointer;
		text-decoration: none;
		transition: color 0.12s, background 0.12s;
	}

	.link-btn:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.04);
		opacity: 1;
	}

	.navbar-cta {
		display: inline-flex;
		align-items: center;
		gap: 5px;
		padding: 8px 14px;
		background: var(--color-white);
		color: var(--color-black);
		border-radius: 7px;
		font-size: 0.8125rem;
		font-weight: 600;
		text-decoration: none;
		letter-spacing: -0.1px;
		transition: opacity 0.15s, transform 0.15s;
	}

	.navbar-cta:hover {
		opacity: 0.92;
		transform: translateY(-1px);
	}

	.menu-toggle {
		display: none;
		width: 34px;
		height: 34px;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: 1px solid var(--border-frost);
		border-radius: 7px;
		color: var(--color-text-primary);
		cursor: pointer;
		transition: background 0.12s;
	}

	.menu-toggle:hover {
		background: rgba(255, 255, 255, 0.04);
	}

	@media (max-width: 900px) {
		.navbar-inner {
			height: 58px;
		}

		.navbar-links {
			display: none;
			position: absolute;
			top: 100%;
			left: 0;
			right: 0;
			background: #0a0a0a;
			border-bottom: 1px solid var(--border-frost);
			flex-direction: column;
			align-items: stretch;
			padding: 8px 16px 16px;
			gap: 2px;
		}

		.navbar-links.open {
			display: flex;
		}

		.nav-link {
			padding: 12px 14px;
			font-size: 0.9375rem;
		}

		.desktop-only {
			display: none;
		}

		.menu-toggle {
			display: inline-flex;
		}
	}
</style>
