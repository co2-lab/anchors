// Command anchors é o CLI único do framework Anchors.
package main

import (
	"errors"
	"fmt"
	"os"
)

// Preenchidas via -ldflags no build de release (ver cli/.goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
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
