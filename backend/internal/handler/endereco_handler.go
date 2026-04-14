package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/caetasousa/diarygo/internal/domain"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// EnderecoHandler agrupa os handlers de enderecos do cliente.
type EnderecoHandler struct {
	svc *service.EnderecoService
}

// NewEnderecoHandler cria um novo EnderecoHandler.
func NewEnderecoHandler(svc *service.EnderecoService) *EnderecoHandler {
	return &EnderecoHandler{svc: svc}
}

// Routes retorna o sub-router chi com as rotas de endereco.
func (h *EnderecoHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.Listar)
	r.Post("/", h.Criar)
	r.Put("/{id}", h.Atualizar)
	r.Delete("/{id}", h.Remover)
	r.Put("/{id}/principal", h.DefinirPrincipal)
	return r
}

// Listar retorna todos os enderecos do cliente autenticado.
//
// @Summary      Listar enderecos
// @Description  Retorna todos os enderecos do cliente autenticado.
// @Tags         enderecos
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.EnderecoResponse
// @Failure      401  {object}  ErroResponse
// @Router       /clientes/me/enderecos [get]
func (h *EnderecoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	resp, err := h.svc.Listar(r.Context(), payload.UsuarioID)
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

// Criar cria um novo endereco para o cliente autenticado.
//
// @Summary      Criar endereco
// @Description  Cria um novo endereco para o cliente autenticado.
// @Tags         enderecos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.EnderecoRequest  true  "Dados do endereco"
// @Success      201   {object}  domain.EnderecoResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Router       /clientes/me/enderecos [post]
func (h *EnderecoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.EnderecoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.svc.Criar(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrClienteNaoEncontrado) {
			RespostaErro(w, http.StatusNotFound, "perfil de cliente nao encontrado")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusCreated, resp)
}

// Atualizar atualiza um endereco do cliente autenticado.
//
// @Summary      Atualizar endereco
// @Description  Atualiza um endereco existente do cliente autenticado.
// @Tags         enderecos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                  true  "ID do endereco"
// @Param        body  body      domain.EnderecoRequest  true  "Dados atualizados"
// @Success      200   {object}  domain.EnderecoResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Failure      403   {object}  ErroResponse
// @Failure      404   {object}  ErroResponse
// @Router       /clientes/me/enderecos/{id} [put]
func (h *EnderecoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	enderecoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespostaErro(w, http.StatusBadRequest, "ID de endereco invalido")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.EnderecoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.svc.Atualizar(r.Context(), payload.UsuarioID, enderecoID, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEnderecoNaoEncontrado):
			RespostaErro(w, http.StatusNotFound, "endereco nao encontrado")
		case errors.Is(err, domain.ErrEnderecoNaoPertenceAoCliente):
			RespostaErro(w, http.StatusForbidden, "acesso negado")
		default:
			RespostaErro(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// Remover remove um endereco do cliente autenticado.
//
// @Summary      Remover endereco
// @Description  Remove um endereco existente do cliente autenticado.
// @Tags         enderecos
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID do endereco"
// @Success      204
// @Failure      401  {object}  ErroResponse
// @Failure      403  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /clientes/me/enderecos/{id} [delete]
func (h *EnderecoHandler) Remover(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	enderecoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespostaErro(w, http.StatusBadRequest, "ID de endereco invalido")
		return
	}

	if err := h.svc.Remover(r.Context(), payload.UsuarioID, enderecoID); err != nil {
		switch {
		case errors.Is(err, domain.ErrEnderecoNaoEncontrado):
			RespostaErro(w, http.StatusNotFound, "endereco nao encontrado")
		case errors.Is(err, domain.ErrEnderecoNaoPertenceAoCliente):
			RespostaErro(w, http.StatusForbidden, "acesso negado")
		default:
			RespostaErro(w, http.StatusInternalServerError, "erro interno")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DefinirPrincipal marca um endereco como principal.
//
// @Summary      Definir endereco principal
// @Description  Marca o endereco informado como principal, desmarcando os outros.
// @Tags         enderecos
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID do endereco"
// @Success      200  {object}  domain.EnderecoResponse
// @Failure      401  {object}  ErroResponse
// @Failure      403  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /clientes/me/enderecos/{id}/principal [put]
func (h *EnderecoHandler) DefinirPrincipal(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	enderecoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespostaErro(w, http.StatusBadRequest, "ID de endereco invalido")
		return
	}

	resp, err := h.svc.DefinirPrincipal(r.Context(), payload.UsuarioID, enderecoID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEnderecoNaoEncontrado):
			RespostaErro(w, http.StatusNotFound, "endereco nao encontrado")
		case errors.Is(err, domain.ErrEnderecoNaoPertenceAoCliente):
			RespostaErro(w, http.StatusForbidden, "acesso negado")
		default:
			RespostaErro(w, http.StatusInternalServerError, "erro interno")
		}
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}
