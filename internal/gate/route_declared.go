package gate

import (
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// route-declared: uma spec de TELA (layer: screen) declara COMO se chega até ela — a
// rota nomeada — e nomeia as telas vizinhas por nome concreto na navegação. É a régua
// de rastreabilidade de navegação: sem a rota, o grafo de telas fica com um nó solto;
// com termos genéricos ("Próxima tela", "Menu principal") a aresta de navegação não
// aponta para lugar nenhum.
//
// Só se aplica a `layer: screen`. Hooks, business-logic, stores, DAOs NÃO têm rota —
// cobrar rota deles (o vício do validador legado que só conhecia screen|component)
// gera falso-positivo. Para qualquer outra camada este checker é Skip.
//
// Substitui as regras `header-route` + `navigation-naming` do `validate-specs.ts` do
// app de referência, agora com consciência de camada (o legado tratava todo não-componente como tela).

// routeLineRE casa a linha de rota do header em prosa: `> **Rota**: ` ou `> **Route**: ` + código entre
// crases não-vazio. Aceita variação de espaço (o header é escrito à mão).
var routeLineRE = regexp.MustCompile("(?m)^>\\s*\\*\\*(?:Rota|Route)\\*\\*:\\s*`[^`]+`")

// navGenericTermRE são os termos genéricos proibidos nas seções de navegação — a
// navegação deve nomear a tela concreta (ex.: `HomeScreen`), não um rótulo vago.
var navGenericTermRE = []*regexp.Regexp{
	regexp.MustCompile(`(?i)Tela\s+de\s+`),
	regexp.MustCompile(`(?i)Screen\s+of\s+`),
	regexp.MustCompile(`(?i)Próxima\s+tela`),
	regexp.MustCompile(`(?i)Next\s+screen`),
	regexp.MustCompile(`(?i)Menu\s+principal`),
	regexp.MustCompile(`(?i)Main\s+menu`),
}

// navHeadingRE is the heading of a navigation section (Entrada/Saída, In/Out, Navigation…),
// the name as a whole word.
var navHeadingRE = regexp.MustCompile(`(?i)^###\s+(?:Entrada|Sa[íi]da|Entry|Exit|In|Out|Incoming|Outgoing|Navigation)\b`)

// navSections returns the body of every navigation section, each up to the next `###`.
//
// A single regex ending a match at `\n###` CONSUMED that heading, so a navigation section
// right after another (`### Entrada` then `### Saída`) was never read, and a generic term
// in it passed. And `In`/`Out` without a word boundary took `### Integração` or
// `### Outline` for navigation.
func navSections(content string) []string {
	var out []string
	var cur *strings.Builder
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "###") {
			if cur != nil {
				out = append(out, cur.String())
				cur = nil
			}
			if navHeadingRE.MatchString(strings.TrimSpace(line)) {
				cur = &strings.Builder{}
			}
			continue
		}
		if cur != nil {
			cur.WriteString(line)
			cur.WriteString("\n")
		}
	}
	if cur != nil {
		out = append(out, cur.String())
	}
	return out
}

func checkRouteDeclared(content string, n mapx.Node) (Verdict, string) {
	if layerOf(n, content) != "screen" {
		// Skip COM motivo: o contador `~1` sozinho deixa quem lê em dúvida se é
		// problema dele. Dizer "não é tela" fecha a dúvida em uma linha.
		return Skip, i18n.T("gate.route_declared.skip_not_screen")
	}
	if !routeLineRE.MatchString(content) {
		return Fail, i18n.T("gate.route.missing")
	}
	// navigation-naming: nas seções Entrada/Saída, proibir termos genéricos.
	for _, section := range navSections(content) {
		for _, line := range strings.Split(section, "\n") {
			if !strings.HasPrefix(strings.TrimSpace(line), "|") || strings.Contains(line, "---") {
				continue // só linhas de tabela de navegação
			}
			for _, term := range navGenericTermRE {
				if term.MatchString(line) {
					return Fail, i18n.T("gate.route.generic")
				}
			}
		}
	}
	return Pass, ""
}

// layerOf devolve o valor do `layer:` declarado no header do arquivo (fonte da verdade
// da identidade), caindo para as tags de camada do nó quando o header é omisso.
func layerOf(n mapx.Node, content string) string {
	if m := headerLayerValueRE.FindStringSubmatch(content); m != nil {
		return strings.TrimSpace(m[1])
	}
	for _, t := range n.Tags {
		if t == "screen" {
			return "screen"
		}
	}
	return ""
}
