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
			absRoot, err := config.AbsRoot(root)
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
		if msg := gitmeta.Explain(gitmeta.Check(root), "instalar o pre-commit"); msg != "" {
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

	// O PRE-PUSH é a rede de segurança do pre-commit: quem commitou antes do
	// congelamento, ou com `--no-verify`, ainda esbarra nele antes de o trabalho sair da
	// máquina. E ele confere SEM cache — o push é raro o bastante para pagar o fetch, e é
	// o último momento em que a informação ainda muda o desfecho.
	pushHook := filepath.Join(hooksDir, "pre-push")
	if existing, rerr := os.ReadFile(pushHook); rerr == nil &&
		!strings.Contains(string(existing), hookMarker) && !force {
		fmt.Printf("⚠  %s existe e não foi escrito pelo anchors — não sobrescrito.\n", pushHook)
	} else if err := os.WriteFile(pushHook, []byte(prePushScript), 0o755); err != nil {
		return fmt.Errorf("escrever %s: %w", pushHook, err)
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
	fmt.Printf("✓ pre-push instalado em %s\n", pushHook)
	fmt.Println("  os dois conferem se o projeto está CONGELADO no remoto antes de deixar o trabalho seguir.")
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

# O PROJETO ESTÁ CONGELADO? — a conferência vem ANTES de qualquer trabalho.
#
# Aqui e não só no pre-push porque o custo de descobrir tarde é a tarde inteira: quem só
# soubesse ao empurrar já teria escrito o código. O commit é o primeiro momento em que o
# Anchors tem a pessoa na frente.
#
# O CACHE é o que torna isto viável. Um 'git fetch' custa ~1.7s, e o pre-commit roda
# dezenas de vezes por dia — pagar isso a cada commit é o tipo de atrito que faz alguém
# desligar o hook, e um hook desligado não protege nada. A rede é consultada no máximo uma
# vez a cada 10 minutos; no resto, vale o que ficou guardado.
#
# O cache vive em .git/, que não é versionado: ele é estado da MÁQUINA, não do projeto.
CACHE="$ROOT/.git/anchors-freeze-cache"
JANELA=600

congelado=""
motivo=""
agora=$(date +%s)
carimbo=0
[ -f "$CACHE" ] && carimbo=$(head -1 "$CACHE" 2>/dev/null || echo 0)

if [ $((agora - carimbo)) -ge "$JANELA" ]; then
  base=""
  for b in develop main master; do
    if git ls-remote --exit-code --heads origin "$b" >/dev/null 2>&1; then base="$b"; break; fi
  done
  if [ -n "$base" ] && git fetch --quiet origin "$base" 2>/dev/null; then
    remota=$(git show "origin/$base:anchors.yaml" 2>/dev/null || true)
    if printf '%s' "$remota" | grep -qE '^enabled:[[:space:]]*false[[:space:]]*$'; then
      congelado="sim"
      motivo=$(printf '%s' "$remota" | sed -n 's/^freeze_reason:[[:space:]]*//p' | head -1)
    fi
    printf '%s\n%s\n%s\n' "$agora" "$congelado" "$motivo" > "$CACHE"
  fi
  # Sem rede: segue. Um hook que barra por não conseguir consultar transformaria trabalho
  # offline em impossível, e o freio de verdade é o ruleset no remoto.
else
  congelado=$(sed -n 2p "$CACHE" 2>/dev/null || true)
  motivo=$(sed -n 3p "$CACHE" 2>/dev/null || true)

  # O CACHE É ASSIMÉTRICO, e isto não é detalhe.
  #
  # Guardar "não congelado" por 10 minutos custa, no pior caso, alguns commits que
  # deveriam ter sido barrados — e o pre-push os pega antes de saírem da máquina.
  #
  # Guardar "CONGELADO" custa o contrário: o projeto é liberado e a pessoa continua
  # barrada por até 10 minutos, sem entender por quê e sem nada que ela possa fazer.
  # Medido: descongelei o remoto e o dev seguiu recusado com o motivo antigo.
  #
  # Então o cache POSITIVO não vale: quando ele diz "congelado", a rede é consultada de
  # novo para confirmar. O custo é pagar o fetch enquanto durar o congelamento — que é
  # exatamente quando ninguém deveria estar commitando de qualquer forma.
  if [ -n "$congelado" ]; then
    base=""
    for b in develop main master; do
      if git ls-remote --exit-code --heads origin "$b" >/dev/null 2>&1; then base="$b"; break; fi
    done
    if [ -n "$base" ] && git fetch --quiet origin "$base" 2>/dev/null; then
      remota=$(git show "origin/$base:anchors.yaml" 2>/dev/null || true)
      if printf '%s' "$remota" | grep -qE '^enabled:[[:space:]]*false[[:space:]]*$'; then
        motivo=$(printf '%s' "$remota" | sed -n 's/^freeze_reason:[[:space:]]*//p' | head -1)
      else
        congelado=""; motivo=""
      fi
      printf '%s\n%s\n%s\n' "$agora" "$congelado" "$motivo" > "$CACHE"
    fi
  fi
fi

if [ -n "$congelado" ]; then
  echo ""
  echo "🛑 COMMIT RECUSADO — o projeto está CONGELADO."
  echo ""
  [ -n "$motivo" ] && echo "   Motivo: $motivo" || \
    echo "   (nenhum motivo declarado no anchors.yaml do origin)"
  echo ""
  echo "   Pare o trabalho que estiver fazendo: ele não vai poder ser entregue até o"
  echo "   descongelamento, e o que se produz agora envelhece contra o conserto que vem."
  echo ""
  echo "   Se você é quem vai CONSERTAR o que causou o congelamento, use --no-verify —"
  echo "   o freio existe para impedir trabalho por inércia, não o próprio conserto."
  echo ""
  exit 1
fi

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

// --- pre-push: o freio ALCANÇA a máquina de quem já clonou ---
//
// O congelamento (`enabled: false` no anchors.yaml) mora no repositório, e o repositório
// remoto é o que está atualizado. Um dev com clone de ontem não tem o campo — e sem este
// hook, ele trabalharia a tarde inteira sem saber que o projeto está parado.
//
// O hook confronta o `anchors.yaml` LOCAL com o do remoto ANTES de deixar o push sair.
// Não é um `pull` automático: puxar por conta própria mudaria a árvore de quem está no
// meio de um trabalho, e um hook que altera o que a pessoa está fazendo é pior que o
// problema que resolve. Ele RECUSA e diz o que fazer.
//
// A verificação é barata: um `git fetch` do arquivo, não do repositório inteiro.
//
// O QUE ELE NÃO GARANTE, e é honesto dizer: um hook local é contornável (`--no-verify`),
// e quem clona depois do congelamento não passa por ele até instalar os hooks. Por isso
// ele é a SEGUNDA camada — a primeira é o ruleset no remoto, que ninguém contorna de
// dentro. Este existe para que a pessoa DESCUBRA cedo, não para ser inviolável.
const prePushScript = `#!/usr/bin/env bash
# managed-by: anchors install-hooks
# Recusa o push quando o projeto está CONGELADO no remoto, ou quando o anchors.yaml
# local está atrasado em relação a ele.
# Reinstale/atualize com: anchors install-hooks --force
set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
CFG="anchors.yaml"
[ -f "$ROOT/$CFG" ] || exit 0

remoto="${1:-origin}"

# O branch de INTEGRAÇÃO é a referência, não o branch atual: o congelamento é declarado
# lá, e um branch de trabalho não o teria mesmo estando em dia.
base=""
for b in develop main master; do
  if git ls-remote --exit-code --heads "$remoto" "$b" >/dev/null 2>&1; then base="$b"; break; fi
done
[ -n "$base" ] || exit 0

# Só o ARQUIVO, não o repositório: barato o bastante para rodar a cada push.
if ! git fetch --quiet "$remoto" "$base" 2>/dev/null; then
  echo "⚠  pre-push: não consegui falar com '$remoto' — seguindo sem conferir o congelamento."
  exit 0
fi

remota=$(git show "$remoto/$base:$CFG" 2>/dev/null || true)
[ -n "$remota" ] || exit 0

# CONGELADO NO REMOTO: barra, e mostra o motivo que está escrito lá.
if printf '%s' "$remota" | grep -qE '^enabled:[[:space:]]*false[[:space:]]*$'; then
  motivo=$(printf '%s' "$remota" | sed -n 's/^freeze_reason:[[:space:]]*//p' | head -1)
  echo ""
  echo "🛑 PUSH RECUSADO — o projeto está CONGELADO."
  echo ""
  [ -n "$motivo" ] && echo "   Motivo: $motivo" || \
    echo "   (nenhum motivo declarado no anchors.yaml do $remoto/$base)"
  echo ""
  echo "   O trabalho não se perde: ele fica no seu branch local até o descongelamento."
  echo "   Acompanhe a issue de congelamento no repositório."
  echo ""
  echo "   Se você é quem vai CONSERTAR o que causou o congelamento, use --no-verify"
  echo "   — o freio existe para impedir trabalho por inércia, não o próprio conserto."
  echo ""
  exit 1
fi

# NÃO congelado, mas o local está ATRASADO: avisa sem barrar.
#
# Barrar aqui seria exigir que todo push viesse de uma árvore sincronizada, e isso quebra
# o trabalho paralelo legítimo. O que importa é a pessoa SABER que o anchors.yaml mudou —
# porque é lá que o congelamento apareceria.
local_cfg=$(cat "$ROOT/$CFG")
if [ "$local_cfg" != "$remota" ]; then
  echo "⚠  pre-push: seu $CFG difere do $remoto/$base."
  echo "   Rode 'git pull --rebase $remoto $base' para ver o que mudou — é lá que um"
  echo "   congelamento seria declarado."
fi
exit 0
`
