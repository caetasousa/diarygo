package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// CategoriaServico representa um tipo de servico oferecido (faxina padrao,
// faxina pesada, pos-obra, etc). Profissional pode habilitar uma ou mais.
type CategoriaServico struct {
	ID               uuid.UUID
	Nome             string
	Descricao        string
	DuracaoMinimaMin int // duracao minima em minutos (usada no calculo de horas)
	Ativa            bool
	CriadaEm         time.Time
	AtualizadaEm     time.Time
}

// Opcional representa um acrescimo que pode ser adicionado ao servico
// (passadoria, limpeza de geladeira, janelas externas, etc).
type Opcional struct {
	ID            uuid.UUID
	Nome          string
	Descricao     string
	ValorExtra    float64 // reais
	TempoExtraMin int
	Ativo         bool
	CriadoEm      time.Time
	AtualizadoEm  time.Time
}

// TabelaPrecos representa o preco-hora e fatores multiplicadores para uma
// combinacao categoria × regiao. Persiste descontos por frequencia e
// acrescimo de fim-de-semana como percentuais.
type TabelaPrecos struct {
	ID                 uuid.UUID
	CategoriaID        uuid.UUID
	RegiaoID           uuid.UUID
	PrecoHora          float64 // R$/hora
	AcrescimoFDS       float64 // %  (ex: 15.0 -> +15% em sabado/domingo)
	DescontoSemanal    float64 // %  (ex: 10.0 -> -10%)
	DescontoQuinzenal  float64 // %
	DescontoDuasSemana float64 // % (2x por semana)
	Ativa              bool
	CriadaEm           time.Time
	AtualizadaEm       time.Time
}

// FrequenciaServico categoriza a cadencia do servico para calculo de desconto.
type FrequenciaServico string

const (
	FrequenciaUnica         FrequenciaServico = "UNICA"
	FrequenciaSemanal       FrequenciaServico = "SEMANAL"
	FrequenciaQuinzenal     FrequenciaServico = "QUINZENAL"
	FrequenciaDuasPorSemana FrequenciaServico = "DUAS_POR_SEMANA"
)

// Valida confere se o valor e uma frequencia conhecida.
func (f FrequenciaServico) Valida() bool {
	switch f {
	case FrequenciaUnica, FrequenciaSemanal, FrequenciaQuinzenal, FrequenciaDuasPorSemana:
		return true
	}
	return false
}

// Erros de dominio do catalogo.
var (
	ErrCategoriaNaoEncontrada   = errors.New("categoria nao encontrada")
	ErrOpcionalNaoEncontrado    = errors.New("opcional nao encontrado")
	ErrTabelaPrecoNaoEncontrada = errors.New("tabela de precos nao encontrada para categoria+regiao")
	ErrFrequenciaInvalida       = errors.New("frequencia invalida")
	ErrDataServicoPassado       = errors.New("data do servico nao pode ser no passado")
	ErrComodosInvalidosCalculo  = errors.New("numero de comodos invalido para calculo")
)

// CalculoPrecoRequest e o payload do cliente para simular um orcamento.
type CalculoPrecoRequest struct {
	CategoriaID  uuid.UUID         `json:"categoria_id"`
	RegiaoID     uuid.UUID         `json:"regiao_id"`
	NumQuartos   int               `json:"num_quartos"`
	NumBanheiros int               `json:"num_banheiros"`
	NumSalas     int               `json:"num_salas"`
	NumCozinhas  int               `json:"num_cozinhas"`
	AreaM2       float64           `json:"area_m2"`
	OpcionaisIDs []uuid.UUID       `json:"opcionais_ids"`
	Frequencia   FrequenciaServico `json:"frequencia"`
	DataServico  time.Time         `json:"data_servico"`
}

// ItemCalculo e uma linha do breakdown exibido no orcamento transparente.
// Tipo pode ser: BASE, COMODO, OPCIONAL, ACRESCIMO, DESCONTO, TOTAL.
type ItemCalculo struct {
	Tipo  string  `json:"tipo"`
	Label string  `json:"label"`
	Valor float64 `json:"valor"`
}

// CalculoPrecoResponse e o breakdown completo devolvido ao cliente.
type CalculoPrecoResponse struct {
	CategoriaID uuid.UUID     `json:"categoria_id"`
	RegiaoID    uuid.UUID     `json:"regiao_id"`
	DuracaoMin  int           `json:"duracao_min"`
	ValorTotal  float64       `json:"valor_total"`
	Itens       []ItemCalculo `json:"itens"`
}

// CategoriaResponse e a resposta publica de /categorias.
type CategoriaResponse struct {
	ID               uuid.UUID `json:"id"`
	Nome             string    `json:"nome"`
	Descricao        string    `json:"descricao"`
	DuracaoMinimaMin int       `json:"duracao_minima_min"`
	Ativa            bool      `json:"ativa"`
}

// OpcionalResponse e a resposta publica de /categorias/:id/opcionais.
type OpcionalResponse struct {
	ID            uuid.UUID `json:"id"`
	Nome          string    `json:"nome"`
	Descricao     string    `json:"descricao"`
	ValorExtra    float64   `json:"valor_extra"`
	TempoExtraMin int       `json:"tempo_extra_min"`
	Ativo         bool      `json:"ativo"`
}

// CatalogoRepository expoe catalogo + precos de forma agregada.
// Mantido como interface unica para in-memory simplificar; em Postgres,
// cada entidade vira um repositorio proprio.
type CatalogoRepository interface {
	ListarCategoriasAtivas(ctx context.Context) ([]*CategoriaServico, error)
	BuscarCategoriaPorID(ctx context.Context, id uuid.UUID) (*CategoriaServico, error)
	ListarOpcionaisAtivos(ctx context.Context) ([]*Opcional, error)
	BuscarOpcionalPorID(ctx context.Context, id uuid.UUID) (*Opcional, error)
	BuscarTabelaPreco(ctx context.Context, categoriaID, regiaoID uuid.UUID) (*TabelaPrecos, error)
}
