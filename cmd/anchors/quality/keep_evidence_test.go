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

const keptSpec = "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\n### PAYMX-B01 — charges @realizes GSSL-R01\n"

// keepRepo is a repository whose spec was proven, then gained a `@realizes`.
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
	if n.Rev != rev || n.Signal.AtRev != rev || len(n.Signal.ProvenCodes) != 1 || n.Signal.EvidenceKept[0].Reason != "only @realizes added" {
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
