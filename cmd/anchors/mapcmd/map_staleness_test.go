package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// projectWithMap monta um projeto mínimo e devolve raiz, grafo e config.
func projectWithMap(t *testing.T, conteudo string) (string, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(
		"version: 2\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "U.spec.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(filepath.Join(root, "anchors.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	g := mapx.Build(files, cfg, nil)
	return root, g, cfg
}

// O MAPA VELHO precisa ser acusado — e "velho" é diferente de "ausente".
//
// O pre-commit já barrava o arquivo REGIDO fora do mapa (o arquivo novo). O que passava era
// o arquivo que ESTÁ no mapa com a `rev` de uma versão anterior: o `map build` rodou, e o
// trabalho continuou depois dele.
//
// MEDIDO no projeto de referência: das 12 reprovações do pipeline `gates`, SETE foram isto —
// a maior causa isolada. Nas sete o agente tinha commitado o mapa. Ninguém esqueceu de
// gerá-lo; todos geraram cedo demais e seguiram editando.
func TestMapaDesatualizadoAcusaOArquivoEditadoDepoisDoBuild(t *testing.T) {
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — a regra.\n")

	if v := StaleMapNodes(root, g, cfg); len(v) != 0 {
		t.Fatalf("mapa recém-construído acusado como velho: %v", v)
	}

	// A EDIÇÃO DEPOIS DO BUILD — o caso dos sete.
	if err := os.WriteFile(filepath.Join(root, "U.spec.md"),
		[]byte("# U\n\nUMA-B01 — a regra.\nUMA-B02 — outra.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	v := StaleMapNodes(root, g, cfg)
	if len(v) != 1 || !strings.HasSuffix(v[0], "U.spec.md") {
		t.Errorf("a edição depois do `map build` não foi acusada: %v", v)
	}
}

// AUSENTE não é VELHO, e misturá-los daria uma mensagem que manda fazer a coisa errada.
//
// Um nó cujo arquivo sumiu é outro caso — o `map build` o remove. Acusá-lo aqui mandaria
// reconstruir o mapa por um arquivo que não existe mais.
func TestMapaDesatualizadoIgnoraArquivoQueSumiu(t *testing.T) {
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — a regra.\n")
	if err := os.Remove(filepath.Join(root, "U.spec.md")); err != nil {
		t.Fatal(err)
	}
	if v := StaleMapNodes(root, g, cfg); len(v) != 0 {
		t.Errorf("arquivo REMOVIDO acusado como mapa velho: %v — são casos diferentes, "+
			"e o conserto de um não é o do outro", v)
	}
}

// A COMPARAÇÃO É DE HASH, não de data.
//
// `updated_at` muda num `git checkout` sem o conteúdo mudar, e não muda numa edição que
// preserve o mtime. Nos dois casos a resposta sairia errada — e a segunda é exatamente o
// cenário que este código existe para pegar.
func TestMapaDesatualizadoNaoOlhaData(t *testing.T) {
	root, g, cfg := projectWithMap(t, "# U\n\nUMA-B01 — a regra.\n")
	arq := filepath.Join(root, "U.spec.md")

	// mtime mexido, CONTEÚDO intacto: não é mapa velho.
	antigo := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(arq, antigo, antigo); err != nil {
		t.Fatal(err)
	}
	if v := StaleMapNodes(root, g, cfg); len(v) != 0 {
		t.Errorf("mtime alterado com conteúdo intacto acusado como velho: %v", v)
	}

	// CONTEÚDO mexido, mtime restaurado ao antigo: É mapa velho.
	if err := os.WriteFile(arq, []byte("# U\n\nUMA-B01 — outra coisa.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(arq, antigo, antigo); err != nil {
		t.Fatal(err)
	}
	if v := StaleMapNodes(root, g, cfg); len(v) != 1 {
		t.Errorf("conteúdo mudado com mtime velho NÃO foi acusado: %v — a comparação "+
			"caiu para data, e é o caso que o hash existe para pegar", v)
	}
}
