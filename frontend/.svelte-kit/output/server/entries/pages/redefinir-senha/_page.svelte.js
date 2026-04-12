import { h as getContext, s as store_get, i as head, j as attr, e as escape_html, u as unsubscribe_stores } from "../../../chunks/renderer.js";
import "clsx";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/root.js";
import "../../../chunks/state.svelte.js";
import "../../../chunks/toasts.js";
const getStores = () => {
  const stores$1 = getContext("__svelte__");
  return {
    /** @type {typeof page} */
    page: {
      subscribe: stores$1.page.subscribe
    },
    /** @type {typeof navigating} */
    navigating: {
      subscribe: stores$1.navigating.subscribe
    },
    /** @type {typeof updated} */
    updated: stores$1.updated
  };
};
const page = {
  subscribe(fn) {
    const store = getStores().page;
    return store.subscribe(fn);
  }
};
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    var $$store_subs;
    let token = store_get($$store_subs ??= {}, "$page", page).url.searchParams.get("token") ?? "";
    let novaSenha = "";
    let confirmar = "";
    let loading = false;
    let errors = {};
    head("1d99bdn", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Redefinir senha — DiaryGo</title>`);
      });
    });
    $$renderer2.push(`<div class="auth-page section svelte-1d99bdn"><div class="container-sm"><div class="auth-header svelte-1d99bdn"><a href="/" class="auth-brand svelte-1d99bdn"><span style="color: var(--color-orange-10)">✦</span> DiaryGo</a> <h1 class="text-heading auth-title svelte-1d99bdn">Nova senha</h1> <p class="text-body">Escolha uma senha forte para sua conta.</p></div> <form class="auth-form card svelte-1d99bdn" novalidate="">`);
    if (!store_get($$store_subs ??= {}, "$page", page).url.searchParams.get("token")) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<div class="form-group"><label for="token" class="label">Token de recuperação</label> <input id="token" type="text" class="input" placeholder="Cole o token recebido"${attr("value", token)}${attr("disabled", loading, true)}/> `);
      if (errors.token) {
        $$renderer2.push("<!--[0-->");
        $$renderer2.push(`<span class="form-error">${escape_html(errors.token)}</span>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--> <div class="form-group"><label for="novaSenha" class="label">Nova senha</label> <input id="novaSenha" type="password" class="input" placeholder="Mínimo 8 caracteres"${attr("value", novaSenha)}${attr("disabled", loading, true)} autocomplete="new-password"/> `);
    if (errors.novaSenha) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="form-error">${escape_html(errors.novaSenha)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="form-group"><label for="confirmar" class="label">Confirmar nova senha</label> <input id="confirmar" type="password" class="input" placeholder="••••••••"${attr("value", confirmar)}${attr("disabled", loading, true)} autocomplete="new-password"/> `);
    if (errors.confirmar) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="form-error">${escape_html(errors.confirmar)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <button type="submit" class="btn btn-white btn-lg" style="width: 100%;"${attr("disabled", loading, true)}>`);
    {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`Redefinir senha`);
    }
    $$renderer2.push(`<!--]--></button></form> <p class="text-body" style="text-align: center; margin-top: var(--space-6);"><a href="/login">← Voltar ao login</a></p></div></div>`);
    if ($$store_subs) unsubscribe_stores($$store_subs);
  });
}
export {
  _page as default
};
