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

// itemPlaceholderRE captura o item de progresso que NÃO cita caminho nenhum.
//
// O `anchors new progress` escreve `- [ ] TODO: um item por spec que esta fase semeia`
// quando a fase não semeia nada — é um convite a preencher, e deveria sair quando alguém
// decide o que a fase faz.
//
// Medido no blue-eyes: ficou no progresso do plano 0017, fase `MTUAO-F02`. E as três
// direções deste gate não o veem — as duas primeiras confrontam itens que CITAM CAMINHO,
// e a terceira olha as sementes do plano (aquela fase não semeia).
//
// O efeito é o oposto do que o gate protege: um `[ ]` eterno faz o plano parecer
// incompleto para sempre. O `anchors next` volta a ele, e quem lê não sabe se falta
// trabalho ou falta limpar o arquivo.
var itemPlaceholderRE = regexp.MustCompile(
	`(?m)^[ \t]*-[ \t]+\[[ xX]\][ \t]*(TODO|FIXME|XXX|\.\.\.)\b.*$`)

// planSeedInListRE captura as specs que o plano SEMEIA — as do item de lista, não as
// mencionadas em prosa.
//
// O `plan-seeds-valid` casa qualquer “ `x.spec.md` “ do arquivo, e para o que ele faz
// isso basta: ele valida a FORMA do caminho, e validar de novo uma spec citada não custa
// nada. Aqui custaria — medido: a revisão `PLTFR-R0003` menciona “ `DataStore.spec.md` “
// em prosa (sem caminho), e o gate acusou o progresso de não listar uma spec que ele
// lista, com outro nome. Três dos quatro achados eram menções assim.
//
// A âncora é o `- [ ]` / `- [x]` no início da linha: o plano promete criar naquele item, e
// a prosa das revisões fala sobre o que já existe.
var planSeedInListRE = regexp.MustCompile("(?m)^[ \t]*-[ \t]+\\[[ xX]\\][ \t]+`([^`]+\\.spec\\.md)`")

// spec citada em QUALQUER lugar do progresso conta como listada: o item pode ter sido
// reescrito, movido de fase, ou anotado — o que importa é o arquivo estar lá.
var specInProgressRE = regexp.MustCompile("`([^`]+\\.spec\\.md)`")

// checkProgressHonest confronta cada item do `-progress.md` contra o disco.
func checkProgressHonest(planContent string, n mapx.Node, root string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
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

	// A TERCEIRA DIREÇÃO: semente que o plano promete e o progresso não lista.
	//
	// As duas primeiras conferem os itens QUE EXISTEM no progresso. Um gate que só olha
	// para dentro do arquivo nunca vê o que falta nele — e foi assim que este gate passou
	// com `✓15` num plano que acabara de ganhar uma spec semeada (o `MutualTls.spec.md`,
	// migrado do 0016 por `PLTFR-R0004`).
	//
	// O efeito é o mesmo que o `[x]` mentiroso: o progresso diz que a fase acabou, e o
	// `anchors next` seguiu para outro plano com trabalho declarado por fazer.
	naoListadas := seedsMissingFromProgress(string(conteudo), planContent)

	// A QUARTA DIREÇÃO: o item que não promete arquivo nenhum.
	var placeholders []string
	for _, m := range itemPlaceholderRE.FindAllString(string(conteudo), -1) {
		placeholders = append(placeholders, strings.TrimSpace(m))
	}

	if len(abertoMasFeito) == 0 && len(marcadoMasAusente) == 0 &&
		len(naoListadas) == 0 && len(placeholders) == 0 {
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
	if len(naoListadas) > 0 {
		fmt.Fprintf(&b, "%d spec(s) que o PLANO semeia e o progresso NÃO lista: %s.\n"+
			"  O progresso diz que a fase acabou, e há trabalho declarado por fazer — "+
			"o `anchors next` segue para outro plano.\n",
			len(naoListadas), strings.Join(naoListadas, ", "))
	}
	if len(placeholders) > 0 {
		fmt.Fprintf(&b, "%d item(ns) PLACEHOLDER que não citam arquivo: %s.\n"+
			"  O `anchors new progress` os escreve quando a fase não semeia nada, e eles "+
			"deveriam sair quando alguém decide o que a fase faz. Um `[ ]` eterno faz o "+
			"plano parecer incompleto para sempre — o `anchors next` volta a ele, e quem "+
			"lê não sabe se falta trabalho ou falta limpar o arquivo.\n",
			len(placeholders), strings.Join(placeholders, " / "))
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

// seedsMissingFromProgress devolve as specs que o plano promete e o progresso não lista.
//
// Compara por CAMINHO, não por linha: o texto do item no progresso pode ser reescrito
// (encurtado, traduzido) sem deixar de ser o mesmo item — o que identifica é o arquivo.
func seedsMissingFromProgress(progresso, plano string) []string {
	noProgresso := map[string]bool{}
	for _, m := range specInProgressRE.FindAllStringSubmatch(progresso, -1) {
		noProgresso[m[1]] = true
	}
	var faltam []string
	vistas := map[string]bool{}
	for _, m := range planSeedInListRE.FindAllStringSubmatch(plano, -1) {
		caminho := m[1]
		// MOLDES não são specs semeadas — é a mesma exceção do `plan-seeds-valid`.
		if strings.HasPrefix(filepath.Base(caminho), "_TEMPLATE_") {
			continue
		}
		if noProgresso[caminho] || vistas[caminho] {
			continue
		}
		vistas[caminho] = true
		faltam = append(faltam, caminho)
	}
	return faltam
}
