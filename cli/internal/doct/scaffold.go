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
//	  arquitetura.md      COMO ele é montado — C4 em quatro níveis, com diagrama Mermaid
//	  comportamento.md    ÍNDICE dos cenários → a página da camada
//	  regras.md           ÍNDICE das regras   → a página da camada
//	  camadas/*.md        O CONTEÚDO — uma página por camada
//	  contratos/*.md      as docs específicas do tipo de projeto (OpenAPI, esquema…)
//
// A MATRIZ é a razão de haver duas visões dos mesmos dados. Uma spec pertence a uma camada
// E tem seções; quem pergunta "todas as regras do sistema" quer o corte por seção, quem
// pergunta "o que existe em screens" quer o corte por camada. Uma hierarquia só serviria
// a um dos dois, e o outro teria de navegar unidade por unidade.
//
// O CONTEÚDO MORA NUM LUGAR SÓ, e os cortes transversais são ÍNDICES.
//
// A primeira versão repetia o texto nas duas visões, e a medida no projeto de referência
// mostrou por que não serve: `regras.md` saiu com 9.608 linhas — o corte transversal de 84
// unidades é o documento inteiro, mais uma vez. E ele cresce com o projeto, sem limite.
//
// Um índice não tem esse problema: uma linha por regra, e o link leva ao texto na página
// da camada. O preço é o clique, e ele é justo — quem abre "todas as regras" procura UMA,
// e antes tinha de rolar por todas.
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

// crase evita o problema que derrubou este arquivo duas vezes: uma crase dentro de uma
// raw string do Go a FECHA, e o resto do template vira código Go inválido.
const crase = "`"

// cerca abre e fecha um bloco de código markdown.
func cerca(lang string) string { return "```" + lang }

// Scaffolds são as páginas que o `anchors docs init` propõe.
//
// PROPÕE, não impõe: são templates comuns, e o projeto os edita. O que o Anchors garante
// é que ninguém precise inventar a organização do zero — inventar a cada projeto é como
// se chega a uma documentação onde cada página segue uma lógica diferente, e quem lê tem
// de descobrir a lógica antes de achar o que procura.
func Scaffolds() []Scaffold {
	return []Scaffold{
		{
			Nome:   "arquitetura.md" + SufixoTemplate,
			Porque: "COMO o sistema é montado — o C4 em quatro níveis, com diagrama Mermaid",
			Corpo: `# Arquitetura

{{/* MERMAID, e não uma imagem.

     O GitHub renderiza blocos de diagrama Mermaid nativamente, e o MkDocs e o Starlight
     também. Um PNG exportado de uma ferramenta de desenho ficaria fora do controle de
     versão útil: o diff não diz o que mudou, e o arquivo-fonte do desenho acaba noutro
     lugar — ou some. Aqui o diagrama É texto, versionado com o resto.

     Os NÍVEIS 1 e 2 são escritos à mão AQUI no template: descrevem o sistema inteiro, e
     nenhuma spec sozinha os conhece. O nível 3 vem das camadas e do mapa. */}}

## Nível 1 — Contexto

Quem usa o sistema, e com que sistemas externos ele fala.

` + cerca("mermaid") + `
graph TB
    user["Pessoa<br/><small>quem usa o sistema</small>"]
    sys["O SISTEMA<br/><small>o que este repositório constrói</small>"]
    ext["Sistema externo<br/><small>de onde vêm os dados</small>"]

    user -->|"usa"| sys
    sys -->|"consulta"| ext

    classDef pessoa fill:#08427b,stroke:#052e56,color:#fff
    classDef sistema fill:#1168bd,stroke:#0b4884,color:#fff
    classDef externo fill:#999,stroke:#6b6b6b,color:#fff
    class user pessoa
    class sys sistema
    class ext externo
` + cerca("") + `

<!-- FORA deste nível: como o sistema é montado por dentro. Isso é o nível 2. -->

## Nível 2 — Contêineres

O que roda separado, e o protocolo de cada conversa. **O protocolo importa**: é ele que
diz o que acontece quando a conversa falha.

` + cerca("mermaid") + `
graph TB
    user["Pessoa"]

    subgraph sistema["O SISTEMA"]
        app["App<br/><small>a interface</small>"]
        api["API<br/><small>as rotas</small>"]
        db[("Banco<br/><small>o que persiste</small>")]
    end

    ext["Sistema externo"]

    user -->|"usa"| app
    app -->|"HTTPS/JSON"| api
    api -->|"lê e escreve"| db
    api -->|"consulta"| ext

    classDef pessoa fill:#08427b,stroke:#052e56,color:#fff
    classDef conteiner fill:#438dd5,stroke:#2e6295,color:#fff
    classDef externo fill:#999,stroke:#6b6b6b,color:#fff
    class user pessoa
    class app,api,db conteiner
    class ext externo
` + cerca("") + `

<!-- FORA deste nível: as peças DENTRO de cada contêiner. Isso é o nível 3. -->

## Nível 3 — Componentes

As camadas do projeto, e o que existe em cada uma. Esta seção é **gerada**: a Estrutura
declara as camadas, e as specs declaram as unidades.

{{/* AS SETAS vêm do MAPA — das arestas ` + crase + `needs` + crase + ` que as unidades declaram e o gate
     confere. Desenhá-las à mão seria garantir que envelheçam: é no nível 3 que uma
     dependência nova aparece primeiro, e ninguém volta ao diagrama para acrescentá-la.

     As fronteiras da Estrutura não serviriam: elas declaram o que uma camada NÃO pode
     alcançar, e a ausência de proibição não é uma dependência. */}}
` + cerca("mermaid") + `
graph LR
{{range layers}}    {{mermaidID .}}["{{.}}<br/><small>{{(size (printf "layer=%s" .)).Units}} unidades</small>"]
{{end}}{{range layerDeps}}    {{mermaidID .From}} -->|"{{.Count}}"| {{mermaidID .To}}
{{end}}
` + cerca("") + `
{{if not layerDeps}}
_Nenhuma dependência entre camadas está declarada no mapa. As caixas acima são o que
existe; as setas aparecem quando as unidades declararem dependência._
{{end}}
{{range $l := layers}}
### {{$l}}

{{$g := size (printf "layer=%s" $l)}}{{$g.Units}} unidades, {{$g.Rules}} regras. O
conteúdo está em [{{layerPage $l}}]({{layerPage $l}}).

{{range specs (printf "layer=%s" $l)}}- **{{.Code}}** — {{.Titulo}}
{{end}}{{end}}

## Nível 4 — Código

<!-- Só onde houver algo não óbvio. O mapa que o Anchors mantém já é a versão de
     máquina, e repeti-lo aqui seria duplicação sem leitor. -->

_Nada a destacar por enquanto._
`,
		},
		{
			Nome:   "comportamento.md" + SufixoTemplate,
			Porque: "O QUE ACONTECE — índice dos cenários, que levam ao conteúdo na página da camada",
			Corpo: `# Comportamento

Todos os cenários do sistema. Cada um leva à unidade que o define.

Um cenário descreve o que o sistema faz numa situação — vem da feature, e é o mesmo que o
teste prova.
{{range $l := layers}}
## {{$l}}
{{range allScenarios (printf "layer=%s" $l)}}
- [{{.Titulo}}]({{scenarioLink $l .}}) ` + crase + `{{.Code}}` + crase + `
{{else}}
_Nenhum cenário ainda: as unidades desta camada não têm feature._
{{end}}{{end}}
`,
		},
		{
			Nome:   "regras.md" + SufixoTemplate,
			Porque: "AS REGRAS de todo o sistema — índice que atravessa as camadas",
			Corpo: `# Regras

Todas as regras do sistema, de todas as unidades. Cada uma leva ao texto na página da sua
camada.

Para ver uma camada inteira de uma vez — com a visão geral de cada unidade e os cenários —
abra a página dela em ` + crase + `camadas/` + crase + `.
{{range $l := layers}}
## {{$l}}
{{range $s := specs (printf "layer=%s" $l)}}
### [{{$s.Code}} — {{$s.Titulo}}]({{layerPage $l}}#{{anchor (printf "%s — %s" $s.Code $s.Titulo)}})
{{range rules $s}}
- [{{.Code}} — {{.Titulo}}]({{ruleLink $l $s .}})
{{end}}{{end}}{{end}}
`,
		},
	}
}

// ScaffoldLayer é a página de UMA camada — onde o conteúdo mora.
//
// Uma por camada, e não uma página com todas: é a que alguém abre para saber o que existe
// em `screens`, e uma página única faria essa pessoa rolar por tudo o que não é screen.
func ScaffoldLayer(camada string) Scaffold {
	f := `"layer=` + camada + `"`
	return Scaffold{
		Nome:   "camadas/" + camada + ".md" + SufixoTemplate,
		Porque: "o que existe na camada `" + camada + "`, com o conteúdo de cada unidade",
		Corpo: `# Camada: ` + camada + `

{{/* O FORMATO vem do LAYOUT, decidido uma vez no ` + crase + `docs build` + crase + ` e igual para todas
     as páginas — esta, o índice de regras e o de comportamento.

     Não há limiar escrito aqui, de propósito. A primeira versão deixava cada template
     escolher o seu, e o resultado foi medido: 483 de 812 links quebrados, porque a
     página resumia por um critério e os links eram montados por outro.

     Para mudar o corte deste projeto:  anchors docs build --max-units 30

     O que muda é o que CABE na página, nunca em quantos arquivos a camada se parte:
     dividir ao cruzar um limiar quebraria todo link externo no dia em que a unidade
     seguinte entrasse. */}}
{{if big ` + f + `}}
> {{layout.Describe (size ` + f + `)}}

{{range specs ` + f + `}}
## {{.Code}} — {{.Titulo}}

{{section . "Visão Geral"}}

{{range rules .}}
- **{{.Code}}** — {{.Titulo}}
{{end}}
{{end}}
{{else}}
{{range specs ` + f + `}}
## {{.Code}} — {{.Titulo}}

{{section . "Visão Geral"}}

{{/* As seções Regras e Invariantes trazem os próprios headings de regra, e por isso
     NÃO ganham um rótulo aqui: um heading "Regras" seguido dos headings das regras as
     põe no mesmo nível, e o índice do documento lista o rótulo como irmão do que ele
     contém. */}}
{{section . "Regras"}}

{{section . "Invariantes"}}
{{/* Cada cenário é um HEADING, e não um item de lista: é para ele que o índice de
     comportamento aponta, e âncora só existe onde há heading. Um índice cujo link não
     resolve é pior que não ter índice — ele parece funcionar. */}}
{{range scenarios .}}
#### {{.Code}} — {{.Titulo}}

` + cerca("gherkin") + `
{{.Corpo}}
` + cerca("") + `
{{end}}
{{end}}
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
