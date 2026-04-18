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

// seed carrega as 12 regioes administrativas oficiais de Goiania/GO.
// A prefeitura divide a cidade em regioes para fins de planejamento; cada
// uma agrega dezenas de setores/bairros. As faixas de CEP sao aproximadas
// com base na cobertura dos Correios e servem para sugerir a regiao a
// partir do CEP do cliente — pode haver excecoes em bairros de divisa.
func (r *RegiaoRepository) seed() {
	const uf = "GO"
	const cidade = "Goiânia"
	regioes := []struct {
		nome              string
		cepInicio, cepFim string
	}{
		{"Região Central", "74000000", "74049999"},
		{"Região Norte", "74300000", "74309999"},
		{"Região Sul", "74080000", "74299999"},
		{"Região Sudoeste", "74310000", "74399999"},
		{"Região Oeste", "74110000", "74149999"},
		{"Região Noroeste", "74400000", "74499999"},
		{"Região Campinas-Centro", "74500000", "74569999"},
		{"Região Macambira", "74570000", "74599999"},
		{"Região Leste", "74600000", "74669999"},
		{"Região Vale do Meia Ponte", "74670000", "74799999"},
		{"Região Sudeste", "74800000", "74899999"},
		{"Região Mendanha", "74900000", "74999999"},
	}
	for _, reg := range regioes {
		id := uuid.New()
		r.regioes[id] = &domain.Regiao{
			ID: id, Nome: reg.nome, Cidade: cidade, Estado: uf,
			CEPInicio: reg.cepInicio, CEPFim: reg.cepFim, Ativa: true, CriadoEm: time.Now(),
		}
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
