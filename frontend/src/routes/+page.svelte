<script lang="ts">
	import { isAuthenticated, currentUser } from '$lib/stores/auth';

	const steps = [
		{
			title: 'Solicite o serviço',
			desc: 'Escolha o tipo de limpeza, informe os cômodos e selecione data e horário. O sistema calcula o valor automaticamente.'
		},
		{
			title: 'Profissional atribuída',
			desc: 'Em instantes, atribuímos a melhor diarista disponível na sua região com base em avaliação e proximidade.'
		},
		{
			title: 'Acompanhe e avalie',
			desc: 'Receba notificações de check-in e conclusão. Após o serviço, avalie a profissional.'
		}
	];

	const services = [
		{ icon: '🧹', name: 'Limpeza Padrão', desc: 'Faxina de rotina completa', duration: 'A partir de 4h' },
		{ icon: '✨', name: 'Limpeza Pesada', desc: 'Limpeza profunda e detalhada', duration: 'A partir de 6h' },
		{ icon: '🏗️', name: 'Pós-Obra', desc: 'Remoção de resíduos de reforma', duration: 'A partir de 8h' },
		{ icon: '📦', name: 'Pré-Mudança', desc: 'Preparar o imóvel para entrada', duration: 'A partir de 6h' },
		{ icon: '🏢', name: 'Comercial', desc: 'Escritórios e salas comerciais', duration: 'A partir de 4h' },
		{ icon: '⚡', name: 'Express', desc: 'Manutenção leve e rápida', duration: 'A partir de 1h30' },
		{ icon: '👔', name: 'Passadoria', desc: 'Passar roupas por hora ou peça', duration: 'A partir de 2h' }
	];
</script>

<svelte:head>
	<title>DiaryGo — Diaristas de confiança</title>
</svelte:head>

<!-- Hero -->
<section class="hero section">
	<div class="container">
		<div class="hero-content">
			<div class="badge badge-orange hero-badge">
				<span>✦</span>
				<span>Plataforma intermediadora de diaristas</span>
			</div>

			<h1 class="text-display hero-title">
				Sua casa<br />
				<span class="title-accent">impecável.</span><br />
				Sempre.
			</h1>

			<p class="text-body-lg hero-subtitle">
				Conectamos você às melhores diaristas da sua região.
				Preço justo, agendamento fácil e profissionais verificadas.
			</p>

			<div class="hero-actions row gap-3">
				{#if $isAuthenticated}
					<a href="/dashboard" class="btn btn-white btn-lg">
						Ir para o dashboard
					</a>
				{:else}
					<a href="/registro" class="btn btn-white btn-lg">
						Contratar agora
					</a>
					<a href="/login" class="btn btn-primary btn-lg">
						Já tenho conta
					</a>
				{/if}
			</div>

			<div class="hero-stats row gap-8">
				<div class="stat">
					<span class="stat-number">4.9</span>
					<span class="stat-label text-small">Avaliação média</span>
				</div>
				<div class="stat-divider"></div>
				<div class="stat">
					<span class="stat-number">500+</span>
					<span class="stat-label text-small">Profissionais ativas</span>
				</div>
				<div class="stat-divider"></div>
				<div class="stat">
					<span class="stat-number">98%</span>
					<span class="stat-label text-small">Clientes satisfeitos</span>
				</div>
			</div>
		</div>
	</div>

	<div class="hero-glow" aria-hidden="true"></div>
</section>

<!-- Como funciona -->
<section class="section">
	<div class="container">
		<div class="section-header">
			<h2 class="text-heading">Como funciona</h2>
			<p class="text-body-lg">Em três passos simples</p>
		</div>

		<div class="steps-grid">
			{#each steps as step, i}
				<div class="step-card card">
					<div class="step-number badge badge-orange">{i + 1}</div>
					<h3 class="text-feature-title">{step.title}</h3>
					<p class="text-body">{step.desc}</p>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- Tipos de serviço -->
<section class="section">
	<div class="container">
		<div class="section-header">
			<h2 class="text-heading">Serviços disponíveis</h2>
			<p class="text-body-lg">Do padrão ao especializado</p>
		</div>

		<div class="services-grid">
			{#each services as svc}
				<div class="service-card card row gap-4">
					<span class="service-icon">{svc.icon}</span>
					<div class="stack gap-1">
						<h3 class="text-subheading">{svc.name}</h3>
						<p class="text-caption">{svc.desc}</p>
						<div class="row gap-2" style="margin-top: 4px;">
							<span class="badge badge-blue">{svc.duration}</span>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- CTA Final -->
<section class="section cta-section">
	<div class="container">
		<div class="cta-box card card-section">
			<div class="cta-glow" aria-hidden="true"></div>
			<h2 class="text-heading">Pronta para começar?</h2>
			<p class="text-body-lg">
				Crie sua conta gratuitamente e agende seu primeiro serviço hoje.
			</p>
			<div class="row gap-3" style="justify-content: center; flex-wrap: wrap;">
				<a href="/registro" class="btn btn-white btn-lg">Sou cliente</a>
				<a href="/registro/profissional" class="btn btn-primary btn-lg">Sou diarista</a>
			</div>
		</div>
	</div>
</section>


<style>
	/* Hero */
	.hero {
		position: relative;
		overflow: hidden;
		padding-top: var(--space-24);
		padding-bottom: var(--space-24);
	}

	.hero-content {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: var(--space-6);
		max-width: 680px;
		position: relative;
		z-index: 1;
	}

	.hero-badge {
		font-size: 0.75rem;
	}

	.hero-title {
		font-size: clamp(3.5rem, 8vw, 6rem);
	}

	.title-accent {
		color: var(--color-orange-10);
	}

	.hero-subtitle {
		max-width: 480px;
	}

	.hero-actions {
		flex-wrap: wrap;
	}

	.hero-stats {
		margin-top: var(--space-4);
		align-items: center;
		flex-wrap: wrap;
		gap: var(--space-6);
	}

	.stat {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.stat-number {
		font-size: 1.5rem;
		font-weight: 600;
		color: var(--color-text-primary);
		letter-spacing: -0.5px;
	}

	.stat-label {
		color: var(--color-text-secondary);
	}

	.stat-divider {
		width: 1px;
		height: 32px;
		background: var(--border-frost);
	}

	.hero-glow {
		position: absolute;
		top: 0;
		left: -200px;
		width: 600px;
		height: 400px;
		background: radial-gradient(ellipse at center, rgba(255, 128, 31, 0.08) 0%, transparent 70%);
		pointer-events: none;
		z-index: 0;
	}

	/* Section header */
	.section-header {
		margin-bottom: var(--space-10);
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.section-header .text-body-lg {
		max-width: 480px;
	}

	/* Steps */
	.steps-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
		gap: var(--space-4);
	}

	.step-card {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}

	.step-number {
		width: fit-content;
	}

	/* Services */
	.services-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: var(--space-4);
	}

	.service-card {
		padding: var(--space-4) var(--space-5);
		align-items: flex-start;
		transition: border-color 0.15s ease;
	}

	.service-card:hover {
		border-color: rgba(214, 235, 253, 0.35);
	}

	.service-icon {
		font-size: 1.5rem;
		flex-shrink: 0;
	}

	/* CTA */
	.cta-box {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: var(--space-6);
		padding: var(--space-16) var(--space-6);
		position: relative;
		overflow: hidden;
	}

	.cta-glow {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: 400px;
		height: 200px;
		background: radial-gradient(ellipse at center, rgba(255, 128, 31, 0.06) 0%, transparent 70%);
		pointer-events: none;
	}
</style>
