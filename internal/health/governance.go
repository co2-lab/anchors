package health

import (
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkGovernanceOpportunities identifica oportunidades de governança e configurações
// subótimas com base nos artefatos existentes no grafo e na configuração do projeto.
//
// O doctor e o check usam este mecanismo para apontar gates canônicos e boas práticas
// que nasceram depois da criação do projeto ou que ainda não foram adotados.
func checkGovernanceOpportunities(g *mapx.Graph, cfg *config.Config) []Finding {
	if cfg == nil || g == nil {
		return nil
	}

	var out []Finding
	hasKind := map[mapx.Kind]bool{}
	for _, n := range g.Nodes {
		hasKind[n.Kind] = true
	}

	gateMap := map[string]config.Gate{}
	for _, gt := range cfg.Gates {
		gateMap[gt.Name] = gt
	}

	// 1. evidence-fresh: projeto tem testes e código, mas não monitora frescor (teste stale)
	if hasKind[mapx.KindTest] && hasKind[mapx.KindCode] {
		if _, ok := gateMap["evidence-fresh"]; !ok {
			out = append(out, Finding{
				Check:    "sugestao-gate",
				Severity: Info,
				Subject:  "evidence-fresh",
				Detail:   i18n.T("health.opportunity.evidence_fresh"),
			})
		}
	}

	// 2. no-secret-leaked: projeto tem código, mas não checa segredos
	if hasKind[mapx.KindCode] {
		if gt, ok := gateMap["no-secret-leaked"]; !ok {
			out = append(out, Finding{
				Check:    "sugestao-gate",
				Severity: Info,
				Subject:  "no-secret-leaked",
				Detail:   i18n.T("health.opportunity.no_secret_leaked"),
			})
		} else if !gt.IsBlocking() {
			// Gate configurado como subótimo: segredo no git não tem volta simples
			out = append(out, Finding{
				Check:    "gate-subotimo",
				Severity: Info,
				Subject:  "no-secret-leaked",
				Detail:   i18n.T("health.opportunity.no_secret_blocking"),
			})
		}
	}

	// 3. dependency-vulnerable: projeto tem código, mas não audita CVEs de dependências
	if hasKind[mapx.KindCode] {
		if _, ok := gateMap["dependency-vulnerable"]; !ok {
			out = append(out, Finding{
				Check:    "sugestao-gate",
				Severity: Info,
				Subject:  "dependency-vulnerable",
				Detail:   i18n.T("health.opportunity.dependency_vulnerable"),
			})
		}
	}

	// 4. no-duplication: projeto tem código, mas não mede copy-paste
	if hasKind[mapx.KindCode] {
		if _, ok := gateMap["no-duplication"]; !ok {
			out = append(out, Finding{
				Check:    "sugestao-gate",
				Severity: Info,
				Subject:  "no-duplication",
				Detail:   i18n.T("health.opportunity.no_duplication"),
			})
		}
	}

	// 5. Configuração subótima de testes: tem nós de teste mas não tem junit configurado
	if hasKind[mapx.KindTest] {
		hasJunit := false
		for _, t := range cfg.Tests {
			if t.JUnit != "" {
				hasJunit = true
				break
			}
		}
		if !hasJunit {
			out = append(out, Finding{
				Check:    "config-subotima",
				Severity: Info,
				Subject:  "tests.junit",
				Detail:   i18n.T("health.opportunity.tests_junit"),
			})
		}
	}

	return out
}

// QuickGovernanceHints devolve uma lista condensada das oportunidades de governança mais
// relevantes para exibição rápida e não obstrutiva (por exemplo, no rodapé do `check`).
func QuickGovernanceHints(g *mapx.Graph, cfg *config.Config) []Finding {
	findings := checkGovernanceOpportunities(g, cfg)
	if len(findings) > 2 {
		return findings[:2]
	}
	return findings
}
