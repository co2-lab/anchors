// @anchors
//   code: CNGTC
//   ref: CTTST

package gate

import (
	"path"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// endpointSectionRE is the API catalog's `## Endpoint` section — the mark of an API unit.
var endpointSectionRE = regexp.MustCompile(`(?m)^##\s+Endpoint\s*$`)

// contract-tested: an API unit is proven against the project's OpenAPI document.
//
// The OpenAPI is compiled from the specs (`anchors docs build`), so it says what the API
// promises; a contract test is what says the implementation keeps it — each language has
// the tool that runs requests against an OpenAPI and validates the answers (Schemathesis,
// Dredd, kin-openapi, jest-openapi, openapi-core, swagger-request-validator). Anchors does
// not choose the tool; it asks the three things that make the proof traceable:
//
//	the scenario   `{CODE}-CT` in the unit's feature, with the project's contract regime
//	the test       a test of the unit that names `{CODE}-CT`
//	the contract   that test loads the OpenAPI document — it validates against the compiled
//	               contract, not against a copy of it written in the test
//
// It runs on the API unit's main code file (a layer tagged `interface`) and reads the spec
// beside it; a spec with no `Endpoint` section is no API unit.
func checkContractTested(_ string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, ""
	}
	base := strings.TrimSuffix(n.ID, path.Ext(n.ID))
	spec, err := readFile(root, base+".spec.md")
	if err != nil {
		return Skip, i18n.T("gate.vr_states.skip_no_spec")
	}
	if !endpointSectionRE.Match(spec) {
		return Skip, i18n.T("gate.contract_tested.skip_not_api")
	}
	m := specCodeRE().FindStringSubmatch(string(spec))
	if m == nil {
		return Skip, i18n.T("gate.vr_states.skip_no_code")
	}
	ct := m[1] + "-CT"
	tag := contractRegimeTag(cfg)

	var gaps []string
	feature, _ := readFile(root, base+".feature")
	if !hasContractScenario(string(feature), ct, tag) {
		gaps = append(gaps, i18n.T("gate.contract_tested.no_scenario", ct, tag))
	}
	named, loads := contractTests(root, g, path.Base(base), ct)
	switch {
	case len(named) == 0:
		gaps = append(gaps, i18n.T("gate.contract_tested.no_test", ct))
	case len(loads) == 0:
		gaps = append(gaps, i18n.T("gate.contract_tested.no_openapi", strings.Join(named, ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// contractRegimeTag is the tag the project gives the contract regime (`derived.regimes`
// maps tag → regime), with `contract-level` as the default.
func contractRegimeTag(cfg *config.Config) string {
	if cfg != nil && cfg.Derived != nil {
		for tag, regime := range cfg.Derived.Regimes {
			if strings.Contains(strings.ToLower(regime), "contract") || strings.Contains(strings.ToLower(regime), "contrato") {
				return strings.TrimPrefix(tag, "@")
			}
		}
	}
	return "contract-level"
}

// hasContractScenario says whether a feature line carries both the contract code and the
// contract regime tag.
func hasContractScenario(feature, ct, tag string) bool {
	for _, line := range strings.Split(feature, "\n") {
		if hasTag(line, "@"+ct) && hasTag(line, "@"+tag) {
			return true
		}
	}
	return false
}

var openapiRefRE = regexp.MustCompile(`(?i)openapi`)

// contractTests are the tests of the unit that name the contract code, and those of them
// that load an OpenAPI document. A test is of the unit when its path names the code, it
// sits beside the unit under its name, or a folder of its path is named after the unit.
func contractTests(root string, g *mapx.Graph, unit, ct string) (named, loads []string) {
	if g == nil {
		return nil, nil
	}
	code := strings.TrimSuffix(ct, "-CT")
	for _, t := range g.Nodes {
		if t.Kind != mapx.KindTest || isImage(t.ID) {
			continue
		}
		b := path.Base(t.ID)
		if !strings.Contains(t.ID, code+"-") && !strings.HasPrefix(b, unit+".") && !strings.HasPrefix(b, unit+"_") &&
			!strings.Contains("/"+t.ID, "/"+unit+"/") {
			continue
		}
		text := t.ID
		if c, err := readFile(root, t.ID); err == nil {
			text += "\n" + string(c)
		}
		if !strings.Contains(text, ct) {
			continue
		}
		named = append(named, t.ID)
		if openapiRefRE.MatchString(text) {
			loads = append(loads, t.ID)
		}
	}
	return named, loads
}
