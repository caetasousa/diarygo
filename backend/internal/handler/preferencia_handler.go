package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/caetasousa/diarygo/internal/domain"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// PreferenciaHandler agrupa os endpoints de favoritas/bloqueios do cliente.
type PreferenciaHandler struct {
	svc *service.PreferenciaService
}

// NewPreferenciaHandler cria um novo PreferenciaHandler.
func NewPreferenciaHandler(svc *service.PreferenciaService) *PreferenciaHandler {
	return &PreferenciaHandler{svc: svc}
}

// RoutesFavoritas retorna as rotas de favoritas sob /clientes/me/favoritas.
func (h *PreferenciaHandler) RoutesFavoritas() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListarFavoritas)
	r.Post("/{profissionalID}", h.Favoritar)
	r.Delete("/{profissionalID}", h.Remover)
	return r
}

// RoutesBloqueios retorna as rotas de bloqueios sob /clientes/me/bloqueios.
func (h *PreferenciaHandler) RoutesBloqueios() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListarBloqueios)
	r.Post("/{profissionalID}", h.Bloquear)
	r.Delete("/{profissionalID}", h.Remover)
	return r
}

// ListarFavoritas godoc
// @Summary      Lista profissionais favoritas do cliente
// @Tags         preferencias
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  domain.PreferenciaResponse
// @Failure      401  {object}  ErroResponse
// @Router       /clientes/me/favoritas [get]
func (h *PreferenciaHandler) ListarFavoritas(w http.ResponseWriter, r *http.Request) {
	h.listar(w, r, domain.PreferenciaFavorita)
}

// ListarBloqueios godoc
// @Summary      Lista profissionais bloqueadas pelo cliente
// @Tags         preferencias
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  domain.PreferenciaResponse
// @Failure      401  {object}  ErroResponse
// @Router       /clientes/me/bloqueios [get]
func (h *PreferenciaHandler) ListarBloqueios(w http.ResponseWriter, r *http.Request) {
	h.listar(w, r, domain.PreferenciaBloqueada)
}

// Favoritar godoc
// @Summary      Marca profissional como favorita
// @Tags         preferencias
// @Produce      json
// @Security     BearerAuth
// @Param        profissionalID  path  string  true  "UUID da profissional"
// @Success      204
// @Failure      401  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /clientes/me/favoritas/{profissionalID} [post]
func (h *PreferenciaHandler) Favoritar(w http.ResponseWriter, r *http.Request) {
	h.registrar(w, r, h.svc.Favoritar)
}

// Bloquear godoc
// @Summary      Bloqueia profissional para o cliente
// @Tags         preferencias
// @Produce      json
// @Security     BearerAuth
// @Param        profissionalID  path  string  true  "UUID da profissional"
// @Success      204
// @Failure      401  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /clientes/me/bloqueios/{profissionalID} [post]
func (h *PreferenciaHandler) Bloquear(w http.ResponseWriter, r *http.Request) {
	h.registrar(w, r, h.svc.Bloquear)
}

// Remover godoc
// @Summary      Remove preferencia (favorita ou bloqueio) de uma profissional
// @Tags         preferencias
// @Produce      json
// @Security     BearerAuth
// @Param        profissionalID  path  string  true  "UUID da profissional"
// @Success      204
// @Failure      401  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /clientes/me/favoritas/{profissionalID} [delete]
func (h *PreferenciaHandler) Remover(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}
	profissionalID, err := uuid.Parse(chi.URLParam(r, "profissionalID"))
	if err != nil {
		RespostaErro(w, http.StatusBadRequest, "id de profissional invalido")
		return
	}
	if err := h.svc.Remover(r.Context(), payload.UsuarioID, profissionalID); err != nil {
		h.mapearErro(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PreferenciaHandler) listar(w http.ResponseWriter, r *http.Request, tipo domain.TipoPreferencia) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}
	resp, err := h.svc.Listar(r.Context(), payload.UsuarioID, tipo)
	if err != nil {
		h.mapearErro(w, err)
		return
	}
	RespostaJSON(w, http.StatusOK, resp)
}

func (h *PreferenciaHandler) registrar(w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, usuarioID, profissionalID uuid.UUID) error) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}
	profissionalID, err := uuid.Parse(chi.URLParam(r, "profissionalID"))
	if err != nil {
		RespostaErro(w, http.StatusBadRequest, "id de profissional invalido")
		return
	}
	if err := fn(r.Context(), payload.UsuarioID, profissionalID); err != nil {
		h.mapearErro(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PreferenciaHandler) mapearErro(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrClienteNaoEncontrado):
		RespostaErro(w, http.StatusNotFound, "perfil de cliente nao encontrado")
	case errors.Is(err, domain.ErrProfissionalNaoEncontrada):
		RespostaErro(w, http.StatusNotFound, "profissional nao encontrada")
	case errors.Is(err, domain.ErrPreferenciaNaoEncontrada):
		RespostaErro(w, http.StatusNotFound, "preferencia nao encontrada")
	case errors.Is(err, domain.ErrTipoPreferenciaInvalido):
		RespostaErro(w, http.StatusBadRequest, "tipo de preferencia invalido")
	default:
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
	}
}
