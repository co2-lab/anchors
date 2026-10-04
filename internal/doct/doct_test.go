// @anchors
//   code: DCTSE
//   ref: DTCDC

package doct

import (
	"fmt"
	"github.com/co2-lab/anchors/internal/config"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const specDeExemplo = `---
code: GLCGL
layer: infra
---

# GoLiveChecklist — a régua que registra o que não foi feito

## Visão Geral

A régua confronta o que foi prometido contra o que existe.

## Regras

### GLCGL-B01 — item sem artefato não passa

Um item que não cita artefato não é conferível.

### GLCGL-B02 — dívida aberta bloqueia

Enquanto a dívida está aberta, a régua não libera.

## Invariantes

### GLCGL-I01 — o resumo nunca mente

O resumo conta o que a régua viu.
`

// projetoDeTeste builds a fake project with the specs on disk and a map that knows them.
func projetoDeTeste(t *testing.T, specs map[string]string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	g := &mapx.Graph{}
	for p, content := range specs {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		code := ""
		if m := headerCodeDeTeste.FindStringSubmatch(content); m != nil {
			code = m[1]
		}
		g.Nodes = append(g.Nodes, mapx.Node{ID: p, Kind: mapx.KindSpec, Code: code, Layer: "spec"})
	}
	return root, g
}

var headerCodeDeTeste = regexp.MustCompile(`(?m)^\s*code:\s*([A-Z0-9]+)`)

func escreveTemplate(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, Dir)
	if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The spec's CONTENT goes into the compiled page — not a link to it. It is why the
// mechanism exists: documentation that only points at files becomes an index.
func TestBuild_specContentEntersTheCompiledPage(t *testing.T) {
	t.Run("DTCDC-B01: The spec's content enters the compiled page", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "infra.md.tmpl",
		`{{range specs "layer=infra"}}## {{.Titulo}}

{{section . "Visão Geral"}}
{{end}}`)

	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(filepath.Join(root, OutDir, "infra.md"))
	if err != nil {
		t.Fatalf("nothing at `%s/infra.md`: %v", OutDir, err)
	}
	out := string(b)
	if !strings.Contains(out, "A régua confronta o que foi prometido") {
		t.Errorf("the section's body did not come in — the page became an index:\n%s", out)
	}
	if !strings.Contains(out, "GoLiveChecklist — a régua") {
		t.Errorf("the title did not come in:\n%s", out)
	}
}

// The layer comes from the spec's HEADER, not from the map: the map files every spec under
// `spec`, and grouping by it would give ONE group with everything.
func TestSpecs_layerComesFromTheSpecHeader(t *testing.T) {
	t.Run("DTCDC-B06: The layer comes from the spec's header", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.fnSpecs("layer=infra")
	if err != nil || len(got) != 1 {
		t.Fatalf("filtering by `infra` (from the header) found %d (%v), want 1", len(got), err)
	}
	if got, _ := c.fnSpecs("layer=spec"); len(got) != 0 {
		t.Error("filtering by the MAP's layer matched something — the header's layer is the one that counts")
	}
}

// The `## Regras` section holds rule `###` headings, and the cut must not stop at the first.
func TestSection_doesNotStopAtSubHeading(t *testing.T) {
	t.Run("DTCDC-B08: A section runs to the next heading of its own level", func(t *testing.T) {})
	s := Spec{raw: specDeExemplo}
	got := fnSection(s, "Regras")
	if !strings.Contains(got, "GLCGL-B02") {
		t.Errorf("the cut stopped before the second rule:\n%s", got)
	}
	if strings.Contains(got, "GLCGL-I01") {
		t.Errorf("the cut ran into `## Invariantes`:\n%s", got)
	}
	if strings.HasPrefix(got, "## ") {
		t.Errorf("the section came with its own heading:\n%s", got)
	}
	if got := fnSection(s, "Ausente"); got != "" {
		t.Errorf("a missing section should be empty, got %q", got)
	}
}

// `rules` returns the rules SEPARATELY — "every rule of the system" is a list of items,
// not the concatenation of `## Regras` sections.
func TestRules_returnsEachRuleSeparately(t *testing.T) {
	t.Run("DTCDC-B09: A spec is split into its separate rules", func(t *testing.T) {})
	rs := fnRules(Spec{raw: specDeExemplo})
	if len(rs) != 3 {
		t.Fatalf("found %d rules, want 3 (B01, B02, I01)", len(rs))
	}
	if rs[0].Code != "GLCGL-B01" || rs[0].Titulo != "item sem artefato não passa" {
		t.Errorf("first rule = %+v", rs[0])
	}
	if strings.Contains(rs[1].Corpo, "Invariantes") {
		t.Errorf("B02's body ran into the next heading: %q", rs[1].Corpo)
	}
	if rs[2].Code != "GLCGL-I01" {
		t.Errorf("the invariant did not come in as a rule: %+v", rs[2])
	}
}

// A template asks for "Visão Geral" and an English spec writes "Overview": the title is read
// as the section it names, in any language of the catalog.
func TestSection_findsTheSameSectionUnderAnotherLanguage(t *testing.T) {
	t.Run("DTCDC-B11: A section asked by its title in one language is found under another", func(t *testing.T) {})
	s := Spec{raw: "# X\n\n## Overview\n\nwhat it does\n\n## Effects\n\n| Effect | Description |\n"}
	if got := fnSection(s, "Visão Geral"); got != "what it does" {
		t.Errorf("Visão Geral on an English spec = %q, want the Overview", got)
	}
	if got := fnSection(s, "Visión General"); got != "what it does" {
		t.Errorf("Visión General on an English spec = %q", got)
	}
	if got := fnSection(s, "Fora de escopo"); got != "" {
		t.Errorf("a title outside the catalog finds only itself, got %q", got)
	}
}

// A spec written in tables and bullets has rules too; the index listed none of them.
func TestRules_readsTheThreeForms(t *testing.T) {
	t.Run("DTCDC-B12: The rules are read in the three catalogued forms, the spec's own codes at their first definition", func(t *testing.T) {})
	raw := "# X\n\n## Effects\n\n| Effect | Description |\n| --- | --- |\n" +
		"| `PSXSH-B01` | a table rule <!-- @no-mark: x --> |\n" +
		"| `OTHER-B01` | another spec's code |\n\n" +
		"- **PSXSH-B02** — a bullet rule\n\n" +
		"### PSXSH-B03 — a heading rule\n\nits body\n\n" +
		"## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `PSXSH-B01` | `field` |\n"
	rs := fnRules(Spec{Code: "PSXSH", raw: raw})
	if len(rs) != 3 {
		t.Fatalf("rules = %+v, want B01, B02 and B03 once each", rs)
	}
	want := []Rule{
		{Code: "PSXSH-B01", Titulo: "a table rule"},
		{Code: "PSXSH-B02", Titulo: "a bullet rule"},
		{Code: "PSXSH-B03", Titulo: "a heading rule", Corpo: "its body", Heading: true},
	}
	for i, w := range want {
		if rs[i] != w {
			t.Errorf("rule %d = %+v, want %+v", i, rs[i], w)
		}
	}
}

// The compiled page carries the MARKER: whoever opens it knows editing it loses the work.
func TestBuild_writesTheMarker(t *testing.T) {
	t.Run("DTCDC-B02: The compiled page opens with the generated marker, its template path and a stamp", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", "nada")
	c, _ := New(root, g)
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(root, OutDir, "x.md"))
	if !markerRE.Match(b) {
		t.Errorf("the compiled page has no marker:\n%s", b)
	}
	if !strings.Contains(string(b), "doct/x.md.tmpl") {
		t.Errorf("the marker does not say where it came from:\n%s", b)
	}
	if markerHashRE.FindSubmatch(b) == nil {
		t.Errorf("the marker carries no inputs stamp:\n%s", b)
	}
}

// The PROTECTION against erasing work: a handwritten `.md` is not overwritten because
// someone created a template with the same name.
func TestBuild_doesNotOverwriteHandwrittenDoc(t *testing.T) {
	t.Run("DTCDC-B03: A handwritten page is never overwritten", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "produto.md.tmpl", "gerado")

	docs := filepath.Join(root, OutDir)
	os.MkdirAll(docs, 0o755)
	byHand := "# Produto\n\nWritten by hand, with content that is in no spec.\n"
	os.WriteFile(filepath.Join(docs, "produto.md"), []byte(byHand), 0o644)

	c, _ := New(root, g)
	res, err := c.Build(false)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := os.ReadFile(filepath.Join(docs, "produto.md"))
	if string(b) != byHand {
		t.Fatalf("the compiler erased handwritten work:\n%s", b)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "produto.md" {
		t.Errorf("the file was skipped in SILENCE: skipped = %v", res.Skipped)
	}
}

// `--dry-run` does not write: it compares without touching the disk.
func TestBuild_dryRunDoesNotWrite(t *testing.T) {
	t.Run("DTCDC-B04: A dry run writes nothing", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", "nada")
	c, _ := New(root, g)
	if _, err := c.Build(true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, OutDir, "x.md")); err == nil {
		t.Error("the dry run wrote to disk")
	}
}

// A broken template FAILS: a syntax error must not become half a document.
func TestBuild_brokenTemplateFails(t *testing.T) {
	t.Run("DTCDC-E03: A broken template fails the build", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", `{{range specs "layer=infra"}}sem end`)
	c, _ := New(root, g)
	if _, err := c.Build(false); err == nil || !strings.Contains(err.Error(), "x.md.tmpl") {
		t.Errorf("a template with no `end` should fail naming it, got %v", err)
	}
}

// `layers` only lists a layer that HAS a spec — a list with empty layers would send the
// reader after sections that do not exist.
func TestLayers_onlyThoseWithSpecs(t *testing.T) {
	t.Run("DTCDC-B07: Only layers that have specs are offered", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	c, _ := New(root, g)
	got := c.fnLayers()
	if len(got) != 1 || got[0] != "infra" {
		t.Errorf("layers = %v, want [infra]", got)
	}
}

// THE WRONG FILTER STOPS THE BUILD. Otherwise the failure would be invisible: the page
// compiles green, and the section that matters simply is not there.
func TestSpecs_wrongFilterFails(t *testing.T) {
	t.Run("DTCDC-E01: A wrong selection filter fails with a message that helps fix it", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	c, _ := New(root, g)

	cases := []struct{ name, filter, wantInMsg string }{
		{"field with a typo", "layar=infra", "unknown field"},
		{"layer that does not exist", "layer=screen", "existing"},
		{"without the `=`", "infra", "`=`"},
		{"code that does not exist", "code=NAOEX", "NAOEX"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := c.fnSpecs(tc.filter)
			if err == nil {
				t.Fatalf("filter %q passed in silence — the section would be empty", tc.filter)
			}
			if !strings.Contains(err.Error(), tc.wantInMsg) {
				t.Errorf("the message does not help fix it: %v", err)
			}
		})
	}
}

// The filter's error climbs to the BUILD and fails it — it does not become half a document.
func TestBuild_wrongFilterFailsTheBuild(t *testing.T) {
	t.Run("DTCDC-E02: A failing selection fails the build without writing the page", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", `{{range specs "layer=naoexiste"}}{{.Code}}{{end}}`)
	c, _ := New(root, g)
	if _, err := c.Build(false); err == nil {
		t.Fatal("the build passed with a filter that matches nothing")
	}
	if _, err := os.Stat(filepath.Join(root, OutDir, "x.md")); err == nil {
		t.Error("an incomplete page was written before failing")
	}
}

// A spec in the MAP and not on disk stops the build: it would leave the documentation
// without a word, and the compiled pages would stay green without it.
func TestNew_specInMapWithoutFileFails(t *testing.T) {
	t.Run("DTCDC-E04: A spec in the map but not on disk fails the compiler", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "pkg/Fantasma.spec.md", Kind: mapx.KindSpec, Code: "FNTSM"},
	}}
	if _, err := New(root, g); err == nil {
		t.Fatal("a spec missing from disk passed in silence")
	}
}

// THE LINK LABEL DOES NOT CARRY THE WAIVER. A rule heading may carry the waiver on its own
// line, and the gate accepts it there on purpose; but the compiled page uses that text as
// the label between brackets, and the result was unreadable. Measured in the reference
// project (#647): 57 occurrences in `docs/regras.md`.
func TestRuleTitleCarriesNoHTMLComment(t *testing.T) {
	t.Run("DTCDC-B10: A rule's title drops the HTML comment on its heading", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{
		"pkg/Com.spec.md": `<!-- @anchors
code: ABCDE
-->

# Com

## Regras

### ABCDE-B01 — a regra diz isto <!-- @no-mark: não há símbolo a marcar -->

O corpo da regra.
`,
	})
	escreveTemplate(t, root, "r.md.tmpl",
		`{{range specs "layer=spec"}}{{range rules .}}- [{{.Code}} — {{.Titulo}}](x)
{{end}}{{end}}`)

	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, OutDir, "r.md"))
	if err != nil {
		t.Fatal(err)
	}
	out := string(b)

	// ONLY THE LINK LINE. The compiled page ALWAYS opens with `<!-- anchors:generated ... -->`,
	// and looking for `<!--` in the whole file would match that header.
	var line string
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, "- [ABCDE-B01") {
			line = l
			break
		}
	}
	if line == "" {
		t.Fatalf("the link line is not in the compiled page:\n%s", out)
	}
	if strings.Contains(line, "<!--") || strings.Contains(line, "-->") {
		t.Errorf("the link label carried the waiver comment:\n  %s", line)
	}
	// AND THE TITLE MUST SURVIVE: cutting the comment must not cut the rule with it.
	if !strings.Contains(line, "a regra diz isto") {
		t.Errorf("the rule's title vanished from the label:\n  %s", line)
	}
}

// The header's date is not an input of the document: the pre-commit dates staged headers
// before the gates run, and hashing the date failed `docs-fresh` on docs just built.
func TestInputsHashIgnoresTheHeaderDate(t *testing.T) {
	t.Run("DTCDC-B13: The header's date is not part of the stamp", func(t *testing.T) {})
	hash := func(raw string) string {
		c := &Compiler{specs: []Spec{{Path: "a.spec.md", raw: raw}}}
		return c.inputsHash([]byte("tmpl"))
	}
	spec := "<!-- @anchors\n  code: AAAAA\n  updated_at: 2026-09-01\n-->\n# A\n\nbody\n"
	base := hash(spec)
	if got := hash(strings.Replace(spec, "2026-09-01", "2026-09-26", 1)); got != base {
		t.Errorf("only the header date changed and the docs were declared out of date")
	}
	if got := hash(strings.Replace(spec, "body", "new body", 1)); got == base {
		t.Error("the content changed and the hash did not")
	}
	// Below the header an `updated_at:` line is content, not the header's date.
	deep := func(date string) string { return spec + strings.Repeat("x\n", 12) + "updated_at: " + date + "\n" }
	if hash(deep("2026-09-01")) == hash(deep("2026-09-26")) {
		t.Error("an updated_at line far below the header is content and must change the hash")
	}
}

// TEMPLATE IN A SUBFOLDER compiles. The shallow `os.ReadDir` skipped `doct/camadas/*.tmpl`
// in silence — half the matrix vanished from the documentation with a green build.
func TestBuild_descendsIntoSubfolder(t *testing.T) {
	t.Run("DTCDC-B05: A template in a subfolder is compiled", func(t *testing.T) {})
	t.Run("DTCDC-I01: Right after a build nothing is stale", func(t *testing.T) {})
	root, g := projetoComFeature(t)
	escreveTemplate(t, root, "top.md.tmpl", `{{range specs}}{{.Code}}{{end}}`)
	escreveTemplate(t, root, "camadas/infra.md"+SufixoTemplate, `{{range specs "layer=infra"}}{{.Code}}{{end}}`)

	c, _ := New(root, g)
	res, err := c.Build(false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res.Written, ",") != "camadas/infra.md,top.md" {
		t.Fatalf("written = %v — the subfolder template was skipped", res.Written)
	}
	b, err := os.ReadFile(filepath.Join(root, OutDir, "camadas", "infra.md"))
	if err != nil {
		t.Fatalf("nothing at docs/camadas/infra.md: %v", err)
	}
	if !strings.Contains(string(b), "GLCGL") {
		t.Errorf("content:\n%s", b)
	}
	// And `Stale` sees the same set: a gate blind to subfolders would never flag the layer
	// pages as out of date.
	if s, _ := c.Stale(); len(s) != 0 {
		t.Errorf("just compiled and Stale reports %v", s)
	}
}

// staleProject builds one template over the example spec and returns the compiler, the
// root and the path of the compiled page.
func staleProject(t *testing.T) (*Compiler, string, string) {
	t.Helper()
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	escreveTemplate(t, root, "x.md.tmpl", `{{range specs}}{{.Code}}{{end}}`)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	return c, root, filepath.Join(root, OutDir, "x.md")
}

func TestStale(t *testing.T) {
	t.Run("DTCDC-B11: A missing page or a page with a different stamp is stale", func(t *testing.T) {
		c, root, page := staleProject(t)
		if s, _ := c.Stale(); len(s) != 0 {
			t.Fatalf("a matching stamp is fresh, got %v", s)
		}
		os.WriteFile(filepath.Join(root, "pkg/GoLive.spec.md"), []byte(specDeExemplo+"\nmore\n"), 0o644)
		c2, _ := New(root, c.Graph)
		if s, _ := c2.Stale(); strings.Join(s, ",") != "x.md" {
			t.Errorf("after a spec change: stale = %v, want [x.md]", s)
		}
		os.Remove(page)
		if s, _ := c2.Stale(); strings.Join(s, ",") != "x.md" {
			t.Errorf("a missing page: stale = %v, want [x.md]", s)
		}
	})
	t.Run("DTCDC-B12: Editing the template makes its page stale", func(t *testing.T) {
		c, root, _ := staleProject(t)
		escreveTemplate(t, root, "x.md.tmpl", `{{range specs}}{{.Titulo}}{{end}}`)
		if s, _ := c.Stale(); strings.Join(s, ",") != "x.md" {
			t.Errorf("after a template edit: stale = %v, want [x.md]", s)
		}
	})
	t.Run("DTCDC-B14: A handwritten page is never reported stale", func(t *testing.T) {
		c, _, page := staleProject(t)
		os.WriteFile(page, []byte("# written by hand\n"), 0o644)
		if s, _ := c.Stale(); len(s) != 0 {
			t.Errorf("a handwritten page was reported stale: %v", s)
		}
	})
	t.Run("DTCDC-X01: Asking for stale pages writes nothing", func(t *testing.T) {
		c, root, page := staleProject(t)
		escreveTemplate(t, root, "x.md.tmpl", `changed`)
		before, _ := os.ReadFile(page)
		if s, _ := c.Stale(); len(s) != 1 {
			t.Fatalf("the page should be stale, got %v", s)
		}
		after, _ := os.ReadFile(page)
		if string(before) != string(after) {
			t.Error("asking for stale pages rewrote the page")
		}
	})
}

func TestUncovered(t *testing.T) {
	t.Run("DTCDC-B15: The specs no template reaches are uncovered", func(t *testing.T) {
		root, g := projetoDeTeste(t, map[string]string{
			"a/Infra.spec.md":  "---\ncode: INFRA\nlayer: infra\n---\n\n# I — i\n",
			"b/Screen.spec.md": "---\ncode: SCRNX\nlayer: screen\n---\n\n# S — s\n",
		})
		escreveTemplate(t, root, "x.md.tmpl", `{{range specs "layer=infra"}}{{.Code}}{{end}}`)
		c, _ := New(root, g)
		got, err := c.Uncovered()
		if err != nil {
			t.Fatal(err)
		}
		if strings.Join(got, ",") != "b/Screen.spec.md" {
			t.Errorf("uncovered = %v, want [b/Screen.spec.md]", got)
		}
	})
}

func TestMarkers(t *testing.T) {
	t.Run("DTCDC-B16: Generated and handwritten markers are recognised", func(t *testing.T) {
		gen := []byte(fmt.Sprintf(GeneratedMarkerHashed, "doct/x.md.tmpl", "0123456789abcdef") + "\n\nx")
		hand := []byte(HandwrittenMarker + "\n# Produto\n")
		plain := []byte("# Nothing said\n")
		if !IsGenerated(gen) || IsHandwritten(gen) {
			t.Error("the generated page was not recognised as generated only")
		}
		if IsGenerated(hand) || !IsHandwritten(hand) {
			t.Error("the handwritten page was not recognised as handwritten only")
		}
		if IsGenerated(plain) || IsHandwritten(plain) {
			t.Error("a page with no marker is neither")
		}
	})
}

func TestNoTemplatesFolder(t *testing.T) {
	t.Run("DTCDC-E05: A project without templates fails the build and has nothing stale or uncovered", func(t *testing.T) {
		root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
		c, _ := New(root, g)
		if _, err := c.Build(false); err == nil || !strings.Contains(err.Error(), "anchors docs init") {
			t.Errorf("build error = %v, want one telling to run docs init", err)
		}
		if s, err := c.Stale(); len(s) != 0 || err != nil {
			t.Errorf("stale = %v, %v; want nothing", s, err)
		}
		if u, err := c.Uncovered(); len(u) != 0 || err != nil {
			t.Errorf("uncovered = %v, %v; want nothing", u, err)
		}
	})
}

func TestSpecByCode(t *testing.T) {
	t.Run("DTCDC-E06: Asking one spec by an unknown code fails", func(t *testing.T) {
		root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
		c, _ := New(root, g)
		if s, err := c.fnSpecByCode("GLCGL"); err != nil || s.Code != "GLCGL" {
			t.Fatalf("GLCGL = %v, %v", s, err)
		}
		if _, err := c.fnSpecByCode("NAOEX"); err == nil || !strings.Contains(err.Error(), "NAOEX") {
			t.Errorf("error = %v, want one naming NAOEX", err)
		}
	})
}

func TestSpecs_onlyFromTheMap(t *testing.T) {
	t.Run("DTCDC-X02: A spec file the map does not list is not loaded", func(t *testing.T) {
		root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
		os.WriteFile(filepath.Join(root, "pkg/Loose.spec.md"),
			[]byte("---\ncode: LOOSE\nlayer: loose\n---\n\n# Loose — l\n"), 0o644)
		c, _ := New(root, g)
		if got := c.fnLayers(); strings.Join(got, ",") != "infra" {
			t.Errorf("layers = %v, want only infra", got)
		}
		if _, err := c.fnSpecByCode("LOOSE"); err == nil {
			t.Error("a spec the map does not list was loaded")
		}
	})
}

// A page compiled before the stamp existed carries the unhashed marker. The fallback
// compiled today's page — which opens with the HASHED marker — and compared the bytes, so
// the marker line alone made every such page stale, current content or not.
func TestStale_unhashedMarkerComparesTheBody(t *testing.T) {
	t.Run("DTCDC-B17: A page with the unhashed marker is stale only when its body differs", func(t *testing.T) {
		c, _, page := staleProject(t)
		b, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		_, body, _ := strings.Cut(string(b), "\n")
		old := fmt.Sprintf(GeneratedMarker, "doct/x.md.tmpl") + "\n" + body
		if err := os.WriteFile(page, []byte(old), 0o644); err != nil {
			t.Fatal(err)
		}
		if s, _ := c.Stale(); len(s) != 0 {
			t.Errorf("an unhashed page with today's body is fresh, got %v", s)
		}
		if err := os.WriteFile(page, []byte(old+"drift\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if s, _ := c.Stale(); strings.Join(s, ",") != "x.md" {
			t.Errorf("an unhashed page with another body is stale, got %v", s)
		}
	})
}

// A container may run layers that have files and no spec yet. The build used to abort on
// them, and leaving them out of the container made the diagram lie by omission.
func TestBuild_layersWithoutSpec(t *testing.T) {
	t.Run("DTCDC-B18: A layer with files and no spec is said, not an error", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	g.Nodes = append(g.Nodes, mapx.Node{ID: "web/a.tsx", Kind: mapx.KindCode, Layer: "presentation"},
		mapx.Node{ID: "web/b.tsx", Kind: mapx.KindCode, Layer: "presentation"})
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	c.Config = &config.Config{
		Layers: map[string]config.Layer{"dao": {}},
		ContainersDecl: []config.Container{
			{Name: "web", Layers: []string{"presentation", "dao"}},
			{Name: "api", Layers: []string{"infra", "presentation"}},
		},
	}
	if got, err := c.fnSpecs("layer=dao"); err != nil || len(got) != 0 {
		t.Errorf("a declared layer with no spec selects nothing, got %v (%v)", got, err)
	}
	if got, err := c.fnSpecs("layer=presentation"); err != nil || len(got) != 0 {
		t.Errorf("a layer with files in the map selects nothing, got %v (%v)", got, err)
	}
	if n := c.fnLayerFiles("presentation"); n != 2 {
		t.Errorf("presentation has 2 files, got %d", n)
	}
	if (&Compiler{}).fnLayerFiles("x") != 0 {
		t.Error("no map, no file")
	}
	c.InitScaffolds(false)
	if _, err := c.Build(false); err != nil {
		t.Fatalf("the build must not abort on a layer with no spec: %v", err)
	}
	b, _ := readFile(c.Root, "architecture.md")
	for _, want := range []string{
		"- **presentation** — 2 file(s)", "- **dao** — 0 file(s)", "_No layer of this container has a spec yet:_",
		"_No spec in this layer — 2 governed file(s)._", `presentation_nospec["2 file(s), no spec"]`,
	} {
		if !strings.Contains(b, want) {
			t.Errorf("the page lacks %q:\n%s", want, b)
		}
	}
}

func TestNewWith_readsThroughTheSource(t *testing.T) {
	t.Run("DTCDC-B19: The compiler reads through its source", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	var asked []string
	c, err := NewWith(root, g, func(rel string) ([]byte, error) {
		asked = append(asked, rel)
		return []byte(strings.Replace(specDeExemplo, "GoLiveChecklist", "FromTheIndex", 1)), nil
	})
	if err != nil || !strings.HasPrefix(c.specs[0].Titulo, "FromTheIndex") || asked[0] != "pkg/GoLive.spec.md" {
		t.Fatalf("the spec is read through the source, got %v %+v %v", err, c.specs, asked)
	}
	if b, err := c.readRel("x/y.md"); err != nil || !strings.Contains(string(b), "FromTheIndex") {
		t.Errorf("every read goes through the source, got %q %v", b, err)
	}
	plain, _ := New(root, g)
	if !strings.HasPrefix(plain.specs[0].Titulo, "GoLiveChecklist") {
		t.Errorf("without a source the tree is read, got %q", plain.specs[0].Titulo)
	}
}

func TestCompiler_templatesThroughTheSource(t *testing.T) {
	t.Run("DTCDC-B20: A compiler with a source compiles the templates it holds", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	for rel, body := range map[string]string{Dir + "/a.md.tmpl": "tree template\n", Dir + "/b.md.tmpl": "only on disk\n"} {
		if err := os.MkdirAll(filepath.Join(root, Dir), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, rel), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	c, err := NewWith(root, g, func(rel string) ([]byte, error) {
		switch rel {
		case Dir + "/a.md.tmpl":
			return []byte("index template\n"), nil
		case "pkg/GoLive.spec.md":
			return []byte(specDeExemplo), nil
		}
		return nil, os.ErrNotExist
	})
	if err != nil {
		t.Fatal(err)
	}
	pages, err := c.Compiled()
	if err != nil || len(pages) != 1 || pages[0].Path != OutDir+"/a.md" || !strings.Contains(string(pages[0].Content), "index template") {
		t.Errorf("only the source's template, as the source holds it, got %+v %v", pages, err)
	}
}

func TestCompiler_compiledPages(t *testing.T) {
	t.Run("DTCDC-B21: The pages out of date, compiled without writing", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pkg/GoLive.spec.md": specDeExemplo})
	if err := os.MkdirAll(filepath.Join(root, Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"fresh", "stale", "hand"} {
		if err := os.WriteFile(filepath.Join(root, Dir, name+".md.tmpl"), []byte(name+" {{range specs}}{{.Code}}{{end}}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, OutDir, "hand.md"), []byte("by hand\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, Dir, "stale.md.tmpl"), []byte("stale, revised {{range specs}}{{.Code}}{{end}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(root, OutDir, "stale.md"))
	pages, err := c.Compiled()
	if err != nil || len(pages) != 1 || pages[0].Path != OutDir+"/stale.md" || !strings.Contains(string(pages[0].Content), "stale, revised GLCGL") {
		t.Errorf("only the page out of date, compiled, got %+v %v", pages, err)
	}
	if after, _ := os.ReadFile(filepath.Join(root, OutDir, "stale.md")); string(after) != string(before) {
		t.Error("nothing is written")
	}
}
