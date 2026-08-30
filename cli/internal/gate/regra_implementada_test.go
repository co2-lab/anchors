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

// O ALVO entre colchetes nomeia a quem a dispensa se aplica.
//
// Sem ele, o marcador vale para a regra da linha — e isso basta quando cada regra está na
// sua. Mas uma spec pode ter seis regras marcadas no código e duas que não têm onde ser
// marcadas: declarar as duas numa linha só, nomeando-as, é mais legível que espalhar o
// marcador por linhas que não falam disso.
func TestNoMarkComAlvoNomeado(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.ts"), []byte("// MTVRX-B01: resolve\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	base := "<!-- @anchors\ncode: MTVRX\n-->\n\n## Regras\n\n" +
		"| Regra | Efeito |\n| --- | --- |\n" +
		"| `MTVRX-B01` | resolve a versão |\n" +
		"| `MTVRX-X01` | NÃO faz I/O |\n" +
		"| `MTVRX-X02` | NÃO toca no relógio |\n"
	n := mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}

	// Sem dispensa, as duas reprovam.
	if v, _ := checkRegraImplementada(base, n, root, nil, nil); v != Fail {
		t.Fatal("preparo: X01 e X02 sem código deveriam reprovar")
	}

	// UMA linha nomeando AS DUAS fecha a conta das duas.
	agrupada := base + "\n> `@no-mark:[MTVRX-X01, MTVRX-X02]` satisfeitas pela AUSÊNCIA de código\n"
	if v, msg := checkRegraImplementada(agrupada, n, root, nil, nil); v != Pass {
		t.Errorf("a declaração agrupada deveria dispensar as duas; veio %v (%s)", v, msg)
	}

	// E o alvo MANDA: nomear X01 numa linha que fala de X02 dispensa X01, não X02.
	soUma := base + "\n> `@no-mark:[MTVRX-X01]` só esta\n"
	v, msg := checkRegraImplementada(soUma, n, root, nil, nil)
	if v != Fail {
		t.Fatalf("X02 continua sem dispensa e deve reprovar; veio %v", v)
	}
	if !strings.Contains(msg, "MTVRX-X02") {
		t.Errorf("a mensagem deveria acusar X02, não X01: %s", msg)
	}
	if strings.Contains(msg, "MTVRX-X01") {
		t.Errorf("X01 foi dispensada e não pode aparecer no achado: %s", msg)
	}
}

// `@no-code` continua valendo: quem já escreveu não pode ver a spec quebrar por uma
// renomeação. O nome mudou porque o antigo mente — o código EXISTE (um `tsconfig.json`
// decide como tudo compila); o que não existe é a MARCAÇÃO.
func TestNomeAntigoNoCodeContinuaValendo(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.ts"), []byte("// MTVRX-B01: resolve\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := "<!-- @anchors\ncode: MTVRX\n-->\n\n## Regras\n\n" +
		"| Regra | Efeito |\n| --- | --- |\n" +
		"| `MTVRX-B01` | resolve a versão |\n" +
		"| `MTVRX-X01` | NÃO faz I/O @no-code: satisfeita pela ausência |\n"
	if v, msg := checkRegraImplementada(spec, mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}, root, nil, nil); v != Pass {
		t.Errorf("`@no-code` deveria continuar dispensando; veio %v (%s)", v, msg)
	}
}

// SEM alvo, o marcador vale para a spec INTEIRA — equivale a `[all]`.
//
// É a saída de quem tem uma unidade em que NADA se marca: nomear todas as regras seria
// repetir a mesma razão N vezes, e repetição em declaração é onde a divergência começa.
// A forma é mais forte que a nomeada, e por isso a razão importa mais: ela dispensa o
// gate inteiro para aquele arquivo.
func TestNoMarkSemAlvoValeParaTodas(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "u.ts"), []byte("nada marcado aqui\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := "<!-- @anchors\ncode: MTVRX\n-->\n\n## Regras\n\n" +
		"| Regra | Efeito |\n| --- | --- |\n" +
		"| `MTVRX-B01` | resolve a versão |\n" +
		"| `MTVRX-X01` | NÃO faz I/O |\n" +
		"\n> `@no-mark:` esta unidade descreve configuração, e nada nela recebe marcação\n"
	n := mapx.Node{Kind: mapx.KindSpec, ID: "u.spec.md"}
	if v, msg := checkRegraImplementada(spec, n, root, nil, nil); v != Pass {
		t.Errorf("sem alvo deveria dispensar todas; veio %v (%s)", v, msg)
	}
	// `[all]` explícito é a mesma coisa, escrita para quem prefere ver a intenção.
	explicito := strings.Replace(spec, "@no-mark:", "@no-mark:[all]", 1)
	if v, msg := checkRegraImplementada(explicito, n, root, nil, nil); v != Pass {
		t.Errorf("`[all]` explícito deveria dispensar todas; veio %v (%s)", v, msg)
	}
}
