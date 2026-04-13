# DiaryGo — Sistema de Contratacao de Diarista

Plataforma intermediadora entre clientes e diaristas. MVP sem pagamento online.
Documento completo de regras de negocio e esquema do banco: `README.md`
Plano de implementacao por etapas: `PLANO.md`
Detalhes adicionais em `docs/` (estrutura, regras, banco, comandos).

## Manutencao do README.md

**REGRA OBRIGATORIA:** toda mudanca no sistema deve ser refletida no `README.md`.
- Nova tabela ou coluna no banco → atualizar secao "Esquema do Banco de Dados"
- Nova regra de negocio → atualizar secao correspondente e o checklist
- Nova funcionalidade → atualizar secao "Funcionalidades por Fase"
- Nova entidade de dominio → atualizar tabela de entidades
- Mudanca de stack ou estrutura de pastas → atualizar secao "Stack e Estrutura do Projeto"

## Stack

Go | chi (router HTTP) | PostgreSQL | Flyway | pgAdmin | Git | Swagger (swaggo/swag)

## Git Flow

- Branches permanentes: `master` (producao) e `developer` (integracao)
- Cada etapa do PLANO.md = uma branch `feature/<nome>` criada a partir de `developer`
- Nunca commitar direto em `master` ou `developer`
- Merge de feature na `developer` com `--no-ff`
- Hotfix criado a partir de `master`; merged em `master` E `developer`
- `git pull --rebase` sempre ao atualizar branches (nunca merge de atualizacao)
- Convencao de commits: `feat:`, `fix:`, `chore:`, `test:`, `docs:`

### Push para o GitHub

**REGRA OBRIGATORIA:** NUNCA executar `git push` sem antes perguntar ao usuario e aguardar confirmacao explicita. Isso vale para qualquer branch, em qualquer situacao, mesmo que o usuario ja tenha pedido push em sessoes anteriores. Cada push exige uma nova confirmacao.

Quando o usuario confirmar, executar conforme o contexto da branch atual:

```bash
# Push da branch atual (feature, hotfix, etc.) — primeiro push
git push --set-upstream origin <branch-atual>

# Push da branch atual — pushes subsequentes
git push

# Apos merge de feature na developer, push da developer
git checkout developer && git push

# Apos merge de release/hotfix na master, push da master
git checkout master && git push
```

**Regras para push:**
- Nunca fazer `git push --force` em `master` ou `developer`
- Em feature branches, `--force-with-lease` e permitido apos rebase
- Sempre confirmar com o usuario antes de fazer push em qualquer branch

## Estrategia de Repository

- Interfaces de repository definidas em `internal/domain/`
- Implementacoes in-memory em `internal/repository/memory/` (Etapas 0-11, sem banco)
- Implementacoes PostgreSQL em `internal/repository/postgres/` (Etapa 12+)
- `main.go` decide qual injetar — services e handlers nunca mudam ao trocar a implementacao
- Handlers HTTP usam chi; anotacoes swaggo em todos os handlers para geracao do Swagger

## Convencoes de Codigo

- Go idiomatico (effective Go, go vet, gofmt)
- Pacotes em minusculo, sem underscores
- Interfaces pequenas (1-3 metodos)
- Erros retornados explicitamente, nunca panic
- `context.Context` como primeiro parametro em funcoes de I/O
- Structs de dominio em `internal/domain/`, sem dependencia de frameworks

## Postura de Qualidade — Sem Preguica

**REGRA OBRIGATORIA:** Nunca deixar erros por corrigir. Sempre:
- Rodar `go test ./...` apos qualquer mudanca de codigo e corrigir **todos** os erros antes de continuar
- Rodar `go vet ./...` e `gofmt -w .` antes de cada commit
- Nao ignorar avisos do compilador, erros de IDE ou falhas de teste
- Se um teste falhar, investigar a causa raiz — nunca comentar ou deletar o teste
- Se um `go vet` apontar problema, corrigir o codigo — nunca suprimir sem justificativa
- Imports nao utilizados, variaveis declaradas e nao usadas, erros ignorados com `_` em I/O: **corrigir sempre**

## Seguranca (OWASP Top 10:2025)

**REGRA OBRIGATORIA:** Toda nova etapa do PLANO.md DEVE aplicar a skill `/owasp-security` durante implementacao. Isso vale tanto para o **backend (Go)** quanto para o **frontend (SvelteKit)**.

Ao implementar autenticacao, endpoints HTTP, validacao de entrada, tratamento de erros ou revisao de seguranca, executar:
```bash
/owasp-security
```

A skill cobre todas as 10 categorias OWASP 2025 com exemplos em Go + chi + PostgreSQL. Para o frontend, aplicar os mesmos principios adaptados:

| Categoria | Backend (Go) | Frontend (SvelteKit) |
|-----------|-------------|----------------------|
| **A01** | Ownership checks, SSRF, CSRF | Guard de rotas, nao expor dados de outros usuarios |
| **A02** | Security headers (chi middleware) | CSP no app.html, headers meta, `skipLibCheck` nao suprime erros reais |
| **A03** | govulncheck, go.sum | npm audit, nao usar `--force` em audit fix sem avaliar |
| **A04** | bcrypt cost >= 12, JWT HS256 | Token em localStorage (tradeoff SPA), nunca logar token no console |
| **A05** | pgx parametrizado $1, $2... | Sem innerHTML dinamico, sem eval(), sem interpolacao direta no DOM |
| **A06** | Rate limiting httprate | maxlength em todos inputs, validacao no cliente E no servidor |
| **A07** | Timing attack bcrypt, NIST 800-63b | autocomplete correto, mensagens de erro genericas, nao revelar campo errado |
| **A08** | Validacao entrada, go.sum verify | Types TypeScript espelham DTOs do backend, sem any silencioso |
| **A09** | slog estruturado, nunca logar CPF/senha | Nunca console.log com token, email ou senha |
| **A10** | fail closed, defer cleanup | Tratamento de erro em todo fetch, nao expor stack trace ao usuario |

Checklist minimo por etapa — **backend**:
- Autenticacao/Autorizacao implementada? → usar skill secoes A01, A07
- Armazenar senhas? → bcrypt cost >= 12 (A04)
- Query ao banco? → pgx parametrizado $1, $2... (A05)
- Handler HTTP novo? → validacao completa + rate limiting (A06)
- Erro possivel? → fail closed, defer cleanup (A10)
- Logando dados? → nunca senhas/tokens/CPF (A09)
- Dependencias adicionadas? → govulncheck na CI (A03)

Checklist minimo por etapa — **frontend**:
- Formulario novo? → maxlength em todos inputs, validacao client-side + server-side (A06)
- Rota protegida? → +page.ts com guard `isAuthenticated` + redirect 302 (A01)
- Dados do usuario exibidos? → escapar via Svelte (nao usar `@html` sem sanitizar) (A05)
- Fetch novo? → try/catch obrigatorio, mensagem generica ao usuario, nao expor erro raw (A10)
- Dependencias adicionadas? → `npm audit` e avaliar severidade (A03)
- app.html atualizado? → CSP, X-Content-Type-Options, Referrer-Policy (A02)
- Console.log? → nunca com token, senha ou dados pessoais (A09)

## Regras Criticas

1. **LC 150/2015:** max 2 visitas/semana mesma profissional no mesmo endereco — bloquear automaticamente
2. Agendamento minimo 24h de antecedencia
3. Cancelamento < 24h = penalizacao no score
4. Nota < 3.5 por 3 servicos consecutivos = suspensao
5. Atribuicao por avaliacao + proximidade geografica
6. Profissional tem 30 min para aceitar; apos, redireciona

## Testes

- Testes unitarios dos services: usam repositories in-memory, sem banco, sem Docker
- Testes de integracao dos repositories postgres: usam banco real (Etapa 13+)
- Nunca usar mocks para PostgreSQL — usar implementacao in-memory ou banco real
- Nomenclatura: `TestMetodo_Cenario_Resultado`

## Comandos Rapidos

```bash
# Backend
cd backend
go run cmd/api/main.go                    # rodar (sem banco nas Etapas 0-11)
go test ./...                             # testar
gofmt -w .                                # formatar
go vet ./...                              # verificar
swag init -g cmd/api/main.go -o ../docs/swagger  # gerar Swagger
docker-compose up -d                      # infra (Etapa 12+)

# Frontend
cd frontend
npm run dev                               # rodar em modo desenvolvimento
npm run build                             # build de producao
npm run preview                           # preview do build
```

## Rodar o Projeto Apos Alteracoes

**REGRA OBRIGATORIA:** Sempre que terminar uma alteracao (backend ou frontend), rodar o projeto e exibir os links de acesso com todas as rotas disponíveis para facilitar testes manuais. Os links devem ser renderizados como markdown clicavel (formato `[texto](url)`) para que o usuario possa clicar diretamente no terminal/IDE e abrir no navegador. Formato obrigatorio de saida:

**Backend rodando em:** [http://localhost:8080](http://localhost:8080)
**Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

**API — Rotas disponiveis:**
- `GET`  [http://localhost:8080/health](http://localhost:8080/health)
- `POST` http://localhost:8080/api/v1/auth/registro/cliente
- `POST` http://localhost:8080/api/v1/auth/registro/profissional
- `POST` http://localhost:8080/api/v1/auth/login
- `POST` http://localhost:8080/api/v1/auth/solicitar-recuperacao-senha
- `POST` http://localhost:8080/api/v1/auth/redefinir-senha
- `GET`  http://localhost:8080/api/v1/me _(requer Bearer token)_

**Frontend rodando em:** [http://localhost:5173](http://localhost:5173)

**Frontend — Paginas disponiveis:**
- [http://localhost:5173/](http://localhost:5173/) — Home
- [http://localhost:5173/login](http://localhost:5173/login) — Login
- [http://localhost:5173/registro](http://localhost:5173/registro) — Registro de cliente
- [http://localhost:5173/registro/profissional](http://localhost:5173/registro/profissional) — Registro de diarista
- [http://localhost:5173/dashboard](http://localhost:5173/dashboard) — Dashboard _(requer login)_
- [http://localhost:5173/recuperar-senha](http://localhost:5173/recuperar-senha) — Recuperar senha
- [http://localhost:5173/redefinir-senha](http://localhost:5173/redefinir-senha) — Redefinir senha

Listar TODAS as rotas — incluindo novas rotas adicionadas na alteracao em destaque.

## Revisao de Codigo — Codex + Engenheiro Senior

**REGRA OBRIGATORIA:** Todo codigo novo ou alterado deve ser revisado sob dois angulos antes de ser considerado pronto:

### 1. Revisao Codex (corretude e padroes)
Simular a perspectiva de um revisor automatizado rigoroso. Verificar:
- Codigo compila e todos os testes passam (`go test ./...` / `npm run check`)
- Sem imports nao usados, variaveis mortas, erros ignorados com `_`
- Convencoes de nomenclatura respeitadas (Go: camelCase/PascalCase; TS: camelCase)
- Nenhuma logica duplicada que deveria ser abstraida
- Tipos TypeScript corretos — sem `any` implicito ou cast forcado
- Nenhum `console.log` ou `fmt.Println` de debug esquecido no codigo

### 2. Revisao Senior — Backend (Go + chi + PostgreSQL)
Simular a perspectiva de um engenheiro Go senior. Verificar:
- Interfaces sao pequenas e focadas (1-3 metodos) — Go idiomatico
- Injecao de dependencia via construtor, sem globals mutaveis
- Erros tratados explicitamente — nunca `_` em operacoes de I/O
- `context.Context` propagado corretamente em toda cadeia de chamadas
- Handlers HTTP: MaxBytesReader, decode, validar, chamar service, mapear erro
- Nenhuma logica de negocio nos handlers — pertence ao service
- SQL (quando aplicavel): parametrizado, sem concatenacao de string
- Seguranca: headers aplicados, rate limiting ativo, JWT validado com metodo fixo

### 2. Revisao Senior — Frontend (SvelteKit + Svelte 5 + TypeScript)
Simular a perspectiva de um engenheiro frontend senior. Verificar:
- Svelte 5: snippets via `{#snippet}` + `{@render}`, layout via `{@render children()}`
- Stores usados corretamente — sem subscribe manual sem unsubscribe
- Formularios: `novalidate` + validacao JS, `autocomplete` correto, `maxlength` em todos inputs
- Navegacao apos submit: `window.location.href` para troca de pagina completa
- Guards de rota em `+page.ts` com `redirect(302, '/login')` — nao em `onMount`
- Sem `@html` com dados do usuario sem sanitizacao
- Erros de API: try/catch em todo fetch, mensagem amigavel ao usuario, nao expor stack trace
- `npm run check` com 0 erros antes de qualquer commit

### 2. Revisao Senior — Infra e DevOps
Simular a perspectiva de um engenheiro de infraestrutura senior. Verificar:
- `.gitignore` cobre todos artefatos gerados (`node_modules`, `.svelte-kit`, `build`, `bin/`, `dist/`)
- Secrets nunca commitados — `.env`, `app.env`, chaves JWT fora do repositorio
- Docker Compose (quando ativo): volumes nomeados, healthcheck definido, portas nao expostas desnecessariamente
- Variaveis de ambiente validadas na inicializacao da aplicacao (`config.Carregar()`)
- Ambientes separados: `development` vs `production` — rate limits, logs e Swagger diferenciados
- Swagger exposto apenas em non-production
- Frontend exposto em `0.0.0.0` apenas em dev — em producao usar reverse proxy (nginx/caddy)
- Node.js e Go em versoes LTS documentadas (`nvm use 20`, `go 1.25`)
- Comandos de build reproduziveis: `go build ./...` e `npm run build` sem erros

## Fase Atual: MVP

Implementacao em andamento conforme `PLANO.md` (13 etapas).
Etapas 0-11: logica completa com repositories in-memory, sem dependencia de banco.
Etapa 12: substituicao para PostgreSQL. Etapa 13: testes de integracao e Swagger.

NAO inclui no MVP: pagamento online, chat, indicacao, planos premium. Arquitetura preparada via interfaces.
