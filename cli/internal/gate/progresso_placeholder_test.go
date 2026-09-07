package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// escreveProgresso monta um plano e o companheiro dele num diretório temporário.
func escreveProgresso(t *testing.T, plano, progresso string) (root, planoPath string) {
	t.Helper()
	root = t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	planoPath = "plans/0017-mutacao.md"
	if err := os.WriteFile(filepath.Join(root, planoPath), []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}
	prog := filepath.Join(root, "plans/0017-mutacao-progress.md")
	if err := os.WriteFile(prog, []byte(progresso), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, planoPath
}

// O PLACEHOLDER, medido no blue-eyes.
//
// O progresso do plano 0017 tinha, na fase F02:
//
//	- [ ] TODO: um item por spec que esta fase semeia
//
// O `anchors new progress` o escreve quando a fase não semeia nada, e ele deveria sair
// quando alguém decide o que a fase faz. Ficou.
//
// As três direções deste gate não o veem: as duas primeiras confrontam itens que CITAM
// CAMINHO (este não cita), e a terceira olha as sementes do plano (esta fase não semeia).
//
// E o efeito é o oposto do que o gate protege: um `[ ]` eterno faz o plano parecer
// incompleto para sempre — o `anchors next` volta a ele, e quem lê não sabe se falta
// trabalho ou falta limpar o arquivo.
func TestProgressHonest_acusaOPlaceholderDeTODO(t *testing.T) {
	plano := "<!-- @anchors\ncode: MTUAO\n-->\n# Plano 0017\n\n## Fases\n\n### MTUAO-F02 — o CI\n\nEsta fase decide quando a mutação roda.\n"
	progresso := "# Progresso — MTUAO\n\n## MTUAO-F02 — o CI\n\n- [ ] TODO: um item por spec que esta fase semeia\n"

	root, planoPath := escreveProgresso(t, plano, progresso)
	n := mapx.Node{ID: planoPath, Kind: mapx.KindPlan, Code: "MTUAO"}

	v, msg := checkProgressHonest(plano, n, root, nil, nil)

	if v != Fail {
		t.Fatalf("o placeholder passou: %v — %s", v, msg)
	}
	if !strings.Contains(strings.ToUpper(msg), "TODO") {
		t.Errorf("a mensagem não nomeia o TODO:\n%s", msg)
	}
}

// Item que cita caminho e existe continua passando: a quarta direção não pode acusar o
// caso normal.
func TestProgressHonest_naoAcusaItemComCaminhoValido(t *testing.T) {
	plano := "<!-- @anchors\ncode: MTUAO\n-->\n# Plano\n\n- [ ] `packages/shared/X.spec.md` — a spec\n"
	progresso := "# Progresso\n\n- [x] `packages/shared/X.spec.md` — a spec\n"

	root, planoPath := escreveProgresso(t, plano, progresso)
	if err := os.MkdirAll(filepath.Join(root, "packages/shared"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "packages/shared/X.spec.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	n := mapx.Node{ID: planoPath, Kind: mapx.KindPlan, Code: "MTUAO"}

	if v, msg := checkProgressHonest(plano, n, root, nil, nil); v != Pass {
		t.Errorf("o caso normal foi acusado: %v — %s", v, msg)
	}
}

// E PROSA não é item: um progresso pode ter texto explicativo com a palavra TODO, e uma
// linha que não é `- [ ]` não é uma promessa de trabalho.
func TestProgressHonest_ignoraProsaComTODO(t *testing.T) {
	plano := "<!-- @anchors\ncode: MTUAO\n-->\n# Plano\n"
	progresso := "# Progresso\n\nEsta fase ainda tem TODO a decidir, e isso está no plano.\n"

	root, planoPath := escreveProgresso(t, plano, progresso)
	n := mapx.Node{ID: planoPath, Kind: mapx.KindPlan, Code: "MTUAO"}

	if v, msg := checkProgressHonest(plano, n, root, nil, nil); v == Fail {
		t.Errorf("prosa com TODO foi acusada como item: %s", msg)
	}
}
