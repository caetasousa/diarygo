# DiaryGo — Plano de Implementação

> Referência: `README.md` (regras de negócio) · `migrations/` (SQL do banco)
> Cada etapa deve ser concluída e testada antes de iniciar a próxima.

## Git Flow

O projeto segue o modelo Gitflow. Branches permanentes:

| Branch | Papel |
|---|---|
| `master` | Histórico oficial de entregas (produção). Nunca recebe commits diretos. |
| `developer` | Integração contínua. Base para todas as features. |

### Branches por tipo

```
feature/<nome>   → criada a partir de developer, merge de volta para developer
release/X.Y.Z    → criada a partir de developer quando um conjunto de features está pronto; merge para master e developer
hotfix/<nome>    → criada a partir de master para correções urgentes em produção; merge para master e developer
```

### Mapeamento Etapa → Feature Branch

Cada etapa do plano corresponde a uma branch feature:

| Etapa | Branch |
|---|---|
| 0 | `feature/estrutura-base` |
| 1 | `feature/autenticacao` |
| 2 | `feature/cadastro-clientes-profissionais` |
| 3 | `feature/catalogo-precificacao` |
| 4 | `feature/solicitacoes` |
| 5 | `feature/matching-atribuicao` |
| 6 | `feature/execucao-checkin-checkout` |
| 7 | `feature/avaliacoes-reputacao` |
| 8 | `feature/recorrencias` |
| 9 | `feature/notificacoes` |
| 10 | `feature/painel-admin` |
| 11 | `feature/historico-transacoes` |
| 12 | `feature/persistencia-postgres` |
| 13 | `feature/testes-integracao-swagger` |

Quando todas as features de uma release estiverem na `developer`, abrir `release/X.Y.Z` para preparar a entrega.

### Fluxo resumido

```bash
# Iniciar feature
git checkout developer
git pull --rebase
git checkout -b feature/<nome>

# Commitar durante o desenvolvimento
git add <arquivos>
git commit -m "feat: descrição sucinta"
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

### Convenção de commits

```
feat: adiciona endpoint de login
fix: corrige cálculo de horas na limpeza express
chore: atualiza dependências
test: adiciona testes de matching com limite LC 150/2015
docs: atualiza README com esquema do banco
```

---

## Estratégia: In-Memory First

Toda a lógica de negócio é desenvolvida com repositories **in-memory** (Etapas 0–11).
O banco de dados é introduzido apenas na **Etapa 12**, substituindo as implementações in-memory por Postgres — sem tocar nos services ou handlers.

```
internal/domain/          ← interfaces de repository (contrato)
internal/repository/
  memory/                 ← implementações in-memory (Etapas 0–11)
  postgres/               ← implementações Postgres  (Etapa 12)
```

O `main.go` decide qual injetar:
```go
// Etapas 0–11
repo := memory.NewClienteRepository()

// Etapa 12 (troca uma linha)
repo := postgres.NewClienteRepository(db)
```

Vantagens: sem Docker obrigatório para desenvolver, testes unitários triviais, services nunca mudam ao trocar o banco.

---

## Etapa 0 — Estrutura do Projeto (sem banco)

**Objetivo:** esqueleto da aplicação funcionando com servidor HTTP, sem dependência de banco de dados.

- [ ] Criar `go.mod` e adicionar dependências iniciais: `go-chi/chi/v5`, `go-chi/docgen`, `swaggo/swag`, JWT, uuid
- [ ] Criar estrutura de diretórios: `cmd/api/`, `internal/{domain,handler,service,repository/memory,repository/postgres,middleware}`, `pkg/`, `migrations/`, `config/`, `docs/swagger/`
- [ ] Criar `cmd/api/main.go` com servidor HTTP usando chi, montando rota `GET /swagger/*` (swagger UI) usando repositories in-memory
- [ ] Criar `config/app.env` com variáveis de ambiente (porta, JWT secret)
- [ ] Criar `V1__schema_inicial.sql` em `migrations/` com o SQL completo do banco (pronto para a Etapa 12)
- [ ] Criar `docker-compose.yml` com PostgreSQL 14, pgAdmin e Flyway (disponível, mas não obrigatório até Etapa 12)

**Critério de conclusão:** `go run cmd/api/main.go` sobe o servidor sem erros; nenhuma dependência de Docker ou banco para rodar.

---

## Etapa 1 — Autenticação e Usuários

**Objetivo:** registro e login funcional para os três tipos de usuário. Dados em memória.

### Domínio
- [x] `internal/domain/usuario.go` — struct `Usuario`, tipos (`CLIENTE`, `PROFISSIONAL`, `ADMIN`), erros de domínio, validações (email, senha NIST 800-63b, tipo)
- [x] `internal/domain/auth.go` — interface `UsuarioRepository` (composição Reader+Writer), DTOs (`RegistroRequest`, `LoginRequest`, `VerificarEmailRequest`), `TokenPayload` com claims JWT

### Repository (in-memory)
- [x] `internal/repository/memory/usuario.go` — implementa `UsuarioRepository` com `map[uuid.UUID]*Usuario` + `map[string]uuid.UUID` (índice por email) + `sync.RWMutex`, retorna cópias para evitar mutação

### Config
- [x] `internal/config/config.go` — carrega variáveis de ambiente, valida JWT_SECRET >= 32 chars (warning em dev, fatal em prod)

### Service
- [x] `internal/service/auth_service.go` — `Registrar`, `Login`, `GerarToken` (JWT HS256), `ValidarToken` (bloqueia alg:none), `VerificarEmail`
- [x] Proteção contra timing attack no Login (bcrypt nos dois caminhos — OWASP A07)
- [x] Custo bcrypt 12 em produção, `BcryptCostTeste` (4) para testes

### Handler
- [x] `POST /auth/registro/cliente` — cria usuário tipo CLIENTE
- [x] `POST /auth/registro/profissional` — cria usuário tipo PROFISSIONAL
- [x] `POST /auth/login` — retorna JWT (Bearer token, 15 min)
- [x] `POST /auth/verificar-email` — ativa email com código (funcional, mas sem envio real de email)
- [x] `internal/handler/response.go` — helpers `RespostaJSON` e `RespostaErro`
- [x] MaxBytesReader(1MB) em todos os handlers (OWASP A06)
- [x] Mapeamento de erros de domínio para status HTTP corretos

### Middleware
- [x] `internal/middleware/auth.go` — `Autenticar` (JWT), `RequererTipo` (autorização), `UsuarioDoContexto` (helper)

### Integração
- [x] `cmd/api/main.go` — DI completa (repo → service → handler), slog estruturado, httprate 10/min em /auth/*, security headers OWASP, rota protegida /api/v1/me

### Testes
- [x] `internal/service/auth_service_test.go` — 16 testes unitários (registro, login, token, verificação email)
- [x] `internal/handler/auth_handler_test.go` — 11 testes HTTP (registro, login, verificação, rotas protegidas)
- [x] `internal/middleware/auth_test.go` — 8 testes de middleware (autenticação, autorização, contexto)

### Pendente: Ajustes de Verificação de Email (Opção B)

Decisão: o campo `EmailVerificado` existe na struct `Usuario` para uso futuro, mas **o login NÃO exige** `email_verificado = true` no MVP. A verificação real de email será implementada na **Etapa 9 (Notificações)**, quando houver infra de envio de email.

- [x] Remover `CodigoVerificacao` do `RegistroResponse` (não retornar código fake na resposta)
- [x] Remover geração de `CodigoVerificacao` do método `Registrar` no service
- [x] Remover endpoint `POST /auth/verificar-email` e método `VerificarEmail` do service
- [x] Remover `VerificarEmailRequest` e `ErrCodigoVerificacaoInvalido` do domínio
- [x] Remover campo `CodigoVerificacao` da struct `Usuario` (não faz sentido sem envio real)
- [x] Atualizar testes para refletir a remoção (remover testes de verificação de email)
- [x] Garantir que `Login` **NÃO** verifica `EmailVerificado` (já é o comportamento atual — apenas documentar)

### Pendente: Recuperação de Senha

Fluxo de recuperação de senha para o MVP (sem envio de email real — em dev, retorna token na resposta; em prod, loga o token):

- [x] `internal/domain/auth.go` — adicionar `SolicitarRecuperacaoRequest{Email}`, `SolicitarRecuperacaoResponse`, `RedefinirSenhaRequest{Token, NovaSenha}`
- [x] `internal/domain/usuario.go` — adicionar `ErrTokenRecuperacaoInvalido` e `ErrTokenRecuperacaoExpirado`
- [x] `internal/domain/usuario.go` — adicionar campos `TokenRecuperacao` e `TokenRecuperacaoExpira` na struct `Usuario`
- [x] `internal/domain/auth.go` — adicionar `BuscarPorTokenRecuperacao` na interface `UsuarioReader`
- [x] `internal/repository/memory/usuario.go` — implementar `BuscarPorTokenRecuperacao`
- [x] `internal/service/auth_service.go` — `SolicitarRecuperacao(ctx, req)`: gera token UUID, salva no usuário com expiração (1h), retorna token em dev
- [x] `internal/service/auth_service.go` — `RedefinirSenha(ctx, req)`: valida token + expiração, hash nova senha com bcrypt, limpa token
- [x] `POST /auth/solicitar-recuperacao-senha` — solicita token de recuperação
- [x] `POST /auth/redefinir-senha` — redefine senha com token válido
- [x] Testes unitários: solicitação (email existente/inexistente), redefinição, token expirado, token inválido, token vazio, nova senha fraca
- [x] Testes HTTP: endpoints de recuperação com cenários de sucesso e erro

**Critério de conclusão:** testes unitários do service e testes HTTP cobrindo registro, login, recuperação de senha, token inválido e rota protegida — tudo sem banco. Login NÃO exige email verificado.

---

## Etapa 2 — Cadastro Completo de Clientes e Profissionais

**Objetivo:** fluxos completos de cadastro com todas as entidades relacionadas. Dados em memória.

### Cliente
- [ ] `internal/domain/cliente.go` — struct, interface `ClienteRepository`, validações (CPF, score)
- [ ] `internal/repository/memory/cliente.go` — implementação in-memory
- [ ] `internal/service/cliente_service.go` — criar, atualizar perfil, buscar
- [ ] `GET/PUT /clientes/me` — perfil do cliente autenticado

### Endereços do Cliente
- [ ] `internal/domain/endereco.go` — struct, interface `EnderecoRepository`, validação de CEP
- [ ] `internal/repository/memory/endereco.go` — implementação in-memory
- [ ] `internal/service/endereco_service.go` — criar, listar, definir principal
- [ ] `POST /clientes/me/enderecos` · `GET /clientes/me/enderecos` · `PUT /clientes/me/enderecos/:id`

### Profissional
- [ ] `internal/domain/profissional.go` — struct, interface `ProfissionalRepository`, enum de status, validações
- [ ] `internal/repository/memory/profissional.go` — implementação in-memory com busca por região/nota
- [ ] `internal/service/profissional_service.go` — criar, atualizar, buscar, calcular nota
- [ ] `GET/PUT /profissionais/me` — perfil da profissional autenticada

### Documentos e Referências
- [ ] `internal/domain/documento.go` · `internal/domain/referencia.go` — structs e interfaces
- [ ] `internal/repository/memory/documento.go` · `internal/repository/memory/referencia.go`
- [ ] `internal/service/credenciamento_service.go` — orquestra aprovação, upload, validação de referências
- [ ] `POST /profissionais/me/documentos` · `POST /profissionais/me/referencias`

### Regiões e Disponibilidade
- [ ] `internal/domain/regiao.go` · `internal/domain/disponibilidade.go` — structs e interfaces
- [ ] `internal/repository/memory/regiao.go` · `internal/repository/memory/disponibilidade.go`
- [ ] `PUT /profissionais/me/regioes` · `PUT /profissionais/me/disponibilidades`

**Critério de conclusão:** profissional consegue se cadastrar completamente; status vai para PENDENTE; campos obrigatórios validados; tudo sem banco.

---

## Etapa 3 — Catálogo: Categorias, Opcionais e Tabela de Preços

**Objetivo:** catálogo de serviços e precificação por região em memória, com seed dos dados padrão.

- [ ] `internal/domain/catalogo.go` — structs `CategoriaServico`, `Opcional`, `TabelaPrecos` + interfaces de repository
- [ ] `internal/repository/memory/catalogo.go` — implementação in-memory com seed dos 7 tipos de serviço e 7 add-ons padrão
- [ ] `internal/service/precificacao_service.go` — `CalcularValorReferencia(solicitacao)` retorna valor estimado e duração
- [ ] `GET /categorias` — lista categorias ativas
- [ ] `GET /categorias/:id/opcionais` — lista opcionais disponíveis
- [ ] `GET /regioes` — lista regiões ativas
- [ ] `GET /precos?categoria=&regiao=` — consulta tabela de preços

**Critério de conclusão:** serviço de precificação calcula corretamente dado cômodos, opcionais e frequência; testado com casos dos 7 tipos de serviço; seed carregado na inicialização do `main.go`.

---

## Etapa 4 — Solicitações de Serviço

**Objetivo:** cliente consegue fazer uma solicitação completa com cálculo automático. Dados em memória.

- [ ] `internal/domain/solicitacao.go` — struct, interface `SolicitacaoRepository`, enum de status, validações (24h antecedência, etc.)
- [ ] `internal/repository/memory/solicitacao.go` — implementação in-memory
- [ ] `internal/service/solicitacao_service.go`:
  - `Criar` — valida dados, calcula valor de referência, persiste
  - `Cancelar` — valida regra das 24h, aplica penalização no score se cabível
  - `BuscarParaCliente`
- [ ] `POST /solicitacoes` — criar nova solicitação
- [ ] `GET /solicitacoes` — listar solicitações do cliente
- [ ] `GET /solicitacoes/:id` — detalhe
- [ ] `DELETE /solicitacoes/:id` — cancelar

**Critério de conclusão:** solicitação criada com valor calculado automaticamente; cancelamento com menos de 24h desconta score; teste cobre criação, validações e cancelamento.

---

## Etapa 5 — Atribuição de Profissional (Motor de Matching)

**Objetivo:** quando uma solicitação é criada, o sistema encontra e notifica a melhor profissional disponível. Dados em memória.

- [ ] `internal/domain/servico.go` — struct `Servico`, interface `ServicoRepository`
- [ ] `internal/repository/memory/servico.go` — implementação in-memory
- [ ] `internal/service/matching_service.go`:
  - `BuscarCandidatas(solicitacao)` — filtra por região, disponibilidade, categoria, nota mínima, verifica LC 150/2015
  - `AtribuirProfissional(solicitacao, profissional)` — cria o `Servico`, atualiza status da solicitação
  - `RedireccionarParaProxima(servico)` — chamado após timeout de 30 min ou recusa
- [ ] Lógica do limite legal implementada no service (consulta in-memory equivalente à `vw_servicos_por_semana`): bloqueia se profissional já tem 2 serviços na semana no mesmo endereço
- [ ] `POST /solicitacoes/:id/aceitar` — profissional aceita
- [ ] `POST /solicitacoes/:id/recusar` — profissional recusa

**Critério de conclusão:** testes unitários cobrem matching com profissional disponível, rotação por recusa, bloqueio pela regra LC 150/2015 e timeout de 30 min — sem banco.

---

## Etapa 6 — Execução do Serviço (Check-in / Check-out)

**Objetivo:** controle da execução em campo com geolocalização. Dados em memória.

- [ ] `internal/domain/historico.go` — struct `HistoricoStatus`, interface `HistoricoRepository`
- [ ] `internal/repository/memory/historico.go` — implementação in-memory
- [ ] `internal/service/execucao_service.go`:
  - `CheckIn(servicoID, lat, lon)` — valida que é a profissional certa, registra horário e localização
  - `CheckOut(servicoID, lat, lon)` — registra conclusão
  - `ConfirmarConclusao(servicoID, atorTipo)` — cliente ou profissional confirma
- [ ] `internal/service/historico_service.go` — registra mudanças de status em memória
- [ ] `POST /servicos/:id/checkin` · `POST /servicos/:id/checkout`
- [ ] `POST /servicos/:id/confirmar` — confirmação de conclusão

**Critério de conclusão:** serviço muda de AGENDADO → EM_ANDAMENTO → CONCLUIDO com timestamps e coordenadas registradas; histórico auditável; testes sem banco.

---

## Etapa 7 — Avaliações e Reputação

**Objetivo:** avaliação mútua pós-serviço e atualização automática de nota. Dados em memória.

- [ ] `internal/domain/avaliacao.go` — structs `AvaliacaoCliente`, `AvaliacaoProfissional`, interfaces de repository, validação de nota 1–5
- [ ] `internal/repository/memory/avaliacao.go` — implementação in-memory
- [ ] `internal/service/avaliacao_service.go`:
  - `AvaliarProfissional` — persiste, recalcula `nota_media` no service (equivalente ao trigger do banco)
  - `AvaliarCliente` — persiste (interna)
  - `VerificarAlertas` — dispara alerta se nota < 4.0, suspensão se nota < 3.5 por 3 consecutivos
- [ ] `POST /servicos/:id/avaliacoes/profissional`
- [ ] `POST /servicos/:id/avaliacoes/cliente`
- [ ] `GET /profissionais/:id/avaliacoes` — avaliações públicas

**Critério de conclusão:** avaliação inserida → nota recalculada no service; suspensão disparada corretamente nos cenários de nota < 3.5 por 3 consecutivos; testes unitários sem banco.

---

## Etapa 8 — Recorrências

**Objetivo:** suporte a serviços recorrentes (semanal, quinzenal, 2x/semana). Dados em memória.

- [ ] `internal/domain/recorrencia.go` — struct, interface `RecorrenciaRepository`, frequências, cálculo de próxima data
- [ ] `internal/repository/memory/recorrencia.go` — implementação in-memory
- [ ] `internal/service/recorrencia_service.go`:
  - `Criar` — baseado em solicitação confirmada
  - `GerarProximasSolicitacoes` — job agendado que cria solicitações futuras
  - `Pausar` · `Cancelar`
- [ ] `GET /recorrencias` · `POST /recorrencias/:id/pausar` · `DELETE /recorrencias/:id`

**Critério de conclusão:** recorrência semanal gera solicitações corretamente; pausa e cancelamento funcionam; rotação de profissional mantém-se dentro do limite da LC 150/2015.

---

## Etapa 9 — Notificações e Verificação Real de Email

**Objetivo:** sistema de notificações multicanal e ativação da verificação real de email. Dados em memória.

### Notificações
- [ ] `internal/domain/notificacao.go` — struct, interface `NotificacaoRepository`, canais (PUSH, EMAIL, SMS, IN_APP)
- [ ] `internal/repository/memory/notificacao.go` — implementação in-memory
- [ ] `internal/service/notificacao_service.go` — interface `Notificador` + implementações stub (MVP envia apenas e-mail/in-app)
- [ ] Integrar notificações nos eventos: solicitação criada, profissional atribuída, véspera de serviço, check-in, conclusão, avaliação recebida
- [ ] `GET /notificacoes` · `PUT /notificacoes/:id/lida`

### Verificação Real de Email (adiada da Etapa 1)
- [ ] Implementar envio real de email com código de verificação (usando o `Notificador`)
- [ ] Recriar `POST /auth/verificar-email` com código enviado por email (não retornado na resposta)
- [ ] Adicionar `CodigoVerificacao` e `CodigoVerificacaoExpira` na struct `Usuario`
- [ ] Opcionalmente ativar bloqueio de login para emails não verificados (decisão a tomar na etapa)
- [ ] Integrar envio real de email na recuperação de senha (substituir o retorno do token na resposta)

**Critério de conclusão:** eventos principais disparam notificações; in-app funcionando; verificação de email funcional com envio real; stubs prontos para substituir por push real na Fase 2.

---

## Etapa 10 — Painel Admin

**Objetivo:** painel básico de administração para gestão da plataforma.

- [ ] `internal/service/admin_service.go` — aprovação/reprovação de profissionais, suspensão, gestão de preços
- [ ] `GET /admin/profissionais` — lista com filtro por status
- [ ] `PUT /admin/profissionais/:id/aprovar` · `PUT /admin/profissionais/:id/reprovar`
- [ ] `PUT /admin/profissionais/:id/suspender` · `PUT /admin/profissionais/:id/descredenciar`
- [ ] `GET /admin/servicos` — lista com filtros (status, data, profissional, cliente)
- [ ] `GET /admin/clientes`
- [ ] `POST /admin/tabela-precos` · `PUT /admin/tabela-precos/:id`
- [ ] `GET /admin/regioes` · `POST /admin/regioes`

**Critério de conclusão:** admin consegue aprovar profissional → status muda → profissional aparece no matching; configuração de preços por região refletida no cálculo de solicitações.

---

## Etapa 11 — Histórico e Transações (Estrutura MVP)

**Objetivo:** registrar transações financeiras como referência (sem cobrança online). Dados em memória.

- [ ] `internal/domain/transacao.go` — struct, interface `TransacaoRepository`
- [ ] `internal/repository/memory/transacao.go` — implementação in-memory
- [ ] `internal/service/transacao_service.go` — `Registrar(servico)` cria transação com status REGISTRADA e metodo DIRETO_EXTERNO
- [ ] `GET /clientes/me/historico` — histórico de serviços com valores
- [ ] `GET /profissionais/me/historico` — histórico com ganhos de referência
- [ ] `GET /admin/transacoes` — visão financeira geral

**Critério de conclusão:** toda conclusão de serviço gera registro de transação; histórico navegável por cliente e profissional; tudo sem banco.

---

## Etapa 12 — Persistência em PostgreSQL (substituição dos repositories)

**Objetivo:** substituir todas as implementações in-memory por implementações Postgres. Services e handlers não mudam.

- [ ] Subir infraestrutura: `docker-compose up -d` → PostgreSQL + pgAdmin + Flyway aplicando `V1__schema_inicial.sql`
- [ ] Adicionar dependência `pgx/v5` ao `go.mod`
- [ ] Criar `internal/repository/postgres/db.go` — pool de conexões com `pgxpool`
- [ ] Para cada repository in-memory, criar o equivalente em `internal/repository/postgres/`:
  - [ ] `usuario.go` · `cliente.go` · `profissional.go` · `admin.go`
  - [ ] `endereco.go` · `regiao.go` · `disponibilidade.go`
  - [ ] `documento.go` · `referencia.go`
  - [ ] `catalogo.go` (categorias, opcionais, tabela de preços)
  - [ ] `solicitacao.go` · `solicitacao_opcional.go`
  - [ ] `servico.go` · `historico.go`
  - [ ] `avaliacao.go` · `recorrencia.go`
  - [ ] `notificacao.go` · `transacao.go` · `dados_bancarios.go`
- [ ] Atualizar `cmd/api/main.go` para injetar repositories Postgres no lugar dos in-memory
- [ ] Confirmar que a lógica do limite legal (LC 150/2015) consulta corretamente a `vw_servicos_por_semana`
- [ ] Confirmar que os triggers do banco (`trg_recalcular_nota`, `trg_*_atualizado`) dispensam a lógica equivalente que estava no service

**Critério de conclusão:** `docker-compose up -d && go run cmd/api/main.go` funcionando com Postgres; todos os fluxos do MVP operam com persistência real; dados sobrevivem ao restart.

---

## Etapa 13 — Testes de Integração, Qualidade e Documentação da API

**Objetivo:** cobertura de testes com banco real e API documentada.

- [ ] Testes de integração para todos os repositories Postgres (banco real, sem mocks)
- [ ] Confirmar que testes unitários dos services (in-memory) continuam verdes
- [ ] `go vet ./...` e `gofmt -w .` sem alertas
- [ ] Anotar todos os handlers com comentários `swaggo` (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Failure`, `@Router`)
- [ ] Rodar `swag init` para gerar `docs/swagger/` e expor em `GET /swagger/*` via chi
- [ ] Revisão de segurança: SQL injection (uso de queries parametrizadas), validação de input, autorização por tipo de usuário

**Critério de conclusão:** `go test ./...` verde; repositories testados contra banco real; todos os endpoints com autorização correta (cliente não acessa rotas de admin/profissional); Swagger UI acessível em `/swagger/index.html` com todos os endpoints documentados.

---

## Pendentes para Fases Futuras (fora do MVP)

| Feature | Fase |
|---|---|
| Gateway de pagamento (cartão/Pix, escrow, repasse) | 2 |
| Chat in-app entre cliente e profissional | 3 |
| Programa de indicação | 3 |
| Planos premium para profissionais | 3 |
| Multas automáticas por cancelamento tardio | 3 |
| Antecipação de recebíveis | 4 |
| Integração contábil | 4 |
| Assistência residencial por assinatura | 4 |
| Parcerias com condomínios | 4 |

---

*Plano atualizado em 11/abril/2026. Seguir a ordem das etapas — cada uma constrói sobre a anterior.*
