// @anchors
//   ref: INPRN

// Package initx infere uma proposta de Estrutura (config.Config) escaneando o
// projeto de forma DETERMINÍSTICA — sem IA (ver DECISIONS: o CLI não embute modelo;
// a tarefa é estrutural). A proposta alimenta o fluxo interativo do `anchors init`,
// que confirma/ajusta com o usuário.
package initx

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// Proposal é o que a inferência descobriu — insumo para as perguntas do init.
type Proposal struct {
	Config     *config.Config
	CodeDirs   []string // diretórios-raiz com código detectados (candidatos a layer)
	GuideDir   string   // pasta de guides, se houver
	GuideFiles []string // guides encontrados (para o governs interativo)
	PlanDir    string   // pasta de planos, se houver (plans/ ou similar)
	HasSpecMD  bool     // achou *.spec.md
	HasFeature bool     // achou *.feature
	HasTest    bool     // achou *.test.*
	Colocated  bool     // spec/feature/test ficam ao lado do código (mesmo stem)?
	CodeExts   []string // extensões de código mais comuns (.ts, .go, .py…)
	// TestHandle — o atributo com que o projeto marca elementos para alcance externo
	// (`testID`, `data-testid`, `contentDescription`). Vazio quando o projeto não
	// marca nada (um backend, uma lib): aí os gates de inventário de handle não têm
	// o que confrontar e pulam, em vez de acusar o repositório inteiro.
	TestHandle string
	// Family is the project's dialect family (`go`, `ts`, `python`…), from the manifest at
	// its root or its most common code extension; empty when neither says.
	Family string
	// TestConventions are the test-naming conventions the project's own test files
	// follow, most used first.
	TestConventions []TestConvention
}

// knownHandles — os atributos de ancoragem de teste dos ecossistemas correntes,
// na ordem em que são procurados. A DETECÇÃO é por contagem no código real, não por
// presença de dependência no manifesto: um projeto pode ter React Native instalado e
// não marcar nada, e é o uso que decide se há contrato a cobrar.
var knownHandles = []string{"testID", "data-testid", "data-test-id", "contentDescription", "accessibilityIdentifier"}

var ignoredDirs = map[string]bool{
	"node_modules": true, ".git": true, "dist": true, "build": true,
	"vendor": true, ".next": true, "coverage": true, ".expo": true, ".anchors": true,
}

// Infer varre root e monta uma Proposal + um Config pré-preenchido.
func Infer(root string) (*Proposal, error) {
	p := &Proposal{}
	extCount := map[string]int{}          // extensão → nº de arquivos de código
	dirCode := map[string]int{}           // diretório-raiz (1º nível relevante) → nº de arquivos de código
	stems := map[string]map[string]bool{} // stem → conjunto de tipos (code/spec/feature/test)
	handleCount := map[string]int{}       // atributo de ancoragem → ocorrências na amostra
	conventionCount := map[TestConvention]int{}
	lidosParaHandle := 0

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if ignoredDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		// ToSlash uma vez, na fronteira: o que sai daqui alimenta isGuide/isPlan/topDir e
		// vira config gravada (GuideDir, CodeDirs), tudo com barra normal por contrato. No
		// Windows o Rel daria "apps\mobile", que não casa padrão nenhum e ainda vazaria
		// o dialeto da máquina para dentro do anchors.yaml.
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		name := d.Name()

		switch {
		case strings.HasSuffix(name, ".spec.md"):
			p.HasSpecMD = true
			mark(stems, stemOf(rel), "spec")
		case strings.HasSuffix(name, ".feature"):
			p.HasFeature = true
			mark(stems, stemOf(rel), "feature")
		case isTest(name):
			p.HasTest = true
			if c, ok := testConventionOf(name); ok {
				conventionCount[c]++
			}
			mark(stems, stemOf(rel), "test")
		case isPlan(rel):
			if p.PlanDir == "" {
				p.PlanDir = filepath.ToSlash(filepath.Dir(rel))
			}
		case isGuide(rel):
			if p.GuideDir == "" {
				p.GuideDir = filepath.ToSlash(filepath.Dir(rel))
			}
			p.GuideFiles = append(p.GuideFiles, rel)
		case isCode(name):
			ext := filepath.Ext(name)
			extCount[ext]++
			dirCode[topDir(rel)]++
			mark(stems, stemOf(rel), "code")
			// Detecção do handle de teste. Amostra os primeiros arquivos de código em
			// vez de ler o repositório inteiro: o atributo, quando existe, é onipresente
			// — e uma varredura completa faria `init` pagar leitura de I/O por um sinal
			// que a amostra já dá.
			if lidosParaHandle < maxAmostraHandle {
				lidosParaHandle++
				if b, e := os.ReadFile(path); e == nil {
					for _, h := range knownHandles {
						handleCount[h] += strings.Count(string(b), h+"=")
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	p.CodeExts = topKeys(extCount, 5)
	p.CodeDirs = codeRoots(dirCode)
	p.Colocated = detectColocation(stems)
	p.TestHandle = handleDominante(handleCount)
	p.Family = detectFamily(root, extCount)
	p.TestConventions = rankConventions(conventionCount)
	sort.Strings(p.GuideFiles)

	p.Config = p.buildConfig()
	return p, nil
}

// maxAmostraHandle — quantos arquivos de código são lidos à procura do atributo de
// ancoragem. O atributo, quando o projeto o usa, aparece cedo e em toda parte; ler
// tudo custaria I/O por um sinal que a amostra já dá.
const maxAmostraHandle = 300

// minOcorrenciasHandle — abaixo disto a presença é anedótica (um exemplo copiado, um
// arquivo de terceiro) e não uma convenção do projeto. Propor um contrato com base
// em duas ocorrências faria o init sugerir dívida onde não há prática.
const minOcorrenciasHandle = 5

// handleDominante escolhe o atributo mais usado na amostra. Empate ou contagem baixa
// devolve vazio — e vazio significa "o projeto não declara", que faz os gates de
// inventário pularem em vez de acusar.
func handleDominante(cont map[string]int) string {
	melhor, n := "", 0
	for _, h := range knownHandles { // ordem estável: empate resolve pelo 1º
		if cont[h] > n {
			melhor, n = h, cont[h]
		}
	}
	if n < minOcorrenciasHandle {
		return ""
	}
	return melhor
}

func mark(stems map[string]map[string]bool, stem, kind string) {
	if stems[stem] == nil {
		stems[stem] = map[string]bool{}
	}
	stems[stem][kind] = true
}

// detectColocation: co-location existe se há stems com code+spec (ou code+test)
// no mesmo caminho — o sinal de que os derivados ficam ao lado do código.
func detectColocation(stems map[string]map[string]bool) bool {
	hits := 0
	for _, kinds := range stems {
		if kinds["code"] && (kinds["spec"] || kinds["test"] || kinds["feature"]) {
			hits++
		}
	}
	return hits >= 3 // alguns casos = padrão, não coincidência
}

func stemOf(rel string) string {
	dir := filepath.Dir(rel)
	base := filepath.Base(rel)
	for _, suf := range []string{".spec.md", ".feature"} {
		if cut, ok := strings.CutSuffix(base, suf); ok {
			return filepath.Join(dir, cut)
		}
	}
	// A test's stem is its unit's: `foo_test.go`, `test_foo.py` and `foo.test.ts` are all
	// `foo`, the same stem as `foo.go` — the conventions are the same list isTest reads.
	if c, ok := testConventionOf(base); ok {
		return filepath.Join(dir, c.stemOfTest(base))
	}
	if i := strings.LastIndex(base, "."); i >= 0 {
		base = base[:i]
	}
	return filepath.Join(dir, base)
}

func isTest(name string) bool {
	_, ok := testConventionOf(name)
	return ok
}

func isGuide(rel string) bool {
	return strings.Contains(rel, "guides/") && strings.HasSuffix(rel, ".md")
}

// isPlan reconhece um plano pela pasta plans/ (ou plan/). Checado ANTES de isGuide
// e do coringa doc para o plano não ser classificado como doc.
func isPlan(rel string) bool {
	if !strings.HasSuffix(rel, ".md") {
		return false
	}
	return strings.Contains(rel, "plans/") || strings.Contains(rel, "plan/")
}

var codeExts = map[string]bool{
	".ts": true, ".tsx": true, ".js": true, ".jsx": true, ".go": true,
	".py": true, ".rb": true, ".java": true, ".rs": true, ".kt": true,
	".swift": true, ".cs": true, ".php": true,
}

func isCode(name string) bool {
	if strings.HasSuffix(name, ".spec.md") || strings.HasSuffix(name, ".feature") || isTest(name) {
		return false
	}
	return codeExts[filepath.Ext(name)]
}

// topDir is the folder a code file belongs to as a layer candidate: its directory, cut at
// two segments (apps/mobile, src/handlers); "." for a file at the root.
func topDir(rel string) string {
	// Barra normal dos dois lados: o rel chega normalizado e o resultado vira entrada de
	// config ("apps/mobile"), que é a mesma em toda máquina. Com filepath.Separator o
	// Split não separava nada no Windows — cada arquivo virava um diretório distinto.
	dir := path.Dir(rel)
	if dir == "." {
		return "."
	}
	parts := strings.Split(dir, "/")
	if len(parts) > 2 {
		parts = parts[:2]
	}
	return strings.Join(parts, "/")
}

// codeRoots are the project's own folders that hold code, each a candidate layer,
// ordered by volume. There is no minimum: a folder of two files is as much the project's
// decision as one of fifty, and a threshold dropped the small layers of a layered project.
func codeRoots(dirCode map[string]int) []string {
	type kv struct {
		dir string
		n   int
	}
	var kvs []kv
	for d, n := range dirCode {
		if n > 0 {
			kvs = append(kvs, kv{d, n})
		}
	}
	// Ties break by name: map order is random, and a tie decided by it proposed the
	// code layers in a different order on each run.
	sort.Slice(kvs, func(i, j int) bool {
		if kvs[i].n != kvs[j].n {
			return kvs[i].n > kvs[j].n
		}
		return kvs[i].dir < kvs[j].dir
	})
	var out []string
	for _, k := range kvs {
		out = append(out, k.dir)
	}
	return out
}

func topKeys(m map[string]int, n int) []string {
	type kv struct {
		k string
		v int
	}
	var kvs []kv
	for k, v := range m {
		kvs = append(kvs, kv{k, v})
	}
	// Ties break by key: two extensions with the same count came out in map order, so
	// the inferred layer pattern (`*.{ts,js}` vs `*.{js,ts}`) changed between runs.
	sort.Slice(kvs, func(i, j int) bool {
		if kvs[i].v != kvs[j].v {
			return kvs[i].v > kvs[j].v
		}
		return kvs[i].k < kvs[j].k
	})
	var out []string
	for i, k := range kvs {
		if i >= n {
			break
		}
		out = append(out, k.k)
	}
	return out
}
