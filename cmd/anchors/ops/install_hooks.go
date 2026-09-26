package ops

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
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
		Short: "Install the git pre-commit that runs the gates over staged files",
		Long: `Writes .git/hooks/pre-commit to run 'anchors check --changed <file>
--no-record' on each staged file. A blocking gate that fails bars the commit.

What passes and what is barred, when the file is not in the map:
  • NOT GOVERNED (matches no layer of 'layers:' — package.json, lockfile, CI):
    ignored. Anchors has no jurisdiction over it.
  • GOVERNED but outside the map (new file, 'map build' did not run): BARS. Outside the
    map no gate confronts it — the triad is not charged and the commit would go on to
    certify work that nobody verified. Run 'anchors map build'.

Incremental: it validates only what the commit touches. It does not write to the map nor open issues
(--no-record). Idempotent — reinstalling is safe; use --force to overwrite an
existing pre-commit that was not written by this command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			return runInstallHooks(absRoot, force)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing pre-commit not managed by anchors")
	return cmd
}

// hookMarker identifica um hook escrito por este comando — permite reinstalar sem
// --force e distinguir de um hook artesanal do usuário.
const hookMarker = "# managed-by: anchors install-hooks"

// legacyCommitMsgMarker is the header the commit-msg hook carried before it carried
// hookMarker. That hook was then taken for the user's on every reinstall without --force,
// and never updated; recognising the old header is what lets existing installs catch up.
const legacyCommitMsgMarker = "# anchors:hook — instalado por 'anchors install-hooks'"

// writtenByAnchors says whether a hook on disk is one this command wrote.
func writtenByAnchors(hook []byte) bool {
	return strings.Contains(string(hook), hookMarker) ||
		strings.Contains(string(hook), legacyCommitMsgMarker)
}

func runInstallHooks(root string, force bool) error {
	// 1. exige anchors.yaml (é um projeto anchors?)
	if _, err := os.Stat(filepath.Join(root, config.DefaultFile)); err != nil {
		return fmt.Errorf("%s not found in %s — run `anchors init` first", config.DefaultFile, root)
	}

	// 2. descobre o diretório de hooks do git (respeita core.hooksPath e worktrees).
	hooksDir, err := gitHooksDir(root)
	if err != nil {
		// O pre-commit VIVE dentro do repositório: sem ele não há onde instalar. Dizer
		// isso aqui evita que o usuário procure o problema no anchors.yaml.
		if msg := gitmeta.Explain(gitmeta.Check(root), "install the pre-commit"); msg != "" {
			return errors.New(msg)
		}
		return fmt.Errorf("locate the git hooks directory: %w", err)
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", hooksDir, err)
	}

	hookPath := filepath.Join(hooksDir, "pre-commit")

	// 3. se já existe e não é nosso, respeita (a menos que --force).
	if existing, rerr := os.ReadFile(hookPath); rerr == nil {
		if !writtenByAnchors(existing) && !force {
			return fmt.Errorf(
				"a pre-commit already exists in %s that was not written by anchors.\n"+
					"  review it and, if you want to replace it, run with --force",
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
		!writtenByAnchors(existing) && !force {
		fmt.Printf("⚠  %s exists and was not written by anchors — not overwritten.\n", msgHook)
	} else if err := os.WriteFile(msgHook, []byte(commitMsgScript), 0o755); err != nil {
		return fmt.Errorf("write %s: %w", msgHook, err)
	}

	// O PRE-PUSH é a rede de segurança do pre-commit: quem commitou antes do
	// congelamento, ou com `--no-verify`, ainda esbarra nele antes de o trabalho sair da
	// máquina. E ele confere SEM cache — o push é raro o bastante para pagar o fetch, e é
	// o último momento em que a informação ainda muda o desfecho.
	pushHook := filepath.Join(hooksDir, "pre-push")
	if existing, rerr := os.ReadFile(pushHook); rerr == nil &&
		!writtenByAnchors(existing) && !force {
		fmt.Printf("⚠  %s exists and was not written by anchors — not overwritten.\n", pushHook)
	} else if err := os.WriteFile(pushHook, []byte(prePushScript), 0o755); err != nil {
		return fmt.Errorf("write %s: %w", pushHook, err)
	}

	if err := os.WriteFile(hookPath, []byte(preCommitScript), 0o755); err != nil {
		return fmt.Errorf("write %s: %w", hookPath, err)
	}

	fmt.Printf("✓ pre-commit installed at %s\n", hookPath)
	if _, lookErr := exec.LookPath("anchors"); lookErr != nil {
		fmt.Println("⚠  the 'anchors' binary is not in the PATH — the hook will fail until you install it")
		fmt.Println("   (cd cli && GOFLAGS=-mod=mod go install ./cmd/anchors) and ensure $(go env GOPATH)/bin is in the PATH")
	}
	fmt.Println("  the hook runs `anchors check --changed` on the staged files; a blocking gate stops the commit.")
	fmt.Println("  a GOVERNED file outside the map also stops it (run `anchors map build`); a non-governed one is ignored.")
	fmt.Printf("✓ pre-push installed at %s\n", pushHook)
	fmt.Println("  both check whether the project is FROZEN on the remote before letting the work proceed.")

	// O MERGE DRIVER do mapa vem junto: sem ele, o git mescla o
	// `anchors.graph.yaml` como texto e apaga carimbo sem conflito e sem aviso.
	//
	// Medido no blue-eyes (co2-lab/anchors#12): um `git merge` removeu 1212 linhas do
	// mapa e 62 carimbos de julgamento. O aviso do `map build` não pega — a perda
	// acontece antes de o Anchors ser chamado.
	//
	// Instalado aqui e não no `init` porque é config LOCAL do clone (`git config`), e
	// o `init` escreve o que é do repositório. Quem clona roda `install-hooks`.
	installMergeDrivers(root)
	return nil
}

// installMergeDrivers registra os drivers de merge dos arquivos DERIVADOS do Anchors.
//
// São dois, e os dois nasceram do mesmo defeito medido: o git mescla como TEXTO um
// arquivo que tem estrutura, e o resultado é dano (no mapa) ou atrito diário (no
// progresso).
func installMergeDrivers(root string) {
	installMergeDriver(root, mergeDriver{
		nome:    "anchors-map",
		attr:    mapx.DefaultPath,
		rotulo:  "anchors: merges the map stamps",
		comando: "anchors map merge %O %A %B",
		assunto: "map",
		porque: "# The MAP is derived, and git merges it as TEXT — erasing judgment stamps\n" +
			"# with no conflict and no warning (measured: 62 at once). The driver merges the\n" +
			"# stamps from both sides. Register it with `anchors install-hooks`.",
		efeito: "`git merge` now MERGES the map stamps instead of merging it as text.",
	})
	installMergeDriver(root, mergeDriver{
		nome:    "anchors-progress",
		attr:    "*-progress.md",
		rotulo:  "anchors: merges the progress of the plans",
		comando: "anchors merge-progress %O %A %B",
		assunto: "progress",
		porque: "# The PROGRESS of a plan is marked by whoever delivers, and two branches that\n" +
			"# deliver specs of the same plan mark neighboring checkboxes: git asks for manual\n" +
			"# resolution every time (measured: three PRs in a row, identical resolution in all three).\n" +
			"# The driver merges both sides, and `[x]` beats `[ ]` — unchecking by merge would erase\n" +
			"# a delivery that already happened.",
		efeito: "`git merge` now MERGES the progress items; `[x]` beats `[ ]`.",
	})
}

// mergeDriver descreve um driver: as duas metades do registro, e a prosa que explica ao
// próximo leitor do `.gitattributes` por que a linha está lá.
type mergeDriver struct {
	nome    string // o nome no git config e no `merge=` do atributo
	attr    string // o padrão de arquivo no .gitattributes
	rotulo  string // merge.<nome>.name
	comando string // merge.<nome>.driver
	assunto string // como ele é chamado nas mensagens ("mapa", "progresso")
	porque  string // o comentário que precede a linha no .gitattributes
	efeito  string // o que muda, dito a quem acabou de instalar
}

// installMergeDriver registra um driver no git local e declara o atributo.
//
// São duas metades e as duas são necessárias: o `.gitattributes` (versionado, diz QUAL
// arquivo usa o driver) e o `git config` (local do clone, diz COMO chamá-lo). Sem a
// segunda, o git avisa que o driver não existe e cai no merge textual — que é o
// comportamento que se está evitando.
//
// Falha em silêncio de propósito: um projeto sem git, ou um `git config` que não roda, tem
// outro problema — e o `install-hooks` já reportou o que importa.
func installMergeDriver(root string, d mergeDriver) {
	linhaAttr := d.attr + " merge=" + d.nome

	if err := exec.Command("git", "-C", root, "config",
		"merge."+d.nome+".name", d.rotulo).Run(); err != nil {
		return
	}
	if err := exec.Command("git", "-C", root, "config",
		"merge."+d.nome+".driver", d.comando).Run(); err != nil {
		return
	}

	attr := filepath.Join(root, ".gitattributes")
	atual, _ := os.ReadFile(attr)
	if strings.Contains(string(atual), linhaAttr) {
		fmt.Printf("✓ %s merge driver already configured (.gitattributes + git config)\n", d.assunto)
		return
	}
	conteudo := string(atual)
	if conteudo != "" && !strings.HasSuffix(conteudo, "\n") {
		conteudo += "\n"
	}
	conteudo += "\n" + d.porque + "\n" + linhaAttr + "\n"
	if err := os.WriteFile(attr, []byte(conteudo), 0o644); err != nil {
		return
	}
	fmt.Printf("✓ %s merge driver installed (.gitattributes + git config)\n", d.assunto)
	fmt.Println("  " + d.efeito)
	fmt.Println("  commit the .gitattributes: it belongs to the repository, and each clone runs `install-hooks`.")
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
  echo "✗ pre-commit: the 'anchors' binary is not in the PATH."
  echo "  Install it (cd cli && GOFLAGS=-mod=mod go install ./cmd/anchors) and ensure \$(go env GOPATH)/bin is in the PATH."
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
# O cache vive no diretório do git, que não é versionado: ele é estado da MÁQUINA, não do
# projeto.
#
# O comando git rev-parse --git-dir E NAO "$ROOT/.git": num WORKTREE, o .git e um ARQUIVO
# que aponta para …/.git/worktrees/<nome>, e escrever dentro dele falha com Not a directory.
#
# MEDIDO: um agente trabalhando em worktree levou
#
#     .git/hooks/pre-commit: line 48: …/be-rev1/.git/anchors-freeze-cache: Not a directory
#
# e commitou com --no-verify — o hook inteiro deixou de rodar por causa do cache. Um hook
# que falha assim é pior que hook nenhum: ele treina quem o usa a contorná-lo.
#
# O --git-dir devolve o caminho certo nos dois casos: .git no clone comum, e o diretorio do
# worktree quando e um.
GITDIR="$(git rev-parse --git-dir)"
CACHE="$GITDIR/anchors-freeze-cache"
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
  echo "🛑 COMMIT REFUSED — the project is FROZEN."
  echo ""
  [ -n "$motivo" ] && echo "   Reason: $motivo" || \
    echo "   (no reason declared in the origin's anchors.yaml)"
  echo ""
  echo "   Stop whatever work you are doing: it will not be deliverable until the"
  echo "   thaw, and what you produce now ages against the fix that is coming."
  echo ""
  echo "   If you are the one who will FIX what caused the freeze, use --no-verify —"
  echo "   the brake exists to stop work by inertia, not the fix itself."
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
  echo "· pre-commit: no staged file is governed by the Structure — nothing to confront."
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
  if [ -x "$GITDIR/hooks/commit-msg" ]; then
    echo "──────────────────────────────────────────────────────────────"
    echo "· gates failed. If it is deliberate, declare it in the commit message:"
    echo "    [skip-<rule>@<CODE>: why]"
    echo "  Without that, the commit-msg blocks it."
    exit 0
  fi
  echo "──────────────────────────────────────────────────────────────"
  echo "✗ commit BLOCKED by the anchors gates. Fix the above and recommit."
  exit 1
fi

# Extensão: hooks locais do projeto (réguas que o anchors não cobre, ex.: lint de
# arquitetura de import). Cada executável em .git/hooks/pre-commit.d/ roda com os
# arquivos staged como argumentos; qualquer um que falhe barra o commit.
# (Portável a bash 3.2 do macOS — sem mapfile.)
HOOK_D="$GITDIR/hooks/pre-commit.d"
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
` + hookMarker + `
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

echo "· commit-msg: confronting with the message in hand"
set +e
(cd "$ROOT" && anchors verify --phase pre-commit --staged --no-record --commit-msg "$MSG_FILE")
STATUS=$?
set -e
if [ "$STATUS" -eq 3 ]; then
  exit 0
elif [ "$STATUS" -ne 0 ]; then
  echo "──────────────────────────────────────────────────────────────"
  echo "✗ commit BLOCKED by the anchors gates."
  echo "  If the failure is deliberate, declare it in the message:"
  echo "    [skip-<rule>@<CODE>: why]"
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
  echo "⚠  pre-push: could not reach '$remoto' — proceeding without checking the freeze."
  exit 0
fi

remota=$(git show "$remoto/$base:$CFG" 2>/dev/null || true)
[ -n "$remota" ] || exit 0

# CONGELADO NO REMOTO: barra, e mostra o motivo que está escrito lá.
if printf '%s' "$remota" | grep -qE '^enabled:[[:space:]]*false[[:space:]]*$'; then
  motivo=$(printf '%s' "$remota" | sed -n 's/^freeze_reason:[[:space:]]*//p' | head -1)
  echo ""
  echo "🛑 PUSH REFUSED — the project is FROZEN."
  echo ""
  [ -n "$motivo" ] && echo "   Reason: $motivo" || \
    echo "   (no reason declared in the anchors.yaml of $remoto/$base)"
  echo ""
  echo "   The work is not lost: it stays on your local branch until the thaw."
  echo "   Follow the freeze issue in the repository."
  echo ""
  echo "   If you are the one who will FIX what caused the freeze, use --no-verify"
  echo "   — the brake exists to stop work by inertia, not the fix itself."
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
  echo "⚠  pre-push: your $CFG differs from $remoto/$base."
  echo "   Run 'git pull --rebase $remoto $base' to see what changed — that is where a"
  echo "   freeze would be declared."
fi

# O BINÁRIO ATENDE AO MÍNIMO QUE O PROJETO DECLARA?
#
# 'min_version' no anchors.yaml do remoto é a alavanca para o momento em que uma correção
# precisa alcançar todo mundo antes que o trabalho continue: um formato de mapa que mudou,
# um gate que passou a pegar algo que passava batido, um defeito que produz dado errado em
# silêncio.
#
# Ele é DECLARADO, e por isso serve. A conferência abaixo, contra o 'gerado_por' do mapa,
# usa um campo DERIVADO — o próximo 'map build' de qualquer agente o reescreve com a versão
# dele, apagando a exigência sem ninguém decidir nada. Medido: o campo voltou a "dev" num
# projeto onde a release corrente era a v0.1.83.
#
# Este BARRA, e a diferença de força é deliberada: o 'gerado_por' avisa sobre uma
# divergência que talvez não importe; o 'min_version' é alguém dizendo "abaixo disto não".
min_ver=$(printf '%s' "$remota" | sed -n 's/^min_version:[[:space:]]*//p' | head -1 | tr -d '"'"'"' ')
local_ver=$(anchors --version 2>/dev/null | sed -n 's/^anchors version \([^ ]*\).*/\1/p')

if [ -n "$min_ver" ] && [ -n "$local_ver" ]; then
  # A comparação é ORDINAL e o shell não a faz sozinho: "0.1.9" > "0.1.84" em ordem de
  # texto, e é MENOR em versão. 'sort -V' resolve — se o menor dos dois for o mínimo, o
  # local atende.
  menor=$(printf '%s\n%s\n' "$min_ver" "$local_ver" | sort -V | head -1)
  if [ "$local_ver" != "$min_ver" ] && [ "$menor" = "$local_ver" ]; then
    echo ""
    echo "🛑 PUSH REFUSED — this project requires 'anchors' $min_ver or newer."
    echo ""
    echo "   Your binary is $local_ver."
    echo ""
    echo "   The minimum is declared in the anchors.yaml of $remoto/$base, and it exists for"
    echo "   when a fix needs to reach everyone before the work proceeds."
    echo ""
    echo "   Update with 'brew upgrade anchors' or"
    echo "   'go install github.com/co2-lab/anchors/cmd/anchors@latest'."
    echo ""
    exit 1
  fi
fi

# O BINÁRIO DIVERGE DO QUE GRAVOU O MAPA? — avisa, não barra.
#
# Uma versão diferente não torna o trabalho errado; torna o mapa suscetível a oscilar. E
# barrar por isso transformaria "atualize quando puder" em "pare agora", caro no meio de
# uma entrega.
remota_ver=$(printf '%s' "$(git show "$remoto/$base:anchors.graph.yaml" 2>/dev/null || true)" \
  | sed -n 's/^gerado_por:[[:space:]]*//p' | head -1)

if [ -n "$remota_ver" ] && [ -n "$local_ver" ] && [ "$remota_ver" != "$local_ver" ]; then
  echo ""
  echo "⚠  pre-push: the map on $remoto/$base was written by 'anchors $remota_ver', and"
  echo "   your binary is $local_ver."
  echo ""
  echo "   Different versions write the map in different ways: the older one UNDOES what"
  echo "   the newer one wrote, and the map oscillates with a conflict on every PR without anyone"
  echo "   doing anything wrong."
  echo ""
  echo "   Update with 'brew upgrade anchors' or"
  echo "   'go install github.com/co2-lab/anchors/cmd/anchors@latest'."
  echo ""
fi
exit 0
`
