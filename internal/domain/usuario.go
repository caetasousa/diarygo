package domain

import (
	"errors"
	"net/mail"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// TipoUsuario representa os perfis do sistema.
type TipoUsuario string

const (
	TipoCliente      TipoUsuario = "CLIENTE"
	TipoProfissional TipoUsuario = "PROFISSIONAL"
	TipoAdmin        TipoUsuario = "ADMIN"
)

// Usuario representa a entidade de autenticacao — espelha a tabela `usuarios` do banco.
type Usuario struct {
	ID                uuid.UUID
	Email             string
	SenhaHash         string
	Tipo              TipoUsuario
	EmailVerificado   bool
	Ativo             bool
	CodigoVerificacao string // nao persiste no SQL; usado para verificacao de email no MVP
	CriadoEm          time.Time
	AtualizadoEm      time.Time
}

// Erros de dominio (sentinela) — usados para mapear status HTTP nos handlers.
var (
	ErrEmailJaExiste             = errors.New("email ja cadastrado")
	ErrCredenciaisInvalidas      = errors.New("credenciais invalidas")
	ErrUsuarioNaoEncontrado      = errors.New("usuario nao encontrado")
	ErrUsuarioInativo            = errors.New("usuario inativo")
	ErrEmailNaoVerificado        = errors.New("email nao verificado")
	ErrTokenInvalido             = errors.New("token invalido")
	ErrCodigoVerificacaoInvalido = errors.New("codigo de verificacao invalido")
)

// ValidarEmail verifica o formato do email (RFC 5322).
func ValidarEmail(email string) error {
	if email == "" {
		return errors.New("email e obrigatorio")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("formato de email invalido")
	}
	return nil
}

// ValidarSenha verifica o comprimento minimo e maximo da senha.
// NIST 800-63b: nao exigir complexidade artificial; apenas comprimento.
func ValidarSenha(senha string) error {
	n := utf8.RuneCountInString(senha)
	if n < 8 {
		return errors.New("senha deve ter no minimo 8 caracteres")
	}
	if n > 72 { // limite do bcrypt
		return errors.New("senha deve ter no maximo 72 caracteres")
	}
	return nil
}

// ValidarTipoUsuario verifica se o tipo e valido para registro publico.
// Admin nao pode ser criado via endpoint publico.
func ValidarTipoUsuario(tipo TipoUsuario) error {
	switch tipo {
	case TipoCliente, TipoProfissional:
		return nil
	case TipoAdmin:
		return errors.New("tipo ADMIN nao pode ser registrado publicamente")
	default:
		return errors.New("tipo de usuario invalido")
	}
}
