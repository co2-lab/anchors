package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain lets a test run the real `main` in a child process: with ANCHORS_MAIN_ARGS set,
// this binary IS `anchors`, and its exit code and stderr are what a hook sees.
func TestMain(m *testing.M) {
	if args, ok := os.LookupEnv("ANCHORS_MAIN_ARGS"); ok {
		os.Args = append([]string{"anchors"}, strings.Split(args, "\x1f")...)
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// runMain runs `anchors <args>` through main() in dir and returns stdout, stderr and the
// exit code.
func runMain(t *testing.T, dir string, args ...string) (string, string, int) {
	t.Helper()
	c := exec.Command(os.Args[0], "-test.run=^$")
	c.Dir = dir
	c.Env = append(os.Environ(), "ANCHORS_MAIN_ARGS="+strings.Join(args, "\x1f"), "ANCHORS_TELEMETRY=off")
	var out, errOut bytes.Buffer
	c.Stdout, c.Stderr = &out, &errOut
	err := c.Run()
	code := 0
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return out.String(), errOut.String(), code
}

func writeTo(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A failure is printed ONCE, prefixed with the localized word for error, and exits 1.
func TestMainPrintsTheErrorOnceAndExitsOne(t *testing.T) {
	t.Run("CLMNC-B01: A failing command prints its error once and exits 1", func(t *testing.T) {})
	_, stderr, code := runMain(t, t.TempDir(), "no-such-command")
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if n := strings.Count(stderr, `unknown command "no-such-command"`); n != 1 {
		t.Errorf("the error appears %d times, want once:\n%s", n, stderr)
	}
	if !strings.HasPrefix(stderr, "error: ") {
		t.Errorf("the error is not prefixed with the localized word:\n%s", stderr)
	}
}

// "Not governed" is not a failure: it exits with its own code, which the hooks read.
func TestMainExitsThreeForANotGovernedFile(t *testing.T) {
	t.Run("CLMNC-B02: A file the project does not govern exits with the not-governed code", func(t *testing.T) {})
	dir := t.TempDir()
	writeTo(t, dir, "anchors.yaml", "version: 1\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n"+
		"gates:\n  - name: spec-has-code\n    id: spec-has-code\n    \"on\": [spec]\n    check: has-code\n    blocking: true\n")
	writeTo(t, dir, "anchors.graph.yaml", "version: 5\nnodes: []\nedges: []\n")
	writeTo(t, dir, "package.json", "{}\n")
	_, stderr, code := runMain(t, dir, "check", "--changed", "package.json", "--no-record")
	if code != 3 {
		t.Errorf("exit code = %d, want 3 (not governed)\n%s", code, stderr)
	}
	if !strings.Contains(stderr, "is not governed") {
		t.Errorf("stderr does not say why:\n%s", stderr)
	}
}

// The build identity reaches every place that records it: the version the CLI reports and
// the generator the map records.
func TestMainInjectsTheBuildIdentity(t *testing.T) {
	t.Run("CLMNC-B03: The build version reaches the reported version and the map's generator", func(t *testing.T) {})
	dir := t.TempDir()
	stdout, _, code := runMain(t, dir, "--version")
	if code != 0 || !strings.Contains(stdout, "dev (commit none, built unknown)") {
		t.Errorf("--version (exit %d):\n%s", code, stdout)
	}
	writeTo(t, dir, "anchors.yaml", "version: 1\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	if _, stderr, code := runMain(t, dir, "map", "build"); code != 0 {
		t.Fatalf("map build (exit %d):\n%s", code, stderr)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "anchors.graph.yaml")); !strings.Contains(string(b), "generated_by: dev") {
		t.Errorf("the map does not record the build version:\n%s", b)
	}
}

// A key the migration renames, in a config of an older format, is told apart from a typo:
// the error says to migrate, not to update the binary.
func TestMainTellsARenamedKeyFromATypo(t *testing.T) {
	t.Run("CLMNC-B04: A renamed key in an older config points at migrate", func(t *testing.T) {})
	dir := t.TempDir()
	writeTo(t, dir, "anchors.yaml", "version: 1\ntrinca_opcional: true\n")
	_, stderr, code := runMain(t, dir, "generated-paths")
	if code != 1 || !strings.Contains(stderr, "anchors migrate") || !strings.Contains(stderr, "`trinca_opcional` was renamed") {
		t.Errorf("exit %d, want the migrate advice:\n%s", code, stderr)
	}
}
