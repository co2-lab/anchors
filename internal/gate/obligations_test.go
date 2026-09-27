package gate

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
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
	t.Run("BLGTN-B03: The pack's duties come first, the inline ones last, with the pack's fields carried over", func(t *testing.T) {})
	t.Run("BLGTN-B04: A pack duty's reason cites the norm it comes from", func(t *testing.T) {})
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
	t.Run("BLGTN-B02: Without packs the duties in force are the inline list", func(t *testing.T) {})
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
	t.Run("BLGTN-B05: A pack that fails to load keeps the inline duties", func(t *testing.T) {})
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
	t.Run("BLGTN-B04: A pack duty's reason cites the norm it comes from", func(t *testing.T) {})
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

func TestEvaluateObligations_oneStatusPerDuty(t *testing.T) {
	t.Run("BLGTN-B07: The report gives one status per duty, in the order handed", func(t *testing.T) {})
	t.Run("BLGTN-B08: A node is a subject only when its header carries the duty's trigger", func(t *testing.T) {})
	t.Run("BLGTN-B09: A waiver with a reason counts as fulfilled and as waived", func(t *testing.T) {})
	t.Run("BLGTN-B10: A duty declared pending counts as debt", func(t *testing.T) {})
	t.Run("BLGTN-B11: Every other subject is missing, and the missing list is sorted", func(t *testing.T) {})
	t.Run("BLGTN-I01: Every subject is counted exactly once", func(t *testing.T) {})
	t.Run("BLGTN-E02: A node whose file cannot be read is not a subject", func(t *testing.T) {})
	root := t.TempDir()
	// purge.ts mentions only the fulfilled models
	writePackFile(t, root, "purge.ts", "const t = [PAID_ONE, PAID_TWO]\n")
	files := map[string]string{
		"models/PaidOne.spec.md": "<!-- @anchors\n  carries: pii\n-->\n",
		"models/PaidTwo.spec.md": "<!-- @anchors\n  carries: pii\n-->\n",
		"models/Waived.spec.md":  "<!-- @anchors\n  carries: pii\n  obligation_waived: pii-purgavel — shared data\n-->\n",
		"models/Owed.spec.md":    "<!-- @anchors\n  carries: pii\n  obligation_pending: pii-purgavel — phase 3\n-->\n",
		"models/Zeta.spec.md":    "<!-- @anchors\n  carries: pii\n-->\n",
		"models/Alpha.spec.md":   "<!-- @anchors\n  carries: pii\n-->\n",
		"models/Neutral.spec.md": "<!-- @anchors\n  code: ABCDE\n-->\n",
	}
	g := &mapx.Graph{}
	for id, body := range files {
		writePackFile(t, root, id, body)
		g.Nodes = append(g.Nodes, mapx.Node{ID: id, Kind: mapx.KindSpec})
	}
	g.Nodes = append(g.Nodes, mapx.Node{ID: "models/Gone.spec.md", Kind: mapx.KindSpec}) // not on disk

	obs := append(obligCfg().Obligations, config.Obligation{Name: "untriggered", MustAppearIn: []string{"purge.ts"}})
	got := EvaluateObligations(root, &config.Config{}, g, obs)

	want := []ObligationStatus{
		{
			Name:      "pii-purgavel",
			Targets:   []string{"purge.ts"},
			Subject:   6,
			Fulfilled: 3, // two that appear, one waived with a reason
			Debt:      1,
			Waived:    1,
			Missing:   []string{"models/Alpha.spec.md", "models/Zeta.spec.md"},
		},
		{Name: "untriggered", Targets: []string{"purge.ts"}}, // no `when`: nobody is subject
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("EvaluateObligations:\n got %+v\nwant %+v", got, want)
	}
}

func TestObligationsInForce_isTheResolvedList(t *testing.T) {
	t.Run("BLGTN-B01: The duties in force are the resolved list", func(t *testing.T) {})
	cfg := obligCfg()
	got := ObligationsInForce(t.TempDir(), cfg)
	if !reflect.DeepEqual(got, cfg.Obligations) {
		t.Errorf("got %+v", got)
	}
}

// The duties in force, seen from outside the package, are the same resolved list the
// gate uses: pack duties included.
func TestObligationsInForce_includesThePackDuties(t *testing.T) {
	t.Run("BLGTN-B01: The duties in force are the resolved list", func(t *testing.T) {})
	root := t.TempDir()
	writePackFile(t, root, "packs/lgpd.yaml", erasurePack)
	cfg := &config.Config{Packs: []string{"packs/lgpd.yaml"}, Obligations: []config.Obligation{{Name: "local"}}}
	got := ObligationsInForce(root, cfg)
	if !reflect.DeepEqual(got, allObligations(root, cfg)) || len(got) != 3 {
		t.Errorf("ObligationsInForce should be the resolved list (pack + inline): %+v", got)
	}
}

// Reading the packs once per root: a pack removed after the first read is still served
// from the cache.
func TestAllObligations_cachedPerRoot(t *testing.T) {
	t.Run("BLGTN-B06: The pack duties are read once per root", func(t *testing.T) {})
	root := t.TempDir()
	writePackFile(t, root, "packs/lgpd.yaml", erasurePack)
	cfg := &config.Config{Packs: []string{"packs/lgpd.yaml"}}
	if got := allObligations(root, cfg); len(got) != 2 {
		t.Fatalf("first read: %+v", got)
	}
	if err := os.Remove(filepath.Join(root, "packs/lgpd.yaml")); err != nil {
		t.Fatal(err)
	}
	if got := allObligations(root, cfg); len(got) != 2 {
		t.Errorf("the second call should come from the cache: %+v", got)
	}
}

// A pack error is a configuration error and must be seen: it goes to stderr.
func TestAllObligations_aBrokenPackIsReportedOnStderr(t *testing.T) {
	t.Run("BLGTN-B05: A pack that fails to load keeps the inline duties", func(t *testing.T) {})
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	allObligations(t.TempDir(), &config.Config{Packs: []string{"packs/missing.yaml"}})
	os.Stderr = old
	w.Close()
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "packs") {
		t.Errorf("the broken pack was not reported on stderr: %q", out)
	}
}

// The report evaluates exactly the duties it is handed: it does not re-read the packs of
// the configuration it receives.
func TestEvaluateObligations_evaluatesOnlyTheHandedDuties(t *testing.T) {
	t.Run("BLGTN-X01: The report evaluates only the duties it is handed", func(t *testing.T) {})
	root := t.TempDir()
	writePackFile(t, root, "packs/lgpd.yaml", erasurePack)
	cfg := &config.Config{Packs: []string{"packs/lgpd.yaml"}}
	g := &mapx.Graph{}
	got := EvaluateObligations(root, cfg, g, []config.Obligation{{Name: "only"}})
	if len(got) != 1 || got[0].Name != "only" {
		t.Errorf("expected exactly the handed duty: %+v", got)
	}
}
