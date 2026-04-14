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

// ProfissionalHandler agrupa os handlers de perfil da profissional.
type ProfissionalHandler struct {
	profSvc *service.ProfissionalService
	credSvc *service.CredenciamentoService
}

// NewProfissionalHandler cria um novo ProfissionalHandler.
func NewProfissionalHandler(profSvc *service.ProfissionalService, credSvc *service.CredenciamentoService) *ProfissionalHandler {
	return &ProfissionalHandler{profSvc: profSvc, credSvc: credSvc}
}

// Routes retorna o sub-router chi com as rotas de profissional.
func (h *ProfissionalHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/me", h.BuscarMe)
	r.Post("/me", h.Criar)
	r.Put("/me", h.Atualizar)

	// Documentos
	r.Post("/me/documentos", h.EnviarDocumento)
	r.Get("/me/documentos", h.ListarDocumentos)

	// Referencias
	r.Post("/me/referencias", h.AdicionarReferencia)
	r.Get("/me/referencias", h.ListarReferencias)

	// Regioes
	r.Put("/me/regioes", h.DefinirRegioes)
	r.Get("/me/regioes", h.ListarRegioes)

	// Disponibilidade
	r.Put("/me/disponibilidades", h.DefinirDisponibilidades)
	r.Get("/me/disponibilidades", h.ListarDisponibilidades)

	return r
}

// BuscarMe retorna o perfil da profissional autenticada.
//
// @Summary      Buscar perfil da profissional
// @Description  Retorna o perfil completo da profissional autenticada.
// @Tags         profissionais
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.ProfissionalResponse
// @Failure      401  {object}  ErroResponse
// @Failure      404  {object}  ErroResponse
// @Router       /profissionais/me [get]
func (h *ProfissionalHandler) BuscarMe(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	resp, err := h.profSvc.BuscarPorUsuarioID(r.Context(), payload.UsuarioID)
	if err != nil {
		if errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
			RespostaErro(w, http.StatusNotFound, "perfil de profissional nao encontrado")
			return
		}
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// Criar cria o perfil da profissional autenticada.
//
// @Summary      Criar perfil da profissional
// @Description  Cria o perfil da profissional (nome, CPF, RG, telefone, foto, MEI). Status inicial: PENDENTE.
// @Tags         profissionais
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.ProfissionalRequest  true  "Dados do perfil"
// @Success      201   {object}  domain.ProfissionalResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Failure      409   {object}  ErroResponse
// @Router       /profissionais/me [post]
func (h *ProfissionalHandler) Criar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.ProfissionalRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.profSvc.Criar(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrProfissionalJaExiste) {
			RespostaErro(w, http.StatusConflict, "perfil de profissional ja existe")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusCreated, resp)
}

// Atualizar atualiza o perfil da profissional autenticada.
//
// @Summary      Atualizar perfil da profissional
// @Description  Atualiza nome, telefone, foto_url e MEI da profissional.
// @Tags         profissionais
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.ProfissionalRequest  true  "Dados atualizados"
// @Success      200   {object}  domain.ProfissionalResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Failure      404   {object}  ErroResponse
// @Router       /profissionais/me [put]
func (h *ProfissionalHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.ProfissionalRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.profSvc.Atualizar(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
			RespostaErro(w, http.StatusNotFound, "perfil de profissional nao encontrado")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// EnviarDocumento envia um documento da profissional autenticada.
//
// @Summary      Enviar documento
// @Description  Registra a URL de um documento da profissional. Status inicial: PENDENTE.
// @Tags         profissionais
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.DocumentoRequest  true  "Tipo e URL do documento"
// @Success      201   {object}  domain.DocumentoResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Router       /profissionais/me/documentos [post]
func (h *ProfissionalHandler) EnviarDocumento(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.DocumentoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.credSvc.EnviarDocumento(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
			RespostaErro(w, http.StatusNotFound, "perfil de profissional nao encontrado")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusCreated, resp)
}

// ListarDocumentos lista os documentos da profissional autenticada.
//
// @Summary      Listar documentos
// @Description  Retorna todos os documentos enviados pela profissional autenticada.
// @Tags         profissionais
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.DocumentoResponse
// @Failure      401  {object}  ErroResponse
// @Router       /profissionais/me/documentos [get]
func (h *ProfissionalHandler) ListarDocumentos(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	resp, err := h.credSvc.ListarDocumentos(r.Context(), payload.UsuarioID)
	if err != nil {
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// AdicionarReferencia adiciona uma referencia profissional.
//
// @Summary      Adicionar referencia
// @Description  Adiciona uma referencia de trabalho anterior da profissional autenticada.
// @Tags         profissionais
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.ReferenciaRequest  true  "Nome e telefone da referencia"
// @Success      201   {object}  domain.ReferenciaResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Router       /profissionais/me/referencias [post]
func (h *ProfissionalHandler) AdicionarReferencia(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.ReferenciaRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.credSvc.AdicionarReferencia(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
			RespostaErro(w, http.StatusNotFound, "perfil de profissional nao encontrado")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusCreated, resp)
}

// ListarReferencias lista as referencias da profissional autenticada.
//
// @Summary      Listar referencias
// @Description  Retorna todas as referencias da profissional autenticada.
// @Tags         profissionais
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.ReferenciaResponse
// @Failure      401  {object}  ErroResponse
// @Router       /profissionais/me/referencias [get]
func (h *ProfissionalHandler) ListarReferencias(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	resp, err := h.credSvc.ListarReferencias(r.Context(), payload.UsuarioID)
	if err != nil {
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// DefinirRegioes define as regioes de atuacao da profissional autenticada.
//
// @Summary      Definir regioes de atuacao
// @Description  Substitui as regioes de atuacao da profissional autenticada.
// @Tags         profissionais
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.DefinirRegioesRequest  true  "Lista de IDs de regioes"
// @Success      200   {array}   domain.RegiaoResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Router       /profissionais/me/regioes [put]
func (h *ProfissionalHandler) DefinirRegioes(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.DefinirRegioesRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.credSvc.DefinirRegioes(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
			RespostaErro(w, http.StatusNotFound, "perfil de profissional nao encontrado")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// ListarRegioes lista as regioes de atuacao da profissional autenticada.
//
// @Summary      Listar regioes de atuacao
// @Description  Retorna as regioes de atuacao da profissional autenticada.
// @Tags         profissionais
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.RegiaoResponse
// @Failure      401  {object}  ErroResponse
// @Router       /profissionais/me/regioes [get]
func (h *ProfissionalHandler) ListarRegioes(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	resp, err := h.credSvc.ListarRegioesAtuacao(r.Context(), payload.UsuarioID)
	if err != nil {
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// DefinirDisponibilidades define a disponibilidade semanal da profissional autenticada.
//
// @Summary      Definir disponibilidade
// @Description  Substitui a disponibilidade semanal da profissional autenticada.
// @Tags         profissionais
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.DefinirDisponibilidadesRequest  true  "Slots de disponibilidade"
// @Success      200   {array}   domain.DisponibilidadeResponse
// @Failure      400   {object}  ErroResponse
// @Failure      401   {object}  ErroResponse
// @Router       /profissionais/me/disponibilidades [put]
func (h *ProfissionalHandler) DefinirDisponibilidades(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req domain.DefinirDisponibilidadesRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "body invalido")
		return
	}

	resp, err := h.credSvc.DefinirDisponibilidades(r.Context(), payload.UsuarioID, req)
	if err != nil {
		if errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
			RespostaErro(w, http.StatusNotFound, "perfil de profissional nao encontrado")
			return
		}
		RespostaErro(w, http.StatusBadRequest, err.Error())
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}

// ListarDisponibilidades lista a disponibilidade semanal da profissional autenticada.
//
// @Summary      Listar disponibilidade
// @Description  Retorna a disponibilidade semanal da profissional autenticada.
// @Tags         profissionais
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.DisponibilidadeResponse
// @Failure      401  {object}  ErroResponse
// @Router       /profissionais/me/disponibilidades [get]
func (h *ProfissionalHandler) ListarDisponibilidades(w http.ResponseWriter, r *http.Request) {
	payload, ok := mw.UsuarioDoContexto(r.Context())
	if !ok {
		RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
		return
	}

	resp, err := h.credSvc.ListarDisponibilidades(r.Context(), payload.UsuarioID)
	if err != nil {
		RespostaErro(w, http.StatusInternalServerError, "erro interno")
		return
	}

	RespostaJSON(w, http.StatusOK, resp)
}
