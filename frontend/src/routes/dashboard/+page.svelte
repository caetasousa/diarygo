<script lang="ts">
	import { onMount } from 'svelte';
	import { currentUser } from '$lib/stores/auth';
	import { api } from '$lib/api/client';
	import type {
		ClienteResponse,
		EnderecoResponse,
		ProfissionalResponse,
		DocumentoResponse,
		ReferenciaResponse,
		RegiaoResponse,
		DisponibilidadeResponse
	} from '$lib/types';
	import {
		Plus,
		ArrowRight,
		TrendingUp,
		Sparkles,
		Star,
		Clock,
		ShieldCheck,
		CalendarDays,
		CircleCheck,
		Circle,
		MapPin,
		Users,
		User as UserIcon,
		Headset,
		Zap,
		House,
		Brush,
		HardHat,
		Package,
		Building2,
		Shirt,
		ChevronRight,
		ArrowUpRight,
		CircleAlert,
		BadgeCheck
	} from 'lucide-svelte';

	const weekDays = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb'];
	const today = new Date();
	const last7 = Array.from({ length: 7 }, (_, i) => {
		const d = new Date(today);
		d.setDate(today.getDate() - (6 - i));
		return {
			label: weekDays[d.getDay()],
			date: d.getDate(),
			isToday: i === 6,
			active: false
		};
	});

	const quickServices = [
		{ icon: House, name: 'Limpeza padrão', tag: '4h' },
		{ icon: Brush, name: 'Limpeza pesada', tag: '6h' },
		{ icon: Package, name: 'Pré-mudança', tag: '6h' },
		{ icon: HardHat, name: 'Pós-obra', tag: '8h' },
		{ icon: Building2, name: 'Comercial', tag: '4h' },
		{ icon: Shirt, name: 'Passadoria', tag: '2h' }
	];

	type ChecklistItem = {
		href: string;
		label: string;
		sub: string;
		done: boolean;
		icon: typeof UserIcon;
	};

	let clienteChecklist = $state<ChecklistItem[]>([
		{ href: '/dashboard/perfil', label: 'Complete seu perfil', sub: 'Nome, CPF e telefone', done: false, icon: UserIcon },
		{ href: '/dashboard/enderecos', label: 'Adicione um endereço', sub: 'Onde o serviço será realizado', done: false, icon: MapPin },
		{ href: '/solicitacoes/novo', label: 'Faça sua primeira solicitação', sub: 'Agende sua diarista', done: false, icon: Sparkles },
	]);

	let profissionalChecklist = $state<ChecklistItem[]>([
		{ href: '/dashboard/perfil', label: 'Dados pessoais', sub: 'Nome, CPF, RG e foto', done: false, icon: UserIcon },
		{ href: '/dashboard/credenciamento', label: 'Credenciamento', sub: 'Documentos e referências', done: false, icon: BadgeCheck },
		{ href: '/dashboard/atuacao', label: 'Atuação', sub: 'Regiões e disponibilidade', done: false, icon: MapPin },
	]);

	const tips = [
		{
			icon: ShieldCheck,
			title: 'Diaristas verificadas',
			desc: 'Toda profissional passa por checagem de documentos, referências e antecedentes.',
			color: 'green'
		},
		{
			icon: Zap,
			title: 'Agendamento rápido',
			desc: 'Encontre e confirme uma diarista disponível em menos de 5 minutos.',
			color: 'orange'
		},
		{
			icon: Star,
			title: 'Avalie e reavalie',
			desc: 'Após cada serviço, você avalia — e pode salvar suas profissionais favoritas.',
			color: 'blue'
		}
	];

	let tipo = $derived($currentUser?.tipo ?? 'CLIENTE');
	let firstName = $derived(($currentUser?.email ?? 'você').split('@')[0]);
	let hour = today.getHours();
	let greeting = $derived(hour < 12 ? 'Bom dia' : hour < 18 ? 'Boa tarde' : 'Boa noite');

	let clienteProgress = $derived(clienteChecklist.filter(i => i.done).length);
	let profissionalProgress = $derived(profissionalChecklist.filter(i => i.done).length);

	function hasText(v: string | undefined | null): boolean {
		return typeof v === 'string' && v.trim().length > 0;
	}

	async function carregarEstadoCliente() {
		const [cliente, enderecos] = await Promise.all([
			api.buscarPerfilCliente().catch(() => null as ClienteResponse | null),
			api.listarEnderecos().catch(() => [] as EnderecoResponse[])
		]);

		const perfilOk =
			!!cliente && hasText(cliente.nome) && hasText(cliente.cpf) && hasText(cliente.telefone);
		const enderecoOk = enderecos.length > 0;

		clienteChecklist[0].done = perfilOk;
		clienteChecklist[1].done = enderecoOk;
	}

	async function carregarEstadoProfissional() {
		const [profissional, documentos, referencias, regioes, disponibilidades] = await Promise.all([
			api.buscarPerfilProfissional().catch(() => null as ProfissionalResponse | null),
			api.listarDocumentos().catch(() => [] as DocumentoResponse[]),
			api.listarReferencias().catch(() => [] as ReferenciaResponse[]),
			api.listarRegioesAtuacao().catch(() => [] as RegiaoResponse[]),
			api.listarDisponibilidades().catch(() => [] as DisponibilidadeResponse[])
		]);

		const perfilOk =
			!!profissional &&
			hasText(profissional.nome) &&
			hasText(profissional.cpf) &&
			hasText(profissional.rg) &&
			hasText(profissional.telefone);
		const credenciamentoOk = documentos.length > 0 && referencias.length >= 2;
		const atuacaoOk = regioes.length > 0 && disponibilidades.length > 0;

		profissionalChecklist[0].done = perfilOk;
		profissionalChecklist[1].done = credenciamentoOk;
		profissionalChecklist[2].done = atuacaoOk;
	}

	onMount(() => {
		const tipoAtual = $currentUser?.tipo;
		if (tipoAtual === 'CLIENTE') {
			carregarEstadoCliente();
		} else if (tipoAtual === 'PROFISSIONAL') {
			carregarEstadoProfissional();
		}
	});

	function formatDate(d: Date) {
		return d.toLocaleDateString('pt-BR', { weekday: 'long', day: 'numeric', month: 'long' });
	}
</script>

<svelte:head>
	<title>Dashboard — DiaryGo</title>
</svelte:head>

<div class="page">

	{#if tipo === 'CLIENTE'}
		{@render ClienteView()}
	{:else if tipo === 'PROFISSIONAL'}
		{@render ProfissionalView()}
	{:else if tipo === 'ADMIN'}
		{@render AdminView()}
	{/if}

</div>


<!-- ════════════════════════════════════════════════════
     CLIENTE
════════════════════════════════════════════════════ -->
{#snippet ClienteView()}

	<!-- ─── Hero card ─── -->
	<section class="hero-card animate-in">
		<div class="hero-glow"></div>
		<div class="hero-body">
			<div class="hero-left">
				<span class="hero-date">{formatDate(today)}</span>
				<h1 class="hero-title">
					{greeting}, <span class="hero-accent">{firstName}</span>.
				</h1>
				<p class="hero-sub">
					Sua casa impecável começa com um clique. Agende uma diarista verificada para hoje, amanhã ou qualquer dia da semana.
				</p>
				<div class="hero-actions">
					<a href="/solicitacoes/novo" class="btn-primary-solid">
						<Plus size={15} strokeWidth={2.5} />
						Solicitar serviço
					</a>
					<a href="/solicitacoes" class="btn-ghost-pill">
						Ver meus serviços
						<ArrowRight size={14} strokeWidth={2} />
					</a>
				</div>
			</div>

			<!-- Hero visual: mini-app card -->
			<div class="hero-visual">
				<div class="mini-app">
					<div class="mini-app-head">
						<span class="mini-dot mini-dot-red"></span>
						<span class="mini-dot mini-dot-yellow"></span>
						<span class="mini-dot mini-dot-green"></span>
					</div>
					<div class="mini-app-body">
						<div class="mini-row">
							<div class="mini-icon-wrap mini-icon-orange">
								<House size={14} strokeWidth={2} />
							</div>
							<div class="mini-text">
								<span class="mini-label">Limpeza padrão</span>
								<span class="mini-sub">3 quartos · 2 banheiros</span>
							</div>
							<span class="mini-price">R$ 180</span>
						</div>
						<div class="mini-divider"></div>
						<div class="mini-row">
							<div class="mini-avatar">M</div>
							<div class="mini-text">
								<span class="mini-label">Maria Silva</span>
								<span class="mini-sub">
									<Star size={10} strokeWidth={2} fill="currentColor" />
									4.9 · a 2.3km
								</span>
							</div>
							<span class="mini-badge">Disponível</span>
						</div>
						<div class="mini-cta">Confirmar agendamento</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- ─── Stats row ─── -->
	<section class="stats-grid animate-in">
		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Serviços realizados</span>
				<div class="stat-icon stat-icon-green">
					<CircleCheck size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">0</span>
				<span class="stat-meta">
					<TrendingUp size={11} strokeWidth={2.2} />
					desde o início
				</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-green-10)')}
		</article>

		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Agendamentos ativos</span>
				<div class="stat-icon stat-icon-blue">
					<CalendarDays size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">0</span>
				<span class="stat-meta stat-meta-muted">sem agendamento</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-blue-10)')}
		</article>

		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Score de confiança</span>
				<div class="stat-icon stat-icon-orange">
					<ShieldCheck size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">100</span>
				<span class="stat-meta stat-meta-good">
					<TrendingUp size={11} strokeWidth={2.2} />
					nota máxima
				</span>
			</div>
			{@render Spark([70,72,78,82,85,90,100], 'var(--color-orange-10)')}
		</article>

		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Tempo médio resposta</span>
				<div class="stat-icon stat-icon-yellow">
					<Clock size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">—</span>
				<span class="stat-meta stat-meta-muted">sem dados</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-yellow-9)')}
		</article>
	</section>

	<!-- ─── Quick services ─── -->
	<section class="quick-section animate-in">
		<div class="section-bar">
			<div>
				<h2 class="section-title">Agende agora</h2>
				<p class="section-sub">Os serviços mais pedidos da semana</p>
			</div>
			<a href="/servicos" class="section-link-pill">
				Todos os serviços
				<ArrowUpRight size={13} strokeWidth={2} />
			</a>
		</div>

		<div class="quick-grid">
			{#each quickServices as svc}
				{@const Icon = svc.icon}
				<a href="/solicitacoes/novo?tipo={svc.name}" class="quick-card">
					<div class="quick-icon">
						<Icon size={18} strokeWidth={1.8} />
					</div>
					<div class="quick-body">
						<span class="quick-name">{svc.name}</span>
						<span class="quick-tag">A partir de {svc.tag}</span>
					</div>
					<ChevronRight size={14} strokeWidth={2} class="quick-arrow" />
				</a>
			{/each}
		</div>
	</section>

	<!-- ─── Main grid ─── -->
	<div class="main-grid">

		<!-- Coluna principal -->
		<div class="col-main">

			<!-- Próximo serviço -->
			<section class="panel animate-in">
				<header class="panel-head">
					<div>
						<h2 class="panel-title">Próximo serviço</h2>
						<p class="panel-sub">Acompanhe seu agendamento em tempo real</p>
					</div>
					<a href="/solicitacoes" class="panel-link">Ver todos<ArrowRight size={13} strokeWidth={2}/></a>
				</header>
				<div class="panel-body empty-illustration">
					<div class="illustration">
						<div class="illus-ring illus-ring-1"></div>
						<div class="illus-ring illus-ring-2"></div>
						<div class="illus-ring illus-ring-3"></div>
						<div class="illus-icon">
							<CalendarDays size={24} strokeWidth={1.5} />
						</div>
					</div>
					<div class="empty-text">
						<h3 class="empty-title">Você ainda não tem nenhum serviço</h3>
						<p class="empty-sub">
							Agende agora e encontre uma diarista disponível na sua região em minutos.
						</p>
					</div>
					<div class="empty-actions">
						<a href="/solicitacoes/novo" class="btn-primary-solid">
							<Plus size={14} strokeWidth={2.5} />
							Agendar primeiro serviço
						</a>
						<a href="/como-funciona" class="btn-ghost-pill">Como funciona</a>
					</div>
				</div>
			</section>

			<!-- Atividade recente -->
			<section class="panel animate-in">
				<header class="panel-head">
					<div>
						<h2 class="panel-title">Atividade recente</h2>
						<p class="panel-sub">Últimos eventos da sua conta</p>
					</div>
				</header>
				<div class="panel-body">
					<div class="timeline">
						<div class="timeline-item">
							<div class="timeline-dot timeline-dot-green">
								<CircleCheck size={11} strokeWidth={2.5} />
							</div>
							<div class="timeline-body">
								<div class="timeline-row">
									<span class="timeline-title">Conta criada</span>
									<span class="timeline-time">hoje</span>
								</div>
								<span class="timeline-sub">Bem-vindo ao DiaryGo — sua jornada começa aqui.</span>
							</div>
						</div>
						<div class="timeline-item timeline-future">
							<div class="timeline-dot timeline-dot-muted">
								<Circle size={9} strokeWidth={2} />
							</div>
							<div class="timeline-body">
								<div class="timeline-row">
									<span class="timeline-title">Complete seu perfil</span>
									<span class="timeline-time">pendente</span>
								</div>
								<span class="timeline-sub">Necessário para solicitar o primeiro serviço.</span>
							</div>
						</div>
						<div class="timeline-item timeline-future">
							<div class="timeline-dot timeline-dot-muted">
								<Circle size={9} strokeWidth={2} />
							</div>
							<div class="timeline-body">
								<div class="timeline-row">
									<span class="timeline-title">Primeira solicitação</span>
									<span class="timeline-time">—</span>
								</div>
								<span class="timeline-sub">Agende sua primeira diarista.</span>
							</div>
						</div>
					</div>
				</div>
			</section>

			<!-- Por que confiar -->
			<section class="panel animate-in">
				<header class="panel-head">
					<div>
						<h2 class="panel-title">Por que DiaryGo</h2>
						<p class="panel-sub">Os três pilares da sua tranquilidade</p>
					</div>
				</header>
				<div class="panel-body">
					<div class="tips-grid">
						{#each tips as tip}
							{@const Icon = tip.icon}
							<article class="tip-card tip-{tip.color}">
								<div class="tip-icon">
									<Icon size={16} strokeWidth={1.8} />
								</div>
								<h3 class="tip-title">{tip.title}</h3>
								<p class="tip-desc">{tip.desc}</p>
							</article>
						{/each}
					</div>
				</div>
			</section>

		</div>

		<!-- Coluna lateral -->
		<aside class="col-side">

			<!-- Setup progress -->
			<section class="side-panel side-panel-accent animate-in">
				<div class="progress-head">
					<div>
						<h3 class="side-title">Configure sua conta</h3>
						<p class="side-sub">{clienteProgress} de {clienteChecklist.length} etapas</p>
					</div>
					<div class="progress-ring">
						<svg width="44" height="44" viewBox="0 0 44 44">
							<circle cx="22" cy="22" r="18" stroke="rgba(255,255,255,0.08)" stroke-width="3" fill="none"/>
							<circle
								cx="22" cy="22" r="18"
								stroke="var(--color-orange-10)"
								stroke-width="3"
								fill="none"
								stroke-linecap="round"
								stroke-dasharray={2 * Math.PI * 18}
								stroke-dashoffset={2 * Math.PI * 18 * (1 - clienteProgress / clienteChecklist.length)}
								transform="rotate(-90 22 22)"
							/>
						</svg>
						<span class="progress-percent">{Math.round(clienteProgress / clienteChecklist.length * 100)}%</span>
					</div>
				</div>
				<div class="checklist">
					{#each clienteChecklist as item}
						{@const Icon = item.icon}
						<a href={item.href} class="check-row" class:done={item.done}>
							<div class="check-mark">
								{#if item.done}
									<CircleCheck size={16} strokeWidth={2} />
								{:else}
									<Icon size={13} strokeWidth={1.8} />
								{/if}
							</div>
							<div class="check-text">
								<span class="check-label">{item.label}</span>
								<span class="check-sub">{item.sub}</span>
							</div>
							<ChevronRight size={13} strokeWidth={2} class="check-arrow" />
						</a>
					{/each}
				</div>
			</section>

			<!-- Esta semana -->
			<section class="side-panel animate-in">
				<div class="side-head-row">
					<h3 class="side-title">Esta semana</h3>
					<span class="side-badge">0 serviços</span>
				</div>
				<div class="week-strip">
					{#each last7 as day}
						<div class="week-cell" class:is-today={day.isToday}>
							<span class="week-lbl">{day.label}</span>
							<div class="week-box" class:filled={day.active}>
								<span class="week-num">{day.date}</span>
							</div>
						</div>
					{/each}
				</div>
				<p class="side-note">
					Nenhum serviço nos últimos 7 dias. Seu histórico semanal aparecerá aqui.
				</p>
			</section>

			<!-- Suporte -->
			<section class="side-panel support-panel animate-in">
				<div class="support-icon-wrap">
					<Headset size={18} strokeWidth={1.8} />
				</div>
				<h3 class="side-title">Precisa de ajuda?</h3>
				<p class="side-note support-note">
					Nosso time responde em até 2 horas nos dias úteis.
				</p>
				<a href="/ajuda" class="btn-ghost-pill support-btn">
					Falar com suporte
					<ArrowRight size={13} strokeWidth={2} />
				</a>
			</section>

		</aside>
	</div>

{/snippet}


<!-- ════════════════════════════════════════════════════
     PROFISSIONAL
════════════════════════════════════════════════════ -->
{#snippet ProfissionalView()}

	<!-- Hero profissional -->
	<section class="hero-card animate-in">
		<div class="hero-glow"></div>
		<div class="hero-body">
			<div class="hero-left">
				<span class="hero-date">{formatDate(today)}</span>
				<h1 class="hero-title">
					{greeting}, <span class="hero-accent">{firstName}</span>.
				</h1>
				<p class="hero-sub">
					Complete seu cadastro para começar a receber solicitações na sua região. Nossa equipe avalia em até 48h úteis.
				</p>
			</div>
		</div>
	</section>

	<!-- Banner de aprovação -->
	<section class="approval-banner animate-in">
		<div class="approval-icon">
			<Clock size={16} strokeWidth={2} />
		</div>
		<div class="approval-body">
			<span class="approval-title">Cadastro em análise</span>
			<span class="approval-sub">
				Sua documentação está sendo revisada pela nossa equipe de qualidade. Você receberá um e-mail assim que for aprovada.
			</span>
		</div>
		<span class="approval-badge">
			<span class="approval-pulse"></span>
			Pendente
		</span>
	</section>

	<!-- Stats profissional -->
	<section class="stats-grid animate-in">
		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Serviços realizados</span>
				<div class="stat-icon stat-icon-green">
					<CircleCheck size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">0</span>
				<span class="stat-meta stat-meta-muted">aguardando aprovação</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-green-10)')}
		</article>

		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Avaliação média</span>
				<div class="stat-icon stat-icon-yellow">
					<Star size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">—</span>
				<span class="stat-meta stat-meta-muted">sem avaliações</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-yellow-9)')}
		</article>

		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Ganhos no mês</span>
				<div class="stat-icon stat-icon-blue">
					<TrendingUp size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">R$ 0</span>
				<span class="stat-meta stat-meta-muted">sem serviços</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-blue-10)')}
		</article>

		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Status</span>
				<div class="stat-icon stat-icon-orange">
					<CircleAlert size={14} strokeWidth={2} />
				</div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value stat-orange" style="font-size: 1.125rem;">EM ANÁLISE</span>
			</div>
			<p class="stat-foot">Aprovação em até 48h úteis</p>
		</article>
	</section>

	<!-- Main grid -->
	<div class="main-grid">

		<div class="col-main">
			<!-- Próximas solicitações -->
			<section class="panel animate-in">
				<header class="panel-head">
					<div>
						<h2 class="panel-title">Próximas solicitações</h2>
						<p class="panel-sub">Serviços disponíveis para você</p>
					</div>
				</header>
				<div class="panel-body empty-illustration">
					<div class="illustration">
						<div class="illus-ring illus-ring-1"></div>
						<div class="illus-ring illus-ring-2"></div>
						<div class="illus-ring illus-ring-3"></div>
						<div class="illus-icon">
							<CalendarDays size={24} strokeWidth={1.5} />
						</div>
					</div>
					<div class="empty-text">
						<h3 class="empty-title">Aguardando aprovação do cadastro</h3>
						<p class="empty-sub">
							Após aprovada, você começa a receber solicitações compatíveis com suas regiões e horários.
						</p>
					</div>
				</div>
			</section>

			<!-- Dicas para novas profissionais -->
			<section class="panel animate-in">
				<header class="panel-head">
					<div>
						<h2 class="panel-title">Dicas para começar bem</h2>
						<p class="panel-sub">Profissionais bem avaliadas fazem isso</p>
					</div>
				</header>
				<div class="panel-body">
					<div class="tips-grid">
						<article class="tip-card tip-green">
							<div class="tip-icon">
								<Star size={16} strokeWidth={1.8} />
							</div>
							<h3 class="tip-title">Capriche nas primeiras avaliações</h3>
							<p class="tip-desc">As três primeiras notas pesam muito no algoritmo de distribuição.</p>
						</article>
						<article class="tip-card tip-orange">
							<div class="tip-icon">
								<Clock size={16} strokeWidth={1.8} />
							</div>
							<h3 class="tip-title">Responda rápido</h3>
							<p class="tip-desc">Aceite em até 30 minutos — depois, a solicitação vai para outra.</p>
						</article>
						<article class="tip-card tip-blue">
							<div class="tip-icon">
								<MapPin size={16} strokeWidth={1.8} />
							</div>
							<h3 class="tip-title">Amplie suas regiões</h3>
							<p class="tip-desc">Quanto mais bairros atende, mais serviços chegam até você.</p>
						</article>
					</div>
				</div>
			</section>
		</div>

		<aside class="col-side">
			<!-- Checklist cadastro -->
			<section class="side-panel side-panel-accent animate-in">
				<div class="progress-head">
					<div>
						<h3 class="side-title">Completar cadastro</h3>
						<p class="side-sub">{profissionalProgress} de {profissionalChecklist.length} etapas</p>
					</div>
					<div class="progress-ring">
						<svg width="44" height="44" viewBox="0 0 44 44">
							<circle cx="22" cy="22" r="18" stroke="rgba(255,255,255,0.08)" stroke-width="3" fill="none"/>
							<circle
								cx="22" cy="22" r="18"
								stroke="var(--color-orange-10)"
								stroke-width="3"
								fill="none"
								stroke-linecap="round"
								stroke-dasharray={2 * Math.PI * 18}
								stroke-dashoffset={2 * Math.PI * 18 * (1 - profissionalProgress / profissionalChecklist.length)}
								transform="rotate(-90 22 22)"
							/>
						</svg>
						<span class="progress-percent">{Math.round(profissionalProgress / profissionalChecklist.length * 100)}%</span>
					</div>
				</div>
				<div class="checklist">
					{#each profissionalChecklist as item}
						{@const Icon = item.icon}
						<a href={item.href} class="check-row" class:done={item.done}>
							<div class="check-mark">
								{#if item.done}
									<CircleCheck size={16} strokeWidth={2} />
								{:else}
									<Icon size={13} strokeWidth={1.8} />
								{/if}
							</div>
							<div class="check-text">
								<span class="check-label">{item.label}</span>
								<span class="check-sub">{item.sub}</span>
							</div>
							<ChevronRight size={13} strokeWidth={2} class="check-arrow" />
						</a>
					{/each}
				</div>
			</section>

			<!-- Semana -->
			<section class="side-panel animate-in">
				<div class="side-head-row">
					<h3 class="side-title">Esta semana</h3>
					<span class="side-badge">0 serviços</span>
				</div>
				<div class="week-strip">
					{#each last7 as day}
						<div class="week-cell" class:is-today={day.isToday}>
							<span class="week-lbl">{day.label}</span>
							<div class="week-box" class:filled={day.active}>
								<span class="week-num">{day.date}</span>
							</div>
						</div>
					{/each}
				</div>
				<p class="side-note">
					Seu histórico de atendimentos aparecerá aqui.
				</p>
			</section>

			<!-- Suporte -->
			<section class="side-panel support-panel animate-in">
				<div class="support-icon-wrap">
					<Headset size={18} strokeWidth={1.8} />
				</div>
				<h3 class="side-title">Dúvidas sobre cadastro?</h3>
				<p class="side-note support-note">
					Nosso time responde em até 2 horas nos dias úteis.
				</p>
				<a href="/ajuda" class="btn-ghost-pill support-btn">
					Falar com suporte
					<ArrowRight size={13} strokeWidth={2} />
				</a>
			</section>
		</aside>
	</div>

{/snippet}


<!-- ════════════════════════════════════════════════════
     ADMIN
════════════════════════════════════════════════════ -->
{#snippet AdminView()}
	<section class="hero-card animate-in">
		<div class="hero-glow"></div>
		<div class="hero-body">
			<div class="hero-left">
				<span class="hero-date">{formatDate(today)}</span>
				<h1 class="hero-title">Painel administrativo</h1>
				<p class="hero-sub">Visão geral da plataforma em tempo real.</p>
			</div>
		</div>
	</section>

	<section class="stats-grid animate-in">
		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Profissionais ativas</span>
				<div class="stat-icon stat-icon-green"><Users size={14} strokeWidth={2}/></div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">—</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-green-10)')}
		</article>
		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Clientes</span>
				<div class="stat-icon stat-icon-blue"><UserIcon size={14} strokeWidth={2}/></div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">—</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-blue-10)')}
		</article>
		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Serviços hoje</span>
				<div class="stat-icon stat-icon-orange"><CalendarDays size={14} strokeWidth={2}/></div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">—</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-orange-10)')}
		</article>
		<article class="stat-card">
			<div class="stat-top">
				<span class="stat-label">Pendentes</span>
				<div class="stat-icon stat-icon-yellow"><CircleAlert size={14} strokeWidth={2}/></div>
			</div>
			<div class="stat-value-row">
				<span class="stat-value">—</span>
			</div>
			{@render Spark([0,0,0,0,0,0,0], 'var(--color-yellow-9)')}
		</article>
	</section>

	<section class="panel animate-in">
		<header class="panel-head">
			<h2 class="panel-title">Acesso rápido</h2>
		</header>
		<div class="panel-body admin-grid">
			<a href="/admin/profissionais" class="admin-card">
				<Users size={20} strokeWidth={1.6}/>
				<span class="admin-card-label">Profissionais</span>
				<span class="admin-card-sub">Gerenciar cadastros e aprovações</span>
			</a>
			<a href="/admin/clientes" class="admin-card">
				<UserIcon size={20} strokeWidth={1.6}/>
				<span class="admin-card-label">Clientes</span>
				<span class="admin-card-sub">Ver, editar e suporte</span>
			</a>
			<a href="/admin/servicos" class="admin-card">
				<CalendarDays size={20} strokeWidth={1.6}/>
				<span class="admin-card-label">Serviços</span>
				<span class="admin-card-sub">Acompanhar agendamentos</span>
			</a>
			<a href="/admin/pendentes" class="admin-card">
				<CircleAlert size={20} strokeWidth={1.6}/>
				<span class="admin-card-label">Fila de aprovação</span>
				<span class="admin-card-sub">Revisar documentação</span>
			</a>
		</div>
	</section>
{/snippet}


<!-- ════════════════════════════════════════════════════
     Sparkline helper
════════════════════════════════════════════════════ -->
{#snippet Spark(values: number[], color: string)}
	{@const max = Math.max(...values, 1)}
	{@const w = 100}
	{@const h = 24}
	{@const step = w / (values.length - 1)}
	{@const points = values.map((v, i) => `${i * step},${h - (v / max) * h}`).join(' ')}
	<svg class="spark" viewBox="0 0 {w} {h}" preserveAspectRatio="none">
		<defs>
			<linearGradient id="g-{color}" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%" stop-color={color} stop-opacity="0.25"/>
				<stop offset="100%" stop-color={color} stop-opacity="0"/>
			</linearGradient>
		</defs>
		<polyline
			points={points}
			fill="none"
			stroke={color}
			stroke-width="1.3"
			stroke-linecap="round"
			stroke-linejoin="round"
			opacity="0.8"
		/>
		<polygon
			points="0,{h} {points} {w},{h}"
			fill="url(#g-{color})"
		/>
	</svg>
{/snippet}


<style>
	/* ── Page shell ──────────────────────────────────── */
	.page {
		padding: 28px 32px 56px;
		max-width: 1280px;
		margin: 0 auto;
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	/* ══════════════════════════════════════════════════
	   Hero card
	   ══════════════════════════════════════════════════ */
	.hero-card {
		position: relative;
		overflow: hidden;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-large);
		background:
			radial-gradient(120% 80% at 0% 0%, rgba(255, 128, 31, 0.06), transparent 50%),
			linear-gradient(180deg, rgba(255,255,255,0.025), rgba(255,255,255,0.008));
		box-shadow: var(--shadow-ring);
	}

	.hero-glow {
		position: absolute;
		top: -100px;
		left: -100px;
		width: 400px;
		height: 400px;
		background: radial-gradient(circle, rgba(255, 128, 31, 0.12) 0%, transparent 60%);
		pointer-events: none;
		filter: blur(20px);
	}

	.hero-body {
		position: relative;
		display: grid;
		grid-template-columns: 1fr auto;
		gap: 32px;
		padding: 32px 36px;
		align-items: center;
	}

	.hero-left {
		display: flex;
		flex-direction: column;
		gap: 10px;
		max-width: 560px;
	}

	.hero-date {
		font-size: 0.6875rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.6px;
	}

	.hero-title {
		font-size: clamp(1.75rem, 3.5vw, 2.5rem);
		font-weight: 400;
		letter-spacing: -1.5px;
		line-height: 1.05;
		color: var(--color-text-primary);
	}

	.hero-accent {
		background: linear-gradient(135deg, var(--color-orange-11), var(--color-orange-10));
		-webkit-background-clip: text;
		background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.hero-sub {
		font-size: 0.9375rem;
		color: var(--color-text-secondary);
		max-width: 520px;
		line-height: 1.55;
	}

	.hero-actions {
		display: flex;
		align-items: center;
		gap: 10px;
		margin-top: 8px;
		flex-wrap: wrap;
	}

	/* Hero visual — mini app mockup */
	.hero-visual {
		position: relative;
	}

	.mini-app {
		width: 280px;
		background: #0a0a0a;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		overflow: hidden;
		box-shadow:
			0 20px 60px rgba(0, 0, 0, 0.5),
			0 0 0 1px rgba(176, 199, 217, 0.08);
		transform: perspective(1000px) rotateY(-6deg) rotateX(2deg);
		transition: transform 0.4s ease;
	}

	.mini-app:hover {
		transform: perspective(1000px) rotateY(-3deg) rotateX(0deg) scale(1.02);
	}

	.mini-app-head {
		display: flex;
		gap: 5px;
		padding: 9px 12px;
		border-bottom: 1px solid var(--border-frost);
	}

	.mini-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
	}
	.mini-dot-red    { background: #ff5f57; }
	.mini-dot-yellow { background: #febc2e; }
	.mini-dot-green  { background: #28c840; }

	.mini-app-body {
		padding: 14px;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.mini-row {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.mini-icon-wrap {
		width: 30px;
		height: 30px;
		border-radius: var(--radius-standard);
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.mini-icon-orange {
		background: var(--color-orange-4);
		color: var(--color-orange-11);
	}

	.mini-avatar {
		width: 30px;
		height: 30px;
		border-radius: 50%;
		background: linear-gradient(135deg, #3b9eff, #0075ff);
		color: #fff;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.75rem;
		font-weight: 600;
		flex-shrink: 0;
	}

	.mini-text {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 1px;
		min-width: 0;
	}

	.mini-label {
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.mini-sub {
		font-size: 0.625rem;
		color: var(--color-text-tertiary);
		display: inline-flex;
		align-items: center;
		gap: 3px;
	}

	.mini-sub :global(svg) {
		color: var(--color-yellow-9);
	}

	.mini-price {
		font-size: 0.8125rem;
		font-weight: 600;
		color: var(--color-text-primary);
		font-feature-settings: 'tnum';
	}

	.mini-badge {
		font-size: 0.5625rem;
		font-weight: 600;
		color: var(--color-green-10);
		background: var(--color-green-3);
		padding: 2px 7px;
		border-radius: var(--radius-pill);
	}

	.mini-divider {
		height: 1px;
		background: var(--border-frost);
	}

	.mini-cta {
		margin-top: 4px;
		background: var(--color-white);
		color: var(--color-black);
		font-size: 0.75rem;
		font-weight: 600;
		padding: 9px;
		border-radius: var(--radius-standard);
		text-align: center;
	}

	/* Primary solid button */
	.btn-primary-solid {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 9px 16px;
		background: var(--color-white);
		color: var(--color-black);
		border: none;
		border-radius: var(--radius-pill);
		font-size: 0.8125rem;
		font-weight: 600;
		letter-spacing: -0.1px;
		text-decoration: none;
		cursor: pointer;
		transition: opacity 0.15s, transform 0.15s;
	}

	.btn-primary-solid:hover {
		opacity: 0.92;
		transform: translateY(-1px);
	}

	.btn-ghost-pill {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 9px 14px;
		background: transparent;
		color: var(--color-text-primary);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		font-size: 0.8125rem;
		font-weight: 500;
		text-decoration: none;
		transition: background 0.15s, border-color 0.15s;
	}

	.btn-ghost-pill:hover {
		background: var(--bg-hover-subtle);
		border-color: var(--border-frost-hover);
		opacity: 1;
	}

	/* ══════════════════════════════════════════════════
	   Stats grid
	   ══════════════════════════════════════════════════ */
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 12px;
	}

	.stat-card {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		padding: 16px 18px 4px;
		display: flex;
		flex-direction: column;
		gap: 10px;
		box-shadow: var(--shadow-ring);
		transition: border-color 0.15s, transform 0.15s;
	}

	.stat-card:hover {
		border-color: var(--border-frost-hover);
		transform: translateY(-1px);
	}

	.stat-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.stat-label {
		font-size: 0.6875rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.stat-icon {
		width: 26px;
		height: 26px;
		border-radius: var(--radius-standard);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.stat-icon-green  { background: var(--color-green-3); color: var(--color-green-10); }
	.stat-icon-blue   { background: var(--color-blue-4);  color: var(--color-blue-10); }
	.stat-icon-orange { background: var(--color-orange-4); color: var(--color-orange-11); }
	.stat-icon-yellow { background: rgba(255,197,61,0.15); color: var(--color-yellow-9); }

	.stat-value-row {
		display: flex;
		align-items: baseline;
		gap: 8px;
	}

	.stat-value {
		font-size: 1.875rem;
		font-weight: 400;
		letter-spacing: -1.5px;
		line-height: 1;
		color: var(--color-text-primary);
		font-feature-settings: 'tnum';
	}

	.stat-orange { color: var(--color-orange-10) !important; }

	.stat-meta {
		display: inline-flex;
		align-items: center;
		gap: 3px;
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		font-weight: 500;
	}

	.stat-meta-good { color: var(--color-green-10); }
	.stat-meta-muted { color: var(--color-text-tertiary); }

	.stat-foot {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	.spark {
		width: 100%;
		height: 24px;
		margin-top: 2px;
	}

	/* ══════════════════════════════════════════════════
	   Quick services
	   ══════════════════════════════════════════════════ */
	.quick-section {
		display: flex;
		flex-direction: column;
		gap: 14px;
	}

	.section-bar {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: 16px;
	}

	.section-title {
		font-size: 1.125rem;
		font-weight: 500;
		letter-spacing: -0.5px;
		color: var(--color-text-primary);
	}

	.section-sub {
		font-size: 0.8125rem;
		color: var(--color-text-tertiary);
		margin-top: 2px;
	}

	.section-link-pill {
		display: inline-flex;
		align-items: center;
		gap: 5px;
		padding: 6px 12px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		font-size: 0.75rem;
		color: var(--color-text-secondary);
		text-decoration: none;
		transition: color 0.12s, background 0.12s, border-color 0.12s;
	}

	.section-link-pill:hover {
		color: var(--color-text-primary);
		background: var(--bg-hover-medium);
		border-color: var(--border-frost-hover);
		opacity: 1;
	}

	.quick-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(210px, 1fr));
		gap: 10px;
	}

	.quick-card {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 13px 14px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		text-decoration: none;
		transition: border-color 0.15s, transform 0.15s, background 0.15s;
	}

	.quick-card:hover {
		border-color: var(--border-frost-hover);
		background: var(--bg-hover-subtle);
		transform: translateY(-1px);
		opacity: 1;
	}

	.quick-icon {
		width: 34px;
		height: 34px;
		border-radius: var(--radius-standard);
		background: var(--bg-hover-subtle);
		border: 1px solid var(--border-frost);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-primary);
		flex-shrink: 0;
	}

	.quick-card:hover .quick-icon {
		background: var(--color-orange-4);
		border-color: var(--color-orange-10);
		color: var(--color-orange-11);
	}

	.quick-body {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 1px;
		min-width: 0;
	}

	.quick-name {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.quick-tag {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	:global(.quick-arrow) {
		color: var(--color-text-tertiary);
		flex-shrink: 0;
		transition: transform 0.15s, color 0.15s;
	}

	.quick-card:hover :global(.quick-arrow) {
		color: var(--color-text-primary);
		transform: translateX(2px);
	}

	/* ══════════════════════════════════════════════════
	   Main grid
	   ══════════════════════════════════════════════════ */
	.main-grid {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 320px;
		gap: 16px;
		align-items: start;
	}

	.col-main {
		display: flex;
		flex-direction: column;
		gap: 16px;
		min-width: 0;
	}

	.col-side {
		display: flex;
		flex-direction: column;
		gap: 14px;
	}

	/* Panels */
	.panel {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		overflow: hidden;
		box-shadow: var(--shadow-ring);
	}

	.panel-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		padding: 16px 20px 14px;
		border-bottom: 1px solid var(--border-frost);
	}

	.panel-title {
		font-size: 0.9375rem;
		font-weight: 500;
		letter-spacing: -0.3px;
		color: var(--color-text-primary);
	}

	.panel-sub {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		margin-top: 2px;
	}

	.panel-link {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		text-decoration: none;
		padding: 4px 8px;
		border-radius: var(--radius-subtle);
		transition: color 0.12s, background 0.12s;
	}

	.panel-link:hover {
		color: var(--color-text-primary);
		background: var(--bg-hover-subtle);
		opacity: 1;
	}

	.panel-body {
		padding: 20px;
	}

	/* Empty illustration */
	.empty-illustration {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		padding: 44px 24px 40px;
		gap: 18px;
	}

	.illustration {
		position: relative;
		width: 120px;
		height: 120px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.illus-ring {
		position: absolute;
		top: 50%;
		left: 50%;
		border-radius: 50%;
		border: 1px solid var(--border-frost);
		transform: translate(-50%, -50%);
	}

	.illus-ring-1 {
		width: 120px;
		height: 120px;
		opacity: 0.35;
		animation: ringPulse 3s ease-in-out infinite;
	}

	.illus-ring-2 {
		width: 88px;
		height: 88px;
		opacity: 0.6;
		animation: ringPulse 3s ease-in-out infinite 0.3s;
	}

	.illus-ring-3 {
		width: 58px;
		height: 58px;
		opacity: 0.9;
		animation: ringPulse 3s ease-in-out infinite 0.6s;
	}

	@keyframes ringPulse {
		0%, 100% { opacity: var(--ring-opacity, 0.4); transform: translate(-50%, -50%) scale(1); }
		50% { opacity: 0.15; transform: translate(-50%, -50%) scale(1.05); }
	}

	.illus-icon {
		position: relative;
		width: 44px;
		height: 44px;
		border-radius: var(--radius-card);
		background: linear-gradient(135deg, var(--color-orange-4), rgba(255, 128, 31, 0.04));
		border: 1px solid var(--color-orange-4);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-orange-11);
	}

	.empty-text {
		display: flex;
		flex-direction: column;
		gap: 4px;
		max-width: 360px;
	}

	.empty-title {
		font-size: 1rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.2px;
	}

	.empty-sub {
		font-size: 0.8125rem;
		color: var(--color-text-tertiary);
		line-height: 1.55;
	}

	.empty-actions {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
		justify-content: center;
	}

	/* Timeline */
	.timeline {
		display: flex;
		flex-direction: column;
		gap: 2px;
		position: relative;
	}

	.timeline-item {
		display: flex;
		gap: 14px;
		padding: 10px 0;
		position: relative;
	}

	.timeline-item:not(:last-child)::after {
		content: '';
		position: absolute;
		left: 9px;
		top: 28px;
		bottom: -10px;
		width: 1px;
		background: var(--border-frost);
	}

	.timeline-dot {
		width: 20px;
		height: 20px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		margin-top: 2px;
		z-index: 1;
	}

	.timeline-dot-green {
		background: var(--color-green-3);
		color: var(--color-green-10);
		border: 1px solid var(--color-green-4);
	}

	.timeline-dot-muted {
		background: rgba(255, 255, 255, 0.03);
		color: var(--color-text-tertiary);
		border: 1px solid var(--border-frost);
	}

	.timeline-body {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.timeline-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
	}

	.timeline-title {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.timeline-future .timeline-title {
		color: var(--color-text-secondary);
	}

	.timeline-time {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		font-weight: 500;
	}

	.timeline-sub {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		line-height: 1.5;
	}

	/* Tips */
	.tips-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 10px;
	}

	.tip-card {
		padding: 16px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.tip-icon {
		width: 30px;
		height: 30px;
		border-radius: var(--radius-standard);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 6px;
	}

	.tip-green .tip-icon  { background: var(--color-green-3);  color: var(--color-green-10); }
	.tip-orange .tip-icon { background: var(--color-orange-4); color: var(--color-orange-11); }
	.tip-blue .tip-icon   { background: var(--color-blue-4);   color: var(--color-blue-10); }

	.tip-title {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.2px;
	}

	.tip-desc {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		line-height: 1.5;
	}

	/* ══════════════════════════════════════════════════
	   Side panels
	   ══════════════════════════════════════════════════ */
	.side-panel {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		padding: 16px 18px;
		box-shadow: var(--shadow-ring);
	}

	.side-panel-accent {
		background:
			radial-gradient(80% 60% at 100% 0%, rgba(255, 128, 31, 0.05), transparent 60%),
			var(--bg-card);
	}

	.progress-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		margin-bottom: 14px;
	}

	.side-title {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.2px;
	}

	.side-sub {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		margin-top: 2px;
	}

	.progress-ring {
		position: relative;
		width: 44px;
		height: 44px;
		flex-shrink: 0;
	}

	.progress-percent {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.6875rem;
		font-weight: 600;
		color: var(--color-text-primary);
		font-feature-settings: 'tnum';
	}

	/* Checklist */
	.checklist {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.check-row {
		display: flex;
		align-items: center;
		gap: 11px;
		padding: 9px 8px;
		border-radius: var(--radius-standard);
		text-decoration: none;
		transition: background 0.12s;
	}

	.check-row:hover {
		background: var(--bg-hover-subtle);
		opacity: 1;
	}

	.check-mark {
		width: 24px;
		height: 24px;
		border-radius: var(--radius-standard);
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-secondary);
		flex-shrink: 0;
	}

	.check-row.done .check-mark {
		background: var(--color-green-3);
		border-color: var(--color-green-4);
		color: var(--color-green-10);
	}

	.check-text {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 1px;
		min-width: 0;
	}

	.check-label {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.1px;
	}

	.check-row.done .check-label {
		color: var(--color-text-secondary);
		text-decoration: line-through;
	}

	.check-sub {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	:global(.check-arrow) {
		color: var(--color-text-tertiary);
		opacity: 0;
		transition: opacity 0.12s, transform 0.12s;
		flex-shrink: 0;
	}

	.check-row:hover :global(.check-arrow) {
		opacity: 1;
		transform: translateX(2px);
	}

	/* Side head */
	.side-head-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 12px;
	}

	.side-badge {
		font-size: 0.625rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
		background: var(--bg-hover-subtle);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		padding: 2px 8px;
		letter-spacing: 0.2px;
	}

	/* Week strip */
	.week-strip {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: 4px;
		margin-bottom: 12px;
	}

	.week-cell {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 4px;
	}

	.week-lbl {
		font-size: 0.5625rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.3px;
	}

	.week-box {
		width: 100%;
		aspect-ratio: 1;
		border-radius: var(--radius-subtle);
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.6875rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
		transition: all 0.15s;
	}

	.week-cell.is-today .week-box {
		background: var(--color-orange-4);
		border-color: var(--color-orange-10);
		color: var(--color-orange-11);
	}

	.week-box.filled {
		background: var(--color-green-3);
		border-color: var(--color-green-4);
		color: var(--color-green-10);
	}

	.week-num {
		font-feature-settings: 'tnum';
	}

	.side-note {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		line-height: 1.5;
	}

	/* Support panel */
	.support-panel {
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		padding: 22px 18px;
		background:
			radial-gradient(60% 50% at 50% 0%, rgba(59, 158, 255, 0.04), transparent 70%),
			var(--bg-card);
	}

	.support-icon-wrap {
		width: 38px;
		height: 38px;
		border-radius: var(--radius-card);
		background: var(--color-blue-4);
		color: var(--color-blue-10);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 4px;
	}

	.support-note {
		max-width: 220px;
		margin: 0 auto;
	}

	.support-btn {
		margin-top: 4px;
	}

	/* ══════════════════════════════════════════════════
	   Approval banner (profissional)
	   ══════════════════════════════════════════════════ */
	.approval-banner {
		display: flex;
		align-items: center;
		gap: 14px;
		padding: 14px 20px;
		background:
			linear-gradient(90deg, var(--color-orange-4), rgba(255, 128, 31, 0.02));
		border: 1px solid var(--color-orange-4);
		border-radius: var(--radius-card);
	}

	.approval-icon {
		width: 34px;
		height: 34px;
		border-radius: var(--radius-standard);
		background: var(--color-orange-4);
		color: var(--color-orange-10);
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.approval-body {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 2px;
		min-width: 0;
	}

	.approval-title {
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.approval-sub {
		font-size: 0.75rem;
		color: var(--color-text-secondary);
		line-height: 1.5;
	}

	.approval-badge {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 5px 12px;
		background: var(--color-orange-4);
		color: var(--color-orange-11);
		border-radius: var(--radius-pill);
		font-size: 0.6875rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.4px;
		flex-shrink: 0;
	}

	.approval-pulse {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-orange-10);
		animation: pulse 2s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.4; }
	}

	/* ══════════════════════════════════════════════════
	   Admin grid
	   ══════════════════════════════════════════════════ */
	.admin-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 10px;
		padding: 20px;
	}

	.admin-card {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 18px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		text-decoration: none;
		color: var(--color-text-primary);
		transition: border-color 0.15s, transform 0.15s, background 0.15s;
	}

	.admin-card:hover {
		border-color: var(--border-frost-hover);
		background: var(--bg-hover-subtle);
		transform: translateY(-1px);
		opacity: 1;
	}

	.admin-card :global(svg) {
		color: var(--color-orange-11);
	}

	.admin-card-label {
		font-size: 0.9375rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.admin-card-sub {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
	}

	/* ══════════════════════════════════════════════════
	   Responsive
	   ══════════════════════════════════════════════════ */
	@media (max-width: 1100px) {
		.main-grid {
			grid-template-columns: 1fr;
		}

		.hero-visual {
			display: none;
		}

		.hero-body {
			grid-template-columns: 1fr;
		}

		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}

	@media (max-width: 640px) {
		.page {
			padding: 20px 16px 40px;
			gap: 14px;
		}

		.hero-body {
			padding: 24px 20px;
		}

		.hero-title {
			font-size: 1.75rem;
		}

		.stats-grid {
			grid-template-columns: 1fr;
		}

		.panel-body {
			padding: 16px;
		}

		.empty-illustration {
			padding: 32px 16px;
		}
	}
</style>
