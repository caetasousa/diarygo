import { i as head, s as store_get, d as ensure_array_like, e as escape_html, u as unsubscribe_stores } from "../../chunks/renderer.js";
import { i as isAuthenticated } from "../../chunks/auth.js";
function _page($$renderer) {
  var $$store_subs;
  const steps = [
    {
      title: "Solicite o serviço",
      desc: "Escolha o tipo de limpeza, informe os cômodos e selecione data e horário. O sistema calcula o valor automaticamente."
    },
    {
      title: "Profissional atribuída",
      desc: "Em instantes, atribuímos a melhor diarista disponível na sua região com base em avaliação e proximidade."
    },
    {
      title: "Acompanhe e avalie",
      desc: "Receba notificações de check-in e conclusão. Após o serviço, avalie a profissional."
    }
  ];
  const services = [
    {
      icon: "🧹",
      name: "Limpeza Padrão",
      desc: "Faxina de rotina completa",
      duration: "A partir de 4h"
    },
    {
      icon: "✨",
      name: "Limpeza Pesada",
      desc: "Limpeza profunda e detalhada",
      duration: "A partir de 6h"
    },
    {
      icon: "🏗️",
      name: "Pós-Obra",
      desc: "Remoção de resíduos de reforma",
      duration: "A partir de 8h"
    },
    {
      icon: "📦",
      name: "Pré-Mudança",
      desc: "Preparar o imóvel para entrada",
      duration: "A partir de 6h"
    },
    {
      icon: "🏢",
      name: "Comercial",
      desc: "Escritórios e salas comerciais",
      duration: "A partir de 4h"
    },
    {
      icon: "⚡",
      name: "Express",
      desc: "Manutenção leve e rápida",
      duration: "A partir de 1h30"
    },
    {
      icon: "👔",
      name: "Passadoria",
      desc: "Passar roupas por hora ou peça",
      duration: "A partir de 2h"
    }
  ];
  head("1uha8ag", $$renderer, ($$renderer2) => {
    $$renderer2.title(($$renderer3) => {
      $$renderer3.push(`<title>DiaryGo — Diaristas de confiança</title>`);
    });
  });
  $$renderer.push(`<section class="hero section svelte-1uha8ag"><div class="container"><div class="hero-content svelte-1uha8ag"><div class="badge badge-orange hero-badge svelte-1uha8ag"><span>✦</span> <span>Plataforma intermediadora de diaristas</span></div> <h1 class="text-display hero-title svelte-1uha8ag">Sua casa<br/> <span class="title-accent svelte-1uha8ag">impecável.</span><br/> Sempre.</h1> <p class="text-body-lg hero-subtitle svelte-1uha8ag">Conectamos você às melhores diaristas da sua região.
				Preço justo, agendamento fácil e profissionais verificadas.</p> <div class="hero-actions row gap-3 svelte-1uha8ag">`);
  if (store_get($$store_subs ??= {}, "$isAuthenticated", isAuthenticated)) {
    $$renderer.push("<!--[0-->");
    $$renderer.push(`<a href="/dashboard" class="btn btn-white btn-lg">Ir para o dashboard</a>`);
  } else {
    $$renderer.push("<!--[-1-->");
    $$renderer.push(`<a href="/registro" class="btn btn-white btn-lg">Contratar agora</a> <a href="/login" class="btn btn-primary btn-lg">Já tenho conta</a>`);
  }
  $$renderer.push(`<!--]--></div> <div class="hero-stats row gap-8 svelte-1uha8ag"><div class="stat svelte-1uha8ag"><span class="stat-number svelte-1uha8ag">4.9</span> <span class="stat-label text-small svelte-1uha8ag">Avaliação média</span></div> <div class="stat-divider svelte-1uha8ag"></div> <div class="stat svelte-1uha8ag"><span class="stat-number svelte-1uha8ag">500+</span> <span class="stat-label text-small svelte-1uha8ag">Profissionais ativas</span></div> <div class="stat-divider svelte-1uha8ag"></div> <div class="stat svelte-1uha8ag"><span class="stat-number svelte-1uha8ag">98%</span> <span class="stat-label text-small svelte-1uha8ag">Clientes satisfeitos</span></div></div></div></div> <div class="hero-glow svelte-1uha8ag" aria-hidden="true"></div></section> <section class="section"><div class="container"><div class="section-header svelte-1uha8ag"><h2 class="text-heading">Como funciona</h2> <p class="text-body-lg svelte-1uha8ag">Em três passos simples</p></div> <div class="steps-grid svelte-1uha8ag"><!--[-->`);
  const each_array = ensure_array_like(steps);
  for (let i = 0, $$length = each_array.length; i < $$length; i++) {
    let step = each_array[i];
    $$renderer.push(`<div class="step-card card svelte-1uha8ag"><div class="step-number badge badge-orange svelte-1uha8ag">${escape_html(i + 1)}</div> <h3 class="text-feature-title">${escape_html(step.title)}</h3> <p class="text-body">${escape_html(step.desc)}</p></div>`);
  }
  $$renderer.push(`<!--]--></div></div></section> <section class="section"><div class="container"><div class="section-header svelte-1uha8ag"><h2 class="text-heading">Serviços disponíveis</h2> <p class="text-body-lg svelte-1uha8ag">Do padrão ao especializado</p></div> <div class="services-grid svelte-1uha8ag"><!--[-->`);
  const each_array_1 = ensure_array_like(services);
  for (let $$index_1 = 0, $$length = each_array_1.length; $$index_1 < $$length; $$index_1++) {
    let svc = each_array_1[$$index_1];
    $$renderer.push(`<div class="service-card card row gap-4 svelte-1uha8ag"><span class="service-icon svelte-1uha8ag">${escape_html(svc.icon)}</span> <div class="stack gap-1"><h3 class="text-subheading">${escape_html(svc.name)}</h3> <p class="text-caption">${escape_html(svc.desc)}</p> <div class="row gap-2" style="margin-top: 4px;"><span class="badge badge-blue">${escape_html(svc.duration)}</span></div></div></div>`);
  }
  $$renderer.push(`<!--]--></div></div></section> <section class="section cta-section"><div class="container"><div class="cta-box card card-section svelte-1uha8ag"><div class="cta-glow svelte-1uha8ag" aria-hidden="true"></div> <h2 class="text-heading">Pronta para começar?</h2> <p class="text-body-lg">Crie sua conta gratuitamente e agende seu primeiro serviço hoje.</p> <div class="row gap-3" style="justify-content: center; flex-wrap: wrap;"><a href="/registro" class="btn btn-white btn-lg">Sou cliente</a> <a href="/registro/profissional" class="btn btn-primary btn-lg">Sou diarista</a></div></div></div></section>`);
  if ($$store_subs) unsubscribe_stores($$store_subs);
}
export {
  _page as default
};
