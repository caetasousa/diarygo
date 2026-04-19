-- =============================================================================
-- DiaryGo — Schema Futuro (Referência para Etapas 3-11)
-- =============================================================================
-- ESTE ARQUIVO NÃO É EXECUTADO PELO FLYWAY. É um arquivo de referência que
-- documenta o desenho planejado das tabelas das Etapas 3-11.
--
-- Cada bloco deve virar uma migration V{N}__*.sql separada conforme a etapa
-- for implementada. Ao extrair um bloco, copie-o para uma nova migration,
-- adicione o cabeçalho obrigatório (etapa, resumo, delta, DROP reverso) e
-- remova-o daqui. Atualize também README.md ("Esquema do Banco de Dados" e
-- "Histórico de Migrations") e CLAUDE.md se necessário.
--
-- As migrations V1 (auth) e V2 (cadastro) já estão em backend/migrations/.
-- =============================================================================


-- =============================================================================
-- Etapa 3 — Catálogo e Precificação — JÁ MIGRADO
-- =============================================================================
-- Movido para backend/migrations/V4__catalogo_precificacao.sql.
-- Preferências (diferencial B) estão em backend/migrations/V3__preferencias_cliente.sql.


-- =============================================================================
-- Etapa 4 — Solicitações de Serviço
-- Virará V5__solicitacoes.sql
-- Depende de: V2 (clientes, enderecos), V4 (categorias_servico, opcionais)
-- =============================================================================

-- CREATE TYPE frequencia_servico AS ENUM ('UNICA', 'SEMANAL', 'DUAS_POR_SEMANA', 'QUINZENAL');
-- CREATE TYPE status_solicitacao AS ENUM (
--     'AGUARDANDO_PROFISSIONAL', 'PROFISSIONAL_ATRIBUIDA', 'CONFIRMADA',
--     'EM_ANDAMENTO', 'CONCLUIDA', 'CANCELADA_CLIENTE', 'CANCELADA_PLATAFORMA'
-- );

-- CREATE TABLE solicitacoes (
--     id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     cliente_id          UUID NOT NULL REFERENCES clientes(id),
--     endereco_id         UUID NOT NULL REFERENCES enderecos(id),
--     categoria_id        UUID NOT NULL REFERENCES categorias_servico(id),
--     frequencia          frequencia_servico NOT NULL DEFAULT 'UNICA',
--     data_servico        DATE NOT NULL,
--     hora_inicio         TIME NOT NULL,
--     duracao_estimada_min INTEGER NOT NULL,
--     valor_referencia    NUMERIC(10,2) NOT NULL,
--     observacoes         TEXT,
--     status              status_solicitacao NOT NULL DEFAULT 'AGUARDANDO_PROFISSIONAL',
--     criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE TABLE solicitacao_opcionais (
--     solicitacao_id  UUID NOT NULL REFERENCES solicitacoes(id) ON DELETE CASCADE,
--     opcional_id     UUID NOT NULL REFERENCES opcionais(id),
--     PRIMARY KEY (solicitacao_id, opcional_id)
-- );

-- CREATE INDEX idx_solicitacoes_cliente ON solicitacoes(cliente_id);
-- CREATE INDEX idx_solicitacoes_status ON solicitacoes(status);
-- CREATE INDEX idx_solicitacoes_data ON solicitacoes(data_servico);
-- CREATE TRIGGER trg_solicitacoes_atualizado BEFORE UPDATE ON solicitacoes
--     FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();


-- =============================================================================
-- Etapa 5 — Matching e Atribuição (Serviços)
-- Virará V5__servicos.sql
-- Depende de: V2 (profissionais), V4 (solicitacoes)
-- Inclui view vw_servicos_por_semana (LC 150/2015)
-- =============================================================================

-- CREATE TYPE status_servico AS ENUM (
--     'AGENDADO', 'EM_ANDAMENTO', 'CONCLUIDO',
--     'CANCELADO', 'NO_SHOW_CLIENTE', 'NO_SHOW_PROFISSIONAL'
-- );

-- CREATE TABLE servicos (
--     id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     solicitacao_id      UUID NOT NULL UNIQUE REFERENCES solicitacoes(id),
--     profissional_id     UUID NOT NULL REFERENCES profissionais(id),
--     recorrencia_id      UUID, -- FK adicionada em V8
--     status              status_servico NOT NULL DEFAULT 'AGENDADO',
--     checkin_em          TIMESTAMPTZ,
--     checkin_lat         NUMERIC(10,7),
--     checkin_lon         NUMERIC(10,7),
--     checkout_em         TIMESTAMPTZ,
--     checkout_lat        NUMERIC(10,7),
--     checkout_lon        NUMERIC(10,7),
--     confirmado_cliente  BOOLEAN NOT NULL DEFAULT FALSE,
--     confirmado_prof     BOOLEAN NOT NULL DEFAULT FALSE,
--     criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE INDEX idx_servicos_profissional ON servicos(profissional_id);
-- CREATE INDEX idx_servicos_status ON servicos(status);
-- CREATE TRIGGER trg_servicos_atualizado BEFORE UPDATE ON servicos
--     FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

-- CREATE VIEW vw_servicos_por_semana AS
-- SELECT
--     s.profissional_id,
--     sol.endereco_id,
--     sol.cliente_id,
--     DATE_TRUNC('week', sol.data_servico) AS semana,
--     COUNT(*) AS total_servicos
-- FROM servicos s
-- JOIN solicitacoes sol ON sol.id = s.solicitacao_id
-- WHERE s.status NOT IN ('CANCELADO', 'NO_SHOW_CLIENTE', 'NO_SHOW_PROFISSIONAL')
-- GROUP BY s.profissional_id, sol.endereco_id, sol.cliente_id,
--          DATE_TRUNC('week', sol.data_servico);


-- =============================================================================
-- Etapa 6 — Execução (histórico de status)
-- Virará V6__historico_status.sql
-- Depende de: V5 (servicos)
-- =============================================================================

-- CREATE TABLE historico_status (
--     id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     servico_id      UUID NOT NULL REFERENCES servicos(id) ON DELETE CASCADE,
--     status_anterior status_servico,
--     status_novo     status_servico NOT NULL,
--     observacao      TEXT,
--     criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );


-- =============================================================================
-- Etapa 7 — Avaliações e Reputação
-- Virará V7__avaliacoes.sql
-- Depende de: V2 (clientes, profissionais), V5 (servicos)
-- =============================================================================

-- CREATE TABLE avaliacoes_cliente (
--     id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     servico_id      UUID NOT NULL UNIQUE REFERENCES servicos(id),
--     cliente_id      UUID NOT NULL REFERENCES clientes(id),
--     profissional_id UUID NOT NULL REFERENCES profissionais(id),
--     nota            SMALLINT NOT NULL CHECK (nota >= 1 AND nota <= 5),
--     comentario      TEXT,
--     pontualidade    SMALLINT CHECK (pontualidade >= 1 AND pontualidade <= 5),
--     qualidade       SMALLINT CHECK (qualidade >= 1 AND qualidade <= 5),
--     educacao        SMALLINT CHECK (educacao >= 1 AND educacao <= 5),
--     criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE TABLE avaliacoes_profissional (
--     id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     servico_id      UUID NOT NULL UNIQUE REFERENCES servicos(id),
--     profissional_id UUID NOT NULL REFERENCES profissionais(id),
--     cliente_id      UUID NOT NULL REFERENCES clientes(id),
--     nota            SMALLINT NOT NULL CHECK (nota >= 1 AND nota <= 5),
--     ambiente        SMALLINT CHECK (ambiente >= 1 AND ambiente <= 5),
--     materiais       SMALLINT CHECK (materiais >= 1 AND materiais <= 5),
--     respeito        SMALLINT CHECK (respeito >= 1 AND respeito <= 5),
--     criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE OR REPLACE FUNCTION fn_recalcular_nota_profissional()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     UPDATE profissionais
--     SET nota_media = (
--             SELECT ROUND(AVG(nota)::NUMERIC, 2)
--             FROM avaliacoes_cliente
--             WHERE profissional_id = NEW.profissional_id
--         ),
--         total_servicos = (
--             SELECT COUNT(*)
--             FROM avaliacoes_cliente
--             WHERE profissional_id = NEW.profissional_id
--         )
--     WHERE id = NEW.profissional_id;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
-- CREATE TRIGGER trg_recalcular_nota AFTER INSERT OR UPDATE ON avaliacoes_cliente
--     FOR EACH ROW EXECUTE FUNCTION fn_recalcular_nota_profissional();


-- =============================================================================
-- Etapa 8 — Recorrências
-- Virará V8__recorrencias.sql
-- Depende de: V4 (solicitacoes), V5 (servicos — adiciona FK recorrencia_id)
-- =============================================================================

-- CREATE TYPE status_recorrencia AS ENUM ('ATIVA', 'PAUSADA', 'CANCELADA');

-- CREATE TABLE recorrencias (
--     id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     solicitacao_origem_id UUID NOT NULL REFERENCES solicitacoes(id),
--     cliente_id          UUID NOT NULL REFERENCES clientes(id),
--     frequencia          frequencia_servico NOT NULL,
--     dia_semana_1        SMALLINT,
--     dia_semana_2        SMALLINT,
--     hora_inicio         TIME NOT NULL,
--     status              status_recorrencia NOT NULL DEFAULT 'ATIVA',
--     proxima_data        DATE,
--     criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- ALTER TABLE servicos ADD CONSTRAINT servicos_recorrencia_fk
--     FOREIGN KEY (recorrencia_id) REFERENCES recorrencias(id);


-- =============================================================================
-- Etapa 9 — Notificações
-- Virará V9__notificacoes.sql
-- Depende de: V1 (usuarios)
-- =============================================================================

-- CREATE TYPE canal_notificacao AS ENUM ('PUSH', 'EMAIL', 'SMS', 'IN_APP');
-- CREATE TYPE status_notificacao AS ENUM ('PENDENTE', 'ENVIADA', 'FALHOU', 'LIDA');

-- CREATE TABLE notificacoes (
--     id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     usuario_id      UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
--     canal           canal_notificacao NOT NULL,
--     titulo          VARCHAR(255) NOT NULL,
--     corpo           TEXT NOT NULL,
--     status          status_notificacao NOT NULL DEFAULT 'PENDENTE',
--     lida            BOOLEAN NOT NULL DEFAULT FALSE,
--     criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE INDEX idx_notificacoes_usuario ON notificacoes(usuario_id);
-- CREATE INDEX idx_notificacoes_status ON notificacoes(status);


-- =============================================================================
-- Etapa 10 — Painel Admin
-- Virará V10__admins.sql
-- Depende de: V1 (usuarios)
-- Observação: coluna documentos.analisado_por já aponta para admins(id) em V2,
-- mas como a tabela admins só existe na V10, em V2 a FK é comentada/diferida.
-- Alternativa: criar admins desde V1 (já está previsto — ver nota em V1/V2).
-- =============================================================================

-- CREATE TABLE admins (
--     id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     usuario_id      UUID NOT NULL UNIQUE REFERENCES usuarios(id) ON DELETE CASCADE,
--     nome            VARCHAR(255) NOT NULL,
--     permissoes      JSONB NOT NULL DEFAULT '{}',
--     criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );


-- =============================================================================
-- Etapa 11 — Transações e Dados Bancários
-- Virará V11__transacoes.sql
-- Depende de: V2 (profissionais, clientes), V5 (servicos)
-- =============================================================================

-- CREATE TYPE status_transacao AS ENUM ('REGISTRADA', 'PENDENTE', 'PAGA', 'ESTORNADA');
-- CREATE TYPE metodo_pagamento AS ENUM ('DIRETO_EXTERNO', 'CARTAO_CREDITO', 'PIX', 'BOLETO');

-- CREATE TABLE dados_bancarios (
--     id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     profissional_id     UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
--     banco               VARCHAR(100),
--     agencia             VARCHAR(20),
--     conta               VARCHAR(30),
--     tipo_conta          VARCHAR(20),
--     chave_pix           VARCHAR(255),
--     principal           BOOLEAN NOT NULL DEFAULT FALSE,
--     criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE TABLE transacoes (
--     id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     servico_id      UUID NOT NULL UNIQUE REFERENCES servicos(id),
--     cliente_id      UUID NOT NULL REFERENCES clientes(id),
--     profissional_id UUID NOT NULL REFERENCES profissionais(id),
--     valor           NUMERIC(10,2) NOT NULL,
--     metodo          metodo_pagamento NOT NULL DEFAULT 'DIRETO_EXTERNO',
--     status          status_transacao NOT NULL DEFAULT 'REGISTRADA',
--     referencia_ext  VARCHAR(255),
--     criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
-- );

-- CREATE TRIGGER trg_transacoes_atualizado BEFORE UPDATE ON transacoes
--     FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();
