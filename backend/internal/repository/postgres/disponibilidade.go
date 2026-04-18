package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// parseHoraHHMM converte "HH:MM" em time.Time (na data zero) — pgx serializa
// como TIME no Postgres.
func parseHoraHHMM(s string) (time.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("postgres: hora invalida %q: %w", s, err)
	}
	return t, nil
}

// DisponibilidadeRepository e a implementacao Postgres de
// domain.DisponibilidadeRepository.
type DisponibilidadeRepository struct {
	pool *pgxpool.Pool
}

func NewDisponibilidadeRepository(pool *pgxpool.Pool) *DisponibilidadeRepository {
	return &DisponibilidadeRepository{pool: pool}
}

// ListarPorProfissionalID retorna as disponibilidades na ordem natural do calendario
// (dia da semana asc, hora de inicio asc).
func (r *DisponibilidadeRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Disponibilidade, error) {
	const q = `
        SELECT id, profissional_id, dia_semana,
               to_char(hora_inicio, 'HH24:MI') AS hora_inicio,
               to_char(hora_fim,    'HH24:MI') AS hora_fim,
               criado_em
          FROM disponibilidades
         WHERE profissional_id = $1
         ORDER BY dia_semana ASC, hora_inicio ASC
    `
	rows, err := r.pool.Query(ctx, q, profissionalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Disponibilidade
	for rows.Next() {
		var d domain.Disponibilidade
		if err := rows.Scan(
			&d.ID, &d.ProfissionalID, &d.DiaSemana, &d.HoraInicio, &d.HoraFim, &d.CriadoEm,
		); err != nil {
			return nil, err
		}
		// Nao ha coluna atualizado_em na tabela; manter zero-value.
		out = append(out, &d)
	}
	return out, rows.Err()
}

// DefinirDisponibilidades substitui atomicamente o conjunto de slots da
// profissional: apaga os existentes e insere os novos em uma transacao.
func (r *DisponibilidadeRepository) DefinirDisponibilidades(ctx context.Context, profissionalID uuid.UUID, slots []*domain.Disponibilidade) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `DELETE FROM disponibilidades WHERE profissional_id = $1`, profissionalID); err != nil {
		return err
	}

	if len(slots) > 0 {
		rows := make([][]any, 0, len(slots))
		for _, s := range slots {
			id := s.ID
			if id == uuid.Nil {
				id = uuid.New()
			}
			hi, err := parseHoraHHMM(s.HoraInicio)
			if err != nil {
				return err
			}
			hf, err := parseHoraHHMM(s.HoraFim)
			if err != nil {
				return err
			}
			rows = append(rows, []any{id, profissionalID, s.DiaSemana, hi, hf})
		}
		_, err := tx.CopyFrom(ctx,
			pgx.Identifier{"disponibilidades"},
			[]string{"id", "profissional_id", "dia_semana", "hora_inicio", "hora_fim"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			if pgErrorCode(err) == codeForeignKeyViolation {
				return domain.ErrProfissionalNaoEncontrada
			}
			return err
		}
	}

	return tx.Commit(ctx)
}
