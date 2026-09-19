package ops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/initx"
)

// gitStep é o PRIMEIRO passo do `anchors init`, antes de escanear: sem versionamento,
// o carimbo de alteração, a cobertura de diff e o pre-commit ficam desligados — e não
// ruidosamente, o que é pior (ver `initx.AvisoGit`).
//
// A falta de git é AVISO aqui, nunca erro: o init segue e escreve o anchors.yaml de
// qualquer jeito. Quem depende de git cobra no seu próprio ponto de uso, onde a
// mensagem pode dizer qual comando precisa do quê.
//
// Devolve false só quando um prompt não pôde rodar (sem TTY) — o chamador aborta antes
// de tocar o disco, pela mesma régua que já protege o anchors.yaml de nascer vazio.
func gitStep(root string) bool {
	estado := initx.DetectGit(root, gitInstalled())
	if estado == initx.GitPronto {
		return true
	}

	fmt.Println(initx.AvisoGit(estado))

	// Git NÃO INSTALADO avisa e para por aqui: não há o que oferecer, porque o
	// `git init` que a pergunta prometeria falharia ao ser aceito. O aviso ainda vale —
	// `init` e `doctor` são os comandos cujo trabalho é justamente antecipar o problema,
	// antes que ele apareça deslocado no meio de outra coisa. O init segue normalmente.
	if !initx.OfferAction(estado) {
		fmt.Println()
		return true
	}

	var pergunta string
	if estado == initx.GitNaoIniciado {
		pergunta = "Initialize the git repository now (git init + .gitignore + first commit)?"
	} else {
		pergunta = "Make the first commit now?"
	}
	if !askConfirmDefault(pergunta, true) {
		if erroDePrompt {
			return false // sem TTY: não age, e o chamador aborta sem escrever nada
		}
		// Recusar é legítimo, e o init segue. Mas seguir em silêncio deixaria o usuário
		// descobrir o buraco só quando um comando entregasse menos do que ele espera —
		// e sem ligar o sintoma à causa. Então nomeia-se AGORA o que fica incompleto.
		fmt.Println("  proceeding without git. Some actions will not complete until there is a repository:")
		fmt.Println("    · the anchors' change stamp (updated_at) is left without a commit date;")
		fmt.Println("    · `anchors coverage --diff` has nothing to compare against;")
		fmt.Println("    · `anchors install-hooks` has nowhere to install the pre-commit;")
		fmt.Println("    · the `github` mode of the workflow needs a repository.")
		fmt.Println("  `git init` later, at any moment, turns all of that back on.")
		fmt.Println()
		return true
	}
	if erroDePrompt {
		return false
	}

	if err := initGit(root, estado); err != nil {
		// Falhar aqui não derruba o init: o anchors.yaml continua valendo a pena. Mas a
		// mensagem precisa nomear a causa, e não sumir no meio do resto.
		fmt.Printf("⚠  could not initialize git: %v\n", err)
		fmt.Println("   proceed without it; `anchors init` does not depend on git to write the anchors.yaml.")
		fmt.Println()
	}
	return true
}

// gitInstalled diz se o binário está no PATH. É a metade da distinção que separa
// "não há git nesta máquina" de "não há git neste projeto".
func gitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// initGit leva a raiz de `estado` até ter um commit. Faz só o que falta: num repo já
// criado (GitSemCommit), não roda `git init` de novo.
func initGit(root string, estado initx.GitState) error {
	if estado == initx.GitNaoIniciado {
		if out, err := runGit(root, "init"); err != nil {
			return fmt.Errorf("git init: %s", out)
		}
		fmt.Println("✓ repository created (git init)")
	}

	// O .gitignore é semeado, nunca sobrescrito: num repo que já tem um, o do usuário
	// vale mais que o nosso padrão.
	ignore := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(ignore); os.IsNotExist(err) {
		if err := os.WriteFile(ignore, []byte(initx.GitignorePadrão), 0o644); err != nil {
			return fmt.Errorf("write .gitignore: %w", err)
		}
		fmt.Println("✓ .gitignore seeded")
	}

	if out, err := runGit(root, "add", "-A"); err != nil {
		return fmt.Errorf("git add: %s", out)
	}

	// `git commit` falha quando não há nada staged — e num diretório vazio isso é o
	// caso normal, não um erro a reportar. `--allow-empty` dá o HEAD que o resto do
	// Anchors precisa (gitmeta.Head, coverage --diff, o pre-commit) sem exigir que o
	// projeto já tenha arquivo nenhum.
	if out, err := runGit(root, "commit", "--allow-empty", "-m", initx.FirstCommitMessage); err != nil {
		// Identidade não configurada é a falha mais provável aqui, e a mensagem crua do
		// git é longa; vale nomear o conserto.
		if strings.Contains(out, "user.email") || strings.Contains(out, "user.name") {
			return fmt.Errorf("git does not know who you are — run:\n"+
				"     git config --global user.name  \"Your Name\"\n"+
				"     git config --global user.email \"you@example.com\"\n"+
				"   (detail: %s)", firstLine(out))
		}
		return fmt.Errorf("git commit: %s", firstLine(out))
	}
	fmt.Printf("✓ first commit (%s)\n\n", initx.FirstCommitMessage)
	return nil
}

// runGit executa git na raiz e devolve saída combinada — o git escreve erro em stderr,
// e sem ele a mensagem de falha chegaria vazia.
func runGit(root string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
