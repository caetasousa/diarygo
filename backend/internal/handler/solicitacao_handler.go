package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// SolicitacaoHandler expoe /solicitacoes para clientes autenticados.
type SolicitacaoHandler struct {
	svc *service.SolicitacaoService
}

// NewSolicitacaoHandler cria um novo SolicitacaoHandler.
func NewSolicitacaoHandler(svc *service.SolicitacaoService) *SolicitacaoHandler {
	return &SolicitacaoHandler{svc: svc}
}

// Routes monta sub-router. Deve ser montado dentro de bloco Autenticar + RequererTipo(CLIENTE).
func (h *SolicitacaoHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Criar)
	r.Get("/", h.Listar)
	r.Get("/{id}", h.Buscar)
	r.Delete("/{id}", h.Cancelar)
	return r
}

// Criar godoc
// @Summary      Cria uma nova solicitação de serviço
// @Description  Valida endereço do cliente, regras de antecedência (24h) e congela o orçamento.
// @Tags         solicitacoes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CriarSolicitacaoRequest  true  "Dados da solicitação"
// @Success      201   {object}  domain.SolicitacaoResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Failure      422   {object}  ErroResponse
// @Router       /solicitacoes [post]
func (h *SolicitacaoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16<<10) // 16KB
	var req domain.CriarSolicitacaoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.svc.Criar(r.Context(), payload.UsuarioID, req)
	if err != nil {
		h.mapearErro(w, err)
		return
	}

	RespostaJSON(w, http.StatusCreated, resp)
}

// Listar godoc
// @Summary      Lista solicitações do cliente autenticado
// @Tags         solicitacoes
// @Produce      json
// @Security     BearerAuth
// @Param        status  query  string  false  "Filtro por status"
// @Param        desde   query  string  false  "Data inicial ISO 8601"
// @Param        ate     query  string  false  "Data final ISO 8601"
// @Success      200  {array}   domain.SolicitacaoResponse
// @Failure      401  {object}  ErroResponse
// @Router       /solicitacoes [get]
func (h *SolicitacaoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	filtro := domain.SolicitacaoFiltro{}
	q := r.URL.Query()
	if s := q.Get("status"); s != "" {
		st := domain.StatusSolicitacao(s)
		if st.Valida() {
			filtro.Status = &st
		}
	}
	if d := q.Get("desde"); d != "" {
		if t, err := time.Parse(time.RFC3339, d); err == nil {
			filtro.Desde = &t
		}
	}
	if a := q.Get("ate"); a != "" {
		if t, err := time.Parse(time.RFC3339, a); err == nil {
			filtro.Ate = &t
		}
	}

	resp, err := h.svc.ListarDoCliente(r.Context(), payload.UsuarioID, filtro)
	if err != nil {
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// Buscar godoc
// @Summary      Busca uma solicitação pelo ID
// @Tags         solicitacoes
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "UUID da solicitação"
// @Success      200  {object}  domain.SolicitacaoResponse
// @Failure      401  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /solicitacoes/{id} [get]
func (h *SolicitacaoHandler) Buscar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespostaErro(w, http.StatusNotFound, "solicitacao nao encontrada")
		return
	}

	resp, err := h.svc.BuscarPorID(r.Context(), payload.UsuarioID, id)
	if err != nil {
		h.mapearErro(w, err)
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// Cancelar godoc
// @Summary      Cancela uma solicitação
// @Description  Cancela a solicitação. Se faltar menos de 24h para o serviço, penaliza o score do cliente em 5 pontos.
// @Tags         solicitacoes
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "UUID da solicitação"
// @Success      204
// @Failure      401  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Failure      422  {object}  ErroResponse
// @Router       /solicitacoes/{id} [delete]
func (h *SolicitacaoHandler) Cancelar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		RespostaErro(w, http.StatusNotFound, "solicitacao nao encontrada")
		return
	}

	if err := h.svc.Cancelar(r.Context(), payload.UsuarioID, id); err != nil {
		h.mapearErro(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SolicitacaoHandler) mapearErro(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrSolicitacaoNaoEncontrada):
		RespostaErro(w, http.StatusNotFound, "solicitacao nao encontrada")
	case errors.Is(err, domain.ErrSolicitacaoEnderecoInvalido):
		RespostaErro(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, domain.ErrSolicitacaoAntecedencia):
		RespostaErro(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, domain.ErrSolicitacaoStatusInvalido):
		RespostaErro(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, domain.ErrObservacaoMuitoLonga):
		RespostaErro(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrCategoriaNaoEncontrada),
		errors.Is(err, domain.ErrRegiaoNaoEncontrada),
		errors.Is(err, domain.ErrOpcionalNaoEncontrado),
		errors.Is(err, domain.ErrTabelaPrecoNaoEncontrada):
		RespostaErro(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrFrequenciaInvalida),
		errors.Is(err, domain.ErrComodosInvalidosCalculo),
		errors.Is(err, domain.ErrComodoNegativo),
		errors.Is(err, domain.ErrComodoInvalido):
		RespostaErro(w, http.StatusBadRequest, err.Error())
	default:
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
	}
}
