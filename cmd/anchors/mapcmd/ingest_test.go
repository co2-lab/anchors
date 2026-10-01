package mapcmd

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testsig"
)

const junitReport = `<?xml version="1.0"?>
<testsuites>
  <testsuite name="login" file="src/login.test.ts">
    <testcase name="LOGIN-B01: checks the password" file="src/login.test.ts"/>
    <testcase name="LOGIN-B02: locks after three tries" file="src/login.test.ts"><failure>boom</failure></testcase>
    <testcase name="later" file="src/login.test.ts"><skipped/></testcase>
  </testsuite>
</testsuites>
`

const lcovReport = "TN:\nSF:src/login.ts\nDA:1,1\nDA:2,0\nDA:3,1\nDA:4,1\nLF:4\nLH:3\nend_of_record\n"

const mutationReport = `{"schemaVersion":"1","thresholds":{"high":80,"low":60},"files":{
  "src/login.ts":{"mutants":[{"status":"Killed","location":{"start":{"line":1}}},
                             {"status":"Killed","location":{"start":{"line":3}}},
                             {"status":"Survived","location":{"start":{"line":4}}}]}}}`

func node(t *testing.T, root, id string) mapx.Node {
	t.Helper()
	for _, n := range loadMap(t, root).Nodes {
		if n.ID == id {
			return n
		}
	}
	t.Fatalf("node %s not in the map", id)
	return mapx.Node{}
}

// The three reports in one pass: execution binds to the test node and proves the
// scenario the spec declares; coverage and mutation bind to the code node.
func TestIngest_junitLcovAndMutationReachTheMap(t *testing.T) {
	t.Run("NGSTI-B01: The three reports in one pass reach the test, spec and code nodes", func(t *testing.T) {})
	root := fixtureProject(t)
	writeProjectFile(t, root, "reports/junit.xml", junitReport)
	writeProjectFile(t, root, "reports/lcov.info", lcovReport)
	writeProjectFile(t, root, "reports/mutation.json", mutationReport)

	out := runCmd(t, newIngestCmd(), "--root", root,
		"--junit", filepath.Join(root, "reports/junit.xml"),
		"--lcov", filepath.Join(root, "reports/lcov.info"),
		"--mutation", filepath.Join(root, "reports/mutation.json"))
	for _, want := range []string{
		"execution: 3 case(s), 1 test file(s) matched, 1 scenario(s) proven",
		"coverage: 1 file(s) in the lcov, 1 code node(s) matched",
		"1 surviving mutant(s)",
		"signals written into the map",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}

	test := node(t, root, "src/login.test.ts").Signal
	if test == nil || test.Passed != 1 || test.Failed != 1 || test.Skipped != 1 {
		t.Errorf("the execution did not reach the test node: %+v", test)
	}
	spec := node(t, root, "src/login.spec.md").Signal
	if spec == nil || strings.Join(spec.ProvenCodes, ",") != "LOGIN-B01" {
		t.Errorf("LOGIN-B01 passed and must be the one proven code: %+v", spec)
	}
	code := node(t, root, "src/login.ts").Signal
	if code == nil || code.CoveredLines != 3 || code.TotalLines != 4 {
		t.Errorf("the lcov did not reach the code node: %+v", code)
	}
	if code.MutantsKilled != 2 || code.MutantsSurvived != 1 {
		t.Errorf("the mutation report did not reach the code node: %+v", code)
	}
}

// A report that names no file of the map says so, instead of silently proving nothing.
func TestIngest_junitThatMatchesNothingWarns(t *testing.T) {
	t.Run("NGSTI-B02: A JUnit report that matches no test node warns", func(t *testing.T) {})
	root := fixtureProject(t)
	writeProjectFile(t, root, "reports/junit.xml", strings.ReplaceAll(junitReport, "src/login.test.ts", "elsewhere/x.test.ts"))
	out := runCmd(t, newIngestCmd(), "--root", root, "--junit", filepath.Join(root, "reports/junit.xml"))
	if !strings.Contains(out, "no test file matched") {
		t.Errorf("a report matching no node must warn:\n%s", out)
	}
}

// A full run of a suite INSIDE the repository replaces what an ad-hoc report from outside
// it left in the map.
func TestIngest_fullRunDropsTheExternalSuite(t *testing.T) {
	t.Run("NGSTI-B08: A full run of a suite inside the repository drops the suites ingested from outside it", func(t *testing.T) {})
	root := fixtureProject(t)
	external := filepath.Join(t.TempDir(), "junit.xml")
	if err := os.WriteFile(external, []byte(junitReport), 0o644); err != nil {
		t.Fatal(err)
	}
	runCmd(t, newIngestCmd(), "--root", root, "--junit", external)

	writeProjectFile(t, root, "reports/junit.xml", junitReport)
	out := runCmd(t, newIngestCmd(), "--root", root, "--junit", filepath.Join(root, "reports/junit.xml"))
	if !strings.Contains(out, "dropped the proof of 1 suite(s) ingested from outside the repository") {
		t.Errorf("the external suite was not dropped:\n%s", out)
	}
	spec := node(t, root, "src/login.spec.md").Signal
	for suite := range spec.ProvenBySuite {
		if strings.HasPrefix(suite, mapx.ExternalSuitePrefix) {
			t.Errorf("the external suite is still in the map: %v", spec.ProvenBySuite)
		}
	}
}

func TestIngest_badReportsAreErrors(t *testing.T) {
	t.Run("NGSTI-E04: A missing or malformed report fails naming the report's format", func(t *testing.T) {})
	t.Run("NGSTI-E05: Ingesting a report without a map fails and asks for the map build", func(t *testing.T) {})
	root := fixtureProject(t)
	bad := filepath.Join(root, "bad.txt")
	writeProjectFile(t, root, "bad.txt", "{ not a report")
	for flag, want := range map[string]string{"--junit": "parse JUnit", "--lcov": "parse lcov", "--mutation": "parse mutation"} {
		if _, err := runCmdErr(newIngestCmd(), t, "--root", root, flag, bad+".missing"); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s with a missing file: expected %q, got %v", flag, want, err)
		}
	}
	if _, err := runCmdErr(newIngestCmd(), t, "--root", t.TempDir(), "--junit", bad); err == nil ||
		!strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: got %v", err)
	}
}

// A project that declares its code length has its JUnit cases read with it: a `[6]`
// project names its rules `LOGINX-B01`, and the permissive default (4 and 5) never read
// one of them — the rule showed no green test with its test passing.
func TestIngest_readsTheDeclaredCodeLength(t *testing.T) {
	t.Run("NGSTI-B03: A project's declared code length governs how JUnit case names are read", func(t *testing.T) {})
	prevLen := config.CodeLengths
	t.Cleanup(func() {
		config.SetCodeLengths(prevLen)
		testsig.SetCodeLenPattern("{4,5}")
	})
	useEnglish(t)
	root := t.TempDir()
	for p, c := range map[string]string{
		"anchors.yaml":      strings.Replace(fixtureYAML, "version: 4\n", "version: 4\ncode_lengths: [6]\n", 1),
		"src/login.spec.md": "<!-- @anchors\n  code: LOGINX\n-->\n# Login\n\nLOGINX-B01 — the password is checked.\n",
		"src/login.ts":      "export const login = () => true\n",
		"src/login.test.ts": "test('LOGINX-B01: checks the password', () => {})\n",
		"guides/CODE.md":    "# Code guide\n",
	} {
		writeProjectFile(t, root, p, c)
	}
	runCmd(t, newMapCmd(), "build", "--root", root)
	writeProjectFile(t, root, "reports/junit.xml", `<?xml version="1.0"?>
<testsuites><testsuite name="login" file="src/login.test.ts">
  <testcase name="LOGINX-B01: checks the password" file="src/login.test.ts"/>
</testsuite></testsuites>
`)
	runCmd(t, newIngestCmd(), "--root", root, "--junit", filepath.Join(root, "reports/junit.xml"))
	spec := node(t, root, "src/login.spec.md").Signal
	if spec == nil || strings.Join(spec.ProvenCodes, ",") != "LOGINX-B01" {
		t.Fatalf("a [6] project's passing case must prove its rule: %+v", spec)
	}
}

func projectWith(t *testing.T, yaml string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

const withSuite = `derived:
  anchor: spec
tests:
  - layer: unit
    run: pnpm test
    junit: .test-results/junit.xml
`

// `anchors test` runs the suite AND ingests in one operation — that is what guarantees
// the signal in the map matches the run that just happened.
//
// Calling `ingest` directly breaks the guarantee unnoticed: one can run the suite, edit
// the code, ingest the old report, and the map asserts a coverage that no longer holds.
func TestIngestManual_warnsButDoesNotBlock(t *testing.T) {
	t.Run("NGSTI-B04: Manual ingestion warns and proceeds by default", func(t *testing.T) {})
	ViaAnchorsTest = false
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	prev := os.Stderr
	os.Stderr = w
	err = warnIfManualIngest(projectWith(t, withSuite))
	os.Stderr = prev
	w.Close()
	warned, _ := io.ReadAll(r)
	if err != nil {
		t.Errorf("the default WARNS and lets it through: there are legitimate uses (a CI that "+
			"ran the suite in another job), and blocking them would leave no way out; got: %v", err)
	}
	if !strings.Contains(string(warned), "MANUAL ingestion") {
		t.Errorf("the manual ingestion must be warned about on stderr, got %q", warned)
	}
}

// When the project declares it, manual ingestion is REFUSED — and the message must say
// what to use instead, or whoever was refused does not know how to go on.
func TestIngestManual_blocksWhenTheProjectAsks(t *testing.T) {
	t.Run("NGSTI-E02: Manual ingestion is refused when the project declares manual ingestion blocks", func(t *testing.T) {})
	ViaAnchorsTest = false
	yaml := strings.Replace(withSuite, "derived:", "workflow:\n  manual_ingest_blocks: true\nderived:", 1)
	err := warnIfManualIngest(projectWith(t, yaml))
	if err == nil {
		t.Fatal("with `manual_ingest_blocks: true` manual ingestion must be refused")
	}
	if !strings.Contains(err.Error(), "anchors test") {
		t.Errorf("the error must say WHAT TO USE instead; got: %v", err)
	}
}

// Coming from `anchors test`, it never complains — not even when the project declares
// the block. If it did, the right command would be barred by the rule that exists to
// promote it.
func TestIngestManual_viaAnchorsTestNeverComplains(t *testing.T) {
	t.Run("NGSTI-B05: Ingestion run by anchors test never complains", func(t *testing.T) {})
	ViaAnchorsTest = true
	defer func() { ViaAnchorsTest = false }()
	yaml := strings.Replace(withSuite, "derived:", "workflow:\n  manual_ingest_blocks: true\nderived:", 1)
	if err := warnIfManualIngest(projectWith(t, yaml)); err != nil {
		t.Errorf("`anchors test` is the RIGHT path — it cannot be barred: %v", err)
	}
}

// WITHOUT `tests:` declared, `anchors test` does not run, and demanding what does not
// exist would leave the project nowhere to go: `ingest` is the only way the signal
// reaches the map. Without a config there is nothing to demand either.
func TestIngestManual_noDeclaredSuiteDemandsNothing(t *testing.T) {
	t.Run("NGSTI-B06: A project with no declared suite or no config is not asked to use anchors test", func(t *testing.T) {})
	ViaAnchorsTest = false
	noSuite := "derived:\n  anchor: spec\nworkflow:\n  manual_ingest_blocks: true\n"
	if err := warnIfManualIngest(projectWith(t, noSuite)); err != nil {
		t.Errorf("with no declared suite there is no alternative to demand; got: %v", err)
	}
	if err := warnIfManualIngest(t.TempDir()); err != nil {
		t.Errorf("with no config there is nothing to demand; got: %v", err)
	}
}

// The report's mtime against the file's: a file changed after the report was written is
// not the file the report measured. The old integration lcov ingested today was stamped
// fresh, and its old line numbers diluted the fresh unit coverage (97% read as 69%).
func TestMarkPredating_fileEditedAfterTheReport(t *testing.T) {
	t.Run("NGSTI-B11: An lcov entry for a file edited after the report was written is marked as predating it", func(t *testing.T) {})
	root := t.TempDir()
	report := filepath.Join(root, "lcov.info")
	edited := filepath.Join(root, "edited.ts")
	untouched := filepath.Join(root, "untouched.ts")
	for _, p := range []string{report, edited, untouched} {
		if err := os.WriteFile(p, []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t0 := time.Now().Add(-2 * time.Hour)
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.Chtimes(untouched, t0, t0))
	must(os.Chtimes(report, t0.Add(time.Minute), t0.Add(time.Minute)))
	must(os.Chtimes(edited, t0.Add(time.Hour), t0.Add(time.Hour)))

	byFile := map[string]mapx.FileCov{"edited.ts": {Total: 1}, "untouched.ts": {Total: 1}, "gone.ts": {Total: 1}}
	markPredating(byFile, root, report)
	if !byFile["edited.ts"].Predates {
		t.Error("a file edited after the report was not flagged")
	}
	if byFile["untouched.ts"].Predates {
		t.Error("a file older than the report was flagged")
	}
	if byFile["gone.ts"].Predates {
		t.Error("a file that is not on disk was flagged")
	}
}

// The suite key must be the same on every machine. A report OUTSIDE the repository
// (CI writing `/tmp/junit.xml`) used to become `../../../../tmp/junit.xml`, whose depth
// depends on where the runner checked the repository out — so each machine wrote its own
// entry and none ever replaced another.
func TestSuiteKey_stableAcrossMachines(t *testing.T) {
	t.Run("NGSTI-B07: The suite key is the report path from the root, or the file name for a report outside the repository", func(t *testing.T) {})
	a := suiteKey("/home/runner/work/anchors/anchors", "/tmp/junit.xml")
	b := suiteKey("/Users/dev/code/anchors", "/tmp/junit.xml")
	if a != b {
		t.Errorf("the same external report got two keys: %q and %q", a, b)
	}
	if a != "external/junit.xml" {
		t.Errorf("external report key = %q, want external/junit.xml", a)
	}

	root := t.TempDir()
	in := suiteKey(root, filepath.Join(root, "apps", "mobile", "junit.xml"))
	if in != "apps/mobile/junit.xml" {
		t.Errorf("a report inside the repo should be keyed by its path from the root, got %q", in)
	}
}

// Two report paths resolving to the SAME node: the result must not depend on map order.
func TestResolveByFile_collisionIsDeterministic(t *testing.T) {
	t.Run("NGSTI-B10: When two report paths name the same node the path equal to the node id wins", func(t *testing.T) {})
	// The workspace-relative path SORTS FIRST here ("src/…" < "web/…"), so the rule that
	// the exact node ID wins is what decides — not alphabetical luck.
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "web/src/x.tsx", Kind: mapx.KindCode}}}
	byFile := map[string]int{
		"src/x.tsx":     1, // workspace-relative
		"web/src/x.tsx": 2, // already the node ID
	}
	for i := 0; i < 50; i++ {
		out := resolveByFile(g, mapx.KindCode, byFile, "/r", "/r/web/lcov.info")
		if out["web/src/x.tsx"] != 2 {
			t.Fatalf("run %d: the path that IS the node ID must win, got %d", i, out["web/src/x.tsx"])
		}
	}
}

// The external proof is dropped only by a FULL run of a suite inside the repository: a
// partial run measured only its cut, and an external report does not drop its own kind.
func TestDropExternal(t *testing.T) {
	t.Run("NGSTI-B09: A partial run or another external report drops no external suite", func(t *testing.T) {})
	graph := func() *mapx.Graph {
		g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a.spec.md", Kind: mapx.KindSpec, Rev: "r1"}}}
		g.IngestExecutionSuite(nil, map[string]bool{"AAAAA-B01": true}, nil,
			map[string][]string{"a.spec.md": {"AAAAA-B01"}}, "", "external/report.xml", "t")
		return g
	}
	for _, c := range []struct {
		name    string
		key     string
		partial bool
		drops   bool
	}{
		{"full in-repo run", ".anchors/junit.xml", false, true},
		{"partial in-repo run", ".anchors/junit.xml", true, false},
		{"another external report", "external/other.xml", false, false},
	} {
		g := graph()
		dropExternal(g, c.key, c.partial)
		_, kept := g.Nodes[0].Signal.ProvenBySuite["external/report.xml"]
		if kept == c.drops {
			t.Errorf("%s: external suite kept=%v, want dropped=%v", c.name, kept, c.drops)
		}
	}
}

// The occurrences are bound to the spec that declares them; a code no spec declares is
// reported, and left out of the map.
func TestIngestLogs_bindsOccurrencesToTheSpec(t *testing.T) {
	t.Run("NGSTI-B12: Log occurrences are bound to the spec that declares the failure, stamped with its revision", func(t *testing.T) {})
	t.Run("NGSTI-B13: Failure codes no spec declares are reported after the log ingestion", func(t *testing.T) {})
	t.Run("NGSTI-X01: A failure code no spec declares is bound to no spec", func(t *testing.T) {})
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
		for _, f := range g.Nodes[i].Failures {
			if f.Rule == "GHOST-E07" {
				t.Errorf("the undeclared code was bound to %s", g.Nodes[i].ID)
			}
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

// Each log ingestion replaces what the previous one bound: the same log ingested twice
// must not double the occurrences.
func TestIngestLogs_replacesInsteadOfAdding(t *testing.T) {
	t.Run("NGSTI-I01: Ingesting the same logs again replaces the earlier occurrences instead of adding to them", func(t *testing.T) {})
	root, _ := logsProject(t)
	runCmd(t, newIngestCmd(), "--logs", "--root", root)
	spec := node(t, root, "src/login.spec.md")
	if len(spec.Failures) != 3 {
		t.Fatalf("a second ingestion must leave 3 failures, got %+v", spec.Failures)
	}
	for _, f := range spec.Failures {
		if f.Rule == "LOGIN-E01" && f.Count != 3 {
			t.Errorf("LOGIN-E01 happened 3 times, a second ingestion made it %d", f.Count)
		}
	}
}

func TestIngestLogs_errors(t *testing.T) {
	t.Run("NGSTI-E03: Log ingestion without declared log paths is refused", func(t *testing.T) {})
	t.Run("NGSTI-E01: Ingest with no report flag refuses with the usage", func(t *testing.T) {})
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

func TestIngest_runTimeIsTheSumOfTheCases(t *testing.T) {
	t.Run("NGSTI-B14: A test file's run time is the sum of its cases", func(t *testing.T) {})
	root := fixtureProject(t)
	timed := strings.Replace(strings.Replace(junitReport,
		`name="LOGIN-B01: checks the password" file="src/login.test.ts"`, `name="LOGIN-B01: checks the password" file="src/login.test.ts" time="0.5"`, 1),
		`name="later" file="src/login.test.ts"`, `name="later" file="src/login.test.ts" time="0.75"`, 1)
	writeProjectFile(t, root, "reports/junit.xml", timed)
	runCmd(t, newIngestCmd(), "--root", root, "--junit", filepath.Join(root, "reports/junit.xml"))
	if s := node(t, root, "src/login.test.ts").Signal; s == nil || s.SecondsBySuite["reports/junit.xml"] != 1.25 {
		t.Fatalf("the test file must record 1.25s under its suite, got %+v", s)
	}
}

func TestSuiteKey(t *testing.T) {
	t.Run("NGSTI-B15: A report's signals are kept under its path from the root", func(t *testing.T) {})
	root := t.TempDir()
	if got := SuiteKey(root, filepath.Join(root, "out", "junit.xml")); got != "out/junit.xml" {
		t.Errorf("a report inside the repository is keyed by its path from the root, got %q", got)
	}
	if got := SuiteKey(root, filepath.Join(filepath.Dir(root), "elsewhere", "report.xml")); got != "external/report.xml" {
		t.Errorf("a report outside it is keyed under external/, got %q", got)
	}
}

// The case that lost proofs: a spec edited, `anchors test`, then `map build`. The run
// proved the new text, but the proof was stamped with the map's old rev of the spec, and
// the rebuilt map read it as stale — scenario-coverage said "no execution ingested".
func TestIngest_aRunStampsTheTreesRevs(t *testing.T) {
	t.Run("NGSTI-B16: A run's proofs are stamped with the tree's revs", func(t *testing.T) {})
	for _, viaRun := range []bool{true, false} {
		root := fixtureProject(t)
		writeProjectFile(t, root, "src/login.spec.md", fixtureSpec+"\nEdited after the map build.\n")
		writeProjectFile(t, root, "reports/junit.xml", junitReport)
		ViaAnchorsTest = viaRun
		err := IngestArtifacts(root, "", filepath.Join(root, "reports/junit.xml"), "", "", "unit", "", "", false)
		ViaAnchorsTest = false
		if err != nil {
			t.Fatal(err)
		}
		runCmd(t, newMapCmd(), "build", "--root", root)
		n := node(t, root, "src/login.spec.md")
		fresh := n.Signal != nil && n.Signal.AtRev == n.Rev
		if viaRun && !fresh {
			t.Errorf("after a run, the spec's proof must be fresh once the map is rebuilt, got %+v (rev %s)", n.Signal, n.Rev)
		}
		if !viaRun && fresh {
			t.Errorf("a manual ingestion must keep the map's revs, so the proof stays stale")
		}
	}
}

// The case that blocked a spec whose tests passed: created after the last `map build`, it
// had no node, and the proof of its first run was dropped.
func TestIngest_aNewSpecKeepsItsFirstProof(t *testing.T) {
	t.Run("NGSTI-B17: A spec created after the map build keeps its first proof", func(t *testing.T) {})
	root := fixtureProject(t)
	writeProjectFile(t, root, "src/pay.spec.md", "<!-- @anchors\n  code: PAYMX\n-->\n# Pay\n\nPAYMX-B01 — the amount is charged.\n")
	writeProjectFile(t, root, "src/pay.ts", "export const pay = () => true\n")
	writeProjectFile(t, root, "src/pay.test.ts", "test('PAYMX-B01: charges', () => {})\n")
	writeProjectFile(t, root, "reports/junit.xml", `<?xml version="1.0"?>
<testsuites><testsuite name="pay" file="src/pay.test.ts">
<testcase name="PAYMX-B01: charges" file="src/pay.test.ts"/>
</testsuite></testsuites>
`)
	ViaAnchorsTest = true
	err := IngestArtifacts(root, "", filepath.Join(root, "reports/junit.xml"), "", "", "unit", "", "", false)
	ViaAnchorsTest = false
	if err != nil {
		t.Fatal(err)
	}
	n := node(t, root, "src/pay.spec.md")
	if n.Signal == nil || !strings.Contains(strings.Join(n.Signal.ProvenCodes, ","), "PAYMX-B01") {
		t.Fatalf("the new spec is in the map with its rule proven, got %+v", n.Signal)
	}
	node(t, root, "src/pay.test.ts")
}

// A report on another drive has no path relative to the root: it is external, keyed by its
// name, the same on every machine.
func TestSuiteKey_anotherDrive(t *testing.T) {
	t.Run("NGSTI-B07: The suite key is the report path from the root, or the file name for a report outside the repository", func(t *testing.T) {})
	if runtime.GOOS != "windows" {
		t.Skip("drives exist only on Windows")
	}
	if got := suiteKey(`C:\work\repo`, `D:\reports\junit.xml`); got != "external/junit.xml" {
		t.Errorf("a report on another drive is external, got %q", got)
	}
}

func TestIngest_suiteCoverageMarksOmitted(t *testing.T) {
	t.Run("NGSTI-B18: A suite's whole coverage run marks the files it left out", func(t *testing.T) {})
	suite := "tests:\n  - layer: test\n    run: \"true\"\n    lcov: reports/lcov.info\n    paths: [\"src/**\"]\n"
	setup := func() string {
		root := fixtureProjectWith(t, suite, fixtureSpec)
		writeProjectFile(t, root, "src/types.ts", "export interface Payload { id: string }\n")
		runCmd(t, newMapCmd(), "build", "--root", root)
		writeProjectFile(t, root, "reports/lcov.info", lcovReport)
		return root
	}
	root := setup()
	out := runCmd(t, newIngestCmd(), "--root", root, "--lcov", filepath.Join(root, "reports/lcov.info"))
	if s := node(t, root, "src/types.ts").Signal; s == nil || s.CoverageOmitted == "" || !strings.Contains(out, "not in its report") {
		t.Errorf("the whole run marks the file of types as omitted and says so, got %+v\n%s", s, out)
	}
	root = setup()
	runCmd(t, newIngestCmd(), "--root", root, "--lcov", filepath.Join(root, "reports/lcov.info"), "--partial")
	if s := node(t, root, "src/types.ts").Signal; s != nil && s.CoverageOmitted != "" {
		t.Errorf("a partial run marks nothing, got %+v", s)
	}
}
