package health

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// raizForaDeRepo devolve um diretório que não está sob nenhum `.git` — senão o teste
// seria falso-positivo ao rodar de dentro do próprio repo do Anchors.
func raizForaDeRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if temRepo(dir) {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}
	return dir
}

// O doctor existe para avisar COM ANTECEDÊNCIA. Git ausente é a mesma classe de
// `ferramenta-ausente`: nada falha ruidosamente, então o silêncio é que custa.
func TestDoctorAvisaProjetoSemRepositorio(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado — este teste cobre o caso 'git existe, repo não'")
	}
	dir := raizForaDeRepo(t)

	fs := checkGitAusente(&config.Config{}, dir)

	if len(fs) != 1 {
		t.Fatalf("esperava 1 achado, veio %d: %+v", len(fs), fs)
	}
	if fs[0].Check != "git-ausente" || fs[0].Severity != Warn {
		t.Errorf("achado errado: %+v", fs[0])
	}
	// A mensagem tem de dizer o CONSERTO — iniciar, não instalar.
	if !strings.Contains(fs[0].Detail, "git init") {
		t.Errorf("a mensagem deveria mandar iniciar o repo: %s", fs[0].Detail)
	}
	if strings.Contains(fs[0].Detail, "PATH") {
		t.Errorf("git ESTÁ instalado aqui — mandar olhar o PATH manda procurar onde não está: %s", fs[0].Detail)
	}
}

// Repo existente não é achado — um "alerta" que nunca exige ação é ruído, e ruído
// recorrente treina a equipe a ignorar o doctor.
func TestDoctorNaoReclamaComRepositorio(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := raizForaDeRepo(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	if fs := checkGitAusente(&config.Config{}, dir); len(fs) != 0 {
		t.Errorf("repo presente não deveria gerar achado: %+v", fs)
	}
}

// Subpasta de um repo JÁ está versionada: reclamar ali mandaria criar repo aninhado.
func TestDoctorNaoReclamaEmSubpastaDeRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := raizForaDeRepo(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "pacotes", "app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if fs := checkGitAusente(&config.Config{}, sub); len(fs) != 0 {
		t.Errorf("subpasta de repo já está versionada: %+v", fs)
	}
}

// Worktree e submódulo têm `.git` como ARQUIVO (ponteiro `gitdir:`). Tratar como
// "sem repo" mandaria reinicializar por cima de um worktree válido.
func TestDoctorAceitaGitComoArquivo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := raizForaDeRepo(t)
	if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /outro/.git/worktrees/x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if fs := checkGitAusente(&config.Config{}, dir); len(fs) != 0 {
		t.Errorf("`.git` como arquivo é worktree/submódulo — repo real existe: %+v", fs)
	}
}

// No modo `github` a falta de repositório deixa de ser débito e vira impedimento: a
// fila de trabalho MORA nas issues de um repositório. A mensagem precisa dizer isso,
// senão o usuário lê "updated_at fica sem data" e adia o que o bloqueia hoje.
func TestDoctorNoModoGitHubDizQueAFilaNaoTemDeOndeVir(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := raizForaDeRepo(t)
	cfg := &config.Config{Workflow: &config.Workflow{
		Mode:   config.ModeGitHub,
		Repo:   "acme/exemplo",
		Labels: []string{"anchors"},
	}}

	fs := checkGitAusente(cfg, dir)

	if len(fs) != 1 {
		t.Fatalf("esperava 1 achado, veio %d", len(fs))
	}
	if !strings.Contains(fs[0].Detail, "fila de trabalho") {
		t.Errorf("no modo github a mensagem tem de nomear o impedimento real: %s", fs[0].Detail)
	}
}

// O caso que o doctor cobre e o `init` não consegue: git NÃO INSTALADO. Aqui o
// conserto é instalar, não iniciar — mandar rodar `git init` numa máquina sem git
// manda o usuário procurar o problema onde ele não está.
func TestDoctorSemBinarioMandaInstalarNaoIniciar(t *testing.T) {
	// `instalado=false` injetado: sem isto, este caso só rodaria numa máquina sem git,
	// e ficaria eternamente em SKIP — justamente a metade da distinção que ninguém
	// reproduz por acidente.
	fs := gitAusente(&config.Config{}, t.TempDir(), false)

	if len(fs) != 1 || fs[0].Subject != "git" {
		t.Fatalf("esperava achado sobre o binário: %+v", fs)
	}
	if strings.Contains(fs[0].Detail, "git init") {
		t.Errorf("sem o binário, `git init` não roda — mandá-lo é mandar ao lugar errado: %s", fs[0].Detail)
	}
}
