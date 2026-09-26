package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// A spec that declares three failures: one still open, one concluded resilient, one
// under observation.
const failingSpec = "<!-- @anchors\n  code: LOGIN\n-->\n# Login\n\nLOGIN-B01 — the password is checked.\n\n" +
	"### LOGIN-E01 — the password store is down\n\n" +
	"### LOGIN-E02 — the token expired @resilient: the client refreshes and retries at once\n\n" +
	"### LOGIN-E03 — the session vanished @observing: not the cache, not the clock\n"

const logsYAML = "logs:\n  paths: [\"logs/*.log\"]\n"

// logsProject builds the fixture with the failing spec, writes a log and ingests it with
// the real `ingest --logs`.
func logsProject(t *testing.T) (string, string) {
	t.Helper()
	root := fixtureProjectWith(t, logsYAML, failingSpec)
	writeProjectFile(t, root, "logs/app.log", strings.Join([]string{
		"10:00 ERROR #[LOGIN-E01] store unreachable",
		"10:01 ERROR #[LOGIN-E01] store unreachable",
		"10:02 ERROR #[LOGIN-E01] store unreachable",
		"10:03 WARN  #[LOGIN-E02] token expired",
		"10:04 ERROR #[LOGIN-E03] session gone",
		"10:05 ERROR #[GHOST-E07] nobody declared this",
		"10:06 INFO  all good",
	}, "\n")+"\n")
	out := runCmd(t, newIngestCmd(), "--logs", "--root", root)
	return root, out
}

// The occurrences are bound to the spec that declares them; a code no spec declares is
// reported, and left out of the map.
func TestIngestLogs_bindsOccurrencesToTheSpec(t *testing.T) {
	root, out := logsProject(t)
	if !strings.Contains(out, "logs: 1 file(s), 7 line(s) — 3 occurrence(s) bound to 3 spec rule(s)") {
		t.Errorf("unexpected summary:\n%s", out)
	}
	if !strings.Contains(out, "GHOST-E07") || !strings.Contains(out, "NO spec declares") {
		t.Errorf("the undeclared failure code was not reported:\n%s", out)
	}
	var spec *mapx.Node
	g := loadMap(t, root)
	for i := range g.Nodes {
		if g.Nodes[i].ID == "src/login.spec.md" {
			spec = &g.Nodes[i]
		}
	}
	if spec == nil || len(spec.Failures) != 3 {
		t.Fatalf("expected 3 failures on the spec node, got %+v", spec)
	}
	for _, f := range spec.Failures {
		if f.Rev != spec.Rev {
			t.Errorf("%s was not stamped with the spec's rev (%q vs %q)", f.Rule, f.Rev, spec.Rev)
		}
		if f.Rule == "LOGIN-E01" && f.Count != 3 {
			t.Errorf("LOGIN-E01 happened 3 times, map says %d", f.Count)
		}
	}
}

func TestIngestLogs_errors(t *testing.T) {
	root := fixtureProject(t) // no `logs:` declared
	if _, err := runCmdErr(newIngestCmd(), t, "--logs", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "logs.paths") {
		t.Errorf("no log declared: expected an error naming logs.paths, got %v", err)
	}
	if _, err := runCmdErr(newIngestCmd(), t, "--root", root); err == nil ||
		!strings.Contains(err.Error(), "--junit") {
		t.Errorf("nothing to ingest: expected the usage error, got %v", err)
	}
}

// Only what carries no conclusion is listed, the most frequent first; --all adds the
// concluded ones with their conclusion.
func TestFailures_listsWhatTheSpecHasNotExplained(t *testing.T) {
	root, _ := logsProject(t)

	out := runCmd(t, newFailuresCmd(), "--root", root)
	if !strings.Contains(out, "failures observed and not yet understood (1):") || !strings.Contains(out, "LOGIN-E01") {
		t.Errorf("LOGIN-E01 is the one open failure:\n%s", out)
	}
	if !strings.Contains(out, "declared in src/login.spec.md") {
		t.Errorf("the spec that declares it is not named:\n%s", out)
	}
	for _, concluded := range []string{"LOGIN-E02", "LOGIN-E03"} {
		if strings.Contains(out, concluded) {
			t.Errorf("%s carries a conclusion and was listed without --all:\n%s", concluded, out)
		}
	}

	all := runCmd(t, newFailuresCmd(), "--all", "--root", root)
	if !strings.Contains(all, "resilient: the client refreshes and retries at once") {
		t.Errorf("--all does not show the resilient conclusion:\n%s", all)
	}
	if !strings.Contains(all, "under observation: not the cache, not the clock") {
		t.Errorf("--all does not show what was ruled out:\n%s", all)
	}
}

// An occurrence measured against an older version of the spec is flagged: concluding
// about it would assert about what was not observed.
func TestFailures_flagsOccurrencesOfAnOlderSpec(t *testing.T) {
	root, _ := logsProject(t)
	g := loadMap(t, root)
	for i := range g.Nodes {
		for j := range g.Nodes[i].Failures {
			g.Nodes[i].Failures[j].Rev = "older"
		}
	}
	if err := mapx.Save(g, filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	out := runCmd(t, newFailuresCmd(), "--root", root)
	if !strings.Contains(out, "the spec changed since the ingestion") {
		t.Errorf("the stale occurrence was not flagged:\n%s", out)
	}
}

// Everything concluded (or nothing ingested): the command says so instead of an empty list.
func TestFailures_noneOpen(t *testing.T) {
	root, _ := logsProject(t)
	resolved := strings.Replace(failingSpec, "store is down", "store is down @resilient: the pool reconnects", 1)
	if err := os.WriteFile(filepath.Join(root, "src/login.spec.md"), []byte(resolved), 0o644); err != nil {
		t.Fatal(err)
	}
	out := runCmd(t, newFailuresCmd(), "--root", root)
	if !strings.Contains(out, "no failure under observation") {
		t.Errorf("with every failure concluded, expected the none message:\n%s", out)
	}
	if _, err := runCmdErr(newFailuresCmd(), t, "--root", t.TempDir()); err == nil {
		t.Error("without a map the command must fail")
	}
}
