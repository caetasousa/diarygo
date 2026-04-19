<script lang="ts">
	import '../../app.css';
	import { page } from '$app/state';
	import { auth, currentUser } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import {
		LayoutDashboard,
		User,
		MapPin,
		BadgeCheck,
		CompassIcon,
		LogOut,
		Sparkles,
		LifeBuoy,
		CalendarCheck,
		ShieldCheck,
		Plus
	} from 'lucide-svelte';

	let { children } = $props();

	type NavItem = {
		href: string;
		label: string;
		icon: typeof LayoutDashboard;
		exact?: boolean;
	};
	type NavGroup = { section: string; items: NavItem[] };

	const navCliente: NavGroup[] = [
		{ section: 'Principal', items: [
			{ href: '/dashboard', label: 'Visão geral', icon: LayoutDashboard, exact: true },
			{ href: '/solicitacoes', label: 'Meus serviços', icon: CalendarCheck },
		]},
		{ section: 'Conta', items: [
			{ href: '/dashboard/perfil', label: 'Meu perfil', icon: User },
			{ href: '/dashboard/enderecos', label: 'Endereços', icon: MapPin },
		]},
	];

	const navProfissional: NavGroup[] = [
		{ section: 'Principal', items: [
			{ href: '/dashboard', label: 'Visão geral', icon: LayoutDashboard, exact: true },
			{ href: '/agenda', label: 'Minha agenda', icon: CalendarCheck },
		]},
		{ section: 'Conta', items: [
			{ href: '/dashboard/perfil', label: 'Meu perfil', icon: User },
			{ href: '/dashboard/credenciamento', label: 'Credenciamento', icon: BadgeCheck },
			{ href: '/dashboard/atuacao', label: 'Atuação', icon: CompassIcon },
		]},
	];

	let tipo = $derived($currentUser?.tipo ?? null);
	let navGroups = $derived(tipo === 'PROFISSIONAL' ? navProfissional : navCliente);
	let currentPath = $derived(page.url.pathname);
	let emailInitial = $derived(($currentUser?.email ?? 'U')[0].toUpperCase());
	let emailLabel = $derived($currentUser?.email ?? '');
	let roleLabel = $derived(
		tipo === 'CLIENTE' ? 'Cliente' :
		tipo === 'PROFISSIONAL' ? 'Profissional' :
		tipo === 'ADMIN' ? 'Administrador' : ''
	);

	function isActive(href: string, exact = false) {
		if (exact) return currentPath === href;
		return currentPath === href || currentPath.startsWith(href + '/');
	}

	function handleLogout() {
		auth.logout();
		goto('/');
	}

</script>

<div class="dash-shell">
	<!-- ═══════════════ Sidebar ═══════════════ -->
	<aside class="sidebar">
		<div class="sidebar-scroll">
			<a href="/" class="sidebar-brand">
				<span class="brand-mark">
					<Sparkles size={14} strokeWidth={2.2} />
				</span>
				<span class="brand-name">DiaryGo</span>
			</a>

			{#if tipo === 'CLIENTE'}
				<a href="/solicitacoes/novo" class="sidebar-cta">
					<Plus size={14} strokeWidth={2.5} />
					Solicitar serviço
				</a>
			{:else if tipo === 'PROFISSIONAL'}
				<div class="sidebar-status">
					<span class="status-dot"></span>
					<div class="status-body">
						<span class="status-title">Em análise</span>
						<span class="status-sub">Aguardando aprovação</span>
					</div>
				</div>
			{/if}

			<nav class="sidebar-nav">
				{#each navGroups as group}
					<div class="nav-group">
						<span class="nav-section">{group.section}</span>
						{#each group.items as item}
							{@const Icon = item.icon}
							<a
								href={item.href}
								class="nav-item"
								class:active={isActive(item.href, item.exact)}
							>
								<Icon size={16} strokeWidth={1.8} />
								<span>{item.label}</span>
							</a>
						{/each}
					</div>
				{/each}
			</nav>
		</div>

		<div class="sidebar-foot">
			<div class="user-card">
				<span class="user-avatar">{emailInitial}</span>
				<div class="user-info">
					<span class="user-email" title={emailLabel}>{emailLabel}</span>
					<span class="user-role">{roleLabel}</span>
				</div>
				<button class="user-logout" onclick={handleLogout} aria-label="Sair">
					<LogOut size={14} strokeWidth={1.8} />
				</button>
			</div>
			<a href="/ajuda" class="foot-link">
				<LifeBuoy size={14} strokeWidth={1.8} />
				<span>Central de ajuda</span>
			</a>
			<div class="foot-trust">
				<ShieldCheck size={12} strokeWidth={1.8} />
				<span>Ambiente seguro</span>
			</div>
		</div>
	</aside>

	<!-- ═══════════════ Main column ═══════════════ -->
	<div class="dash-column">
		<!-- Page content -->
		<main class="dash-main">
			{@render children()}
		</main>
	</div>
</div>

<ToastContainer />

<style>
	.dash-shell {
		display: flex;
		min-height: 100vh;
		background: var(--color-black);
	}

	/* ═══════════════ Sidebar ═══════════════ */
	.sidebar {
		width: 240px;
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		border-right: 1px solid var(--border-frost);
		position: sticky;
		top: 0;
		height: 100vh;
		background: linear-gradient(180deg, rgba(255,255,255,0.012) 0%, transparent 40%);
	}

	.sidebar-scroll {
		display: flex;
		flex-direction: column;
		gap: 18px;
		padding: 20px 14px 14px;
		overflow-y: auto;
		flex: 1;
		min-height: 0;
	}

	.sidebar-brand {
		display: flex;
		align-items: center;
		gap: 9px;
		text-decoration: none;
		color: var(--color-text-primary);
		font-weight: 600;
		font-size: 0.9375rem;
		letter-spacing: -0.3px;
		padding: 4px 6px;
	}

	.brand-mark {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 22px;
		height: 22px;
		border-radius: var(--radius-subtle);
		background: linear-gradient(135deg, var(--color-orange-10), #ff6b1a);
		color: #1a0a00;
		flex-shrink: 0;
	}

	/* Primary CTA */
	.sidebar-cta {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		padding: 8px 12px;
		background: var(--color-white);
		color: var(--color-black);
		border-radius: var(--radius-pill);
		font-size: 0.8125rem;
		font-weight: 600;
		letter-spacing: -0.2px;
		text-decoration: none;
		transition: opacity 0.15s ease, transform 0.15s ease;
	}

	.sidebar-cta:hover {
		opacity: 0.92;
		transform: translateY(-1px);
	}

	/* Status pill */
	.sidebar-status {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 10px 12px;
		background: var(--color-orange-4);
		border: 1px solid var(--color-orange-4);
		border-radius: var(--radius-standard);
	}

	.status-dot {
		width: 7px;
		height: 7px;
		border-radius: 50%;
		background: var(--color-orange-10);
		box-shadow: 0 0 0 3px rgba(255, 128, 31, 0.2);
		animation: pulse 2s ease-in-out infinite;
		flex-shrink: 0;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.4; }
	}

	.status-body {
		display: flex;
		flex-direction: column;
		gap: 1px;
		min-width: 0;
	}

	.status-title {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-orange-11);
	}

	.status-sub {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	/* Nav */
	.sidebar-nav {
		display: flex;
		flex-direction: column;
		gap: 14px;
	}

	.nav-group {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}

	.nav-section {
		font-size: 0.625rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.8px;
		padding: 4px 10px 6px;
	}

	.nav-item {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 7px 10px;
		border-radius: var(--radius-standard);
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		text-decoration: none;
		transition: color 0.12s ease, background 0.12s ease;
		letter-spacing: -0.1px;
		position: relative;
	}

	.nav-item:hover {
		color: var(--color-text-primary);
		background: var(--bg-hover-subtle);
		opacity: 1;
	}

	.nav-item.active {
		color: var(--color-text-primary);
		background: var(--bg-hover-medium);
	}

	.nav-item.active::before {
		content: '';
		position: absolute;
		left: -14px;
		top: 50%;
		transform: translateY(-50%);
		width: 2px;
		height: 14px;
		background: var(--color-orange-10);
		border-radius: 0 2px 2px 0;
	}

	.nav-item :global(svg) {
		flex-shrink: 0;
		opacity: 0.7;
	}

	.nav-item.active :global(svg),
	.nav-item:hover :global(svg) {
		opacity: 1;
	}

	/* Sidebar foot */
	.sidebar-foot {
		padding: 12px 14px 16px;
		border-top: 1px solid var(--border-frost);
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.foot-link {
		display: flex;
		align-items: center;
		gap: 9px;
		padding: 6px 8px;
		border-radius: var(--radius-subtle);
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		text-decoration: none;
		transition: color 0.12s ease, background 0.12s ease;
	}

	.foot-link:hover {
		color: var(--color-text-primary);
		background: var(--bg-hover-subtle);
		opacity: 1;
	}

	.foot-trust {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 0 8px;
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	.foot-trust :global(svg) {
		color: var(--color-green-10);
	}

	/* ═══════════════ Main column ═══════════════ */
	.dash-column {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
	}

	/* User card na sidebar */
	.user-card {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 8px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-standard);
		margin-bottom: 8px;
	}

	.user-avatar {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--color-orange-10), #ff6b1a);
		color: #1a0a00;
		font-size: 0.75rem;
		font-weight: 700;
		flex-shrink: 0;
	}

	.user-info {
		display: flex;
		flex-direction: column;
		gap: 0;
		flex: 1;
		min-width: 0;
	}

	.user-email {
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-text-primary);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.user-role {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	.user-logout {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 26px;
		height: 26px;
		background: transparent;
		border: 1px solid transparent;
		border-radius: var(--radius-subtle);
		color: var(--color-text-tertiary);
		cursor: pointer;
		transition: color 0.12s, background 0.12s, border-color 0.12s;
		flex-shrink: 0;
	}

	.user-logout:hover {
		color: var(--color-red-10);
		background: var(--color-red-5);
		border-color: var(--border-frost);
	}

	/* ═══════════════ Main ═══════════════ */
	.dash-main {
		flex: 1;
		min-width: 0;
		overflow-x: hidden;
	}

	/* ═══════════════ Mobile ═══════════════ */
	@media (max-width: 900px) {
		.sidebar {
			width: 64px;
		}

		.sidebar-scroll {
			padding: 16px 8px 10px;
			gap: 14px;
		}

		.brand-name,
		.nav-item span,
		.nav-section,
		.status-body,
		.sidebar-cta,
		.foot-link span,
		.foot-trust span,
		.user-info {
			display: none;
		}

		.user-card {
			justify-content: center;
			padding: 6px;
		}

		.sidebar-cta {
			display: flex;
			padding: 8px;
			aspect-ratio: 1;
		}

		.nav-item {
			justify-content: center;
			padding: 9px;
		}

		.nav-item.active::before {
			display: none;
		}

		.sidebar-status {
			padding: 8px;
			justify-content: center;
		}

		.sidebar-foot {
			padding: 10px 8px 14px;
			align-items: center;
		}

		.foot-link {
			padding: 6px;
		}
	}

</style>
