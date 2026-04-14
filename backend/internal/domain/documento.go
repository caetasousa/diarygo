package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// TipoDocumento representa o tipo de documento enviado pela profissional.
type TipoDocumento string

const (
	DocRGFrente              TipoDocumento = "RG_FRENTE"
	DocRGVerso               TipoDocumento = "RG_VERSO"
	DocCPF                   TipoDocumento = "CPF"
	DocComprovanteResidencia TipoDocumento = "COMPROVANTE"
	DocFoto                  TipoDocumento = "FOTO"
	DocOutro                 TipoDocumento = "OUTRO"
)

// StatusDocumento representa o estado de analise do documento.
type StatusDocumento string

const (
	DocPendente  StatusDocumento = "PENDENTE"
	DocAprovado  StatusDocumento = "APROVADO"
	DocReprovado StatusDocumento = "REPROVADO"
)

// Documento representa um documento enviado por uma profissional.
type Documento struct {
	ID              uuid.UUID
	ProfissionalID  uuid.UUID
	Tipo            TipoDocumento
	URL             string
	Status          StatusDocumento
	ObservacaoAdmin string
	CriadoEm        time.Time
	AtualizadoEm    time.Time
}

// Erros de dominio de documento.
var (
	ErrDocumentoNaoEncontrado  = errors.New("documento nao encontrado")
	ErrTipoDocumentoInvalido   = errors.New("tipo de documento invalido")
	ErrURLDocumentoObrigatoria = errors.New("URL do documento e obrigatoria")
)

// DocumentoReader define operacoes de leitura.
type DocumentoReader interface {
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Documento, error)
	ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*Documento, error)
}

// DocumentoWriter define operacoes de escrita.
type DocumentoWriter interface {
	Criar(ctx context.Context, d *Documento) error
	Atualizar(ctx context.Context, d *Documento) error
}

// DocumentoRepository e a interface completa.
type DocumentoRepository interface {
	DocumentoReader
	DocumentoWriter
}

// DocumentoRequest representa o payload de envio de documento.
type DocumentoRequest struct {
	Tipo TipoDocumento `json:"tipo"`
	URL  string        `json:"url"`
}

// DocumentoResponse e a resposta com os dados do documento.
type DocumentoResponse struct {
	ID     uuid.UUID       `json:"id"`
	Tipo   TipoDocumento   `json:"tipo"`
	URL    string          `json:"url"`
	Status StatusDocumento `json:"status"`
}

// TiposDocumentoValidos retorna true se o tipo for reconhecido.
func TipoDocumentoValido(t TipoDocumento) bool {
	switch t {
	case DocRGFrente, DocRGVerso, DocCPF, DocComprovanteResidencia, DocFoto, DocOutro:
		return true
	}
	return false
}
