package telemetry

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// --- O AVISO, e por que ele vive no PersistentPreRunE ---
//
// A telemetry é ligada por padrão, e o que torna essa escolha honesta é o aviso chegar a
// quem não pediu por ele. Quatro candidatos foram considerados:
//
//	init     → só alcança projeto NOVO. Quem instala o binário num projeto que outra
//	           pessoa configurou nunca roda `init` — e é o caso comum num time.
//	next     → só alcança quem PEDE CARD. Um revisor que roda `check` e `judge` não vê.
//	doctor   → é opcional, e muita gente nunca o roda.
//	pre-run  → TODOS os 43 comandos passam por ele, e ele roda ANTES do comando.
//
// A última propriedade é a que decide: o aviso aparece antes de o primeiro evento sair.
//
// UMA VEZ POR MÁQUINA. O marcador vive em `.anchors/`, que não é versionado — cada pessoa
// vê uma vez, e o aviso não vira ruído a partir do segundo comando. Repetir seria pior que
// não avisar: quem lê a mesma coisa toda vez para de ler.

// noticeFile marca que esta máquina já viu. At `.anchors/` porque é estado da MÁQUINA:
// versioná-lo faria a primeira pessoa a commitar calar o aviso para todo o time.
const noticeFile = ".anchors/telemetry-noticed"

// AlreadyNoticed responde se esta máquina já viu o aviso neste projeto.
func AlreadyNoticed(root string) bool {
	_, err := os.Stat(filepath.Join(root, noticeFile))
	return err == nil
}

// MarkNoticed registra que o aviso foi mostrado.
//
// Falha em silêncio de propósito: se o diretório não puder ser escrito, o pior que acontece
// é o aviso aparecer de novo — e aparecer duas vezes é melhor que travar um comando.
func MarkNoticed(root string) {
	p := filepath.Join(root, noticeFile)
	if os.MkdirAll(filepath.Dir(p), 0o755) != nil {
		return
	}
	_ = os.WriteFile(p, []byte("a telemetry foi anunciada nesta máquina\n"), 0o644)
}

// Notice escreve o aviso, uma vez.
//
// O TEXTO diz três coisas, nessa ordem: o que é coletado, o que NÃO é, e como desligar. A
// ordem importa — quem lê "coletamos dados" e não encontra o desligamento na mesma tela
// assume o pior, e assume certo na maioria dos produtos.
func Notice(w io.Writer, root string) {
	if AlreadyNoticed(root) {
		return
	}
	fmt.Fprint(w, `
┌─ telemetry ──────────────────────────────────────────────────────────────┐
│ O Anchors envia eventos de DECISÃO — quantos candidatos o claim viu, qual  │
│ gate reprovou, em que estado um turno terminou. É como o produto encontra  │
│ padrões que ninguém vê de dentro de um projeto só.                        │
│                                                                           │
│ NÃO envia: conteúdo de arquivo, spec, diff, mensagem de erro, nome de      │
│ repositório, de usuário ou de branch. Só números e o vocabulário do        │
│ próprio Anchors.                                                          │
│                                                                           │
│ Para desligar, qualquer um destes:                                        │
│     export ANCHORS_TELEMETRY=off                                          │
│     anchors telemetry off                                                 │
│     telemetry: off      (no anchors.yaml, vale para o projeto)            │
│                                                                           │
│ Este aviso aparece uma vez por máquina.                                   │
└───────────────────────────────────────────────────────────────────────────┘

`)
	MarkNoticed(root)
}
