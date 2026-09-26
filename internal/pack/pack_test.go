package pack

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const lgpdFixture = `name: lgpd
domain: privacy
jurisdiction: br
authority: "Lei 13.709/2018"
requires: [erasure_handler]
obligations:
  - name: lgpd-erasure
    article: "Art. 18, VI"
    when: "carries: personal-data"
    must_appear_in: ["{{erasure_handler}}", "{{ audit_log }}/x.go", "static.go"]
`

const gdprFixture = `name: gdpr
domain: privacy
jurisdiction: EU
obligations:
  - name: gdpr-erasure
    when: "carries: personal-data"
    must_appear_in: ["{{erasure_handler}}"]
`

const wcagFixture = `name: wcag
domain: accessibility
jurisdiction: global
obligations:
  - name: wcag-alt
    when: "renders: image"
    must_appear_in: ["src/a11y.ts"]
`

func writePack(t *testing.T, root, rel, body string) string {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad(t *testing.T) {
	root := t.TempDir()
	p, err := Load(writePack(t, root, "packs/privacy/lgpd.yaml", lgpdFixture))
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "lgpd" || p.Domain != "privacy" || p.Jurisdiction != "br" ||
		p.Obligations[0].Article != "Art. 18, VI" || !reflect.DeepEqual(p.Requires, []string{"erasure_handler"}) {
		t.Fatalf("Load = %+v", p)
	}

	for name, tc := range map[string]struct{ body, want string }{
		"no name":        {"domain: x\nobligations:\n  - name: a\n", "without `name`"},
		"no obligations": {"name: empty\n", "empty: pack without obligations"},
		"bad yaml":       {"name: [unclosed\n", "bad.yaml:"},
	} {
		_, err := Load(writePack(t, root, "bad.yaml", tc.body))
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: Load error = %v, want it to contain %q", name, err, tc.want)
		}
	}
	if _, err := Load(filepath.Join(root, "missing.yaml")); !os.IsNotExist(err) {
		t.Errorf("a missing file must surface the read error, got %v", err)
	}
}

func TestResolveRef(t *testing.T) {
	for ref, want := range map[string]string{
		"privacy/lgpd":          filepath.Join("/r", "packs", "privacy", "lgpd.yaml"),
		"./internal/mine.yaml":  filepath.Join("/r", "internal", "mine.yaml"),
		"custom/other.yml":      filepath.Join("/r", "custom", "other.yml"),
		"/abs/elsewhere/p.yaml": "/abs/elsewhere/p.yaml",
	} {
		if got := resolveRef("/r", ref); got != want {
			t.Errorf("resolveRef(%q) = %q, want %q", ref, got, want)
		}
	}
}

func TestLoadAll_resolvesPlaceholdersAndSorts(t *testing.T) {
	root := t.TempDir()
	writePack(t, root, "packs/privacy/lgpd.yaml", lgpdFixture)
	writePack(t, root, "packs/accessibility/wcag.yaml", wcagFixture)
	values := map[string]string{"erasure_handler": "api/erase.go", "audit_log": "api/audit"}

	packs, avisos, err := LoadAll(root, []string{"privacy/lgpd", "accessibility/wcag"}, values, []string{" BR "})
	if err != nil {
		t.Fatal(err)
	}
	if len(avisos) != 0 {
		t.Fatalf("no warnings expected, got %v", avisos)
	}
	if len(packs) != 2 || packs[0].Name != "lgpd" || packs[1].Name != "wcag" {
		t.Fatalf("LoadAll must return the packs sorted by name, got %v", packs)
	}
	got := packs[0].Obligations[0].MustAppearIn
	if want := []string{"api/erase.go", "api/audit/x.go", "static.go"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("placeholders resolved to %v, want %v", got, want)
	}
}

func TestLoadAll_skipsUndeclaredJurisdictionWithWarning(t *testing.T) {
	root := t.TempDir()
	writePack(t, root, "packs/privacy/gdpr.yaml", gdprFixture)
	writePack(t, root, "packs/accessibility/wcag.yaml", wcagFixture)
	values := map[string]string{"erasure_handler": "x.go"}

	packs, avisos, err := LoadAll(root, []string{"privacy/gdpr", "accessibility/wcag"}, values, []string{"br"})
	if err != nil {
		t.Fatal(err)
	}
	if len(packs) != 1 || packs[0].Name != "wcag" {
		t.Fatalf("only the global pack applies to a br-only project, got %v", packs)
	}
	if len(avisos) != 1 || !strings.Contains(avisos[0], "pack `gdpr` belongs to jurisdiction `EU`") {
		t.Fatalf("the skipped pack must be warned about, got %v", avisos)
	}

	// No jurisdictions declared: no filter at all.
	packs, avisos, err = LoadAll(root, []string{"privacy/gdpr"}, values, nil)
	if err != nil || len(packs) != 1 || len(avisos) != 0 {
		t.Fatalf("without jurisdictions every pack loads: %v %v %v", packs, avisos, err)
	}
}

func TestLoadAll_refusesUnresolvedPlaceholders(t *testing.T) {
	root := t.TempDir()
	writePack(t, root, "packs/privacy/lgpd.yaml", lgpdFixture)

	_, _, err := LoadAll(root, []string{"privacy/lgpd"}, map[string]string{}, nil)
	if err == nil {
		t.Fatal("a pack whose placeholders the project does not declare must not load")
	}
	if !strings.Contains(err.Error(), "pack `lgpd` needs the project to declare `audit_log`, `erasure_handler`") {
		t.Fatalf("the error must list the missing values, sorted: %v", err)
	}

	if _, _, err := LoadAll(root, []string{"privacy/missing"}, nil, nil); err == nil {
		t.Fatal("a missing pack must be an error")
	}
}
