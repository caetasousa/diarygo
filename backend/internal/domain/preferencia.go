package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// TipoPreferencia distingue favorita de bloqueada na relacao cliente-profissional.
type TipoPreferencia string

const (
	PreferenciaFavorita  TipoPreferencia = "FAVORITA"
	PreferenciaBloqueada TipoPreferencia = "BLOQUEADA"
)

// Valida verifica se o valor e um TipoPreferencia conhecido.
func (t TipoPreferencia) Valida() bool {
	return t == PreferenciaFavorita || t == PreferenciaBloqueada
}

// ClienteProfissionalPreferencia representa a escolha do cliente em relacao
// a uma profissional: FAVORITA (preferir no matching) ou BLOQUEADA (nunca atribuir).
// Combinacao (ClienteID, ProfissionalID) e unica — mudar de FAVORITA para BLOQUEADA
// (ou vice-versa) substitui a entrada existente.
type ClienteProfissionalPreferencia struct {
	ID             uuid.UUID
	ClienteID      uuid.UUID
	ProfissionalID uuid.UUID
	Tipo           TipoPreferencia
	CriadoEm       time.Time
}

// Erros de dominio de preferencia.
var (
	ErrPreferenciaNaoEncontrada = errors.New("preferencia nao encontrada")
	ErrTipoPreferenciaInvalido  = errors.New("tipo de preferencia invalido")
)

// PreferenciaReader define operacoes de leitura do repositorio de preferencias.
type PreferenciaReader interface {
	BuscarPorClienteEProfissional(ctx context.Context, clienteID, profissionalID uuid.UUID) (*ClienteProfissionalPreferencia, error)
	ListarPorCliente(ctx context.Context, clienteID uuid.UUID, tipo TipoPreferencia) ([]*ClienteProfissionalPreferencia, error)
}

// PreferenciaWriter define operacoes de escrita do repositorio de preferencias.
type PreferenciaWriter interface {
	Upsert(ctx context.Context, p *ClienteProfissionalPreferencia) error
	Remover(ctx context.Context, clienteID, profissionalID uuid.UUID) error
}

// PreferenciaRepository e a interface completa do repositorio de preferencias.
type PreferenciaRepository interface {
	PreferenciaReader
	PreferenciaWriter
}

// PreferenciaResponse e a resposta publica da API ao listar preferencias.
type PreferenciaResponse struct {
	ProfissionalID uuid.UUID       `json:"profissional_id"`
	Nome           string          `json:"nome"`
	FotoURL        string          `json:"foto_url,omitempty"`
	NotaMedia      float64         `json:"nota_media"`
	Tipo           TipoPreferencia `json:"tipo"`
	CriadoEm       time.Time       `json:"criado_em"`
}
