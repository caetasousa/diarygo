package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/google/uuid"
)

// setupPreferencias prepara um cenario com 1 cliente + 2 profissionais ja persistidos
// e devolve o service pronto para uso, alem dos IDs relevantes.
func setupPreferencias(t *testing.T) (*service.PreferenciaService, uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	clienteRepo := memory.NewClienteRepository()
	profRepo := memory.NewProfissionalRepository()
	prefRepo := memory.NewPreferenciaRepository()

	usuarioID := uuid.New()
	cliente := &domain.Cliente{
		ID:        uuid.New(),
		UsuarioID: usuarioID,
		Nome:      "Ana Costa",
		CPF:       "52998224725",
		Telefone:  "62999991111",
		Score:     100,
	}
	if err := clienteRepo.Criar(ctx, cliente); err != nil {
		t.Fatalf("setup cliente: %v", err)
	}

	prof1 := &domain.Profissional{
		ID: uuid.New(), UsuarioID: uuid.New(), Nome: "Joana Silva",
		CPF: "11144477735", Telefone: "62999992222", Status: domain.StatusAprovada,
		NotaMedia: 4.7, TotalServicos: 12,
	}
	prof2 := &domain.Profissional{
		ID: uuid.New(), UsuarioID: uuid.New(), Nome: "Maria Oliveira",
		CPF: "22233344486", Telefone: "62999993333", Status: domain.StatusAprovada,
		NotaMedia: 4.2, TotalServicos: 5,
	}
	if err := profRepo.Criar(ctx, prof1); err != nil {
		t.Fatalf("setup prof1: %v", err)
	}
	if err := profRepo.Criar(ctx, prof2); err != nil {
		t.Fatalf("setup prof2: %v", err)
	}

	svc := service.NewPreferenciaService(prefRepo, clienteRepo, profRepo)
	return svc, usuarioID, prof1.ID, prof2.ID
}

func TestPreferencia_Favoritar_Sucesso(t *testing.T) {
	svc, usuarioID, prof1, _ := setupPreferencias(t)
	ctx := context.Background()

	if err := svc.Favoritar(ctx, usuarioID, prof1); err != nil {
		t.Fatalf("favoritar: %v", err)
	}

	favs, err := svc.Listar(ctx, usuarioID, domain.PreferenciaFavorita)
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if len(favs) != 1 {
		t.Fatalf("esperava 1 favorita, got %d", len(favs))
	}
	if favs[0].ProfissionalID != prof1 {
		t.Errorf("profissionalID esperado %s, got %s", prof1, favs[0].ProfissionalID)
	}
	if favs[0].Tipo != domain.PreferenciaFavorita {
		t.Errorf("tipo esperado FAVORITA, got %s", favs[0].Tipo)
	}
	if favs[0].Nome == "" || favs[0].NotaMedia == 0 {
		t.Errorf("response deveria ter nome e nota da profissional: %+v", favs[0])
	}
}

func TestPreferencia_Bloquear_SubstituiFavorita(t *testing.T) {
	svc, usuarioID, prof1, _ := setupPreferencias(t)
	ctx := context.Background()

	if err := svc.Favoritar(ctx, usuarioID, prof1); err != nil {
		t.Fatalf("favoritar: %v", err)
	}
	if err := svc.Bloquear(ctx, usuarioID, prof1); err != nil {
		t.Fatalf("bloquear: %v", err)
	}

	favs, _ := svc.Listar(ctx, usuarioID, domain.PreferenciaFavorita)
	if len(favs) != 0 {
		t.Errorf("apos bloquear, favoritas deveria estar vazio, got %d", len(favs))
	}
	bloqueios, _ := svc.Listar(ctx, usuarioID, domain.PreferenciaBloqueada)
	if len(bloqueios) != 1 {
		t.Fatalf("esperava 1 bloqueio, got %d", len(bloqueios))
	}
}

func TestPreferencia_Remover_Sucesso(t *testing.T) {
	svc, usuarioID, prof1, _ := setupPreferencias(t)
	ctx := context.Background()

	_ = svc.Favoritar(ctx, usuarioID, prof1)
	if err := svc.Remover(ctx, usuarioID, prof1); err != nil {
		t.Fatalf("remover: %v", err)
	}
	favs, _ := svc.Listar(ctx, usuarioID, domain.PreferenciaFavorita)
	if len(favs) != 0 {
		t.Errorf("esperava lista vazia apos remover, got %d", len(favs))
	}
}

func TestPreferencia_Remover_NaoEncontrada(t *testing.T) {
	svc, usuarioID, prof1, _ := setupPreferencias(t)
	ctx := context.Background()

	err := svc.Remover(ctx, usuarioID, prof1)
	if !errors.Is(err, domain.ErrPreferenciaNaoEncontrada) {
		t.Errorf("esperava ErrPreferenciaNaoEncontrada, got %v", err)
	}
}

func TestPreferencia_Favoritar_ProfissionalInexistente(t *testing.T) {
	svc, usuarioID, _, _ := setupPreferencias(t)
	ctx := context.Background()

	err := svc.Favoritar(ctx, usuarioID, uuid.New())
	if !errors.Is(err, domain.ErrProfissionalNaoEncontrada) {
		t.Errorf("esperava ErrProfissionalNaoEncontrada, got %v", err)
	}
}

func TestPreferencia_Listar_SegregaPorTipo(t *testing.T) {
	svc, usuarioID, prof1, prof2 := setupPreferencias(t)
	ctx := context.Background()

	_ = svc.Favoritar(ctx, usuarioID, prof1)
	_ = svc.Bloquear(ctx, usuarioID, prof2)

	favs, _ := svc.Listar(ctx, usuarioID, domain.PreferenciaFavorita)
	bloqueios, _ := svc.Listar(ctx, usuarioID, domain.PreferenciaBloqueada)

	if len(favs) != 1 || favs[0].ProfissionalID != prof1 {
		t.Errorf("favoritas esperavam [prof1], got %+v", favs)
	}
	if len(bloqueios) != 1 || bloqueios[0].ProfissionalID != prof2 {
		t.Errorf("bloqueios esperavam [prof2], got %+v", bloqueios)
	}
}

func TestPreferencia_Listar_TipoInvalido(t *testing.T) {
	svc, usuarioID, _, _ := setupPreferencias(t)
	_, err := svc.Listar(context.Background(), usuarioID, domain.TipoPreferencia("INVALIDO"))
	if !errors.Is(err, domain.ErrTipoPreferenciaInvalido) {
		t.Errorf("esperava ErrTipoPreferenciaInvalido, got %v", err)
	}
}
