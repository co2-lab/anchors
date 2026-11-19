// @anchors
//   code: VRSND
//   ref: VRSNN

package config

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
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
	// Released lists more released versions besides the tags: the releases made before the
	// project wrote tags.
	Released *Listing `yaml:"released,omitempty"`
	// Branches are the branches whose rule differs from the default, first match wins.
	Branches BranchRules `yaml:"branches,omitempty"`
	// Changelog are the files the product's changelog is written to.
	Changelog *ProductChangelog `yaml:"changelog,omitempty"`
}

// ProductChangelog are the outputs a product's changelog is written to. None: the project's
// `changelog:` block, or its defaults.
type ProductChangelog struct {
	Outputs []ChangelogOutput `yaml:"outputs,omitempty"`
}

// ChangelogOutput is one file (or directory) the changelog is written to: its path, mode and
// template, as the `changelog:` block has them, and the format `changelog-valid` holds it to.
type ChangelogOutput struct {
	Changelog `yaml:",inline"`
	// Source is what the releases are read from: `commits` (the default) or `intents`.
	Source string           `yaml:"source,omitempty"`
	Format *ChangelogFormat `yaml:"format,omitempty"`
}

// FromIntents says whether the output is written from the change intents.
func (o ChangelogOutput) FromIntents() bool { return o.Source == "intents" }

// Intents is where the change intents live and what they may say.
type Intents struct {
	// Path is the directory of the intent files (`*.md`). Empty: .anchors/changes.
	Path string `yaml:"path,omitempty" path:"write"`
	// Types are the types an intent may declare. Empty: feat and fix.
	Types []string `yaml:"types,omitempty"`
	// MaxChars bounds an intent's text.
	MaxChars int `yaml:"max_chars,omitempty"`
	// Consume is what a release does with its intents: `delete` them (the default), or
	// `archive` them under released/<version>/, where the changelog reads them again.
	Consume string `yaml:"consume,omitempty"`
}

// PathOrDefault is the intents' directory.
func (i *Intents) PathOrDefault() string {
	if i == nil || i.Path == "" {
		return ".anchors/changes"
	}
	return i.Path
}

// TypesOrDefault are the types an intent may declare.
func (i *Intents) TypesOrDefault() []string {
	if i == nil || len(i.Types) == 0 {
		return []string{"feat", "fix"}
	}
	return i.Types
}

// ChangelogFormat is what an output must look like, beyond the markers, the order and the
// agreement with the tags, which hold for every output.
type ChangelogFormat struct {
	// Heading is the pattern of a block's first line, a regular expression with {version}
	// and {date} in it.
	Heading string `yaml:"heading,omitempty"`
	// File is the name of a per_version file, with {version}. Empty: {version}.md.
	File string `yaml:"file,omitempty"`
	// MaxChars bounds each release block.
	MaxChars int `yaml:"max_chars,omitempty"`
	// RequiredSections are the `### ` sections every release block has.
	RequiredSections []string `yaml:"required_sections,omitempty"`
	// Required: changed asks the output to have changed since the latest release tag.
	Required string `yaml:"required,omitempty"`
}

// Listing selects versions from a file: the keys of a map, kept when their entry matches
// every field of `where`. It interprets nothing: the file, the map and the fields are the
// project's.
type Listing struct {
	// File is the file the map is read from, JSON or YAML.
	File string `yaml:"file,omitempty" path:"read"`
	// Keys is the dot path of the map whose keys are versions; empty: the document itself.
	Keys string `yaml:"keys,omitempty"`
	// Where are the fields an entry must hold, by dot path, compared as text.
	Where map[string]string `yaml:"where,omitempty"`
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
	// Range bounds the versions the branch may produce: `1.0.x`, or `from_base`, the
	// release the branch was cut from at its last part.
	Range string `yaml:"-"`
	// Changelog marks a branch the changelog is written on.
	Changelog bool `yaml:"-"`
	// Release marks a branch that releases: its changelog lists the releases, one per tag.
	Release bool `yaml:"-"`
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
		if err := strictKeys(v, "branches."+k.Value, "next", "range", "changelog", "release"); err != nil {
			return err
		}
		var r struct {
			Next      *Next  `yaml:"next"`
			Range     string `yaml:"range"`
			Changelog bool   `yaml:"changelog"`
			Release   bool   `yaml:"release"`
		}
		if err := v.Decode(&r); err != nil {
			return err
		}
		*b = append(*b, BranchRule{Pattern: k.Value, Next: r.Next, Range: r.Range, Changelog: r.Changelog, Release: r.Release})
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

// ruleFor is the first rule whose pattern the branch matches; nil when none does.
func (p *Product) ruleFor(branch string) *BranchRule {
	for i, r := range p.Branches {
		if BranchMatches(r.Pattern, branch) {
			return &p.Branches[i]
		}
	}
	return nil
}

// NextFor is the strategy a branch resolves to: the first rule whose pattern it matches, or
// the product's default when that rule sets none or no rule matches; and the pattern that
// decided, empty for the default.
func (p *Product) NextFor(branch string) (version.Strategy, string) {
	if r := p.ruleFor(branch); r != nil && r.Next != nil {
		return r.Next.Resolved(), r.Pattern
	}
	return p.Next.Resolved(), ""
}

// RangeFor is the range of the first rule the branch matches, and its pattern; empty when
// the branch is not bounded.
func (p *Product) RangeFor(branch string) (string, string) {
	if r := p.ruleFor(branch); r != nil && r.Range != "" {
		return r.Range, r.Pattern
	}
	return "", ""
}

// Releases says whether a branch releases the product: its rule says `release: true`, or no
// rule of the product names where it releases, and then every branch does.
func (p *Product) Releases(branch string) bool {
	named := false
	for _, r := range p.Branches {
		named = named || r.Release
	}
	if !named {
		return true
	}
	r := p.ruleFor(branch)
	return r != nil && r.Release
}

// WritesChangelog says whether the changelog is written on a branch: its rule says
// `changelog: true`, or no rule names one, and then every branch does.
func (p *Product) WritesChangelog(branch string) bool {
	named := false
	for _, r := range p.Branches {
		named = named || r.Changelog
	}
	if !named {
		return true
	}
	r := p.ruleFor(branch)
	return r != nil && r.Changelog
}

// ReleaseBranch is the pattern of the first rule that releases, the branch a changelog reads
// when none is named; empty when no rule names one.
func (p *Product) ReleaseBranch() string {
	for _, r := range p.Branches {
		if r.Release {
			return r.Pattern
		}
	}
	return ""
}

// FromBase is the range that bounds a branch to the release it was cut from.
const FromBase = "from_base"

// InRange says whether a version is inside a range: an explicit one (`1.0.x`) keeps the
// parts written before the `x`; `from_base` keeps every part of the base release but the
// last. The reason names the bound when it does not hold.
func InRange(s version.Scheme, rng, base string, v version.V) (bool, string, error) {
	var prefix []int
	bound := rng
	if rng == FromBase {
		if base == "" {
			return false, "the branch was cut from no release (range: from_base)", nil
		}
		b, err := s.Parse(base)
		if err != nil {
			return false, "", fmt.Errorf("the base release: %w", err)
		}
		prefix = b.Nums[:len(b.Nums)-1]
		bound = rangeOf(prefix) + " (from_base: " + base + ")"
	} else {
		var err error
		if prefix, err = parseRange(rng); err != nil {
			return false, "", err
		}
	}
	if len(prefix) > len(v.Nums) {
		return false, v.Raw + " is outside " + bound, nil
	}
	for i, n := range prefix {
		if v.Nums[i] != n {
			return false, v.Raw + " is outside " + bound, nil
		}
	}
	return true, v.Raw + " is inside " + bound, nil
}

// parseRange reads `1.0.x`: the numbers before the `x`, at least one.
func parseRange(rng string) ([]int, error) {
	parts := strings.Split(rng, ".")
	if len(parts) < 2 || parts[len(parts)-1] != "x" {
		return nil, fmt.Errorf("range %q: write the parts kept and an x for the rest (1.0.x), or from_base", rng)
	}
	out := make([]int, 0, len(parts)-1)
	for _, p := range parts[:len(parts)-1] {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("range %q: %q is not a number", rng, p)
		}
		out = append(out, n)
	}
	return out, nil
}

func rangeOf(prefix []int) string {
	var b strings.Builder
	for _, n := range prefix {
		b.WriteString(strconv.Itoa(n) + ".")
	}
	return b.String() + "x"
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

// Pick is the product named, or the only one when no name is given; `how` says where a name
// is written, for the error that asks for one.
func (v Versioning) Pick(name, how string) (string, *Product, error) {
	names := make([]string, 0, len(v))
	for n := range v {
		names = append(names, n)
	}
	sort.Strings(names)
	switch {
	case len(v) == 0:
		return "", nil, fmt.Errorf("anchors.yaml declares no `versioning:` block: nothing is versioned")
	case name != "":
		if p, ok := v[name]; ok {
			return name, p, nil
		}
		return "", nil, fmt.Errorf("versioning has no product %q (%s)", name, strings.Join(names, ", "))
	case len(v) > 1:
		return "", nil, fmt.Errorf("versioning has %d products — name one with %s (%s)", len(v), how, strings.Join(names, ", "))
	}
	return names[0], v[names[0]], nil
}

// ListingOf is a requirement's listing with the file and the map the product's `released:`
// declares, where the requirement names none.
func (p *Product) ListingOf(l Listing) Listing {
	if l.File == "" && p.Released != nil {
		l.File, l.Keys = p.Released.File, p.Released.Keys
	}
	return l
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
		if r.Next != nil {
			if err := validateNext(r.Next, s, "branches."+r.Pattern+".next", len(p.Source) > 0); err != nil {
				return err
			}
		}
		if r.Range != "" && r.Range != FromBase {
			if _, err := parseRange(r.Range); err != nil {
				return fmt.Errorf("branches.%s: %w", r.Pattern, err)
			}
		}
	}
	if p.Released != nil && p.Released.File == "" {
		return fmt.Errorf("released: file is required")
	}
	if p.Changelog != nil {
		for i, o := range p.Changelog.Outputs {
			if err := validateOutput(o); err != nil {
				return fmt.Errorf("changelog.outputs[%d]: %w", i, err)
			}
		}
	}
	return nil
}

func validateOutput(o ChangelogOutput) error {
	if o.Path == "" {
		return fmt.Errorf("path is required")
	}
	if o.Intents != nil {
		return fmt.Errorf("intents are the project's (changelog.intents), not an output's")
	}
	if o.Source != "" && o.Source != "commits" && o.Source != "intents" {
		return fmt.Errorf("source %q: write commits or intents", o.Source)
	}
	if o.Mode != "" && o.Mode != "incremental" && o.Mode != "per_version" {
		return fmt.Errorf("mode %q: write incremental or per_version", o.Mode)
	}
	f := o.Format
	if f == nil {
		return nil
	}
	if f.Heading != "" {
		if !strings.Contains(f.Heading, "{version}") {
			return fmt.Errorf("format.heading %q: write {version} where the version goes", f.Heading)
		}
		if _, err := HeadingRE(f.Heading); err != nil {
			return fmt.Errorf("format.heading %q: %v", f.Heading, err)
		}
	}
	if f.File != "" && !strings.Contains(f.File, "{version}") {
		return fmt.Errorf("format.file %q: write {version} where the version goes", f.File)
	}
	if f.File != "" && o.ModeOrDefault() != "per_version" {
		return fmt.Errorf("format.file is for per_version, and the mode is %s", o.ModeOrDefault())
	}
	if f.Required != "" && f.Required != "changed" {
		return fmt.Errorf("format.required %q: the one requirement is changed", f.Required)
	}
	if f.MaxChars < 0 {
		return fmt.Errorf("format.max_chars %d: a limit is positive", f.MaxChars)
	}
	return nil
}

// HeadingRE compiles a heading pattern: {version} and {date} become groups.
func HeadingRE(pattern string) (*regexp.Regexp, error) {
	p := strings.ReplaceAll(pattern, "{version}", `(?P<version>\S+?)`)
	p = strings.ReplaceAll(p, "{date}", `(?P<date>\d{4}-\d{2}-\d{2})`)
	return regexp.Compile(p)
}

func validateNext(n *Next, s version.Scheme, field string, hasSource bool) error {
	st := n.Resolved()
	switch st.Name {
	case version.Manual:
		if !hasSource {
			return fmt.Errorf("%s: manual reads the version from the sources, and the product declares none", field)
		}
	case version.None, version.Intents:
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
		return fmt.Errorf("%s: strategy %q: write manual, increment, semantic, intents or none", field, st.Name)
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
	if n = walk(n, src.Path); n == nil {
		return sv, fmt.Errorf("%s: no %s in the document", src.File, src.Path)
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

// walk follows a dot path from a node: a key of a map, a number for a list item; nil when
// the document has nothing there. An empty path is the node itself.
func walk(n *yaml.Node, path string) *yaml.Node {
	if path == "" {
		return n
	}
	for _, key := range strings.Split(path, ".") {
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
			return nil
		}
		n = next
	}
	return n
}

// Versions are the keys of the listing's map whose entry holds every field of `where`. A
// file that cannot be read, or with no map at the path, is an error: a listing not found is
// never an empty one.
func (l Listing) Versions(read func(file string) ([]byte, error)) ([]string, error) {
	b, err := read(l.File)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", l.File, err)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("%s: not JSON nor YAML: %v", l.File, err)
	}
	n := &doc
	if n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		n = n.Content[0]
	}
	if n = walk(n, l.Keys); n == nil || n.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s: no map at %s", l.File, or(l.Keys, "the top of the document"))
	}
	var out []string
	for i := 0; i+1 < len(n.Content); i += 2 {
		entry, ok := n.Content[i+1], true
		for field, want := range l.Where {
			got := walk(entry, field)
			ok = ok && got != nil && got.Kind == yaml.ScalarNode && got.Value == want
		}
		if ok {
			out = append(out, n.Content[i].Value)
		}
	}
	return out, nil
}

// String says what the listing selects: `release-policy.json versions, where target = review`.
func (l Listing) String() string {
	s := l.File
	if l.Keys != "" {
		s += " " + l.Keys
	}
	if len(l.Where) > 0 {
		fields := make([]string, 0, len(l.Where))
		for f, v := range l.Where {
			fields = append(fields, f+" = "+v)
		}
		sort.Strings(fields)
		s += ", where " + strings.Join(fields, " and ")
	}
	return s
}

func or(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
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
