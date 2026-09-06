package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// --- o PROGRESSO tem de dizer a verdade sobre o disco ---
//
// O `-progress.md` é o único artefato do Anchors que nada confronta, e isso é deliberado:
// ele fica fora do mapa porque existe para MUDAR, e gates que cobram justificativa de
// mudança não podem alcançá-lo (ver `scan.IsProgressFile`).
//
// Mas "fora do mapa" virou "fora de qualquer verificação", e as duas coisas não precisam
// ser a mesma. O item que cita um CAMINHO de arquivo é verificável de forma trivial: o
// arquivo existe ou não existe.
//
// Medido no blue-eyes: ao criar os 17 progressos, transportei o estado dos checkboxes que
// viviam nos planos — e os planos estavam desatualizados. O progresso do `0002` dizia 6
// itens abertos com 7 das 8 specs já no disco. CINCO itens mentiam, e eu transportei a
// mentira fielmente.
//
// O dano tem duas direções, e a segunda é pior:
//
//   • `[ ]` com arquivo existindo   → retrabalho: alguém refaz o que está feito
//   • `[x]` com arquivo ausente     → plano declarado pronto com trabalho por fazer,
//                                     que é o mesmo defeito que o `Closes` explícito
//                                     fechou por outra porta
//
// O gate é ancorado no PLANO (que está no mapa) e confronta o companheiro. É a única forma
// de alcançar um arquivo que, por desenho, não é nó.

// itemDeProgressoRE captura um item de progresso que cita um caminho entre crases.
//
// Só o que cita CAMINHO é confrontável. Um item em prosa ("revisar com o time") fica de
// fora por construção — não há o que conferir, e reprovar por isso cobraria de quem
// escreveu um progresso legítimo.
var itemDeProgressoRE = regexp.MustCompile(
	"(?m)^[ \t]*-[ \t]+\\[([ xX])\\][ \t]+`([^`]+\\.(?:md|ts|tsx|js|jsx|json|ya?ml|feature|py|go))`")

// checkProgressHonest confronta cada item do `-progress.md` contra o disco.
func checkProgressHonest(_ string, n mapx.Node, root string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, "o progresso acompanha o PLANO — é dele que o gate parte"
	}
	// O companheiro vive ao lado, com sufixo fixo. A definição canônica é do `scan`, que
	// é quem precisa mantê-lo fora do mapa — repeti-la aqui deixaria as duas divergirem.
	prog := progressPathOf(n.ID)
	conteudo, err := os.ReadFile(filepath.Join(root, prog))
	if err != nil {
		// Ausência NÃO é falha deste gate: o plano pode ter nascido antes do mecanismo
		// (`anchors new progress --for` existe para isso), e acusar aqui misturaria duas
		// perguntas diferentes — "o progresso existe?" e "o progresso é verdade?".
		return Skip, "o plano não tem `" + prog + "` ao lado — crie com `anchors new progress --for " + n.ID + "`"
	}

	var abertoMasFeito, marcadoMasAusente []string
	for _, m := range itemDeProgressoRE.FindAllStringSubmatch(string(conteudo), -1) {
		marcado := strings.ToLower(m[1]) == "x"
		caminho := m[2]
		_, errStat := os.Stat(filepath.Join(root, caminho))
		existe := errStat == nil
		switch {
		case existe && !marcado:
			abertoMasFeito = append(abertoMasFeito, caminho)
		case !existe && marcado:
			marcadoMasAusente = append(marcadoMasAusente, caminho)
		}
	}

	if len(abertoMasFeito) == 0 && len(marcadoMasAusente) == 0 {
		return Pass, ""
	}

	var b strings.Builder
	// A ordem importa: o `[x]` sem arquivo vem primeiro porque é o dano mais caro —
	// ele declara pronto o que não está, e quem lê o board decide com base nisso.
	if len(marcadoMasAusente) > 0 {
		fmt.Fprintf(&b, "%d item(ns) marcado(s) `[x]` cujo arquivo NÃO EXISTE: %s.\n"+
			"  O progresso declara pronto o que não está — quem lê o board decide com "+
			"base nisso.\n",
			len(marcadoMasAusente), strings.Join(marcadoMasAusente, ", "))
	}
	if len(abertoMasFeito) > 0 {
		fmt.Fprintf(&b, "%d item(ns) `[ ]` cujo arquivo JÁ EXISTE: %s.\n"+
			"  O trabalho foi feito e ninguém marcou — o próximo a olhar refaz.\n",
			len(abertoMasFeito), strings.Join(abertoMasFeito, ", "))
	}
	b.WriteString("\nMarque no `-progress.md`, NUNCA no plano: alterar o plano significa " +
		"que a DECISÃO mudou, e cobra revisão (`{CODIGO}-R000N`).")
	return Fail, b.String()
}

// progressPathOf devolve o companheiro de um plano.
//
// Deriva do `scan.ProgressPathFor` para que o sufixo tenha UMA definição: o `scan` é quem
// precisa manter o arquivo fora do mapa, e uma segunda constante aqui poderia divergir
// dela em silêncio — o gate passaria a confrontar um arquivo que o scanner indexa, ou a
// procurar um que não existe.
func progressPathOf(plano string) string {
	return scan.ProgressPathFor(plano)
}
