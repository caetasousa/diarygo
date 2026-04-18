package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// RegiaoRepository e a implementacao Postgres de domain.RegiaoRepository.
type RegiaoRepository struct {
	pool *pgxpool.Pool
}

func NewRegiaoRepository(pool *pgxpool.Pool) *RegiaoRepository {
	return &RegiaoRepository{pool: pool}
}

const colsRegiao = `
    id, nome, cidade, estado,
    COALESCE(cep_inicio, '') AS cep_inicio,
    COALESCE(cep_fim, '')    AS cep_fim,
    ativa, criado_em
`

func (r *RegiaoRepository) Criar(ctx context.Context, reg *domain.Regiao) error {
	const q = `
        INSERT INTO regioes (id, nome, cidade, estado, cep_inicio, cep_fim, ativa)
        VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), $7)
        RETURNING criado_em
    `
	err := r.pool.QueryRow(ctx, q,
		reg.ID, reg.Nome, reg.Cidade, reg.Estado, reg.CEPInicio, reg.CEPFim, reg.Ativa,
	).Scan(&reg.CriadoEm)
	return err
}

func (r *RegiaoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Regiao, error) {
	const q = `SELECT ` + colsRegiao + ` FROM regioes WHERE id = $1`
	var reg domain.Regiao
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&reg.ID, &reg.Nome, &reg.Cidade, &reg.Estado, &reg.CEPInicio, &reg.CEPFim,
		&reg.Ativa, &reg.CriadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRegiaoNaoEncontrada
		}
		return nil, err
	}
	return &reg, nil
}

func (r *RegiaoRepository) ListarAtivas(ctx context.Context) ([]*domain.Regiao, error) {
	const q = `SELECT ` + colsRegiao + ` FROM regioes WHERE ativa = TRUE ORDER BY estado, cidade, nome`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Regiao
	for rows.Next() {
		var reg domain.Regiao
		if err := rows.Scan(
			&reg.ID, &reg.Nome, &reg.Cidade, &reg.Estado, &reg.CEPInicio, &reg.CEPFim,
			&reg.Ativa, &reg.CriadoEm,
		); err != nil {
			return nil, err
		}
		out = append(out, &reg)
	}
	return out, rows.Err()
}

// ProfissionalRegiaoRepository e a implementacao Postgres de
// domain.ProfissionalRegiaoRepository (associacao N:N profissional <-> regiao).
type ProfissionalRegiaoRepository struct {
	pool *pgxpool.Pool
}

func NewProfissionalRegiaoRepository(pool *pgxpool.Pool) *ProfissionalRegiaoRepository {
	return &ProfissionalRegiaoRepository{pool: pool}
}

// DefinirRegioes substitui atomicamente o conjunto de regioes da profissional:
// apaga as existentes e insere as novas em uma mesma transacao.
func (r *ProfissionalRegiaoRepository) DefinirRegioes(ctx context.Context, profissionalID uuid.UUID, regiaoIDs []uuid.UUID) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `DELETE FROM regioes_atuacao WHERE profissional_id = $1`, profissionalID); err != nil {
		return err
	}

	if len(regiaoIDs) > 0 {
		// Insere em batch via unnest — uma unica ida ao banco.
		const q = `
            INSERT INTO regioes_atuacao (profissional_id, regiao_id)
            SELECT $1, regiao_id FROM unnest($2::uuid[]) AS t(regiao_id)
            ON CONFLICT DO NOTHING
        `
		if _, err := tx.Exec(ctx, q, profissionalID, regiaoIDs); err != nil {
			if pgErrorCode(err) == codeForeignKeyViolation {
				return domain.ErrRegiaoIDInvalida
			}
			return err
		}
	}

	return tx.Commit(ctx)
}

// ListarPorProfissionalID retorna todas as regioes (mesmo inativas) associadas
// a uma profissional. Usado para exibir o cadastro completo.
func (r *ProfissionalRegiaoRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Regiao, error) {
	const q = `
        SELECT r.id, r.nome, r.cidade, r.estado,
               COALESCE(r.cep_inicio, ''), COALESCE(r.cep_fim, ''),
               r.ativa, r.criado_em
          FROM regioes r
          JOIN regioes_atuacao ra ON ra.regiao_id = r.id
         WHERE ra.profissional_id = $1
         ORDER BY r.estado, r.cidade, r.nome
    `
	rows, err := r.pool.Query(ctx, q, profissionalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Regiao
	for rows.Next() {
		var reg domain.Regiao
		if err := rows.Scan(
			&reg.ID, &reg.Nome, &reg.Cidade, &reg.Estado, &reg.CEPInicio, &reg.CEPFim,
			&reg.Ativa, &reg.CriadoEm,
		); err != nil {
			return nil, err
		}
		out = append(out, &reg)
	}
	return out, rows.Err()
}
