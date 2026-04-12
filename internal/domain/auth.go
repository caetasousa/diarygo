package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UsuarioReader define operacoes de leitura do repositorio de usuarios.
type UsuarioReader interface {
	BuscarPorEmail(ctx context.Context, email string) (*Usuario, error)
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Usuario, error)
}

// UsuarioWriter define operacoes de escrita do repositorio de usuarios.
type UsuarioWriter interface {
	Criar(ctx context.Context, usuario *Usuario) error
	Atualizar(ctx context.Context, usuario *Usuario) error
}

// UsuarioRepository e a interface completa do repositorio (composicao Go idiomatica).
type UsuarioRepository interface {
	UsuarioReader
	UsuarioWriter
}

// RegistroRequest representa o payload de entrada para registro de usuario.
type RegistroRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// LoginRequest representa o payload de entrada para login.
type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// VerificarEmailRequest representa o payload para verificacao de email.
type VerificarEmailRequest struct {
	Email  string `json:"email"`
	Codigo string `json:"codigo"`
}

// TokenResponse e a resposta retornada apos login bem-sucedido.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // segundos
}

// RegistroResponse e a resposta retornada apos registro bem-sucedido.
type RegistroResponse struct {
	ID                uuid.UUID   `json:"id"`
	Email             string      `json:"email"`
	Tipo              TipoUsuario `json:"tipo"`
	CodigoVerificacao string      `json:"codigo_verificacao,omitempty"` // apenas em ENV=development
}

// TokenPayload representa as claims customizadas do JWT.
type TokenPayload struct {
	UsuarioID uuid.UUID   `json:"sub"`
	Email     string      `json:"email"`
	Tipo      TipoUsuario `json:"tipo"`
	jwt.RegisteredClaims
}
