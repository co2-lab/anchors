// @anchors
//   code: GTINA
//   ref: GTSTG

package initx

import (
	"os"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/i18n"
)

// GitState é o que o `init` precisa saber sobre o versionamento ANTES de escanear.
// Git não é detalhe de conforto do Anchors: é o substrato de onde vêm o carimbo de
// alteração (`gitmeta`), a cobertura de diff, o pre-commit e — no modo `github` — o
// próprio `repo` da fila de trabalho. Um projeto sem git carrega e roda, mas boa parte
// do framework fica desligada em silêncio.
//
// "Não há git" tem DOIS significados, e o Anchors age diferente em cada um. Colapsar
// os dois num estado só produziria a pior mensagem possível: oferecer `git init` numa
// máquina sem git (a oferta falha na hora), ou mandar instalar git para quem já o tem
// (o usuário procura o problema onde ele não está). São estados separados porque a
// AÇÃO é separada.
type GitState int

const (
	// GitNaoInstalado — o binário `git` não está no PATH. O Anchors não tem o que
	// oferecer: não há `git init` a rodar. O aviso manda instalar, e não pergunta nada.
	GitNaoInstalado GitState = iota
	// GitNaoIniciado — git existe na máquina, mas a raiz não é repositório. Aqui a
	// oferta faz sentido: `git init` + `.gitignore` + primeiro commit.
	GitNaoIniciado
	// GitSemCommit — repo criado, mas sem HEAD. Estado à parte de propósito: `git log`,
	// `git diff --cached` e o pre-commit ainda não funcionam, então dizer "tem git"
	// aqui seria dizer que funciona quando não funciona. Só falta o commit.
	GitSemCommit
	// GitPronto — repo com pelo menos um commit. Nada a fazer.
	GitPronto
)

// DetectGit classifica a raiz. `gitNoPath` é injetado (e não consultado aqui) para
// manter esta função pura e testável sem depender do que está instalado na máquina de
// quem roda o teste — a suíte precisa cobrir o caso "sem git" mesmo rodando numa
// máquina com git.
func DetectGit(root string, gitNoPath bool) GitState {
	if !gitNoPath {
		return GitNaoInstalado
	}
	dotGit := filepath.Join(root, ".git")
	fi, err := os.Stat(dotGit)
	if err != nil {
		if _, dentro := repoAcima(root); dentro {
			return GitPronto // dentro de um repo existente: não é problema do init
		}
		return GitNaoIniciado
	}
	if !fi.IsDir() {
		return GitPronto // worktree/submódulo: `.git` é ponteiro, o repo real é outro
	}
	if hasCommit(dotGit) {
		return GitPronto
	}
	return GitSemCommit
}

// hasCommit diz se o repo já tem HEAD apontando para algo. Lê o disco em vez de rodar
// `git rev-parse`: num repo recém-criado o `.git/refs/heads` está vazio e não há
// arquivo de ref nenhum — é o sinal mais direto de "ainda não há commit".
func hasCommit(dotGit string) bool {
	heads := filepath.Join(dotGit, "refs", "heads")
	var achou bool
	_ = filepath.WalkDir(heads, func(_ string, d os.DirEntry, err error) error {
		if err == nil && d != nil && !d.IsDir() {
			achou = true
		}
		return nil
	})
	if achou {
		return true
	}
	// Refs empacotadas (`git gc` já rodou): heads/ pode estar vazio e o histórico
	// existir em packed-refs.
	b, err := os.ReadFile(filepath.Join(dotGit, "packed-refs"))
	return err == nil && len(b) > 0
}

// repoAcima procura um `.git` num diretório ANCESTRAL. Sem isso, rodar `anchors init`
// numa subpasta de um repo existente ofereceria criar um repo aninhado — que é quase
// sempre erro, e dos caros: o histórico do subprojeto some do repo de cima.
func repoAcima(root string) (string, bool) {
	dir := filepath.Dir(root)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, true
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			return "", false
		}
		dir = pai
	}
}

// AvisoGit é o texto do aviso para cada estado, ou "" quando não há o que avisar.
// Fica aqui, junto da decisão, para que o teste cubra a mensagem — é ela que ensina
// o usuário POR QUE git importa, e uma mensagem vaga aqui vira um passo pulado.
//
// Through i18n, in the project's `lang:` — the three texts were hard-coded in Portuguese.
func AvisoGit(e GitState) string {
	switch e {
	case GitNaoInstalado:
		return i18n.T("init.git.not_installed")
	case GitNaoIniciado:
		return i18n.T("init.git.not_initialized")
	case GitSemCommit:
		return i18n.T("init.git.no_commit")
	default:
		return ""
	}
}

// OfferAction diz se o Anchors tem algo a PROPOR neste estado. É a distinção que
// separa "não há git" instalado de não iniciado: sem o binário não há oferta nenhuma
// a fazer, e perguntar seria prometer uma ação que falharia ao ser aceita.
func OfferAction(e GitState) bool {
	return e == GitNaoIniciado || e == GitSemCommit
}

// GitignorePadrão é o .gitignore semeado no primeiro commit. Cobre o que o próprio
// Anchors gera e o lixo de sistema — NÃO tenta adivinhar a stack (node_modules, target,
// …): a linguagem ainda não foi lida neste ponto do init, e um .gitignore com regras
// de uma stack errada é pior que um curto.
//
// In English, like everything Anchors writes into a project's files.
const GitignorePadrão = `# system
.DS_Store
Thumbs.db

# generated by Anchors — the WORKING area is not committed.
#
# .anchors/ keeps what describes a RUN, not the project: the pending judgment queue,
# the mirror of the last check output, the cache. Versioning it puts into the diff
# noise that changes by itself on every command.
#
# The judgment queue above all must not be committed: it is work of the CURRENT commit,
# and check blocks while anything is pending. A pending judgment file in a PR is a sign
# that someone went around the gate.
#
# What describes the PROJECT lives OUTSIDE — the SBOM, for example, is born at the root.
# So this rule needs no exception, and an exception in .gitignore is where silence lives.
.anchors/

# the map's write lock: held for the time of one write, and left behind only by a process
# that died holding it.
anchors.graph.yaml.lock*
`

// FirstCommitMessage é o assunto do commit inicial. Diz o que é, não o que faz:
// quem lê o log um ano depois quer saber que aqui o projeto começou.
const FirstCommitMessage = "chore: start the repository"
