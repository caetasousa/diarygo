package memory

import (
	"context"
	"sync"
	"time"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/google/uuid"
)

// CatalogoRepository implementa domain.CatalogoRepository em memoria.
// Faz o seed de categorias, opcionais e tabela de precos por (categoria x regiao)
// usando as regioes carregadas do RegiaoRepository injetado.
type CatalogoRepository struct {
	mu         sync.RWMutex
	categorias map[uuid.UUID]*domain.CategoriaServico
	opcionais  map[uuid.UUID]*domain.Opcional
	tabelas    map[string]*domain.TabelaPrecos // chave: "categoriaID|regiaoID"
	regiaoRepo *RegiaoRepository
}

// NewCatalogoRepository cria o repositorio in-memory e aplica o seed do MVP.
func NewCatalogoRepository(regiaoRepo *RegiaoRepository) *CatalogoRepository {
	r := &CatalogoRepository{
		categorias: make(map[uuid.UUID]*domain.CategoriaServico),
		opcionais:  make(map[uuid.UUID]*domain.Opcional),
		tabelas:    make(map[string]*domain.TabelaPrecos),
		regiaoRepo: regiaoRepo,
	}
	r.seed()
	return r
}

func tabelaKey(categoriaID, regiaoID uuid.UUID) string {
	return categoriaID.String() + "|" + regiaoID.String()
}

// seed carrega 7 categorias, 7 opcionais e uma tabela de precos para cada
// combinacao categoria x regiao ativa. Valores-base baseados em pesquisa de
// mercado local (Goiania) — podem ser ajustados pelo admin na Etapa 10.
func (r *CatalogoRepository) seed() {
	now := time.Now()

	type catSeed struct {
		Nome             string
		Descricao        string
		DuracaoMinimaMin int
		PrecoHoraBase    float64 // preco-hora de referencia (regiao Central)
	}
	categorias := []catSeed{
		{"Faxina Padrao", "Limpeza geral de manutencao da casa.", 180, 28.00},
		{"Faxina Pesada", "Limpeza profunda com foco em gordura, mofo e areas criticas.", 300, 35.00},
		{"Pos-obra", "Limpeza apos reforma: poeira fina, respingos de tinta, entulho.", 360, 42.00},
		{"Organizacao e Closet", "Ordenacao de armarios, roupas, gavetas e ambientes.", 240, 32.00},
		{"Passadoria", "Passar roupas acumuladas, com ou sem dobra.", 120, 26.00},
		{"Limpeza de Janelas", "Janelas, vidros, box e espelhos.", 90, 30.00},
		{"Limpeza de Cozinha", "Foco exclusivo em cozinha: fogao, coifa, geladeira, pias.", 120, 30.00},
	}

	type opSeed struct {
		Nome          string
		Descricao     string
		ValorExtra    float64
		TempoExtraMin int
	}
	opcionais := []opSeed{
		{"Passadoria extra", "Adicionar passar roupas alem do escopo.", 25.00, 45},
		{"Limpeza de geladeira", "Retirar, higienizar e recolocar itens.", 20.00, 30},
		{"Limpeza de fogao e coifa", "Desengordurar fogao, coifa e filtro.", 25.00, 40},
		{"Janelas externas", "Limpeza de janelas pelo lado de fora.", 30.00, 30},
		{"Limpeza de area externa", "Varanda, quintal, area de servico.", 20.00, 30},
		{"Lavagem de roupa de cama", "Trocar e lavar roupa de cama.", 18.00, 30},
		{"Produtos de limpeza", "Profissional leva kit de produtos.", 25.00, 0},
	}

	for _, c := range categorias {
		id := uuid.New()
		r.categorias[id] = &domain.CategoriaServico{
			ID:               id,
			Nome:             c.Nome,
			Descricao:        c.Descricao,
			DuracaoMinimaMin: c.DuracaoMinimaMin,
			Ativa:            true,
			CriadaEm:         now,
			AtualizadaEm:     now,
		}
	}
	for _, o := range opcionais {
		id := uuid.New()
		r.opcionais[id] = &domain.Opcional{
			ID:            id,
			Nome:          o.Nome,
			Descricao:     o.Descricao,
			ValorExtra:    o.ValorExtra,
			TempoExtraMin: o.TempoExtraMin,
			Ativo:         true,
			CriadoEm:      now,
			AtualizadoEm:  now,
		}
	}

	// Seed da tabela de precos: preco_hora da categoria ajustado por regiao.
	// Multiplicador por regiao reflete diferenca de perfil socioeconomico.
	regiaoMultiplicador := map[string]float64{
		"Região Central":            1.10,
		"Região Sul":                1.15,
		"Região Leste":              1.05,
		"Região Oeste":              1.00,
		"Região Sudoeste":           0.95,
		"Região Sudeste":            0.98,
		"Região Norte":              0.95,
		"Região Noroeste":           0.92,
		"Região Campinas-Centro":    1.02,
		"Região Macambira":          0.98,
		"Região Vale do Meia Ponte": 0.95,
		"Região Mendanha":           0.92,
	}

	regioes, _ := r.regiaoRepo.ListarAtivas(context.Background())
	for _, cat := range r.categorias {
		// Recuperar preco-hora base pela correspondencia de nome (facil nos testes)
		var base float64
		for _, c := range categorias {
			if c.Nome == cat.Nome {
				base = c.PrecoHoraBase
				break
			}
		}
		for _, reg := range regioes {
			mult, ok := regiaoMultiplicador[reg.Nome]
			if !ok {
				mult = 1.0
			}
			tab := &domain.TabelaPrecos{
				ID:                 uuid.New(),
				CategoriaID:        cat.ID,
				RegiaoID:           reg.ID,
				PrecoHora:          round2(base * mult),
				AcrescimoFDS:       15.0,
				DescontoSemanal:    10.0,
				DescontoQuinzenal:  5.0,
				DescontoDuasSemana: 15.0,
				Ativa:              true,
				CriadaEm:           now,
				AtualizadaEm:       now,
			}
			r.tabelas[tabelaKey(cat.ID, reg.ID)] = tab
		}
	}
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100.0
}

// ListarCategoriasAtivas retorna copias das categorias ativas.
func (r *CatalogoRepository) ListarCategoriasAtivas(ctx context.Context) ([]*domain.CategoriaServico, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.CategoriaServico, 0, len(r.categorias))
	for _, c := range r.categorias {
		if c.Ativa {
			copia := *c
			res = append(res, &copia)
		}
	}
	return res, nil
}

// BuscarCategoriaPorID retorna uma copia da categoria.
func (r *CatalogoRepository) BuscarCategoriaPorID(ctx context.Context, id uuid.UUID) (*domain.CategoriaServico, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.categorias[id]
	if !ok {
		return nil, domain.ErrCategoriaNaoEncontrada
	}
	copia := *c
	return &copia, nil
}

// ListarOpcionaisAtivos retorna copias dos opcionais ativos.
// No MVP, opcionais sao compartilhados entre todas as categorias.
func (r *CatalogoRepository) ListarOpcionaisAtivos(ctx context.Context) ([]*domain.Opcional, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Opcional, 0, len(r.opcionais))
	for _, o := range r.opcionais {
		if o.Ativo {
			copia := *o
			res = append(res, &copia)
		}
	}
	return res, nil
}

// BuscarOpcionalPorID retorna uma copia do opcional.
func (r *CatalogoRepository) BuscarOpcionalPorID(ctx context.Context, id uuid.UUID) (*domain.Opcional, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.opcionais[id]
	if !ok {
		return nil, domain.ErrOpcionalNaoEncontrado
	}
	copia := *o
	return &copia, nil
}

// BuscarTabelaPreco retorna a linha da tabela para (categoria, regiao).
func (r *CatalogoRepository) BuscarTabelaPreco(ctx context.Context, categoriaID, regiaoID uuid.UUID) (*domain.TabelaPrecos, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tabelas[tabelaKey(categoriaID, regiaoID)]
	if !ok {
		return nil, domain.ErrTabelaPrecoNaoEncontrada
	}
	copia := *t
	return &copia, nil
}
