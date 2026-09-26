package mapcmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/testsig"
)

// A project that declares its code length has its JUnit cases read with it: a `[6]`
// project names its rules `LOGINX-B01`, and the permissive default (4 and 5) never read
// one of them — the rule showed no green test with its test passing.
func TestIngest_readsTheDeclaredCodeLength(t *testing.T) {
	prevLen := config.CodeLengths
	t.Cleanup(func() {
		config.SetCodeLengths(prevLen)
		testsig.SetCodeLenPattern("{4,5}")
	})
	useEnglish(t)
	root := t.TempDir()
	for p, c := range map[string]string{
		"anchors.yaml":      strings.Replace(fixtureYAML, "version: 4\n", "version: 4\ncode_lengths: [6]\n", 1),
		"src/login.spec.md": "<!-- @anchors\n  code: LOGINX\n-->\n# Login\n\nLOGINX-B01 — the password is checked.\n",
		"src/login.ts":      "export const login = () => true\n",
		"src/login.test.ts": "test('LOGINX-B01: checks the password', () => {})\n",
		"guides/CODE.md":    "# Code guide\n",
	} {
		writeProjectFile(t, root, p, c)
	}
	runCmd(t, newMapCmd(), "build", "--root", root)
	writeProjectFile(t, root, "reports/junit.xml", `<?xml version="1.0"?>
<testsuites><testsuite name="login" file="src/login.test.ts">
  <testcase name="LOGINX-B01: checks the password" file="src/login.test.ts"/>
</testsuite></testsuites>
`)
	runCmd(t, newIngestCmd(), "--root", root, "--junit", filepath.Join(root, "reports/junit.xml"))
	spec := node(t, root, "src/login.spec.md").Signal
	if spec == nil || strings.Join(spec.ProvenCodes, ",") != "LOGINX-B01" {
		t.Fatalf("a [6] project's passing case must prove its rule: %+v", spec)
	}
}
