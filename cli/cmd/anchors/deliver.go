package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"

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

	cmd := &cobra.Command{
		Use:   "deliver",
		Short: "Registra a entrega de uma etapa — o gatilho do review",
		Long: `Fecha uma etapa do ciclo registrando o que foi entregue em ` + "`changes/`" + `.

O registro tem três partes, e as duas últimas são o que o torna útil:

  --intent      o que você diz ter feito (a metade DECLARADA do confronto)
  --decision    o que a régua NÃO decidiu e você escolheu sozinho
  --uncovered   o que você SABE que não está provado

As duas últimas podem parecer opcionais e não são: decisão silenciosa é a origem da
maior parte dos defeitos que atravessam os gates (tudo existe, tudo se referencia, e
alguém escolheu sozinho). Declará-las é o que separa dívida assumida de esquecimento.

Exemplo:

  anchors deliver --stage code --unit src/business-logic/pricing.ts \
    --file src/business-logic/pricing.ts --file src/business-logic/pricing.feature \
    --intent "implementa PRICX-B01..B04; o carry-forward usa comparação lexicográfica" \
    --decision "mês fora de YYYY-MM: assumi entrada válida (a spec não decide)" \
    --uncovered "PRICX-B04 não tem cenário — o caso de borda depende da fase 2"

Depois disto, o watcher enfileira a task de review.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if stage == "" || unit == "" {
				return fmt.Errorf("informe --stage (spec|code|feature|test|plan) e --unit <arquivo da unidade>")
			}
			switch stage {
			case "spec", "code", "feature", "test", "plan":
			default:
				return fmt.Errorf("etapa desconhecida %q — use: spec, code, feature, test ou plan", stage)
			}
			if strings.TrimSpace(intent) == "" {
				return fmt.Errorf("informe --intent: o registro sem a intenção declarada não dá " +
					"ao revisor o que confrontar contra o disco")
			}
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			if date == "" {
				return fmt.Errorf("informe --date AAAA-MM-DD (o Anchors não lê o relógio: " +
					"a data é carimbada por quem registra, para o registro ser reproduzível)")
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
			relUnit := relTo(absRoot, unit)
			if peca, ok := pecaExistente(absRoot, relUnit); ok {
				if peca != relUnit {
					fmt.Printf("   (registrando a unidade por `%s`, a peça que já existe nesta etapa)\n", peca)
				}
			} else {
				return fmt.Errorf("%q não existe no disco, e nenhuma peça desta unidade "+
					"(spec/feature/teste/código) existe — `deliver` registra o que foi "+
					"entregue, e não há entrega sem arquivo; confira o caminho", unit)
			}
			if len(files) == 0 {
				files = []string{unit}
			}

			c := change.Change{
				Stage: stage, Unit: relTo(absRoot, unit), Files: files,
				Intent: intent, Decisions: decisions, Uncovered: uncovered,
				Date: date, Agent: agent,
			}
			p, err := change.Save(absRoot, c)
			if err != nil {
				return fmt.Errorf("gravar o registro: %w", err)
			}
			rel, _ := filepath.Rel(absRoot, p)
			fmt.Printf("✓ entrega registrada: %s\n", rel)

			// O REVIEW é a etapa que fecha o ciclo, e a que mais some. Medido: um agente
			// registrou 7 entregas corretamente e o review nunca aconteceu — ele não tinha
			// o watcher ligado, e a fila é um mecanismo de PULL: sem alguém puxando, o
			// registro fica parado. Dizer "próximo passo" em letra pequena não bastou.
			//
			// Então o comando não sugere: ele INSTRUI, com o comando pronto, e diz o que
			// acontece se for pulado.
			fmt.Println("\n── PRÓXIMO PASSO: REVISE esta unidade ───────────────────────")
			fmt.Printf("   anchors work review --for %s\n\n", c.Unit)
			fmt.Println("   Os gates verdes NÃO provam que está certo — eles confrontam o que é")
			fmt.Println("   DECLARÁVEL. Em três rodadas de um E2E real, 7 defeitos graves (perda")
			fmt.Println("   silenciosa de dado, regra sem teste que a prove, contradição entre duas")
			fmt.Println("   regras da mesma spec) passaram com tudo verde. Nenhum foi achado por")
			fmt.Println("   gate; os 7 vieram de revisão adversarial.")
			if !watcherAtivo(absRoot) {
				fmt.Println("\n   (o watcher não está rodando — com `anchors watch start` esta")
				fmt.Println("    entrega entra na fila sozinha, e `anchors next` a puxa)")
			}
			fmt.Println("─────────────────────────────────────────────────────────────")
			// Confronta o declarado contra o disco ANTES de o registro virar material do
			// revisor. Ver deliver_confront.go: as duas checagens nasceram de divergências
			// reais medidas na primeira rodada em que este fluxo funcionou.
			confrontarEntrega(absRoot, files, c.Unit)

			if len(decisions) == 0 && len(uncovered) == 0 {
				fmt.Println("  nota: você declarou ZERO decisões livres e ZERO lacunas de prova.\n" +
					"  Isso AFIRMA que a régua decidiu tudo e que cada regra tem teste que a\n" +
					"  exercita — o revisor vai confrontar exatamente essa afirmação.")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&stage, "stage", "", "OBRIGATÓRIO — etapa entregue: spec|code|feature|test|plan")
	cmd.Flags().StringVar(&unit, "unit", "", "OBRIGATÓRIO — o arquivo que identifica a unidade de propósito")
	cmd.Flags().StringSliceVar(&files, "file", nil, "arquivo tocado (repetível)")
	cmd.Flags().StringVar(&intent, "intent", "", "OBRIGATÓRIO — o que você diz ter feito")
	cmd.Flags().StringSliceVar(&decisions, "decision", nil, "escolha que a régua não decidiu (repetível)")
	cmd.Flags().StringSliceVar(&uncovered, "uncovered", nil, "o que você sabe que não está provado (repetível)")
	cmd.Flags().StringVar(&date, "date", "", "AAAA-MM-DD (OBRIGATÓRIO — o Anchors não lê o relógio)")
	cmd.Flags().StringVar(&agent, "agent", "", "quem entregou (opcional)")
	return cmd
}

// watcherAtivo diz se o daemon está rodando neste projeto. Serve só para não sugerir
// ligar o que já está ligado — ruído em instrução é o que faz a instrução ser ignorada.
func watcherAtivo(root string) bool {
	_, err := os.Stat(filepath.Join(root, ".anchors", "watch.meta"))
	return err == nil
}

// pecaExistente devolve a peça da unidade que está no disco: o próprio alvo, se existir,
// ou a primeira peça derivada dele que exista. É o que permite `deliver --stage spec`
// funcionar na primeira entrega, quando só o `.spec.md` nasceu.
func pecaExistente(root, rel string) (string, bool) {
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
