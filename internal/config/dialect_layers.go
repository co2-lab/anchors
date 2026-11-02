// @anchors
//   code: DLLYD
//   ref: DLCTI

package config

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

// A project may speak more than one language: apps in one, services in another, in the same
// repository so the gates keep confronting the spec, the contract and the code across them
// (DESIGN-dialect-per-layer.md). `dialects:` declares each lexicon once, by name, and each
// layer names the one its files are written in. A file's dialect is its layer's: the layer
// already says what a file is, and it says what it is written in — a second way of naming
// files, by glob, would drift from the layers.
//
// `dialect:` keeps its two roles. Alone, it is the project's one dialect, every layer's.
// Beside `dialects:`, it holds only what is about no language — the Gherkin language and the
// naming conventions —, and every lexicon lives in `dialects:`.

// projectFields are the fields of `dialect:` about no language: every feature is written in
// one Gherkin language, and a name promises the whole set in the same words whatever the
// language of its file (DLPLY-D03).
var projectFields = map[string]bool{"gherkin_language": true, "set_promise": true, "set_slice": true}

// LayerOf is the layer and the kind of a path, by the project's patterns: the layer whose
// pattern matches it and none of its `exclude`, the declared `priority` first, then the
// longest pattern, then the name.
func (c *Config) LayerOf(rel string) (layer, kind string) {
	if c == nil {
		return "", ""
	}
	// Normalized here, not by each caller: a path in the native Windows form matches no
	// pattern with a folder, and the file is left with no layer — out of the map, never
	// confronted, with no line of error.
	rel = strings.ReplaceAll(rel, `\`, "/")
	best := layerRank{prio: math.MinInt, size: -1}
	for name, l := range c.Layers {
		if !matchesGlob(l.Pattern, rel) || matchesAny(l.Exclude, rel) {
			continue
		}
		r := layerRank{prio: l.Priority, size: len(l.Pattern), name: name}
		if layer == "" || r.beats(best) {
			best, layer, kind = r, name, l.Kind
		}
	}
	return layer, kind
}

// layerRank is the tie-break between two layers matching a file, in order: declared
// priority, pattern length, name.
type layerRank struct {
	prio int
	size int
	name string
}

func (a layerRank) beats(b layerRank) bool {
	if a.prio != b.prio {
		return a.prio > b.prio
	}
	if a.size != b.size {
		return a.size > b.size
	}
	return a.name < b.name // a map has no order: the name keeps it stable
}

func matchesGlob(pattern, rel string) bool {
	ok, err := doublestar.Match(pattern, rel)
	return err == nil && ok
}

func matchesAny(globs []string, rel string) bool {
	for _, g := range globs {
		if matchesGlob(g, rel) {
			return true
		}
	}
	return false
}

// DialectAt is the dialect a file is read with: its layer's lexicon, with the project's
// Gherkin language and naming conventions. A project with one dialect gives it to every file.
// A file whose layer names no dialect — a spec, a feature, a doc, or a file of no layer —
// gets the project's fields alone: no lexicon of another layer reads it.
func (c *Config) DialectAt(rel string) Dialect {
	if c == nil || len(c.Dialects) == 0 {
		return c.DialectFor()
	}
	layer, _ := c.LayerOf(rel)
	return c.dialectOfLayer(layer)
}

// DialectOfLayer is DialectAt for the files of a layer.
func (c *Config) DialectOfLayer(layer string) Dialect {
	if c == nil || len(c.Dialects) == 0 {
		return c.DialectFor()
	}
	return c.dialectOfLayer(layer)
}

func (c *Config) dialectOfLayer(layer string) Dialect {
	name := c.Layers[layer].Dialect
	if name == "" && len(c.Dialects) == 1 {
		name = c.DialectNames()[0]
	}
	return c.lexicon(name)
}

// lexicon is the dialect of that name beside `dialects:`, with the project's fields; a name
// `dialects:` does not declare — a layer that names none — has the project's fields alone.
func (c *Config) lexicon(name string) Dialect {
	p := c.projectDialect()
	lex := c.Dialects[name]
	if lex == nil {
		return p
	}
	d := *lex
	d.GherkinLanguage, d.SetPromise, d.SetSlice = p.GherkinLanguage, p.SetPromise, p.SetSlice
	return d
}

// projectDialect is what `dialect:` holds beside `dialects:`: the fields about no language.
func (c *Config) projectDialect() Dialect {
	var d Dialect
	if c != nil && c.Dialect != nil {
		d.GherkinLanguage, d.SetPromise, d.SetSlice = c.Dialect.GherkinLanguage, c.Dialect.SetPromise, c.Dialect.SetSlice
	}
	return d
}

// DialectNames are the dialects the project declares, sorted: the names of `dialects:`, or
// "" alone for the project's one `dialect:`.
func (c *Config) DialectNames() []string {
	if c == nil || len(c.Dialects) == 0 {
		return []string{""}
	}
	names := make([]string, 0, len(c.Dialects))
	for n := range c.Dialects {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Lexicon is the dialect of that name, with the project's fields; "" is the project's one
// `dialect:`.
func (c *Config) Lexicon(name string) Dialect {
	if name == "" || c == nil || len(c.Dialects) == 0 {
		return c.DialectFor()
	}
	return c.lexicon(name)
}

// dialectField is the path a presupposed field is declared at in each dialect: the project's
// `dialect.<field>`, or `dialects.<name>.<field>`.
func dialectField(name, field string) string {
	if name == "" {
		return "dialect." + field
	}
	return "dialects." + name + "." + field
}

// presupposedInDialects splits a presupposition of a dialect's field (`dialect.<field>`) by
// dialect, beside `dialects:`: each dialect declares it, waives it in its own `opt_out`, or
// misses it. Every dialect waiving it means it does not apply; one missing it, that the gate
// is pending for its files, named at that dialect.
func (c *Config) presupposedInDialects(p string) (applies bool, missing []string, waived bool) {
	field := strings.TrimPrefix(p, "dialect.")
	if projectFields[strings.SplitN(field, ".", 2)[0]] {
		return c.Declares(p), nil, false
	}
	nwaived := 0
	for _, name := range c.DialectNames() {
		switch {
		case c.Declares(dialectField(name, field)):
		case c.Dialects[name].WaivedField(field[strings.LastIndex(field, ".")+1:]):
			nwaived++
		default:
			missing = append(missing, dialectField(name, field))
		}
	}
	return len(missing) == 0 && nwaived < len(c.Dialects), missing, nwaived == len(c.Dialects)
}

// validarDialetos refuses a `dialects:` the layers cannot read by: a lexicon left in
// `dialect:` beside it, a dialect declaring a project field, a layer naming a dialect nobody
// declared — or naming one with no `dialects:` —, and, with several declared, a layer of
// code or tests naming none (DLPLY-D01): a default would read it with another language's
// lexicon and pass it green.
func (c *Config) validarDialetos() error {
	if len(c.Dialects) == 0 {
		for _, l := range sortedLayerNames(c) {
			if c.Layers[l].Dialect != "" {
				return fmt.Errorf("layers.%s.dialect: names the dialect %q, and the project declares no `dialects:`", l, c.Layers[l].Dialect)
			}
		}
		return nil
	}
	if c.Dialect != nil {
		for _, f := range declaredFields(c.Dialect) {
			if !projectFields[f] {
				return fmt.Errorf("dialect.%s: beside `dialects:`, `dialect:` holds only gherkin_language, set_promise and set_slice; a lexicon goes in a dialect of `dialects:`", f)
			}
		}
	}
	for _, name := range c.DialectNames() {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("dialects: a dialect needs a name")
		}
		for _, f := range declaredFields(c.Dialects[name]) {
			if projectFields[f] {
				return fmt.Errorf("dialects.%s.%s: the project's, not a dialect's; declare it in `dialect:`", name, f)
			}
		}
	}
	for _, l := range sortedLayerNames(c) {
		layer := c.Layers[l]
		if layer.Dialect != "" {
			if _, ok := c.Dialects[layer.Dialect]; !ok {
				return fmt.Errorf("layers.%s.dialect: names %q, which `dialects:` does not declare (it declares %s)", l, layer.Dialect, strings.Join(c.DialectNames(), ", "))
			}
			continue
		}
		if len(c.Dialects) > 1 && (layer.Kind == "code" || layer.Kind == "test") {
			return fmt.Errorf("layers.%s: a layer of %s names the dialect its files are written in when the project declares several (%s)", l, layer.Kind, strings.Join(c.DialectNames(), ", "))
		}
	}
	return nil
}

func sortedLayerNames(c *Config) []string {
	names := make([]string, 0, len(c.Layers))
	for n := range c.Layers {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// declaredFields are the top-level fields a dialect declares, by their configuration names.
func declaredFields(d *Dialect) []string {
	if d == nil {
		return nil
	}
	var out []string
	for f, v := range fieldsByName(d) {
		if v {
			out = append(out, f)
		}
	}
	sort.Strings(out)
	return out
}

// fieldsByName are the top-level fields of a dialect by their configuration names, each
// saying whether it is declared: what its YAML writes, every field being `omitempty`.
func fieldsByName(d *Dialect) map[string]bool {
	out := map[string]bool{}
	b, err := yaml.Marshal(d)
	if err != nil {
		return out
	}
	var m map[string]any
	if yaml.Unmarshal(b, &m) != nil {
		return out
	}
	for k := range m {
		out[k] = true
	}
	return out
}

// DialectPathAt is where the dialect of a file is declared, for a message naming its field:
// `dialect`, or `dialects.<name>` for the dialect of the file's layer.
func (c *Config) DialectPathAt(rel string) string {
	if c == nil || len(c.Dialects) == 0 {
		return "dialect"
	}
	layer, _ := c.LayerOf(rel)
	name := c.Layers[layer].Dialect
	if name == "" && len(c.Dialects) == 1 {
		name = c.DialectNames()[0]
	}
	if name == "" {
		return "dialect"
	}
	return "dialects." + name
}

// FamilyOfFile is the family whose example reads the file's language, by its comments'
// example — `go` for a `.go` file —; "" when no family's does. The configuration and scripts
// every family comments alike tell no language.
func FamilyOfFile(rel string) string {
	rel = strings.TrimPrefix(strings.ReplaceAll(rel, `\`, "/"), "./")
	shared := map[string]int{}
	for _, fam := range KnownDialectFamilies() {
		for _, c := range FamilyComments(fam) {
			shared[c.Files]++
		}
	}
	for _, fam := range KnownDialectFamilies() {
		for _, c := range FamilyComments(fam) {
			if shared[c.Files] == 1 && c.matches(rel) {
				return fam
			}
		}
	}
	return ""
}
