// @anchors
//   ref: DCSCD

package doct

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
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
//
// IN THE PROJECT'S LANGUAGE. The pages are read by the project's readers, so what they
// show — headings, paragraphs, diagram labels, the notes of an empty state — and the page
// names follow `lang:`. They were Portuguese whatever the project declared (reported from
// baas-proxy). The notes to whoever edits the template (`{{/* … */}}`) are in English, like
// everything else Anchors writes into a project's files. The skeleton is ONE, and the
// language is a table of texts over it: three copies of the template would drift apart on
// the first fix applied to one of them.
func Scaffolds(lang string) []Scaffold {
	t := scaffoldTextFor(lang)
	return []Scaffold{
		{Nome: t["arch.file"] + ".md" + SufixoTemplate, Porque: t["arch.why"], Corpo: t.fill(archBody)},
		{Nome: t["behavior.file"] + ".md" + SufixoTemplate, Porque: t["behavior.why"], Corpo: t.fill(behaviorBody)},
		{Nome: t["rules.file"] + ".md" + SufixoTemplate, Porque: t["rules.why"], Corpo: t.fill(rulesBody)},
	}
}

// ScaffoldLayer é a página de UMA camada — onde o conteúdo mora.
//
// Uma por camada, e não uma página com todas: é a que alguém abre para saber o que existe
// em `screens`, e uma página única faria essa pessoa rolar por tudo o que não é screen.
func ScaffoldLayer(camada, lang string) Scaffold {
	t := scaffoldTextFor(lang)
	f := `"layer=` + camada + `"`
	body := strings.NewReplacer("«filter»", f, "«layer»", camada).Replace(t.fill(layerBody))
	return Scaffold{
		Nome:   LayerDir(lang) + "/" + camada + ".md" + SufixoTemplate,
		Porque: strings.ReplaceAll(t["layer.why"], "«layer»", camada),
		Corpo:  body,
	}
}

// scaffoldText is the visible text of the skeleton in one language, by key.
type scaffoldText map[string]string

// fill puts the language's texts into a skeleton: `«key»` becomes the text, and the
// section titles become the language's titles from the section catalog.
func (t scaffoldText) fill(skeleton string) string {
	out := skeleton
	// Until nothing changes: a text may carry another key (`«layers.dir»` inside the
	// rules intro), and map order decides which is replaced first.
	for prev := ""; prev != out; {
		prev = out
		for k, v := range t {
			out = strings.ReplaceAll(out, "«"+k+"»", v)
		}
	}
	return strings.NewReplacer("``", crase, "«fence-mermaid»", cerca("mermaid"), "«fence-gherkin»", cerca("gherkin"), "«fence»", cerca("")).Replace(out)
}

// scaffoldTextFor is the table of a language, English for one with no table. The section
// titles come from the catalog the gates use, so the page asks for the sections the
// project's specs write.
func scaffoldTextFor(lang string) scaffoldText {
	base := scaffoldTexts["en"]
	l := "en"
	switch {
	case strings.HasPrefix(lang, "pt"):
		l = "pt-BR"
	case strings.HasPrefix(lang, "es"):
		l = "es"
	}
	t := scaffoldText{}
	for k, v := range base {
		t[k] = v
	}
	for k, v := range scaffoldTexts[l] {
		t[k] = v
	}
	for _, k := range []string{"overview", "rules", "effects", "invariants"} {
		t["section."+k] = i18n.TIn(l, "section.title."+k)
	}
	t["layers.dir"] = LayerDir(l)
	return t
}

const archBody = `# «arch.title»

{{/* THE C4 AS THE C4 IS.

     The model's central rule is that each level ZOOMS INTO ONE BOX of the previous one.
     Level 3 is the zoom of ONE container — not of the system. A component diagram mixing
     the app, the API and the infrastructure is level 3 of nothing: it is the failure the
     C4 exists to avoid, one drawing with everything inside.

     A CONTAINER is what runs or stores data — application, service, DATABASE, queue, file
     system. It is not a synonym of "a process we wrote": the database is a container, it
     shows at level 2 with the protocol of the conversation, and being a third party's it
     gets no level 3 — we have no components inside it.

     MERMAID, not an image: GitHub renders it natively. A PNG exported from a drawing tool
     would sit outside useful version control — the diff does not say what changed, and the
     source file ends up elsewhere, or is lost.

     Levels 1 and 2 are written by hand HERE: they describe the whole system, and no spec
     alone knows them. Level 3 comes from the containers declared in the configuration. */}}

## «arch.l1»

«arch.l1.desc»

«fence-mermaid»
graph TB
    user["«arch.person»<br/><small>«arch.person.desc»</small>"]
    sys["«arch.system»<br/><small>«arch.system.desc»</small>"]
    ext["«arch.external»<br/><small>«arch.external.desc»</small>"]

    user -->|"«arch.uses»"| sys
    sys -->|"«arch.queries»"| ext

    classDef person fill:#08427b,stroke:#052e56,color:#fff
    classDef system fill:#1168bd,stroke:#0b4884,color:#fff
    classDef external fill:#999,stroke:#6b6b6b,color:#fff
    class user person
    class sys system
    class ext external
«fence»

<!-- «arch.l1.out» -->

## «arch.l2»

«arch.l2.desc»

{{/* The boxes come from the containers declared in the configuration; the DATABASE is
     among them, because a container is what runs OR STORES data. The arrows and the
     protocols come from "talks", in the same declaration. */}}
«fence-mermaid»
graph TB
    user["«arch.person»"]

    subgraph sistema["«arch.system»"]
{{range containers}}{{if not .External}}        {{.ID}}["{{.Name}}<br/><small>{{.Description}}</small>"]
{{end}}{{end}}    end

{{range containers}}{{if .External}}    {{.ID}}["{{.Name}}<br/><small>{{.Description}}</small>"]
{{end}}{{end}}
{{range containers}}{{$de := .ID}}{{range .Talks}}    {{$de}} -->|"{{.Protocol}}{{if .Why}}<br/>{{.Why}}{{end}}"| {{mermaidID .To}}
{{end}}{{end}}
    classDef person fill:#08427b,stroke:#052e56,color:#fff
    classDef container fill:#438dd5,stroke:#2e6295,color:#fff
    classDef external fill:#999,stroke:#6b6b6b,color:#fff
    class user person
{{range containers}}{{if .External}}    class {{.ID}} external
{{else}}    class {{.ID}} container
{{end}}{{end}}
«fence»
{{if not containers}}
> «arch.no_containers»
{{end}}

<!-- «arch.l2.out» -->

## «arch.l3»

«arch.l3.desc»
{{range internalContainers}}
### {{.Name}}

{{.Description}}

{{if .Units}}
«fence-mermaid»
graph TB
    subgraph {{.ID}}["{{.Name}}"]
{{range .Layers}}        subgraph {{mermaidID .}}["{{.}}"]
{{range $u := specs (printf "layer=%s" .)}}            {{mermaidID $u.Code}}["{{$u.Code}}<br/><small>{{$u.Titulo}}</small>"]
{{else}}            {{mermaidID .}}_nospec["{{layerFiles .}} «arch.files_no_spec»"]
{{end}}        end
{{end}}    end
«fence»

{{range .Layers}}
#### {{.}}

{{range specs (printf "layer=%s" .)}}- **{{.Code}}** — {{.Titulo}}
{{else}}_«arch.layer_no_spec.a» {{layerFiles .}} «arch.layer_no_spec.b»_
{{end}}{{else}}
_«arch.container_no_layers»_
{{end}}
{{else}}{{with .Layers}}
_«arch.no_layer_has_spec»_

{{range .}}- **{{.}}** — {{layerFiles .}} «arch.files»
{{end}}{{else}}
_«arch.no_unit_runs_here»_
{{end}}{{end}}
{{end}}
{{with orphanLayers}}
### «arch.orphans»

«arch.orphans.desc»

{{range .}}- ` + "``" + `{{.}}` + "``" + `
{{end}}{{end}}

## «arch.l4»

<!-- «arch.l4.note» -->

_«arch.l4.empty»_
`

const behaviorBody = `# «behavior.title»

«behavior.intro»
{{range $l := layers}}
## {{$l}}
{{range allScenarios (printf "layer=%s" $l)}}
- [{{.Titulo}}]({{scenarioLink $l .}}) ` + "``" + `{{.Code}}` + "``" + `
{{else}}
_«behavior.none»_
{{end}}{{end}}
`

const rulesBody = `# «rules.title»

«rules.intro»
{{range $l := layers}}
## {{$l}}
{{range $s := specs (printf "layer=%s" $l)}}
### [{{$s.Code}} — {{$s.Titulo}}]({{layerPage $l}}#{{anchor (printf "%s — %s" $s.Code $s.Titulo)}})
{{range rules $s}}
- [{{.Code}} — {{.Titulo}}]({{ruleLink $l $s .}})
{{end}}{{end}}{{end}}
`

const layerBody = `# «layer.title» «layer»

{{/* The FORMAT comes from the LAYOUT, decided once in ` + "``" + `docs build` + "``" + ` and the same for every
     page — this one, the rules index and the behaviour index.

     There is no threshold written here, on purpose. The first version let each template
     pick its own, and the result was measured: 483 of 812 links broken, because the page
     summarised by one criterion and the links were built by another.

     To change this project's cut:  anchors docs build --max-units 30

     What changes is what FITS on the page, never how many files the layer splits into:
     splitting on crossing a threshold would break every external link the day the next
     unit came in. */}}
{{if big «filter»}}
> {{layout.Describe (size «filter»)}}

{{range specs «filter»}}
## {{.Code}} — {{.Titulo}}

{{section . "«section.overview»"}}

{{range rules .}}
- **{{.Code}}** — {{.Titulo}}
{{end}}
{{end}}
{{else}}
{{range specs «filter»}}
## {{.Code}} — {{.Titulo}}

{{section . "«section.overview»"}}

{{/* The rule sections carry their own rule headings, and so get NO label here: a
     heading over the rule headings puts them at the same level, and the document's table
     of contents lists the label as a sibling of what it holds. Both spec formats are
     asked for — the older one's rules section, the current one's effects — and a
     section the spec does not have is empty. */}}
{{section . "«section.rules»"}}

{{section . "«section.effects»"}}

{{section . "«section.invariants»"}}
{{/* Each scenario is a HEADING, not a list item: it is what the behaviour index points
     at, and an anchor only exists where there is a heading. An index whose link does not
     resolve is worse than no index — it looks like it works. */}}
{{range scenarios .}}
#### {{.Code}} — {{.Titulo}}

«fence-gherkin»
{{.Corpo}}
«fence»
{{end}}
{{end}}
{{end}}
`

// scaffoldTexts are the visible texts of the skeleton, per language. English is the base:
// a key a language lacks falls back to it.
var scaffoldTexts = map[string]scaffoldText{
	"en": {
		"arch.file":                "architecture",
		"arch.why":                 "HOW the system is built — the C4: context, containers, and one level 3 per container",
		"arch.title":               "Architecture",
		"arch.l1":                  "Level 1 — Context",
		"arch.l1.desc":             "The system as a single box: who uses it, and which external systems it talks to.",
		"arch.person":              "Person",
		"arch.person.desc":         "who uses the system",
		"arch.system":              "THE SYSTEM",
		"arch.system.desc":         "what this repository builds",
		"arch.external":            "External system",
		"arch.external.desc":       "where the data comes from",
		"arch.uses":                "uses",
		"arch.queries":             "queries",
		"arch.l1.out":              "OUTSIDE this level: how the system is built inside. That is level 2.",
		"arch.l2":                  "Level 2 — Containers",
		"arch.l2.desc":             "The zoom into the box \"THE SYSTEM\": what runs or stores data apart, and **which protocol each\npair talks**. The protocol is what says what happens when the conversation fails.",
		"arch.no_containers":       "**No container declared.** Level 2 and the level 3 ones stay empty until the configuration\n> declares what runs apart. See the containers block in ``anchors.yaml``.",
		"arch.l2.out":              "OUTSIDE this level: the pieces INSIDE each container. That is level 3 — and there is one\n     diagram per container, because each level zooms into ONE box of the previous one.",
		"arch.l3":                  "Level 3 — Components",
		"arch.l3.desc":             "One diagram **per container**: each zooms into one box of level 2. External containers do\nnot show here — we have no components inside them, and drawing them would claim a\nknowledge we do not have.",
		"arch.files_no_spec":       "file(s), no spec",
		"arch.layer_no_spec.a":     "No spec in this layer —",
		"arch.layer_no_spec.b":     "governed file(s).",
		"arch.container_no_layers": "This container declares no layer.",
		"arch.no_layer_has_spec":   "No layer of this container has a spec yet:",
		"arch.files":               "file(s)",
		"arch.no_unit_runs_here":   "No layer declares a unit that runs here. Either the container runs no code of this\nrepository, or it still has to be tied to a layer in ``containers.layers``.",
		"arch.orphans":             "Layers outside every container",
		"arch.orphans.desc":        "These layers exist in the project and no container declares them — they show in no level 3\ndiagram. Either they still have to be declared, or they run nowhere:",
		"arch.l4":                  "Level 4 — Code",
		"arch.l4.note":             "Only where there is something not obvious. The map Anchors keeps is already the machine\n     version, and repeating it here would be duplication with no reader.",
		"arch.l4.empty":            "Nothing to highlight for now.",
		"behavior.file":            "behavior",
		"behavior.why":             "WHAT HAPPENS — the index of the scenarios, which lead to the content on the layer's page",
		"behavior.title":           "Behavior",
		"behavior.intro":           "Every scenario of the system. Each one leads to the unit that defines it.\n\nA scenario describes what the system does in a situation — it comes from the feature, and it\nis what the test proves.",
		"behavior.none":            "No scenario yet: the units of this layer have no feature.",
		"rules.file":               "rules",
		"rules.why":                "THE RULES of the whole system — an index across the layers",
		"rules.title":              "Rules",
		"rules.intro":              "Every rule of the system, from every unit. Each one leads to its text on its layer's page.\n\nTo see a whole layer at once — with each unit's overview and its scenarios — open its page\nin ``«layers.dir»/``.",
		"layer.why":                "what exists in the layer `«layer»`, with the content of each unit",
		"layer.title":              "Layer:",
	},
	"pt-BR": {
		"arch.file":                "arquitetura",
		"arch.why":                 "COMO o sistema é montado — o C4: contexto, contêineres, e um nível 3 por contêiner",
		"arch.title":               "Arquitetura",
		"arch.l1":                  "Nível 1 — Contexto",
		"arch.l1.desc":             "O sistema como uma caixa só: quem o usa, e com que sistemas externos ele fala.",
		"arch.person":              "Pessoa",
		"arch.person.desc":         "quem usa o sistema",
		"arch.system":              "O SISTEMA",
		"arch.system.desc":         "o que este repositório constrói",
		"arch.external":            "Sistema externo",
		"arch.external.desc":       "de onde vêm os dados",
		"arch.uses":                "usa",
		"arch.queries":             "consulta",
		"arch.l1.out":              "FORA deste nível: como o sistema é montado por dentro. Isso é o nível 2.",
		"arch.l2":                  "Nível 2 — Contêineres",
		"arch.l2.desc":             "O zoom da caixa \"O SISTEMA\": o que roda ou armazena separado, e **com que protocolo cada\npar conversa**. O protocolo é o que diz o que acontece quando a conversa falha.",
		"arch.no_containers":       "**Nenhum contêiner declarado.** O nível 2 e os de nível 3 saem vazios até a Estrutura\n> declarar o que roda separado. Ver o bloco de contêineres no ``anchors.yaml``.",
		"arch.l2.out":              "FORA deste nível: as peças DENTRO de cada contêiner. Isso é o nível 3 — e há um\n     diagrama por contêiner, porque cada nível amplia UMA caixa do anterior.",
		"arch.l3":                  "Nível 3 — Componentes",
		"arch.l3.desc":             "Um diagrama **por contêiner**: cada um amplia uma caixa do nível 2. Os externos não\naparecem aqui — não temos componentes dentro deles, e desenhá-los afirmaria um\nconhecimento que não temos.",
		"arch.files_no_spec":       "arquivo(s), nenhuma spec",
		"arch.layer_no_spec.a":     "Nenhuma spec nesta camada —",
		"arch.layer_no_spec.b":     "arquivo(s) regido(s).",
		"arch.container_no_layers": "Este contêiner não declara camadas.",
		"arch.no_layer_has_spec":   "Nenhuma camada deste contêiner tem spec ainda:",
		"arch.files":               "arquivo(s)",
		"arch.no_unit_runs_here":   "Nenhuma camada declara unidade que rode aqui. Ou o contêiner não executa código deste\nrepositório, ou falta ligá-lo a uma camada em ``containers.layers``.",
		"arch.orphans":             "Camadas fora de todo contêiner",
		"arch.orphans.desc":        "Estas camadas existem no projeto e nenhum contêiner as declara — não aparecem em diagrama\nde nível 3 nenhum. Ou falta declará-las, ou elas não rodam em lugar nenhum:",
		"arch.l4":                  "Nível 4 — Código",
		"arch.l4.note":             "Só onde houver algo não óbvio. O mapa que o Anchors mantém já é a versão de\n     máquina, e repeti-lo aqui seria duplicação sem leitor.",
		"arch.l4.empty":            "Nada a destacar por enquanto.",
		"behavior.file":            "comportamento",
		"behavior.why":             "O QUE ACONTECE — índice dos cenários, que levam ao conteúdo na página da camada",
		"behavior.title":           "Comportamento",
		"behavior.intro":           "Todos os cenários do sistema. Cada um leva à unidade que o define.\n\nUm cenário descreve o que o sistema faz numa situação — vem da feature, e é o mesmo que o\nteste prova.",
		"behavior.none":            "Nenhum cenário ainda: as unidades desta camada não têm feature.",
		"rules.file":               "regras",
		"rules.why":                "AS REGRAS de todo o sistema — índice que atravessa as camadas",
		"rules.title":              "Regras",
		"rules.intro":              "Todas as regras do sistema, de todas as unidades. Cada uma leva ao texto na página da sua\ncamada.\n\nPara ver uma camada inteira de uma vez — com a visão geral de cada unidade e os cenários —\nabra a página dela em ``«layers.dir»/``.",
		"layer.why":                "o que existe na camada `«layer»`, com o conteúdo de cada unidade",
		"layer.title":              "Camada:",
	},
	"es": {
		"arch.file":                "arquitectura",
		"arch.why":                 "CÓMO está montado el sistema — el C4: contexto, contenedores, y un nivel 3 por contenedor",
		"arch.title":               "Arquitectura",
		"arch.l1":                  "Nivel 1 — Contexto",
		"arch.l1.desc":             "El sistema como una sola caja: quién lo usa, y con qué sistemas externos habla.",
		"arch.person":              "Persona",
		"arch.person.desc":         "quien usa el sistema",
		"arch.system":              "EL SISTEMA",
		"arch.system.desc":         "lo que este repositorio construye",
		"arch.external":            "Sistema externo",
		"arch.external.desc":       "de donde vienen los datos",
		"arch.uses":                "usa",
		"arch.queries":             "consulta",
		"arch.l1.out":              "FUERA de este nivel: cómo está montado el sistema por dentro. Eso es el nivel 2.",
		"arch.l2":                  "Nivel 2 — Contenedores",
		"arch.l2.desc":             "El zoom de la caja \"EL SISTEMA\": lo que corre o almacena por separado, y **con qué protocolo\nconversa cada par**. El protocolo es lo que dice qué pasa cuando la conversación falla.",
		"arch.no_containers":       "**Ningún contenedor declarado.** El nivel 2 y los de nivel 3 salen vacíos hasta que la\n> configuración declare lo que corre por separado. Ver el bloque de contenedores en ``anchors.yaml``.",
		"arch.l2.out":              "FUERA de este nivel: las piezas DENTRO de cada contenedor. Eso es el nivel 3 — y hay un\n     diagrama por contenedor, porque cada nivel amplía UNA caja del anterior.",
		"arch.l3":                  "Nivel 3 — Componentes",
		"arch.l3.desc":             "Un diagrama **por contenedor**: cada uno amplía una caja del nivel 2. Los externos no\naparecen aquí — no tenemos componentes dentro de ellos, y dibujarlos afirmaría un\nconocimiento que no tenemos.",
		"arch.files_no_spec":       "archivo(s), ninguna spec",
		"arch.layer_no_spec.a":     "Ninguna spec en esta capa —",
		"arch.layer_no_spec.b":     "archivo(s) regido(s).",
		"arch.container_no_layers": "Este contenedor no declara capas.",
		"arch.no_layer_has_spec":   "Ninguna capa de este contenedor tiene spec todavía:",
		"arch.files":               "archivo(s)",
		"arch.no_unit_runs_here":   "Ninguna capa declara una unidad que corra aquí. O el contenedor no ejecuta código de este\nrepositorio, o falta ligarlo a una capa en ``containers.layers``.",
		"arch.orphans":             "Capas fuera de todo contenedor",
		"arch.orphans.desc":        "Estas capas existen en el proyecto y ningún contenedor las declara — no aparecen en ningún\ndiagrama de nivel 3. O falta declararlas, o no corren en ningún lugar:",
		"arch.l4":                  "Nivel 4 — Código",
		"arch.l4.note":             "Solo donde haya algo no obvio. El mapa que Anchors mantiene ya es la versión de\n     máquina, y repetirlo aquí sería duplicación sin lector.",
		"arch.l4.empty":            "Nada que destacar por ahora.",
		"behavior.file":            "comportamiento",
		"behavior.why":             "LO QUE PASA — índice de los escenarios, que llevan al contenido en la página de la capa",
		"behavior.title":           "Comportamiento",
		"behavior.intro":           "Todos los escenarios del sistema. Cada uno lleva a la unidad que lo define.\n\nUn escenario describe lo que el sistema hace en una situación — viene de la feature, y es lo\nmismo que el test prueba.",
		"behavior.none":            "Ningún escenario todavía: las unidades de esta capa no tienen feature.",
		"rules.file":               "reglas",
		"rules.why":                "LAS REGLAS de todo el sistema — índice que atraviesa las capas",
		"rules.title":              "Reglas",
		"rules.intro":              "Todas las reglas del sistema, de todas las unidades. Cada una lleva al texto en la página de\nsu capa.\n\nPara ver una capa entera de una vez — con la visión general de cada unidad y los escenarios —\nabre su página en ``«layers.dir»/``.",
		"layer.why":                "lo que existe en la capa `«layer»`, con el contenido de cada unidad",
		"layer.title":              "Capa:",
	},
}

// InitScaffolds escreve os templates iniciais, inclusive uma página por camada existente.
//
// NÃO sobrescreve sem `--force`: um template já editado carrega a moldura que o time
// escreveu, e refazê-la a cada `init` seria apagar exatamente a parte que não é gerada.
func (c *Compiler) InitScaffolds(force bool) (escritos, pulados []string, err error) {
	todos := Scaffolds(c.lang())
	for _, l := range c.fnLayers() {
		todos = append(todos, ScaffoldLayer(l, c.lang()))
	}
	// A project with API units gets its OpenAPI, compiled from them (see openapi.go).
	if c.hasAPISpecs() {
		todos = append(todos, ScaffoldOpenAPI(filepath.Base(c.Root)))
	}
	// A project whose specs declare environment variables gets their page.
	if c.hasEnvSpecs() {
		todos = append(todos, ScaffoldEnvironment(c.lang()))
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

// ScaffoldOpenAPI is the template of the project's OpenAPI document. The title, the version
// and the servers are the team's to edit; the operations come from the specs.
func ScaffoldOpenAPI(title string) Scaffold {
	return Scaffold{
		Nome:   "openapi.yaml.tmpl",
		Corpo:  fmt.Sprintf("{{ openapi %q \"0.1.0\" }}\n", title),
		Porque: "The API's contract, compiled from the specs that have an Endpoint section. Edit the title, the version and the servers (more arguments) here; the operations, their contracts and errors come from the specs.",
	}
}

// hasAPISpecs says whether any spec the compiler loaded has an Endpoint section.
func (c *Compiler) hasAPISpecs() bool {
	for _, sp := range c.specs {
		if len(sectionTable(sp, "Endpoint")) > 0 {
			return true
		}
	}
	return false
}
