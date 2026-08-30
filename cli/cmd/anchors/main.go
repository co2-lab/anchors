// Command anchors é o CLI único do framework Anchors.
package main

import (
	"errors"
	"fmt"
	"github.com/co2-lab/anchors/internal/mapx"
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
	mapx.GeradoPor = version
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		// "não é regido" sai com código PRÓPRIO: quem automatiza (pre-commit, CI)
		// precisa distinguir "não tenho jurisdição sobre este arquivo" de "este
		// arquivo reprovou". Sem isso, só resta grepar a mensagem — e foi assim que
		// o pre-commit passou a deixar arquivo regido novo escapar sem trinca.
		var nr errNaoRegido
		if errors.As(err, &nr) {
			os.Exit(ExitNaoRegido)
		}
		os.Exit(1)
	}
}
