-- =============================================================================
-- V3 — Preferências do Cliente (Favoritas e Bloqueios)
-- =============================================================================
-- Etapa: 3 (diferencial B) — Favoritas e Bloqueios
-- Resumo: permite ao cliente marcar profissionais como FAVORITA ou BLOQUEADA.
--         A tabela é única para ambos os tipos: a coluna tipo é o discriminador
--         e a UNIQUE (cliente_id, profissional_id) garante que um cliente não
--         pode simultaneamente favoritar e bloquear a mesma profissional — a
--         ação mais recente sobrescreve a anterior (upsert no service).
--
-- Delta vs. V2:
--   + enum tipo_preferencia ('FAVORITA', 'BLOQUEADA')
--   + tabela cliente_profissional_preferencia
--   + índice idx_pref_cliente_tipo (filtro por cliente + tipo ao listar)
--   + índice idx_pref_profissional (para auditoria futura: "quantas bloquearam")
--   + trigger trg_pref_atualizado reutiliza fn_set_atualizado_em() da V2
--
-- Espelhamento service ↔ banco:
--   tipo ENUM                                 → domain.TipoPreferencia.Valida()
--   UNIQUE (cliente_id, profissional_id)      → PreferenciaService.Upsert (mantém ID original)
--   FK cliente/profissional ON DELETE CASCADE → service valida existência antes
--
-- Rollback (dev apenas — produção nunca roda):
--   DROP TRIGGER IF EXISTS trg_pref_atualizado ON cliente_profissional_preferencia;
--   DROP TABLE IF EXISTS cliente_profissional_preferencia;
--   DROP TYPE  IF EXISTS tipo_preferencia;
-- =============================================================================

CREATE TYPE tipo_preferencia AS ENUM ('FAVORITA', 'BLOQUEADA');

CREATE TABLE cliente_profissional_preferencia (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cliente_id      UUID NOT NULL REFERENCES clientes(id) ON DELETE CASCADE,
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    tipo            tipo_preferencia NOT NULL,
    motivo          TEXT,
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cliente_id, profissional_id)
);

CREATE INDEX idx_pref_cliente_tipo    ON cliente_profissional_preferencia(cliente_id, tipo);
CREATE INDEX idx_pref_profissional    ON cliente_profissional_preferencia(profissional_id);

CREATE TRIGGER trg_pref_atualizado
    BEFORE UPDATE ON cliente_profissional_preferencia
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();
