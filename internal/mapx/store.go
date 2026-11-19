// @anchors
//   code: STMPA
//   ref: GRPRG

package mapx

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
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
	data, err := marshalGraph(g)
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
	//
	// A NEWER RELEASE is the exception: it stamps itself even when nothing else changed.
	// Otherwise the map kept the release that last changed it, and the `check` warning
	// ("written by 0.1.248, this is 0.1.258 — run `map build` with the new one") could not
	// be cleared: `map build` wrote nothing, so the warning stayed. An older release never
	// restamps, and neither does `dev`, which no release orders against — so the field
	// still cannot oscillate.
	//
	// The file on disk is read once, and compared as bytes: a project's map of 30 MB was
	// read twice and copied into strings four times on every save.
	if atual, err := os.ReadFile(path); err == nil && sameIgnoringGeneratedBy(atual, completo) &&
		!newerRelease(GeneratedBy, generatedByIn(atual)) {
		return nil
	}
	return os.WriteFile(path, completo, 0o644)
}

// marshalGraph writes the graph as yaml.Marshal writes it, a node and an edge at a time.
// yaml.v3's emitter keeps every event of a document until the document ends, so a map
// marshalled whole held gigabytes: a project's 30 MB map allocated 3.5 GB to be saved,
// and `anchors test` held 3.2 GB after its runner ended. Each item is marshalled in the
// same place of the document — the list under its key —, so the bytes are the same.
func marshalGraph(g *Graph) ([]byte, error) {
	var out bytes.Buffer
	head, err := yaml.Marshal(struct {
		Version   int    `yaml:"version"`
		GeradoPor string `yaml:"generated_by,omitempty"`
		ConfigRev string `yaml:"config_rev,omitempty"`
	}{g.Version, g.GeradoPor, g.ConfigRev})
	if err != nil {
		return nil, err
	}
	out.Write(head)
	if err := marshalList(&out, "nodes", len(g.Nodes), func(i int) any {
		return struct {
			L []Node `yaml:"nodes"`
		}{g.Nodes[i : i+1]}
	}); err != nil {
		return nil, err
	}
	if err := marshalList(&out, "edges", len(g.Edges), func(i int) any {
		return struct {
			L []Edge `yaml:"edges"`
		}{g.Edges[i : i+1]}
	}); err != nil {
		return nil, err
	}
	if g.Flow != nil {
		flow, err := yaml.Marshal(struct {
			Flow *FlowGraph `yaml:"flow,omitempty"`
		}{g.Flow})
		if err != nil {
			return nil, err
		}
		out.Write(flow)
	}
	return out.Bytes(), nil
}

// marshalList writes `key:` and then each item, marshalled alone under the same key, with
// the key's line dropped; an empty list is written as yaml.Marshal writes it.
func marshalList(out *bytes.Buffer, key string, n int, one func(int) any) error {
	if n == 0 {
		out.WriteString(key + ": []\n")
		return nil
	}
	out.WriteString(key + ":\n")
	prefix := []byte(key + ":\n")
	for i := 0; i < n; i++ {
		b, err := yaml.Marshal(one(i))
		if err != nil {
			return err
		}
		if !bytes.HasPrefix(b, prefix) {
			return fmt.Errorf("map: %s[%d] did not marshal as a list item", key, i)
		}
		out.Write(b[len(prefix):])
	}
	return nil
}

// sameIgnoringGeneratedBy says whether two maps are the same but for who generated them.
// The field is a line of the head, before the first list.
func sameIgnoringGeneratedBy(a, b []byte) bool {
	ha, ra := splitHead(a)
	hb, rb := splitHead(b)
	return bytes.Equal(ra, rb) && bytes.Equal(generatedByLineRE.ReplaceAll(ha, nil), generatedByLineRE.ReplaceAll(hb, nil))
}

// splitHead parts a map at its first list: the head holds the scalar fields.
func splitHead(b []byte) ([]byte, []byte) {
	if i := bytes.Index(b, []byte("\nnodes:")); i >= 0 {
		return b[:i+1], b[i+1:]
	}
	return b, nil
}

// Load lê um grafo existente do disco.
func Load(path string) (*Graph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadBytes(data, path)
}

// LoadBytes reads a map from its bytes — one git keeps, for instance — with the same
// format check as Load; `name` says where it came from in the errors.
func LoadBytes(data []byte, name string) (*Graph, error) {
	path := name
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
	dropEmptySignals(&g)
	return &g, nil
}

var generatedByLineRE = regexp.MustCompile(`(?m)^generated_by:.*\n`)

// generatedByOnDisk is the release the map on disk says wrote it; empty when there is none.
func generatedByOnDisk(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return generatedByIn(b)
}

// generatedByIn is the release a map's bytes say wrote it; empty when they say none.
func generatedByIn(b []byte) string {
	head, _ := splitHead(b)
	m := generatedByLineRE.Find(head)
	if m == nil {
		return ""
	}
	return strings.Trim(strings.TrimSpace(strings.TrimPrefix(string(m), "generated_by:")), `"'`)
}

// newerRelease says whether `running` is a release later than `stored`. Either one not
// being a release (`dev`, empty) answers no.
func newerRelease(running, stored string) bool {
	o, err := config.VersionOrder(running, stored)
	return err == nil && o > 0
}

// dropEmptySignals turns a signal that holds nothing into no signal. A signal means
// something was measured, and the gates read its presence that way; an empty one — the
// record of kept evidence lived in the signal for two releases, and a file with nothing
// measured got one — made scenario-coverage read "execution ingested, nothing proven".
func dropEmptySignals(g *Graph) {
	for i := range g.Nodes {
		if s := g.Nodes[i].Signal; s != nil && reflect.DeepEqual(*s, TestSignal{}) {
			g.Nodes[i].Signal = nil
		}
	}
}
