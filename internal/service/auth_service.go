package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// BcryptCost e o custo usado para hash de senhas.
// Em testes, use BcryptCostTeste (4) para acelerar — cost 12 leva ~250ms por hash.
const (
	BcryptCost      = 12
	BcryptCostTeste = 4
)

// dummyHash e usado para evitar timing attack quando o usuario nao existe.
// Deve ser um hash bcrypt valido (qualquer senha).
var dummyHash []byte

func init() {
	h, err := bcrypt.GenerateFromPassword([]byte("dummy-para-timing-attack"), BcryptCostTeste)
	if err != nil {
		panic(fmt.Sprintf("falha ao gerar dummy hash: %v", err))
	}
	dummyHash = h
}

// AuthService orquestra as operacoes de autenticacao.
type AuthService struct {
	repo       domain.UsuarioRepository
	jwtSecret  []byte
	jwtExpiry  time.Duration
	env        string
	bcryptCost int
}

// NewAuthService cria um novo AuthService.
func NewAuthService(repo domain.UsuarioRepository, jwtSecret string, jwtExpiry time.Duration, env string) *AuthService {
	cost := BcryptCost
	if env != "production" {
		// Em desenvolvimento/testes, usar cost menor para acelerar
		// mas ainda seguro o suficiente
		cost = BcryptCost
	}
	return &AuthService{
		repo:       repo,
		jwtSecret:  []byte(jwtSecret),
		jwtExpiry:  jwtExpiry,
		env:        env,
		bcryptCost: cost,
	}
}

// NewAuthServiceComCost cria AuthService com custo bcrypt customizado (util para testes).
func NewAuthServiceComCost(repo domain.UsuarioRepository, jwtSecret string, jwtExpiry time.Duration, env string, cost int) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtSecret:  []byte(jwtSecret),
		jwtExpiry:  jwtExpiry,
		env:        env,
		bcryptCost: cost,
	}
}

// Registrar valida, cria e persiste um novo usuario.
func (s *AuthService) Registrar(ctx context.Context, req domain.RegistroRequest, tipo domain.TipoUsuario) (*domain.RegistroResponse, error) {
	// Normalizar email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Validacoes de dominio
	if err := domain.ValidarEmail(req.Email); err != nil {
		return nil, err
	}
	if err := domain.ValidarSenha(req.Senha); err != nil {
		return nil, err
	}
	if err := domain.ValidarTipoUsuario(tipo); err != nil {
		return nil, err
	}

	// Hash da senha — OWASP A04: bcrypt cost >= 12
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Senha), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao processar senha: %w", err)
	}

	u := &domain.Usuario{
		ID:                uuid.New(),
		Email:             req.Email,
		SenhaHash:         string(hash),
		Tipo:              tipo,
		EmailVerificado:   false,
		Ativo:             true,
		CodigoVerificacao: uuid.New().String(),
	}

	if err := s.repo.Criar(ctx, u); err != nil {
		return nil, err
	}

	slog.Info("usuario registrado",
		"usuario_id", u.ID,
		"tipo", u.Tipo,
	)

	resp := &domain.RegistroResponse{
		ID:    u.ID,
		Email: u.Email,
		Tipo:  u.Tipo,
	}

	// Retornar codigo de verificacao apenas em development (sem infra de email no MVP)
	if s.env != "production" {
		resp.CodigoVerificacao = u.CodigoVerificacao
	}

	return resp, nil
}

// Login autentica um usuario e retorna um JWT.
// OWASP A07: executa bcrypt nos dois caminhos (usuario existe e nao existe)
// para evitar timing attack que revelaria emails cadastrados.
func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.TokenResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	u, err := s.repo.BuscarPorEmail(ctx, req.Email)
	if err != nil {
		// Timing attack protection: bcrypt mesmo com usuario inexistente
		bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Senha)) //nolint:errcheck
		slog.Warn("tentativa de login com email inexistente", "email", req.Email)
		return nil, domain.ErrCredenciaisInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.SenhaHash), []byte(req.Senha)); err != nil {
		slog.Warn("tentativa de login com senha incorreta", "usuario_id", u.ID)
		return nil, domain.ErrCredenciaisInvalidas
	}

	if !u.Ativo {
		return nil, domain.ErrUsuarioInativo
	}

	token, err := s.GerarToken(u)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	slog.Info("login realizado com sucesso", "usuario_id", u.ID, "tipo", u.Tipo)
	return token, nil
}

// VerificarEmail valida o codigo e marca o email como verificado.
func (s *AuthService) VerificarEmail(ctx context.Context, req domain.VerificarEmailRequest) error {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	u, err := s.repo.BuscarPorEmail(ctx, req.Email)
	if err != nil {
		return domain.ErrUsuarioNaoEncontrado
	}

	if u.CodigoVerificacao == "" || u.CodigoVerificacao != req.Codigo {
		return domain.ErrCodigoVerificacaoInvalido
	}

	u.EmailVerificado = true
	u.CodigoVerificacao = "" // limpar apos uso

	if err := s.repo.Atualizar(ctx, u); err != nil {
		return fmt.Errorf("erro ao atualizar usuario: %w", err)
	}

	slog.Info("email verificado com sucesso", "usuario_id", u.ID)
	return nil
}

// GerarToken cria um JWT assinado com as claims do usuario.
func (s *AuthService) GerarToken(u *domain.Usuario) (*domain.TokenResponse, error) {
	expiresAt := time.Now().Add(s.jwtExpiry)

	claims := domain.TokenPayload{
		UsuarioID: u.ID,
		Email:     u.Email,
		Tipo:      u.Tipo,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "diarygo",
			Subject:   u.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("erro ao assinar token: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.jwtExpiry.Seconds()),
	}, nil
}

// ValidarToken valida um JWT e retorna as claims.
// OWASP A04: aceita SOMENTE HS256 — bloqueia ataque "alg: none".
func (s *AuthService) ValidarToken(tokenString string) (*domain.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&domain.TokenPayload{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo nao permitido: %v", t.Header["alg"])
			}
			return s.jwtSecret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, domain.ErrTokenInvalido
	}

	claims, ok := token.Claims.(*domain.TokenPayload)
	if !ok || !token.Valid {
		return nil, domain.ErrTokenInvalido
	}

	return claims, nil
}
