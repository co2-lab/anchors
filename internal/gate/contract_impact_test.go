// @anchors
//   code: CITCN
//   ref: CTRIM

package gate

import (
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/shell"
)

const paySpec = "# Pay\n\n## Data contract\n\n| Field | Origin | Format |\n| --- | --- | --- |\n" +
	"| `amount` | DEP1 | cents |\n| `currency` | DEP1 | ISO |\n| `fee` | fixed | cents |\n| `note` | user | text |\n\n" +
	"## States\n\n| State | Entered when |\n| --- | --- |\n| `PAYMT-S01` | the charge is sent |\n\n" +
	"## Validations\n\n| Rule | Field | Condition | Behavior |\n| --- | --- | --- | --- |\n| `PAYMT-V01` | `amount` | > 0 | 400 |\n\n" +
	"## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `PAYMT-B01` | `amount` |\n| `PAYMT-B02` | `currency` |\n| `PAYMT-B03` | `fee` |\n| `PAYMT-B04` | `PAYMT-S01` |\n"

const checkoutSpec = "# Checkout\n\n## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `CHKOT-B01` | `amount.total` |\n| `CHKOT-B02` | `status` |\n"

// impactRepo commits the two specs, the code and a test, and returns the root, the map and
// the configuration.
func impactRepo(t *testing.T) (string, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", append([]string{"-C", root}, args...)...)
		c.Env = append(shell.WithoutRepoEnv(os.Environ()), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := c.CombinedOutput(); err != nil {
			t.Skipf("git: %v %s", err, out)
		}
	}
	git("init", "-q")
	writeFile(t, root, "src/pay.spec.md", paySpec)
	writeFile(t, root, "src/checkout.spec.md", checkoutSpec)
	writeFile(t, root, "src/pay.go", "package pay\n")
	writeFile(t, root, "src/pay_test.go", "package pay\nfunc TestPay(t *testing.T) {\n\tt.Run(\"PAYMT-B01: charges the amount\", func(t *testing.T) {})\n}\n")
	writeFile(t, root, "src/pay2_test.go", "package pay\nfunc TestOther(t *testing.T) {\n\tt.Run(\"PAYMT-B010: another rule\", func(t *testing.T) {})\n}\n")
	git("add", ".")
	git("commit", "-q", "-m", "init")
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: "src/pay.spec.md", Kind: mapx.KindSpec}, {ID: "src/checkout.spec.md", Kind: mapx.KindSpec},
			{ID: "src/pay.go", Kind: mapx.KindCode}, {ID: "src/pay_test.go", Kind: mapx.KindTest}, {ID: "src/pay2_test.go", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{{From: "src/pay.spec.md", To: "src/pay.go", Type: mapx.EdgeSpecifies},
			{From: "src/checkout.spec.md", To: "src/pay.go", Type: mapx.EdgeDependsOn}},
	}
	return root, g, &config.Config{Dialect: &config.Dialect{Family: "go"}}
}

func edited(root string, t *testing.T) string {
	t.Helper()
	s := strings.Replace(paySpec, "| `amount` | DEP1 | cents |", "| `amount` | DEP1 | decimal |", 1)
	s = strings.Replace(s, "| `fee` | fixed | cents |\n", "| `tip` | user | cents |\n", 1)
	s = strings.Replace(s, "| `note` | user | text |", "| `note` | user | markdown |", 1) // changed, and no rule uses it
	s = strings.Replace(s, "| `PAYMT-S01` | the charge is sent |", "| `PAYMT-S01` | the charge is confirmed |", 1)
	writeFile(t, root, "src/pay.spec.md", s)
	return s
}

func TestContractImpact_changedFields(t *testing.T) {
	t.Run("CTRIM-B01: Changed and removed fields change, added ones do not", func(t *testing.T) {})
	root, g, cfg := impactRepo(t)
	var fields []string
	for _, imp := range ContractImpacts(root, "src/pay.spec.md", edited(root, t), g, cfg) {
		fields = append(fields, imp.Field)
	}
	if !reflect.DeepEqual(fields, []string{"amount", "fee", "paymt-s01"}) {
		t.Fatalf("amount edited, fee removed and the state a rule uses edited; not currency, tip nor the unused note; got %v", fields)
	}
}

func TestContractImpact_rulesAndTests(t *testing.T) {
	t.Run("CTRIM-B02: A changed field names its rules, here and in the dependents, and their tests", func(t *testing.T) {})
	root, g, cfg := impactRepo(t)
	imps := ContractImpacts(root, "src/pay.spec.md", edited(root, t), g, cfg)
	if len(imps) == 0 || imps[0].Field != "amount" {
		t.Fatalf("amount is impacted, got %+v", imps)
	}
	if !reflect.DeepEqual(imps[0].Rules, []string{"CHKOT-B01", "PAYMT-B01", "PAYMT-V01"}) || !reflect.DeepEqual(imps[0].Tests, []string{"src/pay_test.go"}) {
		t.Errorf("the rules here and in the dependent, and the test citing one; got %+v", imps[0])
	}
}

func TestContractImpact_gate(t *testing.T) {
	t.Run("CTRIM-B03: The gate is divergence with the impact, and passes without", func(t *testing.T) {})
	root, g, cfg := impactRepo(t)
	n := mapx.Node{ID: "src/pay.spec.md", Kind: mapx.KindSpec}
	if v, msg := checkContractImpact(paySpec, n, root, g, cfg); v != Pass {
		t.Errorf("unchanged passes, got %v: %s", v, msg)
	}
	v, msg := checkContractImpact(edited(root, t), n, root, g, cfg)
	if v != Diverge || !strings.Contains(msg, "`amount`") || !strings.Contains(msg, "CHKOT-B01") || !strings.Contains(msg, "src/pay_test.go") {
		t.Fatalf("pending naming the field, its rules and tests; got %v: %s", v, msg)
	}
	if v, _ := checkContractImpact(paySpec, mapx.Node{Kind: mapx.KindFeature}, root, g, cfg); v != Skip {
		t.Errorf("a feature is skipped, got %v", v)
	}
	if imps := ContractImpacts(t.TempDir(), "src/pay.spec.md", paySpec, g, cfg); imps != nil {
		t.Errorf("no repository is no impact, got %+v", imps)
	}
}

func TestContractImpact_impactedTests(t *testing.T) {
	t.Run("CTRIM-B04: The impacted tests are listed for the selection", func(t *testing.T) {})
	root, g, cfg := impactRepo(t)
	if got := ImpactedTests(root, g, cfg); len(got) != 0 {
		t.Fatalf("nothing changed, nothing impacted; got %v", got)
	}
	edited(root, t)
	if got := ImpactedTests(root, g, cfg); !reflect.DeepEqual(got, []string{"src/pay_test.go"}) {
		t.Errorf("the test of the affected rule is listed, got %v", got)
	}
}

func TestContractImpact_answeredByTheChangesRevision(t *testing.T) {
	t.Run("CTRIM-B05: The rules the change's own revision revises or checks answer the impact", func(t *testing.T) {})
	root, g, cfg := impactRepo(t)
	s := edited(root, t)
	all := ContractImpacts(root, "src/pay.spec.md", s, g, cfg)
	if len(all) == 0 {
		t.Fatal("the edit has impacts to answer")
	}
	var rules []string
	for _, imp := range all {
		for _, r := range imp.Rules {
			rules = append(rules, "`"+r+"`")
		}
	}
	answered := s + "\n### PAYMT-R0001 — the amount moved to decimal\n\n**Revises:** " + rules[0] + ".\n**Checked:** " + strings.Join(rules, ", ") + ".\n"
	writeFile(t, root, "src/pay.spec.md", answered)
	if got := ContractImpacts(root, "src/pay.spec.md", answered, g, cfg); len(got) != 0 {
		t.Errorf("every rule the change's revision names is answered: %+v", got)
	}
	partly := s + "\n### PAYMT-R0001 — only one\n\n**Revises:** " + rules[0] + ".\n"
	if len(rules) > 1 {
		if got := ContractImpacts(root, "src/pay.spec.md", partly, g, cfg); len(got) == 0 {
			t.Error("an impact with a rule nobody answered stays")
		}
	}
}

func TestContractImpact_answeredWhereTheRuleLives(t *testing.T) {
	t.Run("CTRIM-B06: A rule is answered by a revision added to its own spec — even one not yet in git —, or by the changed spec's naming its full code; a short code answers only its own spec's rule", func(t *testing.T) {})
	root, g, cfg := impactRepo(t)
	s := edited(root, t)
	// A third spec reads the amount too, and is new: not yet in git.
	cart := "# Cart\n\n## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `CARTX-B01` | `amount` |\n"
	writeFile(t, root, "src/cart.spec.md", cart)
	g.Nodes = append(g.Nodes, mapx.Node{ID: "src/cart.spec.md", Kind: mapx.KindSpec})
	g.Edges = append(g.Edges, mapx.Edge{From: "src/cart.spec.md", To: "src/pay.go", Type: mapx.EdgeDependsOn})
	openOn := func(spec string) []string {
		for _, imp := range ContractImpacts(root, "src/pay.spec.md", spec, g, cfg) {
			if imp.Field == "amount" {
				return imp.Rules
			}
		}
		return nil
	}
	if got := openOn(s); strings.Join(got, ",") != "CARTX-B01,CHKOT-B01,PAYMT-B01,PAYMT-V01" {
		t.Fatalf("the amount is read by rules in three specs: %v", got)
	}
	// The changed spec answers its own B01 by the short code: that says nothing of the
	// checkout's B01 nor the cart's.
	own := s + "\n### PAYMT-R0001 — decimal\n\n**Revises:** `B01`, `V01`.\n"
	writeFile(t, root, "src/pay.spec.md", own)
	if got := openOn(own); got == nil {
		t.Error("a short code answers only its own spec's rule: the checkout's and the cart's stay open")
	}
	// Each dependent answers its rule in its own revision; the cart is not in git, so all of
	// its revisions are new.
	writeFile(t, root, "src/checkout.spec.md", checkoutSpec+"\n### CHKOT-R0001 — reads the decimal\n\n**Checked:** `B01`.\n")
	writeFile(t, root, "src/cart.spec.md", cart+"\n### CARTX-R0001 — reads the decimal\n\n**Revises:** `CARTX-B01`.\n")
	if got := openOn(own); got != nil {
		t.Errorf("each rule answered where it lives: %v", got)
	}
	// Or the changed spec names the other spec's rule by its full code.
	writeFile(t, root, "src/checkout.spec.md", checkoutSpec)
	full := own + "**Checked:** `CHKOT-B01`.\n"
	writeFile(t, root, "src/pay.spec.md", full)
	if got := openOn(full); got != nil {
		t.Errorf("a full code in the changed spec's revision answers another spec's rule: %v", got)
	}
}
