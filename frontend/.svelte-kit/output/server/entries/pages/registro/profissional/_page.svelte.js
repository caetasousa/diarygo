import { i as head, d as ensure_array_like, e as escape_html, j as attr } from "../../../../chunks/renderer.js";
import "@sveltejs/kit/internal";
import "../../../../chunks/exports.js";
import "../../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../../chunks/root.js";
import "../../../../chunks/state.svelte.js";
import "../../../../chunks/toasts.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let email = "";
    let senha = "";
    let confirmar = "";
    let loading = false;
    let errors = {};
    const requirements = [
      {
        icon: "📄",
        label: "RG (frente e verso) — enviado após aprovação inicial"
      },
      {
        icon: "🏠",
        label: "Comprovante de residência recente"
      },
      { icon: "📸", label: "Foto do rosto (selfie nítida)" },
      {
        icon: "📞",
        label: "Mínimo 2 referências profissionais com telefone"
      },
      {
        icon: "📍",
        label: "CEP das regiões onde deseja atender"
      }
    ];
    head("1q9le52", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Cadastro de Diarista — DiaryGo</title>`);
      });
    });
    $$renderer2.push(`<div class="auth-page section svelte-1q9le52"><div class="container-sm"><div class="auth-header svelte-1q9le52"><a href="/" class="auth-brand svelte-1q9le52"><span style="color: var(--color-orange-10)">✦</span> DiaryGo</a> <div class="row gap-3" style="align-items: center; flex-wrap: wrap;"><h1 class="text-heading auth-title svelte-1q9le52">Seja diarista</h1> <span class="badge badge-orange">Profissional</span></div> <p class="text-body">Cadastre-se e comece a receber solicitações na sua região.
				Seu perfil será analisado em até 48h úteis.</p></div> <div class="info-card card svelte-1q9le52" style="margin-bottom: var(--space-6);"><h3 class="text-subheading" style="margin-bottom: var(--space-4);">O que você precisará</h3> <ul class="requirements-list svelte-1q9le52"><!--[-->`);
    const each_array = ensure_array_like(requirements);
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let req = each_array[$$index];
      $$renderer2.push(`<li class="requirement-item row gap-3"><span class="req-icon svelte-1q9le52">${escape_html(req.icon)}</span> <span class="text-caption">${escape_html(req.label)}</span></li>`);
    }
    $$renderer2.push(`<!--]--></ul></div> <form class="auth-form card svelte-1q9le52" novalidate=""><div class="form-group"><label for="email" class="label">Email profissional</label> <input id="email" type="email" class="input" placeholder="seu@email.com"${attr("value", email)}${attr("disabled", loading, true)} autocomplete="email"/> `);
    if (errors.email) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="form-error">${escape_html(errors.email)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="form-group"><label for="senha" class="label">Senha</label> <input id="senha" type="password" class="input" placeholder="Mínimo 8 caracteres"${attr("value", senha)}${attr("disabled", loading, true)} autocomplete="new-password"/> `);
    if (errors.senha) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="form-error">${escape_html(errors.senha)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="form-group"><label for="confirmar" class="label">Confirmar senha</label> <input id="confirmar" type="password" class="input" placeholder="••••••••"${attr("value", confirmar)}${attr("disabled", loading, true)} autocomplete="new-password"/> `);
    if (errors.confirmar) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="form-error">${escape_html(errors.confirmar)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <button type="submit" class="btn btn-white btn-lg" style="width: 100%;"${attr("disabled", loading, true)}>`);
    {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`Enviar cadastro`);
    }
    $$renderer2.push(`<!--]--></button> <p class="text-caption" style="text-align: center; color: var(--color-text-tertiary);">Após o envio, sua documentação será analisada em até 48h úteis.</p></form> <div class="divider" style="margin: var(--space-6) 0;"></div> <p class="text-body" style="text-align: center;">Já tem conta? <a href="/login">Entrar</a> · Precisa de serviço? <a href="/registro">Sou cliente</a></p></div></div>`);
  });
}
export {
  _page as default
};
