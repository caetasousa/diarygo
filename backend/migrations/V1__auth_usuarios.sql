-- =============================================================================
-- V1 — Autenticação e Usuários
-- =============================================================================
-- Etapa: 1 — Autenticação e gestão de contas
-- Resumo: cria a fundação de identidade do sistema (tabela `usuarios`), os
--         tipos de conta (CLIENTE, PROFISSIONAL, ADMIN), extensões Postgres
--         usadas globalmente e a função `fn_set_atualizado_em()` reutilizada
--         por triggers de todas as migrations seguintes. Inclui também as
--         colunas de recuperação de senha ("esqueci minha senha"), com índice
--         parcial para custo próximo de zero.
--
-- Princípio de validação — defesa em profundidade:
--   O banco é a ÚLTIMA linha de defesa (UNIQUE, NOT NULL, tipo ENUM). O
--   service (internal/service/*) é a PRIMEIRA — valida formato, normaliza
--   entrada e retorna erros de domínio amigáveis antes do INSERT/UPDATE.
--   Toda constraint aqui tem sua contraparte no service:
--     - email UNIQUE                → AuthService normaliza (lower+trim) e
--                                     ValidarEmail checa formato
--     - senha_hash NOT NULL         → AuthService aplica bcrypt cost 12
--                                     sobre senha já validada pelo NIST 800-63b
--                                     (8-72 chars)
--     - tipo ENUM                   → AuthService só aceita CLIENTE/PROFISSIONAL
--                                     no registro público; ADMIN via seed
--     - token_recuperacao VARCHAR(64) → AuthService gera com crypto/rand
--                                       (32 bytes hex = 64 chars); expiração
--                                       validada antes de aceitar reset
--
-- Delta vs. migrations anteriores: primeira migration — define o baseline.
--   + Extensões: uuid-ossp, unaccent
--   + Enum: tipo_usuario
--   + Tabela: usuarios (inclui token_recuperacao e token_recuperacao_expira)
--   + Função: fn_set_atualizado_em()
--   + Trigger: trg_usuarios_atualizado
--   + Índice parcial: idx_usuarios_token_recuperacao
--   (Não há índice redundante em email — o UNIQUE já cria B-tree automático.)
--
-- Rollback (apenas em dev — nunca em produção):
--   DROP TRIGGER IF EXISTS trg_usuarios_atualizado ON usuarios;
--   DROP INDEX IF EXISTS idx_usuarios_token_recuperacao;
--   DROP TABLE IF EXISTS usuarios;
--   DROP FUNCTION IF EXISTS fn_set_atualizado_em();
--   DROP TYPE IF EXISTS tipo_usuario;
--   DROP EXTENSION IF EXISTS "unaccent";
--   DROP EXTENSION IF EXISTS "uuid-ossp";
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "unaccent";

CREATE TYPE tipo_usuario AS ENUM ('CLIENTE', 'PROFISSIONAL', 'ADMIN');

CREATE OR REPLACE FUNCTION fn_set_atualizado_em()
RETURNS TRIGGER AS $$
BEGIN
    NEW.atualizado_em = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE usuarios (
    id                       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email                    VARCHAR(255) NOT NULL UNIQUE,
    senha_hash               VARCHAR(255) NOT NULL,
    tipo                     tipo_usuario NOT NULL,
    email_verificado         BOOLEAN NOT NULL DEFAULT FALSE,
    ativo                    BOOLEAN NOT NULL DEFAULT TRUE,
    token_recuperacao        VARCHAR(64),
    token_recuperacao_expira TIMESTAMPTZ,
    criado_em                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índice parcial: só indexa linhas com token ativo. Mantém custo de
-- escrita/vacuum próximo de zero — token_recuperacao é NULL na maioria.
CREATE INDEX idx_usuarios_token_recuperacao
    ON usuarios(token_recuperacao)
    WHERE token_recuperacao IS NOT NULL;

CREATE TRIGGER trg_usuarios_atualizado
    BEFORE UPDATE ON usuarios
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();
