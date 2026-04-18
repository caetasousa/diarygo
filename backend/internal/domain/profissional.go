package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// StatusProfissional representa o estado do credenciamento da profissional.
type StatusProfissional string

const (
	StatusPendente       StatusProfissional = "PENDENTE"
	StatusAprovada       StatusProfissional = "APROVADA"
	StatusReprovada      StatusProfissional = "REPROVADA"
	StatusSuspensa       StatusProfissional = "SUSPENSA"
	StatusDescredenciada StatusProfissional = "DESCREDENCIADA"
)

// Profissional representa o perfil completo de uma diarista.
// Espelha a tabela `profissionais` do banco.
type Profissional struct {
	ID            uuid.UUID
	UsuarioID     uuid.UUID
	Nome          string
	CPF           string // apenas digitos
	RG            string
	Telefone      string // apenas digitos
	FotoURL       string
	MEI           bool
	Status        StatusProfissional
	NotaMedia     float64 // 0.0-5.0
	TotalServicos int
	CriadoEm      time.Time
	AtualizadoEm  time.Time
}

// Erros de dominio da profissional.
var (
	ErrProfissionalNaoEncontrada  = errors.New("profissional nao encontrada")
	ErrProfissionalJaExiste       = errors.New("profissional ja cadastrada para este usuario")
	ErrProfissionalNaoAprovada    = errors.New("profissional nao esta aprovada")
	ErrProfissionalSuspensa       = errors.New("profissional esta suspensa")
	ErrProfissionalDescredenciada = errors.New("profissional esta descredenciada")
	ErrNotaMediaInvalida          = errors.New("nota media deve estar entre 1.0 e 5.0")
	ErrTotalServicosInvalido      = errors.New("total de servicos nao pode ser negativo")
)

// ProfissionalReader define operacoes de leitura.
type ProfissionalReader interface {
	BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*Profissional, error)
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Profissional, error)
}

// ProfissionalWriter define operacoes de escrita.
type ProfissionalWriter interface {
	Criar(ctx context.Context, p *Profissional) error
	Atualizar(ctx context.Context, p *Profissional) error
}

// ProfissionalRepository e a interface completa.
type ProfissionalRepository interface {
	ProfissionalReader
	ProfissionalWriter
}

// ProfissionalRequest representa o payload de atualizacao do perfil da profissional.
type ProfissionalRequest struct {
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	RG       string `json:"rg"`
	Telefone string `json:"telefone"`
	FotoURL  string `json:"foto_url"`
	MEI      bool   `json:"mei"`
}

// ValidarNotaMedia garante 1.0..5.0 espelhando o CHECK do banco.
// Primeira linha de defesa em updates pos-avaliacao (Etapa 6).
// Nota: a entidade recem-criada tem NotaMedia=0 e nao deve passar por esta
// validacao — so apos a primeira avaliacao registrada.
func ValidarNotaMedia(nota float64) error {
	if nota < 1.0 || nota > 5.0 {
		return ErrNotaMediaInvalida
	}
	return nil
}

// ValidarTotalServicos garante >= 0 espelhando o CHECK do banco.
func ValidarTotalServicos(total int) error {
	if total < 0 {
		return ErrTotalServicosInvalido
	}
	return nil
}

// ProfissionalResponse e a resposta com os dados da profissional.
type ProfissionalResponse struct {
	ID            uuid.UUID          `json:"id"`
	UsuarioID     uuid.UUID          `json:"usuario_id"`
	Nome          string             `json:"nome"`
	CPF           string             `json:"cpf"`
	RG            string             `json:"rg"`
	Telefone      string             `json:"telefone"`
	FotoURL       string             `json:"foto_url"`
	MEI           bool               `json:"mei"`
	Status        StatusProfissional `json:"status"`
	NotaMedia     float64            `json:"nota_media"`
	TotalServicos int                `json:"total_servicos"`
}
