package settings

import (
	"reflect"
	"strings"
	"testing"
)

// ONLY THE PO AND THE ARCHITECT decide the product. It is the rule that motivated roles: in a
// project with several contributors not everyone may decide for the product, and an agent that
// asks whoever runs it gets an answer that may not be the owner's.
func TestRole_whoDecidesTheProduct(t *testing.T) {
	t.Run("AGRLG-B02: Only the product owner and the architect decide the product", func(t *testing.T) {})
	decide := map[Role]bool{RolePO: true, RoleArchitect: true}

	for _, r := range KnownRoles() {
		got := r.Can(CapDecideProduct)
		if got != decide[r] {
			t.Errorf("%s: decides the product = %v, want %v", r, got, decide[r])
		}
	}
}

// The STRUCTURE is the architect's, not the PO's: "should this screen show the value?" is
// product; "may this layer import from that one?" is structure.
func TestRole_theStructureIsTheArchitects(t *testing.T) {
	t.Run("AGRLG-B03: The structure is the architect's", func(t *testing.T) {})
	if !RoleArchitect.Can(CapDecideStructure) {
		t.Error("the architect must decide the structure")
	}
	if RolePO.Can(CapDecideStructure) {
		t.Error("the PO does not decide the structure — it is another question")
	}
}

// QA AND REVIEWERS execute; they do not decide.
func TestRole_qaAndReviewersDoNotDecide(t *testing.T) {
	t.Run("AGRLG-B04: QA, reviewers and dev execute and do not decide", func(t *testing.T) {})
	for _, r := range []Role{RoleQA, RoleReviewer, RoleReviewerSec, RoleReviewerPerf, RoleDev} {
		if r.Can(CapDecideProduct) || r.Can(CapDecideStructure) {
			t.Errorf("%s must not decide anything — it executes", r)
		}
		if !r.Can(CapReview) {
			t.Errorf("%s must be able to review", r)
		}
	}
}

// EACH ROLE THAT REVIEWS HAS A LENS, and they are DISTINCT. A reviewer without a lens tends to
// do the review it KNOWS how to do — not the one that is missing.
func TestRole_lensesAreDistinct(t *testing.T) {
	t.Run("AGRLG-B05: Each reviewing role has its own lens", func(t *testing.T) {})
	withLens := []Role{RoleQA, RoleReviewer, RoleReviewerSec, RoleReviewerPerf}
	seen := map[string]Role{}

	for _, r := range withLens {
		lens := r.Lens()
		if strings.TrimSpace(lens) == "" {
			t.Errorf("%s reviews and declares no lens", r)
			continue
		}
		if other, repeated := seen[lens]; repeated {
			t.Errorf("%s and %s have the SAME lens — so one of them need not exist", r, other)
		}
		seen[lens] = r
	}

	// And the dev has NO lens: it reviews as part of delivering, without a declared focus.
	if RoleDev.Lens() != "" {
		t.Error("the dev must have no lens — the reviewer roles have it")
	}
}

// AN UNKNOWN ROLE CAN DO NOTHING. A hand-written `settings.yaml` with an invented role unlocks
// no capability — the closed default is the same as the boolean this mechanism replaced.
func TestRole_unknownCanDoNothing(t *testing.T) {
	t.Run("AGRLG-X01: An unknown role can do nothing", func(t *testing.T) {})
	for _, r := range []Role{"", "tech-lead", "gerente", "DEV"} {
		for _, c := range []Capability{
			CapDecideProduct, CapDecideStructure, CapWritePlan,
			CapWriteSpec, CapImplement, CapProve, CapReview,
		} {
			if Role(r).Can(c) {
				t.Errorf("role %q unlocked %q", r, c)
			}
		}
	}
}

// The name is read in both languages and in abbreviations: the question is asked in the
// terminal, and whoever answers types `po`, not `product-owner`.
func TestParseRole(t *testing.T) {
	t.Run("AGRLG-B06: Typed roles are recognised by name and abbreviation", func(t *testing.T) {})
	cases := map[string]Role{
		"dev": RoleDev, "developer": RoleDev, "desenvolvedor": RoleDev,
		"qa": RoleQA, "tester": RoleQA,
		"po": RolePO, "product-owner": RolePO, "produto": RolePO,
		"arch": RoleArchitect, "arquiteto": RoleArchitect,
		"sec": RoleReviewerSec, "seguranca": RoleReviewerSec, "segurança": RoleReviewerSec,
		"perf": RoleReviewerPerf, "performance": RoleReviewerPerf,
		"reviewer": RoleReviewer, "revisor": RoleReviewer,
		" DEV \n": RoleDev,
	}
	for in, want := range cases {
		if got := ParseRole(in); got != want {
			t.Errorf("ParseRole(%q) = %q, want %q", in, got, want)
		}
	}
	// What is not recognised returns empty — and the caller asks again, instead of assuming a
	// role. Assuming here would give capabilities to whoever did not ask for them.
	for _, invalid := range []string{"", "tech-lead", "gerente", "sim"} {
		if got := ParseRole(invalid); got != "" {
			t.Errorf("ParseRole(%q) = %q, want empty", invalid, got)
		}
	}
}

// Every role declares a TITLE and what it DOES — what the command shows when choosing.
func TestKnownRoles_allPresentThemselves(t *testing.T) {
	t.Run("AGRLG-B01: Every known role presents itself", func(t *testing.T) {})
	rs := KnownRoles()
	want := []Role{RoleArchitect, RoleDev, RolePO, RoleQA, RoleReviewer, RoleReviewerPerf, RoleReviewerSec}
	if !reflect.DeepEqual(rs, want) || !reflect.DeepEqual(KnownRoles(), rs) {
		t.Fatalf("roles = %v, want the seven known roles in a stable order %v", rs, want)
	}
	for _, r := range rs {
		if strings.TrimSpace(r.Title()) == "" {
			t.Errorf("%s without a title", r)
		}
		if strings.TrimSpace(r.Does()) == "" {
			t.Errorf("%s does not say what it does", r)
		}
		if len(r.Caps()) == 0 {
			t.Errorf("%s has no capability — so it does nothing", r)
		}
	}
}
