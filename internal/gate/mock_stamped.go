package gate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// mock-carimbado: o dublê carrega a marca do trecho que ele substitui, e o gate
// RECALCULA essa marca contra o módulo real.
//
// É a camada agnóstica do par que ataca mock drift. O `mock-tipado` resolve o problema
// onde a linguagem ajuda (tipos estruturais fazem o compilador conferir a superfície do
// módulo), e só ali: a maior parte dos ecossistemas não tem `Partial<typeof X>`, e nem
// em TypeScript ele alcança a forma do VALOR devolvido — medido, 202 testes dublavam um
// hook do React Query com dois campos onde o real devolve ~15, e o tipo não distingue
// "parcial deliberado" de "defasado".
//
// O carimbo não interpreta o código: ele o LÊ. Hash de texto pega qualquer mudança —
// assinatura, corpo, tipo, constante — e funciona em Python, Ruby, Go ou JS puro, sem
// extrator por linguagem.
//
// O gate RECALCULA em vez de validar formato. Um carimbo que ninguém confronta é
// teatro: quem edita o teste o regeneraria para casar com o próprio mock, e ele passaria
// a certificar a si mesmo. Recalcular é o que torna o mecanismo à prova de quem escreve.
func checkMockStamped(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest {
		return Skip, i18n.T("gate.mock_stamped.skip_not_test")
	}
	detector, err := doubleDetector(cfg)
	if err != nil {
		// Regex inválido é erro de CONFIGURAÇÃO, não do arquivo sob análise — e precisa
		// falhar alto: silenciá-lo faria o gate varrer zero dublês e reportar verde.
		return Fail, fmt.Sprintf(i18n.T("gate.mock_stamped.fail_mock_detect_compile"), err)
	}
	if detector == nil {
		return Skip, i18n.T("gate.mock_stamped.skip_no_mock_detect")
	}

	carimbos := declaredStamps(content)

	// AUSÊNCIA de carimbo é acusada, não pulada — e isto é a metade mais importante do
	// gate.
	//
	// Um carimbo divergente ACUSA; um carimbo ausente é SILÊNCIO, que é o mesmo "falha
	// aberto" que o `trinca-completa` existe para fechar. Se a ausência passasse, o
	// carimbo viraria opcional na prática: quem não põe nunca é cobrado, e o mecanismo
	// protegeria apenas quem já escolheu ser protegido — isso é convenção, não gate.
	//
	// A tentação de pular vem de olhar projeto LEGADO, onde ligar a cobrança produz
	// centenas de achados de uma vez. Mas desenhar por esse caso transformaria um
	// problema de migração numa propriedade permanente do framework: todo projeto
	// futuro herdaria a frouxidão para acomodar um projeto antigo. Num projeto que
	// NASCE com o Anchors não há dívida — o primeiro dublê é escrito depois do gate
	// existir, e cobrar custa um carimbo.
	//
	// O legado se resolve com o vocabulário que já existe: `blocking: false` durante a
	// adoção, e opt-out por unidade com razão escrita onde a decisão for deliberada.
	var semCarimbo []string
	for _, modulo := range detectedDoubles(content, detector) {
		if !isGovernedModule(modulo, g) {
			continue // biblioteca de terceiro não é cobrada (ver `ehModuloRegido`)
		}
		if !moduleHasStamp(modulo, carimbos) {
			semCarimbo = append(semCarimbo, modulo)
		}
	}
	if len(semCarimbo) > 0 {
		return Fail, fmt.Sprintf(i18n.T("gate.mock_stamped.fail_missing_stamps"),
			len(semCarimbo), strings.Join(semCarimbo, ", "))
	}
	if len(carimbos) == 0 {
		return Skip, i18n.T("gate.mock_stamped.skip_no_governed_modules")
	}

	var divergentes []string
	for _, c := range carimbos {
		atual, err := recomputeStamp(root, c)
		if err != nil {
			divergentes = append(divergentes, fmt.Sprintf("%s (%v)", c.anchor, err))
			continue
		}
		if atual != c.hash {
			divergentes = append(divergentes, fmt.Sprintf(i18n.T("gate.mock_stamped.stamp_diff_line"),
				c.file, c.anchor, c.hash, atual))
		}
	}
	if len(divergentes) == 0 {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(i18n.T("gate.mock_stamped.fail_drift"),
		len(divergentes), strings.Join(divergentes, "\n  - "))
}

// moduleHasStamp liga o dublê ao carimbo pelo CAMINHO do arquivo carimbado.
//
// O especificador do dublê (`@/src/hooks/useX`) e o caminho do carimbo
// (`apps/mobile/src/hooks/useX.ts`) descrevem o mesmo arquivo por vias diferentes —
// alias e caminho de disco. Casar pelo sufixo sem extensão resolve os dois sem
// precisar de um resolvedor de alias, que seria específico do ecossistema.
//
// The specifier is compared AS WRITTEN and without a final extension, and both forms
// count. An import specifier usually has no extension, and cutting one anyway cut part of
// the NAME: `@/src/stores/auth.store` became `src/stores/auth`, which no stamp file
// (`auth.store.ts` → `auth.store`) matches. Measured in a real project after stamping all
// its doubles: 93 tests still failed as "double without stamp", all on modules with a dot
// in the name (`auth.store`, `ui.store`). ESM specifiers that DO carry an extension
// (`./x.js`) are why the stripped form is kept too.
func moduleHasStamp(modulo string, carimbos []declaredStamp) bool {
	alvo := strings.TrimPrefix(modulo, "./")
	for strings.HasPrefix(alvo, "../") {
		alvo = strings.TrimPrefix(alvo, "../")
	}
	if strings.HasPrefix(alvo, "@") {
		if i := strings.Index(alvo, "/"); i >= 0 {
			alvo = alvo[i+1:]
		}
	}
	forms := []string{alvo, withoutExtension(alvo)}
	for _, c := range carimbos {
		arq := withoutExtension(c.file)
		for _, f := range forms {
			if arq == f || strings.HasSuffix(arq, "/"+f) {
				return true
			}
		}
	}
	return false
}

// declaredStamp — o que o teste afirma sobre o trecho que dubla.
type declaredStamp struct {
	file   string // caminho do módulo, relativo à raiz
	anchor string // a LINHA INTEIRA que abre o trecho (conteúdo, nunca número)
	count  int    // quantas linhas a partir da âncora entram no hash
	hash   string // o hash gravado
}

// stampRE casa a anotação:
//
//	// @contract: caminho/do/modulo.ts | export function useX( | 10 | 361280fb
//
// O separador é `|` porque a âncora é uma linha de código e pode conter vírgula, dois
// pontos e parênteses — qualquer separador mais comum a partiria no meio.
var stampRE = regexp.MustCompile(
	`@contract:\s*([^|\n]+?)\s*\|\s*(.+?)\s*\|\s*(\d+)\s*\|\s*([0-9a-f]+)`)

func declaredStamps(content string) []declaredStamp {
	var out []declaredStamp
	for _, m := range stampRE.FindAllStringSubmatch(content, -1) {
		count, err := strconv.Atoi(m[3])
		if err != nil || count <= 0 {
			continue
		}
		out = append(out, declaredStamp{
			file: m[1], anchor: m[2], count: count, hash: m[4],
		})
	}
	return out
}

// recomputeStamp lê o módulo real e devolve o hash do trecho HOJE.
//
// A âncora é procurada por CONTEÚDO — é o que torna o carimbo imune a deslocamento.
// Duas ocorrências da mesma linha tornam o alvo ambíguo, e o gate prefere acusar a
// escolher uma: um carimbo que aponta para "alguma das duas" não prova nada.
func recomputeStamp(root string, c declaredStamp) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, c.file))
	if err != nil {
		return "", fmt.Errorf(i18n.T("gate.mock_stamped.err_module_not_found"), c.file)
	}
	linhas := strings.Split(string(b), "\n")

	idx := -1
	ocorrencias := 0
	for i, l := range linhas {
		if strings.TrimRight(l, "\r") == c.anchor {
			ocorrencias++
			if idx < 0 {
				idx = i
			}
		}
	}
	switch {
	case ocorrencias == 0:
		// A âncora sumiu: renome, remoção ou reescrita da linha. É achado, não erro de
		// ferramenta — e falha EXPLÍCITA é melhor que silêncio, porque o dublê
		// certamente está desatualizado.
		return "", fmt.Errorf("%s", i18n.T("gate.mock_stamped.err_anchor_not_found"))
	case ocorrencias > 1:
		return "", fmt.Errorf(i18n.T("gate.mock_stamped.err_anchor_ambiguous"), ocorrencias)
	}

	fim := idx + c.count
	if fim > len(linhas) {
		fim = len(linhas)
	}
	return snippetHash(strings.Join(linhas[idx:fim], "\n")), nil
}

// snippetHash — sha256 truncado em 8 hex. Truncado porque o carimbo mora numa linha de
// comentário e é lido por humano; 32 bits bastam para detectar mudança acidental, que é
// o que este gate persegue (não há adversário forjando colisão contra o próprio teste).
func snippetHash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])[:8]
}

// doubleDetector compila o regex que ESTE projeto usa para escrever um dublê.
//
// Devolve (nil, nil) quando o projeto não declara — e é o que desliga o gate. A
// alternativa (embutir o padrão jest/vitest como default) faria o gate rodar num
// projeto Python, casar zero dublês e reportar VERDE sobre o que não conferiu.
func doubleDetector(cfg *config.Config) (*regexp.Regexp, error) {
	if cfg == nil || cfg.Derived == nil {
		return nil, nil
	}
	p := strings.TrimSpace(cfg.Derived.MockDetect)
	if p == "" {
		return nil, nil
	}
	re, err := regexp.Compile(p)
	if err != nil {
		return nil, err
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("%s", i18n.T("gate.mock_stamped.err_capture_group_required"))
	}
	return re, nil
}

// detectedDoubles aplica o padrão do projeto e devolve os módulos dublados.
func detectedDoubles(content string, re *regexp.Regexp) []string {
	var out []string
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
			out = append(out, m[1])
		}
	}
	return out
}

// TestsStamping returns the test nodes whose `@contract` stamps point at `file`.
//
// It is what makes the person who CHANGES a module see the doubles it breaks. The gate
// recomputes a stamp when the TEST is checked, and a `check --changed` of the module
// alone never checked its tests: the drift surfaced days later, to whoever next touched
// one of them, far from the change that caused it. With these tests in the impact of the
// changed module, the same pre-commit that commits the change lists every double that now
// reproduces the old contract — and the author updates them in the same commit.
func TestsStamping(g *mapx.Graph, root, file string) []string {
	if g == nil || file == "" {
		return nil
	}
	var out []string
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindTest {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, n.ID))
		if err != nil || !strings.Contains(string(b), "@contract:") {
			continue
		}
		for _, c := range declaredStamps(string(b)) {
			if filepath.ToSlash(c.file) == filepath.ToSlash(file) {
				out = append(out, n.ID)
				break
			}
		}
	}
	return out
}

// StampRefresh is one `@contract` stamp that `RefreshStamps` found diverging from the
// module it points at.
type StampRefresh struct {
	Test   string // the test holding the stamp
	Line   int    // 1-based line of the stamp in the test
	File   string // the module the stamp points at
	Anchor string
	Count  int
	Old    string // the hash written in the stamp
	New    string // the hash of the snippet today; empty when it cannot be recomputed
	Err    string // why the stamp could not be refreshed (anchor gone or ambiguous)
}

// RefreshStamps updates the hash of every stamp that points at `file` and no longer matches
// it, in every test of the map, and returns what it found.
//
// It is the command side of "who changes a function updates its doubles": run by the
// author of the change, it answers with the list of doubles that reproduced the old
// contract — the ones to adjust in the same commit. A stamp whose anchor line is gone or
// ambiguous is NOT refreshed: the member was renamed or removed, and only a person can say
// which new line the double now stands for.
//
// With write=false it only reports.
func RefreshStamps(g *mapx.Graph, root, file string, write bool) ([]StampRefresh, error) {
	var out []StampRefresh
	for _, t := range TestsStamping(g, root, file) {
		path := filepath.Join(root, t)
		b, err := os.ReadFile(path)
		if err != nil {
			return out, err
		}
		lines := strings.Split(string(b), "\n")
		changed := false
		for i, l := range lines {
			m := stampRE.FindStringSubmatchIndex(l)
			if m == nil {
				continue
			}
			c := declaredStamp{file: l[m[2]:m[3]], anchor: l[m[4]:m[5]], hash: l[m[8]:m[9]]}
			c.count, _ = strconv.Atoi(l[m[6]:m[7]])
			if filepath.ToSlash(c.file) != filepath.ToSlash(file) || c.count <= 0 {
				continue
			}
			r := StampRefresh{Test: t, Line: i + 1, File: c.file, Anchor: c.anchor, Count: c.count, Old: c.hash}
			atual, err := recomputeStamp(root, c)
			if err != nil {
				r.Err = err.Error()
				out = append(out, r)
				continue
			}
			if atual == c.hash {
				continue
			}
			r.New = atual
			out = append(out, r)
			lines[i] = l[:m[8]] + atual + l[m[9]:]
			changed = true
		}
		if changed && write {
			if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
				return out, err
			}
		}
	}
	return out, nil
}

// StampSnippet returns the `count` lines that start at the unique `anchor` line of
// `content` — the block a stamp covers — or false when the anchor is absent or repeated.
func StampSnippet(content, anchor string, count int) (string, bool) {
	linhas := strings.Split(content, "\n")
	idx := -1
	for i, l := range linhas {
		if strings.TrimRight(l, "\r") == anchor {
			if idx >= 0 {
				return "", false
			}
			idx = i
		}
	}
	if idx < 0 {
		return "", false
	}
	fim := idx + count
	if fim > len(linhas) {
		fim = len(linhas)
	}
	return strings.Join(linhas[idx:fim], "\n"), true
}
