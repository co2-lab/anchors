package mapx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// O OVERRIDE POR CAMADA precisa casar a camada da UNIDADE, não a do arquivo.
//
// Uma spec casa `**/*.spec.md`, e o `Layer` dela é `spec`. Um `when: screen` comparado
// contra ele nunca casa — e o override por camada seria inútil justamente para a âncora,
// que é quem o consulta para resolver os derivados.
//
// Medido no projeto de referência: as camadas de UI exigem `.tsx` no pattern, o
// `derived.files` global diz `.ts`, e o override que reconciliava os dois não era aplicado.
// O `triad-complete` respondia "falta o código" com o arquivo no disco — e a Estrutura
// declarava duas coisas incompatíveis sem nada acusar.
func TestLayerOfUnit_oHeaderVenceOArquivo(t *testing.T) {
	casos := []struct {
		nome     string
		f        scan.File
		esperado string
	}{
		{"a spec declara a camada da unidade", scan.File{Layer: "spec", HeaderLayer: "screen"}, "screen"},
		{"sem header, o arquivo já é a unidade", scan.File{Layer: "screen"}, "screen"},
		{"header vazio não sobrescreve", scan.File{Layer: "lambdas", HeaderLayer: ""}, "lambdas"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := layerOfUnit(c.f); got != c.esperado {
				t.Errorf("layerOfUnit = %q, queria %q", got, c.esperado)
			}
		})
	}
}

// E o efeito no MAPA: o derivado de uma screen é `.tsx`, porque o override casa.
func TestBuild_oOverridePorCamadaAlcancaASpec(t *testing.T) {
	cfg := &config.Config{
		Layers: map[string]config.Layer{
			"spec":   {Pattern: "**/*.spec.md", Kind: "spec"},
			"screen": {Pattern: "app/screens/**/*.tsx", Kind: "code"},
		},
		Derived: &config.Derived{
			Anchor: "spec",
			Files: map[string]config.Padroes{
				"code": {"{{dir}}/{{name}}.ts"},
			},
			Overrides: []config.DerivedOverride{
				{When: "screen", Files: map[string]config.Padroes{
					"code": {"{{dir}}/{{name}}.tsx"},
				}},
			},
		},
	}
	files := []scan.File{
		{Path: "app/screens/Tela.spec.md", Layer: "spec", Kind: "spec",
			HeaderCode: "TELAX", HeaderLayer: "screen"},
		{Path: "app/screens/Tela.tsx", Layer: "screen", Kind: "code"},
	}

	g := Build(files, cfg, nil)

	var achou bool
	for _, e := range g.Edges {
		if e.Type == EdgeSpecifies && e.To == "app/screens/Tela.tsx" {
			achou = true
		}
	}
	if !achou {
		t.Errorf("a spec não aponta o `.tsx` — o override por camada não alcançou a âncora.\nArestas: %+v", g.Edges)
	}
}
