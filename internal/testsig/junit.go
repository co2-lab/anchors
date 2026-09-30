// Package testsig ingere SINAIS de qualidade de teste que o projeto já gera — o
// Anchors NÃO roda o teste, consome o artefato do runner (JUnit para execução, lcov
// para cobertura). Mantém o agnosticismo (D3): qualquer stack que emita esses
// formatos padrão é suportada sem o CLI saber rodar jest/pytest/go.
package testsig

import (
	"encoding/xml"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// CaseResult é um caso de teste individual, do JUnit.
type CaseResult struct {
	Name    string // nome do caso (ex.: "SPCRX-V01: eixo vertical...")
	Class   string // classe/suite (ex.: "Spacer (atom)")
	File    string // arquivo, se o reporter emitiu (atributo file)
	Failed  bool
	Skipped bool
	// Seconds is the case's run time, from its `time` attribute; 0 when absent or unreadable.
	Seconds float64
}

// SuiteResult agrupa casos por arquivo/suite.
type ExecReport struct {
	Cases []CaseResult
}

// --- estruturas do JUnit XML (o subconjunto universal) ---

type junitTestsuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []junitTestsuite `xml:"testsuite"`
	// alguns reporters emitem <testsuite> na raiz, sem <testsuites>
}

type junitTestsuite struct {
	XMLName xml.Name         `xml:"testsuite"`
	Name    string           `xml:"name,attr"`
	File    string           `xml:"file,attr"`
	Cases   []junitTestcase  `xml:"testcase"`
	Suites  []junitTestsuite `xml:"testsuite"` // aninhamento
}

type junitTestcase struct {
	Name      string       `xml:"name,attr"`
	Classname string       `xml:"classname,attr"`
	File      string       `xml:"file,attr"`
	Time      string       `xml:"time,attr"`
	Failure   *junitDetail `xml:"failure"`
	Error     *junitDetail `xml:"error"`
	Skipped   *junitDetail `xml:"skipped"`
}

type junitDetail struct {
	Message string `xml:"message,attr"`
}

// ParseJUnit lê um arquivo JUnit XML e devolve os casos, aceitando tanto a raiz
// <testsuites> quanto <testsuite> solto, e suites aninhadas.
func ParseJUnit(path string) (*ExecReport, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rep := &ExecReport{}

	// tenta <testsuites> na raiz
	var roots junitTestsuites
	if err := xml.Unmarshal(data, &roots); err == nil && len(roots.Suites) > 0 {
		for _, s := range roots.Suites {
			collectSuite(s, rep)
		}
		return rep, nil
	}
	// tenta <testsuite> solto na raiz
	var single junitTestsuite
	if err := xml.Unmarshal(data, &single); err == nil {
		collectSuite(single, rep)
		return rep, nil
	}
	return rep, nil
}

func collectSuite(s junitTestsuite, rep *ExecReport) {
	for _, c := range s.Cases {
		file := c.File
		if file == "" {
			file = s.File
		}
		rep.Cases = append(rep.Cases, CaseResult{
			Name:    c.Name,
			Class:   c.Classname,
			File:    file,
			Failed:  c.Failure != nil || c.Error != nil,
			Skipped: c.Skipped != nil,
			Seconds: caseSeconds(c.Time),
		})
	}
	for _, nested := range s.Suites {
		collectSuite(nested, rep)
	}
}

// scenarioCodeRE — extrai códigos de cenário do NOME de um caso (ex.: "SPCRX-V01: ...").
// Mesma gramática do resto do projeto (TRACEABILITY §3).
// caseCodeRE is built from the CURRENT vocabulary (`SetRuleLetters`, `SetCodeLenPattern`)
// and cached until it changes. A package-level `var` froze the defaults at load, so the
// setters changed nothing: a project's own letter (`-Z01`) was never read from a case
// name, and its rule showed no green test.
var (
	caseCodeMu  sync.Mutex
	caseCodeKey string
	caseCodeVal *regexp.Regexp
)

func caseCodeRE() *regexp.Regexp {
	caseCodeMu.Lock()
	defer caseCodeMu.Unlock()
	if key := ruleLetters + "\x00" + codeLenPattern; key != caseCodeKey || caseCodeVal == nil {
		caseCodeKey, caseCodeVal, caseScenarioVal = key, mustCodeRE(), mustScenarioRE()
	}
	return caseCodeVal
}

var caseScenarioVal *regexp.Regexp

func caseScenarioRE() *regexp.Regexp {
	caseCodeRE()
	caseCodeMu.Lock()
	defer caseCodeMu.Unlock()
	return caseScenarioVal
}

// ScenarioCodesInCase are the scenario codes a case's name mentions, each with its variant
// suffix when it carries one (`RDCH-B02#02`).
func ScenarioCodesInCase(name string) []string {
	return caseScenarioRE().FindAllString(name, -1)
}

// CodesInCase devolve os códigos de cenário mencionados no nome de um caso.
func CodesInCase(name string) []string {
	return caseCodeRE().FindAllString(name, -1)
}

// PassedCodes devolve os códigos de cenário que aparecem em casos que PASSARAM
// (não falharam nem foram pulados) — a base da cobertura semântica por cenário.
func (r *ExecReport) PassedCodes() map[string]bool {
	out := map[string]bool{}
	for _, c := range r.Cases {
		if c.Failed || c.Skipped {
			continue
		}
		for _, code := range ScenarioCodesInCase(c.Name) {
			out[strings.ToUpper(code)] = true
		}
	}
	return out
}

// SeenCodes are the scenario codes named by ANY case of the report — passed, failed or
// skipped. A partial run uses it to know which scenarios it actually measured: a code it
// did not see keeps the proof it had (see mapx.IngestExecutionSuite).
func (r *ExecReport) SeenCodes() map[string]bool {
	out := map[string]bool{}
	for _, c := range r.Cases {
		for _, code := range ScenarioCodesInCase(c.Name) {
			out[strings.ToUpper(code)] = true
		}
	}
	return out
}

// caseSeconds reads a case's `time` attribute: seconds, possibly fractional. A missing,
// malformed or negative value is 0 — a case the report did not time, not an error in the
// report.
func caseSeconds(v string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil || f < 0 {
		return 0
	}
	return f
}
