// @anchors
//   ref: GRING

package testsig

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Ingestão de MUTAÇÃO pelo formato do **gremlins** (github.com/go-gremlins/gremlins).
//
// POR QUE ESTE ARQUIVO EXISTE
//
// O formato canônico do Anchors é o Mutation Testing Elements (ver mutation.go) — um
// formato que várias ferramentas emitem e nenhuma possui. Em Go, porém, ele não existe
// na prática. Levantamento do ecossistema em 2026-08-24:
//
//	gremlins             391★, o runner dominante  → formato PRÓPRIO, não MTE
//	avito/go-mutesting   259★                      → formato próprio
//	zimmski/go-mutesting 673★, parado desde 2024   → sem JSON
//	gtramontina/ooze     284★                      → só texto
//	szhekpisov/gomutants   6★, criado abr/2026     → MTE nativo
//
// e nenhum conversor gremlins→MTE público. O `--output` do gremlins não escolhe
// formato: é só o caminho do arquivo (há um único `json.Marshal` no projeto inteiro).
// O próprio projeto MTE não lista nenhum framework de Go.
//
// Ou seja: o conselho que o `doctor` e o `coverage` dão a qualquer stack — "rode a
// ferramenta de mutação e ingira o relatório" — levava, em Go, a um beco: a ferramenta
// dominante existe, roda, e o relatório dela não entrava. Aceitar este formato é o que
// torna o gate `mutation-score` alcançável para projetos Go reais.
//
// O QUE SE PERDE (e por que não impede a medida)
//
// O gremlins emite MENOS que o MTE: não traz o texto-fonte, nem a posição final do
// mutante, nem id estável. Nada disso entra na conta do score — o `MutationReport` usa
// status e linha inicial, que é exatamente o que o gremlins tem. A ingestão é portanto
// COMPLETA para o que o Anchors mede; o que falta seria necessário só para renderizar o
// relatório visual do Stryker, que não é papel do Anchors.

// gremlinsReport é o subconjunto do formato do gremlins que nos interessa.
//
// A diferença ESTRUTURAL para o MTE é `files`: aqui é um ARRAY de objetos com
// `file_name`; no MTE é um OBJETO indexado pelo caminho. São mutuamente exclusivos para
// um decodificador JSON — é o que permite detectar o formato errado com mensagem útil.
type gremlinsReport struct {
	GoModule string `json:"go_module"`
	Files    []struct {
		Filename  string `json:"file_name"`
		Mutations []struct {
			Status string `json:"status"`
			Type   string `json:"type"`
			Line   int    `json:"line"`
		} `json:"mutations"`
	} `json:"files"`
}

// parseGremlins converte um relatório do gremlins no mesmo MutationReport que o parser
// de MTE produz — daí para frente nada no Anchors sabe de qual ferramenta veio.
//
// Os limiares (Low/High) ficam ZERO: o gremlins carrega `threshold-efficacy` e
// `threshold-mcover` na sua própria configuração e NÃO os escreve no relatório. Zero
// significa "ausente", e o engine cai no default — o mesmo caminho de um relatório MTE
// sem `thresholds`. Inventar um limiar aqui seria fabricar régua que ninguém declarou.
func parseGremlins(b []byte) (*MutationReport, error) {
	var raw gremlinsReport
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("invalid mutation report (the gate declares "+
			"`format: gremlins`, which expects `files` as a LIST of `file_name`; a "+
			"Mutation Testing Elements report brings `files` as an object — check "+
			"whether `format:` matches the tool): %w", err)
	}
	if len(raw.Files) == 0 {
		return nil, fmt.Errorf("gremlins report with no files — check whether the " +
			"run produced mutants (`gremlins unleash --output <file>`)")
	}

	rep := &MutationReport{Files: map[string]FileMutation{}}
	for _, f := range raw.Files {
		var fm FileMutation
		rodados := 0
		for _, m := range f.Mutations {
			switch normalizeGremlinsStatus(m.Status) {
			case "killed":
				// KILLED: o teste percebeu a alteração.
				fm.Killed++
				rodados++
			case "timedout":
				// TIMED OUT: morte por travamento, killed como no parser de MTE — e contado
				// à parte, porque sob carga é o que infla o score (ver TimedOut).
				fm.Killed++
				fm.TimedOut++
				rodados++
			case "survived":
				// LIVED: the test ran the mutated line and did not notice.
				fm.Survived++
				fm.SurvivedAt = append(fm.SurvivedAt, m.Line)
				rodados++
			case "nocoverage":
				// NOT COVERED: no test ran the line. Counted apart and out of the score, as
				// the canonical reading does with NoCoverage (see FileMutation.NoCoverage).
				// It used to count as a survivor here while the comment claimed the two
				// readings agreed — the same file scored 50 from gremlins and 100 from MTE.
				fm.NoCoverage++
				fm.NoCoverageAt = append(fm.NoCoverageAt, m.Line)
			}
			// NOT VIABLE (não compilou) e RUNNABLE/SKIPPED (não chegaram a rodar) ficam
			// FORA do denominador: não dizem nada sobre a qualidade do teste. É a mesma
			// regra que exclui compileerror/runtimeerror no MTE.
		}
		if rodados > 0 {
			fm.Score = float64(fm.Killed) / float64(rodados) * 100
		} else {
			// No mutant ran: 0/0 is 100, as in the canonical reading (parseMTE says why).
			// It used to stay 0, so the same file failed from one tool and passed from the
			// other; the NoCoverage count keeps the 100 from lying.
			fm.Score = 100
		}
		rep.Files[normalizeMutationPath(f.Filename)] = fm
	}
	return rep, nil
}

// normalizeGremlinsStatus traduz o vocabulário do gremlins para as classes que mudam a
// conta (killed, survived, nocoverage; o resto fica fora). Os literais vêm de `internal/mutator/mutator.go` (método String()):
// NOT COVERED, RUNNABLE, SKIPPED, LIVED, KILLED, NOT VIABLE, TIMED OUT.
//
// A normalização remove espaço e caixa porque o status viaja como texto humano ("NOT
// COVERED"), não como enum — e uma versão futura pode variar a grafia.
func normalizeGremlinsStatus(s string) string {
	switch strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), " ", "")) {
	case "killed":
		return "killed"
	case "timedout":
		return "timedout"
	case "lived":
		return "survived"
	case "notcovered":
		return "nocoverage"
	default:
		// notviable, runnable, skipped — e qualquer status futuro que não saibamos
		// classificar. Ficar fora do denominador é a escolha CONSERVADORA: um status
		// desconhecido não infla nem desinfla o score.
		return "ignored"
	}
}
