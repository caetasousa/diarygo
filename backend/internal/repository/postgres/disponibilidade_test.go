//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/postgres"
	"github.com/caetasousa/diarygo/internal/testutil"
	"github.com/caetasousa/diarygo/internal/testutil/factories"
)

func TestDisponibilidadeRepository_Definir_Substitui(t *testing.T) {
	testutil.Truncate(t, testDB)
	ctx := context.Background()
	repo := postgres.NewDisponibilidadeRepository(testDB)

	prof := factories.NewProfissional(t, testDB)

	// Primeira chamada: define 3 slots.
	slots := []*domain.Disponibilidade{
		{ProfissionalID: prof.ID, DiaSemana: 1, HoraInicio: "08:00", HoraFim: "12:00"},
		{ProfissionalID: prof.ID, DiaSemana: 1, HoraInicio: "14:00", HoraFim: "18:00"},
		{ProfissionalID: prof.ID, DiaSemana: 3, HoraInicio: "09:00", HoraFim: "17:00"},
	}
	if err := repo.DefinirDisponibilidades(ctx, prof.ID, slots); err != nil {
		t.Fatalf("DefinirDisponibilidades: %v", err)
	}

	lista, err := repo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		t.Fatalf("Listar: %v", err)
	}
	if len(lista) != 3 {
		t.Fatalf("esperava 3 slots, got %d", len(lista))
	}
	// Ordem: dia_semana asc, hora_inicio asc
	if lista[0].DiaSemana != 1 || lista[0].HoraInicio != "08:00" {
		t.Errorf("ordem errada: %+v", lista[0])
	}
	if lista[2].DiaSemana != 3 {
		t.Errorf("esperava dia 3 no final: %+v", lista[2])
	}

	// Segunda chamada substitui: 1 slot novo.
	novo := []*domain.Disponibilidade{
		{ProfissionalID: prof.ID, DiaSemana: 5, HoraInicio: "10:00", HoraFim: "14:00"},
	}
	if err := repo.DefinirDisponibilidades(ctx, prof.ID, novo); err != nil {
		t.Fatalf("redefinir: %v", err)
	}
	lista, _ = repo.ListarPorProfissionalID(ctx, prof.ID)
	if len(lista) != 1 || lista[0].DiaSemana != 5 {
		t.Errorf("substituicao errada, got %+v", lista)
	}
}
