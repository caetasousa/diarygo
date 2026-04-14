package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// DocumentoRepository implementa domain.DocumentoRepository em memoria.
type DocumentoRepository struct {
	mu              sync.RWMutex
	documentos      map[uuid.UUID]*domain.Documento
	porProfissional map[uuid.UUID][]uuid.UUID // profissionalID -> []documentoID
}

// NewDocumentoRepository cria um novo repositorio in-memory vazio.
func NewDocumentoRepository() *DocumentoRepository {
	return &DocumentoRepository{
		documentos:      make(map[uuid.UUID]*domain.Documento),
		porProfissional: make(map[uuid.UUID][]uuid.UUID),
	}
}

// Criar persiste um novo documento.
func (r *DocumentoRepository) Criar(ctx context.Context, d *domain.Documento) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copia := *d
	copia.CriadoEm = time.Now()
	copia.AtualizadoEm = time.Now()

	r.documentos[copia.ID] = &copia
	r.porProfissional[copia.ProfissionalID] = append(r.porProfissional[copia.ProfissionalID], copia.ID)

	return nil
}

// BuscarPorID retorna uma copia do documento pelo UUID.
func (r *DocumentoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Documento, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, existe := r.documentos[id]
	if !existe {
		return nil, domain.ErrDocumentoNaoEncontrado
	}

	copia := *d
	return &copia, nil
}

// ListarPorProfissionalID retorna todos os documentos de uma profissional.
func (r *DocumentoRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Documento, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.porProfissional[profissionalID]
	resultado := make([]*domain.Documento, 0, len(ids))
	for _, id := range ids {
		if d, ok := r.documentos[id]; ok {
			copia := *d
			resultado = append(resultado, &copia)
		}
	}

	return resultado, nil
}

// Atualizar substitui o documento existente.
func (r *DocumentoRepository) Atualizar(ctx context.Context, d *domain.Documento) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	antigo, existe := r.documentos[d.ID]
	if !existe {
		return domain.ErrDocumentoNaoEncontrado
	}

	copia := *d
	copia.CriadoEm = antigo.CriadoEm
	copia.AtualizadoEm = time.Now()
	r.documentos[copia.ID] = &copia

	return nil
}
