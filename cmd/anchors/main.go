// Command anchors é o CLI único do framework Anchors.
package main

import (
	"errors"
	"fmt"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/migra"
	"os"
)

// Preenchidas via -ldflags no build de release (ver cli/.goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
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
	defer flushTelemetry()

	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		// "não é regido" sai com código PRÓPRIO: quem automatiza (pre-commit, CI)
		// precisa distinguir "não tenho jurisdição sobre este arquivo" de "este
		// arquivo reprovou". Sem isso, só resta grepar a mensagem — e foi assim que
		// o pre-commit passou a deixar arquivo regido novo escapar sem trinca.
		var nr errNotGoverned
		if errors.As(err, &nr) {
			flushTelemetry()
			os.Exit(ExitNotGoverned)
		}
		flushTelemetry()
		os.Exit(1)
	}
}
