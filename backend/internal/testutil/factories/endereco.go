//go:build integration

package factories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

type EnderecoOpt func(*domain.Endereco)

func EnderecoPrincipal(p bool) EnderecoOpt { return func(e *domain.Endereco) { e.Principal = p } }
func EnderecoCidade(c string) EnderecoOpt  { return func(e *domain.Endereco) { e.Cidade = c } }
func EnderecoCEP(cep string) EnderecoOpt   { return func(e *domain.Endereco) { e.CEP = cep } }

// NewEndereco inserts an Endereco bound to an existing clienteID. ClienteID is required.
func NewEndereco(t *testing.T, db *pgxpool.Pool, clienteID uuid.UUID, opts ...EnderecoOpt) *domain.Endereco {
	t.Helper()
	e := &domain.Endereco{
		ID:           uuid.New(),
		ClienteID:    clienteID,
		Logradouro:   "Rua Teste",
		Numero:       "100",
		Bairro:       "Centro",
		Cidade:       "São Paulo",
		Estado:       "SP",
		CEP:          "01310100",
		Principal:    true,
		NumQuartos:   2,
		NumBanheiros: 1,
		NumSalas:     1,
		NumCozinhas:  1,
	}
	for _, opt := range opts {
		opt(e)
	}

	_, err := db.Exec(context.Background(), `
		INSERT INTO enderecos (id, cliente_id, logradouro, numero, bairro, cidade, estado, cep,
		                       principal, num_quartos, num_banheiros, num_salas, num_cozinhas)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, e.ID, e.ClienteID, e.Logradouro, e.Numero, e.Bairro, e.Cidade, e.Estado, e.CEP,
		e.Principal, e.NumQuartos, e.NumBanheiros, e.NumSalas, e.NumCozinhas)
	if err != nil {
		t.Fatalf("factories.NewEndereco: insert: %v", err)
	}
	return e
}
