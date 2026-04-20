package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/google/uuid"
)

// setupSolicitacao monta cenário completo com cliente, endereço, catálogo e regiões seedados.
// Devolve svc, usuarioID do cliente principal, usuarioID de outro cliente,
// enderecoID (do primeiro cliente), categoriaID e regiaoID do seed,
// e uma função overrideAgora para controlar o clock nos testes.
func setupSolicitacao(t *testing.T) (
	svc *service.SolicitacaoService,
	clienteRepo *memory.ClienteRepository,
	usuarioID uuid.UUID,
	outroUsuarioID uuid.UUID,
	enderecoID uuid.UUID,
	categoriaID uuid.UUID,
	regiaoID uuid.UUID,
	overrideAgora func(time.Time),
) {
	t.Helper()
	ctx := context.Background()

	clienteRepo = memory.NewClienteRepository()
	enderecoRepo := memory.NewEnderecoRepository()
	regiaoRepo := memory.NewRegiaoRepository()
	catalogoRepo := memory.NewCatalogoRepository(regiaoRepo)
	solicRepo := memory.NewSolicitacaoRepository()

	cats, _ := catalogoRepo.ListarCategoriasAtivas(ctx)
	if len(cats) == 0 {
		t.Fatal("setup: nenhuma categoria seedada")
	}
	regioes, _ := regiaoRepo.ListarAtivas(ctx)
	if len(regioes) == 0 {
		t.Fatal("setup: nenhuma regiao seedada")
	}
	categoriaID = cats[0].ID
	regiaoID = regioes[0].ID

	// Cliente principal
	usuarioID = uuid.New()
	clienteID := uuid.New()
	c := &domain.Cliente{
		ID: clienteID, UsuarioID: usuarioID,
		Nome: "Ana Lima", CPF: "52998224725", Telefone: "62999991234", Score: 100,
	}
	if err := clienteRepo.Criar(ctx, c); err != nil {
		t.Fatalf("setup criar cliente: %v", err)
	}

	// Outro cliente para testes de ownership
	outroUsuarioID = uuid.New()
	outro := &domain.Cliente{
		ID: uuid.New(), UsuarioID: outroUsuarioID,
		Nome: "Carlos Ramos", CPF: "11144477735", Telefone: "62999994567", Score: 100,
	}
	if err := clienteRepo.Criar(ctx, outro); err != nil {
		t.Fatalf("setup criar outro cliente: %v", err)
	}

	// Endereço do cliente principal
	enderecoID = uuid.New()
	e := &domain.Endereco{
		ID:           enderecoID,
		ClienteID:    clienteID,
		Logradouro:   "Rua das Flores",
		Numero:       "10",
		Bairro:       "Setor Sul",
		Cidade:       "Goiania",
		Estado:       "GO",
		CEP:          "74125010",
		NumQuartos:   2,
		NumBanheiros: 1,
		NumSalas:     1,
		NumCozinhas:  1,
	}
	if err := enderecoRepo.Criar(ctx, e); err != nil {
		t.Fatalf("setup criar endereco: %v", err)
	}

	precSvc := service.NewPrecificacaoService(catalogoRepo, regiaoRepo)

	agoraPtr := time.Now()
	agoraFn := func() time.Time { return agoraPtr }
	svc = service.NewSolicitacaoServiceComAgora(solicRepo, enderecoRepo, clienteRepo, precSvc, agoraFn)
	overrideAgora = func(t time.Time) { agoraPtr = t }

	return
}

func novaReq(enderecoID, categoriaID, regiaoID uuid.UUID, dataServico time.Time) domain.CriarSolicitacaoRequest {
	return domain.CriarSolicitacaoRequest{
		EnderecoID:   enderecoID,
		CategoriaID:  categoriaID,
		RegiaoID:     regiaoID,
		NumQuartos:   1,
		NumBanheiros: 1,
		NumSalas:     1,
		NumCozinhas:  1,
		Frequencia:   domain.FrequenciaUnica,
		DataServico:  dataServico,
	}
}

// === Testes de criação ===

func TestSolicitacao_Criar_Sucesso_CongelaBreakdown(t *testing.T) {
	svc, _, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()
	agora := time.Now()
	override(agora)
	data := agora.Add(25 * time.Hour)

	resp, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data))
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	if resp.ValorTotal <= 0 {
		t.Errorf("valor total deveria ser > 0, got %.2f", resp.ValorTotal)
	}
	if len(resp.Breakdown) == 0 {
		t.Error("breakdown nao deveria estar vazio")
	}
	ultimo := resp.Breakdown[len(resp.Breakdown)-1]
	if ultimo.Tipo != "TOTAL" {
		t.Errorf("ultimo item do breakdown deveria ser TOTAL, got %s", ultimo.Tipo)
	}
	if ultimo.Valor != resp.ValorTotal {
		t.Errorf("item TOTAL (%.2f) != ValorTotal (%.2f)", ultimo.Valor, resp.ValorTotal)
	}
	if resp.Status != domain.StatusAguardando {
		t.Errorf("status deveria ser AGUARDANDO, got %s", resp.Status)
	}
}

func TestSolicitacao_Criar_AntecedenciaMenor24h_Bloqueia(t *testing.T) {
	svc, _, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()
	agora := time.Now()
	override(agora)
	data := agora.Add(23 * time.Hour)

	_, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data))
	if err != domain.ErrSolicitacaoAntecedencia {
		t.Errorf("esperava ErrSolicitacaoAntecedencia, got %v", err)
	}
}

func TestSolicitacao_Criar_AntecedenciaExata24h_Permite(t *testing.T) {
	svc, _, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()
	agora := time.Now()
	override(agora)
	data := agora.Add(24*time.Hour + time.Second) // borda acima de 24h

	_, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data))
	if err != nil {
		t.Errorf("esperava sucesso na borda de 24h, got %v", err)
	}
}

func TestSolicitacao_Criar_EnderecoDeOutroCliente_Rejeita(t *testing.T) {
	svc, _, _, outroUsuarioID, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()
	agora := time.Now()
	override(agora)
	data := agora.Add(25 * time.Hour)

	// Tenta criar solicitação com endereço do primeiro cliente, mas usando o token do segundo
	_, err := svc.Criar(ctx, outroUsuarioID, novaReq(endID, catID, regID, data))
	if err != domain.ErrSolicitacaoEnderecoInvalido {
		t.Errorf("esperava ErrSolicitacaoEnderecoInvalido, got %v", err)
	}
}

func TestSolicitacao_Criar_ObservacaoMaiorQue500_Rejeita(t *testing.T) {
	svc, _, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()
	agora := time.Now()
	override(agora)
	data := agora.Add(25 * time.Hour)

	req := novaReq(endID, catID, regID, data)
	req.Observacao = strings.Repeat("a", 501)

	_, err := svc.Criar(ctx, usuarioID, req)
	if err != domain.ErrObservacaoMuitoLonga {
		t.Errorf("esperava ErrObservacaoMuitoLonga, got %v", err)
	}
}

func TestSolicitacao_Criar_OpcionalInexistente_PropagaErroDePrecificacao(t *testing.T) {
	svc, _, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()
	agora := time.Now()
	override(agora)
	data := agora.Add(25 * time.Hour)

	req := novaReq(endID, catID, regID, data)
	req.OpcionaisIDs = []uuid.UUID{uuid.New()}

	_, err := svc.Criar(ctx, usuarioID, req)
	if err == nil {
		t.Error("esperava erro de opcional nao encontrado, got nil")
	}
}

// === Testes de cancelamento ===

func TestSolicitacao_Cancelar_MenosDe24h_PenalizaScore(t *testing.T) {
	svc, clienteRepo, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()

	// Criar solicitação para daqui 2h, mas fingindo que "agora" é 48h atrás
	agora := time.Now()
	dataServico := agora.Add(2 * time.Hour)
	override(agora.Add(-48 * time.Hour))

	resp, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, dataServico))
	if err != nil {
		t.Fatalf("criar: %v", err)
	}

	// Volta "agora" para o presente real (< 24h antes do serviço)
	override(agora)
	if err := svc.Cancelar(ctx, usuarioID, resp.ID); err != nil {
		t.Fatalf("cancelar: %v", err)
	}

	// Score deve ter caído 5 pontos (100 -> 95)
	c, _ := clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if c.Score != 95 {
		t.Errorf("score deveria ser 95, got %d", c.Score)
	}
}

func TestSolicitacao_Cancelar_MaisDe24h_NaoPenaliza(t *testing.T) {
	svc, clienteRepo, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()

	agora := time.Now()
	override(agora)
	data := agora.Add(48 * time.Hour)

	resp, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data))
	if err != nil {
		t.Fatalf("criar: %v", err)
	}

	// Cancela com 30h de antecedência (> 24h)
	override(data.Add(-30 * time.Hour))
	if err := svc.Cancelar(ctx, usuarioID, resp.ID); err != nil {
		t.Fatalf("cancelar: %v", err)
	}

	c, _ := clienteRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if c.Score != 100 {
		t.Errorf("score nao deveria ter mudado, got %d", c.Score)
	}
}

func TestSolicitacao_Cancelar_JaCancelada_Rejeita(t *testing.T) {
	svc, _, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()

	agora := time.Now()
	override(agora)
	data := agora.Add(48 * time.Hour)

	resp, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data))
	if err != nil {
		t.Fatalf("criar: %v", err)
	}

	if err := svc.Cancelar(ctx, usuarioID, resp.ID); err != nil {
		t.Fatalf("primeiro cancelar: %v", err)
	}

	// Segundo cancelamento deve falhar com status inválido
	err = svc.Cancelar(ctx, usuarioID, resp.ID)
	if err != domain.ErrSolicitacaoStatusInvalido {
		t.Errorf("esperava ErrSolicitacaoStatusInvalido, got %v", err)
	}
}

func TestSolicitacao_Cancelar_DeOutroCliente_RetornaNaoEncontrada(t *testing.T) {
	svc, _, usuarioID, outroUsuarioID, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()

	agora := time.Now()
	override(agora)
	data := agora.Add(48 * time.Hour)

	resp, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data))
	if err != nil {
		t.Fatalf("criar: %v", err)
	}

	err = svc.Cancelar(ctx, outroUsuarioID, resp.ID)
	if err != domain.ErrSolicitacaoNaoEncontrada {
		t.Errorf("esperava ErrSolicitacaoNaoEncontrada, got %v", err)
	}
}

// === Testes de busca ===

func TestSolicitacao_BuscarPorID_DeOutroCliente_RetornaNaoEncontrada(t *testing.T) {
	svc, _, usuarioID, outroUsuarioID, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()

	agora := time.Now()
	override(agora)
	data := agora.Add(48 * time.Hour)

	resp, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data))
	if err != nil {
		t.Fatalf("criar: %v", err)
	}

	_, err = svc.BuscarPorID(ctx, outroUsuarioID, resp.ID)
	if err != domain.ErrSolicitacaoNaoEncontrada {
		t.Errorf("esperava ErrSolicitacaoNaoEncontrada, got %v", err)
	}
}

// === Testes de listagem ===

func TestSolicitacao_ListarDoCliente_FiltroPorStatus(t *testing.T) {
	svc, _, usuarioID, _, endID, catID, regID, override := setupSolicitacao(t)
	ctx := context.Background()

	agora := time.Now()
	override(agora)

	// Criar 2 solicitações
	for i := 0; i < 2; i++ {
		data := agora.Add(time.Duration(48+i) * time.Hour)
		if _, err := svc.Criar(ctx, usuarioID, novaReq(endID, catID, regID, data)); err != nil {
			t.Fatalf("criar %d: %v", i, err)
		}
	}

	todas, _ := svc.ListarDoCliente(ctx, usuarioID, domain.SolicitacaoFiltro{})
	if len(todas) != 2 {
		t.Fatalf("esperava 2 solicitacoes, got %d", len(todas))
	}

	// Cancelar a mais recente (índice 0, ordem desc)
	if err := svc.Cancelar(ctx, usuarioID, todas[0].ID); err != nil {
		t.Fatalf("cancelar: %v", err)
	}

	status := domain.StatusAguardando
	aguardando, _ := svc.ListarDoCliente(ctx, usuarioID, domain.SolicitacaoFiltro{Status: &status})
	if len(aguardando) != 1 {
		t.Errorf("esperava 1 AGUARDANDO, got %d", len(aguardando))
	}

	cancelada := domain.StatusCancelada
	canceladas, _ := svc.ListarDoCliente(ctx, usuarioID, domain.SolicitacaoFiltro{Status: &cancelada})
	if len(canceladas) != 1 {
		t.Errorf("esperava 1 CANCELADA, got %d", len(canceladas))
	}
}
