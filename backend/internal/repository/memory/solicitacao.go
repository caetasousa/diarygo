package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// SolicitacaoRepository implementa domain.SolicitacaoRepository em memória.
type SolicitacaoRepository struct {
	mu                sync.RWMutex
	porID             map[uuid.UUID]*domain.Solicitacao
	opsPorSolicitacao map[uuid.UUID][]domain.SolicitacaoOpcional
	porCliente        map[uuid.UUID][]uuid.UUID // clienteID -> []solicitacaoID em ordem de criação
}

// NewSolicitacaoRepository cria um repositório in-memory vazio.
func NewSolicitacaoRepository() *SolicitacaoRepository {
	return &SolicitacaoRepository{
		porID:             make(map[uuid.UUID]*domain.Solicitacao),
		opsPorSolicitacao: make(map[uuid.UUID][]domain.SolicitacaoOpcional),
		porCliente:        make(map[uuid.UUID][]uuid.UUID),
	}
}

// Criar persiste uma nova solicitação com seus opcionais.
func (r *SolicitacaoRepository) Criar(ctx context.Context, s *domain.Solicitacao, ops []domain.SolicitacaoOpcional) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	agora := time.Now().UTC()
	copia := *s
	copia.CriadaEm = agora
	copia.AtualizadaEm = agora
	// deep copy do breakdown
	copia.Breakdown = make([]domain.ItemCalculo, len(s.Breakdown))
	copy(copia.Breakdown, s.Breakdown)

	r.porID[copia.ID] = &copia
	r.opsPorSolicitacao[copia.ID] = append([]domain.SolicitacaoOpcional(nil), ops...)
	r.porCliente[copia.ClienteID] = append(r.porCliente[copia.ClienteID], copia.ID)

	return nil
}

// BuscarPorID retorna cópia da solicitação e seus opcionais.
func (r *SolicitacaoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Solicitacao, []domain.SolicitacaoOpcional, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, existe := r.porID[id]
	if !existe {
		return nil, nil, domain.ErrSolicitacaoNaoEncontrada
	}

	copia := *s
	copia.Breakdown = make([]domain.ItemCalculo, len(s.Breakdown))
	copy(copia.Breakdown, s.Breakdown)

	ops := append([]domain.SolicitacaoOpcional(nil), r.opsPorSolicitacao[id]...)
	return &copia, ops, nil
}

// ListarPorCliente retorna solicitações do cliente aplicando filtros, em ordem decrescente de criação.
func (r *SolicitacaoRepository) ListarPorCliente(ctx context.Context, clienteID uuid.UUID, filtro domain.SolicitacaoFiltro) ([]*domain.Solicitacao, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.porCliente[clienteID]
	resultado := make([]*domain.Solicitacao, 0, len(ids))

	for _, id := range ids {
		s, ok := r.porID[id]
		if !ok {
			continue
		}
		if filtro.Status != nil && s.Status != *filtro.Status {
			continue
		}
		if filtro.Desde != nil && s.CriadaEm.Before(*filtro.Desde) {
			continue
		}
		if filtro.Ate != nil && s.CriadaEm.After(*filtro.Ate) {
			continue
		}
		copia := *s
		copia.Breakdown = make([]domain.ItemCalculo, len(s.Breakdown))
		copy(copia.Breakdown, s.Breakdown)
		resultado = append(resultado, &copia)
	}

	// Ordem decrescente por CriadaEm
	sort.Slice(resultado, func(i, j int) bool {
		return resultado[i].CriadaEm.After(resultado[j].CriadaEm)
	})

	return resultado, nil
}

// AtualizarStatus altera o status de uma solicitação existente.
func (r *SolicitacaoRepository) AtualizarStatus(ctx context.Context, id uuid.UUID, status domain.StatusSolicitacao, canceladaEm *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, existe := r.porID[id]
	if !existe {
		return domain.ErrSolicitacaoNaoEncontrada
	}
	if !status.Valida() {
		return domain.ErrSolicitacaoStatusInvalido
	}

	s.Status = status
	s.AtualizadaEm = time.Now().UTC()
	if canceladaEm != nil {
		s.CanceladaEm = canceladaEm
	}

	return nil
}
