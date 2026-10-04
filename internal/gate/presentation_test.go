// @anchors
//   code: PRTSB
//   ref: PRSNT

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

var presSpecNode = mapx.Node{Kind: mapx.KindSpec, Layer: "component"}

const presSpec = "## Properties\n\n| Property | Accepted values | Required | Default | Description |\n| --- | --- | --- | --- | --- |\n" +
	"| `status` | `pending`, `active`, `removed` | yes | - | the seat |\n| `isYou` | boolean | no | false | x |\n| `size` | `md` | no | md | one value is no enumeration |\n| `tone` | `'sm' \\| 'lg'` | no | sm | x |\n\n" +
	"## Data States\n\n### `role`\n\n| Value | Condition | What the interface shows |\n| --- | --- | --- |\n| admin | x | x |\n| member | x | x |\n\n" +
	"### `plan`\n\nprose under the heading is not a value\n\n| Value | Condition | What the interface shows |\n| --- | --- | --- |\n| free | x | x |\n| pro | x | x |\n\n" +
	"### `tier`\n\n| Value | Condition | What the interface shows |\n| --- | --- | --- |\n| gold | x | x |\n| silver | x | x |\n\n" +
	"## Test Identifiers\n\n| Identifier | Element | Used in | Action |\n| --- | --- | --- | --- |\n| `row-badge` | badge | x | x |\n\n" +
	"## Presentation validations\n\n| Rule | Prop/State | Condition | Appearance |\n| --- | --- | --- | --- |\n" +
	"| `OGMRX-P01` | `status` | == 'pending' | `row-badge` shows \"PENDENTE\" |\n" +
	"| `OGMRX-P02` | `status` | == 'active' | normal `row` |\n" +
	"| `OGMRX-P03` | `role` | otherwise | plain `row-badge` in the 'muted' tone |\n" +
	"| `OGMRX-P06` | `plan` | == 'free' | `row-badge` hidden |\n" +
	"| `OGMRX-P07` | `plan` | == 'pro' | `row-badge` gold |\n" +
	"| `OGMRX-P08` | `size` | always | `row-badge` medium |\n" +
	"| `OGMRX-P11` | `tier` | == 'gold' | a `halo` around the `row-badge` |\n" +
	"| `OGMRX-P12` | `tone` | == 'sm' | small `row-badge` |\n" +
	"| `OGMRX-P09` | `role` | otherwise | plain `row-badge` in the 'muted' tone |\n" +
	"| `OGMRX-P10` | `status` |\n" +
	"| `OGMRX-P04` | `status` | == 'active' | dimmed `row-badge` (ver OGMRX-M01 \"ativo\") |\n" +
	"| `OGMRX-P05` | TODO | TODO | TODO |\n"

func TestPresentationRows(t *testing.T) {
	t.Run("PRSNT-B01: The rows are read by position, and a spec with none is skipped", func(t *testing.T) {})
	rows := presentationRows(presSpec, nil, "component")
	if len(rows) != 10 || rows[1] != (presentationRow{"OGMRX-P02", "`status`", "== 'active'", "normal `row`"}) {
		t.Fatalf("the ten filled rows with their cells — not the TODO one, not the short one — got %+v", rows)
	}
	none := "## Rules\n\n### OGMRX-B01 — x\n"
	for _, v := range []Verdict{
		first(checkPresentationExhaustive(none, presSpecNode, "", nil, nil)),
		first(checkPresentationConflict(none, presSpecNode, "", nil, nil)),
		first(checkPresentationCopySingleSource(none, presSpecNode, "", nil, nil)),
		first(checkPresentationObservable(none, presSpecNode, "", nil, nil)),
		first(checkPresentationObservable(presSpec, mapx.Node{Kind: mapx.KindFeature}, "", nil, nil)),
	} {
		if v != Skip {
			t.Errorf("a spec with no rows, or a node that is not a spec, is skipped; got %v", v)
		}
	}
}

func first(v Verdict, _ string) Verdict { return v }

func TestPresentationExhaustive(t *testing.T) {
	t.Run("PRSNT-B02: Every declared value has an appearance", func(t *testing.T) {})
	v, msg := checkPresentationExhaustive(presSpec, presSpecNode, "", nil, nil)
	if v != Fail || !strings.Contains(msg, "`status`: removed") || !strings.Contains(msg, "`tier`: silver") || !strings.Contains(msg, "`tone`: lg") ||
		strings.Contains(msg, "role") || strings.Contains(msg, "isyou") || strings.Contains(msg, "plan") || strings.Contains(msg, "size") {
		t.Fatalf("only the status value no row names, got %v: %s", v, msg)
	}
	full := strings.Replace(presSpec, "| `OGMRX-P05` | TODO | TODO | TODO |\n", "| `OGMRX-P05` | `status` | == 'removed' | hidden `row-badge` |\n"+
		"| `OGMRX-P13` | `tier` | == 'silver' | plain `row-badge` |\n| `OGMRX-P14` | `tone` | == 'lg' | big `row-badge` |\n", 1)
	if v, msg := checkPresentationExhaustive(full, presSpecNode, "", nil, nil); v != Pass {
		t.Errorf("every value named passes, got %v: %s", v, msg)
	}
}

func TestPresentationConflict(t *testing.T) {
	t.Run("PRSNT-B03: One prop and one condition lead to one appearance", func(t *testing.T) {})
	v, msg := checkPresentationConflict(presSpec, presSpecNode, "", nil, nil)
	if v != Fail || !strings.Contains(msg, "OGMRX-P02 × OGMRX-P04") || strings.Contains(msg, "P09") {
		t.Fatalf("the two appearances of one condition are named, got %v: %s", v, msg)
	}
	one := strings.Replace(presSpec, "== 'active' | dimmed", "== 'inactive' | dimmed", 1)
	if v, msg := checkPresentationConflict(one, presSpecNode, "", nil, nil); v != Pass {
		t.Errorf("distinct conditions pass, got %v: %s", v, msg)
	}
}

func TestPresentationCopySingleSource(t *testing.T) {
	t.Run("PRSNT-B04: The text shown is a message code, not copy", func(t *testing.T) {})
	v, msg := checkPresentationCopySingleSource(presSpec, presSpecNode, "", nil, nil)
	if v != Fail || !strings.Contains(msg, `OGMRX-P01 → "PENDENTE"`) || strings.Contains(msg, "P04") || strings.Contains(msg, "P02") || strings.Contains(msg, "P03") {
		t.Fatalf("only the quoted copy with no message code, got %v: %s", v, msg)
	}
	clean := strings.Replace(presSpec, "shows \"PENDENTE\"", "shows `OGMRX-M02`", 1)
	if v, msg := checkPresentationCopySingleSource(clean, presSpecNode, "", nil, nil); v != Pass {
		t.Errorf("no copy passes, got %v: %s", v, msg)
	}
}

func TestPresentationObservable(t *testing.T) {
	t.Run("PRSNT-B05: What changes is something a test can point at", func(t *testing.T) {})
	v, msg := checkPresentationObservable(presSpec, presSpecNode, "", nil, nil)
	if v != Fail || !strings.Contains(msg, "OGMRX-P02") || strings.Contains(msg, "P01") || strings.Contains(msg, "P03") || strings.Contains(msg, "P11") {
		t.Fatalf("only the row citing no declared identifier — one cited second counts —, got %v: %s", v, msg)
	}
	seen := strings.Replace(presSpec, "normal `row` |", "normal `row-badge` |", 1)
	if v, msg := checkPresentationObservable(seen, presSpecNode, "", nil, nil); v != Pass {
		t.Errorf("every row observable passes, got %v: %s", v, msg)
	}
}
