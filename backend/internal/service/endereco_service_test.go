package service_test

import (
	"context"
	"testing"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/google/uuid"
)

func setupEnderecoService(t *testing.T) (*service.EnderecoService, *service.ClienteService, uuid.UUID) {
	t.Helper()
	clienteRepo := memory.NewClienteRepository()
	enderecoRepo := memory.NewEnderecoRepository()
	clienteSvc := service.NewClienteService(clienteRepo)
	enderecoSvc := service.NewEnderecoService(enderecoRepo, clienteRepo)

	usuarioID := uuid.New()
	_, err := clienteSvc.Criar(context.Background(), usuarioID, domain.ClienteRequest{
		Nome: "Teste Usuario", CPF: "529.982.247-25", Telefone: "11911111111",
	})
	if err != nil {
		t.Fatalf("falha ao criar cliente de teste: %v", err)
	}
	return enderecoSvc, clienteSvc, usuarioID
}

func enderecoValido() domain.EnderecoRequest {
	return domain.EnderecoRequest{
		Logradouro:   "Rua das Flores",
		Numero:       "123",
		Bairro:       "Centro",
		Cidade:       "São Paulo",
		Estado:       "SP",
		CEP:          "01310-100",
		NumQuartos:   2,
		NumBanheiros: 1,
		NumSalas:     1,
		NumCozinhas:  1,
	}
}

// TestCriarEndereco_Sucesso verifica criacao de endereco valido.
func TestCriarEndereco_Sucesso(t *testing.T) {
	svc, _, usuarioID := setupEnderecoService(t)
	ctx := context.Background()

	resp, err := svc.Criar(ctx, usuarioID, enderecoValido())
	if err != nil {
		t.Fatalf("esperava sucesso, got: %v", err)
	}
	if resp.Logradouro != "Rua das Flores" {
		t.Errorf("logradouro incorreto: %s", resp.Logradouro)
	}
	if resp.CEP != "01310100" {
		t.Errorf("CEP esperado sem mascara '01310100', got '%s'", resp.CEP)
	}
	if resp.Estado != "SP" {
		t.Errorf("estado esperado 'SP', got '%s'", resp.Estado)
	}
}

// TestCriarEndereco_CEPInvalido verifica que CEP invalido e rejeitado.
func TestCriarEndereco_CEPInvalido(t *testing.T) {
	svc, _, usuarioID := setupEnderecoService(t)
	req := enderecoValido()
	req.CEP = "12345"

	_, err := svc.Criar(context.Background(), usuarioID, req)
	if err == nil {
		t.Fatal("esperava erro de CEP invalido")
	}
}

// TestCriarEndereco_EstadoInvalido verifica que UF invalida e rejeitada.
func TestCriarEndereco_EstadoInvalido(t *testing.T) {
	svc, _, usuarioID := setupEnderecoService(t)
	req := enderecoValido()
	req.Estado = "XX"

	_, err := svc.Criar(context.Background(), usuarioID, req)
	if err == nil {
		t.Fatal("esperava erro de estado invalido")
	}
}

// TestCriarEndereco_QuartosZero verifica que numero de quartos zero e rejeitado.
func TestCriarEndereco_QuartosZero(t *testing.T) {
	svc, _, usuarioID := setupEnderecoService(t)
	req := enderecoValido()
	req.NumQuartos = 0

	_, err := svc.Criar(context.Background(), usuarioID, req)
	if err == nil {
		t.Fatal("esperava erro de comodo invalido")
	}
}

// TestListarEnderecos_Sucesso verifica listagem de multiplos enderecos.
func TestListarEnderecos_Sucesso(t *testing.T) {
	svc, _, usuarioID := setupEnderecoService(t)
	ctx := context.Background()

	if _, err := svc.Criar(ctx, usuarioID, enderecoValido()); err != nil {
		t.Fatalf("criacao 1 falhou: %v", err)
	}
	req2 := enderecoValido()
	req2.Logradouro = "Av. Paulista"
	if _, err := svc.Criar(ctx, usuarioID, req2); err != nil {
		t.Fatalf("criacao 2 falhou: %v", err)
	}

	lista, err := svc.Listar(ctx, usuarioID)
	if err != nil {
		t.Fatalf("esperava sucesso na listagem, got: %v", err)
	}
	if len(lista) != 2 {
		t.Errorf("esperava 2 enderecos, got %d", len(lista))
	}
}

// TestDefinirPrincipal_Sucesso verifica que apenas um endereco fica como principal.
func TestDefinirPrincipal_Sucesso(t *testing.T) {
	svc, _, usuarioID := setupEnderecoService(t)
	ctx := context.Background()

	e1, _ := svc.Criar(ctx, usuarioID, enderecoValido())
	req2 := enderecoValido()
	req2.Logradouro = "Av. Brasil"
	e2, _ := svc.Criar(ctx, usuarioID, req2)

	// Definir e1 como principal
	_, err := svc.DefinirPrincipal(ctx, usuarioID, e1.ID)
	if err != nil {
		t.Fatalf("definir principal falhou: %v", err)
	}

	// Definir e2 como principal — deve desmarcar e1
	_, err = svc.DefinirPrincipal(ctx, usuarioID, e2.ID)
	if err != nil {
		t.Fatalf("trocar principal falhou: %v", err)
	}

	lista, _ := svc.Listar(ctx, usuarioID)
	principaisCount := 0
	for _, e := range lista {
		if e.Principal {
			principaisCount++
		}
	}
	if principaisCount != 1 {
		t.Errorf("esperava exatamente 1 endereco principal, got %d", principaisCount)
	}
}

// TestRemoverEndereco_Sucesso verifica que endereco e removido da lista.
func TestRemoverEndereco_Sucesso(t *testing.T) {
	svc, _, usuarioID := setupEnderecoService(t)
	ctx := context.Background()

	e, _ := svc.Criar(ctx, usuarioID, enderecoValido())

	if err := svc.Remover(ctx, usuarioID, e.ID); err != nil {
		t.Fatalf("remover falhou: %v", err)
	}

	lista, _ := svc.Listar(ctx, usuarioID)
	if len(lista) != 0 {
		t.Errorf("esperava lista vazia apos remocao, got %d enderecos", len(lista))
	}
}

// TestRemoverEndereco_NaoPertenceAoCliente verifica protecao de ownership.
func TestRemoverEndereco_NaoPertenceAoCliente(t *testing.T) {
	clienteRepo := memory.NewClienteRepository()
	enderecoRepo := memory.NewEnderecoRepository()
	clienteSvc := service.NewClienteService(clienteRepo)
	svc := service.NewEnderecoService(enderecoRepo, clienteRepo)
	ctx := context.Background()

	// Criar dois usuarios com clientes distintos
	u1 := uuid.New()
	u2 := uuid.New()
	clienteSvc.Criar(ctx, u1, domain.ClienteRequest{Nome: "User1", CPF: "529.982.247-25", Telefone: "11911111111"}) //nolint:errcheck
	clienteSvc.Criar(ctx, u2, domain.ClienteRequest{Nome: "User2", CPF: "528.981.457-95", Telefone: "11922222222"}) //nolint:errcheck

	e1, _ := svc.Criar(ctx, u1, enderecoValido())

	// u2 tenta remover o endereco de u1 — deve falhar
	err := svc.Remover(ctx, u2, e1.ID)
	if err == nil {
		t.Fatal("esperava erro de ownership, mas removeu endereco de outro cliente")
	}
}
