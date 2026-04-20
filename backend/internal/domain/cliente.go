package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Cliente representa o perfil completo do usuario cliente.
// Espelha a tabela `clientes` do banco.
type Cliente struct {
	ID           uuid.UUID
	UsuarioID    uuid.UUID
	Nome         string
	CPF          string // armazenado sem mascara: apenas digitos
	Telefone     string // armazenado sem mascara: apenas digitos
	Score        int    // 0-100, score inicial 100
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// Erros de dominio do cliente.
var (
	ErrClienteNaoEncontrado = errors.New("cliente nao encontrado")
	ErrClienteJaExiste      = errors.New("cliente ja cadastrado para este usuario")
	ErrCPFJaCadastrado      = errors.New("CPF ja cadastrado")
	ErrCPFInvalido          = errors.New("CPF invalido")
	ErrNomeObrigatorio      = errors.New("nome e obrigatorio")
	ErrTelefoneInvalido     = errors.New("telefone invalido")
	ErrScoreInvalido        = errors.New("score deve estar entre 0 e 100")
)

// ClienteReader define operacoes de leitura do repositorio de clientes.
type ClienteReader interface {
	BuscarPorUsuarioID(ctx context.Context, usuarioID uuid.UUID) (*Cliente, error)
	BuscarPorID(ctx context.Context, id uuid.UUID) (*Cliente, error)
	BuscarPorCPF(ctx context.Context, cpf string) (*Cliente, error)
}

// ClienteWriter define operacoes de escrita do repositorio de clientes.
type ClienteWriter interface {
	Criar(ctx context.Context, c *Cliente) error
	Atualizar(ctx context.Context, c *Cliente) error
	// AjustarScore soma delta ao score do cliente, com clamp em [0, 100].
	// Usado para penalidades (delta negativo) e bônus (positivo).
	AjustarScore(ctx context.Context, id uuid.UUID, delta int) error
}

// ClienteRepository e a interface completa do repositorio de clientes.
type ClienteRepository interface {
	ClienteReader
	ClienteWriter
}

// ClienteRequest representa o payload de criacao/atualizacao do perfil do cliente.
type ClienteRequest struct {
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Telefone string `json:"telefone"`
}

// ClienteResponse e a resposta com os dados do perfil do cliente.
type ClienteResponse struct {
	ID        uuid.UUID `json:"id"`
	UsuarioID uuid.UUID `json:"usuario_id"`
	Nome      string    `json:"nome"`
	CPF       string    `json:"cpf"`      // retorna sem mascara
	Telefone  string    `json:"telefone"` // retorna sem mascara
	Score     int       `json:"score"`
}

// ValidarCPF verifica o CPF usando algoritmo de digito verificador.
// Aceita CPF apenas com digitos (sem mascara).
func ValidarCPF(cpf string) error {
	if len(cpf) != 11 {
		return ErrCPFInvalido
	}
	for _, c := range cpf {
		if c < '0' || c > '9' {
			return ErrCPFInvalido
		}
	}

	// Verificar se todos os digitos sao iguais (invalido)
	igual := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			igual = false
			break
		}
	}
	if igual {
		return ErrCPFInvalido
	}

	// Calcular primeiro digito verificador
	soma := 0
	for i := 0; i < 9; i++ {
		soma += int(cpf[i]-'0') * (10 - i)
	}
	resto := soma % 11
	d1 := 0
	if resto >= 2 {
		d1 = 11 - resto
	}
	if int(cpf[9]-'0') != d1 {
		return ErrCPFInvalido
	}

	// Calcular segundo digito verificador
	soma = 0
	for i := 0; i < 10; i++ {
		soma += int(cpf[i]-'0') * (11 - i)
	}
	resto = soma % 11
	d2 := 0
	if resto >= 2 {
		d2 = 11 - resto
	}
	if int(cpf[10]-'0') != d2 {
		return ErrCPFInvalido
	}

	return nil
}

// ValidarTelefone verifica se o telefone tem entre 10 e 11 digitos.
func ValidarTelefone(telefone string) error {
	if len(telefone) < 10 || len(telefone) > 11 {
		return ErrTelefoneInvalido
	}
	for _, c := range telefone {
		if c < '0' || c > '9' {
			return ErrTelefoneInvalido
		}
	}
	return nil
}

// ValidarNome verifica se o nome nao esta vazio e tem tamanho razoavel.
func ValidarNome(nome string) error {
	if nome == "" {
		return ErrNomeObrigatorio
	}
	if len(nome) > 100 {
		return errors.New("nome deve ter no maximo 100 caracteres")
	}
	return nil
}

// ValidarScore garante o intervalo 0..100 espelhando o CHECK do banco.
// Primeira linha de defesa: reprova entrada invalida antes do UPDATE.
func ValidarScore(score int) error {
	if score < 0 || score > 100 {
		return ErrScoreInvalido
	}
	return nil
}
