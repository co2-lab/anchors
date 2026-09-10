// Command anchors é o CLI único do framework Anchors.
package main

import (
	"errors"
	"fmt"
	"github.com/co2-lab/anchors/internal/config"
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
	mapx.GeneratedBy = version
	// O MAPA guarda o nome do gate no carimbo de julgamento, e os nomes migraram do
	// português para o inglês. Sem esta ponte, um carimbo que diz `regra-cumprida` fica
	// órfão quando o gate passa a se chamar `rule-fulfilled` — e o julgamento volta para
	// a fila, com a resposta ali, visível no arquivo, sem nada os ligando.
	//
	// Medido no blue-eyes: 40 julgamentos órfãos de uma vez, e as 10 specs do projeto
	// pedindo julgamento que alguém já tinha dado.
	//
	// Ligado no `main` porque `mapx` importa `config` — a injeção na direção contrária
	// seria ciclo.
	mapx.ResolveGateName = func(n string) string {
		c, _ := config.CanonicalName(n)
		return c
	}
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		// "não é regido" sai com código PRÓPRIO: quem automatiza (pre-commit, CI)
		// precisa distinguir "não tenho jurisdição sobre este arquivo" de "este
		// arquivo reprovou". Sem isso, só resta grepar a mensagem — e foi assim que
		// o pre-commit passou a deixar arquivo regido novo escapar sem trinca.
		var nr errNotGoverned
		if errors.As(err, &nr) {
			os.Exit(ExitNotGoverned)
		}
		os.Exit(1)
	}
}
