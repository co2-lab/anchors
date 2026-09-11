package mapx

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// DefaultPath é onde o mapa material vive na raiz do projeto.
const DefaultPath = "anchors.graph.yaml"

// GeneratedBy é a versão do binário que escreve o mapa, preenchida pelo `main` no início.
//
// Existe porque um binário DESATUALIZADO não avisa — ele grava o formato que conhece, e
// desfaz o que a versão nova escreveu. Medido: depois de renomear um campo do carimbo, o
// `anchors` do PATH (build local anterior) revertia as 26 linhas a cada `check`, e o mapa
// ficava oscilando entre os dois formatos, com conflito a cada PR.
//
// O `--version` não denunciava: os dois builds locais se identificam como "dev". Só o
// diff do mapa mostrava, e para isso alguém precisa estar olhando.
var GeneratedBy string

// Save escreve o grafo como YAML material e versionável.
func Save(g *Graph, path string) error {
	// O FORMATO é de quem ESCREVE, não do chamador.
	//
	// O `Save` já carimbava a versão do binário e deixava o `version:` por conta de quem
	// montou o grafo — e um `Graph{}` sem o campo gravava `version: 0`, que na leitura
	// vira "formato 1: precisa migrar". O arquivo nasceria pedindo migração para um
	// formato em que ele já está.
	g.Version = FormatoAtual
	g.GeradoPor = GeneratedBy
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
	if equalIgnoringGeneratedBy(path, completo) {
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
	// O FORMATO é conferido ANTES de o mapa ser usado, e não depois.
	//
	// Um arquivo que este binário não sabe ler não deve ser parcialmente interpretado: os
	// campos que ele reconhece carregam, os que não reconhece somem, e a próxima gravação
	// escreve o que sobrou. É assim que se perde dado sem nada acusar — e o `julgamentos`
	// renomeado para `judgments` é o exemplo: os carimbos evaporariam e o `check` recobraria
	// o que alguém já respondeu.
	if err := ConfereFormato(path, g.Version); err != nil {
		return nil, err
	}
	return &g, nil
}

// equalIgnoringGeneratedBy diz se o mapa em disco é o mesmo, desconsiderando quem o gerou.
func equalIgnoringGeneratedBy(path string, novo []byte) bool {
	atual, err := os.ReadFile(path)
	if err != nil {
		return false // não existe ainda: há o que escrever
	}
	return withoutGeneratedBy(atual) == withoutGeneratedBy(novo)
}

var generatedByLineRE = regexp.MustCompile(`(?m)^generated_by:.*\n`)

func withoutGeneratedBy(b []byte) string {
	return generatedByLineRE.ReplaceAllString(string(b), "")
}
