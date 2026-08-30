package mapx

import (
	"os"
	"regexp"

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
	// SE SÓ O `gerado_por` MUDOU, não reescreve. Sem isto o campo causaria justamente a
	// oscilação que existe para denunciar: o build local grava "dev", o CI grava "0.1.10",
	// e cada um desfaz o do outro a cada execução.
	//
	// A comparação é do CONTEÚDO — o mapa é derivado, e reescrevê-lo idêntico não é
	// operação neutra: muda o mtime, e é o diff que importa.
	header := []byte("# anchors.graph.yaml — o mapa de dependências (gerado por `anchors map build`)\n" +
		"# Material e versionado. Fonte de verdade do grafo; índices em memória são cache.\n")
	completo := append(header, data...)
	// A comparação é com o arquivo COMPLETO, header incluído. A primeira versão comparava
	// só o YAML contra o arquivo em disco — e como o disco tem o header, nunca eram
	// iguais: a guarda não guardava nada, e só o teste isolado mostrou.
	if igualIgnorandoGeradoPor(path, completo) {
		return nil
	}
	return os.WriteFile(path, completo, 0o644)
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

// igualIgnorandoGeradoPor diz se o mapa em disco é o mesmo, desconsiderando quem o gerou.
func igualIgnorandoGeradoPor(path string, novo []byte) bool {
	atual, err := os.ReadFile(path)
	if err != nil {
		return false // não existe ainda: há o que escrever
	}
	return semGeradoPor(atual) == semGeradoPor(novo)
}

var geradoPorLinhaRE = regexp.MustCompile(`(?m)^gerado_por:.*\n`)

func semGeradoPor(b []byte) string {
	return geradoPorLinhaRE.ReplaceAllString(string(b), "")
}
