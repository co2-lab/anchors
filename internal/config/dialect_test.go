// @anchors
//   code: DLTSA
//   ref: DLCTI

package config

import (
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestDialectFor_withoutDeclarationOnlyTheNamingDefaults(t *testing.T) {
	t.Run("DLCTI-B01: A project that declares no dialect gets only the naming defaults", func(t *testing.T) {})
	for name, c := range map[string]*Config{"nil config": nil, "no dialect block": {}} {
		d := c.DialectFor()
		if d.SetPromise != defaultSetPromise || d.SetSlice != defaultSetSlice {
			t.Errorf("%s: the naming defaults must be the floor, got %q / %q", name, d.SetPromise, d.SetSlice)
		}
		if d.ExportedFunc != "" || d.Loop != "" || d.Cursor != "" || len(d.HandlePatterns) != 0 {
			t.Errorf("%s: no family means no language lexicon, got %+v", name, d)
		}
	}
}

func TestDialectFor_familyFillsWhatIsNotDeclared(t *testing.T) {
	t.Run("DLCTI-B02: The family fills every field the project left empty, and a declared field wins", func(t *testing.T) {})
	c := &Config{Dialect: &Dialect{
		Family:         "go",
		Loop:           `\bloop\b`,
		HandlePatterns: []string{`\bmine\b`},
	}}
	d := c.DialectFor()
	base := dialectFamilies["go"]
	if d.Loop != `\bloop\b` || !reflect.DeepEqual(d.HandlePatterns, []string{`\bmine\b`}) {
		t.Errorf("the declared fields must win over the family, got Loop=%q Handle=%v", d.Loop, d.HandlePatterns)
	}
	if d.ExportedFunc != base.ExportedFunc || d.Cursor != base.Cursor || d.ParamName != base.ParamName ||
		d.HTTPStatus != base.HTTPStatus || d.HTTPStatusDynamic != base.HTTPStatusDynamic ||
		!reflect.DeepEqual(d.LogPatterns, base.LogPatterns) {
		t.Errorf("the family must fill the empty fields, got %+v", d)
	}
	if c.Dialect.ExportedFunc != "" {
		t.Errorf("the project's declaration must not be rewritten, got %q", c.Dialect.ExportedFunc)
	}
}

func TestDialectFor_familyIgnoresCase(t *testing.T) {
	t.Run("DLCTI-B03: The family name is matched ignoring case", func(t *testing.T) {})
	d := (&Config{Dialect: &Dialect{Family: "TS"}}).DialectFor()
	if d.ExportedFunc != dialectFamilies["ts"].ExportedFunc {
		t.Errorf("family TS must resolve to ts, got ExportedFunc=%q", d.ExportedFunc)
	}
}

func TestDialectFor_unknownFamilyContributesNothing(t *testing.T) {
	t.Run("DLCTI-B04: An unknown family contributes nothing, and the known ones are listed in order", func(t *testing.T) {})
	d := (&Config{Dialect: &Dialect{Family: "cobol"}}).DialectFor()
	if d.ExportedFunc != "" || d.Loop != "" || d.Cursor != "" || len(d.HandlePatterns) != 0 {
		t.Errorf("an unknown family must leave the lexicon empty, got %+v", d)
	}
	known := KnownDialectFamilies()
	if !sort.StringsAreSorted(known) || len(known) != len(dialectFamilies) {
		t.Errorf("KnownDialectFamilies = %v, want every family, sorted", known)
	}
	for _, f := range known {
		if _, ok := dialectFamilies[f]; !ok {
			t.Errorf("KnownDialectFamilies lists %q, which has no lexicon", f)
		}
	}
}

func TestDialectFor_namingDefaultsAreOverridable(t *testing.T) {
	t.Run("DLCTI-B05: The naming conventions apply to any family and a declared one replaces them", func(t *testing.T) {})
	d := (&Config{Dialect: &Dialect{Family: "python", SetPromise: `^listar`}}).DialectFor()
	if d.SetPromise != `^listar` {
		t.Errorf("a declared set_promise must win, got %q", d.SetPromise)
	}
	if d.SetSlice != defaultSetSlice {
		t.Errorf("an undeclared set_slice must fall back to the default, got %q", d.SetSlice)
	}
	d = (&Config{Dialect: &Dialect{SetSlice: `^buscarAlguns`}}).DialectFor()
	if d.SetSlice != `^buscarAlguns` || d.SetPromise != defaultSetPromise {
		t.Errorf("a declared set_slice must win and set_promise fall back, got %q / %q", d.SetSlice, d.SetPromise)
	}
}

func TestGherkinFor(t *testing.T) {
	t.Run("DLCTI-B06: The Gherkin language defaults to English, is found in any case, and one outside the table keeps its code with English keywords", func(t *testing.T) {})
	for _, tc := range []struct {
		declared, lang, scenario string
	}{
		{"", "en", "Scenario"},
		{"pt", "pt", "Cenário"},
		{"PT", "pt", "Cenário"},
		{"eo", "eo", "Scenario"}, // Esperanto: Gherkin knows it, the table does not
		// A table key with capitals: the lookup once lower-cased the declaration first,
		// so zh-CN was never reached and got English keywords.
		{"zh-CN", "zh-CN", "场景"},
		{"ZH-cn", "zh-CN", "场景"},
		{"en-AU", "en-AU", "Scenario"}, // outside the table: kept as declared
	} {
		lang, kw := Dialect{GherkinLanguage: tc.declared}.GherkinFor()
		if lang != tc.lang || kw.Scenario != tc.scenario {
			t.Errorf("GherkinFor(%q) = %q/%q, want %q/%q", tc.declared, lang, kw.Scenario, tc.lang, tc.scenario)
		}
	}
}

func TestGherkinScenarioAlternatives(t *testing.T) {
	t.Run("DLCTI-B07: Every way to open a scenario, in every language, deduplicated and longest first", func(t *testing.T) {})
	alts := GherkinScenarioAlternatives()
	seen := map[string]bool{}
	for i, a := range alts {
		if seen[a] {
			t.Errorf("%q appears twice", a)
		}
		seen[a] = true
		if i > 0 && len(alts[i-1]) < len(a) {
			t.Errorf("%q (longer) comes after %q", a, alts[i-1])
		}
	}
	for _, want := range []string{"Scenario", "Scenario Outline", "Example", "Esquema do Cenário", "Cenario", "シナリオ"} {
		if !seen[want] {
			t.Errorf("missing %q", want)
		}
	}
	pos := func(s string) int {
		for i, a := range alts {
			if a == s {
				return i
			}
		}
		return -1
	}
	if pos("Esquema do Cenário") > pos("Cenário") {
		t.Error("the outline must be tried before the scenario it contains")
	}
}

func TestGherkinThenAlternatives(t *testing.T) {
	t.Run("DLCTI-B08: Every result keyword, in every language, sorted", func(t *testing.T) {})
	got := GherkinThenAlternatives()
	if !sort.StringsAreSorted(got) {
		t.Errorf("not sorted: %v", got)
	}
	want := map[string]bool{}
	for _, kw := range gherkinByLang {
		want[kw.Then] = true
	}
	if len(got) != len(want) {
		t.Errorf("got %d keywords, want %d: %v", len(got), len(want), got)
	}
	for _, g := range got {
		if !want[g] {
			t.Errorf("%q is not a result keyword of the table", g)
		}
	}
}

func TestDialectCompile(t *testing.T) {
	t.Run("DLCTI-B09: An empty or invalid pattern compiles to nothing", func(t *testing.T) {})
	var d Dialect
	if d.Compile("") != nil {
		t.Error("an empty pattern must compile to nil")
	}
	if d.Compile(`(?<=x)y`) != nil {
		t.Error("an invalid pattern must compile to nil")
	}
	if re := d.Compile(`\bfor\b`); re == nil || !re.MatchString("for x") {
		t.Error("a valid pattern must compile and match")
	}
}

func TestWaivedField(t *testing.T) {
	t.Run("DLCTI-B10: The opt-out is read by the YAML field name, ignoring case and spaces", func(t *testing.T) {})
	d := Dialect{OptOut: []string{" Collection_Query ", "CollectionQuery"}}
	if !d.WaivedField("collection_query") {
		t.Error("collection_query was waived (case and spaces aside)")
	}
	if d.WaivedField("exported_func") {
		t.Error("exported_func was not waived")
	}
	if (Dialect{OptOut: []string{"HTTPStatus"}}).WaivedField("http_status") {
		t.Error("the Go field name is not the YAML name and must not waive it")
	}
	if (Dialect{}).WaivedField("collection_query") {
		t.Error("with no opt-out nothing is waived")
	}
}

func TestDefaultSetPromise(t *testing.T) {
	t.Run("DLCTI-B11: The set-promise verb is recognised after a provider prefix and never inside a word", func(t *testing.T) {})
	re := (&Config{}).DialectFor().Compile(defaultSetPromise)
	if re == nil {
		t.Fatal("the default set_promise must compile")
	}
	for _, name := range []string{"listUsers", "list_users", "GetAllItems", "cognitoListDevices", "s3ListObjects", "dynamo_query_all", "listarTodos"} {
		if !re.MatchString(name) {
			t.Errorf("%q promises the whole set and was not recognised", name)
		}
	}
	for _, name := range []string{"allocateSlot", "callbackUrl", "getUser", "enlistment"} {
		if re.MatchString(name) {
			t.Errorf("%q does not promise the set and was recognised", name)
		}
	}
}

func TestDefaultSetSlice(t *testing.T) {
	t.Run("DLCTI-B12: The set-slice convention is a query verb opening the name, followed by a slice word", func(t *testing.T) {})
	re := (&Config{}).DialectFor().Compile(defaultSetSlice)
	for _, name := range []string{"listRecent", "get_first", "fetchTopItems", "buscarPrimeiros"} {
		if !re.MatchString(name) {
			t.Errorf("%q declares a slice and was not recognised", name)
		}
	}
	for _, name := range []string{"listUsers", "recentList", "getAll", "cachedListRecent"} {
		if re.MatchString(name) {
			t.Errorf("%q declares no slice and was recognised", name)
		}
	}
}

func TestDialectFamilies_everyPatternCompiles(t *testing.T) {
	t.Run("DLCTI-I01: Every pattern a family or a naming default brings compiles", func(t *testing.T) {})
	for _, fam := range KnownDialectFamilies() {
		d := (&Config{Dialect: &Dialect{Family: fam}}).DialectFor()
		pats := []string{d.ExportedFunc, d.ParamName, d.Loop, d.Cursor, d.CollectionQuery, d.HTTPStatus, d.HTTPStatusDynamic, d.SetPromise, d.SetSlice}
		pats = append(pats, d.HandlePatterns...)
		pats = append(pats, d.LogPatterns...)
		for _, p := range pats {
			if p != "" && d.Compile(p) == nil {
				t.Errorf("family %s: pattern %q does not compile, and a gate would read it as undeclared", fam, p)
			}
		}
	}
}

func TestGherkinFor_keywordsAreReadable(t *testing.T) {
	t.Run("DLCTI-I02: The keywords written for any language are among those every reader recognises", func(t *testing.T) {})
	scen := map[string]bool{}
	for _, a := range GherkinScenarioAlternatives() {
		scen[a] = true
	}
	then := map[string]bool{}
	for _, a := range GherkinThenAlternatives() {
		then[a] = true
	}
	for lang := range gherkinByLang {
		_, kw := Dialect{GherkinLanguage: lang}.GherkinFor()
		if !scen[kw.Scenario] || !scen[kw.Outline] || !then[kw.Then] {
			t.Errorf("%s: %q/%q/%q written but not recognised by the readers", lang, kw.Scenario, kw.Outline, kw.Then)
		}
	}
}

func TestDialectFamilies_noVendorQuery(t *testing.T) {
	t.Run("DLCTI-X01: No family brings a collection query, which is the project's to declare", func(t *testing.T) {})
	for _, fam := range KnownDialectFamilies() {
		if q := (&Config{Dialect: &Dialect{Family: fam}}).DialectFor().CollectionQuery; q != "" {
			t.Errorf("family %s brings collection_query %q", fam, q)
		}
	}
}

// The Go family recognises both shapes of error handling, the plain check and the one with
// an init statement.
func TestGoHandlePatternsSeeBothShapes(t *testing.T) {
	t.Run("DLCTI-B13: The Go family recognises both shapes of error handling", func(t *testing.T) {})
	d := (&Config{Dialect: &Dialect{Family: "go"}}).DialectFor()
	for _, line := range []string{"if err != nil {", "if err := os.Remove(p); err != nil {"} {
		found := false
		for _, p := range d.HandlePatterns {
			if regexp.MustCompile(p).MatchString(line) {
				found = true
			}
		}
		if !found {
			t.Errorf("the Go handle patterns must see %q", line)
		}
	}
}

// GORM keeps the error in a field and Go code returns sentinel errors; neither carries the
// word `err`, and a repository that handled and propagated both failed the failure gates.
func TestGoPatternsSeeFieldAndSentinelErrors(t *testing.T) {
	t.Run("DLCTI-B19: The Go family sees an error in a field and a sentinel error", func(t *testing.T) {})
	d := (&Config{Dialect: &Dialect{Family: "go"}}).DialectFor()
	any := func(pats []string, line string) bool {
		for _, p := range pats {
			if regexp.MustCompile(p).MatchString(line) {
				return true
			}
		}
		return false
	}
	for _, line := range []string{"if result.Error != nil {", "return ErrWalletLinkNotFound"} {
		if !any(d.HandlePatterns, line) {
			t.Errorf("the Go handle patterns must see %q", line)
		}
	}
	for _, line := range []string{"return nil, result.Error", "return nil, ErrWalletLinkNotFound"} {
		if !any(d.LogPatterns, line) {
			t.Errorf("the Go log patterns must see %q", line)
		}
	}
	if any(d.LogPatterns, "return e.Error()") {
		t.Error("turning the error into text is not propagating it")
	}
}

func TestFamiliesSayHowATestIsWritten(t *testing.T) {
	t.Run("DLCTI-B14: The Go and TS families say how a test is written", func(t *testing.T) {})
	re := func(family string) *regexp.Regexp {
		c := &Config{Dialect: &Dialect{Family: family}}
		d := c.DialectFor()
		if d.Tests == nil || d.Tests.Pattern == "" {
			t.Fatalf("family %s must say how a test is written", family)
		}
		return regexp.MustCompile(d.Tests.Pattern)
	}
	goRe := re("go")
	if !goRe.MatchString(`t.Run("X", f)`) || goRe.MatchString(`it('X', f)`) {
		t.Error("go reads t.Run and only it")
	}
	tsRe := re("ts")
	for _, call := range []string{`it('x'`, `test("x"`, "describe(`x`", `it.only('x'`, `test.skip('x'`,
		`it.each([[1, 2], [f(3), {a: [4]}]])('x'`, "it.each(\n  rows,\n)('x'"} {
		if loc := tsRe.FindStringIndex(call); loc == nil || loc[0] != 0 || call[loc[1]] != '\'' && call[loc[1]] != '"' && call[loc[1]] != '`' {
			t.Errorf("ts must read %q up to its title", call)
		}
	}
	own := &Config{Dialect: &Dialect{Family: "ts", Tests: &TestsSource{Script: "node list.mjs"}}}
	if d := own.DialectFor(); d.Tests == nil || d.Tests.Script != "node list.mjs" || d.Tests.Pattern != "" {
		t.Errorf("the project's own tests must win over the family's, got %+v", d.Tests)
	}
	if d := (&Config{Dialect: &Dialect{Family: "python"}}).DialectFor(); d.Tests != nil {
		t.Errorf("a family that does not say how a test is written leaves it undeclared, got %+v", d.Tests)
	}
}

func TestTSFamilyRecognisesCatchWithoutBinding(t *testing.T) {
	t.Run("DLCTI-B15: The TS family recognises catch with or without its binding", func(t *testing.T) {})
	d := (&Config{Dialect: &Dialect{Family: "ts"}}).DialectFor()
	handles := func(src string) bool {
		for _, p := range d.HandlePatterns {
			if regexp.MustCompile(p).MatchString(src) {
				return true
			}
		}
		return false
	}
	for src, want := range map[string]bool{
		"try { f() } catch (e) { log(e) }":  true,
		"try { f() } catch { return null }": true,
		"try { f() } catch{ return null }":  true,
		"const catchPhrase = 'x'":           false,
	} {
		if got := handles(src); got != want {
			t.Errorf("%q: want %v, got %v", src, want, got)
		}
	}
}

func TestDialectAssertionByFamily(t *testing.T) {
	t.Run("DLCTI-B16: The families say what an assertion is, and a project may declare only its own", func(t *testing.T) {})
	goD := (&Config{Dialect: &Dialect{Family: "go"}}).DialectFor()
	if goD.Tests == nil || !regexp.MustCompile(goD.Tests.Assertion).MatchString(`t.Errorf("x")`) ||
		regexp.MustCompile(goD.Tests.Assertion).MatchString(`t.Run("x")`) {
		t.Fatalf("go asserts with t.Errorf and not with t.Run, got %+v", goD.Tests)
	}
	tsD := (&Config{Dialect: &Dialect{Family: "ts"}}).DialectFor()
	for _, a := range []string{"expect(x).toBe(1)", "expect.assertions(1)", "assert(ok)", "assert.equal(a, b)"} {
		if !regexp.MustCompile(tsD.Tests.Assertion).MatchString(a) {
			t.Errorf("ts must read %q as an assertion", a)
		}
	}
	onlyAssertion := (&Config{Dialect: &Dialect{Family: "go", Tests: &TestsSource{Assertion: `\bcheck\(`}}}).DialectFor()
	if onlyAssertion.Tests.Pattern != goD.Tests.Pattern || onlyAssertion.Tests.Assertion != `\bcheck\(` {
		t.Errorf("declaring only the assertion keeps the family's pattern and the project's assertion, got %+v", onlyAssertion.Tests)
	}
	if goD.Tests.Assertion == "" || (&Config{Dialect: &Dialect{Family: "go"}}).DialectFor().Tests.Assertion != goD.Tests.Assertion {
		t.Error("resolving must not change the family's own assertion")
	}
	ownPattern := (&Config{Dialect: &Dialect{Family: "go", Tests: &TestsSource{Pattern: `\bSpec\(`}}}).DialectFor()
	if ownPattern.Tests.Assertion != "" || ownPattern.Tests.Pattern != `\bSpec\(` {
		t.Errorf("a project reading its tests its own way takes no assertion from the family, got %+v", ownPattern.Tests)
	}
	ownScript := (&Config{Dialect: &Dialect{Family: "ts", Tests: &TestsSource{Script: "node l.mjs"}}}).DialectFor()
	if ownScript.Tests.Assertion != "" || ownScript.Tests.Pattern != "" {
		t.Errorf("a project with its own script takes nothing from the family, got %+v", ownScript.Tests)
	}
}

func TestDialectDefinitionByFamily(t *testing.T) {
	t.Run("DLCTI-B17: The families say how code defines a name", func(t *testing.T) {})
	names := func(family, code string) []string {
		re := regexp.MustCompile((&Config{Dialect: &Dialect{Family: family}}).DialectFor().Definition)
		owner := re.SubexpIndex("owner")
		var out []string
		for _, m := range re.FindAllStringSubmatch(code, -1) {
			for i, g := range m[1:] {
				if g != "" && i+1 != owner {
					out = append(out, g)
					break
				}
			}
		}
		return out
	}
	for _, c := range []struct {
		family, code string
		want         []string
	}{
		{"go", "func Build() {}\nfunc (g *Graph) AddFiles(x int) {}\ntype Graph struct{}\nfunc Map[T any]() {}\n\tfunc inner() {}", []string{"Build", "AddFiles", "Graph", "Map"}},
		{"ts", "export function calc(a) {}\nexport default async function main() {}\nclass Cart {}\n  const total = (a) => a\nexport const fetchAll = async (x) => x\nconst n = 3\nfunction* gen() {}", []string{"calc", "main", "Cart", "total", "fetchAll", "gen"}},
		{"python", "def pay(x):\n  async def inner():\nclass Account:", []string{"pay", "inner", "Account"}},
	} {
		if got := names(c.family, c.code); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: want %v, got %v", c.family, c.want, got)
		}
	}
	// A Go method names its receiver's type in the `owner` group.
	goRe := regexp.MustCompile((&Config{Dialect: &Dialect{Family: "go"}}).DialectFor().Definition)
	if m := goRe.FindStringSubmatch("func (f *fakePinger) Ping(ctx context.Context) error {"); m == nil || m[goRe.SubexpIndex("owner")] != "fakePinger" {
		t.Errorf("the go definition should capture the receiver as owner, got %q", m)
	}
	own := (&Config{Dialect: &Dialect{Family: "go", Definition: `(?m)^proc (\w+)`}}).DialectFor()
	if own.Definition != `(?m)^proc (\w+)` {
		t.Errorf("the project's definition wins, got %q", own.Definition)
	}
	if d := (&Config{Dialect: &Dialect{Family: "rust"}}).DialectFor(); d.Definition != "" {
		t.Errorf("a family with no definition leaves it undeclared, got %q", d.Definition)
	}
}

func TestGherkinExamplesAlternatives(t *testing.T) {
	t.Run("DLCTI-B18: Every examples keyword, in every language, sorted", func(t *testing.T) {})
	got := GherkinExamplesAlternatives()
	if !sort.StringsAreSorted(got) {
		t.Errorf("sorted, got %v", got)
	}
	seen := map[string]int{}
	for _, k := range got {
		seen[k]++
	}
	for _, want := range []string{"Examples", "Exemplos", "Ejemplos", "Beispiele", "例"} {
		if seen[want] != 1 {
			t.Errorf("%q once, got %d in %v", want, seen[want], got)
		}
	}
	if len(seen) != len(got) {
		t.Errorf("each keyword once, got %v", got)
	}
}

func TestEnvReadPattern(t *testing.T) {
	t.Run("DLCTI-B20: Each family reads environment variables its own way", func(t *testing.T) {})
	reads := map[string]string{
		"go": `os.Getenv("A")`, "ts": `process.env.A`, "python": `os.environ["A"]`, "java": `System.getenv("A")`,
		"kotlin": `System.getenv("A")`, "rust": `env::var("A")`, "ruby": `ENV["A"]`, "php": `getenv('A')`,
		"csharp": `Environment.GetEnvironmentVariable("A")`,
	}
	name := func(pat, text string) string {
		m := regexp.MustCompile(pat).FindStringSubmatch(text)
		for _, g := range m[1:] {
			if g != "" {
				return g
			}
		}
		return ""
	}
	union := (Dialect{}).EnvReadPattern()
	for fam, text := range reads {
		pat := (&Config{Dialect: &Dialect{Family: fam}}).DialectFor().EnvReadPattern()
		if got := name(pat, text); got != "A" {
			t.Errorf("%s reads %q as %q", fam, text, got)
		}
		if got := name(union, text); got != "A" {
			t.Errorf("with no family, %q is read as %q", text, got)
		}
	}
	if p := (Dialect{EnvRead: `cfg\("(\w+)"\)`}).EnvReadPattern(); p != `cfg\("(\w+)"\)` {
		t.Errorf("the project's own pattern wins: %s", p)
	}
}

func TestFalliblePatterns_familyFillsWhatTheProjectDidNotDeclare(t *testing.T) {
	t.Run("DLCTI-B21: The fallible patterns are the project's, and the family's only when the project declares none", func(t *testing.T) {})
	own := []FalliblePattern{{Call: `\bload\(`, Handled: `err`, Window: 5}}
	c := &Config{Dialect: &Dialect{FalliblePatterns: own}}
	if got := c.DialectFor().FalliblePatterns; len(got) != 1 || got[0].Call != `\bload\(` || got[0].Window != 5 {
		t.Errorf("the project's patterns stay: %+v", got)
	}
	if got := (&Config{}).DialectFor().FalliblePatterns; len(got) != 0 {
		t.Errorf("no family, no project patterns, nothing: %+v", got)
	}
	if DefaultFallibleWindow <= 0 {
		t.Error("the default window is a positive number of lines")
	}
}

func TestFalliblePatterns_familyKnowsThemAndDoesNotImposeThem(t *testing.T) {
	t.Run("DLCTI-B22: A family knows its fallible patterns and does not impose them on a project that declared none", func(t *testing.T) {})
	found := false
	for _, p := range FamilyFalliblePatterns("TS") {
		if strings.Contains(p.Call, "Query") {
			found = true
		}
	}
	if !found {
		t.Errorf("the ts family knows React Query: %+v", FamilyFalliblePatterns("ts"))
	}
	for _, fam := range []string{"go", "python"} {
		if len(FamilyFalliblePatterns(fam)) == 0 {
			t.Errorf("%s knows its fallible patterns", fam)
		}
	}
	if FamilyFalliblePatterns("cobol") != nil {
		t.Error("an unknown family knows none")
	}
	if got := (&Config{Dialect: &Dialect{Family: "ts"}}).DialectFor().FalliblePatterns; len(got) != 0 {
		t.Errorf("a project that declared none is measured with none until it adopts them: %+v", got)
	}
}

func TestDialect_importResolveAndNavigationCall(t *testing.T) {
	t.Run("DLCTI-B23: import_resolve and navigation_call are read from the dialect", func(t *testing.T) {})
	c, err := load(t, "version: 7\nlayers: {}\ndialect:\n  import_resolve:\n    extensions: [.ts, /index.ts]\n    aliases:\n      \"@/\": \"src/\"\n  navigation_call: 'navigate\\(''(?P<route>\\w+)'''\n")
	if err != nil {
		t.Fatal(err)
	}
	d := c.DialectFor()
	if d.ImportResolve == nil || len(d.ImportResolve.Extensions) != 2 || d.ImportResolve.Aliases["@/"] != "src/" || !strings.Contains(d.NavigationCall, "route") {
		t.Errorf("got %+v %q", d.ImportResolve, d.NavigationCall)
	}
}
