// Package testsig ingere SINAIS de qualidade de teste que o projeto já gera — o
// Anchors NÃO roda o teste, consome o artefato do runner (JUnit para execução, lcov
// para cobertura). Mantém o agnosticismo (D3): qualquer stack que emita esses
// formatos padrão é suportada sem o CLI saber rodar jest/pytest/go.
package testsig

import (
	"encoding/xml"
	"os"
	"strings"
)

// CaseResult é um caso de teste individual, do JUnit.
type CaseResult struct {
	Name    string // nome do caso (ex.: "SPCRX-V01: eixo vertical...")
	Class   string // classe/suite (ex.: "Spacer (atom)")
	File    string // arquivo, se o reporter emitiu (atributo file)
	Failed  bool
	Skipped bool
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
		})
	}
	for _, nested := range s.Suites {
		collectSuite(nested, rep)
	}
}

// scenarioCodeRE — extrai códigos de cenário do NOME de um caso (ex.: "SPCRX-V01: ...").
// Mesma gramática do resto do projeto (TRACEABILITY §3).
var caseCodeRE = mustCodeRE()

// CodesInCase devolve os códigos de cenário mencionados no nome de um caso.
func CodesInCase(name string) []string {
	return caseCodeRE.FindAllString(name, -1)
}

// PassedCodes devolve os códigos de cenário que aparecem em casos que PASSARAM
// (não falharam nem foram pulados) — a base da cobertura semântica por cenário.
func (r *ExecReport) PassedCodes() map[string]bool {
	out := map[string]bool{}
	for _, c := range r.Cases {
		if c.Failed || c.Skipped {
			continue
		}
		for _, code := range CodesInCase(c.Name) {
			out[strings.ToUpper(code)] = true
		}
	}
	return out
}
