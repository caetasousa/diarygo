//go:build integration

package factories

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

var regiaoCounter uint64

type RegiaoOpt func(*domain.Regiao)

func RegiaoAtiva(a bool) RegiaoOpt  { return func(r *domain.Regiao) { r.Ativa = a } }
func RegiaoNome(n string) RegiaoOpt { return func(r *domain.Regiao) { r.Nome = n } }
func RegiaoCidadeEstado(c, uf string) RegiaoOpt {
	return func(r *domain.Regiao) { r.Cidade = c; r.Estado = uf }
}

// NewRegiao inserts a Regiao row with unique name per call.
func NewRegiao(t *testing.T, db *pgxpool.Pool, opts ...RegiaoOpt) *domain.Regiao {
	t.Helper()
	n := atomic.AddUint64(&regiaoCounter, 1)
	r := &domain.Regiao{
		ID:        uuid.New(),
		Nome:      fmt.Sprintf("Regiao Teste %d", n),
		Cidade:    "São Paulo",
		Estado:    "SP",
		CEPInicio: "01000000",
		CEPFim:    "05999999",
		Ativa:     true,
	}
	for _, opt := range opts {
		opt(r)
	}

	_, err := db.Exec(context.Background(), `
		INSERT INTO regioes (id, nome, cidade, estado, cep_inicio, cep_fim, ativa)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, r.ID, r.Nome, r.Cidade, r.Estado, r.CEPInicio, r.CEPFim, r.Ativa)
	if err != nil {
		t.Fatalf("factories.NewRegiao: insert: %v", err)
	}
	return r
}
