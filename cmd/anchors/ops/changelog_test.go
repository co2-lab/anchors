// @anchors
//   code: CHTSC
//   ref: CLGCM

package ops

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/shell"
)

// historyRepo builds a repository: a subject commits, `tag:NAME` tags the last commit,
// `file:PATH=TEXT` writes a file (not committed).
func historyRepo(t *testing.T, steps ...string) string {
	t.Helper()
	dir := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(shell.WithoutRepoEnv(os.Environ()), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q", "-b", "main")
	for _, s := range steps {
		switch {
		case strings.HasPrefix(s, "tag:"):
			git("tag", strings.TrimPrefix(s, "tag:"))
		case strings.HasPrefix(s, "file:"):
			p, text, _ := strings.Cut(strings.TrimPrefix(s, "file:"), "=")
			if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, p)), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, p), []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
		default:
			git("commit", "-q", "--allow-empty", "-m", s)
		}
	}
	return dir
}

func changelogRun(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	t.Cleanup(func() { _ = i18n.Set("") })
	c := newChangelogCmd()
	var out, errb bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&errb)
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), errb.String(), err
}

func TestChangelog_Latest(t *testing.T) {
	t.Run("CLGCM-B01: The latest release by default", func(t *testing.T) {})
	root := historyRepo(t, "feat: one", "tag:v0.1.0", "feat: two", "tag:v0.2.0")
	out, _, err := changelogRun(t, "--root", root)
	if err != nil || !strings.Contains(out, "## v0.2.0") || strings.Contains(out, "v0.1.0") || !strings.Contains(out, "- two (") {
		t.Fatalf("only the latest release, got %v\n%s", err, out)
	}
}

func TestChangelog_Range(t *testing.T) {
	t.Run("CLGCM-B02: From a tag, or every release", func(t *testing.T) {})
	root := historyRepo(t, "feat: one", "tag:v0.1.0", "feat: two", "tag:v0.2.0", "fix: three", "tag:v0.3.0")
	out, _, err := changelogRun(t, "--root", root, "--from", "v0.1.0")
	if err != nil || strings.Contains(out, "## v0.1.0") || strings.Index(out, "## v0.3.0") > strings.Index(out, "## v0.2.0") || strings.Index(out, "## v0.2.0") < 0 {
		t.Errorf("the releases after the tag, newest first, got %v\n%s", err, out)
	}
	out, _, _ = changelogRun(t, "--root", root, "--all")
	if !(strings.Index(out, "## v0.3.0") < strings.Index(out, "## v0.2.0") && strings.Index(out, "## v0.2.0") < strings.Index(out, "## v0.1.0")) {
		t.Errorf("every release, newest first, got\n%s", out)
	}
}

func TestChangelog_Unreleased(t *testing.T) {
	t.Run("CLGCM-B03: What is not released yet", func(t *testing.T) {})
	root := historyRepo(t, "feat: one", "tag:v0.1.0", "feat: next")
	out, _, err := changelogRun(t, "--root", root, "--unreleased")
	if err != nil || !strings.HasPrefix(out, "## Unreleased") || !strings.Contains(out, "## v0.1.0") {
		t.Errorf("the unreleased heading on top, got %v\n%s", err, out)
	}
	out, _, err = changelogRun(t, "--root", historyRepo(t, "feat: first"))
	if err != nil || !strings.HasPrefix(out, "## Unreleased") || !strings.Contains(out, "first") {
		t.Errorf("with no tag the history is unreleased, got %v\n%s", err, out)
	}
}

func TestChangelog_WriteIncremental(t *testing.T) {
	t.Run("CLGCM-B04: Writing the incremental file", func(t *testing.T) {})
	root := historyRepo(t, "file:anchors.yaml=version: 1\nlang: pt-BR\n", "feat: one", "tag:v0.1.0", "fix: dois", "tag:v0.2.0")
	out, _, err := changelogRun(t, "--root", root, "--write", "--all")
	if err != nil || !strings.Contains(out, "wrote CHANGELOG.md") {
		t.Fatalf("writes CHANGELOG.md, got %v %s", err, out)
	}
	b, _ := os.ReadFile(filepath.Join(root, "CHANGELOG.md"))
	got := string(b)
	if !strings.HasPrefix(got, "# Changelog técnico\n\n"+changelog.Marker("v0.2.0")) || !strings.Contains(got, "### Correções") || !strings.Contains(got, changelog.Marker("v0.1.0")) {
		t.Fatalf("the Portuguese title and headings, newest first, got\n%s", got)
	}
	out, _, err = changelogRun(t, "--root", root, "--write", "--all")
	if again, _ := os.ReadFile(filepath.Join(root, "CHANGELOG.md")); err != nil || string(again) != got || !strings.Contains(out, "already holds every release") {
		t.Errorf("a second run leaves the file as it was, got %v %s", err, out)
	}
}

func TestChangelog_WritePerVersion(t *testing.T) {
	t.Run("CLGCM-B05: One file per release", func(t *testing.T) {})
	root := historyRepo(t, "file:anchors.yaml=version: 1\nchangelog:\n  mode: per_version\n  path: releases\n",
		"file:releases/v0.1.0.md=by hand\n", "file:releases/unreleased.md=old\n",
		"feat: one", "tag:v0.1.0", "feat: two", "tag:v0.2.0", "fix: three")
	if _, _, err := changelogRun(t, "--root", root, "--write", "--all", "--unreleased"); err != nil {
		t.Fatal(err)
	}
	read := func(n string) string { b, _ := os.ReadFile(filepath.Join(root, "releases", n)); return string(b) }
	if read("v0.1.0.md") != "by hand\n" {
		t.Errorf("an existing release file is kept, got %q", read("v0.1.0.md"))
	}
	if !strings.HasPrefix(read("v0.2.0.md"), "## v0.2.0") {
		t.Errorf("a missing release gets its file, got %q", read("v0.2.0.md"))
	}
	if u := read("unreleased.md"); strings.Contains(u, "old") || !strings.Contains(u, "three") {
		t.Errorf("the unreleased file is rewritten, got %q", u)
	}
	out, _, _ := changelogRun(t, "--root", root, "--write", "--all")
	if !strings.Contains(out, "every release already has its file") {
		t.Errorf("nothing left to write is reported, got %q", out)
	}
}

func TestChangelog_TemplateAndLang(t *testing.T) {
	t.Run("CLGCM-B06: The project's template and language", func(t *testing.T) {})
	root := historyRepo(t, "file:anchors.yaml=version: 1\nlang: es\nchangelog:\n  template: cl.tmpl\n",
		"file:cl.tmpl=[{{.Version}}] {{t \"changelog.features\"}}: {{range .Features}}{{.Subject}}{{end}}\n",
		"feat: uno", "tag:v1")
	out, _, err := changelogRun(t, "--root", root)
	if err != nil || out != "[v1] Funcionalidades: uno\n" {
		t.Fatalf("the project's template, headings in its language, got %v %q", err, out)
	}
	root = historyRepo(t, "file:anchors.yaml=version: 1\nlang: es\n", "fix: x", "tag:v1")
	if out, _, _ := changelogRun(t, "--root", root); !strings.Contains(out, "### Correcciones") {
		t.Errorf("the built-in headings follow lang, got %q", out)
	}
}

func TestChangelog_TemplateAtAURL(t *testing.T) {
	t.Run("CLGCM-B06: The project's template and language", func(t *testing.T) {})
	t.Run("CLGCM-E02: A template that cannot be read", func(t *testing.T) {})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cl.tmpl" {
			_, _ = w.Write([]byte("<{{.Version}}>\n"))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()
	root := historyRepo(t, "file:anchors.yaml=version: 1\nchangelog:\n  template: "+srv.URL+"/cl.tmpl\n", "feat: a", "tag:v1")
	if out, _, err := changelogRun(t, "--root", root); err != nil || out != "<v1>\n" {
		t.Fatalf("a template read from a URL: %v %q", err, out)
	}
	root = historyRepo(t, "file:anchors.yaml=version: 1\nchangelog:\n  template: "+srv.URL+"/nope.tmpl\n", "feat: a", "tag:v1")
	if _, _, err := changelogRun(t, "--root", root); err == nil || !strings.Contains(err.Error(), "changelog.template") {
		t.Errorf("a template URL that answers nothing: %v", err)
	}
}

func TestChangelog_NoConfig(t *testing.T) {
	t.Run("CLGCM-B07: No anchors.yaml", func(t *testing.T) {})
	root := historyRepo(t, "feat: one", "tag:v1")
	if _, _, err := changelogRun(t, "--root", root, "--write"); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(filepath.Join(root, "CHANGELOG.md")); err != nil || !strings.Contains(string(b), "## v1") {
		t.Errorf("CHANGELOG.md with the defaults, got %v %q", err, b)
	}
}

func TestChangelog_NothingToList(t *testing.T) {
	t.Run("CLGCM-B08: Nothing to list", func(t *testing.T) {})
	out, errOut, err := changelogRun(t, "--root", historyRepo(t, "chore: x", "tag:v1"))
	if err != nil || out != "" || !strings.Contains(errOut, "nothing to list") {
		t.Fatalf("nothing printed, the reason on stderr, got %v %q %q", err, out, errOut)
	}
}

func TestChangelog_UnknownFrom(t *testing.T) {
	t.Run("CLGCM-E01: An unknown start tag", func(t *testing.T) {})
	_, _, err := changelogRun(t, "--root", historyRepo(t, "feat: a", "tag:v1"), "--from", "v9")
	if err == nil || !strings.Contains(err.Error(), `"v9"`) {
		t.Fatalf("fails naming the tag, got %v", err)
	}
}

func TestChangelog_MissingTemplate(t *testing.T) {
	t.Run("CLGCM-E02: A template that cannot be read", func(t *testing.T) {})
	root := historyRepo(t, "file:anchors.yaml=version: 1\nchangelog:\n  template: nope.tmpl\n", "feat: a", "tag:v1")
	if _, _, err := changelogRun(t, "--root", root); err == nil || !strings.Contains(err.Error(), "changelog.template") {
		t.Fatalf("fails naming changelog.template, got %v", err)
	}
}

func TestChangelog_BrokenConfig(t *testing.T) {
	t.Run("CLGCM-E03: A broken anchors.yaml", func(t *testing.T) {})
	root := historyRepo(t, "file:anchors.yaml=version: 1\nchangelog:\n  mode: weekly\n", "feat: a", "tag:v1")
	if _, _, err := changelogRun(t, "--root", root); err == nil || !strings.Contains(err.Error(), "anchors.yaml") {
		t.Fatalf("fails naming the file, got %v", err)
	}
}

const productChangelog = `version: 8
versioning:
  app:
    scheme: semver
    source: [{ file: app/version.json, path: version }]
    tag: "app-v{version}"
    next: semantic
    branches:
      production:
        release: true
        changelog: true
      review:
        changelog: true
    changelog:
      outputs:
        - path: CHANGELOG.md
        - path: notes
          mode: per_version
          template: notes.tmpl
`

func TestChangelog_product(t *testing.T) {
	t.Run("CLGCM-B09: With versioning, the releases lie between the product's tags, ordered by its scheme and headed by the version", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+productChangelog, "file:app/version.json={\"version\": \"1.10.0\"}",
			"feat: nine", "tag:app-v1.9.0", "feat: other product", "tag:api-v5.0.0", "feat: ten", "tag:app-v1.10.0")
		out, _, err := changelogRun(t, "--root", root, "--branch", "production", "--all")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(out, "## 1.10.0 — ") || !strings.Contains(out, "## 1.9.0 — ") || strings.Contains(out, "api-v5") || strings.Contains(out, "app-v") {
			t.Errorf("headed by the versions, newest by the scheme first, no other product:\n%s", out)
		}
		if i, j := strings.Index(out, "## 1.10.0"), strings.Index(out, "- other product"); j < i {
			t.Errorf("the commits between two tags of the product are its release, whatever other tags sit between:\n%s", out)
		}
		if out, _, err := changelogRun(t, "--root", root, "--branch", "production", "--from", "1.9.0"); err != nil || !strings.HasPrefix(out, "## 1.10.0") || strings.Contains(out, "## 1.9.0") {
			t.Errorf("--from takes a version or a tag: %v\n%s", err, out)
		}
	})
	t.Run("CLGCM-B10: A branch that releases lists its releases, another the candidate under the version its strategy computes", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+productChangelog, "file:app/version.json={\"version\": \"1.1.0\"}",
			"feat: one", "tag:app-v1.0.0", "feat: two")
		out, _, err := changelogRun(t, "--root", root, "--branch", "review")
		if err != nil || !strings.HasPrefix(out, "## 1.1.0 — Unreleased") || strings.Contains(out, "## 1.0.0") || !strings.Contains(out, "- two (") {
			t.Errorf("review: the candidate alone: %v\n%s", err, out)
		}
		out, _, err = changelogRun(t, "--root", root)
		if err != nil || !strings.HasPrefix(out, "## 1.0.0 — ") {
			t.Errorf("by default, the first branch that releases: %v\n%s", err, out)
		}
	})
	t.Run("CLGCM-B11: --write writes every output of the product, and only on a branch the changelog is written on", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+productChangelog, "file:app/version.json={\"version\": \"1.0.0\"}",
			"file:notes.tmpl=Release {{.Version}}\n", "feat: one", "tag:app-v1.0.0")
		if _, _, err := changelogRun(t, "--root", root, "--branch", "production", "--write"); err != nil {
			t.Fatal(err)
		}
		if got := readFile(t, root, "CHANGELOG.md"); !strings.Contains(got, changelog.Marker("1.0.0")) || !strings.Contains(got, "## 1.0.0 — ") {
			t.Errorf("CHANGELOG.md:\n%s", got)
		}
		if got := readFile(t, root, "notes/1.0.0.md"); got != "Release 1.0.0\n" {
			t.Errorf("notes: %q", got)
		}
		if _, _, err := changelogRun(t, "--root", root, "--branch", "develop", "--write"); err == nil || !strings.Contains(err.Error(), "not written on develop") {
			t.Errorf("develop: %v", err)
		}
	})
	t.Run("CLGCM-E04: --from that is no tag of the product, or a template of an output that cannot be read, fails naming it", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+productChangelog, "file:app/version.json={\"version\": \"1.0.0\"}", "feat: one", "tag:app-v1.0.0", "tag:v9.0.0")
		if _, _, err := changelogRun(t, "--root", root, "--branch", "production", "--from", "v9.0.0"); err == nil || !strings.Contains(err.Error(), "no tag of app") {
			t.Errorf("from: %v", err)
		}
		if _, _, err := changelogRun(t, "--root", root, "--branch", "production", "--write"); err == nil || !strings.Contains(err.Error(), "versioning.app.changelog.outputs[1].template") {
			t.Errorf("template: %v", err)
		}
	})
}

func TestChangelog_intents(t *testing.T) {
	t.Run("CLGCM-B12: An output whose source is intents renders the intents waiting as the candidate, and a release from those archived under released/<version>/", func(t *testing.T) {
		cfg := strings.Replace(productChangelog, "        - path: notes\n          mode: per_version\n          template: notes.tmpl\n", "        - path: NOTES.md\n          source: intents\n", 1)
		cfg = strings.Replace(cfg, "version: 8\n", "version: 8\nchangelog:\n  intents: { types: [feat, fix] }\n", 1)
		cfg = strings.Replace(cfg, "next: semantic", "next: intents", 1)
		root := historyRepo(t, "file:anchors.yaml="+cfg, "file:app/version.json={\"version\": \"1.0.0\"}",
			"file:.anchors/changes/released/1.0.0/a.md=---\nproduct: app\nbump: minor\ntype: feat\n---\nThe first screen.\n",
			"file:.anchors/changes/b.md=---\nproduct: app\nbump: patch\ntype: fix\n---\nA crash on login.\n",
			"feat: one", "tag:app-v1.0.0", "fix: two")
		if _, _, err := changelogRun(t, "--root", root, "--branch", "production", "--write", "--all", "--unreleased"); err != nil {
			t.Fatal(err)
		}
		got := readFile(t, root, "NOTES.md")
		for _, want := range []string{changelog.Marker("") + "\n## 1.0.1 — Unreleased\n\n### Fixes\n\n- A crash on login.", changelog.Marker("1.0.0") + "\n## 1.0.0 — ", "### Features\n\n- The first screen."} {
			if !strings.Contains(got, want) {
				t.Errorf("missing %q in\n%s", want, got)
			}
		}
		if strings.Contains(readFile(t, root, "CHANGELOG.md"), "A crash on login") {
			t.Errorf("the commits' output does not read intents")
		}
	})
}
