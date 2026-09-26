package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
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
