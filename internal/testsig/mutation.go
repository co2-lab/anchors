// @anchors
//   ref: MTINM

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
	// When there is NO denominator — every mutant was ignored, none was covered, or the
	// tool generated none — there is no score: zero, which the map does not write. Nothing
	// was measured, and a number there would say otherwise. See `Ignored`.
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
	// TimedOut are the mutants whose run hit the time limit. They stay counted as killed
	// (a hang is the test noticing the mutation), and are also counted here, because a
	// run under load times out mutants the test would NOT have caught: in the reference
	// app a file read 94.87% with 65 of 78 mutants timed out, and 47.44% measured clean.
	TimedOut int
	// SurvivedAt são as linhas onde um mutante sobreviveu — o que o autor precisa ver
	// para consertar o teste. Sem isso o score é um número sem ação.
	SurvivedAt []int
	// NoCoverageAt are the lines of the mutants no test ran. Together with a branch the
	// coverage report says no test took on the same line, they point at a branch that is
	// likely dead: not only untested, but unreachable by what the tests do.
	NoCoverageAt []int
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
		return nil, fmt.Errorf("unknown mutation report format: %q "+
			"(accepted: `mutation-testing-elements` — the default — and `gremlins`; "+
			"declare it in `gates: - name: mutation-score / format:`)", format)
	}
}

func parseMTE(b []byte) (*MutationReport, error) {
	var raw mtElements
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("invalid mutation report (expected the "+
			"Mutation Testing Elements format, schemaVersion 1.x): %w", err)
	}
	// A report in the format, with no file, is a run that had NOTHING TO MUTATE — a file
	// that is only `export const client = new Client({})`. It was an error, and the file
	// stayed pending forever: the tool ran, said so, and the map never heard. Without even
	// the schema version, it is no report of this format at all.
	if len(raw.Files) == 0 && strings.TrimSpace(raw.SchemaVersion) == "" {
		return nil, fmt.Errorf("mutation report with no files — check whether the " +
			"tool emitted the standard JSON format (Mutation Testing Elements)")
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
				fm.NoCoverageAt = append(fm.NoCoverageAt, m.Location.Start.Line)
			case "timeout":
				// Timeout é morte por travamento — o teste percebeu a mutação. Contado
				// também à parte: sob carga, é o que infla o score (ver TimedOut).
				fm.Killed++
				fm.TimedOut++
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
			// NO MUTANT RAN: all ignored, none covered, or none generated (a table, a
			// type, a re-export). There is no score — 0/0 is not 100. A 100 here was a
			// number with no measurement behind it, and the map recorded it as if the
			// tests had proven the file (measured: 15 mutants, all ignored, written as
			// `mutation_score: 100`). The counters say what happened, and the gate reads
			// them: it answers "does not apply", and the coverage gate owns a file no test
			// runs.
			fm.Score = 0
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
