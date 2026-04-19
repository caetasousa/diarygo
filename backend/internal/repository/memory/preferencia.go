package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// preferenciaKey identifica unicamente uma preferencia pelo par (cliente, profissional).
type preferenciaKey struct {
	ClienteID      uuid.UUID
	ProfissionalID uuid.UUID
}

// PreferenciaRepository implementa domain.PreferenciaRepository em memoria.
type PreferenciaRepository struct {
	mu          sync.RWMutex
	preferencia map[preferenciaKey]*domain.ClienteProfissionalPreferencia
	porCliente  map[uuid.UUID][]preferenciaKey
}

// NewPreferenciaRepository cria um novo repositorio in-memory vazio.
func NewPreferenciaRepository() *PreferenciaRepository {
	return &PreferenciaRepository{
		preferencia: make(map[preferenciaKey]*domain.ClienteProfissionalPreferencia),
		porCliente:  make(map[uuid.UUID][]preferenciaKey),
	}
}

// Upsert insere ou substitui a preferencia do par (cliente, profissional).
// Se o par ja existia, preserva o ID original e atualiza apenas o Tipo.
func (r *PreferenciaRepository) Upsert(ctx context.Context, p *domain.ClienteProfissionalPreferencia) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := preferenciaKey{ClienteID: p.ClienteID, ProfissionalID: p.ProfissionalID}
	if existente, ok := r.preferencia[key]; ok {
		existente.Tipo = p.Tipo
		// CriadoEm preservado; essa operacao e uma atualizacao de tipo
		return nil
	}

	copia := *p
	if copia.ID == uuid.Nil {
		copia.ID = uuid.New()
	}
	if copia.CriadoEm.IsZero() {
		copia.CriadoEm = time.Now()
	}
	r.preferencia[key] = &copia
	r.porCliente[p.ClienteID] = append(r.porCliente[p.ClienteID], key)
	return nil
}

// Remover apaga a preferencia do par (cliente, profissional). Sem erro se nao existir.
func (r *PreferenciaRepository) Remover(ctx context.Context, clienteID, profissionalID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := preferenciaKey{ClienteID: clienteID, ProfissionalID: profissionalID}
	if _, ok := r.preferencia[key]; !ok {
		return domain.ErrPreferenciaNaoEncontrada
	}
	delete(r.preferencia, key)

	keys := r.porCliente[clienteID]
	for i, k := range keys {
		if k == key {
			r.porCliente[clienteID] = append(keys[:i], keys[i+1:]...)
			break
		}
	}
	return nil
}

// BuscarPorClienteEProfissional retorna a preferencia do par (ou erro nao encontrada).
func (r *PreferenciaRepository) BuscarPorClienteEProfissional(ctx context.Context, clienteID, profissionalID uuid.UUID) (*domain.ClienteProfissionalPreferencia, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.preferencia[preferenciaKey{clienteID, profissionalID}]
	if !ok {
		return nil, domain.ErrPreferenciaNaoEncontrada
	}
	copia := *p
	return &copia, nil
}

// ListarPorCliente retorna todas as preferencias do cliente filtradas por tipo.
// Se tipo vazio, retorna ambas (favoritas + bloqueadas).
func (r *PreferenciaRepository) ListarPorCliente(ctx context.Context, clienteID uuid.UUID, tipo domain.TipoPreferencia) ([]*domain.ClienteProfissionalPreferencia, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := r.porCliente[clienteID]
	resultado := make([]*domain.ClienteProfissionalPreferencia, 0, len(keys))
	for _, k := range keys {
		p := r.preferencia[k]
		if tipo != "" && p.Tipo != tipo {
			continue
		}
		copia := *p
		resultado = append(resultado, &copia)
	}
	return resultado, nil
}
