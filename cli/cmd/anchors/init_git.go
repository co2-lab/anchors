package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/initx"
)

// etapaGit é o PRIMEIRO passo do `anchors init`, antes de escanear: sem versionamento,
// o carimbo de alteração, a cobertura de diff e o pre-commit ficam desligados — e não
// ruidosamente, o que é pior (ver `initx.AvisoGit`).
//
// A falta de git é AVISO aqui, nunca erro: o init segue e escreve o anchors.yaml de
// qualquer jeito. Quem depende de git cobra no seu próprio ponto de uso, onde a
// mensagem pode dizer qual comando precisa do quê.
//
// Devolve false só quando um prompt não pôde rodar (sem TTY) — o chamador aborta antes
// de tocar o disco, pela mesma régua que já protege o anchors.yaml de nascer vazio.
func etapaGit(root string) bool {
	estado := initx.DetectaGit(root, gitInstalado())
	if estado == initx.GitPronto {
		return true
	}

	fmt.Println(initx.AvisoGit(estado))

	// Git NÃO INSTALADO avisa e para por aqui: não há o que oferecer, porque o
	// `git init` que a pergunta prometeria falharia ao ser aceito. O aviso ainda vale —
	// `init` e `doctor` são os comandos cujo trabalho é justamente antecipar o problema,
	// antes que ele apareça deslocado no meio de outra coisa. O init segue normalmente.
	if !initx.OfereceAcao(estado) {
		fmt.Println()
		return true
	}

	var pergunta string
	if estado == initx.GitNaoIniciado {
		pergunta = "Iniciar o repositório git agora (git init + .gitignore + primeiro commit)?"
	} else {
		pergunta = "Fazer o primeiro commit agora?"
	}
	if !askConfirmDefault(pergunta, true) {
		if erroDePrompt {
			return false // sem TTY: não age, e o chamador aborta sem escrever nada
		}
		// Recusar é legítimo, e o init segue. Mas seguir em silêncio deixaria o usuário
		// descobrir o buraco só quando um comando entregasse menos do que ele espera —
		// e sem ligar o sintoma à causa. Então nomeia-se AGORA o que fica incompleto.
		fmt.Println("  seguindo sem git. Algumas ações não vão se completar até haver repositório:")
		fmt.Println("    · o carimbo de alteração das âncoras (updated_at) fica sem data de commit;")
		fmt.Println("    · `anchors coverage --diff` não tem contra o que comparar;")
		fmt.Println("    · `anchors install-hooks` não tem onde instalar o pre-commit;")
		fmt.Println("    · o modo `github` do fluxo de trabalho precisa de repositório.")
		fmt.Println("  `git init` depois, a qualquer momento, religa tudo isso.")
		fmt.Println()
		return true
	}
	if erroDePrompt {
		return false
	}

	if err := iniciaGit(root, estado); err != nil {
		// Falhar aqui não derruba o init: o anchors.yaml continua valendo a pena. Mas a
		// mensagem precisa nomear a causa, e não sumir no meio do resto.
		fmt.Printf("⚠  não deu para iniciar o git: %v\n", err)
		fmt.Println("   siga sem ele; `anchors init` não depende de git para escrever o anchors.yaml.")
		fmt.Println()
	}
	return true
}

// gitInstalado diz se o binário está no PATH. É a metade da distinção que separa
// "não há git nesta máquina" de "não há git neste projeto".
func gitInstalado() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// iniciaGit leva a raiz de `estado` até ter um commit. Faz só o que falta: num repo já
// criado (GitSemCommit), não roda `git init` de novo.
func iniciaGit(root string, estado initx.EstadoGit) error {
	if estado == initx.GitNaoIniciado {
		if out, err := rodaGit(root, "init"); err != nil {
			return fmt.Errorf("git init: %s", out)
		}
		fmt.Println("✓ repositório criado (git init)")
	}

	// O .gitignore é semeado, nunca sobrescrito: num repo que já tem um, o do usuário
	// vale mais que o nosso padrão.
	ignore := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(ignore); os.IsNotExist(err) {
		if err := os.WriteFile(ignore, []byte(initx.GitignorePadrão), 0o644); err != nil {
			return fmt.Errorf("escrever .gitignore: %w", err)
		}
		fmt.Println("✓ .gitignore semeado")
	}

	if out, err := rodaGit(root, "add", "-A"); err != nil {
		return fmt.Errorf("git add: %s", out)
	}

	// `git commit` falha quando não há nada staged — e num diretório vazio isso é o
	// caso normal, não um erro a reportar. `--allow-empty` dá o HEAD que o resto do
	// Anchors precisa (gitmeta.Head, coverage --diff, o pre-commit) sem exigir que o
	// projeto já tenha arquivo nenhum.
	if out, err := rodaGit(root, "commit", "--allow-empty", "-m", initx.MensagemPrimeiroCommit); err != nil {
		// Identidade não configurada é a falha mais provável aqui, e a mensagem crua do
		// git é longa; vale nomear o conserto.
		if strings.Contains(out, "user.email") || strings.Contains(out, "user.name") {
			return fmt.Errorf("o git não sabe quem é você — rode:\n"+
				"     git config --global user.name  \"Seu Nome\"\n"+
				"     git config --global user.email \"voce@exemplo.com\"\n"+
				"   (detalhe: %s)", primeiraLinha(out))
		}
		return fmt.Errorf("git commit: %s", primeiraLinha(out))
	}
	fmt.Printf("✓ primeiro commit (%s)\n\n", initx.MensagemPrimeiroCommit)
	return nil
}

// rodaGit executa git na raiz e devolve saída combinada — o git escreve erro em stderr,
// e sem ele a mensagem de falha chegaria vazia.
func rodaGit(root string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func primeiraLinha(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
