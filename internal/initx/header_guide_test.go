package initx

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/i18n"
)

func TestRenderHeaderGuideDialect(t *testing.T) {
	t.Run("HDGDH-B01: The comment dialect follows the preset", func(t *testing.T) {})
	hash := map[string]bool{"django": true, "fastapi": true, "python-lib": true, "rails": true, "phoenix": true}
	for _, p := range Presets {
		out := RenderHeaderGuide(p, nil)
		want, other := "// @anchors", "# @anchors"
		if hash[p.Name] {
			want, other = other, want
		}
		if !strings.Contains(out, want) || strings.Contains(out, "\n"+other) {
			t.Errorf("preset %s: the examples should use %q", p.Name, want)
		}
	}
	if !strings.Contains(RenderHeaderGuide(Preset{}, nil), "// @anchors") {
		t.Error("with no preset the dialect is //")
	}
}

func TestRenderHeaderGuideModules(t *testing.T) {
	t.Run("HDGDH-B02: The grouping example names the first module, or auth", func(t *testing.T) {})
	t.Run("HDGDH-B03: The module list appears only when there are modules", func(t *testing.T) {})
	ts, _ := PresetByName("node-ts")

	with := RenderHeaderGuide(ts, []string{"billing", "auth"})
	if !strings.Contains(with, "@feature: billing\n") {
		t.Error("the grouping example should name the first module")
	}
	if !strings.Contains(with, "In this project: billing, auth.") {
		t.Error("the project's real modules should be listed")
	}

	without := RenderHeaderGuide(ts, nil)
	if !strings.Contains(without, "@feature: auth\n") {
		t.Error("with no modules the example should be auth")
	}
	if strings.Contains(without, "In this project:") {
		t.Error("with no modules there is no module list")
	}
}

func TestRenderHeaderGuideAlwaysHasEssentials(t *testing.T) {
	t.Run("HDGDH-B04: The essentials are always present", func(t *testing.T) {})
	t.Run("HDGDH-B05: The compliance-points section is always present with five points", func(t *testing.T) {})
	// even with no preset (Preset{}), the guide mentions the minimum: code + the gate
	out := RenderHeaderGuide(Preset{}, nil)
	for _, want := range []string{"code:", "updated_at:", "header-valid", "anchors guide header"} {
		if !strings.Contains(out, want) {
			t.Errorf("the header guide should mention %q", want)
		}
	}
	title := i18n.TIn(i18n.Current(), "section.title.compliance_points")
	if !strings.Contains(out, "\n## "+title+"\n") {
		t.Errorf("the guide should have the section %q", title)
	}
	for _, ck := range []string{"CK1", "CK2", "CK3", "CK4", "CK5"} {
		if !strings.Contains(out, "**"+ck+"**") {
			t.Errorf("the compliance section should list %s", ck)
		}
	}
}

func TestRenderHeaderGuideTitle(t *testing.T) {
	t.Run("HDGDH-B06: The title falls back to project", func(t *testing.T) {})
	next, _ := PresetByName("nextjs")
	if !strings.HasPrefix(RenderHeaderGuide(next, nil), "# Header guide — "+next.Title+"\n") {
		t.Errorf("the guide should be titled after %q", next.Title)
	}
	if !strings.HasPrefix(RenderHeaderGuide(Preset{Name: "x"}, nil), "# Header guide — project\n") {
		t.Error("a preset with no title should give the title \"project\"")
	}
}

func TestRenderHeaderGuideWritesNothing(t *testing.T) {
	t.Run("HDGDH-X01: Rendering writes nothing to disk", func(t *testing.T) {})
	dir := t.TempDir()
	t.Chdir(dir)
	if out := RenderHeaderGuide(Preset{}, []string{"auth"}); out == "" {
		t.Fatal("the guide should come back as text")
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("rendering wrote to disk: %v", entries)
	}
}

// The guide `anchors init` SEEDS must pass the gate init itself declares — in any language.
//
// The coupling: `RenderHeaderGuide` writes the compliance section translated by the project's
// `lang:`, and `guide-checklist` (internal/gate) looks for it by title. While the gate's regex
// knew only the Portuguese form, a `lang: en` project was born failing a file init had just
// written — the framework charging for what it does not itself generate, and the worst first
// impression possible.
//
// This test duplicates the gate's ruler on purpose (the regex below mirrors
// `checklistHeadingRE`): putting the confrontation HERE, on the PRODUCER's side, is what makes
// the break show up for whoever touches the guide — and the sibling test in internal/gate does
// the same for whoever touches the gate. The two ends of the coupling, each with its alarm.
func TestSeededGuideHasTheChecklistInEveryLanguage(t *testing.T) {
	t.Run("HDGDH-I01: The seeded guide passes the checklist heading in every language", func(t *testing.T) {})
	// Mirrors checklistHeadingRE: the section is recognised in any language of the catalog.
	titles := i18n.AllTranslations("section.title.compliance_points")
	if len(titles) == 0 {
		t.Fatal("no titles in the catalog for `section.title.compliance_points` — the gate would have nothing to match")
	}
	escaped := make([]string, 0, len(titles))
	for _, x := range titles {
		escaped = append(escaped, regexp.QuoteMeta(x))
	}
	sectionRE := regexp.MustCompile(`(?mi)^##+\s+(?:` + strings.Join(escaped, "|") + `)\b`)
	itemRE := regexp.MustCompile(`(?m)\bCK\d+\b`)

	original := i18n.Current()
	t.Cleanup(func() { _ = i18n.Set(original) })

	for _, lang := range i18n.SupportedLangs {
		if err := i18n.Set(lang); err != nil {
			t.Fatalf("language %s: %v", lang, err)
		}
		g := RenderHeaderGuide(Preset{}, nil)
		if !sectionRE.MatchString(g) {
			t.Errorf("lang=%s: the seeded guide lacks the compliance section `guide-checklist` requires", lang)
		}
		// The section alone is not enough: without at least one CK point, the judgment gate
		// falls back to heuristics — which is exactly what the checklist exists to prevent.
		if !itemRE.MatchString(g) {
			t.Errorf("lang=%s: the section exists but has no CK point — nothing to confront item by item", lang)
		}
	}
}
