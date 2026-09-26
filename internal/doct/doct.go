// Package doct compila DOCUMENTAÇÃO a partir de templates que referenciam as specs.
//
// O PROBLEMA. Uma documentação útil mostra o conteúdo, não uma lista de links — "abrir
// link por link para ver o conteúdo de uma página deixa de ser documentação e passa a ser
// indexação". Mas o conteúdo já existe nas specs, e copiá-lo para dentro de um `.md`
// escrito à mão cria duplicação que desatualiza em silêncio.
//
// E há uma matriz: as specs têm SEÇÕES (Regras, Invariantes, Fora de escopo) e pertencem a
// CAMADAS (screen, lambdas, hook…). Quem documenta quer as duas visões — "todas as regras
// do sistema" e "o que existe em screens" — sobre os mesmos dados.
//
// POR QUE MARKDOWN NÃO RESOLVE SOZINHO. Não há inclusão nativa no CommonMark nem no GFM:
// `{% include %}` é Jekyll (só no Pages), `--8<--` é MkDocs, `<iframe>` é sanitizado pelo
// GitHub. Sem build, a única saída seria duplicar à mão.
//
// A SAÍDA É O BUILD, com três camadas e o conteúdo morando em UMA:
//
//	*.spec.md        a fonte — o conteúdo mora aqui, e só aqui
//	doct/*.md.tmpl   o template — o agente escreve a moldura e REFERENCIA trechos
//	docs/*.md        o compilado — conteúdo real, gerado, ninguém edita
//
// A duplicação existe só no compilado, e ali ela não dói: ele é gerado, e o `docs-fresh`
// acusa quando está defasado. A duplicação que dói é a que desatualiza sem ninguém ver.
//
// O `text/template` é da stdlib de propósito: uma dependência nova só para documentação
// seria custo permanente, e o Anchors já usa a mesma linguagem no `derived.files`
// (`{{dir}}`, `{{name}}`).
package doct

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// Dir é onde os templates vivem.
//
// `doct` e não `doc[t]`: colchete no nome de pasta é classe de caracteres em glob de shell
// e em `.gitignore` — `doc[t]/` casaria `doct/`, e todo padrão que mencionasse a pasta
// precisaria escapar. Alguém esqueceria. E numa URL exige encoding (`doc%5Bt%5D`), que
// quebra link no GitHub.
const Dir = "doct"

// OutDir é onde o compilado é escrito.
const OutDir = "docs"

// SufixoTemplate: `screens.md.tmpl` → `docs/screens.md`.
const SufixoTemplate = ".tmpl"

// GeneratedMarker abre todo arquivo compilado.
//
// Serve a duas coisas, e a segunda é a que protege: quem abre o arquivo sabe que editá-lo
// é perder o trabalho; e o compilador RECUSA sobrescrever um `.md` que não o tenha — um
// documento escrito à mão não é destruído porque alguém criou um template com o mesmo
// nome.
const GeneratedMarker = "<!-- anchors:generated from %s — DO NOT EDIT: run `anchors docs build` -->"

// GeneratedMarkerHashed e' o marcador COM a impressao digital das entradas que
// produziram a pagina: o template e as specs que ele consome.
//
// Existe para que conferir se a doc esta em dia deixe de exigir recompila-la. Sem o
// carimbo, a unica forma de responder "esta defasada?" e' montar o documento inteiro e
// comparar — e o `docs-fresh` faz essa pergunta a cada varredura, mesmo quando nada
// mudou. Com o carimbo, a resposta e' comparar dois hashes.
//
// O hash cobre as DUAS entradas de proposito. Cobrir so' as specs deixaria passar um
// template editado (a doc mudaria de forma sem nenhuma spec mudar); cobrir so' o
// template deixaria passar o caso comum, que e' a spec revisada.
const GeneratedMarkerHashed = "<!-- anchors:generated from %s — inputs:%s — DO NOT EDIT: run `anchors docs build` -->"

var markerRE = regexp.MustCompile(`^<!-- anchors:generated `)

// markerHashRE extrai a impressao digital de um marcador que a tenha. A ausencia nao e'
// erro: paginas compiladas por uma versao anterior do Anchors nao a carregam, e para
// elas o caminho continua sendo recompilar e comparar.
var markerHashRE = regexp.MustCompile(`^<!-- anchors:generated from [^—]+— inputs:([0-9a-f]{16}) —`)

// HandwrittenMarker abre a página que NÃO é gerada — a doc de produto é o caso típico.
//
// É COMENTÁRIO, e é o ponto: a nota se dirige a quem EDITA o arquivo, não a quem lê o
// documento. Escrita como texto visível — um blockquote sob o título, que foi a primeira
// forma que tomou — ela rouba o primeiro lugar da página para falar de mecânica de
// ferramenta, e o leitor que veio saber o que o produto faz lê antes uma explicação sobre
// templates.
//
// Markdown não renderiza `<!-- -->` em lugar nenhum — GitHub, MkDocs, Starlight — e é o
// mesmo mecanismo do `GeneratedMarker`, pela mesma razão: a instrução fica visível para
// quem abre o arquivo e invisível para quem abre a página.
const HandwrittenMarker = `<!-- anchors:handwritten — this page is NOT generated.

     It has no template in ` + "`" + Dir + "/`" + `, and ` + "`anchors docs build`" + ` does not touch it: the compiler
     only overwrites a file carrying the ` + "`anchors:generated`" + ` marker. Edit right here. -->`

var handwrittenRE = regexp.MustCompile(`^<!-- anchors:handwritten`)

// IsGenerated diz se o conteúdo de um `.md` foi produzido pelo compilador.
func IsGenerated(conteudo []byte) bool { return markerRE.Match(conteudo) }

// IsHandwritten diz se a página se DECLARA escrita à mão.
//
// Não é o mesmo que "não é gerada": um `.md` sem marcador nenhum também não é gerado, mas
// não disse nada sobre si. A distinção importa para o `docs list`, que separa o que alguém
// decidiu manter à mão do que simplesmente ainda não foi olhado.
func IsHandwritten(conteudo []byte) bool { return handwrittenRE.Match(conteudo) }

// Spec é o que um template vê de uma unidade.
type Spec struct {
	Code   string
	Layer  string
	Path   string
	Titulo string
	// conteúdo bruto, para as funções de seção
	raw string
}

// Rule é um requisito individual: o código, o título e o corpo.
type Rule struct {
	Code   string
	Titulo string
	Corpo  string
}

// Compiler compila os templates de um projeto.
type Compiler struct {
	Root  string
	Graph *mapx.Graph
	// Config é a Estrutura do projeto — o compilador precisa dela para os CONTÊINERES,
	// que são declarados lá e não deriváveis do mapa.
	Config *config.Config
	// Layout é a decisão de formato, resolvida UMA VEZ e distribuída a todos os
	// templates. Ver `layout.go` — três lugares decidindo por conta foi o que produziu
	// 483 links quebrados.
	Layout Layout
	specs  []Spec
	// consumed registra quais specs os templates PEDIRAM, para a cobertura. Vive no
	// compilador e nao num retorno porque quem as pede sao as funcoes de template,
	// chamadas de dentro do `Execute` — nao ha por onde devolver.
	consumed map[string]bool
	// Os INDICES do grafo, montados uma vez (ver `buildIndexes`). O grafo nao muda
	// durante a compilacao, e reconstruir a resposta a cada consulta era o que fazia uma
	// compilacao completa custar 7.7s.
	kindByID       map[string]mapx.Kind
	featuresBySpec map[string][]string
	// sizeCache memoriza o tamanho de cada recorte — ver `fnSize`, onde esta a medida
	// que justifica o cache.
	sizeCache map[string]Size
}

func New(root string, g *mapx.Graph) (*Compiler, error) {
	c := &Compiler{Root: root, Graph: g, Layout: DefaultLayout()}
	// A Estrutura é OPCIONAL aqui: o erro de carregá-la já é reportado por quem chama o
	// compilador, e um template que não usa contêineres não deve falhar porque o
	// `anchors.yaml` não os declara.
	c.Config, _ = config.Load(filepath.Join(root, config.DefaultFile))
	if err := c.loadSpecs(); err != nil {
		return nil, err
	}
	return c, nil
}

// titleRE casa o H1 da spec: `# GoLiveChecklist — a régua que registra o que não foi feito`.
var titleRE = regexp.MustCompile(`(?m)^#\s+(.+?)\s*$`)

// loadSpecs lê do MAPA quais specs existem, e o disco para o conteúdo.
//
// Do mapa e não de um walk próprio: o mapa é quem sabe a camada de cada arquivo (`layers:`
// do anchors.yaml) e a identidade dele. Um walk aqui teria de reimplementar as duas
// coisas, e divergiria na primeira mudança da Estrutura.
func (c *Compiler) loadSpecs() error {
	for _, n := range c.Graph.Nodes {
		if n.Kind != mapx.KindSpec {
			continue
		}
		b, err := os.ReadFile(filepath.Join(c.Root, n.ID))
		if err != nil {
			// PARA. Uma spec que o mapa conhece e o disco não tem sairia da documentação
			// sem uma palavra — e o compilado, verde, ficaria sem ela. É o modo de falha
			// que o mecanismo inteiro existe para evitar, entrando pela porta dos fundos.
			return fmt.Errorf("spec %s is in the map and not on disk: %w "+
				"(run `anchors map build`)", n.ID, err)
		}
		raw := string(b)
		titulo := ""
		if m := titleRE.FindStringSubmatch(raw); m != nil {
			titulo = m[1]
		}
		c.specs = append(c.specs, Spec{
			Code: n.Code, Layer: headerLayer(raw, n.Layer), Path: n.ID,
			Titulo: titulo, raw: raw,
		})
	}
	sort.Slice(c.specs, func(i, j int) bool { return c.specs[i].Path < c.specs[j].Path })
	return nil
}

// headerLayerRE lê a camada do HEADER da spec.
var headerLayerRE = regexp.MustCompile(`(?m)^\s*layer:\s*([a-zA-Z0-9_-]+)`)

// headerLayer devolve a camada que a SPEC declara, e não a do mapa.
//
// A distinção importa e foi medida: no projeto de referência as 84 specs têm `layer: spec`
// no mapa (elas casam o pattern `**/*.spec.md`), e a camada REAL da unidade — screen,
// lambdas, hook — está no header de cada uma. Agrupar pela do mapa produziria um único
// grupo com 84 itens, que é o oposto do que a visão por camada existe para dar.
func headerLayer(raw, doMapa string) string {
	if m := headerLayerRE.FindStringSubmatch(raw); m != nil {
		return m[1]
	}
	return doMapa
}

// Funcs são as funções que um template pode chamar.
func (c *Compiler) Funcs() template.FuncMap {
	return template.FuncMap{
		"specs":      c.fnSpecs,
		"layers":     c.fnLayers,
		"section":    fnSection,
		"rules":      fnRules,
		"specByCode": c.fnSpecByCode,
		// Os CENÁRIOS: a feature é o único artefato escrito para quem não programa, e é
		// ela que diz o que ACONTECE — a spec diz a regra em abstrato.
		"scenarios":    c.fnScenarios,
		"allScenarios": c.fnAllScenarios,
		// O LINK do índice. `anchor` produz a âncora do heading na convenção que GitHub,
		// MkDocs e Starlight compartilham; `layerPage` monta o caminho da página de
		// camada num lugar só, para que mover a pasta não quebre todos os links.
		"anchor":    fnAnchor,
		"layerPage": fnLayerPage,
		// O LINK de uma regra/cenário conhece o FORMATO da página de destino: onde a
		// camada é grande a página resume, e a âncora da regra não existe lá.
		"ruleLink":     c.fnRuleLink,
		"scenarioLink": c.fnScenarioLink,
		// O TAMANHO do recorte, para o template escolher o formato: tudo numa página
		// quando são poucas unidades, uma página por unidade quando são muitas.
		"size": c.fnSize,
		// O LAYOUT — a mesma resposta para todas as páginas.
		"layout": func() Layout { return c.Layout },
		"big":    c.fnBig,
		// As DEPENDÊNCIAS entre camadas, deduzidas das arestas `needs` do mapa — o
		// diagrama mostra o que o projeto tem, não o que alguém desenhou uma vez.
		"layerDeps": c.fnLayerDeps,
		// Os CONTÊINERES — o nível 2 do C4, e o que dá um diagrama de nível 3 a cada um.
		"containers":         c.fnContainers,
		"internalContainers": c.fnInternalContainers,
		"orphanLayers":       c.fnOrphanLayers,
		"mermaidID":          fnMermaidID,
	}
}

// fnSpecs devolve as specs que casam um filtro simples (`layer=screen`, ou vazio = todas).
//
// O filtro é uma STRING e não opções nomeadas porque o `text/template` não tem argumento
// nomeado: `specs "layer=screen"` lê melhor que `specs "screen"` — que não diria por qual
// campo está filtrando — e melhor que uma função por campo, que multiplicaria a API.
//
// ERRA ALTO, e é o ponto: um filtro que não casa nada devolveria uma seção VAZIA num
// documento que compila verde. A falha seria invisível — o `.md` existe, tem título, tem
// as outras seções, e a que interessa simplesmente não está lá. Ninguém procura o que não
// sabe que falta.
//
// Três formas de errar, todas silenciosas se não fossem erro:
//
//	specs "layar=screen"   campo com typo      → switch sem case
//	specs "layer=Screen"   camada inexistente  → nenhuma spec casa
//	specs "screen"         sem o `=`           → Cut falha
//
// A terceira é a mais provável, porque é como se escreveria se a API fosse por campo.
func (c *Compiler) fnSpecs(filtro ...string) ([]Spec, error) {
	if len(filtro) == 0 || strings.TrimSpace(filtro[0]) == "" {
		c.markConsumed(c.specs)
		return c.specs, nil
	}
	campo, valor, ok := strings.Cut(filtro[0], "=")
	if !ok {
		return nil, fmt.Errorf("filter %q without `=`: use `field=value` (e.g. `layer=screen`)", filtro[0])
	}
	campo, valor = strings.TrimSpace(campo), strings.TrimSpace(valor)
	var out []Spec
	switch campo {
	case "layer":
		for _, s := range c.specs {
			if s.Layer == valor {
				out = append(out, s)
			}
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("no spec in layer %q — existing: %s",
				valor, strings.Join(c.fnLayers(), ", "))
		}
	case "code":
		for _, s := range c.specs {
			if s.Code == valor {
				out = append(out, s)
			}
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("no spec with code %q", valor)
		}
	default:
		return nil, fmt.Errorf("unknown field %q: use `layer` or `code`", campo)
	}
	c.markConsumed(out)
	return out, nil
}

// markConsumed anota quais specs um template pediu — o insumo da COBERTURA.
//
// A pergunta que isto responde nao e' a do `docs-fresh` ("a pagina esta em dia?"), e sim
// "a documentacao ALCANCA todas as specs?". Sao defeitos diferentes, e o segundo e' o
// silencioso: os templates filtram (`specs "layer=gate"`), entao uma spec de camada que
// nenhum template seleciona nao aparece em pagina nenhuma — e nada acusa, porque todas
// as paginas que existem estao corretas.
//
// Hoje nao ha orfa neste repositorio, e e' coincidencia do estado atual: as 51 specs sao
// de `gate`, e ha um template que pede exatamente `layer=gate`. A primeira spec fora
// dessa camada some da documentacao sem uma palavra.
func (c *Compiler) markConsumed(specs []Spec) {
	if c.consumed == nil {
		c.consumed = map[string]bool{}
	}
	for _, sp := range specs {
		c.consumed[sp.Path] = true
	}
}

// Uncovered devolve as specs que NENHUM template consumiu, depois de compilar todos.
//
// Compila de verdade (e nao le os `.md` ja gerados) porque a pergunta e' sobre o que os
// templates ALCANCAM, e isso so' se sabe executando-os: o filtro vive dentro do
// template, e le-lo por regex reimplementaria a linguagem de template por fora.
func (c *Compiler) Uncovered() ([]string, error) {
	tmpls, err := c.templates()
	if err != nil {
		return nil, nil // sem `doct/`, nao ha cobertura a cobrar
	}
	c.consumed = map[string]bool{}
	for _, rel := range tmpls {
		saida := strings.TrimSuffix(rel, SufixoTemplate)
		if _, err := c.compile(filepath.Join(c.Root, Dir, filepath.FromSlash(rel)), saida); err != nil {
			return nil, err
		}
	}
	var out []string
	for _, sp := range c.specs {
		if !c.consumed[sp.Path] {
			out = append(out, sp.Path)
		}
	}
	sort.Strings(out)
	return out, nil
}

// fnLayers devolve as camadas que TÊM spec, ordenadas.
func (c *Compiler) fnLayers() []string {
	visto := map[string]bool{}
	var out []string
	for _, s := range c.specs {
		if s.Layer != "" && !visto[s.Layer] {
			visto[s.Layer] = true
			out = append(out, s.Layer)
		}
	}
	sort.Strings(out)
	return out
}

// fnSpecByCode busca UMA spec pelo código. Erra alto pelo mesmo motivo que o `fnSpecs`:
// devolver `nil` faria o template renderizar o nada, e o documento sairia com um buraco
// exatamente onde alguém quis destacar uma unidade.
func (c *Compiler) fnSpecByCode(code string) (*Spec, error) {
	for i := range c.specs {
		if c.specs[i].Code == code {
			return &c.specs[i], nil
		}
	}
	return nil, fmt.Errorf("no spec with code %q", code)
}

// fnSection devolve o texto de uma seção `## <nome>` da spec, sem o cabeçalho dela.
//
// Sem o cabeçalho porque o template decide o NÍVEL do título no documento: a mesma seção
// pode ser `##` num documento por camada e `###` num agrupado por seção. Devolver o
// cabeçalho junto obrigaria o template a removê-lo.
func fnSection(s Spec, nome string) string {
	return sliceSection(s.raw, "## ", nome)
}

// sliceSection devolve o corpo de um heading até o próximo heading do MESMO nível.
//
// Do mesmo nível e não de qualquer nível: uma seção `## Regras` contém `### CODE-B01`, e
// cortar no primeiro `###` devolveria a seção vazia.
func sliceSection(raw, prefixo, nome string) string {
	linhas := strings.Split(raw, "\n")
	inicio := -1
	for i, l := range linhas {
		if strings.TrimSpace(l) == prefixo+nome {
			inicio = i + 1
			break
		}
	}
	if inicio < 0 {
		return ""
	}
	fim := len(linhas)
	for i := inicio; i < len(linhas); i++ {
		if strings.HasPrefix(linhas[i], prefixo) {
			fim = i
			break
		}
	}
	return strings.TrimSpace(strings.Join(linhas[inicio:fim], "\n"))
}

// ruleRE casa a definição de uma regra: `### GLCGL-B01 — o título`.
var ruleRE = regexp.MustCompile(`(?m)^###\s+([A-Z0-9]{4,5}-[A-Z]\d{2})\s*[—–-]\s*(.+?)\s*$`)

// fnRules devolve as regras individuais de uma spec: código, título e corpo.
//
// Existe além do `section` porque a visão por seção quer as regras SOLTAS — "todas as
// regras do sistema" é uma lista de itens, não a concatenação de seções `## Regras`.
func fnRules(s Spec) []Rule {
	ms := ruleRE.FindAllStringSubmatchIndex(s.raw, -1)
	var out []Rule
	for i, m := range ms {
		fim := len(s.raw)
		if i+1 < len(ms) {
			fim = ms[i+1][0]
		}
		// Um heading de nível MAIOR (`## Invariantes`) também encerra a regra.
		corpo := s.raw[m[1]:fim]
		if j := strings.Index(corpo, "\n## "); j >= 0 {
			corpo = corpo[:j]
		}
		out = append(out, Rule{
			Code: s.raw[m[2]:m[3]],
			// O TÍTULO É RÓTULO DE LINK, e comentário HTML não pertence a ele.
			//
			// O heading de uma regra pode carregar a dispensa na própria linha —
			// `### ABCDE-B01 — o que ela diz <!-- @no-mark: razão -->` —, e o gate a
			// aceita ali de propósito (ver `FRMTT-I02`: à direita do símbolo é onde o
			// formatador não a move). Mas o compilado usa este texto como o rótulo entre
			// colchetes, e o resultado era ilegível:
			//
			//	- [DSHBR-B01 — a escolha de fixar <!-- @no-mark: ... -->](camadas/...)
			//
			// MEDIDO no projeto de referência (#647): 57 ocorrências no `docs/regras.md`,
			// espalhadas por todas as camadas — não é defeito de uma spec, é do gerador.
			//
			// A dispensa continua onde estava: ela é lida do RAW pelos gates, e este corte
			// é só do rótulo. O `Corpo` também não a leva, porque começa depois do heading.
			Titulo: strings.TrimSpace(comentarioHTMLRE.ReplaceAllString(s.raw[m[4]:m[5]], "")),
			Corpo:  strings.TrimSpace(corpo),
		})
	}
	return out
}

// Result é o que uma compilação produziu.
type Result struct {
	Written []string
	Skipped []string // `.md` sem marcador: escrito à mão, não sobrescrito
}

// Build compila todos os templates de `doct/` para `docs/`.
func (c *Compiler) Build(dryRun bool) (Result, error) {
	var res Result
	tmpls, err := c.templates()
	if err != nil {
		return res, err
	}
	for _, rel := range tmpls {
		saida := strings.TrimSuffix(rel, SufixoTemplate)
		destino := filepath.Join(c.Root, OutDir, filepath.FromSlash(saida))

		// PROTEÇÃO: só sobrescreve o que ELE gerou.
		//
		// Um `docs/produto.md` escrito à mão não pode ser destruído porque alguém criou um
		// `doct/produto.md.tmpl` — o compilador para e reporta, em vez de apagar trabalho.
		if b, err := os.ReadFile(destino); err == nil && !markerRE.Match(b) {
			res.Skipped = append(res.Skipped, saida)
			continue
		}

		out, err := c.compile(filepath.Join(c.Root, Dir, filepath.FromSlash(rel)), saida)
		if err != nil {
			return res, err
		}
		if !dryRun {
			if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
				return res, err
			}
			if err := os.WriteFile(destino, out, 0o644); err != nil {
				return res, err
			}
		}
		res.Written = append(res.Written, saida)
	}
	return res, nil
}

// templates lista os templates de `doct/`, INCLUSIVE os de subpasta.
//
// Desce na árvore, e não é detalhe: a estrutura proposta põe uma página por camada em
// `doct/camadas/`, e um `os.ReadDir` raso as pularia todas em silêncio — metade da matriz
// sumiria da documentação com o build verde. É o mesmo modo de falha que este mecanismo
// existe para fechar, entrando pela porta do próprio compilador.
func (c *Compiler) templates() ([]string, error) {
	dir := filepath.Join(c.Root, Dir)
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("no `%s/`: create the templates first (`anchors docs init`)", Dir)
	}
	var out []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), SufixoTemplate) {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(out)
	return out, err
}

// inputsHash e' a impressao digital do template MAIS de todas as specs que o compilador
// carregou — as entradas exatas de `compile`.
//
// Ordem estavel: `loadSpecs` ja ordena as specs por caminho, e o caminho entra no hash
// junto do conteudo. Sem isso, duas execucoes com a mesma verdade produziriam hashes
// diferentes e o carimbo acusaria defasagem que nao existe.
//
// Truncado em 8 bytes (16 hex): o hash aqui detecta MUDANCA, nao defende contra
// adversario. Colisao acidental em 2^64 nao e' o risco desta engenharia, e o marcador
// fica legivel.
func (c *Compiler) inputsHash(tmpl []byte) string {
	h := sha256.New()
	h.Write(tmpl)
	h.Write([]byte{0})
	for _, sp := range c.specs {
		h.Write([]byte(sp.Path))
		h.Write([]byte{0})
		h.Write([]byte(withoutHeaderDate(sp.raw)))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)[:8])
}

// headerDateRE is the `updated_at:` line of an @anchors header.
var headerDateRE = regexp.MustCompile(`^\s*(?://|#|\*)?\s*updated_at:.*$`)

// headerDateLines is how far from the top the @anchors header may carry its date — the
// same reach `anchors touch` uses to find the header it dates.
const headerDateLines = 10

// withoutHeaderDate drops the header's `updated_at:` line from a spec before hashing.
//
// No template renders the date, so it is not an input of the document. Hashing it made
// `docs-fresh` fail on the first commit of the day of any spec change: the pre-commit
// dates the staged headers (`touch.pre_commit`) and only then runs the gates, so the
// docs built a minute before were already "out of date" — the same content, a new date.
// Measured in this repository: 42 blocking `docs-fresh` failures on a commit whose docs
// had just been built.
func withoutHeaderDate(raw string) string {
	lines := strings.SplitN(raw, "\n", headerDateLines+1)
	for i := 0; i < len(lines) && i < headerDateLines; i++ {
		if headerDateRE.MatchString(lines[i]) {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}
	return strings.Join(lines, "\n")
}

func (c *Compiler) compile(tmplPath, saida string) ([]byte, error) {
	b, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}
	t, err := template.New(filepath.Base(tmplPath)).Funcs(c.Funcs()).Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("template %s: %w", tmplPath, err)
	}
	var buf bytes.Buffer
	// O marcador aponta o template de origem com o caminho COMPLETO: com páginas em
	// subpasta, o basename sozinho não diz de qual `doct/camadas/*.tmpl` a doc veio.
	rel, errRel := filepath.Rel(c.Root, tmplPath)
	if errRel != nil {
		rel = tmplPath
	}
	rel = filepath.ToSlash(rel)
	fmt.Fprintf(&buf, GeneratedMarkerHashed+"\n\n", rel, c.inputsHash(b))
	if err := t.Execute(&buf, nil); err != nil {
		return nil, fmt.Errorf("template %s: %w", tmplPath, err)
	}
	return buf.Bytes(), nil
}

// Stale devolve os documentos cujo conteúdo no disco difere do que os templates
// produziriam agora.
//
// Compara SEM escrever, de propósito: o gate roda em hook e em CI, e um gate que conserta
// o que deveria apontar torna o resultado dependente de já ter rodado — a segunda execução
// sempre passaria, e o defeito só apareceria em quem clonasse o repositório.
func (c *Compiler) Stale() ([]string, error) {
	tmpls, err := c.templates()
	if err != nil {
		return nil, nil // sem `doct/`, não há compilado a conferir
	}
	var out []string
	for _, rel := range tmpls {
		saida := strings.TrimSuffix(rel, SufixoTemplate)
		tmplPath := filepath.Join(c.Root, Dir, filepath.FromSlash(rel))

		atual, err := os.ReadFile(filepath.Join(c.Root, OutDir, filepath.FromSlash(saida)))
		if err != nil {
			out = append(out, saida) // não existe: está defasado por ausência
			continue
		}

		// ATALHO PELO CARIMBO: se a pagina declara a impressao digital das entradas que a
		// produziram, e ela bate com as entradas de AGORA, a pagina esta em dia — sem
		// montar o documento.
		//
		// E' o que tira o custo do caso comum. Compilar para descobrir que nada mudou era
		// o trabalho que o `docs-fresh` repetia a cada varredura: medido neste
		// repositorio, ~8s por compilacao completa, paga mesmo quando nenhuma spec fora
		// tocada.
		//
		// A ausencia do carimbo NAO e' defasagem: paginas geradas por uma versao anterior
		// nao o carregam, e acusa-las mandaria o autor rodar um build que nao conserta
		// nada que ele tenha feito. Para elas o caminho antigo continua valendo — compila
		// e compara — e o proximo `docs build` grava o carimbo.
		if m := markerHashRE.FindSubmatch(atual); m != nil {
			b, errT := os.ReadFile(tmplPath)
			if errT != nil {
				return nil, errT
			}
			if string(m[1]) == c.inputsHash(b) {
				continue
			}
			out = append(out, saida)
			continue
		}

		esperado, err := c.compile(tmplPath, saida)
		if err != nil {
			return nil, err
		}
		// Um `.md` sem marcador é escrito à mão, e o `Build` não o sobrescreve. Cobrar
		// atualização de um arquivo que o compilador se recusa a escrever mandaria o autor
		// rodar um comando que não muda nada — o aviso certo vem do `Build`, que diz que o
		// template existe e o documento não foi gerado.
		if !markerRE.Match(atual) {
			continue
		}
		if !bytes.Equal(atual, esperado) {
			out = append(out, saida)
		}
	}
	return out, nil
}

// comentarioHTMLRE casa um comentário HTML inteiro, inclusive o que atravessa linhas.
//
// `(?s)` porque um `<!-- ... -->` pode quebrar de linha quando a razão da dispensa é longa
// — e um regex sem ele deixaria metade do comentário no rótulo, que é pior que o defeito
// original: o `-->` sozinho não diz nem que houve dispensa.
var comentarioHTMLRE = regexp.MustCompile(`(?s)<!--.*?-->`)
