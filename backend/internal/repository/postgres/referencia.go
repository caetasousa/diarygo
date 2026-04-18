package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// ReferenciaRepository e a implementacao Postgres de domain.ReferenciaRepository.
type ReferenciaRepository struct {
	pool *pgxpool.Pool
}

func NewReferenciaRepository(pool *pgxpool.Pool) *ReferenciaRepository {
	return &ReferenciaRepository{pool: pool}
}

const colsReferencia = `id, profissional_id, nome_contato, telefone_contato, status, criado_em, atualizado_em`

func (r *ReferenciaRepository) Criar(ctx context.Context, ref *domain.Referencia) error {
	const q = `
        INSERT INTO referencias (id, profissional_id, nome_contato, telefone_contato, status)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING criado_em, atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		ref.ID, ref.ProfissionalID, ref.NomeContato, ref.TelefoneContato, string(ref.Status),
	).Scan(&ref.CriadoEm, &ref.AtualizadoEm)
	if err != nil {
		if pgErrorCode(err) == codeForeignKeyViolation {
			return domain.ErrProfissionalNaoEncontrada
		}
		return err
	}
	return nil
}

func (r *ReferenciaRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Referencia, error) {
	const q = `SELECT ` + colsReferencia + ` FROM referencias WHERE id = $1`
	var (
		ref    domain.Referencia
		status string
	)
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&ref.ID, &ref.ProfissionalID, &ref.NomeContato, &ref.TelefoneContato, &status,
		&ref.CriadoEm, &ref.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrReferenciaNaoEncontrada
		}
		return nil, err
	}
	ref.Status = domain.StatusReferencia(status)
	return &ref, nil
}

func (r *ReferenciaRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Referencia, error) {
	const q = `
        SELECT ` + colsReferencia + `
          FROM referencias
         WHERE profissional_id = $1
         ORDER BY criado_em ASC
    `
	rows, err := r.pool.Query(ctx, q, profissionalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Referencia
	for rows.Next() {
		var (
			ref    domain.Referencia
			status string
		)
		if err := rows.Scan(
			&ref.ID, &ref.ProfissionalID, &ref.NomeContato, &ref.TelefoneContato, &status,
			&ref.CriadoEm, &ref.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		ref.Status = domain.StatusReferencia(status)
		out = append(out, &ref)
	}
	return out, rows.Err()
}
