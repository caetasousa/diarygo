package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// DocumentoRepository e a implementacao Postgres de domain.DocumentoRepository.
type DocumentoRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentoRepository(pool *pgxpool.Pool) *DocumentoRepository {
	return &DocumentoRepository{pool: pool}
}

const colsDocumento = `
    id, profissional_id, tipo, url, status,
    COALESCE(observacao, '') AS observacao,
    criado_em, atualizado_em
`

func (r *DocumentoRepository) Criar(ctx context.Context, d *domain.Documento) error {
	const q = `
        INSERT INTO documentos (id, profissional_id, tipo, url, status, observacao)
        VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''))
        RETURNING criado_em, atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		d.ID, d.ProfissionalID, string(d.Tipo), d.URL, string(d.Status), d.ObservacaoAdmin,
	).Scan(&d.CriadoEm, &d.AtualizadoEm)
	if err != nil {
		if pgErrorCode(err) == codeForeignKeyViolation {
			return domain.ErrProfissionalNaoEncontrada
		}
		return err
	}
	return nil
}

func (r *DocumentoRepository) Atualizar(ctx context.Context, d *domain.Documento) error {
	const q = `
        UPDATE documentos SET
            tipo = $2, url = $3, status = $4, observacao = NULLIF($5, '')
        WHERE id = $1
        RETURNING atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		d.ID, string(d.Tipo), d.URL, string(d.Status), d.ObservacaoAdmin,
	).Scan(&d.AtualizadoEm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrDocumentoNaoEncontrado
		}
		return err
	}
	return nil
}

func (r *DocumentoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Documento, error) {
	const q = `SELECT ` + colsDocumento + ` FROM documentos WHERE id = $1`
	var (
		d      domain.Documento
		tipo   string
		status string
	)
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&d.ID, &d.ProfissionalID, &tipo, &d.URL, &status, &d.ObservacaoAdmin,
		&d.CriadoEm, &d.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDocumentoNaoEncontrado
		}
		return nil, err
	}
	d.Tipo = domain.TipoDocumento(tipo)
	d.Status = domain.StatusDocumento(status)
	return &d, nil
}

func (r *DocumentoRepository) ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*domain.Documento, error) {
	const q = `
        SELECT ` + colsDocumento + `
          FROM documentos
         WHERE profissional_id = $1
         ORDER BY criado_em ASC
    `
	rows, err := r.pool.Query(ctx, q, profissionalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Documento
	for rows.Next() {
		var (
			d      domain.Documento
			tipo   string
			status string
		)
		if err := rows.Scan(
			&d.ID, &d.ProfissionalID, &tipo, &d.URL, &status, &d.ObservacaoAdmin,
			&d.CriadoEm, &d.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		d.Tipo = domain.TipoDocumento(tipo)
		d.Status = domain.StatusDocumento(status)
		out = append(out, &d)
	}
	return out, rows.Err()
}
