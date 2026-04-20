package service

import (
	"context"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// SolicitacaoService gerencia o ciclo de vida de solicitações de serviço.
type SolicitacaoService struct {
	solicRepo    domain.SolicitacaoRepository
	enderecoRepo domain.EnderecoRepository
	clienteRepo  domain.ClienteRepository
	precificacao *PrecificacaoService
	agora        func() time.Time // injetável em testes
}

// NewSolicitacaoService cria um novo SolicitacaoService.
func NewSolicitacaoService(
	solic domain.SolicitacaoRepository,
	end domain.EnderecoRepository,
	cli domain.ClienteRepository,
	prec *PrecificacaoService,
) *SolicitacaoService {
	return &SolicitacaoService{
		solicRepo:    solic,
		enderecoRepo: end,
		clienteRepo:  cli,
		precificacao: prec,
		agora:        time.Now,
	}
}

// NewSolicitacaoServiceComAgora é usado em testes para injetar um clock controlável.
func NewSolicitacaoServiceComAgora(
	solic domain.SolicitacaoRepository,
	end domain.EnderecoRepository,
	cli domain.ClienteRepository,
	prec *PrecificacaoService,
	agora func() time.Time,
) *SolicitacaoService {
	return &SolicitacaoService{
		solicRepo:    solic,
		enderecoRepo: end,
		clienteRepo:  cli,
		precificacao: prec,
		agora:        agora,
	}
}

// Criar valida e persiste uma nova solicitação com orçamento congelado.
// Recebe usuarioID (do JWT) e resolve para clienteID internamente.
func (s *SolicitacaoService) Criar(ctx context.Context, usuarioID uuid.UUID, req domain.CriarSolicitacaoRequest) (*domain.SolicitacaoResponse, error) {
	if err := domain.ValidarObservacao(req.Observacao); err != nil {
		return nil, err
	}

	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	clienteID := cliente.ID

	// Ownership do endereço
	end, err := s.enderecoRepo.BuscarPorID(ctx, req.EnderecoID)
	if err != nil || end.ClienteID != clienteID {
		return nil, domain.ErrSolicitacaoEnderecoInvalido
	}

	// Antecedência mínima de 24h
	if err := domain.ValidarAntecedencia24h(req.DataServico, s.agora()); err != nil {
		return nil, err
	}

	// Calcular orçamento e capturar snapshots de opcionais
	calcReq := domain.CalculoPrecoRequest{
		CategoriaID:  req.CategoriaID,
		RegiaoID:     req.RegiaoID,
		NumQuartos:   req.NumQuartos,
		NumBanheiros: req.NumBanheiros,
		NumSalas:     req.NumSalas,
		NumCozinhas:  req.NumCozinhas,
		OpcionaisIDs: req.OpcionaisIDs,
		Frequencia:   req.Frequencia,
		DataServico:  req.DataServico,
	}
	resp, snapshots, err := s.precificacao.CalcularComOpcionais(ctx, calcReq)
	if err != nil {
		return nil, err
	}

	id := uuid.New()
	for i := range snapshots {
		snapshots[i].SolicitacaoID = id
	}

	sol := &domain.Solicitacao{
		ID:           id,
		ClienteID:    clienteID,
		EnderecoID:   req.EnderecoID,
		CategoriaID:  req.CategoriaID,
		RegiaoID:     req.RegiaoID,
		NumQuartos:   req.NumQuartos,
		NumBanheiros: req.NumBanheiros,
		NumSalas:     req.NumSalas,
		NumCozinhas:  req.NumCozinhas,
		Frequencia:   req.Frequencia,
		DataServico:  req.DataServico,
		Observacao:   req.Observacao,
		ValorTotal:   resp.ValorTotal,
		DuracaoMin:   resp.DuracaoMin,
		Breakdown:    resp.Itens,
		Status:       domain.StatusAguardando,
	}

	if err := s.solicRepo.Criar(ctx, sol, snapshots); err != nil {
		return nil, err
	}

	return toSolicitacaoResponse(sol), nil
}

// Cancelar cancela uma solicitação, penalizando o score se < 24h do serviço.
// Recebe usuarioID (do JWT) e resolve para clienteID internamente.
func (s *SolicitacaoService) Cancelar(ctx context.Context, usuarioID, solicitacaoID uuid.UUID) error {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return domain.ErrSolicitacaoNaoEncontrada
	}
	clienteID := cliente.ID

	sol, _, err := s.solicRepo.BuscarPorID(ctx, solicitacaoID)
	if err != nil || sol.ClienteID != clienteID {
		return domain.ErrSolicitacaoNaoEncontrada
	}

	if !sol.Status.PodeCancelar() {
		return domain.ErrSolicitacaoStatusInvalido
	}

	// Penalidade de score se menos de 24h antes do serviço
	if sol.DataServico.Sub(s.agora()) < 24*time.Hour {
		_ = s.clienteRepo.AjustarScore(ctx, clienteID, -5)
	}

	agora := s.agora()
	return s.solicRepo.AtualizarStatus(ctx, solicitacaoID, domain.StatusCancelada, &agora)
}

// BuscarPorID retorna uma solicitação pelo ID com verificação de ownership.
// Retorna ErrSolicitacaoNaoEncontrada se não pertencer ao usuário (A01 — não vazar existência).
func (s *SolicitacaoService) BuscarPorID(ctx context.Context, usuarioID, solicitacaoID uuid.UUID) (*domain.SolicitacaoResponse, error) {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, domain.ErrSolicitacaoNaoEncontrada
	}

	sol, _, err := s.solicRepo.BuscarPorID(ctx, solicitacaoID)
	if err != nil || sol.ClienteID != cliente.ID {
		return nil, domain.ErrSolicitacaoNaoEncontrada
	}
	return toSolicitacaoResponse(sol), nil
}

// ListarDoCliente lista as solicitações do cliente com filtros opcionais.
func (s *SolicitacaoService) ListarDoCliente(ctx context.Context, usuarioID uuid.UUID, filtro domain.SolicitacaoFiltro) ([]*domain.SolicitacaoResponse, error) {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		if err == domain.ErrClienteNaoEncontrado {
			return []*domain.SolicitacaoResponse{}, nil
		}
		return nil, err
	}

	sols, err := s.solicRepo.ListarPorCliente(ctx, cliente.ID, filtro)
	if err != nil {
		return nil, err
	}
	resp := make([]*domain.SolicitacaoResponse, 0, len(sols))
	for _, sol := range sols {
		resp = append(resp, toSolicitacaoResponse(sol))
	}
	return resp, nil
}

func toSolicitacaoResponse(s *domain.Solicitacao) *domain.SolicitacaoResponse {
	return &domain.SolicitacaoResponse{
		ID:           s.ID,
		EnderecoID:   s.EnderecoID,
		CategoriaID:  s.CategoriaID,
		RegiaoID:     s.RegiaoID,
		NumQuartos:   s.NumQuartos,
		NumBanheiros: s.NumBanheiros,
		NumSalas:     s.NumSalas,
		NumCozinhas:  s.NumCozinhas,
		Frequencia:   s.Frequencia,
		DataServico:  s.DataServico,
		Observacao:   s.Observacao,
		ValorTotal:   s.ValorTotal,
		DuracaoMin:   s.DuracaoMin,
		Breakdown:    s.Breakdown,
		Status:       s.Status,
		CriadaEm:     s.CriadaEm,
		CanceladaEm:  s.CanceladaEm,
	}
}
