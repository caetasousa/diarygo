-- =============================================================================
-- V2 — Cadastro de Clientes e Profissionais
-- =============================================================================
-- Etapa: 2 — Cadastro completo de clientes e profissionais
-- Resumo: cria o perfil completo de clientes e profissionais com suas
--         dependências imediatas — endereços do cliente, documentos e
--         referências do profissional, regiões de atuação e horários de
--         disponibilidade. Após esta migration é possível completar o
--         onboarding de uma conta criada pela V1.
--
-- Princípio de validação — defesa em profundidade:
--   O banco é a ÚLTIMA linha de defesa; o service (internal/service/*) é a
--   PRIMEIRA. Cada CHECK/UNIQUE/ENUM abaixo TEM contraparte no service, que
--   valida, normaliza (trim, upper, remove máscaras) e retorna erro de
--   domínio amigável antes do SQL rodar. O CHECK só dispara se o service
--   for contornado (bug, migration futura, acesso direto ao banco).
--
--   Espelhamento service ↔ banco:
--     clientes.cpf CHECK '^[0-9]{11}$'         → domain.ValidarCPF (dígitos verificadores + 11 chars)
--     clientes.score 0..100                    → domain.ValidarScore (aplicado em updates futuros)
--     profissionais.cpf CHECK '^[0-9]{11}$'    → domain.ValidarCPF
--     profissionais.nota_media 1.0..5.0 | NULL → domain.ValidarNotaMedia (atualização pós-avaliação)
--     profissionais.total_servicos >= 0        → domain.ValidarTotalServicos
--     profissionais.status ENUM                → constantes StatusPendente/Aprovada/...
--     enderecos.cep CHECK '^[0-9]{8}$'         → domain.ValidarCEP
--     enderecos.cep faixa 74000000..74999999   → domain.ValidarAreaAtendimento (MVP Goiânia)
--     enderecos.cidade/estado = Goiânia/GO     → domain.ValidarAreaAtendimento
--     enderecos.num_* >= 0                     → domain.ValidarComodos
--     enderecos.area_m2 > 0 | NULL             → domain.ValidarAreaM2
--     enderecos UNIQUE (principal por cliente) → EnderecoService.desmarcarPrincipal
--     regioes.cep_fim >= cep_inicio            → domain.ValidarFaixaCEP
--     regioes.estado CHAR(2)                   → domain.ValidarEstado
--     regioes UNIQUE (nome,cidade,estado)      → RegiaoService checa antes de criar
--     regioes_atuacao PK composta              → CredenciamentoService.DefinirRegioes dedupe
--     disponibilidades.dia_semana 0..6         → checado no service (DefinirDisponibilidades)
--     disponibilidades.hora_fim > hora_inicio  → checado no service
--     disponibilidades UNIQUE (slots dup.)     → CredenciamentoService.DefinirDisponibilidades dedupe
--     documentos.tipo ENUM                     → domain.TipoDocumentoValido
--     documentos.status ENUM                   → constantes DocPendente/Aprovado/Reprovado
--     documentos.analisado_por only if !PEND.  → checado no AdminService (Etapa 10)
--
-- Delta vs. V1:
--   + Enums:    status_profissional, status_documento, status_referencia,
--               tipo_documento
--   + Tabelas:  clientes, profissionais, enderecos, documentos, referencias,
--               regioes, regioes_atuacao, disponibilidades
--   + Índices em FKs (requisito para leitura em escala — Postgres não cria
--     índice em coluna FK automaticamente):
--               idx_enderecos_cliente, idx_documentos_profissional,
--               idx_referencias_profissional,
--               idx_disponibilidades_profissional,
--               idx_regioes_atuacao_regiao
--   + Índice por status: idx_profissionais_status (matching por estado)
--   + Índice único parcial: ux_enderecos_cliente_principal
--     (máximo 1 endereço principal por cliente — invariância garantida no DB)
--   + CHECKs de formato (CPF, CEP) e de consistência de horário
--   + Triggers atualizado_em em todas as tabelas mutáveis
--   Observação: NÃO há índices em (email), (clientes.cpf), (profissionais.cpf)
--     porque o UNIQUE de cada coluna já cria B-tree automático — índice
--     adicional apenas duplicaria custo de escrita/vacuum.
--
-- Observação: `documentos.analisado_por` fica como UUID sem FK até a V10
--   (Painel Admin) criar a tabela `admins`. Quando V10 rodar, um ALTER TABLE
--   adicionará a constraint REFERENCES admins(id).
--
-- Rollback (apenas em dev — nunca em produção):
--   DROP TABLE IF EXISTS disponibilidades, regioes_atuacao, regioes,
--                        referencias, documentos, enderecos,
--                        profissionais, clientes CASCADE;
--   DROP TYPE IF EXISTS tipo_documento, status_referencia,
--                       status_documento, status_profissional;
-- =============================================================================

CREATE TYPE status_profissional AS ENUM (
    'PENDENTE', 'APROVADA', 'REPROVADA', 'SUSPENSA', 'DESCREDENCIADA'
);
CREATE TYPE status_documento AS ENUM ('PENDENTE', 'APROVADO', 'REPROVADO');
CREATE TYPE status_referencia AS ENUM ('PENDENTE', 'CONFIRMADA', 'NAO_CONFIRMADA');
CREATE TYPE tipo_documento AS ENUM (
    'RG_FRENTE', 'RG_VERSO', 'CPF', 'COMPROVANTE', 'FOTO', 'OUTRO'
);

-- -----------------------------------------------------------------------------
-- clientes
-- -----------------------------------------------------------------------------
CREATE TABLE clientes (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id    UUID NOT NULL UNIQUE REFERENCES usuarios(id) ON DELETE CASCADE,
    nome          VARCHAR(255) NOT NULL,
    cpf           CHAR(11) NOT NULL UNIQUE
                  CHECK (cpf ~ '^[0-9]{11}$'),
    telefone      VARCHAR(20),
    score         SMALLINT NOT NULL DEFAULT 100
                  CHECK (score >= 0 AND score <= 100),
    criado_em     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- profissionais
-- -----------------------------------------------------------------------------
-- nota_media: NULL significa "sem avaliação ainda". Só recebe valor (1.0-5.0)
--   após a primeira avaliação ser registrada (Etapa 6).
CREATE TABLE profissionais (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id     UUID NOT NULL UNIQUE REFERENCES usuarios(id) ON DELETE CASCADE,
    nome           VARCHAR(255) NOT NULL,
    cpf            CHAR(11) NOT NULL UNIQUE
                   CHECK (cpf ~ '^[0-9]{11}$'),
    rg             VARCHAR(20),
    telefone       VARCHAR(20),
    foto_url       VARCHAR(500),
    status         status_profissional NOT NULL DEFAULT 'PENDENTE',
    nota_media     NUMERIC(3,2)
                   CHECK (nota_media IS NULL
                          OR (nota_media >= 1.0 AND nota_media <= 5.0)),
    total_servicos INTEGER NOT NULL DEFAULT 0
                   CHECK (total_servicos >= 0),
    mei            BOOLEAN NOT NULL DEFAULT FALSE,
    criado_em      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profissionais_status ON profissionais(status);

-- -----------------------------------------------------------------------------
-- enderecos (do cliente)
-- -----------------------------------------------------------------------------
CREATE TABLE enderecos (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cliente_id    UUID NOT NULL REFERENCES clientes(id) ON DELETE CASCADE,
    logradouro    VARCHAR(255) NOT NULL,
    numero        VARCHAR(20),
    complemento   VARCHAR(100),
    bairro        VARCHAR(100) NOT NULL,
    cidade        VARCHAR(100) NOT NULL,
    estado        CHAR(2) NOT NULL,
    cep           CHAR(8) NOT NULL CHECK (cep ~ '^[0-9]{8}$'),
    lat           NUMERIC(10,7),
    lon           NUMERIC(10,7),
    principal     BOOLEAN NOT NULL DEFAULT FALSE,
    num_quartos   SMALLINT NOT NULL DEFAULT 0 CHECK (num_quartos   >= 0),
    num_banheiros SMALLINT NOT NULL DEFAULT 0 CHECK (num_banheiros >= 0),
    num_salas     SMALLINT NOT NULL DEFAULT 0 CHECK (num_salas     >= 0),
    num_cozinhas  SMALLINT NOT NULL DEFAULT 0 CHECK (num_cozinhas  >= 0),
    area_m2       NUMERIC(8,2) CHECK (area_m2 IS NULL OR area_m2 > 0),
    criado_em     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_enderecos_cliente ON enderecos(cliente_id);

-- Invariância: no máximo 1 endereço principal por cliente.
-- Índice parcial → custo zero em linhas não-principais.
CREATE UNIQUE INDEX ux_enderecos_cliente_principal
    ON enderecos(cliente_id) WHERE principal = TRUE;

-- -----------------------------------------------------------------------------
-- documentos (da profissional)
-- -----------------------------------------------------------------------------
-- analisado_por: UUID sem FK — constraint será adicionada na V10 quando
-- a tabela admins existir. Ver observação no cabeçalho.
CREATE TABLE documentos (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    tipo            tipo_documento NOT NULL,
    url             VARCHAR(500) NOT NULL,
    status          status_documento NOT NULL DEFAULT 'PENDENTE',
    analisado_por   UUID,
    observacao      TEXT,
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Só popula analisado_por quando saiu do estado PENDENTE.
    CONSTRAINT ck_documentos_analisado_por_status
        CHECK (analisado_por IS NULL OR status <> 'PENDENTE')
);

CREATE INDEX idx_documentos_profissional ON documentos(profissional_id);

-- -----------------------------------------------------------------------------
-- referencias (da profissional)
-- -----------------------------------------------------------------------------
CREATE TABLE referencias (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profissional_id  UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    nome_contato     VARCHAR(255) NOT NULL,
    telefone_contato VARCHAR(20) NOT NULL,
    status           status_referencia NOT NULL DEFAULT 'PENDENTE',
    criado_em        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_referencias_profissional ON referencias(profissional_id);

-- -----------------------------------------------------------------------------
-- regioes de atendimento
-- -----------------------------------------------------------------------------
CREATE TABLE regioes (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome          VARCHAR(255) NOT NULL,
    cidade        VARCHAR(255) NOT NULL,
    estado        CHAR(2) NOT NULL,
    cep_inicio    CHAR(8) CHECK (cep_inicio IS NULL OR cep_inicio ~ '^[0-9]{8}$'),
    cep_fim       CHAR(8) CHECK (cep_fim    IS NULL OR cep_fim    ~ '^[0-9]{8}$'),
    ativa         BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_regioes_cep_ordem
        CHECK (cep_inicio IS NULL OR cep_fim IS NULL OR cep_fim >= cep_inicio),
    CONSTRAINT uq_regioes_nome_cidade_estado UNIQUE (nome, cidade, estado)
);

-- -----------------------------------------------------------------------------
-- regioes_atuacao (N:N profissional <-> regiao)
-- -----------------------------------------------------------------------------
-- PK composta (profissional_id, regiao_id) cobre lookups por profissional.
-- Para lookups por regiao (ex.: matching) precisamos de índice explícito.
CREATE TABLE regioes_atuacao (
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    regiao_id       UUID NOT NULL REFERENCES regioes(id) ON DELETE CASCADE,
    PRIMARY KEY (profissional_id, regiao_id)
);

CREATE INDEX idx_regioes_atuacao_regiao ON regioes_atuacao(regiao_id);

-- -----------------------------------------------------------------------------
-- disponibilidades (slots semanais da profissional)
-- -----------------------------------------------------------------------------
CREATE TABLE disponibilidades (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    dia_semana      SMALLINT NOT NULL CHECK (dia_semana >= 0 AND dia_semana <= 6),
    hora_inicio     TIME NOT NULL,
    hora_fim        TIME NOT NULL,
    criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_disponibilidades_hora CHECK (hora_fim > hora_inicio),
    -- Impede dois slots idênticos para a mesma profissional no mesmo dia.
    CONSTRAINT uq_disponibilidades_prof_dia_inicio_fim
        UNIQUE (profissional_id, dia_semana, hora_inicio, hora_fim)
);

CREATE INDEX idx_disponibilidades_profissional ON disponibilidades(profissional_id);

-- -----------------------------------------------------------------------------
-- Triggers de atualizado_em
-- -----------------------------------------------------------------------------
CREATE TRIGGER trg_clientes_atualizado
    BEFORE UPDATE ON clientes
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_profissionais_atualizado
    BEFORE UPDATE ON profissionais
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_enderecos_atualizado
    BEFORE UPDATE ON enderecos
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_documentos_atualizado
    BEFORE UPDATE ON documentos
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_referencias_atualizado
    BEFORE UPDATE ON referencias
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();

CREATE TRIGGER trg_regioes_atualizado
    BEFORE UPDATE ON regioes
    FOR EACH ROW EXECUTE FUNCTION fn_set_atualizado_em();
