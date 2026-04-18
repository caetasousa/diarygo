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

func TestDocumentoRepository_CriarListar(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewDocumentoRepository(testDB)

	prof := factories.NewProfissional(t, testDB)

	docs := []*domain.Documento{
		{ID: uuid.New(), ProfissionalID: prof.ID, Tipo: domain.DocRGFrente, URL: "https://x/1", Status: domain.DocPendente},
		{ID: uuid.New(), ProfissionalID: prof.ID, Tipo: domain.DocCPF, URL: "https://x/2", Status: domain.DocPendente},
	}
	for _, d := range docs {
		if err := repo.Criar(ctx, d); err != nil {
			t.Fatalf("Criar: %v", err)
		}
	}

	lista, err := repo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		t.Fatalf("ListarPorProfissionalID: %v", err)
	}
	if len(lista) != 2 {
		t.Fatalf("esperava 2 documentos, got %d", len(lista))
	}
}

func TestDocumentoRepository_Atualizar_StatusObservacao(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewDocumentoRepository(testDB)

	prof := factories.NewProfissional(t, testDB)
	d := &domain.Documento{ID: uuid.New(), ProfissionalID: prof.ID, Tipo: domain.DocFoto, URL: "https://x", Status: domain.DocPendente}
	if err := repo.Criar(ctx, d); err != nil {
		t.Fatalf("Criar: %v", err)
	}

	d.Status = domain.DocAprovado
	d.ObservacaoAdmin = "ok"
	if err := repo.Atualizar(ctx, d); err != nil {
		t.Fatalf("Atualizar: %v", err)
	}

	lido, _ := repo.BuscarPorID(ctx, d.ID)
	if lido.Status != domain.DocAprovado || lido.ObservacaoAdmin != "ok" {
		t.Errorf("update nao persistiu: %+v", lido)
	}
}
