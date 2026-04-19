package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/google/uuid"
)

func setupPrecificacao(t *testing.T) (*service.PrecificacaoService, uuid.UUID, uuid.UUID) {
	t.Helper()
	regiaoRepo := memory.NewRegiaoRepository()
	catalogoRepo := memory.NewCatalogoRepository(regiaoRepo)
	svc := service.NewPrecificacaoService(catalogoRepo, regiaoRepo)

	regioes, _ := regiaoRepo.ListarAtivas(context.Background())
	if len(regioes) == 0 {
		t.Fatalf("setup: nenhuma regiao seedada")
	}

	cats, _ := catalogoRepo.ListarCategoriasAtivas(context.Background())
	if len(cats) == 0 {
		t.Fatalf("setup: nenhuma categoria seedada")
	}

	// Pega a primeira regiao e categoria deterministicamente
	return svc, cats[0].ID, regioes[0].ID
}

func TestListarCategorias_RetornaSete(t *testing.T) {
	svc, _, _ := setupPrecificacao(t)
	cats, err := svc.ListarCategorias(context.Background())
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if len(cats) != 7 {
		t.Errorf("esperava 7 categorias, got %d", len(cats))
	}
}

func TestListarOpcionais_ExigeCategoriaValida(t *testing.T) {
	svc, _, _ := setupPrecificacao(t)
	_, err := svc.ListarOpcionaisDaCategoria(context.Background(), uuid.New())
	if err == nil || err != domain.ErrCategoriaNaoEncontrada {
		t.Errorf("esperava ErrCategoriaNaoEncontrada, got %v", err)
	}
}

func TestListarOpcionais_RetornaSete(t *testing.T) {
	svc, catID, _ := setupPrecificacao(t)
	ops, err := svc.ListarOpcionaisDaCategoria(context.Background(), catID)
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if len(ops) != 7 {
		t.Errorf("esperava 7 opcionais, got %d", len(ops))
	}
}

func TestCalcular_Basico_RetornaBreakdown(t *testing.T) {
	svc, catID, regID := setupPrecificacao(t)
	resp, err := svc.CalcularValorReferencia(context.Background(), domain.CalculoPrecoRequest{
		CategoriaID:  catID,
		RegiaoID:     regID,
		NumQuartos:   1,
		NumBanheiros: 1,
		Frequencia:   domain.FrequenciaUnica,
	})
	if err != nil {
		t.Fatalf("calcular: %v", err)
	}
	if resp.ValorTotal <= 0 {
		t.Errorf("valor total deveria ser > 0, got %.2f", resp.ValorTotal)
	}
	if resp.DuracaoMin <= 0 {
		t.Errorf("duracao deveria ser > 0, got %d", resp.DuracaoMin)
	}
	// Breakdown deve sempre terminar em TOTAL
	if resp.Itens[len(resp.Itens)-1].Tipo != service.ItemTOTAL {
		t.Errorf("ultimo item do breakdown deve ser TOTAL, got %s", resp.Itens[len(resp.Itens)-1].Tipo)
	}
}

func TestCalcular_CategoriaInvalida(t *testing.T) {
	svc, _, regID := setupPrecificacao(t)
	_, err := svc.CalcularValorReferencia(context.Background(), domain.CalculoPrecoRequest{
		CategoriaID: uuid.New(),
		RegiaoID:    regID,
		NumQuartos:  1,
	})
	if err != domain.ErrCategoriaNaoEncontrada {
		t.Errorf("esperava ErrCategoriaNaoEncontrada, got %v", err)
	}
}

func TestCalcular_RegiaoInvalida(t *testing.T) {
	svc, catID, _ := setupPrecificacao(t)
	_, err := svc.CalcularValorReferencia(context.Background(), domain.CalculoPrecoRequest{
		CategoriaID: catID,
		RegiaoID:    uuid.New(),
		NumQuartos:  1,
	})
	if err != domain.ErrTabelaPrecoNaoEncontrada {
		t.Errorf("esperava ErrTabelaPrecoNaoEncontrada, got %v", err)
	}
}

func TestCalcular_NumQuartosZero_Rejeita(t *testing.T) {
	svc, catID, regID := setupPrecificacao(t)
	_, err := svc.CalcularValorReferencia(context.Background(), domain.CalculoPrecoRequest{
		CategoriaID: catID,
		RegiaoID:    regID,
		NumQuartos:  0,
	})
	if err != domain.ErrComodosInvalidosCalculo {
		t.Errorf("esperava ErrComodosInvalidosCalculo, got %v", err)
	}
}

func TestCalcular_FrequenciaSemanal_AplicaDesconto(t *testing.T) {
	svc, catID, regID := setupPrecificacao(t)
	req := domain.CalculoPrecoRequest{
		CategoriaID:  catID,
		RegiaoID:     regID,
		NumQuartos:   2,
		NumBanheiros: 1,
	}

	req.Frequencia = domain.FrequenciaUnica
	unica, err := svc.CalcularValorReferencia(context.Background(), req)
	if err != nil {
		t.Fatalf("unica: %v", err)
	}

	req.Frequencia = domain.FrequenciaSemanal
	semanal, err := svc.CalcularValorReferencia(context.Background(), req)
	if err != nil {
		t.Fatalf("semanal: %v", err)
	}

	if semanal.ValorTotal >= unica.ValorTotal {
		t.Errorf("semanal (%.2f) deveria ser < unica (%.2f)", semanal.ValorTotal, unica.ValorTotal)
	}

	// Deve existir linha de DESCONTO no breakdown
	achou := false
	for _, it := range semanal.Itens {
		if it.Tipo == service.ItemDESCONTO {
			achou = true
			if it.Valor >= 0 {
				t.Errorf("valor do desconto deveria ser negativo, got %.2f", it.Valor)
			}
		}
	}
	if !achou {
		t.Errorf("breakdown semanal deveria ter linha DESCONTO, got %+v", semanal.Itens)
	}
}

func TestCalcular_FimDeSemana_AplicaAcrescimo(t *testing.T) {
	svc, catID, regID := setupPrecificacao(t)

	// Encontrar proximo sabado futuro
	agora := time.Now()
	dias := int(time.Saturday - agora.Weekday())
	if dias <= 0 {
		dias += 7
	}
	sabado := agora.AddDate(0, 0, dias)

	req := domain.CalculoPrecoRequest{
		CategoriaID:  catID,
		RegiaoID:     regID,
		NumQuartos:   1,
		NumBanheiros: 1,
		DataServico:  sabado,
		Frequencia:   domain.FrequenciaUnica,
	}
	resp, err := svc.CalcularValorReferencia(context.Background(), req)
	if err != nil {
		t.Fatalf("calcular: %v", err)
	}

	achou := false
	for _, it := range resp.Itens {
		if it.Tipo == service.ItemACRESCIMO {
			achou = true
			if it.Valor <= 0 {
				t.Errorf("acrescimo FDS deveria ser positivo, got %.2f", it.Valor)
			}
		}
	}
	if !achou {
		t.Errorf("breakdown em sabado deveria ter ACRESCIMO, got %+v", resp.Itens)
	}
}

func TestCalcular_DataPassado_Rejeita(t *testing.T) {
	svc, catID, regID := setupPrecificacao(t)
	req := domain.CalculoPrecoRequest{
		CategoriaID: catID,
		RegiaoID:    regID,
		NumQuartos:  1,
		DataServico: time.Now().AddDate(0, 0, -10),
	}
	_, err := svc.CalcularValorReferencia(context.Background(), req)
	if err != domain.ErrDataServicoPassado {
		t.Errorf("esperava ErrDataServicoPassado, got %v", err)
	}
}

func TestCalcular_OpcionalSomaValorETempo(t *testing.T) {
	regiaoRepo := memory.NewRegiaoRepository()
	catalogoRepo := memory.NewCatalogoRepository(regiaoRepo)
	svc := service.NewPrecificacaoService(catalogoRepo, regiaoRepo)
	ctx := context.Background()

	regioes, _ := regiaoRepo.ListarAtivas(ctx)
	cats, _ := catalogoRepo.ListarCategoriasAtivas(ctx)
	ops, _ := catalogoRepo.ListarOpcionaisAtivos(ctx)
	if len(ops) == 0 {
		t.Fatal("sem opcionais seedados")
	}

	req := domain.CalculoPrecoRequest{
		CategoriaID:  cats[0].ID,
		RegiaoID:     regioes[0].ID,
		NumQuartos:   1,
		NumBanheiros: 1,
	}
	semOpcional, _ := svc.CalcularValorReferencia(ctx, req)

	req.OpcionaisIDs = []uuid.UUID{ops[0].ID}
	comOpcional, _ := svc.CalcularValorReferencia(ctx, req)

	if comOpcional.ValorTotal <= semOpcional.ValorTotal {
		t.Errorf("com opcional (%.2f) deveria ser > sem (%.2f)", comOpcional.ValorTotal, semOpcional.ValorTotal)
	}
	if comOpcional.DuracaoMin < semOpcional.DuracaoMin {
		t.Errorf("duracao com opcional (%d) deveria ser >= sem (%d)", comOpcional.DuracaoMin, semOpcional.DuracaoMin)
	}
}

func TestCalcular_TodasRegioesTemTabela(t *testing.T) {
	regiaoRepo := memory.NewRegiaoRepository()
	catalogoRepo := memory.NewCatalogoRepository(regiaoRepo)
	svc := service.NewPrecificacaoService(catalogoRepo, regiaoRepo)
	ctx := context.Background()

	regioes, _ := regiaoRepo.ListarAtivas(ctx)
	cats, _ := catalogoRepo.ListarCategoriasAtivas(ctx)

	for _, c := range cats {
		for _, r := range regioes {
			_, err := svc.CalcularValorReferencia(ctx, domain.CalculoPrecoRequest{
				CategoriaID:  c.ID,
				RegiaoID:     r.ID,
				NumQuartos:   1,
				NumBanheiros: 1,
			})
			if err != nil {
				t.Errorf("calcular(%s, %s): %v", c.Nome, r.Nome, err)
			}
		}
	}
}

func TestCalcular_FrequenciaInvalida(t *testing.T) {
	svc, catID, regID := setupPrecificacao(t)
	_, err := svc.CalcularValorReferencia(context.Background(), domain.CalculoPrecoRequest{
		CategoriaID: catID,
		RegiaoID:    regID,
		NumQuartos:  1,
		Frequencia:  domain.FrequenciaServico("MENSAL"),
	})
	if err != domain.ErrFrequenciaInvalida {
		t.Errorf("esperava ErrFrequenciaInvalida, got %v", err)
	}
}
