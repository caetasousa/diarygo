# Banco de Dados

Documentação da estratégia de persistência, migrations e operação local.

---

## Estratégia: In-Memory First + Migrations Incrementais

Durante as Etapas 0–11 o backend roda com **repositories in-memory** (`internal/repository/memory/`). Ou seja: `go run cmd/api/main.go` e `go test ./...` **não exigem PostgreSQL nem Docker** no desenvolvimento do dia-a-dia.

As **migrations Flyway**, porém, já crescem etapa por etapa desde a Etapa 1. Cada etapa que introduz, altera ou remove tabelas cria **uma** migration em `backend/migrations/V{N}__{nome_etapa_snake}.sql`.

Os **testes de integração** de repositories e handlers usam banco real (PostgreSQL efêmero via testcontainers-go) e rodam sob a build tag `//go:build integration`. Ver [backend/internal/testutil/](../backend/internal/testutil/).

Na **Etapa 12** as implementações in-memory são substituídas pelos repositories PostgreSQL — services e handlers não mudam.

```
┌───────────────────┐
│  handler (chi)    │  ← não muda nunca
├───────────────────┤
│  service          │  ← não muda nunca (primeira linha de defesa)
├───────────────────┤
│  domain.Repo      │  ← interface, não muda
├───────────────────┤
│  memory.Repo      │  repository/memory/    (Etapas 0-11, unit tests)
│        ou         │
│  postgres.Repo    │  repository/postgres/  (Etapa 12+, integration tests)
└───────────────────┘
```

---

## Princípio de validação — defesa em profundidade

O banco é a **última** linha de defesa; o service é a **primeira**. Toda constraint (UNIQUE, CHECK, ENUM, FK) tem contraparte no service em `internal/service/`. O CHECK do banco só dispara se o service for contornado (bug, migration futura, acesso direto ao Postgres).

Cada cabeçalho `V*.sql` enumera o espelhamento `service ↔ banco`. Ver também a seção "Histórico de Migrations" em [README.md](../README.md#17-esquema-do-banco-de-dados).

---

## PostgreSQL

- **Versão:** PostgreSQL 14 (alinhado ao testcontainers e ao docker-compose)
- **Migrations:** versionadas em [backend/migrations/](../backend/migrations/), aplicadas pelo Flyway
- **Admin UI:** pgAdmin em [http://localhost:5050](http://localhost:5050) para inspeção visual
- **Schema completo:** documentado em [README.md §17](../README.md#17-esquema-do-banco-de-dados)

---

## Regra: uma migration por etapa

**REGRA OBRIGATÓRIA** (ver [CLAUDE.md](../CLAUDE.md) seção "Migrations"):

- Formato do nome: `V{N}__{nome_etapa_snake}.sql` (ex: `V3__catalogo_precificacao.sql`)
- **Nunca editar** migration já aplicada em `developer`/`master`. Correções viram nova migration (`V{N+1}__fix_xxx.sql`)
- Ordem canônica dentro do arquivo: **extensões → enums → tabelas** (FKs resolvidas) **→ índices → triggers → views**
- Cabeçalho obrigatório no topo do arquivo SQL:
  - Etapa associada (`Etapa: N — Nome`)
  - Resumo em 1–2 linhas
  - Delta vs. migrations anteriores (tabelas novas/alteradas, enums, índices, triggers)
  - Rollback equivalente (bloco comentado) para reset em dev — em produção **nunca** se roda rollback
  - Espelhamento service ↔ banco (defesa em profundidade)

Ao abrir PR de feature, atualizar em conjunto:

1. Migration SQL (`backend/migrations/V{N}__*.sql`)
2. `README.md` seções "Esquema do Banco" e "Histórico de Migrations"
3. Domain Go (`backend/internal/domain/`)
4. Services Go (`backend/internal/service/`) — validadores espelhando CHECKs
5. Factories (`backend/internal/testutil/factories/`)
6. `AllTables` em `backend/internal/testutil/truncate.go` (se houver tabela nova)
7. Repository Postgres (`backend/internal/repository/postgres/`) + teste de integração

O schema planejado das Etapas 3–11 (ainda não migrado) vive em [docs/schema-futuro.sql](schema-futuro.sql) como referência — **não** executado pelo Flyway. Ao iniciar cada etapa, extrair o bloco correspondente e criar a V{N} definitiva.

---

## Migrations aplicadas

Ver [README.md — Histórico de Migrations](../README.md#17-esquema-do-banco-de-dados) para o detalhamento completo (etapa, data, tabelas novas/alteradas, enums, índices, espelhamento service↔banco).

| Versão | Arquivo | Etapa |
|---|---|---|
| V1 | [`V1__auth_usuarios.sql`](../backend/migrations/V1__auth_usuarios.sql) | Etapa 1 — Autenticação e Usuários |
| V2 | [`V2__cadastro_clientes_profissionais.sql`](../backend/migrations/V2__cadastro_clientes_profissionais.sql) | Etapa 2 — Cadastro de Clientes e Profissionais |

---

## Subir o banco local

```bash
docker compose up -d
```

Isso sobe três serviços:

| Serviço | Porta host | Descrição |
|---|---|---|
| `postgres` | `5432` | PostgreSQL 14-alpine (db/user/senha: `diarygo`) |
| `pgadmin` | `5050` | pgAdmin 4 para inspeção visual |
| `flyway` | — | Aplica `V1`, `V2`, … automaticamente ao subir |

O Flyway monta [backend/migrations/](../backend/migrations/) e aplica todas as `V{N}__*.sql` em ordem.

### Conectar ao Postgres

**Do host (psql local ou DBeaver):**

| Campo | Valor |
|---|---|
| Host | `localhost` |
| Port | `5432` |
| Database | `diarygo` |
| User | `diarygo` |
| Password | `diarygo` |

**Via container (mais rápido para SQL ad-hoc):**

```bash
docker compose exec postgres psql -U diarygo -d diarygo
```

**Via pgAdmin (GUI):** [http://localhost:5050](http://localhost:5050), credenciais `admin@diarygo.com.br / admin`. Ao adicionar servidor, use `postgres` como hostname (não `localhost` — é a rede interna do Docker).

---

## Reset em dev

Como **não há ambiente de produção** ainda, reset completo é aceitável:

```bash
docker compose down -v      # para tudo e zera o volume
docker compose up -d        # sobe limpo e reaplica V1, V2, ...
```

Nunca rode `down -v` em ambiente compartilhado. Em produção (futura), correções vêm de novas migrations `V{N+1}__fix_*.sql`, não de rollback.

---

## Testes de integração

```bash
cd backend
go test -tags=integration ./internal/testutil/...           # smoke test da infra
go test -tags=integration ./internal/repository/postgres/... # só repositories
go test -tags=integration ./...                             # suite completa
```

Detalhes e reutilização de container: ver [docs/comandos.md](comandos.md#integra%C3%A7%C3%A3o-requer-docker).

---

## Convenções do schema

- **IDs:** `UUID PRIMARY KEY DEFAULT uuid_generate_v4()` (extensão `uuid-ossp` habilitada na V1)
- **Timestamps:** `criado_em` e `atualizado_em` como `TIMESTAMPTZ NOT NULL DEFAULT NOW()`; `atualizado_em` mantido por trigger `fn_set_atualizado_em()` definida na V1
- **Enums:** criados com `CREATE TYPE ... AS ENUM (...)` e reutilizados em várias tabelas
- **UNIQUE em colunas únicas** (email, CPF) — não criar índice adicional porque o UNIQUE já cria B-tree automaticamente
- **FK com índice explícito:** Postgres **não** cria índice em coluna FK automaticamente. Toda coluna `*_id` referenciando outra tabela ganha `CREATE INDEX idx_<tabela>_<ref>`
- **Índice parcial** para colunas raramente preenchidas (ex: `token_recuperacao`) — mantém custo de escrita/vacuum próximo de zero
- **CHECKs de formato** (CPF `^[0-9]{11}$`, CEP `^[0-9]{8}$`, dia_semana 0..6) espelhando validadores do domain
- **CASCADE em FK** quando o filho não faz sentido sem o pai (endereços do cliente, documentos da profissional)
