// Package telemetry emite os EVENTOS DE DECISÃO do Anchors.
//
// POR QUE EXISTE. O que um agente decide é invisível depois que o comando termina. Medido
// numa sessão com seis agentes: três pegaram a mesma issue e ninguém soube até aparecerem
// três PRs idênticos; um encerrou o turno com o card `in-progress` e o trabalho ficou
// parado horas; 25 PRs verdes esperaram revisão com zero cards na fila de revisão.
//
// Nenhum desses é ERRO — todos os comandos retornaram 0. São decisões que, vistas em
// sequência, formam um padrão. É isso que este pacote registra.
//
// O QUE NÃO VAI NO EVENTO, e a lista é curta de propósito:
//
//   - conteúdo de arquivo — nem spec, nem diff, nem mensagem de erro com caminho;
//   - nome de repositório, de usuário, de branch;
//   - qualquer texto livre que alguém tenha escrito.
//
// O que sobra são NÚMEROS e VOCABULÁRIO DO PRODUTO: quantos candidatos o claim viu, qual
// gate reprovou, em que estado o turno terminou. O código de identidade de uma unidade
// (`RLSGR`, `ALFDL`) é o limite — ele é do projeto de quem usa, e por isso a emissão o
// trata como opcional (ver `Config.NoCodes`).
package telemetry

import (
	"time"
)

// Name é o tipo do evento. Fechado de propósito: um conjunto aberto viraria texto livre, e
// texto livre é onde vazamento entra sem ninguém notar.
type Name string

const (
	// ClaimServed — o claim entregou um card. Os atributos dizem QUANTOS candidatos ele
	// viu e quantos pulou, por quê: é o que revela três agentes disputando a mesma fila.
	ClaimServed Name = "claim.served"
	// ClaimEmpty — o claim não achou trabalho. Com `pulados > 0` significa que HAVIA
	// trabalho e algo o barrou, que é diferente de fila vazia.
	ClaimEmpty Name = "claim.empty"
	// CheckFinished — quantos gates passaram, reprovaram, e se barrou.
	CheckFinished Name = "check.finished"
	// EscalatedToUser — o trabalho parou esperando uma decisão de pessoa.
	EscalatedToUser Name = "escalate.raised"
	// TurnEnded — o comando terminou e o card ficou EM ALGUM ESTADO. É o evento que
	// teria pego "aguardando a nova rodada": o turno acabou, o card seguiu `in-progress`,
	// e ninguém soube.
	TurnEnded Name = "turn.ended"
)

// Event é uma decisão registrada.
type Event struct {
	Name Name
	At   time.Time
	// Attrs são pares chave-valor. Números e vocabulário do produto — nunca texto livre.
	Attrs map[string]any
}

// New build um evento com o instante de now.
//
// O relógio é injetável (`now`) porque o pacote não inventa tempo: quem chama carimba, e
// o teste passa um relógio fixo em vez de dormir.
func New(n Name, attrs map[string]any, now func() time.Time) Event {
	if attrs == nil {
		attrs = map[string]any{}
	}
	return Event{Name: n, At: now(), Attrs: attrs}
}
