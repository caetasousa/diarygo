package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Regiao representa uma area geografica de atuacao cadastrada no sistema.
type Regiao struct {
	ID        uuid.UUID
	Nome      string
	Cidade    string
	Estado    string
	CEPInicio string // 8 digitos
	CEPFim    string // 8 digitos
	Ativa     bool
	CriadoEm  time.Time
}

// Erros de dominio de regiao.
var (
	ErrRegiaoNaoEncontrada = errors.New("regiao nao encontrada")
	ErrRegiaoInativa       = errors.New("regiao nao esta ativa")
)

// RegiaoReader define operacoes de leitura.
type RegiaoReader interface {
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Regiao, error)
	ListarAtivas(ctx context.Context) ([]*Regiao, error)
}

// RegiaoWriter define operacoes de escrita.
type RegiaoWriter interface {
	Criar(ctx context.Context, r *Regiao) error
}

// RegiaoRepository e a interface completa.
type RegiaoRepository interface {
	RegiaoReader
	RegiaoWriter
}

// RegiaoResponse e a resposta com os dados da regiao.
type RegiaoResponse struct {
	ID        uuid.UUID `json:"id"`
	Nome      string    `json:"nome"`
	Cidade    string    `json:"cidade"`
	Estado    string    `json:"estado"`
	CEPInicio string    `json:"cep_inicio"`
	CEPFim    string    `json:"cep_fim"`
}

// ProfissionalRegiaoRepository gerencia a associacao N:N profissional <-> regioes.
type ProfissionalRegiaoRepository interface {
	DefinirRegioes(ctx context.Context, profissionalID uuid.UUID, regiaoIDs []uuid.UUID) error
	ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*Regiao, error)
}

// DefinirRegioesRequest representa o payload para definir regioes de atuacao.
type DefinirRegioesRequest struct {
	RegiaoIDs []uuid.UUID `json:"regiao_ids"`
}

// Erros adicionais para associacao de regioes.
var ErrRegiaoIDInvalida = errors.New("regiao ID invalida")
