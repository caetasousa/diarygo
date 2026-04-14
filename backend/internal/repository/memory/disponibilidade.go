package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// DisponibilidadeRepository implementa domain.DisponibilidadeRepository em memoria.
type DisponibilidadeRepository struct {
	mu              sync.RWMutex
	slots           map[uuid.UUID]*domain.Disponibilidade
	porProfissional map[uuid.UUID][]uuid.UUID
}

// NewDisponibilidadeRepository cria um novo repositorio in-memory vazio.
func NewDisponibilidadeRepository() *DisponibilidadeRepository {
	return &DisponibilidadeRepository{
		slots:           make(map[uuid.UUID]*domain.Disponibilidade),
		porProfissional: make(map[uuid.UUID][]uuid.UUID),
	}
}

// DefinirDisponibilidades substitui todos os slots da profissional.
func (r *DisponibilidadeRepository) DefinirDisponibilidades(ctx context.Context, profissionalID uuid.UUID, novos []*domain.Disponibilidade) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Remover slots anteriores
	for _, sid := range r.porProfissional[profissionalID] {
		delete(r.slots, sid)
	}

	// Inserir novos slots
	ids := make([]uuid.UUID, 0, len(novos))
	for _, slot := range novos {
		copia := *slot
		copia.ProfissionalID = profissionalID
		copia.CriadoEm = time.Now()
		copia.AtualizadoEm = time.Now()
		r.slots[copia.ID] = &copia
		ids = append(ids, copia.ID)
	}
	r.porProfissional[profissionalID] = ids

	return nil
}

// ListarPorProfissionalID retorna todos os slots de disponibilidade da profissional.
func (r *DisponibilidadeRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Disponibilidade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.porProfissional[profissionalID]
	resultado := make([]*domain.Disponibilidade, 0, len(ids))
	for _, id := range ids {
		if s, ok := r.slots[id]; ok {
			copia := *s
			resultado = append(resultado, &copia)
		}
	}

	return resultado, nil
}
