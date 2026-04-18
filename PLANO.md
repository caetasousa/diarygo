# DiaryGo — Plano de Implementacao

> Referencia: `README.md` (regras de negocio) . `migrations/` (SQL do banco)
> Cada etapa deve ser concluida e testada antes de iniciar a proxima.
> **Cada etapa inclui backend E frontend** — o frontend acompanha o backend desde a primeira funcionalidade.

## Git Flow

O projeto segue o modelo Gitflow. Branches permanentes:

| Branch | Papel |
|---|---|
| `master` | Historico oficial de entregas (producao). Nunca recebe commits diretos. |
| `developer` | Integracao continua. Base para todas as features. |

### Branches por tipo

```
feature/<nome>   -> criada a partir de developer, merge de volta para developer
release/X.Y.Z    -> criada a partir de developer quando um conjunto de features esta pronto; merge para master e developer
hotfix/<nome>    -> criada a partir de master para correcoes urgentes em producao; merge para master e developer
```

### Mapeamento Etapa -> Feature Branch

Cada etapa do plano corresponde a uma branch feature:

| Etapa | Branch | Status |
|---|---|---|
| 0 | `feature/estrutura-base` | Concluida |
| 1 | `feature/autenticacao` | Concluida |
| 2 | `feature/cadastro-clientes-profissionais` | Pendente |
| 3 | `feature/catalogo-precificacao` | Pendente |
| 4 | `feature/solicitacoes` | Pendente |
| 5 | `feature/matching-atribuicao` | Pendente |
| 6 | `feature/execucao-checkin-checkout` | Pendente |
| 7 | `feature/avaliacoes-reputacao` | Pendente |
| 8 | `feature/recorrencias` | Pendente |
| 9 | `feature/notificacoes` | Pendente |
| 10 | `feature/painel-admin` | Pendente |
| 11 | `feature/historico-transacoes` | Pendente |
| 12 | `feature/persistencia-postgres` | Pendente |
| 13 | `feature/testes-integracao-swagger` | Pendente |

Quando todas as features de uma release estiverem na `developer`, abrir `release/X.Y.Z` para preparar a entrega.

### Fluxo resumido

```bash
# Iniciar feature
git checkout developer
git pull --rebase
git checkout -b feature/<nome>

# Commitar durante o desenvolvimento
git add <arquivos>
git commit -m "feat: descricao sucinta"
git push --set-upstream origin feature/<nome>   # primeiro push
git push                                         # pushes seguintes

# Manter feature atualizada com developer
git checkout developer && git pull --rebase
git checkout feature/<nome>
git rebase developer
git push --force-with-lease

# Finalizar feature (merge na developer)
git checkout developer
git merge --no-ff feature/<nome>
git push

# Hotfix (a partir da master)
git checkout master && git pull --rebase
git checkout -b hotfix/<nome>
# ... corrige e commita ...
git checkout master  && git merge --no-ff hotfix/<nome> && git push
git checkout developer && git merge --no-ff hotfix/<nome> && git push
```

### Convencao de commits

```
feat: adiciona endpoint de login
fix: corrige calculo de horas na limpeza express
chore: atualiza dependencias
test: adiciona testes de matching com limite LC 150/2015
docs: atualiza README com esquema do banco
```

---

## Estrategia: In-Memory First

Toda a logica de negocio e desenvolvida com repositories **in-memory** (Etapas 0-11).
O banco de dados e introduzido apenas na **Etapa 12**, substituindo as implementacoes in-memory por Postgres — sem tocar nos services ou handlers.

```
internal/domain/          <- interfaces de repository (contrato)
internal/repository/
  memory/                 <- implementacoes in-memory (Etapas 0-11)
  postgres/               <- implementacoes Postgres  (Etapa 12)
```

O `main.go` decide qual injetar:
```go
// Etapas 0-11
repo := memory.NewClienteRepository()

// Etapa 12 (troca uma linha)
repo := postgres.NewClienteRepository(db)
```

Vantagens: sem Docker obrigatorio para desenvolver, testes unitarios triviais, services nunca mudam ao trocar o banco.

---

## Estrategia: Frontend Integrado por Etapa

O frontend acompanha o backend desde a Etapa 0. Cada etapa entrega a funcionalidade completa — da API ate a interface do usuario.

### Padroes Estabelecidos (referencia de excelencia)

O codigo das Etapas 0-1 define os padroes que todas as etapas seguintes devem respeitar:

**Componentes e Paginas:**
- Svelte 5 com `{#snippet}` + `{@render}` para composicao
- Layout via `{@render children()}` no `+layout.svelte`
- Guards de rota em `+page.ts` com `redirect(302, '/login')` — nunca em `onMount`
- Navegacao apos submit: `window.location.href` para troca de pagina completa
- Formularios: `novalidate` + validacao JS, `autocomplete` correto, `maxlength` em todos inputs

**Estado e API:**
- Stores Svelte para estado global (auth, toasts)
- `ApiClient` singleton com metodos tipados — nunca `fetch` direto nas paginas
- Tipos TypeScript espelhando DTOs do backend — sem `any` implicito
- Try/catch em todo fetch, mensagem amigavel ao usuario, nunca expor stack trace

**Estilo:**
- CSS global em `app.css` com design system (variaveis, utilitarios, componentes)
- Componentes reutilizaveis em `src/lib/components/`
- Responsividade em todas as paginas (mobile-first)

**Seguranca:**
- CSP no `app.html` com politica restritiva
- Token em localStorage (tradeoff SPA), nunca logar token no console
- Sem `@html` com dados do usuario sem sanitizacao
- Input validation no cliente E no servidor

---

## Etapa 0 — Estrutura do Projeto (sem banco) [CONCLUIDA]

**Objetivo:** esqueleto da aplicacao funcionando com servidor HTTP e frontend SvelteKit, sem dependencia de banco de dados.

### Backend
- [x] Criar `go.mod` e adicionar dependencias iniciais: `go-chi/chi/v5`, `httprate`, `swaggo/swag`, JWT, uuid
- [x] Criar estrutura de diretorios: `cmd/api/`, `internal/{domain,handler,service,repository/memory,repository/postgres,middleware}`, `pkg/`, `migrations/`, `config/`, `docs/swagger/`
- [x] Criar `cmd/api/main.go` com servidor HTTP usando chi, montando rota `GET /swagger/*` (swagger UI) usando repositories in-memory
- [x] Criar `config/app.env` com variaveis de ambiente (porta, JWT secret)
- [x] Security headers OWASP A02 (X-Content-Type-Options, X-Frame-Options, CSP, HSTS, Referrer-Policy, Permissions-Policy)
- [x] Rate limiting global com `httprate`
- [x] Structured logging com `slog`
- [x] Rota `GET /health` para verificacao de saude

### Frontend
- [x] Criar projeto SvelteKit 5 com TypeScript, Vite, Playwright, Vitest
- [x] Criar `app.html` com meta tags de seguranca (CSP, X-Content-Type-Options, Referrer-Policy)
- [x] Criar `app.css` com design system completo (variaveis, tipografia, cores, componentes, utilitarios, responsividade)
- [x] Criar `vite.config.ts` com proxy `/api` -> `http://localhost:8080`
- [x] Criar `src/lib/api/client.ts` — ApiClient singleton com `setToken`, headers, tratamento de erro
- [x] Criar `src/lib/types/index.ts` — tipos TypeScript espelhando DTOs do backend
- [x] Criar `src/lib/stores/auth.ts` — store de autenticacao (login, logout, JWT parsing, localStorage)
- [x] Criar `src/lib/stores/toasts.ts` — store de notificacoes toast (success, error, info)
- [x] Criar `src/lib/components/Navbar.svelte` — barra de navegacao responsiva com estado de autenticacao
- [x] Criar `src/lib/components/Toast.svelte` + `ToastContainer.svelte` — sistema de notificacoes
- [x] Criar `+layout.svelte` com Navbar e ToastContainer

### Infra
- [x] Criar `V1__schema_inicial.sql` em `migrations/` com SQL completo do banco (pronto para Etapa 12)
- [x] Criar `docker-compose.yml` com PostgreSQL 14, pgAdmin e Flyway (disponivel, mas nao obrigatorio ate Etapa 12)
- [x] `.gitignore` cobrindo artefatos gerados (node_modules, .svelte-kit, build, bin/, dist/)

**Criterio de conclusao:** `go run cmd/api/main.go` sobe o servidor sem erros; `npm run dev` sobe o frontend com proxy funcional; nenhuma dependencia de Docker ou banco para rodar.

---

## Etapa 1 — Autenticacao e Usuarios [CONCLUIDA]

**Objetivo:** registro e login funcional para os tres tipos de usuario. Dados em memoria. Frontend com fluxo completo de autenticacao.

### Backend — Dominio
- [x] `internal/domain/usuario.go` — struct `Usuario`, tipos (`CLIENTE`, `PROFISSIONAL`, `ADMIN`), erros de dominio, validacoes (email, senha NIST 800-63b, tipo)
- [x] `internal/domain/auth.go` — interface `UsuarioRepository` (composicao Reader+Writer), DTOs (`RegistroRequest`, `LoginRequest`), `TokenPayload` com claims JWT

### Backend — Repository (in-memory)
- [x] `internal/repository/memory/usuario.go` — implementa `UsuarioRepository` com `map[uuid.UUID]*Usuario` + `map[string]uuid.UUID` (indice por email) + `sync.RWMutex`, retorna copias para evitar mutacao

### Backend — Config
- [x] `internal/config/config.go` — carrega variaveis de ambiente, valida JWT_SECRET >= 32 chars (warning em dev, fatal em prod)

### Backend — Service
- [x] `internal/service/auth_service.go` — `Registrar`, `Login`, `GerarToken` (JWT HS256), `ValidarToken` (bloqueia alg:none)
- [x] Protecao contra timing attack no Login (bcrypt nos dois caminhos — OWASP A07)
- [x] Custo bcrypt 12 em producao, `BcryptCostTeste` (4) para testes
- [x] `SolicitarRecuperacao` — gera token UUID com expiracao 1h, sem enumeracao de usuarios (OWASP A07)
- [x] `RedefinirSenha` — valida token + expiracao, hash nova senha com bcrypt, limpa token

### Backend — Handler
- [x] `POST /auth/registro/cliente` — cria usuario tipo CLIENTE
- [x] `POST /auth/registro/profissional` — cria usuario tipo PROFISSIONAL
- [x] `POST /auth/login` — retorna JWT (Bearer token, 15 min)
- [x] `POST /auth/solicitar-recuperacao-senha` — solicita token de recuperacao
- [x] `POST /auth/redefinir-senha` — redefine senha com token valido (protegido com guard de token)
- [x] `internal/handler/response.go` — helpers `RespostaJSON` e `RespostaErro`
- [x] MaxBytesReader(1MB) em todos os handlers (OWASP A06)
- [x] Mapeamento de erros de dominio para status HTTP corretos

### Backend — Middleware
- [x] `internal/middleware/auth.go` — `Autenticar` (JWT), `RequererTipo` (autorizacao), `UsuarioDoContexto` (helper)

### Backend — Integracao
- [x] `cmd/api/main.go` — DI completa (repo -> service -> handler), slog estruturado, httprate 10/min em /auth/*, security headers OWASP, rota protegida /api/v1/me

### Backend — Testes
- [x] `internal/service/auth_service_test.go` — 16 testes unitarios (registro, login, token, recuperacao, redefinicao)
- [x] `internal/handler/auth_handler_test.go` — 11 testes HTTP (registro, login, rotas protegidas, recuperacao)
- [x] `internal/middleware/auth_test.go` — 8 testes de middleware (autenticacao, autorizacao, contexto)

### Frontend — Paginas
- [x] `/login` — formulario de login com email/senha, validacao client-side, link para recuperacao, links para registro
- [x] `/registro` — registro de cliente com email, senha, confirmacao, indicador de forca de senha (barra 4 segmentos), validacao 8-72 chars
- [x] `/registro/profissional` — registro de profissional com checklist de requisitos, card informativo de documentacao, aviso de 48h
- [x] `/recuperar-senha` — formulario de email, estado de sucesso, exibicao do token em modo dev
- [x] `/redefinir-senha` — nova senha + confirmacao, token via query string, `+page.ts` com guard (redirect se sem token)
- [x] `/dashboard` — conteudo role-based (CLIENTE, PROFISSIONAL, ADMIN), grid de stats, secoes especificas por tipo
- [x] `/` — landing page com hero, "como funciona" (3 passos), vitrine de servicos (7 tipos), CTA final

### Frontend — API Client
- [x] Metodos: `registrarCliente`, `registrarProfissional`, `login`, `solicitarRecuperacao`, `redefinirSenha`, `me`

### Decisoes Arquiteturais
- [x] Verificacao de email adiada para Etapa 9 (quando houver infra de envio) — login NAO exige email verificado
- [x] Recuperacao de senha funcional: em dev retorna token na resposta; em prod loga o token

**Criterio de conclusao:** testes unitarios do service e testes HTTP cobrindo registro, login, recuperacao de senha, token invalido e rota protegida — tudo sem banco. Frontend com fluxo completo de auth (registro -> login -> dashboard -> recuperacao -> redefinicao). Login NAO exige email verificado.

---

## Etapa 2 — Cadastro Completo de Clientes e Profissionais

**Objetivo:** fluxos completos de cadastro com todas as entidades relacionadas. Dados em memoria. Frontend com formularios de perfil completo.

### Backend — Cliente
- [ ] `internal/domain/cliente.go` — struct `Cliente` (ID, UsuarioID, Nome, CPF, Telefone, Score), interface `ClienteRepository`, validacoes (CPF com digito verificador, score 0-100)
- [ ] `internal/repository/memory/cliente.go` — implementacao in-memory com indices por UsuarioID e CPF
- [ ] `internal/service/cliente_service.go` — `Criar` (a partir do usuario autenticado), `Atualizar`, `BuscarPorUsuarioID`
- [ ] `GET /clientes/me` — perfil do cliente autenticado
- [ ] `PUT /clientes/me` — atualizar perfil (nome, telefone)
- [ ] Testes unitarios: criacao, atualizacao, validacao CPF, busca

### Backend — Enderecos do Cliente
- [ ] `internal/domain/endereco.go` — struct `Endereco` (com lat/lon, num_quartos, num_banheiros, num_salas, num_cozinhas, area_m2, principal), interface `EnderecoRepository`, validacao de CEP
- [ ] `internal/repository/memory/endereco.go` — implementacao in-memory com indice por ClienteID
- [ ] `internal/service/endereco_service.go` — `Criar`, `Listar`, `Atualizar`, `DefinirPrincipal`, `Remover`
- [ ] `POST /clientes/me/enderecos` — criar endereco
- [ ] `GET /clientes/me/enderecos` — listar enderecos
- [ ] `PUT /clientes/me/enderecos/:id` — atualizar endereco
- [ ] `DELETE /clientes/me/enderecos/:id` — remover endereco
- [ ] Testes unitarios: CRUD completo, validacao CEP, logica de endereco principal

### Backend — Profissional
- [ ] `internal/domain/profissional.go` — struct `Profissional` (com status PENDENTE|APROVADA|REPROVADA|SUSPENSA|DESCREDENCIADA, nota_media, total_servicos, mei), interface `ProfissionalRepository`, validacoes
- [ ] `internal/repository/memory/profissional.go` — implementacao in-memory com busca por regiao/nota
- [ ] `internal/service/profissional_service.go` — `Criar`, `Atualizar`, `BuscarPorUsuarioID`, `BuscarPorID`
- [ ] `GET /profissionais/me` — perfil da profissional autenticada
- [ ] `PUT /profissionais/me` — atualizar perfil (nome, telefone, foto_url, mei)
- [ ] Testes unitarios: criacao, atualizacao, status inicial PENDENTE, validacoes

### Backend — Documentos e Referencias
- [ ] `internal/domain/documento.go` — struct `Documento` (tipo: RG_FRENTE|RG_VERSO|CPF|COMPROVANTE|FOTO|OUTRO, status: PENDENTE|APROVADO|REPROVADO), interface `DocumentoRepository`
- [ ] `internal/domain/referencia.go` — struct `Referencia` (nome_contato, telefone_contato, status: PENDENTE|CONFIRMADA|NAO_CONFIRMADA), interface `ReferenciaRepository`
- [ ] `internal/repository/memory/documento.go` + `internal/repository/memory/referencia.go`
- [ ] `internal/service/credenciamento_service.go` — orquestra aprovacao, upload de URL de documento, validacao de referencias
- [ ] `POST /profissionais/me/documentos` — enviar documento (URL)
- [ ] `GET /profissionais/me/documentos` — listar documentos enviados
- [ ] `POST /profissionais/me/referencias` — adicionar referencia
- [ ] `GET /profissionais/me/referencias` — listar referencias
- [ ] Testes unitarios: envio de documentos, validacao de tipos, fluxo de credenciamento

### Backend — Regioes e Disponibilidade
- [ ] `internal/domain/regiao.go` — struct `Regiao` (nome, cidade, estado, cep_inicio, cep_fim, ativa), interface `RegiaoRepository`
- [ ] `internal/domain/disponibilidade.go` — struct `Disponibilidade` (dia_semana 0-6, hora_inicio, hora_fim), interface `DisponibilidadeRepository`
- [ ] `internal/repository/memory/regiao.go` + `internal/repository/memory/disponibilidade.go`
- [ ] `GET /regioes` — listar regioes ativas (publica)
- [ ] `PUT /profissionais/me/regioes` — definir regioes de atuacao
- [ ] `GET /profissionais/me/regioes` — listar regioes da profissional
- [ ] `PUT /profissionais/me/disponibilidades` — definir disponibilidade semanal
- [ ] `GET /profissionais/me/disponibilidades` — listar disponibilidade
- [ ] Testes unitarios: associacao de regioes, validacao de horarios, dia da semana

### Frontend — Completar Perfil do Cliente
- [ ] `/dashboard/perfil` — formulario de perfil do cliente com campos: nome completo, CPF (mascara e validacao de digito verificador), telefone (mascara)
- [ ] Componente `CpfInput.svelte` — input com mascara XXX.XXX.XXX-XX e validacao em tempo real
- [ ] Componente `TelefoneInput.svelte` — input com mascara (XX) XXXXX-XXXX
- [ ] Fluxo: apos primeiro login como CLIENTE, redirecionar para completar perfil se dados obrigatorios faltam
- [ ] Exibir estado do perfil no dashboard (completo / incompleto)

### Frontend — Enderecos do Cliente
- [ ] `/dashboard/enderecos` — lista de enderecos com cards (principal destacado), botoes de editar/remover/definir principal
- [ ] `/dashboard/enderecos/novo` — formulario de novo endereco com campos: CEP (mascara e busca automatica), logradouro, numero, complemento, bairro, cidade, estado, detalhes do imovel (quartos, banheiros, salas, cozinhas, area_m2)
- [ ] Componente `CepInput.svelte` — input com mascara XXXXX-XXX e integracao com ViaCEP para autopreenchimento
- [ ] Componente `EnderecoCard.svelte` — card reutilizavel para exibir endereco
- [ ] Validacao: CEP formato valido, campos obrigatorios, numero de comodos > 0

### Frontend — Completar Perfil da Profissional
- [ ] `/dashboard/perfil` — formulario de perfil da profissional: nome, CPF, RG, telefone, foto (URL), MEI (checkbox)
- [ ] Exibir status do credenciamento (PENDENTE, APROVADA, REPROVADA) com badge colorido
- [ ] Secao de documentos obrigatorios com checklist visual: RG frente, RG verso, comprovante de residencia
- [ ] Upload de documentos (URL no MVP — sem storage de arquivos)
- [ ] Secao de referencias: formulario para adicionar referencia (nome + telefone), lista de referencias com status

### Frontend — Regioes e Disponibilidade
- [ ] `/dashboard/regioes` — seletor de regioes de atuacao com lista de regioes disponiveis, checkbox para selecionar/deselecionar
- [ ] `/dashboard/disponibilidade` — grade semanal interativa (seg-dom) com seletor de horarios (hora_inicio, hora_fim) por dia
- [ ] Componente `GradeSemanal.svelte` — grade visual de disponibilidade com toggle por dia e seletor de horario
- [ ] Componente `RegiaoSelector.svelte` — lista de regioes com busca e selecao multipla

### Frontend — API Client (novos metodos)
- [ ] `buscarPerfilCliente()`, `atualizarPerfilCliente(req)` — GET/PUT /clientes/me
- [ ] `criarEndereco(req)`, `listarEnderecos()`, `atualizarEndereco(id, req)`, `removerEndereco(id)` — CRUD /clientes/me/enderecos
- [ ] `buscarPerfilProfissional()`, `atualizarPerfilProfissional(req)` — GET/PUT /profissionais/me
- [ ] `enviarDocumento(req)`, `listarDocumentos()` — POST/GET /profissionais/me/documentos
- [ ] `adicionarReferencia(req)`, `listarReferencias()` — POST/GET /profissionais/me/referencias
- [ ] `listarRegioes()` — GET /regioes
- [ ] `definirRegioesAtuacao(req)`, `listarRegioesAtuacao()` — PUT/GET /profissionais/me/regioes
- [ ] `definirDisponibilidades(req)`, `listarDisponibilidades()` — PUT/GET /profissionais/me/disponibilidades

### Frontend — Tipos TypeScript (novos)
- [ ] `Cliente`, `ClienteRequest`, `Endereco`, `EnderecoRequest`
- [ ] `Profissional`, `ProfissionalRequest`, `Documento`, `DocumentoRequest`
- [ ] `Referencia`, `ReferenciaRequest`, `Regiao`, `Disponibilidade`, `DisponibilidadeRequest`

**Criterio de conclusao:** profissional consegue completar cadastro completo via frontend (perfil + documentos + referencias + regioes + disponibilidade); status vai para PENDENTE; cliente completa perfil e gerencia enderecos; campos obrigatorios validados no frontend E backend; tudo sem banco. `go test ./...` e `npm run check` passando.

---

## Etapa 3 — Catalogo: Categorias, Opcionais e Tabela de Precos

**Objetivo:** catalogo de servicos e precificacao por regiao em memoria, com seed dos dados padrao. Frontend com vitrine de servicos e consulta de precos.

### Backend
- [ ] `internal/domain/catalogo.go` — structs `CategoriaServico`, `Opcional`, `TabelaPrecos` + interfaces de repository
- [ ] `internal/repository/memory/catalogo.go` — implementacao in-memory com seed dos 7 tipos de servico e 7 add-ons padrao
- [ ] `internal/service/precificacao_service.go` — `CalcularValorReferencia(solicitacao)` retorna valor estimado e duracao
- [ ] `GET /categorias` — lista categorias ativas (publica)
- [ ] `GET /categorias/:id/opcionais` — lista opcionais disponiveis (publica)
- [ ] `GET /precos?categoria=&regiao=` — consulta tabela de precos (publica)
- [ ] Testes unitarios: calculo com comodos, opcionais, frequencia; desconto semanal/quinzenal; acrescimo fds

### Frontend — Vitrine de Servicos
- [ ] Atualizar pagina home (`/`) para carregar categorias da API dinamicamente (substituir dados estaticos)
- [ ] `/servicos` — pagina de catalogo com todos os tipos de servico, descricao, duracao minima, preco a partir de
- [ ] `/servicos/:id` — detalhe da categoria com lista de opcionais e calculadora de preco interativa

### Frontend — Calculadora de Preco
- [ ] Componente `CalculadoraPreco.svelte` — formulario interativo:
  - Selecao de categoria de servico
  - Selecao de regiao (dropdown)
  - Numero de comodos (quartos, banheiros, salas, cozinhas)
  - Opcionais (checkboxes com valor extra visivel)
  - Frequencia (unica, semanal, quinzenal, 2x/semana) com desconto exibido
  - Valor estimado atualizado em tempo real
  - Duracao estimada
- [ ] Componente `CategoriaCard.svelte` — card de categoria com icone, nome, descricao, duracao, preco base
- [ ] Componente `OpcionalCheckbox.svelte` — checkbox com nome, descricao, valor extra e tempo extra

### Frontend — API Client (novos metodos)
- [ ] `listarCategorias()` — GET /categorias
- [ ] `listarOpcionais(categoriaId)` — GET /categorias/:id/opcionais
- [ ] `consultarPrecos(categoriaId, regiaoId)` — GET /precos

### Frontend — Tipos TypeScript (novos)
- [ ] `CategoriaServico`, `Opcional`, `TabelaPrecos`, `CalculoPrecoRequest`, `CalculoPrecoResponse`

**Criterio de conclusao:** servico de precificacao calcula corretamente dado comodos, opcionais e frequencia; testado com casos dos 7 tipos de servico; seed carregado na inicializacao do `main.go`. Frontend exibe catalogo dinamico e calculadora de preco funcional.

---

## Etapa 4 — Solicitacoes de Servico

**Objetivo:** cliente consegue fazer uma solicitacao completa com calculo automatico. Dados em memoria. Frontend com fluxo multi-step de solicitacao.

### Backend
- [ ] `internal/domain/solicitacao.go` — struct `Solicitacao` (com status AGUARDANDO|ATRIBUIDA|CONFIRMADA|EM_ANDAMENTO|CONCLUIDA|CANCELADA), interface `SolicitacaoRepository`, validacoes (24h antecedencia, etc.)
- [ ] `internal/domain/solicitacao_opcional.go` — tabela associativa N:N solicitacao <-> opcionais
- [ ] `internal/repository/memory/solicitacao.go` — implementacao in-memory com indices por ClienteID e status
- [ ] `internal/service/solicitacao_service.go`:
  - `Criar` — valida dados, valida endereco pertence ao cliente, calcula valor de referencia via `precificacao_service`, persiste
  - `Cancelar` — valida regra das 24h, aplica penalizacao no score se cabivel
  - `BuscarParaCliente` — lista solicitacoes do cliente com filtros
  - `BuscarPorID` — detalhe com opcionais
- [ ] `POST /solicitacoes` — criar nova solicitacao
- [ ] `GET /solicitacoes` — listar solicitacoes do cliente autenticado
- [ ] `GET /solicitacoes/:id` — detalhe da solicitacao
- [ ] `DELETE /solicitacoes/:id` — cancelar solicitacao
- [ ] Testes unitarios: criacao com calculo, validacao 24h, cancelamento com penalizacao, validacao de endereco/categoria

### Frontend — Fluxo de Solicitacao (multi-step)
- [ ] `/solicitar` — formulario multi-step:
  - **Passo 1:** Selecao de endereco (lista de enderecos do cliente, opcao de cadastrar novo)
  - **Passo 2:** Selecao de categoria de servico + opcionais (com detalhes do imovel carregados do endereco)
  - **Passo 3:** Data, horario e frequencia (calendario com validacao de 24h antecedencia, seletor de horario)
  - **Passo 4:** Resumo com valor estimado, duracao, detalhes do endereco — botao de confirmar
- [ ] Componente `StepIndicator.svelte` — indicador de progresso multi-step (1/4, 2/4, etc.) reutilizavel
- [ ] Componente `EnderecoSelector.svelte` — cards de enderecos do cliente com radio button para selecao
- [ ] Componente `CalendarioAgendamento.svelte` — seletor de data com bloqueio de datas < 24h
- [ ] Componente `ResumoSolicitacao.svelte` — card de resumo com todos os detalhes antes de confirmar

### Frontend — Lista de Solicitacoes
- [ ] `/dashboard/solicitacoes` — lista de solicitacoes do cliente com filtro por status, ordenacao por data
- [ ] Componente `SolicitacaoCard.svelte` — card com status (badge colorido), data, categoria, valor, endereco resumido
- [ ] `/dashboard/solicitacoes/:id` — detalhe da solicitacao com timeline de status, opcao de cancelar

### Frontend — Cancelamento
- [ ] Modal de confirmacao de cancelamento com aviso de penalizacao se < 24h
- [ ] Componente `ModalConfirmacao.svelte` — modal reutilizavel com titulo, mensagem, botoes confirmar/cancelar

### Frontend — API Client (novos metodos)
- [ ] `criarSolicitacao(req)` — POST /solicitacoes
- [ ] `listarSolicitacoes(filtros?)` — GET /solicitacoes
- [ ] `buscarSolicitacao(id)` — GET /solicitacoes/:id
- [ ] `cancelarSolicitacao(id)` — DELETE /solicitacoes/:id

### Frontend — Tipos TypeScript (novos)
- [ ] `Solicitacao`, `SolicitacaoRequest`, `SolicitacaoDetalhe`, `StatusSolicitacao`

**Criterio de conclusao:** solicitacao criada com valor calculado automaticamente via frontend multi-step; cancelamento com menos de 24h desconta score; lista de solicitacoes filtravel; tudo sem banco. `go test ./...` e `npm run check` passando.

---

## Etapa 5 — Atribuicao de Profissional (Motor de Matching)

**Objetivo:** quando uma solicitacao e criada, o sistema encontra e notifica a melhor profissional disponivel. Dados em memoria. Frontend com interface de aceite/recusa para profissional.

### Backend
- [ ] `internal/domain/servico.go` — struct `Servico` (com status AGENDADO|EM_ANDAMENTO|CONCLUIDO|CANCELADO|NO_SHOW, campos de checkin/checkout com lat/lon), interface `ServicoRepository`
- [ ] `internal/repository/memory/servico.go` — implementacao in-memory com indices por ProfissionalID e SolicitacaoID
- [ ] `internal/service/matching_service.go`:
  - `BuscarCandidatas(solicitacao)` — filtra por regiao, disponibilidade, categoria, nota minima, verifica LC 150/2015
  - `AtribuirProfissional(solicitacao, profissional)` — cria o `Servico`, atualiza status da solicitacao para ATRIBUIDA
  - `RedireccionarParaProxima(servico)` — chamado apos timeout de 30 min ou recusa
- [ ] Logica do limite legal implementada no service (consulta in-memory equivalente a `vw_servicos_por_semana`): bloqueia se profissional ja tem 2 servicos na semana no mesmo endereco
- [ ] `POST /solicitacoes/:id/aceitar` — profissional aceita (requer autenticacao + tipo PROFISSIONAL)
- [ ] `POST /solicitacoes/:id/recusar` — profissional recusa
- [ ] `GET /profissionais/me/solicitacoes-pendentes` — lista solicitacoes atribuidas aguardando aceite
- [ ] Testes unitarios: matching com profissional disponivel, rotacao por recusa, bloqueio pela regra LC 150/2015, timeout 30 min, prioridade por nota + proximidade

### Frontend — Painel da Profissional (Solicitacoes)
- [ ] `/dashboard/solicitacoes-pendentes` — lista de solicitacoes atribuidas a profissional aguardando aceite
- [ ] Componente `SolicitacaoPendenteCard.svelte` — card com detalhes da solicitacao (categoria, endereco resumido, data, horario, valor, duracao), timer de 30min, botoes aceitar/recusar
- [ ] Componente `TimerCountdown.svelte` — countdown visual de 30 minutos com barra de progresso
- [ ] Modal de confirmacao ao recusar (motivo opcional)
- [ ] Atualizar dashboard da profissional para exibir contador de solicitacoes pendentes

### Frontend — Visualizacao do Matching (Cliente)
- [ ] Atualizar detalhe da solicitacao (`/dashboard/solicitacoes/:id`) para exibir:
  - Status "Buscando profissional..." quando AGUARDANDO
  - Dados da profissional atribuida (nome, nota, foto) quando ATRIBUIDA/CONFIRMADA
  - Timeline de status atualizada

### Frontend — API Client (novos metodos)
- [ ] `listarSolicitacoesPendentes()` — GET /profissionais/me/solicitacoes-pendentes
- [ ] `aceitarSolicitacao(id)` — POST /solicitacoes/:id/aceitar
- [ ] `recusarSolicitacao(id)` — POST /solicitacoes/:id/recusar

### Frontend — Tipos TypeScript (novos)
- [ ] `Servico`, `StatusServico`, `SolicitacaoPendente`

**Criterio de conclusao:** testes unitarios cobrem matching com profissional disponivel, rotacao por recusa, bloqueio pela regra LC 150/2015 e timeout de 30 min — sem banco. Frontend permite profissional aceitar/recusar e cliente acompanhar status.

---

## Etapa 6 — Execucao do Servico (Check-in / Check-out)

**Objetivo:** controle da execucao em campo com geolocalizacao. Dados em memoria. Frontend com interface de execucao para profissional e acompanhamento para cliente.

### Backend
- [ ] `internal/domain/historico.go` — struct `HistoricoStatus` (status_anterior, status_novo, observacao), interface `HistoricoRepository`
- [ ] `internal/repository/memory/historico.go` — implementacao in-memory
- [ ] `internal/service/execucao_service.go`:
  - `CheckIn(servicoID, profissionalID, lat, lon)` — valida que e a profissional certa, status AGENDADO, registra horario e localizacao
  - `CheckOut(servicoID, profissionalID, lat, lon)` — valida status EM_ANDAMENTO, registra conclusao
  - `ConfirmarConclusao(servicoID, atorID, atorTipo)` — cliente ou profissional confirma
- [ ] `internal/service/historico_service.go` — registra mudancas de status em memoria
- [ ] `POST /servicos/:id/checkin` — check-in com coordenadas
- [ ] `POST /servicos/:id/checkout` — check-out com coordenadas
- [ ] `POST /servicos/:id/confirmar` — confirmacao de conclusao
- [ ] `GET /servicos/:id` — detalhe do servico (com historico de status)
- [ ] `GET /profissionais/me/servicos` — servicos da profissional (filtro por status)
- [ ] `GET /clientes/me/servicos` — servicos do cliente (filtro por status)
- [ ] Testes unitarios: fluxo completo AGENDADO -> EM_ANDAMENTO -> CONCLUIDO, check-in/out com coordenadas, confirmacao bilateral, historico registrado

### Frontend — Execucao (Profissional)
- [ ] `/dashboard/servicos` — lista de servicos da profissional com filtro por status (agendados, em andamento, concluidos)
- [ ] `/dashboard/servicos/:id` — detalhe do servico com:
  - Informacoes do endereco e categoria
  - Timeline visual de status (AGENDADO -> EM_ANDAMENTO -> CONCLUIDO)
  - Botao "Check-in" (pede permissao de geolocalizacao do navegador)
  - Botao "Check-out" (apos check-in)
  - Botao "Confirmar conclusao"
- [ ] Componente `TimelineStatus.svelte` — timeline vertical com icones, horarios e status
- [ ] Componente `BotaoGeolocalizacao.svelte` — botao que captura lat/lon via `navigator.geolocation`, exibe status de carregamento
- [ ] Componente `ServicoCard.svelte` — card de servico com status, data, horario, endereco, categoria

### Frontend — Acompanhamento (Cliente)
- [ ] `/dashboard/servicos` — lista de servicos do cliente com filtro por status
- [ ] `/dashboard/servicos/:id` — detalhe do servico com timeline atualizada, botao de confirmar conclusao (lado do cliente)
- [ ] Exibir horarios de check-in/check-out quando disponiveis

### Frontend — API Client (novos metodos)
- [ ] `checkinServico(id, lat, lon)` — POST /servicos/:id/checkin
- [ ] `checkoutServico(id, lat, lon)` — POST /servicos/:id/checkout
- [ ] `confirmarServico(id)` — POST /servicos/:id/confirmar
- [ ] `buscarServico(id)` — GET /servicos/:id
- [ ] `listarServicosProfissional(filtros?)` — GET /profissionais/me/servicos
- [ ] `listarServicosCliente(filtros?)` — GET /clientes/me/servicos

### Frontend — Tipos TypeScript (novos)
- [ ] `CheckinRequest`, `CheckoutRequest`, `HistoricoStatus`

**Criterio de conclusao:** servico muda de AGENDADO -> EM_ANDAMENTO -> CONCLUIDO com timestamps e coordenadas registradas; historico auditavel; frontend permite check-in/out com geolocalizacao; testes sem banco. `go test ./...` e `npm run check` passando.

---

## Etapa 7 — Avaliacoes e Reputacao

**Objetivo:** avaliacao mutua pos-servico e atualizacao automatica de nota. Dados em memoria. Frontend com formularios de avaliacao e exibicao de reputacao.

### Backend
- [ ] `internal/domain/avaliacao.go` — structs `AvaliacaoCliente` (nota, comentario, pontualidade, qualidade, educacao), `AvaliacaoProfissional` (nota, ambiente, materiais, respeito), interfaces de repository, validacao nota 1-5
- [ ] `internal/repository/memory/avaliacao.go` — implementacao in-memory com indices por ServicoID e ProfissionalID
- [ ] `internal/service/avaliacao_service.go`:
  - `AvaliarProfissional` — valida servico CONCLUIDO, persiste, recalcula `nota_media` no service (equivalente ao trigger do banco)
  - `AvaliarCliente` — persiste avaliacao interna
  - `VerificarAlertas` — dispara alerta se nota < 4.0, suspensao se nota < 3.5 por 3 consecutivos
  - `ListarAvaliacoesProfissional` — avaliacoes publicas com paginacao
- [ ] `POST /servicos/:id/avaliacoes/profissional` — cliente avalia profissional (requer CLIENTE)
- [ ] `POST /servicos/:id/avaliacoes/cliente` — profissional avalia cliente (requer PROFISSIONAL)
- [ ] `GET /profissionais/:id/avaliacoes` — avaliacoes publicas da profissional
- [ ] Testes unitarios: avaliacao -> nota recalculada, suspensao por 3 consecutivas < 3.5, validacao nota 1-5, avaliacao duplicada bloqueada

### Frontend — Avaliacao (Cliente avalia Profissional)
- [ ] Apos conclusao de servico, exibir prompt de avaliacao no detalhe do servico
- [ ] Componente `FormularioAvaliacao.svelte` — formulario com:
  - Estrelas interativas para nota geral (1-5) — componente `RatingStars.svelte`
  - Estrelas para criterios individuais: pontualidade, qualidade, educacao
  - Campo de comentario (textarea com maxlength)
  - Preview antes de enviar
- [ ] Componente `RatingStars.svelte` — seletor de estrelas reutilizavel (hover, click, readonly)

### Frontend — Avaliacao (Profissional avalia Cliente)
- [ ] Formulario similar no detalhe do servico da profissional: nota, ambiente, materiais, respeito

### Frontend — Exibicao de Reputacao
- [ ] Componente `PerfilPublicoProfissional.svelte` — nota media, total de servicos, lista de avaliacoes recentes
- [ ] `/profissionais/:id` — pagina publica da profissional com avaliacoes (sem dados sensiveis)
- [ ] Exibir nota media e selo no card da profissional (onde aplicavel)
- [ ] Badge de alerta no dashboard da profissional se nota < 4.0

### Frontend — API Client (novos metodos)
- [ ] `avaliarProfissional(servicoId, req)` — POST /servicos/:id/avaliacoes/profissional
- [ ] `avaliarCliente(servicoId, req)` — POST /servicos/:id/avaliacoes/cliente
- [ ] `listarAvaliacoesProfissional(profissionalId)` — GET /profissionais/:id/avaliacoes

### Frontend — Tipos TypeScript (novos)
- [ ] `AvaliacaoCliente`, `AvaliacaoClienteRequest`, `AvaliacaoProfissional`, `AvaliacaoProfissionalRequest`

**Criterio de conclusao:** avaliacao inserida -> nota recalculada no service; suspensao disparada corretamente nos cenarios de nota < 3.5 por 3 consecutivos; frontend com formularios de estrelas e perfil publico; testes unitarios sem banco. `go test ./...` e `npm run check` passando.

---

## Etapa 8 — Recorrencias

**Objetivo:** suporte a servicos recorrentes (semanal, quinzenal, 2x/semana). Dados em memoria. Frontend com configuracao e gestao de recorrencias.

### Backend
- [ ] `internal/domain/recorrencia.go` — struct `Recorrencia` (frequencia, dia_semana_1, dia_semana_2, hora_inicio, status ATIVA|PAUSADA|CANCELADA, proxima_data), interface `RecorrenciaRepository`, calculo de proxima data
- [ ] `internal/repository/memory/recorrencia.go` — implementacao in-memory
- [ ] `internal/service/recorrencia_service.go`:
  - `Criar` — baseado em solicitacao confirmada com frequencia != UNICA
  - `GerarProximasSolicitacoes` — cria solicitacoes futuras a partir das recorrencias ativas
  - `Pausar` — muda status para PAUSADA
  - `Retomar` — muda status para ATIVA, recalcula proxima_data
  - `Cancelar` — muda status para CANCELADA
- [ ] `GET /recorrencias` — listar recorrencias do cliente
- [ ] `GET /recorrencias/:id` — detalhe da recorrencia
- [ ] `POST /recorrencias/:id/pausar` — pausar recorrencia
- [ ] `POST /recorrencias/:id/retomar` — retomar recorrencia
- [ ] `DELETE /recorrencias/:id` — cancelar recorrencia
- [ ] Testes unitarios: geracao de proximas solicitacoes, pausa/retomada, cancelamento, limite LC 150/2015 com rotacao

### Frontend — Gestao de Recorrencias
- [ ] `/dashboard/recorrencias` — lista de recorrencias com status (badge), frequencia, proximo servico, endereco
- [ ] `/dashboard/recorrencias/:id` — detalhe da recorrencia com:
  - Frequencia e dias selecionados
  - Proximo servico agendado
  - Historico de servicos gerados
  - Botoes: pausar, retomar, cancelar
- [ ] Componente `RecorrenciaCard.svelte` — card com frequencia, status, proximo servico
- [ ] No fluxo de solicitacao (Etapa 4), quando frequencia != UNICA, exibir confirmacao de recorrencia apos primeiro servico

### Frontend — Integracao com Solicitacao
- [ ] Atualizar formulario de solicitacao para destacar descontos por frequencia
- [ ] Exibir info de recorrencia no detalhe da solicitacao quando aplicavel

### Frontend — API Client (novos metodos)
- [ ] `listarRecorrencias()` — GET /recorrencias
- [ ] `buscarRecorrencia(id)` — GET /recorrencias/:id
- [ ] `pausarRecorrencia(id)` — POST /recorrencias/:id/pausar
- [ ] `retomarRecorrencia(id)` — POST /recorrencias/:id/retomar
- [ ] `cancelarRecorrencia(id)` — DELETE /recorrencias/:id

### Frontend — Tipos TypeScript (novos)
- [ ] `Recorrencia`, `StatusRecorrencia`, `FrequenciaServico`

**Criterio de conclusao:** recorrencia semanal gera solicitacoes corretamente; pausa e cancelamento funcionam; rotacao de profissional mantem-se dentro do limite da LC 150/2015; frontend com gestao completa de recorrencias. `go test ./...` e `npm run check` passando.

---

## Etapa 9 — Notificacoes e Verificacao Real de Email

**Objetivo:** sistema de notificacoes multicanal e ativacao da verificacao real de email. Dados em memoria. Frontend com central de notificacoes.

### Backend — Notificacoes
- [ ] `internal/domain/notificacao.go` — struct `Notificacao` (canal PUSH|EMAIL|SMS|IN_APP, titulo, corpo, status PENDENTE|ENVIADA|FALHOU|LIDA), interface `NotificacaoRepository`
- [ ] `internal/repository/memory/notificacao.go` — implementacao in-memory com indice por UsuarioID
- [ ] `internal/service/notificacao_service.go` — interface `Notificador` + implementacoes stub (MVP envia apenas in-app)
- [ ] Integrar notificacoes nos eventos:
  - Solicitacao criada -> notifica profissional atribuida
  - Profissional aceita -> notifica cliente
  - Vespera de servico -> notifica cliente (com nome/foto da profissional)
  - Check-in -> notifica cliente
  - Conclusao -> notifica ambos (pedido de avaliacao)
  - Avaliacao recebida -> notifica avaliado
  - Recorrencia proxima -> notifica cliente
- [ ] `GET /notificacoes` — lista notificacoes do usuario autenticado (paginada)
- [ ] `PUT /notificacoes/:id/lida` — marcar como lida
- [ ] `PUT /notificacoes/todas-lidas` — marcar todas como lidas
- [ ] `GET /notificacoes/nao-lidas/contagem` — contagem de nao lidas
- [ ] Testes unitarios: criacao de notificacao por evento, marcacao de lida, listagem paginada

### Backend — Verificacao Real de Email (adiada da Etapa 1)
- [ ] Implementar envio real de email com codigo de verificacao usando `Notificador`
- [ ] Recriar `POST /auth/verificar-email` com codigo enviado por email (nao retornado na resposta)
- [ ] Adicionar `CodigoVerificacao` e `CodigoVerificacaoExpira` na struct `Usuario`
- [ ] Integrar envio real de email na recuperacao de senha (substituir retorno do token na resposta)
- [ ] Decisao: ativar ou nao bloqueio de login para emails nao verificados

### Frontend — Central de Notificacoes
- [ ] Componente `NotificacaoBadge.svelte` — badge no Navbar com contagem de nao lidas (polling periodico ou WebSocket futuro)
- [ ] `/dashboard/notificacoes` — lista de notificacoes com:
  - Agrupamento por data (hoje, ontem, esta semana, anterior)
  - Notificacoes nao lidas destacadas
  - Botao "marcar todas como lidas"
  - Click na notificacao -> navega para o recurso relacionado e marca como lida
- [ ] Componente `NotificacaoItem.svelte` — item de notificacao com icone (por tipo de evento), titulo, corpo, horario relativo ("ha 5 min")

### Frontend — Verificacao de Email
- [ ] Atualizar fluxo de registro para exibir mensagem "verifique seu email" apos cadastro
- [ ] `/verificar-email` — pagina para inserir codigo de verificacao recebido por email
- [ ] Atualizar pagina de recuperacao de senha para nao exibir token (envio real por email)

### Frontend — Integracao com Navbar
- [ ] Atualizar `Navbar.svelte` para incluir icone de sino com badge de notificacoes
- [ ] Dropdown rapido com ultimas 5 notificacoes + link "ver todas"

### Frontend — API Client (novos metodos)
- [ ] `listarNotificacoes(pagina?)` — GET /notificacoes
- [ ] `marcarNotificacaoLida(id)` — PUT /notificacoes/:id/lida
- [ ] `marcarTodasLidas()` — PUT /notificacoes/todas-lidas
- [ ] `contagemNaoLidas()` — GET /notificacoes/nao-lidas/contagem
- [ ] `verificarEmail(req)` — POST /auth/verificar-email

### Frontend — Tipos TypeScript (novos)
- [ ] `Notificacao`, `CanalNotificacao`, `StatusNotificacao`, `VerificarEmailRequest`

**Criterio de conclusao:** eventos principais disparam notificacoes in-app; central de notificacoes funcional no frontend com badge no Navbar; verificacao de email funcional com envio real; stubs prontos para substituir por push real na Fase 2. `go test ./...` e `npm run check` passando.

---

## Etapa 10 — Painel Admin

**Objetivo:** painel completo de administracao para gestao da plataforma. Backend com endpoints de admin. Frontend com interface administrativa.

### Backend
- [ ] `internal/service/admin_service.go`:
  - Gestao de profissionais: listar (filtro por status), aprovar, reprovar, suspender, descredenciar
  - Gestao de clientes: listar, visualizar detalhes
  - Gestao de servicos: listar (filtros status, data, profissional, cliente)
  - Gestao de precos: CRUD tabela de precos por regiao
  - Gestao de regioes: CRUD regioes
  - Gestao de categorias: CRUD categorias e opcionais
- [ ] `GET /admin/profissionais` — lista com filtro por status (paginada)
- [ ] `GET /admin/profissionais/:id` — detalhe com documentos e referencias
- [ ] `PUT /admin/profissionais/:id/aprovar` — aprovar profissional
- [ ] `PUT /admin/profissionais/:id/reprovar` — reprovar com motivo
- [ ] `PUT /admin/profissionais/:id/suspender` — suspender com motivo
- [ ] `PUT /admin/profissionais/:id/descredenciar` — descredenciar
- [ ] `GET /admin/clientes` — lista paginada
- [ ] `GET /admin/clientes/:id` — detalhe do cliente
- [ ] `GET /admin/servicos` — lista com filtros (status, data, profissional, cliente)
- [ ] `GET /admin/servicos/:id` — detalhe do servico
- [ ] `POST /admin/tabela-precos` — criar entrada na tabela de precos
- [ ] `PUT /admin/tabela-precos/:id` — atualizar preco
- [ ] `GET /admin/regioes` — listar regioes
- [ ] `POST /admin/regioes` — criar regiao
- [ ] `PUT /admin/regioes/:id` — atualizar regiao
- [ ] Todos os endpoints protegidos com `RequererTipo(ADMIN)`
- [ ] Testes unitarios: aprovacao muda status, profissional aprovada aparece no matching, CRUD de precos reflete no calculo

### Frontend — Layout Admin
- [ ] `/admin` — layout dedicado para admin com sidebar de navegacao:
  - Dashboard (visao geral)
  - Profissionais (gestao)
  - Clientes (gestao)
  - Servicos (monitoramento)
  - Precos (configuracao)
  - Regioes (configuracao)
- [ ] Componente `AdminSidebar.svelte` — sidebar com links, icones, item ativo destacado
- [ ] Componente `AdminHeader.svelte` — header com titulo da pagina, breadcrumb
- [ ] Guard de rota: `+page.ts` em todas as rotas /admin com `RequererTipo(ADMIN)` no frontend

### Frontend — Dashboard Admin
- [ ] `/admin` (index) — visao geral com:
  - Cards de metricas: total de profissionais (por status), total de clientes, servicos do dia, servicos do mes
  - Lista de profissionais pendentes de aprovacao (acesso rapido)
  - Ultimos servicos

### Frontend — Gestao de Profissionais
- [ ] `/admin/profissionais` — tabela com: nome, CPF, status (badge), nota, total servicos, data cadastro
  - Filtros: status (dropdown), busca por nome/CPF
  - Paginacao
- [ ] `/admin/profissionais/:id` — detalhe completo:
  - Dados pessoais, regioes, disponibilidade
  - Documentos com visualizacao e status (aprovado/reprovado/pendente)
  - Referencias com status
  - Historico de servicos
  - Avaliacoes recebidas
  - Botoes de acao: aprovar, reprovar (com modal de motivo), suspender, descredenciar
- [ ] Componente `TabelaAdmin.svelte` — tabela reutilizavel com sort, filtro, paginacao
- [ ] Componente `ModalMotivo.svelte` — modal para inserir motivo de reprovacao/suspensao

### Frontend — Gestao de Clientes
- [ ] `/admin/clientes` — tabela com: nome, CPF, email, score, total servicos, data cadastro
- [ ] `/admin/clientes/:id` — detalhe com enderecos, historico de servicos, avaliacoes

### Frontend — Gestao de Servicos
- [ ] `/admin/servicos` — tabela com: data, categoria, cliente, profissional, status, valor
  - Filtros: status, data (range), profissional, cliente
  - Paginacao
- [ ] `/admin/servicos/:id` — detalhe completo com timeline

### Frontend — Configuracao de Precos
- [ ] `/admin/precos` — tabela de precos por categoria x regiao, editavel inline
- [ ] Formulario de novo preco: selecao de categoria + regiao + valores (preco/hora, acrescimo fds, descontos)
- [ ] `/admin/regioes` — CRUD de regioes com formulario de criacao/edicao

### Frontend — API Client (novos metodos)
- [ ] Todos os metodos de admin: listar/aprovar/reprovar/suspender profissionais, listar clientes, listar servicos, CRUD precos, CRUD regioes

### Frontend — Tipos TypeScript (novos)
- [ ] `AdminProfissionalFiltro`, `AdminServicoFiltro`, `PaginacaoRequest`, `PaginacaoResponse`

**Criterio de conclusao:** admin consegue aprovar profissional via frontend -> status muda -> profissional aparece no matching; configuracao de precos por regiao refletida no calculo de solicitacoes; todas as telas admin funcionais com filtros e paginacao. `go test ./...` e `npm run check` passando.

---

## Etapa 11 — Historico e Transacoes (Estrutura MVP)

**Objetivo:** registrar transacoes financeiras como referencia (sem cobranca online). Dados em memoria. Frontend com historico navegavel e visao financeira.

### Backend
- [ ] `internal/domain/transacao.go` — struct `Transacao` (valor, metodo DIRETO_EXTERNO, status REGISTRADA), interface `TransacaoRepository`
- [ ] `internal/domain/dados_bancarios.go` — struct `DadosBancarios` (banco, agencia, conta, tipo_conta, chave_pix, principal), interface `DadosBancariosRepository`
- [ ] `internal/repository/memory/transacao.go` + `internal/repository/memory/dados_bancarios.go` — implementacoes in-memory
- [ ] `internal/service/transacao_service.go` — `Registrar(servico)` cria transacao com status REGISTRADA e metodo DIRETO_EXTERNO
- [ ] `internal/service/dados_bancarios_service.go` — CRUD de dados bancarios da profissional (opcional no MVP)
- [ ] Integrar: toda conclusao de servico gera registro de transacao automaticamente
- [ ] `GET /clientes/me/historico` — historico de servicos com valores (paginado)
- [ ] `GET /profissionais/me/historico` — historico com ganhos de referencia (paginado)
- [ ] `GET /profissionais/me/dados-bancarios` — listar dados bancarios
- [ ] `POST /profissionais/me/dados-bancarios` — adicionar dados bancarios
- [ ] `PUT /profissionais/me/dados-bancarios/:id` — atualizar
- [ ] `GET /admin/transacoes` — visao financeira geral (filtros: data, profissional, cliente)
- [ ] `GET /admin/dashboard/financeiro` — metricas financeiras (total movimentado, por regiao, por categoria)
- [ ] Testes unitarios: geracao de transacao na conclusao, historico correto por papel, dados bancarios CRUD

### Frontend — Historico do Cliente
- [ ] `/dashboard/historico` — lista de servicos concluidos com:
  - Data, categoria, profissional (nome + nota), valor, status do pagamento
  - Filtro por periodo (data inicio/fim)
  - Resumo: total gasto no periodo, quantidade de servicos
- [ ] Componente `HistoricoCard.svelte` — card de servico concluido com valor e profissional

### Frontend — Historico da Profissional
- [ ] `/dashboard/historico` (contexto profissional) — lista de servicos concluidos com:
  - Data, categoria, cliente, valor de referencia
  - Resumo: total de ganhos no periodo, quantidade de servicos
- [ ] `/dashboard/dados-bancarios` — CRUD de dados bancarios/Pix
  - Formulario: banco, agencia, conta, tipo conta, chave Pix
  - Lista de contas cadastradas com opcao de definir principal

### Frontend — Painel Financeiro Admin
- [ ] `/admin/financeiro` — dashboard financeiro com:
  - Cards de metricas: total movimentado (mes/semana), ticket medio, servicos concluidos
  - Grafico simples de evolucao (ultimos 30 dias) — pode ser CSS puro ou lib leve
  - Tabela de transacoes com filtros (data, profissional, cliente, regiao)
- [ ] `/admin/transacoes` — lista detalhada de transacoes

### Frontend — API Client (novos metodos)
- [ ] `historicoCliente(filtros?)` — GET /clientes/me/historico
- [ ] `historicoProfissional(filtros?)` — GET /profissionais/me/historico
- [ ] `listarDadosBancarios()` — GET /profissionais/me/dados-bancarios
- [ ] `criarDadosBancarios(req)` — POST /profissionais/me/dados-bancarios
- [ ] `atualizarDadosBancarios(id, req)` — PUT /profissionais/me/dados-bancarios/:id
- [ ] `listarTransacoesAdmin(filtros?)` — GET /admin/transacoes
- [ ] `dashboardFinanceiro()` — GET /admin/dashboard/financeiro

### Frontend — Tipos TypeScript (novos)
- [ ] `Transacao`, `DadosBancarios`, `DadosBancariosRequest`, `DashboardFinanceiro`, `HistoricoFiltro`

**Criterio de conclusao:** toda conclusao de servico gera registro de transacao; historico navegavel por cliente e profissional; dados bancarios opcionais; painel financeiro admin funcional; tudo sem banco. `go test ./...` e `npm run check` passando.

---

## Etapa 12 — Persistencia em PostgreSQL (substituicao dos repositories)

**Objetivo:** substituir todas as implementacoes in-memory por implementacoes Postgres. Services e handlers nao mudam. Frontend nao muda — a troca e transparente.

### Backend — Infraestrutura
- [ ] Subir infraestrutura: `docker-compose up -d` -> PostgreSQL + pgAdmin + Flyway aplicando as migrations versionadas (`V1__auth_usuarios.sql`, `V2__cadastro_clientes_profissionais.sql`, V3…V11 conforme etapas concluidas)
- [ ] Adicionar dependencia `pgx/v5` ao `go.mod`
- [ ] Criar `internal/repository/postgres/db.go` — pool de conexoes com `pgxpool`, health check, graceful shutdown

### Backend — Repositories Postgres
Para cada repository in-memory, criar o equivalente em `internal/repository/postgres/`:
- [ ] `usuario.go` — queries parametrizadas com `$1, $2...` (OWASP A05), email normalizado
- [ ] `cliente.go` — CRUD com join em usuarios quando necessario
- [ ] `profissional.go` — busca por regiao/nota com indices
- [ ] `endereco.go` — CRUD com validacao de pertencimento ao cliente
- [ ] `regiao.go` + `disponibilidade.go`
- [ ] `documento.go` + `referencia.go`
- [ ] `catalogo.go` — categorias, opcionais, tabela de precos
- [ ] `solicitacao.go` + `solicitacao_opcional.go`
- [ ] `servico.go` + `historico.go`
- [ ] `avaliacao.go`
- [ ] `recorrencia.go`
- [ ] `notificacao.go`
- [ ] `transacao.go` + `dados_bancarios.go`

### Backend — Integracao
- [ ] Atualizar `cmd/api/main.go` para injetar repositories Postgres no lugar dos in-memory (configuravel via ENV)
- [ ] Confirmar que a logica do limite legal (LC 150/2015) consulta corretamente a `vw_servicos_por_semana`
- [ ] Confirmar que os triggers do banco (`trg_recalcular_nota`, `trg_*_atualizado`) dispensam a logica equivalente que estava no service
- [ ] Validar que todas as queries usam parametros (`$1, $2...`), nunca concatenacao de string

### Frontend
- [ ] **Nenhuma mudanca** — o frontend ja consome a API via proxy; a troca de repository e transparente
- [ ] Validar que todos os fluxos continuam funcionando apos a troca (teste manual completo)

**Criterio de conclusao:** `docker-compose up -d && go run cmd/api/main.go` funcionando com Postgres; todos os fluxos do MVP operam com persistencia real; dados sobrevivem ao restart; frontend sem alteracoes. `go test ./...` verde.

---

## Etapa 13 — Testes de Integracao, Qualidade e Documentacao da API

**Objetivo:** fechar lacunas de cobertura, rodar revisao final de qualidade e deixar a API documentada. Desde as Etapas 1-12 cada migration e repository Postgres ganha testes de integracao proprios (via `backend/internal/testutil` + testcontainers-go); esta etapa consolida o que ficou pendente, nao cria a infraestrutura do zero.

### Backend — Testes de Integracao (consolidacao)
- [ ] Preencher lacunas de cobertura nos repositories Postgres (comparar com a matriz de services — qualquer repo sem `*_test.go` com tag `integration` entra aqui)
- [ ] Confirmar que testes unitarios dos services (in-memory) continuam verdes
- [ ] Testes end-to-end dos fluxos principais: registro -> login -> completar perfil -> solicitar -> matching -> execucao -> avaliacao
- [ ] Verificar que `TESTCONTAINERS_REUSE_ENABLE=true` acelera runs locais sem causar flakes

### Backend — Qualidade
- [ ] `go vet ./...` sem alertas
- [ ] `gofmt -w .` sem mudancas
- [ ] Zero imports nao utilizados, variaveis mortas, erros ignorados com `_`
- [ ] Revisao de seguranca final: SQL parametrizado, validacao de input, autorizacao por tipo

### Backend — Swagger
- [ ] Anotar todos os handlers com comentarios `swaggo` (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security`)
- [ ] Rodar `swag init -g cmd/api/main.go -o ../docs/swagger` para gerar documentacao
- [ ] Validar que `GET /swagger/index.html` exibe todos os endpoints documentados
- [ ] Agrupar endpoints por tags: Auth, Clientes, Profissionais, Solicitacoes, Servicos, Avaliacoes, Recorrencias, Notificacoes, Admin, Financeiro

### Frontend — Qualidade
- [ ] `npm run check` com 0 erros
- [ ] Sem `any` implicito ou cast forcado
- [ ] Sem `console.log` com token, senha ou dados pessoais
- [ ] Todos os formularios com `maxlength`, `novalidate`, validacao JS
- [ ] Todas as rotas protegidas com guard em `+page.ts`
- [ ] Responsividade testada em viewport mobile (375px) e desktop (1440px)

### Frontend — Testes
- [ ] Testes unitarios dos stores (auth, toasts) com Vitest
- [ ] Testes de componentes criticos (formularios, avaliacao) com Testing Library
- [ ] Testes E2E dos fluxos principais com Playwright: registro -> login -> dashboard, solicitacao, avaliacao

### Revisao Final (Codex + Senior)
- [ ] Revisao de codigo completa seguindo checklist de `CLAUDE.md` (Revisao Codex + Senior Backend + Senior Frontend + Senior Infra)
- [ ] Todos os erros de compilacao e teste resolvidos
- [ ] Security headers verificados no backend E no frontend
- [ ] Rate limiting ativo em todos os endpoints de autenticacao
- [ ] CSP restritiva no frontend

**Criterio de conclusao:** `go test ./...` verde; repositories testados contra banco real; todos os endpoints com autorizacao correta; Swagger UI acessivel com todos os endpoints documentados; frontend com `npm run check` limpo e testes E2E passando; revisao de seguranca aprovada.

---

## Pendentes para Fases Futuras (fora do MVP)

| Feature | Fase |
|---|---|
| Gateway de pagamento (cartao/Pix, escrow, repasse) | 2 |
| Chat in-app entre cliente e profissional | 3 |
| Programa de indicacao | 3 |
| Planos premium para profissionais | 3 |
| Multas automaticas por cancelamento tardio | 3 |
| Antecipacao de recebiveis | 4 |
| Integracao contabil | 4 |
| Assistencia residencial por assinatura | 4 |
| Parcerias com condominios | 4 |
| App mobile nativo (React Native / Flutter) | 4 |

---

*Plano atualizado em 13/abril/2026. Seguir a ordem das etapas — cada uma constroi sobre a anterior. Cada etapa entrega funcionalidade completa: backend + frontend + testes.*
