package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

const erasurePack = `name: lgpd
domain: privacy
authority: "Lei 13.709/2018"
obligations:
  - name: lgpd-erasure
    article: "Art. 18, VI"
    when: "carries: pii"
    must_appear_in: ["purge.ts"]
    identified_by: screaming-snake
    because: "personal data must be erasable"
  - name: lgpd-log
    when: "carries: pii"
    must_appear_in: ["audit.ts"]
`

func writePackFile(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAllObligations_packFirstThenInline(t *testing.T) {
	root := t.TempDir()
	writePackFile(t, root, "packs/lgpd.yaml", erasurePack)
	inline := config.Obligation{Name: "local", When: "carries: pii", MustAppearIn: []string{"x.ts"}}
	cfg := &config.Config{Packs: []string{"packs/lgpd.yaml"}, Obligations: []config.Obligation{inline}}

	for _, round := range []string{"loaded", "cached"} {
		got := allObligations(root, cfg)
		if len(got) != 3 {
			t.Fatalf("%s: pack (2) + inline (1), got %+v", round, got)
		}
		if got[0].Name != "lgpd-erasure" || got[1].Name != "lgpd-log" || got[2].Name != "local" {
			t.Errorf("%s: the pack comes first, the inline declaration last: %+v", round, got)
		}
		if want := "personal data must be erasable (Lei 13.709/2018, Art. 18, VI)"; got[0].Because != want {
			t.Errorf("%s: the reason cites the norm: %q", round, got[0].Because)
		}
		if got[0].IdentifiedBy != "screaming-snake" || got[0].When != "carries: pii" || got[0].MustAppearIn[0] != "purge.ts" {
			t.Errorf("%s: the pack's fields carry over: %+v", round, got[0])
		}
		if want := "Lei 13.709/2018"; got[1].Because != want {
			t.Errorf("%s: without a reason or article, the authority alone is the reason: %q", round, got[1].Because)
		}
	}
}

func TestAllObligations_withoutPacksIsTheInlineList(t *testing.T) {
	if allObligations(t.TempDir(), nil) != nil {
		t.Error("no config, no obligations")
	}
	inline := []config.Obligation{{Name: "local"}}
	got := allObligations(t.TempDir(), &config.Config{Obligations: inline})
	if len(got) != 1 || got[0].Name != "local" {
		t.Errorf("got %+v", got)
	}
}

func TestAllObligations_aBrokenPackKeepsTheInlineDuties(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{
		Packs:       []string{"packs/missing.yaml"},
		Obligations: []config.Obligation{{Name: "local"}},
	}
	got := allObligations(root, cfg)
	if len(got) != 1 || got[0].Name != "local" {
		t.Errorf("a pack that fails to load drops only its own duties: %+v", got)
	}
}

func TestWithSource(t *testing.T) {
	cases := []struct{ because, authority, article, want string }{
		{"", "", "", ""},
		{"why", "", "", "why"},
		{"", "Law", "", "Law"},
		{"", "", "Art. 1", "Art. 1"},
		{"", "Law", "Art. 1", "Law, Art. 1"},
		{"why", "Law", "", "why (Law)"},
		{"why", "", "Art. 1", "why (Art. 1)"},
		{"why", "Law", "Art. 1", "why (Law, Art. 1)"},
	}
	for _, c := range cases {
		if got := withSource(c.because, c.authority, c.article); got != c.want {
			t.Errorf("withSource(%q, %q, %q) = %q, want %q", c.because, c.authority, c.article, got, c.want)
		}
	}
}
