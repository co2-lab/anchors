package ops

import (
	"strings"
	"testing"
)

// Product doctrine does NOT inherit the spec's body.
//
// `sectionBody` resolves the body through `section.body.<key>` in the translation
// catalog, and the template's literal only counts when the key is absent there. With bare
// keys (`rules`, `overview`) the doctrine silently inherited the SPEC's text — measured:
// the command emitted `### XXXXX-B01 — TODO rule` and "what the unit does", spec text in
// an artifact that has no unit at all.
func TestNewProduct_doesNotInheritTheSpecBody(t *testing.T) {
	tpl, ok := templates["product"]
	if !ok {
		t.Fatal("the `product` kind is not registered")
	}
	for _, s := range tpl.sections {
		if !strings.HasPrefix(s.Key, "doctrine_") {
			t.Errorf("section %q without the `doctrine_` prefix: it will inherit another kind's body", s.Key)
		}
	}
	body := renderArtifact(tpl, "CreditLimit", "CRLMT", "product/x.doctrine.md", t.TempDir(),
		map[string]bool{"doctrine_title": true, "doctrine_rules": true}, nil, nil)
	if !strings.Contains(body, "CRLMT-R01") {
		t.Errorf("the doctrine should emit `-R01` rules, got:\n%s", body)
	}
	if strings.Contains(body, "-B01 — TODO rule") {
		t.Errorf("the doctrine inherited the spec body:\n%s", body)
	}
	// Doctrine has no TARGET: the layer belongs to the target a spec describes, and a
	// `layer: TODO` here would be a field nobody can fill — the eternal placeholder that
	// `placeholder-filled` exists to accuse.
	if strings.Contains(tpl.headerFn("CRLMT", "product/x.doctrine.md"), "layer:") {
		t.Error("the doctrine header must not declare `layer:`")
	}
}
