// @anchors
//   code: PTTSP
//   ref: CFPCD

package code

import "testing"

func TestGenerateFromPath_normalNoCollision(t *testing.T) {
	t.Run("CFPCD-B04: A normal file name gets the generated code when free", func(t *testing.T) {})
	// a normal basename, nothing taken → the generated code, unchanged.
	got := GenerateFromPath("packages/backend/models/metadata.spec.md", "metadata", false, map[string]bool{})
	want := Generate("metadata")
	if got != want {
		t.Fatalf("without a collision it should be the generated %q, got %q", want, got)
	}
}

func TestGenerateFromPath_collisionStaysUnique(t *testing.T) {
	t.Run("CFPCD-I01: Same-named units get distinct, deterministic codes", func(t *testing.T) {})
	// models/metadata takes the generated code; repositories/metadata COLLIDES → gets a UNIQUE
	// variation (blind tie-break, deterministic given the map). The SYMMETRIC tie-break (both
	// prefixed) is a separate recode feature — not here.
	taken := map[string]bool{}
	first := GenerateFromPath("packages/backend/models/metadata.spec.md", "metadata", false, taken)
	taken[first] = true

	second := GenerateFromPath("packages/backend/repositories/metadata.spec.md", "metadata", false, taken)
	if second == first {
		t.Fatalf("collision not resolved: both %q", second)
	}
	if taken[second] {
		t.Errorf("the tie-break code %q is not unique", second)
	}
	// determinism given the map: same input + same taken → same result.
	again := GenerateFromPath("packages/backend/repositories/metadata.spec.md", "metadata", false, taken)
	if again != second {
		t.Errorf("not deterministic given the map: %q then %q", second, again)
	}
}

func TestGenerateFromPath_genericBasenameUsesParent(t *testing.T) {
	t.Run("CFPCD-B02: A generic file name takes its code from the parent folder", func(t *testing.T) {})
	// functions/manage-metadata/handler.ts → the unit is "manage-metadata".
	got := GenerateFromPath("packages/backend/functions/manage-metadata/handler.ts", "handler", true, map[string]bool{})
	want := GenerateUnique("manage-metadata", map[string]bool{})
	if got != want {
		t.Fatalf("a generic basename should use the parent folder (%q), got %q", want, got)
	}
	// and it must NOT be the code of "handler"
	if got == Generate("handler") {
		t.Errorf("the generic code of 'handler' leaked: %q", got)
	}
	// with no parent folder, the file's own name is all there is.
	if got := GenerateFromPath("handler.ts", "handler", true, map[string]bool{}); got != Generate("handler") {
		t.Errorf("a root handler.ts → %q, want the code of handler %q", got, Generate("handler"))
	}
}

func TestGenerateFromPath_artifactSuffixesShareTheBase(t *testing.T) {
	t.Run("CFPCD-B03: The code comes from the unit's name the caller gives", func(t *testing.T) {})
	want := GenerateFromPath("a/Login.tsx", "Login", false, map[string]bool{})
	for _, p := range []string{"a/Login.spec.md", "a/Login.feature", "a/Login.test.tsx"} {
		if got := GenerateFromPath(p, "Login", false, map[string]bool{}); got != want {
			t.Errorf("%s → %q, want the code of Login %q", p, got, want)
		}
	}
}

func TestGenerateFromPath_genericIsTheCallers(t *testing.T) {
	t.Run("CFPCD-B01: Whether a name is generic is the caller's to say", func(t *testing.T) {})
	if got := GenerateFromPath("functions/manage-metadata/handler.ts", "handler", false, map[string]bool{}); got != Generate("handler") {
		t.Errorf("not generic for this project: the file's own name, got %q", got)
	}
}
