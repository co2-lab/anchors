package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
)

// fakeGH puts a `gh` on the PATH that answers without touching GitHub and logs every
// call (argv, and the stdin of `--input -`) to the returned file. `perm` is what the
// collaborator permission query answers; `issue list` answers one card.
func fakeGH(t *testing.T, perm string) string {
	t.Helper()
	bin := t.TempDir()
	log := filepath.Join(bin, "calls.log")
	script := `#!/bin/sh
echo "gh $*" >> "` + log + `"
case "$*" in
  "auth status") exit 0 ;;
  "api user --jq .login") echo me; exit 0 ;;
  *collaborators/me/permission*) echo ` + perm + `; exit 0 ;;
  *branches/staging/protection*) echo "Branch not found"; exit 1 ;;
  "api --method PUT "*) echo "stdin: $(cat)" >> "` + log + `"; exit 0 ;;
  *"--jq .enforce_admins.enabled"*) echo false; exit 0 ;;
  "label create "*) exit 0 ;;
  "issue list "*) printf '12\tFix login\tanchors:in-progress\n'; exit 0 ;;
esac
echo "unexpected gh call: $*" >&2
exit 2
`
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return log
}

const doctorGitHubYAML = `version: 2
layers: {}
workflow:
  mode: github
  repo: acme/app
  labels: [anchors]
  protected_branches: [main, staging]
`

func TestDoctorPrintsTheDiagnosisGroupedByCheck(t *testing.T) {
	englishOutput(t)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "LONE.md", Kind: mapx.KindGuide},
		{ID: "a.spec.md", Kind: mapx.KindSpec},
	}}
	dir := qProject(t, "version: 2\nlayers: {}\n", map[string]string{"LONE.md": "# l\n", "a.spec.md": "# a\n"}, g)

	out, err := runQ(t, newDoctorCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"anchors doctor — 2 nodes, 0 edges, 0 layers",
		"⚠ guide-sem-governo (1)\n  ⚠ LONE.md — ",
		"⚠ identidade-ausente (1)\n  ⚠ a.spec.md — ",
		"(diagnosis — nothing was blocked; you decide what to reconcile)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

// A group is marked warn when ANY of its findings is warn, even if an info came first.
func TestPrintReportGroupMarkTakesTheHighestSeverity(t *testing.T) {
	englishOutput(t)
	rep := health.Report{Nodes: 3, Edges: 1, Layers: 2, Findings: []health.Finding{
		{Check: "mixed", Severity: health.Info, Subject: "a", Detail: "note"},
		{Check: "mixed", Severity: health.Warn, Subject: "b", Detail: "problem"},
		{Check: "calm", Severity: health.Info, Subject: "c", Detail: "fine"},
	}}
	out := captureStdout(t, func() { printReport(rep) })
	for _, want := range []string{
		"⚠ mixed (2)\n    a — note\n  ⚠ b — problem\n",
		"ℹ calm (1)\n",
		"summary: 1 attention point(s), 3 finding(s) in total",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}

	out = captureStdout(t, func() { printReport(health.Report{}) })
	if !strings.Contains(out, "✓ no systemic loose end found") {
		t.Errorf("no finding: the report says it is clean:\n%s", out)
	}
}

func TestDoctorFixOutsideGitHubModeDoesNothing(t *testing.T) {
	englishOutput(t)
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, &mapx.Graph{})

	out, err := runQ(t, newDoctorCmd(), "--root", dir, "--fix")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "--fix: nothing to do") {
		t.Errorf("local mode has no GitHub environment to repair:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, initx.DirWorkflows)); err == nil {
		t.Error("local mode must not seed pipelines")
	}
}

func TestDoctorNeedsConfigAndMap(t *testing.T) {
	dir := t.TempDir()
	if _, err := runQ(t, newDoctorCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "load config") {
		t.Errorf("no anchors.yaml; got %v", err)
	}
	dir = qProject(t, "version: 2\nlayers: {}\n", nil, nil)
	if _, err := runQ(t, newDoctorCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "load map") {
		t.Errorf("no map; got %v", err)
	}
}

// --fix in github mode, against a fake `gh`: it seeds the pipelines, protects the
// declared branches with the body on stdin (a missing branch is reported, not fatal),
// disables the unreachable approval, creates the state labels, and warns about the
// local queue and delivery files that must not exist in this mode.
func TestRepairEnvironmentInGitHubMode(t *testing.T) {
	englishOutput(t)
	log := fakeGH(t, "write")
	dir := qProject(t, doctorGitHubYAML, map[string]string{
		".anchors/tasks/pending__001-code-x.yaml": "id: 001\n",
		"changes/delivery.md":                     "# delivered\n",
	}, &mapx.Graph{})

	out, err := runQ(t, newDoctorCmd(), "--root", dir, "--fix")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"⚠ 1 file(s) in .anchors/tasks/ — and mode is `github`",
		"⚠ 1 delivery record(s) in `changes/` — and mode is `github`",
		"pipeline(s) created in " + initx.DirWorkflows,
		"✓ main protected — nothing enters without a PR",
		"· staging does not exist yet — protect it when created",
		"✓ approval requirement DISABLED",
		"state label(s) ensured in acme/app",
		"the board is OPTIONAL",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	for _, w := range initx.WorkflowsDoFluxo {
		if _, err := os.Stat(filepath.Join(dir, initx.DirWorkflows, w.Arquivo)); err != nil {
			t.Errorf("pipeline %s not seeded: %v", w.Arquivo, err)
		}
	}
	calls := readQ(t, log)
	for _, want := range []string{
		"gh auth status",
		"gh api --method PUT repos/acme/app/branches/main/protection --input -",
		// the body must reach gh — an empty stdin is a 422 and an open door
		`stdin: {"required_status_checks":null,"enforce_admins":false,"required_pull_request_reviews":{"required_approving_review_count":1},"restrictions":null}`,
		"gh label create anchors --repo acme/app --color 5319e7 --force",
		"gh label create " + initx.LabelNeedsUser + " --repo acme/app --color d73a4a --force",
	} {
		if !strings.Contains(calls, want) {
			t.Errorf("gh was not called with %q:\n%s", want, calls)
		}
	}

	// IDEMPOTENT: a second run creates no pipeline.
	out, err = runQ(t, newDoctorCmd(), "--root", dir, "--fix")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "--fix: the pipelines already exist and are up to date.") {
		t.Errorf("a set-up repository changes nothing:\n%s", out)
	}
}

// Without `gh`, --fix refuses with the instruction, before trying five repairs that
// would each fail on their own.
func TestRepairEnvironmentRequiresGH(t *testing.T) {
	englishOutput(t)
	t.Setenv("PATH", t.TempDir())
	cfg := &config.Config{Workflow: &config.Workflow{Mode: config.ModeGitHub, Repo: "acme/app", Labels: []string{"anchors"}}}
	dir := t.TempDir()

	var err error
	captureStdout(t, func() { err = repairEnvironment(dir, cfg) })
	if err == nil || !strings.Contains(err.Error(), "gh is not installed") {
		t.Errorf("no gh: must refuse naming it; got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, initx.DirWorkflows)); statErr == nil {
		t.Error("nothing is seeded before the credential check passes")
	}
	if err := protectBranches(cfg); err == nil {
		t.Error("protectBranches without gh must fail")
	}
	if err := createStateLabels(cfg); err == nil {
		t.Error("createStateLabels without gh must fail")
	}
}

// --check-pipelines answers one question from the working directory's project.
func TestCheckPipelines(t *testing.T) {
	englishOutput(t)
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, nil)
	t.Chdir(dir)
	out, err := runQ(t, newDoctorCmd(), "--check-pipelines")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "· local mode — no workflow pipeline to check.") {
		t.Errorf("local mode has no pipeline:\n%s", out)
	}

	gh := qProject(t, doctorGitHubYAML, nil, nil)
	t.Chdir(gh)
	out, err = runQ(t, newDoctorCmd(), "--check-pipelines")
	if err != nil {
		t.Fatalf("by default a missing pipeline warns and CI continues; got %v", err)
	}
	if !strings.Contains(out, "MISSING") || !strings.Contains(out, "(warning — CI continues.") {
		t.Errorf("missing pipelines are named, with the opt-in to block:\n%s", out)
	}

	blocks := qProject(t, doctorGitHubYAML+"  stale_pipeline_blocks: true\n", nil, nil)
	t.Chdir(blocks)
	if _, err := runQ(t, newDoctorCmd(), "--check-pipelines"); err == nil || !strings.Contains(err.Error(), "stale_pipeline_blocks: true") {
		t.Errorf("with stale_pipeline_blocks the check fails; got %v", err)
	}
	if _, _, err := initx.SemeiaWorkflows(blocks, &config.Config{Workflow: &config.Workflow{Mode: config.ModeGitHub}}); err != nil {
		t.Fatal(err)
	}
	out, err = runQ(t, newDoctorCmd(), "--check-pipelines")
	if err != nil {
		t.Fatalf("seeded pipelines are current; got %v", err)
	}
	if !strings.Contains(out, "workflow pipelines are in place and up to date") {
		t.Errorf("all current:\n%s", out)
	}
}
