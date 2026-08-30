package mapx

import (
	"os"

	"gopkg.in/yaml.v3"
)

// DefaultPath é onde o mapa material vive na raiz do projeto.
const DefaultPath = "anchors.graph.yaml"

// GeradoPor é a versão do binário que escreve o mapa, preenchida pelo `main` no início.
//
// Existe porque um binário DESATUALIZADO não avisa — ele grava o formato que conhece, e
// desfaz o que a versão nova escreveu. Medido: depois de renomear um campo do carimbo, o
// `anchors` do PATH (build local anterior) revertia as 26 linhas a cada `check`, e o mapa
// ficava oscilando entre os dois formatos, com conflito a cada PR.
//
// O `--version` não denunciava: os dois builds locais se identificam como "dev". Só o
// diff do mapa mostrava, e para isso alguém precisa estar olhando.
var GeradoPor string

// Save escreve o grafo como YAML material e versionável.
func Save(g *Graph, path string) error {
	g.GeradoPor = GeradoPor
	data, err := yaml.Marshal(g)
	if err != nil {
		return err
	}
	header := []byte("# anchors.graph.yaml — o mapa de dependências (gerado por `anchors map build`)\n" +
		"# Material e versionado. Fonte de verdade do grafo; índices em memória são cache.\n")
	return os.WriteFile(path, append(header, data...), 0o644)
}

// Load lê um grafo existente do disco.
func Load(path string) (*Graph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var g Graph
	if err := yaml.Unmarshal(data, &g); err != nil {
		return nil, err
	}
	return &g, nil
}
