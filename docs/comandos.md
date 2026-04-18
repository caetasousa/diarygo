# Comandos Úteis

Referência de todos os comandos do dia-a-dia. Para rodar o projeto pela primeira vez, veja o [Quickstart no README](../README.md#-quickstart-3-comandos).

---

## Índice

- [Backend — Go](#backend--go)
- [Frontend — SvelteKit](#frontend--sveltekit)
- [Testes](#testes)
- [Banco de Dados](#banco-de-dados)
- [Swagger](#swagger)
- [Git](#git)
- [Troubleshooting](#troubleshooting)

---

## Backend — Go

```bash
cd backend

# Rodar a API
go run cmd/api/main.go

# Compilar para binário
go build -o bin/api cmd/api/main.go

# Qualidade de código (antes de commitar)
gofmt -w .
go vet ./...
go test ./...

# Análise de vulnerabilidades (OWASP A03)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Dependências
go mod tidy        # limpar imports não usados
go mod download    # baixar todas as deps
go mod verify      # checar go.sum
```

**Variáveis de ambiente (opcional):** crie `backend/config/app.env` ou exporte no shell. Defaults seguros em dev.

| Variável | Default dev | Obrigatório em produção |
|---|---|---|
| `PORT` | `8080` | não |
| `ENV` | `development` | sim (`production`) |
| `JWT_SECRET` | padding inseguro + warn | **sim** (≥32 chars) |
| `JWT_EXPIRATION_MINUTES` | `15` | não |

---

## Frontend — SvelteKit

```bash
cd frontend

# Primeira vez
npm install

# Dev (hot reload)
npm run dev

# Build de produção
npm run build
npm run preview        # servir o build local

# Qualidade
npm run check          # svelte-check (TypeScript)
npm run lint           # eslint
npm run format         # prettier
```

---

## Testes

### Unit (rápido, sem Docker)

```bash
cd backend
go test ./...                     # todos os pacotes
go test -v ./internal/service     # pacote específico, verbose
go test -cover ./...              # com cobertura
go test -run TestLogin ./...      # filtrar por nome
```

### Integração (requer Docker)

Sobe um PostgreSQL efêmero via testcontainers-go, aplica as migrations de `backend/migrations/` e roda. Um container por pacote; `TRUNCATE` entre casos.

```bash
cd backend

# Suite completa de integração
go test -tags=integration ./...

# Só smoke test da infra (sobe container, aplica V1+V2, trunca, baixa)
go test -tags=integration ./internal/testutil/...

# Só os repositories Postgres
go test -tags=integration ./internal/repository/postgres/...

# Reutilizar container entre runs (salta ~3s de startup)
export TESTCONTAINERS_REUSE_ENABLE=true
go test -tags=integration ./...
```

**Requisito:** Docker em execução. Cobre migrations `backend/migrations/V{N}__*.sql` aplicadas via `golang-migrate` (mesmos arquivos que o Flyway consome em produção).

---

## Banco de Dados

### Stack local (Postgres + pgAdmin + Flyway)

```bash
# Subir (aplica V1, V2 automaticamente)
docker compose up -d

# Ver logs do Flyway (confirmar que migrations aplicaram)
docker compose logs flyway

# Parar (mantém volume)
docker compose down

# Parar E zerar dados (reset total — reaplicará V1, V2 na próxima subida)
docker compose down -v
```

### Inspecionar schema

```bash
# Listar tabelas
docker compose exec postgres psql -U diarygo -d diarygo -c '\dt'

# Descrever tabela
docker compose exec postgres psql -U diarygo -d diarygo -c '\d usuarios'

# SQL interativo
docker compose exec postgres psql -U diarygo -d diarygo

# Contar linhas por tabela
docker compose exec postgres psql -U diarygo -d diarygo -c "
  SELECT schemaname, relname AS tabela, n_live_tup AS linhas
  FROM pg_stat_user_tables
  ORDER BY n_live_tup DESC;"
```

### pgAdmin

Abra [http://localhost:5050](http://localhost:5050) → login `admin@diarygo.com.br` / `admin` → **Add New Server**:

| Campo | Valor |
|---|---|
| General → Name | `DiaryGo Local` |
| Connection → Host name/address | `postgres` |
| Connection → Port | `5432` |
| Connection → Maintenance DB | `diarygo` |
| Connection → Username | `diarygo` |
| Connection → Password | `diarygo` |

> Use `postgres` como host (não `localhost`) porque o pgAdmin roda dentro do Docker e resolve pelo nome do serviço.

### Flyway manual (se precisar forçar migrate)

```bash
# Re-rodar Flyway (normalmente é automático)
docker compose run --rm flyway migrate

# Ver histórico de migrations aplicadas
docker compose run --rm flyway info
```

---

## Swagger

```bash
# Instalar swag (uma vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar documentação a partir das anotações dos handlers
cd backend
swag init -g cmd/api/main.go -o ../docs/swagger
```

Após gerar, recarregue a API (`go run cmd/api/main.go`) e acesse [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

---

## Git

```bash
# Fluxo padrão (convenção do projeto)
git checkout developer
git pull --rebase origin developer
git checkout -b feature/nome-da-etapa

# ... trabalho ...

git add <arquivos-específicos>     # nunca `git add .`
git commit -m "feat: descrição curta"
git push --set-upstream origin feature/nome-da-etapa

# Merge de feature na developer (com --no-ff preserva histórico da branch)
git checkout developer
git merge --no-ff feature/nome-da-etapa
```

Convenção de commits: `feat:` `fix:` `chore:` `test:` `docs:` (ver [CLAUDE.md](../CLAUDE.md)).

---

## Troubleshooting

### Backend não sobe

```bash
# Porta 8080 ocupada
lsof -i:8080
# Matar o processo ou rodar em outra porta:
PORT=8081 go run cmd/api/main.go

# go.sum corrompido
go clean -modcache
go mod download

# Erro de build após mudança em domain/
go build ./...     # ver mensagem completa
```

### Frontend não sobe

```bash
# Versão errada do Node
node -v     # precisa ser 20+
nvm install 20 && nvm use 20

# node_modules corrompido
rm -rf node_modules package-lock.json
npm install

# Porta 5173 ocupada
npm run dev -- --port 5174
```

### Docker / Postgres

```bash
# Flyway falha com "schema corrupted" ou "migration out of order"
docker compose down -v        # zera o volume
docker compose up -d          # sobe tudo limpo

# Container "diarygo_postgres" já existe
docker compose down
docker rm diarygo_postgres

# pgAdmin não conecta (hostname error)
# Use "postgres" como hostname, não "localhost" — é a rede interna do Docker

# Testcontainers "pulling image" demora na primeira vez
docker pull postgres:14-alpine    # adiantar o pull

# Testcontainers cria container novo a cada run mesmo com TESTCONTAINERS_REUSE_ENABLE=true
# Docker Desktop precisa estar rodando (não rootless) — Ryuk precisa de permissão
```

### Testes

```bash
# Teste de integração trava
docker ps                           # container órfão?
docker rm -f $(docker ps -aq --filter "label=org.testcontainers=true")

# go test falha com "cannot find package"
go mod tidy

# Teste passa local mas falha em CI
go test -race ./...                 # detectar race conditions
```

---

## Referência rápida

| Quero… | Comando |
|---|---|
| Rodar API | `cd backend && go run cmd/api/main.go` |
| Rodar frontend | `cd frontend && npm run dev` |
| Rodar testes unit | `cd backend && go test ./...` |
| Rodar testes de integração | `cd backend && go test -tags=integration ./...` |
| Subir banco local | `docker compose up -d` |
| Ver tabelas aplicadas | `docker compose exec postgres psql -U diarygo -d diarygo -c '\dt'` |
| Zerar banco | `docker compose down -v && docker compose up -d` |
| Formatar tudo (backend) | `cd backend && gofmt -w .` |
| Gerar Swagger | `cd backend && swag init -g cmd/api/main.go -o ../docs/swagger` |
| Checar vulnerabilidades | `cd backend && govulncheck ./...` |
