package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O `-progress.md` é o único artefato do Anchors que nada confrontava — e "fora do mapa"
// virou "fora de qualquer verificação", que não é a mesma coisa.
//
// Medido no blue-eyes: ao criar os 17 progressos, transportei o estado dos checkboxes que
// viviam nos planos, e os planos estavam desatualizados. O progresso do `0002` dizia 6
// itens abertos com 7 das 8 specs já no disco. CINCO itens mentiam, e eu transportei a
// mentira fielmente.
func TestProgressHonest(t *testing.T) {
	monta := func(t *testing.T, progresso string, existentes ...string) (string, mapx.Node) {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
			t.Fatal(err)
		}
		plano := "plans/0002-plataforma.md"
		if err := os.WriteFile(filepath.Join(root, plano), []byte("# Plano\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if progresso != "" {
			prog := filepath.Join(root, "plans", "0002-plataforma-progress.md")
			if err := os.WriteFile(prog, []byte(progresso), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		for _, f := range existentes {
			p := filepath.Join(root, f)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte("x\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root, mapx.Node{ID: plano, Kind: mapx.KindPlan}
	}

	t.Run("[x] com arquivo AUSENTE reprova — o dano mais caro", func(t *testing.T) {
		root, n := monta(t, "- [x] `packages/shared/Nunca.spec.md` — nunca existiu\n")
		v, d := checkProgressHonest("", n, root, nil, nil)
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if !strings.Contains(d, "NÃO EXISTE") || !strings.Contains(d, "Nunca.spec.md") {
			t.Errorf("o laudo não nomeia o problema:\n%s", d)
		}
	})

	t.Run("[ ] com arquivo EXISTINDO reprova — produz retrabalho", func(t *testing.T) {
		root, n := monta(t,
			"- [ ] `packages/shared/Feito.spec.md` — entregue e não marcado\n",
			"packages/shared/Feito.spec.md")
		v, d := checkProgressHonest("", n, root, nil, nil)
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if !strings.Contains(d, "JÁ EXISTE") {
			t.Errorf("o laudo não nomeia o problema:\n%s", d)
		}
	})

	t.Run("progresso HONESTO passa", func(t *testing.T) {
		root, n := monta(t,
			"- [x] `a.spec.md` — feito\n- [ ] `b.spec.md` — por fazer\n",
			"a.spec.md")
		if v, d := checkProgressHonest("", n, root, nil, nil); v != Pass {
			t.Errorf("veredito %v, queria Pass: %s", v, d)
		}
	})

	// Item em PROSA não é confrontável, e reprovar por ele cobraria de quem escreveu um
	// progresso legítimo — a fase pode ser "revisar com o time".
	t.Run("item sem caminho é ignorado", func(t *testing.T) {
		root, n := monta(t, "- [ ] revisar com o time antes de seguir\n- [x] alinhado na daily\n")
		if v, d := checkProgressHonest("", n, root, nil, nil); v != Pass {
			t.Errorf("veredito %v, queria Pass — prosa não é confrontável: %s", v, d)
		}
	})

	// AUSÊNCIA do progresso não é falha DESTE gate: o plano pode ter nascido antes do
	// mecanismo. São duas perguntas diferentes — "existe?" e "é verdade?".
	t.Run("sem progresso ao lado, Skip com o caminho da correção", func(t *testing.T) {
		root, n := monta(t, "")
		v, d := checkProgressHonest("", n, root, nil, nil)
		if v != Skip {
			t.Fatalf("veredito %v, queria Skip: %s", v, d)
		}
		if !strings.Contains(d, "anchors new progress --for") {
			t.Errorf("o Skip não diz como corrigir:\n%s", d)
		}
	})

	t.Run("nó que não é plano faz Skip", func(t *testing.T) {
		root := t.TempDir()
		n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
		if v, _ := checkProgressHonest("", n, root, nil, nil); v != Skip {
			t.Errorf("veredito %v, queria Skip", v)
		}
	})
}

// O sufixo tem UMA definição, no `scan` — é ele que precisa manter o arquivo fora do
// mapa. Duas constantes divergiriam em silêncio, e o gate passaria a procurar um arquivo
// que não existe.
func TestProgressHonest_caminhoDerivaDoScan(t *testing.T) {
	casos := map[string]string{
		"plans/0002-plataforma.md": "plans/0002-plataforma-progress.md",
		"plans/0017-mutacao.md":    "plans/0017-mutacao-progress.md",
		"sem-extensao":             "sem-extensao-progress.md",
	}
	for plano, quer := range casos {
		if got := progressPathOf(plano); got != quer {
			t.Errorf("progressPathOf(%q) = %q, queria %q", plano, got, quer)
		}
	}
}
