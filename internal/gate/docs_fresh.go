package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- o COMPILADO tem de refletir a spec ---
//
// A documentação é gerada: `doct/*.md.tmpl` referencia trechos das specs, e `anchors docs
// build` produz `docs/*.md`. O conteúdo mora na spec e só nela — a cópia que vive no
// compilado é derivada, e derivada envelhece.
//
// Sem gate, o modo de falha é silencioso e conhecido: alguém altera uma regra na spec,
// esquece de recompilar, e a documentação segue afirmando a regra ANTIGA. Ninguém vê,
// porque o `.md` está lá, bonito, com conteúdo real. É pior que a página vazia — uma doc
// obviamente incompleta manda procurar a fonte; uma doc desatualizada convence.
//
// INFORMATIVO, e por decisão explícita: "não acho que precise ser bloqueante porque o
// bloqueio pode ser resolvido com um build". Um gate que reprova algo que um comando
// conserta sozinho gasta a atenção de quem revisa com trabalho de máquina — o pipeline
// roda o build no merge, e o gate serve para o autor ver antes.
//
// Ancorado na SPEC porque é ela que muda. O documento compilado não é nó do mapa (é
// gerado, e cobrar revisão dele seria cobrar revisão de uma saída de compilador), então
// não há de onde partir a não ser da fonte.

// A RESPOSTA E' A MESMA PARA TODAS AS SPECS — e custava uma recompilacao por spec.
//
// `Stale()` compila os templates INTEIROS (a documentacao das 51 specs) e compara com o
// que esta em `docs/`. O resultado depende de `(root, graph)` e de mais nada: confrontar
// a spec A ou a spec B produz a mesma lista de documentos defasados.
//
// O gate, porem, e' ancorado na spec (ver acima), entao roda UMA VEZ POR ALVO. Medido
// com `check --all --timing` neste repositorio, antes deste cache:
//
//	docs-fresh   6m37.21s   51 alvo(s)   pior 8.99s      <- 97% da varredura inteira
//
// Eram 51 compilacoes completas do mesmo conjunto, jogando fora 50. O cache guarda a
// primeira e devolve as outras 50.
//
// A chave inclui o GRAFO, e nao so' a raiz: duas varreduras da mesma sessao (o `check`
// roda antes e depois do `--fix`) enxergam mapas diferentes, e reusar a resposta da
// primeira responderia sobre um estado que ja mudou. Comparar o ponteiro basta porque o
// mapa e' carregado uma vez por varredura e nao e' mutado durante ela.
var (
	staleMu    sync.Mutex
	staleRoot  string
	staleGraph *mapx.Graph
	staleOut   []string
	staleErr   error
	staleOK    bool
)

// staleDocs returns the stale documents, compiling at most once per scan.
func staleDocs(root string, g *mapx.Graph) ([]string, error) {
	staleMu.Lock()
	defer staleMu.Unlock()
	if staleOK && staleRoot == root && staleGraph == g {
		return staleOut, staleErr
	}
	c, err := doct.New(root, g)
	if err != nil {
		staleRoot, staleGraph, staleOut, staleErr, staleOK = root, g, nil, err, true
		return nil, err
	}
	out, err := c.Stale()
	staleRoot, staleGraph, staleOut, staleErr, staleOK = root, g, out, err, true
	return out, err
}

// resetDocsCache clears the memory between scans. Existe para os TESTES: em
// producao a chave (raiz + grafo) ja invalida sozinha, mas um teste roda varias
// varreduras no mesmo processo e precisa comecar de um estado conhecido.
func resetDocsCache() {
	staleMu.Lock()
	defer staleMu.Unlock()
	staleOK, staleRoot, staleGraph, staleOut, staleErr = false, "", nil, nil, nil
}

// checkDocsFresh confronta o `docs/*.md` contra o que os templates produziriam AGORA.
func checkDocsFresh(_ string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.docs_fresh.skip_not_spec")
	}
	if _, err := os.Stat(filepath.Join(root, doct.Dir)); err != nil {
		return Skip, fmt.Sprintf(i18n.T("gate.docs_fresh.skip_no_doct_dir"), doct.Dir)
	}

	// COMPARA em memória. Escrever aqui seria o gate consertando o que ele deveria
	// apontar: o `check` roda em hook e em CI, e um gate que altera arquivos torna o
	// resultado dependente de ter rodado antes — a segunda execução sempre passaria.
	defasados, err := staleDocs(root, g)
	if err != nil {
		// O erro de LER as specs e o de COMPILAR chegam pelo mesmo caminho agora, e a
		// distincao importa: nao ter o que ler e' Skip (o gate nao mediu), e compilar e
		// falhar e' Fail (mediu e deu errado). O texto do erro de leitura vem de
		// `doct.New`, que so' falha quando uma spec do mapa nao esta no disco.
		if strings.Contains(err.Error(), "is in the map and not on disk") {
			return Skip, fmt.Sprintf(i18n.T("gate.docs_fresh.skip_read_specs_err"), err.Error())
		}
		return Fail, fmt.Sprintf(i18n.T("gate.docs_fresh.fail_compile_err"), err.Error())
	}
	if len(defasados) == 0 {
		return Pass, ""
	}

	return Fail, fmt.Sprintf(i18n.T("gate.docs_fresh.stale_docs"),
		len(defasados), strings.Join(defasados, ", "), doct.OutDir)
}
