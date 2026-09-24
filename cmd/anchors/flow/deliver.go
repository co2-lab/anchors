package flow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"

	"github.com/co2-lab/anchors/internal/change"
	"github.com/spf13/cobra"
)

// `anchors deliver` fecha uma etapa: registra O QUE foi entregue, para que o review tenha
// escopo e a intenção declarada tenha com o que ser confrontada.
//
// O ciclo antes disto terminava quando o código nascia. Medido em três rodadas de um E2E
// real: 7 defeitos graves — perda silenciosa de dado do usuário, tabela inexistente na
// infra, regra sem teste que a prove, contradição entre duas regras da mesma spec —
// passaram com TODOS os gates verdes. Nenhum foi achado por gate. Os 7 vieram de revisão
// adversarial, e o review só aconteceu porque alguém lembrou de pedir.
//
// O registro resolve as duas metades disso: o watcher vê o arquivo (ninguém precisa
// lembrar) e o revisor recebe escopo (não precisa adivinhar o que mudou).
func newDeliverCmd() *cobra.Command {
	var root, stage, unit, intent, agent, date string
	var files, decisions, uncovered []string
	var card int

	cmd := &cobra.Command{
		Use:   "deliver",
		Short: "Record the delivery of a stage — the review trigger",
		Long: `Closes a stage of the cycle by recording what was delivered in ` + "`changes/`" + `.

The record has three parts, and the last two are what makes it useful:

  --intent      what you say you did (the DECLARED half of the confrontation)
  --decision    what the ruler did NOT decide and you chose alone
  --uncovered   what you KNOW is not proved

The last two may look optional and are not: a silent decision is the origin of most
of the defects that cross the gates (everything exists, everything references
everything, and someone chose alone). Declaring them is what separates assumed debt
from forgetting.

Example:

  anchors deliver --stage code --unit src/business-logic/pricing.ts \
    --file src/business-logic/pricing.ts --file src/business-logic/pricing.feature \
    --intent "implements PRICX-B01..B04; the carry-forward uses lexicographic comparison" \
    --decision "month outside YYYY-MM: I assumed valid input (the spec does not decide)" \
    --uncovered "PRICX-B04 has no scenario — the edge case depends on phase 2"

After this, the watcher queues the review task.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if stage == "" || unit == "" {
				return fmt.Errorf("provide --stage (spec|code|feature|test|plan) and --unit <file of the unit>")
			}
			switch stage {
			case "spec", "code", "feature", "test", "plan":
			default:
				return fmt.Errorf("unknown stage %q — use: spec, code, feature, test or plan", stage)
			}
			if strings.TrimSpace(intent) == "" {
				return fmt.Errorf("provide --intent: the record without the declared intent gives " +
					"the reviewer nothing to confront against the disk")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if date == "" {
				return fmt.Errorf("provide --date YYYY-MM-DD (Anchors does not read the clock: " +
					"the date is stamped by whoever records, so the record is reproducible)")
			}
			// A unidade tem de EXISTIR. `deliver` registra o que foi entregue; entrega de
			// arquivo inexistente é registro de trabalho que não aconteceu — e o comando
			// ainda emitia "PRÓXIMO PASSO: revise esta unidade", enfileirando review de
			// nada.
			//
			// O `anchors check --changed` já recusava o mesmo caminho ("não existe no disco
			// nem no mapa"). Dois comandos do mesmo binário, o mesmo argumento, políticas
			// contrárias, e nenhuma régua dizendo qual manda — a forma mais barata do modo
			// de falha mais caro deste framework.
			// A unidade tem de existir NO DISCO — mas a peça que existe depende da ETAPA.
			// Na etapa `spec`, o `.ts` do alvo por definição ainda não nasceu (o próprio
			// `work` o chama de "onde nasce" e proíbe escrevê-lo), então exigir o alvo
			// trava a primeira entrega de toda unidade nova. Foi o que aconteceu: o
			// worker teve de passar o `.spec.md` como `--unit`, partindo o ledger da
			// unidade em duas identidades diferentes dentro de `changes/`.
			//
			// A régua: aceita o alvo OU qualquer peça da mesma unidade que já exista. O
			// que não se aceita é registrar entrega de unidade que não existe em peça
			// nenhuma — aí não há trabalho, e o `deliver` ainda mandaria revisá-lo.
			relUnit := common.RelTo(absRoot, unit)
			if peca, ok := existingPiece(absRoot, relUnit); ok {
				if peca != relUnit {
					fmt.Printf("   (recording the unit by `%s`, the piece that already exists in this stage)\n", peca)
				}
			} else {
				return fmt.Errorf("%q does not exist on disk, and no piece of this unit "+
					"(spec/feature/test/code) exists — `deliver` records what was "+
					"delivered, and there is no delivery without a file; check the path", unit)
			}
			if len(files) == 0 {
				files = []string{unit}
			}

			c := change.Change{
				Stage: stage, Unit: common.RelTo(absRoot, unit), Files: files,
				Intent: intent, Decisions: decisions, Uncovered: uncovered,
				Date: date, Agent: agent,
			}
			// O MODO decide ONDE o registro vive, e os dois são excludentes.
			//
			// No local o arquivo É o mecanismo: o watcher o vê aparecer e enfileira o
			// review. No modo `github` esse watcher não é quem move nada — quem move o
			// card é o pipeline, e quem revisa está lendo a ISSUE.
			//
			// Medido no projeto de referência: 73 registros no repositório, nenhum
			// revisado, e as issues sem a informação. O revisor não sabia que o arquivo
			// existia; o watcher, que o veria, não estava rodando — porque naquele modo
			// ele não é o mecanismo. É a mesma regra que o `next` aprendeu na v0.1.55.
			cfg, _ := config.Load(filepath.Join(absRoot, config.DefaultFile))
			// A VENDORED pipeline has no card here: its code, its spec and its board live
			// upstream, in the Anchors project that seeded it. Demanding the unit's code
			// refused every delivery touching one — measured in blue-eyes, where the map
			// gives `anchors-claim.yml` no code, as it should.
			//
			// In `github` mode there is no card to carry the record, and writing it to
			// `changes/` instead would be the fallback between modes this command refuses
			// (see above). So nothing is recorded, and the output says why and what makes
			// the file the project's. In local mode `changes/` IS the mechanism, and the
			// record is written as for any unit.
			upstream := upstreamUnit(absRoot, c.Unit)
			if upstream && cfg != nil && cfg.GitHubMode() && card == 0 {
				fmt.Printf("✓ nothing to record: `%s` is an Anchors-seeded pipeline, owned upstream —\n"+
					"  it has no local code and no card. A change to it belongs to Anchors (`anchors doctor --fix`\n"+
					"  replaces the file whole); to make it this project's, remove its `%s` line.\n",
					c.Unit, scan.UpstreamMarker)
				return nil
			}
			if cfg != nil && cfg.GitHubMode() {
				if err := deliverToBoard(absRoot, cfg, c, card); err != nil {
					return err
				}
			} else {
				p, err := change.Save(absRoot, c)
				if err != nil {
					return fmt.Errorf("write the record: %w", err)
				}
				rel, _ := filepath.Rel(absRoot, p)
				fmt.Printf("✓ delivery recorded: %s\n", rel)
			}

			// O REVIEW é a etapa que fecha o ciclo, e a que mais some. Medido: um agente
			// registrou 7 entregas corretamente e o review nunca aconteceu — ele não tinha
			// o watcher ligado, e a fila é um mecanismo de PULL: sem alguém puxando, o
			// registro fica parado. Dizer "próximo passo" em letra pequena não bastou.
			//
			// Então o comando não sugere: ele INSTRUI, com o comando pronto, e diz o que
			// acontece se for pulado.
			fmt.Println("\n── NEXT STEP: REVIEW this unit ──────────────────────────────")
			fmt.Printf("   anchors work review --for %s\n\n", c.Unit)
			fmt.Println("   Green gates do NOT prove it is right — they confront what is")
			fmt.Println("   DECLARABLE. In three rounds of a real E2E, 7 serious defects (silent")
			fmt.Println("   data loss, a rule with no test that proves it, contradiction between two")
			fmt.Println("   rules of the same spec) passed with everything green. None was found by a")
			fmt.Println("   gate; the 7 came from adversarial review.")
			if !watcherActive(absRoot) {
				fmt.Println("\n   (the watcher is not running — with `anchors watch start` this")
				fmt.Println("    delivery enters the queue on its own, and `anchors next` pulls it)")
			}
			fmt.Println("─────────────────────────────────────────────────────────────")
			// Confronta o declarado contra o disco ANTES de o registro virar material do
			// revisor. Ver deliver_confront.go: as duas checagens nasceram de divergências
			// reais medidas na primeira rodada em que este fluxo funcionou.
			confrontDelivery(absRoot, files, c.Unit)

			if len(decisions) == 0 && len(uncovered) == 0 {
				fmt.Println("  note: you declared ZERO free decisions and ZERO proof gaps.\n" +
					"  That ASSERTS that the ruler decided everything and that every rule has a test that\n" +
					"  exercises it — the reviewer will confront exactly that assertion.")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&stage, "stage", "", "REQUIRED — stage delivered: spec|code|feature|test|plan")
	cmd.Flags().StringVar(&unit, "unit", "", "REQUIRED — the file that identifies the unit of purpose")
	cmd.Flags().StringSliceVar(&files, "file", nil, "file touched (repeatable)")
	cmd.Flags().StringVar(&intent, "intent", "", "REQUIRED — what you claim to have done")
	// `StringArray`, e NÃO `StringSlice`: estas duas flags recebem PROSA, e o
	// `StringSlice` do pflag divide o valor na vírgula.
	//
	// Medido no blue-eyes, entregando a spec do DataStore:
	//
	//	--decision "oito regras e dois invariantes, na letra B/I que as vizinhas usam"
	//
	// virou DUAS decisões no registro — "oito regras e dois invariantes" e " na letra
	// B/I que as vizinhas usam", a segunda começando com espaço e sem sujeito. De cinco
	// decisões declaradas saíram nove itens, quatro deles fragmentos.
	//
	// O dano é sobre o que o registro existe para fazer: o revisor confronta cada
	// decisão contra o disco, e meia frase não é confrontável. Pior, a contagem infla —
	// "nove decisões" descreve um trabalho que tomou cinco.
	//
	// Vírgula em prosa é pontuação, não separador. `--file` continua `StringSlice`
	// porque caminho de arquivo não tem vírgula, e ali a divisão é conveniência real.
	cmd.Flags().StringArrayVar(&decisions, "decision", nil, "choice the ruler did not decide (repeatable)")
	cmd.Flags().StringArrayVar(&uncovered, "uncovered", nil, "what you know is not proven (repeatable)")
	cmd.Flags().StringVar(&date, "date", "", "YYYY-MM-DD (REQUIRED — Anchors does not read the clock)")
	cmd.Flags().StringVar(&agent, "agent", "", "who delivered (optional)")
	cmd.Flags().IntVar(&card, "card", 0, "github mode: the card to record on, instead of finding it by the unit's code "+
		"(a plan card with no [CODE] in its title, or a change to generated files only)")
	return cmd
}

// watcherActive diz se o daemon está rodando neste projeto. Serve só para não sugerir
// ligar o que já está ligado — ruído em instrução é o que faz a instrução ser ignorada.
func watcherActive(root string) bool {
	_, err := os.Stat(filepath.Join(root, ".anchors", "watch.meta"))
	return err == nil
}

// existingPiece devolve a peça da unidade que está no disco: o próprio alvo, se existir,
// ou a primeira peça derivada dele que exista. É o que permite `deliver --stage spec`
// funcionar na primeira entrega, quando só o `.spec.md` nasceu.
func existingPiece(root, rel string) (string, bool) {
	if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
		return rel, true
	}
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx", ".ts", ".tsx"} {
		cand := base + suf
		if cand == rel {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, cand)); err == nil {
			return cand, true
		}
	}
	return "", false
}

// upstreamUnit says whether the delivered unit is a vendored Anchors pipeline — read from
// the file itself, not from the map, which may be older than the file.
func upstreamUnit(root, unit string) bool {
	b, err := os.ReadFile(filepath.Join(root, unit))
	if err != nil {
		return false
	}
	return scan.IsUpstreamOwned(unit, b)
}

// deliverToBoard posta o registro de entrega como comentário na issue da unidade.
//
// FALHA se não achar a issue, e não cai para o arquivo em silêncio: um registro que
// deveria estar na issue e foi parar no disco é invisível para quem revisa — exatamente o
// defeito que este caminho existe para fechar. Melhor o comando parar e dizer o que falta.
//
// `--card` names the card directly. The code lookup cannot find two real cases, and both
// were measured in blue-eyes: a PLAN card (`[plano] …`) has no `[CODE]` in its title, and a
// change to generated files only (`docs/`, the map) has no unit in the map. Agents fell
// back to writing the record by hand, outside the format the reviewer reads.
func deliverToBoard(root string, cfg *config.Config, c change.Change, cardNumber int) error {
	cli := board.Client{Repo: cfg.Workflow.Repo, Labels: cfg.Workflow.Labels}
	if cardNumber > 0 {
		card, err := cli.FindOpenByNumber(cardNumber)
		if err != nil {
			return err
		}
		if err := cli.Comment(card.Number, c.Render()); err != nil {
			return err
		}
		fmt.Printf("✓ delivery recorded on issue #%d — %s\n", card.Number, card.Title)
		return nil
	}
	codigo := codeOfUnit(root, c.Unit)
	if codigo == "" {
		return fmt.Errorf("could not find the CODE of unit `%s` in the map.\n"+
			"  In `github` mode the delivery record goes to the ISSUE, and it is the code that\n"+
			"  identifies it. Run `anchors map build` and check whether the unit is there —\n"+
			"  or name the card with `--card <n>`", c.Unit)
	}
	card, err := cli.FindByCode(codigo)
	if err != nil {
		return fmt.Errorf("%w.\n"+
			"  In `github` mode the delivery is recorded on the card issue, and it must\n"+
			"  exist and be OPEN. If the card was already closed, reopen it — or record the\n"+
			"  delivery on the card that covers this work with `--card <n>`", err)
	}
	if err := cli.Comment(card.Number, c.Render()); err != nil {
		return err
	}
	fmt.Printf("✓ delivery recorded on issue #%d — %s\n", card.Number, card.Title)
	return nil
}

// codeOfUnit devolve o código da unidade a que um arquivo pertence.
//
// Do MAPA, porque é ele que sabe: o código pode estar no header da spec, ser inferido do
// texto, ou vir da âncora irmã de um derivado — três regras que o `mapx` já resolve, e
// reimplementá-las aqui as faria divergir na primeira mudança.
func codeOfUnit(root, unit string) string {
	g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	if err != nil {
		return ""
	}
	unit = filepath.ToSlash(unit)
	// Casa o arquivo exato primeiro; depois qualquer peça da mesma unidade (o `--unit`
	// aceita a spec quando o alvo ainda não nasceu, e o código é o mesmo).
	stem, _ := mapx.StemOfAnchor(unit)
	for _, n := range g.Nodes {
		if n.ID == unit && n.Code != "" {
			return n.Code
		}
	}
	for _, n := range g.Nodes {
		if s, _ := mapx.StemOfAnchor(n.ID); s == stem && n.Code != "" {
			return n.Code
		}
	}
	return ""
}
