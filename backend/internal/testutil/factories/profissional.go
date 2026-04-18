//go:build integration

package factories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

type ProfissionalOpt func(*domain.Profissional)

func ProfissionalComNome(n string) ProfissionalOpt {
	return func(p *domain.Profissional) { p.Nome = n }
}
func ProfissionalComCPF(cpf string) ProfissionalOpt {
	return func(p *domain.Profissional) { p.CPF = cpf }
}
func ProfissionalComStatus(s domain.StatusProfissional) ProfissionalOpt {
	return func(p *domain.Profissional) { p.Status = s }
}
func ProfissionalParaUsuario(id uuid.UUID) ProfissionalOpt {
	return func(p *domain.Profissional) { p.UsuarioID = id }
}

// NewProfissional inserts a valid Profissional row. If no UsuarioID is supplied,
// a new Usuario (tipo=PROFISSIONAL) is created automatically.
func NewProfissional(t *testing.T, db *pgxpool.Pool, opts ...ProfissionalOpt) *domain.Profissional {
	t.Helper()
	p := &domain.Profissional{
		ID:       uuid.New(),
		Nome:     "Profissional Teste",
		CPF:      GenerateCPF(),
		Telefone: "11988888888",
		Status:   domain.StatusPendente,
		MEI:      false,
	}
	for _, opt := range opts {
		opt(p)
	}
	if p.UsuarioID == uuid.Nil {
		u := NewUsuario(t, db, ComTipo(domain.TipoProfissional))
		p.UsuarioID = u.ID
	}

	_, err := db.Exec(context.Background(), `
		INSERT INTO profissionais (id, usuario_id, nome, cpf, telefone, status, mei)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, p.ID, p.UsuarioID, p.Nome, p.CPF, p.Telefone, string(p.Status), p.MEI)
	if err != nil {
		t.Fatalf("factories.NewProfissional: insert: %v", err)
	}
	return p
}
