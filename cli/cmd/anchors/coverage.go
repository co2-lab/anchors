package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testsig"
	"github.com/spf13/cobra"
)

func readFileString(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}

// `anchors coverage` lê os sinais ingeridos e responde as duas perguntas de confiança:
//   - cobertura por CENÁRIO: cada código de cenário da spec tem um teste que passou?
//     (o diferencial — cobertura semântica, "o requisito está provado?")
//   - cobertura de LINHA: quais arquivos de código estão abaixo do limiar (do lcov).
//
// Sinais STALE (o arquivo mudou desde a ingestão) são marcados — não confie neles.
func newCoverageCmd() *cobra.Command {
	var root, mapPath, diffRef, diffFile, lcov string
	var threshold float64
	var delta bool
	cmd := &cobra.Command{
		Use:   "coverage [<spec>]",
		Short: "Mostra a cobertura por cenário, por linha, do DIFF e o delta",
		Long: `Lê os sinais de 'anchors ingest' e responde as perguntas de confiança:

  anchors coverage                    panorama: cenários sem teste verde + linha < limiar
  anchors coverage <spec>             os cenários de UMA spec e quais estão provados
  anchors coverage --diff main --lcov cov.info   o que MUDEI (vs main) está coberto?
  anchors coverage --delta            a cobertura CAIU desde a ingestão anterior?

Três perguntas complementares:
- por CENÁRIO: cada requisito da spec tem um teste que passou? (semântica)
- do DIFF: as linhas que você mudou (git diff) estão cobertas? (pega bug na linha nova)
- DELTA: a cobertura global caiu vs. a medição anterior? (pega regressão de cobertura —
  código todo coberto no diff, mas derrubou a cobertura de outra parte)`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			// --diff não usa o grafo (é um cruzamento diff × lcov) — não exige o mapa.
			if diffRef != "" || diffFile != "" {
				return coverageDiff(absRoot, diffRef, diffFile, lcov, threshold)
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}
			if delta {
				return coverageDelta(g)
			}
			if len(args) == 1 {
				return coverageForSpec(g, absRoot, relTo(absRoot, args[0]))
			}
			return coveragePanorama(g, absRoot, threshold)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().Float64Var(&threshold, "threshold", 70, "limiar de cobertura (%)")
	cmd.Flags().StringVar(&diffRef, "diff", "", "cobertura do DIFF vs esta ref git (ex: main)")
	cmd.Flags().StringVar(&diffFile, "diff-file", "", "cobertura do DIFF a partir de um unified diff (fallback não-git)")
	cmd.Flags().StringVar(&lcov, "lcov", "", "o lcov a cruzar com o diff (obrigatório com --diff)")
	cmd.Flags().BoolVar(&delta, "delta", false, "a cobertura caiu desde a ingestão anterior?")
	return cmd
}

// coverageForSpec: lê os códigos de cenário declarados na spec (do arquivo) e cruza
// com os provados (o Signal do próprio nó da spec, gravado na ingestão).
func coverageForSpec(g *mapx.Graph, root, specID string) error {
	var node *mapx.Node
	for i := range g.Nodes {
		if g.Nodes[i].ID == specID {
			node = &g.Nodes[i]
			break
		}
	}
	if node == nil {
		return fmt.Errorf("spec %q não está no mapa", specID)
	}
	declared, err := codesInFile(filepath.Join(root, specID))
	if err != nil {
		return fmt.Errorf("ler a spec: %w", err)
	}
	if len(declared) == 0 {
		fmt.Printf("%s não declara códigos de cenário (nada a cobrir)\n", specID)
		return nil
	}
	proven := map[string]bool{}
	if node.Signal != nil {
		for _, c := range node.Signal.ProvenCodes {
			proven[c] = true
		}
	}
	sort.Strings(declared)
	provados, faltando := 0, []string{}
	fmt.Printf("cobertura por cenário — %s:\n\n", specID)
	for _, code := range declared {
		if proven[code] {
			fmt.Printf("  ✓ %s — provado por teste verde\n", code)
			provados++
		} else {
			fmt.Printf("  ✗ %s — sem teste que passe\n", code)
			faltando = append(faltando, code)
		}
	}
	fmt.Printf("\n%d/%d cenário(s) provado(s)", provados, len(declared))
	if node.SignalStale() {
		fmt.Printf("  ⚠ SINAL STALE (a spec mudou desde a ingestão — reingira)")
	}
	fmt.Println()
	if len(faltando) > 0 {
		fmt.Printf("faltam testes verdes para: %v\n", faltando)
	}
	return nil
}

func coveragePanorama(g *mapx.Graph, root string, threshold float64) error {
	// cenários sem prova: specs cujos códigos declarados não estão todos provados
	var specGaps, lineGaps, stale int
	fmt.Println("== cobertura por cenário (specs com requisito sem teste verde) ==")
	specShown := 0
	for i := range g.Nodes {
		n := g.Nodes[i]
		if n.Kind != mapx.KindSpec {
			continue
		}
		declared, _ := codesInFile(filepath.Join(root, n.ID))
		if len(declared) == 0 {
			continue
		}
		proven := map[string]bool{}
		if n.Signal != nil {
			for _, c := range n.Signal.ProvenCodes {
				proven[c] = true
			}
		}
		var missing []string
		for _, c := range declared {
			if !proven[c] {
				missing = append(missing, c)
			}
		}
		if len(missing) > 0 {
			specGaps++
			if specShown < 15 {
				sort.Strings(missing)
				fmt.Printf("  %s — %d/%d sem prova: %v\n", n.ID, len(missing), len(declared), missing)
				specShown++
			}
		}
	}
	if specGaps == 0 {
		fmt.Println("  (nenhuma — todos os cenários declarados têm teste verde, OU nada foi ingerido)")
	} else if specGaps > specShown {
		fmt.Printf("  … e mais %d spec(s) com lacuna de cenário\n", specGaps-specShown)
	}

	// quantos nós de código TÊM cada sinal — o denominador que separa "medido e bom" de
	// "não medido".
	var comCov, comMut, códigos int
	for i := range g.Nodes {
		if g.Nodes[i].Kind != mapx.KindCode {
			continue
		}
		códigos++
		if sg := g.Nodes[i].Signal; sg != nil {
			if sg.TotalLines > 0 {
				comCov++
			}
			if sg.MutantsKilled > 0 || sg.MutantsSurvived > 0 {
				comMut++
			}
		}
	}

	fmt.Printf("\n== cobertura de linha (< %.0f%%) ==\n", threshold)
	lineShown := 0
	for i := range g.Nodes {
		n := g.Nodes[i]
		if n.Kind != mapx.KindCode || n.Signal == nil || n.Signal.TotalLines == 0 {
			continue
		}
		if n.Signal.LineCoverage < threshold {
			lineGaps++
			if lineShown < 15 {
				mark := ""
				if n.SignalStale() {
					mark = " ⚠stale"
					stale++
				}
				fmt.Printf("  %s — %.0f%%%s\n", n.ID, n.Signal.LineCoverage, mark)
				lineShown++
			}
		}
	}
	if lineGaps == 0 {
		// Distingue os dois casos que antes se confundiam num só texto: "tudo coberto" e
		// "nada medido" levam ao mesmo silêncio, e são conclusões opostas.
		if comCov == 0 {
			fmt.Println("  ⚠ nenhuma cobertura de linha ingerida — este relatório não está " +
				"dizendo que está tudo coberto, e sim que nada foi medido. " +
				"Rode `anchors ingest --lcov <arquivo>`")
		} else {
			fmt.Printf("  ✓ nenhum dos %d arquivo(s) medido(s) abaixo do limiar\n", comCov)
		}
	} else if lineGaps > lineShown {
		fmt.Printf("  … e mais %d arquivo(s) abaixo do limiar\n", lineGaps-lineShown)
	}

	// MUTAÇÃO — a pergunta que as duas seções acima não respondem. Cobertura de linha
	// diz que a linha executou; mutação diz se alguém verificou o resultado.
	fmt.Printf("\n== mutação: o teste PROVA a linha? ==\n")
	mutGaps, mutShown, sobreviventes := 0, 0, 0
	for i := range g.Nodes {
		n := g.Nodes[i]
		if n.Kind != mapx.KindCode || n.Signal == nil {
			continue
		}
		// Mesma régua do gate: sem mutante mensurável não há o que listar. O `Ignored`
		// não muda a decisão aqui — muda no gate, que precisa distinguir "nunca rodou".
		if n.Signal.MutantsKilled == 0 && n.Signal.MutantsSurvived == 0 {
			continue
		}
		sobreviventes += n.Signal.MutantsSurvived
		if n.Signal.MutationScore < threshold {
			mutGaps++
			if mutShown < 15 {
				mark := ""
				if n.SignalStale() {
					mark = " ⚠stale"
				}
				fmt.Printf("  %s — %.0f%% (%d sobrevivente(s))%s\n",
					n.ID, n.Signal.MutationScore, n.Signal.MutantsSurvived, mark)
				mutShown++
			}
		}
	}
	switch {
	case comMut == 0:
		fmt.Printf("  ⚠ nenhum sinal de mutação ingerido (%d arquivo(s) de código). "+
			"Uma suíte pode estar 100%% verde e 100%% coberta e ainda assim não derrubar "+
			"teste algum quando você apaga uma guarda. Se o stack tiver ferramenta "+
			"(Stryker/PIT/Infection/mutmut), rode e `anchors ingest --mutation <relatório>`\n",
			códigos)
	case mutGaps == 0:
		// Acima do limiar NÃO é o mesmo que sem sobreviventes: cada sobrevivente ainda é
		// uma linha que ninguém prova. O limiar é tolerância, não aprovação.
		if sobreviventes > 0 {
			fmt.Printf("  ~ %d arquivo(s) medido(s) acima do limiar, mas com %d mutante(s) "+
				"sobrevivente(s) — linhas que nenhum teste prova\n", comMut, sobreviventes)
		} else {
			fmt.Printf("  ✓ %d arquivo(s) medido(s), nenhum mutante sobreviveu\n", comMut)
		}
		if comMut < códigos {
			fmt.Printf("  (sem esse sinal: %d de %d arquivo(s) de código)\n", códigos-comMut, códigos)
		}
	default:
		if mutGaps > mutShown {
			fmt.Printf("  … e mais %d arquivo(s) abaixo do limiar\n", mutGaps-mutShown)
		}
		if comMut < códigos {
			fmt.Printf("  (medidos: %d de %d arquivo(s) de código)\n", comMut, códigos)
		}
	}

	fmt.Printf("\nresumo: %d spec(s) com cenário sem prova", specGaps)
	// "0 abaixo do limiar" e "nada medido" são conclusões opostas; o resumo não pode
	// imprimir o mesmo número para as duas.
	if comCov > 0 {
		fmt.Printf("; %d arquivo(s) < %.0f%% de linha", lineGaps, threshold)
	} else {
		fmt.Printf("; cobertura de linha NÃO medida")
	}
	if comMut > 0 {
		fmt.Printf("; %d mutante(s) sobrevivente(s) em %d arquivo(s) medido(s)", sobreviventes, comMut)
	} else {
		fmt.Printf("; mutação NÃO medida")
	}
	if stale > 0 {
		fmt.Printf("; %d sinal(is) stale (reingira)", stale)
	}
	fmt.Println()
	return nil
}

// coverageDiff responde "as linhas que eu mudei estão cobertas?" — cruza o diff
// (git ou arquivo) com o lcov. Só linhas INSTRUMENTADAS do diff contam (comentários/
// blanks não). Reporta as linhas mudadas e descobertas + o % de cobertura do diff.
func coverageDiff(root, ref, diffFile, lcov string, threshold float64) error {
	if lcov == "" {
		return fmt.Errorf("--lcov <arquivo> é obrigatório com --diff (o Anchors cruza o diff com a cobertura)")
	}
	var changed testsig.ChangedLines
	var err error
	if diffFile != "" {
		changed, err = testsig.ParseDiffFile(diffFile)
	} else {
		changed, err = testsig.GitDiff(root, ref)
	}
	if err != nil {
		// A cobertura DE DIFF pergunta "as linhas que você mudou estão cobertas?" — a
		// pergunta inteira depende de haver histórico contra o que comparar.
		if msg := gitmeta.Explica(gitmeta.Verifica(root), "obter o diff"); msg != "" {
			return errors.New(msg + "\n  (o `--diff-file` aceita um unified diff de qualquer fonte, se você tiver um)")
		}
		return fmt.Errorf("obter o diff: %w", err)
	}
	cov, err := testsig.ParseLCOV(lcov)
	if err != nil {
		return fmt.Errorf("parse lcov: %w", err)
	}
	// índice cobertura por arquivo (casa por sufixo de caminho)
	covByFile := map[string]testsig.FileCoverage{}
	for _, fc := range cov.Files {
		covByFile[filepath.ToSlash(fc.File)] = fc
	}

	fmt.Printf("cobertura do diff")
	if ref != "" {
		fmt.Printf(" (vs %s)", ref)
	}
	fmt.Println(":")
	var totalInstr, totalUncov int
	shown := false
	for file, lines := range changed {
		fc, ok := matchCoverage(file, covByFile)
		if !ok {
			continue // arquivo mudado sem cobertura ingerida (ex.: não é código)
		}
		instr := fc.InstrumentedIn(lines)
		if instr == 0 {
			continue
		}
		uncov := fc.UncoveredIn(lines)
		totalInstr += instr
		totalUncov += len(uncov)
		shown = true
		pct := float64(instr-len(uncov)) / float64(instr) * 100
		fmt.Printf("  %s — %d linha(s) mudada(s) instrumentada(s), %.0f%% coberta(s)", file, instr, pct)
		if len(uncov) > 0 {
			sort.Ints(uncov)
			fmt.Printf(" — SEM teste: %v", uncov)
		}
		fmt.Println()
	}
	if !shown {
		fmt.Println("  (nenhuma linha de código instrumentada no diff — nada a cobrir)")
		return nil
	}
	diffPct := float64(totalInstr-totalUncov) / float64(totalInstr) * 100
	fmt.Printf("\ncobertura do diff: %.0f%% (%d/%d linhas mudadas cobertas)\n", diffPct, totalInstr-totalUncov, totalInstr)
	if diffPct < threshold {
		fmt.Printf("✗ abaixo do limiar (%.0f%%) — as mudanças introduzem código sem teste\n", threshold)
		os.Exit(1)
	}
	fmt.Println("✓ o que você mudou está coberto")
	return nil
}

// coverageDelta responde "a cobertura caiu?" — compara, por nó, a cobertura atual
// contra a anterior (o baseline preservado no Signal na última ingestão).
func coverageDelta(g *mapx.Graph) error {
	fmt.Println("delta de cobertura (vs. ingestão anterior):")
	var dropped, improved int
	worst := 0.0
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindCode || n.Signal == nil || n.Signal.PrevLineCoverage == 0 {
			continue
		}
		delta := n.Signal.LineCoverage - n.Signal.PrevLineCoverage
		if delta < -0.01 {
			dropped++
			fmt.Printf("  ⚠ %s — %.0f%% → %.0f%% (%.0f)\n", n.ID, n.Signal.PrevLineCoverage, n.Signal.LineCoverage, delta)
			if delta < worst {
				worst = delta
			}
		} else if delta > 0.01 {
			improved++
		}
	}
	if dropped == 0 {
		fmt.Println("  ✓ nenhum arquivo perdeu cobertura desde a última ingestão")
		if improved > 0 {
			fmt.Printf("  (%d arquivo(s) melhoraram)\n", improved)
		}
		return nil
	}
	fmt.Printf("\n✗ %d arquivo(s) perderam cobertura (pior: %.0f pontos)\n", dropped, worst)
	os.Exit(1)
	return nil
}

// matchCoverage acha a cobertura de um arquivo do diff, casando por sufixo de caminho
// (o lcov pode ter caminho absoluto ou com prefixo).
func matchCoverage(diffFile string, covByFile map[string]testsig.FileCoverage) (testsig.FileCoverage, bool) {
	if fc, ok := covByFile[diffFile]; ok {
		return fc, true
	}
	for cf, fc := range covByFile {
		if strings.HasSuffix(cf, "/"+diffFile) || strings.HasSuffix(diffFile, "/"+cf) {
			return fc, true
		}
	}
	return testsig.FileCoverage{}, false
}

// codesInFile lê um arquivo e extrai os códigos de cenário que ele declara.
func codesInFile(path string) ([]string, error) {
	data, err := readFileString(path)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []string
	for _, c := range testsig.CodesInCase(data) {
		if !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	return out, nil
}
