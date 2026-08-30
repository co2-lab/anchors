package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// --- a mensagem de commit é matéria-prima do changelog ---
//
// O changelog nasce dos COMMITS, e por isso o padrão tem de valer ANTES de ele existir:
// commit já feito não se conserta, e um histórico onde metade das mensagens não casa o
// formato produz um changelog com buracos que ninguém consegue preencher depois.
//
// Medido no projeto de referência: o squash do PR do plano de mutação entrou como
// "[MTUAO] Plano 0017 — mutação", porque o GitHub usa o TÍTULO DO PR como mensagem do
// squash e o título estava no formato do card. Esse commit — o que introduziu o plano —
// simplesmente não apareceria no changelog.
//
// A régua é o Conventional Commits, e não uma invenção nossa: é o que as ferramentas de
// changelog já sabem ler. Inventar formato aqui daria trabalho duas vezes.

// tiposConvencionais são os tipos aceitos, e o que cada um significa PARA O CHANGELOG.
//
// `feat` e `fix` são os que aparecem para o usuário; o resto é histórico interno, que a
// maioria das ferramentas agrupa ou omite. A lista é fechada de propósito: tipo livre
// vira sinônimo ("bugfix", "hotfix", "correção") e o agrupamento se desfaz.
var tiposConvencionais = []string{
	"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert",
}

func assuntoRE() *regexp.Regexp {
	return regexp.MustCompile(`^(` + strings.Join(tiposConvencionais, "|") +
		`)(\([^)]+\))?(!)?: .+`)
}

// mensagensQueOGitGera não são escritas por ninguém — barrá-las quebraria operações
// normais do git em vez de melhorar o histórico.
func geradaPeloGit(assunto string) bool {
	for _, p := range []string{"Merge ", "Revert ", "fixup! ", "squash! ", "Reapply "} {
		if strings.HasPrefix(assunto, p) {
			return true
		}
	}
	return false
}

func newCommitMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-msg <arquivo>",
		Short: "Confronta a mensagem de commit com o formato que o changelog vai ler",
		Long: `Lê o arquivo de mensagem (o que o git passa ao hook commit-msg) e confronta
o assunto com o Conventional Commits.

O changelog nasce dos commits. Um assunto fora do formato não some do histórico —
some do CHANGELOG, e só se descobre quando alguém vai gerar a primeira versão e o
que faltou já está a centenas de commits de distância.

Mensagens geradas pelo git (merge, revert, fixup) passam: ninguém as escreveu.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("ler a mensagem: %w", err)
			}
			assunto := primeiraLinhaUtil(string(b))
			if assunto == "" {
				return nil // mensagem vazia: quem barra é o git, e a mensagem dele é melhor
			}
			if geradaPeloGit(assunto) {
				return nil
			}
			if assuntoRE().MatchString(assunto) {
				return nil
			}
			cmd.SilenceUsage = true
			return fmt.Errorf("o assunto não casa o formato que o changelog lê:\n"+
				"    %s\n\n"+
				"  Use `tipo(escopo): o que mudou` — o escopo é opcional, o `!` marca quebra:\n"+
				"    feat(board): a contagem do refresh aparece no botão\n"+
				"    fix: o gate acusava quem não mudou\n"+
				"    feat!: `--changed` passa a exigir caminho relativo\n\n"+
				"  Tipos: %s\n\n"+
				"  Isto não é estética: o changelog nasce dos commits, e o que está fora do\n"+
				"  formato não aparece nele. Commit já feito não se conserta.",
				assunto, strings.Join(tiposConvencionais, ", "))
		},
	}
	return cmd
}

// primeiraLinhaUtil devolve o assunto, pulando os comentários que o git põe no arquivo.
func primeiraLinhaUtil(s string) string {
	for _, l := range strings.Split(s, "\n") {
		t := strings.TrimSpace(l)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		return t
	}
	return ""
}
