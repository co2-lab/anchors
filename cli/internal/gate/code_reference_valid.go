package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// code-reference-valid: um código citado por uma spec precisa EXISTIR.
//
// A âncora que mente é a falha que o framework existe para impedir, e ela tem uma forma
// que nenhum outro gate pegava: a spec cita `XXXXX-B07` de outra unidade — numa Tabela de
// Dependências, numa nota, numa referência cruzada — e essa unidade não existe em lugar
// nenhum. A citação parece rastreabilidade e não é: aponta para o vazio.
//
// Aconteceu, medido: uma spec de schema afirmava "Índices que o schema criou (2026-08-11)"
// e referenciava quatro códigos (`MTVRX-B01`, `MTVRX-B07`, `RDMDX-B03`, `RDMDX-B04`) de specs
// que não existiam. O arquivo que ela descreve não tinha uma linha da entidade. Passou por
// TODOS os gates: a spec tinha código, header, seções, e o gate de dependências não olha
// prosa. Um leitor futuro — humano ou agente — a lê como registro do que foi feito.
//
// Por que é diferente de `dependency-honored`: aquele confronta a Tabela de Dependências
// contra o CÓDIGO (o método prometido é usado?). Este confronta a citação contra o
// UNIVERSO DE IDENTIDADES (a unidade citada existe?). Uma spec pode citar corretamente um
// símbolo de um arquivo que existe, e ao mesmo tempo citar o código de uma spec que nunca
// foi escrita.
//
// O gate NÃO cobra o próprio código da spec (ela é dona dele) nem os códigos de cenário
// que ela define — só as referências a OUTRAS unidades.
func checkCodeReferenceValid(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a citação cruzada é da spec — é ela que referencia outras unidades"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate precisa do universo de identidades"
	}

	meu := codeDoHeader(content)
	donos := donosDeCodigo(g, root)
	if len(donos) == 0 {
		return Pending, "nenhuma identidade no mapa — rode `anchors map build`"
	}

	orfas := map[string]bool{}
	for _, m := range refCodeRE().FindAllStringSubmatch(content, -1) {
		code := m[1]
		if code == meu || donos[code] {
			continue
		}
		orfas[code] = true
	}
	if len(orfas) == 0 {
		return Pass, ""
	}

	var lista []string
	for c := range orfas {
		lista = append(lista, "`"+c+"`")
	}
	sort.Strings(lista)
	return Fail, fmt.Sprintf("cita %d código(s) que não existem no projeto: %s. "+
		"Uma citação assim PARECE rastreabilidade e aponta para o vazio — quem lê depois "+
		"(pessoa ou agente) a toma como registro do que foi feito. Ou a unidade citada "+
		"precisa nascer, ou a citação está errada e deve sair",
		len(lista), strings.Join(lista, ", "))
}

// refCodeRE casa uma referência a um requisito: `XXXXX-B07`. Não distingue definição de
// citação de propósito — o código que a própria spec define é o dela, e é filtrado pelo
// `meu`; qualquer outro é referência a uma unidade que precisa existir.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func refCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`\b([A-Z0-9]` + config.CodeLengthPattern() + `)-[A-Z]\d{2}\b`)
}

// headerCodeCaptureRE é o `headerCodeRE` (internal_checks.go) com o código CAPTURADO:
// aquele só detecta a presença do campo; aqui é preciso o valor. Mesmos prefixos de
// comentário, para valer nos mesmos dialetos de header.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func headerCodeCaptureRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*code:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
}

func codeDoHeader(content string) string {
	if m := headerCodeCaptureRE().FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

// donosDeCodigo varre as specs do mapa e coleta o código que cada uma DECLARA no header —
// o universo de identidades vivas do projeto. Lê do disco porque o mapa guarda o nó, não
// o conteúdo.
func donosDeCodigo(g *mapx.Graph, root string) map[string]bool {
	out := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Code != "" {
			out[n.Code] = true
			continue
		}
		if n.Kind != mapx.KindSpec {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, n.ID))
		if err != nil {
			continue
		}
		if c := codeDoHeader(string(b)); c != "" {
			out[c] = true
		}
	}
	return out
}
