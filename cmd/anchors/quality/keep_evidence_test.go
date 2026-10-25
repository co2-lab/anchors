// @anchors
//   code: KETKP
//   ref: KPEVD

package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

const keptSpec = "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\nRealizes the payments doctrine.\n\n### PAYMX-B01 — charges @realizes GSSL-R01\n"

// keepRepo is a repository whose spec was proven, then gained a `@realizes` and a sentence
// saying so — outside its rules, so the rebuild drops the proof.
func keepRepo(t *testing.T, rebuild bool) (syncRepo, string) {
	t.Helper()
	r := newSyncRepo(t, true)
	touchWrite(t, r.root, "src/pay.spec.md", keptSpec)
	mapPath := filepath.Join(r.root, mapx.DefaultPath)
	if rebuild {
		old, err := mapx.Load(mapPath)
		if err != nil {
			t.Fatal(err)
		}
		files, _ := scan.Walk(r.root, r.cfg)
		g := mapx.Build(files, r.cfg, gitmeta.AllCommitDates(r.root))
		mapx.PreserveStamps(g, old)
		if err := mapx.Save(g, mapPath); err != nil {
			t.Fatal(err)
		}
	}
	return r, mapPath
}

func keptSignal(t *testing.T, mapPath string) (*mapx.Node, string) {
	t.Helper()
	g, err := mapx.Load(mapPath)
	if err != nil {
		t.Fatal(err)
	}
	return g.Node("src/pay.spec.md"), scan.ShortHash([]byte(keptSpec))
}

func TestKeepEvidence_movesWithTheReason(t *testing.T) {
	t.Run("KPEVD-B01: The evidence moves to the current content with the reason", func(t *testing.T) {})
	r, mapPath := keepRepo(t, false)
	out, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "only @realizes added")
	if err != nil {
		t.Fatal(err)
	}
	n, rev := keptSignal(t, mapPath)
	if n.Rev != rev || n.Signal.AtRev != rev || len(n.Signal.ProvenCodes) != 1 || n.EvidenceKept[0].Reason != "only @realizes added" {
		t.Errorf("the proof is at the current content with the reason, got %+v %+v", n, n.Signal)
	}
	if !strings.Contains(out, "src/pay.spec.md — evidence kept:") || !strings.Contains(out, "1 file(s) kept their evidence") {
		t.Errorf("the command says what it carried:\n%s", out)
	}
}

func TestKeepEvidence_fromHead(t *testing.T) {
	t.Run("KPEVD-B02: A rebuilt map's lost evidence comes from HEAD", func(t *testing.T) {})
	r, mapPath := keepRepo(t, true)
	if n, _ := keptSignal(t, mapPath); n.Signal != nil {
		t.Fatalf("the rebuild must have dropped the proof, got %+v", n.Signal)
	}
	if _, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "only @realizes added"); err != nil {
		t.Fatal(err)
	}
	if n, rev := keptSignal(t, mapPath); n.Signal == nil || n.Signal.AtRev != rev || n.Signal.ProvenCodes[0] != "PAYMX-B01" {
		t.Errorf("the proof comes back from HEAD at the current content, got %+v", n.Signal)
	}
	if headMap(t.TempDir()) != nil {
		t.Error("no repository, no map at HEAD")
	}
}

func TestKeepEvidence_nothingToKeep(t *testing.T) {
	t.Run("KPEVD-B03: Nothing to keep is said", func(t *testing.T) {})
	r, _ := keepRepo(t, false)
	if _, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "x"); err != nil {
		t.Fatal(err)
	}
	out, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "again")
	if err != nil || !strings.Contains(out, "src/pay.spec.md — nothing to keep") || !strings.Contains(out, "0 file(s) kept") ||
		strings.Contains(out, "refreshed under the same declaration") {
		t.Errorf("a second run has nothing to keep, got %v:\n%s", err, out)
	}
}

func TestKeepEvidence_refreshesTheStamps(t *testing.T) {
	t.Run("KPEVD-B04: The contract stamps are refreshed under the same declaration", func(t *testing.T) {})
	r, _ := keepRepo(t, false)
	out, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "x")
	if err != nil || !strings.Contains(out, "refreshed under the same declaration") ||
		!strings.Contains(out, "src/pay.spec.md — no double is stamped against a previous version of it.") {
		t.Errorf("the stamps pointing at the file are refreshed and listed, got %v:\n%s", err, out)
	}
}

func TestKeepEvidence_refusals(t *testing.T) {
	t.Run("KPEVD-E01: A missing reason, file or map fails", func(t *testing.T) {})
	r, mapPath := keepRepo(t, false)
	before, _ := os.ReadFile(mapPath)
	if _, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "  "); err == nil || !strings.Contains(err.Error(), "--reason is required") {
		t.Errorf("a blank reason fails, got %v", err)
	}
	if _, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "src/ghost.ts", "--reason", "x"); err == nil || !strings.Contains(err.Error(), "src/ghost.ts") {
		t.Errorf("an unreadable file fails naming it, got %v", err)
	}
	touchWrite(t, r.root, "notes.md", "x")
	if _, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "notes.md", "--reason", "x"); err == nil || !strings.Contains(err.Error(), "not in the map") {
		t.Errorf("a file outside the map fails, got %v", err)
	}
	if after, _ := os.ReadFile(mapPath); string(after) != string(before) {
		t.Error("a refusal writes nothing")
	}
	if _, err := runQ(t, newKeepEvidenceCmd(), "--root", t.TempDir(), "a.md", "--reason", "x"); err == nil || !strings.Contains(err.Error(), "map build") {
		t.Errorf("no map points at the map build, got %v", err)
	}
}

func TestKeepEvidence_suitesNotRunAgainComeFromHead(t *testing.T) {
	t.Run("KPEVD-B05: The suites a later run did not replace come back from HEAD and are carried with the rest", func(t *testing.T) {})
	r, mapPath := keepRepo(t, false)
	// HEAD's map: the spec proven by the unit suite and by an e2e suite, at the old content.
	g, err := mapx.Load(mapPath)
	if err != nil {
		t.Fatal(err)
	}
	n := g.Node("src/pay.spec.md")
	old := n.Rev
	n.Signal = &mapx.TestSignal{
		ProvenCodes:      []string{"PAYMX-B01"},
		ProvenBySuite:    map[string][]string{"unit.xml": {"PAYMX-B01"}, "e2e/flow.xml": {"PAYMX-B01"}},
		ProvenRevBySuite: map[string]string{"unit.xml": old, "e2e/flow.xml": old},
		AtRev:            old,
	}
	if err := mapx.Save(g, mapPath); err != nil {
		t.Fatal(err)
	}
	r.git("add", mapx.DefaultPath)
	r.git("commit", "-qm", "proven", "--", mapx.DefaultPath)
	// The map is rebuilt at the new content, and only the unit suite runs again.
	_, newRev := keptSignal(t, mapPath)
	g, _ = mapx.Load(mapPath)
	n = g.Node("src/pay.spec.md")
	n.Rev = newRev
	n.Signal = &mapx.TestSignal{
		ProvenCodes:      []string{"PAYMX-B01"},
		ProvenBySuite:    map[string][]string{"unit.xml": {"PAYMX-B01"}},
		ProvenRevBySuite: map[string]string{"unit.xml": newRev},
		AtRev:            newRev,
	}
	if err := mapx.Save(g, mapPath); err != nil {
		t.Fatal(err)
	}
	out, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "only VR metadata")
	if err != nil {
		t.Fatal(err)
	}
	got, _ := keptSignal(t, mapPath)
	s := got.Signal
	if s == nil || len(s.ProvenBySuite) != 2 || s.ProvenRevBySuite["e2e/flow.xml"] != newRev || s.ProvenRevBySuite["unit.xml"] != newRev ||
		s.AtRev != newRev || !strings.Contains(out, "evidence kept:") {
		t.Errorf("the e2e proof comes back from HEAD and is at the new content beside the unit run: %+v\n%s", s, out)
	}
}

func TestKeepEvidence_clearsTheScenariosOfAChangedRule(t *testing.T) {
	t.Run("KPEVD-B06: The scenarios of a changed rule, marked stale by the rebuild, are kept too", func(t *testing.T) {})
	r, mapPath := keepRepo(t, false)
	g, err := mapx.Load(mapPath)
	if err != nil {
		t.Fatal(err)
	}
	n := g.Node("src/pay.spec.md")
	rev := scan.ShortHash([]byte(keptSpec))
	n.Rev = rev
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"PAYMX-B01"}, AtRev: rev, StaleCodes: []string{"PAYMX-B01"}}
	if err := mapx.Save(g, mapPath); err != nil {
		t.Fatal(err)
	}
	out, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "the rule's wording only")
	if err != nil || !strings.Contains(out, "src/pay.spec.md — scenarios of a changed rule kept: PAYMX-B01") || !strings.Contains(out, "1 file(s) kept") {
		t.Fatalf("the command says what it kept, got %v:\n%s", err, out)
	}
	if n, _ := keptSignal(t, mapPath); len(n.Signal.StaleCodes) != 0 {
		t.Errorf("nothing is stale any more, got %+v", n.Signal)
	}
}

func TestKeepEvidence_survivesADateAfterwards(t *testing.T) {
	t.Run("KPEVD-B07: The evidence kept survives a later change that proves nothing, as a header's date", func(t *testing.T) {})
	// The map still at the content before the change: the command moves the node's revision
	// and must move its evidence revisions with it, or the next rebuild drops what it kept.
	r, mapPath := keepRepo(t, false)
	if _, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.spec.md", "--reason", "only @realizes added"); err != nil {
		t.Fatal(err)
	}
	dated := strings.Replace(keptSpec, "updated_at: 2026-09-01", "updated_at: 2026-10-24", 1)
	touchWrite(t, r.root, "src/pay.spec.md", dated)
	old, err := mapx.Load(mapPath)
	if err != nil {
		t.Fatal(err)
	}
	files, _ := scan.Walk(r.root, r.cfg)
	g := mapx.Build(files, r.cfg, gitmeta.AllCommitDates(r.root))
	mapx.PreserveStamps(g, old)
	rev := scan.ShortHash([]byte(dated))
	if n := g.Node("src/pay.spec.md"); n.Rev != rev || n.Signal == nil || n.Signal.AtRev != rev || len(n.Signal.ProvenCodes) != 1 {
		t.Errorf("the proof kept goes along with the date, got %+v %+v", n, n.Signal)
	}
}

func TestKeepEvidence_mutationNotRunAgainComesFromHead(t *testing.T) {
	t.Run("KPEVD-B08: A mutation a later run did not replace comes back from HEAD and is carried with the rest", func(t *testing.T) {})
	r, mapPath := keepRepo(t, false)
	// HEAD's map: the code file's mutation, at its old content.
	g, err := mapx.Load(mapPath)
	if err != nil {
		t.Fatal(err)
	}
	n := g.Node("src/pay.ts")
	old := n.Rev
	n.Signal = &mapx.TestSignal{MutantsKilled: 9, MutantsSurvived: 1, MutationScore: 90, AtRev: old, MutationAtRev: old,
		MutationByScope: map[string]mapx.MutationScope{"isolated": {Killed: 9, Survived: 1, Score: 90, AtRev: old}}}
	if err := mapx.Save(g, mapPath); err != nil {
		t.Fatal(err)
	}
	r.git("add", mapx.DefaultPath)
	r.git("commit", "-qm", "mutated", "--", mapx.DefaultPath)
	// Its import changes, the map is rebuilt, and only the coverage runs again: a new signal
	// with no mutation.
	edited := "export const pay = 1 // imports moved\n"
	touchWrite(t, r.root, "src/pay.ts", edited)
	rev := scan.ShortHash([]byte(edited))
	g, _ = mapx.Load(mapPath)
	n = g.Node("src/pay.ts")
	n.Rev = rev
	n.Signal = &mapx.TestSignal{CoveredLines: 1, TotalLines: 1, LineCoverage: 100, AtRev: rev,
		CoverageBySuite: map[string]mapx.SuiteCoverage{"lcov.info": {CoveredLines: 1, TotalLines: 1, AtRev: rev}}}
	if err := mapx.Save(g, mapPath); err != nil {
		t.Fatal(err)
	}
	out, err := runQ(t, newKeepEvidenceCmd(), "--root", r.root, "src/pay.ts", "--lines", "--reason", "only the import specifiers moved")
	if err != nil {
		t.Fatal(err)
	}
	g, _ = mapx.Load(mapPath)
	s := g.Node("src/pay.ts").Signal
	if s == nil || s.MutationAtRev != rev || s.MutationScore != 90 || s.MutationByScope["isolated"].AtRev != rev ||
		s.CoverageBySuite["lcov.info"].AtRev != rev || !strings.Contains(out, "evidence kept:") {
		t.Errorf("the mutation comes back from HEAD at the new content beside the coverage run: %+v\n%s", s, out)
	}
	// A mutation measured again at the new content stays its own.
	s.MutationScore, s.MutationAtRev = 50, rev
	if s.AdoptMissingMutation(&mapx.TestSignal{MutationScore: 90, MutationAtRev: old}) || s.MutationScore != 50 {
		t.Errorf("a signal with a mutation of its own adopts none, got %+v", s)
	}
}
