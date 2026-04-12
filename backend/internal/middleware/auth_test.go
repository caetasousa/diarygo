package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
)

const (
	testSecret = "segredo-de-teste-com-32-caracteres-ok"
	testEmail  = "mw@diarygo.com.br"
	testSenha  = "Senha123"
)

func novaService(t *testing.T) *service.AuthService {
	t.Helper()
	repo := memory.NewUsuarioRepository()
	return service.NewAuthServiceComCost(repo, testSecret, 15*time.Minute, "development", service.BcryptCostTeste)
}

func gerarToken(t *testing.T, svc *service.AuthService) string {
	t.Helper()
	ctx := context.Background()
	svc.Registrar(ctx, domain.RegistroRequest{
		Email: testEmail,
		Senha: testSenha,
	}, domain.TipoCliente)
	resp, err := svc.Login(ctx, domain.LoginRequest{
		Email: testEmail,
		Senha: testSenha,
	})
	if err != nil {
		t.Fatalf("falha ao gerar token: %v", err)
	}
	return resp.AccessToken
}

// --- Autenticar ---

func TestAutenticar_TokenValido_Retorna200(t *testing.T) {
	svc := novaService(t)
	token := gerarToken(t, svc)

	handler := mw.Autenticar(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d", w.Code)
	}
}

func TestAutenticar_SemHeader_Retorna401(t *testing.T) {
	svc := novaService(t)

	handler := mw.Autenticar(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

func TestAutenticar_TokenInvalido_Retorna401(t *testing.T) {
	svc := novaService(t)

	handler := mw.Autenticar(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token.invalido.aqui")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

func TestAutenticar_HeaderMalformado_Retorna401(t *testing.T) {
	svc := novaService(t)

	handler := mw.Autenticar(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "tokenSemBearer")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

// --- RequererTipo ---

func TestRequererTipo_TipoPermitido_Passa(t *testing.T) {
	svc := novaService(t)
	token := gerarToken(t, svc)

	// Token e CLIENTE — testar que CLIENTE passa
	handler := mw.Autenticar(svc)(
		mw.RequererTipo(domain.TipoCliente)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d", w.Code)
	}
}

func TestRequererTipo_TipoNegado_Retorna403(t *testing.T) {
	svc := novaService(t)
	token := gerarToken(t, svc)

	// Token e CLIENTE — ADMIN nao deve passar
	handler := mw.Autenticar(svc)(
		mw.RequererTipo(domain.TipoAdmin)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperava 403, got %d", w.Code)
	}
}

// --- UsuarioDoContexto ---

func TestUsuarioDoContexto_Presente_RetornaPayload(t *testing.T) {
	svc := novaService(t)
	token := gerarToken(t, svc)

	var capturedPayload *domain.TokenPayload

	handler := mw.Autenticar(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := mw.UsuarioDoContexto(r.Context())
		if !ok {
			t.Error("esperava payload no contexto")
		}
		capturedPayload = payload
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if capturedPayload == nil {
		t.Error("payload nao capturado")
	}
}

func TestUsuarioDoContexto_Ausente_RetornaFalse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, ok := mw.UsuarioDoContexto(req.Context())
	if ok {
		t.Error("esperava false para contexto sem payload")
	}
}
