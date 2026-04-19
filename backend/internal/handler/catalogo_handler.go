package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CatalogoHandler expoe /categorias, /categorias/:id/opcionais e /precos/calcular.
type CatalogoHandler struct {
	svc *service.PrecificacaoService
}

// NewCatalogoHandler cria um novo CatalogoHandler.
func NewCatalogoHandler(svc *service.PrecificacaoService) *CatalogoHandler {
	return &CatalogoHandler{svc: svc}
}

// RoutesCategorias monta sub-router para /categorias (publico).
func (h *CatalogoHandler) RoutesCategorias() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListarCategorias)
	r.Get("/{categoriaID}/opcionais", h.ListarOpcionais)
	return r
}

// RoutesPrecos monta sub-router para /precos (publico).
func (h *CatalogoHandler) RoutesPrecos() chi.Router {
	r := chi.NewRouter()
	r.Post("/calcular", h.Calcular)
	return r
}

// ListarCategorias godoc
// @Summary      Lista categorias de servico ativas
// @Tags         catalogo
// @Produce      json
// @Success      200  {array}  domain.CategoriaResponse
// @Router       /categorias [get]
func (h *CatalogoHandler) ListarCategorias(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.ListarCategorias(r.Context())
	if err != nil {
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}
	RespostaJSON(w, http.StatusOK, resp)
}

// ListarOpcionais godoc
// @Summary      Lista opcionais disponiveis para uma categoria
// @Tags         catalogo
// @Produce      json
// @Param        categoriaID  path  string  true  "UUID da categoria"
// @Success      200  {array}  domain.OpcionalResponse
// @Failure      404  {object}  ErroResponse
// @Router       /categorias/{categoriaID}/opcionais [get]
func (h *CatalogoHandler) ListarOpcionais(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "categoriaID"))
	if err != nil {
		RespostaErro(w, http.StatusBadRequest, "id de categoria invalido")
		return
	}
	resp, err := h.svc.ListarOpcionaisDaCategoria(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCategoriaNaoEncontrada) {
			RespostaErro(w, http.StatusNotFound, "categoria nao encontrada")
			return
		}
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}
	RespostaJSON(w, http.StatusOK, resp)
}

// Calcular godoc
// @Summary      Calcula orcamento com breakdown (diferencial C)
// @Description  Publico. Retorna valor total, duracao e itens[] do breakdown para exibicao.
// @Tags         catalogo
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CalculoPrecoRequest  true  "Parametros do calculo"
// @Success      200   {object}  domain.CalculoPrecoResponse
// @Failure      400   {object}  ErroResponse
// @Failure      404   {object}  ErroResponse
// @Router       /precos/calcular [post]
func (h *CatalogoHandler) Calcular(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.CalculoPrecoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}
	resp, err := h.svc.CalcularValorReferencia(r.Context(), req)
	if err != nil {
		h.mapearErroCalculo(w, err)
		return
	}
	RespostaJSON(w, http.StatusOK, resp)
}

func (h *CatalogoHandler) mapearErroCalculo(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCategoriaNaoEncontrada):
		RespostaErro(w, http.StatusNotFound, "categoria nao encontrada")
	case errors.Is(err, domain.ErrRegiaoNaoEncontrada):
		RespostaErro(w, http.StatusNotFound, "regiao nao encontrada")
	case errors.Is(err, domain.ErrOpcionalNaoEncontrado):
		RespostaErro(w, http.StatusNotFound, "opcional nao encontrado")
	case errors.Is(err, domain.ErrTabelaPrecoNaoEncontrada):
		RespostaErro(w, http.StatusNotFound, "tabela de precos nao encontrada para a regiao")
	case errors.Is(err, domain.ErrFrequenciaInvalida),
		errors.Is(err, domain.ErrDataServicoPassado),
		errors.Is(err, domain.ErrComodosInvalidosCalculo),
		errors.Is(err, domain.ErrComodoNegativo),
		errors.Is(err, domain.ErrComodoInvalido):
		RespostaErro(w, http.StatusBadRequest, err.Error())
	default:
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
	}
}
