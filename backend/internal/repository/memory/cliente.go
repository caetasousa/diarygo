package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// ClienteRepository implementa domain.ClienteRepository em memoria.
type ClienteRepository struct {
	mu         sync.RWMutex
	clientes   map[uuid.UUID]*domain.Cliente
	porUsuario map[uuid.UUID]uuid.UUID // usuarioID -> clienteID
	porCPF     map[string]uuid.UUID    // cpf -> clienteID
}

// NewClienteRepository cria um novo repositorio in-memory vazio.
func NewClienteRepository() *ClienteRepository {
	return &ClienteRepository{
		clientes:   make(map[uuid.UUID]*domain.Cliente),
		porUsuario: make(map[uuid.UUID]uuid.UUID),
		porCPF:     make(map[string]uuid.UUID),
	}
}

// Criar persiste um novo cliente. Retorna erro se o usuario ou CPF ja tiver cliente.
func (r *ClienteRepository) Criar(ctx context.Context, c *domain.Cliente) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, existe := r.porUsuario[c.UsuarioID]; existe {
		return domain.ErrClienteJaExiste
	}
	if c.CPF != "" {
		if _, existe := r.porCPF[c.CPF]; existe {
			return domain.ErrCPFJaCadastrado
		}
	}

	copia := *c
	copia.CriadoEm = time.Now()
	copia.AtualizadoEm = time.Now()

	r.clientes[copia.ID] = &copia
	r.porUsuario[copia.UsuarioID] = copia.ID
	if copia.CPF != "" {
		r.porCPF[copia.CPF] = copia.ID
	}

	return nil
}

// BuscarPorUsuarioID retorna uma copia do cliente pelo ID do usuario associado.
func (r *ClienteRepository) BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*domain.Cliente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, existe := r.porUsuario[usuarioID]
	if !existe {
		return nil, domain.ErrClienteNaoEncontrado
	}

	c := r.clientes[id]
	copia := *c
	return &copia, nil
}

// BuscarPorID retorna uma copia do cliente pelo UUID.
func (r *ClienteRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Cliente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, existe := r.clientes[id]
	if !existe {
		return nil, domain.ErrClienteNaoEncontrado
	}

	copia := *c
	return &copia, nil
}

// BuscarPorCPF retorna uma copia do cliente pelo CPF.
func (r *ClienteRepository) BuscarPorCPF(ctx context.Context, cpf string) (*domain.Cliente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, existe := r.porCPF[cpf]
	if !existe {
		return nil, domain.ErrClienteNaoEncontrado
	}

	c := r.clientes[id]
	copia := *c
	return &copia, nil
}

// AjustarScore soma delta ao score do cliente com clamp em [0, 100].
func (r *ClienteRepository) AjustarScore(ctx context.Context, id uuid.UUID, delta int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, existe := r.clientes[id]
	if !existe {
		return domain.ErrClienteNaoEncontrado
	}

	novo := c.Score + delta
	if novo < 0 {
		novo = 0
	}
	if novo > 100 {
		novo = 100
	}
	c.Score = novo
	c.AtualizadoEm = time.Now()

	return nil
}

// Atualizar substitui o cliente existente. Atualiza indice de CPF se mudou.
func (r *ClienteRepository) Atualizar(ctx context.Context, c *domain.Cliente) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	antigo, existe := r.clientes[c.ID]
	if !existe {
		return domain.ErrClienteNaoEncontrado
	}

	// Se o CPF mudou, verificar conflito e atualizar indice
	if c.CPF != antigo.CPF {
		if c.CPF != "" {
			if id, ok := r.porCPF[c.CPF]; ok && id != c.ID {
				return domain.ErrCPFJaCadastrado
			}
		}
		if antigo.CPF != "" {
			delete(r.porCPF, antigo.CPF)
		}
		if c.CPF != "" {
			r.porCPF[c.CPF] = c.ID
		}
	}

	copia := *c
	copia.CriadoEm = antigo.CriadoEm
	copia.AtualizadoEm = time.Now()
	r.clientes[copia.ID] = &copia

	return nil
}
