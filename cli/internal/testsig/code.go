package testsig

import "regexp"

// ruleLetters são as letras de tipo de regra em vigor. Começam nas canônicas e são
// substituídas por SetRuleLetters quando a config do projeto é carregada — o vocabulário
// é do PROJETO (`rule_types`), não uma lista fixa do engine.
//
// Por que importa aqui: este regex extrai o scenario-code dos NOMES dos casos no relatório
// JUnit. Preso às canônicas, um cenário de letra declarada pelo projeto (ex.: `-I01`, de
// Invariant) nunca é reconhecido como provado — e o requisito aparece "sem teste verde"
// mesmo tendo um teste que passa.
var ruleLetters = "SRVAXBNMD"

// codeLenPattern espelha `config.CodeLengthPattern()`. Duplicado pelo mesmo motivo que
// `ruleLetters`: o pacote testsig não depende de scan nem de config, e o comprimento do
// código é do PROJETO — preso a um número fixo, um relatório com códigos de outro tamanho
// não casaria nenhum caso e todo requisito apareceria "sem teste verde".
var codeLenPattern = "{4,5}"

// SetCodeLenPattern ajusta o padrão de comprimento (ver config.CodeLengthPattern).
func SetCodeLenPattern(p string) {
	if p != "" {
		codeLenPattern = p
	}
}

// SetRuleLetters ajusta o vocabulário de letras usado ao ler o relatório de execução.
func SetRuleLetters(letters string) {
	if letters != "" {
		ruleLetters = letters
	}
}

// mustCodeRE devolve o regex de código de cenário — MESMA gramática do scan
// (TRACEABILITY §3), duplicada aqui para o pacote testsig não depender de scan.
func mustCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`\b[A-Z0-9]` + codeLenPattern + `-(?:[` + regexp.QuoteMeta(ruleLetters) + `]\d{2}(?:-[a-z][a-z0-9-]*)?|DS-[A-Za-z0-9-]+|VR)\b`)
}
