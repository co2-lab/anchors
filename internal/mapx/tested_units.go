package mapx

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
)

// TestedUnits maps each test file of the map to the code files it tests, found by the
// project's own derivation (`derived.files` and the overrides of the code's layer).
//
// The map links a unit's test only through the triad: a spec, its feature, the feature's
// test. A unit with no spec — a util, a model, what a `regime: declarativo` layer holds —
// has no edge to its test at all, and a gate asking "what does this test test?" had no
// answer (in the reference app, 130 tests with no feature had 12 edges among them). The
// derivation answers it with nothing the engine assumes about the language: the same
// templates that say where a unit's test lives, applied to each code file.
//
// With `anchor: code` the code file is the unit, and its dir, name and extension fill the
// templates as they do in the map build. With `anchor: spec` the templates hang off the
// spec's name, so the code file's own path is matched against the `code` templates to
// recover the variables first; a code template is read literally around its variables.
func TestedUnits(g *Graph, cfg *config.Config) map[string][]string {
	if g == nil || cfg == nil || cfg.Derived == nil {
		return nil
	}
	tests := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Kind == KindTest {
			tests[n.ID] = true
		}
	}
	out := map[string][]string{}
	for _, n := range g.Nodes {
		if n.Kind != KindCode {
			continue
		}
		byLayer := templatesFor(cfg, n.Layer)
		for _, v := range unitVars(n, cfg, byLayer["code"]) {
			for _, tmpl := range byLayer["test"] {
				want := resolveTemplateM(tmpl, v.dir, v.name, v.ext, v.module)
				// The exact path first: a directory named `[slug]` is a path, not a
				// character class. Only a path the map does not have is tried as a glob.
				if tests[want] {
					out[want] = appendOnce(out[want], n.ID)
					continue
				}
				if !strings.ContainsAny(want, "*?[{") {
					continue
				}
				for t := range tests {
					if ok, _ := doublestar.Match(want, t); ok {
						out[t] = appendOnce(out[t], n.ID)
					}
				}
			}
		}
	}
	for t := range out {
		sort.Strings(out[t])
	}
	return out
}

// templatesFor is the project's derivation for a unit of `layer`: the default templates,
// each kind replaced by an override whose `when` is that layer.
func templatesFor(cfg *config.Config, layer string) map[string]config.Padroes {
	out := map[string]config.Padroes{}
	for kind, ps := range cfg.Derived.PadroesDe() {
		out[kind] = ps
	}
	for _, ov := range cfg.Derived.Overrides {
		if ov.Code != "" || ov.When != layer {
			continue
		}
		for kind, ps := range ov.PadroesDe() {
			out[kind] = ps
		}
	}
	return out
}

// templateVars are the values a derivation template is filled with.
type templateVars struct{ dir, name, ext, module string }

// unitVars recovers the template variables of a code file: from its own path when code is
// the anchor, from matching it against the `code` templates otherwise.
func unitVars(n Node, cfg *config.Config, codeTemplates config.Padroes) []templateVars {
	dir := filepath.ToSlash(filepath.Dir(n.ID))
	if cfg.Derived.Anchor == "code" {
		name, ext := StemOfAnchor(n.ID)
		return []templateVars{{dir: dir, name: name, ext: ext, module: filepath.Base(dir)}}
	}
	var out []templateVars
	for _, tmpl := range codeTemplates {
		if v, ok := reverseTemplate(tmpl, n.ID); ok {
			out = append(out, v)
		}
	}
	return out
}

var templateVarRE = regexp.MustCompile(`\{\{(dir|name|ext|module)\}\}`)

// reverseTemplate matches a path against a template and returns the variables that
// produce it. A variable absent from the template is left as the path implies (the dir
// and module of the file). Everything else in the template is read literally: a
// directory named `[slug]` matches itself, and glob characters match only themselves.
func reverseTemplate(tmpl, path string) (templateVars, bool) {
	var b strings.Builder
	b.WriteString("^")
	var order []string
	last := 0
	for _, loc := range templateVarRE.FindAllStringSubmatchIndex(tmpl, -1) {
		b.WriteString(regexp.QuoteMeta(tmpl[last:loc[0]]))
		v := tmpl[loc[2]:loc[3]]
		switch v {
		case "dir":
			b.WriteString(`(.+)`)
		case "ext":
			b.WriteString(`([^/.]+)`)
		default: // name, module
			b.WriteString(`([^/]+)`)
		}
		order = append(order, v)
		last = loc[1]
	}
	b.WriteString(regexp.QuoteMeta(tmpl[last:]) + "$")
	m := regexp.MustCompile(b.String()).FindStringSubmatch(path)
	if m == nil {
		return templateVars{}, false
	}
	dir := filepath.ToSlash(filepath.Dir(path))
	v := templateVars{dir: dir, module: filepath.Base(dir)}
	for i, name := range order {
		switch name {
		case "dir":
			v.dir = m[i+1]
		case "name":
			v.name = m[i+1]
		case "ext":
			v.ext = m[i+1]
		case "module":
			v.module = m[i+1]
		}
	}
	return v, true
}

func appendOnce(xs []string, x string) []string {
	for _, y := range xs {
		if y == x {
			return xs
		}
	}
	return append(xs, x)
}
