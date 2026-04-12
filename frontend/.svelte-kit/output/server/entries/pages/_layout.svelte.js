import { n as noop, a as attr_class, s as store_get, e as escape_html, u as unsubscribe_stores, f as fallback, b as bind_props, c as stringify, d as ensure_array_like, g as slot } from "../../chunks/renderer.js";
import "@sveltejs/kit/internal";
import "../../chunks/exports.js";
import "../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../chunks/root.js";
import "../../chunks/state.svelte.js";
import { i as isAuthenticated, c as currentUser } from "../../chunks/auth.js";
import "clsx";
import { t as toasts } from "../../chunks/toasts.js";
function createEventDispatcher() {
  return noop;
}
function Navbar($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    var $$store_subs;
    let menuOpen = false;
    $$renderer2.push(`<header class="navbar svelte-rfuq4y"><div class="container navbar-inner svelte-rfuq4y"><a href="/" class="navbar-brand svelte-rfuq4y"><span class="brand-icon svelte-rfuq4y">✦</span> <span class="brand-name svelte-rfuq4y">DiaryGo</span></a> <nav${attr_class("navbar-links svelte-rfuq4y", void 0, { "open": menuOpen })}><a href="/" class="nav-link text-nav svelte-rfuq4y">Início</a> `);
    if (store_get($$store_subs ??= {}, "$isAuthenticated", isAuthenticated)) {
      $$renderer2.push("<!--[0-->");
      if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "CLIENTE") {
        $$renderer2.push("<!--[0-->");
        $$renderer2.push(`<a href="/dashboard" class="nav-link text-nav svelte-rfuq4y">Dashboard</a> <a href="/solicitacoes" class="nav-link text-nav svelte-rfuq4y">Serviços</a>`);
      } else if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "PROFISSIONAL") {
        $$renderer2.push("<!--[1-->");
        $$renderer2.push(`<a href="/dashboard" class="nav-link text-nav svelte-rfuq4y">Dashboard</a> <a href="/agenda" class="nav-link text-nav svelte-rfuq4y">Agenda</a>`);
      } else if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "ADMIN") {
        $$renderer2.push("<!--[2-->");
        $$renderer2.push(`<a href="/admin" class="nav-link text-nav svelte-rfuq4y">Admin</a>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]-->`);
    } else {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`<a href="/como-funciona" class="nav-link text-nav svelte-rfuq4y">Como funciona</a> <a href="/profissionais" class="nav-link text-nav svelte-rfuq4y">Para diaristas</a>`);
    }
    $$renderer2.push(`<!--]--></nav> <div class="navbar-actions svelte-rfuq4y">`);
    if (store_get($$store_subs ??= {}, "$isAuthenticated", isAuthenticated)) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span class="badge badge-blue text-small">${escape_html(store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo)}</span> <button class="btn btn-primary btn-sm">Sair</button>`);
    } else {
      $$renderer2.push("<!--[-1-->");
      $$renderer2.push(`<a href="/login" class="btn btn-primary btn-sm">Entrar</a> <a href="/registro" class="btn btn-white btn-sm">Começar</a>`);
    }
    $$renderer2.push(`<!--]--> <button class="menu-toggle svelte-rfuq4y" aria-label="Menu"><span class="svelte-rfuq4y"></span> <span class="svelte-rfuq4y"></span> <span class="svelte-rfuq4y"></span></button></div></div></header>`);
    if ($$store_subs) unsubscribe_stores($$store_subs);
  });
}
function Toast($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let message = $$props["message"];
    let type = fallback($$props["type"], "info");
    let duration = fallback($$props["duration"], 4e3);
    const dispatch = createEventDispatcher();
    let visible = true;
    {
      if (duration > 0) {
        setTimeout(
          () => {
            visible = false;
            setTimeout(() => dispatch("close"), 300);
          },
          duration
        );
      }
    }
    if (visible) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<div${attr_class(`toast toast-${stringify(type)} animate-in`, "svelte-1cpok13")} role="alert"><span class="toast-icon svelte-1cpok13">`);
      if (type === "success") {
        $$renderer2.push("<!--[0-->");
        $$renderer2.push(`✓`);
      } else if (type === "error") {
        $$renderer2.push("<!--[1-->");
        $$renderer2.push(`✕`);
      } else {
        $$renderer2.push("<!--[-1-->");
        $$renderer2.push(`ℹ`);
      }
      $$renderer2.push(`<!--]--></span> <span class="toast-message svelte-1cpok13">${escape_html(message)}</span> <button class="toast-close svelte-1cpok13">✕</button></div>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]-->`);
    bind_props($$props, { message, type, duration });
  });
}
function ToastContainer($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    var $$store_subs;
    $$renderer2.push(`<div class="toast-container svelte-cqwvc2"><!--[-->`);
    const each_array = ensure_array_like(store_get($$store_subs ??= {}, "$toasts", toasts));
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let toast = each_array[$$index];
      Toast($$renderer2, { message: toast.message, type: toast.type });
    }
    $$renderer2.push(`<!--]--></div>`);
    if ($$store_subs) unsubscribe_stores($$store_subs);
  });
}
function _layout($$renderer, $$props) {
  Navbar($$renderer);
  $$renderer.push(`<!----> <main class="svelte-12qhfyh"><!--[-->`);
  slot($$renderer, $$props, "default", {});
  $$renderer.push(`<!--]--></main> `);
  ToastContainer($$renderer);
  $$renderer.push(`<!---->`);
}
export {
  _layout as default
};
