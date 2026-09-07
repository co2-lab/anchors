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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

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
const GeneratedMarker = "<!-- anchors:generated de %s — NÃO EDITE: rode `anchors docs build` -->"

var markerRE = regexp.MustCompile(`^<!-- anchors:generated `)

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
	specs []Spec
}

func New(root string, g *mapx.Graph) (*Compiler, error) {
	c := &Compiler{Root: root, Graph: g}
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
			return fmt.Errorf("spec %s está no mapa e não no disco: %w "+
				"(rode `anchors map build`)", n.ID, err)
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
		return c.specs, nil
	}
	campo, valor, ok := strings.Cut(filtro[0], "=")
	if !ok {
		return nil, fmt.Errorf("filtro %q sem `=`: use `campo=valor` (ex.: `layer=screen`)", filtro[0])
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
			return nil, fmt.Errorf("nenhuma spec na camada %q — existem: %s",
				valor, strings.Join(c.fnLayers(), ", "))
		}
	case "code":
		for _, s := range c.specs {
			if s.Code == valor {
				out = append(out, s)
			}
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("nenhuma spec com o código %q", valor)
		}
	default:
		return nil, fmt.Errorf("campo %q desconhecido: use `layer` ou `code`", campo)
	}
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
	return nil, fmt.Errorf("nenhuma spec com o código %q", code)
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
			Code:   s.raw[m[2]:m[3]],
			Titulo: s.raw[m[4]:m[5]],
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
		return nil, fmt.Errorf("sem `%s/`: crie os templates antes (`anchors docs init`)", Dir)
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
	fmt.Fprintf(&buf, GeneratedMarker+"\n\n", rel)
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
		esperado, err := c.compile(filepath.Join(c.Root, Dir, filepath.FromSlash(rel)), saida)
		if err != nil {
			return nil, err
		}
		atual, err := os.ReadFile(filepath.Join(c.Root, OutDir, filepath.FromSlash(saida)))
		if err != nil {
			out = append(out, saida) // não existe: está defasado por ausência
			continue
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
