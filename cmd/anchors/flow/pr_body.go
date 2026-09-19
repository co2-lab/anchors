package flow

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// --- o vínculo card↔PR é do Anchors; a palavra-chave é da plataforma ---
//
// O GitHub só entende as palavras de vínculo em INGLÊS. Não há configuração e não há
// tradução: um projeto que escreva "Refere #44" não cria vínculo nenhum, em silêncio.
//
// A primeira tentativa foi um gate que EXIGIA a palavra em inglês, e isso estava errado
// pelo mesmo motivo que tiramos match de prosa dos gates: o Anchors é multi-idioma, e uma
// régua que obriga o texto do PR a estar noutra língua não é régua do Anchors — é uma
// exigência da plataforma vazando para dentro da doutrina.
//
// A inversão: o vínculo é declarado no VOCABULÁRIO do Anchors — o card que o agente pegou
// (`anchors-owner:`) e os achados que nasceram sob ele (`anchors:under-<n>`) —, e a linha
// que a plataforma entende é GERADA a partir dele. Quem escreve o PR não precisa saber a
// palavra; quem muda de plataforma muda o gerador, não a doutrina.

// linkSyntax é como cada plataforma quer receber "este PR entrega aquele card".
//
// VINCULAR, E NÃO FECHAR — e a diferença é o fluxo inteiro.
//
// Isto gerava `Closes #N`, e o `Closes` faz DUAS coisas de uma vez: o pipeline o lê para
// saber que o PR CONCLUI a implementação (e então move o card para `ready-to-test`), e o
// GitHub o lê para FECHAR a issue no merge. As duas afirmações foram tratadas como uma, e
// não são: "a implementação acabou" não é "o card acabou".
//
// A esteira não termina em `ready-to-test`. Vêm `in-test`, `ready-to-release` e
// `production` — colunas do CD, que o Anchors não rastreia hoje (ele vai até o CI). Fechar
// o card no merge declara entregue um trabalho com três estados à frente, e quem ia testar
// perde de vista o que precisa testar.
//
// MEDIDO no projeto de referência: 71 cards fechados em `ready-to-test` contra 7 abertos,
// 35 num só dia. A coluna que deveria ACUMULAR o que espera teste mostrava só o resíduo —
// os que nunca receberam `Closes`. A esteira entregava, e nada disso era visível.
//
// `Refs #N` cria a referência na timeline do card, que é o que o pipeline precisa para
// achar o vínculo, e NÃO fecha. Quem fecha o card é quem termina a esteira.
//
// Um mapa, e não um `if`: acrescentar uma plataforma é acrescentar uma linha, e o gerador
// não precisa saber quantas existem.
var linkSyntax = map[string]string{
	"github": "Refs #%s",
	// GitLab aceita as mesmas palavras, mas com `#` só no mesmo projeto — a diferença
	// aparece quando o card vive noutro repositório.
	"gitlab": "Refs #%s",
}

func newPRBodyCmd() *cobra.Command {
	var root, cards string
	var sob bool
	cmd := &cobra.Command{
		Use:   "pr-body",
		Short: "Write the lines linking this work's cards, in the platform's syntax",
		Long: `Prints the link lines for the PR body.

LINKS without closing: the card advances to 'ready-to-test' and stays OPEN, because
the pipeline continues in the CD ('in-test', 'ready-to-release', 'production') and
Anchors does not track those columns. Whoever closes the card is whoever finishes
the pipeline.

The link is declared in the Anchors vocabulary: the card you took
('anchors-owner:') and the findings born under it ('anchors:under-<n>'). The
platform keyword is GENERATED from that.

You do not need to know that GitHub only accepts those words in English — and in a
project written in another language, that is exactly what gets wrong in silence: the
PR merges and the card links to nothing.

    anchors pr-body --cards 44          # the card and everything born under it
    anchors pr-body                     # discovers it through ANCHORS_AGENT`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return err
			}
			if !cfg.GitHubMode() {
				cmd.SilenceUsage = true
				return fmt.Errorf("`pr-body` exists in github mode: in local mode there is no " +
					"card to link, and the work is recorded by moving the folder in `issues/`")
			}
			sintaxe, ok := linkSyntax[cfg.Workflow.Mode]
			if !ok {
				cmd.SilenceUsage = true
				return fmt.Errorf("I do not know the link syntax of `%s` — the known ones "+
					"are: %s", cfg.Workflow.Mode, strings.Join(knownPlatforms(), ", "))
			}

			raizes := requestedCards(cards, cfg)
			if len(raizes) == 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("no card: provide `--cards 44` or set `ANCHORS_AGENT` " +
					"so Anchors finds what you took")
			}

			// Cada card ARRASTA o que nasceu sob ele: o achado foi descoberto fazendo
			// aquele trabalho, e os dois se entregam juntos. Esquecer um deixaria o board
			// afirmando que há trabalho pendente que já foi feito.
			raizesPedidas := map[string]bool{}
			todos := map[string]bool{}
			for _, c := range raizes {
				raizesPedidas[c] = true
				todos[c] = true
				for _, sob := range cardsUnder(cfg, c) {
					todos[sob] = true
				}
			}
			ordenados := make([]string, 0, len(todos))
			for c := range todos {
				ordenados = append(ordenados, c)
			}
			sort.Slice(ordenados, func(i, j int) bool {
				a, _ := strconv.Atoi(ordenados[i])
				b, _ := strconv.Atoi(ordenados[j])
				return a < b
			})
			for _, c := range ordenados {
				// `--so-sob` serve a quem JÁ declarou os cards raiz e quer saber o que
				// mais precisa entrar: o CI, que lê os cards do corpo do PR e precisa
				// conferir se os achados sob eles também estão lá.
				if sob && raizesPedidas[c] {
					continue
				}
				fmt.Printf(sintaxe+"\n", c)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&cards, "cards", "", "card number(s), comma-separated")
	cmd.Flags().BoolVar(&sob, "so-sob", false,
		"only the findings born under the given cards, not the cards themselves")
	return cmd
}

func knownPlatforms() []string {
	out := make([]string, 0, len(linkSyntax))
	for k := range linkSyntax {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// requestedCards resolve o que foi passado em `--cards`, ou descobre pelo agente.
func requestedCards(cards string, cfg *config.Config) []string {
	if s := strings.TrimSpace(cards); s != "" {
		var out []string
		for _, c := range strings.Split(s, ",") {
			if c = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(c), "#")); c != "" {
				out = append(out, c)
			}
		}
		return out
	}
	var out []string
	for _, c := range common.AgentCards(cfg) {
		out = append(out, c.Numero)
	}
	return out
}

// cardsUnder lista os achados que nasceram durante o trabalho de um card.
func cardsUnder(cfg *config.Config, card string) []string {
	out, err := exec.Command("gh", "issue", "list",
		"--repo", cfg.Workflow.Repo,
		"--state", "open",
		"--label", initx.LabelSob(card),
		"--limit", "100",
		"--json", "number",
	).Output()
	if err != nil {
		return nil
	}
	var itens []struct {
		Number int `json:"number"`
	}
	if json.Unmarshal(out, &itens) != nil {
		return nil
	}
	nums := make([]string, 0, len(itens))
	for _, i := range itens {
		nums = append(nums, strconv.Itoa(i.Number))
	}
	return nums
}
