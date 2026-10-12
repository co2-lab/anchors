// @anchors
//   code: PHGPR
//   ref: PRHNP

package gate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/progress"
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
// Medido no app de referência: ao criar os 17 progressos, transportei o estado dos checkboxes que
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

// progressItemRE captura um item de progresso que cita um caminho entre crases.
//
// Só o que cita CAMINHO é confrontável. Um item em prosa ("revisar com o time") fica de
// fora por construção — não há o que conferir, e reprovar por isso cobraria de quem
// escreveu um progresso legítimo.

// placeholderItemRE captura o item de progresso que NÃO cita caminho nenhum.
//
// O `anchors new progress` escreve `- [ ] TODO: um item por spec que esta fase semeia`
// quando a fase não semeia nada — é um convite a preencher, e deveria sair quando alguém
// decide o que a fase faz.
//
// Medido no app de referência: ficou no progresso do plano 0017, fase `MTUAO-W02`. E as três
// direções deste gate não o veem — as duas primeiras confrontam itens que CITAM CAMINHO,
// e a terceira olha as sementes do plano (aquela fase não semeia).
//
// O efeito é o oposto do que o gate protege: um `[ ]` eterno faz o plano parecer
// incompleto para sempre. O `anchors next` volta a ele, e quem lê não sabe se falta
// trabalho ou falta limpar o arquivo.

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
		return Skip, i18n.T("gate.progress_honest.skip_not_plan")
	}
	// The plan's progress is its demand in anchors.progress.md (DESIGN-progress-tracking.md):
	// the demand whose `plan:` is this plan.
	file, err := progress.Load(root)
	if err != nil {
		return Pending, i18n.T("gate.progress_honest.unreadable", err)
	}
	var d *progress.Demand
	for _, x := range file.Demands {
		if x.Plan == n.ID {
			d = x
			break
		}
	}
	if d == nil {
		return Skip, i18n.T("gate.progress_honest.skip_no_demand", n.ID)
	}

	var abertoMasFeito, marcadoMasAusente, placeholders []string
	var cited strings.Builder
	for _, t := range d.AllTasks() {
		if t.State == progress.Cancelled {
			continue
		}
		cited.WriteString(t.Title + "\n")
		if placeholderTitleRE.MatchString(t.Title) {
			placeholders = append(placeholders, t.ID+" "+t.Title)
			continue
		}
		m := taskPathRE.FindStringSubmatch(t.Title)
		if m == nil {
			continue
		}
		caminho := m[1]
		_, errStat := os.Stat(filepath.Join(root, caminho))
		existe := errStat == nil
		switch {
		case existe && t.State != progress.Done:
			abertoMasFeito = append(abertoMasFeito, t.ID+" "+caminho)
		case !existe && t.State == progress.Done:
			marcadoMasAusente = append(marcadoMasAusente, t.ID+" "+caminho)
		}
	}

	naoListadas := seedsMissingFromProgress(cited.String(), planContent)

	if len(abertoMasFeito) == 0 && len(marcadoMasAusente) == 0 &&
		len(naoListadas) == 0 && len(placeholders) == 0 {
		return Pass, ""
	}

	var b strings.Builder
	if len(marcadoMasAusente) > 0 {
		b.WriteString(i18n.T("gate.progress_honest.marked_but_missing",
			len(marcadoMasAusente), strings.Join(marcadoMasAusente, ", ")))
	}
	if len(abertoMasFeito) > 0 {
		b.WriteString(i18n.T("gate.progress_honest.open_but_exists",
			len(abertoMasFeito), strings.Join(abertoMasFeito, ", ")))
	}
	if len(naoListadas) > 0 {
		b.WriteString(i18n.T("gate.progress_honest.seeds_missing",
			len(naoListadas), strings.Join(naoListadas, ", ")))
	}
	if len(placeholders) > 0 {
		b.WriteString(i18n.T("gate.progress_honest.placeholder_items",
			len(placeholders), strings.Join(placeholders, " / ")))
	}
	b.WriteString(i18n.T("gate.progress_honest.guidance_mark_in_progress"))
	return Fail, b.String()
}

// taskPathRE is the file a task cites, between backticks: the work it delivers.
var taskPathRE = regexp.MustCompile("`([^`]+\\.(?:md|ts|tsx|js|jsx|json|ya?ml|feature|py|go|ml|mli|rs|rb|java|kt|cs|swift|php))`")

// placeholderTitleRE is a task that cites no work yet: a TODO left from a template.
var placeholderTitleRE = regexp.MustCompile(`^(TODO|FIXME|XXX|\.\.\.)\b`)

// checkProgressIntact: anchors.progress.md is the file Anchors wrote — it parses, and its
// sum matches its content. A hand edit breaks the format Anchors edits deterministically.
func checkProgressIntact(_ string, _ mapx.Node, root string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
	if _, err := os.Stat(progress.Path(root)); err != nil {
		return Skip, i18n.T("gate.progress_intact.skip_absent")
	}
	if _, err := progress.Load(root); err != nil {
		return Fail, err.Error()
	}
	return Pass, ""
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
