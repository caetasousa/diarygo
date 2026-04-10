-- =============================================================================
-- DiaryGo — Schema Inicial (Etapa 12+)
-- Aplicado via Flyway. Não editar após aplicado — criar nova migration.
-- =============================================================================

-- Extensões
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "unaccent";

-- =============================================================================
-- ENUMS
-- =============================================================================

CREATE TYPE tipo_usuario AS ENUM ('CLIENTE', 'PROFISSIONAL', 'ADMIN');
CREATE TYPE status_profissional AS ENUM ('PENDENTE', 'APROVADA', 'REPROVADA', 'SUSPENSA', 'DESCREDENCIADA');
CREATE TYPE status_documento AS ENUM ('PENDENTE', 'APROVADO', 'REPROVADO');
CREATE TYPE status_referencia AS ENUM ('PENDENTE', 'CONFIRMADA', 'NAO_CONFIRMADA');
CREATE TYPE frequencia_servico AS ENUM ('UNICA', 'SEMANAL', 'DUAS_POR_SEMANA', 'QUINZENAL');
CREATE TYPE status_solicitacao AS ENUM (
    'AGUARDANDO_PROFISSIONAL', 'PROFISSIONAL_ATRIBUIDA', 'CONFIRMADA',
    'EM_ANDAMENTO', 'CONCLUIDA', 'CANCELADA_CLIENTE', 'CANCELADA_PLATAFORMA'
);
CREATE TYPE status_servico AS ENUM (
    'AGENDADO', 'EM_ANDAMENTO', 'CONCLUIDO',
    'CANCELADO', 'NO_SHOW_CLIENTE', 'NO_SHOW_PROFISSIONAL'
);
CREATE TYPE status_recorrencia AS ENUM ('ATIVA', 'PAUSADA', 'CANCELADA');
CREATE TYPE canal_notificacao AS ENUM ('PUSH', 'EMAIL', 'SMS', 'IN_APP');
CREATE TYPE status_notificacao AS ENUM ('PENDENTE', 'ENVIADA', 'FALHOU', 'LIDA');
CREATE TYPE status_transacao AS ENUM ('REGISTRADA', 'PENDENTE', 'PAGA', 'ESTORNADA');
CREATE TYPE metodo_pagamento AS ENUM ('DIRETO_EXTERNO', 'CARTAO_CREDITO', 'PIX', 'BOLETO');
CREATE TYPE tipo_documento AS ENUM ('RG_FRENTE', 'RG_VERSO', 'CPF', 'COMPROVANTE_RESIDENCIA', 'FOTO_ROSTO', 'OUTRO');

-- =============================================================================
-- FUNÇÃO helper: atualiza atualizado_em
-- =============================================================================

CREATE OR REPLACE FUNCTION fn_set_atualizado_em()
RETURNS TRIGGER AS $$
BEGIN
    NEW.atualizado_em = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- TABELAS CORE
-- =============================================================================

CREATE TABLE usuarios (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(255) NOT NULL UNIQUE,
    senha_hash      VARCHAR(255) NOT NULL,
    tipo            tipo_usuario NOT NULL,
    email_verificado BOOLEAN NOT NULL DEFAULT FALSE,
    ativo           BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE clientes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id      UUID NOT NULL UNIQUE REFERENCES usuarios(id) ON DELETE CASCADE,
    nome            VARCHAR(255) NOT NULL,
    cpf             CHAR(11) NOT NULL UNIQUE,
    telefone        VARCHAR(20),
    score           SMALLINT NOT NULL DEFAULT 100 CHECK (score >= 0 AND score <= 100),
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE profissionais (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id          UUID NOT NULL UNIQUE REFERENCES usuarios(id) ON DELETE CASCADE,
    nome                VARCHAR(255) NOT NULL,
    cpf                 CHAR(11) NOT NULL UNIQUE,
    rg                  VARCHAR(20),
    telefone            VARCHAR(20),
    foto_url            VARCHAR(500),
    status              status_profissional NOT NULL DEFAULT 'PENDENTE',
    nota_media          NUMERIC(3,2) CHECK (nota_media >= 1.0 AND nota_media <= 5.0),
    total_servicos      INTEGER NOT NULL DEFAULT 0,
    mei                 BOOLEAN NOT NULL DEFAULT FALSE,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE admins (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id      UUID NOT NULL UNIQUE REFERENCES usuarios(id) ON DELETE CASCADE,
    nome            VARCHAR(255) NOT NULL,
    permissoes      JSONB NOT NULL DEFAULT '{}',
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- DOCUMENTOS E REFERÊNCIAS
-- =============================================================================

CREATE TABLE documentos (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profissional_id     UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    tipo                tipo_documento NOT NULL,
    url                 VARCHAR(500) NOT NULL,
    status              status_documento NOT NULL DEFAULT 'PENDENTE',
    analisado_por       UUID REFERENCES admins(id),
    observacao          TEXT,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE referencias (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profissional_id     UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    nome_contato        VARCHAR(255) NOT NULL,
    telefone_contato    VARCHAR(20) NOT NULL,
    status              status_referencia NOT NULL DEFAULT 'PENDENTE',
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- REGIÕES E DISPONIBILIDADE
-- =============================================================================

CREATE TABLE regioes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome            VARCHAR(255) NOT NULL,
    cidade          VARCHAR(255) NOT NULL,
    estado          CHAR(2) NOT NULL,
    cep_inicio      CHAR(8),
    cep_fim         CHAR(8),
    ativa           BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE regioes_atuacao (
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    regiao_id       UUID NOT NULL REFERENCES regioes(id) ON DELETE CASCADE,
    PRIMARY KEY (profissional_id, regiao_id)
);

CREATE TABLE disponibilidades (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profissional_id     UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    dia_semana          SMALLINT NOT NULL CHECK (dia_semana >= 0 AND dia_semana <= 6), -- 0=Dom, 6=Sáb
    hora_inicio         TIME NOT NULL,
    hora_fim            TIME NOT NULL,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- CATÁLOGO
-- =============================================================================

CREATE TABLE categorias_servico (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome                VARCHAR(255) NOT NULL UNIQUE,
    descricao           TEXT,
    duracao_minima_min  INTEGER NOT NULL,
    ativa               BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE categorias_profissional (
    profissional_id     UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    categoria_id        UUID NOT NULL REFERENCES categorias_servico(id) ON DELETE CASCADE,
    PRIMARY KEY (profissional_id, categoria_id)
);

CREATE TABLE opcionais (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome                VARCHAR(255) NOT NULL UNIQUE,
    descricao           TEXT,
    valor_extra         NUMERIC(10,2) NOT NULL DEFAULT 0,
    tempo_extra_min     INTEGER NOT NULL DEFAULT 0,
    ativo               BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tabela_precos (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    categoria_id            UUID NOT NULL REFERENCES categorias_servico(id),
    regiao_id               UUID NOT NULL REFERENCES regioes(id),
    preco_hora              NUMERIC(10,2) NOT NULL,
    acrescimo_fds           NUMERIC(5,2) NOT NULL DEFAULT 0,    -- percentual ex: 10.00 = 10%
    desconto_semanal        NUMERIC(5,2) NOT NULL DEFAULT 10.00,
    desconto_quinzenal      NUMERIC(5,2) NOT NULL DEFAULT 5.00,
    desconto_duas_semana    NUMERIC(5,2) NOT NULL DEFAULT 15.00,
    ativa                   BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (categoria_id, regiao_id)
);

-- =============================================================================
-- ENDEREÇOS DO CLIENTE
-- =============================================================================

CREATE TABLE enderecos (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cliente_id      UUID NOT NULL REFERENCES clientes(id) ON DELETE CASCADE,
    logradouro      VARCHAR(255) NOT NULL,
    numero          VARCHAR(20),
    complemento     VARCHAR(100),
    bairro          VARCHAR(100) NOT NULL,
    cidade          VARCHAR(100) NOT NULL,
    estado          CHAR(2) NOT NULL,
    cep             CHAR(8) NOT NULL,
    lat             NUMERIC(10,7),
    lon             NUMERIC(10,7),
    principal       BOOLEAN NOT NULL DEFAULT FALSE,
    num_quartos     SMALLINT NOT NULL DEFAULT 0,
    num_banheiros   SMALLINT NOT NULL DEFAULT 0,
    num_salas       SMALLINT NOT NULL DEFAULT 0,
    num_cozinhas    SMALLINT NOT NULL DEFAULT 0,
    area_m2         NUMERIC(8,2),
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- DADOS BANCÁRIOS
-- =============================================================================

CREATE TABLE dados_bancarios (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profissional_id     UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    banco               VARCHAR(100),
    agencia             VARCHAR(20),
    conta               VARCHAR(30),
    tipo_conta          VARCHAR(20),
    chave_pix           VARCHAR(255),
    principal           BOOLEAN NOT NULL DEFAULT FALSE,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- SOLICITAÇÕES
-- =============================================================================

CREATE TABLE solicitacoes (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cliente_id          UUID NOT NULL REFERENCES clientes(id),
    endereco_id         UUID NOT NULL REFERENCES enderecos(id),
    categoria_id        UUID NOT NULL REFERENCES categorias_servico(id),
    frequencia          frequencia_servico NOT NULL DEFAULT 'UNICA',
    data_servico        DATE NOT NULL,
    hora_inicio         TIME NOT NULL,
    duracao_estimada_min INTEGER NOT NULL,
    valor_referencia    NUMERIC(10,2) NOT NULL,
    observacoes         TEXT,
    status              status_solicitacao NOT NULL DEFAULT 'AGUARDANDO_PROFISSIONAL',
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE solicitacao_opcionais (
    solicitacao_id  UUID NOT NULL REFERENCES solicitacoes(id) ON DELETE CASCADE,
    opcional_id     UUID NOT NULL REFERENCES opcionais(id),
    PRIMARY KEY (solicitacao_id, opcional_id)
);

-- =============================================================================
-- RECORRÊNCIAS
-- =============================================================================

CREATE TABLE recorrencias (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    solicitacao_origem_id UUID NOT NULL REFERENCES solicitacoes(id),
    cliente_id          UUID NOT NULL REFERENCES clientes(id),
    frequencia          frequencia_servico NOT NULL,
    dia_semana_1        SMALLINT,
    dia_semana_2        SMALLINT,
    hora_inicio         TIME NOT NULL,
    status              status_recorrencia NOT NULL DEFAULT 'ATIVA',
    proxima_data        DATE,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- SERVIÇOS (execução)
-- =============================================================================

CREATE TABLE servicos (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    solicitacao_id      UUID NOT NULL UNIQUE REFERENCES solicitacoes(id),
    profissional_id     UUID NOT NULL REFERENCES profissionais(id),
    recorrencia_id      UUID REFERENCES recorrencias(id),
    status              status_servico NOT NULL DEFAULT 'AGENDADO',
    checkin_em          TIMESTAMPTZ,
    checkin_lat         NUMERIC(10,7),
    checkin_lon         NUMERIC(10,7),
    checkout_em         TIMESTAMPTZ,
    checkout_lat        NUMERIC(10,7),
    checkout_lon        NUMERIC(10,7),
    confirmado_cliente  BOOLEAN NOT NULL DEFAULT FALSE,
    confirmado_prof     BOOLEAN NOT NULL DEFAULT FALSE,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE historico_status (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    servico_id      UUID NOT NULL REFERENCES servicos(id) ON DELETE CASCADE,
    status_anterior status_servico,
    status_novo     status_servico NOT NULL,
    observacao      TEXT,
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- AVALIAÇÕES
-- =============================================================================

CREATE TABLE avaliacoes_cliente (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    servico_id      UUID NOT NULL UNIQUE REFERENCES servicos(id),
    cliente_id      UUID NOT NULL REFERENCES clientes(id),
    profissional_id UUID NOT NULL REFERENCES profissionais(id),
    nota            SMALLINT NOT NULL CHECK (nota >= 1 AND nota <= 5),
    comentario      TEXT,
    pontualidade    SMALLINT CHECK (pontualidade >= 1 AND pontualidade <= 5),
    qualidade       SMALLINT CHECK (qualidade >= 1 AND qualidade <= 5),
    educacao        SMALLINT CHECK (educacao >= 1 AND educacao <= 5),
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE avaliacoes_profissional (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    servico_id      UUID NOT NULL UNIQUE REFERENCES servicos(id),
    profissional_id UUID NOT NULL REFERENCES profissionais(id),
    cliente_id      UUID NOT NULL REFERENCES clientes(id),
    nota            SMALLINT NOT NULL CHECK (nota >= 1 AND nota <= 5),
    ambiente        SMALLINT CHECK (ambiente >= 1 AND ambiente <= 5),
    materiais       SMALLINT CHECK (materiais >= 1 AND materiais <= 5),
    respeito        SMALLINT CHECK (respeito >= 1 AND respeito <= 5),
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- TRANSAÇÕES
-- =============================================================================

CREATE TABLE transacoes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    servico_id      UUID NOT NULL UNIQUE REFERENCES servicos(id),
    cliente_id      UUID NOT NULL REFERENCES clientes(id),
    profissional_id UUID NOT NULL REFERENCES profissionais(id),
    valor           NUMERIC(10,2) NOT NULL,
    metodo          metodo_pagamento NOT NULL DEFAULT 'DIRETO_EXTERNO',
    status          status_transacao NOT NULL DEFAULT 'REGISTRADA',
    referencia_ext  VARCHAR(255),
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- NOTIFICAÇÕES
-- =============================================================================

CREATE TABLE notificacoes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id      UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    canal           canal_notificacao NOT NULL,
    titulo          VARCHAR(255) NOT NULL,
    corpo           TEXT NOT NULL,
    status          status_notificacao NOT NULL DEFAULT 'PENDENTE',
    lida            BOOLEAN NOT NULL DEFAULT FALSE,
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- VIEW DE CONTROLE LEGAL — LC 150/2015
-- =============================================================================

CREATE VIEW vw_servicos_por_semana AS
SELECT
    s.profissional_id,
    sol.endereco_id,
    sol.cliente_id,
    DATE_TRUNC('week', sol.data_servico) AS semana,
    COUNT(*) AS total_servicos
FROM servicos s
JOIN solicitacoes sol ON sol.id = s.solicitacao_id
WHERE s.status NOT IN ('CANCELADO', 'NO_SHOW_CLIENTE', 'NO_SHOW_PROFISSIONAL')
GROUP BY s.profissional_id, sol.endereco_id, sol.cliente_id, DATE_TRUNC('week', sol.data_servico);

-- =============================================================================
-- TRIGGERS — atualizado_em
-- =============================================================================

CREATE TRIGGER trg_usuarios_atualizado
    BEFORE UPDATE ON usuarios
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_clientes_atualizado
    BEFORE UPDATE ON clientes
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_profissionais_atualizado
    BEFORE UPDATE ON profissionais
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_solicitacoes_atualizado
    BEFORE UPDATE ON solicitacoes
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_servicos_atualizado
    BEFORE UPDATE ON servicos
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_transacoes_atualizado
    BEFORE UPDATE ON transacoes
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

-- =============================================================================
-- TRIGGER — recalcular nota_media da profissional
-- =============================================================================

CREATE OR REPLACE FUNCTION fn_recalcular_nota_profissional()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE profissionais
    SET
        nota_media = (
            SELECT ROUND(AVG(nota)::NUMERIC, 2)
            FROM avaliacoes_cliente
            WHERE profissional_id = NEW.profissional_id
        ),
        total_servicos = (
            SELECT COUNT(*)
            FROM avaliacoes_cliente
            WHERE profissional_id = NEW.profissional_id
        )
    WHERE id = NEW.profissional_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_recalcular_nota
    AFTER INSERT OR UPDATE ON avaliacoes_cliente
    FOR EACH ROW EXECUTE FUNCTION fn_recalcular_nota_profissional();

-- =============================================================================
-- ÍNDICES
-- =============================================================================

CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_clientes_cpf ON clientes(cpf);
CREATE INDEX idx_profissionais_cpf ON profissionais(cpf);
CREATE INDEX idx_profissionais_status ON profissionais(status);
CREATE INDEX idx_solicitacoes_cliente ON solicitacoes(cliente_id);
CREATE INDEX idx_solicitacoes_status ON solicitacoes(status);
CREATE INDEX idx_solicitacoes_data ON solicitacoes(data_servico);
CREATE INDEX idx_servicos_profissional ON servicos(profissional_id);
CREATE INDEX idx_servicos_status ON servicos(status);
CREATE INDEX idx_notificacoes_usuario ON notificacoes(usuario_id);
CREATE INDEX idx_notificacoes_status ON notificacoes(status);
