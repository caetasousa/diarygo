//go:build integration

package factories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

type ClienteOpt func(*domain.Cliente)

func ClienteComNome(nome string) ClienteOpt    { return func(c *domain.Cliente) { c.Nome = nome } }
func ClienteComCPF(cpf string) ClienteOpt      { return func(c *domain.Cliente) { c.CPF = cpf } }
func ClienteComScore(s int) ClienteOpt         { return func(c *domain.Cliente) { c.Score = s } }
func ClienteComTelefone(tel string) ClienteOpt { return func(c *domain.Cliente) { c.Telefone = tel } }
func ClienteParaUsuario(id uuid.UUID) ClienteOpt {
	return func(c *domain.Cliente) { c.UsuarioID = id }
}

// NewCliente inserts a valid Cliente row. If no UsuarioID is supplied via
// ClienteParaUsuario, a new Usuario (tipo=CLIENTE) is created automatically.
func NewCliente(t *testing.T, db *pgxpool.Pool, opts ...ClienteOpt) *domain.Cliente {
	t.Helper()
	c := &domain.Cliente{
		ID:       uuid.New(),
		Nome:     "Cliente Teste",
		CPF:      GenerateCPF(),
		Telefone: "11999999999",
		Score:    100,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.UsuarioID == uuid.Nil {
		u := NewUsuario(t, db, ComTipo(domain.TipoCliente))
		c.UsuarioID = u.ID
	}

	_, err := db.Exec(context.Background(), `
		INSERT INTO clientes (id, usuario_id, nome, cpf, telefone, score)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, c.ID, c.UsuarioID, c.Nome, c.CPF, c.Telefone, c.Score)
	if err != nil {
		t.Fatalf("factories.NewCliente: insert: %v", err)
	}
	return c
}
