package health

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

func hasFinding(r Report, check, subject string) bool {
	for _, f := range r.Findings {
		if f.Check == check && f.Subject == subject {
			return true
		}
	}
	return false
}

// A toy project on disk + a graph consistent with it, except for the problems we want
// the doctor to catch.
func setup(t *testing.T) (string, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	// files that EXIST
	write := func(rel string) {
		full := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		_ = os.WriteFile(full, []byte("x"), 0o644)
	}
	write("Login.spec.md")
	write("Login.tsx")
	write("Orphan.tsx") // code without a spec
	write("guides/G.md")

	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "Login.spec.md", Kind: mapx.KindSpec, Code: "LOGIX"},
			{ID: "Login.tsx", Kind: mapx.KindCode},
			{ID: "Orphan.tsx", Kind: mapx.KindCode},                   // code without a spec: acceptable, NOT reported
			{ID: "guides/G.md", Kind: mapx.KindGuide},                 // governs nothing → guide-sem-governo
			{ID: "Sumido.spec.md", Kind: mapx.KindSpec, Code: "SUMI"}, // in the map, no file → no-fantasma
		},
		Edges: []mapx.Edge{
			{From: "Login.spec.md", To: "Login.tsx", Type: mapx.EdgeSpecifies},
			{From: "Login.spec.md", To: "Sumido.spec.md", Type: mapx.EdgeGoverns},
		},
	}
	cfg := &config.Config{
		Layers: map[string]config.Layer{
			"spec":    {Kind: "spec"},
			"code":    {Kind: "code"},
			"guide":   {Kind: "guide"},
			"feature": {Kind: "feature"}, // declared but WITHOUT files → camada-vazia
		},
		Gates: []config.Gate{
			{Name: "spec-complete", On: []string{"spec"}, Check: "non-empty"},
		},
	}
	return root, g, cfg
}

func TestDiagnose_mapFidelity(t *testing.T) {
	t.Run("DCTRO-B03: A node whose file is gone is a ghost", func(t *testing.T) {})
	root, g, cfg := setup(t)
	r := Diagnose(g, cfg, root)
	if !hasFinding(r, "no-fantasma", "Sumido.spec.md") {
		t.Error("it should detect a ghost node (in the map, no file)")
	}
}

func TestDiagnose_deadEdge(t *testing.T) {
	t.Run("DCTRO-B04: An edge to an unknown node is dead", func(t *testing.T) {})
	root, g, cfg := setup(t)
	g.Edges = append(g.Edges, mapx.Edge{From: "Login.spec.md", To: "Nowhere.tsx", Type: mapx.EdgeSpecifies})
	r := Diagnose(g, cfg, root)
	if !hasFinding(r, "aresta-morta", "Login.spec.md → Nowhere.tsx") {
		t.Errorf("an edge to a node the map does not know is dead: %+v", r.Findings)
	}
	// Edges between known nodes are not dead, even when one node's file is gone.
	if hasFinding(r, "aresta-morta", "Login.spec.md → Sumido.spec.md") {
		t.Error("an edge between known nodes is not dead")
	}
}

func TestDiagnose_orphans(t *testing.T) {
	t.Run("DCTRO-X01: Code without a spec is not reported", func(t *testing.T) {})
	root, g, cfg := setup(t)
	r := Diagnose(g, cfg, root)
	// Code WITHOUT a spec is NOT an orphan — it is the normal case (utils, constants, hooks).
	for _, f := range r.Findings {
		if f.Subject == "Orphan.tsx" {
			t.Errorf("code without a spec is acceptable; it should not be reported: %+v", f)
		}
	}
	// Login.spec has code LOGIX → NOT identidade-ausente
	if hasFinding(r, "identidade-ausente", "Login.spec.md") {
		t.Error("Login.spec has a code; it should not be identidade-ausente")
	}
}

func TestDiagnose_specWithoutRealization(t *testing.T) {
	t.Run("DCTRO-B05: A spec that points at nothing has no realization", func(t *testing.T) {})
	root, g, cfg := setup(t)
	r := Diagnose(g, cfg, root)
	if !hasFinding(r, "spec-sem-realizacao", "Sumido.spec.md") {
		t.Errorf("a spec with no outgoing edge has no realization: %+v", r.Findings)
	}
	if hasFinding(r, "spec-sem-realizacao", "Login.spec.md") {
		t.Error("Login.spec specifies Login.tsx; it is realized")
	}
}

func TestDiagnose_looseLayers(t *testing.T) {
	t.Run("DCTRO-B12: A declared layer with no node of its kind is empty", func(t *testing.T) {})
	t.Run("DCTRO-B13: A guide that governs nothing is reported", func(t *testing.T) {})
	root, g, cfg := setup(t)
	r := Diagnose(g, cfg, root)
	if !hasFinding(r, "camada-vazia", "feature") {
		t.Error("it should detect the 'feature' layer declared but empty")
	}
	if !hasFinding(r, "guide-sem-governo", "guides/G.md") {
		t.Error("it should detect a guide that governs nothing")
	}
}

func TestDiagnose_missingIdentity(t *testing.T) {
	t.Run("DCTRO-B06: A spec without a code has no identity", func(t *testing.T) {})
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "NoCode.spec.md"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "NoCode.tsx"), []byte("x"), 0o644)
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "NoCode.spec.md", Kind: mapx.KindSpec, Code: ""}, // no code
			{ID: "NoCode.tsx", Kind: mapx.KindCode},
		},
		Edges: []mapx.Edge{{From: "NoCode.spec.md", To: "NoCode.tsx", Type: mapx.EdgeSpecifies}},
	}
	r := Diagnose(g, &config.Config{Layers: map[string]config.Layer{}}, root)
	if !hasFinding(r, "identidade-ausente", "NoCode.spec.md") {
		t.Error("a spec without a code should be identidade-ausente (the invisible orphan)")
	}
}

func findingSeverity(r Report, check, subject string) (Severity, bool) {
	for _, f := range r.Findings {
		if f.Check == check && f.Subject == subject {
			return f.Severity, true
		}
	}
	return "", false
}

func TestCheckDuplicateCodes(t *testing.T) {
	t.Run("DCTRO-B07: A code owned by units of different domains is a duplicate identity", func(t *testing.T) {})
	t.Run("DCTRO-B08: Units of the same domain may share a code", func(t *testing.T) {})
	t.Run("DCTRO-B09: Only specs, features, tests and code own a code", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		// a single unit with code UNIQ — no collision
		{ID: "features/a/screens/A.spec.md", Kind: mapx.KindSpec, Code: "UNIQ"},
		{ID: "features/a/screens/A.tsx", Kind: mapx.KindCode, Code: "UNIQ"},
		// intra-feature: lib + screens of the SAME feature 'b' with INTR
		{ID: "features/b/lib/x.ts", Kind: mapx.KindCode, Code: "INTR"},
		{ID: "features/b/screens/X.tsx", Kind: mapx.KindCode, Code: "INTR"},
		// cross-domain: distinct features with CROSS → Warn
		{ID: "features/auth/screens/S.tsx", Kind: mapx.KindCode, Code: "CROSS"},
		{ID: "features/dash/screens/H.tsx", Kind: mapx.KindCode, Code: "CROSS"},
		// a doc is NOT an owner: it only references CROSS
		{ID: "guides/report.md", Kind: mapx.KindDoc, Code: "CROSS"},
		// a doc in another domain referencing UNIQ does not make it a duplicate
		{ID: "docs/unique.md", Kind: mapx.KindDoc, Code: "UNIQ"},
	}}
	r := Diagnose(g, &config.Config{Layers: map[string]config.Layer{}}, t.TempDir())

	if hasFinding(r, "identidade-duplicada", "UNIQ") {
		t.Error("UNIQ has a single owner — it should not be a duplicate")
	}
	// intra-feature (lib+screen of the same feature) is legitimate sharing — NOT reported
	// (a finding that never demands action is noise).
	if hasFinding(r, "identidade-duplicada", "INTR") {
		t.Error("INTR (intra-feature) is legitimate; it should not give a finding")
	}
	// only the cross-domain collision is reported, as Warn, naming both owners.
	if sev, ok := findingSeverity(r, "identidade-duplicada", "CROSS"); !ok || sev != Warn {
		t.Errorf("CROSS (cross-domain) should be Warn, got sev=%v ok=%v", sev, ok)
	}
	for _, f := range r.Findings {
		if f.Subject == "CROSS" && (!strings.Contains(f.Detail, "features/auth/screens/S") ||
			!strings.Contains(f.Detail, "features/dash/screens/H") || strings.Contains(f.Detail, "guides/report")) {
			t.Errorf("the finding must name exactly the owning units: %s", f.Detail)
		}
	}
}

func TestCheckDuplicateCodes_unitAndOptOut(t *testing.T) {
	t.Run("DCTRO-B10: The files of one unit count as one owner", func(t *testing.T) {
		g := &mapx.Graph{Nodes: []mapx.Node{
			// At the root each path is its own top-level domain, so only the unit's
			// identity keeps the four files from colliding.
			{ID: "Login.spec.md", Kind: mapx.KindSpec, Code: "LOGIN"},
			{ID: "Login.feature", Kind: mapx.KindFeature, Code: "LOGIN"},
			{ID: "Login.test.tsx", Kind: mapx.KindTest, Code: "LOGIN"},
			{ID: "Login.tsx", Kind: mapx.KindCode, Code: "LOGIN"},
		}}
		if fs := checkDuplicateCodes(g); len(fs) != 0 {
			t.Fatalf("one unit's triad is one owner: %+v", fs)
		}
	})
	t.Run("DCTRO-B11: A file that declares a shared code is not an owner", func(t *testing.T) {
		g := &mapx.Graph{Nodes: []mapx.Node{
			{ID: "apps/Login.tsx", Kind: mapx.KindCode, Code: "LOGIN"},
			{ID: "packages/e2e/flow.test.ts", Kind: mapx.KindTest, Code: "LOGIN", SharedCode: true},
		}}
		if fs := checkDuplicateCodes(g); len(fs) != 0 {
			t.Fatalf("a shared-code file does not collide: %+v", fs)
		}
		g.Nodes[1].SharedCode = false
		if fs := checkDuplicateCodes(g); len(fs) != 1 {
			t.Fatalf("without the opt-out the two domains collide: %+v", fs)
		}
	})
}

func TestCheckGateCoverage(t *testing.T) {
	t.Run("DCTRO-B14: A spec, feature or test kind that no gate confronts is reported", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.spec.md", Kind: mapx.KindSpec}, {ID: "a.feature", Kind: mapx.KindFeature},
		{ID: "a_test.go", Kind: mapx.KindTest}, {ID: "a.go", Kind: mapx.KindCode},
		{ID: "README.md", Kind: mapx.KindDoc},
	}}
	cfg := &config.Config{Gates: []config.Gate{{Name: "g", On: []string{"feature"}}}}
	var got []string
	for _, f := range checkGateCoverage(g, cfg) {
		if f.Check != "kind-sem-gate" || f.Severity != Warn {
			t.Errorf("wrong finding: %+v", f)
		}
		got = append(got, f.Subject)
	}
	if strings.Join(got, ",") != "spec,test" {
		t.Fatalf("spec and test lack a gate (feature has one, code and doc need none), got %v", got)
	}
}

// A typo in `skip_on` is the most misleading error: it matches no perspective, the gate
// keeps running in both, and the author believes it switched it off in one. A silent
// failure on the dangerous side — whoever wrote it wanted LESS execution and got more.
func TestCheckSkipOnValid(t *testing.T) {
	t.Run("DCTRO-B15: An unknown perspective in skip_on names the gate and the value", func(t *testing.T) {})
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "with-typo", SkipOn: []string{"chnage"}},
		{Name: "right", SkipOn: []string{config.PerspectiveAll}},
		{Name: "undeclared"},
	}}

	fs := checkSkipOnValid(cfg)
	if len(fs) != 1 {
		t.Fatalf("only the typo should be accused, got %d: %+v", len(fs), fs)
	}
	if fs[0].Subject != "with-typo" {
		t.Errorf("the finding must name the gate: %+v", fs[0])
	}
	if !strings.Contains(fs[0].Detail, "chnage") {
		t.Errorf("the finding must quote the invalid value: %q", fs[0].Detail)
	}
}

// A gate skipped for a missing tool is indistinguishable from a gate that passed, for
// whoever only looks at the green of `check`. The doctor is where that difference must
// show — without this warning the project would be uncovered exactly as long as nobody
// noticed.
func TestDoctorWarnsMissingToolAsWarn(t *testing.T) {
	t.Run("DCTRO-B16: A gate whose tool is not on the PATH is reported with its install hint", func(t *testing.T) {})
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "no-secret-leaked", NeedsTool: "binario-inexistente-xyz", InstallHint: "brew install foo"},
		{Name: "gate-without-tool"},
		{Name: "gate-with-tool", NeedsTool: "sh"},
	}}

	fs := checkMissingTools(cfg)
	if len(fs) != 1 {
		t.Fatalf("only the gate with a MISSING tool should give a finding; got %d: %+v", len(fs), fs)
	}
	if fs[0].Severity != Warn {
		t.Errorf("the declared coverage is not the real one — that is Warn, got %q", fs[0].Severity)
	}
	if fs[0].Subject != "no-secret-leaked" {
		t.Errorf("the finding must point at the disabled GATE, got %q", fs[0].Subject)
	}
	// The warning must end in an action: without the hint the reader knows the problem
	// and not the fix.
	if !strings.Contains(fs[0].Detail, "brew install foo") {
		t.Errorf("the warning must carry the install_hint; got: %q", fs[0].Detail)
	}
}

// --- git ---

// rootOutsideRepo returns a directory under no `.git` — otherwise the test would be a
// false positive when run from inside Anchors' own repo.
func rootOutsideRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if hasRepo(dir) {
		t.Skipf("the temporary directory %s is inside a git repo", dir)
	}
	return dir
}

// The doctor exists to warn AHEAD. A missing git is the same class as
// `ferramenta-ausente`: nothing fails loudly, so the silence is what costs.
func TestDoctorWarnsProjectWithoutRepository(t *testing.T) {
	t.Run("DCTRO-B18: A project under no repository is told to initialize one", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed — this test covers the 'git exists, repo does not' case")
	}
	dir := rootOutsideRepo(t)

	fs := checkGitMissing(&config.Config{}, dir)

	if len(fs) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(fs), fs)
	}
	if fs[0].Check != "git-ausente" || fs[0].Severity != Warn {
		t.Errorf("wrong finding: %+v", fs[0])
	}
	// The message must say the FIX — initialize, not install.
	if !strings.Contains(fs[0].Detail, "git init") {
		t.Errorf("the message should say to initialize the repo: %s", fs[0].Detail)
	}
	if strings.Contains(fs[0].Detail, "PATH") {
		t.Errorf("git IS installed here — pointing at the PATH sends the reader where the problem is not: %s", fs[0].Detail)
	}
}

// An existing repo is not a finding — an "alert" that never demands action is noise.
func TestDoctorDoesNotComplainWithRepository(t *testing.T) {
	t.Run("DCTRO-B19: A .git in the root or an ancestor, as a folder or a file, is a repository", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := rootOutsideRepo(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	if fs := checkGitMissing(&config.Config{}, dir); len(fs) != 0 {
		t.Errorf("a present repo should not give a finding: %+v", fs)
	}
}

// A subfolder of a repo IS already versioned: complaining there would ask for a nested repo.
func TestDoctorDoesNotComplainInRepoSubfolder(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := rootOutsideRepo(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "packages", "app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if fs := checkGitMissing(&config.Config{}, sub); len(fs) != 0 {
		t.Errorf("a repo subfolder is already versioned: %+v", fs)
	}
}

// Worktrees and submodules have `.git` as a FILE (a `gitdir:` pointer). Treating it as
// "no repo" would ask to reinitialize on top of a valid worktree.
func TestDoctorAcceptsGitAsFile(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := rootOutsideRepo(t)
	if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /other/.git/worktrees/x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if fs := checkGitMissing(&config.Config{}, dir); len(fs) != 0 {
		t.Errorf("`.git` as a file is a worktree/submodule — a real repo exists: %+v", fs)
	}
}

// In `github` mode the missing repository stops being debt and becomes a blocker: the
// work queue LIVES in a repository's issues. The message must say so.
func TestDoctorInGitHubModeSaysTheQueueHasNowhereToComeFrom(t *testing.T) {
	t.Run("DCTRO-B20: In GitHub mode a missing repository names the work queue", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := rootOutsideRepo(t)
	cfg := &config.Config{Workflow: &config.Workflow{
		Mode:   config.ModeGitHub,
		Repo:   "acme/exemplo",
		Labels: []string{"anchors"},
	}}

	fs := checkGitMissing(cfg, dir)

	if len(fs) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(fs))
	}
	if !strings.Contains(fs[0].Detail, "fila de trabalho") && !strings.Contains(fs[0].Detail, "work queue") {
		t.Errorf("in github mode the message must name the real blocker: %s", fs[0].Detail)
	}
}

// The case the doctor covers and `init` cannot: git NOT INSTALLED. The fix is to install,
// not to initialize — asking for `git init` on a machine without git sends the user to
// the wrong place.
func TestDoctorWithoutBinaryAsksToInstallNotInitialize(t *testing.T) {
	t.Run("DCTRO-B17: A missing git binary is told to install, not to initialize", func(t *testing.T) {})
	// `installed=false` injected: without it this case would only run on a machine
	// without git, and stay SKIP forever.
	fs := gitMissing(&config.Config{}, t.TempDir(), false)

	if len(fs) != 1 || fs[0].Subject != "git" {
		t.Fatalf("expected a finding about the binary: %+v", fs)
	}
	if strings.Contains(fs[0].Detail, "git init") {
		t.Errorf("without the binary `git init` does not run — asking for it is the wrong place: %s", fs[0].Detail)
	}
}

// --- plan needs ---

func plans(ns ...mapx.Node) *mapx.Graph {
	g := &mapx.Graph{}
	for _, n := range ns {
		n.Kind = mapx.KindPlan
		g.Nodes = append(g.Nodes, n)
	}
	// the edges the build would create: only to an existing target
	exists := map[string]bool{}
	for _, n := range g.Nodes {
		exists[n.ID] = true
	}
	for _, n := range g.Nodes {
		for _, target := range n.Needs {
			if exists[target] {
				g.Edges = append(g.Edges, mapx.Edge{From: n.ID, To: target, Type: mapx.EdgeNeeds})
			}
		}
	}
	return g
}

// A `needs` to a missing plan is INVISIBLE in the graph — the build drops the edge. And
// the effect is the declared work order simply not holding.
func TestNeedsToMissingPlanIsReported(t *testing.T) {
	t.Run("DCTRO-B21: A needs to a plan that does not exist is broken", func(t *testing.T) {})
	g := plans(
		mapx.Node{ID: "plans/0001-base.md"},
		mapx.Node{ID: "plans/0002-feature.md", Needs: []string{"plans/0099-nao-existe.md"}},
	)

	fs := checkPlanNeeds(g)

	if len(fs) != 1 || fs[0].Check != "needs-quebrado" {
		t.Fatalf("expected needs-quebrado, got %+v", fs)
	}
	if !strings.Contains(fs[0].Detail, "0099-nao-existe") {
		t.Errorf("the message must name the missing target: %s", fs[0].Detail)
	}
}

// A cycle: none of the plans can ever start. Without this check the board holds cards
// nobody takes and nobody knows why.
func TestNeedsCycleIsReportedWithThePath(t *testing.T) {
	t.Run("DCTRO-B22: A cycle of needs gives one finding naming only the cycle, the same on every run", func(t *testing.T) {})
	g := plans(
		mapx.Node{ID: "plans/a.md", Needs: []string{"plans/b.md"}},
		mapx.Node{ID: "plans/b.md", Needs: []string{"plans/a.md"}},
	)

	fs := checkPlanNeeds(g)

	if len(fs) != 1 || fs[0].Check != "needs-ciclo" {
		t.Fatalf("expected needs-ciclo, got %+v", fs)
	}
	// The path is what makes the finding actionable.
	if !strings.Contains(fs[0].Detail, "→") {
		t.Errorf("the message must show the cycle's path: %s", fs[0].Detail)
	}
}

// A legitimate chain (A ← B ← C) is not a cycle, and a project that declares a correct
// order cannot get a finding.
func TestLegitimateChainIsNotAFinding(t *testing.T) {
	t.Run("DCTRO-B23: A chain in order, or no plan at all, gives nothing", func(t *testing.T) {})
	g := plans(
		mapx.Node{ID: "plans/0001-fundacao.md"},
		mapx.Node{ID: "plans/0002-backend.md", Needs: []string{"plans/0001-fundacao.md"}},
		mapx.Node{ID: "plans/0003-tela.md", Needs: []string{"plans/0002-backend.md"}},
	)

	if fs := checkPlanNeeds(g); len(fs) != 0 {
		t.Errorf("a chain in order is not a finding: %+v", fs)
	}
}

// A project without any plan: nothing to check.
func TestNoPlansReportsNothing(t *testing.T) {
	if fs := checkPlanNeeds(&mapx.Graph{}); len(fs) != 0 {
		t.Errorf("without plans there are no needs to check: %+v", fs)
	}
}

// --- signals ---

func TestCheckMissingSignals(t *testing.T) {
	code := func(id string, s *mapx.TestSignal) mapx.Node {
		return mapx.Node{ID: id, Kind: mapx.KindCode, Signal: s}
	}
	test := func(id string, s *mapx.TestSignal) mapx.Node {
		return mapx.Node{ID: id, Kind: mapx.KindTest, Signal: s}
	}
	find := func(fs []Finding, subject string) (Finding, bool) {
		for _, f := range fs {
			if f.Check == "sinal-ausente" && f.Subject == subject {
				return f, true
			}
		}
		return Finding{}, false
	}

	t.Run("DCTRO-B24: Tests without results and code without coverage are warnings", func(t *testing.T) {
		fs := checkMissingSignals(&mapx.Graph{Nodes: []mapx.Node{test("a_test.go", nil), code("a.go", nil)}})
		if f, ok := find(fs, "execução (JUnit)"); !ok || f.Severity != Warn {
			t.Errorf("tests without execution results are a warning: %+v", fs)
		}
		if f, ok := find(fs, "cobertura (lcov)"); !ok || f.Severity != Warn {
			t.Errorf("code without coverage is a warning: %+v", fs)
		}
		ran := &mapx.TestSignal{Passed: 1, TotalLines: 10}
		fs = checkMissingSignals(&mapx.Graph{Nodes: []mapx.Node{test("a_test.go", ran), code("a.go", ran)}})
		if _, ok := find(fs, "execução (JUnit)"); ok {
			t.Errorf("one test with results is enough: %+v", fs)
		}
		if _, ok := find(fs, "cobertura (lcov)"); ok {
			t.Errorf("one file with coverage is enough: %+v", fs)
		}
	})

	t.Run("DCTRO-B25: Missing or partial mutation is informational", func(t *testing.T) {
		fs := checkMissingSignals(&mapx.Graph{Nodes: []mapx.Node{code("a.go", nil), code("b.go", nil)}})
		if f, ok := find(fs, "mutação"); !ok || f.Severity != Info {
			t.Errorf("no mutation signal is informational: %+v", fs)
		}
		killed := &mapx.TestSignal{MutantsKilled: 3, TotalLines: 5}
		fs = checkMissingSignals(&mapx.Graph{Nodes: []mapx.Node{code("a.go", killed), code("b.go", nil)}})
		f, ok := find(fs, "mutação (parcial)")
		if !ok || f.Severity != Info || !strings.Contains(f.Detail, "1 of 2") && !strings.Contains(f.Detail, "1 de 2") {
			t.Errorf("partial mutation is informational with m of n: %+v", fs)
		}
		fs = checkMissingSignals(&mapx.Graph{Nodes: []mapx.Node{code("a.go", killed), code("b.go", killed)}})
		if len(fs) != 0 {
			t.Errorf("every file with coverage and mutation gives nothing: %+v", fs)
		}
	})
}

// --- the report ---

func TestDiagnose_sortsAndCounts(t *testing.T) {
	t.Run("DCTRO-B01: The report is sorted by check then subject, and counts the map", func(t *testing.T) {})
	root, g, cfg := setup(t)
	g.Nodes = append(g.Nodes, mapx.Node{ID: "A.spec.md", Kind: mapx.KindSpec, Code: "AAAA"})
	r := Diagnose(g, cfg, root)
	if r.Nodes != len(g.Nodes) || r.Edges != len(g.Edges) || r.Layers != len(cfg.Layers) {
		t.Fatalf("counts = %d/%d/%d, want %d/%d/%d", r.Nodes, r.Edges, r.Layers, len(g.Nodes), len(g.Edges), len(cfg.Layers))
	}
	if len(r.Findings) < 3 {
		t.Fatalf("setup: the toy project has several findings: %+v", r.Findings)
	}
	for i := 1; i < len(r.Findings); i++ {
		a, b := r.Findings[i-1], r.Findings[i]
		if a.Check > b.Check || a.Check == b.Check && a.Subject > b.Subject {
			t.Fatalf("findings out of order at %d: %+v then %+v", i, a, b)
		}
	}
}

func TestReport_Warnings(t *testing.T) {
	t.Run("DCTRO-B02: Warnings keeps only the warnings", func(t *testing.T) {})
	r := Report{Findings: []Finding{
		{"a", Warn, "x", ""}, {"b", Info, "y", ""}, {"c", Warn, "z", ""},
	}}
	w := r.Warnings()
	if len(w) != 2 || w[0].Check != "a" || w[1].Check != "c" {
		t.Fatalf("Warnings = %+v, want the two warnings in order", w)
	}
}

// A project with nothing wrong gives the doctor nothing to warn about: every check that
// finds nothing adds nothing.
func TestDiagnose_healthyProjectHasNoWarnings(t *testing.T) {
	t.Run("DCTRO-I01: A healthy project has no warnings", func(t *testing.T) {})
	root := rootOutsideRepo(t)
	os.MkdirAll(filepath.Join(root, ".git"), 0o755)
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	for _, f := range []string{"pkg/a.spec.md", "pkg/a.go", "pkg/a_test.go", "guide.md"} {
		os.WriteFile(filepath.Join(root, f), []byte("# A\n"), 0o644)
	}
	full := &mapx.TestSignal{Passed: 1, TotalLines: 3, MutantsKilled: 1}
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "pkg/a.spec.md", Kind: mapx.KindSpec, Code: "AAAA"},
			{ID: "pkg/a.go", Kind: mapx.KindCode, Code: "AAAA", Signal: full},
			{ID: "pkg/a_test.go", Kind: mapx.KindTest, Code: "AAAA", Signal: full},
			{ID: "guide.md", Kind: mapx.KindGuide},
		},
		Edges: []mapx.Edge{
			{From: "pkg/a.spec.md", To: "pkg/a.go", Type: mapx.EdgeSpecifies},
			{From: "guide.md", To: "pkg/a.spec.md", Type: mapx.EdgeGoverns},
		},
	}
	cfg := &config.Config{
		Layers: map[string]config.Layer{"spec": {Kind: "spec"}, "code": {Kind: "code"}, "test": {Kind: "test"}},
		Gates:  []config.Gate{{Name: "g", On: []string{"spec", "test"}}},
	}
	if w := Diagnose(g, cfg, root).Warnings(); len(w) != 0 {
		t.Fatalf("a healthy project has no warnings: %+v", w)
	}
}

// The DFS started from each plan in map order, and the reported path began at the plan
// the walk started from: with `x → a ⇄ b` the finding changed between runs and could read
// "x → a → b → a", naming a plan that is not in the cycle.
func TestNeedsCycleNamesOnlyTheCycleDeterministically(t *testing.T) {
	t.Run("DCTRO-B22: A cycle of needs gives one finding naming only the cycle, the same on every run", func(t *testing.T) {
		g := plans(
			mapx.Node{ID: "plans/c.md", Needs: []string{"plans/d.md"}},
			mapx.Node{ID: "plans/a.md", Needs: []string{"plans/b.md"}},
			mapx.Node{ID: "plans/0.md", Needs: []string{"plans/a.md"}}, // walked first: leads into the cycle
			mapx.Node{ID: "plans/b.md", Needs: []string{"plans/a.md"}},
			mapx.Node{ID: "plans/d.md", Needs: []string{"plans/c.md"}},
		)
		want := checkPlanNeeds(g)
		if len(want) != 1 || want[0].Subject != "plans/b.md" ||
			want[0].Detail != i18n.T("health.needs_cycle", "plans/a.md → plans/b.md → plans/a.md") {
			t.Fatalf("want the a⇄b cycle, found first in path order; got %+v", want)
		}
		for i := 0; i < 50; i++ {
			if got := checkPlanNeeds(g); len(got) != 1 || got[0] != want[0] {
				t.Fatalf("run %d gave %+v, want %+v", i, got, want[0])
			}
		}
	})
}
