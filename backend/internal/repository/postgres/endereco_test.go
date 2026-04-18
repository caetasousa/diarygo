//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/postgres"
	"github.com/caetasousa/diarygo/internal/testutil"
	"github.com/caetasousa/diarygo/internal/testutil/factories"
)

func TestEnderecoRepository_CRUD(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewEnderecoRepository(testDB)

	cli := factories.NewCliente(t, testDB)

	e := &domain.Endereco{
		ID:           uuid.New(),
		ClienteID:    cli.ID,
		Logradouro:   "Av. Paulista",
		Numero:       "1000",
		Bairro:       "Bela Vista",
		Cidade:       "São Paulo",
		Estado:       "SP",
		CEP:          "01310100",
		Latitude:     -23.5613,
		Longitude:    -46.6558,
		NumQuartos:   2,
		NumBanheiros: 1,
		NumSalas:     1,
		NumCozinhas:  1,
		AreaM2:       75.5,
		Principal:    true,
	}
	if err := repo.Criar(ctx, e); err != nil {
		t.Fatalf("Criar: %v", err)
	}

	lista, err := repo.ListarPorClienteID(ctx, cli.ID)
	if err != nil {
		t.Fatalf("ListarPorClienteID: %v", err)
	}
	if len(lista) != 1 {
		t.Fatalf("esperava 1 endereco, got %d", len(lista))
	}
	if lista[0].Latitude == 0 {
		t.Error("latitude nao persistida")
	}

	e.Numero = "1500"
	e.Principal = false
	if err := repo.Atualizar(ctx, e); err != nil {
		t.Fatalf("Atualizar: %v", err)
	}
	lido, _ := repo.BuscarPorID(ctx, e.ID)
	if lido.Numero != "1500" || lido.Principal {
		t.Errorf("update nao persistiu: %+v", lido)
	}

	if err := repo.Remover(ctx, e.ID); err != nil {
		t.Fatalf("Remover: %v", err)
	}
	_, err = repo.BuscarPorID(ctx, e.ID)
	if !errors.Is(err, domain.ErrEnderecoNaoEncontrado) {
		t.Errorf("esperado ErrEnderecoNaoEncontrado, got %v", err)
	}
}

func TestEnderecoRepository_OrdemPrincipalPrimeiro(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewEnderecoRepository(testDB)

	cli := factories.NewCliente(t, testDB)

	for i, principal := range []bool{false, true, false} {
		e := &domain.Endereco{
			ID: uuid.New(), ClienteID: cli.ID,
			Logradouro: "R X", Numero: "1", Bairro: "B", Cidade: "C", Estado: "SP",
			CEP: "01310100", Principal: principal,
		}
		if err := repo.Criar(ctx, e); err != nil {
			t.Fatalf("Criar %d: %v", i, err)
		}
	}
	lista, err := repo.ListarPorClienteID(ctx, cli.ID)
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if len(lista) != 3 || !lista[0].Principal {
		t.Errorf("esperava principal primeiro, got %+v", lista)
	}
}
