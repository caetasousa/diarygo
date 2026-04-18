# DiaryGo — Plano de Implementacao

> Referencia: `README.md` (regras de negocio) . `backend/migrations/` (SQL do banco).
> Cada etapa inclui **backend + frontend + testes** e so termina quando `go test ./...`, `go test -tags=integration ./...` e `npm run check` estao verdes.
> Escopo MVP: **somente Goiania/GO**, sem pagamento online, sem chat real-time, sem storage binario (imagens sao URL externa ate integrarmos CDN).

## Git Flow

Modelo Gitflow com duas branches permanentes:

| Branch | Papel |
|---|---|
| `master` | Historico oficial de entregas (producao). Nunca recebe commits diretos. |
| `developer` | Integracao continua. Base para todas as features. |

### Branches por tipo

```
feature/<nome>   -> criada a partir de developer, merge de volta para developer
release/X.Y.Z    -> criada a partir de developer quando um conjunto de features esta pronto; merge para master e developer
hotfix/<nome>    -> criada a partir de master para correcoes urgentes; merge para master e developer
```

### Mapeamento Etapa -> Feature Branch

| Etapa | Branch | Status |
|---|---|---|
| 0 | `feature/estrutura-base` | **Concluida** |
| 1 | `feature/autenticacao` | **Concluida** |
| 2 | `feature/cadastro-clientes-profissionais` | **Concluida** (inclui consolidacao dashboard + Goiania-first) |
| 3 | `feature/catalogo-precificacao` | Pendente |
| 4 | `feature/solicitacoes` | Pendente |
| 5 | `feature/matching-atribuicao` | Pendente |
| 6 | `feature/execucao-checkin-checkout` | Pendente |
| 7 | `feature/avaliacoes-reputacao` | Pendente |
| 8 | `feature/recorrencias` | Pendente |
| 9 | `feature/notificacoes` | Pendente (inclui verificacao real de email) |
| 10 | `feature/painel-admin` | Pendente |
| 11 | `feature/historico-transacoes` | Pendente |
| 12 | `feature/persistencia-postgres` | Pendente |
| 13 | `feature/finalizacao-mvp` | Pendente (chat 1-campo + recibo PDF + E2E + polish) |

### Convencao de commits

```
feat: adiciona endpoint de login
fix:  corrige calculo de horas na limpeza express
chore: atualiza dependencias
test: adiciona testes de matching com limite LC 150/2015
docs: atualiza README com esquema do banco
```

---

## Estrategia: In-Memory First

Toda a logica de negocio e desenvolvida com repositories **in-memory** (Etapas 0-11).
O banco e introduzido apenas na **Etapa 12**, substituindo as implementacoes in-memory por Postgres — sem tocar em services ou handlers.

```
internal/domain/          <- interfaces de repository (contrato)
internal/repository/
  memory/                 <- implementacoes in-memory (Etapas 0-11)
  postgres/               <- implementacoes Postgres  (Etapa 12)
```

Vantagens: sem Docker obrigatorio, testes unitarios triviais, services nunca mudam ao trocar o banco.

---

## Estrategia: Frontend Integrado por Etapa

Cada etapa entrega a funcionalidade **completa da API ate a interface**. Padroes ja estabelecidos (Etapas 0-2) sao lei para as seguintes:

- **Svelte 5** com `{#snippet}` + `{@render}`; layout via `{@render children()}`
- **Guards de rota** em `+page.ts` com `redirect(302, '/login')` — nunca em `onMount`
- **Formularios** com `novalidate`, `maxlength` em todos inputs, validacao JS + server-side
- **ApiClient singleton** tipado (`src/lib/api/client.ts`) — nunca `fetch` direto nas paginas
- **Stores Svelte** para estado global (auth, toasts)
- **CSP restritiva** em `app.html` com libera apenas origens conhecidas (ViaCEP, BrasilAPI, imagens `https:`)
- **Try/catch** em todo fetch, mensagem amigavel ao usuario, nunca expor stack trace
- Componentes reutilizaveis em `src/lib/components/` (hoje: Avatar, ImagePreview, StatusBadge, CepInput, TelefoneInput, CpfInput, EnderecoCard, Toast, Navbar, GradeSemanal, RegiaoSelector)

---

## Escopo MVP: Decisoes de Arquitetura

Decisoes ja tomadas e aplicadas, que guiam todas as proximas etapas:

### Goiania-first
- CEP validado em faixa **74000000..74999999** no service; UF fixa `GO`; cidade normalizada para "goiania" (tolerante a acento/case).
- Seed in-memory carrega as **12 regioes administrativas oficiais de Goiania** (Central, Norte, Sul, Sudoeste, Oeste, Noroeste, Campinas-Centro, Macambira, Leste, Vale do Meia Ponte, Sudeste, Mendanha).
- Formulario de endereco pre-seta `Goiania/GO`. CepInput avisa em amarelo quando CEP nao comeca com `74`.
- Arquitetura preparada para expansao: basta parametrizar `CEPGoianiaInicio`/`CEPGoianiaFim` para atender outras cidades.

### Sem storage binario
- **Nenhum campo de imagem envia blob/base64 ao backend**. Foto de perfil da profissional e (futuramente) referencias com foto usam **URL externa** (Drive, Dropbox, Imgur).
- Frontend renderiza preview via `ImagePreview`/`Avatar`. Contrato do backend nao muda quando a CDN for integrada — so o componente de input vira uploader.
- Descartado no MVP: OCR de documento, upload binario, storage no Postgres.

### Dashboard consolidado
- **Profissional** tem 3 sub-rotas: `/dashboard/perfil` (dados pessoais + foto URL) | `/dashboard/credenciamento` (referencias) | `/dashboard/atuacao` (regioes + disponibilidade).
- **Cliente** tem 2 sub-rotas: `/dashboard/perfil` | `/dashboard/enderecos`.
- Sidebar com cartao do usuario + logout no rodape. Topbar removida (busca fake e notificacoes vazias nao agregavam).
- Rotas antigas (`/documentos`, `/referencias`, `/regioes`, `/disponibilidade`) redirecionam 302 para preservar bookmarks.

### Credenciamento sem upload
- Dados pessoais da profissional (Nome, CPF, RG, Telefone, MEI) sao **digitados** em `/dashboard/perfil`. Nao pedimos foto do documento.
- Credenciamento = **referencias profissionais** (nome + telefone de contato), no minimo 2 confirmadas para aprovacao.
- Admin aprova manualmente via painel (Etapa 10).

---

## Diferenciais Competitivos (integrados nas etapas)

Tres mecanicas que elevam o produto acima do padrao GetNinjas/Parafuzo, embutidas nas etapas relevantes ao inves de virarem feature separada:

### (A) Score composto de ranking — integrado na Etapa 5

Matching nao decide so por nota media. Score por candidata:

```
score = 0.35 * nota_media_normalizada      // 0..1 onde 5.0 -> 1.0
      + 0.20 * taxa_aceite                 // aceites / total_atribuicoes
      + 0.15 * (1 - taxa_no_show)          // penaliza faltas
      + 0.15 * proximidade_regiao          // match exato = 1, fora = 0
      + 0.15 * bonus_relacionamento        // cliente ja teve servico concluido com ela
```

Pesos expostos como constantes no `matching_service.go` para ajuste fino. Empate -> quem atendeu menos o cliente recentemente (fairness entre profissionais).

### (B) Favoritas e bloqueios — integrado na Etapa 2 (dominio) + Etapa 5 (matching)

Cliente pode:
- **Favoritar** profissional apos servico concluido. Em proxima solicitacao, matching tenta favoritas primeiro (se disponiveis na data/regiao). Se nenhuma disponivel, cai no score composto normal.
- **Bloquear** profissional. Matching nunca atribui uma bloqueada a esse cliente.

Tabela `cliente_profissional_preferencia` (cliente_id, profissional_id, tipo FAVORITA|BLOQUEADA, criado_em).

### (C) Orcamento transparente — integrado na Etapa 3

Calculo exposto como **breakdown** no frontend, nao caixa-preta:

```
Base regional (Região Sul):     R$ 25,00/h
Comodos (3 quartos, 2 banh.):   +R$ 15,00
Opcionais (passadoria):         +R$ 20,00
Frequencia (semanal):           -10%
                                ────────
Estimativa:                     R$ 162,00
Duracao:                        4h30
```

DTO `CalculoPrecoResponse` inclui array `itens[]` com `{label, valor, tipo}`. Frontend renderiza como recibo.

---

## Etapa 0 — Estrutura do Projeto [CONCLUIDA]

Esqueleto Go + chi + SvelteKit com security headers (OWASP A02), rate limiting, `slog` estruturado, design system CSS, ApiClient singleton, proxy `/api` -> `:8080`. Docker Compose disponivel (nao obrigatorio ate Etapa 12). Migration `V1__auth_usuarios.sql` criada.

---

## Etapa 1 — Autenticacao e Usuarios [CONCLUIDA]

Registro/login com JWT HS256, bcrypt cost 12 (4 em teste), protecao contra timing attack, recuperacao de senha sem enumeracao, middleware `Autenticar` + `RequererTipo`, 35 testes unitarios + HTTP cobrindo todos os fluxos. Frontend com todas as paginas de auth + dashboard role-based. **Verificacao real de email mantida na Etapa 9** (quando houver infra de envio).

---

## Etapa 2 — Cadastro Completo de Clientes e Profissionais [CONCLUIDA]

### Entregue

**Backend:**
- Domain + repository + service + handler para: Cliente, Endereco, Profissional, Documento, Referencia, Regiao, Disponibilidade, ProfissionalRegiao
- Validacoes Goiania-first: `ValidarAreaAtendimento(cep, cidade, estado)` no service de endereco
- Seed das 12 regioes administrativas oficiais de Goiania
- Endpoints autenticados: `GET/POST/PUT /clientes/me`, CRUD `/clientes/me/enderecos`, `GET/POST/PUT /profissionais/me`, CRUD `/profissionais/me/{documentos,referencias,regioes,disponibilidades}`, publica `GET /regioes`
- Migration `V2__cadastro_clientes_profissionais.sql` (pronta para Etapa 12)
- Testes unitarios completos dos services + infra de integracao com testcontainers-go

**Frontend:**
- `/dashboard/perfil` com Avatar + ImagePreview + stats do profissional
- `/dashboard/enderecos` com CepInput (ViaCEP + fallback BrasilAPI), pre-set Goiania, validacao de Bairro e Numero obrigatorios, detalhes de imovel em branco por padrao
- `/dashboard/credenciamento` consolida so Referencias (sem upload de documento)
- `/dashboard/atuacao` consolida Regioes + Disponibilidade com "Salvar alteracoes" unico
- Componentes: `Avatar`, `ImagePreview`, `StatusBadge`, `CepInput` com fallback e aviso fora-de-Goiania
- Sidebar com cartao do usuario no rodape + logout; topbar removida
- CSP libera ViaCEP, BrasilAPI, imagens `https:`

### Pendencias que foram absorvidas em etapas posteriores

- **Documentos com imagem** -> descartado no MVP (vira URL so quando CDN entrar, em etapa futura fora do plano atual)
- **Favoritas/bloqueios de cliente** -> Etapa 2 dominio + Etapa 5 matching (ver diferencial (B))

### Nova sub-tarefa residual (puxar antes da Etapa 3)

- [ ] `internal/domain/preferencia.go` — struct `ClienteProfissionalPreferencia` (ClienteID, ProfissionalID, Tipo: FAVORITA|BLOQUEADA, CriadoEm) + interface `PreferenciaRepository`
- [ ] `internal/repository/memory/preferencia.go` — com indices por ClienteID
- [ ] `internal/service/preferencia_service.go` — `Favoritar`, `Bloquear`, `Remover`, `ListarFavoritas`, `ListarBloqueadas`
- [ ] Endpoints: `POST /clientes/me/favoritas/:profissionalId`, `DELETE /clientes/me/favoritas/:profissionalId`, `POST /clientes/me/bloqueios/:profissionalId`, `DELETE /clientes/me/bloqueios/:profissionalId`, `GET /clientes/me/favoritas`, `GET /clientes/me/bloqueios`
- [ ] Migration `V3__preferencias_cliente.sql` (ainda nao aplicada — so quando chegar a Etapa 12)
- [ ] Frontend: apenas persistir o dominio. UI de favoritar/bloquear entra no final da Etapa 6 (apos conclusao de servico)

---

## Etapa 3 — Catalogo: Categorias, Opcionais, Precos e Orcamento Transparente

**Objetivo:** catalogo de servicos e precificacao por regiao com **breakdown do calculo** exibido ao cliente.

### Backend
- [ ] `internal/domain/catalogo.go` — `CategoriaServico`, `Opcional`, `TabelaPrecos` + interfaces
- [ ] `internal/repository/memory/catalogo.go` — seed dos 7 tipos de servico e 7 add-ons padrao, tabela de precos por categoria × regiao (12 regioes de Goiania)
- [ ] `internal/service/precificacao_service.go`:
  - `CalcularValorReferencia(req)` retorna `CalculoPrecoResponse{valor_total, duracao, itens[]}` onde `itens[]` e o breakdown passo a passo
  - Regras: base regional, incremento por comodo, opcionais, frequencia (unica / semanal / quinzenal / 2x_semana), acrescimo fim-de-semana
- [ ] `GET /categorias` (publica) — categorias ativas
- [ ] `GET /categorias/:id/opcionais` (publica) — opcionais da categoria
- [ ] `POST /precos/calcular` (publica, body com categoria + regiao + comodos + opcionais + frequencia + data) — **este e o endpoint que alimenta a tela de orcamento transparente**
- [ ] Testes: calculo por categoria, desconto semanal/quinzenal (-10%/-5%), acrescimo fds (+15%), breakdown com itens corretos, todas as 12 regioes

### Frontend
- [ ] Home (`/`) passa a carregar categorias da API (substitui dados estaticos)
- [ ] `/servicos` — vitrine de servicos com card por categoria (icone, descricao, duracao min, "a partir de R$...")
- [ ] `/servicos/:id` — detalhe + **calculadora de preco com breakdown em tempo real** (o diferencial (C))
- [ ] Componentes: `CategoriaCard`, `OpcionalCheckbox`, `CalculadoraPreco` (com secao "Como chegamos nesse valor" listando cada item do breakdown), `RegiaoSelector` (pre-seleciona regiao pelo CEP do endereco principal se o usuario estiver logado)

**Criterio:** calculo correto para os 7 tipos × 12 regioes; frontend exibe quebra detalhada do preco; sem banco.

---

## Etapa 4 — Solicitacoes de Servico

**Objetivo:** cliente faz solicitacao completa em multi-step com valor calculado e breakdown transparente.

### Backend
- [ ] `internal/domain/solicitacao.go` — `Solicitacao` com status `AGUARDANDO|ATRIBUIDA|CONFIRMADA|EM_ANDAMENTO|CONCLUIDA|CANCELADA`
- [ ] `internal/domain/solicitacao_opcional.go` — associativa N:N
- [ ] `internal/repository/memory/solicitacao.go` — indices por `ClienteID`, `Status`
- [ ] `internal/service/solicitacao_service.go`:
  - `Criar` — valida endereco do cliente, regra de **24h de antecedencia**, calcula valor via `precificacao_service`, persiste com snapshot do breakdown
  - `Cancelar` — regra `< 24h` -> penalizacao -5 no score (dominio de Cliente)
  - `BuscarPorID`, `ListarDoCliente(filtros)`
- [ ] Endpoints: `POST /solicitacoes`, `GET /solicitacoes`, `GET /solicitacoes/:id`, `DELETE /solicitacoes/:id`
- [ ] Testes: criacao com calculo congelado, 24h bloqueia, cancelamento penaliza score

### Frontend
- [ ] `/solicitacoes/novo` — fluxo 4-step: Endereco (cards dos enderecos do cliente) → Categoria + Opcionais → Data/Horario/Frequencia → Resumo com breakdown + botao Confirmar
- [ ] `/dashboard/solicitacoes` — lista filtravel por status; timeline visual do status
- [ ] `/dashboard/solicitacoes/:id` — detalhe com breakdown congelado + opcao cancelar
- [ ] Componentes: `StepIndicator`, `EnderecoSelector` (radio de cards), `CalendarioAgendamento` (bloqueia < 24h), `ResumoSolicitacao` (reutiliza breakdown da Etapa 3), `ModalConfirmacao` (avisa penalizacao se < 24h)

**Criterio:** fluxo end-to-end do cliente — cadastra endereco em /dashboard/enderecos, solicita em /solicitacoes/novo com valor e breakdown corretos, recebe na lista, cancela com aviso.

---

## Etapa 5 — Matching com Score Composto + Preferencias

**Objetivo:** motor de atribuicao que prioriza favoritas, respeita bloqueios, aplica score composto e enforce LC 150/2015.

### Backend
- [ ] `internal/domain/servico.go` — `Servico` com status `AGENDADO|EM_ANDAMENTO|CONCLUIDO|CANCELADO|NO_SHOW`, campos de checkin/checkout (lat/lon/horario)
- [ ] `internal/repository/memory/servico.go` — indices por `ProfissionalID`, `SolicitacaoID`, `ClienteID`
- [ ] `internal/service/matching_service.go`:
  - `BuscarCandidatas(solicitacao)` — filtra por regiao, disponibilidade no dia/horario, categoria habilitada, status `APROVADA`, remove bloqueadas pelo cliente, respeita **limite LC 150/2015** (max 2 servicos/semana da mesma profissional no mesmo endereco)
  - `RankearCandidatas(candidatas, cliente)` — **score composto** (diferencial A) com pesos em constantes do package
  - `AtribuirProfissional(solicitacao)`:
    1. tenta favoritas do cliente disponiveis -> ordena por score -> melhor
    2. se vazio, rankeia todas por score -> melhor
    3. cria `Servico` com status `AGENDADO` e solicitacao vira `ATRIBUIDA`
    4. agenda timeout de 30min para aceite
  - `AceitarAtribuicao(servicoID, profissionalID)` — solicitacao vira `CONFIRMADA`
  - `RecusarAtribuicao(servicoID, profissionalID)` — chama `RedireccionarParaProxima`
  - `RedireccionarParaProxima(solicitacao, excluir[])` — re-rankeia excluindo quem ja recusou/timeout
- [ ] Endpoints: `POST /servicos/:id/aceitar`, `POST /servicos/:id/recusar`, `GET /profissionais/me/atribuicoes-pendentes`
- [ ] Testes: favoritas preferidas, bloqueadas nunca atribuidas, LC 150/2015 bloqueia, timeout 30min rotaciona, score composto com cenarios (nota alta vs. aceite alto vs. relacionamento), tie-breaker por fairness

### Frontend
- [ ] `/dashboard/atribuicoes` (profissional) — cards de solicitacoes atribuidas aguardando aceite, com `TimerCountdown` de 30min, detalhes do servico, botoes aceitar/recusar
- [ ] Cliente ve no detalhe da solicitacao: "Buscando profissional..." -> dados da profissional atribuida (nome, nota, foto) com timeline atualizada
- [ ] Componentes: `TimerCountdown`, `SolicitacaoAtribuidaCard`

**Criterio:** cliente com favorita disponivel recebe-a preferencialmente; bloqueada nunca atribuida; nao estoura LC 150/2015; timeout rotaciona corretamente.

---

## Etapa 6 — Execucao: Check-in/Check-out + Feedback de Preferencia

**Objetivo:** controle de execucao em campo com geolocalizacao + UI de favoritar/bloquear apos conclusao.

### Backend
- [ ] `internal/domain/historico.go` — `HistoricoStatus` auditavel
- [ ] `internal/service/execucao_service.go`:
  - `CheckIn(servicoID, profissionalID, lat, lon)` — valida profissional certa e status `AGENDADO`
  - `CheckOut(servicoID, profissionalID, lat, lon)` — valida `EM_ANDAMENTO`
  - `ConfirmarConclusao(servicoID, atorID, atorTipo)` — bilateral; quando ambos confirmam, status `CONCLUIDO`
- [ ] Endpoints: `POST /servicos/:id/checkin`, `POST /servicos/:id/checkout`, `POST /servicos/:id/confirmar`, `GET /servicos/:id`, `GET /profissionais/me/servicos`, `GET /clientes/me/servicos`
- [ ] Testes: fluxo `AGENDADO -> EM_ANDAMENTO -> CONCLUIDO`, historico registrado, confirmacao bilateral

### Frontend
- [ ] `/dashboard/servicos` (profissional) — lista filtravel por status
- [ ] `/dashboard/servicos/:id` — timeline + botoes geolocalizacao:
  - `BotaoGeolocalizacao` captura `navigator.geolocation` com permissao, fallback para input manual se negada
  - Fluxo: Check-in -> Check-out -> Confirmar
- [ ] `/dashboard/servicos` (cliente) — timeline em tempo real (polling 10s enquanto em andamento) + botao Confirmar
- [ ] **Apos servico CONCLUIDO**, cliente ve na pagina do servico dois botoes: "Favoritar profissional" / "Bloquear profissional" (diferencial B — UI minima, um clique)
- [ ] Componentes: `TimelineStatus`, `BotaoGeolocalizacao`, `ServicoCard`

**Criterio:** execucao completa de ponta a ponta com geolocalizacao, cliente consegue favoritar/bloquear em um clique apos conclusao.

---

## Etapa 7 — Avaliacoes e Reputacao

**Objetivo:** avaliacao mutua pos-servico, recalculo de `nota_media` e atualizacao dos sinais usados pelo score composto (Etapa 5).

### Backend
- [ ] `internal/domain/avaliacao.go` — `AvaliacaoCliente` (pontualidade, qualidade, educacao), `AvaliacaoProfissional` (ambiente, materiais, respeito); ambas com nota geral 1-5 + comentario
- [ ] `internal/service/avaliacao_service.go`:
  - `AvaliarProfissional` — valida servico `CONCLUIDO`, atualiza `nota_media` e `total_servicos` do profissional
  - `AvaliarCliente` — interna, alimenta score de confiabilidade do cliente
  - `VerificarAlertas` — nota < 4.0 dispara alerta interno; nota < 3.5 em 3 servicos consecutivos suspende automaticamente (status `SUSPENSA`)
  - Atualiza **taxa_aceite** e **taxa_no_show** do profissional (metricas usadas pelo score composto)
- [ ] Endpoints: `POST /servicos/:id/avaliacoes/profissional` (CLIENTE), `POST /servicos/:id/avaliacoes/cliente` (PROFISSIONAL), `GET /profissionais/:id/avaliacoes` (publica, sem dados sensiveis)
- [ ] Testes: recalculo de media, suspensao automatica, metricas atualizadas corretamente

### Frontend
- [ ] Apos `CONCLUIDO`, prompt de avaliacao aparece no detalhe do servico
- [ ] Componentes: `RatingStars` (hover, click, readonly), `FormularioAvaliacao` com estrelas por criterio + comentario
- [ ] `/profissionais/:id` — pagina publica com nota media, total, avaliacoes recentes (sem CPF/email)
- [ ] Badge de alerta no dashboard da profissional se nota < 4.0

**Criterio:** 3 notas < 3.5 consecutivas -> profissional suspensa e removida do matching; score composto reflete novas metricas.

---

## Etapa 8 — Recorrencias

**Objetivo:** servicos recorrentes (semanal, quinzenal, 2x/semana) com rotacao de profissional respeitando LC 150/2015.

### Backend
- [ ] `internal/domain/recorrencia.go` — `Recorrencia` com `Frequencia`, `DiaSemana1/2`, `HoraInicio`, `ProximaData`, `Status ATIVA|PAUSADA|CANCELADA`
- [ ] `internal/service/recorrencia_service.go`:
  - `Criar` a partir de solicitacao confirmada com frequencia != UNICA
  - `GerarProximasSolicitacoes` — job que, diariamente, cria solicitacoes futuras das recorrencias ativas (14 dias a frente por padrao)
  - `Pausar`/`Retomar`/`Cancelar`
- [ ] Endpoints: `GET/POST /recorrencias`, `POST /recorrencias/:id/{pausar,retomar}`, `DELETE /recorrencias/:id`
- [ ] Testes: geracao correta, pausa nao gera, rotacao por LC 150/2015 quando mesma profissional nao pode

### Frontend
- [ ] `/dashboard/recorrencias` — lista com badge de status, frequencia, proximo servico
- [ ] `/dashboard/recorrencias/:id` — detalhe com historico + pausar/retomar/cancelar
- [ ] Componente `RecorrenciaCard`
- [ ] Fluxo de solicitacao (Etapa 4) destaca descontos da frequencia

**Criterio:** recorrencia semanal gera 2 solicitacoes futuras em 14 dias; rotacao escolhe nova profissional respeitando o limite legal.

---

## Etapa 9 — Notificacoes + Verificacao Real de Email

**Objetivo:** central de notificacoes in-app + ativacao da verificacao de email (adiada da Etapa 1).

### Backend
- [ ] `internal/domain/notificacao.go` — `Notificacao` (canal `PUSH|EMAIL|SMS|IN_APP`, titulo, corpo, status)
- [ ] `internal/service/notificacao_service.go` — interface `Notificador` + stub (MVP envia apenas IN_APP)
- [ ] Integracao de eventos:
  - Solicitacao criada -> notifica profissional atribuida
  - Profissional aceita -> notifica cliente
  - Vespera do servico -> notifica cliente (nome + foto da profissional)
  - Check-in -> notifica cliente
  - Conclusao -> notifica ambos (pedido de avaliacao)
  - Recorrencia proxima -> notifica cliente
- [ ] Endpoints: `GET /notificacoes`, `PUT /notificacoes/:id/lida`, `PUT /notificacoes/todas-lidas`, `GET /notificacoes/nao-lidas/contagem`
- [ ] Testes: evento dispara notificacao, contagem correta, marcar lida

### Backend — Verificacao Real de Email (resgatada da Etapa 1)
- [ ] Implementar envio real via `Notificador` (driver SMTP configuravel)
- [ ] Recriar `POST /auth/verificar-email` com codigo enviado — sem retornar na response
- [ ] Adicionar `CodigoVerificacao` + `CodigoVerificacaoExpira` no `Usuario`
- [ ] Recuperacao de senha passa a enviar token so por email (sem expor na response)
- [ ] Decidir ao fim da etapa: bloquear login para nao-verificados? (recomendacao: manter opcional via ENV `REQUER_EMAIL_VERIFICADO`)

### Frontend
- [ ] `/dashboard/notificacoes` — lista agrupada por data (hoje / ontem / semana / anterior) com botao "marcar todas lidas"
- [ ] Badge de nao-lidas na sidebar com polling a cada 30s
- [ ] Clique na notificacao navega para o recurso + marca como lida
- [ ] `/verificar-email` — input de codigo recebido
- [ ] Componentes: `NotificacaoItem` (com icone por tipo de evento + horario relativo), `NotificacaoBadge`

**Criterio:** eventos principais notificam em tempo real; email de verificacao funciona; recuperacao de senha nao mais expoe token.

---

## Etapa 10 — Painel Admin

**Objetivo:** interface completa de administracao — aprovar profissionais, configurar precos, monitorar servicos.

### Backend
- [ ] `internal/service/admin_service.go`:
  - Gestao de profissionais: listar (filtro por status), aprovar, reprovar (com motivo), suspender, descredenciar
  - Gestao de clientes: listar, visualizar detalhes (inclusive score e penalidades)
  - Gestao de servicos: listar com filtros (status, data, profissional, cliente)
  - CRUD tabela de precos por regiao
  - CRUD regioes (ativar/desativar)
  - CRUD categorias e opcionais
- [ ] Endpoints: todos em `/admin/*`, protegidos com `RequererTipo(ADMIN)`
  - `/admin/profissionais` (GET com filtro, `:id` GET detalhe, `:id/aprovar`, `:id/reprovar`, `:id/suspender`, `:id/descredenciar`)
  - `/admin/clientes` (GET lista, `:id` detalhe)
  - `/admin/servicos` (GET com filtros)
  - `/admin/tabela-precos` (CRUD)
  - `/admin/regioes` (CRUD)
  - `/admin/categorias` (CRUD)
- [ ] Testes: aprovacao muda status -> profissional entra no matching, preco editado reflete em novo calculo

### Frontend
- [ ] Layout dedicado `/admin` com sidebar propria: Dashboard / Profissionais / Clientes / Servicos / Precos / Regioes / Categorias
- [ ] `/admin` — dashboard com metricas (pendentes de aprovacao, servicos do dia/mes, ticket medio)
- [ ] `/admin/profissionais` — tabela com filtros e paginacao; detalhe com documentos, referencias, historico, avaliacoes, botoes de acao + modal de motivo
- [ ] `/admin/precos` — edicao inline por categoria × regiao
- [ ] Componentes: `AdminSidebar`, `TabelaAdmin` (sort + filtro + paginacao), `ModalMotivo`
- [ ] Guard `RequererTipo(ADMIN)` em todas as rotas `/admin/*`

**Criterio:** admin aprova profissional -> ela aparece no matching; edita preco -> reflete no calculo; suspende -> sai do pool.

---

## Etapa 11 — Historico e Transacoes (Estrutura MVP)

**Objetivo:** registrar transacoes financeiras como referencia (sem cobranca online) + painel financeiro.

### Backend
- [ ] `internal/domain/transacao.go` — `Transacao` (valor, metodo `DIRETO_EXTERNO`, status `REGISTRADA`)
- [ ] `internal/domain/dados_bancarios.go` — `DadosBancarios` (banco, agencia, conta, tipo, chave Pix, principal)
- [ ] `internal/service/transacao_service.go` — `Registrar(servico)` dispara na conclusao
- [ ] `internal/service/dados_bancarios_service.go` — CRUD (opcional; profissional pode cadastrar se quiser)
- [ ] Endpoints:
  - `GET /clientes/me/historico` (paginado, filtros)
  - `GET /profissionais/me/historico` (paginado, filtros)
  - `GET/POST/PUT /profissionais/me/dados-bancarios`
  - `GET /admin/transacoes`
  - `GET /admin/dashboard/financeiro` (total movimentado, por regiao, por categoria, ticket medio)
- [ ] Testes: toda conclusao gera transacao, historico correto por papel, dados bancarios CRUD

### Frontend
- [ ] `/dashboard/historico` — contexto cliente (total gasto no periodo + lista de servicos concluidos)
- [ ] `/dashboard/historico` — contexto profissional (total ganho no periodo + servicos)
- [ ] `/dashboard/dados-bancarios` (profissional) — CRUD de contas/Pix
- [ ] `/admin/financeiro` — cards de metricas + grafico simples (30 dias) + tabela filtravel
- [ ] Componentes: `HistoricoCard`, grafico CSS puro ou biblioteca leve (evitar Chart.js full se possivel — D3 minimo ou `<svg>` inline)

**Criterio:** conclusao gera transacao; historicos navegaveis; admin ve painel financeiro com dados reais.

---

## Etapa 12 — Persistencia em PostgreSQL

**Objetivo:** substituir todos os repositories in-memory por implementacoes Postgres. Services, handlers e frontend nao mudam.

### Backend — Infraestrutura
- [ ] `docker-compose up -d` sobe Postgres + pgAdmin + Flyway aplicando as migrations `V1` ate `V{N}` (uma por etapa — `V3__preferencias_cliente.sql`, `V4__catalogo_precificacao.sql`, ..., `V11__transacoes.sql`)
- [ ] `internal/repository/postgres/db.go` — pool `pgxpool`, health check, graceful shutdown

### Backend — Repositories Postgres
Um arquivo por dominio em `internal/repository/postgres/`, todas as queries parametrizadas com `$1, $2...` (OWASP A05):
- [ ] `usuario.go`, `cliente.go`, `profissional.go`, `endereco.go`, `preferencia.go`
- [ ] `regiao.go`, `disponibilidade.go`, `documento.go`, `referencia.go`
- [ ] `catalogo.go` (categorias + opcionais + tabela de precos), `solicitacao.go` + `solicitacao_opcional.go`
- [ ] `servico.go` + `historico.go`, `avaliacao.go`, `recorrencia.go`, `notificacao.go`
- [ ] `transacao.go` + `dados_bancarios.go`

### Backend — Integracao
- [ ] `cmd/api/main.go` aceita `USE_POSTGRES=true` e injeta repositories Postgres
- [ ] Consulta de limite LC 150/2015 migra para usar `vw_servicos_por_semana` (view SQL)
- [ ] Triggers de banco assumem recalculos que estavam no service (`trg_recalcular_nota`, `trg_atualizado_em`) — services simplificam
- [ ] Testes de integracao (ja configurados com testcontainers desde a Etapa 2) passam para todos os repositories

### Frontend
- [ ] **Nenhuma mudanca de codigo**. Regressao manual completa em todos os fluxos apos a troca.

**Criterio:** `USE_POSTGRES=true go run cmd/api/main.go` funciona identico ao in-memory; dados sobrevivem a restart; todos os testes de integracao verdes.

---

## Etapa 13 — Finalizacao do MVP: Chat Simples + Recibo PDF + E2E + Polish

**Objetivo:** amarrar as pontas do MVP com dois itens de profissionalismo (diferenciais D e E) e garantir qualidade final.

### (D) Chat limitado — mensagem unica do cliente
- [ ] `internal/domain/mensagem_solicitacao.go` — campo simples `observacao_cliente` (text, max 500 chars) adicionado a `Solicitacao`. **Nao e chat real-time**; e uma nota que o cliente passa junto com a solicitacao ("porteiro e fulano; chave com vizinho; tem gato").
- [ ] Migration incremental para adicionar a coluna (se Etapa 12 ja rodou)
- [ ] Frontend: no passo 4 do fluxo de solicitacao, campo `textarea` opcional "Observacoes para a profissional" (maxlength 500)
- [ ] Profissional ve a observacao no detalhe do servico em secao destacada
- [ ] Testes: limite de 500 chars, sanitizacao (sem `@html`)

### (E) Recibo PDF simples
- [ ] `internal/service/recibo_service.go` — gera PDF com dados do servico (cliente, profissional, endereco, data, itens do breakdown, valor, assinatura digital simples = hash SHA256 dos campos)
- [ ] Biblioteca: `github.com/jung-kurt/gofpdf` ou `github.com/go-pdf/fpdf` (ambos sem dependencia nativa)
- [ ] Endpoint: `GET /servicos/:id/recibo` retorna `application/pdf`
- [ ] Frontend: botao "Baixar recibo" no detalhe de servico `CONCLUIDO` — tanto para cliente quanto profissional
- [ ] Testes: PDF gerado com campos corretos; hash permite verificar integridade

### Qualidade final (ja e obrigacao continua pelo CLAUDE.md, esta etapa so valida)
- [ ] `go test ./...` verde
- [ ] `go test -tags=integration ./...` verde
- [ ] `go vet ./...` + `gofmt -l .` limpos
- [ ] `npm run check` com 0 erros
- [ ] Swagger atualizado (`swag init -g cmd/api/main.go -o ../docs/swagger`) com todos os handlers anotados
- [ ] Checklist OWASP Top 10:2025 aplicado (via `/owasp-security`)

### Testes E2E (Playwright)
- [ ] Fluxo cliente: registro -> completar perfil -> cadastrar endereco -> solicitar servico -> acompanhar status -> avaliar apos conclusao -> favoritar profissional -> baixar recibo
- [ ] Fluxo profissional: registro -> completar perfil -> adicionar referencias -> definir regioes + disponibilidade -> aprovar atribuicao -> check-in -> check-out -> confirmar -> ver historico
- [ ] Fluxo admin: aprovar profissional -> alterar preco -> verificar dashboard financeiro

### Polish UX
- [ ] Revisao de responsividade (mobile 375px + desktop 1440px) em todas as paginas
- [ ] Loading states e empty states consistentes (padrao ja estabelecido na Etapa 2)
- [ ] Mensagens de erro amigaveis (nunca expor stack trace)
- [ ] Auditoria de acessibilidade minima: labels em todos os inputs, contraste AA, navegacao por teclado

**Criterio:** cliente consegue baixar recibo de servico concluido; observacao e visivel para a profissional; E2E cobrindo os 3 fluxos principais verde; Swagger completo; zero `any` no frontend; zero alerta de `govulncheck`/`npm audit` alta/critica.

---

## Fora do MVP (Fases futuras)

| Feature | Fase | Comentario |
|---|---|---|
| Integracao com CDN para upload real de fotos | 2 | Hoje e URL externa; o contrato do backend nao muda — so o componente de input vira uploader |
| Gateway de pagamento (cartao/Pix, escrow, repasse) | 2 | Estrutura de `Transacao` na Etapa 11 ja prepara |
| Chat in-app real-time | 3 | Hoje temos "observacao unica" (Etapa 13); chat com WebSocket e outro bicho |
| Programa de indicacao | 3 | |
| Planos premium para profissionais | 3 | |
| Multas automaticas por cancelamento tardio | 3 | Hoje e so penalidade no score |
| Antecipacao de recebiveis | 4 | Depende de gateway |
| Integracao contabil | 4 | |
| Assistencia residencial por assinatura | 4 | |
| Parcerias com condominios | 4 | |
| App mobile nativo (React Native / Flutter) | 4 | |
| Expansao para outras cidades | Depende de demanda | Arquitetura ja parametrizavel (CEP, regiao) |

---

*Plano atualizado em 18/abril/2026. Etapas 0-2 concluidas. Proxima: Etapa 2-residual (preferencias) + Etapa 3 (catalogo com orcamento transparente).*
