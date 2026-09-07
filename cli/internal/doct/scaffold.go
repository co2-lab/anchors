package doct

import (
	"os"
	"path/filepath"
)

// --- a ESTRUTURA da documentação ---
//
// As seções não são uma lista solta: elas respondem a perguntas diferentes, de leitores
// diferentes, e a organização é o que faz cada leitor achar a sua sem ler as outras.
//
//	docs/
//	  produto.md          O QUE o sistema faz e para quem — vem dos planos, escrito à mão
//	  arquitetura.md      COMO ele é montado — C4, os quatro níveis
//	  comportamento.md    O QUE ACONTECE — os cenários das features, todos
//	  regras.md           AS REGRAS — todas as specs, por seção (a visão-matriz A)
//	  camadas/*.md        POR ONDE — uma página por camada (a visão-matriz B)
//	  contratos/*.md      as docs específicas do tipo de projeto (OpenAPI, esquema…)
//
// A MATRIZ é a razão de haver duas visões dos mesmos dados. Uma spec pertence a uma camada
// E tem seções; quem pergunta "todas as regras do sistema" quer o corte por seção, quem
// pergunta "o que existe em screens" quer o corte por camada. Uma hierarquia só serviria
// a um dos dois, e o outro teria de navegar unidade por unidade.
//
// A duplicação entre as duas visões é real e é DELIBERADA: ela existe só no compilado, que
// é gerado a cada build e conferido pelo `docs-fresh`. A duplicação que dói é a escrita à
// mão, que envelhece sem ninguém ver.
//
// A ORDEM das páginas é de fora para dentro — produto, arquitetura, comportamento, regras,
// camadas. É a ordem em que alguém que chega precisa delas, e não a ordem em que o time
// as escreveu.

// Scaffold é um template inicial que o `anchors docs init` escreve.
type Scaffold struct {
	Nome  string
	Corpo string
	// Porque explica ao autor o que a página responde — vai como comentário no topo do
	// template, onde quem for editá-lo vai ler.
	Porque string
}

// Scaffolds são as páginas que o `anchors docs init` propõe.
//
// PROPÕE, não impõe: são templates comuns, e o projeto os edita. O que o Anchors garante
// é que ninguém precise inventar a organização do zero — inventar a cada projeto é como
// se chega a uma documentação onde cada página segue uma lógica diferente.
func Scaffolds() []Scaffold {
	return []Scaffold{
		{
			Nome:   "arquitetura.md" + SufixoTemplate,
			Porque: "COMO o sistema é montado — o C4, em quatro níveis",
			Corpo: `# Arquitetura

> Os níveis 1 e 2 do C4 são escritos à mão AQUI no template: eles descrevem o sistema
> inteiro, e nenhuma spec sozinha os conhece. Os níveis 3 e 4 vêm das specs.

## Nível 1 — Contexto

<!-- Quem usa o sistema, e com que sistemas externos ele fala. -->

## Nível 2 — Contêineres

<!-- O que roda separado, e o protocolo de cada conversa. -->

## Nível 3 — Componentes

As camadas do projeto, e o que existe em cada uma.

{{range layers}}
### {{.}}
{{range specs (printf "layer=%s" .)}}
- **{{.Code}}** — {{.Titulo}}
{{end}}
{{end}}
`,
		},
		{
			Nome:   "comportamento.md" + SufixoTemplate,
			Porque: "O QUE ACONTECE — os cenários das features, que é o comportamento observável",
			Corpo: `# Comportamento

Cada cenário descreve o que o sistema faz numa situação. Vêm das features, e são os
mesmos que os testes provam — não uma descrição paralela deles.
{{range layers}}
## {{.}}
{{range $c := allScenarios (printf "layer=%s" .)}}
### {{$c.Titulo}}  ` + "`{{$c.Code}}`" + `

` + "```gherkin" + `
{{$c.Corpo}}
` + "```" + `
{{else}}
_Nenhum cenário ainda: as unidades desta camada não têm feature._
{{end}}{{end}}
`,
		},
		{
			Nome:   "regras.md" + SufixoTemplate,
			Porque: "AS REGRAS de todo o sistema, por seção — o corte que atravessa as camadas",
			Corpo: `# Regras

Todas as regras do sistema, de todas as unidades. Para ver o que existe em uma camada
específica, veja ` + "`camadas/`" + `.

{{range specs}}
## {{.Code}} — {{.Titulo}}

{{section . "Visão Geral"}}

{{range rules .}}
### {{.Code}} — {{.Titulo}}

{{.Corpo}}
{{end}}
{{end}}
`,
		},
	}
}

// ScaffoldLayer é a página de UMA camada — a outra metade da matriz.
//
// Uma por camada, e não uma página com todas: a página por camada é a que alguém abre
// para saber o que existe em `screens`, e uma página única faria essa pessoa rolar por
// tudo o que não é screen.
func ScaffoldLayer(camada string) Scaffold {
	return Scaffold{
		Nome:   "camadas/" + camada + ".md" + SufixoTemplate,
		Porque: "o que existe na camada `" + camada + "`, com o conteúdo de cada unidade",
		Corpo: `# Camada: ` + camada + `

{{range specs "layer=` + camada + `"}}
## {{.Code}} — {{.Titulo}}

{{section . "Visão Geral"}}

{{/* As seções Regras e Invariantes trazem os próprios headings de regra, e por isso
     NÃO ganham um rótulo aqui: um heading "Regras" seguido dos headings das regras as
     põe no mesmo nível, e o índice do documento lista o rótulo como irmão do que ele
     contém. */}}
{{section . "Regras"}}

{{section . "Invariantes"}}
{{with scenarios .}}
### Cenários desta unidade
{{range .}}
- ` + "`{{.Code}}`" + ` {{.Titulo}}
{{end}}{{end}}
{{end}}
`,
	}
}

// InitScaffolds escreve os templates iniciais, inclusive uma página por camada existente.
//
// NÃO sobrescreve sem `--force`: um template já editado carrega a moldura que o time
// escreveu, e refazê-la a cada `init` seria apagar exatamente a parte que não é gerada.
func (c *Compiler) InitScaffolds(force bool) (escritos, pulados []string, err error) {
	todos := Scaffolds()
	for _, l := range c.fnLayers() {
		todos = append(todos, ScaffoldLayer(l))
	}
	for _, s := range todos {
		destino := filepath.Join(c.Root, Dir, s.Nome)
		if _, e := os.Stat(destino); e == nil && !force {
			pulados = append(pulados, s.Nome)
			continue
		}
		if e := os.MkdirAll(filepath.Dir(destino), 0o755); e != nil {
			return escritos, pulados, e
		}
		corpo := "{{/* " + s.Porque + " */}}\n" + s.Corpo
		if e := os.WriteFile(destino, []byte(corpo), 0o644); e != nil {
			return escritos, pulados, e
		}
		escritos = append(escritos, s.Nome)
	}
	return escritos, pulados, nil
}
