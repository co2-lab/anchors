package scan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestProgresso_reconheceOSufixo(t *testing.T) {
	casos := map[string]bool{
		"plans/0017-mutacao-progress.md": true,
		"plans/0017-mutacao.md":          false,
		"0001-progress.md":               true,
		// Não é o sufixo: `progress` no meio do nome não faz um arquivo de estado.
		"docs/progress-notes.md": false,
		"plans/progress.md":      false,
	}
	for caminho, quer := range casos {
		if got := EhArquivoDeProgresso(caminho); got != quer {
			t.Errorf("%q: %v, queria %v", caminho, got, quer)
		}
	}
}

// O PROGRESSO NÃO ENTRA NO MAPA — é o ponto inteiro da separação.
//
// Um arquivo que existe para MUDAR não pode ser confrontado pelos gates que cobram
// justificativa de mudança: o `plano-alterado-justificado` o acusaria a cada `[x]` novo,
// e a única saída seria escrever uma revisão dizendo "o trabalho andou".
//
// O teste roda o `Walk` de verdade, com o plano e o companheiro na mesma pasta e a MESMA
// camada casando os dois — que é a situação real. Um teste que só chamasse
// `EhArquivoDeProgresso` provaria que a função responde certo, não que a varredura a usa.
func TestProgresso_ficaForaDoMapa(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	plano := "<!-- @anchors\n  code: MTUAO\n  layer: plan\n-->\n# Plano\n\n### MTUAO-F01 — fase\n"
	if err := os.WriteFile(filepath.Join(dir, "plans", "0017-mutacao.md"), []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}
	prog := "# Progresso — MTUAO\n\n## MTUAO-F01\n\n- [x] feito\n"
	if err := os.WriteFile(filepath.Join(dir, "plans", "0017-mutacao-progress.md"), []byte(prog), 0o644); err != nil {
		t.Fatal(err)
	}

	// A camada casa OS DOIS arquivos: é o que torna o teste honesto. Se o glob excluísse
	// o companheiro, o teste passaria sem que a guarda do `Walk` fizesse nada.
	cfg := &config.Config{
		Layers: map[string]config.Layer{"plan": {Pattern: "plans/*.md", Kind: "plan"}},
	}

	files, err := Walk(dir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var achouPlano, achouProgresso bool
	for _, f := range files {
		switch f.Path {
		case "plans/0017-mutacao.md":
			achouPlano = true
		case "plans/0017-mutacao-progress.md":
			achouProgresso = true
		}
	}
	if !achouPlano {
		t.Errorf("o PLANO tem de entrar no mapa — sem ele o teste não prova nada sobre o "+
			"companheiro; arquivos vistos: %v", caminhos(files))
	}
	if achouProgresso {
		t.Errorf("o arquivo de progresso ENTROU no mapa: os gates voltariam a cobrar "+
			"justificativa a cada `[x]`; arquivos vistos: %v", caminhos(files))
	}
}

func caminhos(fs []File) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.Path)
	}
	return out
}
