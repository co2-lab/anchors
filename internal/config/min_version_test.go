package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// `min_version` exists for the moment an Anchors fix must reach everyone before work goes
// on. The previous check compared the binary with the map's `gerado_por` — a DERIVED field,
// rewritten by any agent's next `map build`. Measured: it went back to "dev" in a project
// whose current release was v0.1.83.
//
// Here the declaration is a person's, in the file the team versions.

func TestAtendeMinVersion_orderIsOrdinalNotEquality(t *testing.T) {
	t.Run("MNVRM-B01: Versions are ordered part by part as numbers", func(t *testing.T) {})
	cases := []struct {
		minimum, running string
		satisfies        bool
	}{
		{"0.1.84", "0.1.84", true},  // exactly the minimum
		{"0.1.84", "0.1.85", true},  // newer satisfies
		{"0.1.84", "0.2.0", true},   // greater minor
		{"0.1.84", "1.0.0", true},   // greater major
		{"0.1.84", "0.1.83", false}, // smaller patch
		{"0.1.84", "0.1.9", false},  // 9 < 84: a NUMBER, not text
		{"0.2.0", "0.1.99", false},  // minor rules over patch
		{"1.0.0", "0.99.99", false}, // major rules over everything
	}
	for _, c := range cases {
		got, err := AtendeMinVersion(c.minimum, c.running)
		if err != nil {
			t.Errorf("min=%s running=%s: unexpected error %v", c.minimum, c.running, err)
			continue
		}
		if got != c.satisfies {
			t.Errorf("min=%s running=%s: want satisfies=%v", c.minimum, c.running, c.satisfies)
		}
	}
}

// Comparing as TEXT would give the wrong answer here, and this is the case that exposes it:
// "0.1.9" > "0.1.84" in lexicographic order, and it is SMALLER as a version.
func TestAtendeMinVersion_doesNotCompareAsText(t *testing.T) {
	t.Run("MNVRM-B01: Versions are ordered part by part as numbers", func(t *testing.T) {})
	if "0.1.9" <= "0.1.84" {
		t.Skip("the test's premise changed: lexicographic order no longer inverts here")
	}
	satisfies, err := AtendeMinVersion("0.1.84", "0.1.9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if satisfies {
		t.Error("0.1.9 is SMALLER than 0.1.84 — the comparison is being made as text")
	}
}

// With no minimum declared, every binary satisfies: most projects do not declare one, and
// requiring the declaration in order to work would invert the default.
func TestAtendeMinVersion_noMinimumEveryoneSatisfies(t *testing.T) {
	t.Run("MNVRM-B03: With no minimum declared, any running binary satisfies it", func(t *testing.T) {})
	for _, running := range []string{"0.0.1", "dev", "", "anything"} {
		satisfies, err := AtendeMinVersion("", running)
		if err != nil {
			t.Errorf("running=%q: with no minimum there should be no error, got %v", running, err)
		}
		if !satisfies {
			t.Errorf("running=%q: with no minimum declared, it should satisfy", running)
		}
	}
}

// `dev` IS NOT ORDERABLE, and pretending it is would be worse than refusing to compare.
//
// A local build may be newer than any release (the tree of whoever develops Anchors) or
// older than all of them. The caller decides what to do with the error; what is not done
// is answering "satisfies" or "does not" without grounds.
func TestAtendeMinVersion_devIsNotOrderable(t *testing.T) {
	t.Run("MNVRM-B04: A running version that cannot be ordered does not satisfy, and says why", func(t *testing.T) {})
	t.Run("MNVRM-X01: A pre-release is never compared as if it were its final release", func(t *testing.T) {})
	for _, running := range []string{"dev", "", "0.1", "0.1.84-rc1", "0.1.85-rc1", "v-nothing", "1.2.3.4"} {
		satisfies, err := AtendeMinVersion("0.1.84", running)
		if err == nil {
			t.Errorf("running=%q should refuse the comparison, and answered without error", running)
		}
		if satisfies {
			t.Errorf("running=%q cannot be ordered and must not satisfy", running)
		}
	}
}

// The tag's `v` is accepted on both sides: the release is called `v0.1.84` and the binary
// reports itself as `0.1.84`, and requiring whoever declares to know which of the two to use
// is the kind of detail that is got wrong once and vanishes.
func TestAtendeMinVersion_acceptsTheTagPrefix(t *testing.T) {
	t.Run("MNVRM-B02: The tag prefix v is accepted on either side", func(t *testing.T) {})
	for _, pair := range [][2]string{
		{"v0.1.84", "0.1.84"},
		{"0.1.84", "v0.1.84"},
		{"v0.1.84", "v0.1.85"},
	} {
		satisfies, err := AtendeMinVersion(pair[0], pair[1])
		if err != nil {
			t.Errorf("min=%s running=%s: error %v", pair[0], pair[1], err)
		}
		if !satisfies {
			t.Errorf("min=%s running=%s: should satisfy", pair[0], pair[1])
		}
	}
}

// The FIELD accepts only `MAJOR.MINOR.PATCH`, and the refusal is at LOAD.
//
// Failing only at comparison would make the effect of a mistyped value be the WARNING
// VANISHING — and the field exists precisely to force an upgrade. A `min_version: latest`
// would produce the silence it should break.
func TestValidarMinVersion_refusesWhatIsNotMajorMinorPatch(t *testing.T) {
	t.Run("MNVRM-B05: A declared minimum that is not MAJOR.MINOR.PATCH is refused, dev included", func(t *testing.T) {})
	bad := []string{
		"dev",        // a development binary is a fact one finds, not a decision one declares
		"latest",     // looks reasonable and compares with nothing
		"0.1",        // the patch is missing
		"0.1.84-rc1", // a pre-release compared as final would answer "satisfies" to whoever has less
		"1.2.3.4",
		"abc",
		"0.1.x",
	}
	for _, v := range bad {
		c := &Config{MinVersion: v}
		if err := c.validarMinVersion(); err == nil {
			t.Errorf("min_version=%q should be refused at load", v)
		}
	}
}

func TestValidarMinVersion_acceptsTheFormatAndEmpty(t *testing.T) {
	t.Run("MNVRM-B06: An absent or well-formed minimum is accepted", func(t *testing.T) {})
	for _, v := range []string{"", "0.1.84", "v0.1.84", "1.0.0", "10.20.30"} {
		c := &Config{MinVersion: v}
		if err := c.validarMinVersion(); err != nil {
			t.Errorf("min_version=%q should be accepted, got: %v", v, err)
		}
	}
}

// The error message must say THE FORMAT and THE CONSEQUENCE. A dry "invalid value" leaves
// the reader guessing, and the field is rare enough that nobody remembers it by heart.
func TestValidarMinVersion_theMessageTeaches(t *testing.T) {
	t.Run("MNVRM-B07: The refusal names the expected format, an example and what a bad value silences", func(t *testing.T) {})
	c := &Config{MinVersion: "latest"}
	err := c.validarMinVersion()
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, want := range []string{"MAJOR.MINOR.PATCH", "0.1.84", "silences", `"latest"`} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the message should contain %q; got: %s", want, err)
		}
	}
}

func TestLoad_refusesABadMinVersion(t *testing.T) {
	t.Run("MNVRM-B05: A declared minimum that is not MAJOR.MINOR.PATCH is refused, dev included", func(t *testing.T) {})
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultFile)
	if err := os.WriteFile(path, []byte("min_version: latest\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "min_version") {
		t.Errorf("loading a project with min_version: latest should fail naming the field, got %v", err)
	}
	if err := os.WriteFile(path, []byte("min_version: 0.1.84\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if c, err := Load(path); err != nil || c.MinVersion != "0.1.84" {
		t.Errorf("a well-formed min_version should load, got %v, %v", c, err)
	}
}

func TestVersionOrder_isAntisymmetric(t *testing.T) {
	t.Run("MNVRM-I01: Swapping the two versions always flips the order", func(t *testing.T) {})
	versions := []string{"0.1.9", "0.1.84", "v0.2.0", "1.0.0", "0.99.99"}
	for _, a := range versions {
		for _, b := range versions {
			ab, err1 := VersionOrder(a, b)
			ba, err2 := VersionOrder(b, a)
			if err1 != nil || err2 != nil {
				t.Fatalf("%s vs %s: %v %v", a, b, err1, err2)
			}
			if ab != -ba {
				t.Errorf("VersionOrder(%s,%s)=%d but VersionOrder(%s,%s)=%d", a, b, ab, b, a, ba)
			}
			if (a == b) != (ab == 0) {
				t.Errorf("VersionOrder(%s,%s)=%d", a, b, ab)
			}
		}
	}
}
