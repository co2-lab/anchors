package ops

import (
	"fmt"
	"github.com/co2-lab/anchors/internal/i18n"
	"os"
	"os/exec"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// --- dois PRs que discordam no mesmo arquivo viram UM trabalho ---
//
// Quando dois agentes escrevem sobre a mesma regra, o merge conflita no CONTEÚDO — e isso
// não se resolve por automação: escolher um lado é escolher sem ler os dois.
//
// O QUE ACONTECIA: o PR ficava parado esperando o autor. E o autor pode ter terminado a
// sessão — medido no projeto de referência, sete PRs parados assim, alguns há dias. A
// espera não tinha dono, e o trabalho de ninguém envelhecia junto.
//
// O QUE ESTE COMANDO FAZ: transforma a espera em trabalho com dono. Abre um card que pede
// a SÍNTESE — não "escolha um lado", mas "leia os dois e entregue o melhor de cada um".
//
// MEDIDO num caso real (PRs #693 e #556): os dois lados estavam CERTOS e eram sobre coisas
// diferentes — um documentava que o corpo da regra contradizia a decisão, o outro que o
// marcador de dispensa fazia o gate responder INDETERMINADO. Não havia lado a escolher: as
// duas coisas entram, e uma delas renumerada.
//
// A RASTREABILIDADE é o ponto, e é o que distingue isto de fechar os PRs e recomeçar:
//
//	o card novo          cita as duas issues originais e os dois PRs
//	as issues originais  recebem comentário apontando para o card novo
//	os dois PRs          são fechados com comentário dizendo onde o trabalho continuou
//	o PR da síntese      (quem o abrir) cita os dois PRs fechados
//
// Nada se perde: quem procurar por qualquer uma das cinco pontas chega às outras quatro.

func newSynthesizeCmd() *cobra.Command {
	var root, prA, prB, arquivos string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "synthesize",
		Short: "Two PRs conflicting in CONTENT become a synthesis card",
		Long: `Opens the card that asks for the synthesis of two conflicting PRs, and closes both.

A content conflict is not resolved by automation: two agents wrote different things
about the same rule, and picking a side is choosing without reading both.

What this command does is give the wait an OWNER. The new card asks for what neither of the
two PRs delivers alone: the best of each, in a result better than both.

    anchors synthesize --pr-a 693 --pr-b 556 --files Dashboards.spec.md

TRACEABILITY is the point: the card cites both issues and both PRs; the issues
get a comment pointing to it; the PRs are closed saying where the work
continued. Nothing is lost.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(absRoot + "/" + config.DefaultFile)
			if err != nil {
				return err
			}
			if !cfg.GitHubMode() {
				cmd.SilenceUsage = true
				return fmt.Errorf("%s", i18n.T("synthesize.github_only"))
			}
			if strings.TrimSpace(prA) == "" {
				cmd.SilenceUsage = true
				return fmt.Errorf("%s", i18n.T("synthesize.need_pr_a"))
			}
			// O SEGUNDO LADO PODE FALTAR, e o card nasce assim mesmo.
			//
			// Quando o conflito é contra o `develop`, o "outro lado" é um trabalho já
			// mesclado — achá-lo exige ler o histórico do arquivo, e adivinhar erraria.
			// Um card com um lado só ainda é melhor que um PR parado sem dono: quem o
			// pegar descobre o outro lado lendo o arquivo, que é o que ele teria de fazer
			// de qualquer forma.
			repo := cfg.Workflow.Repo
			a := strings.TrimPrefix(strings.TrimSpace(prA), "#")
			b := strings.TrimPrefix(strings.TrimSpace(prB), "#")

			// OS CARDS DOS DOIS PRs, para o card novo herdar a procedência de ambos.
			//
			// É o mesmo `Refs`/`Closes` que o `pr-body` escreve e que o `pr-checks` lê —
			// aqui ele responde "que trabalho cada um destes PRs entregava?".
			cardA, cardB := common.CardDoPR(repo, a), common.CardDoPR(repo, b)
			tituloA := prTitle(repo, a)
			tituloB := prTitle(repo, b)

			corpo := corpoDaSintese(a, b, cardA, cardB, tituloA, tituloB, arquivos)
			titulo := i18n.T("synthesize.title_one", a)
			if b != "" {
				titulo = i18n.T("synthesize.title_two", a, b)
			}

			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), i18n.T("synthesize.dry_run_header"))
				fmt.Fprintln(cmd.OutOrStdout(), titulo)
				fmt.Fprintln(cmd.OutOrStdout())
				fmt.Fprintln(cmd.OutOrStdout(), corpo)
				return nil
			}

			tmp, err := os.CreateTemp("", "anchors-sintese-*.md")
			if err != nil {
				return err
			}
			defer os.Remove(tmp.Name())
			if _, err := tmp.WriteString(corpo); err != nil {
				return err
			}
			tmp.Close()

			// AS LABELS: `to-do` para o claim servir, e `under-<card>` para cada card de
			// origem — a síntese pertence aos dois trabalhos, não a um.
			labels := []string{cfg.Workflow.Labels[0], "anchors:to-do"}
			for _, c := range []string{cardA, cardB} {
				if c == "" {
					continue
				}
				rotulo := initx.LabelSob(c)
				_ = exec.Command("gh", "label", "create", rotulo,
					"--repo", repo, "--color", "c5def5",
					"--description", "finding that was born during the work of card #"+c).Run()
				labels = append(labels, rotulo)
			}

			argv := []string{"issue", "create", "--repo", repo,
				"--title", titulo, "--body-file", tmp.Name()}
			for _, l := range labels {
				argv = append(argv, "--label", l)
			}
			out, err := exec.Command("gh", argv...).CombinedOutput()
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("%s", i18n.T("synthesize.open_failed", err, strings.TrimSpace(string(out))))
			}
			url := strings.TrimSpace(string(out))
			novo := common.NumeroDaIssue(url)
			fmt.Fprintln(cmd.OutOrStdout(), i18n.T("synthesize.card", url))

			// AS QUATRO PONTAS DA RASTREABILIDADE, e cada falha é AVISADA em vez de
			// silenciosa: o card já existe, e perder um link é pior descoberto depois.
			var fechados, naoFechados []string
			pares := [][2]string{{a, b}}
			if b != "" {
				pares = append(pares, [2]string{b, a})
			}
			for _, par := range pares {
				pr, outro := par[0], par[1]
				if outro == "" {
					outro = i18n.T("synthesize.other_side_branch")
				} else {
					outro = "#" + outro
				}
				msg := i18n.T("synthesize.comment_pr", outro, url)
				if o, err := exec.Command("gh", "pr", "comment", pr, "--repo", repo,
					"--body", msg).CombinedOutput(); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), i18n.T("synthesize.warn_no_comment", pr, strings.TrimSpace(string(o))))
				}
				if o, err := exec.Command("gh", "pr", "close", pr, "--repo", repo).CombinedOutput(); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), i18n.T("synthesize.warn_not_closed", pr, strings.TrimSpace(string(o))))
					naoFechados = append(naoFechados, "#"+pr)
				} else {
					fechados = append(fechados, "#"+pr)
				}
			}
			for _, c := range []string{cardA, cardB} {
				if c == "" {
					continue
				}
				msg := i18n.T("synthesize.comment_card", url)
				if o, err := exec.Command("gh", "issue", "comment", c, "--repo", repo,
					"--body", msg).CombinedOutput(); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), i18n.T("synthesize.warn_card_no_link", c, strings.TrimSpace(string(o))))
				}
			}
			// The closing line names what gh ACTUALLY closed. It counted the PRs it meant to
			// close: "PRs #693 and # were closed" with `--pr-a` alone, and "closed" after a
			// `gh pr close` that failed and had just been warned about.
			if novo != "" {
				switch len(fechados) {
				case 0:
				case 1:
					fmt.Fprintln(cmd.OutOrStdout(), i18n.T("synthesize.closed_one", fechados[0]))
				default:
					fmt.Fprintln(cmd.OutOrStdout(), i18n.T("synthesize.closed_many", strings.Join(fechados, ", ")))
				}
				if len(naoFechados) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), i18n.T("synthesize.still_open", strings.Join(naoFechados, ", ")))
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&prA, "pr-a", "", "the first PR in conflict")
	cmd.Flags().StringVar(&prB, "pr-b", "", "the second PR in conflict")
	cmd.Flags().StringVar(&arquivos, "files", "", "the files where the conflict is, comma-separated")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show the card without opening or closing anything")
	return cmd
}

// prTitle lê o título de um PR, para o card da síntese dizer o que cada lado entregava.
//
// Vazio em qualquer erro: o título é contexto, e um card sem ele ainda é acionável — quem
// o ler tem os números dos PRs.
func prTitle(repo, pr string) string {
	out, err := exec.Command("gh", "pr", "view", pr, "--repo", repo,
		"--json", "title", "--jq", ".title // \"\"").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// corpoDaSintese escreve o card que pede o melhor dos dois PRs.
//
// O CARD NÃO DIZ QUAL LADO VENCE, e é isso que o distingue de "resolva o conflito": quem o
// pegar tem de LER os dois e entregar o que nenhum entrega sozinho. Num caso real medido
// (os PRs #693 e #556), os dois lados estavam certos e eram sobre coisas diferentes — não
// havia lado a escolher.
func corpoDaSintese(a, b, cardA, cardB, tituloA, tituloB, arquivos string) string {
	var s strings.Builder
	if b != "" {
		s.WriteString(i18n.T("synthesize.body.head_two"))
	} else {
		s.WriteString(i18n.T("synthesize.body.head_one"))
	}
	s.WriteString(i18n.T("synthesize.body.ask"))

	s.WriteString(i18n.T("synthesize.body.what_each"))
	s.WriteString(fmt.Sprintf("| #%s | %s | %s |\n", a, cardRef(cardA), ouTraco(tituloA)))
	if b != "" {
		s.WriteString(fmt.Sprintf("| #%s | %s | %s |\n\n", b, cardRef(cardB), ouTraco(tituloB)))
	} else {
		s.WriteString(i18n.T("synthesize.body.other_unknown"))
	}

	if strings.TrimSpace(arquivos) != "" {
		s.WriteString(i18n.T("synthesize.body.where"))
		for _, f := range strings.Split(arquivos, ",") {
			if f = strings.TrimSpace(f); f != "" {
				s.WriteString("- `" + f + "`\n")
			}
		}
		s.WriteString("\n")
	}

	s.WriteString(i18n.T("synthesize.body.how"))
	if b != "" {
		s.WriteString(i18n.T("synthesize.body.read_two"))
		s.WriteString(i18n.T("synthesize.body.open_two", a, b))
	} else {
		s.WriteString(i18n.T("synthesize.body.read_one"))
		s.WriteString(i18n.T("synthesize.body.open_one", a))
	}
	s.WriteString(i18n.T("synthesize.body.steps_rest"))
	s.WriteString(i18n.T("synthesize.body.dont"))
	return s.String()
}

func cardRef(c string) string {
	if c == "" {
		return "—"
	}
	return "#" + c
}

func ouTraco(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}
