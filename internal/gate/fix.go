// @anchors
//   code: FXGTF
//   ref: FXIXX

package gate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/migra"
	"github.com/co2-lab/anchors/internal/scan"
)

// O self-healer: alguns checks são REPARÁVEIS — a correção é mecânica e segura (ex.:
// escrever a data do commit no `updated_at`). `anchors check --fix` aplica esses
// reparos. Só checks com um fixer REGISTRADO são reparáveis; os demais só reportam.

// FixResult descreve um reparo aplicado (ou tentado).
type FixResult struct {
	Gate   string
	Target string
	Fixed  bool
	Detail string
}

// fixers: check → função que repara o arquivo. Recebe o conteúdo e o contexto (o mapa,
// para o reparo que precisa saber a unidade do arquivo), devolve o novo conteúdo, se
// mudou algo e o que foi feito. Só reparos determinísticos e seguros.
var fixers = map[string]func(content string, n mapx.Node, root string, g *mapx.Graph) (string, bool, string){
	"updated-at-atual": func(content string, n mapx.Node, root string, _ *mapx.Graph) (string, bool, string) {
		out, changed := fixUpdatedAt(content, n, root)
		return out, changed, i18n.T("gate.fix.updated_at_fixed")
	},
	"header-valid": fixMissingHeader,
}

// Fixable diz se um check tem reparo automático.
func Fixable(check string) bool {
	_, ok := fixers[check]
	return ok
}

// Fix roda o FIXER de cada gate reparável sobre todo nó a que o gate se aplica e grava
// as correções em disco. Devolve o que foi consertado (ou tentado). Não consulta o
// veredito do gate: é o fixer quem decide se há o que corrigir, e um nó já correto
// não muda e não entra no resultado.
func Fix(gates []config.Gate, nodes []mapx.Node, root string, g *mapx.Graph) []FixResult {
	var out []FixResult
	for _, gt := range gates {
		fix, ok := fixers[gt.Check]
		if !ok {
			continue
		}
		for _, n := range nodes {
			if !applies(gt, n, root) {
				continue
			}
			path := filepath.Join(root, n.ID)
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			newContent, changed, detail := fix(string(content), n, root, g)
			if !changed {
				continue
			}
			if err := os.WriteFile(path, []byte(newContent), 0o644); err != nil {
				// Os detalhes passam por i18n: eram literais em português, impressos
				// assim em projeto de qualquer `lang:`.
				out = append(out, FixResult{gt.Name, n.ID, false, i18n.T("gate.fix.write_failed", err)})
				continue
			}
			out = append(out, FixResult{gt.Name, n.ID, true, detail})
		}
	}
	return out
}

// fixUpdatedAt reescreve o `updated_at` do header com a data do último commit — só se
// diverge e o arquivo está commitado (nada a corrigir num arquivo com edição pendente,
// cuja data ainda vai mudar). Não cria o campo se ele não existe (isso é trabalho do
// autor + o gate header-conforme); só CORRIGE um valor errado.
func fixUpdatedAt(content string, n mapx.Node, root string) (string, bool) {
	m := updatedAtRE.FindStringSubmatchIndex(content)
	if m == nil {
		return content, false // sem campo — nada a corrigir aqui (o autor cria)
	}
	// a data CORRETA: hoje se há edição não-commitada (alteração em curso), senão a
	// data do último commit. Espelha a regra do checkUpdatedAt.
	var correct string
	mudou, sabido := gitmeta.UncommittedChanges(root, n.ID)
	if !sabido {
		return content, false // sem git: não há data correta a apurar, e chutar seria pior
	}
	if mudou {
		correct = gitmeta.Today()
	} else if d, ok := gitmeta.LastCommitDate(root, n.ID); ok {
		correct = d
	} else {
		return content, false // sem commit e sem edição — nada a comparar
	}
	// m[2]:m[3] é o grupo capturado (a data). Substitui só ela se diverge.
	if content[m[2]:m[3]] == correct {
		return content, false
	}
	return content[:m[2]] + correct + content[m[3]:], true
}

// fixMissingHeader writes the `@anchors` header a governed file lacks, at its top: the
// `ref:` of the unit the map ties it to — the spec that specifies it, or through its
// feature the spec of the test — and, for a file that belongs to no unit (a guide, a
// document), its `layer:`. A file that has a header is left as it is, even a wrong one:
// correcting an identity someone wrote is not mechanical. A file with neither a unit nor
// a layer is not repaired: its identity has to be decided.
//
// It exists because the header became due on every governed file, and a project that
// updates would otherwise have hundreds to write by hand — the map already knows each
// file's unit.
func fixMissingHeader(content string, n mapx.Node, _ string, g *mapx.Graph) (string, bool, string) {
	if isBinary(content) || isExecutableScript(n) {
		return content, false, ""
	}
	// With no header at the top, one is written there — even when a header-looking block
	// stands lower: in a guide or a template that is an example inside a string or a code
	// fence, and the top header makes it read as the text it is.
	block := scan.AnchorsHeader([]byte(content))
	var b string
	if block != nil {
		b = string(block)
	}
	var fields []string
	// The file's OWN code, generated from its name and its type, unique in the map.
	if !headerOwnCodeRE().MatchString(b) && g != nil {
		taken := map[string]bool{}
		for _, x := range g.Nodes {
			for _, c := range []string{x.Code, x.FileCode} {
				if c != "" {
					taken[c] = true
				}
			}
		}
		f := scan.File{Path: n.ID, Layer: n.Layer, Kind: string(n.Kind)}
		fields = append(fields, "code: "+code.GenerateUnique(migra.FileCodeName(f), taken))
	}
	if !headerCodeRE().MatchString(b) && !headerRefRE().MatchString(b) && !headerLayerRE.MatchString(b) {
		var field string
		if g != nil {
			if codes := g.UnitCodesOf(n.ID); len(codes) > 0 {
				field = "ref: " + strings.Join(codes, ", ")
			}
		}
		if field == "" && n.Layer != "" && (n.Kind == mapx.KindGuide || n.Kind == mapx.KindDoc || n.Support) {
			field = "layer: " + n.Layer
		}
		if field != "" {
			fields = append(fields, field)
		}
	}
	// A file whose identity cannot be told gets nothing: a code alone would read as a unit
	// of its own.
	hasIdentity := headerCodeRE().MatchString(b) || headerRefRE().MatchString(b) || headerLayerRE.MatchString(b)
	if len(fields) == 0 || (!hasIdentity && len(fields) == 1 && strings.HasPrefix(fields[0], "code: ")) {
		return content, false, ""
	}
	what := strings.Join(fields, ", ")
	c := config.LineCommentFor(n.ID)
	indent := c + "   "
	if c == "<!--" {
		indent = "  "
	}
	var lines strings.Builder
	for _, f := range fields {
		lines.WriteString("\n" + indent + f)
	}
	// A header at the top gets the missing lines right below `@anchors`: nothing someone
	// wrote is changed, the identity is added.
	if block != nil {
		at := strings.Index(content, b)
		opener := at + strings.Index(b, "\n")
		if opener < at {
			opener = at + len(block)
		}
		return content[:opener] + lines.String() + content[opener:], true, i18n.T("gate.fix.header_written", what)
	}
	var header string
	if c == "<!--" {
		header = "<!-- @anchors" + lines.String() + "\n-->\n\n"
	} else {
		header = c + " @anchors" + lines.String() + "\n\n"
	}
	// A shebang stays the first line: the system reads it there.
	if strings.HasPrefix(content, "#!") {
		if i := strings.Index(content, "\n"); i >= 0 {
			return content[:i+1] + header + content[i+1:], true, i18n.T("gate.fix.header_written", what)
		}
	}
	return header + content, true, i18n.T("gate.fix.header_written", what)
}
