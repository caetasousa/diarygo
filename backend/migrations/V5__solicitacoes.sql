-- =============================================================================
-- V5 — Solicitações de Serviço
-- =============================================================================
-- Etapa: 4 — Solicitações de Serviço
-- Resumo: cria as tabelas de solicitacao e solicitacao_opcional, que armazenam
--         o pedido do cliente com snapshot imutável do orçamento (breakdown JSONB,
--         valor_total, duracao_min). O snapshot congela o preço combinado mesmo que
--         a tabela de preços mude depois.
--
-- Delta vs. V4:
--   + enum frequencia_servico    (UNICA | SEMANAL | QUINZENAL | DUAS_POR_SEMANA)
--   + enum status_solicitacao    (AGUARDANDO | ATRIBUIDA | CONFIRMADA | EM_ANDAMENTO | CONCLUIDA | CANCELADA)
--   + tabela solicitacao         (FK cliente, endereco, categoria, regiao; snapshot de orçamento)
--   + tabela solicitacao_opcional (N:N solicitacao ↔ opcional com snapshot de preço/tempo)
--   + índice idx_solicitacao_cliente_status (listagem por cliente filtrada por status)
--   + índice idx_solicitacao_data           (busca por data do serviço)
--   + trigger trg_solicitacao_atualizado    (atualiza atualizado_em automaticamente)
--
-- Espelhamento service ↔ banco:
--   frequencia_servico ENUM               → domain.FrequenciaServico + Valida()
--   status_solicitacao ENUM               → domain.StatusSolicitacao + Valida()
--   observacao CHECK char_length <= 500   → domain.ValidarObservacao (conta runas no service)
--   data_servico >= NOW() + 24h           → domain.ValidarAntecedencia24h (checado no service)
--   FK endereco ON DELETE RESTRICT        → service verifica ownership antes de persistir
--   breakdown_json JSONB                  → service serializa []ItemCalculo; in-memory usa campo direto
--
-- Rollback (dev apenas — produção nunca roda):
--   DROP TRIGGER IF EXISTS trg_solicitacao_atualizado ON solicitacao;
--   DROP TABLE  IF EXISTS solicitacao_opcional;
--   DROP TABLE  IF EXISTS solicitacao;
--   DROP TYPE   IF EXISTS status_solicitacao;
--   DROP TYPE   IF EXISTS frequencia_servico;
-- =============================================================================

-- Frequência do serviço (reutilizada futuramente em recorrencias)
CREATE TYPE frequencia_servico AS ENUM (
    'UNICA',
    'SEMANAL',
    'QUINZENAL',
    'DUAS_POR_SEMANA'
);

CREATE TYPE status_solicitacao AS ENUM (
    'AGUARDANDO',
    'ATRIBUIDA',
    'CONFIRMADA',
    'EM_ANDAMENTO',
    'CONCLUIDA',
    'CANCELADA'
);

CREATE TABLE solicitacao (
    id              UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    cliente_id      UUID        NOT NULL REFERENCES clientes(id)           ON DELETE RESTRICT,
    endereco_id     UUID        NOT NULL REFERENCES enderecos(id)          ON DELETE RESTRICT,
    categoria_id    UUID        NOT NULL REFERENCES categorias_servico(id) ON DELETE RESTRICT,
    regiao_id       UUID        NOT NULL REFERENCES regioes(id)            ON DELETE RESTRICT,

    num_quartos     INTEGER     NOT NULL CHECK (num_quartos  BETWEEN 1 AND 20),
    num_banheiros   INTEGER     NOT NULL CHECK (num_banheiros BETWEEN 0 AND 20),
    num_salas       INTEGER     NOT NULL CHECK (num_salas    BETWEEN 0 AND 20),
    num_cozinhas    INTEGER     NOT NULL CHECK (num_cozinhas BETWEEN 0 AND 20),

    frequencia      frequencia_servico  NOT NULL,
    data_servico    TIMESTAMPTZ NOT NULL,
    -- observacao: aceita até 500 chars unicode; char_length conta caracteres, não bytes
    observacao      TEXT        NOT NULL DEFAULT '' CHECK (char_length(observacao) <= 500),

    -- Snapshot congelado do orçamento no momento da criação
    valor_total     NUMERIC(10,2) NOT NULL CHECK (valor_total >= 0),
    duracao_min     INTEGER       NOT NULL CHECK (duracao_min > 0),
    breakdown_json  JSONB         NOT NULL,  -- serialização de []ItemCalculo

    status          status_solicitacao NOT NULL DEFAULT 'AGUARDANDO',
    criada_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizada_em   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelada_em    TIMESTAMPTZ
);

CREATE TABLE solicitacao_opcional (
    solicitacao_id  UUID          NOT NULL REFERENCES solicitacao(id) ON DELETE CASCADE,
    opcional_id     UUID          NOT NULL REFERENCES opcionais(id)   ON DELETE RESTRICT,
    nome_snapshot   VARCHAR(255)  NOT NULL,
    valor_snapshot  NUMERIC(10,2) NOT NULL,
    tempo_snapshot  INTEGER       NOT NULL,
    PRIMARY KEY (solicitacao_id, opcional_id)
);

CREATE INDEX idx_solicitacao_cliente_status ON solicitacao(cliente_id, status);
CREATE INDEX idx_solicitacao_data           ON solicitacao(data_servico);

CREATE TRIGGER trg_solicitacao_atualizado
    BEFORE UPDATE ON solicitacao
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();
