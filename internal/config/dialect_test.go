// @anchors
//   code: DLTSA
//   ref: DLCTI

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// example is the family's example lexicon, or nothing for a family with none.
func example(family string) Dialect {
	if d := FamilyLexicon(family); d != nil {
		return *d
	}
	return Dialect{}
}

func TestDialectFor_onlyWhatIsDeclared(t *testing.T) {
	t.Run("DLCTI-B01: A project's dialect is what it declares, and nothing fills it", func(t *testing.T) {})
	for name, c := range map[string]*Config{"nil config": nil, "no dialect block": {}, "a family alone": {Dialect: &Dialect{Family: "go"}}} {
		d := c.DialectFor()
		if d.SetPromise != "" || d.SetSlice != "" || d.ExportedFunc != "" || d.Loop != "" || d.Cursor != "" ||
			len(d.HandlePatterns) != 0 || d.Tests != nil || d.Definition != "" || d.EnvRead != "" {
			t.Errorf("%s: nothing comes from a family or a default, got %+v", name, d)
		}
	}
	c := &Config{Dialect: &Dialect{Family: "go", Loop: `\bloop\b`}}
	if d := c.DialectFor(); d.Loop != `\bloop\b` || d.ExportedFunc != "" {
		t.Errorf("the declared fields are read as declared, got %+v", d)
	}
}

func TestFamilyLexicon_isAnExample(t *testing.T) {
	t.Run("DLCTI-B02: Each family has an example of the lexicon, never applied by itself", func(t *testing.T) {})
	ex := FamilyLexicon("go")
	base := dialectFamilies["go"]
	if ex == nil || ex.ExportedFunc != base.ExportedFunc || ex.Cursor != base.Cursor || ex.ParamName != base.ParamName ||
		ex.HTTPStatus != base.HTTPStatus || ex.HTTPStatusDynamic != base.HTTPStatusDynamic ||
		!reflect.DeepEqual(ex.LogPatterns, base.LogPatterns) || ex.StatementEnd == "" || ex.SignatureEnd == "" ||
		len(ex.Limits) == 0 || len(ex.GuardPatterns) == 0 {
		t.Errorf("the example carries the family's lexicon, got %+v", ex)
	}
	ex.HandlePatterns[0] = "changed"
	if FamilyLexicon("go").HandlePatterns[0] == "changed" {
		t.Error("an example handed out is a copy: changing it changes no family")
	}
}

func TestFamilyLexicon_ignoresCase(t *testing.T) {
	t.Run("DLCTI-B03: The family name is matched ignoring case", func(t *testing.T) {})
	if ex := FamilyLexicon("TS"); ex == nil || ex.ExportedFunc != dialectFamilies["ts"].ExportedFunc {
		t.Errorf("family TS must resolve to ts, got %+v", ex)
	}
}

func TestFamilyLexicon_unknownFamilyHasNone(t *testing.T) {
	t.Run("DLCTI-B04: An unknown family has no example, and the known ones are listed in order", func(t *testing.T) {})
	if ex := FamilyLexicon("cobol"); ex != nil {
		t.Errorf("an unknown family has no example, got %+v", ex)
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

func TestFamilyLexicon_namingConventions(t *testing.T) {
	t.Run("DLCTI-B05: Every family's example carries the two naming conventions", func(t *testing.T) {})
	for _, fam := range KnownDialectFamilies() {
		if ex := FamilyLexicon(fam); ex.SetPromise != defaultSetPromise || ex.SetSlice != defaultSetSlice {
			t.Errorf("%s: the example names the set conventions, got %q / %q", fam, ex.SetPromise, ex.SetSlice)
		}
	}
	if d := (&Config{Dialect: &Dialect{SetSlice: `^buscarAlguns`}}).DialectFor(); d.SetSlice != `^buscarAlguns` || d.SetPromise != "" {
		t.Errorf("a project declares its own conventions, or none, got %q / %q", d.SetSlice, d.SetPromise)
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
	re := Dialect{}.Compile(defaultSetPromise)
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
	re := Dialect{}.Compile(defaultSetSlice)
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
		d := example(fam)
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
		if q := example(fam).CollectionQuery; q != "" {
			t.Errorf("family %s brings collection_query %q", fam, q)
		}
	}
}

// The Go family recognises both shapes of error handling, the plain check and the one with
// an init statement.
func TestGoHandlePatternsSeeBothShapes(t *testing.T) {
	t.Run("DLCTI-B13: The Go family recognises both shapes of error handling", func(t *testing.T) {})
	d := example("go")
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
	d := example("go")
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
		d := example(family)
		if d.Tests == nil || d.Tests.Pattern == "" {
			t.Fatalf("family %s must say how a test is written", family)
		}
		return regexp.MustCompile(d.Tests.Pattern)
	}
	// title is what a `title` group of the pattern captured.
	title := func(re *regexp.Regexp, s string) string {
		m := re.FindStringSubmatch(s)
		for i, n := range re.SubexpNames() {
			if m != nil && strings.HasPrefix(n, "title") && m[i] != "" {
				return m[i]
			}
		}
		return ""
	}
	goRe := re("go")
	if title(goRe, `t.Run("X", f)`) != "X" || goRe.MatchString(`it('X', f)`) {
		t.Error("go reads t.Run and only it")
	}
	tsRe := re("ts")
	for _, call := range []string{`it('x'`, `test("x"`, "describe(`x`", `it.only('x'`, `test.skip('x'`,
		`it.each([[1, 2], [f(3), {a: [4]}]])('x'`, "it.each(\n  rows,\n)('x'"} {
		if title(tsRe, call+"x"+string(call[len(call)-1])+", f)") != "x" {
			t.Errorf("ts must read %q and its title", call)
		}
	}
	own := &Config{Dialect: &Dialect{Family: "ts", Tests: &TestsSource{Script: "node list.mjs"}}}
	if d := own.DialectFor(); d.Tests == nil || d.Tests.Script != "node list.mjs" || d.Tests.Pattern != "" || d.Tests.Assertion != "" {
		t.Errorf("the project's own tests are read as declared, got %+v", d.Tests)
	}
	if d := example("python"); d.Tests != nil {
		t.Errorf("a family that does not say how a test is written has no example of it, got %+v", d.Tests)
	}
}

func TestTSFamilyRecognisesCatchWithoutBinding(t *testing.T) {
	t.Run("DLCTI-B15: The TS family recognises catch with or without its binding", func(t *testing.T) {})
	d := example("ts")
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
	goD := example("go")
	if goD.Tests == nil || !regexp.MustCompile(goD.Tests.Assertion).MatchString(`t.Errorf("x")`) ||
		regexp.MustCompile(goD.Tests.Assertion).MatchString(`t.Run("x")`) {
		t.Fatalf("go asserts with t.Errorf and not with t.Run, got %+v", goD.Tests)
	}
	tsD := example("ts")
	for _, a := range []string{"expect(x).toBe(1)", "expect.assertions(1)", "assert(ok)", "assert.equal(a, b)"} {
		if !regexp.MustCompile(tsD.Tests.Assertion).MatchString(a) {
			t.Errorf("ts must read %q as an assertion", a)
		}
	}
	onlyAssertion := (&Config{Dialect: &Dialect{Family: "go", Tests: &TestsSource{Assertion: `\bcheck\(`}}}).DialectFor()
	if onlyAssertion.Tests.Pattern != "" || onlyAssertion.Tests.Assertion != `\bcheck\(` {
		t.Errorf("a project declaring only its assertion has only its assertion, got %+v", onlyAssertion.Tests)
	}
	if example("go").Tests.Assertion != goD.Tests.Assertion {
		t.Error("handing out the example must not change the family's own assertion")
	}
}

func TestDialectDefinitionByFamily(t *testing.T) {
	t.Run("DLCTI-B17: The families say how code defines a name", func(t *testing.T) {})
	names := func(family, code string) []string {
		re := regexp.MustCompile(example(family).Definition)
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
	goRe := regexp.MustCompile(example("go").Definition)
	if m := goRe.FindStringSubmatch("func (f *fakePinger) Ping(ctx context.Context) error {"); m == nil || m[goRe.SubexpIndex("owner")] != "fakePinger" {
		t.Errorf("the go definition should capture the receiver as owner, got %q", m)
	}
	own := (&Config{Dialect: &Dialect{Family: "go", Definition: `(?m)^proc (\w+)`}}).DialectFor()
	if own.Definition != `(?m)^proc (\w+)` {
		t.Errorf("the project's definition is read as declared, got %q", own.Definition)
	}
	if d := example("rust"); d.Definition != "" {
		t.Errorf("a family with no definition has no example of it, got %q", d.Definition)
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
	for fam, text := range reads {
		pat := example(fam).EnvReadPattern()
		if got := name(pat, text); got != "A" {
			t.Errorf("%s reads %q as %q", fam, text, got)
		}
	}
	if p := (Dialect{}).EnvReadPattern(); p != "" {
		t.Errorf("with no env_read nothing is read, got %s", p)
	}
	if p := (Dialect{EnvRead: `cfg\("(\w+)"\)`}).EnvReadPattern(); p != `cfg\("(\w+)"\)` {
		t.Errorf("the project's own pattern is read: %s", p)
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
	// Every pattern of every family says the function a failure cites, in its name group,
	// and the ts ones take a type argument.
	for _, fam := range []string{"ts", "go", "python"} {
		for _, p := range FamilyFalliblePatterns(fam) {
			if re := regexp.MustCompile(p.Call); re.SubexpIndex("name") < 0 {
				t.Errorf("%s: %q has no name group", fam, p.Call)
			}
		}
	}
	for _, p := range FamilyFalliblePatterns("ts") {
		re := regexp.MustCompile(p.Call)
		if m := re.FindStringSubmatch("const q = useQuery<Budget>(key)"); m != nil && m[re.SubexpIndex("name")] != "useQuery" {
			t.Errorf("useQuery<Budget>( is named useQuery: %v", m)
		}
	}
	matched := false
	for _, p := range FamilyFalliblePatterns("ts") {
		if regexp.MustCompile(p.Call).MatchString("useQuery<Budget>(key)") {
			matched = true
		}
	}
	if !matched {
		t.Error("the ts family recognises a generic useQuery")
	}
	if FamilyFalliblePatterns("cobol") != nil {
		t.Error("an unknown family knows none")
	}
	if got := example("ts").FalliblePatterns; len(got) != 0 {
		t.Errorf("a project that declared none is measured with none until it adopts them: %+v", got)
	}
}

func TestFalliblePatterns_familySaysWhatPropagates(t *testing.T) {
	t.Run("DLCTI-B26: A family's fallible patterns say when a call propagates", func(t *testing.T) {})
	propagates := func(fam, text string) bool {
		for _, p := range FamilyFalliblePatterns(fam) {
			re := regexp.MustCompile(p.Call)
			if m := re.FindStringSubmatchIndex(text); m != nil {
				i := re.SubexpIndex("propagates")
				return i >= 0 && m[2*i] >= 0
			}
		}
		t.Fatalf("%s: no pattern matches %q", fam, text)
		return false
	}
	for _, c := range []struct {
		fam, text string
		want      bool
	}{
		{"ts", "  return useQuery<Budget>(key)", true},
		{"ts", "const load = async () => await fetch(url)", true},
		{"go", "\treturn http.Get(url)", true},
		{"python", "    return requests.get(url)", true},
		{"ts", "const q = useQuery(key)", false},
		{"go", "resp, err := http.Get(url)", false},
	} {
		if got := propagates(c.fam, c.text); got != c.want {
			t.Errorf("%s %q: propagates=%v, want %v", c.fam, c.text, got, c.want)
		}
	}
}

func TestDialect_importResolveAndNavigationCall(t *testing.T) {
	t.Run("DLCTI-B23: import_resolve and navigation_call are read from the dialect", func(t *testing.T) {})
	c, err := load(t, "version: 8\nlayers: {}\ndialect:\n  import_resolve:\n    extensions: [.ts, /index.ts]\n    aliases:\n      \"@/\": \"src/\"\n  navigation_call: 'navigate\\(''(?P<route>\\w+)'''\n")
	if err != nil {
		t.Fatal(err)
	}
	d := c.DialectFor()
	if d.ImportResolve == nil || len(d.ImportResolve.Extensions) != 2 || d.ImportResolve.Aliases["@/"] != "src/" || !strings.Contains(d.NavigationCall, "route") {
		t.Errorf("got %+v %q", d.ImportResolve, d.NavigationCall)
	}
	sep, err := load(t, "version: 8\nlayers: {}\ndialect:\n  import_resolve:\n    separator: '.'\n")
	if err != nil || sep.DialectFor().ImportResolve.Separator != "." {
		t.Errorf("the separator is read: %v", err)
	}
	if got := example("ts").NavigationCall; got != "" || !strings.Contains(FamilyNavigationCall("ts"), "(?P<route>") || FamilyNavigationCall("go") != "" {
		t.Errorf("no family gives a navigation call by itself, and ts has an example: %q", got)
	}
}

func TestDialectFor_ocaml(t *testing.T) {
	t.Run("DLCTI-B24: The OCaml family reads OCaml", func(t *testing.T) {})
	d := example("OCaml")
	src := "open Lexer\nopen! Base.Option\nlet parse x =\n  let local = x in local\nlet rec walk' t = t\nlet home = Sys.getenv_opt \"HOME\"\n"
	var exported []string
	for _, m := range regexp.MustCompile(d.ExportedFunc).FindAllStringSubmatch(src, -1) {
		exported = append(exported, m[1])
	}
	if !reflect.DeepEqual(exported, []string{"parse", "walk'", "home"}) {
		t.Errorf("exported %v", exported)
	}
	if m := regexp.MustCompile(d.EnvRead).FindStringSubmatch(src); m == nil || m[1] != "HOME" {
		t.Errorf("env read: %v", m)
	}
	if !regexp.MustCompile(d.Loop).MatchString("for i = 1 to 3 do () done") {
		t.Error("the loop is OCaml's")
	}
	var opens []string
	for _, m := range FamilyChain("ocaml").Imports.FindAll(src, "src/p.ml") {
		opens = append(opens, m.Groups["path"].Value)
	}
	if !reflect.DeepEqual(opens, []string{"Lexer", "Base", "Base", "Sys"}) { // opened, then qualifying Option: one import, deduplicated by the gate
		t.Errorf("the example reads the modules opened and qualified: %v", opens)
	}
	if len(d.Imports) != 0 {
		t.Error("the family's import patterns are an example, never applied by themselves")
	}
}

func TestImportResolve_byName(t *testing.T) {
	t.Run("DLCTI-B25: The import resolution declares imports by name", func(t *testing.T) {})
	path := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(path, []byte("version: 8\ndialect:\n  family: ocaml\n  import_resolve:\n    by_name: true\n    name_case: lower_first\n    extensions: [.ml]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if r := c.DialectFor().ImportResolve; r == nil || !r.ByName || r.NameCase != "lower_first" || r.Extensions[0] != ".ml" {
		t.Errorf("resolution: %+v", r)
	}
}

func TestDialect_testHandles(t *testing.T) {
	t.Run("DLCTI-B27: The test handle patterns are read and checked on load", func(t *testing.T) {})
	c, err := load(t, "version: 8\nlayers: {}\ndialect:\n  test_handles:\n    exposed: ['testID=\"(?P<testid>[a-z-]+)\"']\n    mentioned: ['-(?P<suffix>row)']\n    queried: ['id: (?P<prefix>[a-z-]+)']\n")
	if err != nil {
		t.Fatal(err)
	}
	if th := c.Dialect.TestHandles; th == nil || len(th.Exposed) != 1 || len(th.Mentioned) != 1 || len(th.Queried) != 1 {
		t.Errorf("the three lists are read: %+v", c.Dialect.TestHandles)
	}
	for _, c := range []struct{ yaml, names string }{
		{"exposed: ['(?P<testid>[a-z']", "test_handles.exposed[0]` does not compile"},
		{"queried: ['id: ([a-z]+)']", "test_handles.queried[0]` captures with none of the groups testid, prefix"},
		{"exposed: ['-(?P<suffix>row)']", "test_handles.exposed[0]` captures with none"},
	} {
		if _, err := load(t, "version: 8\nlayers: {}\ndialect:\n  test_handles:\n    "+c.yaml+"\n"); err == nil || !strings.Contains(err.Error(), c.names) {
			t.Errorf("%s: the config is refused naming the pattern, got %v", c.yaml, err)
		}
	}
}

func TestFamilyTestHandles_anExampleNeverApplied(t *testing.T) {
	t.Run("DLCTI-B28: The ts family has an example of the test handle patterns, never applied by itself", func(t *testing.T) {})
	ex := FamilyTestHandles("TS", "data-testid")
	if ex == nil || len(ex.Exposed) == 0 || len(ex.Queried) == 0 || !strings.Contains(ex.Exposed[0].Re, "data-testid") {
		t.Fatalf("the ts example reads the project's attribute: %+v", ex)
	}
	if !regexp.MustCompile(ex.Exposed[0].Re).MatchString(`<div data-testid="abcd-root">`) {
		t.Errorf("the example finds a data-testid: %s", ex.Exposed[0].Re)
	}
	if FamilyTestHandles("go", "testID") != nil || FamilyTestHandles("ts", " ") != nil {
		t.Error("another family, or no attribute, gives no example")
	}
	if got := example("ts").TestHandles; got != nil {
		t.Errorf("a project that declared none is measured with none: %+v", got)
	}
	// Reported from MIF: a conditional nested three deep, one branch a template. Every branch
	// is read; the condition's literal (`'email'`) is not.
	src := "testID={\n  i === 0\n    ? ':import-source-gallery'\n    : i === 1\n      ? ':import-source-camera'\n      : kind === 'email'\n        ? `:urpt-row-${i}-delete`\n        : ':import-manual-entry-button'\n}"
	var got []string
	for _, m := range FamilyTestHandles("ts", "testID").Exposed.FindAll(src, "x.tsx") {
		for _, g := range []string{"testid", "prefix"} {
			if c, ok := m.Groups[g]; ok {
				got = append(got, g+"="+c.Value)
			}
		}
	}
	if want := "testid=:import-source-gallery,testid=:import-source-camera,prefix=:urpt-row-,testid=:import-manual-entry-button"; strings.Join(got, ",") != want {
		t.Errorf("every branch of the nested conditional, none of its condition:\n got %v\nwant %s", got, want)
	}
}

func TestPattern_engine(t *testing.T) {
	t.Run("DLCTI-B29: A dialect pattern runs over the whole content, as an expression or a region", func(t *testing.T) {})
	src := "a\nimport {\n  One,\n  Two as Deux,\n  Three\n} from './x'\n"
	across := Pattern{Re: `import\s*\{[^}]*\}\s*from\s*'(?P<path>[^']+)'`}
	ms := across.FindAll(src)
	if len(ms) != 1 || LineOf(src, ms[0].Start) != 2 || LineOf(src, ms[0].Groups["path"].Start) != 6 || ms[0].Listed {
		t.Fatalf("one match from line 2, its path on line 6: %+v", ms)
	}
	region := Pattern{Within: `\{(?P<within>[^}]*)\}`, Each: `(?P<name>\w+)\s*(?:,|$)`}
	var names []string
	for _, m := range region.FindAll(src) {
		if !m.Listed || m.Region != ms[0].Start+len("import ") {
			t.Errorf("the region's names share its start and are listed: %+v", m)
		}
		names = append(names, m.Groups["name"].Value+"@"+strconv.Itoa(LineOf(src, m.Start)))
	}
	if strings.Join(names, ",") != "One@3,Deux@4,Three@5" {
		t.Errorf("each name of the region, on its line: %v", names)
	}
	goOnly := Patterns{{Re: `(?P<path>x)`, Files: "**/*.go"}}
	if len(goOnly.FindAll("x", "a/b.ts")) != 0 || len(goOnly.FindAll("x", "a/b.go")) != 1 {
		t.Error("files narrows the pattern to its glob")
	}
	var loaded struct {
		P Patterns `yaml:"p"`
	}
	if err := yaml.Unmarshal([]byte("p:\n  - '(?P<path>a)'\n  - within: '(?P<within>b)'\n    each: '(?P<path>c)'\n    files: '**/*.go'\n"), &loaded); err != nil ||
		len(loaded.P) != 2 || loaded.P[0].Re != "(?P<path>a)" || loaded.P[1].Each != "(?P<path>c)" || loaded.P[1].Files != "**/*.go" {
		t.Errorf("written as text or as a map: %v %+v", err, loaded.P)
	}
	if out, _ := yaml.Marshal(loaded.P); !strings.HasPrefix(string(out), "- (?P<path>a)\n- within:") {
		t.Errorf("a pattern that is only an expression is written as text:\n%s", out)
	}
	for _, c := range []struct {
		p    Pattern
		says string
	}{
		{Pattern{Re: `(?P<path>[`}, "does not compile"},
		{Pattern{Within: `(?P<within>x)`}, "a region needs both"},
		{Pattern{Within: `(x)`, Each: `(?P<path>x)`}, "has no `within` group"},
		{Pattern{Re: `(x)`}, "captures with none of the groups path"},
		{Pattern{Re: `(?P<path>x)`, Files: "[a"}, "not a valid glob"},
	} {
		if err := (Patterns{c.p}).Validate("dialect.imports", []string{"path"}); err == nil || !strings.Contains(err.Error(), c.says) || !strings.Contains(err.Error(), "dialect.imports[0]") {
			t.Errorf("%+v: refused naming it (%q), got %v", c.p, c.says, err)
		}
	}
}

func TestDialect_chainPatterns(t *testing.T) {
	t.Run("DLCTI-B30: The chain patterns are the dialect's, and a family gives none by itself", func(t *testing.T) {})
	c, err := load(t, "version: 8\nlayers: {}\ndialect:\n  imports: ['from (?P<path>\\S+)']\n  import_names: ['import (?P<name>\\w+)']\n  declarations: ['def (?P<name>\\w+)', 'export default(?P<default>)']\n")
	if err != nil {
		t.Fatal(err)
	}
	if d := c.DialectFor(); len(d.Imports) != 1 || len(d.ImportNames) != 1 || len(d.Declarations) != 2 {
		t.Errorf("the three lists are read: %+v", d)
	}
	if _, err := load(t, "version: 8\nlayers: {}\ndialect:\n  imports: ['from (\\S+)']\n"); err == nil || !strings.Contains(err.Error(), "dialect.imports[0]") {
		t.Errorf("an import pattern with no path group refuses the config: %v", err)
	}
	if d := example("ts"); len(d.Imports) != 0 || d.ImportResolve != nil || d.TypeImportPattern != "" {
		t.Errorf("the family gives no import pattern nor resolution by itself: %+v", d)
	}
}

func TestFamilyChain_examples(t *testing.T) {
	t.Run("DLCTI-B31: The families have an example of the chain patterns", func(t *testing.T) {})
	paths := func(family, src, rel string) []string {
		var out []string
		for _, m := range FamilyChain(family).Imports.FindAll(src, rel) {
			out = append(out, m.Groups["path"].Value)
		}
		return out
	}
	ts := "import {\n  A,\n  B,\n} from './ab'\nimport './polyfills'\nconst { c } = await import('./c')\nconst x = 'from'\n"
	if got := paths("ts", ts, "src/x.ts"); strings.Join(got, ",") != "./ab,./polyfills,./c" {
		t.Errorf("ts: %v", got)
	}
	if got := paths("ts", "- runFlow:\n    file: login.yaml\n", "e2e/a.yaml"); strings.Join(got, ",") != "login.yaml" {
		t.Errorf("ts flows: %v", got)
	}
	goSrc := "package x\n\nimport \"fmt\"\n\nimport (\n\t\"os\"\n\tx \"example.com/m/x\"\n)\n\nvar s = []string{\n\t\"not\",\n}\n"
	if got := paths("go", goSrc, "a.go"); strings.Join(got, ",") != "fmt,os,example.com/m/x" {
		t.Errorf("go: %v", got)
	}
	if got := paths("python", "from app.core import (\n    a,\n    b,\n)\nimport os\n", "a.py"); strings.Join(got, ",") != "app.core,os" {
		t.Errorf("python: %v", got)
	}
	if got := FamilyChain("python").NamesBrought("from app.core import (\n    a,\n    b,\n)", "a.py"); strings.Join(got, ",") != "a,b" {
		t.Errorf("python names: %v", got)
	}
	if got := paths("ocaml", "open Lexer\nlet x = Eval.apply 1\n", "a.ml"); strings.Join(got, ",") != "Lexer,Eval" {
		t.Errorf("ocaml: %v", got)
	}
	if FamilyChain("cobol") != nil {
		t.Error("an unknown family has no example")
	}
}

func TestFamilyMocksAndRoutePattern(t *testing.T) {
	t.Run("DLCTI-B32: The families' examples of the mock shapes and the route pattern", func(t *testing.T) {})
	ts := FamilyMocks("ts")
	if ts == nil || !regexp.MustCompile(ts.MockDetect).MatchString(`jest.mock('./x'`) ||
		!regexp.MustCompile(ts.MockFactory).MatchString(`jest.mock('./x', () => ({}))`) || regexp.MustCompile(ts.MockFactory).MatchString(`jest.mock('./x')`) ||
		!regexp.MustCompile(ts.MockTyped).MatchString(`(): Partial<X> => ({})`) {
		t.Errorf("the ts mocks: %+v", ts)
	}
	if py := FamilyMocks("python"); py == nil || !regexp.MustCompile(py.MockDetect).MatchString(`patch('app.x')`) || py.MockFactory != "" {
		t.Errorf("the python mocks: %+v", py)
	}
	if FamilyMocks("go") != nil || FamilyRoutePattern("go") != "" {
		t.Error("go has no example of either")
	}
	route := regexp.MustCompile(FamilyRoutePattern("ts"))
	for _, src := range []string{`<Stack.Screen name="Perfil" />`, "type P = {\n  Ajustes: undefined\n}", `addResource('signup')`} {
		if !route.MatchString(src) {
			t.Errorf("the ts route pattern reads %q", src)
		}
	}
}

func TestImportResolve_scopedAliases(t *testing.T) {
	t.Run("DLCTI-B33: Aliases scoped to the importers a glob matches", func(t *testing.T) {})
	r := &ImportResolve{
		Aliases: map[string]string{"@/": "apps/mobile/", "@backend/": "packages/backend/"},
		Scoped: []ScopedAliases{
			{Files: "apps/landing-page/**", Aliases: map[string]string{"@/": "apps/landing-page/src/"}},
			{Files: "apps/**", Aliases: map[string]string{"@/": "never/"}},
		},
	}
	if got := r.AliasesFor("apps/mobile/src/x.ts"); got["@/"] != "never/" {
		t.Errorf("a scoped entry matching the importer wins over the global: %v", got)
	}
	got := r.AliasesFor("apps/landing-page/src/a/b.test.tsx")
	if got["@/"] != "apps/landing-page/src/" || got["@backend/"] != "packages/backend/" {
		t.Errorf("the first matching entry wins, the global ones kept: %v", got)
	}
	if got := r.AliasesFor("packages/backend/x.ts"); got["@/"] != "apps/mobile/" {
		t.Errorf("no entry matches: the global aliases: %v", got)
	}
	if (*ImportResolve)(nil).AliasesFor("x") != nil {
		t.Error("no resolution, no aliases")
	}
	for _, bad := range []string{
		"  import_resolve:\n    scoped:\n      - aliases: {\"@/\": x/}\n",
		"  import_resolve:\n    scoped:\n      - files: 'apps/**'\n",
	} {
		p := filepath.Join(t.TempDir(), "anchors.yaml")
		if err := os.WriteFile(p, []byte("version: 8\ndialect:\n"+bad), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(p); err == nil || !strings.Contains(err.Error(), "import_resolve.scoped[0]") {
			t.Errorf("refused naming the entry: %v", err)
		}
	}
}
