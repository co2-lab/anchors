package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		Short: "Dois PRs que conflitam no CONTEÚDO viram um card de síntese",
		Long: `Abre o card que pede a síntese de dois PRs em conflito, e fecha os dois.

Conflito de conteúdo não se resolve por automação: dois agentes escreveram coisas
diferentes sobre a mesma regra, e escolher um lado é escolher sem ler os dois.

O que este comando faz é dar DONO à espera. O card novo pede o que nenhum dos dois
PRs entrega sozinho: o melhor de cada um, num resultado melhor que os dois.

    anchors synthesize --pr-a 693 --pr-b 556 --files Dashboards.spec.md

A RASTREABILIDADE é o ponto: o card cita as duas issues e os dois PRs; as issues
recebem comentário apontando para ele; os PRs são fechados dizendo onde o trabalho
continuou. Nada se perde.`,
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
				return fmt.Errorf("`synthesize` existe no modo github: no modo local não há " +
					"PR a sintetizar")
			}
			if strings.TrimSpace(prA) == "" {
				cmd.SilenceUsage = true
				return fmt.Errorf("informe ao menos o PR que ficou de fora: `--pr-a <n>`")
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
			cardA, cardB := cardDoPR(repo, a), cardDoPR(repo, b)
			tituloA := prTitle(repo, a)
			tituloB := prTitle(repo, b)

			corpo := corpoDaSintese(a, b, cardA, cardB, tituloA, tituloB, arquivos)
			titulo := "[síntese] o que o PR #" + a + " entrega, reconciliado com o que já entrou"
			if b != "" {
				titulo = "[síntese] o que os PRs #" + a + " e #" + b + " entregam, num só"
			}

			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "— card que seria aberto —")
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
					"--description", "achado que nasceu durante o trabalho do card #"+c).Run()
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
				return fmt.Errorf("abrir o card da síntese: %v — %s", err, strings.TrimSpace(string(out)))
			}
			url := strings.TrimSpace(string(out))
			novo := numeroDaIssue(url)
			fmt.Fprintf(cmd.OutOrStdout(), "card da síntese: %s\n", url)

			// AS QUATRO PONTAS DA RASTREABILIDADE, e cada falha é AVISADA em vez de
			// silenciosa: o card já existe, e perder um link é pior descoberto depois.
			pares := [][2]string{{a, b}}
			if b != "" {
				pares = append(pares, [2]string{b, a})
			}
			for _, par := range pares {
				pr, outro := par[0], par[1]
				if outro == "" {
					outro = "o que já está no branch de integração"
				} else {
					outro = "#" + outro
				}
				msg := fmt.Sprintf("🔀 **Fechado em favor da síntese.**\n\n"+
					"Este PR conflita no CONTEÚDO com %s — os dois escreveram sobre o mesmo "+
					"lugar, e escolher um lado seria escolher sem ler o outro.\n\n"+
					"O trabalho continua em %s, que pede o melhor dos dois num resultado só. "+
					"O conteúdo daqui NÃO se perde: é o que aquele card tem de entregar.",
					outro, url)
				if o, err := exec.Command("gh", "pr", "comment", pr, "--repo", repo,
					"--body", msg).CombinedOutput(); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "· aviso: o #%s não recebeu o comentário: %s\n",
						pr, strings.TrimSpace(string(o)))
				}
				if o, err := exec.Command("gh", "pr", "close", pr, "--repo", repo).CombinedOutput(); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "· aviso: o #%s não foi fechado: %s\n",
						pr, strings.TrimSpace(string(o)))
				}
			}
			for _, c := range []string{cardA, cardB} {
				if c == "" {
					continue
				}
				msg := fmt.Sprintf("🔀 O PR deste card conflitava no conteúdo com outro, e o "+
					"trabalho continua na síntese: %s\n\nEste card segue aberto — ele só fecha "+
					"quando a síntese entregar o que ele pedia.", url)
				if o, err := exec.Command("gh", "issue", "comment", c, "--repo", repo,
					"--body", msg).CombinedOutput(); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "· aviso: o card #%s não recebeu o link: %s\n",
						c, strings.TrimSpace(string(o)))
				}
			}
			if novo != "" {
				fmt.Fprintf(cmd.OutOrStdout(),
					"  os PRs #%s e #%s foram fechados, e as quatro pontas se apontam\n", a, b)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&prA, "pr-a", "", "o primeiro PR em conflito")
	cmd.Flags().StringVar(&prB, "pr-b", "", "o segundo PR em conflito")
	cmd.Flags().StringVar(&arquivos, "files", "", "os arquivos onde o conflito está, separados por vírgula")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "mostra o card sem abrir nem fechar nada")
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
		s.WriteString("🔀 **Dois PRs escreveram sobre o mesmo lugar, e os dois foram fechados.**\n\n")
	} else {
		s.WriteString("🔀 **Este PR escreveu sobre um lugar que já mudou, e foi fechado.**\n\n")
	}
	s.WriteString("Este card pede a SÍNTESE: leia os dois lados, e entregue o melhor de cada um ")
	s.WriteString("num resultado melhor que os dois. Não é \"escolha um lado\" — se fosse, a ")
	s.WriteString("automação já teria escolhido.\n\n")

	s.WriteString("## O que cada um entregava\n\n")
	s.WriteString(fmt.Sprintf("| PR | card | o que era |\n|---|---|---|\n"))
	s.WriteString(fmt.Sprintf("| #%s | %s | %s |\n", a, cardRef(cardA), ouTraco(tituloA)))
	if b != "" {
		s.WriteString(fmt.Sprintf("| #%s | %s | %s |\n\n", b, cardRef(cardB), ouTraco(tituloB)))
	} else {
		s.WriteString("\n> **O outro lado não foi identificado.** O conflito é contra o que já ")
		s.WriteString("está no branch de integração — descubra o que mudou ali lendo o histórico ")
		s.WriteString("do arquivo (`git log -p <arquivo>`). Quem fez aquela mudança não precisa ")
		s.WriteString("ser interrompido: o trabalho dele já entrou.\n\n")
	}

	if strings.TrimSpace(arquivos) != "" {
		s.WriteString("## Onde eles discordam\n\n")
		for _, f := range strings.Split(arquivos, ",") {
			if f = strings.TrimSpace(f); f != "" {
				s.WriteString("- `" + f + "`\n")
			}
		}
		s.WriteString("\n")
	}

	s.WriteString("## Como entregar\n\n")
	if b != "" {
		s.WriteString("1. Leia os dois PRs fechados — o conteúdo deles está lá, e é o que você tem de preservar.\n")
	} else {
		s.WriteString("1. Leia o PR fechado — o conteúdo dele está lá, e é o que você tem de preservar.\n")
	}
	if b != "" {
		s.WriteString("2. Abra UM PR com a síntese, citando os dois: `Refs #" + a + "` e `Refs #" + b + "`.\n")
	} else {
		s.WriteString("2. Abra UM PR com a síntese, citando o que foi fechado: `Refs #" + a + "`.\n")
	}
	s.WriteString("3. Se os dois lados forem sobre coisas DIFERENTES, os dois entram — o conflito ")
	s.WriteString("era de posição no arquivo, não de conteúdo.\n")
	s.WriteString("4. Se forem sobre a MESMA coisa, decida com a razão escrita: qual argumento ")
	s.WriteString("sobrevive, e por quê. Registre como revisão na spec.\n\n")

	s.WriteString("> **O que NÃO fazer:** pegar um lado e descartar o outro sem dizer por quê. ")
	s.WriteString("O trabalho dos dois está fechado aqui; se algo se perder, ninguém vai saber ")
	s.WriteString("o que era.\n")
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
