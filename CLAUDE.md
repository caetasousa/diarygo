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

**REGRA OBRIGATORIA:** Toda nova etapa do PLANO.md DEVE aplicar a skill `/owasp-security` durante implementacao.

Ao implementar autenticacao, endpoints HTTP, validacao de entrada, tratamento de erros ou revisao de seguranca, executar:
```bash
/owasp-security
```

A skill cobre todas as 10 categorias OWASP 2025 com exemplos em Go + chi + PostgreSQL:

| Categoria | Foco |
|-----------|------|
| **A01** | Broken Access Control (ownership, SSRF, CSRF) |
| **A02** | Security Misconfiguration (headers, config, Swagger) |
| **A03** | Software Supply Chain Failures (govulncheck, go.sum) |
| **A04** | Cryptographic Failures (bcrypt, JWT, TLS) |
| **A05** | Injection (SQL parametrizado, XSS, command injection) |
| **A06** | Insecure Design (rate limiting, threat modeling) |
| **A07** | Authentication Failures (timing attacks, NIST 800-63b) |
| **A08** | Data Integrity Failures (validacao, go.sum verify) |
| **A09** | Logging & Alerting Failures (slog, eventos criticos) |
| **A10** | Exceptional Conditions (fail closed, rollback, defer) |

Checklist minimo por etapa:
- Autenticacao/Autorizacao implementada? → usar skill secoes A01, A07
- Armazenar senhas? → bcrypt cost >= 12 (A04)
- Query ao banco? → pgx parametrizado $1, $2... (A05)
- Handler HTTP novo? → validacao completa + rate limiting (A06)
- Erro possivel? → fail closed, defer cleanup (A10)
- Logando dados? → nunca senhas/tokens/CPF (A09)
- Dependencias adicionadas? → govulncheck na CI (A03)

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

**REGRA OBRIGATORIA:** Sempre que terminar uma alteracao (backend ou frontend), rodar o projeto e exibir os links de acesso com todas as rotas disponíveis para facilitar testes manuais. Formato obrigatorio de saida:

```
Backend rodando em: http://localhost:8080
Swagger UI:         http://localhost:8080/swagger/index.html

API — Rotas disponiveis:
  GET  http://localhost:8080/health
  POST http://localhost:8080/api/v1/auth/registro/cliente
  POST http://localhost:8080/api/v1/auth/registro/profissional
  POST http://localhost:8080/api/v1/auth/login
  POST http://localhost:8080/api/v1/auth/solicitar-recuperacao-senha
  POST http://localhost:8080/api/v1/auth/redefinir-senha
  GET  http://localhost:8080/api/v1/me  (requer Bearer token)

Frontend rodando em: http://localhost:5173

Frontend — Paginas disponiveis:
  http://localhost:5173/                         Home
  http://localhost:5173/login                    Login
  http://localhost:5173/registro                 Registro de cliente
  http://localhost:5173/registro/profissional    Registro de diarista
  http://localhost:5173/dashboard                Dashboard (requer login)
  http://localhost:5173/recuperar-senha          Recuperar senha
  http://localhost:5173/redefinir-senha          Redefinir senha
```

Listar TODAS as rotas — incluindo novas rotas adicionadas na alteracao em destaque.

## Fase Atual: MVP

Implementacao em andamento conforme `PLANO.md` (13 etapas).
Etapas 0-11: logica completa com repositories in-memory, sem dependencia de banco.
Etapa 12: substituicao para PostgreSQL. Etapa 13: testes de integracao e Swagger.

NAO inclui no MVP: pagamento online, chat, indicacao, planos premium. Arquitetura preparada via interfaces.
