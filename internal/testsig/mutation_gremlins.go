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
		return nil, fmt.Errorf("relatório de mutação inválido (o gate declara "+
			"`format: gremlins`, que espera `files` como LISTA de `file_name`; um "+
			"relatório Mutation Testing Elements traz `files` como objeto — confira "+
			"se o `format:` bate com a ferramenta): %w", err)
	}
	if len(raw.Files) == 0 {
		return nil, fmt.Errorf("relatório do gremlins sem arquivos — confira se a " +
			"execução produziu mutantes (`gremlins unleash --output <arquivo>`)")
	}

	rep := &MutationReport{Files: map[string]FileMutation{}}
	for _, f := range raw.Files {
		var fm FileMutation
		rodados := 0
		for _, m := range f.Mutations {
			switch normalizeGremlinsStatus(m.Status) {
			case "killed":
				// KILLED e TIMED OUT: o teste percebeu a alteração. Timeout é morte por
				// travamento — a mesma leitura que o parser de MTE faz.
				fm.Killed++
				rodados++
			case "survived":
				// LIVED (sobreviveu ao teste) e NOT COVERED (nenhum teste sequer o
				// executou) contam igual: em ambos a linha não foi provada. É a mesma
				// decisão que o MTE toma com Survived/NoCoverage.
				fm.Survived++
				fm.SurvivedAt = append(fm.SurvivedAt, m.Line)
				rodados++
			}
			// NOT VIABLE (não compilou) e RUNNABLE/SKIPPED (não chegaram a rodar) ficam
			// FORA do denominador: não dizem nada sobre a qualidade do teste. É a mesma
			// regra que exclui compileerror/runtimeerror no MTE.
		}
		if rodados > 0 {
			fm.Score = float64(fm.Killed) / float64(rodados) * 100
		}
		rep.Files[normalizeMutationPath(f.Filename)] = fm
	}
	return rep, nil
}

// normalizeGremlinsStatus traduz o vocabulário do gremlins para as duas classes que
// mudam a conta. Os literais vêm de `internal/mutator/mutator.go` (método String()):
// NOT COVERED, RUNNABLE, SKIPPED, LIVED, KILLED, NOT VIABLE, TIMED OUT.
//
// A normalização remove espaço e caixa porque o status viaja como texto humano ("NOT
// COVERED"), não como enum — e uma versão futura pode variar a grafia.
func normalizeGremlinsStatus(s string) string {
	switch strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), " ", "")) {
	case "killed", "timedout":
		return "killed"
	case "lived", "notcovered":
		return "survived"
	default:
		// notviable, runnable, skipped — e qualquer status futuro que não saibamos
		// classificar. Ficar fora do denominador é a escolha CONSERVADORA: um status
		// desconhecido não infla nem desinfla o score.
		return "ignored"
	}
}
