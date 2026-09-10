package mapx

import "path/filepath"

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
			if len(pc) > 0 {
				ensureSignal(n)
				n.Signal.ProvenCodes = pc
				n.Signal.AtRev = n.Rev
				n.Signal.IngestedAt = now
				matchedCodes += len(pc)
			}
		}
	}
	return
}

// IngestCoverage grava a cobertura de linha nos nós de CÓDIGO, casando por caminho.
func (g *Graph) IngestCoverage(byFile map[string]FileCov, now string) (matched int) {
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
