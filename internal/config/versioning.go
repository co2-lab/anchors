// @anchors
//   code: VRSND
//   ref: VRSNN

package config

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/version"
	"gopkg.in/yaml.v3"
)

// Versioning maps a name to a product: one versioned thing with its own release line. A
// repository with a single product still names it, so a second one is not a breaking change.
type Versioning map[string]*Product

// Product is how one product writes, orders and moves its version.
type Product struct {
	// Scheme orders the versions: semver, calver, number or pattern.
	Scheme string `yaml:"scheme"`
	// Format is calver's fields (`YYYY.0M.MICRO`).
	Format string `yaml:"format,omitempty"`
	// Pattern is the regular expression whose groups are a pattern version's numbers.
	Pattern string `yaml:"pattern,omitempty"`
	// Source are the files that hold the version; two or more must agree. None: the version
	// is the tag alone.
	Source []VersionSource `yaml:"source,omitempty"`
	// Tag is the release record's name, with `{version}`. Empty: `v{version}`.
	Tag string `yaml:"tag,omitempty"`
	// Next is the default strategy of every branch. Empty: manual.
	Next *Next `yaml:"next,omitempty"`
	// Branches are the branches whose rule differs from the default, first match wins.
	Branches BranchRules `yaml:"branches,omitempty"`
}

// VersionSource is one file that holds the version: at a JSON or YAML path, or in the one
// group of a regular expression over any text.
type VersionSource struct {
	File  string `yaml:"file" path:"write"`
	Path  string `yaml:"path,omitempty"`
	Regex string `yaml:"regex,omitempty"`
}

// Next is how a branch decides its next version: the strategy and its options.
type Next struct {
	Strategy string
	Part     string
	Initial  bool
	Types    map[string]string
}

// UnmarshalYAML reads a strategy written as its name (`next: none`) or as a map
// (`next: { strategy: increment, part: patch }`).
func (n *Next) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		n.Strategy = node.Value
		return nil
	}
	var m struct {
		Strategy string            `yaml:"strategy"`
		Part     string            `yaml:"part"`
		Initial  bool              `yaml:"initial"`
		Types    map[string]string `yaml:"types"`
	}
	if err := strictKeys(node, "next", "strategy", "part", "initial", "types"); err != nil {
		return err
	}
	if err := node.Decode(&m); err != nil {
		return err
	}
	*n = Next{Strategy: m.Strategy, Part: m.Part, Initial: m.Initial, Types: m.Types}
	return nil
}

// Resolved is the strategy the version package computes with; manual when none is written.
func (n *Next) Resolved() version.Strategy {
	if n == nil {
		return version.Strategy{Name: version.Manual}
	}
	name := n.Strategy
	if name == "" {
		name = version.Manual
	}
	return version.Strategy{Name: name, Part: n.Part, Initial: n.Initial, Types: n.Types}
}

// BranchRule is the rule of the branches a pattern names.
type BranchRule struct {
	Pattern string `yaml:"-"` // read by BranchRules
	Next    *Next  `yaml:"-"`
}

// BranchRules keep the order they are written in: the first pattern a branch matches wins.
type BranchRules []BranchRule

// UnmarshalYAML reads `branches:` as a map, in the order it is written.
func (b *BranchRules) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("line %d: branches: a map of branch pattern → its rule", node.Line)
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		k, v := node.Content[i], node.Content[i+1]
		if err := strictKeys(v, "branches."+k.Value, "next"); err != nil {
			return err
		}
		var r struct {
			Next *Next `yaml:"next"`
		}
		if err := v.Decode(&r); err != nil {
			return err
		}
		*b = append(*b, BranchRule{Pattern: k.Value, Next: r.Next})
	}
	return nil
}

// strictKeys refuses a key the block does not know: a nested decode does not inherit the
// strictness of the file's.
func strictKeys(node *yaml.Node, field string, known ...string) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("line %d: %s: a map is expected", node.Line, field)
	}
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		ok := false
		for _, w := range known {
			ok = ok || k.Value == w
		}
		if !ok {
			return fmt.Errorf("line %d: %s: field %s not found (known: %s)", k.Line, field, k.Value, strings.Join(known, ", "))
		}
	}
	return nil
}

// NextFor is the strategy a branch resolves to: the first rule whose pattern it matches, or
// the product's default; and the pattern that decided, empty for the default.
func (p *Product) NextFor(branch string) (version.Strategy, string) {
	for _, r := range p.Branches {
		if BranchMatches(r.Pattern, branch) && r.Next != nil {
			return r.Next.Resolved(), r.Pattern
		}
	}
	return p.Next.Resolved(), ""
}

// TagFor is the tag that records a version of the product.
func (p *Product) TagFor(v string) string {
	return strings.ReplaceAll(p.TagTemplate(), "{version}", v)
}

// TagTemplate is the product's tag, `v{version}` when it declares none.
func (p *Product) TagTemplate() string {
	if p.Tag == "" {
		return "v{version}"
	}
	return p.Tag
}

// VersionOfTag is the version a tag records for the product, and whether it is one of its.
func (p *Product) VersionOfTag(tag string) (string, bool) {
	pre, post, _ := strings.Cut(p.TagTemplate(), "{version}")
	if !strings.HasPrefix(tag, pre) || !strings.HasSuffix(tag, post) || len(tag) <= len(pre)+len(post) {
		return "", false
	}
	return tag[len(pre) : len(tag)-len(post)], true
}

// SchemeOf builds the product's scheme.
func (p *Product) SchemeOf() (version.Scheme, error) {
	return version.NewScheme(p.Scheme, p.Format, p.Pattern)
}

// validateVersioning refuses a product the Anchors could not version: a scheme that reads
// nothing, a source that names no place, a tag with no version in it, a strategy that does
// not exist or does not fit the scheme, a branch pattern that matches nothing.
func validateVersioning(v Versioning) error {
	names := make([]string, 0, len(v))
	for n := range v {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := validateProduct(v[name]); err != nil {
			return fmt.Errorf("versioning.%s: %w", name, err)
		}
	}
	return nil
}

func validateProduct(p *Product) error {
	if p == nil {
		return fmt.Errorf("a product declares at least its scheme")
	}
	s, err := p.SchemeOf()
	if err != nil {
		return err
	}
	for i, src := range p.Source {
		if src.File == "" {
			return fmt.Errorf("source[%d]: file is required", i)
		}
		if (src.Path == "") == (src.Regex == "") {
			return fmt.Errorf("source[%d] (%s): write path: or regex:, one of them", i, src.File)
		}
		if src.Regex != "" {
			re, err := regexp.Compile(src.Regex)
			if err != nil {
				return fmt.Errorf("source[%d] (%s): regex %q: %v", i, src.File, src.Regex, err)
			}
			if re.NumSubexp() != 1 {
				return fmt.Errorf("source[%d] (%s): regex %q has %d groups: one group holds the version", i, src.File, src.Regex, re.NumSubexp())
			}
		}
	}
	if !strings.Contains(p.TagTemplate(), "{version}") || strings.Count(p.TagTemplate(), "{version}") > 1 {
		return fmt.Errorf("tag %q: write {version} once, where the version goes", p.Tag)
	}
	if err := validateNext(p.Next, s, "next", len(p.Source) > 0); err != nil {
		return err
	}
	for _, r := range p.Branches {
		if err := validBranchPattern(r.Pattern); err != nil {
			return fmt.Errorf("branches: %w", err)
		}
		if err := validateNext(r.Next, s, "branches."+r.Pattern+".next", len(p.Source) > 0); err != nil {
			return err
		}
	}
	return nil
}

func validateNext(n *Next, s version.Scheme, field string, hasSource bool) error {
	st := n.Resolved()
	switch st.Name {
	case version.Manual:
		if !hasSource {
			return fmt.Errorf("%s: manual reads the version from the sources, and the product declares none", field)
		}
	case version.None:
	case version.Increment:
		if !s.ValidPart(st.Part) {
			return fmt.Errorf("%s: %s has no part %q (%s)", field, s.Kind, st.Part, strings.Join(s.Parts(), ", "))
		}
	case version.Semantic:
		if s.Kind != version.SemVer {
			return fmt.Errorf("%s: semantic needs scheme semver, not %s", field, s.Kind)
		}
		for t, l := range st.Types {
			if !version.ValidLevel(l) {
				return fmt.Errorf("%s.types.%s: %q is not a level (major, minor, patch, none)", field, t, l)
			}
		}
	default:
		return fmt.Errorf("%s: strategy %q: write manual, increment, semantic or none", field, st.Name)
	}
	if st.Name != version.Increment && st.Part != "" {
		return fmt.Errorf("%s: part is an option of increment, not of %s", field, st.Name)
	}
	if st.Name != version.Semantic && (st.Initial || len(st.Types) > 0) {
		return fmt.Errorf("%s: initial and types are options of semantic, not of %s", field, st.Name)
	}
	return nil
}

// SourceVersion is the version one source holds, and where.
type SourceVersion struct {
	File, Version string
	start, end    int // the version's bytes in the file
}

// ReadSources reads the version every source holds. A source that cannot be read, or that
// holds no version where it says, is an error naming the file: a version not found is never
// "no version".
func (p *Product) ReadSources(read func(file string) ([]byte, error)) ([]SourceVersion, error) {
	var out []SourceVersion
	for _, src := range p.Source {
		b, err := read(src.File)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", src.File, err)
		}
		sv, err := locate(src, b)
		if err != nil {
			return nil, err
		}
		out = append(out, sv)
	}
	return out, nil
}

// Agreed is the one version the sources hold, or an error naming each when they differ.
func Agreed(svs []SourceVersion) (string, error) {
	if len(svs) == 0 {
		return "", nil
	}
	for _, s := range svs[1:] {
		if s.Version != svs[0].Version {
			var parts []string
			for _, s := range svs {
				parts = append(parts, s.File+" "+s.Version)
			}
			return "", fmt.Errorf("the sources disagree: %s", strings.Join(parts, ", "))
		}
	}
	return svs[0].Version, nil
}

// Rewrite is the source's content with the version replaced, and nothing else of the file
// touched: the bytes before and after the version are kept as they are.
func (p *Product) Rewrite(file string, content []byte, v string) ([]byte, error) {
	for _, src := range p.Source {
		if src.File != file {
			continue
		}
		sv, err := locate(src, content)
		if err != nil {
			return nil, err
		}
		out := append([]byte(nil), content[:sv.start]...)
		out = append(out, v...)
		return append(out, content[sv.end:]...), nil
	}
	return nil, fmt.Errorf("%s is not a source of the product", file)
}

// locate finds the version in a source: by its regex group, or at its path in a JSON or YAML
// document — yaml reads both, and gives the scalar's place, so the file is rewritten in place.
func locate(src VersionSource, b []byte) (SourceVersion, error) {
	sv := SourceVersion{File: src.File}
	if src.Regex != "" {
		m := regexp.MustCompile(src.Regex).FindSubmatchIndex(b)
		if m == nil || m[2] < 0 {
			return sv, fmt.Errorf("%s: regex %q matches nothing", src.File, src.Regex)
		}
		sv.start, sv.end, sv.Version = m[2], m[3], string(b[m[2]:m[3]])
		return sv, nil
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return sv, fmt.Errorf("%s: not JSON nor YAML: %v", src.File, err)
	}
	n := &doc
	if n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		n = n.Content[0]
	}
	for _, key := range strings.Split(src.Path, ".") {
		next := (*yaml.Node)(nil)
		switch n.Kind {
		case yaml.MappingNode:
			for i := 0; i+1 < len(n.Content); i += 2 {
				if n.Content[i].Value == key {
					next = n.Content[i+1]
				}
			}
		case yaml.SequenceNode:
			var idx int
			if _, err := fmt.Sscanf(key, "%d", &idx); err == nil && idx >= 0 && idx < len(n.Content) {
				next = n.Content[idx]
			}
		}
		if next == nil {
			return sv, fmt.Errorf("%s: no %s in the document", src.File, src.Path)
		}
		n = next
	}
	if n.Kind != yaml.ScalarNode || n.Value == "" {
		return sv, fmt.Errorf("%s: %s holds no version", src.File, src.Path)
	}
	start := offsetOf(b, n.Line, n.Column)
	if n.Style == yaml.DoubleQuotedStyle || n.Style == yaml.SingleQuotedStyle {
		start++
	}
	end := start + len(n.Value)
	if start < 0 || end > len(b) || string(b[start:end]) != n.Value {
		return sv, fmt.Errorf("%s: %s is written in a form the Anchors cannot rewrite in place; use regex:", src.File, src.Path)
	}
	sv.start, sv.end, sv.Version = start, end, n.Value
	return sv, nil
}

// offsetOf is the byte offset of a 1-based line and column.
func offsetOf(b []byte, line, col int) int {
	off := 0
	for l := 1; l < line; l++ {
		i := bytes.IndexByte(b[off:], '\n')
		if i < 0 {
			return -1
		}
		off += i + 1
	}
	// yaml counts columns in runes
	rs := []rune(string(b[off:]))
	if col-1 > len(rs) {
		return -1
	}
	return off + len(string(rs[:col-1]))
}
