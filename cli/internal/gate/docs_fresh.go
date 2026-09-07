package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
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

// checkDocsFresh confronta o `docs/*.md` contra o que os templates produziriam AGORA.
func checkDocsFresh(_ string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "o compilado deriva das SPECS — é delas que o gate parte"
	}
	if _, err := os.Stat(filepath.Join(root, doct.Dir)); err != nil {
		return Skip, "o projeto não usa `" + doct.Dir + "/` — a documentação não é compilada"
	}

	c, err := doct.New(root, g)
	if err != nil {
		return Skip, "não foi possível ler as specs: " + err.Error()
	}

	// COMPARA em memória. Escrever aqui seria o gate consertando o que ele deveria
	// apontar: o `check` roda em hook e em CI, e um gate que altera arquivos torna o
	// resultado dependente de ter rodado antes — a segunda execução sempre passaria.
	defasados, err := c.Stale()
	if err != nil {
		return Fail, "os templates não compilam: " + err.Error() +
			"\n  Enquanto não compilarem, a documentação inteira está congelada no último build."
	}
	if len(defasados) == 0 {
		return Pass, ""
	}

	return Fail, fmt.Sprintf(
		"%d documento(s) fora de data: %s.\n"+
			"  O conteúdo mora na spec; o `%s/` é cópia derivada, e esta envelheceu — "+
			"a documentação afirma a versão ANTIGA, com conteúdo real e convincente.\n"+
			"  Rode `anchors docs build`.",
		len(defasados), strings.Join(defasados, ", "), doct.OutDir)
}
