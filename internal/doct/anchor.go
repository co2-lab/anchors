package doct

import (
	"sort"
	"strings"
	"unicode"
)

// --- o LINK que leva ao conteúdo ---
//
// Um índice só serve se o item levar ao texto. E o alvo é uma âncora de heading, que
// nenhum gerador inventa: GitHub, MkDocs e Starlight seguem a mesma convenção — o texto do
// heading em minúsculas, pontuação removida, espaços viram hífen.
//
// Implementar isso aqui e não deixar para o autor do template é o ponto: uma âncora escrita
// à mão erra em silêncio. O link existe, é clicável, e não vai a lugar nenhum — o navegador
// simplesmente fica onde está. É o pior tipo de erro num índice, porque parece funcionar.

// GitHubAnchor devolve a âncora de um heading, na convenção do GitHub.
//
// As regras, na ordem em que se aplicam:
//
//	minúsculas          "GLCGL-B01" → "glcgl-b01"
//	pontuação removida  "a régua, e o resto" → "a régua e o resto"
//	espaço vira hífen   "a régua e o resto" → "a-régua-e-o-resto"
//
// Os ACENTOS FICAM. É contraintuitivo e foi verificado: o GitHub preserva `é` na âncora
// (`#a-régua`), e removê-los produziria um link que não resolve. Bibliotecas que fazem
// "slug" costumam transliterar por padrão, e é por isso que esta função existe em vez de
// uma delas.
func GitHubAnchor(heading string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(heading)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteRune('-')
			// Todo o resto — pontuação, travessão, crase, parêntese — desaparece.
		}
	}
	return b.String()
}

// fnAnchor é a função de template: `{{anchor .Titulo}}`.
func fnAnchor(heading string) string { return GitHubAnchor(heading) }

// fnLayerPage devolve o caminho da página de uma camada, relativo a `docs/`.
//
// Existe para que o template não monte o caminho com `printf`: um índice em `docs/` e as
// páginas em `docs/camadas/` já são dois níveis, e a primeira vez que alguém mover a pasta
// todos os links quebram de uma vez. Aqui há um lugar só para consertar.
func fnLayerPage(layer string) string { return "camadas/" + layer + ".md" }

// fnRuleLink devolve o link COMPLETO para uma regra, âncora incluída.
//
// Existe porque o índice não pode adivinhar o formato da página de destino, e essa foi a
// lição mais cara desta rodada: a página de camada mostra as regras como HEADING quando a
// camada é pequena, e como item de lista quando é grande (uma página com trinta unidades
// completas não se lê). O índice montava a âncora da regra sempre — e nas camadas grandes
// ela não existia.
//
// Medido: 483 de 812 links quebrados, todos clicáveis, todos parando no mesmo lugar. Um
// índice assim é pior que não ter índice.
//
// A regra aqui é uma só: QUEM DECIDE O FORMATO É QUEM MONTA O LINK. Onde a regra tem
// heading próprio, o link leva a ela; onde não tem, leva à UNIDADE que a contém — que é
// onde o leitor a encontra de qualquer jeito, e é uma promessa que a página cumpre.
func (c *Compiler) fnRuleLink(layer string, s Spec, r Rule) string {
	destino := fnLayerPage(layer)
	if c.layerIsBig(layer) {
		return destino + "#" + GitHubAnchor(s.Code+" — "+s.Titulo)
	}
	return destino + "#" + GitHubAnchor(r.Code+" — "+r.Titulo)
}

// fnScenarioLink faz o mesmo para um cenário.
//
// Na camada grande o cenário não aparece: a página resume, e o corpo Gherkin fica fora.
// O link leva à unidade, e o índice diz onde o cenário está definido — o que é verdade, e
// é mais do que um link para lugar nenhum diria.
func (c *Compiler) fnScenarioLink(layer string, sc Scenario) string {
	destino := fnLayerPage(layer)
	if !c.layerIsBig(layer) {
		return destino + "#" + GitHubAnchor(sc.Code+" — "+sc.Titulo)
	}
	if s := c.specOf(sc.Spec); s != nil {
		return destino + "#" + GitHubAnchor(s.Code+" — "+s.Titulo)
	}
	return destino
}

// layerIsBig consulta o LAYOUT — a mesma fonte que os templates consultam.
//
// É o que garante que o link e a página concordem: os dois perguntam ao mesmo objeto, e
// não há um segundo lugar onde discordar. Ver `layout.go`.
func (c *Compiler) layerIsBig(layer string) bool {
	g, _ := c.fnSize("layer=" + layer)
	return c.Layout.Big(g)
}

// fnBig é a pergunta do template: `{{if big "layer=lambdas"}}`.
//
// Sem argumento de limite, de propósito. Um template que escolhesse o seu recriaria o
// defeito: a página resumiria por um critério e os links seriam montados por outro.
func (c *Compiler) fnBig(filtro ...string) (bool, error) {
	g, err := c.fnSize(filtro...)
	if err != nil {
		return false, err
	}
	return c.Layout.Big(g), nil
}

func (c *Compiler) specOf(code string) *Spec {
	for i := range c.specs {
		if c.specs[i].Code == code {
			return &c.specs[i]
		}
	}
	return nil
}

// --- QUANTO cabe numa página ---
//
// Uma camada com três unidades numa página é conveniente; com trinta, é um documento que
// ninguém rola até o fim. O corte certo depende do projeto, e o Anchors não tem como
// adivinhá-lo — mas tem o DADO, e é o que ele dá ao template.
//
// A decisão fica no template, e não no compilador, por uma razão que só aparece depois:
// se o compilador dividisse sozinho, o número de páginas mudaria quando uma camada
// cruzasse o limiar, e todo link externo para `camadas/lambdas.md#CODIGO` quebraria em
// silêncio no dia em que a trigésima unidade entrasse. Com a decisão no template, ela está
// escrita, versionada, e muda quando alguém decide mudá-la.

// Size descreve o tamanho de um recorte, para o template decidir o formato.
type Size struct {
	Units  int // quantas unidades
	Rules  int // quantas regras somadas
	Lines  int // quantas linhas de texto as specs somam
	Scenes int // quantos cenários
}

// O tamanho é DADO; a decisão sobre ele é do `Layout`.
//
// A primeira versão punha um `Big(limite int)` aqui, para o template decidir. Parecia
// flexível e era o defeito: cada template escolhia o seu limite, e os links — montados
// pelo compilador — usavam outro. Ver `layout.go`.

// fnSize mede um recorte — o mesmo filtro que o `specs` aceita.
func (c *Compiler) fnSize(filtro ...string) (Size, error) {
	sel, err := c.fnSpecs(filtro...)
	if err != nil {
		// Recorte VAZIO é tamanho zero, não erro: perguntar o tamanho de uma camada sem
		// unidade é legítimo — a resposta é "não tem nada" — e propagar o erro do `specs`
		// aqui derrubaria o build de um template que só queria decidir um formato.
		return Size{}, nil
	}
	var out Size
	out.Units = len(sel)
	for _, s := range sel {
		out.Rules += len(fnRules(s))
		out.Lines += strings.Count(s.raw, "\n")
		out.Scenes += len(c.fnScenarios(s))
	}
	return out, nil
}

// --- as DEPENDÊNCIAS entre camadas, medidas ---
//
// Um diagrama de caixas sem seta não é diagrama: é uma lista num formato caro. E as setas
// não podem ser desenhadas à mão no nível 3, porque é justamente ali que o desenho
// envelhece primeiro — alguém acrescenta uma dependência no código e ninguém volta ao
// diagrama.
//
// O `boundaries:` do `anchors.yaml` não serve: ele declara o que uma camada NÃO pode
// alcançar (regex de proibição), e a ausência de proibição não é uma dependência.
//
// O que serve é o MAPA. As arestas `needs` e `depends-on` são declaradas nas próprias
// unidades e conferidas por gate — o diagrama passa a mostrar o que o projeto realmente
// tem, e muda sozinho quando o projeto muda.

// LayerDep é uma dependência entre duas camadas, com o peso.
type LayerDep struct {
	From, To string
	Count    int // quantas unidades sustentam esta seta
}

// fnLayerDeps devolve as dependências entre camadas, deduzidas do mapa.
//
// Agrega por CAMADA, e o `Count` é o que sobra da agregação: trinta unidades de `screen`
// dependendo de `shared` viram uma seta, e o número diz que ela é grossa. Sem ele, uma
// dependência única e uma sistêmica pareceriam a mesma coisa.
func (c *Compiler) fnLayerDeps() []LayerDep {
	camadaDe := map[string]string{}
	for _, s := range c.specs {
		camadaDe[s.Path] = s.Layer
		camadaDe[s.Code] = s.Layer
	}

	peso := map[[2]string]int{}
	for _, e := range c.Graph.Edges {
		if e.Type != "needs" && e.Type != "depends-on" {
			continue
		}
		de, ok1 := camadaDe[e.From]
		para, ok2 := camadaDe[e.To]
		// Aresta que sai ou chega fora das specs — um plano, um teste — não é
		// dependência de camada. Deixá-la entrar desenharia setas de nós que a página
		// não mostra.
		if !ok1 || !ok2 || de == "" || para == "" || de == para {
			continue
		}
		peso[[2]string{de, para}]++
	}

	var out []LayerDep
	for k, n := range peso {
		out = append(out, LayerDep{From: k[0], To: k[1], Count: n})
	}
	// Ordem estável: sem ela o diagrama muda de forma a cada build, e o `docs-fresh`
	// acusaria defasagem onde nada mudou.
	sort.Slice(out, func(i, j int) bool {
		if out[i].From != out[j].From {
			return out[i].From < out[j].From
		}
		return out[i].To < out[j].To
	})
	return out
}

// MermaidID devolve um identificador seguro para nó de diagrama Mermaid.
//
// O Mermaid quebra em hífen e em acento no ID do nó, e o erro é de PARSE: o bloco inteiro
// deixa de renderizar e aparece como texto cru. Uma camada chamada `feature-hook` — que
// existe no projeto de referência — bastaria para derrubar o diagrama todo.
func MermaidID(nome string) string {
	var b strings.Builder
	for _, r := range nome {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return "n_" + b.String()
}

func fnMermaidID(nome string) string { return MermaidID(nome) }
