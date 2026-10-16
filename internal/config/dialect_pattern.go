// @anchors
//   code: DLPTD
//   ref: DLCTI

package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

// Pattern is how a dialect says where something stands in a file, in the project's own
// syntax: the core reads none. It is a regular expression, run over the whole content —
// a statement may span lines —, whose named groups say what it captures. A list of things
// cannot come from one group (a regular expression captures once per match), so a pattern
// may instead be a region and what repeats in it:
//
//   - 'import\s+(?P<default>\w+)\s+from'               # one capture per match
//   - within: 'import\s*\{(?P<within>[^}]*)\}\s*from'   # the region, in its `within` group
//     each: '(?P<name>\w+)\s*(?:,|$)'                   # run over the region, repeatedly
//
// Files narrows the files it applies to (a glob, as the layers write theirs); Relative says
// its paths are relative to the importer even without `./` — a flow names the flow it runs;
// Runs says a statement it reads plainly runs, unless the type pattern reads it as types —
// so a type flag on it is wrong.
type Pattern struct {
	Re       string `yaml:"re,omitempty"`
	Within   string `yaml:"within,omitempty"`
	Each     string `yaml:"each,omitempty"`
	Files    string `yaml:"files,omitempty"`
	Relative bool   `yaml:"relative,omitempty"`
	Runs     bool   `yaml:"runs,omitempty"`
}

// Patterns is a list of patterns, read in the order their matches stand in the source.
type Patterns []Pattern

// UnmarshalYAML takes a pattern as text — the regular expression — or as a map.
func (p *Pattern) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		return n.Decode(&p.Re)
	}
	type plain Pattern
	var v plain
	if err := n.Decode(&v); err != nil {
		return err
	}
	*p = Pattern(v)
	return nil
}

// MarshalYAML writes a pattern that is only a regular expression as text.
func (p Pattern) MarshalYAML() (any, error) {
	if p.Within == "" && p.Each == "" && p.Files == "" && !p.Relative && !p.Runs {
		return p.Re, nil
	}
	type plain Pattern
	return plain(p), nil
}

// Capture is what a group captured, and where: its offsets in the content.
type Capture struct {
	Value      string
	Start, End int
}

// Match is one match of a pattern: its span in the content, and the groups that took part.
type Match struct {
	Start, End int
	Groups     map[string]Capture
	// Region is where the statement it belongs to starts: the region's match for a pattern
	// with `within`, its own start otherwise. The matches of one list share it.
	Region int
	// Listed: it came from a region — one of a list.
	Listed bool
	// Relative: its pattern says its paths are relative to the importer.
	Relative bool
	// Runs: its pattern says the statement plainly runs.
	Runs bool
}

// Has says whether the group took part in the match.
func (m Match) Has(group string) bool {
	_, ok := m.Groups[group]
	return ok
}

var patternCache sync.Map // string → *regexp.Regexp, or nil when it does not compile

func compiled(expr string) *regexp.Regexp {
	if v, ok := patternCache.Load(expr); ok {
		re, _ := v.(*regexp.Regexp)
		return re
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		re = nil
	}
	patternCache.Store(expr, re)
	return re
}

// AppliesTo says whether the pattern reads this file (its path relative to the root, with
// `/`).
func (p Pattern) AppliesTo(rel string) bool {
	if strings.TrimSpace(p.Files) == "" {
		return true
	}
	ok, _ := doublestar.Match(p.Files, rel)
	return ok
}

// FindAll are the pattern's matches in the content. A pattern that does not compile finds
// nothing: the config refuses it on load (Validate).
func (p Pattern) FindAll(src string) []Match {
	if p.Within == "" {
		out := findAll(compiled(p.Re), src, 0, len(src))
		for i := range out {
			out[i].Region, out[i].Relative, out[i].Runs = out[i].Start, p.Relative, p.Runs
		}
		return out
	}
	within, each := compiled(p.Within), compiled(p.Each)
	if within == nil || each == nil {
		return nil
	}
	wi := within.SubexpIndex("within")
	if wi < 0 {
		return nil
	}
	var out []Match
	for _, m := range within.FindAllStringSubmatchIndex(src, -1) {
		if m[2*wi] < 0 {
			continue
		}
		for _, e := range findAll(each, src, m[2*wi], m[2*wi+1]) {
			e.Region, e.Listed, e.Relative, e.Runs = m[0], true, p.Relative, p.Runs
			out = append(out, e)
		}
	}
	return out
}

// findAll runs a regular expression over src[from:to], the offsets kept in src's terms.
func findAll(re *regexp.Regexp, src string, from, to int) []Match {
	if re == nil {
		return nil
	}
	region := src[from:to]
	names := re.SubexpNames()
	var out []Match
	for _, m := range re.FindAllStringSubmatchIndex(region, -1) {
		mt := Match{Start: from + m[0], End: from + m[1], Groups: map[string]Capture{}}
		for i, name := range names {
			if name == "" || m[2*i] < 0 {
				continue
			}
			mt.Groups[name] = Capture{Value: region[m[2*i]:m[2*i+1]], Start: from + m[2*i], End: from + m[2*i+1]}
		}
		out = append(out, mt)
	}
	return out
}

// FindAll are the matches of every pattern that reads this file, in the order they stand in
// the content.
func (ps Patterns) FindAll(src, rel string) []Match {
	var out []Match
	for _, p := range ps {
		if p.AppliesTo(rel) {
			out = append(out, p.FindAll(src)...)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Start < out[j].Start })
	return out
}

// capturing is the expression whose groups a pattern captures with: its own, or `each`.
func (p Pattern) capturing() string {
	if p.Within != "" || p.Each != "" {
		return p.Each
	}
	return p.Re
}

// Validate says what is wrong with a pattern of `field`: an expression that does not
// compile, a region with no `within` group, a capture with none of the groups its list
// reads, a glob that is no glob. A pattern with none of them would find nothing, and its gate
// would answer as if the code held nothing.
func (p Pattern) Validate(field string, groups []string) error {
	if (p.Within == "") != (p.Each == "") {
		return fmt.Errorf("`%s` declares one of `within` and `each`: a region needs both", field)
	}
	if p.Within != "" {
		if p.Re != "" {
			return fmt.Errorf("`%s` declares `re` and a region: declare one", field)
		}
		re, err := regexp.Compile(p.Within)
		if err != nil {
			return fmt.Errorf("`%s.within` does not compile: %w", field, err)
		}
		if re.SubexpIndex("within") < 0 {
			return fmt.Errorf("`%s.within` has no `within` group: name the region, `(?P<within>…)`", field)
		}
	}
	re, err := regexp.Compile(p.capturing())
	if err != nil {
		return fmt.Errorf("`%s` does not compile: %w", field, err)
	}
	if strings.TrimSpace(p.capturing()) == "" {
		return fmt.Errorf("`%s` is empty", field)
	}
	named := false
	for _, g := range groups {
		if re.SubexpIndex(g) >= 0 {
			named = true
		}
	}
	if !named {
		return fmt.Errorf("`%s` captures with none of the groups %s: name what it captures, e.g. `(?P<%s>…)`",
			field, strings.Join(groups, ", "), groups[0])
	}
	if p.Files != "" && !doublestar.ValidatePattern(p.Files) {
		return fmt.Errorf("`%s.files` is not a valid glob: %s", field, p.Files)
	}
	return nil
}

// Validate checks each pattern of the list (see Pattern.Validate).
func (ps Patterns) Validate(field string, groups []string) error {
	for i, p := range ps {
		if err := p.Validate(fmt.Sprintf("%s[%d]", field, i), groups); err != nil {
			return err
		}
	}
	return nil
}

// LineOf is the 1-based line of an offset in the content.
func LineOf(src string, off int) int {
	if off > len(src) {
		off = len(src)
	}
	return strings.Count(src[:off], "\n") + 1
}

// NamesBrought are the names an import statement brings, by the dialect's import_names: each
// `name`, and `default` for the module's default.
func (d Dialect) NamesBrought(statement, rel string) []string {
	var out []string
	for _, m := range d.ImportNames.FindAll(statement, rel) {
		switch {
		case m.Has("default"):
			out = append(out, "default")
		case m.Groups["name"].Value != "":
			out = append(out, m.Groups["name"].Value)
		}
	}
	return out
}
