package mapcmd

import (
	"strings"
	"testing"
)

// A change to the spec propagates DOWN to the code and the test it specifies, and is
// governed by nobody; a change to the code propagates nowhere and is confronted UP with
// the spec and the guide.
func TestImpact_bothDirections(t *testing.T) {
	root := fixtureProject(t)

	spec := runCmd(t, newImpactCmd(), "src/login.spec.md", "--root", root)
	if !strings.Contains(spec, "impact of: src/login.spec.md") {
		t.Errorf("the origin is not named:\n%s", spec)
	}
	down := spec[strings.Index(spec, "↓"):strings.Index(spec, "↑")]
	for _, id := range []string{"src/login.ts", "src/login.test.ts"} {
		if !strings.Contains(down, id) {
			t.Errorf("%s should be redone when the spec changes:\n%s", id, spec)
		}
	}
	if !strings.Contains(spec, "(nothing — it is not governed by anyone)") {
		t.Errorf("the spec is governed by nobody here:\n%s", spec)
	}

	code := runCmd(t, newImpactCmd(), "src/login.ts", "--root", root)
	if !strings.Contains(code, "(nothing — no child depends on it)") {
		t.Errorf("nothing depends on the code:\n%s", code)
	}
	up := code[strings.Index(code, "↑"):]
	for _, id := range []string{"src/login.spec.md", "guides/CODE.md"} {
		if !strings.Contains(up, id) {
			t.Errorf("the code must be confronted with %s:\n%s", id, code)
		}
	}
}

func TestImpact_errors(t *testing.T) {
	root := fixtureProject(t)
	if _, err := runCmdErr(newImpactCmd(), t, "src/ghost.ts", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "is not in the map") {
		t.Errorf("a file outside the map: got %v", err)
	}
	if _, err := runCmdErr(newImpactCmd(), t, "src/login.ts", "--root", t.TempDir()); err == nil ||
		!strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: got %v", err)
	}
}
