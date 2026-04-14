package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// ReferenciaRepository implementa domain.ReferenciaRepository em memoria.
type ReferenciaRepository struct {
	mu              sync.RWMutex
	referencias     map[uuid.UUID]*domain.Referencia
	porProfissional map[uuid.UUID][]uuid.UUID
}

// NewReferenciaRepository cria um novo repositorio in-memory vazio.
func NewReferenciaRepository() *ReferenciaRepository {
	return &ReferenciaRepository{
		referencias:     make(map[uuid.UUID]*domain.Referencia),
		porProfissional: make(map[uuid.UUID][]uuid.UUID),
	}
}

// Criar persiste uma nova referencia.
func (r *ReferenciaRepository) Criar(ctx context.Context, ref *domain.Referencia) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copia := *ref
	copia.CriadoEm = time.Now()
	copia.AtualizadoEm = time.Now()

	r.referencias[copia.ID] = &copia
	r.porProfissional[copia.ProfissionalID] = append(r.porProfissional[copia.ProfissionalID], copia.ID)

	return nil
}

// BuscarPorID retorna uma copia da referencia pelo UUID.
func (r *ReferenciaRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Referencia, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ref, existe := r.referencias[id]
	if !existe {
		return nil, domain.ErrReferenciaNaoEncontrada
	}

	copia := *ref
	return &copia, nil
}

// ListarPorProfissionalID retorna todas as referencias de uma profissional.
func (r *ReferenciaRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Referencia, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.porProfissional[profissionalID]
	resultado := make([]*domain.Referencia, 0, len(ids))
	for _, id := range ids {
		if ref, ok := r.referencias[id]; ok {
			copia := *ref
			resultado = append(resultado, &copia)
		}
	}

	return resultado, nil
}
