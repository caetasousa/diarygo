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

func TestProfissionalRepository_CRUD(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewProfissionalRepository(testDB)

	u := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoProfissional))
	p := &domain.Profissional{
		ID:        uuid.New(),
		UsuarioID: u.ID,
		Nome:      "Profissional Teste",
		CPF:       factories.GenerateCPF(),
		Telefone:  "11988888888",
		Status:    domain.StatusPendente,
		MEI:       false,
	}
	if err := repo.Criar(ctx, p); err != nil {
		t.Fatalf("Criar: %v", err)
	}
	if p.CriadoEm.IsZero() {
		t.Error("CriadoEm nao preenchido")
	}

	lido, err := repo.BuscarPorUsuarioID(ctx, u.ID)
	if err != nil {
		t.Fatalf("BuscarPorUsuarioID: %v", err)
	}
	if lido.Status != domain.StatusPendente {
		t.Errorf("status: %s", lido.Status)
	}

	lido.Status = domain.StatusAprovada
	lido.NotaMedia = 4.8
	lido.TotalServicos = 3
	if err := repo.Atualizar(ctx, lido); err != nil {
		t.Fatalf("Atualizar: %v", err)
	}
	again, _ := repo.BuscarPorID(ctx, p.ID)
	if again.Status != domain.StatusAprovada || again.NotaMedia != 4.8 || again.TotalServicos != 3 {
		t.Errorf("update nao persistiu: %+v", again)
	}
}

func TestProfissionalRepository_CPFDuplicado(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewProfissionalRepository(testDB)

	cpf := factories.GenerateCPF()
	u1 := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoProfissional))
	u2 := factories.NewUsuario(t, testDB, factories.ComTipo(domain.TipoProfissional))

	if err := repo.Criar(ctx, &domain.Profissional{ID: uuid.New(), UsuarioID: u1.ID, Nome: "A", CPF: cpf, Status: domain.StatusPendente}); err != nil {
		t.Fatalf("Criar p1: %v", err)
	}
	err := repo.Criar(ctx, &domain.Profissional{ID: uuid.New(), UsuarioID: u2.ID, Nome: "B", CPF: cpf, Status: domain.StatusPendente})
	if !errors.Is(err, domain.ErrCPFJaCadastrado) {
		t.Errorf("esperado ErrCPFJaCadastrado, got %v", err)
	}
}
