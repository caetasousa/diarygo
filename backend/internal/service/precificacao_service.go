package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// Tipos de ItemCalculo no breakdown — constantes para consumo do frontend.
const (
	ItemBASE      = "BASE"
	ItemCOMODO    = "COMODO"
	ItemOPCIONAL  = "OPCIONAL"
	ItemACRESCIMO = "ACRESCIMO"
	ItemDESCONTO  = "DESCONTO"
	ItemTOTAL     = "TOTAL"
)

// Constantes de negocio — expostas para testes e ajuste fino pelo admin.
const (
	TempoExtraPorQuarto   = 20 // min adicionais por quarto alem do primeiro
	TempoExtraPorBanheiro = 20
	TempoExtraPorSala     = 15
	TempoExtraPorCozinha  = 25
)

// PrecificacaoService calcula orcamentos com breakdown transparente.
type PrecificacaoService struct {
	catalogo domain.CatalogoRepository
	regiao   domain.RegiaoRepository
}

// NewPrecificacaoService cria um novo PrecificacaoService.
func NewPrecificacaoService(catalogo domain.CatalogoRepository, regiao domain.RegiaoRepository) *PrecificacaoService {
	return &PrecificacaoService{catalogo: catalogo, regiao: regiao}
}

// ListarCategorias retorna todas as categorias ativas em DTO.
func (s *PrecificacaoService) ListarCategorias(ctx context.Context) ([]*domain.CategoriaResponse, error) {
	cats, err := s.catalogo.ListarCategoriasAtivas(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*domain.CategoriaResponse, 0, len(cats))
	for _, c := range cats {
		res = append(res, &domain.CategoriaResponse{
			ID:               c.ID,
			Nome:             c.Nome,
			Descricao:        c.Descricao,
			DuracaoMinimaMin: c.DuracaoMinimaMin,
			Ativa:            c.Ativa,
		})
	}
	return res, nil
}

// ListarOpcionaisDaCategoria retorna opcionais disponiveis.
// No MVP todos sao compartilhados — requer categoria valida para espelhar
// o contrato futuro (categoria-opcional N:N).
func (s *PrecificacaoService) ListarOpcionaisDaCategoria(ctx context.Context, categoriaID uuid.UUID) ([]*domain.OpcionalResponse, error) {
	if _, err := s.catalogo.BuscarCategoriaPorID(ctx, categoriaID); err != nil {
		return nil, err
	}
	ops, err := s.catalogo.ListarOpcionaisAtivos(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*domain.OpcionalResponse, 0, len(ops))
	for _, o := range ops {
		res = append(res, &domain.OpcionalResponse{
			ID:            o.ID,
			Nome:          o.Nome,
			Descricao:     o.Descricao,
			ValorExtra:    o.ValorExtra,
			TempoExtraMin: o.TempoExtraMin,
			Ativo:         o.Ativo,
		})
	}
	return res, nil
}

// CalcularValorReferencia produz o orcamento com breakdown completo.
// Regras:
//   - duracao_min = max(duracao_minima_categoria, sum(tempos por comodo) + sum(opcionais.tempo))
//   - valor_base = preco_hora * horas (proporcional a duracao)
//   - opcionais somam valor e tempo
//   - fim-de-semana aplica acrescimo (+%) apos opcionais
//   - frequencia aplica desconto (-%) no final
//
// O breakdown lista cada item na ordem: BASE → COMODO → OPCIONAIS → ACRESCIMO → DESCONTO → TOTAL.
func (s *PrecificacaoService) CalcularValorReferencia(ctx context.Context, req domain.CalculoPrecoRequest) (*domain.CalculoPrecoResponse, error) {
	if req.CategoriaID == uuid.Nil {
		return nil, domain.ErrCategoriaNaoEncontrada
	}
	if req.RegiaoID == uuid.Nil {
		return nil, domain.ErrRegiaoNaoEncontrada
	}
	if req.NumQuartos < 1 {
		return nil, domain.ErrComodosInvalidosCalculo
	}
	if err := domain.ValidarComodos(req.NumQuartos, req.NumBanheiros, req.NumSalas, req.NumCozinhas); err != nil {
		return nil, err
	}
	if req.Frequencia != "" && !req.Frequencia.Valida() {
		return nil, domain.ErrFrequenciaInvalida
	}
	// Data no passado so rejeita se foi explicitamente informada (zero value = aceita).
	if !req.DataServico.IsZero() && req.DataServico.Before(time.Now().Truncate(24*time.Hour)) {
		return nil, domain.ErrDataServicoPassado
	}

	cat, err := s.catalogo.BuscarCategoriaPorID(ctx, req.CategoriaID)
	if err != nil {
		return nil, err
	}
	tab, err := s.catalogo.BuscarTabelaPreco(ctx, req.CategoriaID, req.RegiaoID)
	if err != nil {
		return nil, err
	}

	itens := make([]domain.ItemCalculo, 0, 8)

	// Duracao baseline = duracao minima da categoria
	duracaoMin := cat.DuracaoMinimaMin

	// Ajuste por comodos alem do 1o quarto (primeiro quarto ja esta na baseline).
	extraQuartos := (req.NumQuartos - 1) * TempoExtraPorQuarto
	extraBanheiros := req.NumBanheiros * TempoExtraPorBanheiro
	extraSalas := req.NumSalas * TempoExtraPorSala
	extraCozinhas := req.NumCozinhas * TempoExtraPorCozinha
	duracaoMin += extraQuartos + extraBanheiros + extraSalas + extraCozinhas

	// Opcionais somam tempo e valor.
	opcionaisValor := 0.0
	opcionaisItens := make([]domain.ItemCalculo, 0, len(req.OpcionaisIDs))
	for _, oid := range req.OpcionaisIDs {
		op, err := s.catalogo.BuscarOpcionalPorID(ctx, oid)
		if err != nil {
			return nil, err
		}
		duracaoMin += op.TempoExtraMin
		opcionaisValor += op.ValorExtra
		opcionaisItens = append(opcionaisItens, domain.ItemCalculo{
			Tipo:  ItemOPCIONAL,
			Label: op.Nome,
			Valor: op.ValorExtra,
		})
	}

	// Calculo de valor base proporcional a duracao em horas.
	horas := float64(duracaoMin) / 60.0
	valorBase := tab.PrecoHora * horas

	itens = append(itens, domain.ItemCalculo{
		Tipo:  ItemBASE,
		Label: fmt.Sprintf("Base %s (%.1fh x R$ %.2f/h)", cat.Nome, horas, tab.PrecoHora),
		Valor: round2(valorBase),
	})

	if extraQuartos+extraBanheiros+extraSalas+extraCozinhas > 0 {
		itens = append(itens, domain.ItemCalculo{
			Tipo:  ItemCOMODO,
			Label: fmt.Sprintf("Comodos (%d quartos, %d banh., %d salas, %d cozinhas)", req.NumQuartos, req.NumBanheiros, req.NumSalas, req.NumCozinhas),
			Valor: 0, // ja embutido na duracao; item informativo
		})
	}

	itens = append(itens, opcionaisItens...)

	subtotal := valorBase + opcionaisValor

	// Acrescimo fim-de-semana.
	if !req.DataServico.IsZero() {
		wd := req.DataServico.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			acrescimo := subtotal * (tab.AcrescimoFDS / 100.0)
			itens = append(itens, domain.ItemCalculo{
				Tipo:  ItemACRESCIMO,
				Label: fmt.Sprintf("Acrescimo fim-de-semana (+%.0f%%)", tab.AcrescimoFDS),
				Valor: round2(acrescimo),
			})
			subtotal += acrescimo
		}
	}

	// Desconto por frequencia.
	if req.Frequencia != "" && req.Frequencia != domain.FrequenciaUnica {
		pct := 0.0
		label := ""
		switch req.Frequencia {
		case domain.FrequenciaSemanal:
			pct = tab.DescontoSemanal
			label = "semanal"
		case domain.FrequenciaQuinzenal:
			pct = tab.DescontoQuinzenal
			label = "quinzenal"
		case domain.FrequenciaDuasPorSemana:
			pct = tab.DescontoDuasSemana
			label = "2x por semana"
		}
		if pct > 0 {
			desc := subtotal * (pct / 100.0)
			itens = append(itens, domain.ItemCalculo{
				Tipo:  ItemDESCONTO,
				Label: fmt.Sprintf("Desconto %s (-%.0f%%)", label, pct),
				Valor: -round2(desc),
			})
			subtotal -= desc
		}
	}

	total := round2(subtotal)
	itens = append(itens, domain.ItemCalculo{
		Tipo:  ItemTOTAL,
		Label: "Total estimado",
		Valor: total,
	})

	return &domain.CalculoPrecoResponse{
		CategoriaID: req.CategoriaID,
		RegiaoID:    req.RegiaoID,
		DuracaoMin:  duracaoMin,
		ValorTotal:  total,
		Itens:       itens,
	}, nil
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100.0
}

// sanityCheckRegion garante que a regiao existe — util para chamadas publicas.
func (s *PrecificacaoService) sanityCheckRegion(ctx context.Context, id uuid.UUID) error {
	_, err := s.regiao.BuscarPorID(ctx, id)
	if err != nil && !errors.Is(err, domain.ErrRegiaoNaoEncontrada) {
		return err
	}
	if err != nil {
		return domain.ErrRegiaoNaoEncontrada
	}
	return nil
}
