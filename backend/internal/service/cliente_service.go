package service

import (
	"context"
	"strings"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// ClienteService gerencia a logica de negocio dos clientes.
type ClienteService struct {
	repo domain.ClienteRepository
}

// NewClienteService cria um novo ClienteService.
func NewClienteService(repo domain.ClienteRepository) *ClienteService {
	return &ClienteService{repo: repo}
}

// Criar cria o perfil de cliente associado ao usuario autenticado.
// O CPF e obrigatorio e deve ser valido.
func (s *ClienteService) Criar(ctx context.Context, usuarioID uuid.UUID, req domain.ClienteRequest) (*domain.ClienteResponse, error) {
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

	c := &domain.Cliente{
		ID:        uuid.New(),
		UsuarioID: usuarioID,
		Nome:      strings.TrimSpace(req.Nome),
		CPF:       cpf,
		Telefone:  telefone,
		Score:     100, // score inicial
	}

	if err := s.repo.Criar(ctx, c); err != nil {
		return nil, err
	}

	return toClienteResponse(c), nil
}

// Atualizar atualiza nome e telefone do cliente. CPF nao pode ser alterado apos criacao.
func (s *ClienteService) Atualizar(ctx context.Context, usuarioID uuid.UUID, req domain.ClienteRequest) (*domain.ClienteResponse, error) {
	c, err := s.repo.BuscarPorUsuarioID(ctx, usuarioID)
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

	c.Nome = strings.TrimSpace(req.Nome)
	c.Telefone = telefone

	if err := s.repo.Atualizar(ctx, c); err != nil {
		return nil, err
	}

	return toClienteResponse(c), nil
}

// BuscarPorUsuarioID retorna o perfil do cliente pelo ID do usuario autenticado.
func (s *ClienteService) BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*domain.ClienteResponse, error) {
	c, err := s.repo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	return toClienteResponse(c), nil
}

// toClienteResponse converte a entidade para o DTO de resposta.
func toClienteResponse(c *domain.Cliente) *domain.ClienteResponse {
	return &domain.ClienteResponse{
		ID:        c.ID,
		UsuarioID: c.UsuarioID,
		Nome:      c.Nome,
		CPF:       c.CPF,
		Telefone:  c.Telefone,
		Score:     c.Score,
	}
}

// limparMascara remove todos os caracteres nao numericos de uma string.
func limparMascara(s string) string {
	var sb strings.Builder
	for _, c := range s {
		if c >= '0' && c <= '9' {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}
