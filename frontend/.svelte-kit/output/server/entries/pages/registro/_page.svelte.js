import { i as head, j as attr, e as escape_html, d as ensure_array_like, a as attr_class } from "../../../chunks/renderer.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/root.js";
import "../../../chunks/state.svelte.js";
import "../../../chunks/toasts.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let email = "";
    let senha = "";
    let confirmar = "";
    let loading = false;
    let errors = {};
    head("j7ul7i", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Criar conta — DiaryGo</title>`);
      });
    });
    $$renderer2.push(`<div class="auth-page section svelte-j7ul7i"><div class="container-sm"><div class="auth-header svelte-j7ul7i"><a href="/" class="auth-brand svelte-j7ul7i"><span style="color: var(--color-orange-10)">✦</span> DiaryGo</a> <div class="row gap-3" style="align-items: center; flex-wrap: wrap;"><h1 class="text-heading auth-title svelte-j7ul7i">Criar conta</h1> <span class="badge badge-green">Cliente</span></div> <p class="text-body">Crie sua conta e agende seu primeiro serviço.</p></div> <form class="auth-form card svelte-j7ul7i" novalidate=""><div class="form-group"><label for="email" class="label">Email</label> <input id="email" type="email" class="input" placeholder="seu@email.com"${attr("value", email)}${attr("disabled", loading, true)} autocomplete="email"/> `);
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
    $$renderer2.push(`<!--]--></div> <div class="password-strength svelte-j7ul7i"><div class="strength-bar svelte-j7ul7i"><!--[-->`);
    const each_array = ensure_array_like(Array(4));
    for (let i = 0, $$length = each_array.length; i < $$length; i++) {
      each_array[i];
      $$renderer2.push(`<div${attr_class("strength-segment svelte-j7ul7i", void 0, {
        "active": senha.length > i * 18,
        "strong": senha.length >= 12
      })}></div>`);
    }
    $$renderer2.push(`<!--]--></div> <span class="text-small" style="color: var(--color-text-tertiary);">`);
    if (senha.length === 0) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`Digite uma senha`);
    } else if (senha.length < 8) {
      $$renderer2.push("<!--[1-->");
      $$renderer2.push(`Muito curta`);
    } else if (senha.length < 12) {
      $$renderer2.push("<!--[2-->");
      $$renderer2.push(`Razoável`);
    } else if (senha.length < 20) {
      $$renderer2.push("<!--[3-->");
      $$renderer2.push(`Boa`);
    } else {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`Excelente`);
    }
    $$renderer2.push(`<!--]--></span></div> <button type="submit" class="btn btn-white btn-lg" style="width: 100%;"${attr("disabled", loading, true)}>`);
    {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`Criar conta de cliente`);
    }
    $$renderer2.push(`<!--]--></button> <p class="text-caption" style="text-align: center; color: var(--color-text-tertiary);">Ao criar conta, você concorda com nossos termos de uso.</p></form> <div class="divider" style="margin: var(--space-6) 0;"></div> <p class="text-body" style="text-align: center;">Já tem conta? <a href="/login">Entrar</a> · É diarista? <a href="/registro/profissional">Cadastre-se aqui</a></p></div></div>`);
  });
}
export {
  _page as default
};
