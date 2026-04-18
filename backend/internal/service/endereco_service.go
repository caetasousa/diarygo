package service

import (
	"context"
	"strings"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// EnderecoService gerencia a logica de negocio dos enderecos.
type EnderecoService struct {
	repo        domain.EnderecoRepository
	clienteRepo domain.ClienteRepository
}

// NewEnderecoService cria um novo EnderecoService.
func NewEnderecoService(repo domain.EnderecoRepository, clienteRepo domain.ClienteRepository) *EnderecoService {
	return &EnderecoService{repo: repo, clienteRepo: clienteRepo}
}

// Criar cria um novo endereco para o cliente autenticado.
func (s *EnderecoService) Criar(ctx context.Context, usuarioID uuid.UUID, req domain.EnderecoRequest) (*domain.EnderecoResponse, error) {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	if err := s.validarRequest(req); err != nil {
		return nil, err
	}

	// Se definido como principal, desmarcar outros
	if req.Principal {
		if err := s.desmarcarPrincipal(ctx, cliente.ID); err != nil {
			return nil, err
		}
	}

	e := &domain.Endereco{
		ID:           uuid.New(),
		ClienteID:    cliente.ID,
		Logradouro:   strings.TrimSpace(req.Logradouro),
		Numero:       strings.TrimSpace(req.Numero),
		Complemento:  strings.TrimSpace(req.Complemento),
		Bairro:       strings.TrimSpace(req.Bairro),
		Cidade:       strings.TrimSpace(req.Cidade),
		Estado:       strings.ToUpper(strings.TrimSpace(req.Estado)),
		CEP:          limparMascara(req.CEP),
		NumQuartos:   req.NumQuartos,
		NumBanheiros: req.NumBanheiros,
		NumSalas:     req.NumSalas,
		NumCozinhas:  req.NumCozinhas,
		AreaM2:       req.AreaM2,
		Principal:    req.Principal,
	}

	if err := s.repo.Criar(ctx, e); err != nil {
		return nil, err
	}

	return toEnderecoResponse(e), nil
}

// Listar retorna todos os enderecos do cliente autenticado.
func (s *EnderecoService) Listar(ctx context.Context, usuarioID uuid.UUID) ([]*domain.EnderecoResponse, error) {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	enderecos, err := s.repo.ListarPorClienteID(ctx, cliente.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.EnderecoResponse, len(enderecos))
	for i, e := range enderecos {
		resp[i] = toEnderecoResponse(e)
	}
	return resp, nil
}

// Atualizar atualiza um endereco existente. Verifica que pertence ao cliente.
func (s *EnderecoService) Atualizar(ctx context.Context, usuarioID uuid.UUID, enderecoID uuid.UUID, req domain.EnderecoRequest) (*domain.EnderecoResponse, error) {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	e, err := s.repo.BuscarPorID(ctx, enderecoID)
	if err != nil {
		return nil, err
	}

	if e.ClienteID != cliente.ID {
		return nil, domain.ErrEnderecoNaoPertenceAoCliente
	}

	if err := s.validarRequest(req); err != nil {
		return nil, err
	}

	// Se definido como principal, desmarcar outros
	if req.Principal && !e.Principal {
		if err := s.desmarcarPrincipal(ctx, cliente.ID); err != nil {
			return nil, err
		}
	}

	e.Logradouro = strings.TrimSpace(req.Logradouro)
	e.Numero = strings.TrimSpace(req.Numero)
	e.Complemento = strings.TrimSpace(req.Complemento)
	e.Bairro = strings.TrimSpace(req.Bairro)
	e.Cidade = strings.TrimSpace(req.Cidade)
	e.Estado = strings.ToUpper(strings.TrimSpace(req.Estado))
	e.CEP = limparMascara(req.CEP)
	e.NumQuartos = req.NumQuartos
	e.NumBanheiros = req.NumBanheiros
	e.NumSalas = req.NumSalas
	e.NumCozinhas = req.NumCozinhas
	e.AreaM2 = req.AreaM2
	e.Principal = req.Principal

	if err := s.repo.Atualizar(ctx, e); err != nil {
		return nil, err
	}

	return toEnderecoResponse(e), nil
}

// DefinirPrincipal marca um endereco como principal e desmarca os outros.
func (s *EnderecoService) DefinirPrincipal(ctx context.Context, usuarioID uuid.UUID, enderecoID uuid.UUID) (*domain.EnderecoResponse, error) {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	e, err := s.repo.BuscarPorID(ctx, enderecoID)
	if err != nil {
		return nil, err
	}

	if e.ClienteID != cliente.ID {
		return nil, domain.ErrEnderecoNaoPertenceAoCliente
	}

	if err := s.desmarcarPrincipal(ctx, cliente.ID); err != nil {
		return nil, err
	}

	e.Principal = true
	if err := s.repo.Atualizar(ctx, e); err != nil {
		return nil, err
	}

	return toEnderecoResponse(e), nil
}

// Remover remove um endereco. Verifica que pertence ao cliente.
func (s *EnderecoService) Remover(ctx context.Context, usuarioID uuid.UUID, enderecoID uuid.UUID) error {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return err
	}

	e, err := s.repo.BuscarPorID(ctx, enderecoID)
	if err != nil {
		return err
	}

	if e.ClienteID != cliente.ID {
		return domain.ErrEnderecoNaoPertenceAoCliente
	}

	return s.repo.Remover(ctx, enderecoID)
}

// desmarcarPrincipal desmarca todos os enderecos principais de um cliente.
func (s *EnderecoService) desmarcarPrincipal(ctx context.Context, clienteID uuid.UUID) error {
	enderecos, err := s.repo.ListarPorClienteID(ctx, clienteID)
	if err != nil {
		return err
	}
	for _, e := range enderecos {
		if e.Principal {
			e.Principal = false
			if err := s.repo.Atualizar(ctx, e); err != nil {
				return err
			}
		}
	}
	return nil
}

// validarRequest valida os campos obrigatorios de um EnderecoRequest.
// Espelha os CHECKs da tabela `enderecos` (cep formato, num_* >= 0,
// area_m2 > 0 quando presente) e adiciona a regra de negocio de que todo
// endereco residencial precisa ter ao menos 1 quarto.
func (s *EnderecoService) validarRequest(req domain.EnderecoRequest) error {
	cep := limparMascara(req.CEP)
	if err := domain.ValidarCEP(cep); err != nil {
		return err
	}
	if strings.TrimSpace(req.Logradouro) == "" {
		return domain.ErrLogradouroObrigatorio
	}
	if strings.TrimSpace(req.Cidade) == "" {
		return domain.ErrCidadeObrigatoria
	}
	estado := strings.ToUpper(strings.TrimSpace(req.Estado))
	if err := domain.ValidarEstado(estado); err != nil {
		return err
	}
	if err := domain.ValidarComodos(req.NumQuartos, req.NumBanheiros, req.NumSalas, req.NumCozinhas); err != nil {
		return err
	}
	if err := domain.ValidarAreaM2(req.AreaM2); err != nil {
		return err
	}
	return nil
}

// toEnderecoResponse converte a entidade para o DTO de resposta.
func toEnderecoResponse(e *domain.Endereco) *domain.EnderecoResponse {
	return &domain.EnderecoResponse{
		ID:           e.ID,
		ClienteID:    e.ClienteID,
		Logradouro:   e.Logradouro,
		Numero:       e.Numero,
		Complemento:  e.Complemento,
		Bairro:       e.Bairro,
		Cidade:       e.Cidade,
		Estado:       e.Estado,
		CEP:          e.CEP,
		NumQuartos:   e.NumQuartos,
		NumBanheiros: e.NumBanheiros,
		NumSalas:     e.NumSalas,
		NumCozinhas:  e.NumCozinhas,
		AreaM2:       e.AreaM2,
		Principal:    e.Principal,
	}
}
