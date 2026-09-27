package flow

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// The two constants must be the SAME string.
//
// `scan` keeps the file out of the map; the command creates it. If they diverge, `new`
// generates a file with one suffix and the map excludes another — and the progress is
// confronted by the gates again, silently reintroducing the defect the separation removed.
// No behaviour test would catch it: both sides would work, each with its own ruler.
func TestProgress_suffixMatchesTheScanner(t *testing.T) {
	t.Run("PLPRP-I01: The suffix is the one the scanner keeps out of the map", func(t *testing.T) {})
	if !scan.IsProgressFile("plans/0001-x" + progressSuffix) {
		t.Fatalf("the command's suffix (%q) is not recognised by the scanner — `new` would "+
			"create a file the map does NOT exclude, and the gates would confront it again",
			progressSuffix)
	}
}

func TestProgress_pathSitsBesideThePlan(t *testing.T) {
	t.Run("PLPRP-B01: The progress file sits beside the plan", func(t *testing.T) {})
	got := progressPath("plans/0017-mutation.md")
	if want := "plans/0017-mutation-progress.md"; got != want {
		t.Fatalf("path: %q, want %q", got, want)
	}
}

// The progress is born with one section per DECLARED PHASE, read from the headers.
//
// The source is the same one the phase gates use. Keeping a second list would make the
// progress talk about phases that do not exist — and the divergence would only show when
// someone renamed a phase.
func TestProgress_oneSectionPerPhaseOfThePlan(t *testing.T) {
	t.Run("PLPRP-B02: One section per phase declared in the plan's headers", func(t *testing.T) {})
	plan := `# Plan 0017

## Phases

### MTUAO-F01 — the tool and the report

- ` + "`packages/shared/MutationHarness.spec.md`" + `

### MTUAO-F02 — CI ingests the signal (depends on MTUAO-F01)

## ABCDE-F03 — two

#### ABCDE-F04 — four

# ABCDE-F05 — one is the title level

##### ABCDE-F06 — five is below the phases
`
	dir := t.TempDir()
	p := filepath.Join(dir, "0017-mutation.md")
	if err := os.WriteFile(p, []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}

	dest, err := writeInitialProgress(p, plan, "MTUAO")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)

	for _, phase := range []string{"MTUAO-F01", "MTUAO-F02", "ABCDE-F03", "ABCDE-F04"} {
		if !strings.Contains(got, "## "+phase) {
			t.Errorf("the section of phase %s is missing:\n%s", phase, got)
		}
	}
	for _, phase := range []string{"ABCDE-F05", "ABCDE-F06"} {
		if strings.Contains(got, phase) {
			t.Errorf("a level-one or level-five header is not a phase, but %s got a section:\n%s", phase, got)
		}
	}
	// The phase TITLE comes along: without it the file is a list of codes, and whoever opens
	// it has to go back to the plan to know what each one is about.
	if !strings.Contains(got, "## MTUAO-F01 — the tool and the report") {
		t.Errorf("the phase title was not copied:\n%s", got)
	}
	// The checkbox lives HERE, and that is the whole point of the separation.
	if strings.Count(got, "- [ ]") != 4 {
		t.Errorf("each phase gets one unchecked item — where `[x]` is marked:\n%s", got)
	}
}

// IT DOES NOT OVERWRITE: the file holds state, and rewriting it would erase recorded work.
func TestProgress_doesNotOverwriteExistingState(t *testing.T) {
	t.Run("PLPRP-B05: An existing progress file is never overwritten", func(t *testing.T) {})
	dir := t.TempDir()
	p := filepath.Join(dir, "0001-x.md")
	if err := os.WriteFile(p, []byte("# Plan\n\n### ABCDE-F01 — phase\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	done := "# Progress — ABCDE\n\n## ABCDE-F01\n\n- [x] done\n"
	if err := os.WriteFile(progressPath(p), []byte(done), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := writeInitialProgress(p, "# Plan\n\n### ABCDE-F01 — phase\n", "ABCDE"); err == nil {
		t.Fatal("it overwrote the existing progress — the `[x]` of whoever worked would be erased")
	}

	b, _ := os.ReadFile(progressPath(p))
	if string(b) != done {
		t.Fatalf("the content changed:\n%s", b)
	}
}

// A plan WITHOUT declared phases still gets the file, with what to do written in it.
func TestProgress_planWithoutPhasesSaysWhatToDo(t *testing.T) {
	t.Run("PLPRP-B04: A plan with no phases gets a note saying what to add", func(t *testing.T) {})
	dir := t.TempDir()
	p := filepath.Join(dir, "0002-y.md")
	plan := "# Plan\n\n## Goal\n\nno phases yet\n"
	if err := os.WriteFile(p, []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}
	dest, err := writeInitialProgress(p, plan, "YYYYY")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(dest)
	if !strings.Contains(string(b), "the plan does not declare phases yet") {
		t.Errorf("a plan without phases should say what to do:\n%s", b)
	}
}

// THE CODE LENGTH comes from the CONFIG, not from a literal.
//
// `code_lengths` is configurable per project. While the regex fixed `{4,5}`, a project with
// three-letter codes had its phases recognised by the gates (which use
// `config.CodeLengthPattern()`) and IGNORED by this command — the progress was born empty,
// with nothing to accuse it.
func TestProgress_honoursTheProjectsCodeLengths(t *testing.T) {
	t.Run("PLPRP-B03: The phase code length follows the project's configuration", func(t *testing.T) {})
	original := config.CodeLengths
	t.Cleanup(func() { config.CodeLengths = original })
	config.CodeLengths = []int{3}

	plan := "# Plan\n\n### ABC-F01 — short code phase\n"
	dir := t.TempDir()
	p := filepath.Join(dir, "0001-short.md")
	if err := os.WriteFile(p, []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}

	dest, err := writeInitialProgress(p, plan, "ABC")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(dest)
	if !strings.Contains(string(b), "## ABC-F01") {
		t.Fatalf("the phase of a project with code_lengths=[3] was not recognised — the gates "+
			"see it and this command does not:\n%s", b)
	}
}

func runNewProgress(t *testing.T, args ...string) error {
	t.Helper()
	cmd := newProgressCmd()
	cmd.SetArgs(args)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	var err error
	stdoutOf(t, func() { err = cmd.Execute() })
	return err
}

// `anchors new progress --for <plan>` exists for the plans born BEFORE the mechanism.
// `anchors new plan` creates the companion along with the plan — but only it, and a project
// that adopted Anchors before this version is left with every plan without a companion.
//
// Measured: 17 plans, ZERO with `-progress.md`, and 17 with checkboxes INSIDE the plan.
func TestNewProgress_createsForAnExistingPlan(t *testing.T) {
	t.Run("PLPRP-B06: new progress creates the file for an existing plan", func(t *testing.T) {})
	t.Run("PLPRP-X01: The progress takes its code from the plan's header", func(t *testing.T) {})
	root := t.TempDir()
	plan := filepath.Join(root, "plans", "0002-platform.md")
	if err := os.MkdirAll(filepath.Dir(plan), 0o755); err != nil {
		t.Fatal(err)
	}
	const content = `<!-- @anchors
  code: PLTFR
  layer: plan
-->
# Platform

### PLTFR-F01 — the contract

### PLTFR-F02 — access to the sources
`
	if err := os.WriteFile(plan, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := runNewProgress(t, "--root", root, "--for", "plans/0002-platform.md"); err != nil {
		t.Fatalf("create the progress: %v", err)
	}

	prog := filepath.Join(root, "plans", "0002-platform-progress.md")
	b, err := os.ReadFile(prog)
	if err != nil {
		t.Fatalf("the file was not created: %v", err)
	}
	got := string(b)
	// The IDENTITY comes from the plan, not from an argument: a progress with a code that
	// diverges from its plan would no longer be locatable.
	if !strings.Contains(got, "# Progress — PLTFR") {
		t.Errorf("the progress did not inherit the plan's code:\n%s", got)
	}
	for _, phase := range []string{"## PLTFR-F01", "## PLTFR-F02"} {
		if !strings.Contains(got, phase) {
			t.Errorf("section %q is missing — the phases come from the plan's headers", phase)
		}
	}

	// IT DOES NOT OVERWRITE: that is what makes running it over a whole project safe.
	if err := runNewProgress(t, "--root", root, "--for", "plans/0002-platform.md"); err == nil {
		t.Error("running it again overwrote the state — it should refuse")
	}
}

// Without `code:` in the header the progress would be born without identity.
func TestNewProgress_refusesAPlanWithoutCode(t *testing.T) {
	t.Run("PLPRP-E01: A plan without a code is refused", func(t *testing.T) {})
	root := t.TempDir()
	plan := filepath.Join(root, "p.md")
	if err := os.WriteFile(plan, []byte("# Plan without header\n\n### ABCDE-F01 — phase\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runNewProgress(t, "--root", root, "--for", "p.md"); err == nil {
		t.Error("it accepted a plan without `code:` — the progress would be born without identity")
	}
	if _, err := os.Stat(progressPath(plan)); !os.IsNotExist(err) {
		t.Errorf("no progress file may be created for a plan without code (stat: %v)", err)
	}
}

func TestNewProgress_requiresAReadablePlan(t *testing.T) {
	t.Run("PLPRP-B07: new progress without the plan is refused", func(t *testing.T) {})
	t.Run("PLPRP-E02: A plan that cannot be read is refused", func(t *testing.T) {})
	root := t.TempDir()
	if err := runNewProgress(t, "--root", root); err == nil || !strings.Contains(err.Error(), "--for") {
		t.Errorf("without --for the command must refuse naming the flag, got %v", err)
	}
	if err := runNewProgress(t, "--root", root, "--for", "plans/none.md"); err == nil || !strings.Contains(err.Error(), "read the plan") {
		t.Errorf("a missing plan must fail reading it, got %v", err)
	}
}
