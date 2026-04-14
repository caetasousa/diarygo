package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// RegiaoRepository implementa domain.RegiaoRepository em memoria.
type RegiaoRepository struct {
	mu      sync.RWMutex
	regioes map[uuid.UUID]*domain.Regiao
}

// NewRegiaoRepository cria um novo repositorio in-memory com seed das regioes iniciais.
func NewRegiaoRepository() *RegiaoRepository {
	r := &RegiaoRepository{
		regioes: make(map[uuid.UUID]*domain.Regiao),
	}
	r.seed()
	return r
}

// seed carrega regioes de exemplo para o MVP.
func (r *RegiaoRepository) seed() {
	regioes := []*domain.Regiao{
		{
			ID: uuid.New(), Nome: "Centro - SP", Cidade: "São Paulo", Estado: "SP",
			CEPInicio: "01000000", CEPFim: "01499999", Ativa: true, CriadoEm: time.Now(),
		},
		{
			ID: uuid.New(), Nome: "Zona Sul - SP", Cidade: "São Paulo", Estado: "SP",
			CEPInicio: "04000000", CEPFim: "04999999", Ativa: true, CriadoEm: time.Now(),
		},
		{
			ID: uuid.New(), Nome: "Zona Norte - SP", Cidade: "São Paulo", Estado: "SP",
			CEPInicio: "02000000", CEPFim: "02999999", Ativa: true, CriadoEm: time.Now(),
		},
		{
			ID: uuid.New(), Nome: "Centro - RJ", Cidade: "Rio de Janeiro", Estado: "RJ",
			CEPInicio: "20000000", CEPFim: "20999999", Ativa: true, CriadoEm: time.Now(),
		},
		{
			ID: uuid.New(), Nome: "Belo Horizonte", Cidade: "Belo Horizonte", Estado: "MG",
			CEPInicio: "30000000", CEPFim: "30999999", Ativa: true, CriadoEm: time.Now(),
		},
	}
	for _, reg := range regioes {
		r.regioes[reg.ID] = reg
	}
}

// Criar persiste uma nova regiao.
func (r *RegiaoRepository) Criar(ctx context.Context, reg *domain.Regiao) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copia := *reg
	copia.CriadoEm = time.Now()
	r.regioes[copia.ID] = &copia

	return nil
}

// BuscarPorID retorna uma copia da regiao pelo UUID.
func (r *RegiaoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Regiao, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reg, existe := r.regioes[id]
	if !existe {
		return nil, domain.ErrRegiaoNaoEncontrada
	}

	copia := *reg
	return &copia, nil
}

// ListarAtivas retorna todas as regioes ativas.
func (r *RegiaoRepository) ListarAtivas(ctx context.Context) ([]*domain.Regiao, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resultado := make([]*domain.Regiao, 0)
	for _, reg := range r.regioes {
		if reg.Ativa {
			copia := *reg
			resultado = append(resultado, &copia)
		}
	}

	return resultado, nil
}

// ProfissionalRegiaoRepository implementa domain.ProfissionalRegiaoRepository em memoria.
type ProfissionalRegiaoRepository struct {
	mu          sync.RWMutex
	associacoes map[uuid.UUID][]uuid.UUID // profissionalID -> []regiaoID
	regiaoRepo  *RegiaoRepository
}

// NewProfissionalRegiaoRepository cria um novo repositorio in-memory.
func NewProfissionalRegiaoRepository(regiaoRepo *RegiaoRepository) *ProfissionalRegiaoRepository {
	return &ProfissionalRegiaoRepository{
		associacoes: make(map[uuid.UUID][]uuid.UUID),
		regiaoRepo:  regiaoRepo,
	}
}

// DefinirRegioes substitui as regioes de atuacao da profissional.
func (r *ProfissionalRegiaoRepository) DefinirRegioes(ctx context.Context, profissionalID uuid.UUID, regiaoIDs []uuid.UUID) error {
	// Validar que todas as regioes existem
	for _, rid := range regiaoIDs {
		if _, err := r.regiaoRepo.BuscarPorID(ctx, rid); err != nil {
			return domain.ErrRegiaoNaoEncontrada
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Copia defensiva do slice
	copia := make([]uuid.UUID, len(regiaoIDs))
	copy(copia, regiaoIDs)
	r.associacoes[profissionalID] = copia

	return nil
}

// ListarPorProfissionalID retorna as regioes de atuacao da profissional.
func (r *ProfissionalRegiaoRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Regiao, error) {
	r.mu.RLock()
	ids := r.associacoes[profissionalID]
	copia := make([]uuid.UUID, len(ids))
	copy(copia, ids)
	r.mu.RUnlock()

	resultado := make([]*domain.Regiao, 0, len(copia))
	for _, rid := range copia {
		reg, err := r.regiaoRepo.BuscarPorID(ctx, rid)
		if err != nil {
			continue
		}
		resultado = append(resultado, reg)
	}

	return resultado, nil
}
