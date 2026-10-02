package initx

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LANGUAGE IS DIALECT, NOT STRUCTURE.
//
// What a language decides is how a test file is NAMED and how a comment is written — the
// things the map needs to read the project. Where the project keeps each layer is the
// project's own decision, and nothing here reads or proposes it.
//
// The conventions below are the naming forms the ecosystems use. Inference counts the
// files that follow each one, so the patterns written to anchors.yaml are the ones the
// project's own files follow; a language's default convention is used only when there is
// no test file at all yet.

// TestConvention is one way a language names a test file.
type TestConvention struct {
	// Prefix and Suffix are what a test file's name starts and ends with (`test_` and
	// `.py`; `` and `_test.go`).
	Prefix, Suffix string
	// Ext is the code extension the test sits beside, without the dot.
	Ext string
	// Family is the dialect family the convention belongs to.
	Family string
}

// testConventions are the naming forms inference recognises, in order: a longer suffix
// before a shorter one it ends with (`.spec.tsx` before `.tsx`), so a file is read once.
var testConventions = []TestConvention{
	{Suffix: "_test.go", Ext: "go", Family: "go"},
	{Suffix: ".test.tsx", Ext: "tsx", Family: "ts"},
	{Suffix: ".test.ts", Ext: "ts", Family: "ts"},
	{Suffix: ".test.jsx", Ext: "jsx", Family: "ts"},
	{Suffix: ".test.js", Ext: "js", Family: "ts"},
	{Suffix: ".spec.tsx", Ext: "tsx", Family: "ts"},
	{Suffix: ".spec.ts", Ext: "ts", Family: "ts"},
	{Suffix: ".spec.jsx", Ext: "jsx", Family: "ts"},
	{Suffix: ".spec.js", Ext: "js", Family: "ts"},
	{Suffix: "_test.py", Ext: "py", Family: "python"},
	{Prefix: "test_", Suffix: ".py", Ext: "py", Family: "python"},
	{Suffix: "_spec.rb", Ext: "rb", Family: "ruby"},
	{Suffix: "_test.rb", Ext: "rb", Family: "ruby"},
	{Suffix: "Test.java", Ext: "java", Family: "java"},
	{Suffix: "Tests.java", Ext: "java", Family: "java"},
	{Suffix: "Test.kt", Ext: "kt", Family: "kotlin"},
	{Suffix: "Tests.cs", Ext: "cs", Family: "csharp"},
	{Suffix: "Test.cs", Ext: "cs", Family: "csharp"},
	{Suffix: "Test.php", Ext: "php", Family: "php"},
}

// familyDefault is the convention a family's ecosystem uses when the project has no test
// file yet to read one from.
var familyDefault = map[string]TestConvention{
	"go":     {Suffix: "_test.go", Ext: "go", Family: "go"},
	"ts":     {Suffix: ".test.ts", Ext: "ts", Family: "ts"},
	"python": {Prefix: "test_", Suffix: ".py", Ext: "py", Family: "python"},
	"ruby":   {Suffix: "_spec.rb", Ext: "rb", Family: "ruby"},
	"java":   {Suffix: "Test.java", Ext: "java", Family: "java"},
	"kotlin": {Suffix: "Test.kt", Ext: "kt", Family: "kotlin"},
	"csharp": {Suffix: "Tests.cs", Ext: "cs", Family: "csharp"},
	"php":    {Suffix: "Test.php", Ext: "php", Family: "php"},
}

// manifestFamily is the dialect family a manifest at the project's root declares.
var manifestFamily = []struct {
	file, family string
}{
	{"go.mod", "go"},
	{"Cargo.toml", "rust"},
	{"pyproject.toml", "python"}, {"setup.py", "python"}, {"requirements.txt", "python"},
	{"Gemfile", "ruby"},
	{"composer.json", "php"},
	{"package.json", "ts"},
}

// extFamily is the dialect family of a code extension, when no manifest says.
var extFamily = map[string]string{
	".go": "go", ".ts": "ts", ".tsx": "ts", ".js": "ts", ".jsx": "ts", ".py": "python",
	".rs": "rust", ".rb": "ruby", ".php": "php", ".cs": "csharp", ".java": "java", ".kt": "kotlin",
}

// coverageHints are how each family's test runner emits what `anchors ingest` reads: the
// JUnit report of the run and the lcov of the lines. Dialect too — the command a language
// runs its tests with — and shown by init so a project does not have to find it alone.
var coverageHints = map[string]string{
	"go":     "go test ./... -coverprofile=cover.out; gcov2lcov < cover.out > lcov.info; go test -json ./... | go-junit-report > junit.xml",
	"ts":     "jest/vitest --coverage with the lcov reporter, and a junit reporter → coverage/lcov.info + junit.xml",
	"python": "pytest --cov --cov-report=lcov --junitxml=junit.xml (pytest-cov)",
	"rust":   "cargo llvm-cov --lcov --output-path lcov.info; cargo nextest run with the junit profile → junit.xml",
	"java":   "the surefire reports (JUnit XML) and jacoco's report, converted to lcov",
	"kotlin": "the surefire reports (JUnit XML) and jacoco's report, converted to lcov",
}

// CoverageHint is how the family's tests emit the reports ingest reads; empty for a
// family with no hint.
func CoverageHint(family string) string { return coverageHints[family] }

// gateScripts are how a family runs a gate step written in its own language.
var gateScripts = map[string]string{
	"go":     "go run ./tools/gates <gate>",
	"ts":     "node tools/gates.mjs <gate>",
	"python": "python -m tools.gates <gate>",
	"ruby":   "ruby tools/gates.rb <gate>",
	"rust":   "cargo run --bin gates -- <gate>",
	"java":   "a JBang or Gradle task: `gradle gates -Pgate=<gate>`",
	"kotlin": "a Gradle task: `gradle gates -Pgate=<gate>`",
	"csharp": "dotnet run --project tools/Gates -- <gate>",
	"php":    "php tools/gates.php <gate>",
}

// GateScriptExample is how a gate step that does more than call one tool is written in the
// family's language, so it runs the same on every system; a family with no example gets
// the generic advice.
//
// Shell is what breaks there: BSD sed has no `\s`, `2>&1` mixes a tool's noise into the
// findings, and Windows has no `sh` (reported from baas-proxy, whose scripts copied from
// this repository broke on macOS).
func GateScriptExample(family string) string {
	if s, ok := gateScripts[family]; ok {
		return s
	}
	return "a program in the project's own language"
}

// testConventionOf is the convention a file name follows, if it is a test file.
func testConventionOf(name string) (TestConvention, bool) {
	for _, c := range testConventions {
		if strings.HasPrefix(name, c.Prefix) && strings.HasSuffix(name, c.Suffix) && len(name) > len(c.Prefix)+len(c.Suffix) {
			return c, true
		}
	}
	return TestConvention{}, false
}

// Glob is the convention's pattern for every test file of the project.
func (c TestConvention) Glob() string { return "**/" + c.Prefix + "*" + c.Suffix }

// Template is where the convention puts a unit's test, beside its code.
func (c TestConvention) Template() string { return "{{dir}}/" + c.Prefix + "{{name}}" + c.Suffix }

// stemOfTest is a test file's unit name: its name without the convention's prefix and
// suffix.
func (c TestConvention) stemOfTest(name string) string {
	return strings.TrimSuffix(strings.TrimPrefix(name, c.Prefix), c.Suffix)
}

// detectFamily is the project's dialect family: the manifest at its root, or the
// family of its most common code extension.
func detectFamily(root string, extCount map[string]int) string {
	for _, m := range manifestFamily {
		if _, err := os.Stat(filepath.Join(root, m.file)); err == nil {
			return m.family
		}
	}
	for _, ext := range topKeys(extCount, len(extCount)) {
		if f, ok := extFamily[ext]; ok {
			return f
		}
	}
	return ""
}

// rankConventions orders the conventions the project's test files follow by how many
// follow each, most first, ties by suffix.
func rankConventions(counts map[TestConvention]int) []TestConvention {
	out := make([]TestConvention, 0, len(counts))
	for c := range counts {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if counts[out[i]] != counts[out[j]] {
			return counts[out[i]] > counts[out[j]]
		}
		return out[i].Prefix+out[i].Suffix < out[j].Prefix+out[j].Suffix
	})
	return out
}

// TestGlobs are the test patterns the project's own files follow, or its family's
// default convention when it has no test file yet; empty when neither is known.
func (p *Proposal) TestGlobs() []string {
	var out []string
	for _, c := range p.conventions() {
		out = append(out, c.Glob())
	}
	return out
}

// TestTemplate is where a unit's test sits beside its code, by the project's dominant
// convention; empty when no convention is known.
func (p *Proposal) TestTemplate() string {
	cs := p.conventions()
	if len(cs) == 0 {
		return ""
	}
	return cs[0].Template()
}

// conventions are the detected conventions, or the family default.
func (p *Proposal) conventions() []TestConvention {
	if len(p.TestConventions) > 0 {
		return p.TestConventions
	}
	if c, ok := familyDefault[p.Family]; ok {
		return []TestConvention{c}
	}
	return nil
}
