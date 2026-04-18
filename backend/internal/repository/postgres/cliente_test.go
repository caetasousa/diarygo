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

func TestClienteRepository_Criar_BuscarPorUsuario(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewClienteRepository(testDB)

	u := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoCliente))

	c := &domain.Cliente{
		ID:        uuid.New(),
		UsuarioID: u.ID,
		Nome:      "Alice",
		CPF:       factories.GenerateCPF(),
		Telefone:  "11999999999",
		Score:     100,
	}
	if err := repo.Criar(ctx, c); err != nil {
		t.Fatalf("Criar: %v", err)
	}

	encontrado, err := repo.BuscarPorUsuarioID(ctx, u.ID)
	if err != nil {
		t.Fatalf("BuscarPorUsuarioID: %v", err)
	}
	if encontrado.ID != c.ID {
		t.Errorf("id difere: got %s, want %s", encontrado.ID, c.ID)
	}

	porCPF, err := repo.BuscarPorCPF(ctx, c.CPF)
	if err != nil {
		t.Fatalf("BuscarPorCPF: %v", err)
	}
	if porCPF.Nome != "Alice" {
		t.Errorf("Nome lido: %s", porCPF.Nome)
	}
}

func TestClienteRepository_Criar_CPFDuplicado(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewClienteRepository(testDB)

	cpf := factories.GenerateCPF()
	u1 := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoCliente))
	u2 := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoCliente))

	if err := repo.Criar(ctx, &domain.Cliente{ID: uuid.New(), UsuarioID: u1.ID, Nome: "A", CPF: cpf, Score: 100}); err != nil {
		t.Fatalf("Criar c1: %v", err)
	}
	err := repo.Criar(ctx, &domain.Cliente{ID: uuid.New(), UsuarioID: u2.ID, Nome: "B", CPF: cpf, Score: 100})
	if !errors.Is(err, domain.ErrCPFJaCadastrado) {
		t.Errorf("esperado ErrCPFJaCadastrado, got %v", err)
	}
}

func TestClienteRepository_Criar_UsuarioDuplicado(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewClienteRepository(testDB)

	u := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoCliente))

	if err := repo.Criar(ctx, &domain.Cliente{ID: uuid.New(), UsuarioID: u.ID, Nome: "A", CPF: factories.GenerateCPF(), Score: 100}); err != nil {
		t.Fatalf("Criar c1: %v", err)
	}
	err := repo.Criar(ctx, &domain.Cliente{ID: uuid.New(), UsuarioID: u.ID, Nome: "B", CPF: factories.GenerateCPF(), Score: 100})
	if !errors.Is(err, domain.ErrClienteJaExiste) {
		t.Errorf("esperado ErrClienteJaExiste, got %v", err)
	}
}

func TestClienteRepository_Atualizar(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewClienteRepository(testDB)

	u := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoCliente))
	c := &domain.Cliente{ID: uuid.New(), UsuarioID: u.ID, Nome: "Old", CPF: factories.GenerateCPF(), Telefone: "11900000000", Score: 100}
	if err := repo.Criar(ctx, c); err != nil {
		t.Fatalf("Criar: %v", err)
	}

	c.Nome = "New"
	c.Score = 80
	if err := repo.Atualizar(ctx, c); err != nil {
		t.Fatalf("Atualizar: %v", err)
	}

	lido, err := repo.BuscarPorID(ctx, c.ID)
	if err != nil {
		t.Fatalf("BuscarPorID: %v", err)
	}
	if lido.Nome != "New" || lido.Score != 80 {
		t.Errorf("update nao persistiu: %+v", lido)
	}
}
