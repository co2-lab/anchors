package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

const docsRulesTmpl = "{{range specs}}{{range rules .}}{{.Code}} {{.Titulo}}\n{{end}}{{end}}"

// docsRepo is the sync repository with a template listing the rules, and its page built
// and committed.
func docsRepo(t *testing.T) syncRepo {
	t.Helper()
	r := newSyncRepo(t, true)
	touchWrite(t, r.root, "doct/rules.md.tmpl", docsRulesTmpl)
	g, err := mapx.Load(filepath.Join(r.root, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	c, err := doct.New(r.root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	r.git("add", ".")
	r.git("commit", "-qm", "docs")
	return r
}

func (r syncRepo) index(t *testing.T, rel string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", r.root, "show", ":./"+rel).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func (r syncRepo) disk(t *testing.T, rel string) string {
	t.Helper()
	b, _ := os.ReadFile(filepath.Join(r.root, filepath.FromSlash(rel)))
	return string(b)
}

// stageSpec changes the spec's rule title and stages it.
func (r syncRepo) stageSpec(t *testing.T, title string) {
	t.Helper()
	touchWrite(t, r.root, "src/pay.spec.md", "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\n### PAYMX-B01 — "+title+"\n")
	r.git("add", "src/pay.spec.md")
}

func TestDocsSync_aStalePageIsCompiledAndStaged(t *testing.T) {
	t.Run("DCSYN-B01: A page left out of date is compiled and staged", func(t *testing.T) {})
	r := docsRepo(t)
	r.stageSpec(t, "charges the amount")
	msg, err := syncDocsForCommit(r.root, r.cfg)
	if err != nil || !strings.Contains(msg, "docs/rules.md") || !strings.Contains(msg, "written and staged") {
		t.Fatalf("the stale page is written and staged, got %q %v", msg, err)
	}
	if !strings.Contains(r.index(t, "docs/rules.md"), "PAYMX-B01 charges the amount") || r.disk(t, "docs/rules.md") != r.index(t, "docs/rules.md") {
		t.Errorf("the page is compiled from the staged spec, on disk and in the index alike:\n%s", r.index(t, "docs/rules.md"))
	}
	if msg, err := syncDocsForCommit(r.root, r.cfg); err != nil || msg != "" {
		t.Errorf("an up-to-date page is not compiled again, got %q %v", msg, err)
	}
}

func TestDocsSync_theTreeAheadStagesApart(t *testing.T) {
	t.Run("DCSYN-B02: With the tree ahead, the page goes straight into the index", func(t *testing.T) {})
	r := docsRepo(t)
	before := r.disk(t, "docs/rules.md")
	r.stageSpec(t, "charges the amount")
	touchWrite(t, r.root, "src/pay.spec.md", "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\n### PAYMX-B01 — half done elsewhere\n")
	touchWrite(t, r.root, "src/other.ts", "export const other = 2 // not staged\n")
	msg, err := syncDocsForCommit(r.root, r.cfg)
	if err != nil || !strings.Contains(msg, "staged apart") {
		t.Fatalf("the page is staged apart, got %q %v", msg, err)
	}
	if got := r.index(t, "docs/rules.md"); !strings.Contains(got, "charges the amount") || strings.Contains(got, "half done") {
		t.Errorf("the index holds the page of the staged spec:\n%s", got)
	}
	if r.disk(t, "docs/rules.md") != before {
		t.Error("the page on disk is as it was")
	}
}

func TestDocsSync_nothingToCompile(t *testing.T) {
	t.Run("DCSYN-B03: Nothing to compile, nothing staged", func(t *testing.T) {})
	off := false
	r := docsRepo(t)
	r.stageSpec(t, "charges the amount")
	cfg := *r.cfg
	cfg.Docs = &config.Docs{PreCommit: &off}
	if msg, err := syncDocsForCommit(r.root, &cfg); err != nil || msg != "" {
		t.Errorf("turned off, nothing is compiled, got %q %v", msg, err)
	}
	if msg, err := syncDocsForCommit(r.root, nil); err != nil || msg != "" {
		t.Errorf("with no configuration, nothing is compiled, got %q %v", msg, err)
	}
	noTemplates := newSyncRepo(t, true)
	noTemplates.stageSpec(t, "charges the amount")
	if msg, err := syncDocsForCommit(noTemplates.root, noTemplates.cfg); err != nil || msg != "" {
		t.Errorf("without doct/, nothing is compiled, got %q %v", msg, err)
	}
	noMap := docsRepo(t)
	noMap.git("rm", "-q", "--cached", mapx.DefaultPath)
	if err := os.Remove(filepath.Join(noMap.root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	noMap.stageSpec(t, "charges the amount")
	if msg, err := syncDocsForCommit(noMap.root, noMap.cfg); err != nil || msg != "" {
		t.Errorf("without a map, nothing is compiled, got %q %v", msg, err)
	}
	if strings.Contains(noMap.index(t, "docs/rules.md"), "charges the amount") {
		t.Error("nothing was staged")
	}
}

func TestDocsSync_neverAHandwrittenOrIgnoredPage(t *testing.T) {
	t.Run("DCSYN-X01: A page written by hand, or one git ignores, is never touched", func(t *testing.T) {})
	r := docsRepo(t)
	touchWrite(t, r.root, "doct/hand.md.tmpl", docsRulesTmpl)
	touchWrite(t, r.root, "docs/hand.md", "written by hand\n")
	touchWrite(t, r.root, "doct/ign.md.tmpl", docsRulesTmpl)
	touchWrite(t, r.root, ".gitignore", "docs/ign.md\n")
	r.git("add", ".")
	r.git("commit", "-qm", "more pages")
	r.stageSpec(t, "charges the amount")
	if _, err := syncDocsForCommit(r.root, r.cfg); err != nil {
		t.Fatal(err)
	}
	if r.disk(t, "docs/hand.md") != "written by hand\n" || r.index(t, "docs/hand.md") != "written by hand\n" {
		t.Error("the page written by hand is not written nor staged")
	}
	if r.disk(t, "docs/ign.md") != "" || r.index(t, "docs/ign.md") != "" {
		t.Error("the page git ignores is not written nor staged")
	}
}

func TestDocsSync_nothingOutsideDocs(t *testing.T) {
	t.Run("DCSYN-X02: Nothing outside docs is written or staged", func(t *testing.T) {})
	r := docsRepo(t)
	r.stageSpec(t, "charges the amount")
	status := func() map[string]bool {
		out := r.git("status", "--porcelain")
		m := map[string]bool{}
		for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
			if len(l) > 3 {
				m[l] = true
			}
		}
		return m
	}
	before := status()
	if _, err := syncDocsForCommit(r.root, r.cfg); err != nil {
		t.Fatal(err)
	}
	for l := range status() {
		if !before[l] && !strings.HasPrefix(l[3:], "docs/") {
			t.Errorf("the sync changed %q, outside docs/", l)
		}
	}
}

func TestDocsSync_theGateFindsThePagesUpToDate(t *testing.T) {
	t.Run("DCSYN-I01: After the sync the gate finds the pages up to date", func(t *testing.T) {})
	r := docsRepo(t)
	r.stageSpec(t, "charges the amount")
	touchWrite(t, r.root, "src/other.ts", "export const other = 2 // not staged\n")
	if _, err := syncDocsForCommit(r.root, r.cfg); err != nil {
		t.Fatal(err)
	}
	read, err := scan.IndexReader(r.root)
	if err != nil {
		t.Fatal(err)
	}
	c, err := doct.NewWith(r.root, stagedMap(r.root, filepath.Join(r.root, mapx.DefaultPath)), read)
	if err != nil {
		t.Fatal(err)
	}
	if stale, err := c.Stale(); err != nil || len(stale) != 0 {
		t.Errorf("confronted with the index, no page is out of date, got %v %v", stale, err)
	}
}

func TestDocsSync_aBrokenTemplateIsAnError(t *testing.T) {
	t.Run("DCSYN-E01: A template that does not compile comes back as an error", func(t *testing.T) {})
	r := docsRepo(t)
	touchWrite(t, r.root, "doct/rules.md.tmpl", "{{range specs}}\n")
	r.git("add", "doct/rules.md.tmpl")
	r.stageSpec(t, "charges the amount")
	before := r.index(t, "docs/rules.md")
	if _, err := syncDocsForCommit(r.root, r.cfg); err == nil || !strings.Contains(err.Error(), "rules.md.tmpl") {
		t.Errorf("the error names the template, got %v", err)
	}
	if r.index(t, "docs/rules.md") != before {
		t.Error("no page is staged")
	}
}
