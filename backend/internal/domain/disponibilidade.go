package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Disponibilidade representa a janela de horario disponivel de uma profissional num dia da semana.
type Disponibilidade struct {
	ID             uuid.UUID
	ProfissionalID uuid.UUID
	DiaSemana      int    // 0=Domingo, 1=Segunda, ..., 6=Sabado
	HoraInicio     string // formato "HH:MM"
	HoraFim        string // formato "HH:MM"
	CriadoEm       time.Time
	AtualizadoEm   time.Time
}

// Erros de dominio de disponibilidade.
var (
	ErrDisponibilidadeNaoEncontrada = errors.New("disponibilidade nao encontrada")
	ErrDiaSemanaInvalido            = errors.New("dia da semana invalido: deve ser 0 (dom) a 6 (sab)")
	ErrHoraInicioInvalida           = errors.New("hora de inicio invalida: use formato HH:MM")
	ErrHoraFimInvalida              = errors.New("hora de fim invalida: use formato HH:MM")
	ErrHoraFimAntesDaInicio         = errors.New("hora de fim deve ser apos a hora de inicio")
)

// DisponibilidadeReader define operacoes de leitura.
type DisponibilidadeReader interface {
	ListarPorProfissionalID(ctx context.Context, profissionalID uuid.UUID) ([]*Disponibilidade, error)
}

// DisponibilidadeWriter define operacoes de escrita.
type DisponibilidadeWriter interface {
	// DefinirDisponibilidades substitui todas as disponibilidades da profissional.
	DefinirDisponibilidades(ctx context.Context, profissionalID uuid.UUID, slots []*Disponibilidade) error
}

// DisponibilidadeRepository e a interface completa.
type DisponibilidadeRepository interface {
	DisponibilidadeReader
	DisponibilidadeWriter
}

// DisponibilidadeSlot e um item do payload de definicao de disponibilidade.
type DisponibilidadeSlot struct {
	DiaSemana  int    `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}

// DefinirDisponibilidadesRequest representa o payload para definir disponibilidade semanal.
type DefinirDisponibilidadesRequest struct {
	Slots []DisponibilidadeSlot `json:"slots"`
}

// DisponibilidadeResponse e a resposta com os dados de disponibilidade.
type DisponibilidadeResponse struct {
	ID         uuid.UUID `json:"id"`
	DiaSemana  int       `json:"dia_semana"`
	HoraInicio string    `json:"hora_inicio"`
	HoraFim    string    `json:"hora_fim"`
}

// ValidarHora verifica se a string esta no formato HH:MM valido.
func ValidarHora(h string) error {
	if len(h) != 5 || h[2] != ':' {
		return errors.New("formato invalido")
	}
	for i, c := range h {
		if i == 2 {
			continue
		}
		if c < '0' || c > '9' {
			return errors.New("formato invalido")
		}
	}
	hh := int(h[0]-'0')*10 + int(h[1]-'0')
	mm := int(h[3]-'0')*10 + int(h[4]-'0')
	if hh > 23 || mm > 59 {
		return errors.New("hora ou minuto fora do intervalo valido")
	}
	return nil
}
