package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// UsuarioRepository implementa domain.UsuarioRepository em memoria.
// Seguro para uso concorrente via sync.RWMutex.
type UsuarioRepository struct {
	mu       sync.RWMutex
	usuarios map[uuid.UUID]*domain.Usuario
	porEmail map[string]uuid.UUID // indice secundario para busca O(1) por email
}

// NewUsuarioRepository cria um novo repositorio in-memory vazio.
func NewUsuarioRepository() *UsuarioRepository {
	return &UsuarioRepository{
		usuarios: make(map[uuid.UUID]*domain.Usuario),
		porEmail: make(map[string]uuid.UUID),
	}
}

// Criar persiste um novo usuario. Retorna ErrEmailJaExiste se o email ja estiver em uso.
func (r *UsuarioRepository) Criar(ctx context.Context, u *domain.Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	emailNorm := strings.ToLower(u.Email)
	if _, existe := r.porEmail[emailNorm]; existe {
		return domain.ErrEmailJaExiste
	}

	// Armazena copia para evitar mutacao externa
	copia := *u
	copia.Email = emailNorm
	copia.CriadoEm = time.Now()
	copia.AtualizadoEm = time.Now()

	r.usuarios[copia.ID] = &copia
	r.porEmail[emailNorm] = copia.ID

	return nil
}

// BuscarPorEmail retorna uma copia do usuario pelo email. Case-insensitive.
func (r *UsuarioRepository) BuscarPorEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	emailNorm := strings.ToLower(email)
	id, existe := r.porEmail[emailNorm]
	if !existe {
		return nil, domain.ErrUsuarioNaoEncontrado
	}

	u := r.usuarios[id]
	copia := *u
	return &copia, nil
}

// BuscarPorID retorna uma copia do usuario pelo UUID.
func (r *UsuarioRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, existe := r.usuarios[id]
	if !existe {
		return nil, domain.ErrUsuarioNaoEncontrado
	}

	copia := *u
	return &copia, nil
}

// BuscarPorTokenRecuperacao retorna uma copia do usuario pelo token de recuperacao de senha.
func (r *UsuarioRepository) BuscarPorTokenRecuperacao(ctx context.Context, token string) (*domain.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.usuarios {
		if u.TokenRecuperacao != "" && u.TokenRecuperacao == token {
			copia := *u
			return &copia, nil
		}
	}

	return nil, domain.ErrUsuarioNaoEncontrado
}

// Atualizar substitui o usuario existente. Retorna ErrUsuarioNaoEncontrado se nao existir.
func (r *UsuarioRepository) Atualizar(ctx context.Context, u *domain.Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, existe := r.usuarios[u.ID]; !existe {
		return domain.ErrUsuarioNaoEncontrado
	}

	copia := *u
	copia.AtualizadoEm = time.Now()
	r.usuarios[copia.ID] = &copia

	return nil
}
