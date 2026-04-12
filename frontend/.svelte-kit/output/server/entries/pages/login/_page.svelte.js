import { i as head, j as attr, e as escape_html } from "../../../chunks/renderer.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/root.js";
import "../../../chunks/state.svelte.js";
import "../../../chunks/auth.js";
import "../../../chunks/toasts.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let email = "";
    let senha = "";
    let loading = false;
    let errors = {};
    head("1x05zx6", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Entrar — DiaryGo</title>`);
      });
    });
    $$renderer2.push(`<div class="auth-page section svelte-1x05zx6"><div class="container-sm"><div class="auth-header svelte-1x05zx6"><a href="/" class="auth-brand svelte-1x05zx6"><span style="color: var(--color-orange-10)">✦</span> DiaryGo</a> <h1 class="text-heading auth-title svelte-1x05zx6">Entrar</h1> <p class="text-body">Acesse sua conta para continuar.</p></div> <form class="auth-form card svelte-1x05zx6" novalidate=""><div class="form-group"><label for="email" class="label">Email</label> <input id="email" type="email" class="input" placeholder="seu@email.com"${attr("value", email)}${attr("disabled", loading, true)} autocomplete="email"/> `);
    if (errors.email) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="form-error">${escape_html(errors.email)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="form-group"><div class="row" style="justify-content: space-between; align-items: center;"><label for="senha" class="label" style="margin: 0;">Senha</label> <a href="/recuperar-senha" class="text-caption" style="color: var(--color-blue-10);">Esqueceu a senha?</a></div> <input id="senha" type="password" class="input" placeholder="••••••••"${attr("value", senha)}${attr("disabled", loading, true)} autocomplete="current-password"/> `);
    if (errors.senha) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="form-error">${escape_html(errors.senha)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <button type="submit" class="btn btn-white btn-lg" style="width: 100%;"${attr("disabled", loading, true)}>`);
    {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`Entrar`);
    }
    $$renderer2.push(`<!--]--></button></form> <div class="divider" style="margin: var(--space-6) 0;"></div> <p class="text-body" style="text-align: center;">Não tem conta? <a href="/registro">Criar conta de cliente</a> · <a href="/registro/profissional">Sou diarista</a></p></div></div>`);
  });
}
export {
  _page as default
};
