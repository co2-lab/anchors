package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// rodaRegraImpl monta o par spec+código em disco e roda o gate. O nome do arquivo de
// código deriva do da spec (é assim que o gate acha o alvo).
func rodaRegraImpl(t *testing.T, spec, codigo string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "u.ts"), []byte(codigo), 0o644))
	cabecalho := "<!-- @anchors\n  code: SBNKX\n-->\n"
	if strings.Contains(spec, "AAAAX-") {
		cabecalho = "<!-- @anchors\n  code: AAAAX\n-->\n"
	}
	return checkRegraImplementada(cabecalho+spec,
		mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}, root, nil, nil)
}

// TestSpecQueFalaSozinhaEhAcusada guarda o defeito que motivou o gate: uma spec de handler
// ganhou 5 regras novas (98 linhas de catálogo) e o handler correspondente não ganhou uma
// linha. Metade de uma feature era código morto declarado como pronto, e os 26 gates
// ficaram verdes — porque a spec existe, o código existe, e os dois se referenciam pelo
// header.
func TestSpecQueFalaSozinhaEhAcusada(t *testing.T) {
	// O caso REALX: a unidade já entrou na prática (o código marca `B01`), e a spec ganhou
	// regras novas que ninguém implementou. É diferente da dívida de migração — aqui há
	// declaração, e ela está incompleta.
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "h.ts"), []byte(
		"// @anchors\n//   ref: MRHMX\n// MRHMX-B01: marca o lançamento\nexport function marcar() { return 1 }\n"), 0o644))

	spec := `<!-- @anchors
  code: MRHMX
-->
## Regras
| Regra | Efeito |
| --- | --- |
| ` + "`MRHMX-B01`" + ` | marca o lançamento |
| ` + "`MRHMX-B04`" + ` | aciona a migração |
| ` + "`MRHMX-B05`" + ` | a ordem é migrar depois mudar o estado |
`
	v, msg := checkRegraImplementada(spec, mapx.Node{Kind: mapx.KindSpec, ID: "h.spec.md"}, root, nil, nil)
	if v != Fail {
		t.Fatalf("spec com regra que o código ignora deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "MRHMX-B04") {
		t.Errorf("a mensagem precisa nomear a regra órfã; veio: %s", msg)
	}
}

// A DISPENSA DECLARADA é o que troca heurística por confronto. Exigir todas as regras
// marcadas seria falso (restrição é satisfeita pela ausência de código); "ao menos uma"
// não diz nada sobre as outras quinze. Quem escreve a spec declara, regra a regra.
func TestDispensaDeclaradaFechaAConta(t *testing.T) {
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "u.ts"), []byte(
		"// MTVRX-B01: resolve a versão vigente\nexport const f = 1\n"), 0o644))

	semDispensa := `<!-- @anchors
  code: MTVRX
-->
| Regra | Efeito |
| --- | --- |
| ` + "`MTVRX-B01`" + ` | resolve a versão |
| ` + "`MTVRX-X01`" + ` | NÃO faz I/O |
`
	v, msg := checkRegraImplementada(semDispensa, mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}, root, nil, nil)
	if v != Fail {
		t.Fatalf("X01 sem código e sem dispensa deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "@no-code") {
		t.Errorf("a mensagem tem de ensinar a saída declarada; veio: %s", msg)
	}

	comDispensa := strings.Replace(semDispensa,
		"| `MTVRX-X01` | NÃO faz I/O |",
		"| `MTVRX-X01` | NÃO faz I/O @no-code: satisfeita pela ausência de import |", 1)
	if v, msg := checkRegraImplementada(comDispensa, mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}, root, nil, nil); v != Pass {
		t.Errorf("dispensa declarada COM razão fecha a conta; veio %v (%s)", v, msg)
	}

	// Marcador nu não dispensa nada — a razão escrita é o que torna a dispensa uma
	// prestação de contas em vez de um jeito de calar o gate.
	nu := strings.Replace(semDispensa, "| `MTVRX-X01` | NÃO faz I/O |", "| `MTVRX-X01` | NÃO faz I/O @no-code: |", 1)
	if v, _ := checkRegraImplementada(nu, mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}, root, nil, nil); v != Fail {
		t.Errorf("`@no-code:` sem razão não pode dispensar; veio %v", v)
	}
}

// A unidade que não declarou NADA é dívida de MIGRAÇÃO, não defeito: nasceu antes da
// prática. Medido no repositório de origem: 3.114 regras em 590 unidades. Acusá-las
// reprovaria 98% do projeto, e um gate assim é desligado no primeiro dia.
func TestUnidadeAnteriorAPraticaEhPendencia(t *testing.T) {
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "v.ts"), []byte("export const f = 1\n"), 0o644))
	spec := "<!-- @anchors\n  code: ANTGX\n-->\n| `ANTGX-B01` | faz algo |\n| `ANTGX-B02` | faz outro |\n"

	v, msg := checkRegraImplementada(spec, mapx.Node{Kind: mapx.KindSpec, ID: "v.spec.md"}, root, nil, nil)
	if v != Pending {
		t.Errorf("nenhuma regra declarada = dívida de migração (Pending); veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "anterior à prática") {
		t.Errorf("a mensagem tem de NOMEAR a dívida, não fingir aprovação; veio: %s", msg)
	}
}

// Sem código no disco quem acusa é o `trinca-completa`; duplicar a cobrança faria dois
// gates apontando o mesmo dedo.
func TestSemCodigoNaoEhAssunto(t *testing.T) {
	spec := "<!-- @anchors\n  code: NOVAX\n-->\n| `NOVAX-B01` | faz algo |\n"
	if v, _ := checkRegraImplementada(spec, mapx.Node{Kind: mapx.KindSpec, ID: "x.spec.md"}, t.TempDir(), nil, nil); v != Skip {
		t.Errorf("sem código, Skip; veio %v", v)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// A PENDÊNCIA DE MIGRAÇÃO TEM DE VENCER.
//
// Ela existe para não reprovar 98% de um projeto que nasceu antes da prática. Mas sem
// uma saída ela nunca acaba — e pendência que não vence é indistinguível de "ninguém
// olhou". Medido: uma spec descrevia "mantém o cartão selecionado" enquanto o código era
// um CRUD sem seleção; o gate VIU a regra ausente, caiu no ramo de migração e devolveu
// pendência. O defeito atravessou os 44 gates.
func TestRegraImplementada_marcacaoExigidaVenceAPendencia(t *testing.T) {
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "u.ts"), []byte(
		"export const useStore = () => ({ addAccount() {} })\n"), 0o644))
	spec := "<!-- @anchors\n  code: SBNKX\n-->\n| `SBNKX-B01` | selecionar cartão | atualiza o contexto |\n"
	n := mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}

	// Sem declaração: migração em curso, pendência (o comportamento de hoje).
	if v, _ := checkRegraImplementada(spec, n, root, nil, nil); v != Pending {
		t.Errorf("sem `rule_marking` a dívida é pendência: %v", v)
	}

	// Declarado: a migração acabou, e a unidade é cobrada como qualquer outra.
	cfg := &config.Config{Derived: &config.Derived{RuleMarking: "required"}}
	v, msg := checkRegraImplementada(spec, n, root, nil, cfg)
	if v != Fail {
		t.Fatalf("com `rule_marking: required` a pendência vira reprovação: %v", v)
	}
	if !strings.Contains(msg, "spec descreve mesmo esta unidade") {
		t.Errorf("a mensagem deve levantar a hipótese de spec errada: %s", msg)
	}
}

// Com a marcação exigida, quem JÁ marca segue passando — a exigência não pune quem
// está em dia.
func TestRegraImplementada_marcacaoExigidaNaoPuneQuemMarca(t *testing.T) {
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "u.ts"), []byte(
		"// SBNKX-B01: inclui a conta\nexport function addAccount() {}\n"), 0o644))
	spec := "<!-- @anchors\n  code: SBNKX\n-->\n| `SBNKX-B01` | incluir conta | acrescenta à lista |\n"
	cfg := &config.Config{Derived: &config.Derived{RuleMarking: "required"}}

	if v, msg := checkRegraImplementada(spec, mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"},
		root, nil, cfg); v != Pass {
		t.Errorf("quem marca a regra passa: %v (%s)", v, msg)
	}
}
