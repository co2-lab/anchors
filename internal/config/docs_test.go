package config

import (
	"reflect"
	"testing"
)

func docsFixture() *Config {
	return &Config{Docs: &Docs{Required: []DocArtifact{
		{Kind: KindSchema, Path: "docs/data.md", Trigger: []string{"DTSTD"}},
		{Kind: KindOpenAPI, Path: "docs/api.yaml", Trigger: []string{"lambdas"}},
		{Kind: KindC4, Path: "docs/architecture.md"}, // no trigger
	}}}
}

// THE LAYER ALONE GETS IT WRONG, and that is why a trigger accepts a unit.
//
// Measured in the reference project: the `infra` layer has NINE units, and only one touches
// the data schema — the other eight are API Gateway, authentication, mTLS, a process ruler.
// Charging the schema to all eight teaches the agent to ignore the warning, which is the
// worst possible outcome: the gate stays and nobody reads it.
func TestRequiredFor_theUnitIsMorePreciseThanTheLayer(t *testing.T) {
	t.Run("DCRQA-B02: A trigger naming a unit code charges that unit and not its neighbours in the same layer", func(t *testing.T) {})
	c := docsFixture()

	// The unit that TOUCHES the schema is charged.
	if d := c.RequiredFor("infra", "DTSTD"); len(d) != 1 || d[0].Kind != KindSchema {
		t.Errorf("DTSTD (the database) should owe the schema, got %v", d)
	}
	// The neighbour unit, in the SAME layer, is not.
	if d := c.RequiredFor("infra", "GLCGL"); len(d) != 0 {
		t.Errorf("GLCGL (a process ruler) was charged %v — it touches no table", d)
	}
}

// The trigger by LAYER still holds: whoever declared `[lambdas]` wants the whole layer.
func TestRequiredFor_theLayerTriggerStillHolds(t *testing.T) {
	t.Run("DCRQA-B01: A trigger naming a layer charges every change in that layer, whatever the unit", func(t *testing.T) {})
	c := docsFixture()
	for _, code := range []string{"", "ANYUN"} {
		if d := c.RequiredFor("lambdas", code); len(d) != 1 || d[0].Kind != KindOpenAPI {
			t.Errorf("code=%q: lambdas should owe the OpenAPI, got %v", code, d)
		}
	}
}

// With no code, the answer is the LAYER's — what `--layer` asks, in the abstract.
func TestRequiredFor_withoutCodeAnswersForTheLayer(t *testing.T) {
	t.Run("DCRQA-B03: Asked without a unit code, the answer is the layer's alone", func(t *testing.T) {})
	c := docsFixture()
	if d := c.RequiredFor("infra"); len(d) != 0 {
		t.Errorf("the `infra` layer has no layer trigger, and got %v", d)
	}
	if d := c.RequiredFor("infra", ""); len(d) != 0 {
		t.Errorf("an empty code is no code, and got %v", d)
	}
	// A blank trigger item names no unit, so an absent code cannot match it.
	blank := DocArtifact{Trigger: []string{" "}}
	if blank.TriggeredBy("infra", "") {
		t.Error("a blank trigger matched a question asked without a code")
	}
}

// A doc WITHOUT a trigger is charged to nobody: the C4 case, which changes when the
// STRUCTURE changes — not when a unit changes.
func TestRequiredFor_aDocWithoutTriggerIsNeverCharged(t *testing.T) {
	t.Run("DCRQA-B04: A documentation with no trigger is never owed by a unit change", func(t *testing.T) {})
	c := docsFixture()
	for _, pair := range [][2]string{{"infra", "DTSTD"}, {"lambdas", "X"}, {"", ""}} {
		for _, d := range c.RequiredFor(pair[0], pair[1]) {
			if d.Kind == KindC4 {
				t.Errorf("the C4 was charged to %v — it has no trigger", pair)
			}
		}
	}
}

func TestTriggeredBy_ignoresCaseAndSpaces(t *testing.T) {
	t.Run("DCRQA-B05: A trigger matches ignoring case and the spaces around it", func(t *testing.T) {})
	d := DocArtifact{Trigger: []string{" Lambdas ", " dtstd"}}
	if !d.TriggeredBy("LAMBDAS", "") {
		t.Error("the layer should match regardless of case and spaces")
	}
	if !d.TriggeredBy("infra", "DTSTD") {
		t.Error("the unit code should match regardless of case and spaces")
	}
}

func TestAllRequiredDocs(t *testing.T) {
	t.Run("DCRQA-B06: Every declared documentation is listed, and a project with no docs block owes none", func(t *testing.T) {})
	c := docsFixture()
	if got := c.AllRequiredDocs(); !reflect.DeepEqual(got, c.Docs.Required) {
		t.Errorf("AllRequiredDocs = %v, want the declared list", got)
	}
	var nilCfg *Config
	for name, cfg := range map[string]*Config{"nil config": nilCfg, "no docs block": {}} {
		if got := cfg.AllRequiredDocs(); got != nil {
			t.Errorf("%s: AllRequiredDocs = %v, want nil", name, got)
		}
		if got := cfg.RequiredFor("lambdas", "DTSTD"); got != nil {
			t.Errorf("%s: RequiredFor = %v, want nil", name, got)
		}
	}
}

func TestRequiredFor_aCodeNeverChargesLess(t *testing.T) {
	t.Run("DCRQA-I01: Naming the unit never removes a documentation the layer alone owes", func(t *testing.T) {})
	c := docsFixture()
	for _, layer := range []string{"lambdas", "infra", ""} {
		byLayer := c.RequiredFor(layer)
		for _, code := range []string{"DTSTD", "OTHER"} {
			withCode := map[string]bool{}
			for _, d := range c.RequiredFor(layer, code) {
				withCode[d.Path] = true
			}
			for _, d := range byLayer {
				if !withCode[d.Path] {
					t.Errorf("layer %q: %s is owed by the layer but not with code %s", layer, d.Path, code)
				}
			}
		}
	}
}

func TestTriggeredBy_neverBySubstring(t *testing.T) {
	t.Run("DCRQA-X01: A trigger is matched whole, never as part of a longer name", func(t *testing.T) {})
	d := DocArtifact{Trigger: []string{"lambdas", "DTSTD"}}
	for _, tc := range [][2]string{{"lambda", ""}, {"lambdas-edge", ""}, {"x", "DTST"}, {"x", "DTSTDX"}} {
		if d.TriggeredBy(tc[0], tc[1]) {
			t.Errorf("TriggeredBy(%q, %q) matched a part of a trigger", tc[0], tc[1])
		}
	}
}
