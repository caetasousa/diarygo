package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// EnderecoRepository implementa domain.EnderecoRepository em memoria.
type EnderecoRepository struct {
	mu         sync.RWMutex
	enderecos  map[uuid.UUID]*domain.Endereco
	porCliente map[uuid.UUID][]uuid.UUID // clienteID -> []enderecoID
}

// NewEnderecoRepository cria um novo repositorio in-memory vazio.
func NewEnderecoRepository() *EnderecoRepository {
	return &EnderecoRepository{
		enderecos:  make(map[uuid.UUID]*domain.Endereco),
		porCliente: make(map[uuid.UUID][]uuid.UUID),
	}
}

// Criar persiste um novo endereco.
func (r *EnderecoRepository) Criar(ctx context.Context, e *domain.Endereco) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copia := *e
	copia.CriadoEm = time.Now()
	copia.AtualizadoEm = time.Now()

	r.enderecos[copia.ID] = &copia
	r.porCliente[copia.ClienteID] = append(r.porCliente[copia.ClienteID], copia.ID)

	return nil
}

// BuscarPorID retorna uma copia do endereco pelo UUID.
func (r *EnderecoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Endereco, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, existe := r.enderecos[id]
	if !existe {
		return nil, domain.ErrEnderecoNaoEncontrado
	}

	copia := *e
	return &copia, nil
}

// ListarPorClienteID retorna todos os enderecos de um cliente.
func (r *EnderecoRepository) ListarPorClienteID(ctx context.Context, clienteID uuid.UUID) ([]*domain.Endereco, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.porCliente[clienteID]
	resultado := make([]*domain.Endereco, 0, len(ids))
	for _, id := range ids {
		if e, ok := r.enderecos[id]; ok {
			copia := *e
			resultado = append(resultado, &copia)
		}
	}

	return resultado, nil
}

// Atualizar substitui o endereco existente.
func (r *EnderecoRepository) Atualizar(ctx context.Context, e *domain.Endereco) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	antigo, existe := r.enderecos[e.ID]
	if !existe {
		return domain.ErrEnderecoNaoEncontrado
	}

	copia := *e
	copia.CriadoEm = antigo.CriadoEm
	copia.AtualizadoEm = time.Now()
	r.enderecos[copia.ID] = &copia

	return nil
}

// Remover remove o endereco pelo ID. Remove tambem do indice por cliente.
func (r *EnderecoRepository) Remover(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	e, existe := r.enderecos[id]
	if !existe {
		return domain.ErrEnderecoNaoEncontrado
	}

	delete(r.enderecos, id)

	// Remover do indice por cliente
	ids := r.porCliente[e.ClienteID]
	novos := ids[:0]
	for _, eid := range ids {
		if eid != id {
			novos = append(novos, eid)
		}
	}
	r.porCliente[e.ClienteID] = novos

	return nil
}
