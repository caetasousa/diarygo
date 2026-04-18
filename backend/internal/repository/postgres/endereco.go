package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// EnderecoRepository e a implementacao Postgres de domain.EnderecoRepository.
type EnderecoRepository struct {
	pool *pgxpool.Pool
}

func NewEnderecoRepository(pool *pgxpool.Pool) *EnderecoRepository {
	return &EnderecoRepository{pool: pool}
}

const colsEndereco = `
    id, cliente_id, logradouro, numero, COALESCE(complemento, '') AS complemento,
    bairro, cidade, estado, cep,
    COALESCE(lat, 0)::float8 AS lat,
    COALESCE(lon, 0)::float8 AS lon,
    principal, num_quartos, num_banheiros, num_salas, num_cozinhas,
    COALESCE(area_m2, 0)::float8 AS area_m2,
    criado_em, atualizado_em
`

// Criar insere um endereco. lat/lon e area_m2 iguais a zero sao gravados como NULL.
func (r *EnderecoRepository) Criar(ctx context.Context, e *domain.Endereco) error {
	var lat, lon, area *float64
	if e.Latitude != 0 {
		v := e.Latitude
		lat = &v
	}
	if e.Longitude != 0 {
		v := e.Longitude
		lon = &v
	}
	if e.AreaM2 != 0 {
		v := e.AreaM2
		area = &v
	}

	const q = `
        INSERT INTO enderecos (
            id, cliente_id, logradouro, numero, complemento, bairro, cidade, estado, cep,
            lat, lon, principal, num_quartos, num_banheiros, num_salas, num_cozinhas, area_m2
        ) VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING criado_em, atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		e.ID, e.ClienteID, e.Logradouro, e.Numero, e.Complemento, e.Bairro, e.Cidade, e.Estado, e.CEP,
		lat, lon, e.Principal, e.NumQuartos, e.NumBanheiros, e.NumSalas, e.NumCozinhas, area,
	).Scan(&e.CriadoEm, &e.AtualizadoEm)
	if err != nil {
		if pgErrorCode(err) == codeForeignKeyViolation {
			return domain.ErrClienteNaoEncontrado
		}
		return err
	}
	return nil
}

func (r *EnderecoRepository) Atualizar(ctx context.Context, e *domain.Endereco) error {
	var lat, lon, area *float64
	if e.Latitude != 0 {
		v := e.Latitude
		lat = &v
	}
	if e.Longitude != 0 {
		v := e.Longitude
		lon = &v
	}
	if e.AreaM2 != 0 {
		v := e.AreaM2
		area = &v
	}

	const q = `
        UPDATE enderecos SET
            logradouro = $2, numero = $3, complemento = NULLIF($4, ''),
            bairro = $5, cidade = $6, estado = $7, cep = $8,
            lat = $9, lon = $10, principal = $11,
            num_quartos = $12, num_banheiros = $13, num_salas = $14, num_cozinhas = $15,
            area_m2 = $16
        WHERE id = $1
        RETURNING atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		e.ID, e.Logradouro, e.Numero, e.Complemento, e.Bairro, e.Cidade, e.Estado, e.CEP,
		lat, lon, e.Principal, e.NumQuartos, e.NumBanheiros, e.NumSalas, e.NumCozinhas, area,
	).Scan(&e.AtualizadoEm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrEnderecoNaoEncontrado
		}
		return err
	}
	return nil
}

func (r *EnderecoRepository) Remover(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM enderecos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrEnderecoNaoEncontrado
	}
	return nil
}

func (r *EnderecoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Endereco, error) {
	const q = `SELECT ` + colsEndereco + ` FROM enderecos WHERE id = $1`
	var e domain.Endereco
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&e.ID, &e.ClienteID, &e.Logradouro, &e.Numero, &e.Complemento,
		&e.Bairro, &e.Cidade, &e.Estado, &e.CEP,
		&e.Latitude, &e.Longitude,
		&e.Principal, &e.NumQuartos, &e.NumBanheiros, &e.NumSalas, &e.NumCozinhas, &e.AreaM2,
		&e.CriadoEm, &e.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEnderecoNaoEncontrado
		}
		return nil, err
	}
	return &e, nil
}

func (r *EnderecoRepository) ListarPorClienteID(ctx context.Context, clienteID uuid.UUID) ([]*domain.Endereco, error) {
	const q = `
        SELECT ` + colsEndereco + `
          FROM enderecos
         WHERE cliente_id = $1
         ORDER BY principal DESC, criado_em ASC
    `
	rows, err := r.pool.Query(ctx, q, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Endereco
	for rows.Next() {
		var e domain.Endereco
		if err := rows.Scan(
			&e.ID, &e.ClienteID, &e.Logradouro, &e.Numero, &e.Complemento,
			&e.Bairro, &e.Cidade, &e.Estado, &e.CEP,
			&e.Latitude, &e.Longitude,
			&e.Principal, &e.NumQuartos, &e.NumBanheiros, &e.NumSalas, &e.NumCozinhas, &e.AreaM2,
			&e.CriadoEm, &e.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}
