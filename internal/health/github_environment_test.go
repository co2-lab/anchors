package health

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/initx"
)

func cfgGitHub() *config.Config {
	return &config.Config{Workflow: &config.Workflow{
		Mode:   config.ModeGitHub,
		Repo:   "acme/exemplo",
		Labels: []string{"anchors"},
	}}
}

// Charging board and pipelines to a project that declared `mode: local` would be
// guaranteed noise — and recurring noise trains the team to ignore the doctor.
func TestEnvironmentChecksNothingInLocalMode(t *testing.T) {
	t.Run("GHEGT-B01: Local mode checks nothing", func(t *testing.T) {})
	dir := t.TempDir()

	if fs := checkGitHubEnv(&config.Config{}, dir); len(fs) != 0 {
		t.Errorf("local mode uses neither board nor pipelines: %+v", fs)
	}
	if fs := checkGitHubEnv(nil, dir); len(fs) != 0 {
		t.Errorf("no configuration checks nothing: %+v", fs)
	}
}

// A missing pipeline gives no error anywhere — it gives SILENCE, and the artefacts stay
// without a card forever. That is why the doctor must warn: it is the only thing that shows.
func TestEnvironmentWarnsMissingPipeline(t *testing.T) {
	t.Run("GHEGT-B02: Each missing pipeline gives a warning that names what stops happening", func(t *testing.T) {})
	dir := t.TempDir()

	fs := checkPipelines(dir, cfgGitHub())

	if len(fs) != len(initx.WorkflowsDoFluxo) {
		t.Fatalf("expected %d findings, got %d", len(initx.WorkflowsDoFluxo), len(fs))
	}
	for _, f := range fs {
		if f.Check != "pipeline-ausente" || f.Severity != Warn {
			t.Errorf("wrong finding: %+v", f)
		}
		// The message must say what DOES NOT HAPPEN — "a file is missing" does not
		// explain why it matters.
		if !strings.Contains(f.Detail, "não acontece") && !strings.Contains(f.Detail, "does not happen") {
			t.Errorf("the message should name what stops happening: %s", f.Detail)
		}
		if !strings.Contains(f.Detail, "--fix") {
			t.Errorf("the message should point at the fix: %s", f.Detail)
		}
	}
}

// The worst case: the pipeline IS there, so it looks configured — but without
// `concurrency` it runs in parallel with itself and brings back the race it existed to
// remove. Its own finding, never counted as "missing".
func TestEnvironmentCatchesPipelineWithoutSerialization(t *testing.T) {
	t.Run("GHEGT-B03: A pipeline without serialization is its own finding, not a missing one", func(t *testing.T) {})
	dir := t.TempDir()
	if _, _, err := initx.SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}
	// Break one: remove the serialization, keep the file.
	target := filepath.Join(dir, initx.DirWorkflows, "anchors-claim.yml")
	if err := os.WriteFile(target, []byte("name: claim\non:\n  workflow_dispatch:\njobs:\n  x:\n    runs-on: ubuntu-latest\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fs := checkPipelines(dir, cfgGitHub())

	if len(fs) != 1 {
		t.Fatalf("expected 1 finding (the broken pipeline), got %d: %+v", len(fs), fs)
	}
	if fs[0].Check != "pipeline-sem-serializacao" {
		t.Errorf("the file exists — it cannot be reported as missing: %+v", fs[0])
	}
	if !strings.Contains(fs[0].Detail, "mesmo card a dois agentes") && !strings.Contains(fs[0].Detail, "same card to two agents") {
		t.Errorf("the message should say what breaks in practice: %s", fs[0].Detail)
	}
}

// After `--fix` no pipeline finding may remain: what the doctor charges and what `--fix`
// creates come from the SAME list, and a disagreement would make the doctor ask forever
// for something it has just created.
func TestFixSatisfiesWhatTheDoctorCharges(t *testing.T) {
	t.Run("GHEGT-I01: After the fix seeds the pipelines no pipeline finding remains", func(t *testing.T) {})
	dir := t.TempDir()
	if fs := checkPipelines(dir, cfgGitHub()); len(fs) == 0 {
		t.Fatal("setup: the empty project should have findings")
	}

	if _, _, err := initx.SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}

	if fs := checkPipelines(dir, cfgGitHub()); len(fs) != 0 {
		t.Errorf("--fix created the pipelines and the doctor still complains: %+v", fs)
	}
}

// The board is NOT checked, deliberately: the state of work is a label, and the Project
// is an optional mirror (BOOTSTRAP.md §7.13). The previous decision — state in the column
// — demanded a PAT with `project` scope from every adopter, an adoption cost for a
// visualization choice.
//
// This test is what prevents the way back: charging the board would make the doctor ask
// for a piece the flow does not use.
func TestDoctorDoesNotChargeTheBoard(t *testing.T) {
	t.Run("GHEGT-X01: The board is never charged", func(t *testing.T) {})
	dir := t.TempDir()
	if _, _, err := initx.SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}

	fs := checkGitHubEnv(cfgGitHub(), dir)

	for _, f := range fs {
		if strings.Contains(strings.ToLower(f.Check), "board") ||
			strings.Contains(strings.ToLower(f.Detail), "project") {
			t.Errorf("the doctor charges the board, which became optional: %+v", f)
		}
	}
}

// The OUTDATED pipeline is how a fix in the flow's design reaches whoever already
// installed it. Without this finding, a defect fixed at the source keeps running in the
// project forever — and nothing warns, because the file IS there.
func TestEnvironmentWarnsOutdatedPipeline(t *testing.T) {
	t.Run("GHEGT-B04: A pipeline behind its template is reported only while it carries the marker", func(t *testing.T) {})
	dir := t.TempDir()
	if _, _, err := initx.SemeiaWorkflows(dir, cfgGitHub()); err != nil {
		t.Fatal(err)
	}
	// Age one pipeline: content different from the template, marker INTACT.
	target := filepath.Join(dir, initx.DirWorkflows, initx.WorkflowsDoFluxo[0].Arquivo)
	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, append(b, []byte("\n# drift\n")...), 0o644); err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, f := range checkPipelines(dir, cfgGitHub()) {
		if f.Check == "pipeline-desatualizado" {
			found = true
		}
	}
	if !found {
		t.Error("the pipeline fell behind and the doctor did not warn")
	}

	// Now WITHOUT the marker: the team adopted the file, and the difference is its own
	// customization — warning here would ask it to undo its own work.
	noMarker := strings.ReplaceAll(string(b), initx.MarcadorDeTemplate, "# edited by the team")
	if err := os.WriteFile(target, []byte(noMarker+"\n# drift\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, f := range checkPipelines(dir, cfgGitHub()) {
		if f.Check == "pipeline-desatualizado" {
			t.Error("a pipeline without the marker is the team's — it cannot be charged as outdated")
		}
	}
}

// The branch protection is the rule the whole flow assumes, and the only one whose absence
// produces no error anywhere: a direct push works and skips the card. The fake `gh`
// answers offline.
func TestEnvironmentChecksBranchProtection(t *testing.T) {
	q := "'api repos/acme/exemplo/branches/main/protection --jq .required_pull_request_reviews != null'"

	t.Run("GHEGT-B05: An unreadable branch protection is an unprotected branch", func(t *testing.T) {
		fakeGH(t) // the read fails, as a 404 for an unprotected branch does
		fs := checkBranchProtection(cfgGitHub())
		want := Finding{"main-sem-protecao", Warn, "acme/exemplo", i18n.T("health.main_unprotected", "main")}
		if len(fs) != 1 || fs[0] != want {
			t.Fatalf("checkBranchProtection = %+v, want [%+v]", fs, want)
		}
	})

	t.Run("GHEGT-B06: A protection without required reviews is partially protected", func(t *testing.T) {
		fakeGH(t, ghAnswer{match: q, out: "false"})
		fs := checkBranchProtection(cfgGitHub())
		want := Finding{"main-sem-protecao", Warn, "acme/exemplo", i18n.T("health.main_partially_unprotected", "main")}
		if len(fs) != 1 || fs[0] != want {
			t.Fatalf("checkBranchProtection = %+v, want [%+v]", fs, want)
		}

		fakeGH(t, ghAnswer{match: q, out: "true"})
		if fs := checkBranchProtection(cfgGitHub()); len(fs) != 0 {
			t.Fatalf("a protection that requires reviews gives no finding: %+v", fs)
		}
	})

	// The protection was read on a hard-coded `main`: a project whose work lands on
	// `develop` was told its protected branch was unprotected, and an unprotected
	// `develop` passed whenever `main` happened to be protected.
	t.Run("GHEGT-B08: The protection is read on the declared integration branch", func(t *testing.T) {
		cfg := cfgGitHub()
		cfg.Workflow.IntegrationBranch = "develop"
		dev := "'api repos/acme/exemplo/branches/develop/protection --jq .required_pull_request_reviews != null'"
		fakeGH(t, ghAnswer{match: dev, out: "true"}, ghAnswer{match: q, out: "false"})
		if fs := checkBranchProtection(cfg); len(fs) != 0 {
			t.Fatalf("a protected develop gives no finding: %+v", fs)
		}
		fakeGH(t, ghAnswer{match: q, out: "true"})
		fs := checkBranchProtection(cfg)
		if len(fs) != 1 || !strings.Contains(fs[0].Detail, "`develop`") {
			t.Fatalf("an unprotected develop is the finding, named in the message: %+v", fs)
		}
	})

	t.Run("GHEGT-B07: Without the platform CLI the branch protection is not asked", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		if fs := checkBranchProtection(cfgGitHub()); len(fs) != 0 {
			t.Fatalf("without gh there must be no finding: %+v", fs)
		}
	})
}
