package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/domain"
)

// UsuarioRepository e a implementacao Postgres de domain.UsuarioRepository.
// Todas as queries usam parametros posicionais ($1, $2 ...) — OWASP A05.
type UsuarioRepository struct {
	pool *pgxpool.Pool
}

// NewUsuarioRepository cria um repository Postgres de usuarios.
func NewUsuarioRepository(pool *pgxpool.Pool) *UsuarioRepository {
	return &UsuarioRepository{pool: pool}
}

const colsUsuario = `
    id, email, senha_hash, tipo, email_verificado, ativo,
    token_recuperacao, token_recuperacao_expira,
    criado_em, atualizado_em
`

// Criar insere um novo usuario. Email e normalizado para lowercase.
// Mapeia unique violation -> domain.ErrEmailJaExiste.
func (r *UsuarioRepository) Criar(ctx context.Context, u *domain.Usuario) error {
	emailNorm := strings.ToLower(strings.TrimSpace(u.Email))

	var tokenRec *string
	if u.TokenRecuperacao != "" {
		tr := u.TokenRecuperacao
		tokenRec = &tr
	}
	var tokenExp *time.Time
	if !u.TokenRecuperacaoExpira.IsZero() {
		te := u.TokenRecuperacaoExpira
		tokenExp = &te
	}

	const q = `
        INSERT INTO usuarios (
            id, email, senha_hash, tipo, email_verificado, ativo,
            token_recuperacao, token_recuperacao_expira
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING criado_em, atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		u.ID, emailNorm, u.SenhaHash, string(u.Tipo), u.EmailVerificado, u.Ativo,
		tokenRec, tokenExp,
	).Scan(&u.CriadoEm, &u.AtualizadoEm)
	if err != nil {
		if pgErrorCode(err) == codeUniqueViolation {
			return domain.ErrEmailJaExiste
		}
		return err
	}
	u.Email = emailNorm
	return nil
}

// Atualizar persiste alteracoes permitidas (senha, flags, tokens). Nao permite
// troca de email/tipo via esta rota (nao ha caso de uso ainda; se surgir, criar
// metodo dedicado para auditar).
func (r *UsuarioRepository) Atualizar(ctx context.Context, u *domain.Usuario) error {
	var tokenRec *string
	if u.TokenRecuperacao != "" {
		tr := u.TokenRecuperacao
		tokenRec = &tr
	}
	var tokenExp *time.Time
	if !u.TokenRecuperacaoExpira.IsZero() {
		te := u.TokenRecuperacaoExpira
		tokenExp = &te
	}

	const q = `
        UPDATE usuarios SET
            senha_hash               = $2,
            email_verificado         = $3,
            ativo                    = $4,
            token_recuperacao        = $5,
            token_recuperacao_expira = $6
        WHERE id = $1
        RETURNING atualizado_em
    `
	err := r.pool.QueryRow(ctx, q,
		u.ID, u.SenhaHash, u.EmailVerificado, u.Ativo, tokenRec, tokenExp,
	).Scan(&u.AtualizadoEm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUsuarioNaoEncontrado
		}
		return err
	}
	return nil
}

// BuscarPorEmail retorna o usuario pelo email (case-insensitive).
func (r *UsuarioRepository) BuscarPorEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	emailNorm := strings.ToLower(strings.TrimSpace(email))
	const q = `SELECT ` + colsUsuario + ` FROM usuarios WHERE email = $1`
	return r.scanOne(ctx, q, emailNorm)
}

// BuscarPorID retorna o usuario pelo UUID.
func (r *UsuarioRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Usuario, error) {
	const q = `SELECT ` + colsUsuario + ` FROM usuarios WHERE id = $1`
	return r.scanOne(ctx, q, id)
}

// BuscarPorTokenRecuperacao retorna o usuario pelo token de recuperacao. Usa o
// indice parcial idx_usuarios_token_recuperacao.
func (r *UsuarioRepository) BuscarPorTokenRecuperacao(ctx context.Context, token string) (*domain.Usuario, error) {
	if token == "" {
		return nil, domain.ErrUsuarioNaoEncontrado
	}
	const q = `SELECT ` + colsUsuario + ` FROM usuarios WHERE token_recuperacao = $1`
	return r.scanOne(ctx, q, token)
}

func (r *UsuarioRepository) scanOne(ctx context.Context, q string, args ...any) (*domain.Usuario, error) {
	var (
		u        domain.Usuario
		tipo     string
		tokenRec *string
		tokenExp *time.Time
	)
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&u.ID, &u.Email, &u.SenhaHash, &tipo, &u.EmailVerificado, &u.Ativo,
		&tokenRec, &tokenExp,
		&u.CriadoEm, &u.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUsuarioNaoEncontrado
		}
		return nil, err
	}
	u.Tipo = domain.TipoUsuario(tipo)
	if tokenRec != nil {
		u.TokenRecuperacao = *tokenRec
	}
	if tokenExp != nil {
		u.TokenRecuperacaoExpira = *tokenExp
	}
	return &u, nil
}
