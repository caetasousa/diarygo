//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/caetasousa/diarygo/internal/repository/postgres"
	"github.com/caetasousa/diarygo/internal/testutil"
	"github.com/caetasousa/diarygo/internal/testutil/factories"
)

func TestRegiaoRepository_CriarListarAtivas(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewRegiaoRepository(testDB)

	ativa := factories.NewRegiao(t, testDB, factories.RegiaoAtiva(true))
	factories.NewRegiao(t, testDB, factories.RegiaoAtiva(false))

	lista, err := repo.ListarAtivas(ctx)
	if err != nil {
		t.Fatalf("ListarAtivas: %v", err)
	}
	if len(lista) != 1 || lista[0].ID != ativa.ID {
		t.Fatalf("esperava 1 regiao ativa, got %d: %+v", len(lista), lista)
	}
}

func TestProfissionalRegiaoRepository_DefinirRegioes(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewProfissionalRegiaoRepository(testDB)

	prof := factories.NewProfissional(t, testDB)
	r1 := factories.NewRegiao(t, testDB)
	r2 := factories.NewRegiao(t, testDB)
	r3 := factories.NewRegiao(t, testDB)

	// Define 2 regioes
	if err := repo.DefinirRegioes(ctx, prof.ID, []uuid.UUID{r1.ID, r2.ID}); err != nil {
		t.Fatalf("DefinirRegioes: %v", err)
	}
	lista, _ := repo.ListarPorProfissionalID(ctx, prof.ID)
	if len(lista) != 2 {
		t.Fatalf("esperava 2 regioes, got %d", len(lista))
	}

	// Substitui por 1 regiao — as antigas devem desaparecer.
	if err := repo.DefinirRegioes(ctx, prof.ID, []uuid.UUID{r3.ID}); err != nil {
		t.Fatalf("DefinirRegioes 2: %v", err)
	}
	lista, _ = repo.ListarPorProfissionalID(ctx, prof.ID)
	if len(lista) != 1 || lista[0].ID != r3.ID {
		t.Errorf("substituicao nao funcionou, got %+v", lista)
	}

	// Lista vazia remove todas.
	if err := repo.DefinirRegioes(ctx, prof.ID, nil); err != nil {
		t.Fatalf("DefinirRegioes vazia: %v", err)
	}
	lista, _ = repo.ListarPorProfissionalID(ctx, prof.ID)
	if len(lista) != 0 {
		t.Errorf("esperava 0 regioes apos limpar, got %d", len(lista))
	}
}
