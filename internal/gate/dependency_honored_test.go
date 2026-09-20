package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestBacktickedSymbols(t *testing.T) {
	t.Run("DEPHN-X02: The gate does not interpret dependency semantics or parameter signatures", func(t *testing.T) {})
	// o scan PRESERVA as crases do campo Método — elas são o sinal de "isto é um
	// SÍMBOLO" (verificável), separando-o da prosa descritiva. Formas com as bordas
	// já removidas seguem aceitas (mapas gravados antes dessa mudança).
	cases := map[string][]string{
		"`resolveVersion`":              {"resolveVersion"},              // forma canônica
		"`resolveVersion`, `applyEdit`": {"resolveVersion", "applyEdit"}, // 2 símbolos
		"resolveVersion":                nil,                             // SEM crase = prosa, não contrato
		"resolveVersion`, `applyEdit":   {"resolveVersion", "applyEdit"}, // bordas antigas
		"extractCaller`, `parseBody":    {"extractCaller", "parseBody"},  // bordas antigas
		"CRUD + consultas":              nil,                             // prosa → nada
		"CRUD de member/grupo":          nil,                             // prosa
		"":                              nil,                             // vazio
		"`RecurrenceRuleLike` (tipo)":   {"RecurrenceRuleLike"},          // símbolo + anotação
	}
	for in, want := range cases {
		got := backtickedSymbols(in)
		if len(got) != len(want) {
			t.Errorf("backtickedSymbols(%q) = %v, quer %v", in, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("backtickedSymbols(%q)[%d] = %q, quer %q", in, i, got[i], want[i])
			}
		}
	}
}

func TestSymbolUsed(t *testing.T) {
	t.Run("DEPHN-I02: Symbol presence is matched strictly on token word boundaries", func(t *testing.T) {})
	code := "import { resolveVersion } from './x'\nconst v = resolveVersion(h, m)\n"
	if !symbolUsed(code, "resolveVersion") {
		t.Error("resolveVersion deveria ser usado")
	}
	if symbolUsed(code, "applyEdit") {
		t.Error("applyEdit NÃO está no código")
	}
	// fronteira de palavra: 'resolve' não casa 'resolveVersion'
	if symbolUsed(code, "resolve") {
		t.Error("'resolve' não deveria casar 'resolveVersion' (fronteira)")
	}
}

func TestBacktickedSymbols_prosaNaoEhContrato(t *testing.T) {
	t.Run("DEPHN-I01: Prose descriptions in dependency methods are never treated as contracts", func(t *testing.T) {})
	// células SEM crase são descrição livre — não viram promessa (o furo que
	// transformava "aportes"/"solicitações" em símbolo prometido).
	for _, prose := range []string{
		"solicitações", "aportes", "parcelas", "CRUD de IRDependent",
		"contagem de orçamentos", "leitura/atualização",
	} {
		if got := backtickedSymbols(prose); got != nil {
			t.Errorf("backtickedSymbols(%q) = %v, quer nil (prosa não é contrato)", prose, got)
		}
	}
	// com crase, continua sendo contrato
	if got := backtickedSymbols("`pingHeartbeat`"); len(got) != 1 || got[0] != "pingHeartbeat" {
		t.Errorf("símbolo crasado deveria ser extraído, got %v", got)
	}
}

func TestNearestSymbol(t *testing.T) {
	t.Run("DEPHN-I03: Near-symbol rename suggestions are strictly conservative", func(t *testing.T) {})
	code := "import { pingHeartbeat } from './x'\nconst r = computeSeatsAmountCents(a)\n"
	if got := nearestSymbol(code, "ping"); got != "pingHeartbeat" {
		t.Errorf("nearestSymbol(ping) = %q, quer pingHeartbeat", got)
	}
	if got := nearestSymbol(code, "computeSeatsAmount"); got != "computeSeatsAmountCents" {
		t.Errorf("nearestSymbol(computeSeatsAmount) = %q, quer computeSeatsAmountCents", got)
	}
	// sem parentesco → sem palpite
	if got := nearestSymbol(code, "totalmenteOutro"); got != "" {
		t.Errorf("esperava sem sugestão, got %q", got)
	}
}

func TestDependencyHonored(t *testing.T) {
	t.Run("DEPHN-B01: An artifact that is not a spec leaves without a verdict", func(t *testing.T) {
		v, msg := checkDependencyHonored("package main", mapx.Node{Kind: mapx.KindCode}, "", nil, nil)
		if v != Skip {
			t.Errorf("expected Skip for non-spec node, got %v (%s)", v, msg)
		}
	})

	t.Run("DEPHN-B02: Without a relational map the verdict is undetermined", func(t *testing.T) {
		v, msg := checkDependencyHonored("# spec", mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if v != Pending {
			t.Errorf("expected Pending when graph is nil, got %v (%s)", v, msg)
		}
	})

	t.Run("DEPHN-B03: A spec declaring no confrontable symbols leaves without a verdict", func(t *testing.T) {
		g := &mapx.Graph{
			Edges: []mapx.Edge{
				{From: "s.spec.md", To: "dep.go", Type: mapx.EdgeDependsOn, Method: "CRUD + queries"},
			},
		}
		v, msg := checkDependencyHonored("# spec", mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, "", g, nil)
		if v != Skip {
			t.Errorf("expected Skip when no confrontable symbols are declared, got %v (%s)", v, msg)
		}
	})

	t.Run("DEPHN-B04: A spec governing no code leaves the verdict undetermined", func(t *testing.T) {
		g := &mapx.Graph{
			Edges: []mapx.Edge{
				{From: "s.spec.md", To: "dep.go", Type: mapx.EdgeDependsOn, Method: "`resolveVersion`"},
			},
		}
		v, msg := checkDependencyHonored("# spec", mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, "", g, nil)
		if v != Pending {
			t.Errorf("expected Pending when no code is specified, got %v (%s)", v, msg)
		}
	})

	t.Run("DEPHN-B05: When every promised symbol appears in governed code, the gate passes", func(t *testing.T) {
		root := t.TempDir()
		codeFile := "app.go"
		content := "package main\n\nfunc Run() {\n\tresolveVersion()\n\tapplyEdit()\n}\n"
		if err := os.WriteFile(filepath.Join(root, codeFile), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		g := &mapx.Graph{
			Edges: []mapx.Edge{
				{From: "s.spec.md", To: "dep.go", Type: mapx.EdgeDependsOn, Method: "`resolveVersion`, `applyEdit`"},
				{From: "s.spec.md", To: codeFile, Type: mapx.EdgeSpecifies},
			},
		}
		v, msg := checkDependencyHonored("# spec", mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, root, g, nil)
		if v != Pass {
			t.Fatalf("expected Pass when all symbols are used, got %v (%s)", v, msg)
		}
		if msg != "" {
			t.Errorf("expected empty message on Pass, got %s", msg)
		}
	})

	t.Run("DEPHN-B06: When a promised symbol is absent from governed code, the gate fails", func(t *testing.T) {
		root := t.TempDir()
		codeFile := "app.go"
		content := "package main\n\nfunc Run() {\n\tresolveVersion()\n}\n"
		if err := os.WriteFile(filepath.Join(root, codeFile), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		g := &mapx.Graph{
			Edges: []mapx.Edge{
				{From: "s.spec.md", To: "dep.go", Type: mapx.EdgeDependsOn, Method: "`resolveVersion`, `applyEdit`"},
				{From: "s.spec.md", To: codeFile, Type: mapx.EdgeSpecifies},
			},
		}
		v, msg := checkDependencyHonored("# spec", mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, root, g, nil)
		if v != Fail {
			t.Fatalf("expected Fail when promised symbol is missing, got %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "applyEdit") {
			t.Errorf("expected failure message to cite missing symbol 'applyEdit', got %s", msg)
		}
		if !strings.Contains(msg, "dep.go") {
			t.Errorf("expected failure message to cite target file 'dep.go', got %s", msg)
		}
	})

	t.Run("DEPHN-B07: When an absent symbol resembles an identifier in code, the verdict suggests the rename", func(t *testing.T) {
		root := t.TempDir()
		codeFile := "app.go"
		content := "package main\n\nfunc Run() {\n\tpingHeartbeat()\n}\n"
		if err := os.WriteFile(filepath.Join(root, codeFile), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		g := &mapx.Graph{
			Edges: []mapx.Edge{
				{From: "s.spec.md", To: "dep.go", Type: mapx.EdgeDependsOn, Method: "`ping`"},
				{From: "s.spec.md", To: codeFile, Type: mapx.EdgeSpecifies},
			},
		}
		v, msg := checkDependencyHonored("# spec", mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, root, g, nil)
		if v != Fail {
			t.Fatalf("expected Fail when promised symbol is missing, got %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "pingHeartbeat") {
			t.Errorf("expected rename suggestion to mention 'pingHeartbeat', got %s", msg)
		}
	})

	t.Run("DEPHN-B08: Symbols appearing only in comments do not fulfill the promise", func(t *testing.T) {
		root := t.TempDir()
		codeFile := "app.go"
		content := "package main\n\n// resolveVersion is not called\nfunc Run() {}\n"
		if err := os.WriteFile(filepath.Join(root, codeFile), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		g := &mapx.Graph{
			Edges: []mapx.Edge{
				{From: "s.spec.md", To: "dep.go", Type: mapx.EdgeDependsOn, Method: "`resolveVersion`"},
				{From: "s.spec.md", To: codeFile, Type: mapx.EdgeSpecifies},
			},
		}
		v, msg := checkDependencyHonored("# spec", mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, root, g, nil)
		if v != Fail {
			t.Fatalf("expected Fail when symbol appears only in comments, got %v (%s)", v, msg)
		}
	})

	t.Run("DEPHN-I01: Prose descriptions in dependency methods are never treated as contracts", func(t *testing.T) {
		if got := backtickedSymbols("CRUD + queries"); got != nil {
			t.Errorf("expected nil for prose, got %v", got)
		}
		if got := backtickedSymbols("`pingHeartbeat`"); len(got) != 1 || got[0] != "pingHeartbeat" {
			t.Errorf("expected ['pingHeartbeat'], got %v", got)
		}
	})

	t.Run("DEPHN-I02: Symbol presence is matched strictly on token word boundaries", func(t *testing.T) {
		code := "const resolveVersion = 1"
		if symbolUsed(code, "resolve") {
			t.Errorf("symbolUsed matched substring 'resolve' inside 'resolveVersion'")
		}
		if !symbolUsed(code, "resolveVersion") {
			t.Errorf("symbolUsed failed to match exact token 'resolveVersion'")
		}
	})

	t.Run("DEPHN-I03: Near-symbol rename suggestions are strictly conservative", func(t *testing.T) {
		code := "const pingHeartbeat = 1"
		if got := nearestSymbol(code, "pin"); got != "" {
			t.Errorf("expected empty suggestion for symbol shorter than 4 chars, got %q", got)
		}
		if got := nearestSymbol(code, "ping"); got != "pingHeartbeat" {
			t.Errorf("expected 'pingHeartbeat', got %q", got)
		}
		if got := nearestSymbol(code, "unrelated"); got != "" {
			t.Errorf("expected empty suggestion for unrelated identifier, got %q", got)
		}
	})

	t.Run("DEPHN-X01: The gate performs static textual confrontation without runtime execution", func(t *testing.T) {
		root := t.TempDir()
		codeFile := "app.go"
		content := "package main\n\nfunc processData() {}\n"
		if err := os.WriteFile(filepath.Join(root, codeFile), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		g := &mapx.Graph{
			Edges: []mapx.Edge{
				{From: "s.spec.md", To: "dep.go", Type: mapx.EdgeDependsOn, Method: "`processData`"},
				{From: "s.spec.md", To: codeFile, Type: mapx.EdgeSpecifies},
			},
		}
		v, msg := checkDependencyHonored("# spec", mapx.Node{ID: "s.spec.md", Kind: mapx.KindSpec}, root, g, nil)
		if v != Pass {
			t.Fatalf("expected Pass for static token occurrence without execution, got %v (%s)", v, msg)
		}
	})

	t.Run("DEPHN-X02: The gate does not interpret dependency semantics or parameter signatures", func(t *testing.T) {
		symbols := backtickedSymbols("`RecurrenceRuleLike` (type)")
		if len(symbols) != 1 || symbols[0] != "RecurrenceRuleLike" {
			t.Errorf("expected single symbol 'RecurrenceRuleLike' ignoring '(type)', got %v", symbols)
		}
	})
}
