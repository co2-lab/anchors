// @anchors
//   code: SYNTX
//   ref: MGSYN

package migra

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/co2-lab/anchors/internal/config"
)

// FORMAT 8 — what the core knew of a language, written into the project.
//
// The binary before format 8 read comments with a table of its own, filled every lexicon
// field a project left empty from its family, and knew the shapes of a route, a mock, a
// statement's end, a limit and a test's block by itself. Format 8 reads only what the
// project declares. Crossing it writes, once, the examples the old binary applied — only
// where the project declared nothing —, so its table does not change on the upgrade, and
// the yaml now says what is read. What the old binary applied and no example can say is
// named, to be declared by hand.

// Syntax is what crossing format 8 adds to a project's anchors.yaml.
type Syntax struct {
	// Dialect holds the dialect fields to add, each one the project left empty.
	Dialect *config.Dialect
	// Tests holds the fields to add inside a `dialect.tests` the project already declares.
	Tests *config.TestsSource
	// Derived holds the `derived` fields to add.
	Derived *DerivedAdd
	// Headless are the test layers of data files — a runner's flows — that carry no header.
	Headless []string
	// Runners are the example runners of the tools the project's suites run with, which the
	// old monitor knew by itself.
	Runners []config.MonitorRunner
	// Manual is what the old binary applied and no example can say: to declare by hand.
	Manual []string
}

// DerivedAdd are the `derived` fields crossing format 8 may add, written only when set.
type DerivedAdd struct {
	RoutePattern string `yaml:"route_pattern,omitempty"`
	MockDetect   string `yaml:"mock_detect,omitempty"`
	MockFactory  string `yaml:"mock_factory,omitempty"`
	MockTyped    string `yaml:"mock_typed,omitempty"`
}

// Empty says whether there is nothing to add.
func (s Syntax) Empty() bool {
	return s.Dialect == nil && s.Tests == nil && s.Derived == nil && len(s.Headless) == 0 && len(s.Runners) == 0
}

// DeclaredSyntax is what crossing format 8 adds to the project, given its configuration and
// the files of its layers (paths relative to the root):
//   - the comments of every example whose glob covers one of its files, Anchors' own formats
//     left out, when it declares none; an extension no example covers is named;
//   - its family's lexicon, every field it left empty — the naming conventions alone with no
//     family, which the old binary applied to any project;
//   - the route pattern and the mock shapes the old binary fell back to, for the gates that
//     read them;
//   - `headless` on a test layer of YAML files, whose header the old binary never asked for;
//   - the example runners of the tools its suites run with, which the old monitor knew.
func DeclaredSyntax(cfg *config.Config, files []string) Syntax {
	var s Syntax
	if cfg == nil {
		return s
	}
	d := config.Dialect{}
	if cfg.Dialect != nil {
		d = *cfg.Dialect
	}
	add := &config.Dialect{}

	if len(d.Comments) == 0 {
		add.Comments, s.Manual = commentsFor(files)
	}

	ex := config.FamilyLexicon(d.Family)
	if ex == nil {
		promise, slice := config.NamingExample()
		fill(&add.SetPromise, d.SetPromise, promise)
		fill(&add.SetSlice, d.SetSlice, slice)
		if d.EnvRead == "" && declares(cfg, "env-declared") {
			s.Manual = append(s.Manual, "dialect.env_read: with no family the old binary read every family's environment reads together — declare your language's")
		}
	} else {
		fill(&add.ExportedFunc, d.ExportedFunc, ex.ExportedFunc)
		fill(&add.ParamName, d.ParamName, ex.ParamName)
		fill(&add.Loop, d.Loop, ex.Loop)
		fill(&add.Cursor, d.Cursor, ex.Cursor)
		fill(&add.SetPromise, d.SetPromise, ex.SetPromise)
		fill(&add.SetSlice, d.SetSlice, ex.SetSlice)
		fill(&add.HTTPStatus, d.HTTPStatus, ex.HTTPStatus)
		fill(&add.HTTPStatusDynamic, d.HTTPStatusDynamic, ex.HTTPStatusDynamic)
		fill(&add.EnvRead, d.EnvRead, ex.EnvRead)
		fill(&add.Definition, d.Definition, ex.Definition)
		fill(&add.StatementEnd, d.StatementEnd, ex.StatementEnd)
		fill(&add.SignatureEnd, d.SignatureEnd, ex.SignatureEnd)
		if len(d.HandlePatterns) == 0 {
			add.HandlePatterns = ex.HandlePatterns
		}
		if len(d.LogPatterns) == 0 {
			add.LogPatterns = ex.LogPatterns
		}
		if len(d.GuardPatterns) == 0 {
			add.GuardPatterns = ex.GuardPatterns
		}
		if len(d.Limits) == 0 {
			add.Limits = ex.Limits
		}
		s.Tests, add.Tests = testsFor(d.Tests, ex.Tests)
	}
	if d.Tests != nil && d.Tests.Pattern != "" && !regexp.MustCompile(`\(\?P<title`).MatchString(d.Tests.Pattern) {
		s.Manual = append(s.Manual, "dialect.tests.pattern: name where the title stands with a `title` group — `(?P<title>…)` —; the old binary read the string literal after the match")
	}
	if !reflectEmpty(add) {
		s.Dialect = add
	}

	der := &DerivedAdd{}
	have := config.Derived{}
	if cfg.Derived != nil {
		have = *cfg.Derived
	}
	if have.RoutePattern == "" && declares(cfg, "route-exists") {
		der.RoutePattern = config.FamilyRoutePattern(d.Family)
		if der.RoutePattern == "" {
			der.RoutePattern = config.FamilyRoutePattern("ts") // the old binary's fallback, whatever the family
		}
	}
	if declares(cfg, "mock-typed") {
		mocks := config.FamilyMocks(d.Family)
		if mocks == nil || mocks.MockFactory == "" {
			mocks = config.FamilyMocks("ts") // the old binary knew Jest's and Vitest's, whatever the family
		}
		fill(&der.MockDetect, have.MockDetect, mocks.MockDetect)
		fill(&der.MockFactory, have.MockFactory, mocks.MockFactory)
		fill(&der.MockTyped, have.MockTyped, mocks.MockTyped)
	}
	if der.RoutePattern != "" || der.MockDetect != "" || der.MockFactory != "" || der.MockTyped != "" {
		s.Derived = der
	}

	for name, l := range cfg.Layers {
		if l.Kind == "test" && !l.Headless && dataFiles(l.Pattern) {
			s.Headless = append(s.Headless, name)
		}
	}
	sort.Strings(s.Headless)
	s.Runners = runnersFor(cfg)
	return s
}

// runnersFor are the example runners some suite of the project runs with — its `run` the
// runner's command line —, minus the ones it declares by name.
func runnersFor(cfg *config.Config) []config.MonitorRunner {
	declared := map[string]bool{}
	if cfg.Monitor != nil {
		for _, r := range cfg.Monitor.Runners {
			declared[r.Name] = true
		}
	}
	var out []config.MonitorRunner
	for _, r := range config.ExampleRunners() {
		if declared[r.Name] {
			continue
		}
		match := regexp.MustCompile(r.Match)
		for _, suites := range [][]config.Suite{cfg.Tests, cfg.Mutation} {
			for _, su := range suites {
				if su.Run != "" && match.MatchString(su.Run) && !declared[r.Name] {
					declared[r.Name] = true
					out = append(out, r)
				}
			}
		}
	}
	return out
}

// Fields names what is added, a line each: `dialect.comments` with the files of its entries,
// every other field by its path.
func (s Syntax) Fields() []string {
	var out []string
	if s.Dialect != nil {
		for _, k := range keysOf("dialect.", s.Dialect) {
			if k == "dialect.comments" {
				var globs []string
				for _, c := range s.Dialect.Comments {
					globs = append(globs, c.Files)
				}
				k += " (" + strings.Join(globs, ", ") + ")"
			}
			out = append(out, k)
		}
	}
	if s.Tests != nil {
		out = append(out, keysOf("dialect.tests.", s.Tests)...)
	}
	if s.Derived != nil {
		out = append(out, keysOf("derived.", s.Derived)...)
	}
	for _, l := range s.Headless {
		out = append(out, "layers."+l+".headless")
	}
	for _, r := range s.Runners {
		out = append(out, "monitor.runners: "+r.Name)
	}
	return out
}

// keysOf are the keys v is written with, in order, after the prefix.
func keysOf(prefix string, v any) []string {
	var doc yaml.Node
	b, err := yaml.Marshal(v)
	if err != nil || yaml.Unmarshal(b, &doc) != nil || len(doc.Content) == 0 {
		return nil
	}
	var out []string
	for i := 0; i+1 < len(doc.Content[0].Content); i += 2 {
		out = append(out, prefix+doc.Content[0].Content[i].Value)
	}
	return out
}

// fill sets dst to the example when the project declared nothing.
func fill(dst *string, declared, example string) {
	if declared == "" {
		*dst = example
	}
}

// testsFor is what of the family's tests a project lacks: the whole example when it declares
// none, else the fields the old binary filled inside its own — the family's way to open a
// test when it declared only its assertion, and the shape of a block, which it applied to
// every project.
func testsFor(declared, example *config.TestsSource) (inside, whole *config.TestsSource) {
	if example == nil {
		return nil, nil
	}
	if declared == nil {
		cp := *example
		return nil, &cp
	}
	add := config.TestsSource{}
	if declared.Pattern == "" && declared.Script == "" {
		add.Pattern = example.Pattern
		fill(&add.Assertion, declared.Assertion, example.Assertion)
	}
	fill(&add.BlockOpen, declared.BlockOpen, example.BlockOpen)
	fill(&add.BlockClose, declared.BlockClose, example.BlockClose)
	fill(&add.EmptyBody, declared.EmptyBody, example.EmptyBody)
	fill(&add.Placeholder, declared.Placeholder, example.Placeholder)
	if add == (config.TestsSource{}) {
		return nil, nil
	}
	return &add, nil
}

// commentsFor are the examples' entries covering the files, and the extensions none covers,
// named to declare by hand.
func commentsFor(files []string) ([]config.Comments, []string) {
	out, uncovered := config.CommentsCovering(files)
	if len(uncovered) == 0 {
		return out, nil
	}
	return out, []string{fmt.Sprintf("dialect.comments: no example covers %s — declare how these files comment, or their headers and flags are not read", strings.Join(uncovered, ", "))}
}

// declares says whether the project declares the gate, by its name or its check.
func declares(cfg *config.Config, check string) bool {
	for _, g := range cfg.Gates {
		if g.Name == check || g.Check == check {
			return true
		}
	}
	return false
}

// dataFiles says whether a layer's pattern names only YAML files: a runner's flows, which
// the old binary never asked for a header.
func dataFiles(pattern string) bool {
	p := strings.TrimSpace(pattern)
	for _, suf := range []string{".yaml", ".yml", ".{yaml,yml}", ".{yml,yaml}", ".y{a,}ml"} {
		if strings.HasSuffix(p, suf) {
			return true
		}
	}
	return false
}

// reflectEmpty says whether no field of the dialect is set.
func reflectEmpty(d *config.Dialect) bool {
	b, err := yaml.Marshal(d)
	return err == nil && strings.TrimSpace(string(b)) == "{}"
}

// WithSyntax writes the additions into the text of an anchors.yaml: each section's fields at
// the end of its block, indented as its children are, the section appended when the file has
// none; a layer gets `headless: true`. The rest of the text — its comments, order and
// quoting — is kept. A section written in flow style cannot take a line, and is named in the
// error.
func WithSyntax(text string, s Syntax) (string, error) {
	var err error
	if s.Dialect != nil {
		if text, err = addFields(text, []string{"dialect"}, s.Dialect); err != nil {
			return "", err
		}
	}
	if s.Tests != nil {
		if text, err = addFields(text, []string{"dialect", "tests"}, s.Tests); err != nil {
			return "", err
		}
	}
	if s.Derived != nil {
		if text, err = addFields(text, []string{"derived"}, s.Derived); err != nil {
			return "", err
		}
	}
	for _, l := range s.Headless {
		if text, err = addFields(text, []string{"layers", l}, map[string]bool{"headless": true}); err != nil {
			return "", err
		}
	}
	if len(s.Runners) > 0 {
		if text, err = addFields(text, []string{"monitor", "runners"}, s.Runners); err != nil {
			return "", err
		}
	}
	return text, nil
}

// addFields writes v's YAML at the end of the block at the path.
func addFields(text string, at []string, v any) (string, error) {
	var node yaml.Node
	if err := node.Encode(v); err != nil {
		return "", err
	}
	flowShortLists(&node)
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&node); err != nil {
		return "", err
	}
	body := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	eol := "\n"
	if strings.Contains(text, "\r\n") {
		eol = "\r\n"
	}
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	// lvl is the indentation of the keys at this depth.
	from, to, lvl := 0, len(lines), 0
	for depth, key := range at {
		keyLine, end, child, flow := blockOf(lines, from, to, key)
		if flow {
			if inline, ok := intoFlow(lines[keyLine], body); ok {
				lines[keyLine] = inline
				return strings.Join(lines, eol), nil
			}
			return "", fmt.Errorf("`%s` is written in flow style: add %s by hand", strings.Join(at[:depth+1], "."), strings.Join(body, "; "))
		}
		if keyLine < 0 {
			// The section is missing: written whole at the end of its parent.
			pad := strings.Repeat(" ", lvl)
			block := []string{pad + key + ":"}
			for _, rest := range at[depth+1:] {
				pad += "  "
				block = append(block, pad+rest+":")
			}
			for _, b := range body {
				block = append(block, pad+"  "+b)
			}
			i := endOf(lines, from, to, lvl)
			if depth == 0 {
				i = to
				for i > 0 && strings.TrimSpace(lines[i-1]) == "" {
					i--
				}
				block = append([]string{""}, block...)
			}
			lines = insertAt(lines, i, block)
			return strings.Join(lines, eol), nil
		}
		from, to, lvl = keyLine+1, end, child
	}
	block := make([]string, len(body))
	for i, b := range body {
		block[i] = strings.Repeat(" ", lvl) + b
	}
	lines = insertAt(lines, endOf(lines, from, to, lvl), block)
	return strings.Join(lines, eol), nil
}

// flowShortLists writes a list of short texts on one line — `line: ['//']`, a block's
// `['/*', '*/']` —, as a person would; a list of long patterns keeps a line each.
func flowShortLists(n *yaml.Node) {
	if n.Kind == yaml.SequenceNode {
		short := true
		for _, c := range n.Content {
			flowShortLists(c)
			if c.Kind == yaml.SequenceNode {
				short = short && c.Style == yaml.FlowStyle
				continue
			}
			short = short && c.Kind == yaml.ScalarNode && len(c.Value) <= 12
		}
		if short {
			n.Style = yaml.FlowStyle
		}
		return
	}
	for _, c := range n.Content {
		flowShortLists(c)
	}
}

// blockOf finds the key among the keys of lines[from:to] — those at the indentation of its
// first one —: its line, where its block ends, the indentation of its children (two more
// when it has none) and whether its value is written on its line — flow style. A missing
// key gives -1.
func blockOf(lines []string, from, to int, key string) (at, end, child int, flow bool) {
	keyRE := regexp.MustCompile(`^( *)(?:"` + regexp.QuoteMeta(key) + `"|'` + regexp.QuoteMeta(key) + `'|` + regexp.QuoteMeta(key) + `):(.*)$`)
	at = -1
	want := -1
	for i := from; i < to; i++ {
		ind, content := indentOf(lines[i])
		if content == "" || strings.HasPrefix(content, "#") {
			continue
		}
		if want < 0 {
			want = ind
		}
		if ind != want {
			continue
		}
		if m := keyRE.FindStringSubmatch(lines[i]); m != nil {
			at = i
			rest := strings.TrimSpace(m[2])
			flow = rest != "" && !strings.HasPrefix(rest, "#")
			break
		}
	}
	if at < 0 || flow {
		return at, 0, 0, flow
	}
	keyIndent, _ := indentOf(lines[at])
	end = to
	child = -1
	for i := at + 1; i < to; i++ {
		ind, content := indentOf(lines[i])
		if content == "" || strings.HasPrefix(content, "#") {
			continue
		}
		if ind <= keyIndent {
			end = i
			break
		}
		if child < 0 {
			child = ind
		}
	}
	if child < 0 {
		child = keyIndent + 2
	}
	return at, end, child, false
}

// endOf is where lines go at the end of a block ending at `to`, whose keys are at `lvl`:
// before the blank lines and the comments that open what comes next, less indented.
func endOf(lines []string, from, to, lvl int) int {
	i := to
	for i > from {
		ind, content := indentOf(lines[i-1])
		if content == "" || (strings.HasPrefix(content, "#") && ind < lvl) {
			i--
			continue
		}
		break
	}
	return i
}

// intoFlow writes the fields into a one-line flow mapping — `name: {kind: test}` —, when
// they are one line each.
func intoFlow(line string, body []string) (string, bool) {
	t := strings.TrimRight(line, " ")
	if !strings.HasSuffix(t, "}") || strings.Count(t, "{") != 1 {
		return "", false
	}
	for _, b := range body {
		if strings.HasPrefix(b, " ") || strings.HasPrefix(b, "-") {
			return "", false
		}
	}
	return strings.TrimSuffix(t, "}") + ", " + strings.Join(body, ", ") + "}", true
}

func indentOf(l string) (int, string) {
	t := strings.TrimLeft(l, " ")
	return len(l) - len(t), strings.TrimSpace(t)
}

func insertAt(lines []string, i int, block []string) []string {
	out := append([]string{}, lines[:i]...)
	out = append(out, block...)
	return append(out, lines[i:]...)
}
