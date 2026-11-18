// @anchors
//   code: SCTSS
//   ref: DCSCD

package doct

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// Every scaffold has a NAME and a BODY. An empty name makes init try to write over the
// `doct/` directory itself — and the message ("is a directory") does not say which scaffold
// is broken.
func TestScaffolds_allHaveNameAndBody(t *testing.T) {
	t.Run("DCSCD-B01: The three fixed templates each have a name, a body and what they answer", func(t *testing.T) {})
	t.Run("DCSCD-B02: A layer's page template is named after the layer", func(t *testing.T) {})
	var names []string
	for _, s := range Scaffolds("en") {
		names = append(names, s.Nome)
	}
	if strings.Join(names, ",") != "architecture.md.tmpl,behavior.md.tmpl,rules.md.tmpl" {
		t.Errorf("fixed scaffolds = %v, want architecture, behaviour and rules", names)
	}
	if got := ScaffoldLayer("infra", "en").Nome; got != "layers/infra.md.tmpl" {
		t.Errorf("layer scaffold name = %q, want layers/infra.md.tmpl", got)
	}
	all := append(Scaffolds("en"), ScaffoldLayer("infra", "en"))
	for i, s := range all {
		if strings.TrimSpace(s.Nome) == "" {
			t.Errorf("scaffold #%d has no name (why: %q)", i, s.Porque)
		}
		if !strings.HasSuffix(s.Nome, SufixoTemplate) {
			t.Errorf("scaffold %q does not end in %s", s.Nome, SufixoTemplate)
		}
		if strings.TrimSpace(s.Corpo) == "" {
			t.Errorf("scaffold %q has no body", s.Nome)
		}
		if strings.TrimSpace(s.Porque) == "" {
			t.Errorf("scaffold %q does not say what the page answers", s.Nome)
		}
	}
}

func TestInitScaffolds_writesFixedAndPerLayer(t *testing.T) {
	t.Run("DCSCD-B03: Init writes the fixed templates and one page per layer, each opening with its purpose", func(t *testing.T) {
		root, g := projetoComFeature(t)
		c, _ := New(root, g)
		written, _, err := c.InitScaffolds(false)
		if err != nil {
			t.Fatal(err)
		}
		want := "architecture.md.tmpl,behavior.md.tmpl,rules.md.tmpl,layers/infra.md.tmpl"
		if strings.Join(written, ",") != want {
			t.Errorf("written = %v, want %s", written, want)
		}
		for _, s := range append(Scaffolds("en"), ScaffoldLayer("infra", "en")) {
			b, err := os.ReadFile(filepath.Join(root, Dir, s.Nome))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(string(b), "{{/* "+s.Porque+" */}}\n") {
				t.Errorf("%s does not open with its purpose comment:\n%.80s", s.Nome, b)
			}
		}
	})
	t.Run("DCSCD-X01: Init writes no compiled page", func(t *testing.T) {
		root, g := projetoComFeature(t)
		c, _ := New(root, g)
		if _, _, err := c.InitScaffolds(false); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(root, OutDir)); err == nil {
			t.Error("init created the compiled documentation folder")
		}
	})
}

// THE SKELETON compiles. A scaffold the compiler itself cannot process would hand the team a
// syntax error as a starting point.
func TestInitScaffolds_theSkeletonCompiles(t *testing.T) {
	t.Run("DCSCD-I01: The skeleton init writes compiles", func(t *testing.T) {})
	t.Run("DCSCD-B05: The small layer page carries each scenario under its own heading, and the index links it", func(t *testing.T) {})
	root, g := projetoComFeature(t)
	c, _ := New(root, g)

	written, _, err := c.InitScaffolds(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) < 4 {
		t.Fatalf("wrote %d templates: %v", len(written), written)
	}

	res, err := c.Build(false)
	if err != nil {
		t.Fatalf("the skeleton Anchors generates does not compile: %v", err)
	}
	if len(res.Written) != len(written) {
		t.Errorf("generated %d docs from %d templates", len(res.Written), len(written))
	}

	// The CONTENT goes into the LAYER page — it is the one the index promises.
	b, err := os.ReadFile(filepath.Join(root, OutDir, "layers", "infra.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "Então a régua não libera") {
		t.Errorf("the scenario body did not reach the layer page:\n%s", b)
	}

	// And the index POINTS there, with an anchor that exists. A link that does not resolve
	// is the worst defect of an index: it is clickable, and the browser stays put.
	idx, err := os.ReadFile(filepath.Join(root, OutDir, "behavior.md"))
	if err != nil {
		t.Fatal(err)
	}
	target := "layers/infra.md#" + GitHubAnchor("GLCGL-B01 — item sem artefato não passa")
	if !strings.Contains(string(idx), target) {
		t.Errorf("the index does not point at `%s`:\n%s", target, idx)
	}
	if !strings.Contains(string(b), "#### GLCGL-B01 — item sem artefato não passa") {
		t.Errorf("the index's anchor does not exist on the target page:\n%s", b)
	}
}

func TestInitScaffolds_bigLayerSummarizes(t *testing.T) {
	t.Run("DCSCD-B06: The big layer page summarizes each unit without scenario steps", func(t *testing.T) {
		root, g := projetoComFeature(t)
		c, _ := New(root, g)
		c.Layout = Layout{MaxUnits: 0, MaxLines: 0}
		if _, _, err := c.InitScaffolds(false); err != nil {
			t.Fatal(err)
		}
		if _, err := c.Build(false); err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(filepath.Join(root, OutDir, "layers", "infra.md"))
		if err != nil {
			t.Fatal(err)
		}
		page := string(b)
		for _, want := range []string{
			c.Layout.Describe(Size{Units: 1, Rules: 3}),   // the summary sentence
			"A régua confronta o que foi prometido",       // the unit's overview
			"**GLCGL-B01** — item sem artefato não passa", // the rule list
		} {
			if !strings.Contains(page, want) {
				t.Errorf("the big layer page lacks %q:\n%s", want, page)
			}
		}
		if strings.Contains(page, "Então a régua não libera") {
			t.Errorf("the big layer page carries scenario steps:\n%s", page)
		}
	})
}

// Running again does not erase what the team edited: an edited template carries the frame
// written by hand, which is precisely the part that is not generated.
func TestInitScaffolds_doesNotOverwriteWithoutForce(t *testing.T) {
	t.Run("DCSCD-B04: An edited template is kept unless forced", func(t *testing.T) {})
	root, g := projetoComFeature(t)
	c, _ := New(root, g)
	c.InitScaffolds(false)

	target := filepath.Join(root, Dir, "architecture.md"+SufixoTemplate)
	edited := "{{/* edited by the team */}}\n# Our architecture\n"
	os.WriteFile(target, []byte(edited), 0o644)

	_, skipped, err := c.InitScaffolds(false)
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(target); string(b) != edited {
		t.Fatalf("init erased the handwritten frame:\n%s", b)
	}
	if len(skipped) == 0 {
		t.Error("skipped in silence — whoever runs it again looks for why nothing changed")
	}

	if _, _, err := c.InitScaffolds(true); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(target); string(b) == edited {
		t.Error("with force the template should be rewritten")
	}
}

func TestInitScaffolds_unwritableFolderFails(t *testing.T) {
	t.Run("DCSCD-E01: A templates folder that cannot be created stops init with the error", func(t *testing.T) {
		root, g := projetoComFeature(t)
		c, _ := New(root, g)
		// A FILE where the templates folder should be: nothing can be created under it.
		if err := os.WriteFile(filepath.Join(root, Dir), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		written, _, err := c.InitScaffolds(false)
		if err == nil {
			t.Fatal("init succeeded with no place to write the templates")
		}
		if len(written) != 0 {
			t.Errorf("written = %v, want none", written)
		}
	})
}

// Level 2 carries the PROTOCOL on each arrow, the database among the boxes, and there is a
// level 3 per internal container — none for the external one.
func TestBuild_theC4FollowsTheModel(t *testing.T) {
	t.Run("DCSCD-B07: The architecture page follows the C4 model", func(t *testing.T) {})
	c := projectWithContainers(t)
	if _, _, err := c.InitScaffolds(false); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, err := readFile(c.Root, "architecture.md")
	if err != nil {
		t.Fatal(err)
	}

	// The protocol: an arrow without it says the two talk and not what happens when the
	// conversation fails — the only thing a level 2 has to say about risk.
	for _, p := range []string{"HTTPS/JSON", "SQL/TLS"} {
		if !strings.Contains(b, p) {
			t.Errorf("protocol %q did not reach the diagram", p)
		}
	}
	if !strings.Contains(b, "o que persiste") {
		t.Error("the database is missing from level 2 — a container is what runs OR STORES")
	}
	if !strings.Contains(b, "### app") || !strings.Contains(b, "### api") {
		t.Errorf("an internal container's level 3 is missing:\n%s", b)
	}
	if strings.Contains(b, "### banco") {
		t.Error("the database got a level 3 section")
	}
}

// With no container declared, the page SAYS so instead of shipping empty diagrams.
func TestBuild_noContainerDeclaredWarns(t *testing.T) {
	t.Run("DCSCD-B08: No container declared is said on the architecture page", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{
		"a/X.spec.md": "---\ncode: XXXXX\nlayer: infra\n---\n\n# X — x\n"})
	c, _ := New(root, g)
	c.Config = &config.Config{}
	c.InitScaffolds(false)
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, _ := readFile(c.Root, "architecture.md")
	if !strings.Contains(b, "No container declared") {
		t.Errorf("the page does not say that the declaration is missing:\n%s", b)
	}
}

// The pages are the project's readers', in the project's language: names, headings, the
// layer folder and the section titles the page asks for. They were Portuguese whatever
// `lang:` declared.
func TestScaffolds_speakTheProjectLanguage(t *testing.T) {
	t.Run("DCSCD-B09: The templates are named and written in the project's language", func(t *testing.T) {})
	for lang, want := range map[string]struct{ names, layer, heading, overview string }{
		"en":    {"architecture,behavior,rules", "layers/infra.md.tmpl", "# Architecture", `{{section . "Overview"}}`},
		"pt-BR": {"arquitetura,comportamento,regras", "camadas/infra.md.tmpl", "# Arquitetura", `{{section . "Visão Geral"}}`},
		"es":    {"arquitectura,comportamiento,reglas", "capas/infra.md.tmpl", "# Arquitectura", `{{section . "Visión General"}}`},
	} {
		var names []string
		fixed := Scaffolds(lang)
		for _, s := range fixed {
			names = append(names, strings.TrimSuffix(s.Nome, ".md"+SufixoTemplate))
		}
		if strings.Join(names, ",") != want.names {
			t.Errorf("%s: names = %v, want %s", lang, names, want.names)
		}
		if !strings.HasPrefix(fixed[0].Corpo, want.heading+"\n") {
			t.Errorf("%s: the architecture page should open with %q", lang, want.heading)
		}
		layer := ScaffoldLayer("infra", lang)
		if layer.Nome != want.layer || !strings.Contains(layer.Corpo, want.overview) {
			t.Errorf("%s: layer page %q should ask %s", lang, layer.Nome, want.overview)
		}
		for _, s := range append(fixed, layer) {
			if strings.Contains(s.Corpo, "«") {
				t.Errorf("%s: %s keeps an unfilled key", lang, s.Nome)
			}
		}
	}
	if Scaffolds("fr")[0].Nome != "architecture.md"+SufixoTemplate {
		t.Error("a language with no table falls back to English")
	}
}

func TestScaffold_aProjectWithAPIUnitsGetsTheOpenAPITemplate(t *testing.T) {
	t.Run("DCSCD-B05: A project with API units gets the OpenAPI template", func(t *testing.T) {})
	for spec, want := range map[string]bool{apiSpec: true, qrSpec: false} {
		root, g := projetoDeTeste(t, map[string]string{"src/a.spec.md": spec})
		c, err := New(root, g)
		if err != nil {
			t.Fatal(err)
		}
		written, _, err := c.InitScaffolds(false)
		if err != nil {
			t.Fatal(err)
		}
		got := strings.Contains(strings.Join(written, ","), "openapi.yaml.tmpl")
		if got != want {
			t.Errorf("openapi.yaml.tmpl written = %v, want %v (%v)", got, want, written)
		}
	}
}
