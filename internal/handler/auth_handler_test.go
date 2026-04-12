package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/handler"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
)

const (
	testSecret = "segredo-de-teste-com-32-caracteres-ok"
	testEmail  = "teste@diarygo.com.br"
	testSenha  = "Senha123"
)

// novoRouter monta um chi.Router completo para testes, sem httprate.
func novoRouter(t *testing.T) (*chi.Mux, *service.AuthService) {
	t.Helper()
	repo := memory.NewUsuarioRepository()
	svc := service.NewAuthServiceComCost(repo, testSecret, 15*time.Minute, "development", service.BcryptCostTeste)
	h := handler.NewAuthHandler(svc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/auth", h.Routes())
		r.Group(func(r chi.Router) {
			r.Use(mw.Autenticar(svc))
			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				payload, ok := mw.UsuarioDoContexto(r.Context())
				if !ok {
					handler.RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
					return
				}
				handler.RespostaJSON(w, http.StatusOK, payload)
			})
		})
	})

	return r, svc
}

func doRequest(router http.Handler, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body) //nolint:errcheck
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// --- Registro ---

func TestRegistroCliente_Sucesso_Retorna201(t *testing.T) {
	r, _ := novoRouter(t)
	w := doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", map[string]string{
		"email": testEmail,
		"senha": testSenha,
	}, "")

	if w.Code != http.StatusCreated {
		t.Errorf("esperava 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp domain.RegistroResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("falha ao decodificar resposta:", err)
	}
	if resp.Email != testEmail {
		t.Errorf("email esperado %s, got %s", testEmail, resp.Email)
	}
}

func TestRegistroCliente_EmailDuplicado_Retorna409(t *testing.T) {
	r, _ := novoRouter(t)
	body := map[string]string{"email": testEmail, "senha": testSenha}

	doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", body, "")
	w := doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", body, "")

	if w.Code != http.StatusConflict {
		t.Errorf("esperava 409, got %d", w.Code)
	}
}

func TestRegistroCliente_BodyInvalido_Retorna400(t *testing.T) {
	r, _ := novoRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/registro/cliente",
		bytes.NewBufferString("nao e json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, got %d", w.Code)
	}
}

func TestRegistroCliente_SenhaFraca_Retorna400(t *testing.T) {
	r, _ := novoRouter(t)
	w := doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", map[string]string{
		"email": testEmail,
		"senha": "curta",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, got %d", w.Code)
	}
}

func TestRegistroProfissional_Sucesso_Retorna201(t *testing.T) {
	r, _ := novoRouter(t)
	w := doRequest(r, http.MethodPost, "/api/v1/auth/registro/profissional", map[string]string{
		"email": "profissional@diarygo.com.br",
		"senha": testSenha,
	}, "")

	if w.Code != http.StatusCreated {
		t.Errorf("esperava 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp domain.RegistroResponse
	json.NewDecoder(w.Body).Decode(&resp) //nolint:errcheck
	if resp.Tipo != domain.TipoProfissional {
		t.Errorf("esperava PROFISSIONAL, got %s", resp.Tipo)
	}
}

// --- Login ---

func TestLogin_Sucesso_Retorna200ComToken(t *testing.T) {
	r, _ := novoRouter(t)
	doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", map[string]string{
		"email": testEmail,
		"senha": testSenha,
	}, "")

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": testEmail,
		"senha": testSenha,
	}, "")

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp domain.TokenResponse
	json.NewDecoder(w.Body).Decode(&resp) //nolint:errcheck
	if resp.AccessToken == "" {
		t.Error("esperava access_token nao vazio")
	}
}

func TestLogin_CredenciaisInvalidas_Retorna401(t *testing.T) {
	r, _ := novoRouter(t)
	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "naoexiste@diarygo.com.br",
		"senha": testSenha,
	}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

// --- Recuperacao de Senha ---

func TestSolicitarRecuperacaoSenha_Sucesso_Retorna200(t *testing.T) {
	r, _ := novoRouter(t)
	doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", map[string]string{
		"email": testEmail,
		"senha": testSenha,
	}, "")

	w := doRequest(r, http.MethodPost, "/api/v1/auth/solicitar-recuperacao-senha", map[string]string{
		"email": testEmail,
	}, "")

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp domain.SolicitarRecuperacaoResponse
	json.NewDecoder(w.Body).Decode(&resp) //nolint:errcheck
	if resp.Token == "" {
		t.Error("esperava token nao vazio em development")
	}
}

func TestSolicitarRecuperacaoSenha_EmailInexistente_Retorna200(t *testing.T) {
	r, _ := novoRouter(t)

	// OWASP A07: nao revelar que o email nao existe
	w := doRequest(r, http.MethodPost, "/api/v1/auth/solicitar-recuperacao-senha", map[string]string{
		"email": "naoexiste@diarygo.com.br",
	}, "")

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d", w.Code)
	}
}

func TestRedefinirSenha_Sucesso_Retorna200(t *testing.T) {
	r, _ := novoRouter(t)
	doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", map[string]string{
		"email": testEmail,
		"senha": testSenha,
	}, "")

	// Solicitar token
	wRec := doRequest(r, http.MethodPost, "/api/v1/auth/solicitar-recuperacao-senha", map[string]string{
		"email": testEmail,
	}, "")
	var recResp domain.SolicitarRecuperacaoResponse
	json.NewDecoder(wRec.Body).Decode(&recResp) //nolint:errcheck

	// Redefinir
	w := doRequest(r, http.MethodPost, "/api/v1/auth/redefinir-senha", map[string]string{
		"token":      recResp.Token,
		"nova_senha": "NovaSenha123",
	}, "")

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRedefinirSenha_TokenInvalido_Retorna401(t *testing.T) {
	r, _ := novoRouter(t)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/redefinir-senha", map[string]string{
		"token":      "token-invalido",
		"nova_senha": "NovaSenha123",
	}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

// --- Rota protegida ---

func TestRotaProtegida_SemToken_Retorna401(t *testing.T) {
	r, _ := novoRouter(t)
	w := doRequest(r, http.MethodGet, "/api/v1/me", nil, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}

func TestRotaProtegida_TokenValido_Retorna200(t *testing.T) {
	r, _ := novoRouter(t)

	// Registrar e fazer login
	doRequest(r, http.MethodPost, "/api/v1/auth/registro/cliente", map[string]string{
		"email": testEmail,
		"senha": testSenha,
	}, "")
	wLogin := doRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": testEmail,
		"senha": testSenha,
	}, "")

	var tokenResp domain.TokenResponse
	json.NewDecoder(wLogin.Body).Decode(&tokenResp) //nolint:errcheck

	w := doRequest(r, http.MethodGet, "/api/v1/me", nil, tokenResp.AccessToken)
	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRotaProtegida_TokenInvalido_Retorna401(t *testing.T) {
	r, _ := novoRouter(t)
	w := doRequest(r, http.MethodGet, "/api/v1/me", nil, "token.invalido.aqui")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, got %d", w.Code)
	}
}
