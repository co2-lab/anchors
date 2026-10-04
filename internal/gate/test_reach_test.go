// @anchors
//   code: TRTTS
//   ref: TSRCH

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const payUnit = "export function charge(amount) { return amount }\nexport const refund = (x) => x\nfunction ok() {}\n"

// reachCfg derives `src/pay.test.ts` from `src/pay.ts`, in the TS family.
func reachCfg(invocations ...string) *config.Config {
	return &config.Config{
		Dialect: &config.Dialect{Family: "ts"},
		Derived: &config.Derived{Anchor: "code", Files: map[string]config.Padroes{"test": {"{{dir}}/{{name}}.test.ts"}}},
		Gates:   []config.Gate{{Name: "r", Check: "test-ref-matches-unit", On: []string{"test"}, Invocations: invocations}},
	}
}

// exercise writes the unit and the test and confronts the test with test-exercises-unit.
func exercise(t *testing.T, cfg *config.Config, test string) (Verdict, string) {
	t.Helper()
	resetTestedUnits()
	root := t.TempDir()
	writeFile(t, root, "src/pay.ts", payUnit)
	writeFile(t, root, "src/pay.test.ts", test)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "src/pay.ts", Kind: mapx.KindCode}, {ID: "src/pay.test.ts", Kind: mapx.KindTest}}}
	return checkTestExercisesUnit(test, g.Nodes[1], root, g, cfg)
}

func TestTestReach_namingWhatTheUnitDefines(t *testing.T) {
	t.Run("TSRCH-B01: A test that names what its unit defines exercises it", func(t *testing.T) {})
	if v, msg := exercise(t, reachCfg(), "test('a', () => expect(charge(1)).toBe(1))\n"); v != Pass {
		t.Errorf("a test calling the unit's function passes, got %v: %s", v, msg)
	}
	v, msg := exercise(t, reachCfg(), "test('a', () => expect(recharge(1)).toBe(1))\n")
	if v != Fail || !strings.Contains(msg, "`src/pay.ts`: the test neither imports") {
		t.Errorf("a name that only contains the unit's is not the unit's, got %v: %s", v, msg)
	}
}

func TestTestReach_importsReachTheModule(t *testing.T) {
	t.Run("TSRCH-B02: An import of the unit's module reaches it", func(t *testing.T) {})
	for _, test := range []string{"import * as p from './pay'\n", "import x from '../lib/src'\n"} {
		if !importsModule(test, "src/pay.ts", nil) {
			t.Errorf("%q imports the unit", test)
		}
	}
	for _, test := range []string{"const s = 'pay'\n", "log('src')\n"} {
		if importsModule(test, "src/pay.ts", nil) {
			t.Errorf("%q names the unit outside an import", test)
		}
	}
	declared := &config.Config{Dialect: &config.Dialect{ImportPattern: `^load\s`}}
	if !importsModule("load \"pay\"\n", "src/pay.ts", declared) || importsModule("load \"pay\"\n", "src/pay.ts", nil) {
		t.Error("the dialect's import_pattern names what an import line is")
	}
	if importsModule("import x from 'lib'\n", "pay.ts", nil) {
		t.Error("a unit at the root has no directory to import")
	}
	if v, msg := exercise(t, reachCfg(), "import { thing } from './pay'\n"); v != Pass {
		t.Errorf("an import reaches the unit even with no name of it, got %v: %s", v, msg)
	}
}

func TestTestReach_aCopyIsNamed(t *testing.T) {
	t.Run("TSRCH-B03: A test that defines the unit's names exercises a copy", func(t *testing.T) {})
	test := "import './pay'\nfunction charge(a) { return a }\nconst refund = (x) => x\ntest('a', () => expect(charge(1)).toBe(1))\n"
	v, msg := exercise(t, reachCfg(), test)
	if v != Fail || !strings.Contains(msg, "`src/pay.ts`: the test defines charge, refund again") || strings.Contains(msg, "neither imports") {
		t.Errorf("the copy fails naming the unit and the names, got %v: %s", v, msg)
	}
	test = "function charge(a) { return a }\ntest('a', () => expect(charge(1)).toBe(1))\n"
	if v, msg := exercise(t, reachCfg(), test); v != Fail || !strings.Contains(msg, "neither imports") {
		t.Errorf("a copied name does not reach the unit, got %v: %s", v, msg)
	}
}

func TestTestReach_shortNamesAreAnyones(t *testing.T) {
	t.Run("TSRCH-B04: Short names are anyone's", func(t *testing.T) {})
	def := definitionOf(reachCfg())
	names := definedNames("function ok() {}\nfunction pay() {}\n", def)
	if _, ok := names["ok"]; ok || len(names) != 1 || names["pay"].bare != "pay" {
		t.Errorf("a two-letter name is left out and a three-letter one kept, got %v", names)
	}
	if v, msg := exercise(t, reachCfg(), "function ok() {}\ntest('a', () => ok())\n"); v != Fail || strings.Contains(msg, "its own copy") {
		t.Errorf("a short name neither copies nor reaches, got %v: %s", v, msg)
	}
}

// refGraph has spec REFXX governing two code files, and the test.
func refGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "fn/create/handler.ts", Kind: mapx.KindCode}, {ID: "fn/create/rules.ts", Kind: mapx.KindCode},
			{ID: "fn/create/handler.spec.md", Kind: mapx.KindSpec, Code: "REFXX", CodeDeclarado: true},
			{ID: "e2e/flow.test.ts", Kind: mapx.KindTest},
		},
		Edges: []mapx.Edge{
			{Type: mapx.EdgeSpecifies, From: "fn/create/handler.spec.md", To: "fn/create/rules.ts"},
			{Type: mapx.EdgeSpecifies, From: "fn/create/handler.spec.md", To: "fn/create/handler.ts"},
		},
	}
}

func refMatch(t *testing.T, cfg *config.Config, test string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "fn/create/handler.ts", "export async function handler(event) {}\n")
	writeFile(t, root, "fn/create/rules.ts", "export function validate(x) {}\n")
	g := refGraph()
	return checkTestRefMatchesUnit(test, g.Nodes[3], root, g, cfg)
}

func TestTestReach_theRefsUnit(t *testing.T) {
	t.Run("TSRCH-B05: The ref's unit must be reached", func(t *testing.T) {})
	if v, msg := refMatch(t, reachCfg(), "// ref: REFXX\ntest('a', () => validate(1))\n"); v != Pass {
		t.Errorf("reaching one governed file passes, got %v: %s", v, msg)
	}
	v, msg := refMatch(t, reachCfg(), "// ref: REFXX\ntest('a', () => other(1))\n")
	if v != Fail || !strings.Contains(msg, "`ref: REFXX` names fn/create/handler.ts, fn/create/rules.ts") {
		t.Errorf("reaching none fails naming the ref and the files sorted, got %v: %s", v, msg)
	}
	noDef := &config.Config{Dialect: &config.Dialect{Family: "rust"}}
	if v, _ := refMatch(t, noDef, "// ref: REFXX\nimport { h } from '../fn/create/handler'\n"); v != Pass {
		t.Errorf("with no definition an import still reaches, got %v", v)
	}
}

func TestTestReach_declaredInvocations(t *testing.T) {
	t.Run("TSRCH-B06: A declared invocation reaches the unit it names", func(t *testing.T) {})
	cfg := reachCfg(`lambdaInvoke\(\s*(?:'(\w+)'|"(\w+)")`)
	for _, call := range []string{`lambdaInvoke("handler")`, `lambdaInvoke('create')`} {
		if v, msg := refMatch(t, cfg, "// ref: REFXX\n"+call+"\n"); v != Pass {
			t.Errorf("%s reaches the unit, got %v: %s", call, v, msg)
		}
	}
	if v, msg := refMatch(t, cfg, "// ref: REFXX\nlambdaInvoke('delete')\nlambdaInvoke('create')\n"); v != Pass {
		t.Errorf("every invocation is read, not only the first, got %v: %s", v, msg)
	}
	if v, _ := refMatch(t, cfg, "// ref: REFXX\nlambdaInvoke('delete')\n"); v != Fail {
		t.Errorf("an invocation naming another unit does not reach it, got %v", v)
	}
	if v, _ := refMatch(t, reachCfg(`(`), "// ref: REFXX\nlambdaInvoke('create')\n"); v != Fail {
		t.Errorf("an invocation that does not compile reaches nothing, got %v", v)
	}
	if invokes("run(fn)", "fn/create/handler.ts", nil) {
		t.Error("with no invocation declared nothing is invoked")
	}
}

func TestTestReach_skips(t *testing.T) {
	t.Run("TSRCH-B07: Nothing to confront is skipped", func(t *testing.T) {})
	cfg := reachCfg()
	for _, check := range []func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string){checkTestExercisesUnit, checkTestRefMatchesUnit} {
		if v, _ := check("", mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}, "", refGraph(), cfg); v != Skip {
			t.Errorf("a spec is skipped, got %v", v)
		}
		if v, _ := check("", mapx.Node{ID: "h.test.ts", Kind: mapx.KindTest, Support: true}, "", refGraph(), cfg); v != Skip {
			t.Errorf("a support file is skipped, got %v", v)
		}
	}
	resetTestedUnits()
	if v, _ := checkTestExercisesUnit("", mapx.Node{ID: "e2e/flow.test.ts", Kind: mapx.KindTest}, "", refGraph(), cfg); v != Skip {
		t.Errorf("a test no unit pairs with is skipped, got %v", v)
	}
	if v, _ := exercise(t, &config.Config{Dialect: &config.Dialect{Family: "rust"}, Derived: cfg.Derived}, "x"); v != Skip {
		t.Errorf("no definition is skipped, got %v", v)
	}
	if v, _ := refMatch(t, cfg, "test('a')\n"); v != Skip {
		t.Errorf("a test with no ref is skipped, got %v", v)
	}
	if v, _ := refMatch(t, cfg, "// ref: NOPEX\n"); v != Skip {
		t.Errorf("a ref governing no code is skipped, got %v", v)
	}
}

func TestTestReach_declared(t *testing.T) {
	t.Run("TSRCH-B08: A declared way to reach the unit passes", func(t *testing.T) {})
	copied := "function charge(a) { return a }\ntest('a', () => expect(charge(1)).toBe(1))\n"
	if v, msg := exercise(t, reachCfg(), "// @no-unit-import: the unit runs in a worker; the test drives the worker\n"+copied); v != Pass {
		t.Errorf("a declared test passes, got %v: %s", v, msg)
	}
	if v, _ := exercise(t, reachCfg(), "// @no-unit-import:\n"+copied); v != Fail {
		t.Errorf("a bare waiver waives nothing, got %v", v)
	}
}

// A fake that implements the unit's interface defines a method with the unit's method name.
// That is the test's own type, not a copy of the unit (reported from baas-proxy: a
// `fakePinger.Ping` beside the unit's `GormPinger.Ping` failed as "defines Ping again").
func TestTestReach_aFakeImplementingTheInterfaceIsNotACopy(t *testing.T) {
	t.Run("TSRCH-B09: A member of the test's own type is not a copy of the unit's", func(t *testing.T) {})
	root := t.TempDir()
	goCfg := &config.Config{Dialect: &config.Dialect{Family: "go"}}
	writeFile(t, root, "src/handlers/probes.go", "package handlers\n\ntype Pinger interface{ Ping(ctx context.Context) error }\n\ntype GormPinger struct{ db *gorm.DB }\n\nfunc (g GormPinger) Ping(ctx context.Context) error { return nil }\n\nfunc ReadyHandler(p Pinger) http.HandlerFunc { return nil }\n")
	fake := "package handlers\n\ntype fakePinger struct{ err error }\n\nfunc (f fakePinger) Ping(ctx context.Context) error { return f.err }\n\nfunc TestReady(t *testing.T) { ReadyHandler(fakePinger{}) }\n"
	if r := reachOf(fake, "src/handlers/probes.go", root, goCfg, definitionOf(goCfg), nil); len(r.Copied) != 0 || !r.Reached {
		t.Errorf("a fake implementing the interface reaches the unit and copies nothing, got %+v", r)
	}
	// The SAME method of the SAME type defined again is still a copy.
	again := "package handlers_test\n\nfunc (g GormPinger) Ping(ctx context.Context) error { return nil }\n"
	if r := reachOf(again, "src/handlers/probes.go", root, goCfg, definitionOf(goCfg), nil); len(r.Copied) != 1 || r.Copied[0] != "GormPinger.Ping" {
		t.Errorf("the unit's own method defined again is a copy, got %+v", r)
	}
	// Python: a method indented under the test's fake class is a member, not a copy.
	pyCfg := &config.Config{Dialect: &config.Dialect{Family: "python"}}
	writeFile(t, root, "app/probes.py", "class DbPinger:\n    def ping(self):\n        return True\n\ndef ready(pinger):\n    return pinger.ping()\n")
	pyFake := "class FakePinger:\n    def ping(self):\n        return False\n\ndef test_ready():\n    assert not ready(FakePinger())\n"
	if r := reachOf(pyFake, "app/probes.py", root, pyCfg, definitionOf(pyCfg), nil); len(r.Copied) != 0 || !r.Reached {
		t.Errorf("a python fake's method is a member, not a copy, got %+v", r)
	}
}
