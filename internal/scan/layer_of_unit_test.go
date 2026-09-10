package scan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// A camada da UNIDADE não é a do arquivo, e a spec é o caso que separa as duas.
//
// Uma spec casa dois padrões: o dela (`**/*.spec.md`, camada `spec`) e o da unidade que
// governa. O `ClassifyPath` devolve o primeiro — certo para o mapa, onde a spec É um nó da
// camada `spec` — e errado para quem pergunta "que camada é esta unidade".
//
// Medido: `anchors docs duties --unit <x>.spec.md` respondia "nenhuma documentação
// obrigatória" para uma lambda que deve o OpenAPI. Com o `.ts` da mesma unidade a resposta
// era certa, e o card aponta a SPEC — o caminho que o agente usa.
func TestLayerOfUnit(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{Layers: map[string]config.Layer{
		"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
		"lambdas": {Pattern: "pkg/lambdas/**/*.ts", Kind: "code"},
	}}

	escreve := func(rel, corpo string) {
		p := filepath.Join(root, rel)
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(corpo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	escreve("pkg/lambdas/push/X.spec.md", "<!-- @anchors\n  code: XPTOX\n  layer: lambdas\n-->\n# X\n")
	escreve("pkg/lambdas/push/X.ts", "export const x = 1;\n")
	escreve("pkg/lambdas/push/SemHeader.spec.md", "# sem header\n")

	casos := []struct{ nome, rel, esperado string }{
		{"a spec declara a camada da unidade no header", "pkg/lambdas/push/X.spec.md", "lambdas"},
		{"o código já tem a camada da unidade pelo padrão", "pkg/lambdas/push/X.ts", "lambdas"},
		// Sem header, cai no `ClassifyPath` — e ali a spec é da camada `spec` mesmo.
		// Inventar outra resposta seria adivinhar.
		{"spec sem header cai no padrão", "pkg/lambdas/push/SemHeader.spec.md", "spec"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := LayerOfUnit(root, c.rel, cfg); got != c.esperado {
				t.Errorf("LayerOfUnit(%s) = %q, queria %q", c.rel, got, c.esperado)
			}
		})
	}

	// E o `ClassifyPath` continua devolvendo a camada do ARQUIVO — o mapa depende disso:
	// a spec é um nó da camada `spec`, e mudar isso reclassificaria 84 nós num projeto real.
	if l, _ := ClassifyPath("pkg/lambdas/push/X.spec.md", cfg); l != "spec" {
		t.Errorf("ClassifyPath da spec = %q, queria `spec` — o mapa depende disso", l)
	}
}
