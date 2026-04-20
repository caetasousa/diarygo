package domain

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// StatusSolicitacao representa o estado do ciclo de vida de uma solicitação.
type StatusSolicitacao string

const (
	StatusAguardando  StatusSolicitacao = "AGUARDANDO"
	StatusAtribuida   StatusSolicitacao = "ATRIBUIDA"
	StatusConfirmada  StatusSolicitacao = "CONFIRMADA"
	StatusEmAndamento StatusSolicitacao = "EM_ANDAMENTO"
	StatusConcluida   StatusSolicitacao = "CONCLUIDA"
	StatusCancelada   StatusSolicitacao = "CANCELADA"
)

// Valida confere se o status é um valor conhecido.
func (s StatusSolicitacao) Valida() bool {
	switch s {
	case StatusAguardando, StatusAtribuida, StatusConfirmada,
		StatusEmAndamento, StatusConcluida, StatusCancelada:
		return true
	}
	return false
}

// cancelaveis lista os status a partir dos quais o cliente pode cancelar.
var cancelaveis = map[StatusSolicitacao]bool{
	StatusAguardando: true,
	StatusAtribuida:  true,
	StatusConfirmada: true,
}

// PodeCancelar informa se uma solicitação neste status aceita cancelamento.
func (s StatusSolicitacao) PodeCancelar() bool {
	return cancelaveis[s]
}

// Solicitacao congela o orçamento no momento da criação (breakdown, valor, duração).
// O snapshot é imutável — se a tabela de preços mudar depois, este registro
// continua refletindo o que foi combinado com o cliente.
type Solicitacao struct {
	ID           uuid.UUID
	ClienteID    uuid.UUID
	EnderecoID   uuid.UUID
	CategoriaID  uuid.UUID
	RegiaoID     uuid.UUID
	NumQuartos   int
	NumBanheiros int
	NumSalas     int
	NumCozinhas  int
	Frequencia   FrequenciaServico
	DataServico  time.Time
	Observacao   string // até 500 runas

	// Snapshot congelado do cálculo (Etapa 3).
	ValorTotal float64
	DuracaoMin int
	Breakdown  []ItemCalculo // runtime; serializado como JSONB no Postgres

	Status       StatusSolicitacao
	CriadaEm     time.Time
	AtualizadaEm time.Time
	CanceladaEm  *time.Time
}

// SolicitacaoOpcional é a tabela associativa com snapshot do opcional no momento da criação.
type SolicitacaoOpcional struct {
	SolicitacaoID uuid.UUID
	OpcionalID    uuid.UUID
	NomeSnapshot  string
	ValorSnapshot float64
	TempoSnapshot int
}

// SolicitacaoRepository define operações de persistência de solicitações.
type SolicitacaoRepository interface {
	Criar(ctx context.Context, s *Solicitacao, ops []SolicitacaoOpcional) error
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Solicitacao, []SolicitacaoOpcional, error)
	ListarPorCliente(ctx context.Context, clienteID uuid.UUID, filtro SolicitacaoFiltro) ([]*Solicitacao, error)
	AtualizarStatus(ctx context.Context, id uuid.UUID, status StatusSolicitacao, canceladaEm *time.Time) error
}

// SolicitacaoFiltro filtra a listagem de solicitações.
type SolicitacaoFiltro struct {
	Status *StatusSolicitacao
	Desde  *time.Time
	Ate    *time.Time
}

// Erros de domínio de solicitação.
var (
	ErrSolicitacaoNaoEncontrada    = errors.New("solicitacao nao encontrada")
	ErrSolicitacaoAntecedencia     = errors.New("agendamento exige minimo de 24h de antecedencia")
	ErrSolicitacaoStatusInvalido   = errors.New("status invalido para essa operacao")
	ErrSolicitacaoEnderecoInvalido = errors.New("endereco nao pertence ao cliente")
	ErrObservacaoMuitoLonga        = errors.New("observacao nao pode exceder 500 caracteres")
)

// ValidarAntecedencia24h valida se a data do serviço está pelo menos 24h no futuro.
// Recebe agora como parâmetro para testabilidade.
func ValidarAntecedencia24h(dataServico, agora time.Time) error {
	if dataServico.Sub(agora) < 24*time.Hour {
		return ErrSolicitacaoAntecedencia
	}
	return nil
}

// ValidarObservacao valida o tamanho em runas (não bytes).
func ValidarObservacao(obs string) error {
	if utf8.RuneCountInString(obs) > 500 {
		return ErrObservacaoMuitoLonga
	}
	return nil
}

// CriarSolicitacaoRequest é o payload de criação de solicitação.
type CriarSolicitacaoRequest struct {
	EnderecoID   uuid.UUID         `json:"endereco_id"`
	CategoriaID  uuid.UUID         `json:"categoria_id"`
	RegiaoID     uuid.UUID         `json:"regiao_id"`
	NumQuartos   int               `json:"num_quartos"`
	NumBanheiros int               `json:"num_banheiros"`
	NumSalas     int               `json:"num_salas"`
	NumCozinhas  int               `json:"num_cozinhas"`
	OpcionaisIDs []uuid.UUID       `json:"opcionais_ids"`
	Frequencia   FrequenciaServico `json:"frequencia"`
	DataServico  time.Time         `json:"data_servico"`
	Observacao   string            `json:"observacao"`
}

// SolicitacaoResponse é o DTO de resposta de solicitação.
type SolicitacaoResponse struct {
	ID           uuid.UUID         `json:"id"`
	EnderecoID   uuid.UUID         `json:"endereco_id"`
	CategoriaID  uuid.UUID         `json:"categoria_id"`
	RegiaoID     uuid.UUID         `json:"regiao_id"`
	NumQuartos   int               `json:"num_quartos"`
	NumBanheiros int               `json:"num_banheiros"`
	NumSalas     int               `json:"num_salas"`
	NumCozinhas  int               `json:"num_cozinhas"`
	Frequencia   FrequenciaServico `json:"frequencia"`
	DataServico  time.Time         `json:"data_servico"`
	Observacao   string            `json:"observacao"`
	ValorTotal   float64           `json:"valor_total"`
	DuracaoMin   int               `json:"duracao_min"`
	Breakdown    []ItemCalculo     `json:"breakdown"`
	Status       StatusSolicitacao `json:"status"`
	CriadaEm     time.Time         `json:"criada_em"`
	CanceladaEm  *time.Time        `json:"cancelada_em,omitempty"`
}
