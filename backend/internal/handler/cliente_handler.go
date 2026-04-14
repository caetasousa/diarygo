package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/caetasousa/diarygo/internal/domain"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
)

// ClienteHandler agrupa os handlers de perfil do cliente.
type ClienteHandler struct {
	svc *service.ClienteService
}

// NewClienteHandler cria um novo ClienteHandler.
func NewClienteHandler(svc *service.ClienteService) *ClienteHandler {
	return &ClienteHandler{svc: svc}
}

// Routes retorna o sub-router chi com as rotas de cliente.
func (h *ClienteHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/me", h.BuscarMe)
	r.Post("/me", h.Criar)
	r.Put("/me", h.Atualizar)
	return r
}

// BuscarMe retorna o perfil do cliente autenticado.
//
// @Summary      Buscar perfil do cliente
// @Description  Retorna o perfil completo do cliente autenticado.
// @Tags         clientes
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.ClienteResponse
// @Failure      401  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /clientes/me [get]
func (h *ClienteHandler) BuscarMe(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	resp, err := h.svc.BuscarPorUsuarioID(r.Context(), payload.UsuarioID)
	if err != nil {
		if errors.Is(err, domain.ErrClienteNaoEncontrado) {
			RespostaErro(w, http.StatusNotFound, "perfil de cliente nao encontrado")
			return
		}
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// Criar cria o perfil do cliente autenticado.
//
// @Summary      Criar perfil do cliente
// @Description  Cria o perfil completo do cliente (nome, CPF, telefone).
// @Tags         clientes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.ClienteRequest  true  "Dados do perfil"
// @Success      201   {object}  domain.ClienteResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Failure      409   {object}  ErroResponse
// @Router       /clientes/me [post]
func (h *ClienteHandler) Criar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.ClienteRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.svc.Criar(r.Context(), payload.UsuarioID, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrClienteJaExiste):
			RespostaErro(w, http.StatusConflict, "perfil de cliente ja existe")
		case errors.Is(err, domain.ErrCPFJaCadastrado):
			RespostaErro(w, http.StatusConflict, "CPF ja cadastrado")
		default:
			RespostaErro(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	RespostaJSON(w, http.StatusCreated, resp)
}

// Atualizar atualiza o perfil do cliente autenticado.
//
// @Summary      Atualizar perfil do cliente
// @Description  Atualiza nome e telefone do cliente (CPF nao pode ser alterado).
// @Tags         clientes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.ClienteRequest  true  "Dados atualizados"
// @Success      200   {object}  domain.ClienteResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Failure      404   {object}  ErroResponse
// @Router       /clientes/me [put]
func (h *ClienteHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.ClienteRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.svc.Atualizar(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrClienteNaoEncontrado) {
			RespostaErro(w, http.StatusNotFound, "perfil de cliente nao encontrado")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}
