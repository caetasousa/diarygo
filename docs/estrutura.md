# Estrutura do Projeto

Árvore anotada com responsabilidade de cada diretório. Para convenções de código, ver [CLAUDE.md](../CLAUDE.md).

---

## Visão geral (raiz)

```
diarygo/
├── backend/              # API Go (chi + Postgres/in-memory + JWT)
├── frontend/             # SPA SvelteKit 5 + TypeScript
├── docs/                 # Documentação (este diretório)
├── docker-compose.yml    # Postgres + pgAdmin + Flyway para dev local
├── CLAUDE.md             # Instruções e convenções obrigatórias
├── PLANO.md              # Plano de implementação por etapa (0-13)
└── README.md             # Regras de negócio, esquema, rotas
```

---

## Backend

```
backend/
├── cmd/
│   └── api/
│       └── main.go               # Entrypoint HTTP — wiring de DI, middlewares, rotas
├── internal/
│   ├── config/
│   │   └── config.go             # Carrega env vars (PORT, JWT_SECRET, ENV, ...)
│   ├── domain/                   # Entidades + interfaces de repo + VALIDADORES
│   │   ├── usuario.go            # domain.Usuario + TipoUsuario + erros
│   │   ├── cliente.go            # domain.Cliente + ValidarCPF, ValidarScore, ...
│   │   ├── profissional.go       # domain.Profissional + ValidarNotaMedia, ...
│   │   ├── endereco.go           # domain.Endereco + ValidarCEP, ValidarComodos, ...
│   │   ├── documento.go          # tipos/status de documento
│   │   ├── referencia.go         # status de referência
│   │   ├── regiao.go             # Regiao + interface ProfissionalRegiaoRepository
│   │   ├── disponibilidade.go    # slots semanais + ValidarHora
│   │   └── auth.go               # JWT claims, requests de login/registro
│   ├── service/                  # Use cases (PRIMEIRA linha de defesa)
│   │   ├── auth_service.go       # Login, registro, bcrypt, JWT, recuperação de senha
│   │   ├── cliente_service.go    # Criar/atualizar cliente
│   │   ├── profissional_service.go # Criar/atualizar profissional
│   │   ├── endereco_service.go   # CRUD de endereços + principal
│   │   └── credenciamento_service.go # Documentos, referências, regiões, disponibilidades
│   ├── handler/                  # Handlers HTTP (chi) — mapeia req→service→resp
│   │   ├── auth_handler.go
│   │   ├── cliente_handler.go
│   │   ├── profissional_handler.go
│   │   ├── endereco_handler.go
│   │   ├── regiao_handler.go
│   │   └── resposta.go           # RespostaJSON, RespostaErro, MapearErro
│   ├── repository/
│   │   ├── memory/               # Implementações in-memory (Etapas 0-11, unit tests)
│   │   └── postgres/             # Implementações pgx (Etapa 12+, integration tests)
│   │       ├── db.go             # NewPool — pgxpool configurado
│   │       ├── errors.go         # pgErrorCode, pgErrorConstraint (SQLSTATE helpers)
│   │       ├── usuario.go        # + demais entidades
│   │       └── *_test.go         # build tag `integration`
│   ├── middleware/               # Autenticar (JWT), RequererTipo, securityHeaders
│   │   └── auth.go
│   └── testutil/                 # Infraestrutura de testes de integração
│       ├── container.go          # StartPostgres via testcontainers-go
│       ├── migrate.go            # ApplyMigrations via golang-migrate
│       ├── truncate.go           # Truncate(AllTables) entre testes
│       └── factories/            # Builders com defaults válidos (NewCliente, NewEndereco, ...)
├── migrations/                   # Flyway SQL — consumido TAMBÉM pelos integration tests
│   ├── V1__auth_usuarios.sql
│   └── V2__cadastro_clientes_profissionais.sql
├── config/
│   └── app.env                   # Config local (gitignored)
├── go.mod
└── go.sum
```

### Regras de dependência

```
cmd/api → handler → service → domain.Repository ← repository/{memory,postgres}
                                   ↑
                                interface
```

- `domain/` **não** importa nada dos outros pacotes internos
- `service/` só depende de `domain/`
- `handler/` depende de `service/` e `domain/`
- `repository/*` implementa interfaces de `domain/`
- `cmd/api/main.go` é o único lugar que conhece todos eles — faz o wiring

---

## Frontend

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api/                  # Cliente HTTP tipado (fetch wrapper)
│   │   │   └── client.ts         # apiGet, apiPost com injeção de Bearer token
│   │   ├── components/           # Componentes Svelte compartilhados
│   │   │   ├── Navbar.svelte
│   │   │   ├── Footer.svelte
│   │   │   ├── Toast.svelte
│   │   │   └── ToastContainer.svelte
│   │   ├── stores/               # Estado global
│   │   │   ├── auth.ts           # Token JWT + payload + persistência localStorage
│   │   │   └── toasts.ts         # Fila de notificações
│   │   └── types/                # DTOs que espelham o backend
│   ├── routes/                   # Páginas SvelteKit (file-based routing)
│   │   ├── +layout.svelte        # Layout raiz (Navbar + ToastContainer)
│   │   ├── +page.svelte          # Home (landing)
│   │   ├── login/+page.svelte
│   │   ├── registro/+page.svelte
│   │   ├── registro/profissional/+page.svelte
│   │   ├── recuperar-senha/+page.svelte
│   │   ├── redefinir-senha/
│   │   │   ├── +page.svelte
│   │   │   └── +page.ts          # Guard: só entra se tiver token
│   │   └── dashboard/
│   │       ├── +layout.svelte    # Navbar lateral do dashboard
│   │       ├── +layout.ts        # Guard: exige JWT (redirect 302 /login)
│   │       ├── +page.svelte      # Dashboard adaptado ao tipo de usuário
│   │       ├── perfil/+page.svelte
│   │       ├── enderecos/+page.svelte
│   │       ├── documentos/+page.svelte
│   │       ├── referencias/+page.svelte
│   │       ├── regioes/+page.svelte
│   │       └── disponibilidade/+page.svelte
│   ├── app.html                  # Template HTML raiz (CSP, meta headers)
│   └── app.css                   # Design system global (Resend-inspired)
├── static/                       # Assets estáticos (favicon, imagens)
├── svelte.config.js
├── vite.config.ts
├── tsconfig.json
└── package.json
```

### Guards de rota

- `+page.ts` / `+layout.ts` com `redirect(302, '/login')` quando `!isAuthenticated`
- **Nunca** usar `onMount` para guard — atrasa e mostra conteúdo protegido brevemente

---

## Docs

```
docs/
├── comandos.md           # Todos os comandos Go, npm, docker, flyway, testes
├── banco-dados.md        # Estratégia de persistência, migrations, pgAdmin, reset
├── estrutura.md          # Este arquivo
├── regras-negocio.md     # Regras críticas (LC 150/2015, score, avaliação)
├── schema-futuro.sql     # DDL das Etapas 3-11 (referência, NÃO aplicado)
└── swagger/              # Gerado por `swag init` — NÃO editar manualmente
    ├── docs.go
    ├── swagger.json
    └── swagger.yaml
```

---

## Arquivos-raiz importantes

| Arquivo | O que contém |
|---|---|
| [README.md](../README.md) | Regras de negócio, esquema do banco, rotas da API, como rodar |
| [CLAUDE.md](../CLAUDE.md) | Convenções obrigatórias: Git flow, qualidade, segurança OWASP, migrations |
| [PLANO.md](../PLANO.md) | Plano de implementação por etapa (0-13) com status |
| [docker-compose.yml](../docker-compose.yml) | Postgres + pgAdmin + Flyway para dev local |
| `.gitignore` | Cobre `node_modules/`, `.svelte-kit/`, `build/`, `bin/`, `app.env`, test-results/ |
