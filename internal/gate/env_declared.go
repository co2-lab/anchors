// @anchors
//   code: EDGNV
//   ref: ENVDC

package gate

import (
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

var (
	variableCols   = []string{"variable", "variável", "variavel"}
	deprecatedCols = []string{"deprecated", "depreciada", "obsoleta"}
)

// env-declared: the environment variables a unit's code reads are the ones its spec
// declares — and only those.
//
// An environment variable is a contract with whoever deploys: a variable the code reads
// and no spec names is found missing in production, and a variable the spec names and the
// code no longer reads is configured forever for nothing. Like `contract-status-declared`
// with the status codes, both sides are confronted. A variable declared deprecated may stop
// being read: it is on its way out.
//
// The code is the files the spec `specifies`, comments removed, read with the dialect's
// `env_read`; a project that declares none leaves the gate pending, naming the field. A
// declared variable those files do not read is looked for in the code they depend on.
// A read by a computed name (`process.env[name]`) is out of reach of the text.
func checkEnvDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec || g == nil {
		return Skip, ""
	}
	declared := map[string]bool{}
	deprecated := map[string]bool{}
	for _, r := range sectionRows(content, "section.title.environment") {
		v, _ := col(r, variableCols...)
		v = strings.Trim(strings.TrimSpace(v), "`")
		if v == "" || strings.HasPrefix(strings.ToUpper(v), "TODO") {
			continue
		}
		declared[v] = true
		if d, _ := col(r, deprecatedCols...); isYes(d) {
			deprecated[v] = true
		}
	}
	// Each file is read with its own dialect's `env_read`: a Go service and a TS app read a
	// variable in ways that look nothing alike.
	read := map[string]bool{}
	var governed []string
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		governed = append(governed, e.To)
		if pat := cfg.DialectAt(e.To).EnvReadPattern(); pat == "" {
			// How code reads a variable is the language's: with no `env_read`, nothing can be
			// confronted.
			return Pending, i18n.T("gate.field_undeclared", cfg.DialectPathAt(e.To)+".env_read")
		}
		envReads(root, e.To, cfg, read)
	}
	// Where each deployed unit gets only the variables declared for it, what its dependencies
	// read is its contract as well, and is asked of its spec (`follow_dependencies`, by the
	// layer of the code it governs).
	entry := gateEntry(cfg, "env-declared")
	var byDeps map[string]bool
	if followsDependencies(entry, governed, g) {
		byDeps = readByDependencies(governed, root, g, cfg)
		for v := range byDeps {
			read[v] = true
		}
	}
	// What the platform gives every unit is asked of none (`provided`).
	for v := range read {
		if isProvided(entry, v) {
			delete(read, v)
		}
	}
	if len(declared) == 0 && len(read) == 0 {
		return Skip, i18n.T("gate.env_declared.skip_none")
	}
	var undeclared, unread []string
	for v := range read {
		if !declared[v] {
			undeclared = append(undeclared, v)
		}
	}
	for v := range declared {
		if read[v] || deprecated[v] {
			continue
		}
		if byDeps == nil {
			byDeps = readByDependencies(governed, root, g, cfg)
		}
		if !byDeps[v] {
			unread = append(unread, v)
		}
	}
	sort.Strings(undeclared)
	sort.Strings(unread)
	var gaps []string
	if len(undeclared) > 0 {
		gaps = append(gaps, i18n.T("gate.env_declared.undeclared", len(undeclared), strings.Join(undeclared, ", ")))
	}
	if len(unread) > 0 {
		gaps = append(gaps, i18n.T("gate.env_declared.unread", len(unread), strings.Join(unread, ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// envReads adds to `read` the variables the file reads by its dialect's `env_read`, comments
// removed; a file it cannot read, or whose dialect declares none, adds nothing.
func envReads(root, file string, cfg *config.Config, read map[string]bool) {
	d := cfg.DialectAt(file)
	pat := d.EnvReadPattern()
	if pat == "" {
		return
	}
	b, err := readFile(root, file)
	if err != nil {
		return
	}
	re, err := regexp.Compile(pat)
	if err != nil {
		return
	}
	for _, m := range re.FindAllStringSubmatch(stripComments(string(b), file, d), -1) {
		for _, name := range m[1:] {
			if name != "" {
				read[name] = true
				break
			}
		}
	}
}

// readByDependencies are the variables the code the unit depends on reads — the code files
// reached from its own through `depends-on`, transitively. A variable the spec declares and a
// shared helper reads is the unit's deploy contract all the same: measured on a project, a handler's
// spec declared the table its quota helper reads (`handler.ts` → `cotaDeIa.ts` →
// `aiUsage.ts`), and the gate refused the row because the handler did not read it itself.
func readByDependencies(from []string, root string, g *mapx.Graph, cfg *config.Config) map[string]bool {
	read := map[string]bool{}
	seen := map[string]bool{}
	for _, f := range from {
		seen[f] = true
	}
	queue := append([]string(nil), from...)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range g.Neighbors(cur).Out {
			// An import of types alone reads nothing at run time.
			if e.Type != mapx.EdgeDependsOn || e.NoSideEffect || seen[e.To] {
				continue
			}
			seen[e.To] = true
			if t := g.Node(e.To); t == nil || t.Kind != mapx.KindCode {
				continue
			}
			envReads(root, e.To, cfg, read)
			queue = append(queue, e.To)
		}
	}
	return read
}

// followsDependencies says whether a code file the spec governs is of a layer the gate entry
// follows the dependencies of.
func followsDependencies(entry config.Gate, governed []string, g *mapx.Graph) bool {
	for _, f := range governed {
		if n := g.Node(f); n != nil && slices.Contains(entry.FollowDependencies, n.Layer) {
			return true
		}
	}
	return false
}

// isProvided says whether the variable is one the gate entry declares the platform gives.
func isProvided(entry config.Gate, name string) bool {
	for _, p := range entry.Provided {
		if ok, _ := path.Match(p, name); ok {
			return true
		}
	}
	return false
}

// isYes reads a yes in a cell: `yes`, `sim`, `sí`, `true`, or one followed by a note
// (`yes: use PAYMENTS_KEY`).
func isYes(cell string) bool {
	c := strings.ToLower(strings.Trim(strings.TrimSpace(cell), "`*"))
	for _, y := range []string{"yes", "sim", "sí", "si", "true", "x", "✓"} {
		if c == y || strings.HasPrefix(c, y+":") || strings.HasPrefix(c, y+" ") || strings.HasPrefix(c, y+",") {
			return true
		}
	}
	return false
}
