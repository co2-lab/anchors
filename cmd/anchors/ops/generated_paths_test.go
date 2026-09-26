package ops

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func generatedPaths(t *testing.T, root string, args ...string) (string, error) {
	t.Helper()
	c := newGeneratedPathsCmd()
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SetArgs(append([]string{"--root", root}, args...))
	err := c.Execute()
	return out.String(), err
}

// The `re` form is what a conflict script feeds to `grep -E`: it must tell the generated
// files from the work files — and the dot must be literal.
func TestGeneratedPathsRegexSeparatesNoiseFromDivergence(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, "version: 1\n")
	out, err := generatedPaths(t, root, "--format", "re")
	if err != nil {
		t.Fatal(err)
	}
	re, err := regexp.Compile(strings.TrimSpace(out))
	if err != nil {
		t.Fatalf("the alternation is not a valid regex: %v (%q)", err, out)
	}
	for path, generated := range map[string]bool{
		"anchors.graph.yaml":         true,
		"docs/arquitetura.md":        true,
		"plans/0001-progress.md":     true,
		"anchorsXgraph.yaml":         false, // an unescaped dot would match this
		"sub/anchors.graph.yaml":     false,
		"internal/gate/x.spec.md":    false,
		"plans/0001-progress.md.bak": false,
	} {
		if re.MatchString(path) != generated {
			t.Errorf("%s: generated = %v, want %v (regex %q)", path, !generated, generated, out)
		}
	}
}

// The default form is one pattern per line — the same patterns.
func TestGeneratedPathsLinesAndConfigRequired(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, "version: 1\n")
	lines, err := generatedPaths(t, root)
	if err != nil {
		t.Fatal(err)
	}
	alt, _ := generatedPaths(t, root, "--format", "re")
	if got := strings.Join(strings.Fields(lines), "|"); got != strings.TrimSpace(alt) {
		t.Errorf("lines %q and the alternation %q disagree", lines, alt)
	}
	if n := len(strings.Fields(lines)); n != 3 {
		t.Errorf("want three patterns (map, docs, progress), got %d:\n%s", n, lines)
	}
	if _, err := generatedPaths(t, t.TempDir()); err == nil {
		t.Error("a project without anchors.yaml must be refused")
	}
}
