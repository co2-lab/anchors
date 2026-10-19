// @anchors
//   code: THGTS
//   ref: TICTS

// Helpers de testID: reconhecer o handle nas duas pontas que o ESCREVEM — o código
// (pelos padrões do projeto, `dialect.test_handles`) e a spec (tabela de inventário).
//
// Este arquivo já foi `testid_declared.go`, sede do gate homônimo. O gate foi
// aposentado por `testid-coerente`, que confronta a unidade inteira em vez de duas
// pontas; os reconhecedores continuam aqui porque são dele que todos partem.
package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// handlesIn returns the handles the patterns find in a source: a `testid` group as the
// handle, a `prefix` group as `head-*` and a `suffix` group as `*-seg-*`. The patterns are
// the project's (`dialect.test_handles`): the core reads no syntax of its own. It knew JSX —
// the attribute, a template, a prop ending in Prefix, a conditional — and a project whose
// handles a helper composes (`getOptionTestID(base, option)`, at MIF) had no way to say so.
//
// A value one pattern captures as a prefix is not also a whole handle by another: a prop
// that gives a child the head of its handles (`testIDPrefix=":otp-input"`) also matches the
// pattern of the attribute, and `otp-input` never shows on screen.
func handlesIn(src, rel string, patterns config.Patterns) []string {
	type span struct{ from, to int }
	type found struct {
		at    span
		value string // as it reads: `id`, `head-*`, `*-seg-*`
		whole bool
	}
	heads := map[span]bool{}
	var all []found
	for _, m := range patterns.FindAll(src, rel) {
		if c, ok := m.Groups["prefix"]; ok && strings.TrimSuffix(c.Value, "-") != "" {
			heads[span{c.Start, c.End}] = true
			all = append(all, found{span{c.Start, c.End}, strings.TrimSuffix(c.Value, "-") + "-*", false})
		}
		if c, ok := m.Groups["suffix"]; ok {
			if seg := strings.TrimRight(c.Value, "-"); measurableSegment(seg) {
				all = append(all, found{span{c.Start, c.End}, "*" + seg + "-*", false})
			}
		}
		if c, ok := m.Groups["testid"]; ok && c.Value != "" {
			all = append(all, found{span{c.Start, c.End}, c.Value, true})
		}
	}
	// In the order they stand in the source, whichever pattern found them.
	sort.SliceStable(all, func(i, j int) bool { return all[i].at.from < all[j].at.from })
	seen := map[string]bool{}
	var out []string
	for _, f := range all {
		if (f.whole && heads[f.at]) || seen[f.value] {
			continue
		}
		seen[f.value] = true
		out = append(out, f.value)
	}
	return out
}

// measurableSegment says whether a derived handle's own segment names anything: a segment
// with no letter — none at all (`${prefix}${key}`), or a separator (`${prefix}-${i}`) — leaves
// the whole handle to values the text does not hold.
func measurableSegment(seg string) bool {
	return strings.IndexFunc(seg, unicode.IsLetter) >= 0
}

// unmeasurableHandles are the places — `file:line` — where the patterns find a handle made of
// the caller's head and a segment that names nothing: there is a handle, and no text says
// which. A gate answers it pending, never as measured.
func unmeasurableHandles(src, rel string, patterns config.Patterns) []string {
	var out []string
	for _, m := range patterns.FindAll(src, rel) {
		if c, ok := m.Groups["suffix"]; ok && !measurableSegment(strings.TrimRight(c.Value, "-")) {
			out = append(out, fmt.Sprintf("%s:%d", rel, config.LineOf(src, m.Start)))
		}
	}
	return out
}

// projectDerivedSegments are the segments the units of the project derive from a head their
// caller gives — `-option-` of `*-option-*` —, read once per map: a caller's handle is
// queried through them.
func projectDerivedSegments(root string, g *mapx.Graph, exposed config.Patterns) []string {
	derivedMu.Lock()
	defer derivedMu.Unlock()
	if derivedFor == g && derivedSegs != nil {
		return derivedSegs
	}
	seen := map[string]bool{}
	segs := []string{}
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindCode {
			continue
		}
		b, err := readFile(root, n.ID)
		if err != nil {
			continue // @resilient: a file that cannot be read derives no segment; the map notices it
		}
		for _, h := range handlesIn(string(b), n.ID, exposed) {
			if seg, ok := derivedSegment(h); ok && !seen[seg] {
				seen[seg] = true
				segs = append(segs, seg)
			}
		}
	}
	sort.Strings(segs)
	derivedFor, derivedSegs = g, segs
	return segs
}

var (
	derivedMu   sync.Mutex
	derivedFor  *mapx.Graph
	derivedSegs []string
)

// testHandles are the patterns the project declares for its test handles, or nil.
func testHandles(cfg *config.Config) *config.TestHandles {
	if cfg == nil || cfg.Dialect == nil {
		return nil
	}
	return cfg.Dialect.TestHandles
}

// exposedPatterns are the patterns of where the code exposes a handle; none when the
// project declares none, and then the handle gates measure nothing.
func exposedPatterns(cfg *config.Config) config.Patterns {
	if th := testHandles(cfg); th != nil {
		return th.Exposed
	}
	return nil
}

// testHandleAttr — o atributo com que ESTE projeto marca elementos para alcance
// externo. Sem declaração devolve vazio, e os gates de inventário pulam: o Anchors
// não adivinha o ecossistema, do mesmo modo que não adivinha idioma nem vendor.
func testHandleAttr(cfg *config.Config) string {
	if cfg == nil || cfg.Derived == nil {
		return ""
	}
	return strings.TrimSpace(cfg.Derived.TestHandle)
}

// declaredTestIDRE captura o id citado em crase na spec. Aceita as três formas
// catalogadas que o projeto já usa (SPEC_TYPES §5): linha de tabela, bullet e
// cabeçalho — o que importa é o id estar em crase, não a moldura em volta.
var declaredTestIDRE = regexp.MustCompile("`(:?[a-zA-Z][a-zA-Z0-9._-]*(?:-\\*)?|\\*-[a-zA-Z][a-zA-Z0-9._-]*-\\*)`")

// surfaceSectionRE delimita a seção de inventário. Sem ela, qualquer crase na prosa
// da spec (nome de classe CSS, de campo, de arquivo) contaria como declaração — e o
// gate passaria a aprovar por acidente, que é pior que reprovar por engano.
//
// O título aceita QUALIFICADOR depois do nome (`## Test IDs (Maestro)`): 42 das 159
// specs do app de referência nomeiam ali a superfície que consome os ids. Exigir fim-de-linha logo
// após o título fazia o gate não enxergar a seção dessas specs e acusá-las de não
// declarar nada — reprovando justamente as que documentam melhor.
var surfaceSectionRE = regexp.MustCompile(`(?im)^#{1,6}\s*(?:superf[íi]cie de teste|test\s*ids?)\b.*$`)

// declaredTestIDs lê o inventário DENTRO da seção de Superfície de Teste.
func declaredTestIDs(spec, attr string) []string {
	loc := surfaceSectionRE.FindStringIndex(spec)
	if loc == nil {
		return nil
	}
	corpo := spec[loc[1]:]
	// A seção termina no próximo cabeçalho de mesmo nível ou acima.
	if fim := regexp.MustCompile(`(?m)^#{1,6}\s`).FindStringIndex(corpo); fim != nil {
		corpo = corpo[:fim[0]]
	}
	visto := map[string]bool{}
	var out []string
	for _, linha := range strings.Split(corpo, "\n") {
		// Numa linha de TABELA, só a PRIMEIRA célula é o id — as demais descrevem o
		// elemento e citam os cenários que o usam, e essas colunas também levam crase.
		// Ler a linha inteira fazia o gate colher `TouchableOpacity` (coluna Elemento)
		// e `ATLNX-VR` (coluna "Usado em") como se fossem testIDs declarados, e depois
		// acusá-los de órfãos — inventando dívida a partir da própria documentação.
		alvo := linha
		if celulas := rowCells(linha); celulas != nil {
			alvo = celulas[0]
		}
		for _, m := range declaredTestIDRE.FindAllStringSubmatch(alvo, -1) {
			// O PRÓPRIO nome do atributo não é um id. O átomo genérico (Button,
			// Avatar, ActionLink) recebe o handle de fora e a spec o documenta como
			// `testID` (prop) — declarar isso é dizer "aceito um handle", não "exponho
			// este". Cobrá-lo faria o gate exigir do átomo um id fixo, que é o oposto
			// do que um componente reusável deve ter.
			if strings.EqualFold(strings.TrimPrefix(m[1], ":"), attr) {
				continue
			}
			if !visto[m[1]] {
				visto[m[1]] = true
				out = append(out, m[1])
			}
		}
	}
	return out
}

// rowCells devolve as células de uma linha de tabela Markdown, ou nil se a
// linha não for tabela. Separador (`| --- |`) não tem conteúdo e devolve nil.
func rowCells(linha string) []string {
	t := strings.TrimSpace(linha)
	if !strings.HasPrefix(t, "|") {
		return nil
	}
	if strings.Trim(t, "| -:") == "" {
		return nil // linha separadora
	}
	partes := strings.Split(strings.Trim(t, "|"), "|")
	for i := range partes {
		partes[i] = strings.TrimSpace(partes[i])
	}
	return partes
}

// covers diz se `padrao` cobre `id`, comparando sem o `:` de marcação — a spec pode
// citar `bdgt-screen` ou `:bdgt-screen` e as duas dizem a mesma coisa. Quem compara
// chama nos DOIS sentidos, porque o curinga casa dos dois lados: exposto `otp-input-*` covers o declarado
// `otp-input-0`, e declarado `abcd-item-*` covers o exposto `abcd-item-3`. Sem isso o
// gate acusaria em ambos os sentidos um id que o outro lado descreve corretamente —
// só que na forma genérica em vez da concreta, ou vice-versa.
func covers(padrao, id string) bool {
	p := strings.TrimPrefix(padrao, ":")
	i := strings.TrimPrefix(id, ":")
	if p == i {
		return true
	}
	// A DERIVED handle — `*-option-*`, its head the caller's — covers an id that carries
	// its segment after a head.
	if seg, ok := derivedSegment(p); ok {
		j := strings.Index(i, seg)
		return j > 0 || (strings.HasPrefix(i, "*") && j >= 0)
	}
	if strings.HasSuffix(p, "-*") {
		return strings.HasPrefix(i, strings.TrimSuffix(p, "*"))
	}
	return false
}

// derivedSegment is the segment a derived handle adds to its caller's head — `-option-` of
// `*-option-*` —, and whether the handle is one.
func derivedSegment(handle string) (string, bool) {
	h := strings.TrimPrefix(handle, ":")
	if !strings.HasPrefix(h, "*") {
		return "", false
	}
	return strings.TrimSuffix(strings.TrimPrefix(h, "*"), "*"), true
}
