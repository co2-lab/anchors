package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// --- o vínculo card↔PR é do Anchors; a palavra-chave é da plataforma ---
//
// O GitHub fecha um card quando o corpo do PR carrega `Closes #N` — e só em INGLÊS. Não
// há configuração e não há tradução: um projeto que escreva "Fecha #44" tem o card aberto
// depois do merge, em silêncio.
//
// A primeira tentativa foi um gate que EXIGIA a palavra em inglês, e isso estava errado
// pelo mesmo motivo que tiramos match de prosa dos gates: o Anchors é multi-idioma, e uma
// régua que obriga o texto do PR a estar noutra língua não é régua do Anchors — é uma
// exigência da plataforma vazando para dentro da doutrina.
//
// A inversão: o vínculo é declarado no VOCABULÁRIO do Anchors — o card que o agente pegou
// (`anchors-owner:`) e os achados que nasceram sob ele (`anchors:sob-<n>`) —, e a linha
// que a plataforma entende é GERADA a partir dele. Quem escreve o PR não precisa saber a
// palavra; quem muda de plataforma muda o gerador, não a doutrina.

// sintaxeDeFechamento é como cada plataforma quer receber "este PR fecha aquele card".
//
// Um mapa, e não um `if`: acrescentar uma plataforma é acrescentar uma linha, e o gerador
// não precisa saber quantas existem.
var sintaxeDeFechamento = map[string]string{
	"github": "Closes #%s",
	// GitLab aceita as mesmas palavras, mas com `#` só no mesmo projeto — a diferença
	// aparece quando o card vive noutro repositório.
	"gitlab": "Closes #%s",
}

func newPRBodyCmd() *cobra.Command {
	var root, cards string
	cmd := &cobra.Command{
		Use:   "pr-body",
		Short: "Escreve as linhas que fecham os cards deste trabalho, na sintaxe da plataforma",
		Long: `Imprime as linhas de fechamento para o corpo do PR.

O vínculo é declarado no vocabulário do Anchors: o card que você pegou
('anchors-owner:') e os achados que nasceram sob ele ('anchors:sob-<n>'). A
palavra-chave da plataforma é GERADA a partir disso.

Você não precisa saber que o GitHub só aceita 'Closes' em inglês — e num projeto
escrito noutro idioma, isso é justamente o que se erra em silêncio: o PR mescla e
o card fica aberto.

    anchors pr-body --cards 44          # o card e tudo que nasceu sob ele
    anchors pr-body                     # descobre pelo ANCHORS_AGENT`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return err
			}
			if !cfg.ModoGitHub() {
				cmd.SilenceUsage = true
				return fmt.Errorf("`pr-body` existe no modo github: no modo local não há " +
					"card a fechar, e o trabalho se registra movendo a pasta em `issues/`")
			}
			sintaxe, ok := sintaxeDeFechamento[cfg.Workflow.Mode]
			if !ok {
				cmd.SilenceUsage = true
				return fmt.Errorf("não sei a sintaxe de fechamento de `%s` — as conhecidas "+
					"são: %s", cfg.Workflow.Mode, strings.Join(plataformasConhecidas(), ", "))
			}

			raizes := cardsPedidos(cards, cfg)
			if len(raizes) == 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("nenhum card: informe `--cards 44` ou defina `ANCHORS_AGENT` " +
					"para que o Anchors ache o que você pegou")
			}

			// Cada card ARRASTA o que nasceu sob ele: o achado foi descoberto fazendo
			// aquele trabalho, e os dois se entregam juntos. Esquecer um deixaria o board
			// afirmando que há trabalho pendente que já foi feito.
			todos := map[string]bool{}
			for _, c := range raizes {
				todos[c] = true
				for _, sob := range cardsSob(cfg, c) {
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
				fmt.Printf(sintaxe+"\n", c)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&cards, "cards", "", "número(s) do(s) card(s), separados por vírgula")
	return cmd
}

func plataformasConhecidas() []string {
	out := make([]string, 0, len(sintaxeDeFechamento))
	for k := range sintaxeDeFechamento {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// cardsPedidos resolve o que foi passado em `--cards`, ou descobre pelo agente.
func cardsPedidos(cards string, cfg *config.Config) []string {
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
	for _, c := range cardsDoAgente(cfg) {
		out = append(out, c.numero)
	}
	return out
}

// cardsSob lista os achados que nasceram durante o trabalho de um card.
func cardsSob(cfg *config.Config, card string) []string {
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
