-- =============================================================================
-- V4 — Catálogo e Precificação
-- =============================================================================
-- Etapa: 3 — Catálogo, Opcionais e Tabela de Preços (diferencial C)
-- Resumo: introduz o vocabulário de serviço da plataforma — categorias (tipos
--         de faxina), opcionais (itens extras cobráveis) e a tabela de preços
--         cruzando categoria × região. Habilita o diferencial C (orçamento
--         transparente): o service PrecificacaoService consome essas tabelas
--         para produzir o breakdown público em /precos/calcular.
--
-- Delta vs. V3:
--   + tabela categorias_servico         (7 categorias no MVP)
--   + tabela categorias_profissional    (N:N profissional ↔ categoria — MVP aceita todas)
--   + tabela opcionais                  (7 opcionais compartilhados no MVP)
--   + tabela tabela_precos              (UNIQUE categoria_id + regiao_id)
--   + índices por categoria/região
--   + triggers trg_*_atualizado reutilizam fn_set_atualizado_em() da V2
--
-- Princípio de precificação:
--   Em V4 o preço é LINEAR: preco_hora × horas + opcionais. Fim-de-semana e
--   frequência entram como acréscimo/desconto aplicados no service (itens do
--   breakdown), não como colunas de cálculo. A tabela guarda os percentuais
--   por categoria×região para que o admin possa ajustá-los sem deploy.
--
-- Espelhamento service ↔ banco:
--   categorias_servico.duracao_minima_min    → service.TempoExtraPorQuarto/Banheiro/Sala/Cozinha (ajustes)
--   tabela_precos.preco_hora > 0             → domain.ValidarPrecoHora (futuro)
--   tabela_precos.acrescimo_fds 0..100       → domain.ValidarPercentualFDS (futuro)
--   tabela_precos.desconto_*   0..100        → service aplica se frequência != UNICA
--   UNIQUE (categoria_id, regiao_id)         → CatalogoRepository.BuscarTabelaPreco (chave composta)
--
-- Rollback (dev apenas — produção nunca roda):
--   DROP TRIGGER IF EXISTS trg_tabela_precos_atualizado ON tabela_precos;
--   DROP TRIGGER IF EXISTS trg_opcionais_atualizado     ON opcionais;
--   DROP TRIGGER IF EXISTS trg_categorias_atualizado    ON categorias_servico;
--   DROP TABLE  IF EXISTS tabela_precos;
--   DROP TABLE  IF EXISTS categorias_profissional;
--   DROP TABLE  IF EXISTS opcionais;
--   DROP TABLE  IF EXISTS categorias_servico;
-- =============================================================================

CREATE TABLE categorias_servico (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome                VARCHAR(255) NOT NULL UNIQUE,
    descricao           TEXT,
    duracao_minima_min  INTEGER      NOT NULL CHECK (duracao_minima_min > 0),
    ativa               BOOLEAN      NOT NULL DEFAULT TRUE,
    criado_em           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_categorias_atualizado
    BEFORE UPDATE ON categorias_servico
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TABLE categorias_profissional (
    profissional_id     UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    categoria_id        UUID NOT NULL REFERENCES categorias_servico(id) ON DELETE CASCADE,
    criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (profissional_id, categoria_id)
);

CREATE INDEX idx_cat_prof_categoria ON categorias_profissional(categoria_id);

CREATE TABLE opcionais (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome                VARCHAR(255) NOT NULL UNIQUE,
    descricao           TEXT,
    valor_extra         NUMERIC(10,2) NOT NULL DEFAULT 0 CHECK (valor_extra >= 0),
    tempo_extra_min     INTEGER       NOT NULL DEFAULT 0 CHECK (tempo_extra_min >= 0),
    ativo               BOOLEAN       NOT NULL DEFAULT TRUE,
    criado_em           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    atualizado_em       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_opcionais_atualizado
    BEFORE UPDATE ON opcionais
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TABLE tabela_precos (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    categoria_id            UUID NOT NULL REFERENCES categorias_servico(id) ON DELETE CASCADE,
    regiao_id               UUID NOT NULL REFERENCES regioes(id)            ON DELETE CASCADE,
    preco_hora              NUMERIC(10,2) NOT NULL CHECK (preco_hora > 0),
    acrescimo_fds           NUMERIC(5,2)  NOT NULL DEFAULT 0     CHECK (acrescimo_fds     >= 0 AND acrescimo_fds     <= 100),
    desconto_semanal        NUMERIC(5,2)  NOT NULL DEFAULT 10.00 CHECK (desconto_semanal  >= 0 AND desconto_semanal  <= 100),
    desconto_quinzenal      NUMERIC(5,2)  NOT NULL DEFAULT 5.00  CHECK (desconto_quinzenal >= 0 AND desconto_quinzenal <= 100),
    desconto_duas_semana    NUMERIC(5,2)  NOT NULL DEFAULT 15.00 CHECK (desconto_duas_semana >= 0 AND desconto_duas_semana <= 100),
    ativa                   BOOLEAN       NOT NULL DEFAULT TRUE,
    criado_em               TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    atualizado_em           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    UNIQUE (categoria_id, regiao_id)
);

CREATE INDEX idx_precos_categoria ON tabela_precos(categoria_id);
CREATE INDEX idx_precos_regiao    ON tabela_precos(regiao_id);

CREATE TRIGGER trg_tabela_precos_atualizado
    BEFORE UPDATE ON tabela_precos
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();
