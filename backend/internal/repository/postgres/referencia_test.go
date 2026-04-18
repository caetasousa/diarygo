//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/postgres"
	"github.com/caetasousa/diarygo/internal/testutil"
	"github.com/caetasousa/diarygo/internal/testutil/factories"
)

func TestReferenciaRepository_CriarListar(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewReferenciaRepository(testDB)

	prof := factories.NewProfissional(t, testDB)
	for i := 0; i < 2; i++ {
		r := &domain.Referencia{
			ID: uuid.New(), ProfissionalID: prof.ID,
			NomeContato: "Contato", TelefoneContato: "11988887777", Status: domain.RefPendente,
		}
		if err := repo.Criar(ctx, r); err != nil {
			t.Fatalf("Criar: %v", err)
		}
	}

	lista, err := repo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		t.Fatalf("ListarPorProfissionalID: %v", err)
	}
	if len(lista) != 2 {
		t.Fatalf("esperava 2 referencias, got %d", len(lista))
	}
}
