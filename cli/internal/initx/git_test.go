package initx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// raizIsolada devolve um diretório que NÃO está dentro de nenhum repo git. Sem isso o
// teste seria falso-positivo: t.TempDir() no macOS cai em /var/folders, mas rodar a
// suíte de dentro do próprio repo do Anchors com um caminho relativo acharia o `.git`
// de cima e todo estado viraria GitPronto.
func raizIsolada(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if _, dentro := repoAcima(dir); dentro {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}
	return dir
}

// A distinção que motiva os quatro estados: "não há git" na máquina é OUTRA coisa de
// "não há git neste projeto", porque a ação do Anchors é diferente — sem o binário não
// existe `git init` a oferecer.
func TestGitNaoInstaladoNaoEhOMesmoQueNaoIniciado(t *testing.T) {
	dir := raizIsolada(t)

	semBinario := DetectaGit(dir, false)
	if semBinario != GitNaoInstalado {
		t.Fatalf("sem o binário, o estado tem de ser GitNaoInstalado, foi %v", semBinario)
	}
	comBinario := DetectaGit(dir, true)
	if comBinario != GitNaoIniciado {
		t.Fatalf("com o binário e sem repo, o estado tem de ser GitNaoIniciado, foi %v", comBinario)
	}

	// E a ação difere: só um dos dois tem o que oferecer.
	if OfereceAcao(GitNaoInstalado) {
		t.Error("não há `git init` a oferecer numa máquina sem git — perguntar prometeria o que falharia")
	}
	if !OfereceAcao(GitNaoIniciado) {
		t.Error("com git instalado e sem repo, a oferta de inicializar é justamente o passo")
	}
}

// Cada estado ensina uma coisa diferente. Uma mensagem que mande "instalar" quem já tem
// git faz o usuário procurar o problema onde ele não está.
func TestAvisoDizOQueFazerEmCadaEstado(t *testing.T) {
	casos := []struct {
		estado EstadoGit
		contém string
	}{
		{GitNaoInstalado, "não está instalado"},
		{GitNaoIniciado, "não está sob git"},
		{GitSemCommit, "nenhum commit"},
	}
	for _, c := range casos {
		av := AvisoGit(c.estado)
		if !strings.Contains(av, c.contém) {
			t.Errorf("aviso de %v deveria conter %q, foi: %s", c.estado, c.contém, av)
		}
	}
	if AvisoGit(GitPronto) != "" {
		t.Error("repo pronto não tem o que avisar")
	}
}

// Repo criado e sem commit é estado PRÓPRIO: dizer "tem git" aqui seria dizer que
// `git log`/`git diff` funcionam, e eles não funcionam sem HEAD.
func TestRepoSemCommitNaoContaComoPronto(t *testing.T) {
	dir := raizIsolada(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git", "refs", "heads"), 0o755); err != nil {
		t.Fatal(err)
	}

	if e := DetectaGit(dir, true); e != GitSemCommit {
		t.Fatalf("repo sem ref nenhuma tem de ser GitSemCommit, foi %v", e)
	}
	if !OfereceAcao(GitSemCommit) {
		t.Error("falta o commit — e é isso que o Anchors tem a oferecer aqui")
	}
}

func TestRepoComCommitEhPronto(t *testing.T) {
	dir := raizIsolada(t)
	heads := filepath.Join(dir, ".git", "refs", "heads")
	if err := os.MkdirAll(heads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(heads, "main"), []byte("abc123\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if e := DetectaGit(dir, true); e != GitPronto {
		t.Fatalf("repo com ref em heads/ é pronto, foi %v", e)
	}
	if OfereceAcao(GitPronto) {
		t.Error("nada a oferecer num repo pronto")
	}
}

// `git gc` empacota as refs: heads/ fica vazio e o histórico vive em packed-refs. Ler
// só heads/ classificaria um repo com anos de histórico como "sem commit".
func TestRepoComRefsEmpacotadasEhPronto(t *testing.T) {
	dir := raizIsolada(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git", "refs", "heads"), 0o755); err != nil {
		t.Fatal(err)
	}
	packed := "# pack-refs with: peeled fully-peeled sorted \nabc123 refs/heads/main\n"
	if err := os.WriteFile(filepath.Join(dir, ".git", "packed-refs"), []byte(packed), 0o644); err != nil {
		t.Fatal(err)
	}

	if e := DetectaGit(dir, true); e != GitPronto {
		t.Fatalf("refs empacotadas contam como commit, foi %v", e)
	}
}

// Rodar `init` numa subpasta de um repo existente NÃO deve oferecer criar repo
// aninhado: o histórico do subprojeto sumiria do repo de cima, e isso não se desfaz
// com revert.
func TestSubpastaDeRepoExistenteNaoOfereceRepoAninhado(t *testing.T) {
	dir := raizIsolada(t)
	heads := filepath.Join(dir, ".git", "refs", "heads")
	if err := os.MkdirAll(heads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(heads, "main"), []byte("abc123\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "pacotes", "app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if e := DetectaGit(sub, true); e != GitPronto {
		t.Fatalf("subpasta de repo existente já está versionada, foi %v", e)
	}
}

// Worktree e submódulo têm `.git` como ARQUIVO (ponteiro `gitdir: …`). Tratar como
// "sem repo" ofereceria reinicializar por cima de um worktree válido.
func TestGitComoArquivoEhPronto(t *testing.T) {
	dir := raizIsolada(t)
	if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /outro/lugar/.git/worktrees/x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if e := DetectaGit(dir, true); e != GitPronto {
		t.Fatalf("`.git` como arquivo é worktree/submódulo — repo real existe, foi %v", e)
	}
}

// O .gitignore semeado cobre o que o PRÓPRIO Anchors gera. Não adivinha stack: neste
// ponto do init o preset ainda não foi escolhido.
func TestGitignoreCobreOQueOAnchorsGera(t *testing.T) {
	for _, esperado := range []string{".DS_Store", ".anchors/cache/"} {
		if !strings.Contains(GitignorePadrão, esperado) {
			t.Errorf(".gitignore padrão deveria ignorar %q", esperado)
		}
	}
	for _, stack := range []string{"node_modules", "target/", "vendor/"} {
		if strings.Contains(GitignorePadrão, stack) {
			t.Errorf(".gitignore não deve chutar stack (%q): o preset ainda não foi escolhido", stack)
		}
	}
}
