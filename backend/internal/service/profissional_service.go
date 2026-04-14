package service

import (
	"context"
	"strings"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// ProfissionalService gerencia a logica de negocio das profissionais.
type ProfissionalService struct {
	repo domain.ProfissionalRepository
}

// NewProfissionalService cria um novo ProfissionalService.
func NewProfissionalService(repo domain.ProfissionalRepository) *ProfissionalService {
	return &ProfissionalService{repo: repo}
}

// Criar cria o perfil de profissional associado ao usuario autenticado.
// Status inicial: PENDENTE.
func (s *ProfissionalService) Criar(ctx context.Context, usuarioID uuid.UUID, req domain.ProfissionalRequest) (*domain.ProfissionalResponse, error) {
	if err := domain.ValidarNome(req.Nome); err != nil {
		return nil, err
	}

	cpf := limparMascara(req.CPF)
	if err := domain.ValidarCPF(cpf); err != nil {
		return nil, err
	}

	telefone := limparMascara(req.Telefone)
	if err := domain.ValidarTelefone(telefone); err != nil {
		return nil, err
	}

	p := &domain.Profissional{
		ID:        uuid.New(),
		UsuarioID: usuarioID,
		Nome:      strings.TrimSpace(req.Nome),
		CPF:       cpf,
		RG:        strings.TrimSpace(req.RG),
		Telefone:  telefone,
		FotoURL:   strings.TrimSpace(req.FotoURL),
		MEI:       req.MEI,
		Status:    domain.StatusPendente,
	}

	if err := s.repo.Criar(ctx, p); err != nil {
		return nil, err
	}

	return toProfissionalResponse(p), nil
}

// Atualizar atualiza o perfil da profissional (nome, telefone, foto, MEI).
// Status e nota nao podem ser alterados pela profissional.
func (s *ProfissionalService) Atualizar(ctx context.Context, usuarioID uuid.UUID, req domain.ProfissionalRequest) (*domain.ProfissionalResponse, error) {
	p, err := s.repo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidarNome(req.Nome); err != nil {
		return nil, err
	}

	telefone := limparMascara(req.Telefone)
	if err := domain.ValidarTelefone(telefone); err != nil {
		return nil, err
	}

	p.Nome = strings.TrimSpace(req.Nome)
	p.RG = strings.TrimSpace(req.RG)
	p.Telefone = telefone
	p.FotoURL = strings.TrimSpace(req.FotoURL)
	p.MEI = req.MEI

	if err := s.repo.Atualizar(ctx, p); err != nil {
		return nil, err
	}

	return toProfissionalResponse(p), nil
}

// BuscarPorUsuarioID retorna o perfil da profissional pelo ID do usuario autenticado.
func (s *ProfissionalService) BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*domain.ProfissionalResponse, error) {
	p, err := s.repo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	return toProfissionalResponse(p), nil
}

// BuscarPorID retorna o perfil da profissional pelo UUID da entidade.
func (s *ProfissionalService) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.ProfissionalResponse, error) {
	p, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toProfissionalResponse(p), nil
}

// toProfissionalResponse converte a entidade para o DTO de resposta.
func toProfissionalResponse(p *domain.Profissional) *domain.ProfissionalResponse {
	return &domain.ProfissionalResponse{
		ID:            p.ID,
		UsuarioID:     p.UsuarioID,
		Nome:          p.Nome,
		CPF:           p.CPF,
		RG:            p.RG,
		Telefone:      p.Telefone,
		FotoURL:       p.FotoURL,
		MEI:           p.MEI,
		Status:        p.Status,
		NotaMedia:     p.NotaMedia,
		TotalServicos: p.TotalServicos,
	}
}
