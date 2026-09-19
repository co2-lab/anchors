package board

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// --- o `gh` NÃO PODE PENDURAR O ANCHORS ---
//
// Toda conversa com o board passa pelo `gh`, e ele é um processo externo que o Anchors não
// controla. Sem credencial, o `gh` costuma falhar rápido — mas há caminhos em que ele
// ESPERA: um prompt de autenticação, um proxy que não responde, uma rede que engole o
// pacote sem devolver erro.
//
// Medido: um dev novo rodou `anchors status` num ambiente sem `gh auth login`, o comando
// ficou pendurado, e ele precisou matá-lo. O relato dele foi "o anchors fica pendurado" —
// e do ponto de vista de quem usa, é exatamente isso: a ferramenta travou.
//
// Um comando que não volta é pior que um que falha. Quem falha diz o que fazer; quem
// pendura consome a sessão e não deixa nem a mensagem de erro.

// ghTimeout é o teto de cada chamada ao `gh`.
//
// Trinta segundos: uma consulta ao board leva menos de dois em rede normal, e um `gh
// workflow run` leva menos de cinco. O que passa disso não é lentidão — é espera por algo
// que não vai chegar (prompt, proxy, rede morta).
//
// Generoso de propósito. Um teto apertado transformaria rede lenta em erro, e o custo de
// errar aqui é assimétrico: falhar cedo demais faz o agente desistir de trabalho que ia
// funcionar; falhar tarde demais custa uma sessão.
const ghTimeout = 30 * time.Second

// runGH executa o `gh` com teto de tempo, e NUNCA deixa o processo esperando entrada.
//
// O stdin fechado é a outra metade da proteção: sem ele, um `gh` que decida pedir
// credencial ficaria bloqueado na leitura, e o timeout mataria o processo sem que ninguém
// soubesse por quê. Com stdin fechado, o `gh` desiste e devolve a mensagem que explica o
// que fazer.
func runGH(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ghTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", args...)
	cmd.Stdin = nil // sem entrada: um prompt falha em vez de esperar
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return out, fmt.Errorf("`gh` did not respond in %s (`gh %s`).\n"+
			"  This is usually credentials: run `gh auth status` and, if needed,\n"+
			"  `gh auth login`. It can also be proxy or network — Anchors does not wait\n"+
			"  longer than this because a command that never returns costs the whole session",
			ghTimeout, strings.Join(args, " "))
	}
	if err != nil {
		return out, fmt.Errorf("%w%s", err, authHint(err, args))
	}
	return out, nil
}

// authHint acrescenta a causa provável quando o `gh` falha por credencial.
//
// O `exit status 4` do `gh` é o código de "não autenticado", e sozinho ele não diz nada:
// medido com um dev novo, o `anchors next` respondeu `exit status 4` e ele não tinha como
// saber que faltava login. Um erro que não nomeia a causa manda a pessoa procurar no lugar
// errado — e o lugar errado costuma ser a própria máquina.
func authHint(err error, args []string) string {
	var ee *exec.ExitError
	if !errors.As(err, &ee) || ee.ExitCode() != exitCodeNotAuthenticated {
		return ""
	}
	return "\n\n  `gh` code 4 is NOT AUTHENTICATED. In `github` mode the board is the queue,\n" +
		"  and without credentials Anchors reads no card and claims no work:\n\n" +
		"      gh auth login\n\n" +
		"  The login is interactive — an agent has no way to complete it."
}

// exitCodeNotAuthenticated é o `exit status` que o `gh` devolve sem credencial.
const exitCodeNotAuthenticated = 4
