package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
)

const (
	testSecret = "segredo-de-teste-com-32-caracteres-ok"
	testEmail  = "teste@diarygo.com.br"
	testSenha  = "Senha123"
)

// novaService cria um AuthService com custo bcrypt baixo para testes rapidos.
func novaService(t *testing.T) *service.AuthService {
	t.Helper()
	repo := memory.NewUsuarioRepository()
	return service.NewAuthServiceComCost(repo, testSecret, 15*time.Minute, "development", service.BcryptCostTeste)
}

func registrarUsuario(t *testing.T, svc *service.AuthService, email string, tipo domain.TipoUsuario) *domain.RegistroResponse {
	t.Helper()
	resp, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: email,
		Senha: testSenha,
	}, tipo)
	if err != nil {
		t.Fatalf("falha ao registrar usuario: %v", err)
	}
	return resp
}

// --- Registrar ---

func TestRegistrar_EmailValido_CriaUsuario(t *testing.T) {
	svc := novaService(t)
	resp, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoCliente)

	if err != nil {
		t.Fatalf("esperava nil, got: %v", err)
	}
	if resp.Email != testEmail {
		t.Errorf("email esperado %s, got %s", testEmail, resp.Email)
	}
	if resp.Tipo != domain.TipoCliente {
		t.Errorf("tipo esperado CLIENTE, got %s", resp.Tipo)
	}
	if resp.CodigoVerificacao == "" {
		t.Error("esperava codigo de verificacao em development")
	}
}

func TestRegistrar_EmailDuplicado_RetornaErro(t *testing.T) {
	svc := novaService(t)
	registrarUsuario(t, svc, testEmail, domain.TipoCliente)

	_, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoCliente)

	if err == nil {
		t.Fatal("esperava erro de email duplicado")
	}
	if !isErro(err, domain.ErrEmailJaExiste) {
		t.Errorf("esperava ErrEmailJaExiste, got: %v", err)
	}
}

func TestRegistrar_SenhaFraca_RetornaErro(t *testing.T) {
	svc := novaService(t)
	_, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: testEmail,
		Senha: "curta",
	}, domain.TipoCliente)

	if err == nil {
		t.Fatal("esperava erro de senha fraca")
	}
}

func TestRegistrar_EmailInvalido_RetornaErro(t *testing.T) {
	svc := novaService(t)
	_, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: "nao-e-um-email",
		Senha: testSenha,
	}, domain.TipoCliente)

	if err == nil {
		t.Fatal("esperava erro de email invalido")
	}
}

func TestRegistrar_TipoCliente_DefineCorreto(t *testing.T) {
	svc := novaService(t)
	resp, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoCliente)

	if err != nil {
		t.Fatal(err)
	}
	if resp.Tipo != domain.TipoCliente {
		t.Errorf("esperava TipoCliente, got %s", resp.Tipo)
	}
}

func TestRegistrar_TipoProfissional_DefineCorreto(t *testing.T) {
	svc := novaService(t)
	resp, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoProfissional)

	if err != nil {
		t.Fatal(err)
	}
	if resp.Tipo != domain.TipoProfissional {
		t.Errorf("esperava TipoProfissional, got %s", resp.Tipo)
	}
}

func TestRegistrar_TipoAdmin_RetornaErro(t *testing.T) {
	svc := novaService(t)
	_, err := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoAdmin)

	if err == nil {
		t.Fatal("esperava erro ao registrar ADMIN via endpoint publico")
	}
}

// --- Login ---

func TestLogin_CredenciaisValidas_RetornaToken(t *testing.T) {
	svc := novaService(t)
	registrarUsuario(t, svc, testEmail, domain.TipoCliente)

	resp, err := svc.Login(context.Background(), domain.LoginRequest{
		Email: testEmail,
		Senha: testSenha,
	})

	if err != nil {
		t.Fatalf("esperava nil, got: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("esperava access token nao vazio")
	}
	if resp.TokenType != "Bearer" {
		t.Errorf("esperava Bearer, got %s", resp.TokenType)
	}
}

func TestLogin_SenhaErrada_RetornaErroGenerico(t *testing.T) {
	svc := novaService(t)
	registrarUsuario(t, svc, testEmail, domain.TipoCliente)

	_, err := svc.Login(context.Background(), domain.LoginRequest{
		Email: testEmail,
		Senha: "SenhaErrada123",
	})

	if err == nil {
		t.Fatal("esperava erro")
	}
	// OWASP A07: mensagem generica — nao revela "senha incorreta" vs "email nao existe"
	if !isErro(err, domain.ErrCredenciaisInvalidas) {
		t.Errorf("esperava ErrCredenciaisInvalidas, got: %v", err)
	}
}

func TestLogin_EmailInexistente_RetornaErroGenerico(t *testing.T) {
	svc := novaService(t)

	// Timing attack: mesmo erro mesmo sem o usuario existir
	_, err := svc.Login(context.Background(), domain.LoginRequest{
		Email: "naoexiste@diarygo.com.br",
		Senha: testSenha,
	})

	if err == nil {
		t.Fatal("esperava erro")
	}
	if !isErro(err, domain.ErrCredenciaisInvalidas) {
		t.Errorf("esperava ErrCredenciaisInvalidas, got: %v", err)
	}
}

func TestLogin_UsuarioInativo_RetornaErro(t *testing.T) {
	repo := memory.NewUsuarioRepository()
	svc := service.NewAuthServiceComCost(repo, testSecret, 15*time.Minute, "development", service.BcryptCostTeste)

	// Registrar e depois desativar manualmente
	resp, _ := svc.Registrar(context.Background(), domain.RegistroRequest{
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoCliente)

	u, _ := repo.BuscarPorID(context.Background(), resp.ID)
	u.Ativo = false
	repo.Atualizar(context.Background(), u) //nolint:errcheck

	_, err := svc.Login(context.Background(), domain.LoginRequest{
		Email: testEmail,
		Senha: testSenha,
	})

	if err == nil {
		t.Fatal("esperava erro")
	}
	if !isErro(err, domain.ErrUsuarioInativo) {
		t.Errorf("esperava ErrUsuarioInativo, got: %v", err)
	}
}

// --- GerarToken / ValidarToken ---

func TestValidarToken_TokenValido_RetornaPayload(t *testing.T) {
	svc := novaService(t)
	resp := registrarUsuario(t, svc, testEmail, domain.TipoCliente)

	// Login para obter token real
	tokenResp, err := svc.Login(context.Background(), domain.LoginRequest{
		Email: testEmail,
		Senha: testSenha,
	})
	if err != nil {
		t.Fatal(err)
	}

	claims, err := svc.ValidarToken(tokenResp.AccessToken)
	if err != nil {
		t.Fatalf("esperava nil, got: %v", err)
	}
	if claims.UsuarioID != resp.ID {
		t.Errorf("ID esperado %s, got %s", resp.ID, claims.UsuarioID)
	}
	if claims.Tipo != domain.TipoCliente {
		t.Errorf("tipo esperado CLIENTE, got %s", claims.Tipo)
	}
}

func TestValidarToken_TokenExpirado_RetornaErro(t *testing.T) {
	svc := service.NewAuthServiceComCost(
		memory.NewUsuarioRepository(),
		testSecret,
		-1*time.Second, // expiracao no passado
		"development",
		service.BcryptCostTeste,
	)

	repo := memory.NewUsuarioRepository()
	svc2 := service.NewAuthServiceComCost(repo, testSecret, 15*time.Minute, "development", service.BcryptCostTeste)
	registrarUsuario(t, svc2, testEmail, domain.TipoCliente)
	tokenResp, _ := svc2.Login(context.Background(), domain.LoginRequest{
		Email: testEmail,
		Senha: testSenha,
	})

	// Usar svc com expiracao negativa para simular token expirado via secret diferente
	// — na pratica, testar com token gerado com expiry negativo
	repo2 := memory.NewUsuarioRepository()
	svcExp := service.NewAuthServiceComCost(repo2, testSecret, -1*time.Second, "development", service.BcryptCostTeste)
	svcExp.Registrar(context.Background(), domain.RegistroRequest{ //nolint:errcheck
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoCliente)
	expiredResp, _ := svcExp.Login(context.Background(), domain.LoginRequest{
		Email: testEmail,
		Senha: testSenha,
	})

	_ = svc // apenas para usar a variavel
	_ = tokenResp

	_, err := svc2.ValidarToken(expiredResp.AccessToken)
	if err == nil {
		t.Fatal("esperava erro para token expirado")
	}
}

func TestValidarToken_SecretErrado_RetornaErro(t *testing.T) {
	svc := novaService(t)
	registrarUsuario(t, svc, testEmail, domain.TipoCliente)
	tokenResp, _ := svc.Login(context.Background(), domain.LoginRequest{
		Email: testEmail,
		Senha: testSenha,
	})

	// Servico com secret diferente nao deve validar o token
	svcOutro := service.NewAuthServiceComCost(
		memory.NewUsuarioRepository(),
		"outro-secret-completamente-diferente-123",
		15*time.Minute,
		"development",
		service.BcryptCostTeste,
	)

	_, err := svcOutro.ValidarToken(tokenResp.AccessToken)
	if err == nil {
		t.Fatal("esperava erro para token com secret errado")
	}
}

// --- VerificarEmail ---

func TestVerificarEmail_CodigoValido_MarcaVerificado(t *testing.T) {
	svc := novaService(t)
	resp := registrarUsuario(t, svc, testEmail, domain.TipoCliente)

	err := svc.VerificarEmail(context.Background(), domain.VerificarEmailRequest{
		Email:  testEmail,
		Codigo: resp.CodigoVerificacao,
	})

	if err != nil {
		t.Fatalf("esperava nil, got: %v", err)
	}
}

func TestVerificarEmail_CodigoInvalido_RetornaErro(t *testing.T) {
	svc := novaService(t)
	registrarUsuario(t, svc, testEmail, domain.TipoCliente)

	err := svc.VerificarEmail(context.Background(), domain.VerificarEmailRequest{
		Email:  testEmail,
		Codigo: "codigo-errado",
	})

	if err == nil {
		t.Fatal("esperava erro")
	}
	if !isErro(err, domain.ErrCodigoVerificacaoInvalido) {
		t.Errorf("esperava ErrCodigoVerificacaoInvalido, got: %v", err)
	}
}

func TestVerificarEmail_EmailInexistente_RetornaErro(t *testing.T) {
	svc := novaService(t)

	err := svc.VerificarEmail(context.Background(), domain.VerificarEmailRequest{
		Email:  "naoexiste@diarygo.com.br",
		Codigo: "qualquer-codigo",
	})

	if err == nil {
		t.Fatal("esperava erro")
	}
	if !isErro(err, domain.ErrUsuarioNaoEncontrado) {
		t.Errorf("esperava ErrUsuarioNaoEncontrado, got: %v", err)
	}
}

// --- helpers ---

func isErro(got, esperado error) bool {
	return strings.Contains(got.Error(), esperado.Error()) || got == esperado
}
