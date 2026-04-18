package service

import (
	"context"
	"strings"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// CredenciamentoService orquestra documentos, referencias, regioes e disponibilidade.
type CredenciamentoService struct {
	profissionalRepo    domain.ProfissionalRepository
	documentoRepo       domain.DocumentoRepository
	referenciaRepo      domain.ReferenciaRepository
	regiaoRepo          domain.RegiaoRepository
	profRegiaoRepo      domain.ProfissionalRegiaoRepository
	disponibilidadeRepo domain.DisponibilidadeRepository
}

// NewCredenciamentoService cria um novo CredenciamentoService.
func NewCredenciamentoService(
	profissionalRepo domain.ProfissionalRepository,
	documentoRepo domain.DocumentoRepository,
	referenciaRepo domain.ReferenciaRepository,
	regiaoRepo domain.RegiaoRepository,
	profRegiaoRepo domain.ProfissionalRegiaoRepository,
	disponibilidadeRepo domain.DisponibilidadeRepository,
) *CredenciamentoService {
	return &CredenciamentoService{
		profissionalRepo:    profissionalRepo,
		documentoRepo:       documentoRepo,
		referenciaRepo:      referenciaRepo,
		regiaoRepo:          regiaoRepo,
		profRegiaoRepo:      profRegiaoRepo,
		disponibilidadeRepo: disponibilidadeRepo,
	}
}

// --- Documentos ---

// EnviarDocumento registra um novo documento da profissional.
func (s *CredenciamentoService) EnviarDocumento(ctx context.Context, usuarioID uuid.UUID, req domain.DocumentoRequest) (*domain.DocumentoResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	if !domain.TipoDocumentoValido(req.Tipo) {
		return nil, domain.ErrTipoDocumentoInvalido
	}
	if strings.TrimSpace(req.URL) == "" {
		return nil, domain.ErrURLDocumentoObrigatoria
	}

	d := &domain.Documento{
		ID:             uuid.New(),
		ProfissionalID: prof.ID,
		Tipo:           req.Tipo,
		URL:            strings.TrimSpace(req.URL),
		Status:         domain.DocPendente,
	}

	if err := s.documentoRepo.Criar(ctx, d); err != nil {
		return nil, err
	}

	return toDocumentoResponse(d), nil
}

// ListarDocumentos retorna todos os documentos da profissional autenticada.
func (s *CredenciamentoService) ListarDocumentos(ctx context.Context, usuarioID uuid.UUID) ([]*domain.DocumentoResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	docs, err := s.documentoRepo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.DocumentoResponse, len(docs))
	for i, d := range docs {
		resp[i] = toDocumentoResponse(d)
	}
	return resp, nil
}

// --- Referencias ---

// AdicionarReferencia registra uma nova referencia da profissional.
func (s *CredenciamentoService) AdicionarReferencia(ctx context.Context, usuarioID uuid.UUID, req domain.ReferenciaRequest) (*domain.ReferenciaResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.NomeContato) == "" {
		return nil, domain.ErrNomeContatoObrigatorio
	}
	telefone := limparMascara(req.TelefoneContato)
	if err := domain.ValidarTelefone(telefone); err != nil {
		return nil, domain.ErrTelefoneContatoObrigatorio
	}

	r := &domain.Referencia{
		ID:              uuid.New(),
		ProfissionalID:  prof.ID,
		NomeContato:     strings.TrimSpace(req.NomeContato),
		TelefoneContato: telefone,
		Status:          domain.RefPendente,
	}

	if err := s.referenciaRepo.Criar(ctx, r); err != nil {
		return nil, err
	}

	return toReferenciaResponse(r), nil
}

// ListarReferencias retorna todas as referencias da profissional autenticada.
func (s *CredenciamentoService) ListarReferencias(ctx context.Context, usuarioID uuid.UUID) ([]*domain.ReferenciaResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	refs, err := s.referenciaRepo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.ReferenciaResponse, len(refs))
	for i, r := range refs {
		resp[i] = toReferenciaResponse(r)
	}
	return resp, nil
}

// --- Regioes ---

// ListarRegioes retorna todas as regioes ativas do sistema (rota publica).
func (s *CredenciamentoService) ListarRegioes(ctx context.Context) ([]*domain.RegiaoResponse, error) {
	regioes, err := s.regiaoRepo.ListarAtivas(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.RegiaoResponse, len(regioes))
	for i, r := range regioes {
		resp[i] = toRegiaoResponse(r)
	}
	return resp, nil
}

// DefinirRegioes define as regioes de atuacao da profissional autenticada.
// Deduplica IDs do payload antes de enviar ao repo — a PK composta
// (profissional_id, regiao_id) ja impede duplicatas no banco, mas o service
// responde com erro mais claro se o cliente enviar IDs repetidos.
func (s *CredenciamentoService) DefinirRegioes(ctx context.Context, usuarioID uuid.UUID, req domain.DefinirRegioesRequest) ([]*domain.RegiaoResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	visto := make(map[uuid.UUID]struct{}, len(req.RegiaoIDs))
	ids := make([]uuid.UUID, 0, len(req.RegiaoIDs))
	for _, id := range req.RegiaoIDs {
		if id == uuid.Nil {
			return nil, domain.ErrRegiaoIDInvalida
		}
		if _, ok := visto[id]; ok {
			continue
		}
		visto[id] = struct{}{}
		ids = append(ids, id)
	}

	if err := s.profRegiaoRepo.DefinirRegioes(ctx, prof.ID, ids); err != nil {
		return nil, err
	}

	regioes, err := s.profRegiaoRepo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.RegiaoResponse, len(regioes))
	for i, r := range regioes {
		resp[i] = toRegiaoResponse(r)
	}
	return resp, nil
}

// ListarRegioesAtuacao retorna as regioes de atuacao da profissional autenticada.
func (s *CredenciamentoService) ListarRegioesAtuacao(ctx context.Context, usuarioID uuid.UUID) ([]*domain.RegiaoResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	regioes, err := s.profRegiaoRepo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.RegiaoResponse, len(regioes))
	for i, r := range regioes {
		resp[i] = toRegiaoResponse(r)
	}
	return resp, nil
}

// --- Disponibilidade ---

// DefinirDisponibilidades substitui a disponibilidade semanal da profissional.
func (s *CredenciamentoService) DefinirDisponibilidades(ctx context.Context, usuarioID uuid.UUID, req domain.DefinirDisponibilidadesRequest) ([]*domain.DisponibilidadeResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	// Chave de deduplicacao espelha o UNIQUE uq_disponibilidades_prof_dia_inicio_fim.
	// Detectar duplicata no payload aqui permite erro claro antes do SQL disparar.
	type slotKey struct {
		dia int
		ini string
		fim string
	}
	visto := make(map[slotKey]struct{}, len(req.Slots))
	slots := make([]*domain.Disponibilidade, 0, len(req.Slots))
	for _, s2 := range req.Slots {
		if s2.DiaSemana < 0 || s2.DiaSemana > 6 {
			return nil, domain.ErrDiaSemanaInvalido
		}
		if err := domain.ValidarHora(s2.HoraInicio); err != nil {
			return nil, domain.ErrHoraInicioInvalida
		}
		if err := domain.ValidarHora(s2.HoraFim); err != nil {
			return nil, domain.ErrHoraFimInvalida
		}
		if s2.HoraFim <= s2.HoraInicio {
			return nil, domain.ErrHoraFimAntesDaInicio
		}
		k := slotKey{dia: s2.DiaSemana, ini: s2.HoraInicio, fim: s2.HoraFim}
		if _, ok := visto[k]; ok {
			return nil, domain.ErrSlotDuplicado
		}
		visto[k] = struct{}{}
		slots = append(slots, &domain.Disponibilidade{
			ID:         uuid.New(),
			DiaSemana:  s2.DiaSemana,
			HoraInicio: s2.HoraInicio,
			HoraFim:    s2.HoraFim,
		})
	}

	if err := s.disponibilidadeRepo.DefinirDisponibilidades(ctx, prof.ID, slots); err != nil {
		return nil, err
	}

	resultado, err := s.disponibilidadeRepo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.DisponibilidadeResponse, len(resultado))
	for i, d := range resultado {
		resp[i] = toDisponibilidadeResponse(d)
	}
	return resp, nil
}

// ListarDisponibilidades retorna a disponibilidade semanal da profissional autenticada.
func (s *CredenciamentoService) ListarDisponibilidades(ctx context.Context, usuarioID uuid.UUID) ([]*domain.DisponibilidadeResponse, error) {
	prof, err := s.profissionalRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	slots, err := s.disponibilidadeRepo.ListarPorProfissionalID(ctx, prof.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.DisponibilidadeResponse, len(slots))
	for i, d := range slots {
		resp[i] = toDisponibilidadeResponse(d)
	}
	return resp, nil
}

// --- helpers de conversao ---

func toDocumentoResponse(d *domain.Documento) *domain.DocumentoResponse {
	return &domain.DocumentoResponse{
		ID:     d.ID,
		Tipo:   d.Tipo,
		URL:    d.URL,
		Status: d.Status,
	}
}

func toReferenciaResponse(r *domain.Referencia) *domain.ReferenciaResponse {
	return &domain.ReferenciaResponse{
		ID:              r.ID,
		NomeContato:     r.NomeContato,
		TelefoneContato: r.TelefoneContato,
		Status:          r.Status,
	}
}

func toRegiaoResponse(r *domain.Regiao) *domain.RegiaoResponse {
	return &domain.RegiaoResponse{
		ID:        r.ID,
		Nome:      r.Nome,
		Cidade:    r.Cidade,
		Estado:    r.Estado,
		CEPInicio: r.CEPInicio,
		CEPFim:    r.CEPFim,
	}
}

func toDisponibilidadeResponse(d *domain.Disponibilidade) *domain.DisponibilidadeResponse {
	return &domain.DisponibilidadeResponse{
		ID:         d.ID,
		DiaSemana:  d.DiaSemana,
		HoraInicio: d.HoraInicio,
		HoraFim:    d.HoraFim,
	}
}
