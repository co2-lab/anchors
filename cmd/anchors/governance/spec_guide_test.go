package governance

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
)

func TestSpecGuideHasACompleteExample(t *testing.T) {
	t.Run("SPGDS-B02: The guide shows the three rule forms and a complete example", func(t *testing.T) {})
	t.Run("SPGDS-B06: The guide starts from the command that generates the skeleton", func(t *testing.T) {})
	// The reason this file exists: the built-in guide (`anchors guide spec`) has rules and
	// ZERO examples. Nobody writes conforming markdown from a description — one copies an
	// example.
	g := renderSpecGuide(&config.Config{}, "LOGIX")
	for _, required := range []string{
		"<!-- @anchors",    // the header, the easiest thing to get wrong
		"code: LOGIX",      // with the code in the right place
		"### LOGIX-B01",    // the heading form of a catalogued rule
		"| `LOGIX-B02` |",  // the table-row form
		"- **LOGIX-B03**",  // the bold-bullet form
		"### LOGIX-S01",    // the example's rules
		"anchors new spec", // the command that solves it, BEFORE the doctrine
		"--list-sections",  // and the one that lists the sections with their criterion
	} {
		if !strings.Contains(g, required) {
			t.Errorf("the guide must contain %q", required)
		}
	}
	if strings.Index(g, "anchors new spec") > strings.Index(g, "## The format the gate requires") {
		t.Error("the command must come before the format")
	}
}

// The example has to match what `anchors new spec` writes in THIS project: the section
// titles come from the same catalogue the generator uses, resolved by `lang:`. Pinning them
// here would show the reader an example that `spec-sections` rejects.
func TestSpecGuideTitlesComeFromTheCatalogue(t *testing.T) {
	t.Run("SPGDS-B03: The section titles come from the catalogue the generator uses", func(t *testing.T) {})
	g := renderSpecGuide(&config.Config{}, "LOGIX")
	for _, key := range []string{"overview", "rules", "open"} {
		title := i18n.TIn(i18n.Current(), "section.title."+key)
		if title == "" {
			t.Fatalf("the catalogue has no title for %q", key)
		}
		if !strings.Contains(g, "## "+title+"\n") {
			t.Errorf("the example should use the catalogue's title %q", title)
		}
	}
	// without `lang:` the framework's default holds (English)
	if !strings.Contains(g, "## Open Decisions") {
		t.Error("with no lang declared the titles are the English ones")
	}
}

func TestSpecGuideDefaultsWithNothingDeclared(t *testing.T) {
	t.Run("SPGDS-B01: An empty example code defaults to LOGI", func(t *testing.T) {})
	t.Run("SPGDS-B05: The code length is stated only when the project declares it", func(t *testing.T) {})
	g := RenderSpecGuide(&config.Config{}, "")
	if !strings.Contains(g, "code: LOGI\n") || !strings.Contains(g, "### LOGI-B01") {
		t.Errorf("an empty example code should fall back to LOGI:\n%s", g)
	}
	if strings.Contains(g, "code_lengths") {
		t.Error("a project that declares no code_lengths gets no length sentence")
	}
	g = RenderSpecGuide(&config.Config{CodeLengths: []int{4, 5}}, "")
	if !strings.Contains(g, "The identity code has 4 or 5 character(s) in this project") {
		t.Errorf("the declared lengths should be stated:\n%s", g)
	}
}

func TestSpecGuideUsesTheProjectDialect(t *testing.T) {
	t.Run("SPGDS-B04: The rule letters are the project's when declared and the canonical ones otherwise", func(t *testing.T) {})
	t.Run("SPGDS-X01: A project that declares its rule types is not offered the canonical letters", func(t *testing.T) {})
	// A generic guide does not serve: the project declares the letters it USES, and an agent
	// reading the canonical ones writes codes the `rule-types` gate rejects.
	cfg := &config.Config{
		RuleTypes:   []config.RuleType{{Letter: "I", Term: "Invariant"}},
		CodeLengths: []int{4},
	}
	g := renderSpecGuide(cfg, "ABCDX")
	if !strings.Contains(g, "| `I` | Invariant |") {
		t.Error("the letters declared in rule_types have to appear in the guide")
	}
	if !strings.Contains(g, "4 character") {
		t.Error("the length declared in code_lengths has to appear")
	}
	// and the canonical list is NOT offered when the project declared its own
	if strings.Contains(g, "canonical letters") {
		t.Error("a project with rule_types must not get the canonical list — it contradicts its own")
	}
}

func TestSpecGuideWithoutRuleTypesOffersTheCanonical(t *testing.T) {
	t.Run("SPGDS-B04: The rule letters are the project's when declared and the canonical ones otherwise", func(t *testing.T) {})
	g := renderSpecGuide(&config.Config{}, "ABCDX")
	if !strings.Contains(g, "canonical letters") || !strings.Contains(g, "`B` behaviour") {
		t.Error("without rule_types, the guide has to say which letters hold")
	}
}
