package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/spf13/cobra"
)

// `anchors install-hooks` fecha o furo do enforcement LOCAL: os gates do
// anchors.yaml podem ser "bloqueantes", mas nada os executa a menos que alguém
// rode `anchors check` na mão. Este comando instala um pre-commit que roda
// `anchors check --changed` sobre os arquivos STAGED — incremental e sem gravar.
//
// É genérico: lê a raiz do repo git e o anchors.yaml; não sabe nada do projeto.
func newInstallHooksCmd() *cobra.Command {
	var root string
	var force bool
	cmd := &cobra.Command{
		Use:   "install-hooks",
		Short: "Instala o git pre-commit que roda os gates sobre os arquivos staged",
		Long: `Escreve .git/hooks/pre-commit para rodar 'anchors check --changed <arquivo>
--no-record' em cada arquivo staged. Um gate bloqueante que reprova barra o commit.

O que passa e o que barra, quando o arquivo não está no mapa:
  • NÃO-REGIDO (não casa nenhuma camada do 'layers:' — package.json, lockfile, CI):
    ignorado. O Anchors não tem jurisdição sobre ele.
  • REGIDO mas fora do mapa (arquivo novo, 'map build' não rodou): BARRA. Fora do
    mapa nenhum gate o confronta — a trinca não é cobrada e o commit passaria a
    certificar trabalho que ninguém verificou. Rode 'anchors map build'.

Incremental: valida só o que o commit toca. Não grava no mapa nem abre issues
(--no-record). Idempotente — reinstalar é seguro; use --force para sobrescrever um
pre-commit existente que não foi escrito por este comando.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			return runInstallHooks(absRoot, force)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&force, "force", false, "sobrescreve um pre-commit existente não gerenciado pelo anchors")
	return cmd
}

// hookMarker identifica um hook escrito por este comando — permite reinstalar sem
// --force e distinguir de um hook artesanal do usuário.
const hookMarker = "# managed-by: anchors install-hooks"

func runInstallHooks(root string, force bool) error {
	// 1. exige anchors.yaml (é um projeto anchors?)
	if _, err := os.Stat(filepath.Join(root, config.DefaultFile)); err != nil {
		return fmt.Errorf("%s não encontrado em %s — rode `anchors init` primeiro", config.DefaultFile, root)
	}

	// 2. descobre o diretório de hooks do git (respeita core.hooksPath e worktrees).
	hooksDir, err := gitHooksDir(root)
	if err != nil {
		// O pre-commit VIVE dentro do repositório: sem ele não há onde instalar. Dizer
		// isso aqui evita que o usuário procure o problema no anchors.yaml.
		if msg := gitmeta.Explica(gitmeta.Verifica(root), "instalar o pre-commit"); msg != "" {
			return errors.New(msg)
		}
		return fmt.Errorf("localizar o diretório de hooks do git: %w", err)
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return fmt.Errorf("criar %s: %w", hooksDir, err)
	}

	hookPath := filepath.Join(hooksDir, "pre-commit")

	// 3. se já existe e não é nosso, respeita (a menos que --force).
	if existing, rerr := os.ReadFile(hookPath); rerr == nil {
		if !strings.Contains(string(existing), hookMarker) && !force {
			return fmt.Errorf(
				"já existe um pre-commit em %s que não foi escrito pelo anchors.\n"+
					"  revise-o e, se quiser substituir, rode com --force",
				hookPath)
		}
	}

	// O COMMIT-MSG acompanha o pre-commit, e existe por um motivo só: é o único hook que
	// RECEBE a mensagem. Medido: o git não grava `.git/COMMIT_EDITMSG` antes do
	// pre-commit — nem com `-m` —, e lê-lo ali devolve a mensagem do commit ANTERIOR.
	//
	// Sem este hook, os marcadores `[skip-regra@CODIGO: motivo]` não teriam efeito: a
	// dispensa escrita na mensagem seria lida do commit errado, em silêncio.
	msgHook := filepath.Join(hooksDir, "commit-msg")
	if existing, rerr := os.ReadFile(msgHook); rerr == nil &&
		!strings.Contains(string(existing), hookMarker) && !force {
		fmt.Printf("⚠  %s existe e não foi escrito pelo anchors — não sobrescrito.\n", msgHook)
	} else if err := os.WriteFile(msgHook, []byte(commitMsgScript), 0o755); err != nil {
		return fmt.Errorf("escrever %s: %w", msgHook, err)
	}

	if err := os.WriteFile(hookPath, []byte(preCommitScript), 0o755); err != nil {
		return fmt.Errorf("escrever %s: %w", hookPath, err)
	}

	fmt.Printf("✓ pre-commit instalado em %s\n", hookPath)
	if _, lookErr := exec.LookPath("anchors"); lookErr != nil {
		fmt.Println("⚠  o binário 'anchors' não está no PATH — o hook vai falhar até instalá-lo")
		fmt.Println("   (cd cli && GOFLAGS=-mod=mod go install ./cmd/anchors) e garanta $(go env GOPATH)/bin no PATH")
	}
	fmt.Println("  o hook roda `anchors check --changed` nos arquivos staged; gate bloqueante barra o commit.")
	fmt.Println("  arquivo REGIDO fora do mapa também barra (rode `anchors map build`); não-regido é ignorado.")
	return nil
}

// gitHooksDir devolve o diretório de hooks efetivo do repo: honra core.hooksPath se
// configurado; senão <git-common-dir>/hooks (correto também em worktrees).
func gitHooksDir(root string) (string, error) {
	// core.hooksPath tem precedência (pode ser relativo à raiz do worktree).
	if out, err := gitOutput(root, "config", "--get", "core.hooksPath"); err == nil {
		hp := strings.TrimSpace(out)
		if hp != "" {
			if filepath.IsAbs(hp) {
				return hp, nil
			}
			return filepath.Join(root, hp), nil
		}
	}
	common, err := gitOutput(root, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", err
	}
	dir := strings.TrimSpace(common)
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(root, dir)
	}
	return filepath.Join(dir, "hooks"), nil
}

func gitOutput(root string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.Output()
	return string(out), err
}

// preCommitScript é o hook instalado. Portátil (bash), sem dependências além do
// binário `anchors` no PATH.
//
// O hook é DELIBERADAMENTE burro: uma chamada a `anchors verify --phase pre-commit`.
// Toda a régua — quais gates a fase cobra, quais ferramentas externas entram, o que é
// benigno — mora no anchors.yaml, não aqui. Um hook que decide vira uma segunda
// configuração, versionada fora do repo e diferente em cada máquina.
const preCommitScript = `#!/usr/bin/env bash
` + hookMarker + `
# Roda o que a fase pre-commit cobra: gates do Anchors + ferramentas externas
# (tsc/eslint/spellcheck) declaradas no anchors.yaml.
# Reinstale/atualize com: anchors install-hooks --force
set -euo pipefail

if ! command -v anchors >/dev/null 2>&1; then
  echo "✗ pre-commit: binário 'anchors' não está no PATH."
  echo "  Instale (cd cli && GOFLAGS=-mod=mod go install ./cmd/anchors) e garanta \$(go env GOPATH)/bin no PATH."
  exit 1
fi

ROOT="$(git rev-parse --show-toplevel)"
STAGED=$(git diff --cached --name-only --diff-filter=ACMR)
[ -z "$STAGED" ] && exit 0

# DISPENSA — na MENSAGEM DO COMMIT:
#
#   [skip-trinca-completa@WRKSP: a feature ainda e um card]
#
# O alvo e o CODIGO do artefato, e so ele fica dispensado: os outros continuam sendo
# confrontados. Dispensar a regra inteira apagaria o gate para o repositorio todo, e uma
# quebra por descuido noutro lugar passaria junto.
#
# Quem barra e o hook commit-msg, nao este: a mensagem NAO EXISTE no pre-commit (o git so
# a grava depois). Este hook reporta cedo, para o problema aparecer antes de escrever a
# mensagem.
#
# ANCHORS_SKIP_RULES="id=motivo" ainda funciona, para um CI que nao controla a mensagem.
#
# E POR REGRA, e nao por commit: dispensar trinca-completa deixa passar a spec que
# nasce sozinha, e os outros gates continuam barrando. Um bypass global calaria também o
# gate que achou defeito de verdade.
#
# O motivo é obrigatório — sem justificativa escrita, uma dispensa é indistinguível de
# alguém fugindo de um gate.

FAIL=0
# Uma invocação para TODOS os arquivos staged: config e mapa carregam uma vez, e os
# gates relacionais confrontam a unidade uma vez (antes, o loop por arquivo repetia
# o mesmo trabalho a cada peça da mesma trinca).
#
# O código 3 é "não tenho jurisdição sobre isto" — NENHUM arquivo staged casa uma
# camada do 'layers:'. Não é reprovação, e barrar aí impediria commitar mudança só
# de configuração (package.json, yarn.lock, .json de ferramenta), que é trabalho
# legítimo que a Estrutura deliberadamente não rege. O código existe justamente
# para o hook distinguir os dois casos; tratá-lo como falha desperdiça a distinção.
set +e
(cd "$ROOT" && anchors verify --phase pre-commit --staged --no-record)
STATUS=$?
set -e
if [ "$STATUS" -eq 3 ]; then
  echo "· pre-commit: nenhum arquivo staged é regido pela Estrutura — nada a confrontar."
elif [ "$STATUS" -ne 0 ]; then
  FAIL=1
fi

if [ "$FAIL" -ne 0 ]; then
  # A DECISÃO FINAL é do 'commit-msg', e não daqui.
  #
  # A mensagem NÃO EXISTE no pre-commit: o git só a grava depois, e '.git/COMMIT_EDITMSG'
  # aqui carrega a do commit ANTERIOR (medido, e confirmado na documentação do githooks).
  # Barrar agora impediria toda dispensa declarada na mensagem de ter efeito — o commit
  # morreria antes de alguém poder ler o '[skip-regra@CODIGO: motivo]'.
  #
  # Então este hook REPORTA e deixa seguir; o 'commit-msg' reconfronta com a mensagem em
  # mãos e barra se a dispensa não cobrir o que reprovou. Quem não usa dispensa nenhuma vê
  # o mesmo resultado, um passo depois.
  if [ -x "$ROOT/.git/hooks/commit-msg" ]; then
    echo "──────────────────────────────────────────────────────────────"
    echo "· gates reprovaram. Se for deliberado, declare na mensagem do commit:"
    echo "    [skip-<regra>@<CODIGO>: por quê]"
    echo "  Sem isso, o commit-msg barra."
    exit 0
  fi
  echo "──────────────────────────────────────────────────────────────"
  echo "✗ commit BARRADO pelos gates do anchors. Corrija acima e recommite."
  exit 1
fi

# Extensão: hooks locais do projeto (réguas que o anchors não cobre, ex.: lint de
# arquitetura de import). Cada executável em .git/hooks/pre-commit.d/ roda com os
# arquivos staged como argumentos; qualquer um que falhe barra o commit.
# (Portável a bash 3.2 do macOS — sem mapfile.)
HOOK_D="$ROOT/.git/hooks/pre-commit.d"
if [ -d "$HOOK_D" ]; then
  STAGED_ARR=()
  while IFS= read -r line; do [ -n "$line" ] && STAGED_ARR+=("$line"); done <<< "$STAGED"
  for h in "$HOOK_D"/*; do
    [ -x "$h" ] || continue
    "$h" "${STAGED_ARR[@]}" || exit 1
  done
fi
exit 0
`

// commitMsgScript roda os gates COM a mensagem em mãos.
//
// O pre-commit não pode fazer isso: medido, o git não grava `.git/COMMIT_EDITMSG` antes
// dele — nem com `-m`. Só o `commit-msg` RECEBE o arquivo, como primeiro argumento.
//
// Por que rodar de novo, e não só aqui: o pre-commit é a barreira que pega o caso comum
// (nada dispensado) o mais cedo possível, antes de o autor escrever a mensagem. Este hook
// existe para o caso em que a mensagem MUDA o veredito — e ele só reexecuta quando há
// marcador, para não pagar o custo duas vezes em todo commit.
const commitMsgScript = `#!/usr/bin/env bash
# anchors:hook — instalado por 'anchors install-hooks'
set -euo pipefail

MSG_FILE="$1"
ROOT="$(git rev-parse --show-toplevel)"

# Roda SEMPRE, e não só quando há marcador: é aqui que o veredito vale. O pre-commit
# reporta cedo (útil para ver o problema antes de escrever a mensagem), mas não pode
# barrar — sem a mensagem, ele não tem como saber se a reprovação foi dispensada.
command -v anchors >/dev/null 2>&1 || exit 0

# A MENSAGEM primeiro, e independente de haver arquivo staged: ela é matéria-prima do
# changelog, e o changelog nasce dos commits. Um assunto fora do formato não some do
# histórico — some do CHANGELOG, e isso só se descobre quando alguém gera a primeira
# versão e o que faltou já está a centenas de commits de distância.
if ! anchors commit-msg "$MSG_FILE"; then
  exit 1
fi

STAGED=$(git diff --cached --name-only --diff-filter=ACMR)
[ -z "$STAGED" ] && exit 0

echo "· commit-msg: confrontando com a mensagem em mãos"
set +e
(cd "$ROOT" && anchors verify --phase pre-commit --staged --no-record --commit-msg "$MSG_FILE")
STATUS=$?
set -e
if [ "$STATUS" -eq 3 ]; then
  exit 0
elif [ "$STATUS" -ne 0 ]; then
  echo "──────────────────────────────────────────────────────────────"
  echo "✗ commit BARRADO pelos gates do anchors."
  echo "  Se a reprovação for deliberada, declare na mensagem:"
  echo "    [skip-<regra>@<CODIGO>: por quê]"
  exit 1
fi
`
