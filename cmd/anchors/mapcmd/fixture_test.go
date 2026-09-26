package mapcmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// fixtureYAML is a small project with the whole chain: a guide that governs the code, a
// spec that specifies it, a test that proves it, a lonely guide (an orphan), one judgment
// gate over the code and one external gate that always fails on it.
const fixtureYAML = `version: 4
layers:
  spec: {pattern: "src/*.spec.md", kind: spec}
  code: {pattern: "src/*.ts", kind: code, tags: [code], exclude: ["src/*.test.ts"]}
  test: {pattern: "src/*.test.ts", kind: test}
  guide: {pattern: "guides/*.md", kind: guide}
derived:
  anchor: code
  files:
    spec: ["{{dir}}/{{name}}.spec.md"]
    test: ["{{dir}}/{{name}}.test.{{ext}}"]
governs:
  - from: guides/CODE.md
    governs: code
gates:
  - name: atomic
    measures: judgment
    on: [code]
    guide: guides/CODE.md
  - name: always-red
    on: [code]
    run: "false"
`

const fixtureSpec = "<!-- @anchors\n  code: LOGIN\n-->\n# Login\n\nLOGIN-B01 — the password is checked.\n"

// fixtureProject writes the project and builds its map with the real `map build`.
func fixtureProject(t *testing.T) string {
	t.Helper()
	return fixtureProjectWith(t, "", fixtureSpec)
}

// fixtureProjectWith is fixtureProject with more anchors.yaml and another spec text.
func fixtureProjectWith(t *testing.T, extraYAML, spec string) string {
	t.Helper()
	useEnglish(t)
	root := t.TempDir()
	for p, c := range map[string]string{
		"anchors.yaml":      fixtureYAML + extraYAML,
		"src/login.spec.md": spec,
		"src/login.ts":      "export const login = () => true\n",
		"src/login.test.ts": "test('LOGIN-B01: checks the password', () => {})\n",
		"guides/CODE.md":    "# Code guide\n",
		"guides/LONELY.md":  "# Nobody points here\n",
	} {
		writeProjectFile(t, root, p, c)
	}
	runCmd(t, newMapCmd(), "build", "--root", root)
	return root
}

func writeProjectFile(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// useEnglish pins the message catalog, so assertions on text do not depend on what an
// earlier test left selected.
func useEnglish(t *testing.T) {
	t.Helper()
	prev := i18n.Current()
	if err := i18n.Set("en"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = i18n.Set(prev) })
}

// runCmd executes a command and returns its stdout; a failure is fatal.
func runCmd(t *testing.T, cmd *cobra.Command, args ...string) string {
	t.Helper()
	out, err := runCmdErr(cmd, t, args...)
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", cmd.Name(), args, err, out)
	}
	return out
}

// runCmdErr executes a command and returns its stdout and error.
func runCmdErr(cmd *cobra.Command, t *testing.T, args ...string) (string, error) {
	t.Helper()
	var err error
	cmd.SetArgs(args)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SilenceUsage = true
	out := capturaStdout(t, func() { err = cmd.Execute() })
	return out, err
}

func loadMap(t *testing.T, root string) *mapx.Graph {
	t.Helper()
	g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	return g
}
