package flow

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// ghRule is one answer of the scripted `gh`: when the joined arguments match the shell
// glob `match`, it prints `out` and exits with `code`. The first matching rule wins.
type ghRule struct {
	match string
	out   string
	code  int
	// stdinTo, when set, is the file the call's standard input is saved to — how
	// `gh issue comment --body-file -` receives its body.
	stdinTo string
}

// scriptedGH puts on the PATH a `gh` that answers by rule and records every call, one
// line per call, in order. Nothing reaches the network: a call no rule matches exits 0
// with no output, which most callers in this package read as "nothing found". Not
// `decided`'s unblock lookup: there, unreadable output is "unknown" and refuses, so
// a test that means "no unblock work" scripts the `[]` gh really prints.
func scriptedGH(t *testing.T, rules ...ghRule) (calls func() []string) {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "calls.txt")
	var s strings.Builder
	s.WriteString("#!/bin/sh\n")
	// ONE LINE PER CALL: a `--body` carries newlines, and they would split one call into
	// several lines of the log.
	s.WriteString("printf '%s' \"$*\" | tr '\\n' ' ' >> '" + log + "'\n")
	s.WriteString("printf '\\n' >> '" + log + "'\n")
	s.WriteString("case \"$*\" in\n")
	for _, r := range rules {
		s.WriteString(shellGlob(r.match) + ")\n")
		if r.stdinTo != "" {
			s.WriteString("cat > '" + r.stdinTo + "'\n")
		}
		if r.out != "" {
			s.WriteString("cat <<'__GH_EOF__'\n" + r.out + "\n__GH_EOF__\n")
		}
		s.WriteString("exit " + strconv.Itoa(r.code) + " ;;\n")
	}
	s.WriteString("esac\nexit 0\n")
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(s.String()), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	// The agent's own cards are looked up through `gh` when ANCHORS_AGENT is set; a test
	// that does not ask for it must not inherit the caller's.
	t.Setenv("ANCHORS_AGENT", "")
	return func() []string {
		b, err := os.ReadFile(log)
		if err != nil {
			return nil
		}
		var out []string
		for _, l := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(l) != "" {
				out = append(out, l)
			}
		}
		return out
	}
}

// shellGlob quotes the literal parts of a pattern, so spaces and quotes in it are text and
// only `*` stays a wildcard.
func shellGlob(match string) string {
	parts := strings.Split(match, "*")
	for i, p := range parts {
		if p != "" {
			parts[i] = "'" + strings.ReplaceAll(p, "'", `'\''`) + "'"
		}
	}
	return strings.Join(parts, "*")
}

// callsWith returns the recorded calls that contain every given fragment.
func callsWith(calls []string, fragments ...string) []string {
	var out []string
	for _, c := range calls {
		ok := true
		for _, f := range fragments {
			if !strings.Contains(c, f) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, c)
		}
	}
	return out
}

// githubProject writes a project in github mode and returns its root.
func githubProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	yaml := "version: 2\nworkflow:\n  mode: github\n  repo: acme/app\n  labels: [anchors]\n" +
		"layers:\n  logic:\n    pattern: \"src/**/*.ts\"\n    kind: code\n"
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// localProject writes a project in local mode and returns its root.
func localProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	yaml := "version: 1\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n"
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// writeFile creates `rel` under root with the given content, making the directories.
func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// gitRepo turns dir into a real git repository with one commit of whatever is in it.
func gitRepo(t *testing.T, dir string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"},
		{"config", "user.email", "t@t"},
		{"config", "user.name", "t"},
		{"config", "commit.gpgsign", "false"},
		{"add", "-A"},
		{"commit", "-q", "--allow-empty", "-m", "base"},
	} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}
}

// cfgGitHub is the in-memory config of githubProject.
func cfgGitHub() *config.Config {
	return &config.Config{Workflow: &config.Workflow{Mode: config.ModeGitHub, Repo: "acme/app", Labels: []string{"anchors"}}}
}
