package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// StatusReferencia representa o estado de confirmacao de uma referencia.
type StatusReferencia string

const (
	RefPendente      StatusReferencia = "PENDENTE"
	RefConfirmada    StatusReferencia = "CONFIRMADA"
	RefNaoConfirmada StatusReferencia = "NAO_CONFIRMADA"
)

// Referencia representa uma referencia profissional fornecida pela diarista.
type Referencia struct {
	ID              uuid.UUID
	ProfissionalID  uuid.UUID
	NomeContato     string
	TelefoneContato string // apenas digitos
	Status          StatusReferencia
	CriadoEm        time.Time
	AtualizadoEm    time.Time
}

// Erros de dominio de referencia.
var (
	ErrReferenciaNaoEncontrada    = errors.New("referencia nao encontrada")
	ErrNomeContatoObrigatorio     = errors.New("nome do contato e obrigatorio")
	ErrTelefoneContatoObrigatorio = errors.New("telefone do contato e obrigatorio")
)

// ReferenciaReader define operacoes de leitura.
type ReferenciaReader interface {
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Referencia, error)
	ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*Referencia, error)
}

// ReferenciaWriter define operacoes de escrita.
type ReferenciaWriter interface {
	Criar(ctx context.Context, r *Referencia) error
}

// ReferenciaRepository e a interface completa.
type ReferenciaRepository interface {
	ReferenciaReader
	ReferenciaWriter
}

// ReferenciaRequest representa o payload de adicao de referencia.
type ReferenciaRequest struct {
	NomeContato     string `json:"nome_contato"`
	TelefoneContato string `json:"telefone_contato"`
}

// ReferenciaResponse e a resposta com os dados da referencia.
type ReferenciaResponse struct {
	ID              uuid.UUID        `json:"id"`
	NomeContato     string           `json:"nome_contato"`
	TelefoneContato string           `json:"telefone_contato"`
	Status          StatusReferencia `json:"status"`
}
