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
// O DEFAULT tem de ser o mesmo do `config.DefaultRuleLetters`, e ele divergiu.
//
// Medido no blue-eyes: um projeto sem `rule_types:` declarado usa o default do `config`
// (`SRVAXBNMDEIQF`, que inclui `I` de Invariant). O `SetRuleLetters` é chamado com esse
// valor e a divergência não apareceria — MAS o `ingest` só o chama quando a config
// carrega, e qualquer caminho que leia o relatório antes disso usa esta constante.
//
// O sintoma foi o que o comentário acima descreve: 11 casos verdes, três cenários
// reconhecidos como provados, e os invariantes (`GLCGL-I01`, `I02`, `I03`) aparecendo
// "sem teste verde" com teste passando.
//
// Duas cópias da mesma lista divergem na primeira letra nova — e esta ficou três atrás
// (`E`, `I`, `Q`, `F`). O comentário do `config` já registra que "é a terceira vez que a
// lista fica para trás de uma letra nova".
var ruleLetters = "SRVAXBNMDEIQF"

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
