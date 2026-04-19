<script lang="ts">
	import { isAuthenticated } from '$lib/stores/auth';
	import {
		ArrowRight,
		ArrowUpRight,
		Star,
		ShieldCheck,
		BadgeCheck,
		Clock,
		MapPin,
		Sparkles,
		CheckCircle,
		Zap,
		Heart,
		DollarSign,
		CalendarCheck,
		TrendingUp,
		House,
		Brush,
		HardHat,
		Package,
		Building2,
		Shirt,
		Search,
		UserCheck,
		Play,
		Plus,
		Minus,
		PiggyBank,
		HeartHandshake,
		Headphones
	} from 'lucide-svelte';

	let openFaq = $state<number | null>(0);

	function toggleFaq(i: number) {
		openFaq = openFaq === i ? null : i;
	}

	const heroStats = [
		{ value: '4.9', label: 'Avaliação média', icon: Star, color: 'yellow' },
		{ value: '+180k', label: 'Serviços realizados', icon: CheckCircle, color: 'green' },
		{ value: '500+', label: 'Diaristas ativas', icon: UserCheck, color: 'blue' },
		{ value: '48', label: 'Cidades atendidas', icon: MapPin, color: 'orange' }
	];

	const services = [
		{
			icon: House, name: 'Limpeza padrão',
			desc: 'Faxina completa de rotina para manter sua casa em dia.',
			from: 'R$ 140', duration: '4h', color: 'orange', popular: true
		},
		{
			icon: Brush, name: 'Limpeza pesada',
			desc: 'Limpeza profunda em pontos que a rotina não alcança.',
			from: 'R$ 220', duration: '6h', color: 'blue'
		},
		{
			icon: Package, name: 'Pré-mudança',
			desc: 'Deixe seu próximo imóvel impecável antes de entrar.',
			from: 'R$ 240', duration: '6h', color: 'green'
		},
		{
			icon: HardHat, name: 'Pós-obra',
			desc: 'Remoção de resíduos e poeira após reforma ou construção.',
			from: 'R$ 320', duration: '8h', color: 'yellow'
		},
		{
			icon: Building2, name: 'Escritórios',
			desc: 'Limpeza comercial para escritórios e salas.',
			from: 'R$ 160', duration: '4h', color: 'blue'
		},
		{
			icon: Shirt, name: 'Passadoria',
			desc: 'Passar roupas por hora, com excelência e cuidado.',
			from: 'R$ 80', duration: '2h', color: 'orange'
		}
	];

	const steps = [
		{
			n: '01',
			icon: Search,
			title: 'Escolha o serviço',
			desc: 'Selecione o tipo de limpeza, os cômodos, data e horário. Nosso sistema calcula automaticamente o valor justo para o serviço.'
		},
		{
			n: '02',
			icon: UserCheck,
			title: 'Atribuímos a melhor profissional',
			desc: 'Em minutos, nosso algoritmo encontra a diarista mais bem avaliada e próxima da sua região, com disponibilidade confirmada.'
		},
		{
			n: '03',
			icon: CalendarCheck,
			title: 'Relaxe e avalie',
			desc: 'Receba notificações de check-in, conclusão e acompanhe tudo pelo app. Ao final, avalie a profissional e salve como favorita.'
		}
	];

	const diferenciais = [
		{
			icon: BadgeCheck,
			title: 'Diaristas verificadas',
			desc: 'Checagem completa de documentos, referências profissionais e antecedentes. Toda profissional passa por entrevista antes de ser aprovada.',
			color: 'green'
		},
		{
			icon: ShieldCheck,
			title: 'Seguro contra imprevistos',
			desc: 'Todo serviço contratado pelo DiaryGo inclui cobertura de seguro para você e para a profissional durante o atendimento.',
			color: 'blue'
		},
		{
			icon: Zap,
			title: 'Agendamento em minutos',
			desc: 'Encontre uma diarista disponível para hoje ou amanhã em menos de 5 minutos. Sem telefonemas, sem espera.',
			color: 'orange'
		},
		{
			icon: Heart,
			title: 'Favorite sua diarista',
			desc: 'Gostou do atendimento? Marque como favorita e agende diretamente com ela sempre que quiser.',
			color: 'red'
		},
		{
			icon: TrendingUp,
			title: 'Preço transparente',
			desc: 'Sem surpresas. O valor é calculado com base no tipo de serviço e tamanho do imóvel — você vê antes de confirmar.',
			color: 'yellow'
		},
		{
			icon: Headphones,
			title: 'Suporte humano',
			desc: 'Nosso time de atendimento responde em até 2 horas nos dias úteis. Com pessoas de verdade, não bots.',
			color: 'blue'
		}
	];

	const profBenefits = [
		{ icon: DollarSign, title: 'Até R$ 4.500/mês', desc: 'Ganhe mais escolhendo seus horários e regiões.' },
		{ icon: CalendarCheck, title: 'Sua agenda', desc: 'Defina quais dias e horários quer atender.' },
		{ icon: PiggyBank, title: 'Pagamento rápido', desc: 'Receba no dia seguinte ao serviço concluído.' },
		{ icon: ShieldCheck, title: 'Suporte e seguro', desc: 'Você tem respaldo em qualquer situação.' }
	];

	const testimonials = [
		{
			name: 'Ana Carolina',
			city: 'São Paulo · SP',
			initial: 'A',
			rating: 5,
			text: 'Contratei pela primeira vez há 8 meses e nunca mais voltei atrás. A Clara virou minha diarista favorita — marco direto com ela toda semana. Processo ridículo de fácil.',
			gradient: 'linear-gradient(135deg, #ff801f, #ff5900)'
		},
		{
			name: 'Rodrigo Menezes',
			city: 'Rio de Janeiro · RJ',
			initial: 'R',
			rating: 5,
			text: 'Mudança de apartamento foi um parto. Agendei a limpeza pós-mudança no domingo e segunda de manhã a casa tava impecável. Vale cada centavo.',
			gradient: 'linear-gradient(135deg, #3b9eff, #0075ff)'
		},
		{
			name: 'Juliana Freitas',
			city: 'Curitiba · PR',
			initial: 'J',
			rating: 5,
			text: 'O que mais me impressionou foi a transparência do preço. Já contratei outros apps e sempre tinha taxa surpresa. Aqui o valor fechado é o que aparece.',
			gradient: 'linear-gradient(135deg, #22ff99, #11ff99)'
		},
		{
			name: 'Paulo Henrique',
			city: 'Belo Horizonte · MG',
			initial: 'P',
			rating: 5,
			text: 'Precisei remarcar duas vezes em cima da hora e o suporte resolveu em minutos. Profissional que veio era super atenciosa e cuidadosa com meus pets.',
			gradient: 'linear-gradient(135deg, #ffc53d, #ff801f)'
		}
	];

	const faqs = [
		{
			q: 'Como é feita a verificação das diaristas?',
			a: 'Toda profissional passa por conferência de documentos pessoais (RG, CPF, comprovante de residência), checagem de antecedentes, no mínimo 2 referências profissionais validadas por contato direto e entrevista com nossa equipe de qualidade. Só 1 em cada 4 candidatas é aprovada na plataforma.'
		},
		{
			q: 'Como funciona o pagamento?',
			a: 'O pagamento é feito pelo app após o serviço ser concluído. Aceitamos cartão de crédito, débito e Pix. Você vê o valor exato antes de confirmar o agendamento — sem taxas surpresa.'
		},
		{
			q: 'Posso cancelar ou remarcar?',
			a: 'Sim. Cancelamentos e remarcações com mais de 24h de antecedência são gratuitos. Dentro de 24h, há uma taxa que cobre o deslocamento da profissional. Em casos de emergência, nosso suporte avalia isenção.'
		},
		{
			q: 'E se eu não gostar do serviço?',
			a: 'Após o atendimento você avalia a profissional. Se houver problema real, nossa equipe entra em contato em até 24h para entender, mediar e, quando aplicável, oferecer um novo serviço sem custo ou reembolso.'
		},
		{
			q: 'Como agendar uma diarista específica?',
			a: 'Após um serviço, você pode favoritar a profissional. No próximo agendamento, é só escolhê-la entre suas favoritas — sujeito à disponibilidade dela para o horário desejado.'
		},
		{
			q: 'O serviço inclui produtos de limpeza?',
			a: 'Por padrão, o cliente fornece os produtos e equipamentos. Se preferir, pode contratar nossa opção com materiais inclusos — ideal para pós-mudança ou quando você prefere praticidade.'
		}
	];

	const pressLogos = ['EXAME', 'ESTADÃO', 'FOLHA', 'G1', 'O GLOBO', 'UOL', 'VALOR', 'CNN'];
</script>

<svelte:head>
	<title>DiaryGo — Diaristas verificadas, agendamento em minutos</title>
	<meta name="description" content="Encontre diaristas verificadas na sua região. Preço transparente, agendamento em minutos e profissionais avaliadas pela comunidade DiaryGo." />
</svelte:head>

<!-- ═════════════════════════════════════════════════════
     1. HERO
═════════════════════════════════════════════════════ -->
<section class="hero">
	<div class="hero-bg">
		<div class="hero-glow glow-orange"></div>
		<div class="hero-glow glow-blue"></div>
		<div class="hero-grid-lines"></div>
	</div>

	<div class="container hero-content">
		<div class="hero-left">
			<div class="hero-eyebrow">
				<span class="eyebrow-dot"></span>
				Mais de 180 mil serviços realizados
			</div>

			<h1 class="hero-title">
				Diaristas de confiança,<br />
				<span class="title-accent">agendadas em minutos</span>.
			</h1>

			<p class="hero-sub">
				A forma mais simples e segura de encontrar diaristas verificadas na sua região.
				Preço transparente, profissionais avaliadas e suporte humano quando precisar.
			</p>

			<div class="hero-actions">
				{#if $isAuthenticated}
					<a href="/dashboard" class="btn-solid">
						Ir para o dashboard
						<ArrowRight size={15} strokeWidth={2.2} />
					</a>
				{:else}
					<a href="/registro" class="btn-solid">
						Contratar agora
						<ArrowRight size={15} strokeWidth={2.2} />
					</a>
					<a href="#como-funciona" class="btn-outline">
						<Play size={13} strokeWidth={2.2} />
						Como funciona
					</a>
				{/if}
			</div>

			<div class="hero-social">
				<div class="avatar-stack">
					<div class="avatar avatar-1">A</div>
					<div class="avatar avatar-2">M</div>
					<div class="avatar avatar-3">J</div>
					<div class="avatar avatar-4">R</div>
				</div>
				<div class="social-text">
					<div class="stars-row">
						{#each Array(5) as _}
							<Star size={12} strokeWidth={2} fill="currentColor" />
						{/each}
					</div>
					<span class="social-label">
						<strong>4.9/5</strong> · mais de 15 mil avaliações
					</span>
				</div>
			</div>
		</div>

		<!-- Hero visual -->
		<div class="hero-right">
			<div class="mockup-frame">
				<div class="mockup-browser">
					<div class="browser-chrome">
						<span class="chrome-dot chrome-dot-red"></span>
						<span class="chrome-dot chrome-dot-yellow"></span>
						<span class="chrome-dot chrome-dot-green"></span>
						<div class="chrome-url">
							<Sparkles size={10} strokeWidth={2} />
							diarygo.com/agendar
						</div>
					</div>
					<div class="browser-body">
						<div class="mock-title">Quando você precisa?</div>
						<div class="mock-pills">
							<span class="mock-pill mock-pill-active">Hoje</span>
							<span class="mock-pill">Amanhã</span>
							<span class="mock-pill">Esta semana</span>
						</div>

						<div class="mock-section-label">Diaristas próximas</div>
						<div class="mock-card">
							<div class="mock-avatar" style="background: linear-gradient(135deg, #ff801f, #ff5900)">C</div>
							<div class="mock-info">
								<div class="mock-name">Clara Oliveira</div>
								<div class="mock-meta">
									<Star size={10} strokeWidth={2} fill="currentColor" />
									4.98 · a 1.2km · 247 serviços
								</div>
							</div>
							<span class="mock-badge mock-badge-green">
								<span class="pulse-dot"></span>
								Hoje
							</span>
						</div>
						<div class="mock-card">
							<div class="mock-avatar" style="background: linear-gradient(135deg, #3b9eff, #0075ff)">M</div>
							<div class="mock-info">
								<div class="mock-name">Maria Santos</div>
								<div class="mock-meta">
									<Star size={10} strokeWidth={2} fill="currentColor" />
									4.92 · a 2.1km · 182 serviços
								</div>
							</div>
							<span class="mock-badge mock-badge-green">
								<span class="pulse-dot"></span>
								Hoje
							</span>
						</div>
						<div class="mock-card">
							<div class="mock-avatar" style="background: linear-gradient(135deg, #22ff99, #11cc77)">L</div>
							<div class="mock-info">
								<div class="mock-name">Letícia Alves</div>
								<div class="mock-meta">
									<Star size={10} strokeWidth={2} fill="currentColor" />
									4.87 · a 3.4km · 96 serviços
								</div>
							</div>
							<span class="mock-badge mock-badge-muted">Amanhã</span>
						</div>
						<button class="mock-cta">Continuar agendamento</button>
					</div>
				</div>

				<!-- Floating cards -->
				<div class="float-card float-rating">
					<Star size={16} strokeWidth={2} fill="currentColor" />
					<div class="float-body">
						<div class="float-value">4.98</div>
						<div class="float-label">Avaliação</div>
					</div>
				</div>

				<div class="float-card float-verified">
					<div class="float-icon-wrap">
						<BadgeCheck size={16} strokeWidth={2} />
					</div>
					<div class="float-body">
						<div class="float-value">Verificada</div>
						<div class="float-label">Documentos e referências</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Stats bar -->
	<div class="container">
		<div class="stats-bar">
			{#each heroStats as s}
				{@const Icon = s.icon}
				<div class="stat-item">
					<div class="stat-icon stat-{s.color}">
						<Icon size={14} strokeWidth={2} />
					</div>
					<div class="stat-text">
						<span class="stat-val">{s.value}</span>
						<span class="stat-lbl">{s.label}</span>
					</div>
				</div>
			{/each}
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     2. LOGO BAR
═════════════════════════════════════════════════════ -->
<section class="press-band">
	<div class="container">
		<p class="press-label">Como visto em</p>
		<div class="press-logos">
			{#each pressLogos as logo}
				<span class="press-logo">{logo}</span>
			{/each}
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     3. SERVIÇOS
═════════════════════════════════════════════════════ -->
<section class="section" id="servicos">
	<div class="container">
		<div class="section-header">
			<div class="section-eyebrow">
				<Sparkles size={11} strokeWidth={2.2} />
				Serviços
			</div>
			<h2 class="section-title">Tudo que sua casa precisa, em um só lugar</h2>
			<p class="section-lead">
				Do faxina de rotina à limpeza especializada. Agende em minutos, acompanhe pelo app.
			</p>
		</div>

		<div class="services-grid">
			{#each services as svc}
				{@const Icon = svc.icon}
				<article class="service-card">
					{#if svc.popular}
						<span class="service-badge">Mais pedido</span>
					{/if}
					<div class="service-icon service-icon-{svc.color}">
						<Icon size={20} strokeWidth={1.6} />
					</div>
					<div class="service-body">
						<h3 class="service-name">{svc.name}</h3>
						<p class="service-desc">{svc.desc}</p>
					</div>
					<div class="service-foot">
						<div class="service-pricing">
							<span class="service-price">{svc.from}</span>
							<span class="service-duration"><Clock size={10} strokeWidth={2}/> {svc.duration}</span>
						</div>
						<a href="/registro" class="service-link" aria-label="Agendar {svc.name}">
							<ArrowUpRight size={14} strokeWidth={2.2} />
						</a>
					</div>
				</article>
			{/each}
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     4. COMO FUNCIONA
═════════════════════════════════════════════════════ -->
<section class="section section-alt" id="como-funciona">
	<div class="container">
		<div class="section-header section-header-center">
			<div class="section-eyebrow">
				<CalendarCheck size={11} strokeWidth={2.2} />
				Como funciona
			</div>
			<h2 class="section-title">Da solicitação à avaliação em 3 passos</h2>
			<p class="section-lead">
				Processo transparente, sem telefonemas, sem surpresas. Em média, nosso tempo de atribuição é de <strong>4 minutos</strong>.
			</p>
		</div>

		<div class="steps-grid">
			{#each steps as step, i}
				{@const Icon = step.icon}
				<div class="step-card">
					<div class="step-head">
						<span class="step-num">{step.n}</span>
						<div class="step-icon">
							<Icon size={18} strokeWidth={1.7} />
						</div>
					</div>
					<h3 class="step-title">{step.title}</h3>
					<p class="step-desc">{step.desc}</p>
					{#if i < steps.length - 1}
						<div class="step-connector">
							<ArrowRight size={14} strokeWidth={2} />
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     5. DIFERENCIAIS
═════════════════════════════════════════════════════ -->
<section class="section" id="diferenciais">
	<div class="container">
		<div class="section-header">
			<div class="section-eyebrow">
				<ShieldCheck size={11} strokeWidth={2.2} />
				Por que DiaryGo
			</div>
			<h2 class="section-title">Confiança que dá pra sentir</h2>
			<p class="section-lead">
				Construímos cada detalhe pensando em quem já passou por experiências ruins contratando diaristas.
			</p>
		</div>

		<div class="dif-grid">
			{#each diferenciais as dif}
				{@const Icon = dif.icon}
				<article class="dif-card">
					<div class="dif-icon dif-icon-{dif.color}">
						<Icon size={18} strokeWidth={1.7} />
					</div>
					<h3 class="dif-title">{dif.title}</h3>
					<p class="dif-desc">{dif.desc}</p>
				</article>
			{/each}
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     6. DEPOIMENTOS
═════════════════════════════════════════════════════ -->
<section class="section section-alt">
	<div class="container">
		<div class="section-header">
			<div class="section-eyebrow">
				<Heart size={11} strokeWidth={2.2} />
				Depoimentos
			</div>
			<h2 class="section-title">Histórias de quem já se apaixonou</h2>
			<p class="section-lead">
				Milhares de clientes contam como o DiaryGo mudou a rotina deles.
			</p>
		</div>

		<div class="test-grid">
			{#each testimonials as t}
				<article class="test-card">
					<div class="test-stars">
						{#each Array(t.rating) as _}
							<Star size={13} strokeWidth={2} fill="currentColor" />
						{/each}
					</div>
					<p class="test-text">"{t.text}"</p>
					<div class="test-foot">
						<div class="test-avatar" style="background: {t.gradient}">{t.initial}</div>
						<div class="test-info">
							<span class="test-name">{t.name}</span>
							<span class="test-city">{t.city}</span>
						</div>
					</div>
				</article>
			{/each}
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     7. PARA DIARISTAS
═════════════════════════════════════════════════════ -->
<section class="section prof-section" id="para-diaristas">
	<div class="prof-bg"></div>
	<div class="container prof-grid">
		<div class="prof-left">
			<div class="section-eyebrow">
				<HeartHandshake size={11} strokeWidth={2.2} />
				Para diaristas
			</div>
			<h2 class="section-title">Trabalhe quando quiser.<br/>Ganhe o que merece.</h2>
			<p class="section-lead">
				Junte-se a mais de 500 profissionais que já construíram sua carreira no DiaryGo. Você define sua agenda, suas regiões e suas condições.
			</p>

			<div class="prof-benefits">
				{#each profBenefits as b}
					{@const Icon = b.icon}
					<div class="prof-benefit">
						<div class="prof-ben-icon">
							<Icon size={15} strokeWidth={1.8} />
						</div>
						<div>
							<div class="prof-ben-title">{b.title}</div>
							<div class="prof-ben-desc">{b.desc}</div>
						</div>
					</div>
				{/each}
			</div>

			<a href="/registro/profissional" class="btn-solid prof-cta">
				Quero me cadastrar
				<ArrowRight size={15} strokeWidth={2.2} />
			</a>
		</div>

		<div class="prof-right">
			<div class="earnings-card">
				<div class="earnings-head">
					<span class="earnings-label">Ganhos no mês</span>
					<span class="earnings-delta">
						<TrendingUp size={11} strokeWidth={2.2} />
						+18%
					</span>
				</div>
				<div class="earnings-value">R$ 3.847<span class="earnings-cents">,50</span></div>

				<div class="earnings-bars">
					{#each [40, 55, 48, 72, 85, 68, 92] as h, i}
						<div class="bar-wrap">
							<div class="bar" style="height: {h}%"></div>
							<span class="bar-lbl">{['S','T','Q','Q','S','S','D'][i]}</span>
						</div>
					{/each}
				</div>

				<div class="earnings-meta">
					<div class="earnings-meta-item">
						<div class="meta-dot meta-dot-orange"></div>
						<div>
							<div class="meta-val">32</div>
							<div class="meta-lbl">Serviços</div>
						</div>
					</div>
					<div class="earnings-meta-item">
						<div class="meta-dot meta-dot-blue"></div>
						<div>
							<div class="meta-val">4.9</div>
							<div class="meta-lbl">Nota média</div>
						</div>
					</div>
					<div class="earnings-meta-item">
						<div class="meta-dot meta-dot-green"></div>
						<div>
							<div class="meta-val">96%</div>
							<div class="meta-lbl">Taxa aceite</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     8. FAQ
═════════════════════════════════════════════════════ -->
<section class="section" id="faq">
	<div class="container faq-container">
		<div class="faq-left">
			<div class="section-eyebrow">
				<Search size={11} strokeWidth={2.2} />
				Perguntas frequentes
			</div>
			<h2 class="section-title">Tudo que você quer saber</h2>
			<p class="section-lead">
				Respostas rápidas para as dúvidas mais comuns. Não achou o que procurava? Nosso time responde em até 2 horas.
			</p>
			<a href="/ajuda" class="btn-outline">
				Central de ajuda
				<ArrowUpRight size={14} strokeWidth={2.2} />
			</a>
		</div>

		<div class="faq-list">
			{#each faqs as faq, i}
				<div class="faq-item" class:open={openFaq === i}>
					<button class="faq-q" onclick={() => toggleFaq(i)} aria-expanded={openFaq === i}>
						<span>{faq.q}</span>
						<span class="faq-toggle">
							{#if openFaq === i}
								<Minus size={14} strokeWidth={2.2} />
							{:else}
								<Plus size={14} strokeWidth={2.2} />
							{/if}
						</span>
					</button>
					{#if openFaq === i}
						<div class="faq-a">{faq.a}</div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>


<!-- ═════════════════════════════════════════════════════
     9. CTA FINAL
═════════════════════════════════════════════════════ -->
<section class="cta-section">
	<div class="container">
		<div class="cta-box">
			<div class="cta-glow"></div>
			<div class="cta-content">
				<div class="section-eyebrow">
					<Sparkles size={11} strokeWidth={2.2} />
					Comece agora
				</div>
				<h2 class="cta-title">Sua casa impecável<br/>está a um clique</h2>
				<p class="cta-lead">
					Crie sua conta gratuitamente e agende seu primeiro serviço hoje. Sem mensalidade, sem taxa de cadastro.
				</p>
				<div class="cta-actions">
					<a href="/registro" class="btn-solid">
						Sou cliente
						<ArrowRight size={15} strokeWidth={2.2} />
					</a>
					<a href="/registro/profissional" class="btn-outline">
						Sou diarista
					</a>
				</div>
			</div>
		</div>
	</div>
</section>


<style>
	/* ═════════════════ shared tokens ═════════════════ */
	.container {
		max-width: 1240px;
		margin: 0 auto;
		padding: 0 24px;
	}

	.section {
		padding: 96px 0;
		position: relative;
	}

	.section-alt {
		background:
			radial-gradient(60% 50% at 50% 0%, rgba(255,255,255,0.015), transparent 70%),
			#000;
		border-top: 1px solid var(--border-frost);
		border-bottom: 1px solid var(--border-frost);
	}

	.section-header {
		max-width: 640px;
		margin-bottom: 56px;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.section-header-center {
		margin-left: auto;
		margin-right: auto;
		text-align: center;
		align-items: center;
	}

	.section-eyebrow {
		display: inline-flex;
		align-items: center;
		gap: 5px;
		padding: 4px 10px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		font-size: 0.6875rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		width: fit-content;
	}

	.section-eyebrow :global(svg) {
		color: var(--color-orange-10);
	}

	.section-title {
		font-size: clamp(1.75rem, 3.5vw, 2.625rem);
		font-weight: 400;
		letter-spacing: -1.5px;
		line-height: 1.1;
		color: var(--color-text-primary);
	}

	.section-lead {
		font-size: 0.9375rem;
		color: var(--color-text-secondary);
		line-height: 1.6;
		max-width: 560px;
	}

	.section-header-center .section-lead {
		margin: 0 auto;
	}

	/* Buttons — pílulas do design.md */
	.btn-solid {
		display: inline-flex;
		align-items: center;
		gap: 7px;
		padding: 11px 20px;
		background: var(--color-white);
		color: var(--color-black);
		border: none;
		border-radius: var(--radius-pill);
		font-size: 0.875rem;
		font-weight: 600;
		letter-spacing: -0.1px;
		text-decoration: none;
		cursor: pointer;
		transition: opacity 0.15s, transform 0.15s;
	}

	.btn-solid:hover {
		opacity: 0.92;
		transform: translateY(-1px);
	}

	.btn-solid:focus-visible {
		box-shadow: var(--shadow-focus);
		outline: none;
	}

	.btn-outline {
		display: inline-flex;
		align-items: center;
		gap: 7px;
		padding: 11px 18px;
		background: transparent;
		color: var(--color-text-primary);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		font-size: 0.875rem;
		font-weight: 500;
		text-decoration: none;
		transition: background 0.15s, border-color 0.15s;
		width: fit-content;
	}

	.btn-outline:hover {
		background: var(--bg-hover-white);
		border-color: var(--border-frost-alt);
		opacity: 1;
	}

	.btn-outline:focus-visible {
		box-shadow: var(--shadow-focus);
		outline: none;
	}

	/* ═════════════════ HERO ═════════════════ */
	.hero {
		position: relative;
		padding: 72px 0 96px;
		overflow: hidden;
	}

	.hero-bg {
		position: absolute;
		inset: 0;
		pointer-events: none;
		z-index: 0;
	}

	.hero-glow {
		position: absolute;
		border-radius: 50%;
		filter: blur(80px);
	}

	.glow-orange {
		width: 600px;
		height: 600px;
		top: -200px;
		left: -200px;
		background: radial-gradient(circle, rgba(255, 128, 31, 0.18), transparent 70%);
	}

	.glow-blue {
		width: 500px;
		height: 500px;
		bottom: -150px;
		right: -150px;
		background: radial-gradient(circle, rgba(59, 158, 255, 0.1), transparent 70%);
	}

	.hero-grid-lines {
		position: absolute;
		inset: 0;
		background-image:
			linear-gradient(rgba(214, 235, 253, 0.04) 1px, transparent 1px),
			linear-gradient(90deg, rgba(214, 235, 253, 0.04) 1px, transparent 1px);
		background-size: 44px 44px;
		mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
		-webkit-mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
	}

	.hero-content {
		position: relative;
		z-index: 1;
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 60px;
		align-items: center;
		margin-bottom: 72px;
	}

	.hero-left {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.hero-eyebrow {
		display: inline-flex;
		align-items: center;
		gap: 7px;
		padding: 5px 12px;
		background: var(--color-orange-4);
		border: 1px solid var(--color-orange-4);
		border-radius: var(--radius-pill);
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-orange-11);
		width: fit-content;
	}

	.eyebrow-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-orange-10);
		box-shadow: 0 0 8px var(--color-orange-10);
		animation: pulse 2s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.4; }
	}

	.hero-title {
		font-size: clamp(2.25rem, 5vw, 3.75rem);
		font-weight: 400;
		letter-spacing: -2px;
		line-height: 1.05;
		color: var(--color-text-primary);
	}

	.title-accent {
		background: linear-gradient(135deg, var(--color-orange-11) 0%, var(--color-orange-10) 50%, #ff5900 100%);
		-webkit-background-clip: text;
		background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.hero-sub {
		font-size: 1rem;
		color: var(--color-text-secondary);
		max-width: 480px;
		line-height: 1.6;
	}

	.hero-actions {
		display: flex;
		gap: 10px;
		margin-top: 8px;
		flex-wrap: wrap;
	}

	.hero-social {
		display: flex;
		align-items: center;
		gap: 14px;
		margin-top: 16px;
	}

	.avatar-stack {
		display: flex;
	}

	.avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.75rem;
		font-weight: 600;
		color: #fff;
		border: 2px solid #000;
		margin-left: -8px;
	}

	.avatar:first-child { margin-left: 0; }
	.avatar-1 { background: linear-gradient(135deg, #ff801f, #ff5900); }
	.avatar-2 { background: linear-gradient(135deg, #3b9eff, #0075ff); }
	.avatar-3 { background: linear-gradient(135deg, #22ff99, #11cc77); color: #003322; }
	.avatar-4 { background: linear-gradient(135deg, #ffc53d, #ff801f); color: #3a1e00; }

	.social-text {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.stars-row {
		display: flex;
		gap: 1px;
		color: var(--color-yellow-9);
	}

	.social-label {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
	}

	.social-label strong {
		color: var(--color-text-primary);
		font-weight: 600;
	}

	/* Hero visual / mockup */
	.hero-right {
		position: relative;
		display: flex;
		justify-content: center;
	}

	.mockup-frame {
		position: relative;
		width: 100%;
		max-width: 440px;
	}

	.mockup-browser {
		background: #0a0a0a;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-large);
		overflow: hidden;
		box-shadow: var(--shadow-mockup), var(--shadow-ring);
		transform: perspective(1400px) rotateY(-6deg) rotateX(3deg);
		transition: transform 0.5s ease;
	}

	.mockup-browser:hover {
		transform: perspective(1400px) rotateY(-3deg) rotateX(1deg) scale(1.02);
	}

	.browser-chrome {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 11px 14px;
		border-bottom: 1px solid var(--border-frost);
	}

	.chrome-dot {
		width: 9px;
		height: 9px;
		border-radius: 50%;
	}

	.chrome-dot-red    { background: #ff5f57; }
	.chrome-dot-yellow { background: #febc2e; }
	.chrome-dot-green  { background: #28c840; }

	.chrome-url {
		margin-left: 12px;
		flex: 1;
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 4px 10px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-subtle);
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		font-family: var(--font-mono);
	}

	.chrome-url :global(svg) {
		color: var(--color-orange-10);
	}

	.browser-body {
		padding: 20px 22px 22px;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.mock-title {
		font-size: 0.9375rem;
		font-weight: 500;
		color: var(--color-text-primary);
		margin-bottom: 2px;
		letter-spacing: -0.2px;
	}

	.mock-pills {
		display: flex;
		gap: 6px;
		margin-bottom: 8px;
	}

	.mock-pill {
		padding: 5px 11px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-pill);
		font-size: 0.6875rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
	}

	.mock-pill-active {
		background: var(--color-orange-4);
		border-color: var(--color-orange-10);
		color: var(--color-orange-11);
	}

	.mock-section-label {
		font-size: 0.6875rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.4px;
	}

	.mock-card {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 10px 12px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-standard);
	}

	.mock-avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.8125rem;
		font-weight: 600;
		color: #fff;
		flex-shrink: 0;
	}

	.mock-info {
		flex: 1;
		min-width: 0;
	}

	.mock-name {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.mock-meta {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		display: inline-flex;
		align-items: center;
		gap: 3px;
	}

	.mock-meta :global(svg) {
		color: var(--color-yellow-9);
	}

	.mock-badge {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		font-size: 0.625rem;
		font-weight: 600;
		padding: 3px 8px;
		border-radius: var(--radius-pill);
	}

	.mock-badge-green {
		background: var(--color-green-3);
		color: var(--color-green-10);
	}

	.mock-badge-muted {
		background: var(--bg-hover-subtle);
		color: var(--color-text-tertiary);
	}

	.pulse-dot {
		width: 5px;
		height: 5px;
		border-radius: 50%;
		background: currentColor;
		animation: pulse 2s ease-in-out infinite;
	}

	.mock-cta {
		margin-top: 6px;
		padding: 11px;
		background: var(--color-white);
		color: var(--color-black);
		border: none;
		border-radius: var(--radius-pill);
		font-size: 0.8125rem;
		font-weight: 600;
		cursor: pointer;
		font-family: var(--font-body);
	}

	/* Floating cards */
	.float-card {
		position: absolute;
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 10px 14px;
		background: #0a0a0a;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-standard);
		box-shadow: var(--shadow-float);
		animation: floatBob 4s ease-in-out infinite;
	}

	@keyframes floatBob {
		0%, 100% { transform: translateY(0); }
		50% { transform: translateY(-8px); }
	}

	.float-rating {
		top: 10%;
		left: -32px;
		color: var(--color-yellow-9);
	}

	.float-verified {
		bottom: 12%;
		right: -28px;
		animation-delay: 1s;
	}

	.float-icon-wrap {
		width: 26px;
		height: 26px;
		border-radius: var(--radius-subtle);
		background: var(--color-green-3);
		color: var(--color-green-10);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.float-body {
		display: flex;
		flex-direction: column;
		gap: 0;
	}

	.float-value {
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--color-text-primary);
		line-height: 1.2;
	}

	.float-label {
		font-size: 0.625rem;
		color: var(--color-text-tertiary);
	}

	/* Stats bar */
	.stats-bar {
		position: relative;
		z-index: 1;
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 2px;
		background: var(--border-frost);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		overflow: hidden;
	}

	.stat-item {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 18px 20px;
		background: rgba(0, 0, 0, 0.9);
	}

	.stat-icon {
		width: 32px;
		height: 32px;
		border-radius: var(--radius-standard);
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.stat-green  { background: var(--color-green-3); color: var(--color-green-10); }
	.stat-blue   { background: var(--color-blue-4);  color: var(--color-blue-10); }
	.stat-orange { background: var(--color-orange-4); color: var(--color-orange-11); }
	.stat-yellow { background: rgba(255,197,61,0.15); color: var(--color-yellow-9); }

	.stat-text {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}

	.stat-val {
		font-size: 1.25rem;
		font-weight: 500;
		letter-spacing: -0.8px;
		color: var(--color-text-primary);
		line-height: 1;
	}

	.stat-lbl {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	/* ═════════════════ PRESS BAND ═════════════════ */
	.press-band {
		padding: 40px 0;
		border-top: 1px solid var(--border-frost);
		border-bottom: 1px solid var(--border-frost);
	}

	.press-label {
		font-size: 0.6875rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.8px;
		text-align: center;
		margin-bottom: 20px;
	}

	.press-logos {
		display: flex;
		align-items: center;
		justify-content: space-around;
		flex-wrap: wrap;
		gap: 32px;
	}

	.press-logo {
		font-family: var(--font-body);
		font-size: 0.9375rem;
		font-weight: 700;
		letter-spacing: 1px;
		color: var(--color-text-tertiary);
		opacity: 0.6;
		transition: opacity 0.2s, color 0.2s;
	}

	.press-logo:hover {
		opacity: 1;
		color: var(--color-text-primary);
	}

	/* ═════════════════ SERVICES ═════════════════ */
	.services-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 14px;
	}

	.service-card {
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 22px;
		background:
			linear-gradient(180deg, var(--bg-card) 0%, rgba(255,255,255,0.01) 100%);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		transition: transform 0.2s, border-color 0.2s;
		min-height: 200px;
	}

	.service-card:hover {
		transform: translateY(-2px);
		border-color: var(--border-frost-alt);
	}

	.service-badge {
		position: absolute;
		top: 14px;
		right: 14px;
		font-size: 0.625rem;
		font-weight: 600;
		padding: 3px 8px;
		background: rgba(255, 128, 31, 0.12);
		color: var(--color-orange-11);
		border: 1px solid var(--color-orange-4);
		border-radius: 99px;
		text-transform: uppercase;
		letter-spacing: 0.3px;
	}

	.service-icon {
		width: 40px;
		height: 40px;
		border-radius: var(--radius-standard);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.service-icon-orange { background: var(--color-orange-4); color: var(--color-orange-11); }
	.service-icon-blue   { background: var(--color-blue-4);   color: var(--color-blue-10); }
	.service-icon-green  { background: var(--color-green-3);  color: var(--color-green-10); }
	.service-icon-yellow { background: rgba(255,197,61,0.15); color: var(--color-yellow-9); }

	.service-body {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.service-name {
		font-size: 1rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.3px;
	}

	.service-desc {
		font-size: 0.8125rem;
		color: var(--color-text-tertiary);
		line-height: 1.55;
	}

	.service-foot {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding-top: 12px;
		border-top: 1px solid var(--border-frost);
	}

	.service-pricing {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.service-price {
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--color-text-primary);
		font-feature-settings: 'tnum';
	}

	.service-duration {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
		display: inline-flex;
		align-items: center;
		gap: 3px;
	}

	.service-link {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-standard);
		color: var(--color-text-secondary);
		transition: background 0.15s, color 0.15s, border-color 0.15s;
	}

	.service-card:hover .service-link {
		background: var(--color-white);
		color: var(--color-black);
		border-color: var(--color-white);
	}

	/* ═════════════════ STEPS ═════════════════ */
	.steps-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 20px;
		position: relative;
	}

	.step-card {
		position: relative;
		padding: 28px 24px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-large);
		display: flex;
		flex-direction: column;
		gap: 12px;
		box-shadow: var(--shadow-ring);
	}

	.step-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 8px;
	}

	.step-num {
		font-size: 0.6875rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
		letter-spacing: 1px;
		font-feature-settings: 'tnum';
	}

	.step-icon {
		width: 38px;
		height: 38px;
		border-radius: var(--radius-standard);
		background: linear-gradient(135deg, var(--color-orange-4), transparent);
		border: 1px solid var(--color-orange-4);
		color: var(--color-orange-11);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.step-title {
		font-size: 1.125rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.4px;
	}

	.step-desc {
		font-size: 0.875rem;
		color: var(--color-text-secondary);
		line-height: 1.6;
	}

	.step-connector {
		position: absolute;
		right: -18px;
		top: 50%;
		transform: translateY(-50%);
		width: 28px;
		height: 28px;
		background: #000;
		border: 1px solid var(--border-frost);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-tertiary);
		z-index: 2;
	}

	/* ═════════════════ DIFERENCIAIS ═════════════════ */
	.dif-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 14px;
	}

	.dif-card {
		padding: 24px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		display: flex;
		flex-direction: column;
		gap: 10px;
		transition: transform 0.2s, border-color 0.2s;
	}

	.dif-card:hover {
		transform: translateY(-2px);
		border-color: var(--border-frost-hover);
	}

	.dif-icon {
		width: 38px;
		height: 38px;
		border-radius: var(--radius-standard);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 6px;
	}

	.dif-icon-green  { background: var(--color-green-3);  color: var(--color-green-10); }
	.dif-icon-blue   { background: var(--color-blue-4);   color: var(--color-blue-10); }
	.dif-icon-orange { background: var(--color-orange-4); color: var(--color-orange-11); }
	.dif-icon-red    { background: var(--color-red-5);    color: var(--color-red-10); }
	.dif-icon-yellow { background: rgba(255,197,61,0.15); color: var(--color-yellow-9); }

	.dif-title {
		font-size: 1rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.3px;
	}

	.dif-desc {
		font-size: 0.8125rem;
		color: var(--color-text-tertiary);
		line-height: 1.6;
	}

	/* ═════════════════ TESTIMONIALS ═════════════════ */
	.test-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 14px;
	}

	.test-card {
		padding: 22px;
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-card);
		display: flex;
		flex-direction: column;
		gap: 14px;
	}

	.test-stars {
		display: flex;
		gap: 2px;
		color: var(--color-yellow-9);
	}

	.test-text {
		font-size: 0.875rem;
		color: var(--color-text-primary);
		line-height: 1.6;
		flex: 1;
		letter-spacing: -0.1px;
	}

	.test-foot {
		display: flex;
		align-items: center;
		gap: 10px;
		padding-top: 12px;
		border-top: 1px solid var(--border-frost);
	}

	.test-avatar {
		width: 36px;
		height: 36px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.875rem;
		font-weight: 600;
		color: #fff;
		flex-shrink: 0;
	}

	.test-info {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}

	.test-name {
		font-size: 0.8125rem;
		font-weight: 500;
		color: var(--color-text-primary);
	}

	.test-city {
		font-size: 0.6875rem;
		color: var(--color-text-tertiary);
	}

	/* ═════════════════ PARA DIARISTAS ═════════════════ */
	.prof-section {
		position: relative;
		overflow: hidden;
	}

	.prof-bg {
		position: absolute;
		inset: 0;
		background:
			radial-gradient(60% 80% at 100% 50%, rgba(34, 255, 153, 0.06), transparent 60%),
			radial-gradient(50% 70% at 0% 50%, rgba(59, 158, 255, 0.05), transparent 60%);
		pointer-events: none;
	}

	.prof-grid {
		position: relative;
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 64px;
		align-items: center;
	}

	.prof-left {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.prof-benefits {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 16px;
		margin-top: 16px;
	}

	.prof-benefit {
		display: flex;
		align-items: flex-start;
		gap: 12px;
	}

	.prof-ben-icon {
		width: 32px;
		height: 32px;
		border-radius: var(--radius-standard);
		background: var(--color-green-3);
		color: var(--color-green-10);
		border: 1px solid var(--color-green-4);
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.prof-ben-title {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--color-text-primary);
		letter-spacing: -0.2px;
	}

	.prof-ben-desc {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
		line-height: 1.5;
	}

	.prof-cta {
		margin-top: 12px;
	}

	/* Earnings card */
	.prof-right {
		display: flex;
		justify-content: center;
	}

	.earnings-card {
		width: 100%;
		max-width: 400px;
		padding: 26px;
		background: #0a0a0a;
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-large);
		box-shadow: var(--shadow-lift), var(--shadow-ring);
	}

	.earnings-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 10px;
	}

	.earnings-label {
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--color-text-tertiary);
	}

	.earnings-delta {
		display: inline-flex;
		align-items: center;
		gap: 3px;
		padding: 3px 8px;
		background: var(--color-green-3);
		color: var(--color-green-10);
		border-radius: var(--radius-pill);
		font-size: 0.6875rem;
		font-weight: 600;
	}

	.earnings-value {
		font-size: 2.25rem;
		font-weight: 400;
		letter-spacing: -2px;
		color: var(--color-text-primary);
		line-height: 1;
		margin-bottom: 24px;
		font-feature-settings: 'tnum';
	}

	.earnings-cents {
		font-size: 1.25rem;
		color: var(--color-text-tertiary);
	}

	.earnings-bars {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: 6px;
		height: 80px;
		margin-bottom: 20px;
		padding-bottom: 16px;
		border-bottom: 1px solid var(--border-frost);
	}

	.bar-wrap {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 6px;
		height: 100%;
		justify-content: flex-end;
	}

	.bar {
		width: 100%;
		background: linear-gradient(180deg, var(--color-orange-10), var(--color-orange-4));
		border-radius: var(--radius-sharp) var(--radius-sharp) 2px 2px;
		min-height: 4px;
	}

	.bar-lbl {
		font-size: 0.5625rem;
		font-weight: 600;
		color: var(--color-text-tertiary);
	}

	.earnings-meta {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 12px;
	}

	.earnings-meta-item {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.meta-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		flex-shrink: 0;
	}

	.meta-dot-orange { background: var(--color-orange-10); }
	.meta-dot-blue   { background: var(--color-blue-10); }
	.meta-dot-green  { background: var(--color-green-10); }

	.meta-val {
		font-size: 0.9375rem;
		font-weight: 500;
		color: var(--color-text-primary);
		line-height: 1;
		font-feature-settings: 'tnum';
	}

	.meta-lbl {
		font-size: 0.625rem;
		color: var(--color-text-tertiary);
	}

	/* ═════════════════ FAQ ═════════════════ */
	.faq-container {
		display: grid;
		grid-template-columns: 0.8fr 1.2fr;
		gap: 64px;
		align-items: start;
	}

	.faq-left {
		display: flex;
		flex-direction: column;
		gap: 14px;
		position: sticky;
		top: 100px;
	}

	.faq-list {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.faq-item {
		background: var(--bg-card);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-standard);
		overflow: hidden;
		transition: border-color 0.2s;
	}

	.faq-item.open {
		border-color: var(--border-frost-hover);
	}

	.faq-q {
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		padding: 18px 20px;
		background: transparent;
		border: none;
		font-family: var(--font-body);
		font-size: 0.9375rem;
		font-weight: 500;
		color: var(--color-text-primary);
		text-align: left;
		cursor: pointer;
		letter-spacing: -0.2px;
		gap: 14px;
	}

	.faq-q:hover {
		color: var(--color-text-primary);
	}

	.faq-toggle {
		flex-shrink: 0;
		width: 26px;
		height: 26px;
		border-radius: var(--radius-subtle);
		background: var(--bg-hover-subtle);
		border: 1px solid var(--border-frost);
		display: inline-flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-secondary);
		transition: all 0.15s;
	}

	.faq-item.open .faq-toggle {
		background: var(--color-orange-4);
		border-color: var(--color-orange-10);
		color: var(--color-orange-11);
	}

	.faq-a {
		padding: 0 20px 18px;
		font-size: 0.875rem;
		color: var(--color-text-secondary);
		line-height: 1.65;
	}

	/* ═════════════════ CTA FINAL ═════════════════ */
	.cta-section {
		padding: 96px 0 120px;
	}

	.cta-box {
		position: relative;
		padding: 72px 48px;
		background:
			linear-gradient(180deg, var(--color-orange-4), transparent);
		border: 1px solid var(--border-frost);
		border-radius: var(--radius-section);
		overflow: hidden;
		text-align: center;
	}

	.cta-glow {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: 600px;
		height: 300px;
		background: radial-gradient(ellipse, rgba(255, 128, 31, 0.1), transparent 60%);
		pointer-events: none;
	}

	.cta-content {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 14px;
	}

	.cta-title {
		font-size: clamp(2rem, 4.5vw, 3.25rem);
		font-weight: 400;
		letter-spacing: -2px;
		line-height: 1.05;
		color: var(--color-text-primary);
	}

	.cta-lead {
		font-size: 0.9375rem;
		color: var(--color-text-secondary);
		max-width: 480px;
		line-height: 1.6;
	}

	.cta-actions {
		display: flex;
		gap: 10px;
		margin-top: 12px;
		flex-wrap: wrap;
		justify-content: center;
	}

	/* ═════════════════ RESPONSIVE ═════════════════ */
	@media (max-width: 1000px) {
		.hero-content {
			grid-template-columns: 1fr;
			gap: 40px;
		}

		.hero-right {
			display: none;
		}

		.stats-bar {
			grid-template-columns: repeat(2, 1fr);
		}

		.steps-grid {
			grid-template-columns: 1fr;
		}

		.step-connector {
			display: none;
		}

		.prof-grid,
		.faq-container {
			grid-template-columns: 1fr;
			gap: 40px;
		}

		.faq-left {
			position: static;
		}
	}

	@media (max-width: 640px) {
		.section {
			padding: 64px 0;
		}

		.hero {
			padding: 48px 0 64px;
		}

		.section-header {
			margin-bottom: 40px;
		}

		.prof-benefits {
			grid-template-columns: 1fr;
		}

		.stats-bar {
			grid-template-columns: 1fr;
		}

		.cta-box {
			padding: 48px 24px;
		}

		.press-logos {
			gap: 20px;
		}

		.press-logo {
			font-size: 0.8125rem;
		}
	}
</style>
