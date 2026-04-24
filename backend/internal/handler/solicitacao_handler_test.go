package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/handler"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// setupSolicitacaoHandler monta router com /solicitacoes autenticado + dados seedados.
// Retorna router, token JWT do cliente autenticado e IDs necessários para criar solicitação.
func setupSolicitacaoHandler(t *testing.T) (
	r chi.Router,
	token string,
	enderecoID uuid.UUID,
	categoriaID uuid.UUID,
	regiaoID uuid.UUID,
) {
	t.Helper()
	ctx := context.Background()

	usuarioRepo := memory.NewUsuarioRepository()
	clienteRepo := memory.NewClienteRepository()
	enderecoRepo := memory.NewEnderecoRepository()
	regiaoRepo := memory.NewRegiaoRepository()
	catalogoRepo := memory.NewCatalogoRepository(regiaoRepo)
	solicRepo := memory.NewSolicitacaoRepository()

	authSvc := service.NewAuthServiceComCost(
		usuarioRepo, testSecret, 15*time.Minute, "development", service.BcryptCostTeste,
	)

	// Cria usuário + cliente
	regResp, err := authSvc.Registrar(ctx, domain.RegistroRequest{
		Email: "handler-test@diarygo.com.br", Senha: "Senha1234!",
	}, domain.TipoCliente)
	if err != nil {
		t.Fatalf("registrar cliente: %v", err)
	}

	cliente := &domain.Cliente{
		ID: uuid.New(), UsuarioID: regResp.ID,
		Nome: "Teste Handler", CPF: "52998224725",
		Telefone: "62999991234", Score: 100,
	}
	if err := clienteRepo.Criar(ctx, cliente); err != nil {
		t.Fatalf("criar cliente: %v", err)
	}

	// Endereço
	enderecoID = uuid.New()
	if err := enderecoRepo.Criar(ctx, &domain.Endereco{
		ID: enderecoID, ClienteID: cliente.ID,
		Logradouro: "Rua Teste", Numero: "100",
		Bairro: "Setor Sul", Cidade: "Goiania", Estado: "GO", CEP: "74125010",
		NumQuartos: 2, NumBanheiros: 1, NumSalas: 1, NumCozinhas: 1,
	}); err != nil {
		t.Fatalf("criar endereco: %v", err)
	}

	// Categoria e região do seed
	cats, _ := catalogoRepo.ListarCategoriasAtivas(ctx)
	if len(cats) == 0 {
		t.Fatal("nenhuma categoria seedada")
	}
	categoriaID = cats[0].ID
	regs, _ := regiaoRepo.ListarAtivas(ctx)
	if len(regs) == 0 {
		t.Fatal("nenhuma regiao seedada")
	}
	regiaoID = regs[0].ID

	// Token
	usr, err := usuarioRepo.BuscarPorID(ctx, regResp.ID)
	if err != nil {
		t.Fatalf("buscar usuario: %v", err)
	}
	tokResp, err := authSvc.GerarToken(usr)
	if err != nil {
		t.Fatalf("gerar token: %v", err)
	}
	token = tokResp.AccessToken

	// Router
	precSvc := service.NewPrecificacaoService(catalogoRepo, regiaoRepo)
	solicSvc := service.NewSolicitacaoService(solicRepo, enderecoRepo, clienteRepo, precSvc)
	h := handler.NewSolicitacaoHandler(solicSvc)

	mux := chi.NewRouter()
	mux.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Autenticar(authSvc))
			r.Use(mw.RequererTipo(domain.TipoCliente))
			r.Mount("/solicitacoes", h.Routes())
		})
	})
	r = mux
	return
}

// criarSolicitacaoHelper faz POST /solicitacoes e retorna o ID criado.
func criarSolicitacaoHelper(
	t *testing.T, r chi.Router, token string,
	endID, catID, regID uuid.UUID,
) uuid.UUID {
	t.Helper()
	req := map[string]interface{}{
		"endereco_id":   endID,
		"categoria_id":  catID,
		"regiao_id":     regID,
		"num_quartos":   1,
		"num_banheiros": 1,
		"num_salas":     1,
		"num_cozinhas":  1,
		"opcionais_ids": []string{},
		"frequencia":    "UNICA",
		"data_servico":  time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	}
	w := doRequest(r, http.MethodPost, "/api/v1/solicitacoes", req, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("criar solicitacao: esperava 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp domain.SolicitacaoResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode resposta: %v", err)
	}
	return resp.ID
}

// === DELETE /solicitacoes/{id} ===

// TestSolicitacaoHandler_Cancelar_Sucesso_Retorna204SemBody garante que o DELETE
// responde 204 No Content com body vazio — contrato que o frontend depende para
// não quebrar com SyntaxError ao tentar parsear JSON de uma resposta sem corpo.
func TestSolicitacaoHandler_Cancelar_Sucesso_Retorna204SemBody(t *testing.T) {
	r, token, endID, catID, regID := setupSolicitacaoHandler(t)
	solID := criarSolicitacaoHelper(t, r, token, endID, catID, regID)

	w := doRequest(r, http.MethodDelete, "/api/v1/solicitacoes/"+solID.String(), nil, token)

	if w.Code != http.StatusNoContent {
		t.Errorf("esperava 204, got %d: %s", w.Code, w.Body.String())
	}
	if body := w.Body.String(); body != "" {
		t.Errorf("esperava body vazio em 204, got %q", body)
	}
}

func TestSolicitacaoHandler_Cancelar_NaoExistente_Retorna404(t *testing.T) {
	r, token, _, _, _ := setupSolicitacaoHandler(t)

	w := doRequest(r, http.MethodDelete, "/api/v1/solicitacoes/"+uuid.New().String(), nil, token)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperava 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSolicitacaoHandler_Cancelar_IDInvalido_Retorna404(t *testing.T) {
	// UUID malformado é tratado como "não encontrado" (404) — mesma resposta
	// de um UUID válido que não existe. Isso é intencional: não diferenciar
	// "ID malformado" de "não existe" evita dar pistas ao atacante (A01).
	r, token, _, _, _ := setupSolicitacaoHandler(t)

	w := doRequest(r, http.MethodDelete, "/api/v1/solicitacoes/nao-eh-uuid", nil, token)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperava 404 em ID malformado, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSolicitacaoHandler_Cancelar_SemToken_Retorna401(t *testing.T) {
	r, _, _, _, _ := setupSolicitacaoHandler(t)

	w := doRequest(r, http.MethodDelete, "/api/v1/solicitacoes/"+uuid.New().String(), nil, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401 sem token, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSolicitacaoHandler_Cancelar_Idempotente_SegundaChamadaRetorna422(t *testing.T) {
	// Após cancelar uma solicitação, tentar cancelar de novo deve retornar 422
	// (status inválido — não está mais cancelável), e não 204. Isso garante
	// que o frontend saiba que a segunda chamada é inválida.
	r, token, endID, catID, regID := setupSolicitacaoHandler(t)
	solID := criarSolicitacaoHelper(t, r, token, endID, catID, regID)

	w1 := doRequest(r, http.MethodDelete, "/api/v1/solicitacoes/"+solID.String(), nil, token)
	if w1.Code != http.StatusNoContent {
		t.Fatalf("primeira: esperava 204, got %d", w1.Code)
	}

	w2 := doRequest(r, http.MethodDelete, "/api/v1/solicitacoes/"+solID.String(), nil, token)
	if w2.Code != http.StatusUnprocessableEntity {
		t.Errorf("segunda: esperava 422, got %d: %s", w2.Code, w2.Body.String())
	}
}

// === Smoke adicional: POST e GET ===

func TestSolicitacaoHandler_Criar_Sucesso_Retorna201(t *testing.T) {
	r, token, endID, catID, regID := setupSolicitacaoHandler(t)
	req := map[string]interface{}{
		"endereco_id":   endID,
		"categoria_id":  catID,
		"regiao_id":     regID,
		"num_quartos":   1,
		"num_banheiros": 1,
		"num_salas":     1,
		"num_cozinhas":  1,
		"opcionais_ids": []string{},
		"frequencia":    "UNICA",
		"data_servico":  time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	}

	w := doRequest(r, http.MethodPost, "/api/v1/solicitacoes", req, token)

	if w.Code != http.StatusCreated {
		t.Errorf("esperava 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSolicitacaoHandler_Listar_SemSolicitacoes_RetornaArrayVazio(t *testing.T) {
	r, token, _, _, _ := setupSolicitacaoHandler(t)

	w := doRequest(r, http.MethodGet, "/api/v1/solicitacoes", nil, token)

	if w.Code != http.StatusOK {
		t.Errorf("esperava 200, got %d: %s", w.Code, w.Body.String())
	}
	// Deve retornar [] e não null — frontend itera sobre array
	body := bytes.TrimSpace(w.Body.Bytes())
	if string(body) != "[]" {
		t.Errorf("esperava []', got %q", string(body))
	}
}
