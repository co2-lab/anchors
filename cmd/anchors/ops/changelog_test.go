// @anchors
//   ref: CLGCM

package ops

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/i18n"
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
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
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
