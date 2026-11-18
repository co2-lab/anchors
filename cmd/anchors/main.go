// @anchors
//   code: MNCMM
//   ref: CLMNC

// Command anchors é o CLI único do framework Anchors.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/cmd/anchors/flow"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/migra"
)

// Preenchidas via -ldflags no build de release (ver cli/.goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	common.Version = version
	common.Commit = commit
	common.Date = date

	// O mapa registra QUEM o escreveu, para que um binário mais velho seja acusado em vez
	// de reverter em silêncio o que a versão nova gravou (ver mapx.GeradoPor).
	mapx.GeneratedBy = version
	// A mensagem de chave desconhecida precisa distinguir TYPO de chave RENOMEADA, e quem
	// sabe disso é o registro de migração. Injetado aqui porque o `config` não pode
	// importar o `migra` — seria ciclo.
	config.RenamedKey = migra.RenamedKey

	// O FLUSH da telemetria, e ele precisa acontecer nos DOIS caminhos de saída.
	//
	// `Emit` dispara o POST numa goroutine para não atrasar o comando — e num CLI isso
	// significa que o processo morre antes de o envio sair. Medido: o coletor local não
	// recebeu NADA até este `defer` existir, e o defeito é invisível em teste de unidade,
	// onde o processo continua vivo depois da chamada.
	//
	// `defer` e não uma chamada no fim: o caminho de ERRO abaixo sai com `os.Exit`, que
	// não roda defers — por isso ele também chama o flush, explicitamente.
	defer common.FlushTelemetry()

	err := newRootCmd().Execute()
	exit := exitCodeOf(err)
	// The run of a long command is recorded for the monitor, with how it ended.
	flow.EndOwnRun(exit)
	if err == nil {
		return
	}
	// A command that ends with an exit code already said what it had to.
	var ec common.ExitCode
	if !errors.As(err, &ec) {
		fmt.Fprintln(os.Stderr, i18n.T("error")+":", err)
	}
	common.FlushTelemetry()
	os.Exit(exit)
}

// exitCodeOf is the code the process exits with: 0 with no error; the code a command ended
// with; its own code for "not governed" — whoever automates (pre-commit, CI) must tell "this
// file is not mine to judge" from "this file failed", and grepping the message is how the
// pre-commit once let a new governed file through without its unit —; and 1 otherwise.
func exitCodeOf(err error) int {
	if err == nil {
		return 0
	}
	var ec common.ExitCode
	if errors.As(err, &ec) {
		return ec.Code
	}
	var nr common.ErrNotGoverned
	if errors.As(err, &nr) {
		return common.ExitNotGoverned
	}
	return 1
}
