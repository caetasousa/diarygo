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
		FileText,
		Users,
		Map,
		CalendarClock,
		LogOut,
		Bell,
		Search,
		ChevronDown,
		Sparkles,
		LifeBuoy,
		CalendarCheck,
		ShieldCheck,
		Settings,
		Plus
	} from 'lucide-svelte';

	let { children } = $props();

	let profileOpen = $state(false);
	let notifOpen = $state(false);

	const navCliente = [
		{ section: 'Principal', items: [
			{ href: '/dashboard', label: 'Visão geral', icon: LayoutDashboard, exact: true },
			{ href: '/solicitacoes', label: 'Meus serviços', icon: CalendarCheck },
		]},
		{ section: 'Conta', items: [
			{ href: '/dashboard/perfil', label: 'Meu perfil', icon: User },
			{ href: '/dashboard/enderecos', label: 'Endereços', icon: MapPin },
		]},
	];

	const navProfissional = [
		{ section: 'Principal', items: [
			{ href: '/dashboard', label: 'Visão geral', icon: LayoutDashboard, exact: true },
			{ href: '/agenda', label: 'Minha agenda', icon: CalendarCheck },
		]},
		{ section: 'Cadastro', items: [
			{ href: '/dashboard/perfil', label: 'Meu perfil', icon: User },
			{ href: '/dashboard/documentos', label: 'Documentos', icon: FileText },
			{ href: '/dashboard/referencias', label: 'Referências', icon: Users },
		]},
		{ section: 'Atuação', items: [
			{ href: '/dashboard/regioes', label: 'Regiões', icon: Map },
			{ href: '/dashboard/disponibilidade', label: 'Disponibilidade', icon: CalendarClock },
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

	function closeMenus() {
		profileOpen = false;
		notifOpen = false;
	}
</script>

<svelte:window onclick={closeMenus} />

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
		<!-- Topbar -->
		<header class="topbar">
		<div
			class="topbar-inner"
			role="presentation"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
		>
			<div class="topbar-search">
				<Search size={15} strokeWidth={1.8} />
				<input type="text" placeholder="Buscar serviços, endereços, diaristas..." />
				<kbd class="kbd">⌘K</kbd>
			</div>

			<div class="topbar-actions">
				<div class="topbar-menu-wrap">
					<button
						class="topbar-icon-btn"
						onclick={() => { notifOpen = !notifOpen; profileOpen = false; }}
						aria-label="Notificações"
					>
						<Bell size={16} strokeWidth={1.8} />
						<span class="notif-dot"></span>
					</button>

					{#if notifOpen}
						<div class="dropdown notif-dropdown">
							<div class="dropdown-head">
								<span class="dropdown-title">Notificações</span>
								<button class="dropdown-action">Marcar lidas</button>
							</div>
							<div class="notif-empty">
								<Bell size={20} strokeWidth={1.4} />
								<span class="notif-empty-title">Nenhuma notificação</span>
								<span class="notif-empty-sub">Avisos sobre agendamentos aparecerão aqui.</span>
							</div>
						</div>
					{/if}
				</div>

				<div class="topbar-divider"></div>

				<div class="topbar-menu-wrap">
					<button
						class="profile-btn"
						onclick={() => { profileOpen = !profileOpen; notifOpen = false; }}
					>
						<span class="profile-avatar">{emailInitial}</span>
						<div class="profile-info">
							<span class="profile-name">{emailLabel}</span>
							<span class="profile-role">{roleLabel}</span>
						</div>
						<ChevronDown size={14} strokeWidth={1.8} />
					</button>

					{#if profileOpen}
						<div class="dropdown profile-dropdown">
							<div class="profile-dropdown-head">
								<span class="profile-avatar profile-avatar-lg">{emailInitial}</span>
								<div class="profile-info">
									<span class="profile-name">{emailLabel}</span>
									<span class="profile-role">{roleLabel}</span>
								</div>
							</div>
							<div class="dropdown-sep"></div>
							<a href="/dashboard/perfil" class="dropdown-item">
								<User size={14} strokeWidth={1.8} />
								<span>Meu perfil</span>
							</a>
							<a href="/dashboard/configuracoes" class="dropdown-item">
								<Settings size={14} strokeWidth={1.8} />
								<span>Preferências</span>
							</a>
							<a href="/ajuda" class="dropdown-item">
								<LifeBuoy size={14} strokeWidth={1.8} />
								<span>Suporte</span>
							</a>
							<div class="dropdown-sep"></div>
							<button class="dropdown-item danger" onclick={handleLogout}>
								<LogOut size={14} strokeWidth={1.8} />
								<span>Sair</span>
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
		</header>

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
		border-radius: 6px;
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
		border-radius: 8px;
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
		background: rgba(255, 128, 31, 0.08);
		border: 1px solid var(--color-orange-4);
		border-radius: 8px;
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
		border-radius: 7px;
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
		background: rgba(255, 255, 255, 0.04);
		opacity: 1;
	}

	.nav-item.active {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.07);
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
		border-radius: 6px;
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		text-decoration: none;
		transition: color 0.12s ease, background 0.12s ease;
	}

	.foot-link:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.04);
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

	/* ═══════════════ Topbar ═══════════════ */
	.topbar {
		border-bottom: 1px solid var(--border-frost);
		background: rgba(0, 0, 0, 0.7);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		position: sticky;
		top: 0;
		z-index: 20;
	}

	.topbar-inner {
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 12px 32px;
	}

	.topbar-search {
		display: flex;
		align-items: center;
		gap: 9px;
		flex: 1;
		max-width: 440px;
		padding: 8px 14px;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid var(--border-frost);
		border-radius: 8px;
		color: var(--color-text-tertiary);
		transition: border-color 0.15s, background 0.15s;
	}

	.topbar-search:focus-within {
		border-color: rgba(214, 235, 253, 0.3);
		background: rgba(255, 255, 255, 0.05);
	}

	.topbar-search input {
		flex: 1;
		background: transparent;
		border: none;
		outline: none;
		color: var(--color-text-primary);
		font-size: 0.8125rem;
		font-family: var(--font-body);
	}

	.topbar-search input::placeholder {
		color: var(--color-text-tertiary);
	}

	.kbd {
		font-family: var(--font-mono);
		font-size: 0.625rem;
		color: var(--color-text-tertiary);
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid var(--border-frost);
		border-radius: 4px;
		padding: 1px 5px;
	}

	.topbar-actions {
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.topbar-divider {
		width: 1px;
		height: 22px;
		background: var(--border-frost);
		margin: 0 4px;
	}

	.topbar-icon-btn {
		position: relative;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 34px;
		height: 34px;
		background: transparent;
		border: 1px solid transparent;
		border-radius: 8px;
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: color 0.12s, background 0.12s, border-color 0.12s;
	}

	.topbar-icon-btn:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.04);
		border-color: var(--border-frost);
	}

	.notif-dot {
		position: absolute;
		top: 8px;
		right: 8px;
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-orange-10);
		border: 1.5px solid var(--color-black);
	}

	/* Profile button */
	.topbar-menu-wrap {
		position: relative;
	}

	.profile-btn {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 5px 10px 5px 5px;
		background: transparent;
		border: 1px solid transparent;
		border-radius: 8px;
		cursor: pointer;
		transition: background 0.12s, border-color 0.12s;
	}

	.profile-btn:hover {
		background: rgba(255, 255, 255, 0.04);
		border-color: var(--border-frost);
	}

	.profile-avatar {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 26px;
		height: 26px;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--color-orange-10), #ff6b1a);
		color: #1a0a00;
		font-size: 0.75rem;
		font-weight: 700;
		flex-shrink: 0;
	}

	.profile-avatar-lg {
		width: 38px;
		height: 38px;
		font-size: 0.9375rem;
	}

	.profile-info {
		display: flex;
		flex-direction: column;
		gap: 0;
		min-width: 0;
		max-width: 180px;
	}

	.profile-name {
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-text-primary);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.profile-role {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	.profile-btn :global(svg) {
		color: var(--color-text-tertiary);
		flex-shrink: 0;
	}

	/* Dropdowns */
	.dropdown {
		position: absolute;
		top: calc(100% + 8px);
		right: 0;
		min-width: 280px;
		background: #0a0a0a;
		border: 1px solid var(--border-frost);
		border-radius: 10px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.6), var(--shadow-ring);
		padding: 6px;
		z-index: 50;
		animation: dropdownIn 0.15s ease-out;
	}

	@keyframes dropdownIn {
		from { opacity: 0; transform: translateY(-4px); }
		to { opacity: 1; transform: translateY(0); }
	}

	.dropdown-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 8px 10px 10px;
	}

	.dropdown-title {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.dropdown-action {
		background: transparent;
		border: none;
		color: var(--color-text-tertiary);
		font-size: 0.6875rem;
		cursor: pointer;
		padding: 2px 4px;
		border-radius: 4px;
	}

	.dropdown-action:hover {
		color: var(--color-text-primary);
	}

	.dropdown-sep {
		height: 1px;
		background: var(--border-frost);
		margin: 4px 0;
	}

	.dropdown-item {
		display: flex;
		align-items: center;
		gap: 10px;
		width: 100%;
		padding: 8px 10px;
		background: transparent;
		border: none;
		border-radius: 6px;
		font-size: 0.8125rem;
		color: var(--color-text-secondary);
		text-decoration: none;
		cursor: pointer;
		text-align: left;
		transition: color 0.12s, background 0.12s;
		font-family: var(--font-body);
	}

	.dropdown-item:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.05);
		opacity: 1;
	}

	.dropdown-item.danger:hover {
		color: var(--color-red-10);
		background: var(--color-red-5);
	}

	.dropdown-item :global(svg) {
		color: var(--color-text-tertiary);
	}

	.dropdown-item:hover :global(svg) {
		color: currentColor;
	}

	.notif-dropdown {
		min-width: 320px;
	}

	.notif-empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: 6px;
		padding: 28px 16px 20px;
	}

	.notif-empty :global(svg) {
		color: var(--color-text-tertiary);
		margin-bottom: 4px;
	}

	.notif-empty-title {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.notif-empty-sub {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		max-width: 220px;
	}

	.profile-dropdown {
		min-width: 260px;
	}

	.profile-dropdown-head {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 10px 10px 12px;
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
		.foot-trust span {
			display: none;
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

	@media (max-width: 640px) {
		.topbar-inner {
			padding: 10px 16px;
		}

		.topbar-search {
			display: none;
		}

		.profile-info {
			display: none;
		}

		.topbar-divider {
			display: none;
		}
	}
</style>
