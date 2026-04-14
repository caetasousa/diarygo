package service_test

import (
	"context"
	"testing"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/google/uuid"
)

type credSetup struct {
	svc       *service.CredenciamentoService
	profSvc   *service.ProfissionalService
	usuarioID uuid.UUID
}

func novoCredenciamentoSetup(t *testing.T) *credSetup {
	t.Helper()
	profRepo := memory.NewProfissionalRepository()
	docRepo := memory.NewDocumentoRepository()
	refRepo := memory.NewReferenciaRepository()
	regiaoRepo := memory.NewRegiaoRepository()
	profRegiaoRepo := memory.NewProfissionalRegiaoRepository(regiaoRepo)
	dispRepo := memory.NewDisponibilidadeRepository()

	profSvc := service.NewProfissionalService(profRepo)
	credSvc := service.NewCredenciamentoService(profRepo, docRepo, refRepo, regiaoRepo, profRegiaoRepo, dispRepo)

	usuarioID := uuid.New()
	if _, err := profSvc.Criar(context.Background(), usuarioID, domain.ProfissionalRequest{
		Nome:     "Cred Tester",
		CPF:      "529.982.247-25",
		Telefone: "11911111111",
	}); err != nil {
		t.Fatalf("criacao de profissional falhou: %v", err)
	}

	return &credSetup{svc: credSvc, profSvc: profSvc, usuarioID: usuarioID}
}

// TestEnviarDocumento_Sucesso verifica criacao de documento com status PENDENTE.
func TestEnviarDocumento_Sucesso(t *testing.T) {
	cs := novoCredenciamentoSetup(t)
	ctx := context.Background()

	resp, err := cs.svc.EnviarDocumento(ctx, cs.usuarioID, domain.DocumentoRequest{
		Tipo: domain.DocRGFrente,
		URL:  "https://cdn.example.com/rg-frente.jpg",
	})
	if err != nil {
		t.Fatalf("esperava sucesso, got: %v", err)
	}
	if resp.Status != domain.DocPendente {
		t.Errorf("status esperado PENDENTE, got %s", resp.Status)
	}
	if resp.Tipo != domain.DocRGFrente {
		t.Errorf("tipo incorreto: %s", resp.Tipo)
	}
}

// TestEnviarDocumento_TipoInvalido verifica rejeicao de tipo desconhecido.
func TestEnviarDocumento_TipoInvalido(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	_, err := cs.svc.EnviarDocumento(context.Background(), cs.usuarioID, domain.DocumentoRequest{
		Tipo: "INVALIDO",
		URL:  "https://cdn.example.com/doc.jpg",
	})
	if err == nil {
		t.Fatal("esperava erro de tipo invalido")
	}
}

// TestEnviarDocumento_URLVazia verifica rejeicao de URL vazia.
func TestEnviarDocumento_URLVazia(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	_, err := cs.svc.EnviarDocumento(context.Background(), cs.usuarioID, domain.DocumentoRequest{
		Tipo: domain.DocCPF,
		URL:  "",
	})
	if err == nil {
		t.Fatal("esperava erro de URL vazia")
	}
}

// TestListarDocumentos_Sucesso verifica listagem de documentos.
func TestListarDocumentos_Sucesso(t *testing.T) {
	cs := novoCredenciamentoSetup(t)
	ctx := context.Background()

	cs.svc.EnviarDocumento(ctx, cs.usuarioID, domain.DocumentoRequest{Tipo: domain.DocRGFrente, URL: "https://cdn.example.com/rg1.jpg"}) //nolint:errcheck
	cs.svc.EnviarDocumento(ctx, cs.usuarioID, domain.DocumentoRequest{Tipo: domain.DocRGVerso, URL: "https://cdn.example.com/rg2.jpg"})  //nolint:errcheck

	lista, err := cs.svc.ListarDocumentos(ctx, cs.usuarioID)
	if err != nil {
		t.Fatalf("listagem falhou: %v", err)
	}
	if len(lista) != 2 {
		t.Errorf("esperava 2 documentos, got %d", len(lista))
	}
}

// TestAdicionarReferencia_Sucesso verifica criacao com status PENDENTE.
func TestAdicionarReferencia_Sucesso(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	resp, err := cs.svc.AdicionarReferencia(context.Background(), cs.usuarioID, domain.ReferenciaRequest{
		NomeContato:     "Dona Helena",
		TelefoneContato: "11988887777",
	})
	if err != nil {
		t.Fatalf("esperava sucesso, got: %v", err)
	}
	if resp.Status != domain.RefPendente {
		t.Errorf("status esperado PENDENTE, got %s", resp.Status)
	}
}

// TestAdicionarReferencia_NomeVazio verifica rejeicao de nome vazio.
func TestAdicionarReferencia_NomeVazio(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	_, err := cs.svc.AdicionarReferencia(context.Background(), cs.usuarioID, domain.ReferenciaRequest{
		NomeContato:     "",
		TelefoneContato: "11988887777",
	})
	if err == nil {
		t.Fatal("esperava erro de nome vazio")
	}
}

// TestListarRegioes_Sucesso verifica que regioes do seed sao retornadas.
func TestListarRegioes_Sucesso(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	lista, err := cs.svc.ListarRegioes(context.Background())
	if err != nil {
		t.Fatalf("esperava sucesso, got: %v", err)
	}
	if len(lista) == 0 {
		t.Error("esperava pelo menos uma regiao do seed")
	}
}

// TestDefinirRegioes_Sucesso verifica associacao de regioes validas.
func TestDefinirRegioes_Sucesso(t *testing.T) {
	cs := novoCredenciamentoSetup(t)
	ctx := context.Background()

	regioes, _ := cs.svc.ListarRegioes(ctx)
	if len(regioes) == 0 {
		t.Skip("sem regioes no seed")
	}

	ids := []uuid.UUID{regioes[0].ID}
	resp, err := cs.svc.DefinirRegioes(ctx, cs.usuarioID, domain.DefinirRegioesRequest{RegiaoIDs: ids})
	if err != nil {
		t.Fatalf("definir regioes falhou: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("esperava 1 regiao associada, got %d", len(resp))
	}
}

// TestDefinirRegioes_IDInvalido verifica rejeicao de regiao inexistente.
func TestDefinirRegioes_IDInvalido(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	_, err := cs.svc.DefinirRegioes(context.Background(), cs.usuarioID, domain.DefinirRegioesRequest{
		RegiaoIDs: []uuid.UUID{uuid.New()}, // ID que nao existe
	})
	if err == nil {
		t.Fatal("esperava erro de regiao invalida")
	}
}

// TestDefinirDisponibilidades_Sucesso verifica substituicao completa dos slots.
func TestDefinirDisponibilidades_Sucesso(t *testing.T) {
	cs := novoCredenciamentoSetup(t)
	ctx := context.Background()

	req := domain.DefinirDisponibilidadesRequest{
		Slots: []domain.DisponibilidadeSlot{
			{DiaSemana: 1, HoraInicio: "08:00", HoraFim: "17:00"},
			{DiaSemana: 3, HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}

	resp, err := cs.svc.DefinirDisponibilidades(ctx, cs.usuarioID, req)
	if err != nil {
		t.Fatalf("definir disponibilidades falhou: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("esperava 2 slots, got %d", len(resp))
	}
}

// TestDefinirDisponibilidades_Substitui verifica que chamada subsequente substitui.
func TestDefinirDisponibilidades_Substitui(t *testing.T) {
	cs := novoCredenciamentoSetup(t)
	ctx := context.Background()

	req1 := domain.DefinirDisponibilidadesRequest{
		Slots: []domain.DisponibilidadeSlot{
			{DiaSemana: 1, HoraInicio: "08:00", HoraFim: "17:00"},
			{DiaSemana: 2, HoraInicio: "08:00", HoraFim: "17:00"},
		},
	}
	cs.svc.DefinirDisponibilidades(ctx, cs.usuarioID, req1) //nolint:errcheck

	req2 := domain.DefinirDisponibilidadesRequest{
		Slots: []domain.DisponibilidadeSlot{
			{DiaSemana: 5, HoraInicio: "09:00", HoraFim: "14:00"},
		},
	}
	resp, err := cs.svc.DefinirDisponibilidades(ctx, cs.usuarioID, req2)
	if err != nil {
		t.Fatalf("segunda definicao falhou: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("esperava 1 slot apos substituicao, got %d", len(resp))
	}
}

// TestDefinirDisponibilidades_HoraFimAntesDaInicio verifica validacao de horario.
func TestDefinirDisponibilidades_HoraFimAntesDaInicio(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	req := domain.DefinirDisponibilidadesRequest{
		Slots: []domain.DisponibilidadeSlot{
			{DiaSemana: 1, HoraInicio: "17:00", HoraFim: "08:00"},
		},
	}
	_, err := cs.svc.DefinirDisponibilidades(context.Background(), cs.usuarioID, req)
	if err == nil {
		t.Fatal("esperava erro de hora_fim antes da hora_inicio")
	}
}

// TestDefinirDisponibilidades_DiaSemanaInvalido verifica rejeicao de dia invalido.
func TestDefinirDisponibilidades_DiaSemanaInvalido(t *testing.T) {
	cs := novoCredenciamentoSetup(t)

	req := domain.DefinirDisponibilidadesRequest{
		Slots: []domain.DisponibilidadeSlot{
			{DiaSemana: 7, HoraInicio: "08:00", HoraFim: "17:00"}, // 7 invalido
		},
	}
	_, err := cs.svc.DefinirDisponibilidades(context.Background(), cs.usuarioID, req)
	if err == nil {
		t.Fatal("esperava erro de dia da semana invalido")
	}
}
