// @anchors
//   code: FSTFR
//   ref: FRSCF

package mapx

import (
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

func frozenCfg() *config.Config {
	return &config.Config{Version: 1, Derived: &config.Derived{
		Anchor: "spec",
		Files:  map[string]config.Padroes{"code": {"{{dir}}/{{name}}.ts"}},
	}}
}

func frozenFiles() []scan.File {
	return []scan.File{
		{Path: "pay/Pay.spec.md", Kind: "spec", HeaderCode: "PAYUN", Frozen: "QLDTQ-F01", Rev: "1"},
		{Path: "pay/Pay.ts", Kind: "code", Rev: "2"},
		{Path: "pay/Pay.test.ts", Kind: "test", HeaderCode: "PAYTS", HeaderRefs: []string{"PAYUN"}, Frozen: "flaky until the sandbox is back", Rev: "3"},
		{Path: "other/Cite.ts", Kind: "code", Codes: []string{"PAYUN-B01"}, Rev: "4"},
		{Path: "other/Cart.spec.md", Kind: "spec", HeaderCode: "CARTU", Rev: "5",
			FrozenRules: map[string]string{"CARTU-B02": "the coupon waits for legal"},
			Realizes:    []scan.Realizes{{From: "CARTU-B01", To: "QLDTQ-R01"}, {From: "CARTU-B02", To: "QLDTQ-R01"}, {From: "CARTU-B03", To: "QLDTQ-R02"}, {From: "CARTU-B04", To: "PLANM-R01"}}},
		{Path: "product/q.doctrine.md", Kind: "product", HeaderCode: "QLDTQ", Rev: "6",
			FrozenRules: map[string]string{"QLDTQ-R01": "QLDTQ-F02"},
			Freezes:     map[string]string{"QLDTQ-F01": "the plan mode ships next quarter", "QLDTQ-F02": "the flag stays off", "QLDTQ-F03": "nothing yet"}},
		{Path: "product/plan.doctrine.md", Kind: "product", HeaderCode: "PLANM", Frozen: "QLDTQ-F03", Rev: "7"},
	}
}

func frozenOf(g *Graph, id string) Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return Node{}
}

func TestFrozen_unitCoversItsFiles(t *testing.T) {
	t.Run("FRSCF-B01: A frozen unit covers its files, and nothing that only cites it", func(t *testing.T) {})
	g := Build(frozenFiles(), frozenCfg(), nil)
	for id, want := range map[string]string{
		"pay/Pay.spec.md": "QLDTQ-F01",
		"pay/Pay.ts":      "QLDTQ-F01",
		"pay/Pay.test.ts": "flaky until the sandbox is back",
		"other/Cite.ts":   "",
	} {
		if got := frozenOf(g, id).Frozen; got != want {
			t.Errorf("%s frozen %q, want %q", id, got, want)
		}
	}
}

func TestFrozen_doctrineCascades(t *testing.T) {
	t.Run("FRSCF-B02: A frozen doctrine freezes the rules that realize it", func(t *testing.T) {})
	g := Build(frozenFiles(), frozenCfg(), nil)
	want := map[string]string{"CARTU-B01": "QLDTQ-F02", "CARTU-B02": "the coupon waits for legal", "CARTU-B04": "QLDTQ-F03"}
	if got := frozenOf(g, "other/Cart.spec.md").FrozenRules; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFrozen_describe(t *testing.T) {
	t.Run("FRSCF-B03: A freeze reads with its reason", func(t *testing.T) {})
	g := Build(frozenFiles(), frozenCfg(), nil)
	for v, want := range map[string]string{
		"QLDTQ-F01":     "QLDTQ-F01 — the plan mode ships next quarter",
		"QLDTQ-F09":     "QLDTQ-F09",
		"a plain cause": "a plain cause",
	} {
		if got := g.DescribeFreeze(v); got != want {
			t.Errorf("%q reads %q, want %q", v, got, want)
		}
	}
	if r := g.FreezeRecords()["QLDTQ-F02"]; r.Doc != "product/q.doctrine.md" || r.Reason != "the flag stays off" {
		t.Errorf("record: %+v", r)
	}
}

func TestFrozen_parts(t *testing.T) {
	t.Run("FRSCF-B04: The frozen parts of the project, and of one freeze", func(t *testing.T) {})
	g := Build(frozenFiles(), frozenCfg(), nil)
	all := g.FrozenParts("")
	want := []FrozenPart{
		{Path: "other/Cart.spec.md", Rule: "CARTU-B01", Value: "QLDTQ-F02"},
		{Path: "other/Cart.spec.md", Rule: "CARTU-B02", Value: "the coupon waits for legal"},
		{Path: "other/Cart.spec.md", Rule: "CARTU-B04", Value: "QLDTQ-F03"},
		{Path: "pay/Pay.spec.md", Value: "QLDTQ-F01"},
		{Path: "pay/Pay.test.ts", Value: "flaky until the sandbox is back"},
		{Path: "pay/Pay.ts", Value: "QLDTQ-F01"},
		{Path: "product/plan.doctrine.md", Value: "QLDTQ-F03"},
		{Path: "product/q.doctrine.md", Rule: "QLDTQ-R01", Value: "QLDTQ-F02"},
	}
	if !reflect.DeepEqual(all, want) {
		t.Errorf("all:\n got %+v\nwant %+v", all, want)
	}
	if one := g.FrozenParts("QLDTQ-F03"); len(one) != 2 || one[0].Rule != "CARTU-B04" || one[1].Path != "product/plan.doctrine.md" {
		t.Errorf("one freeze: %+v", one)
	}
	if g.FrozenRules()["QLDTQ-R01"] != "QLDTQ-F02" || g.FrozenAt("./pay/Pay.ts") != "QLDTQ-F01" || g.FrozenAt("other/Cite.ts") != "" {
		t.Error("rules and paths")
	}
}

func TestFrozen_noMarkNoFreeze(t *testing.T) {
	t.Run("FRSCF-I01: No mark, no freeze", func(t *testing.T) {})
	g := Build(testFiles(), testCfg(), nil)
	if len(g.FrozenParts("")) != 0 || len(g.FrozenRules()) != 0 {
		t.Errorf("no mark, no freeze: %+v", g.FrozenParts(""))
	}
}

func TestFrozen_folderOfTheConfig(t *testing.T) {
	t.Run("FRSCF-B05: A frozen folder freezes its files and the units it owns", func(t *testing.T) {})
	cfg := frozenCfg()
	cfg.FrozenPaths = []config.FrozenPath{{Freeze: "IRDOC-F01", Paths: []string{"ir/**"}}}
	g := Build([]scan.File{
		{Path: "ir/Tax.spec.md", Kind: "spec", HeaderCode: "TAXIR", Rev: "1"},
		{Path: "ir/Tax.ts", Kind: "code", Rev: "2"},
		{Path: "ir/Old.ts", Kind: "code", Frozen: "kept for the audit", Rev: "3"},
		{Path: "flows/tax.yaml", Kind: "test", HeaderCode: "TAXFL", HeaderRefs: []string{"TAXIR"}, Rev: "4"},
		{Path: "pay/Pay.ts", Kind: "code", Rev: "5"},
	}, cfg, nil)
	for id, want := range map[string]string{
		"ir/Tax.spec.md": "IRDOC-F01",
		"ir/Tax.ts":      "IRDOC-F01",
		"ir/Old.ts":      "kept for the audit",
		"flows/tax.yaml": "IRDOC-F01",
		"pay/Pay.ts":     "",
	} {
		if got := frozenOf(g, id).Frozen; got != want {
			t.Errorf("%s frozen %q, want %q", id, got, want)
		}
	}
}
