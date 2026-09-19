package health

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func TestCheckGovernanceOpportunities_EvidenceFreshSugestao(t *testing.T) {
	// Projeto com testes e código, mas sem evidence-fresh
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "DataTable.tsx", Kind: mapx.KindCode},
			{ID: "DataTable.test.tsx", Kind: mapx.KindTest},
		},
	}
	cfg := &config.Config{
		Gates: []config.Gate{
			{Name: "tests-green", Check: "tests-pass"},
		},
	}

	achados := checkGovernanceOpportunities(g, cfg)
	achou := false
	for _, f := range achados {
		if f.Check == "sugestao-gate" && f.Subject == "evidence-fresh" {
			achou = true
			if f.Severity != Info {
				t.Errorf("severidade deveria ser Info, veio %v", f.Severity)
			}
		}
	}
	if !achou {
		t.Error("esperava sugestão de evidence-fresh para projeto com testes e código")
	}
}

func TestCheckGovernanceOpportunities_NoSecretLeakedNaoBloqueante(t *testing.T) {
	// Projeto com no-secret-leaked mas com blocking: false (subótimo)
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "main.go", Kind: mapx.KindCode},
		},
	}
	blockingFalse := false
	cfg := &config.Config{
		Gates: []config.Gate{
			{Name: "no-secret-leaked", Blocking: &blockingFalse},
		},
	}

	achados := checkGovernanceOpportunities(g, cfg)
	achou := false
	for _, f := range achados {
		if f.Check == "gate-subotimo" && f.Subject == "no-secret-leaked" {
			achou = true
		}
	}
	if !achou {
		t.Error("esperava aviso de gate subótimo para no-secret-leaked não bloqueante")
	}
}

func TestCheckGovernanceOpportunities_TestsSemJunit(t *testing.T) {
	// Projeto com testes mas sem junit configurado
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "auth.test.ts", Kind: mapx.KindTest},
		},
	}
	cfg := &config.Config{
		Tests: []config.Suite{
			{Layer: "unit", Run: "pnpm test"},
		},
	}

	achados := checkGovernanceOpportunities(g, cfg)
	achou := false
	for _, f := range achados {
		if f.Check == "config-subotima" && f.Subject == "tests.junit" {
			achou = true
		}
	}
	if !achou {
		t.Error("esperava aviso de configuração subótima para testes sem junit")
	}
}
