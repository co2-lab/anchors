package doct

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestInstructionFor_knownKinds(t *testing.T) {
	t.Run("DCKND-B01: Every known kind has a title, items to cover and a trap", func(t *testing.T) {
		kinds := KnownKinds()
		want := []string{"openapi", "c4", "schema", "components", "adr", "runbook"}
		if strings.Join(kinds, ",") != strings.Join(want, ",") {
			t.Fatalf("KnownKinds() = %v, want %v", kinds, want)
		}
		for _, k := range kinds {
			i := InstructionFor(config.DocArtifact{Kind: k})
			if strings.TrimSpace(i.Titulo) == "" || len(i.Pede) == 0 || strings.TrimSpace(i.Armadilha) == "" {
				t.Errorf("kind %q is incomplete: %+v", k, i)
			}
		}
	})
	t.Run("DCKND-I01: The known kinds and the instructions are the same set", func(t *testing.T) {
		for _, k := range KnownKinds() {
			if i := InstructionFor(config.DocArtifact{Kind: k}); i.Titulo == k {
				t.Errorf("kind %q fell back to the minimal instruction", k)
			}
		}
		if len(instructions) != len(KnownKinds()) {
			t.Errorf("%d instructions for %d known kinds", len(instructions), len(KnownKinds()))
		}
	})
}

func TestInstructionFor_lookup(t *testing.T) {
	t.Run("DCKND-B02: The kind is found ignoring case and surrounding spaces", func(t *testing.T) {
		got := InstructionFor(config.DocArtifact{Kind: " OpenAPI "})
		if got.Titulo != instructions[config.KindOpenAPI].Titulo {
			t.Errorf("\" OpenAPI \" gave %q, want the openapi instruction", got.Titulo)
		}
	})
	t.Run("DCKND-B03: An unknown kind gets a minimal instruction titled by the kind or the path", func(t *testing.T) {
		a := InstructionFor(config.DocArtifact{Kind: "glossary", Path: "docs/g.md"})
		b := InstructionFor(config.DocArtifact{Path: "docs/x.md"})
		if a.Titulo != "glossary" || b.Titulo != "docs/x.md" {
			t.Errorf("titles = %q, %q; want the kind, then the path", a.Titulo, b.Titulo)
		}
		for _, i := range []Instruction{a, b} {
			if len(i.Pede) != 0 || i.Armadilha != "" {
				t.Errorf("a minimal instruction has no items nor trap: %+v", i)
			}
		}
	})
}

func TestDuty(t *testing.T) {
	t.Run("DCKND-B04: The duty text lists title and path, why, items and trap in order", func(t *testing.T) {
		d := config.DocArtifact{Kind: "openapi", Path: "docs/api.yaml", Why: "the mobile app reads it"}
		got := Duty(d)
		i := instructions[config.KindOpenAPI]
		lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		if len(lines) != 2+len(i.Pede)+1 {
			t.Fatalf("got %d lines, want title + why + %d items + trap:\n%s", len(lines), len(i.Pede), got)
		}
		if strings.TrimSpace(lines[0]) != i.Titulo+" — `docs/api.yaml`" {
			t.Errorf("first line = %q", lines[0])
		}
		if strings.TrimSpace(lines[1]) != "the mobile app reads it" {
			t.Errorf("second line = %q, want the why", lines[1])
		}
		for n, p := range i.Pede {
			if strings.TrimSpace(lines[2+n]) != "· "+p {
				t.Errorf("item line %d = %q", n, lines[2+n])
			}
		}
		if !strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "ARMADILHA: ") {
			t.Errorf("last line = %q, want the trap", lines[len(lines)-1])
		}
		// Without a why and for an unknown kind, only the title line remains.
		if got := Duty(config.DocArtifact{Kind: "glossary", Path: "g.md"}); strings.Count(got, "\n") != 1 {
			t.Errorf("an unknown kind with no why should be one line, got %q", got)
		}
	})
	t.Run("DCKND-X01: No documentation kind is refused", func(t *testing.T) {
		if got := Duty(config.DocArtifact{Kind: "never-heard-of", Path: "docs/n.md"}); !strings.Contains(got, "docs/n.md") {
			t.Errorf("the duty of an unknown kind should still name the document, got %q", got)
		}
	})
}
