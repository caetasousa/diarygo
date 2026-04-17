<script lang="ts">
	import { Sparkles, ShieldCheck, Lock, Mail } from 'lucide-svelte';

	const year = new Date().getFullYear();

	const columns = [
		{
			title: 'Plataforma',
			links: [
				{ href: '/#servicos', label: 'Serviços' },
				{ href: '/#como-funciona', label: 'Como funciona' },
				{ href: '/#diferenciais', label: 'Diferenciais' },
				{ href: '/#faq', label: 'Perguntas frequentes' }
			]
		},
		{
			title: 'Para diaristas',
			links: [
				{ href: '/#para-diaristas', label: 'Por que trabalhar conosco' },
				{ href: '/registro/profissional', label: 'Cadastre-se' },
				{ href: '/login', label: 'Área da profissional' }
			]
		},
		{
			title: 'Empresa',
			links: [
				{ href: '/sobre', label: 'Sobre' },
				{ href: '/imprensa', label: 'Imprensa' },
				{ href: '/contato', label: 'Contato' },
				{ href: '/ajuda', label: 'Central de ajuda' }
			]
		},
		{
			title: 'Legal',
			links: [
				{ href: '/termos', label: 'Termos de uso' },
				{ href: '/privacidade', label: 'Privacidade' },
				{ href: '/lgpd', label: 'LGPD' },
				{ href: '/cookies', label: 'Política de cookies' }
			]
		}
	];
</script>

<footer class="footer">
	<div class="footer-top">
		<div class="footer-brand-col">
			<a href="/" class="footer-brand">
				<span class="brand-mark">
					<Sparkles size={14} strokeWidth={2.2} />
				</span>
				<span class="brand-name">DiaryGo</span>
			</a>
			<p class="footer-tag">
				A plataforma que conecta você às melhores diaristas verificadas da sua região.
			</p>
			<div class="trust-badges">
				<div class="trust-badge">
					<ShieldCheck size={12} strokeWidth={2} />
					<span>Profissionais verificadas</span>
				</div>
				<div class="trust-badge">
					<Lock size={12} strokeWidth={2} />
					<span>Pagamento seguro</span>
				</div>
			</div>
		</div>

		<div class="footer-grid">
			{#each columns as col}
				<div class="footer-col">
					<h4 class="footer-col-title">{col.title}</h4>
					<ul class="footer-links">
						{#each col.links as link}
							<li>
								<a href={link.href} class="footer-link">{link.label}</a>
							</li>
						{/each}
					</ul>
				</div>
			{/each}
		</div>
	</div>

	<div class="footer-bottom">
		<div class="footer-bottom-inner">
			<span class="copyright">© {year} DiaryGo. Todos os direitos reservados.</span>
			<div class="socials">
				<a href="/" class="social-btn" aria-label="Instagram">
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="20" rx="5" ry="5"/><path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"/><line x1="17.5" y1="6.5" x2="17.51" y2="6.5"/></svg>
				</a>
				<a href="/" class="social-btn" aria-label="Facebook">
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"/></svg>
				</a>
				<a href="/" class="social-btn" aria-label="LinkedIn">
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-4 0v7h-4v-7a6 6 0 0 1 6-6z"/><rect x="2" y="9" width="4" height="12"/><circle cx="4" cy="4" r="2"/></svg>
				</a>
				<a href="mailto:contato@diarygo.com" class="social-btn" aria-label="E-mail">
					<Mail size={14} strokeWidth={1.8} />
				</a>
			</div>
		</div>
	</div>
</footer>

<style>
	.footer {
		border-top: 1px solid var(--border-frost);
		background:
			radial-gradient(80% 60% at 50% 0%, rgba(255, 128, 31, 0.03), transparent 70%),
			#000;
		margin-top: 80px;
	}

	.footer-top {
		max-width: 1240px;
		margin: 0 auto;
		padding: 64px 24px 48px;
		display: grid;
		grid-template-columns: 1.3fr 2.2fr;
		gap: 64px;
		align-items: start;
	}

	.footer-brand-col {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 320px;
	}

	.footer-brand {
		display: inline-flex;
		align-items: center;
		gap: 9px;
		text-decoration: none;
		color: var(--color-text-primary);
	}

	.brand-mark {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		border-radius: 7px;
		background: linear-gradient(135deg, var(--color-orange-10), #ff6b1a);
		color: #1a0a00;
	}

	.brand-name {
		font-size: 1rem;
		font-weight: 600;
		letter-spacing: -0.3px;
	}

	.footer-tag {
		font-size: 0.8125rem;
		color: var(--color-text-tertiary);
		line-height: 1.55;
	}

	.trust-badges {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
		margin-top: 4px;
	}

	.trust-badge {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 6px 10px;
		background: rgba(255, 255, 255, 0.025);
		border: 1px solid var(--border-frost);
		border-radius: 99px;
		font-size: 0.6875rem;
		color: var(--color-text-secondary);
		font-weight: 500;
	}

	.trust-badge :global(svg) {
		color: var(--color-green-10);
	}

	.footer-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 32px;
	}

	.footer-col-title {
		font-size: 0.6875rem;
		font-weight: 600;
		color: var(--color-text-primary);
		text-transform: uppercase;
		letter-spacing: 0.8px;
		margin-bottom: 16px;
	}

	.footer-links {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.footer-link {
		font-size: 0.8125rem;
		color: var(--color-text-tertiary);
		text-decoration: none;
		transition: color 0.12s;
	}

	.footer-link:hover {
		color: var(--color-text-primary);
		opacity: 1;
	}

	.footer-bottom {
		border-top: 1px solid var(--border-frost);
	}

	.footer-bottom-inner {
		max-width: 1240px;
		margin: 0 auto;
		padding: 20px 24px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		flex-wrap: wrap;
	}

	.copyright {
		font-size: 0.75rem;
		color: var(--color-text-tertiary);
	}

	.socials {
		display: flex;
		gap: 6px;
	}

	.social-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid var(--border-frost);
		border-radius: 7px;
		color: var(--color-text-tertiary);
		text-decoration: none;
		transition: color 0.12s, background 0.12s, border-color 0.12s;
	}

	.social-btn:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.06);
		border-color: rgba(214, 235, 253, 0.35);
		opacity: 1;
	}

	@media (max-width: 960px) {
		.footer-top {
			grid-template-columns: 1fr;
			gap: 40px;
			padding: 48px 24px 32px;
		}

		.footer-grid {
			grid-template-columns: repeat(2, 1fr);
			gap: 28px;
		}
	}

	@media (max-width: 480px) {
		.footer-grid {
			grid-template-columns: 1fr;
			gap: 24px;
		}
	}
</style>
