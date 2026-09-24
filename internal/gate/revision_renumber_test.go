package gate

import (
	"strings"
	"testing"
)

const renumberHead = "<!-- @anchors\n  code: PRICX\n-->\n# Pricing\n\n"

func rev(n int, what string) string {
	return "> **PRICX-R000" + string(rune('0'+n)) + ":** " + what + "\n"
}

func renames(plan []Renumbering) map[string]string {
	m := map[string]string{}
	for _, r := range plan {
		m[r.Old] = r.New
	}
	return m
}

func TestPlanRenumber(t *testing.T) {
	t.Run("RVRNR-B01: A revision the branch added whose number the base already uses moves to the next free number", func(t *testing.T) {
		mb := renumberHead + rev(1, "first") + rev(2, "second")
		plan, err := PlanRenumber([]RevisionFile{{
			Path:      "Pricing.spec.md",
			MergeBase: mb,
			Base:      mb + rev(3, "the other PR's change"),
			Working:   mb + rev(3, "this branch's change"),
		}})
		if err != nil {
			t.Fatal(err)
		}
		if got := renames(plan); len(got) != 1 || got["PRICX-R0003"] != "PRICX-R0004" {
			t.Fatalf("expected R0003 → R0004, got %v", plan)
		}
	})

	t.Run("RVRNR-B02: After a rebase, only the branch's revision moves, not the base's with the same number", func(t *testing.T) {
		base := renumberHead + rev(1, "first") + rev(2, "second") + rev(3, "the other PR's change")
		plan, err := PlanRenumber([]RevisionFile{{
			Path:      "Pricing.spec.md",
			MergeBase: base, // after the rebase the branch forks from the base itself
			Base:      base,
			Working:   base + rev(3, "this branch's change"),
		}})
		if err != nil {
			t.Fatal(err)
		}
		if len(plan) != 1 || plan[0].Old != "PRICX-R0003" || plan[0].New != "PRICX-R0004" {
			t.Fatalf("expected one R0003 → R0004, got %v", plan)
		}
		out, n := RewriteRevisionCitations(base+rev(3, "this branch's change"), base, renames(plan))
		if n != 1 {
			t.Fatalf("expected one citation rewritten, got %d", n)
		}
		if !strings.Contains(out, rev(3, "the other PR's change")) || !strings.Contains(out, rev(4, "this branch's change")) {
			t.Fatalf("the base's R0003 must stay and the branch's must become R0004:\n%s", out)
		}
	})

	t.Run("RVRNR-B03: When one added revision collides, every added revision of that code moves in order", func(t *testing.T) {
		mb := renumberHead + rev(1, "first") + rev(2, "second")
		plan, err := PlanRenumber([]RevisionFile{{
			Path:      "Pricing.spec.md",
			MergeBase: mb,
			Base:      mb + rev(3, "the other PR's change"),
			Working:   mb + rev(3, "this branch, first") + rev(4, "this branch, second"),
		}})
		if err != nil {
			t.Fatal(err)
		}
		got := renames(plan)
		if len(got) != 2 || got["PRICX-R0003"] != "PRICX-R0004" || got["PRICX-R0004"] != "PRICX-R0005" {
			t.Fatalf("expected R0003 → R0004 and R0004 → R0005, got %v", plan)
		}
	})

	t.Run("RVRNR-B04: A revision the branch added whose number the base does not use stays", func(t *testing.T) {
		mb := renumberHead + rev(1, "first")
		plan, err := PlanRenumber([]RevisionFile{{
			Path:      "Pricing.spec.md",
			MergeBase: mb,
			Base:      mb,
			Working:   mb + rev(2, "this branch's change"),
		}})
		if err != nil {
			t.Fatal(err)
		}
		if len(plan) != 0 {
			t.Fatalf("nothing collides, nothing should move: %v", plan)
		}
	})

	t.Run("RVRNR-B06: The branch adding the same number twice is refused, naming the code", func(t *testing.T) {
		mb := renumberHead + rev(1, "first")
		_, err := PlanRenumber([]RevisionFile{{
			Path:      "Pricing.spec.md",
			MergeBase: mb,
			Base:      mb + rev(2, "base"),
			Working:   mb + rev(2, "one") + rev(2, "two"),
		}})
		if err == nil || !strings.Contains(err.Error(), "PRICX-R0002") {
			t.Fatalf("expected a refusal naming PRICX-R0002, got %v", err)
		}
	})

	t.Run("RVRNR-I01: A revision the base has is never renumbered, even with its explanation edited", func(t *testing.T) {
		mb := renumberHead + rev(1, "first") + rev(2, "second")
		plan, err := PlanRenumber([]RevisionFile{{
			Path:      "Pricing.spec.md",
			MergeBase: mb,
			Base:      mb,
			Working:   renumberHead + rev(1, "first") + rev(2, "second, with a typo fixed"),
		}})
		if err != nil {
			t.Fatal(err)
		}
		if len(plan) != 0 {
			t.Fatalf("an edited base revision is not an added one: %v", plan)
		}
	})
}

func TestRewriteRevisionCitations(t *testing.T) {
	t.Run("RVRNR-B05: A citation is rewritten only on a line the branch added", func(t *testing.T) {
		mb := "// see PRICX-R0003 of the other PR\nfunc a() {}\n"
		working := mb + "// PRICX-R0003: the carry-forward\nfunc b() {}\n"
		out, n := RewriteRevisionCitations(working, mb, map[string]string{"PRICX-R0003": "PRICX-R0004"})
		if n != 1 {
			t.Fatalf("expected one rewrite, got %d:\n%s", n, out)
		}
		if !strings.Contains(out, "// see PRICX-R0003 of the other PR") {
			t.Errorf("a merge-base line must keep its code:\n%s", out)
		}
		if !strings.Contains(out, "// PRICX-R0004: the carry-forward") {
			t.Errorf("the branch's line must be rewritten:\n%s", out)
		}
	})

	t.Run("RVRNR-I02: Rewriting is one pass, and a chain of renames never cascades", func(t *testing.T) {
		working := "@PRICX-R0003 and @PRICX-R0004\n"
		out, n := RewriteRevisionCitations(working, "", map[string]string{
			"PRICX-R0003": "PRICX-R0004", "PRICX-R0004": "PRICX-R0005"})
		if out != "@PRICX-R0004 and @PRICX-R0005\n" || n != 2 {
			t.Fatalf("expected each code moved once, got %q (%d)", out, n)
		}
	})

	t.Run("RVRNR-X01: Does not rewrite a revision cited without its unit code", func(t *testing.T) {
		working := "the R0003 fixed it; PRICX-R00031 is not a revision code either\n"
		out, n := RewriteRevisionCitations(working, "", map[string]string{"PRICX-R0003": "PRICX-R0004"})
		if out != working || n != 0 {
			t.Fatalf("a short citation must be left alone, got %q (%d)", out, n)
		}
	})
}
