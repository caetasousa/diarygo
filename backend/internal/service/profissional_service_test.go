package service_test

import (
	"context"
	"testing"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/google/uuid"
)

func novoProfissionalService() *service.ProfissionalService {
	return service.NewProfissionalService(memory.NewProfissionalRepository())
}

func reqProfissionalValida() domain.ProfissionalRequest {
	return domain.ProfissionalRequest{
		Nome:     "Ana Diarista",
		CPF:      "529.982.247-25",
		RG:       "123456789",
		Telefone: "(11) 91234-5678",
		MEI:      false,
	}
}

// TestCriarProfissional_Sucesso verifica criacao com status PENDENTE.
func TestCriarProfissional_Sucesso(t *testing.T) {
	svc := novoProfissionalService()
	ctx := context.Background()
	usuarioID := uuid.New()

	resp, err := svc.Criar(ctx, usuarioID, reqProfissionalValida())
	if err != nil {
		t.Fatalf("esperava sucesso, got: %v", err)
	}
	if resp.Status != domain.StatusPendente {
		t.Errorf("status inicial esperado PENDENTE, got %s", resp.Status)
	}
	if resp.CPF != "52998224725" {
		t.Errorf("CPF esperado sem mascara, got '%s'", resp.CPF)
	}
	if resp.NotaMedia != 0.0 {
		t.Errorf("nota_media inicial esperada 0.0, got %f", resp.NotaMedia)
	}
}

// TestCriarProfissional_UsuarioDuplicado verifica que o mesmo usuario nao cria dois perfis.
func TestCriarProfissional_UsuarioDuplicado(t *testing.T) {
	svc := novoProfissionalService()
	ctx := context.Background()
	usuarioID := uuid.New()

	if _, err := svc.Criar(ctx, usuarioID, reqProfissionalValida()); err != nil {
		t.Fatalf("primeiro registro falhou: %v", err)
	}

	_, err := svc.Criar(ctx, usuarioID, reqProfissionalValida())
	if err == nil {
		t.Fatal("esperava erro de profissional ja existente")
	}
}

// TestCriarProfissional_CPFInvalido verifica rejeicao de CPF invalido.
func TestCriarProfissional_CPFInvalido(t *testing.T) {
	svc := novoProfissionalService()
	req := reqProfissionalValida()
	req.CPF = "000.000.000-00"

	_, err := svc.Criar(context.Background(), uuid.New(), req)
	if err == nil {
		t.Fatal("esperava erro de CPF invalido")
	}
}

// TestAtualizarProfissional_Sucesso verifica que nome e telefone sao atualizados.
func TestAtualizarProfissional_Sucesso(t *testing.T) {
	svc := novoProfissionalService()
	ctx := context.Background()
	usuarioID := uuid.New()

	if _, err := svc.Criar(ctx, usuarioID, reqProfissionalValida()); err != nil {
		t.Fatalf("criacao falhou: %v", err)
	}

	resp, err := svc.Atualizar(ctx, usuarioID, domain.ProfissionalRequest{
		Nome:     "Ana Diarista Atualizada",
		CPF:      "529.982.247-25",
		Telefone: "11999999999",
		MEI:      true,
	})
	if err != nil {
		t.Fatalf("atualizacao falhou: %v", err)
	}
	if resp.Nome != "Ana Diarista Atualizada" {
		t.Errorf("nome incorreto: %s", resp.Nome)
	}
	if !resp.MEI {
		t.Error("MEI deveria ser true apos atualizacao")
	}
	// Status nao deve mudar
	if resp.Status != domain.StatusPendente {
		t.Errorf("status nao deve mudar via Atualizar, got %s", resp.Status)
	}
}

// TestAtualizarProfissional_NaoEncontrada verifica erro quando profissional nao existe.
func TestAtualizarProfissional_NaoEncontrada(t *testing.T) {
	svc := novoProfissionalService()

	_, err := svc.Atualizar(context.Background(), uuid.New(), reqProfissionalValida())
	if err == nil {
		t.Fatal("esperava erro de profissional nao encontrada")
	}
}

// TestBuscarProfissional_Sucesso verifica busca por usuarioID.
func TestBuscarProfissional_Sucesso(t *testing.T) {
	svc := novoProfissionalService()
	ctx := context.Background()
	usuarioID := uuid.New()

	if _, err := svc.Criar(ctx, usuarioID, reqProfissionalValida()); err != nil {
		t.Fatalf("criacao falhou: %v", err)
	}

	resp, err := svc.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		t.Fatalf("busca falhou: %v", err)
	}
	if resp.UsuarioID != usuarioID {
		t.Errorf("usuarioID incorreto: %s", resp.UsuarioID)
	}
}
