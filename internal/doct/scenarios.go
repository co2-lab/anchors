package doct

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/mapx"
)

// --- os CENÁRIOS entram na documentação ---
//
// A feature é o único artefato da trinca escrito para ser lido por gente que não programa:
// Gherkin existe para isso. E é o que responde a pergunta que a spec não responde — a spec
// diz a REGRA ("dívida aberta bloqueia"), o cenário diz o que acontece ("dado uma dívida
// aberta, quando confiro a régua, então não libera").
//
// Deixá-lo fora da documentação seria manter o comportamento observável do sistema num
// arquivo `.feature` que só o time de desenvolvimento abre — enquanto a documentação, que
// é para os outros, descreve regras em abstrato.
//
// O cenário chega aqui pela ÂNCORA: a feature é irmã da spec no mapa (aresta `covered-by`),
// e é assim que a doc de uma unidade sabe quais cenários são dela.

// Scenario é um cenário de feature, como a documentação precisa dele.
type Scenario struct {
	Code   string   // `GLCGL-B01` ou `GLCGL-B01#02`
	Titulo string   // o texto depois de `Cenário:`
	Tags   []string // as tags da tag-line, sem o `@` (regime, marcadores do projeto)
	Corpo  string   // os passos Gherkin, como escritos
	Spec   string   // o código da spec de onde ele vem
}

var (
	// A tag-line: `@GLCGL-B01 @nivel-unit`. O primeiro código é a identidade.
	cenarioTagRE = regexp.MustCompile(`@([A-Za-z0-9][A-Za-z0-9_#-]*)`)
	// O título aceita as duas grafias — a feature pode estar em `# language: pt` ou não,
	// e uma documentação que ignora metade dos cenários por causa do idioma do arquivo
	// não é a documentação de ninguém.
	cenarioTituloRE = regexp.MustCompile(`(?i)^\s*(?:Cen[áa]rio|Scenario)(?:\s+Outline|\s+Esquema.*?)?:\s*(.+?)\s*$`)
	// O código de identidade tem a forma `ABCDE-B01` (com `#NN` opcional).
	cenarioCodigoRE = regexp.MustCompile(`^[A-Z0-9]{4,6}-(?:[A-Z]{1,2}\d{2}|DS-[A-Za-z0-9-]+|VR)(?:#\d{2})?$`)
)

// parseScenarios lê os cenários de uma feature.
//
// Captura o CORPO junto, e é o que distingue esta leitura da do gate: o gate confronta o
// cenário contra o teste e só precisa do código; a documentação precisa dos passos, porque
// é o corpo que descreve o comportamento — o título sozinho é uma manchete.
func parseScenarios(content, specCode string) []Scenario {
	linhas := strings.Split(content, "\n")
	var out []Scenario
	var tagsPend []string
	var codePend string
	atual := -1

	for _, ln := range linhas {
		trim := strings.TrimSpace(ln)

		if strings.HasPrefix(trim, "@") {
			tagsPend, codePend = nil, ""
			for _, m := range cenarioTagRE.FindAllStringSubmatch(trim, -1) {
				if codePend == "" && cenarioCodigoRE.MatchString(m[1]) {
					codePend = m[1]
					continue
				}
				tagsPend = append(tagsPend, m[1])
			}
			continue
		}

		if m := cenarioTituloRE.FindStringSubmatch(ln); m != nil {
			out = append(out, Scenario{
				Code: codePend, Titulo: m[1], Tags: tagsPend, Spec: specCode,
			})
			atual = len(out) - 1
			tagsPend, codePend = nil, ""
			continue
		}

		// Um heading de Gherkin (`Funcionalidade:`, `Contexto:`) encerra o cenário: os
		// passos de um `Contexto:` não pertencem ao cenário que veio antes dele.
		if atual >= 0 && endsScenario(trim) {
			atual = -1
			continue
		}
		if atual >= 0 && trim != "" {
			out[atual].Corpo += ln + "\n"
		}
	}
	for i := range out {
		out[i].Corpo = strings.TrimRight(out[i].Corpo, "\n")
	}
	return out
}

var headingGherkinRE = regexp.MustCompile(`(?i)^\s*(?:Funcionalidade|Feature|Contexto|Background|Regra|Rule):`)

func endsScenario(trim string) bool { return headingGherkinRE.MatchString(trim) }

// fnScenarios devolve os cenários da feature ligada a uma spec.
//
// Parte da SPEC porque é dela que a documentação parte: o template pede "os cenários desta
// unidade", e a unidade é identificada pela spec. A feature é encontrada pela aresta do
// mapa — não por convenção de nome, que quebraria no primeiro projeto que a organizasse
// de outro jeito.
func (c *Compiler) fnScenarios(s Spec) []Scenario {
	var out []Scenario
	for _, f := range c.featuresOf(s) {
		b, err := os.ReadFile(filepath.Join(c.Root, f))
		if err != nil {
			continue
		}
		out = append(out, parseScenarios(string(b), s.Code)...)
	}
	return out
}

// featuresOf devolve as features ligadas a uma spec, pelo mapa.
func (c *Compiler) featuresOf(s Spec) []string {
	var out []string
	for _, e := range c.Graph.Edges {
		if e.From != s.Path && e.To != s.Path {
			continue
		}
		outro := e.To
		if e.To == s.Path {
			outro = e.From
		}
		if c.kindOf(outro) == mapx.KindFeature {
			out = append(out, outro)
		}
	}
	return out
}

func (c *Compiler) kindOf(id string) mapx.Kind {
	for _, n := range c.Graph.Nodes {
		if n.ID == id {
			return n.Kind
		}
	}
	return ""
}

// fnAllScenarios devolve os cenários de TODAS as specs de um recorte.
//
// É a visão "todos os comportamentos do sistema", irmã de `rules`: quem lê a documentação
// para saber o que o sistema faz quer a lista, não uma navegação unidade por unidade.
func (c *Compiler) fnAllScenarios(filtro ...string) ([]Scenario, error) {
	specs, err := c.fnSpecs(filtro...)
	if err != nil {
		return nil, err
	}
	var out []Scenario
	for _, s := range specs {
		out = append(out, c.fnScenarios(s)...)
	}
	return out, nil
}
