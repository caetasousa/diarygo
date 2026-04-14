package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// ProfissionalRepository implementa domain.ProfissionalRepository em memoria.
type ProfissionalRepository struct {
	mu            sync.RWMutex
	profissionais map[uuid.UUID]*domain.Profissional
	porUsuario    map[uuid.UUID]uuid.UUID // usuarioID -> profissionalID
}

// NewProfissionalRepository cria um novo repositorio in-memory vazio.
func NewProfissionalRepository() *ProfissionalRepository {
	return &ProfissionalRepository{
		profissionais: make(map[uuid.UUID]*domain.Profissional),
		porUsuario:    make(map[uuid.UUID]uuid.UUID),
	}
}

// Criar persiste uma nova profissional.
func (r *ProfissionalRepository) Criar(ctx context.Context, p *domain.Profissional) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, existe := r.porUsuario[p.UsuarioID]; existe {
		return domain.ErrProfissionalJaExiste
	}

	copia := *p
	copia.CriadoEm = time.Now()
	copia.AtualizadoEm = time.Now()

	r.profissionais[copia.ID] = &copia
	r.porUsuario[copia.UsuarioID] = copia.ID

	return nil
}

// BuscarPorUsuarioID retorna uma copia da profissional pelo ID do usuario.
func (r *ProfissionalRepository) BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*domain.Profissional, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, existe := r.porUsuario[usuarioID]
	if !existe {
		return nil, domain.ErrProfissionalNaoEncontrada
	}

	p := r.profissionais[id]
	copia := *p
	return &copia, nil
}

// BuscarPorID retorna uma copia da profissional pelo UUID.
func (r *ProfissionalRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Profissional, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, existe := r.profissionais[id]
	if !existe {
		return nil, domain.ErrProfissionalNaoEncontrada
	}

	copia := *p
	return &copia, nil
}

// Atualizar substitui a profissional existente.
func (r *ProfissionalRepository) Atualizar(ctx context.Context, p *domain.Profissional) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	antigo, existe := r.profissionais[p.ID]
	if !existe {
		return domain.ErrProfissionalNaoEncontrada
	}

	copia := *p
	copia.CriadoEm = antigo.CriadoEm
	copia.AtualizadoEm = time.Now()
	r.profissionais[copia.ID] = &copia

	return nil
}
