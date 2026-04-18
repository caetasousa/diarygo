package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// ProfissionalRepository e a implementacao Postgres de domain.ProfissionalRepository.
type ProfissionalRepository struct {
	pool *pgxpool.Pool
}

func NewProfissionalRepository(pool *pgxpool.Pool) *ProfissionalRepository {
	return &ProfissionalRepository{pool: pool}
}

const colsProfissional = `
    id, usuario_id, nome, cpf, rg, telefone, foto_url, status,
    COALESCE(nota_media, 0)::float8 AS nota_media,
    total_servicos, mei, criado_em, atualizado_em
`

// Criar insere uma profissional. Mapeia unique violations:
//   - profissionais_usuario_id_key -> ErrProfissionalJaExiste
//   - profissionais_cpf_key        -> ErrCPFJaCadastrado
//
// nota_media: se vier zero do dominio (valor default do float), inserimos NULL
// para nao violar o CHECK (nota_media >= 1.0). nota_media so recebe valor
// efetivo apos a primeira avaliacao (Etapa 6).
func (r *ProfissionalRepository) Criar(ctx context.Context, p *domain.Profissional) error {
	var notaMedia *float64
	if p.NotaMedia > 0 {
		n := p.NotaMedia
		notaMedia = &n
	}

	const q = `
        INSERT INTO profissionais (
            id, usuario_id, nome, cpf, rg, telefone, foto_url, status,
            nota_media, total_servicos, mei
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING criado_em, atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		p.ID, p.UsuarioID, p.Nome, p.CPF, p.RG, p.Telefone, p.FotoURL,
		string(p.Status), notaMedia, p.TotalServicos, p.MEI,
	).Scan(&p.CriadoEm, &p.AtualizadoEm)
	if err != nil {
		if pgErrorCode(err) == codeUniqueViolation {
			switch pgErrorConstraint(err) {
			case "profissionais_usuario_id_key":
				return domain.ErrProfissionalJaExiste
			case "profissionais_cpf_key":
				return domain.ErrCPFJaCadastrado
			}
			return domain.ErrCPFJaCadastrado
		}
		return err
	}
	return nil
}

// Atualizar persiste alteracoes do perfil (exceto usuario_id).
func (r *ProfissionalRepository) Atualizar(ctx context.Context, p *domain.Profissional) error {
	var notaMedia *float64
	if p.NotaMedia > 0 {
		n := p.NotaMedia
		notaMedia = &n
	}

	const q = `
        UPDATE profissionais SET
            nome = $2, cpf = $3, rg = $4, telefone = $5, foto_url = $6,
            status = $7, nota_media = $8, total_servicos = $9, mei = $10
        WHERE id = $1
        RETURNING atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		p.ID, p.Nome, p.CPF, p.RG, p.Telefone, p.FotoURL,
		string(p.Status), notaMedia, p.TotalServicos, p.MEI,
	).Scan(&p.AtualizadoEm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrProfissionalNaoEncontrada
		}
		if pgErrorCode(err) == codeUniqueViolation {
			return domain.ErrCPFJaCadastrado
		}
		return err
	}
	return nil
}

func (r *ProfissionalRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Profissional, error) {
	return r.scanOne(ctx, `SELECT `+colsProfissional+` FROM profissionais WHERE id = $1`, id)
}

func (r *ProfissionalRepository) BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*domain.Profissional, error) {
	return r.scanOne(ctx, `SELECT `+colsProfissional+` FROM profissionais WHERE usuario_id = $1`, usuarioID)
}

func (r *ProfissionalRepository) scanOne(ctx context.Context, q string, args ...any) (*domain.Profissional, error) {
	var (
		p      domain.Profissional
		status string
	)
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&p.ID, &p.UsuarioID, &p.Nome, &p.CPF, &p.RG, &p.Telefone, &p.FotoURL, &status,
		&p.NotaMedia, &p.TotalServicos, &p.MEI, &p.CriadoEm, &p.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProfissionalNaoEncontrada
		}
		return nil, err
	}
	p.Status = domain.StatusProfissional(status)
	return &p, nil
}
