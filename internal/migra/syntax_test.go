// @anchors
//   code: SYNTT
//   ref: MGSYN

package migra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// load reads a configuration text as the project would.
func load(t *testing.T, text string) *config.Config {
	t.Helper()
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := config.Load(p)
	if err != nil {
		t.Fatalf("%v\n%s", err, text)
	}
	return c
}

func TestDeclaredSyntax_commentsOfTheFiles(t *testing.T) {
	t.Run("MGSYN-B01: The comments of the examples covering the project's files are added, and an extension none covers is named", func(t *testing.T) {})
	s := DeclaredSyntax(&config.Config{}, []string{"a.go", "b.ts", "flows/x.yaml", "c.zig", "a.spec.md", "a.feature"})
	if s.Dialect == nil {
		t.Fatal("comments are added")
	}
	var globs []string
	for _, c := range s.Dialect.Comments {
		globs = append(globs, c.Files)
	}
	got := strings.Join(globs, " ")
	for _, want := range []string{"**/*.{ts,tsx,js,jsx,mjs,cjs,mts,cts}", "**/*.go", "**/*.{yaml,yml,toml,sh,bash,zsh}"} {
		if !strings.Contains(got, want) {
			t.Errorf("an entry for %s, got %s", want, got)
		}
	}
	if len(globs) != 3 {
		t.Errorf("only the entries the files use, got %s", got)
	}
	if len(s.Manual) != 1 || !strings.Contains(s.Manual[0], ".zig") {
		t.Errorf("the extension no example covers is named, got %v", s.Manual)
	}
	declared := &config.Config{Dialect: &config.Dialect{Comments: []config.Comments{{Files: "**/*.go", Line: []string{"//"}}}}}
	if s := DeclaredSyntax(declared, []string{"a.go", "b.ts"}); s.Dialect != nil && len(s.Dialect.Comments) > 0 {
		t.Errorf("a project declaring its comments gets none, got %+v", s.Dialect.Comments)
	}
}

func TestDeclaredSyntax_lexiconOfTheFamily(t *testing.T) {
	t.Run("MGSYN-B02: The family's lexicon fills what the project left empty", func(t *testing.T) {})
	cfg := &config.Config{Dialect: &config.Dialect{Family: "go", Comments: config.FamilyComments("go"), Loop: `\bmine\b`}}
	s := DeclaredSyntax(cfg, []string{"a.go"})
	ex := config.FamilyLexicon("go")
	if s.Dialect == nil || s.Dialect.Loop != "" || s.Dialect.ExportedFunc != ex.ExportedFunc ||
		s.Dialect.StatementEnd != ex.StatementEnd || len(s.Dialect.HandlePatterns) == 0 || s.Dialect.Tests == nil ||
		s.Dialect.SetPromise != ex.SetPromise {
		t.Fatalf("the declared loop is kept and the rest comes from the example, got %+v", s.Dialect)
	}
	none := DeclaredSyntax(&config.Config{Dialect: &config.Dialect{Comments: config.FamilyComments("go")}}, nil)
	promise, slice := config.NamingExample()
	if none.Dialect == nil || none.Dialect.SetPromise != promise || none.Dialect.SetSlice != slice || none.Dialect.ExportedFunc != "" {
		t.Errorf("with no family only the naming conventions, got %+v", none.Dialect)
	}
	env := DeclaredSyntax(&config.Config{Gates: []config.Gate{{Name: "env-declared"}}}, nil)
	if len(env.Manual) != 1 || !strings.Contains(env.Manual[0], "dialect.env_read") {
		t.Errorf("env-declared with no family is named, got %v", env.Manual)
	}
}

func TestDeclaredSyntax_insideTheProjectsTests(t *testing.T) {
	t.Run("MGSYN-B03: Inside the project's own tests only what the old binary filled is added", func(t *testing.T) {})
	ex := config.FamilyLexicon("ts").Tests
	onlyAssertion := &config.Config{Dialect: &config.Dialect{Family: "ts", Tests: &config.TestsSource{Assertion: `\bcheck\(`}}}
	s := DeclaredSyntax(onlyAssertion, nil)
	if s.Tests == nil || s.Tests.Pattern != ex.Pattern || s.Tests.Assertion != "" || s.Tests.BlockOpen != ex.BlockOpen || s.Dialect.Tests != nil {
		t.Errorf("an assertion alone takes the family's pattern and the block, inside its tests, got %+v / %+v", s.Tests, s.Dialect.Tests)
	}
	script := &config.Config{Dialect: &config.Dialect{Family: "ts", Tests: &config.TestsSource{Script: "node l.mjs"}}}
	if s := DeclaredSyntax(script, nil); s.Tests == nil || s.Tests.Pattern != "" || s.Tests.Assertion != "" || s.Tests.BlockClose != ex.BlockClose {
		t.Errorf("a script takes only the block, got %+v", s.Tests)
	}
	untitled := &config.Config{Dialect: &config.Dialect{Family: "go", Tests: &config.TestsSource{Pattern: `t\.Run\(`}}}
	if s := DeclaredSyntax(untitled, nil); len(s.Manual) != 1 || !strings.Contains(s.Manual[0], "title") {
		t.Errorf("a pattern with no title group is named, got %v", s.Manual)
	}
}

func TestDeclaredSyntax_shapesOfTheGates(t *testing.T) {
	t.Run("MGSYN-B04: The route and mock shapes are added for the gates that read them, and a layer of flows becomes headless", func(t *testing.T) {})
	cfg := &config.Config{
		Dialect: &config.Dialect{Family: "go", Comments: config.FamilyComments("go")},
		Gates:   []config.Gate{{Name: "route-exists"}, {Name: "mock-typed"}},
		Derived: &config.Derived{MockDetect: `mine\(`},
		Layers: map[string]config.Layer{
			"flows": {Pattern: "e2e/**/*.yaml", Kind: "test"},
			"done":  {Pattern: "x/**/*.yml", Kind: "test", Headless: true},
			"unit":  {Pattern: "**/*_test.go", Kind: "test"},
		},
	}
	s := DeclaredSyntax(cfg, nil)
	mocks := config.FamilyMocks("ts")
	if s.Derived == nil || s.Derived.RoutePattern != config.FamilyRoutePattern("ts") || s.Derived.MockDetect != "" ||
		s.Derived.MockFactory != mocks.MockFactory || s.Derived.MockTyped != mocks.MockTyped {
		t.Errorf("the old fallbacks, the declared detection kept, got %+v", s.Derived)
	}
	if strings.Join(s.Headless, ",") != "flows" {
		t.Errorf("only the flows layer not yet headless, got %v", s.Headless)
	}
	if s := DeclaredSyntax(&config.Config{Dialect: cfg.Dialect}, nil); s.Derived != nil {
		t.Errorf("no gate, no shape, got %+v", s.Derived)
	}
}

const projectYAML = `# the project
version: 8
dialect:
  family: go
  # how a test is written
  tests:
    script: node l.mjs

layers:
  flows: {pattern: "e2e/**/*.yaml", kind: test}
  unit:
    pattern: "**/*_test.go"
    kind: test

# the gates
gates:
  - name: route-exists
`

func TestWithSyntax_writesIntoTheText(t *testing.T) {
	t.Run("MGSYN-B05: The additions go at the end of each block, the rest of the text kept, and the file reads back", func(t *testing.T) {})
	cfg := load(t, projectYAML)
	s := DeclaredSyntax(cfg, []string{"a.go"})
	out, err := WithSyntax(projectYAML, s)
	if err != nil {
		t.Fatal(err)
	}
	for _, keep := range []string{"# the project\n", "  # how a test is written\n", "\n# the gates\ngates:\n"} {
		if !strings.Contains(out, keep) {
			t.Errorf("the text is kept: %q missing in\n%s", keep, out)
		}
	}
	if !strings.Contains(out, `flows: {pattern: "e2e/**/*.yaml", kind: test, headless: true}`) {
		t.Errorf("a flow-style layer gets headless inline:\n%s", out)
	}
	back := load(t, out)
	d := back.DialectFor()
	if len(d.Comments) == 0 || d.ExportedFunc == "" || d.Tests.Script != "node l.mjs" || d.Tests.BlockOpen == "" ||
		back.Derived == nil || back.Derived.RoutePattern == "" || !back.Layers["flows"].Headless {
		t.Errorf("everything reads back where it belongs:\n%s", out)
	}
	if again := DeclaredSyntax(back, []string{"a.go"}); !again.Empty() {
		t.Errorf("a second pass adds nothing, got %+v", again.Fields())
	}
	if !strings.Contains(strings.Join(s.Fields(), "\n"), "dialect.comments (**/*.go") || !strings.Contains(strings.Join(s.Fields(), "\n"), "layers.flows.headless") {
		t.Errorf("the fields are named, got %v", s.Fields())
	}
}

func TestWithSyntax_sectionsMissingOrInFlow(t *testing.T) {
	t.Run("MGSYN-B06: A missing section is written whole, and one in flow style that cannot take the fields is refused", func(t *testing.T) {})
	text := "version: 8\ngates:\n  - name: route-exists\n"
	out, err := WithSyntax(text, DeclaredSyntax(load(t, text), []string{"a.go"}))
	if err != nil {
		t.Fatal(err)
	}
	if back := load(t, out); len(back.DialectFor().Comments) == 0 || back.Derived == nil || back.Derived.RoutePattern == "" {
		t.Errorf("dialect and derived are appended whole:\n%s", out)
	}
	flow := "version: 8\ndialect: {family: go}\n"
	if _, err := WithSyntax(flow, DeclaredSyntax(load(t, flow), []string{"a.go"})); err == nil || !strings.Contains(err.Error(), "`dialect` is written in flow style") {
		t.Errorf("a flow-style dialect is refused naming it, got %v", err)
	}
}

func TestDeclaredSyntax_runnersOfTheSuites(t *testing.T) {
	t.Run("MGSYN-B07: The example runners the project's suites run with are added, under monitor.runners", func(t *testing.T) {})
	text := "version: 8\ntests:\n  - layer: unit\n    run: npx jest --ci\n  - layer: e2e\n    run: maestro test flows/\nmutation:\n  - layer: unit\n    run: npx stryker run\nmonitor:\n  runners:\n    - name: maestro\n      match: 'maestro'\n"
	cfg := load(t, text)
	s := DeclaredSyntax(cfg, nil)
	var names []string
	for _, r := range s.Runners {
		names = append(names, r.Name)
	}
	if strings.Join(names, ",") != "jest,stryker" {
		t.Fatalf("the runners of its suites, the declared one kept, got %v", names)
	}
	out, err := WithSyntax(text, s)
	if err != nil {
		t.Fatal(err)
	}
	back := load(t, out)
	if len(back.Monitor.Runners) != 3 || back.Monitor.Runners[1].Name != "jest" {
		t.Errorf("appended to the project's runners:\n%s", out)
	}
	if !strings.Contains(strings.Join(s.Fields(), "\n"), "monitor.runners: jest") {
		t.Errorf("named: %v", s.Fields())
	}
	none := "version: 8\ntests:\n  - layer: unit\n    run: go test ./...\n"
	out, err = WithSyntax(none, DeclaredSyntax(load(t, none), nil))
	if err != nil || len(load(t, out).Monitor.Runners) != 1 {
		t.Errorf("a project with no monitor block gets one:\n%s %v", out, err)
	}
}
