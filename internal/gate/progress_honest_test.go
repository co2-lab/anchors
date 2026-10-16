// @anchors
//   code: PHTPR
//   ref: PRHNP

package gate

import (
	"github.com/co2-lab/anchors/internal/testkit"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/progress"
	"github.com/co2-lab/anchors/internal/scan"
)

// trackPlan writes, in anchors.progress.md, the plan's demand as the legacy text describes it:
// each `- [ ]`/`- [x]` line a task (done when checked), each heading a stage. The tests keep
// the progress they always described; the file is the one the gate now reads.
func trackPlan(t *testing.T, root, plan, legacy string) {
	t.Helper()
	item := regexp.MustCompile(`^\s*- \[([ xX])\]\s+(.+?)\s*$`)
	stage := regexp.MustCompile(`^#{2,4}\s+(.+?)\s*$`)
	if err := progress.Update(root, "", func(f *progress.File) error {
		d, err := f.AddDemand("plan", plan)
		if err != nil {
			return err
		}
		parent := d
		for _, line := range strings.Split(legacy, "\n") {
			if m := stage.FindStringSubmatch(line); m != nil {
				if parent, err = f.AddStage(d, m[1]); err != nil {
					return err
				}
				continue
			}
			m := item.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			id, err := f.AddTask(parent, m[2], 1, false)
			if err != nil {
				return err
			}
			if strings.ToLower(m[1]) == "x" {
				if err := f.Done(id, "agent", "2026-10-12"); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

// O `-progress.md` é o único artefato do Anchors que nada confrontava — e "fora do mapa"
// virou "fora de qualquer verificação", que não é a mesma coisa.
//
// Medido no app de referência: ao criar os 17 progressos, transportei o estado dos checkboxes que
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
			trackPlan(t, root, plano, progresso)
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

	t.Run("PRHNP-B04: A ticked item whose file does not exist is failed", func(t *testing.T) {
		root, n := monta(t, "- [x] `packages/shared/Nunca.spec.md` — nunca existiu\n")
		v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if !strings.Contains(d, "NÃO EXISTE") && !strings.Contains(d, "DOES NOT EXIST") &&
			!strings.Contains(d, "NO EXISTE") {
			t.Errorf("o laudo não nomeia o problema:\n%s", d)
		}
	})

	t.Run("PRHNP-B11: The verdict names each offending path", func(t *testing.T) {
		root, n := monta(t, "- [x] `packages/shared/Nunca.spec.md` — nunca existiu\n")
		v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		// Quem lê não pode ter de comparar o arquivo com o disco à mão.
		if !strings.Contains(d, "Nunca.spec.md") {
			t.Errorf("o laudo não nomeia o caminho acusado:\n%s", d)
		}
	})

	t.Run("PRHNP-B05: An open item whose file already exists is failed", func(t *testing.T) {
		root, n := monta(t,
			"- [ ] `packages/shared/Feito.spec.md` — entregue e não marcado\n",
			"packages/shared/Feito.spec.md")
		v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if !strings.Contains(d, "JÁ EXISTE") && !strings.Contains(d, "ALREADY EXISTS") &&
			!strings.Contains(d, "YA EXISTE") {
			t.Errorf("o laudo não nomeia o problema:\n%s", d)
		}
	})

	// A ORDEM importa: o `[x]` sem arquivo vem primeiro porque é o dano mais caro — ele
	// declara pronto o que não está, e quem lê o board decide com base nisso.
	t.Run("PRHNP-B08: The ticked-but-absent finding is reported first", func(t *testing.T) {
		root, n := monta(t,
			"- [ ] `packages/shared/Feito.spec.md` — já entregue\n"+
				"- [x] `packages/shared/Nunca.spec.md` — nunca existiu\n",
			"packages/shared/Feito.spec.md")
		v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		iAusente, iFeito := strings.Index(d, "Nunca.spec.md"), strings.Index(d, "Feito.spec.md")
		if iAusente < 0 || iFeito < 0 {
			t.Fatalf("o laudo precisa citar os dois achados:\n%s", d)
		}
		if iAusente > iFeito {
			t.Errorf("o `[x]` sem arquivo tem de vir ANTES do retrabalho:\n%s", d)
		}
	})

	t.Run("PRHNP-B10: A progress that agrees with the disk on every item passes", func(t *testing.T) {
		root, n := monta(t,
			"- [x] `a.spec.md` — feito\n- [ ] `b.spec.md` — por fazer\n",
			"a.spec.md")
		if v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil)); v != Pass {
			t.Errorf("veredito %v, queria Pass: %s", v, d)
		}
	})

	// Item em PROSA não é confrontável, e reprovar por ele cobraria de quem escreveu um
	// progresso legítimo — a fase pode ser "revisar com o time".
	t.Run("PRHNP-B09: An item in prose citing no path is not charged", func(t *testing.T) {
		root, n := monta(t, "- [ ] revisar com o time antes de seguir\n- [x] alinhado na daily\n")
		if v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil)); v != Pass {
			t.Errorf("veredito %v, queria Pass — prosa não é confrontável: %s", v, d)
		}
	})

	// O `anchors new progress` escreve o molde quando a fase não semeia nada. Um `[ ]`
	// eterno faz o plano parecer incompleto para sempre: o `anchors next` volta a ele, e
	// quem lê não sabe se falta trabalho ou falta limpar o arquivo.
	t.Run("PRHNP-B07: A checkbox item promising no file at all is failed", func(t *testing.T) {
		root, n := monta(t, "- [ ] TODO: um item por spec que esta fase semeia\n- [ ] FIXME: a outra fase\n")
		v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		// Cada molde é acusado, não só o primeiro: o segundo `[ ]` eterno engana do mesmo jeito.
		for _, molde := range []string{"TODO", "FIXME"} {
			if !strings.Contains(d, molde) {
				t.Errorf("o laudo não mostra o item molde %q:\n%s", molde, d)
			}
		}
	})

	// AUSÊNCIA do progresso não é falha DESTE gate: o plano pode ter nascido antes do
	// mecanismo. São duas perguntas diferentes — "existe?" e "é verdade?".
	t.Run("PRHNP-B02: An untracked plan is skipped, not failed", func(t *testing.T) {
		root, n := monta(t, "")
		if v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil)); v != Skip {
			t.Fatalf("veredito %v, queria Skip: %s", v, d)
		}
	})

	t.Run("PRHNP-B03: The skip for an untracked plan says how to track it", func(t *testing.T) {
		root, n := monta(t, "")
		v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Skip {
			t.Fatalf("veredito %v, queria Skip: %s", v, d)
		}
		if !strings.Contains(d, "anchors progress add demand") {
			t.Errorf("o Skip não diz como corrigir:\n%s", d)
		}
	})

	t.Run("PRHNP-X01: The gate does not charge that the plan is tracked", func(t *testing.T) {
		root, n := monta(t, "")
		if v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil)); v == Fail {
			t.Errorf("a ausência do companheiro é outra pergunta — não pode reprovar aqui: %s", d)
		}
	})

	// O que o item DIZ não é conferido: a régua é o disco, que é determinístico e
	// dispensa opinião.
	t.Run("PRHNP-X02: The gate does not judge the content of an item beyond the path", func(t *testing.T) {
		root, n := monta(t,
			"- [x] `a.spec.md` — a spec do relatorio mensal de faturamento\n",
			"a.spec.md")
		if v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil)); v != Pass {
			t.Errorf("veredito %v, queria Pass — a descrição não é confrontada: %s", v, d)
		}
	})

	// Uma direção sem achado não entra no laudo: um "0 item(s)" lê como defeito que não há.
	t.Run("PRHNP-B12: The verdict carries only the directions that found something", func(t *testing.T) {
		vazio := func(chave string) string { return i18n.T(chave, 0, "") }

		root, n := monta(t,
			"- [ ] `packages/shared/Feito.spec.md` — entregue e não marcado\n",
			"packages/shared/Feito.spec.md")
		v, d := checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		for _, chave := range []string{
			"gate.progress_honest.marked_but_missing",
			"gate.progress_honest.seeds_missing",
			"gate.progress_honest.placeholder_items",
		} {
			if strings.Contains(d, vazio(chave)) {
				t.Errorf("a direção %s não achou nada e entrou no laudo:\n%s", chave, d)
			}
		}

		root, n = monta(t, "- [x] `packages/shared/Nunca.spec.md` — nunca existiu\n")
		v, d = checkProgressHonest("", n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if strings.Contains(d, vazio("gate.progress_honest.open_but_exists")) {
			t.Errorf("a direção do retrabalho não achou nada e entrou no laudo:\n%s", d)
		}
	})

	t.Run("PRHNP-B01: An artifact that is not a plan leaves without a verdict", func(t *testing.T) {
		root := t.TempDir()
		n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
		if v, _ := checkProgressHonest("", n, root, nil, testkit.WithComments(nil)); v != Skip {
			t.Errorf("veredito %v, queria Skip", v)
		}
	})
}

// The progress file's name has ONE definition, the scanner's — it is what keeps the file out
// of the map. Two constants would diverge in silence, and the gate would read a file nobody
// writes.
func TestProgressHonest_caminhoDerivaDoScan(t *testing.T) {
	t.Run("PRHNP-I01: The progress file's name has one definition, derived from the scanner", func(t *testing.T) {})
	if progress.FileName != scan.ProgressFile || progress.Path("r") != filepath.Join("r", scan.ProgressFile) {
		t.Errorf("the gate's file %q diverged from the scanner's %q", progress.FileName, scan.ProgressFile)
	}
}

// The progress file stays OUT of the map on purpose: it exists to CHANGE, and a gate that
// reached it as a node would charge a justification for every change of the artifact that
// records work.
func TestProgressHonest_companheiroForaDoMapa(t *testing.T) {
	t.Run("PRHNP-X03: The gate does not put the progress file into the map", func(t *testing.T) {})
	if !scan.IsProgressFile(progress.FileName) {
		t.Fatalf("%q must be recognized by the scan — that is how it stays out of the map", progress.FileName)
	}
	if scan.IsProgressFile("plans/0002-plataforma.md") {
		t.Errorf("the PLAN is no progress file: it stays a node, or the gate loses its anchor")
	}
}

// A TERCEIRA DIREÇÃO: semente que o plano promete e o progresso não lista.
//
// As duas primeiras conferem os itens QUE EXISTEM no progresso. Um gate que só olha para
// dentro do arquivo nunca vê o que falta nele — e foi assim que este gate passou com
// `✓15` num plano que acabara de ganhar uma spec semeada (o `MutualTls.spec.md`, migrado
// do 0016 por `PLTFR-R0004`).
//
// O efeito é o mesmo do `[x]` mentiroso: o progresso diz que a fase acabou, e o
// `anchors next` seguiu para outro plano com trabalho declarado por fazer.
func TestProgressHonest_sementeForaDoProgresso(t *testing.T) {
	monta := func(t *testing.T, plano, progresso string) (string, mapx.Node, string) {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
			t.Fatal(err)
		}
		rel := "plans/0002-plataforma.md"
		if err := os.WriteFile(filepath.Join(root, rel), []byte(plano), 0o644); err != nil {
			t.Fatal(err)
		}
		trackPlan(t, root, rel, progresso)
		return root, mapx.Node{ID: rel, Kind: mapx.KindPlan}, plano
	}

	t.Run("PRHNP-B06: A spec the plan seeds and the progress does not list is failed", func(t *testing.T) {
		// Várias sementes e vários itens: a ausente vem DEPOIS das listadas no plano, e o
		// progresso lista mais de uma — cada semente e cada item contam, não só o primeiro.
		root, n, conteudo := monta(t,
			"- [ ] `packages/infra/DataStore.spec.md` — a tabela\n"+
				"- [ ] `packages/infra/Cache.spec.md` — o cache\n"+
				"- [ ] `packages/infra/MutualTls.spec.md` — o canal\n",
			"- [ ] `packages/infra/DataStore.spec.md` — a tabela\n"+
				"- [ ] `packages/infra/Cache.spec.md` — o cache\n")
		v, d := checkProgressHonest(conteudo, n, root, nil, testkit.WithComments(nil))
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if !strings.Contains(d, "MutualTls.spec.md") {
			t.Errorf("o laudo não nomeia a spec ausente:\n%s", d)
		}
		for _, listada := range []string{"DataStore", "Cache"} {
			if strings.Contains(d, listada) {
				t.Errorf("acusou %s, uma spec que ESTÁ no progresso:\n%s", listada, d)
			}
		}
	})

	// O FALSO POSITIVO que apareceu no app de referência: a revisão `PLTFR-R0003` menciona
	// `` `DataStore.spec.md` `` em prosa, SEM caminho. O regex antigo casava qualquer
	// `x.spec.md` do arquivo, e o gate acusou o progresso de não listar uma spec que ele
	// lista — três dos quatro achados eram menções assim.
	//
	// A âncora é o `- [ ]` no início da linha: o plano PROMETE criar naquele item, e a
	// prosa das revisões fala sobre o que já existe.
	t.Run("PRHNP-I03: A spec mentioned in the plan's prose is not a seed", func(t *testing.T) {
		root, n, conteudo := monta(t,
			"> **PLTFR-R0003:** o `DataStore.spec.md` deixa de guardar estado de incidente.\n\n"+
				"- [ ] `packages/infra/DataStore.spec.md` — a tabela\n",
			"- [ ] `packages/infra/DataStore.spec.md` — a tabela\n")
		if v, d := checkProgressHonest(conteudo, n, root, nil, testkit.WithComments(nil)); v != Pass {
			t.Errorf("veredito %v, queria Pass — a menção em prosa não é promessa: %s", v, d)
		}
	})

	// O item pode ter sido reescrito (encurtado, movido de fase) sem deixar de ser o
	// mesmo item — o que identifica é o CAMINHO, não a linha.
	t.Run("PRHNP-I02: A seed is matched by path, never by the item's text", func(t *testing.T) {
		root, n, conteudo := monta(t,
			"- [ ] `packages/infra/DataStore.spec.md` — a tabela DynamoDB e o que vive nela\n",
			"- [x] `packages/infra/DataStore.spec.md`\n")
		if v, d := checkProgressHonest(conteudo, n, root, nil, testkit.WithComments(nil)); v == Fail &&
			(strings.Contains(d, "NÃO lista") || strings.Contains(d, "DOES NOT list") ||
				strings.Contains(d, "NO lista")) {
			t.Errorf("acusou por descrição diferente — o caminho é o que identifica:\n%s", d)
		}
	})

	t.Run("PRHNP-I04: A template file is never a seeded spec", func(t *testing.T) {
		root, n, conteudo := monta(t,
			"- [ ] `packages/_TEMPLATE_Area.spec.md` — o gabarito\n",
			"# Progresso — PLTFR\n")
		if v, d := checkProgressHonest(conteudo, n, root, nil, testkit.WithComments(nil)); v != Pass {
			t.Errorf("veredito %v, queria Pass — molde não é spec semeada: %s", v, d)
		}
	})
}

// writeProgress lays a plan and its companion out in a temporary directory.
func writeProgress(t *testing.T, plan, progress string) (root, planPath string) {
	t.Helper()
	root = t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	planPath = "plans/0017-mutacao.md"
	if err := os.WriteFile(filepath.Join(root, planPath), []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}
	trackPlan(t, root, planPath, progress)
	return root, planPath
}

// THE PLACEHOLDER, measured in the reference app.
//
// The progress of plan 0017 had, in phase F02:
//
//   - [ ] TODO: um item por spec que esta fase semeia
//
// `anchors new progress` writes it when the phase seeds nothing, and it should go once
// someone decides what the phase does. It stayed.
//
// The three directions of this gate do not see it: the first two confront items that CITE
// A PATH (this one does not), and the third looks at the plan's seeds (this phase seeds
// nothing).
//
// And the effect is the opposite of what the gate protects: an eternal `[ ]` makes the plan
// look unfinished forever — `anchors next` comes back to it, and whoever reads does not
// know whether work is missing or the file needs cleaning.
func TestProgressHonest_accusesTheTODOPlaceholder(t *testing.T) {
	t.Run("PRHNP-B07: A checkbox item promising no file at all is failed", func(t *testing.T) {})
	plan := "<!-- @anchors\ncode: MTUAO\n-->\n# Plano 0017\n\n## Fases\n\n### MTUAO-W02 — o CI\n\nEsta fase decide quando a mutação roda.\n"
	progress := "# Progresso — MTUAO\n\n## MTUAO-W02 — o CI\n\n- [ ] TODO: um item por spec que esta fase semeia\n"

	root, planPath := writeProgress(t, plan, progress)
	n := mapx.Node{ID: planPath, Kind: mapx.KindPlan, Code: "MTUAO"}

	v, msg := checkProgressHonest(plan, n, root, nil, testkit.WithComments(nil))

	if v != Fail {
		t.Fatalf("the placeholder passed: %v — %s", v, msg)
	}
	if !strings.Contains(strings.ToUpper(msg), "TODO") {
		t.Errorf("the message does not name the TODO:\n%s", msg)
	}
}

// An item that cites an existing path still passes: the fourth direction must not accuse
// the normal case.
func TestProgressHonest_doesNotAccuseAnItemWithAValidPath(t *testing.T) {
	t.Run("PRHNP-B10: A progress that agrees with the disk on every item passes", func(t *testing.T) {})
	plan := "<!-- @anchors\ncode: MTUAO\n-->\n# Plano\n\n- [ ] `packages/shared/X.spec.md` — a spec\n"
	progress := "# Progresso\n\n- [x] `packages/shared/X.spec.md` — a spec\n"

	root, planPath := writeProgress(t, plan, progress)
	if err := os.MkdirAll(filepath.Join(root, "packages/shared"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "packages/shared/X.spec.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	n := mapx.Node{ID: planPath, Kind: mapx.KindPlan, Code: "MTUAO"}

	if v, msg := checkProgressHonest(plan, n, root, nil, testkit.WithComments(nil)); v != Pass {
		t.Errorf("the normal case was accused: %v — %s", v, msg)
	}
}

// And PROSE is not an item: a progress file may have explanatory text with the word TODO,
// and a line that is not `- [ ]` is not a promise of work.
func TestProgressHonest_ignoresProseWithTODO(t *testing.T) {
	t.Run("PRHNP-B09: An item in prose citing no path is not charged", func(t *testing.T) {})
	plan := "<!-- @anchors\ncode: MTUAO\n-->\n# Plano\n"
	progress := "# Progresso\n\nEsta fase ainda tem TODO a decidir, e isso está no plano.\n"

	root, planPath := writeProgress(t, plan, progress)
	n := mapx.Node{ID: planPath, Kind: mapx.KindPlan, Code: "MTUAO"}

	if v, msg := checkProgressHonest(plan, n, root, nil, testkit.WithComments(nil)); v == Fail {
		t.Errorf("prose with TODO was accused as an item: %s", msg)
	}
}

func TestProgressHonest_unreadableIsPending(t *testing.T) {
	t.Run("PRHNP-E01: An unreadable progress file is pending", func(t *testing.T) {})
	root := t.TempDir()
	trackPlan(t, root, "plans/0001.md", "- [ ] `a.spec.md`\n")
	b, _ := os.ReadFile(progress.Path(root))
	if err := os.WriteFile(progress.Path(root), append(b, []byte("hand edit\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, d := checkProgressHonest("", mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}, root, nil, testkit.WithComments(nil)); v != Pending || !strings.Contains(d, "edited by hand") {
		t.Errorf("%v %s", v, d)
	}
}

func TestProgressIntact(t *testing.T) {
	t.Run("PRHNP-B13: The progress file is the one Anchors wrote", func(t *testing.T) {})
	root := t.TempDir()
	if v, _ := checkProgressIntact("", mapx.Node{}, root, nil, testkit.WithComments(nil)); v != Skip {
		t.Errorf("no file: %v", v)
	}
	trackPlan(t, root, "plans/0001.md", "- [ ] `a.spec.md`\n")
	if v, d := checkProgressIntact("", mapx.Node{}, root, nil, testkit.WithComments(nil)); v != Pass {
		t.Errorf("written by anchors: %v %s", v, d)
	}
	b, _ := os.ReadFile(progress.Path(root))
	_ = os.WriteFile(progress.Path(root), []byte(strings.Replace(string(b), "plan", "PLAN", 1)), 0o644)
	if v, d := checkProgressIntact("", mapx.Node{}, root, nil, testkit.WithComments(nil)); v != Fail || !strings.Contains(d, "edited by hand") {
		t.Errorf("hand edit: %v %s", v, d)
	}
}
