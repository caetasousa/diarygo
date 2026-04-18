//go:build integration

package factories

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/caetasousa/diarygo/internal/domain"
)

var emailCounter uint64

// UsuarioOpt customizes the default Usuario built by NewUsuario.
type UsuarioOpt func(*domain.Usuario)

func ComEmail(email string) UsuarioOpt        { return func(u *domain.Usuario) { u.Email = email } }
func ComTipo(t domain.TipoUsuario) UsuarioOpt { return func(u *domain.Usuario) { u.Tipo = t } }
func ComSenha(plain string) UsuarioOpt {
	return func(u *domain.Usuario) {
		hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
		if err != nil {
			panic(fmt.Sprintf("factories: hash senha: %v", err))
		}
		u.SenhaHash = string(hash)
	}
}

// NewUsuario inserts a valid Usuario row and returns the domain struct.
// Defaults: unique email per call, tipo=CLIENTE, senha="Senha@1234" (hashed),
// email_verificado=true, ativo=true.
func NewUsuario(t *testing.T, db *pgxpool.Pool, opts ...UsuarioOpt) *domain.Usuario {
	t.Helper()
	n := atomic.AddUint64(&emailCounter, 1)
	hash, _ := bcrypt.GenerateFromPassword([]byte("Senha@1234"), 12)

	u := &domain.Usuario{
		ID:              uuid.New(),
		Email:           fmt.Sprintf("user%d@diarygo.test", n),
		SenhaHash:       string(hash),
		Tipo:            domain.TipoCliente,
		EmailVerificado: true,
		Ativo:           true,
	}
	for _, opt := range opts {
		opt(u)
	}

	_, err := db.Exec(context.Background(), `
		INSERT INTO usuarios (id, email, senha_hash, tipo, email_verificado, ativo)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, u.ID, u.Email, u.SenhaHash, string(u.Tipo), u.EmailVerificado, u.Ativo)
	if err != nil {
		t.Fatalf("factories.NewUsuario: insert: %v", err)
	}
	return u
}
