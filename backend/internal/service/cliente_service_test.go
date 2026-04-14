package service_test

import (
	"context"
	"testing"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/google/uuid"
)

func novoClienteService() *service.ClienteService {
	return service.NewClienteService(memory.NewClienteRepository())
}

// TestCriarCliente_Sucesso verifica que um cliente e criado com score inicial 100.
func TestCriarCliente_Sucesso(t *testing.T) {
	svc := novoClienteService()
	ctx := context.Background()
	usuarioID := uuid.New()

	resp, err := svc.Criar(ctx, usuarioID, domain.ClienteRequest{
		Nome:     "Maria Silva",
		CPF:      "529.982.247-25",
		Telefone: "(11) 91234-5678",
	})

	if err != nil {
		t.Fatalf("esperava sucesso, got: %v", err)
	}
	if resp.Score != 100 {
		t.Errorf("score inicial esperado 100, got %d", resp.Score)
	}
	if resp.CPF != "52998224725" {
		t.Errorf("CPF esperado sem mascara '52998224725', got '%s'", resp.CPF)
	}
	if resp.UsuarioID != usuarioID {
		t.Errorf("usuarioID esperado %s, got %s", usuarioID, resp.UsuarioID)
	}
}

// TestCriarCliente_CPFInvalido verifica que CPF invalido e rejeitado.
func TestCriarCliente_CPFInvalido(t *testing.T) {
	svc := novoClienteService()
	ctx := context.Background()

	_, err := svc.Criar(ctx, uuid.New(), domain.ClienteRequest{
		Nome:     "Joao Santos",
		CPF:      "111.111.111-11",
		Telefone: "11912345678",
	})

	if err == nil {
		t.Fatal("esperava erro de CPF invalido")
	}
}

// TestCriarCliente_CPFDuplicado verifica que CPF duplicado retorna erro.
func TestCriarCliente_CPFDuplicado(t *testing.T) {
	svc := novoClienteService()
	ctx := context.Background()

	req := domain.ClienteRequest{Nome: "Ana Souza", CPF: "529.982.247-25", Telefone: "11912345678"}
	if _, err := svc.Criar(ctx, uuid.New(), req); err != nil {
		t.Fatalf("primeiro registro falhou: %v", err)
	}

	_, err := svc.Criar(ctx, uuid.New(), req)
	if err == nil {
		t.Fatal("esperava erro de CPF duplicado")
	}
}

// TestCriarCliente_UsuarioDuplicado verifica que o mesmo usuario nao pode ter dois perfis.
func TestCriarCliente_UsuarioDuplicado(t *testing.T) {
	svc := novoClienteService()
	ctx := context.Background()
	usuarioID := uuid.New()

	if _, err := svc.Criar(ctx, usuarioID, domain.ClienteRequest{
		Nome: "Carlos Lima", CPF: "529.982.247-25", Telefone: "11911111111",
	}); err != nil {
		t.Fatalf("primeiro registro falhou: %v", err)
	}

	_, err := svc.Criar(ctx, usuarioID, domain.ClienteRequest{
		Nome: "Carlos Lima 2", CPF: "528.981.457-95", Telefone: "11922222222",
	})
	if err == nil {
		t.Fatal("esperava erro de cliente ja existente para o usuario")
	}
}

// TestAtualizarCliente_Sucesso verifica que nome e telefone sao atualizados.
func TestAtualizarCliente_Sucesso(t *testing.T) {
	svc := novoClienteService()
	ctx := context.Background()
	usuarioID := uuid.New()

	if _, err := svc.Criar(ctx, usuarioID, domain.ClienteRequest{
		Nome: "Pedro Alves", CPF: "529.982.247-25", Telefone: "11911111111",
	}); err != nil {
		t.Fatalf("criacao falhou: %v", err)
	}

	resp, err := svc.Atualizar(ctx, usuarioID, domain.ClienteRequest{
		Nome: "Pedro Alves Jr", Telefone: "11999999999",
	})
	if err != nil {
		t.Fatalf("esperava sucesso na atualizacao, got: %v", err)
	}
	if resp.Nome != "Pedro Alves Jr" {
		t.Errorf("nome esperado 'Pedro Alves Jr', got '%s'", resp.Nome)
	}
	if resp.Telefone != "11999999999" {
		t.Errorf("telefone esperado '11999999999', got '%s'", resp.Telefone)
	}
}

// TestAtualizarCliente_NaoEncontrado verifica erro quando cliente nao existe.
func TestAtualizarCliente_NaoEncontrado(t *testing.T) {
	svc := novoClienteService()
	ctx := context.Background()

	_, err := svc.Atualizar(ctx, uuid.New(), domain.ClienteRequest{
		Nome: "Ninguem", Telefone: "11900000000",
	})
	if err == nil {
		t.Fatal("esperava erro de cliente nao encontrado")
	}
}

// TestBuscarCliente_Sucesso verifica que o perfil e retornado corretamente.
func TestBuscarCliente_Sucesso(t *testing.T) {
	svc := novoClienteService()
	ctx := context.Background()
	usuarioID := uuid.New()

	if _, err := svc.Criar(ctx, usuarioID, domain.ClienteRequest{
		Nome: "Lucia Ferreira", CPF: "529.982.247-25", Telefone: "11911111111",
	}); err != nil {
		t.Fatalf("criacao falhou: %v", err)
	}

	resp, err := svc.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		t.Fatalf("esperava sucesso na busca, got: %v", err)
	}
	if resp.Nome != "Lucia Ferreira" {
		t.Errorf("nome esperado 'Lucia Ferreira', got '%s'", resp.Nome)
	}
}

// TestValidarCPF verifica casos do algoritmo de digito verificador.
func TestValidarCPF(t *testing.T) {
	casos := []struct {
		cpf    string
		valido bool
	}{
		{"52998224725", true},
		{"11111111111", false},  // todos iguais
		{"00000000000", false},  // todos zeros
		{"12345678901", false},  // digito incorreto
		{"1234567890", false},   // curto demais
		{"529982247251", false}, // longo demais
		{"529a8224725", false},  // nao numerico
	}

	for _, tc := range casos {
		err := domain.ValidarCPF(tc.cpf)
		if tc.valido && err != nil {
			t.Errorf("CPF %q deveria ser valido, got erro: %v", tc.cpf, err)
		}
		if !tc.valido && err == nil {
			t.Errorf("CPF %q deveria ser invalido, mas passou", tc.cpf)
		}
	}
}
