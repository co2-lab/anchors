package mapx

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// A ingestão de sinais de teste no grafo (o Anchors consome o artefato do runner e o
// amarra aos nós). Casa os caminhos dos relatórios aos nós por SUFIXO de caminho —
// relatórios costumam ter caminhos absolutos ou relativos a subdirs, então casamos
// pelo fim do caminho, o mais específico possível.

// IngestExecution grava, nos nós de TESTE, quantos casos passaram/falharam/pularam,
// e nos nós de CÓDIGO/SPEC os códigos de cenário PROVADOS (que aparecem num caso que
// passou). `byFile` mapeia arquivo-de-teste → (passed, failed, skipped). `proven` é o
// conjunto global de códigos provados. `now` é carimbado por quem chama.
type ExecByFile struct {
	Passed, Failed, Skipped int
}

// declaredByNode: para os nós que precisam de cobertura semântica (specs), quais
// códigos de cenário ELES declaram (lido do arquivo pelo comando — o mapx não toca
// disco). IngestExecution cruza esses declarados com os provados.
func (g *Graph) IngestExecution(byFile map[string]ExecByFile, proven map[string]bool, declaredByNode map[string][]string, layer, now string) (matchedFiles, matchedCodes int) {
	return g.IngestExecutionSuite(byFile, proven, nil, declaredByNode, layer, "", now)
}

// IngestExecutionSuite is IngestExecution with the SUITE named (the report the execution
// came from). With a suite, each spec's proof is stored under its key and `ProvenCodes`
// becomes the union of all of them — re-ingesting the SAME suite still erases what it
// stopped proving, but one suite's ingestion no longer speaks for the others. An empty
// suite keeps the previous behaviour: the report is the whole measurement.
//
// An old proof with no suite (ingested before this field existed) has no known owner;
// the first suite ingestion touching the spec replaces it with the measured union.
//
// `seen` marks a PARTIAL run (`anchors test --changed`): the codes its cases named. A
// full run is the whole measurement, and what it did not prove stops being proven — that
// is how a deleted test loses its proof. A partial run measured only its cut: each spec
// keeps its earlier proof for the codes this run did not see, and takes this run's result
// for the codes it did. Before, a partial run of 4 test files erased the proof of every
// spec outside the cut (reported from MIF: MoneyDetailScreen lost 12 green scenarios to a
// run that never executed its test). nil = full run.
func (g *Graph) IngestExecutionSuite(byFile map[string]ExecByFile, proven, seen map[string]bool, declaredByNode map[string][]string, layer, suite, now string) (matchedFiles, matchedCodes int) {
	if layer == "" {
		layer = "unit" // camada default quando não informada
	}
	for i := range g.Nodes {
		n := &g.Nodes[i]
		// execução: casa nós de teste pelo caminho, ACUMULANDO por camada
		if n.Kind == KindTest {
			for file, ex := range byFile {
				if pathMatches(n.ID, file) {
					ensureSignal(n)
					if n.Signal.ByLayer == nil {
						n.Signal.ByLayer = map[string]LayerExec{}
					}
					// substitui a camada (reingerir a mesma camada atualiza; outras
					// camadas permanecem — é o merge).
					n.Signal.ByLayer[layer] = LayerExec{Passed: ex.Passed, Failed: ex.Failed, Skipped: ex.Skipped}
					n.Signal.Passed, n.Signal.Failed, n.Signal.Skipped = sumLayers(n.Signal.ByLayer)
					n.Signal.AtRev = n.Rev
					// Carimba também as revs do FECHO — o que este teste alcança descendo.
					// É o que permite dizer depois "a evidência venceu porque o util que ele
					// compõe mudou", e não só "o próprio arquivo de teste mudou".
					n.Signal.ClosureRev = g.EvidenceClosure(n.ID)
					n.Signal.IngestedAt = now
					matchedFiles++
					break
				}
			}
		}
		// cobertura semântica: dos cenários que ESTE nó declara, quais estão provados?
		if declared, ok := declaredByNode[n.ID]; ok {
			var pc []string
			for _, code := range declared {
				if proven[code] {
					pc = append(pc, code)
				}
			}
			if seen != nil {
				// PARTIAL: what this run did not see keeps its earlier proof.
				var antes []string
				if n.Signal != nil {
					if suite == "" {
						antes = n.Signal.ProvenCodes
					} else {
						antes = n.Signal.ProvenBySuite[suite]
					}
				}
				tocou := false
				for _, code := range declared {
					if seen[code] {
						tocou = true
						break
					}
				}
				if !tocou {
					continue // outside the cut: nothing measured here, nothing changes
				}
				for _, code := range antes {
					if !seen[code] && !proven[code] {
						pc = append(pc, code)
					}
				}
				sort.Strings(pc)
			}
			// GRAVA MESMO VAZIO. "Nenhum cenário provado" é informação, não ausência
			// dela — e a ingestão é a medição inteira, não um acréscimo à anterior.
			//
			// Com `len(pc) > 0` o nó guardava para sempre a última prova que teve.
			// MEDIDO: ao corrigir o filtro que fazia a spec declarar a regra do
			// vizinho, sete nós continuaram carregando `proven_codes` alheios; o
			// `map build` e o `ingest` rodaram de novo e não limparam nada.
			//
			// O caso geral é pior que o resíduo: uma unidade perde o último teste
			// verde (apagado, ou o código do cenário renomeado), e o mapa segue
			// dizendo que a regra está provada. O `stale` não cobra a unidade,
			// porque a prova velha responde por ela.
			if suite == "" {
				if len(pc) > 0 || (n.Signal != nil && len(n.Signal.ProvenCodes) > 0) {
					ensureSignal(n)
					n.Signal.ProvenCodes = pc
					n.Signal.AtRev = n.Rev
					n.Signal.IngestedAt = now
					matchedCodes += len(pc)
				}
				continue
			}
			var previous []string
			if n.Signal != nil {
				previous = n.Signal.ProvenBySuite[suite]
			}
			if len(pc) > 0 || len(previous) > 0 {
				ensureSignal(n)
				if n.Signal.ProvenBySuite == nil {
					n.Signal.ProvenBySuite = map[string][]string{}
				}
				if n.Signal.ProvenRevBySuite == nil {
					n.Signal.ProvenRevBySuite = map[string]string{}
				}
				if len(pc) > 0 {
					n.Signal.ProvenBySuite[suite] = pc
					n.Signal.ProvenRevBySuite[suite] = n.Rev
				} else {
					// The suite stopped proving anything here: it leaves the union, and
					// leaves no empty entry behind in the versioned map.
					delete(n.Signal.ProvenBySuite, suite)
					delete(n.Signal.ProvenRevBySuite, suite)
				}
				n.Signal.ProvenCodes = unionProven(n.Signal.ProvenBySuite)
				n.Signal.AtRev = unionRev(n.Signal, n.Rev)
				n.Signal.IngestedAt = now
				matchedCodes += len(pc)
			}
		}
	}
	return
}

func (g *Graph) ingestCoverageBySuite(byFile map[string]FileCov, suite, now string) (matched int) {
	for i := range g.Nodes {
		n := &g.Nodes[i]
		if n.Kind != KindCode {
			continue
		}
		for file, cov := range byFile {
			if !pathMatches(n.ID, file) {
				continue
			}
			ensureSignal(n)
			if n.Signal.TotalLines > 0 {
				n.Signal.PrevLineCoverage = n.Signal.LineCoverage
			}
			if n.Signal.CoverageBySuite == nil {
				n.Signal.CoverageBySuite = map[string]SuiteCoverage{}
			}
			var instrumented, covered []int
			for line, hit := range cov.Lines {
				instrumented = append(instrumented, line)
				if hit {
					covered = append(covered, line)
				}
			}
			n.Signal.CoverageBySuite[suite] = SuiteCoverage{
				Instrumented: encodeRanges(instrumented),
				Covered:      encodeRanges(covered),
				CoveredLines: cov.Covered,
				TotalLines:   cov.Total,
				AtRev:        n.Rev,
			}
			if cov.Predates {
				e := n.Signal.CoverageBySuite[suite]
				e.AtRev = ""
				n.Signal.CoverageBySuite[suite] = e
			}
			c, t := unionCoverage(n.Signal.CoverageBySuite, n.Rev)
			n.Signal.CoveredLines, n.Signal.TotalLines = c, t
			if t > 0 {
				n.Signal.LineCoverage = float64(c) / float64(t) * 100
			}
			n.Signal.AtRev = oldestSuiteRev(n.Signal.CoverageBySuite, n.Rev)
			n.Signal.IngestedAt = now
			matched++
			break
		}
	}
	return
}

// unionCoverage computes covered/total over the union of the suites' lines.
//
// A suite that reported only totals (LF/LH, no per-line detail) cannot join a line union.
// When any suite lacks detail, the union falls back to the suite with the most covered
// lines — an UNDER-estimate, never an over-estimate: summing counts would count a line
// twice when two suites cover it.
//
// Only suites measured at the CURRENT rev take part. A suite recorded before the file was
// edited speaks of another text: its line numbers no longer point at these lines, and
// joining them mixes two numberings. Measured in the reference project: an edited handler,
// unit fresh at 73/75, integration stale at 0/76 in the old numbering — the union read
// 73/106 = 69% and failed the threshold on a file the fresh suite covered at 97%. The
// stale suite stays stored (it is refreshed when it runs again) and still keeps the node
// stale through `oldestSuiteRev`; it just says nothing about the current lines.
func unionCoverage(bySuite map[string]SuiteCoverage, currentRev string) (covered, total int) {
	fresh := map[string]SuiteCoverage{}
	for s, sc := range bySuite {
		if sc.AtRev == currentRev {
			fresh[s] = sc
		}
	}
	bySuite = fresh
	instr, cov := map[int]bool{}, map[int]bool{}
	detailed := true
	for _, sc := range bySuite {
		if sc.Instrumented == "" && sc.TotalLines > 0 {
			detailed = false
		}
		for _, l := range decodeRanges(sc.Instrumented) {
			instr[l] = true
		}
		for _, l := range decodeRanges(sc.Covered) {
			cov[l] = true
		}
	}
	if detailed {
		return len(cov), len(instr)
	}
	for _, sc := range bySuite {
		if sc.CoveredLines > covered || (sc.CoveredLines == covered && sc.TotalLines > total) {
			covered, total = sc.CoveredLines, sc.TotalLines
		}
	}
	return covered, total
}

// oldestSuiteRev: the union is only as fresh as its stalest suite (see unionRev).
func oldestSuiteRev(bySuite map[string]SuiteCoverage, current string) string {
	suites := make([]string, 0, len(bySuite))
	for s := range bySuite {
		suites = append(suites, s)
	}
	sort.Strings(suites)
	for _, s := range suites {
		if rev := bySuite[s].AtRev; rev != current {
			if rev == "" {
				return "unknown"
			}
			return rev
		}
	}
	return current
}

// encodeRanges writes sorted line numbers as compact ranges: `1-5,9,12-20`.
func encodeRanges(lines []int) string {
	if len(lines) == 0 {
		return ""
	}
	sort.Ints(lines)
	var b strings.Builder
	start, prev := lines[0], lines[0]
	flush := func() {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		if start == prev {
			b.WriteString(strconv.Itoa(start))
		} else {
			b.WriteString(strconv.Itoa(start) + "-" + strconv.Itoa(prev))
		}
	}
	for _, l := range lines[1:] {
		if l == prev || l == prev+1 {
			prev = l
			continue
		}
		flush()
		start, prev = l, l
	}
	flush()
	return b.String()
}

// decodeRanges reads what encodeRanges wrote.
func decodeRanges(s string) []int {
	var out []int
	for _, part := range strings.Split(s, ",") {
		if part == "" {
			continue
		}
		lo, hi, isRange := strings.Cut(part, "-")
		a, err := strconv.Atoi(lo)
		if err != nil {
			continue
		}
		b := a
		if isRange {
			if b, err = strconv.Atoi(hi); err != nil {
				continue
			}
		}
		for l := a; l <= b; l++ {
			out = append(out, l)
		}
	}
	return out
}

// unionProven: the codes proven by ANY suite, without repetition and in a stable order
// (the map is versioned; map order would change the file on every ingestion).
func unionProven(bySuite map[string][]string) []string {
	visto := map[string]bool{}
	var out []string
	for _, codes := range bySuite {
		for _, c := range codes {
			if !visto[c] {
				visto[c] = true
				out = append(out, c)
			}
		}
	}
	sort.Strings(out)
	return out
}

// unionRev is the rev the union of proofs can honestly claim.
//
// The union is only as fresh as its OLDEST contributor. If every suite that still proves
// something measured the current rev, the union is current. If any of them measured an
// earlier rev, that rev is returned, and `SignalStale` reports the node stale until that
// suite runs again — which is the truth: part of what the map calls proven was measured
// against a spec that no longer exists.
//
// An entry with no recorded rev (written before `ProvenRevBySuite` existed) has unknown
// freshness, and unknown is not fresh: it returns a rev that can never match.
func unionRev(sig *TestSignal, current string) string {
	suites := make([]string, 0, len(sig.ProvenBySuite))
	for s := range sig.ProvenBySuite {
		suites = append(suites, s)
	}
	sort.Strings(suites) // stable: the map is versioned
	for _, s := range suites {
		rev, ok := sig.ProvenRevBySuite[s]
		if !ok {
			return "unknown"
		}
		if rev != current {
			return rev
		}
	}
	return current
}

// IngestCoverage grava a cobertura de linha nos nós de CÓDIGO, casando por caminho.
func (g *Graph) IngestCoverage(byFile map[string]FileCov, now string) (matched int) {
	return g.IngestCoverageSuite(byFile, "", now)
}

// IngestCoverageSuite is IngestCoverage with the SUITE named (the lcov report it came
// from). With a suite, each suite's measurement is kept on its own and the node's coverage
// is the UNION: a line covered by the unit suite stays covered when the integration suite
// is ingested after it. An empty suite keeps the previous behaviour — the report is the
// whole measurement.
func (g *Graph) IngestCoverageSuite(byFile map[string]FileCov, suite, now string) (matched int) {
	if suite != "" {
		return g.ingestCoverageBySuite(byFile, suite, now)
	}
	for i := range g.Nodes {
		n := &g.Nodes[i]
		if n.Kind != KindCode {
			continue
		}
		for file, cov := range byFile {
			if pathMatches(n.ID, file) {
				ensureSignal(n)
				// preserva a cobertura anterior como baseline do delta (só quando já
				// havia uma medição real, para não criar um baseline falso de 0).
				if n.Signal.TotalLines > 0 {
					n.Signal.PrevLineCoverage = n.Signal.LineCoverage
				}
				n.Signal.CoveredLines = cov.Covered
				n.Signal.TotalLines = cov.Total
				if cov.Total > 0 {
					n.Signal.LineCoverage = float64(cov.Covered) / float64(cov.Total) * 100
				}
				n.Signal.AtRev = n.Rev
				n.Signal.IngestedAt = now
				matched++
				break
			}
		}
	}
	return
}

// IngestMutation amarra o sinal de MUTAÇÃO aos nós de código, no mesmo molde do lcov:
// o projeto rodou a ferramenta dele, o Anchors só lê o relatório padrão e o pendura no
// nó. Ver a doutrina em internal/testsig/mutation.go.
func (g *Graph) IngestMutation(byFile map[string]FileMutation, now string) (matched int) {
	return g.IngestMutationScoped(byFile, "", now, 0, 0)
}

// IngestMutationScoped grava o sinal de mutação num ESCOPO nomeado, além dos totais.
//
// O escopo diz QUE SUÍTE rodou contra os mutantes:
//
//	isolated — só o teste da própria unidade
//	full     — os testes de todos os que a importam
//
// Sem ele, a segunda ingestão do mesmo arquivo apagava a primeira, e a informação
// mais útil — a DIFERENÇA entre os dois — era impossível de calcular. Escopo vazio
// mantém o comportamento anterior (grava só os totais), para quem já ingere hoje.
//
// Os totais continuam sendo os da ÚLTIMA ingestão: eles são o que os checks legados
// leem, e mudá-los para uma média silenciaria a régua que o projeto já calibrou.
func (g *Graph) IngestMutationScoped(byFile map[string]FileMutation, scope, now string, low, high float64) (matched int) {
	for i := range g.Nodes {
		n := &g.Nodes[i]
		if n.Kind != KindCode {
			continue
		}
		for file, mu := range byFile {
			if pathMatches(n.ID, file) {
				ensureSignal(n)
				n.Signal.MutantsKilled = mu.Killed
				n.Signal.MutantsSurvived = mu.Survived
				n.Signal.MutantsNoCoverage = mu.NoCoverage
				n.Signal.MutantsIgnored = mu.Ignored
				n.Signal.MutationScore = mu.Score
				n.Signal.MutationLow, n.Signal.MutationHigh = low, high
				if scope != "" {
					if n.Signal.MutationByScope == nil {
						n.Signal.MutationByScope = map[string]MutationScope{}
					}
					n.Signal.MutationByScope[scope] = MutationScope{
						AtRev:  n.Rev,
						Killed: mu.Killed, Survived: mu.Survived,
						NoCoverage: mu.NoCoverage, Ignored: mu.Ignored,
						Score: mu.Score,
					}
				}
				n.Signal.AtRev = n.Rev
				n.Signal.IngestedAt = now
				matched++
				break
			}
		}
	}
	return
}

// FileMutation é o resultado de mutação de um arquivo (desacopla o mapx do testsig).
type FileMutation struct {
	Killed   int
	Survived int
	// NoCoverage: mutantes que nenhum teste executou. Fora do score — quem cobra
	// "não há teste aqui" é o gate de cobertura.
	NoCoverage int
	// Ignored: mutantes que a ferramenta descartou antes de rodar. Fora do score, e
	// gravados a parte — sao eles que separam "100% porque tudo foi provado" de "100%
	// porque nao havia o que provar".
	Ignored int
	Score   float64
}

// FileCov é a cobertura de um arquivo (desacopla o mapx do pacote testsig).
type FileCov struct {
	Covered, Total int
	// Lines: per instrumented line, covered or not. Empty when the report gave only totals.
	Lines map[int]bool
	// Predates: the report was written BEFORE the file's current content — its line
	// numbers describe another text. The entry is kept, with no rev, so it neither joins
	// the union nor passes for fresh.
	Predates bool
}

// SignalStale diz se o sinal de um nó envelheceu — o arquivo mudou de rev desde a
// ingestão. Um sinal stale não deve ser confiado (o teste rodou numa versão antiga).
func (n Node) SignalStale() bool {
	return n.Signal != nil && n.Signal.AtRev != "" && n.Signal.AtRev != n.Rev
}

func sumLayers(byLayer map[string]LayerExec) (p, f, s int) {
	for _, le := range byLayer {
		p += le.Passed
		f += le.Failed
		s += le.Skipped
	}
	return
}

func ensureSignal(n *Node) {
	if n.Signal == nil {
		n.Signal = &TestSignal{}
	}
}

// pathMatches: o caminho do nó (relativo à raiz) casa o caminho do relatório? Casa
// por sufixo normalizado — o relatório pode ter caminho absoluto ou com prefixo de
// subdir; o nó é sempre relativo à raiz. Exige que o sufixo comece numa fronteira de
// caminho (para "a/b.go" não casar "xa/b.go").
func pathMatches(nodeID, reportPath string) bool {
	a := filepath.ToSlash(nodeID)
	b := filepath.ToSlash(reportPath)
	if a == b {
		return true
	}
	// o relatório termina com o caminho do nó?
	if hasPathSuffix(b, a) {
		return true
	}
	// ou o nó termina com o do relatório (relatório mais curto)?
	if hasPathSuffix(a, b) {
		return true
	}
	return false
}

func hasPathSuffix(full, suffix string) bool {
	if len(suffix) > len(full) {
		return false
	}
	if full == suffix {
		return true
	}
	if full[len(full)-len(suffix):] != suffix {
		return false
	}
	// a posição antes do sufixo deve ser um separador (fronteira de caminho)
	return full[len(full)-len(suffix)-1] == '/'
}

// ResolveReportPaths decides, BEFORE ingestion, which node each path of the report
// belongs to — and returns the path rewritten to the node's exact ID.
//
// Why it exists: matching is by SUFFIX, and monorepo runners write paths relative to
// their own workspace (`src/components/atoms/SectionLabel.tsx`). With two workspaces
// holding the same file, the suffix matches BOTH nodes and both received the signal — the
// landing's coverage showed up on the mobile component of the same name, and the next
// ingestion flipped it. MEASURED in MIF (2026-09-23): 100 files in the lcov "matched" 101
// nodes.
//
// Tie-break: the node sharing the longest directory prefix with the report ITSELF
// (`reportRel`, relative to the root). The `unit.xml` under
// `apps/landing-page/test-output/` belongs to the landing. If it still ties, the path is
// left with NO owner and comes back in `ambiguous`: assigning it to one of the two on a
// guess would assert a proof nobody measured.
func (g *Graph) ResolveReportPaths(kind Kind, paths []string, reportRel string) (resolved map[string]string, ambiguous []string) {
	resolved = map[string]string{}
	hint := filepath.ToSlash(reportRel)
	for _, p := range paths {
		var cands []string
		for _, n := range g.Nodes {
			if n.Kind == kind && pathMatches(n.ID, p) {
				cands = append(cands, n.ID)
			}
		}
		switch len(cands) {
		case 0:
			resolved[p] = p // nothing matches: kept as it came (ingestion just ties nothing)
		case 1:
			resolved[p] = cands[0]
		default:
			best, bestLen, tie := "", -1, false
			for _, c := range cands {
				l := commonDirPrefix(c, hint)
				switch {
				case l > bestLen:
					best, bestLen, tie = c, l, false
				case l == bestLen:
					tie = true
				}
			}
			if tie {
				ambiguous = append(ambiguous, p)
				continue
			}
			resolved[p] = best
		}
	}
	return
}

// commonDirPrefix conta quantos segmentos de diretório iniciais `a` e `b` compartilham.
func commonDirPrefix(a, b string) int {
	sa, sb := strings.Split(filepath.ToSlash(a), "/"), strings.Split(filepath.ToSlash(b), "/")
	n := 0
	for n < len(sa)-1 && n < len(sb)-1 && sa[n] == sb[n] {
		n++
	}
	return n
}
