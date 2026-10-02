// @anchors
//   ref: CDGNC

package code

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// withSlots runs the test at a fixed generated length and restores it afterwards.
func withSlots(t *testing.T, n int) {
	t.Helper()
	prev := Slots
	Slots = n
	t.Cleanup(func() { Slots = prev })
}

func TestGenerateShape(t *testing.T) {
	t.Run("CDGNC-B05: A single word takes consonants before vowels", func(t *testing.T) {})
	withSlots(t, 5)
	// I do not demand the reference app's EXACT code for every name; I demand the SHAPE:
	// Slots chars, upper case, deterministic.
	for _, name := range []string{"Spacer", "Divider", "HomeScreen", "Button", "LoginScreen"} {
		c := Generate(name)
		if len(c) != Slots {
			t.Errorf("%s → %q: want %d chars", name, c, Slots)
		}
		for i := 0; i < len(c); i++ {
			if !(c[i] >= 'A' && c[i] <= 'Z') && !(c[i] >= '0' && c[i] <= '9') {
				t.Errorf("%s → %q: a char that is not upper case or a digit", name, c)
			}
		}
	}
	for name, want := range map[string]string{
		"Spacer": "SPCRA", "Login": "LGNOI", "AlertsScreen": "LRTSA",
		// z and Z are consonants, 0 and 9 are digits: the ends of each range count.
		"Buzz": "BZZUX", "BUZZ": "BZZUX", "X0": "X0XXX", "X9": "X9XXX",
		// a vowel repeated in another case is taken once: "Anna" has one A, not two.
		"Anna": "NNAXX",
	} {
		if got := Generate(name); got != want {
			t.Errorf("%s → %q, want %q", name, got, want)
		}
	}
}

func TestStripGeneric(t *testing.T) {
	t.Run("CDGNC-B01: Generic suffixes are dropped unless they are the whole name", func(t *testing.T) {})
	cases := map[string]string{
		"LoginScreen": "Login", "AlertSheet": "Alert", "ScrollableLayout": "Scrollable",
		"Button": "Button", "Screen": "Screen", // only the generic → kept
	}
	for in, want := range cases {
		if got := StripGeneric(in); got != want {
			t.Errorf("StripGeneric(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestGenerate_wordSplitting(t *testing.T) {
	t.Run("CDGNC-B02: Words are split at separators, camel case and acronyms", func(t *testing.T) {})
	withSlots(t, 5)
	for name, want := range map[string]string{"user_profile-edit": "UPESR", "ABCParser": "ABPRB"} {
		if got := Generate(name); got != want {
			t.Errorf("%s → %q, want %q", name, got, want)
		}
	}
	// The ends of each letter range decide a boundary as much as the middle does.
	for name, want := range map[string][]string{
		"fooAbc":    {"foo", "Abc"},
		"fooZed":    {"foo", "Zed"},
		"dataBase":  {"data", "Base"},
		"fizzBuzz":  {"fizz", "Buzz"},
		"fooB":      {"foo", "B"},
		"AParser":   {"A", "Parser"},
		"XYZParser": {"XYZ", "Parser"},
		"XMLTzar":   {"XML", "Tzar"},
		"ABC":       {"ABC"},
		// only a lowercase ASCII letter after the capital closes an acronym; é is not one
		"ABCé": {"ABCé"},
	} {
		if got := tokenize(name); !reflect.DeepEqual(got, want) {
			t.Errorf("tokenize(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestGenerate_longNameTakesInitials(t *testing.T) {
	t.Run("CDGNC-B03: A long name takes the initials", func(t *testing.T) {})
	withSlots(t, 5)
	if got := Generate("UserProfileEditFormDraft"); got != "UPEFD" {
		t.Errorf("UserProfileEditFormDraft → %q, want UPEFD", got)
	}
}

func TestGenerateMultiWordUsesEachWord(t *testing.T) {
	t.Run("CDGNC-B04: A two-word name takes letters from each word", func(t *testing.T) {})
	withSlots(t, 5)
	// TransactionDetail (2 words) → must draw from both: T...D...
	c := Generate("TransactionDetail")
	if c != "TRDTT" {
		t.Errorf("TransactionDetail → %q, want TRDTT", c)
	}
	if c[0] != 'T' || !contains([]byte(c), 'D') {
		t.Errorf("TransactionDetail → %q: must start with T and hold the D of the 2nd word", c)
	}
	// A word short of consonants completes its share with vowels, and never takes more than
	// its share: the next word keeps its own.
	withSlots(t, 6)
	if got := Generate("AudioPlayer"); got != "ADUPLY" {
		t.Errorf("AudioPlayer at 6 → %q, want ADUPLY", got)
	}
}

func TestGenerateUniqueResolvesCollision(t *testing.T) {
	t.Run("CDGNC-I01: Generation is deterministic and a resolved code is free", func(t *testing.T) {})
	base := Generate("Spacer")
	taken := map[string]bool{base: true}
	got := GenerateUnique("Spacer", taken)
	if got == base {
		t.Fatalf("collision not resolved: %q == %q", got, base)
	}
	if taken[got] {
		t.Fatalf("the resolved code %q is still taken", got)
	}
	if len(got) != Slots {
		t.Fatalf("the resolved code %q does not have %d chars", got, Slots)
	}
	if again := GenerateUnique("Spacer", taken); again != got {
		t.Fatalf("not deterministic: %q then %q", got, again)
	}
}

func TestGenerateUnique_prefixAndCanonical(t *testing.T) {
	t.Run("CDGNC-B08: A collision varies the last position and keeps the prefix", func(t *testing.T) {})
	withSlots(t, 5)
	if got := GenerateUnique("Divider", map[string]bool{}); got != Generate("Divider") {
		t.Errorf("without a collision it should return the generated code, got %q", got)
	}
	taken := map[string]bool{"AULGN": true}
	if got := GenerateUniqueWithPrefix("Login", "AU", taken); got != "AULGA" {
		t.Errorf("collision on AULGN → %q, want AULGA (last position varied, prefix kept)", got)
	}
	// With the whole last position taken, the one before varies — but never the prefix.
	for c := byte('A'); c <= 'Z'; c++ {
		taken["AULG"+string(c)] = true
	}
	if got := GenerateUniqueWithPrefix("Login", "AU", taken); got[:2] != "AU" || got == "AULGN" || taken[got] {
		t.Errorf("a full last position → %q, want a free code that keeps the AU prefix", got)
	}
	// Every position outside the prefix full: the prefix is still never varied.
	for pos := 2; pos < 5; pos++ {
		for c := byte('A'); c <= 'Z'; c++ {
			v := []byte("AULGN")
			v[pos] = c
			taken[string(v)] = true
		}
	}
	if got := GenerateUniqueWithPrefix("Login", "AU", taken); got != "AULGN" {
		t.Errorf("only the prefix positions left → %q, want AULGN: the prefix is never varied", got)
	}
	// Z is tried too: with A..Y taken in the last position, the last position becomes Z.
	upToY := map[string]bool{}
	for c := byte('A'); c < 'Z'; c++ {
		upToY["AULG"+string(c)] = true
	}
	if got := GenerateUniqueWithPrefix("Login", "AU", upToY); got != "AULGZ" {
		t.Errorf("A..Y taken in the last position → %q, want AULGZ", got)
	}
	// The first position after the prefix is varied too, when every later one is full.
	laterFull := map[string]bool{}
	for pos := 3; pos < 5; pos++ {
		for c := byte('A'); c <= 'Z'; c++ {
			v := []byte("AULGN")
			v[pos] = c
			laterFull[string(v)] = true
		}
	}
	if got := GenerateUniqueWithPrefix("Login", "AU", laterFull); got != "AUAGN" {
		t.Errorf("positions after the third full → %q, want AUAGN", got)
	}
}

func TestGenerateUnique_saturatedNamespaceReturnsTheGeneratedCode(t *testing.T) {
	t.Run("CDGNC-X01: A saturated namespace returns the generated code", func(t *testing.T) {})
	base := Generate("Spacer")
	taken := map[string]bool{base: true}
	for pos := 0; pos < len(base); pos++ {
		for c := byte('A'); c <= 'Z'; c++ {
			v := []byte(base)
			v[pos] = c
			taken[string(v)] = true
		}
	}
	if got := GenerateUnique("Spacer", taken); got != base {
		t.Errorf("saturated namespace → %q, want the generated %q", got, base)
	}
}

func TestShortNamePadsWithX(t *testing.T) {
	t.Run("CDGNC-B06: Short codes are padded with X and existing codes are completed the same way", func(t *testing.T) {})
	withSlots(t, 5)
	// a short name → completed with X ("Ok" has few letters)
	if c := Generate("Ok"); c != "KOXXX" {
		t.Fatalf("Ok → %q, want KOXXX", c)
	}
	if got := Pad("mtvr"); got != "MTVRX" {
		t.Errorf("Pad(mtvr) = %q, want MTVRX", got)
	}
	if got := Pad("ABCDEFG"); got != "ABCDE" {
		t.Errorf("Pad(ABCDEFG) = %q, want ABCDE", got)
	}
}

func TestGenerateWithPrefix(t *testing.T) {
	t.Run("CDGNC-B07: A module prefix starts the code", func(t *testing.T) {})
	withSlots(t, 5)
	if got := GenerateWithPrefix("Login", "au"); got != "AULGN" {
		t.Errorf("Login with prefix au → %q, want AULGN", got)
	}
	if got := GenerateWithPrefix("Login", "abcdefg"); got != "ABCDE" {
		t.Errorf("Login with prefix abcdefg → %q, want ABCDE", got)
	}
}

func TestModulePrefix(t *testing.T) {
	t.Run("CDGNC-B09: The module prefix is the initial and the first consonant", func(t *testing.T) {})
	for in, want := range map[string]string{"auth": "AT", "family": "FM", "i": "IX"} {
		if got := ModulePrefix(in); got != want {
			t.Errorf("ModulePrefix(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSlots_followTheProjectsSmallestLength(t *testing.T) {
	t.Run("CDGNC-B10: The generated length follows the smallest declared length", func(t *testing.T) {})
	withSlots(t, 5)
	prevLengths := config.CodeLengths
	t.Cleanup(func() { config.SetCodeLengths(prevLengths) })

	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("code_lengths: [5, 4]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(p); err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if Slots != 4 || Generate("Login") != "LGNO" {
		t.Fatalf("after loading code_lengths [5, 4]: Slots=%d, Login → %q; want 4, LGNO", Slots, Generate("Login"))
	}
	SetSlots([]int{1})
	SetSlots(nil)
	if Slots != 4 {
		t.Errorf("lengths below 2 or none must leave the length unchanged, got %d", Slots)
	}
	SetSlots([]int{2})
	if Slots != 2 {
		t.Errorf("a declared length of 2 is the smallest accepted, got %d", Slots)
	}
}
