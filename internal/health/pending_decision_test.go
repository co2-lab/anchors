// @anchors
//   ref: PNDCP

package health

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// specWithOpenDecisions is a spec whose open decisions section lists n questions.
func specWithOpenDecisions(n int) string {
	s := "# Unit\n\n## Decisões em aberto\n\n| Código | Pergunta | Quem decide |\n| --- | --- | --- |\n"
	for i := 1; i <= n; i++ {
		s += fmt.Sprintf("| `PARCX-Q%02d` | Question %d? | Product |\n", i, i)
	}
	return s
}

// The doctor listed a missing mutation signal as a point of attention while declaring
// "0 points" with a rule waiting for a decision. The hierarchy was upside down: what is
// still to be MEASURED showed, what is still to be DECIDED did not — and the decision is
// what most needs to survive the session, because it depends on a person and can take weeks.
func TestDoctorReportsPendingDecision(t *testing.T) {
	t.Run("PNDCP-B01: A spec with two open decisions gives one warning carrying the count", func(t *testing.T) {})
	t.Run("PNDCP-B02: A section closed with none is not pending", func(t *testing.T) {})
	dir := t.TempDir()
	spec := "# Unidade\n\n## Decisões em aberto\n\n| Código | Pergunta | Quem decide |\n| --- | --- | --- |\n" +
		"| `PARCX-Q01` | Fuso do vencimento? | Produto |\n| `PARCX-Q02` | Retenção dos logs? | Segurança |\n"
	if err := os.WriteFile(filepath.Join(dir, "u.spec.md"), []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}}}

	fs := checkPendingDecisions(g, dir, nil)

	if len(fs) != 1 {
		t.Fatalf("expected 1 finding (the spec with open decisions), got %d", len(fs))
	}
	if fs[0].Check != "decisao-pendente" || fs[0].Severity != Warn {
		t.Errorf("wrong finding: %+v", fs[0])
	}
	if !strings.Contains(fs[0].Detail, "2 decisão") && !strings.Contains(fs[0].Detail, "2 decision") {
		t.Errorf("it should count both questions: %s", fs[0].Detail)
	}

	// A spec that CLOSED the section with "none" is not pending — it is the honest opt-out.
	closed := "# Unidade\n\n## Decisões em aberto\n\nnenhuma\n"
	if err := os.WriteFile(filepath.Join(dir, "u.spec.md"), []byte(closed), 0o644); err != nil {
		t.Fatal(err)
	}
	if fs := checkPendingDecisions(g, dir, nil); len(fs) != 0 {
		t.Errorf("a section closed with `none` is not pending: %+v", fs)
	}
}

func TestPendingDecisions_heaviestFirst(t *testing.T) {
	t.Run("PNDCP-B03: The spec with the most open decisions comes first", func(t *testing.T) {})
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "light.spec.md"), []byte(specWithOpenDecisions(1)), 0o644)
	os.WriteFile(filepath.Join(dir, "heavy.spec.md"), []byte(specWithOpenDecisions(3)), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "light.spec.md", Kind: mapx.KindSpec},
		{ID: "heavy.spec.md", Kind: mapx.KindSpec},
	}}
	fs := checkPendingDecisions(g, dir, nil)
	if len(fs) != 2 || fs[0].Subject != "heavy.spec.md" || fs[1].Subject != "light.spec.md" {
		t.Fatalf("the heaviest spec must come first: %+v", fs)
	}
}

func TestPendingDecisions_nothingToReport(t *testing.T) {
	t.Run("PNDCP-B04: A nil map or a map without open decisions gives nothing", func(t *testing.T) {})
	if fs := checkPendingDecisions(nil, t.TempDir(), nil); fs != nil {
		t.Fatalf("a nil map gives nothing: %+v", fs)
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "notes.md"), []byte(specWithOpenDecisions(2)), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "notes.md", Kind: mapx.KindDoc}}}
	if fs := checkPendingDecisions(g, dir, nil); fs != nil {
		t.Fatalf("only specs are read: %+v", fs)
	}
}

func TestPendingDecisions_sameCountAsTheCheck(t *testing.T) {
	t.Run("PNDCP-I01: The doctor's count is the check's count", func(t *testing.T) {})
	t.Run("PNDCP-X01: The count follows the check's rule, not a reading of its own", func(t *testing.T) {})
	dir := t.TempDir()
	body := specWithOpenDecisions(3)
	os.WriteFile(filepath.Join(dir, "u.spec.md"), []byte(body), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}}}
	fs := checkPendingDecisions(g, dir, nil)
	want := i18n.T("health.pending_decisions", gate.OpenDecisions(body, nil, ""))
	if len(fs) != 1 || fs[0].Detail != want {
		t.Fatalf("the doctor's count must be the check's: got %+v, want detail %q", fs, want)
	}

	closed := "# Unit\n\n## Decisões em aberto\n\nnenhuma\n"
	os.WriteFile(filepath.Join(dir, "u.spec.md"), []byte(closed), 0o644)
	if n := gate.OpenDecisions(closed, nil, ""); n != 0 {
		t.Fatalf("setup: the check reads %d open decisions in a closed section", n)
	}
	if fs := checkPendingDecisions(g, dir, nil); len(fs) != 0 {
		t.Fatalf("where the check counts zero the doctor reports nothing: %+v", fs)
	}
}

func TestPendingDecisions_unreadableSpecIsSkipped(t *testing.T) {
	t.Run("PNDCP-E01: A spec missing on disk is skipped and the others are still reported", func(t *testing.T) {})
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "here.spec.md"), []byte(specWithOpenDecisions(1)), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "gone.spec.md", Kind: mapx.KindSpec},
		{ID: "here.spec.md", Kind: mapx.KindSpec},
	}}
	fs := checkPendingDecisions(g, dir, nil)
	if len(fs) != 1 || fs[0].Subject != "here.spec.md" {
		t.Fatalf("only the readable spec is reported: %+v", fs)
	}
}
