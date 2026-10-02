// @anchors
//   ref: GRPRG

package mapx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// THE `generated_by` FIELD MUST NOT CAUSE THE OSCILLATION IT EXISTS TO REVEAL.
//
// The local build writes "dev" and CI writes the published release. If every save rewrote
// the field, the two would undo each other on every run — exactly the defect the field was
// created to expose.
//
// The first implementation compared only the YAML against the file on disk, which carries the
// header: they were never equal, and the guard guarded nothing. Only an isolated test showed it.
func TestGeneratedByAloneDoesNotRewriteTheMap(t *testing.T) {
	t.Run("GRPRG-B03: A save that changes only the writer's release leaves the file untouched", func(t *testing.T) {})
	t.Run("GRPRG-B04: A real change rewrites the file with the running release", func(t *testing.T) {})
	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.graph.yaml")
	g := &Graph{Version: FormatoAtual, Nodes: []Node{{ID: "a", Rev: "r1"}}}

	defer func(old string) { GeneratedBy = old }(GeneratedBy)
	GeneratedBy = "0.1.10"
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(p)

	// Another binary, the SAME map: the file must stay intact.
	GeneratedBy = "dev"
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(p)
	if string(before) != string(after) {
		t.Error("only `generated_by` changed — rewriting makes the map flip between the local " +
			"and the published build, the defect the field exists to expose")
	}

	// A REAL change is written, the new generated_by included.
	g.Nodes = append(g.Nodes, Node{ID: "b", Rev: "r1"})
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	final, _ := os.ReadFile(p)
	if !strings.Contains(string(final), "generated_by: dev") {
		t.Error("when the map really changes, the writer recorded is the one running")
	}
}

// The FORMAT belongs to whoever WRITES. A `Graph{}` with no version used to be saved as
// `version: 0`, which reads back as "format 1: needs migrating" — a file born asking for a
// migration to the format it is already in.
func TestSaveStampsTheCurrentFormatAndTheHeader(t *testing.T) {
	t.Run("GRPRG-B01: Saving stamps the current format and the running binary's release", func(t *testing.T) {})
	t.Run("GRPRG-B02: The saved file starts with the fixed comment header", func(t *testing.T) {})
	defer func(old string) { GeneratedBy = old }(GeneratedBy)
	GeneratedBy = "0.1.10"
	p := filepath.Join(t.TempDir(), "anchors.graph.yaml")
	if err := Save(&Graph{Nodes: []Node{{ID: "a", Rev: "r1"}}}, p); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(p)
	text := string(data)
	if !strings.HasPrefix(text, "# anchors.graph.yaml — ") || !strings.Contains(strings.SplitN(text, "\n", 2)[0], "anchors map build") {
		t.Errorf("the file must open with the fixed comment header; got:\n%s", text)
	}
	if !strings.Contains(text, "\nversion: 6\n") || !strings.Contains(text, "\ngenerated_by: 0.1.10\n") {
		t.Errorf("the save must stamp the current format and the running release; got:\n%s", text)
	}
	if _, err := Load(p); err != nil {
		t.Errorf("a freshly saved map must load with no migration asked: %v", err)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Run("GRPRG-I01: A saved graph loads back unchanged", func(t *testing.T) {})
	defer func(old string) { GeneratedBy = old }(GeneratedBy)
	GeneratedBy = "dev"
	g := &Graph{
		Nodes: []Node{
			{ID: "a.spec.md", Kind: KindSpec, Rev: "r1", Code: "AAAAA", CodeDeclarado: true},
			{ID: "a.go", Kind: KindCode, Rev: "r2", Signal: &TestSignal{CoveredLines: 3, TotalLines: 4, LineCoverage: 75, AtRev: "r2"}},
		},
		Edges: []Edge{{
			From: "a.spec.md", To: "a.go", Type: EdgeSpecifies, Origin: OriginConvention,
			Stamp:       &Stamp{ValidatedFromRev: "r1", ValidatedToRev: "r2", ChangedAt: "2026-09-26", Verdict: "ok"},
			Julgamentos: []Judgment{{Gate: "g", Verdict: "ok", ValidatedFromRev: "r1", ValidatedToRev: "r2", ChangedAt: "2026-09-26"}},
		}},
	}
	p := filepath.Join(t.TempDir(), "anchors.graph.yaml")
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	back, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, g) {
		t.Errorf("the graph changed on the way through the file:\nsaved  %+v\nloaded %+v", g, back)
	}
}

// The FORMAT is checked BEFORE the map is used: a file this binary cannot read must not be
// partially interpreted, or the next save writes only what survived.
func TestLoadRefusesAnUnreadableFormat(t *testing.T) {
	t.Run("GRPRG-B05: Loading a map in an unreadable format is refused", func(t *testing.T) {})
	t.Run("GRPRG-X01: An unreadable map yields no graph at all", func(t *testing.T) {})
	dir := t.TempDir()
	for _, v := range []string{"1", "9"} {
		p := filepath.Join(dir, "m"+v+".yaml")
		os.WriteFile(p, []byte("version: "+v+"\nnodes:\n  - id: a\n    kind: code\n    rev: r1\nedges: []\n"), 0o644)
		g, err := Load(p)
		var fe *ErroDeFormato
		if !errors.As(err, &fe) {
			t.Errorf("format %s: expected the format refusal, got %v", v, err)
		}
		if g != nil {
			t.Errorf("format %s: no graph may come back from an unreadable map, got %+v", v, g)
		}
	}
}

func TestLoadReturnsReadAndParseErrors(t *testing.T) {
	t.Run("GRPRG-E01: Loading a missing file returns the read error", func(t *testing.T) {})
	t.Run("GRPRG-E02: Loading text that is not the map returns the parse error", func(t *testing.T) {})
	dir := t.TempDir()
	if g, err := Load(filepath.Join(dir, "absent.yaml")); err == nil || g != nil {
		t.Errorf("a missing file must return an error and no graph, got %v / %v", g, err)
	}
	bad := filepath.Join(dir, "bad.yaml")
	os.WriteFile(bad, []byte("nodes: [unclosed\n"), 0o644)
	if g, err := Load(bad); err == nil || g != nil {
		t.Errorf("text that is not the map must return an error and no graph, got %v / %v", g, err)
	}
}

func TestSaveReturnsTheWriteError(t *testing.T) {
	t.Run("GRPRG-E04: Saving into a missing directory returns the write error", func(t *testing.T) {})
	p := filepath.Join(t.TempDir(), "missing", "anchors.graph.yaml")
	if err := Save(&Graph{}, p); err == nil {
		t.Error("a save into a directory that does not exist must return the error")
	}
}

func TestLoadBytes(t *testing.T) {
	t.Run("GRPRG-B06: A map read from bytes is read like one from disk", func(t *testing.T) {})
	g, err := LoadBytes([]byte(fmt.Sprintf("version: %d\nnodes:\n  - id: a\n    kind: code\n    rev: r1\nedges: []\n", FormatoAtual)), "HEAD:m")
	if err != nil || len(g.Nodes) != 1 || g.Nodes[0].Rev != "r1" {
		t.Fatalf("a current map reads, got %+v %v", g, err)
	}
	g, err = LoadBytes([]byte("version: 1\nnodes: []\n"), "HEAD:m")
	var fe *ErroDeFormato
	if g != nil || !errors.As(err, &fe) || !strings.Contains(err.Error(), "HEAD:m") {
		t.Errorf("an unreadable format is refused naming where it came from, got %+v %v", g, err)
	}
	if _, err := LoadBytes([]byte(":\n-"), "x"); err == nil {
		t.Error("bytes that are not a map are refused")
	}
}

func TestLoad_emptySignalIsNoSignal(t *testing.T) {
	t.Run("GRPRG-B07: An empty signal loads as no signal", func(t *testing.T) {})
	g, err := LoadBytes([]byte(fmt.Sprintf("version: %d\nnodes:\n  - id: a\n    kind: spec\n    rev: r1\n    signal: {}\n  - id: b\n    kind: test\n    rev: r1\n    signal:\n      passed: 2\nedges: []\n", FormatoAtual)), "m")
	if err != nil {
		t.Fatal(err)
	}
	if g.Nodes[0].Signal != nil || g.Nodes[1].Signal == nil || g.Nodes[1].Signal.Passed != 2 {
		t.Errorf("the empty signal is dropped and the measured one kept, got %+v %+v", g.Nodes[0].Signal, g.Nodes[1].Signal)
	}
}
