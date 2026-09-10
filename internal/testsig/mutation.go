package testsig

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Ingestão de MUTAÇÃO pelo formato aberto **Mutation Testing Elements**
// (`schemaVersion: 1.x`) — o mesmo papel que o JUnit tem para execução e o lcov para
// cobertura: um formato que várias ferramentas emitem e nenhuma possui.
//
// Por que mutação importa: cobertura de linha diz que a linha EXECUTOU; não diz que
// alguém VERIFICOU o resultado. Um teste pode cobrir 100% das linhas e não provar
// nada. Mutação responde a pergunta certa — altere a linha; se o teste continuar
// verde, ele não prova aquela linha. É a única medida objetiva de "o teste prova algo",
// e pega a classe de erro que satisfaz todos os outros gates.
//
// O Anchors NÃO roda mutação e não conhece ferramenta. Emitem este schema, entre
// outros: Stryker (JS/TS/C#/Scala), PIT via plugin (Java), Infection (PHP),
// mutmut/cosmic-ray (Python). Um projeto sem ferramenta de mutação simplesmente não
// ingere o sinal — e o gate correspondente fica Pending, dizendo o que falta.

// MutationReport é o resultado agregado por ARQUIVO de origem.
type MutationReport struct {
	Files map[string]FileMutation
	// Low e o minimo ACEITAVEL e High o DESEJAVEL, como o projeto os declarou na
	// ferramenta. Zero = ausente no relatorio; quem decide passa a ser o default do
	// engine.
	Low, High float64
}

// FileMutation são os mutantes de um arquivo, já classificados.
type FileMutation struct {
	Killed   int
	Survived int
	// Score é killed / (killed + survived + timeout) × 100 — o denominador exclui os
	// mutantes que nem rodaram (erro de compilação/runtime), porque eles não dizem
	// nada sobre a qualidade do teste.
	//
	// Quando NÃO HÁ denominador — todo mutante do arquivo foi ignorado, ou a ferramenta
	// não gerou nenhum —, o score é 100: não sobra nada a provar ali. Ver `Ignored`.
	Score float64
	// NoCoverage são os mutantes que NENHUM TESTE EXECUTOU. Ficam fora do score, e a
	// razão é separação de responsabilidade: "existe teste que execute esta linha?" é a
	// pergunta do gate de COBERTURA, não do de mutação. O de mutação pergunta a seguinte —
	// "dado que executa, o teste VERIFICA o resultado?" — e ela só faz sentido depois que
	// a primeira foi respondida.
	//
	// Antes disto o NoCoverage contava como sobrevivente, com o argumento de que mutante
	// não executado é, por definição, não provado. O argumento é verdadeiro e mesmo assim
	// leva ao lugar errado: MEDIDO no app de referência em 25/08, 187 arquivos marcavam 0% e a
	// esmagadora maioria dos mutantes deles era NoCoverage — camadas provadas por
	// integração, que a config de mutação exclui de propósito. O resultado prático era um
	// gate com 187 achados que ele não é dono de resolver, e no meio deles se perdiam os
	// poucos zeros REAIS (o `models/holidays.ts`, com 32 sem cobertura e 15 sobreviventes
	// de verdade: teste que roda e não verifica).
	//
	// Um gate que reporta o que não é dele treina quem lê a ignorá-lo.
	NoCoverage int
	// Ignored são os mutantes que a FERRAMENTA descartou antes de rodar: `ignoreStatic`,
	// `// Stryker disable`, e equivalentes. Eles não entram no score porque não houve
	// experimento — mas são gravados, e por um motivo prático: sem eles, "arquivo 100%
	// porque tudo foi provado" e "arquivo 100% porque não havia o que provar" viram o
	// mesmo número, e quem lê uma lista ordenada não consegue separar os dois.
	Ignored int
	// SurvivedAt são as linhas onde um mutante sobreviveu — o que o autor precisa ver
	// para consertar o teste. Sem isso o score é um número sem ação.
	SurvivedAt []int
}

// mtElements é o subconjunto do schema que nos interessa.
type mtElements struct {
	SchemaVersion string `json:"schemaVersion"`
	// Os limiares do PROJETO viajam no proprio relatorio: `thresholds` e campo
	// OBRIGATORIO do schema Mutation Testing Elements, com `high` e `low` obrigatorios
	// dentro dele. Ler daqui e tao agnostico quanto ler o status dos mutantes — e
	// evita a duplicacao que declara-los no anchors.yaml criaria: a regua ja existe na
	// config da ferramenta e chega junto com a medicao.
	//
	// `break` NAO entra: e extensao do Stryker, fora do schema. Ele governa o codigo de
	// saida da ferramenta, que e outra decisao (do CI), nao a regua do gate.
	Thresholds *struct {
		High *float64 `json:"high"`
		Low  *float64 `json:"low"`
	} `json:"thresholds"`
	Files map[string]struct {
		Mutants []struct {
			Status   string `json:"status"`
			Location struct {
				Start struct {
					Line int `json:"line"`
				} `json:"start"`
			} `json:"location"`
		} `json:"mutants"`
	} `json:"files"`
}

// ParseMutation lê um relatório no formato canônico (Mutation Testing Elements).
//
// Mantida para quem não declara formato — é o caminho de todo projeto que já existia.
// Para escolher o formato, use ParseMutationFormat.
func ParseMutation(path string) (*MutationReport, error) {
	return ParseMutationFormat(path, "")
}

// ParseMutationFormat lê um relatório de mutação no formato pedido.
//
// `format` vem do gate `mutation-score` do anchors.yaml (`config.Gate.Format`). Vazio
// resolve para o canônico — o silêncio nunca muda o comportamento de quem já rodava.
func ParseMutationFormat(path, format string) (*MutationReport, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	switch strings.TrimSpace(strings.ToLower(format)) {
	case "", "mutation-testing-elements", "mte", "stryker":
		return parseMTE(b)
	case "gremlins":
		return parseGremlins(b)
	default:
		return nil, fmt.Errorf("formato de relatório de mutação desconhecido: %q "+
			"(aceitos: `mutation-testing-elements` — o default — e `gremlins`; "+
			"declare em `gates: - name: mutation-score / format:`)", format)
	}
}

func parseMTE(b []byte) (*MutationReport, error) {
	var raw mtElements
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("relatório de mutação inválido (esperado o formato "+
			"Mutation Testing Elements, schemaVersion 1.x): %w", err)
	}
	if len(raw.Files) == 0 {
		return nil, fmt.Errorf("relatório de mutação sem arquivos — confira se a " +
			"ferramenta emitiu o formato JSON padrão (Mutation Testing Elements)")
	}

	rep := &MutationReport{Files: map[string]FileMutation{}}
	if raw.Thresholds != nil {
		if raw.Thresholds.Low != nil {
			rep.Low = *raw.Thresholds.Low
		}
		if raw.Thresholds.High != nil {
			rep.High = *raw.Thresholds.High
		}
	}
	for path, f := range raw.Files {
		var fm FileMutation
		rodados := 0
		for _, m := range f.Mutants {
			switch strings.ToLower(m.Status) {
			case "killed":
				fm.Killed++
				rodados++
			case "survived":
				// O teste EXECUTOU a linha mutada e não percebeu a diferença. É o único
				// achado que o gate de mutação é dono de cobrar.
				fm.Survived++
				fm.SurvivedAt = append(fm.SurvivedAt, m.Location.Start.Line)
				rodados++
			case "nocoverage", "no coverage":
				// Nenhum teste executou. Contado à parte e fora do score — ver o campo.
				fm.NoCoverage++
			case "timeout":
				// Timeout é morte por travamento — o teste percebeu a mutação.
				fm.Killed++
				rodados++
			case "ignored":
				// A ferramenta descartou o mutante ANTES de rodar. Não é teste faltando
				// nem teste fraco: é o instrumentador dizendo que ali não há experimento
				// a fazer. Fica fora do score e é contado à parte.
				fm.Ignored++
			}
			// compileerror / runtimeerror: o mutante não rodou; não diz nada sobre o
			// teste, então fica fora do denominador.
		}
		if rodados > 0 {
			fm.Score = float64(fm.Killed) / float64(rodados) * 100
		} else {
			// Nenhum mutante RODOU: todos foram ignorados, nenhum teste executou o
			// arquivo, ou a ferramenta não gerou mutante nenhum (tabela, tipo, reexport).
			// Nos três casos não há dívida DESTE gate — 0/0 vira 100, e não "sem sinal".
			// Quando o motivo é ausência de teste, quem cobra é o gate de cobertura.
			//
			// Antes disso o arquivo ficava num limbo: o sinal era gravado vazio, e o gate,
			// que testava `killed == 0 && survived == 0`, mandava "rode a ferramenta de
			// mutação" para um arquivo em que ela já tinha rodado e ignorado tudo
			// corretamente. O contador de `Ignored` é o que impede o 100 de mentir.
			fm.Score = 100
		}
		rep.Files[normalizeMutationPath(path)] = fm
	}
	return rep, nil
}

// normalizeMutationPath deixa o caminho no formato do mapa (relativo, sem prefixo de
// diretório de execução). Ferramentas diferentes emitem caminhos diferentes — absoluto,
// relativo à raiz do pacote, ou com `./`.
func normalizeMutationPath(p string) string {
	p = strings.TrimPrefix(p, "./")
	if i := strings.Index(p, "/src/"); i >= 0 {
		return p[i+1:]
	}
	return p
}
