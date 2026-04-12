import { s as store_get, i as head, a as attr_class, c as stringify, e as escape_html, d as ensure_array_like, u as unsubscribe_stores, j as attr } from "../../../chunks/renderer.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/root.js";
import "../../../chunks/state.svelte.js";
import { c as currentUser } from "../../../chunks/auth.js";
function ClienteDashboard($$renderer) {
  $$renderer.push(`<div class="dash-section svelte-x1i5gj"><div class="section-title-row svelte-x1i5gj"><h2 class="text-feature-title">Próximos serviços</h2> <a href="/solicitacoes" class="btn btn-ghost btn-sm">Ver todos</a></div> <div class="empty-state card svelte-x1i5gj"><span class="empty-icon svelte-x1i5gj">🧹</span> <p class="text-subheading">Nenhum serviço agendado</p> <p class="text-body">Solicite seu primeiro serviço e mantenha sua casa impecável.</p> <a href="/solicitacoes/novo" class="btn btn-white">Solicitar agora</a></div></div>`);
}
function ProfissionalDashboard($$renderer) {
  $$renderer.push(`<div class="dash-section svelte-x1i5gj"><div class="pending-approval card svelte-x1i5gj" style="border-color: var(--color-orange-4);"><div class="row gap-4" style="align-items: flex-start;"><span style="font-size: 1.5rem;">⏳</span> <div class="stack gap-2"><h3 class="text-subheading">Cadastro em análise</h3> <p class="text-body">Sua documentação está sendo analisada pela equipe DiaryGo.
						Você receberá uma notificação em até 48h úteis.</p> <span class="badge badge-orange">PENDENTE</span></div></div></div> <div class="section-title-row svelte-x1i5gj" style="margin-top: var(--space-8);"><h2 class="text-feature-title">Próximos serviços</h2></div> <div class="empty-state card svelte-x1i5gj"><span class="empty-icon svelte-x1i5gj">📅</span> <p class="text-subheading">Nenhum serviço agendado</p> <p class="text-body">Após aprovação, você começará a receber solicitações na sua região.</p></div></div>`);
}
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    var $$store_subs;
    let stats;
    const adminActions = [
      {
        icon: "👔",
        title: "Profissionais",
        desc: "Aprovar e gerenciar cadastros",
        href: "/admin/profissionais"
      },
      {
        icon: "👥",
        title: "Clientes",
        desc: "Ver e gerenciar clientes",
        href: "/admin/clientes"
      },
      {
        icon: "📋",
        title: "Serviços",
        desc: "Acompanhar todos os serviços",
        href: "/admin/servicos"
      },
      {
        icon: "💰",
        title: "Tabela de preços",
        desc: "Configurar preços por região",
        href: "/admin/precos"
      }
    ];
    const statsByTipo = {
      CLIENTE: [
        {
          label: "Serviços realizados",
          value: "0",
          icon: "✓",
          color: "green"
        },
        {
          label: "Agendados",
          value: "0",
          icon: "📅",
          color: "blue"
        },
        {
          label: "Score",
          value: "100",
          icon: "⭐",
          color: "yellow"
        }
      ],
      PROFISSIONAL: [
        {
          label: "Serviços realizados",
          value: "0",
          icon: "✓",
          color: "green"
        },
        {
          label: "Nota média",
          value: "—",
          icon: "⭐",
          color: "yellow"
        },
        {
          label: "Status",
          value: "PENDENTE",
          icon: "⏳",
          color: "orange"
        }
      ],
      ADMIN: [
        {
          label: "Usuários",
          value: "—",
          icon: "👥",
          color: "blue"
        },
        {
          label: "Serviços hoje",
          value: "—",
          icon: "📋",
          color: "green"
        },
        {
          label: "Pendentes",
          value: "—",
          icon: "⚠️",
          color: "orange"
        }
      ]
    };
    stats = store_get($$store_subs ??= {}, "$currentUser", currentUser) ? statsByTipo[store_get($$store_subs ??= {}, "$currentUser", currentUser).tipo] ?? [] : [];
    function AdminDashboard($$renderer3) {
      $$renderer3.push(`<div class="dash-section svelte-x1i5gj"><div class="admin-grid svelte-x1i5gj"><!--[-->`);
      const each_array = ensure_array_like(adminActions);
      for (let $$index_1 = 0, $$length = each_array.length; $$index_1 < $$length; $$index_1++) {
        let action = each_array[$$index_1];
        $$renderer3.push(`<a${attr("href", action.href)} class="admin-action card row gap-4 svelte-x1i5gj"><span class="admin-action-icon svelte-x1i5gj">${escape_html(action.icon)}</span> <div class="stack gap-1"><h3 class="text-subheading">${escape_html(action.title)}</h3> <p class="text-caption">${escape_html(action.desc)}</p></div></a>`);
      }
      $$renderer3.push(`<!--]--></div></div>`);
    }
    head("x1i5gj", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Dashboard — DiaryGo</title>`);
      });
    });
    $$renderer2.push(`<div class="section"><div class="container"><div class="dashboard-header svelte-x1i5gj"><div><div class="row gap-3" style="align-items: center; margin-bottom: var(--space-2);"><h1 class="text-heading dash-title svelte-x1i5gj">Dashboard</h1> `);
    if (store_get($$store_subs ??= {}, "$currentUser", currentUser)) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<span${attr_class(`badge badge-${stringify(store_get($$store_subs ??= {}, "$currentUser", currentUser).tipo === "CLIENTE" ? "green" : store_get($$store_subs ??= {}, "$currentUser", currentUser).tipo === "PROFISSIONAL" ? "orange" : "blue")}`)}>${escape_html(store_get($$store_subs ??= {}, "$currentUser", currentUser).tipo)}</span>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <p class="text-body">`);
    if (store_get($$store_subs ??= {}, "$currentUser", currentUser)) {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`Olá, <strong style="color: var(--color-text-primary)">${escape_html(store_get($$store_subs ??= {}, "$currentUser", currentUser).email)}</strong>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></p></div> `);
    if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "CLIENTE") {
      $$renderer2.push("<!--[0-->");
      $$renderer2.push(`<a href="/solicitacoes/novo" class="btn btn-white">+ Solicitar serviço</a>`);
    } else if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "ADMIN") {
      $$renderer2.push("<!--[1-->");
      $$renderer2.push(`<a href="/admin/profissionais" class="btn btn-white">Gerenciar profissionais</a>`);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="stats-grid svelte-x1i5gj"><!--[-->`);
    const each_array_1 = ensure_array_like(stats);
    for (let $$index = 0, $$length = each_array_1.length; $$index < $$length; $$index++) {
      let stat = each_array_1[$$index];
      $$renderer2.push(`<div class="stat-card card svelte-x1i5gj"><div class="row gap-3" style="justify-content: space-between; align-items: flex-start;"><span class="stat-card-label text-caption svelte-x1i5gj">${escape_html(stat.label)}</span> <span${attr_class(`badge badge-${stringify(stat.color)}`, "svelte-x1i5gj")}>${escape_html(stat.icon)}</span></div> <span class="stat-card-value svelte-x1i5gj">${escape_html(stat.value)}</span></div>`);
    }
    $$renderer2.push(`<!--]--></div> `);
    if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "CLIENTE") {
      $$renderer2.push("<!--[0-->");
      ClienteDashboard($$renderer2);
    } else if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "PROFISSIONAL") {
      $$renderer2.push("<!--[1-->");
      ProfissionalDashboard($$renderer2);
    } else if (store_get($$store_subs ??= {}, "$currentUser", currentUser)?.tipo === "ADMIN") {
      $$renderer2.push("<!--[2-->");
      AdminDashboard($$renderer2);
    } else {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div></div>`);
    if ($$store_subs) unsubscribe_stores($$store_subs);
  });
}
export {
  _page as default
};
