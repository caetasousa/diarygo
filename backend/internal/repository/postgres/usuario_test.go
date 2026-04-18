//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/postgres"
	"github.com/caetasousa/diarygo/internal/testutil"
)

func TestUsuarioRepository_CriarBuscarAtualizar(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewUsuarioRepository(testDB)

	u := &domain.Usuario{
		ID:              uuid.New(),
		Email:           "  ALICE@Example.com  ",
		SenhaHash:       "$2a$12$fakehash",
		Tipo:            domain.TipoCliente,
		EmailVerificado: false,
		Ativo:           true,
	}
	if err := repo.Criar(ctx, u); err != nil {
		t.Fatalf("Criar: %v", err)
	}
	if u.Email != "alice@example.com" {
		t.Errorf("email nao normalizado: %q", u.Email)
	}
	if u.CriadoEm.IsZero() {
		t.Error("CriadoEm nao preenchido pelo Criar")
	}

	encontrado, err := repo.BuscarPorEmail(ctx, "Alice@Example.com")
	if err != nil {
		t.Fatalf("BuscarPorEmail: %v", err)
	}
	if encontrado.ID != u.ID {
		t.Errorf("ID: got %s, want %s", encontrado.ID, u.ID)
	}

	porID, err := repo.BuscarPorID(ctx, u.ID)
	if err != nil {
		t.Fatalf("BuscarPorID: %v", err)
	}
	if porID.Email != "alice@example.com" {
		t.Errorf("email lido: %s", porID.Email)
	}

	porID.TokenRecuperacao = "abc123"
	porID.TokenRecuperacaoExpira = time.Now().Add(30 * time.Minute)
	if err := repo.Atualizar(ctx, porID); err != nil {
		t.Fatalf("Atualizar: %v", err)
	}

	viaToken, err := repo.BuscarPorTokenRecuperacao(ctx, "abc123")
	if err != nil {
		t.Fatalf("BuscarPorTokenRecuperacao: %v", err)
	}
	if viaToken.ID != u.ID {
		t.Errorf("BuscarPorTokenRecuperacao retornou id errado")
	}
}

func TestUsuarioRepository_Criar_EmailDuplicado(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewUsuarioRepository(testDB)

	base := &domain.Usuario{
		ID: uuid.New(), Email: "dup@example.com", SenhaHash: "h", Tipo: domain.TipoCliente, Ativo: true,
	}
	if err := repo.Criar(ctx, base); err != nil {
		t.Fatalf("Criar base: %v", err)
	}

	dup := &domain.Usuario{
		ID: uuid.New(), Email: "DUP@example.com", SenhaHash: "h", Tipo: domain.TipoProfissional, Ativo: true,
	}
	err := repo.Criar(ctx, dup)
	if !errors.Is(err, domain.ErrEmailJaExiste) {
		t.Errorf("esperado ErrEmailJaExiste, got %v", err)
	}
}

func TestUsuarioRepository_BuscarPorID_NaoEncontrado(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewUsuarioRepository(testDB)

	_, err := repo.BuscarPorID(ctx, uuid.New())
	if !errors.Is(err, domain.ErrUsuarioNaoEncontrado) {
		t.Errorf("esperado ErrUsuarioNaoEncontrado, got %v", err)
	}
}
