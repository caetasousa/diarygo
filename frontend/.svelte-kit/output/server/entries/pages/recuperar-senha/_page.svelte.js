import { i as head, j as attr } from "../../../chunks/renderer.js";
import "../../../chunks/toasts.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let email = "";
    let loading = false;
    head("hu0xp4", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Recuperar senha — DiaryGo</title>`);
      });
    });
    $$renderer2.push(`<div class="auth-page section svelte-hu0xp4"><div class="container-sm"><div class="auth-header svelte-hu0xp4"><a href="/" class="auth-brand svelte-hu0xp4"><span style="color: var(--color-orange-10)">✦</span> DiaryGo</a> <h1 class="text-heading auth-title svelte-hu0xp4">Recuperar senha</h1> <p class="text-body">Informe seu email e enviaremos instruções para redefinir sua senha.</p></div> `);
    {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`<form class="auth-form card svelte-hu0xp4" novalidate=""><div class="form-group"><label for="email" class="label">Email</label> <input id="email" type="email" class="input" placeholder="seu@email.com"${attr("value", email)}${attr("disabled", loading, true)} autocomplete="email"/></div> <button type="submit" class="btn btn-white btn-lg" style="width: 100%;"${attr("disabled", !email, true)}>`);
      {
        $$renderer2.push("<!--[-1-->");
        $$renderer2.push(`Enviar instruções`);
      }
      $$renderer2.push(`<!--]--></button></form> <p class="text-body" style="text-align: center; margin-top: var(--space-6);">Lembrou a senha? <a href="/login">Entrar</a></p>`);
    }
    $$renderer2.push(`<!--]--></div></div>`);
  });
}
export {
  _page as default
};
