package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/spf13/cobra"
)

// `anchors verify` é o COMANDO ÚNICO da fase: uma invocação que roda tudo o que
// aquele momento cobra — os gates do Anchors E as ferramentas de terceiro (tsc,
// eslint, spellcheck), declaradas como gates externos no anchors.yaml.
//
// O que ele resolve. O pre-commit chamava `anchors check --changed <f>` num LOOP,
// um processo por arquivo (~1,2s cada: 63 arquivos = ~76s), e as demais ferramentas
// ficavam fora — penduradas em `pre-commit.d/` ou só no `yarn verify`, cada uma com
// sua própria noção de quando rodar. O resultado era o pior dos dois mundos: lento
// no commit e incompleto, com a régua espalhada por três lugares que ninguém
// mantinha em sincronia.
//
// Aqui a fase é o argumento (`--phase pre-commit`) e o anchors.yaml é a única fonte
// sobre o que cada gate mede, quanto custa e em que momento é cobrado.
func newVerifyCmd() *cobra.Command {
	var root, phase, category string
	var changed []string
	var staged, all, skipSlow, noRecord bool

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Roda TUDO o que a fase cobra: gates do Anchors + ferramentas externas",
		Long: `Uma invocação por FASE, em vez de um comando por ferramenta.

Roda os gates declarados no anchors.yaml que valem para a fase pedida — inclusive os
externos (` + "`run:`" + `), que delegam a tsc/eslint/spellcheck. O projeto declara em
QUE MOMENTO cada gate é cobrado (` + "`when:`" + `), QUANTO ele custa (` + "`cost:`" + `)
e O QUE ele mede (` + "`category:`" + `); este comando só obedece.

  anchors verify --phase pre-commit --staged   # o que o hook chama
  anchors verify --phase ci --all              # a foto completa
  anchors verify --category types              # só uma família

Escopo dos gates externos (` + "`scope:`" + `):
  node     (default) uma execução por arquivo — o gate mede um arquivo isolado
  batch    UMA execução recebendo os arquivos em "$@" — eslint, prettier
  project  UMA execução sem alvos — tsc, knip, madge (olham o projeto inteiro)

Um gate de escopo project/batch só roda se HOUVER arquivo relevante no recorte: um
commit só de README não dispara o typecheck do monorepo.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if staged {
				lista, err := arquivosStaged(root)
				if err != nil {
					return err
				}
				if len(lista) == 0 {
					fmt.Println("nada staged — nada a verificar.")
					return nil
				}
				changed = append(changed, lista...)
			}
			if !all && len(changed) == 0 {
				return fmt.Errorf("informe --staged, --changed <arquivo> ou --all")
			}

			// `verify` é uma fachada: delega ao MESMO pipeline do `check`, para não
			// existirem duas verdades sobre o que é passar. O que ele acrescenta é a
			// fase (e a coleta do staged), não uma segunda régua.
			sub := []string{"check", "--root", root}
			if all {
				sub = append(sub, "--all")
			} else {
				for _, c := range changed {
					sub = append(sub, "--changed", c)
				}
			}
			if phase != "" {
				sub = append(sub, "--phase", phase)
			}
			if category != "" {
				sub = append(sub, "--category", category)
			}
			if skipSlow {
				sub = append(sub, "--skip-slow")
			}
			if noRecord {
				sub = append(sub, "--no-record")
			}
			// A fase automática nunca espera IA: gate de julgamento não pode barrar um
			// commit nem enfileirar lixo repetido a cada volta.
			if phase != "" && phase != "manual" {
				sub = append(sub, "--deterministic")
			}
			return rodarSubcomando(sub)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&phase, "phase", "", "a fase a cobrar (pre-commit|pre-push|ci|manual)")
	cmd.Flags().StringVar(&category, "category", "", "cobra só os gates desta natureza (types|style|traceability…)")
	cmd.Flags().StringSliceVar(&changed, "changed", nil, "arquivo(s) a verificar (repetível)")
	cmd.Flags().BoolVar(&staged, "staged", false, "usa os arquivos staged no git (o modo do pre-commit)")
	cmd.Flags().BoolVar(&all, "all", false, "varre todos os nós (a foto completa; caro)")
	cmd.Flags().BoolVar(&skipSlow, "skip-slow", false, "pula os gates declarados `cost: slow`")
	cmd.Flags().BoolVar(&noRecord, "no-record", false, "só reporta: não carimba o mapa nem abre issues")
	return cmd
}

// arquivosStaged lista o que está no índice do git (ACMR — sem deleções, que não há
// como verificar). É a mesma lista que o pre-commit usava, agora obtida pelo próprio
// anchors: o hook deixa de precisar saber a sintaxe do git.
func arquivosStaged(root string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACMR")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		// `--staged` é o modo do pre-commit: sem git não há índice, e o erro cru do git
		// (`exit status 128`) não diz qual das duas faltas é.
		if msg := gitmeta.Explica(gitmeta.Verifica(root), "listar os arquivos staged"); msg != "" {
			return nil, errors.New(msg)
		}
		return nil, fmt.Errorf("listar arquivos staged: %w", err)
	}
	var lista []string
	for _, l := range strings.Split(string(out), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lista = append(lista, l)
		}
	}
	return lista, nil
}

// rodarSubcomando reexecuta o próprio binário. Reusar o pipeline do `check` por
// processo (em vez de refatorar o RunE dele para uma função compartilhada) mantém
// UMA implementação do que é verificar — e o custo é um fork, não N.
func rodarSubcomando(args []string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	c := exec.Command(exe, args...)
	c.Stdout, c.Stderr, c.Stdin = os.Stdout, os.Stderr, os.Stdin
	return traduzSaidaDoFilho(c.Run())
}

// traduzSaidaDoFilho converte o resultado bruto de um subprocesso no erro que o
// `main` sabe interpretar.
//
// Existe separada para ser TESTÁVEL: exercitar isto pelo `rodarSubcomando`
// exigiria reexecutar o binário do anchors contra um projeto de verdade em
// disco, e o que está sob teste é a tradução do código, não o comando.
func traduzSaidaDoFilho(err error) error {
	// O código de saída do FILHO precisa atravessar a fronteira do processo.
	//
	// `c.Run()` devolve um `*exec.ExitError` genérico, e o `main` — que converte
	// `errNaoRegido` em `ExitNaoRegido` — não o reconhece. Resultado: o `check`
	// saía com 3 ("não tenho jurisdição"), o `verify` traduzia para 1, e o
	// pre-commit barrava um commit só de configuração (package.json, yarn.lock),
	// que é exatamente o caso que o código 3 existe para permitir.
	var ee *exec.ExitError
	if errors.As(err, &ee) && ee.ExitCode() == ExitNaoRegido {
		return errNaoRegido{target: "os arquivos staged"}
	}
	return err
}
