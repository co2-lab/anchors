package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
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

	t.Run("PRHNP-B04: A ticked item whose file does not exist is failed", func(t *testing.T) {
		root, n := monta(t, "- [x] `packages/shared/Nunca.spec.md` — nunca existiu\n")
		v, d := checkProgressHonest("", n, root, nil, nil)
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
		v, d := checkProgressHonest("", n, root, nil, nil)
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
		v, d := checkProgressHonest("", n, root, nil, nil)
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
		v, d := checkProgressHonest("", n, root, nil, nil)
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
		if v, d := checkProgressHonest("", n, root, nil, nil); v != Pass {
			t.Errorf("veredito %v, queria Pass: %s", v, d)
		}
	})

	// Item em PROSA não é confrontável, e reprovar por ele cobraria de quem escreveu um
	// progresso legítimo — a fase pode ser "revisar com o time".
	t.Run("PRHNP-B09: An item in prose citing no path is not charged", func(t *testing.T) {
		root, n := monta(t, "- [ ] revisar com o time antes de seguir\n- [x] alinhado na daily\n")
		if v, d := checkProgressHonest("", n, root, nil, nil); v != Pass {
			t.Errorf("veredito %v, queria Pass — prosa não é confrontável: %s", v, d)
		}
	})

	// O `anchors new progress` escreve o molde quando a fase não semeia nada. Um `[ ]`
	// eterno faz o plano parecer incompleto para sempre: o `anchors next` volta a ele, e
	// quem lê não sabe se falta trabalho ou falta limpar o arquivo.
	t.Run("PRHNP-B07: A checkbox item promising no file at all is failed", func(t *testing.T) {
		root, n := monta(t, "- [ ] TODO: um item por spec que esta fase semeia\n")
		v, d := checkProgressHonest("", n, root, nil, nil)
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if !strings.Contains(d, "TODO") {
			t.Errorf("o laudo não mostra o item molde:\n%s", d)
		}
	})

	// AUSÊNCIA do progresso não é falha DESTE gate: o plano pode ter nascido antes do
	// mecanismo. São duas perguntas diferentes — "existe?" e "é verdade?".
	t.Run("PRHNP-B02: A plan with no companion progress file is skipped, not failed", func(t *testing.T) {
		root, n := monta(t, "")
		if v, d := checkProgressHonest("", n, root, nil, nil); v != Skip {
			t.Fatalf("veredito %v, queria Skip: %s", v, d)
		}
	})

	t.Run("PRHNP-B03: The skip for a missing companion says how to create it", func(t *testing.T) {
		root, n := monta(t, "")
		v, d := checkProgressHonest("", n, root, nil, nil)
		if v != Skip {
			t.Fatalf("veredito %v, queria Skip: %s", v, d)
		}
		if !strings.Contains(d, "anchors new progress --for") {
			t.Errorf("o Skip não diz como corrigir:\n%s", d)
		}
	})

	t.Run("PRHNP-X01: The gate does not charge the existence of the progress file", func(t *testing.T) {
		root, n := monta(t, "")
		if v, d := checkProgressHonest("", n, root, nil, nil); v == Fail {
			t.Errorf("a ausência do companheiro é outra pergunta — não pode reprovar aqui: %s", d)
		}
	})

	// O que o item DIZ não é conferido: a régua é o disco, que é determinístico e
	// dispensa opinião.
	t.Run("PRHNP-X02: The gate does not judge the content of an item beyond the path", func(t *testing.T) {
		root, n := monta(t,
			"- [x] `a.spec.md` — a spec do relatorio mensal de faturamento\n",
			"a.spec.md")
		if v, d := checkProgressHonest("", n, root, nil, nil); v != Pass {
			t.Errorf("veredito %v, queria Pass — a descrição não é confrontada: %s", v, d)
		}
	})

	t.Run("PRHNP-B01: An artifact that is not a plan leaves without a verdict", func(t *testing.T) {
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
	t.Run("PRHNP-I01: The companion's path has one definition, derived from the scanner", func(t *testing.T) {})
	casos := map[string]string{
		"plans/0002-plataforma.md": "plans/0002-plataforma-progress.md",
		"plans/0017-mutacao.md":    "plans/0017-mutacao-progress.md",
		"sem-extensao":             "sem-extensao-progress.md",
	}
	for plano, quer := range casos {
		if got := progressPathOf(plano); got != quer {
			t.Errorf("progressPathOf(%q) = %q, queria %q", plano, got, quer)
		}
		// E a resposta é a MESMA que a do scan — que é quem detém a definição.
		if got, doScan := progressPathOf(plano), scan.ProgressPathFor(plano); got != doScan {
			t.Errorf("o gate divergiu do scan em %q: %q vs %q", plano, got, doScan)
		}
	}
}

// O companheiro fica FORA do mapa de propósito: ele existe para MUDAR, e um gate que o
// alcançasse como nó cobraria justificativa de cada edição do artefato que registra
// trabalho.
func TestProgressHonest_companheiroForaDoMapa(t *testing.T) {
	t.Run("PRHNP-X03: The gate does not put the progress file into the map", func(t *testing.T) {})
	prog := progressPathOf("plans/0002-plataforma.md")
	if !scan.IsProgressFile(prog) {
		t.Fatalf("%q devia ser reconhecido como companheiro pelo scan — é assim que ele fica fora do mapa", prog)
	}
	if scan.IsProgressFile("plans/0002-plataforma.md") {
		t.Errorf("o PLANO não é companheiro: ele precisa continuar sendo nó, ou o gate perde a âncora")
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
		prog := filepath.Join(root, "plans", "0002-plataforma-progress.md")
		if err := os.WriteFile(prog, []byte(progresso), 0o644); err != nil {
			t.Fatal(err)
		}
		return root, mapx.Node{ID: rel, Kind: mapx.KindPlan}, plano
	}

	t.Run("PRHNP-B06: A spec the plan seeds and the progress does not list is failed", func(t *testing.T) {
		root, n, conteudo := monta(t,
			"- [ ] `packages/infra/MutualTls.spec.md` — o canal\n"+
				"- [ ] `packages/infra/DataStore.spec.md` — a tabela\n",
			"- [ ] `packages/infra/DataStore.spec.md` — a tabela\n")
		v, d := checkProgressHonest(conteudo, n, root, nil, nil)
		if v != Fail {
			t.Fatalf("veredito %v, queria Fail: %s", v, d)
		}
		if !strings.Contains(d, "MutualTls.spec.md") {
			t.Errorf("o laudo não nomeia a spec ausente:\n%s", d)
		}
		if strings.Contains(d, "DataStore") {
			t.Errorf("acusou uma spec que ESTÁ no progresso:\n%s", d)
		}
	})

	// O FALSO POSITIVO que apareceu no blue-eyes: a revisão `PLTFR-R0003` menciona
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
		if v, d := checkProgressHonest(conteudo, n, root, nil, nil); v != Pass {
			t.Errorf("veredito %v, queria Pass — a menção em prosa não é promessa: %s", v, d)
		}
	})

	// O item pode ter sido reescrito (encurtado, movido de fase) sem deixar de ser o
	// mesmo item — o que identifica é o CAMINHO, não a linha.
	t.Run("PRHNP-I02: A seed is matched by path, never by the item's text", func(t *testing.T) {
		root, n, conteudo := monta(t,
			"- [ ] `packages/infra/DataStore.spec.md` — a tabela DynamoDB e o que vive nela\n",
			"- [x] `packages/infra/DataStore.spec.md`\n")
		if v, d := checkProgressHonest(conteudo, n, root, nil, nil); v == Fail &&
			(strings.Contains(d, "NÃO lista") || strings.Contains(d, "DOES NOT list") ||
				strings.Contains(d, "NO lista")) {
			t.Errorf("acusou por descrição diferente — o caminho é o que identifica:\n%s", d)
		}
	})

	t.Run("PRHNP-I04: A template file is never a seeded spec", func(t *testing.T) {
		root, n, conteudo := monta(t,
			"- [ ] `packages/_TEMPLATE_Area.spec.md` — o gabarito\n",
			"# Progresso — PLTFR\n")
		if v, d := checkProgressHonest(conteudo, n, root, nil, nil); v != Pass {
			t.Errorf("veredito %v, queria Pass — molde não é spec semeada: %s", v, d)
		}
	})
}
