package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
)

// AuthHandler agrupa os handlers de autenticacao.
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler cria um novo AuthHandler.
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Routes retorna o sub-router chi com todas as rotas de autenticacao.
func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/registro/cliente", h.RegistrarCliente)
	r.Post("/registro/profissional", h.RegistrarProfissional)
	r.Post("/login", h.Login)
	r.Post("/verificar-email", h.VerificarEmail)
	return r
}

// RegistrarCliente cria um novo usuario do tipo CLIENTE.
//
// @Summary      Registrar novo cliente
// @Description  Cria um usuario com tipo CLIENTE. Retorna codigo de verificacao em development.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.RegistroRequest  true  "Email e senha"
// @Success      201   {object}  domain.RegistroResponse
// @Failure      400   {object}  ErroResponse
// @Failure      409   {object}  ErroResponse
// @Failure      500   {object}  ErroResponse
// @Router       /auth/registro/cliente [post]
func (h *AuthHandler) RegistrarCliente(w http.ResponseWriter, r *http.Request) {
	h.registrar(w, r, domain.TipoCliente)
}

// RegistrarProfissional cria um novo usuario do tipo PROFISSIONAL.
//
// @Summary      Registrar nova profissional
// @Description  Cria um usuario com tipo PROFISSIONAL com status PENDENTE (aguarda aprovacao).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.RegistroRequest  true  "Email e senha"
// @Success      201   {object}  domain.RegistroResponse
// @Failure      400   {object}  ErroResponse
// @Failure      409   {object}  ErroResponse
// @Failure      500   {object}  ErroResponse
// @Router       /auth/registro/profissional [post]
func (h *AuthHandler) RegistrarProfissional(w http.ResponseWriter, r *http.Request) {
	h.registrar(w, r, domain.TipoProfissional)
}

// registrar e o metodo interno compartilhado pelos handlers de registro.
func (h *AuthHandler) registrar(w http.ResponseWriter, r *http.Request, tipo domain.TipoUsuario) {
	// OWASP A06: limitar tamanho do body a 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req domain.RegistroRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido: "+err.Error())
		return
	}

	resp, err := h.svc.Registrar(r.Context(), req, tipo)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailJaExiste):
			RespostaErro(w, http.StatusConflict, "email ja cadastrado")
		case isValidacaoErro(err):
			RespostaErro(w, http.StatusBadRequest, err.Error())
		default:
			RespostaErro(w, http.StatusInternalServerError, "erro interno")
		}
		return
	}

	RespostaJSON(w, http.StatusCreated, resp)
}

// Login autentica um usuario e retorna um JWT.
//
// @Summary      Fazer login
// @Description  Autentica com email e senha. Retorna access token JWT.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.LoginRequest  true  "Credenciais"
// @Success      200   {object}  domain.TokenResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Failure      403   {object}  ErroResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCredenciaisInvalidas):
			// OWASP A07: mensagem generica — nunca revelar "email nao existe" vs "senha errada"
			RespostaErro(w, http.StatusUnauthorized, "credenciais invalidas")
		case errors.Is(err, domain.ErrUsuarioInativo):
			RespostaErro(w, http.StatusForbidden, "usuario inativo")
		default:
			RespostaErro(w, http.StatusInternalServerError, "erro interno")
		}
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// VerificarEmail ativa o email do usuario com o codigo recebido.
//
// @Summary      Verificar email
// @Description  Ativa o email do usuario usando o codigo enviado no registro (em development: retornado na resposta de registro).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.VerificarEmailRequest  true  "Email e codigo de verificacao"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  ErroResponse
// @Failure      404   {object}  ErroResponse
// @Router       /auth/verificar-email [post]
func (h *AuthHandler) VerificarEmail(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req domain.VerificarEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	if err := h.svc.VerificarEmail(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, domain.ErrUsuarioNaoEncontrado):
			RespostaErro(w, http.StatusNotFound, "usuario nao encontrado")
		case errors.Is(err, domain.ErrCodigoVerificacaoInvalido):
			RespostaErro(w, http.StatusBadRequest, "codigo de verificacao invalido")
		default:
			RespostaErro(w, http.StatusInternalServerError, "erro interno")
		}
		return
	}

	RespostaJSON(w, http.StatusOK, map[string]string{"mensagem": "email verificado com sucesso"})
}

// isValidacaoErro verifica se o erro e de validacao de dominio (nao e um sentinel conhecido).
func isValidacaoErro(err error) bool {
	sentinels := []error{
		domain.ErrEmailJaExiste,
		domain.ErrCredenciaisInvalidas,
		domain.ErrUsuarioNaoEncontrado,
		domain.ErrUsuarioInativo,
		domain.ErrEmailNaoVerificado,
		domain.ErrTokenInvalido,
		domain.ErrCodigoVerificacaoInvalido,
	}
	for _, s := range sentinels {
		if errors.Is(err, s) {
			return false
		}
	}
	return true
}
