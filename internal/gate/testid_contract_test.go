package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cfgHandle returns a config that DECLARES the anchoring attribute — without it the
// inventory gates skip, which is the opt-in behaviour.
func cfgHandle(attr string) *config.Config {
	return &config.Config{Derived: &config.Derived{Anchor: "code", TestHandle: attr}}
}

// inventoryFixture: a spec linked (specifies) to a unit whose content the test
// controls. There is no consumer surface: no test, no flow, no neighbour.
func inventoryFixture(t *testing.T, unitSrc string) (mapx.Node, *mapx.Graph, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(unitSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies}},
	}
	return spec, g, root
}

// consumerFixture: spec → feature → test (two hops, because `tested-by` starts at the
// FEATURE). `testSrc` is the test's content, the consumer of the handle. The code exposes
// `:abcd-screen`.
func consumerFixture(t *testing.T, testSrc string) (mapx.Node, *mapx.Graph, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.test.tsx"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}, {ID: "x.feature", Kind: mapx.KindFeature}, {ID: "x.test.tsx", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{
			{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies},
			{From: "x.spec.md", To: "x.feature", Type: mapx.EdgeCoveredBy},
			{From: "x.feature", To: "x.test.tsx", Type: mapx.EdgeTestedBy},
		},
	}
	return spec, g, root
}

const sectionOK = "## Superfície de Teste\n\n| id | role |\n| -- | ----- |\n| `:abcd-screen` | root |\n"

// e2eConfig declares the e2e surface in its two steps: `surfaces[e2e]` gives the KEY,
// `files[key]` gives the path.
func e2eConfig() *config.Config {
	return &config.Config{Derived: &config.Derived{
		Anchor: "code", TestHandle: "testID",
		Surfaces: map[string]string{"e2e": "e2e"},
		Files:    map[string]config.Padroes{"e2e": {"e2e/{{name}}.yaml"}},
	}}
}

// --- when the gate confronts at all ---

func TestTestIDContract_skipsWhatItCannotConfront(t *testing.T) {
	t.Run("TICTS-B01: The gate skips what is not a spec, and is Pending without a map", func(t *testing.T) {})
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	if v, _ := checkTestIDCoherent(sectionOK, mapx.Node{ID: "x.spec.md", Kind: mapx.KindCode}, root, g, cfgHandle("testID")); v != Skip {
		t.Errorf("a code node should Skip, got %v", v)
	}
	if v, _ := checkTestIDCoherent(sectionOK, n, root, nil, cfgHandle("testID")); v != Pending {
		t.Errorf("no map should be Pending, got %v", v)
	}
}

func TestTestIDDeclared_withoutDeclaredHandleSkips(t *testing.T) {
	t.Run("TICTS-B02: Without a declared handle attribute the gate skips", func(t *testing.T) {})
	// The point of the opt-in: without `derived.test_handle` the gate has no attribute to
	// look for. Inferring `testID` by default would have the gate report GREEN over what
	// it never checked — in an Android project, in any backend.
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	if v, _ := checkTestIDCoherent("", n, root, g, &config.Config{}); v != Skip {
		t.Errorf("without test_handle the gate must skip: %v", v)
	}
}

func TestTestIDHonored_withoutDeclaredHandleSkips(t *testing.T) {
	t.Run("TICTS-B02: Without a declared handle attribute the gate skips", func(t *testing.T) {})
	n, g, root := consumerFixture(t, `render(<X />)`)
	if v, _ := checkTestIDCoherent(sectionOK, n, root, g, &config.Config{}); v != Skip {
		t.Errorf("without test_handle the gate must skip: %v", v)
	}
}

// A spec with no linked code, or whose linked code cannot be read, has no end that
// exposes the handle: nothing to confront.
func TestTestIDContract_withoutReadableCodeSkips(t *testing.T) {
	t.Run("TICTS-B03: A spec with no readable linked code skips", func(t *testing.T) {})
	t.Run("TICTS-E01: Linked code that cannot be read is not an end", func(t *testing.T) {})
	root := t.TempDir()
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	noEdge := &mapx.Graph{Nodes: []mapx.Node{spec}}
	if v, _ := checkTestIDCoherent(sectionOK, spec, root, noEdge, cfgHandle("testID")); v != Skip {
		t.Errorf("no linked code should Skip, got %v", v)
	}
	missing := &mapx.Graph{
		Nodes: []mapx.Node{spec},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "gone.tsx", Type: mapx.EdgeSpecifies}},
	}
	if v, msg := checkTestIDCoherent(sectionOK, spec, root, missing, cfgHandle("testID")); v != Skip || !strings.Contains(msg, "without linked code") {
		t.Errorf("unreadable linked code should Skip as no code: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_unitWithoutHandleIsNoOffence(t *testing.T) {
	t.Run("TICTS-B04: A unit that exposes no handle and declares none skips", func(t *testing.T) {})
	// Not every component has a test surface. Charging an inventory of a unit that marks
	// no element would become noise over every pure presentation unit.
	n, g, root := inventoryFixture(t, `<View />`)
	if v, _ := checkTestIDCoherent("", n, root, g, cfgHandle("testID")); v != Skip {
		t.Errorf("a unit exposing no handle has no contract to declare: %v", v)
	}
}

// --- the four ends ---

func TestTestIDDeclared_completeInventoryPasses(t *testing.T) {
	t.Run("TICTS-B05: A handle exposed, declared and queried passes", func(t *testing.T) {})
	t.Run("TICTS-B09: With no consumer surface the queried end is not charged", func(t *testing.T) {})
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("an inventory matching the code should pass: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_queriedIdPasses(t *testing.T) {
	t.Run("TICTS-B05: A handle exposed, declared and queried passes", func(t *testing.T) {})
	n, g, root := consumerFixture(t, `getByTestId(':abcd-screen')`)
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("an id queried by the test should pass: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_exposedWithoutDeclaringFails(t *testing.T) {
	t.Run("TICTS-B06: A handle the code exposes and the spec does not declare fails", func(t *testing.T) {})
	// Direction code → spec: an uncontracted surface. It is where an identity divergence
	// gets in without anyone seeing.
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" /><View testID=":abcd-hidden" />`)
	v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("exposed without declaring should fail: %v", v)
	}
	if !strings.Contains(msg, "abcd-hidden") {
		t.Errorf("the message must name the undeclared id: %s", msg)
	}
}

func TestTestIDContract_withoutInventoryChargesTheExposed(t *testing.T) {
	t.Run("TICTS-B06: A handle the code exposes and the spec does not declare fails", func(t *testing.T) {})
	// With one gate for the whole contract there is nothing to duplicate: skipping a spec
	// with no inventory would let through exactly the case where the code exposes a handle
	// nobody declared, the uncontracted surface.
	n, g, root := consumerFixture(t, `render(<X />)`)
	v, msg := checkTestIDCoherent("", n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("a spec with no inventory and code exposing a handle must fail: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, ":abcd-screen") {
		t.Errorf("the report must name the exposed, undeclared handle; got: %s", msg)
	}
}

func TestTestIDDeclared_declaredWithoutExposingFails(t *testing.T) {
	t.Run("TICTS-B07: A handle the spec declares and the code does not expose fails", func(t *testing.T) {})
	// Direction spec → code: a dead contract. Whoever writes the test looks for what does
	// not exist. It is the half a one-direction gate does not catch.
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	spec := sectionOK + "| `:abcd-ghost` | gone in the refactor |\n"
	v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("declared without exposing should fail: %v", v)
	}
	if !strings.Contains(msg, "abcd-ghost") {
		t.Errorf("the message must name the missing id: %s", msg)
	}
}

func TestTestIDHonored_orphanIdFails(t *testing.T) {
	t.Run("TICTS-B08: A handle no consumer queries fails when a consumer surface exists", func(t *testing.T) {})
	// Declared, exposed and nobody queries it: cost with no return. And worse than useless —
	// it looks like coverage, because the surface is there and the inventory complete.
	n, g, root := consumerFixture(t, `render(<X />)`)
	v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("an id nobody queries should fail: %v", v)
	}
	if !strings.Contains(msg, "abcd-screen") {
		t.Errorf("the message must name the orphan id: %s", msg)
	}
}

// Where the project has no consumer surface, the queried column says so with a dash:
// an ✗ there would accuse the author of a missing configuration.
func TestTestIDContract_noConsumerShowsADash(t *testing.T) {
	t.Run("TICTS-B09: With no consumer surface the queried end is not charged", func(t *testing.T) {})
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" /><View testID=":abcd-hidden" />`)
	_, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID"))
	if !strings.Contains(msg, "queried —") || strings.Contains(msg, "queried ✗\n") {
		t.Errorf("with no consumer the queried end must read —: %s", msg)
	}
}

// The mark is a writing convention, not part of the identity: `:abcd-screen` in the code
// and `abcd-screen` in the spec are one handle, reported on one line, spelled as the code
// writes it.
func TestTestIDContract_oneLinePerHandle(t *testing.T) {
	t.Run("TICTS-I01: One handle is one line of the report, whatever the spelling at each end", func(t *testing.T) {})
	n, g, root := consumerFixture(t, `render(<X />)`)
	spec := "## Test IDs\n\n- `abcd-screen`\n"
	v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("nobody queries the handle, expected Fail: %v", v)
	}
	if strings.Count(msg, "abcd-screen") != 1 || !strings.Contains(msg, ":abcd-screen") {
		t.Errorf("the handle must appear once, spelled as the code writes it: %s", msg)
	}
}

// The fourth edge — a consumer querying a handle no code exposes — is the project's
// question (testid-queried-exists), not one spec's.
func TestTestIDContract_queriedButNotExposedIsNotChargedHere(t *testing.T) {
	t.Run("TICTS-X01: A handle only a consumer mentions is not charged to the spec", func(t *testing.T) {})
	n, g, root := consumerFixture(t, `getByTestId(':abcd-screen'); getByTestId(':other-screen-ghost')`)
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a handle only the consumer mentions is not this spec's: %v (%s)", v, msg)
	}
}

// --- how the code exposes a handle ---

func TestTestIDDeclared_attributeOfAnotherEcosystem(t *testing.T) {
	t.Run("TICTS-B10: The handle attribute is the one the project declares", func(t *testing.T) {})
	// The same gate in a web project: the attribute is `data-testid`. If the regex were
	// hardwired to `testID`, this case would pass empty (false green).
	n, g, root := inventoryFixture(t, `<div data-testid="abcd-root" /><div data-testid="abcd-item" />`)
	spec := "## Superfície de Teste\n\n- `abcd-root`\n"
	v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("data-testid"))
	if v != Fail || !strings.Contains(msg, "abcd-item") {
		t.Errorf("it should accuse the undeclared web id: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_derivedPropIsAHandle(t *testing.T) {
	t.Run("TICTS-B11: A literal handle counts, also through a derived prop or an object key", func(t *testing.T) {})
	// `backTestID=":spending-button-back"` and `confirmTestID: ':x'` (object property): the
	// value ends up as a real handle in the child. Ignoring them had the gate call missing
	// an id the spec documents and the flows use.
	n, g, root := inventoryFixture(t,
		`<Header backTestID=":abcd-back" /><Alert buttons={[{ confirmTestID: ':abcd-ok' }]} />`)
	spec := "## Test IDs\n\n- `abcd-back`\n- `abcd-ok`\n"
	if v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a derived prop carries a handle: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_templateDeclaresThePrefix(t *testing.T) {
	t.Run("TICTS-B12: A template handle, or a prefix prop, exposes the prefix as a wildcard", func(t *testing.T) {})
	// The suffix is runtime DATA (`${id}`). Requiring the spec to declare `abcd-item-42`
	// would require it to declare the data; the prefix is declared with the `-*` mark.
	n, g, root := inventoryFixture(t, "<View testID={`:abcd-item-${id}`} />")
	spec := "## Superfície de Teste\n\n- `:abcd-item-*`\n"
	if v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a dynamic prefix declared with `-*` should pass: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_prefixBuiltInTheChild(t *testing.T) {
	t.Run("TICTS-B12: A template handle, or a prefix prop, exposes the prefix as a wildcard", func(t *testing.T) {})
	// `testIDPrefix=":otp-input"` becomes `otp-input-0`…`otp-input-5` in the child. The
	// literal is the HEAD, not an id: registering it as is would charge the spec an
	// `otp-input` that never appears, and call missing the `otp-input-N` six flows use.
	n, g, root := inventoryFixture(t, `<Otp testIDPrefix=":otp-input" />`)
	spec := "## Test IDs\n\n- `otp-input-0`\n- `otp-input-5`\n"
	if v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a built prefix covers the concrete ids: %v (%s)", v, msg)
	}
}

func TestExposedTestIDs_collectsEveryOccurrence(t *testing.T) {
	t.Run("TICTS-B12: A template handle, or a prefix prop, exposes the prefix as a wildcard", func(t *testing.T) {})
	t.Run("TICTS-B13: Only the branches of a conditional handle are handles, never its condition", func(t *testing.T) {})
	src := "<View testID={`row-${id}`} /><Text testID={`cell-${i}`} />\n" +
		"<Icon testID={on ? 'icon-on' : 'icon-off'} />\n"
	got := exposedTestIDs(src, "testID")
	want := []string{"row-*", "cell-*", "icon-on", "icon-off"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("every template and both ternary branches:\n got %v\nwant %v", got, want)
	}
}

func TestTestIDDeclared_ternaryIgnoresTheCondition(t *testing.T) {
	t.Run("TICTS-B13: Only the branches of a conditional handle are handles, never its condition", func(t *testing.T) {})
	// `testID={c2.k === 'push' ? firstToggleTestID : undefined}` — the literal of the
	// CONDITION is domain state, not a handle. Collecting it would register `push` as
	// exposed and the gate would charge the spec an id that does not exist.
	n, g, root := inventoryFixture(t,
		`<V testID={k === 'push' ? ':abcd-toggle' : undefined} /><V testID={i === 0 ? ':abcd-first' : undefined} />`)
	spec := "## Test IDs\n\n- `abcd-toggle`\n- `abcd-first`\n"
	if v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("the condition is not a handle; both branches are: %v (%s)", v, msg)
	}
}

// --- how the spec declares a handle ---

func TestTestIDDeclared_titleWithQualifier(t *testing.T) {
	t.Run("TICTS-B14: The inventory is read only inside the test surface section", func(t *testing.T) {})
	// `## Test IDs (Maestro)` — 42 of the 159 specs of the reference app name in the title
	// the surface that consumes the ids. Requiring the line to end right after the title
	// had the gate miss the section and accuse the best-documented specs of declaring
	// nothing.
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	spec := "## Test IDs (Maestro)\n\n| testID | Element |\n| -- | -- |\n| `abcd-screen` | root |\n"
	if v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a title with a qualifier is the same section: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_backtickOutsideTheSectionDoesNotCount(t *testing.T) {
	t.Run("TICTS-B14: The inventory is read only inside the test surface section", func(t *testing.T) {})
	// The section bounds the inventory. Without it, any backtick in the prose (a CSS
	// class, a field name) would count as a declaration — and the gate would pass by
	// accident, which is worse than failing by mistake.
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	spec := "## Notes\n\nThe root uses `:abcd-screen` and the class `text-mute`.\n"
	v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("a mention in prose is not an inventory: %v (%s)", v, msg)
	}
}

// The section ends at the next heading: a backtick in the section after it is prose.
func TestDeclaredTestIDs_sectionEndsAtTheNextHeading(t *testing.T) {
	t.Run("TICTS-B14: The inventory is read only inside the test surface section", func(t *testing.T) {})
	spec := "## Test IDs\n\n- `abcd-screen`\n\n## Notes\n\n- `abcd-later`\n"
	if got := declaredTestIDs(spec, "testID"); !reflect.DeepEqual(got, []string{"abcd-screen"}) {
		t.Errorf("only the ids before the next heading are declared: %v", got)
	}
}

func TestTestIDDeclared_onlyTheFirstColumnIsTheID(t *testing.T) {
	t.Run("TICTS-B15: In a table only the first cell is the id; on any other line every quoted id counts", func(t *testing.T) {})
	// The reference app's inventory table has "Element" and "Used in (flow)" columns, and
	// both carry backticks: `TouchableOpacity root`, `ATLNX-VR`. Reading the whole line
	// had the gate collect those cells as declared testIDs — and then accuse them of
	// being orphans, inventing debt out of the documentation.
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	spec := "## Test IDs (Maestro)\n\n" +
		"| testID | Element | Used in |\n| -- | -- | -- |\n" +
		"| `abcd-screen` | Root of the `TouchableOpacity` | `ABCDX-VR` |\n"
	v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID"))
	if v != Pass {
		t.Errorf("only the first cell is the id; the others describe it: %v (%s)", v, msg)
	}
	if ids := declaredTestIDs(spec, "testID"); len(ids) != 1 || ids[0] != "abcd-screen" {
		t.Errorf("the inventory should hold only the first column's id, got: %v", ids)
	}
}

func TestDeclaredTestIDs_everyIdOnALine(t *testing.T) {
	t.Run("TICTS-B15: In a table only the first cell is the id; on any other line every quoted id counts", func(t *testing.T) {})
	spec := "## Test IDs\n\n- `abcd-screen`, `abcd-close`\n"
	got := declaredTestIDs(spec, "testID")
	if !reflect.DeepEqual(got, []string{"abcd-screen", "abcd-close"}) {
		t.Errorf("got %v", got)
	}
}

func TestTestIDDeclared_attributeNameIsNotAnID(t *testing.T) {
	t.Run("TICTS-B16: The attribute's own name is not a declared id", func(t *testing.T) {})
	// The generic atom (Button, Avatar) RECEIVES the handle from outside, and the spec
	// documents it as `testID` (prop). Declaring that says "I accept a handle", not "I
	// expose this one" — charging it would demand a fixed id from the atom, the opposite
	// of being reusable.
	n, g, root := inventoryFixture(t, `<Touchable testID={testID} />`)
	spec := "## Test IDs\n\n| testID | Element |\n| -- | -- |\n| `testID` (prop) | Root |\n"
	if v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID")); v != Skip && v != Pass {
		t.Errorf("the attribute's name is not a declared id: %v (%s)", v, msg)
	}
	if ids := declaredTestIDs(spec, "testID"); len(ids) != 0 {
		t.Errorf("the attribute's name was declared as an id: %v", ids)
	}
}

func TestTestIDDeclared_wildcardMatchesBothWays(t *testing.T) {
	t.Run("TICTS-B17: A wildcard at either end covers the concrete ids it opens", func(t *testing.T) {})
	// An exposed generic covers a declared concrete, and vice versa — the generic form is
	// born both in the code (template) and in the spec.
	if !covers("abcd-item-*", "abcd-item-3") {
		t.Error("an exposed wildcard should cover the concrete id")
	}
	if covers("abcd-item-3", "abcd-item-*") {
		t.Error("a concrete id does NOT cover the wildcard (only the reverse)")
	}
	if !covers("abcd-x", "abcd-x") || covers("abcd-x", "abcd-y") {
		t.Error("exact equality should hold, and only it")
	}
	// Through the gate: the spec declares the wildcard, the code exposes concrete ids.
	n, g, root := inventoryFixture(t, `<V testID=":abcd-item-1" /><V testID=":abcd-item-2" />`)
	if v, msg := checkTestIDCoherent("## Test IDs\n\n- `abcd-item-*`\n", n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a declared wildcard covers the exposed concrete ids: %v (%s)", v, msg)
	}
}

// --- who consumes a handle ---

func TestTestIDHonored_dynamicPrefixIsEnough(t *testing.T) {
	t.Run("TICTS-B18: A handle is queried when a consumer mentions it; a wildcard, when it mentions the prefix", func(t *testing.T) {})
	// The flow matches `abcd-item-.*` or builds `abcd-item-${id}`; the suffix is runtime
	// data. Requiring the whole id would make every template a false orphan.
	n, g, root := consumerFixture(t, "getByTestId(`:abcd-item-${id}`)")
	// The code must EXPOSE the template — the default fixture exposes `:abcd-screen`,
	// another handle. Without this the case would measure "declared and not exposed".
	if err := os.WriteFile(filepath.Join(root, "x.tsx"),
		[]byte("<View testID={`:abcd-item-${id}`} />"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := "## Superfície de Teste\n\n- `:abcd-item-*`\n"
	if v, msg := checkTestIDCoherent(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a queried prefix honours the dynamic contract: %v (%s)", v, msg)
	}
}

// The consumer may have been written before the marking convention: it mentions the id
// without the mark.
func TestTestIDContract_queriedWithoutTheMark(t *testing.T) {
	t.Run("TICTS-B18: A handle is queried when a consumer mentions it; a wildcard, when it mentions the prefix", func(t *testing.T) {})
	n, g, root := consumerFixture(t, `getByTestId('abcd-screen')`)
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a query without the mark is still a query: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_e2eFlowCountsAsConsumer(t *testing.T) {
	t.Run("TICTS-B19: The e2e flows the project declares are consumers", func(t *testing.T) {})
	// The end-to-end flow is the MAIN consumer of the handle and lives OUTSIDE the graph.
	// If only the unit test counted, every id used only by the E2E would be reported as an
	// orphan — accusing exactly who honours the contract.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.test.tsx"), []byte(`render(<X />)`), 0o644); err != nil {
		t.Fatal(err)
	}
	flows := filepath.Join(root, "e2e")
	if err := os.MkdirAll(flows, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(flows, "f.yaml"), []byte("- tapOn:\n    id: ':abcd-screen'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}, {ID: "x.feature", Kind: mapx.KindFeature}, {ID: "x.test.tsx", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{
			{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies},
			{From: "x.spec.md", To: "x.feature", Type: mapx.EdgeCoveredBy},
			{From: "x.feature", To: "x.test.tsx", Type: mapx.EdgeTestedBy},
		},
	}
	if v, msg := checkTestIDCoherent(sectionOK, spec, root, g, e2eConfig()); v != Pass {
		t.Errorf("an id used only by the e2e flow is not an orphan: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_surfaceWithoutFilesInventsNoPath(t *testing.T) {
	t.Run("TICTS-B19: The e2e flows the project declares are consumers", func(t *testing.T) {})
	// The trap this test locks: `surfaces[e2e]` returns the surface's KEY, not a path.
	// Treating it as a directory has the gate walk a missing path, find zero consumers and
	// accuse every id only the E2E uses. Without `files[e2e]` there is nowhere to look,
	// and the honest answer is to claim nothing about a surface it cannot read.
	n, g, root := consumerFixture(t, `getByTestId(':abcd-screen')`)
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "code", TestHandle: "testID",
		Surfaces: map[string]string{"e2e": "e2e"}, // no files["e2e"]
	}}
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfg); v != Pass {
		t.Errorf("a surface with no declared path must not become an accusation: %v (%s)", v, msg)
	}
}

// When the surface has no file template, the first override that declares it gives the
// path.
func TestTestIDContract_e2eSurfaceFromAnOverride(t *testing.T) {
	t.Run("TICTS-B19: The e2e flows the project declares are consumers", func(t *testing.T) {})
	n, g, root := inventoryFixture(t, `<View testID=":abcd-screen" />`)
	if err := os.MkdirAll(filepath.Join(root, "flows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "flows", "f.yaml"), []byte("id: ':abcd-other'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "code", TestHandle: "testID",
		Surfaces:  map[string]string{"e2e": "e2e"},
		Overrides: []config.DerivedOverride{{Files: map[string]config.Padroes{"e2e": {"flows/{{name}}.yaml"}}}},
	}}
	v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "queried ✗\n") {
		t.Errorf("the override's flows are a consumer surface that does not query the id: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_sharedTestCountsAsConsumer(t *testing.T) {
	t.Run("TICTS-B20: The test files beside the spec, and in its sibling folders, are consumers", func(t *testing.T) {})
	// `PhiScreens.test.tsx` proves six screens. The `tested-by` edge links the feature to
	// the file, but not every spec of the group reaches it through the graph — and the
	// handle, though queried, would show up as an orphan. The neighbour test recovers it.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Shared.test.tsx"),
		[]byte(`getByTestId(':abcd-screen')`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	// The code is there (the gate needs the end that EXPOSES), but the edge to the TEST
	// is absent — that is what this case measures: the consumer reached by neighbourhood,
	// not by the graph.
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies}},
	} // no edge to a test — that is the point
	if v, msg := checkTestIDCoherent(sectionOK, spec, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("a neighbour test querying the id honours the contract: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_neighbourCountsEvenWithFlows(t *testing.T) {
	t.Run("TICTS-B20: The test files beside the spec, and in its sibling folders, are consumers", func(t *testing.T) {})
	// A guard `len(consumers)==0` used to cancel the safety net: a project declaring an
	// e2e surface always returns the flows' content, so the list was never empty and the
	// neighbour test was never read. The gate accused 14 ids the shared tests query.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Shared.test.tsx"),
		[]byte(`getByTestId(':abcd-screen')`), 0o644); err != nil {
		t.Fatal(err)
	}
	flows := filepath.Join(root, "e2e")
	if err := os.MkdirAll(flows, 0o755); err != nil {
		t.Fatal(err)
	}
	// A flow that exists but does NOT cite the id — what kept `consumers` non-empty.
	if err := os.WriteFile(filepath.Join(flows, "f.yaml"), []byte("- tapOn:\n    id: ':other'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies}},
	}
	if v, msg := checkTestIDCoherent(sectionOK, spec, root, g, e2eConfig()); v != Pass {
		t.Errorf("the neighbour test must count even with flows: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_componentExercisedByTheScreen(t *testing.T) {
	t.Run("TICTS-B20: The test files beside the spec, and in its sibling folders, are consumers", func(t *testing.T) {})
	// `CategoryCard` lives in `trends/components/` and whoever exercises it is
	// `TrendsScreen.test.tsx`, in `trends/screens/` — the screen that renders it. Without
	// going up to the module, the gate accuses as orphan a handle the test queries.
	root := t.TempDir()
	comp := filepath.Join(root, "trends", "components")
	screen := filepath.Join(root, "trends", "screens")
	for _, d := range []string{comp, screen} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(screen, "TrendsScreen.test.tsx"),
		[]byte(`getByTestId(':abcd-screen')`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "trends/components/CategoryCard.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	// The edge starts from the spec's REAL id, which has a path here. The edge to the TEST
	// is still absent: this case measures the consumer reached by neighbourhood.
	code := filepath.Join("trends", "components", "CategoryCard.tsx")
	if err := os.WriteFile(filepath.Join(root, code), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: code, Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: spec.ID, To: code, Type: mapx.EdgeSpecifies}},
	}
	if v, msg := checkTestIDCoherent(sectionOK, spec, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("the sibling screen's test exercises the component: %v (%s)", v, msg)
	}
}

// A surface declared with a glob is read from its static root. The glob stayed in the
// path (`e2e/**/*.yaml`), the walk found no directory, and every id only the flows use
// was accused as an orphan — the same trap as the key-for-path one above.
func TestTestIDHonored_globSurfaceIsReadFromItsRoot(t *testing.T) {
	t.Run("TICTS-B21: A surface declared with a glob is read from the directory before the first wildcard", func(t *testing.T) {})
	for pattern, want := range map[string]string{
		"e2e/**/*.yaml":           "e2e",
		"e2e/login-*.yaml":        "e2e",
		"e2e/{{name}}.yaml":       "e2e",
		"apps/x-{{module}}/flows": "apps",
		"e2e/flows":               "e2e/flows",
		"**/*.yaml":               ".",
	} {
		if got := firstStaticSegment(pattern); got != want {
			t.Errorf("firstStaticSegment(%q) = %q, want %q", pattern, got, want)
		}
	}
	n, g, root := consumerFixture(t, `render(<X />)`)
	if err := os.MkdirAll(filepath.Join(root, "e2e", "login"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "e2e", "login", "f.yaml"), []byte("- tapOn:\n    id: ':abcd-screen'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := e2eConfig()
	cfg.Derived.Files = map[string]config.Padroes{"e2e": {"e2e/**/*.yaml"}}
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfg); v != Pass {
		t.Errorf("the flow under the glob queries the id: %v (%s)", v, msg)
	}
}

// A handle is queried only when a consumer names THAT handle. The mention was a bare
// substring: a test that only queries `:abcd-screen-header` counted as querying
// `:abcd-screen`, and an id nobody uses passed as consumed.
func TestTestIDHonored_longerIdIsNotAQuery(t *testing.T) {
	t.Run("TICTS-B22: A consumer that names only a longer id does not query the shorter one", func(t *testing.T) {})
	n, g, root := consumerFixture(t, `getByTestId(':abcd-screen-header'); getByTestId('xabcd-screen')`)
	v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID"))
	if v != Fail || !strings.Contains(msg, "abcd-screen") {
		t.Errorf("abcd-screen is queried by nobody: %v (%s)", v, msg)
	}
	for blob, want := range map[string]bool{
		`getByTestId('abcd-screen')`:        true,
		`getByTestId(":abcd-screen")`:       true,
		"id: abcd-screen\n":                 true,
		`getByTestId('abcd-screen-header')`: false,
		`getByTestId('abcd-screens')`:       false,
		`getByTestId('my-abcd-screen')`:     false,
		"getByTestId(`abcd-screen-${i}`)":   false,
	} {
		if got := queriedID(blob, ":abcd-screen"); got != want {
			t.Errorf("queriedID(%q, :abcd-screen) = %v, want %v", blob, got, want)
		}
	}
	if !queriedID("id: 'abcd-item-.*'", ":abcd-item-*") || queriedID("id: 'xabcd-item-3'", ":abcd-item-*") {
		t.Errorf("a wildcard is queried by its prefix, and only at an id boundary")
	}
}

// farConsumerFixture: spec → feature → test, with the test in a folder no neighbour
// scan reaches (`tests/far/`), plus a test beside the spec that queries nothing — so
// a consumer surface exists either way, and only the two-hop link can make the handle
// queried. The code exposes `:abcd-screen`; the feature's content is `featureSrc`.
func farConsumerFixture(t *testing.T, farTestSrc, featureSrc string) (mapx.Node, *mapx.Graph, string) {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"mod/screens/x.tsx":          `<View testID=":abcd-screen" />`,
		"mod/screens/x.feature":      featureSrc,
		"mod/screens/other.test.tsx": `render(<Y />)`,
		"tests/far/x.test.tsx":       farTestSrc,
	}
	for rel, body := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	spec := mapx.Node{ID: "mod/screens/x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec,
			{ID: "mod/screens/x.tsx", Kind: mapx.KindCode},
			{ID: "mod/screens/x.feature", Kind: mapx.KindFeature},
			{ID: "tests/far/x.test.tsx", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{
			{From: "mod/screens/x.spec.md", To: "mod/screens/x.tsx", Type: mapx.EdgeSpecifies},
			{From: "mod/screens/x.spec.md", To: "mod/screens/x.feature", Type: mapx.EdgeCoveredBy},
			{From: "mod/screens/x.feature", To: "tests/far/x.test.tsx", Type: mapx.EdgeTestedBy},
		},
	}
	return spec, g, root
}

// The existing consumer fixtures put the linked test beside the spec, where the
// neighbour scan reads it anyway — so the two-hop link itself was never proven.
func TestTestIDContract_linkedTestIsAConsumerWhereverItLives(t *testing.T) {
	t.Run("TICTS-B23: The test linked to the spec's feature is a consumer wherever it lives", func(t *testing.T) {})
	n, g, root := farConsumerFixture(t, `getByTestId(':abcd-screen')`, "Feature: x\n")
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("the test reached through the feature queries the handle: %v (%s)", v, msg)
	}
	// The counter-proof: the far test is what makes it queried.
	n, g, root = farConsumerFixture(t, `render(<X />)`, "Feature: x\n")
	if v, msg := checkTestIDCoherent(sectionOK, n, root, g, cfgHandle("testID")); v != Fail {
		t.Errorf("with the far test querying nothing, the handle is an orphan: %v (%s)", v, msg)
	}
}

// The feature column is information: a handle the feature describes reads ✓, one it
// does not reads ✗ — and neither decides the verdict.
func TestTestIDContract_reportShowsTheFeatureEnd(t *testing.T) {
	t.Run("TICTS-B24: The report shows whether the feature describes each handle", func(t *testing.T) {})
	n, g, root := farConsumerFixture(t, `render(<X />)`, "Then the element \":abcd-screen\" is visible\n")
	if err := os.WriteFile(filepath.Join(root, "mod/screens/x.tsx"),
		[]byte(`<View testID=":abcd-screen" /><View testID=":abcd-other" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	v, msg := checkTestIDCoherent("", n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("neither handle is declared, expected Fail: %v (%s)", v, msg)
	}
	lineOf := func(id string) string {
		for i, l := range strings.Split(msg, "\n") {
			if strings.TrimSpace(l) == id {
				return strings.Split(msg, "\n")[i+1]
			}
		}
		t.Fatalf("no report line for %s: %s", id, msg)
		return ""
	}
	if l := lineOf(":abcd-screen"); !strings.Contains(l, "feature ✓") {
		t.Errorf("the feature describes :abcd-screen, the line must read feature ✓: %q", l)
	}
	if l := lineOf(":abcd-other"); !strings.Contains(l, "feature ✗") {
		t.Errorf("the feature does not describe :abcd-other, the line must read feature ✗: %q", l)
	}
}
