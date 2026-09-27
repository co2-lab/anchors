package recode

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestRewrite_headerDialects(t *testing.T) {
	t.Run("RCRWR-B02: The header code is replaced in every header style", func(t *testing.T) {})
	cases := map[string]string{
		"md":      "<!-- @anchors\n  code: TCDTX\n-->\n",
		"ts":      "// @anchors\n//   ref: TCDTX\n",
		"feature": "# @anchors\n#   ref: TCDTX\n",
	}
	for name, in := range cases {
		got, n := Rewrite(in, "TCDTX", "TCTXX")
		if strings.Contains(got, "TCDTX") {
			t.Errorf("%s: TCDTX left over: %q", name, got)
		}
		if !strings.Contains(got, "TCTXX") || n == 0 {
			t.Errorf("%s: not replaced (n=%d): %q", name, n, got)
		}
	}
}

func TestRewrite_refListPreservesOthers(t *testing.T) {
	t.Run("RCRWR-B03: Only the old code changes in a ref list", func(t *testing.T) {})
	// a ref list: only TCDTX changes, MNDTX stays.
	in := "//   ref: TCDTX, MNDTX\n"
	got, _ := Rewrite(in, "TCDTX", "TCTXX")
	if !strings.Contains(got, "TCTXX") || !strings.Contains(got, "MNDTX") {
		t.Errorf("list not preserved: %q", got)
	}
	if strings.Contains(got, "TCDTX") {
		t.Errorf("TCDTX left over in the list: %q", got)
	}
}

func TestRewrite_scenarioCodeVariety(t *testing.T) {
	t.Run("RCRWR-B04: Scenario codes keep their suffixes", func(t *testing.T) {})
	in := "TCDTX-B01 TCDTX-S04 TCDTX-A04 TCDTX-FP01b TCDTX-DS-receipt-none TCDTX-VR TCDTX-RA03"
	got, n := Rewrite(in, "TCDTX", "TCTXX")
	if strings.Contains(got, "TCDTX") {
		t.Errorf("some TCDTX-* left over: %q", got)
	}
	for _, want := range []string{"TCTXX-B01", "TCTXX-FP01b", "TCTXX-DS-receipt-none", "TCTXX-VR", "TCTXX-RA03"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in: %q", want, got)
		}
	}
	if n != 7 {
		t.Errorf("want 7 replacements, got %d", n)
	}
}

func TestRewrite_bareCrossRef(t *testing.T) {
	t.Run("RCRWR-B05: A bare mention is replaced and its neighbour kept", func(t *testing.T) {})
	in := "Ver a regra em TCDTX (a tela de detalhe)."
	got, _ := Rewrite(in, "TCDTX", "TCTXX")
	if got != "Ver a regra em TCTXX (a tela de detalhe)." {
		t.Errorf("bare ref not replaced: %q", got)
	}
}

func TestRewrite_wordBoundary_doesNotTouchLongerCode(t *testing.T) {
	t.Run("RCRWR-X01: Longer codes and neighbours are never touched", func(t *testing.T) {})
	// The target must NOT match inside a longer code nor a neighbour that contains it.
	// The examples here were 4 chars and the move to 5 made them collide with the target —
	// `TCDTX-B01` became the target itself, so "must not touch" turned into "must".
	in := "TCDTXX-B01 e XTCDTX e TCDTXABCD"
	got, n := Rewrite(in, "TCDTX", "TCTXX")
	if got != in || n != 0 {
		t.Errorf("touched a longer code/neighbour (n=%d): %q", n, got)
	}
}

func TestRewrite_idempotentOnNewAbsent(t *testing.T) {
	t.Run("RCRWR-I01: A text without the old code is left unchanged", func(t *testing.T) {})
	in := "// @anchors\n//   ref: WXYZX\n\nWXYZ-B01"
	got, n := Rewrite(in, "TCDTX", "TCTXX")
	if got != in || n != 0 {
		t.Errorf("changed a text with no TCDTX (n=%d)", n)
	}
	once, _ := Rewrite("<!-- @anchors\n  code: TCDTX\n-->\nTCDTX-B01 and TCDTX.", "TCDTX", "TCTXX")
	if twice, n := Rewrite(once, "TCDTX", "TCTXX"); twice != once || n != 0 {
		t.Errorf("a second pass changed the rewritten text (n=%d): %q", n, twice)
	}
}

func TestFind_classifies(t *testing.T) {
	t.Run("RCRWR-B06: The dry run classifies each occurrence", func(t *testing.T) {})
	in := "// @anchors\n//   ref: TCDTX\n\nit('TCDTX-S02: x')\n// ver TCDTX adiante\n"
	occ := Find(in, "TCDTX")
	kinds := map[string]int{}
	for _, o := range occ {
		kinds[o.Kind]++
	}
	if kinds["scenario-code"] != 1 {
		t.Errorf("want 1 scenario-code, got %d (%v)", kinds["scenario-code"], occ)
	}
	if kinds["header"] != 1 {
		t.Errorf("want 1 header, got %d", kinds["header"])
	}
	if kinds["bare-ref"] != 1 {
		t.Errorf("want 1 bare-ref, got %d", kinds["bare-ref"])
	}
}

func TestValidCode(t *testing.T) {
	t.Run("RCRWR-B01: A code is well formed only in the project's lengths and alphabet", func(t *testing.T) {})
	// The valid length is the PROJECT's (`code_lengths`), not the engine's — so the test
	// declares what it exercises, as a project would. Without this the case depends on the
	// default and breaks when the default changes (it did: `MN01`, 4 chars, stopped being
	// valid when the canonical length became 5).
	config.SetCodeLengths([]int{4, 5})
	t.Cleanup(func() { config.SetCodeLengths([]int{5}) })
	for _, ok := range []string{"TCDTX", "MN01", "ABCDX", "ABCDE"} {
		if !ValidCode(ok) {
			t.Errorf("%s should be valid", ok)
		}
	}
	// Still invalid: outside the length range, lower case, and a character that is not an
	// upper-case alphanumeric.
	for _, bad := range []string{"TCD", "ABCDEF", "tcdt", "TC-T"} {
		if ValidCode(bad) {
			t.Errorf("%s should NOT be valid", bad)
		}
	}
}
