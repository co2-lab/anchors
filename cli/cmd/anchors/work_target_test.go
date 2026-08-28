package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func cfgUnidade() *config.Config {
	return &config.Config{Layers: map[string]config.Layer{
		"logic":   {Pattern: "src/**/*.ts", Kind: "code"},
		"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
		"feature": {Pattern: "**/*.feature", Kind: "feature"},
		"test":    {Pattern: "**/*.test.ts", Kind: "test"},
	}}
}

// O alvo do `work` é a UNIDADE, não uma peça derivada. Apontar para a spec é o engano
// previsível (é o artefato que já existe), e o resultado era silenciosamente absurdo:
// caminhos como `x.spec.spec.md` e `x.spec.feature`, sem aviso nenhum.
func TestWorkRedirecionaPecaDerivada(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src"), 0o755)
	unidade := "src/metadataVersioning.ts"
	os.WriteFile(filepath.Join(dir, unidade), []byte("export const x = 1\n"), 0o644)

	for _, peca := range []string{
		"src/metadataVersioning.spec.md",
		"src/metadataVersioning.feature",
		"src/metadataVersioning.test.ts",
	} {
		t.Run(peca, func(t *testing.T) {
			got, achou := unidadeDaPecaDerivada(dir, peca, cfgUnidade(), nil)
			if !achou {
				t.Fatalf("não reconheceu %q como peça derivada", peca)
			}
			if got != unidade {
				t.Fatalf("redirecionou para %q, queria %q", got, unidade)
			}
		})
	}
}

// O próprio arquivo de código NÃO é redirecionado — ele já é a unidade.
func TestWorkNaoRedirecionaAUnidade(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src"), 0o755)
	os.WriteFile(filepath.Join(dir, "src/x.ts"), []byte("x\n"), 0o644)
	if _, achou := unidadeDaPecaDerivada(dir, "src/x.ts", cfgUnidade(), nil); achou {
		t.Fatal("o código é a unidade — não deveria redirecionar")
	}
}

// Com mapa, a aresta `specifies` é a fonte precisa: diz exatamente qual código a spec
// descreve, mesmo quando a convenção de nome não bastaria.
func TestWorkUsaOMapaQuandoExiste(t *testing.T) {
	g := &mapx.Graph{Edges: []mapx.Edge{
		{From: "a/nome-diferente.spec.md", To: "b/outroNome.ts", Type: "specifies"},
	}}
	got, achou := unidadeDaPecaDerivada(t.TempDir(), "a/nome-diferente.spec.md", cfgUnidade(), g)
	if !achou || got != "b/outroNome.ts" {
		t.Fatalf("o mapa deveria resolver o alvo: got=%q achou=%v", got, achou)
	}
}

// Peça derivada cuja unidade não existe: não inventa alvo (quem chama segue com o
// original e o resto do prompt explica o que falta).
func TestWorkSemUnidadeNaoInventa(t *testing.T) {
	if _, achou := unidadeDaPecaDerivada(t.TempDir(), "src/fantasma.spec.md", cfgUnidade(), nil); achou {
		t.Fatal("sem código no disco e sem mapa, não há alvo a deduzir")
	}
}
