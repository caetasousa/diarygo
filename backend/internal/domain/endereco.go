package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Endereco representa um endereco associado a um cliente.
// Espelha a tabela `enderecos` do banco.
type Endereco struct {
	ID           uuid.UUID
	ClienteID    uuid.UUID
	Logradouro   string
	Numero       string
	Complemento  string
	Bairro       string
	Cidade       string
	Estado       string // UF, 2 chars
	CEP          string // apenas digitos, 8 chars
	Latitude     float64
	Longitude    float64
	NumQuartos   int
	NumBanheiros int
	NumSalas     int
	NumCozinhas  int
	AreaM2       float64
	Principal    bool
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Erros de dominio de endereco.
var (
	ErrEnderecoNaoEncontrado        = errors.New("endereco nao encontrado")
	ErrEnderecoNaoPertenceAoCliente = errors.New("endereco nao pertence ao cliente")
	ErrCEPInvalido                  = errors.New("CEP invalido")
	ErrLogradouroObrigatorio        = errors.New("logradouro e obrigatorio")
	ErrCidadeObrigatoria            = errors.New("cidade e obrigatoria")
	ErrEstadoInvalido               = errors.New("estado invalido: deve ser UF com 2 letras")
	ErrComodoInvalido               = errors.New("numero de comodos deve ser maior que zero")
	ErrComodoNegativo               = errors.New("numero de comodos nao pode ser negativo")
	ErrAreaM2Invalida               = errors.New("area em m2 deve ser maior que zero")
)

// EnderecoReader define operacoes de leitura do repositorio de enderecos.
type EnderecoReader interface {
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Endereco, error)
	ListarPorClienteID(ctx context.Context, clienteID uuid.UUID) ([]*Endereco, error)
}

// EnderecoWriter define operacoes de escrita do repositorio de enderecos.
type EnderecoWriter interface {
	Criar(ctx context.Context, e *Endereco) error
	Atualizar(ctx context.Context, e *Endereco) error
	Remover(ctx context.Context, id uuid.UUID) error
}

// EnderecoRepository e a interface completa do repositorio de enderecos.
type EnderecoRepository interface {
	EnderecoReader
	EnderecoWriter
}

// EnderecoRequest representa o payload de criacao/atualizacao de endereco.
type EnderecoRequest struct {
	Logradouro   string  `json:"logradouro"`
	Numero       string  `json:"numero"`
	Complemento  string  `json:"complemento"`
	Bairro       string  `json:"bairro"`
	Cidade       string  `json:"cidade"`
	Estado       string  `json:"estado"`
	CEP          string  `json:"cep"`
	NumQuartos   int     `json:"num_quartos"`
	NumBanheiros int     `json:"num_banheiros"`
	NumSalas     int     `json:"num_salas"`
	NumCozinhas  int     `json:"num_cozinhas"`
	AreaM2       float64 `json:"area_m2"`
	Principal    bool    `json:"principal"`
}

// EnderecoResponse e a resposta com os dados do endereco.
type EnderecoResponse struct {
	ID           uuid.UUID `json:"id"`
	ClienteID    uuid.UUID `json:"cliente_id"`
	Logradouro   string    `json:"logradouro"`
	Numero       string    `json:"numero"`
	Complemento  string    `json:"complemento"`
	Bairro       string    `json:"bairro"`
	Cidade       string    `json:"cidade"`
	Estado       string    `json:"estado"`
	CEP          string    `json:"cep"`
	NumQuartos   int       `json:"num_quartos"`
	NumBanheiros int       `json:"num_banheiros"`
	NumSalas     int       `json:"num_salas"`
	NumCozinhas  int       `json:"num_cozinhas"`
	AreaM2       float64   `json:"area_m2"`
	Principal    bool      `json:"principal"`
}

// ValidarCEP verifica se o CEP tem exatamente 8 digitos numericos.
func ValidarCEP(cep string) error {
	if len(cep) != 8 {
		return ErrCEPInvalido
	}
	for _, c := range cep {
		if c < '0' || c > '9' {
			return ErrCEPInvalido
		}
	}
	return nil
}

// ufsValidas contem os estados brasileiros validos.
var ufsValidas = map[string]bool{
	"AC": true, "AL": true, "AP": true, "AM": true, "BA": true,
	"CE": true, "DF": true, "ES": true, "GO": true, "MA": true,
	"MT": true, "MS": true, "MG": true, "PA": true, "PB": true,
	"PR": true, "PE": true, "PI": true, "RJ": true, "RN": true,
	"RS": true, "RO": true, "RR": true, "SC": true, "SP": true,
	"SE": true, "TO": true,
}

// ValidarEstado verifica se o estado e uma UF brasileira valida.
func ValidarEstado(estado string) error {
	if len(estado) != 2 {
		return ErrEstadoInvalido
	}
	if !ufsValidas[estado] {
		return ErrEstadoInvalido
	}
	return nil
}

// ValidarComodos garante que cada contador de comodo seja >= 0, espelhando
// os CHECKs num_quartos/num_banheiros/num_salas/num_cozinhas do banco.
// Regra de negocio adicional: exige pelo menos 1 quarto (endereco residencial).
func ValidarComodos(quartos, banheiros, salas, cozinhas int) error {
	if quartos < 0 || banheiros < 0 || salas < 0 || cozinhas < 0 {
		return ErrComodoNegativo
	}
	if quartos < 1 {
		return ErrComodoInvalido
	}
	return nil
}

// ValidarAreaM2 garante > 0 quando presente, espelhando o CHECK do banco.
// Zero significa "nao informado" e eh aceito.
func ValidarAreaM2(area float64) error {
	if area == 0 {
		return nil
	}
	if area <= 0 {
		return ErrAreaM2Invalida
	}
	return nil
}

// ValidarFaixaCEP garante cep_fim >= cep_inicio quando ambos presentes,
// espelhando o CHECK ck_regioes_cep_ordem. Ambos podem ser vazios.
func ValidarFaixaCEP(cepInicio, cepFim string) error {
	if cepInicio == "" || cepFim == "" {
		return nil
	}
	if err := ValidarCEP(cepInicio); err != nil {
		return err
	}
	if err := ValidarCEP(cepFim); err != nil {
		return err
	}
	if cepFim < cepInicio {
		return errors.New("cep_fim deve ser maior ou igual a cep_inicio")
	}
	return nil
}
