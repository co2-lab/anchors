// @anchors
//   ref: FLRSA

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

// Only what carries no conclusion is listed, the most frequent first; --all adds the
// concluded ones with their conclusion.
func TestFailures_listsWhatTheSpecHasNotExplained(t *testing.T) {
	t.Run("FLRSA-B01: Only the ingested failures whose rule carries no conclusion are listed, with the spec that declares them", func(t *testing.T) {})
	t.Run("FLRSA-B03: With all the concluded failures are listed too, with their conclusion", func(t *testing.T) {})
	t.Run("FLRSA-I01: What the log ingestion binds to a spec is what the failures review lists", func(t *testing.T) {})
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
	t.Run("FLRSA-B04: An occurrence measured against another version of the spec is flagged", func(t *testing.T) {})
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
	t.Run("FLRSA-B05: With nothing open the review says so instead of printing an empty list", func(t *testing.T) {})
	t.Run("FLRSA-E01: Without a map the failures review fails", func(t *testing.T) {})
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

// The most frequent failure first, whatever order the map holds them in: it is the one
// that costs the most.
func TestFailures_mostFrequentFirst(t *testing.T) {
	t.Run("FLRSA-B02: The open failures are listed from the most frequent to the least", func(t *testing.T) {})
	root := fixtureProjectWith(t, "", strings.NewReplacer(" @resilient: the client refreshes and retries at once", "",
		" @observing: not the cache, not the clock", "").Replace(failingSpec))
	g := loadMap(t, root)
	for i := range g.Nodes {
		if g.Nodes[i].ID == "src/login.spec.md" {
			g.Nodes[i].Failures = []mapx.FailureSignal{
				{Rule: "LOGIN-E02", Count: 1},
				{Rule: "LOGIN-E03", Count: 2},
				{Rule: "LOGIN-E01", Count: 5},
			}
		}
	}
	if err := mapx.Save(g, filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	out := runCmd(t, newFailuresCmd(), "--root", root)
	e1, e3, e2 := strings.Index(out, "LOGIN-E01"), strings.Index(out, "LOGIN-E03"), strings.Index(out, "LOGIN-E02")
	if e1 < 0 || e3 < 0 || e2 < 0 || !(e1 < e3 && e3 < e2) {
		t.Errorf("expected E01 (5), E03 (2), E02 (1) in that order:\n%s", out)
	}
}

// A spec that can no longer be read is skipped: its occurrences are not listed, and the
// review still answers for the rest.
func TestFailures_unreadableSpecIsSkipped(t *testing.T) {
	t.Run("FLRSA-E02: A spec whose file can no longer be read is skipped by the review", func(t *testing.T) {})
	root, _ := logsProject(t)
	if err := os.Remove(filepath.Join(root, "src/login.spec.md")); err != nil {
		t.Fatal(err)
	}
	out := runCmd(t, newFailuresCmd(), "--root", root)
	if strings.Contains(out, "LOGIN-E01") || !strings.Contains(out, "no failure under observation") {
		t.Errorf("the unreadable spec must be skipped:\n%s", out)
	}
}

// The review asks for the conclusion; it does not record anything itself.
func TestFailures_writesNothing(t *testing.T) {
	t.Run("FLRSA-X01: The failures review leaves the map and the specs unchanged", func(t *testing.T) {})
	root, _ := logsProject(t)
	files := []string{filepath.Join(root, mapx.DefaultPath), filepath.Join(root, "src/login.spec.md")}
	before := map[string]string{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		before[f] = string(b)
	}
	runCmd(t, newFailuresCmd(), "--all", "--root", root)
	for _, f := range files {
		b, _ := os.ReadFile(f)
		if string(b) != before[f] {
			t.Errorf("the review changed %s", f)
		}
	}
}
