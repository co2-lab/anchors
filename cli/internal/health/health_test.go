package health

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func hasFinding(r Report, check, subject string) bool {
	for _, f := range r.Findings {
		if f.Check == check && f.Subject == subject {
			return true
		}
	}
	return false
}

// projeto de brinquedo em disco + grafo coerente com ele, exceto pelos problemas
// que queremos que o doctor pegue.
func setup(t *testing.T) (string, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	// arquivos que EXISTEM
	write := func(rel string) {
		full := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		_ = os.WriteFile(full, []byte("x"), 0o644)
	}
	write("Login.spec.md")
	write("Login.tsx")
	write("Orphan.tsx") // código sem spec
	write("guides/G.md")

	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "Login.spec.md", Kind: mapx.KindSpec, Code: "LOGIX"},
			{ID: "Login.tsx", Kind: mapx.KindCode},
			{ID: "Orphan.tsx", Kind: mapx.KindCode},                   // código sem spec: aceitável, NÃO reportado
			{ID: "guides/G.md", Kind: mapx.KindGuide},                 // não rege nada → guide-sem-governo
			{ID: "Sumido.spec.md", Kind: mapx.KindSpec, Code: "SUMI"}, // no mapa, sem arquivo → no-fantasma
		},
		Edges: []mapx.Edge{
			{From: "Login.spec.md", To: "Login.tsx", Type: mapx.EdgeSpecifies},
			{From: "Login.spec.md", To: "Sumido.spec.md", Type: mapx.EdgeGoverns}, // aresta p/ nó fantasma
		},
	}
	cfg := &config.Config{
		Layers: map[string]config.Layer{
			"spec":    {Kind: "spec"},
			"code":    {Kind: "code"},
			"guide":   {Kind: "guide"},
			"feature": {Kind: "feature"}, // declarada mas SEM arquivos → camada-vazia
		},
		Gates: []config.Gate{
			{Name: "spec-completa", On: []string{"spec"}, Check: "non-empty"},
			// nenhum gate para "feature" nem "test" → kind-sem-gate (mas não há test/feature nós)
		},
	}
	return root, g, cfg
}

func TestDiagnose_mapFidelity(t *testing.T) {
	root, g, cfg := setup(t)
	r := Diagnose(g, cfg, root)
	if !hasFinding(r, "no-fantasma", "Sumido.spec.md") {
		t.Error("deveria detectar nó fantasma (no mapa, sem arquivo)")
	}
}

func TestDiagnose_orphans(t *testing.T) {
	root, g, cfg := setup(t)
	r := Diagnose(g, cfg, root)
	// código SEM spec NÃO é órfão — é o caso normal (utils, constants, hooks).
	if hasFinding(r, "codigo-sem-spec", "Orphan.tsx") {
		t.Error("código sem spec é aceitável; não deveria ser reportado")
	}
	// Login.spec tem código LOGIX → NÃO é identidade-ausente
	if hasFinding(r, "identidade-ausente", "Login.spec.md") {
		t.Error("Login.spec tem código; não deveria ser identidade-ausente")
	}
}

func TestDiagnose_looseLayers(t *testing.T) {
	root, g, cfg := setup(t)
	r := Diagnose(g, cfg, root)
	if !hasFinding(r, "camada-vazia", "feature") {
		t.Error("deveria detectar camada 'feature' declarada mas vazia")
	}
	if !hasFinding(r, "guide-sem-governo", "guides/G.md") {
		t.Error("deveria detectar guide que não rege nada")
	}
}

func TestDiagnose_identidadeAusente(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "NoCode.spec.md"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "NoCode.tsx"), []byte("x"), 0o644)
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "NoCode.spec.md", Kind: mapx.KindSpec, Code: ""}, // sem código
			{ID: "NoCode.tsx", Kind: mapx.KindCode},
		},
		Edges: []mapx.Edge{{From: "NoCode.spec.md", To: "NoCode.tsx", Type: mapx.EdgeSpecifies}},
	}
	r := Diagnose(g, &config.Config{Layers: map[string]config.Layer{}}, root)
	if !hasFinding(r, "identidade-ausente", "NoCode.spec.md") {
		t.Error("spec sem código deveria ser identidade-ausente (órfão invisível)")
	}
}

func findingSeverity(r Report, check, subject string) (Severity, bool) {
	for _, f := range r.Findings {
		if f.Check == check && f.Subject == subject {
			return f.Severity, true
		}
	}
	return "", false
}

func TestCheckDuplicateCodes(t *testing.T) {
	g := &mapx.Graph{Nodes: []mapx.Node{
		// unidade única com código UNIQ — não colide
		{ID: "features/a/screens/A.spec.md", Kind: mapx.KindSpec, Code: "UNIQ"},
		{ID: "features/a/screens/A.tsx", Kind: mapx.KindCode, Code: "UNIQ"},
		// intra-feature: lib + screens da MESMA feature 'b' com INTR → Info
		{ID: "features/b/lib/x.ts", Kind: mapx.KindCode, Code: "INTR"},
		{ID: "features/b/screens/X.tsx", Kind: mapx.KindCode, Code: "INTR"},
		// cross-domain: features distintas com CROSS → Warn
		{ID: "features/auth/screens/S.tsx", Kind: mapx.KindCode, Code: "CROSS"},
		{ID: "features/dash/screens/H.tsx", Kind: mapx.KindCode, Code: "CROSS"},
		// doc NÃO é dono: só referencia CROSS, não deve virar dono
		{ID: "guides/report.md", Kind: mapx.KindDoc, Code: "CROSS"},
	}}
	r := Diagnose(g, &config.Config{Layers: map[string]config.Layer{}}, t.TempDir())

	if hasFinding(r, "identidade-duplicada", "UNIQ") {
		t.Error("UNIQ tem um só dono — não deveria ser duplicada")
	}
	// intra-feature (lib+screen da mesma feature) é compartilhamento legítimo — NÃO
	// reportado (um achado que nunca exige ação é ruído).
	if hasFinding(r, "identidade-duplicada", "INTR") {
		t.Error("INTR (intra-feature) é legítimo; não deveria gerar finding")
	}
	// só a colisão cross-domain é reportada, como Warn.
	if sev, ok := findingSeverity(r, "identidade-duplicada", "CROSS"); !ok || sev != Warn {
		t.Errorf("CROSS (cross-domain) deveria ser Warn, veio sev=%v ok=%v", sev, ok)
	}
}

// Typo em `skip_on` é o erro que mais engana: não casa perspectiva nenhuma, o gate segue
// rodando nas duas, e o autor acredita tê-lo desligado numa delas. Falha silenciosa do
// lado perigoso — quem escreveu queria MENOS execução e recebeu mais.
func TestCheckSkipOnValido(t *testing.T) {
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "com-typo", SkipOn: []string{"chnage"}},
		{Name: "certo", SkipOn: []string{config.PerspectiveAll}},
		{Name: "sem-declaracao"},
	}}

	fs := checkSkipOnValido(cfg)
	if len(fs) != 1 {
		t.Fatalf("só o typo deveria ser acusado, veio %d: %+v", len(fs), fs)
	}
	if fs[0].Subject != "com-typo" {
		t.Errorf("o achado deve nomear o gate: %+v", fs[0])
	}
	if !strings.Contains(fs[0].Detail, "chnage") {
		t.Errorf("o achado deve citar o valor inválido: %q", fs[0].Detail)
	}
}
