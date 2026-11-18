// @anchors
//   code: MDTSA
//   ref: GRMDG

package mapx

import (
	"reflect"
	"sort"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStale(t *testing.T) {
	t.Run("GRMDG-B01: A relation never stamped is stale", func(t *testing.T) {})
	t.Run("GRMDG-B02: A stamp is fresh while both ends keep their revisions", func(t *testing.T) {})
	g := &Graph{
		Nodes: []Node{{ID: "a", Rev: "r2"}, {ID: "b", Rev: "r5"}},
	}
	// edge with no stamp → stale
	if !g.Stale(Edge{From: "a", To: "b"}) {
		t.Error("an edge with no stamp should be stale")
	}
	// the stamp matches the current revs → not stale
	fresh := Edge{From: "a", To: "b", Stamp: &Stamp{ValidatedFromRev: "r2", ValidatedToRev: "r5"}}
	if g.Stale(fresh) {
		t.Error("an edge with a current stamp should not be stale")
	}
	// one end moved → stale
	old := Edge{From: "a", To: "b", Stamp: &Stamp{ValidatedFromRev: "r2", ValidatedToRev: "r4"}}
	if !g.Stale(old) {
		t.Error("an edge whose target moved should be stale")
	}
	oldFrom := Edge{From: "a", To: "b", Stamp: &Stamp{ValidatedFromRev: "r1", ValidatedToRev: "r5"}}
	if !g.Stale(oldFrom) {
		t.Error("an edge whose source moved should be stale")
	}
}

func TestStaleReadsRevisionsOnly(t *testing.T) {
	t.Run("GRMDG-X01: The stamp's date and verdict play no part in staleness", func(t *testing.T) {})
	g := &Graph{Nodes: []Node{{ID: "a", Rev: "r2"}, {ID: "b", Rev: "r5"}}}
	e := Edge{From: "a", To: "b", Stamp: &Stamp{ValidatedFromRev: "r2", ValidatedToRev: "r5", Verdict: "issue", ChangedAt: "2020-01-01"}}
	if g.Stale(e) {
		t.Error("a stamp at the current revisions is fresh, whatever its verdict and date")
	}
}

func TestMutationScopeStale(t *testing.T) {
	t.Run("GRMDG-B03: A mutation scope is stale only against a recorded, different revision", func(t *testing.T) {})
	measured := MutationScope{AtRev: "r1"}
	unknown := MutationScope{}
	if !measured.Stale("r2") {
		t.Error("a scope measured at r1 is stale against r2")
	}
	if measured.Stale("r1") {
		t.Error("a scope measured at r1 is fresh against r1")
	}
	if unknown.Stale("r2") || unknown.Stale("r1") {
		t.Error("a scope with no recorded revision is never called stale")
	}
}

// keysOf writes v in the map format and returns its top-level keys, sorted.
func keysOf(t *testing.T, v any) []string {
	t.Helper()
	data, err := yaml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sorted(s ...string) []string {
	sort.Strings(s)
	return s
}

// The keys are the file format: a renamed key is a format change, not a refactor.
func TestModelKeysAreFixed(t *testing.T) {
	t.Run("GRMDG-I01: The map file's keys are fixed and English", func(t *testing.T) {})
	stamp := Stamp{ValidatedFromRev: "a", ValidatedToRev: "b", ChangedAt: "d", Verdict: "ok", Gate: "g"}
	judgment := Judgment{Gate: "g", Verdict: "ok", ValidatedFromRev: "a", ValidatedToRev: "b", ChangedAt: "d"}
	edge := Edge{From: "a", To: "b", Julgamentos: []Judgment{judgment}, Type: EdgeDependsOn, Origin: OriginDeclared,
		Method: "m", Dep: "DEP1", Stamp: &stamp}
	node := Node{ID: "a", Kind: KindCode, Rev: "r", UpdatedAt: "d", Code: "AAAAA", Layer: "l", CodeDeclarado: true,
		Tags: []string{"t"}, Regime: "x", NoPropagation: true, SharedCode: true, Needs: []string{"n"}, Parent: "P",
		Upstream: true, Revises: []string{"v"}, Signal: &TestSignal{Passed: 1}, Failures: []FailureSignal{{Rule: "A-E01", Count: 1}}}
	graph := Graph{Version: 4, GeradoPor: "dev", Nodes: []Node{node}, Edges: []Edge{edge}, Flow: &FlowGraph{}}

	for _, c := range []struct {
		name string
		v    any
		want []string
	}{
		{"graph", graph, sorted("version", "generated_by", "nodes", "edges", "flow")},
		{"edge", edge, sorted("from", "to", "judgments", "type", "origin", "method", "dep", "stamp")},
		{"stamp", stamp, sorted("validated_from_rev", "validated_to_rev", "changed_at", "verdict", "gate")},
		{"judgment", judgment, sorted("gate", "verdict", "validated_from_rev", "validated_to_rev", "changed_at")},
		{"node", node, sorted("id", "kind", "rev", "updated_at", "code", "layer", "code_declared", "tags", "regime",
			"no_propagation", "shared_code", "needs", "parent", "upstream", "revises", "signal", "failures")},
	} {
		if got := keysOf(t, c.v); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s keys = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestModelOmitsEmptyOptionalFields(t *testing.T) {
	t.Run("GRMDG-I02: Empty optional fields are left out", func(t *testing.T) {})
	if got := keysOf(t, Node{ID: "a", Kind: KindCode, Rev: "r"}); !reflect.DeepEqual(got, sorted("id", "kind", "rev")) {
		t.Errorf("a bare node writes %v", got)
	}
	if got := keysOf(t, Edge{From: "a", To: "b", Type: EdgeGoverns, Origin: OriginDeclared}); !reflect.DeepEqual(got, sorted("from", "to", "type", "origin")) {
		t.Errorf("a bare edge writes %v", got)
	}
}

func TestModelVocabulary(t *testing.T) {
	t.Run("GRMDG-I03: The node kinds and relation types are the fixed vocabulary", func(t *testing.T) {})
	kinds := []Kind{KindSpec, KindFeature, KindTest, KindCode, KindDoc, KindGuide, KindPlan, KindProduct, KindFlag}
	wantKinds := []Kind{"spec", "feature", "test", "code", "doc", "guide", "plan", "product", "flag"}
	if !reflect.DeepEqual(kinds, wantKinds) {
		t.Errorf("kinds = %v, want %v", kinds, wantKinds)
	}
	types := []EdgeType{EdgeGoverns, EdgeSpecifies, EdgeCoveredBy, EdgeTestedBy, EdgeReferences, EdgeDependsOn,
		EdgeSeeds, EdgeNeeds, EdgeRealizes, EdgeGatedBy}
	wantTypes := []EdgeType{"governs", "specifies", "covered-by", "tested-by", "references", "depends-on",
		"seeds", "needs", "realizes", "gated-by"}
	if !reflect.DeepEqual(types, wantTypes) {
		t.Errorf("edge types = %v, want %v", types, wantTypes)
	}
}

func TestEdgeNavigatesTo(t *testing.T) {
	t.Run("GRMDG-B04: The navigation between screens is an edge type of its own", func(t *testing.T) {})
	if EdgeNavigatesTo != "navigates-to" || EdgeNavigatesTo == EdgeDependsOn {
		t.Errorf("the navigation edge is its own type: %q", EdgeNavigatesTo)
	}
}
