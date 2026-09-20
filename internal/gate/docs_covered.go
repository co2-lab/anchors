package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- toda spec tem de CHEGAR a alguma página ---
//
// O gate irmão (`docs-fresh`) pergunta se a página compilada reflete a spec de onde veio.
// Este pergunta o que aquele não tem como ver: a spec chega a ALGUMA página?
//
// São defeitos diferentes, e este é o silencioso. Os templates em `doct/` filtram por
// camada (`{{range specs "layer=gate"}}`), então uma spec que nenhum filtro seleciona
// compila para lugar nenhum. E nada acusa — porque todas as páginas que existem estão
// corretas, e `docs-fresh` confere exatamente isso. A unidade existe, tem trinca
// completa, passa por todos os gates relacionais, e não está documentada.
//
// Medido neste repositório ao escrever o gate: 51 specs, 51 documentadas, zero órfãs. É
// coincidência do estado atual — todas as specs são de `gate`, e há um template que pede
// `layer=gate`. A PRIMEIRA spec de outra camada sai da documentação sem uma palavra, e é
// justamente quando o projeto cresce para uma segunda camada que ninguém está olhando.
//
// INFORMATIVO pelo mesmo critério do irmão: quem conserta é um humano decidindo se amplia
// o filtro de um template ou cria a página da camada, e nenhuma das duas é trabalho de
// máquina. Mas ao contrário do irmão, `anchors docs build` NÃO conserta isto — o build
// gera as páginas que os templates pedem, e o defeito é justamente não haver quem peça.
var (
	uncoveredMu    sync.Mutex
	uncoveredRoot  string
	uncoveredGraph *mapx.Graph
	uncoveredOut   []string
	uncoveredErr   error
	uncoveredOK    bool
)

// uncoveredSpecs devolve as specs que nenhum template alcança, compilando no máximo uma
// vez por varredura — a resposta é a mesma para todos os alvos, pela mesma razão que no
// `docs-fresh` (ver o cache de lá).
func uncoveredSpecs(root string, g *mapx.Graph) ([]string, error) {
	uncoveredMu.Lock()
	defer uncoveredMu.Unlock()
	if uncoveredOK && uncoveredRoot == root && uncoveredGraph == g {
		return uncoveredOut, uncoveredErr
	}
	c, err := doct.New(root, g)
	if err != nil {
		uncoveredRoot, uncoveredGraph, uncoveredOut, uncoveredErr, uncoveredOK = root, g, nil, err, true
		return nil, err
	}
	out, err := c.Uncovered()
	uncoveredRoot, uncoveredGraph, uncoveredOut, uncoveredErr, uncoveredOK = root, g, out, err, true
	return out, err
}

// resetDocsCoverageCache clears the memory between scans — for the tests, like its sibling.
func resetDocsCoverageCache() {
	uncoveredMu.Lock()
	defer uncoveredMu.Unlock()
	uncoveredOK, uncoveredRoot, uncoveredGraph, uncoveredOut, uncoveredErr = false, "", nil, nil, nil
}

// checkDocsCovered confronta cada spec contra o que os templates ALCANÇAM.
func checkDocsCovered(_ string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.docs_covered.skip_not_spec")
	}
	if _, err := os.Stat(filepath.Join(root, doct.Dir)); err != nil {
		return Skip, fmt.Sprintf(i18n.T("gate.docs_covered.skip_no_doct_dir"), doct.Dir)
	}
	orphans, err := uncoveredSpecs(root, g)
	if err != nil {
		// Compilar e falhar é problema do template, e o gate irmão já o reporta com o
		// laudo do compilador. Repetir a mesma falha em dois gates faria o autor pensar
		// que há dois defeitos.
		return Skip, ""
	}
	// O veredito é sobre ESTE alvo: a lista é do conjunto, mas quem responde é a spec
	// confrontada. Reportar as órfãs das outras aqui acusaria um arquivo pelo que falta
	// noutro.
	for _, o := range orphans {
		if o == n.ID {
			return Fail, fmt.Sprintf(i18n.T("gate.docs_covered.uncovered"), 1, o)
		}
	}
	return Pass, ""
}
