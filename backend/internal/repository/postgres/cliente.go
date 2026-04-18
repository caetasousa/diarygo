package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// ClienteRepository e a implementacao Postgres de domain.ClienteRepository.
type ClienteRepository struct {
	pool *pgxpool.Pool
}

// NewClienteRepository cria um repository Postgres de clientes.
func NewClienteRepository(pool *pgxpool.Pool) *ClienteRepository {
	return &ClienteRepository{pool: pool}
}

const colsCliente = `id, usuario_id, nome, cpf, telefone, score, criado_em, atualizado_em`

// Criar insere um cliente. Mapeia violacoes de unicidade:
//   - usuarios_id unique -> ErrClienteJaExiste
//   - cpf unique         -> ErrCPFJaCadastrado
func (r *ClienteRepository) Criar(ctx context.Context, c *domain.Cliente) error {
	const q = `
        INSERT INTO clientes (id, usuario_id, nome, cpf, telefone, score)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING criado_em, atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		c.ID, c.UsuarioID, c.Nome, c.CPF, c.Telefone, c.Score,
	).Scan(&c.CriadoEm, &c.AtualizadoEm)
	if err != nil {
		if pgErrorCode(err) == codeUniqueViolation {
			switch pgErrorConstraint(err) {
			case "clientes_usuario_id_key":
				return domain.ErrClienteJaExiste
			case "clientes_cpf_key":
				return domain.ErrCPFJaCadastrado
			}
			return domain.ErrCPFJaCadastrado
		}
		return err
	}
	return nil
}

// Atualizar atualiza nome, cpf, telefone e score. O usuario_id nao e alteravel.
func (r *ClienteRepository) Atualizar(ctx context.Context, c *domain.Cliente) error {
	const q = `
        UPDATE clientes
           SET nome = $2, cpf = $3, telefone = $4, score = $5
         WHERE id = $1
        RETURNING atualizado_em
    `
	err := r.pool.QueryRow(ctx, q, c.ID, c.Nome, c.CPF, c.Telefone, c.Score).
		Scan(&c.AtualizadoEm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrClienteNaoEncontrado
		}
		if pgErrorCode(err) == codeUniqueViolation {
			return domain.ErrCPFJaCadastrado
		}
		return err
	}
	return nil
}

func (r *ClienteRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Cliente, error) {
	return r.scanOne(ctx, `SELECT `+colsCliente+` FROM clientes WHERE id = $1`, id)
}

func (r *ClienteRepository) BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*domain.Cliente, error) {
	return r.scanOne(ctx, `SELECT `+colsCliente+` FROM clientes WHERE usuario_id = $1`, usuarioID)
}

func (r *ClienteRepository) BuscarPorCPF(ctx context.Context, cpf string) (*domain.Cliente, error) {
	return r.scanOne(ctx, `SELECT `+colsCliente+` FROM clientes WHERE cpf = $1`, cpf)
}

func (r *ClienteRepository) scanOne(ctx context.Context, q string, args ...any) (*domain.Cliente, error) {
	var c domain.Cliente
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&c.ID, &c.UsuarioID, &c.Nome, &c.CPF, &c.Telefone, &c.Score,
		&c.CriadoEm, &c.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrClienteNaoEncontrado
		}
		return nil, err
	}
	return &c, nil
}
