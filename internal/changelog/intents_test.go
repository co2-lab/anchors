// @anchors
//   code: INTNC
//   ref: INTNA

package changelog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const anIntent = "---\nproduct: app\nbump: minor\ntype: feat\n---\nReceipts can now be shared from the list.\n"

func TestIntent_parse(t *testing.T) {
	t.Run("INTNA-B01: An intent reads its product, bump and type from the front matter, and its text after it", func(t *testing.T) {
		in, err := ParseIntent("share.md", []byte(strings.ReplaceAll(anIntent, "\n", "\r\n")))
		if err != nil || in.File != "share.md" || in.Product != "app" || in.Bump != "minor" || in.Type != "feat" || in.Text != "Receipts can now be shared from the list." {
			t.Errorf("%+v %v", in, err)
		}
	})
	t.Run("INTNA-E01: An intent with no front matter, one not closed, or a key it does not know is an error naming the file", func(t *testing.T) {
		for body, says := range map[string]string{
			"Receipts.\n":                            "opens with a --- front matter",
			"---\nproduct: app\nReceipts.\n":         "not closed",
			"---\nproduct: app\nscope: ui\n---\nx\n": "scope",
			"---\nproduct: [app\n---\nx\n":           "the front matter",
		} {
			if _, err := ParseIntent("x.md", []byte(body)); err == nil || !strings.Contains(err.Error(), "x.md") || !strings.Contains(err.Error(), says) {
				t.Errorf("%q: %v, want %q", body, err, says)
			}
		}
	})
}

func TestIntent_read(t *testing.T) {
	t.Run("INTNA-B02: A directory's intents are read sorted by name, each one that does not parse is an error of its own, and no directory holds none", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "b.md"), []byte(anIntent), 0o644)
		os.WriteFile(filepath.Join(dir, "a.md"), []byte(strings.Replace(anIntent, "minor", "patch", 1)), 0o644)
		os.WriteFile(filepath.Join(dir, "c.md"), []byte("no front matter"), 0o644)
		os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o644)
		os.MkdirAll(filepath.Join(dir, "released", "1.0.0"), 0o755)
		os.WriteFile(filepath.Join(dir, "released", "1.0.0", "old.md"), []byte(anIntent), 0o644)
		got, errs := ReadIntents(dir)
		if len(got) != 2 || got[0].File != "a.md" || got[1].File != "b.md" || len(errs) != 1 || !strings.Contains(errs[0].Error(), "c.md") {
			t.Errorf("%+v %v", got, errs)
		}
		if got, errs := ReadIntents(filepath.Join(dir, "nope")); len(got) != 0 || len(errs) != 0 {
			t.Errorf("no directory: %v %v", got, errs)
		}
	})
}

func TestIntent_render(t *testing.T) {
	t.Run("INTNA-B03: Intents are grouped by type in the order declared, another type after them, and rendered with a heading per type", func(t *testing.T) {
		intents := []Intent{{Type: "fix", Text: "A crash on login."}, {Type: "perf", Text: "Faster lists."}, {Type: "feat", Text: "Sharing."}}
		r := GroupIntents("", "", intents, []string{"feat", "fix", "docs"})
		if len(r.Sections) != 3 || r.Sections[0].Type != "feat" || r.Sections[1].Type != "fix" || r.Sections[2].Type != "perf" {
			t.Fatalf("%+v", r.Sections)
		}
		r.Candidate = "1.1.0"
		tr := func(k string) string {
			return map[string]string{"changelog.unreleased": "Unreleased", "changelog.type.feat": "Features", "changelog.type.fix": "Fixes"}[k]
		}
		got, err := RenderIntents(DefaultIntentsTemplate, r, tr)
		want := "## 1.1.0 — Unreleased\n\n### Features\n\n- Sharing.\n\n### Fixes\n\n- A crash on login.\n\n### \n\n- Faster lists.\n"
		if err != nil || got != want {
			t.Errorf("%v\n%q\nwant\n%q", err, got, want)
		}
		if _, err := RenderIntents("{{", r, tr); err == nil || !strings.Contains(err.Error(), "changelog template") {
			t.Errorf("a template that does not parse: %v", err)
		}
	})
}
