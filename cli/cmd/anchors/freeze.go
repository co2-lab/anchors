package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/spf13/cobra"
)

// --- o BOTÃO DE PÂNICO: congelar e descongelar o projeto ---
//
// O caso é este: um problema aparece — o plano aponta para uma spec que não existe, uma
// decisão de arquitetura se revelou errada, uma credencial vazou — e ninguém deve
// trabalhar até ele ser resolvido. Sem um freio, o time continua produzindo contra uma
// base que vai mudar, e o trabalho de horas envelhece antes de ser entregue.
//
// TRÊS CAMADAS, e nenhuma sozinha basta:
//
//  1. `enabled: false` no anchors.yaml — alcança a MÁQUINA de quem já clonou. É o que os
//     hooks de pre-commit e pre-push consultam. Sem isto, o freio da plataforma não
//     impede ninguém de continuar produzindo estado local.
//
//  2. o RULESET no remoto — barra push e merge para todo mundo. É a única camada que não
//     se contorna de dentro: um `--no-verify` fura os hooks, não fura o servidor.
//
//  3. a ISSUE com o motivo — é onde a pessoa vai ler o que aconteceu e acompanhar o
//     conserto. Um freio sem explicação faz quem esbarra nele tentar contornar.
//
// Elas se ligam numa operação só porque separá-las produz estado inconsistente: um
// ruleset ativo sem issue explicando, ou um `enabled: false` que ninguém empurrou.

func newFreezeCmd() *cobra.Command {
	var root, motivo string
	var semRuleset, semPush bool

	cmd := &cobra.Command{
		Use:   "freeze",
		Short: "Congela o projeto: ninguém trabalha até o `anchors thaw`",
		Long: `Para o trabalho no projeto inteiro, em três camadas.

  1. escreve 'enabled: false' e o motivo no anchors.yaml, e EMPURRA
  2. cria um ruleset no remoto que barra push e merge em todos os branches
  3. abre uma issue com o motivo, para quem esbarrar no freio saber o que houve

O --motivo é OBRIGATÓRIO. Um congelamento sem razão escrita é indistinguível de
configuração quebrada, e quem esbarra nele tenta contornar em vez de ler.

O commit e o push são feitos com --no-verify, de propósito: os hooks que este
comando acabou de ativar recusariam o próprio congelamento.

QUEM CONSERTA PASSA. O admin contorna o ruleset, e '--no-verify' contorna os
hooks — o freio existe para impedir trabalho por INÉRCIA, não o conserto que o
destrava.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			if strings.TrimSpace(motivo) == "" {
				return fmt.Errorf("--motivo é obrigatório: escreva o que aconteceu, porque é " +
					"o texto que TODO comando recusado vai mostrar")
			}
			cfgPath := filepath.Join(absRoot, config.DefaultFile)
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("carregar %s: %w", config.DefaultFile, err)
			}
			if cfg.Congelado() {
				fmt.Printf("○ o projeto JÁ está congelado: %s\n", cfg.MotivoDoCongelamento())
				fmt.Println("  para trocar o motivo, rode `anchors thaw` e congele de novo.")
				return nil
			}

			// 1) O ARQUIVO — escrito por texto, não por serialização do Config.
			//
			// Regravar o YAML inteiro reordenaria chaves e comeria comentários, e o
			// anchors.yaml deste projeto é documentação tanto quanto configuração. Duas
			// linhas no topo bastam, e o diff fica legível para quem revisar depois.
			if err := escreveCongelamento(cfgPath, motivo); err != nil {
				return err
			}
			fmt.Printf("✓ %s: enabled: false\n", config.DefaultFile)

			// 2) O COMMIT e o PUSH — `--no-verify` é obrigatório aqui.
			//
			// Os hooks consultam o remoto, e o remoto ainda não tem o congelamento. Mas o
			// pre-commit deste projeto roda os gates, e um congelamento urgente não pode
			// depender de a suíte estar verde: o motivo do congelamento pode ser
			// justamente que ela não está.
			if !semPush {
				if err := commitaEEmpurra(absRoot, motivo); err != nil {
					fmt.Printf("⚠ não consegui empurrar: %v\n", err)
					fmt.Println("  o arquivo local está congelado; empurre à mão para o freio alcançar o time:")
					fmt.Println("    git add anchors.yaml && git commit --no-verify -m 'freeze' && git push --no-verify")
				} else {
					fmt.Println("✓ commitado e empurrado (--no-verify: os hooks recusariam o próprio freeze)")
				}
			}

			// 3) O RULESET e a ISSUE — só no modo github, onde há remoto a trancar.
			if cfg.ModoGitHub() && cfg.Workflow.Repo != "" {
				if !semRuleset {
					if err := criaRuleset(cfg.Workflow.Repo, motivo); err != nil {
						fmt.Printf("⚠ ruleset não criado: %v\n", err)
						fmt.Println("  o freio local vale; o remoto continua aceitando push de quem contornar os hooks.")
					} else {
						fmt.Println("✓ ruleset `anchors-freeze` ativo — push e merge barrados no remoto")
					}
				}
				if url, err := abreIssueDeCongelamento(cfg.Workflow.Repo, motivo); err != nil {
					fmt.Printf("⚠ issue não aberta: %v\n", err)
				} else {
					fmt.Printf("✓ issue do congelamento: %s\n", url)
				}
			}

			fmt.Println()
			fmt.Println("🛑 PROJETO CONGELADO.")
			fmt.Printf("   Motivo: %s\n", motivo)
			fmt.Println()
			fmt.Println("   Quem for CONSERTAR passa: `--no-verify` nos hooks, e o admin no ruleset.")
			fmt.Println("   Para liberar: `anchors thaw`")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&motivo, "motivo", "", "por que o projeto está sendo congelado (obrigatório)")
	cmd.Flags().BoolVar(&semRuleset, "sem-ruleset", false, "não cria o ruleset no remoto (só o freio local)")
	cmd.Flags().BoolVar(&semPush, "sem-push", false, "não commita nem empurra o anchors.yaml")
	return cmd
}

func newThawCmd() *cobra.Command {
	var root string
	var semPush bool

	cmd := &cobra.Command{
		Use:   "thaw",
		Short: "Descongela o projeto: o trabalho volta",
		Long: `Desfaz o 'anchors freeze': remove o congelamento do anchors.yaml, apaga o
ruleset e fecha a issue.

O commit e o push saem com --no-verify pelo mesmo motivo do freeze: enquanto o
remoto ainda diz 'congelado', os hooks recusariam o próprio descongelamento.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfgPath := filepath.Join(absRoot, config.DefaultFile)
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("carregar %s: %w", config.DefaultFile, err)
			}
			if !cfg.Congelado() {
				fmt.Println("○ o projeto não está congelado — nada a fazer.")
				return nil
			}

			if err := removeCongelamento(cfgPath); err != nil {
				return err
			}
			fmt.Printf("✓ %s: congelamento removido\n", config.DefaultFile)

			if !semPush {
				if err := commitaEEmpurra(absRoot, "thaw"); err != nil {
					fmt.Printf("⚠ não consegui empurrar: %v\n", err)
					fmt.Println("  empurre à mão — enquanto o remoto disser congelado, o time segue barrado.")
				} else {
					fmt.Println("✓ commitado e empurrado")
				}
			}

			if cfg.ModoGitHub() && cfg.Workflow.Repo != "" {
				if err := removeRuleset(cfg.Workflow.Repo); err != nil {
					fmt.Printf("⚠ ruleset não removido: %v\n", err)
					fmt.Println("  o remoto continua barrando push e merge — remova à mão.")
				} else {
					fmt.Println("✓ ruleset removido")
				}
				if err := fechaIssueDeCongelamento(cfg.Workflow.Repo); err != nil {
					fmt.Printf("⚠ issue não fechada: %v\n", err)
				} else {
					fmt.Println("✓ issue do congelamento fechada")
				}
			}

			fmt.Println()
			fmt.Println("✓ PROJETO LIBERADO — o trabalho volta.")
			fmt.Println()
			fmt.Println("  Quem tem o cache do hook quente pode levar até 10 minutos para")
			fmt.Println("  perceber — ou rodar `git fetch` para forçar.")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&semPush, "sem-push", false, "não commita nem empurra o anchors.yaml")
	return cmd
}

// escreveCongelamento põe `enabled: false` e o motivo no TOPO do arquivo.
//
// No topo porque é a primeira coisa que quem abrir o arquivo tem de ver, e porque
// `version:` já mora ali — é onde as declarações de projeto inteiro vivem.
func escreveCongelamento(path, motivo string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// O motivo entra entre aspas: ele é texto livre e costuma ter `:`, que sem aspas
	// quebraria o YAML.
	bloco := fmt.Sprintf("enabled: false\nfreeze_reason: %q\n", motivo)
	return os.WriteFile(path, append([]byte(bloco), b...), 0o644)
}

// removeCongelamento apaga as duas linhas, e SÓ elas.
func removeCongelamento(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var out []string
	for _, l := range strings.Split(string(b), "\n") {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "enabled:") || strings.HasPrefix(t, "freeze_reason:") {
			continue
		}
		out = append(out, l)
	}
	return os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}

// commitaEEmpurra grava o anchors.yaml no remoto, contornando os hooks.
//
// `--no-verify` nos dois: o pre-commit roda os gates, e um congelamento urgente não pode
// depender de a suíte estar verde — o motivo do congelamento pode ser justamente que ela
// não está. E o pre-push consultaria o remoto, que ainda diz o contrário do que estamos
// gravando.
func commitaEEmpurra(root, motivo string) error {
	msg := "chore(anchors): freeze — " + primeiraLinhaDoMotivo(motivo)
	if motivo == "thaw" {
		msg = "chore(anchors): thaw — o projeto volta a aceitar trabalho"
	}
	for _, argv := range [][]string{
		{"add", config.DefaultFile},
		{"commit", "--no-verify", "-m", msg},
		{"push", "--no-verify"},
	} {
		c := exec.Command("git", argv...)
		c.Dir = root
		if out, err := c.CombinedOutput(); err != nil {
			return fmt.Errorf("git %s: %s", argv[0], strings.TrimSpace(string(out)))
		}
	}
	return nil
}

// nomeDoRuleset é fixo: é assim que o `thaw` encontra o que o `freeze` criou.
const nomeDoRuleset = "anchors-freeze"

// criaRuleset tranca push e merge em TODOS os branches.
//
// Um ruleset e não branch protection: ele cobre `~ALL` num objeto só, liga e desliga sem
// tocar na configuração de cada branch, e o histórico de quem o criou fica no audit log
// da organização.
//
// `bypass_actors` com o admin (role_id 5) é deliberado: quem vai consertar precisa
// mesclar o PR que conserta, e um freio que impede o próprio conserto se torna o
// problema. O contorno é explícito e fica registrado.
func criaRuleset(repo, motivo string) error {
	corpo := fmt.Sprintf(`{
  "name": %q,
  "target": "branch",
  "enforcement": "active",
  "bypass_actors": [{"actor_id": 5, "actor_type": "RepositoryRole", "bypass_mode": "always"}],
  "conditions": {"ref_name": {"include": ["~ALL"], "exclude": []}},
  "rules": [{"type": "update"}, {"type": "deletion"}, {"type": "non_fast_forward"}]
}`, nomeDoRuleset)

	c := exec.Command("gh", "api", "--method", "POST",
		"repos/"+repo+"/rulesets", "--input", "-")
	c.Stdin = strings.NewReader(corpo)
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("%v — %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func removeRuleset(repo string) error {
	id, err := idDoRuleset(repo)
	if err != nil {
		return err
	}
	if id == "" {
		return nil // não existe: nada a remover, e isso não é erro
	}
	out, err := exec.Command("gh", "api", "--method", "DELETE",
		"repos/"+repo+"/rulesets/"+id).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v — %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func idDoRuleset(repo string) (string, error) {
	out, err := exec.Command("gh", "api", "repos/"+repo+"/rulesets",
		"--jq", fmt.Sprintf(`.[] | select(.name == %q) | .id`, nomeDoRuleset)).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// tituloDaIssueDeCongelamento é fixo: é assim que o `thaw` a encontra.
const tituloDaIssueDeCongelamento = "[congelado] o projeto está parado"

func abreIssueDeCongelamento(repo, motivo string) (string, error) {
	corpo := "🛑 **O projeto está CONGELADO.**\n\n" +
		"**Motivo:** " + motivo + "\n\n" +
		"Enquanto esta issue estiver aberta:\n\n" +
		"- o `anchors claim` não entrega card novo\n" +
		"- os hooks de `pre-commit` e `pre-push` recusam, citando este motivo\n" +
		"- o ruleset `" + nomeDoRuleset + "` barra push e merge no remoto\n\n" +
		"**Seu trabalho não se perde:** o que está no seu branch local continua lá.\n\n" +
		"**Quem for consertar passa:** `--no-verify` contorna os hooks, e o admin contorna " +
		"o ruleset. O freio existe para impedir trabalho por inércia, não o conserto que " +
		"o destrava.\n\n" +
		"Para liberar: `anchors thaw`."

	out, err := exec.Command("gh", "issue", "create",
		"--repo", repo,
		"--title", tituloDaIssueDeCongelamento,
		"--body", corpo,
	).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v — %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func fechaIssueDeCongelamento(repo string) error {
	out, err := exec.Command("gh", "issue", "list",
		"--repo", repo, "--state", "open", "--limit", "50",
		"--search", tituloDaIssueDeCongelamento,
		"--json", "number,title",
		"--jq", fmt.Sprintf(`[.[] | select(.title == %q)] | .[0].number // empty`, tituloDaIssueDeCongelamento),
	).Output()
	if err != nil {
		return err
	}
	n := strings.TrimSpace(string(out))
	if n == "" {
		return nil
	}
	o, err := exec.Command("gh", "issue", "close", n, "--repo", repo,
		"--reason", "completed",
		"--comment", "✓ Projeto descongelado — o trabalho volta.").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v — %s", err, strings.TrimSpace(string(o)))
	}
	return nil
}
