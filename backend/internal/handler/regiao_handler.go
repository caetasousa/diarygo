package handler

import (
	"net/http"

	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
)

// RegiaoHandler agrupa os handlers de regioes (rota publica).
type RegiaoHandler struct {
	credSvc *service.CredenciamentoService
}

// NewRegiaoHandler cria um novo RegiaoHandler.
func NewRegiaoHandler(credSvc *service.CredenciamentoService) *RegiaoHandler {
	return &RegiaoHandler{credSvc: credSvc}
}

// Routes retorna o sub-router chi com as rotas de regioes.
func (h *RegiaoHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListarAtivas)
	return r
}

// ListarAtivas retorna todas as regioes ativas do sistema.
//
// @Summary      Listar regioes ativas
// @Description  Retorna todas as regioes de atuacao ativas disponíveis no sistema (rota publica).
// @Tags         regioes
// @Produce      json
// @Success      200  {array}   domain.RegiaoResponse
// @Failure      500  {object}  ErroResponse
// @Router       /regioes [get]
func (h *RegiaoHandler) ListarAtivas(w http.ResponseWriter, r *http.Request) {
	resp, err := h.credSvc.ListarRegioes(r.Context())
	if err != nil {
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}
