package governance

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
)

// A pack on disk (the authority groups the report) plus one inline obligation (grouped
// as "declared in the project"). Three models carry personal data: one is erased by the
// handler, one declares the debt, one is simply forgotten.
const complianceYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "**/*.go"
packs: ["./packs/local-privacy.yaml"]
pack_values:
  erasure_handler: "handlers/erase.go"
obligations:
  - name: audit-logged
    when: "audited: yes"
    must_appear_in: ["handlers/log.go"]
  - name: retention
    when: "retained: yes"
    must_appear_in: ["handlers/log.go"]
`

const localPack = `name: local-privacy
domain: privacy
authority: "Local Privacy Act"
requires:
  - erasure_handler
obligations:
  - name: erasure
    article: "Art. 9"
    when: "carries: personal-data"
    must_appear_in: ["{{erasure_handler}}"]
`

func complianceFiles() map[string]string {
	return map[string]string{
		"packs/local-privacy.yaml": localPack,
		"handlers/erase.go":        "package handlers\n// deletes user\n",
		"handlers/log.go":          "package handlers\n",
		"models/user.go":           "// carries: personal-data\npackage models\n",
		"models/order.go":          "// carries: personal-data\n// obligation_pending: erasure — next sprint\npackage models\n",
		"models/profile.go":        "// carries: personal-data\npackage models\n",
		"models/event.go":          "// audited: yes\npackage models\n",
		"models/archive.go":        "// retained: yes\n// obligation_pending: retention — Q3 cleanup\npackage models\n",
	}
}

func complianceGraph() *mapx.Graph {
	var nodes []mapx.Node
	for _, id := range []string{"handlers/erase.go", "handlers/log.go", "models/user.go", "models/order.go", "models/profile.go", "models/event.go", "models/archive.go"} {
		nodes = append(nodes, mapx.Node{ID: id, Kind: mapx.KindCode, Layer: "code"})
	}
	return &mapx.Graph{Nodes: nodes}
}

func TestComplianceReportsEachDutyByNorm(t *testing.T) {
	dir := govProject(t, complianceYAML, complianceFiles(), complianceGraph())

	out, err := runCmd(t, newComplianceCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	// the pack's duty, grouped under its authority, cited with its article:
	// 3 subject (user, order, profile), 1 complying (user), 1 debt (order)
	if !strings.Contains(out, "── Local Privacy Act\n") {
		t.Errorf("a pack duty is grouped under the pack's authority:\n%s", out)
	}
	if !strings.Contains(out, "   ✗ erasure") || !strings.Contains(out, "Art. 9") ||
		!strings.Contains(out, "  3 subject,   1 complying, 1 assumed debt(s)") {
		t.Errorf("erasure: 3 subject, 1 complying, 1 debt, cited by article:\n%s", out)
	}
	// the inline duty: its only subject is missing, and nobody declared debt — the
	// report warns it is probably a disconnected target, not total violation
	if !strings.Contains(out, "── declared in the project\n") {
		t.Errorf("an inline obligation is grouped as declared in the project:\n%s", out)
	}
	if !strings.Contains(out, "⚠ NONE complies — check whether `handlers/log.go`") {
		t.Errorf("zero complying with subjects must point at the target path:\n%s", out)
	}
	// a duty whose only subject ASSUMED the debt is known and scheduled, not disconnected
	if strings.Count(out, "⚠ NONE complies") != 1 {
		t.Errorf("only audit-logged warns; retention carries declared debt:\n%s", out)
	}
	if !strings.Contains(out, "total: 3 subject duty(ies) across 5 node(s), 1 fulfilled") {
		t.Errorf("the totals add both norms:\n%s", out)
	}
	if !strings.Contains(out, "(use --verbose") {
		t.Errorf("with gaps and no --verbose, the hint is printed:\n%s", out)
	}
	if strings.Contains(out, "✗ models/profile.go") {
		t.Errorf("without --verbose the missing nodes are not listed:\n%s", out)
	}
	// packs on disk not adopted are listed; local-privacy is not an embedded pack
	if !strings.Contains(out, "packs available and not adopted:") || !strings.Contains(out, "privacy/lgpd") {
		t.Errorf("the embedded packs the project did not adopt should be listed:\n%s", out)
	}
}

func TestComplianceVerboseListsTheMissingNodes(t *testing.T) {
	dir := govProject(t, complianceYAML, complianceFiles(), complianceGraph())

	out, err := runCmd(t, newComplianceCmd(), "--root", dir, "--verbose")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "       ✗ models/profile.go\n") {
		t.Errorf("--verbose lists the node that does not comply:\n%s", out)
	}
	if strings.Contains(out, "✗ models/user.go") || strings.Contains(out, "✗ models/order.go") {
		t.Errorf("the complying node and the one with declared debt are not missing:\n%s", out)
	}
	if strings.Contains(out, "(use --verbose") {
		t.Errorf("the hint is for the non-verbose run only:\n%s", out)
	}
}

func TestComplianceWithoutDuties(t *testing.T) {
	dir := govProject(t, "version: 2\nlayers: {}\n", nil, &mapx.Graph{})

	out, err := runCmd(t, newComplianceCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "No duty declared.\n") {
		t.Errorf("a project without packs or obligations says so:\n%s", out)
	}
	if !strings.Contains(out, "packs available and not adopted:") {
		t.Errorf("and shows what it could adopt:\n%s", out)
	}
}

func TestComplianceFailsOnAPackMissingItsValues(t *testing.T) {
	yaml := strings.Replace(complianceYAML, "pack_values:\n  erasure_handler: \"handlers/erase.go\"\n", "", 1)
	dir := govProject(t, yaml, complianceFiles(), complianceGraph())

	_, err := runCmd(t, newComplianceCmd(), "--root", dir)
	if err == nil || !strings.Contains(err.Error(), "erasure_handler") {
		t.Errorf("a pack whose required value is missing must fail naming it; got %v", err)
	}
}

// printAvailable separates "does not apply to me" from "I forgot": an adopted pack,
// written either as a short name or as its path under ./packs/, is not offered again.
func TestPrintAvailableSkipsAdoptedPacks(t *testing.T) {
	var all []string
	for _, names := range initx.AvailablePacks() {
		all = append(all, names...)
	}
	if len(all) < 2 {
		t.Fatalf("the test needs at least two embedded packs, got %v", all)
	}

	adopted := []string{"./packs/" + all[0] + ".yaml"}
	out := captureStdout(t, func() { printAvailable(&config.Config{Packs: adopted}) })
	offered := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		if rest, ok := strings.CutPrefix(line, "packs available and not adopted: "); ok {
			for _, n := range strings.Split(rest, ", ") {
				offered[n] = true
			}
		}
	}
	if offered[all[0]] {
		t.Errorf("the adopted pack %q must not be offered:\n%s", all[0], out)
	}
	if len(offered) != len(all)-1 {
		t.Errorf("every other pack must be offered (%d), got %d:\n%s", len(all)-1, len(offered), out)
	}

	out = captureStdout(t, func() { printAvailable(&config.Config{Packs: all}) })
	if out != "" {
		t.Errorf("with every pack adopted there is nothing to offer; got:\n%s", out)
	}
}
