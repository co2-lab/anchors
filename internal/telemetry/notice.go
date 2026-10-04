// @anchors
//   code: NTAPN
//   ref: TLNTT

package telemetry

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
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
// ONCE PER PROJECT ON THIS MACHINE. O marcador vive em `.anchors/` da raiz do projeto, que
// não é versionado — cada pessoa vê uma vez por projeto (por clone), e o aviso não vira ruído
// a partir do segundo comando. Repetir seria pior que não avisar: quem lê a mesma coisa toda
// vez para de ler. (The notice once said "once per machine"; the marker was always per
// project root, and the text now says so.)

// noticeFile marca que este projeto, nesta máquina, já viu. At `.anchors/` porque é estado
// local: versioná-lo faria a primeira pessoa a commitar calar o aviso para todo o time.
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
	_ = os.WriteFile(p, []byte("the telemetry notice was shown for this project\n"), 0o644)
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
	// The text comes from the catalog, in the project's language: it was a Portuguese
	// literal, shown as is to every project whatever its `lang:`. The frame has no right
	// border, so a translation of any line length still closes it.
	var b strings.Builder
	b.WriteString("\n┌─ telemetry ─────────────────────────────────────────────────────────────\n")
	for _, l := range strings.Split(i18n.T("telemetry.notice"), "\n") {
		b.WriteString(strings.TrimRight("│ "+l, " ") + "\n")
	}
	b.WriteString("└─────────────────────────────────────────────────────────────────────────\n\n")
	fmt.Fprint(w, b.String())
	MarkNoticed(root)
}
