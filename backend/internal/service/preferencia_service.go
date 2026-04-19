package service

import (
	"context"
	"errors"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// PreferenciaService gerencia favoritas e bloqueios do cliente em relacao a profissionais.
type PreferenciaService struct {
	repo             domain.PreferenciaRepository
	clienteRepo      domain.ClienteRepository
	profissionalRepo domain.ProfissionalRepository
}

// NewPreferenciaService cria um novo PreferenciaService.
func NewPreferenciaService(
	repo domain.PreferenciaRepository,
	clienteRepo domain.ClienteRepository,
	profissionalRepo domain.ProfissionalRepository,
) *PreferenciaService {
	return &PreferenciaService{
		repo:             repo,
		clienteRepo:      clienteRepo,
		profissionalRepo: profissionalRepo,
	}
}

// Favoritar marca a profissional como favorita do cliente.
// Substitui uma entrada BLOQUEADA existente.
func (s *PreferenciaService) Favoritar(ctx context.Context, usuarioID, profissionalID uuid.UUID) error {
	return s.registrar(ctx, usuarioID, profissionalID, domain.PreferenciaFavorita)
}

// Bloquear marca a profissional como bloqueada para o cliente.
// Substitui uma entrada FAVORITA existente.
func (s *PreferenciaService) Bloquear(ctx context.Context, usuarioID, profissionalID uuid.UUID) error {
	return s.registrar(ctx, usuarioID, profissionalID, domain.PreferenciaBloqueada)
}

// Remover apaga a preferencia (seja favorita ou bloqueada) daquela profissional.
func (s *PreferenciaService) Remover(ctx context.Context, usuarioID, profissionalID uuid.UUID) error {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return err
	}
	return s.repo.Remover(ctx, cliente.ID, profissionalID)
}

// Listar retorna as preferencias do cliente filtradas por tipo (FAVORITA ou BLOQUEADA).
// Enriquece com dados publicos da profissional (nome, foto, nota).
func (s *PreferenciaService) Listar(ctx context.Context, usuarioID uuid.UUID, tipo domain.TipoPreferencia) ([]*domain.PreferenciaResponse, error) {
	if tipo != "" && !tipo.Valida() {
		return nil, domain.ErrTipoPreferenciaInvalido
	}

	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	prefs, err := s.repo.ListarPorCliente(ctx, cliente.ID, tipo)
	if err != nil {
		return nil, err
	}

	resp := make([]*domain.PreferenciaResponse, 0, len(prefs))
	for _, p := range prefs {
		prof, err := s.profissionalRepo.BuscarPorID(ctx, p.ProfissionalID)
		if err != nil {
			// Profissional descadastrada — ignora a entrada sem quebrar a listagem.
			if errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
				continue
			}
			return nil, err
		}
		resp = append(resp, &domain.PreferenciaResponse{
			ProfissionalID: prof.ID,
			Nome:           prof.Nome,
			FotoURL:        prof.FotoURL,
			NotaMedia:      prof.NotaMedia,
			Tipo:           p.Tipo,
			CriadoEm:       p.CriadoEm,
		})
	}
	return resp, nil
}

// registrar faz o upsert da preferencia com o tipo informado,
// validando que o cliente existe e que a profissional existe.
func (s *PreferenciaService) registrar(ctx context.Context, usuarioID, profissionalID uuid.UUID, tipo domain.TipoPreferencia) error {
	cliente, err := s.clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return err
	}
	if _, err := s.profissionalRepo.BuscarPorID(ctx, profissionalID); err != nil {
		return err
	}
	return s.repo.Upsert(ctx, &domain.ClienteProfissionalPreferencia{
		ClienteID:      cliente.ID,
		ProfissionalID: profissionalID,
		Tipo:           tipo,
	})
}
