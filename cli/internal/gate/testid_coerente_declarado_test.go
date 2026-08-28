package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cfgHandle devolve uma config que DECLARA o atributo de ancoragem — sem ela os
// gates de inventário pulam, que é o comportamento opt-in.
func cfgHandle(attr string) *config.Config {
	return &config.Config{Derived: &config.Derived{Anchor: "code", TestHandle: attr}}
}

// inventarioFixture: uma spec ligada (specifies) a uma unidade cujo conteúdo o teste
// controla. `specBody` é a spec inteira, com ou sem a seção de inventário.
func inventarioFixture(t *testing.T, unitSrc string) (mapx.Node, *mapx.Graph, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(unitSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies}},
	}
	return spec, g, root
}

const secaoOK = "## Superfície de Teste\n\n| id | papel |\n| -- | ----- |\n| `:abcd-screen` | raiz |\n"

func TestTestIDDeclared_inventarioCompletoPassa(t *testing.T) {
	n, g, root := inventarioFixture(t, `<View testID=":abcd-screen" />`)
	if v, msg := checkTestIDCoerente(secaoOK, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("inventário batendo com o código deveria passar: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_expostoSemDeclararReprova(t *testing.T) {
	// Sentido código → spec: superfície não-contratada. É por onde a divergência de
	// identidade entra sem ninguém ver.
	n, g, root := inventarioFixture(t, `<View testID=":abcd-screen" /><View testID=":abcd-oculto" />`)
	v, msg := checkTestIDCoerente(secaoOK, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("exposto sem declarar deveria reprovar: %v", v)
	}
	if !strings.Contains(msg, "abcd-oculto") {
		t.Errorf("a mensagem deve nomear o id não declarado: %s", msg)
	}
}

func TestTestIDDeclared_declaradoSemExporReprova(t *testing.T) {
	// Sentido spec → código: contrato morto. Quem for escrever o teste procura o que
	// não existe. É a metade que um gate de sentido único não pega.
	n, g, root := inventarioFixture(t, `<View testID=":abcd-screen" />`)
	spec := secaoOK + "| `:abcd-fantasma` | sumiu no refactor |\n"
	v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("declarado sem expor deveria reprovar: %v", v)
	}
	if !strings.Contains(msg, "abcd-fantasma") {
		t.Errorf("a mensagem deve nomear o id inexistente: %s", msg)
	}
}

func TestTestIDDeclared_semHandleDeclaradoPula(t *testing.T) {
	// O ponto do opt-in: sem `derived.test_handle` o gate não tem atributo a procurar.
	// Inferir `testID` por default faria o gate reportar VERDE sobre o que não
	// conferiu — num projeto Android, em qualquer backend.
	n, g, root := inventarioFixture(t, `<View testID=":abcd-screen" />`)
	if v, _ := checkTestIDCoerente("", n, root, g, &config.Config{}); v != Skip {
		t.Errorf("sem test_handle o gate deve pular: %v", v)
	}
}

func TestTestIDDeclared_atributoDeOutroEcossistema(t *testing.T) {
	// O mesmo gate, num projeto web: o atributo é `data-testid`. Se o regex estivesse
	// cravado em `testID`, este caso passaria vazio (falso verde).
	n, g, root := inventarioFixture(t, `<div data-testid="abcd-root" /><div data-testid="abcd-item" />`)
	spec := "## Superfície de Teste\n\n- `abcd-root`\n"
	v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("data-testid"))
	if v != Fail || !strings.Contains(msg, "abcd-item") {
		t.Errorf("deveria acusar o id web não declarado: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_templateDeclaraPrefixo(t *testing.T) {
	// O sufixo é DADO de runtime (`${id}`). Exigir que a spec declare `abcd-item-42`
	// seria exigir que ela declare o dado; declara-se o prefixo com a marca `-*`.
	n, g, root := inventarioFixture(t, "<View testID={`:abcd-item-${id}`} />")
	spec := "## Superfície de Teste\n\n- `:abcd-item-*`\n"
	if v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("prefixo dinâmico declarado com `-*` deveria passar: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_unidadeSemHandleNaoEhOfensa(t *testing.T) {
	// Nem todo componente tem superfície de teste. Cobrar inventário de quem não
	// marca elemento algum viraria ruído sobre toda unidade de apresentação pura.
	n, g, root := inventarioFixture(t, `<View />`)
	if v, _ := checkTestIDCoerente("", n, root, g, cfgHandle("testID")); v != Skip {
		t.Errorf("unidade que não expõe handle não tem contrato a declarar: %v", v)
	}
}

func TestTestIDDeclared_tituloComQualificador(t *testing.T) {
	// `## Test IDs (Maestro)` — 42 das 159 specs do app de referência nomeiam no título a superfície
	// que consome os ids. Exigir fim-de-linha logo após o título fazia o gate não ver
	// a seção e acusar de "não declara nada" justamente as specs mais bem documentadas.
	n, g, root := inventarioFixture(t, `<View testID=":abcd-screen" />`)
	spec := "## Test IDs (Maestro)\n\n| testID | Elemento |\n| -- | -- |\n| `abcd-screen` | raiz |\n"
	if v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("título com qualificador é a mesma seção: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_craseForaDaSecaoNaoConta(t *testing.T) {
	// A seção delimita o inventário. Sem isso, qualquer crase na prosa (classe CSS,
	// nome de campo) contaria como declaração — e o gate aprovaria por acidente, que
	// é pior que reprovar por engano.
	n, g, root := inventarioFixture(t, `<View testID=":abcd-screen" />`)
	spec := "## Notas\n\nA raiz usa `:abcd-screen` e a classe `text-mute`.\n"
	v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("menção em prosa não é inventário: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_soAPrimeiraColunaEhOID(t *testing.T) {
	// A tabela de inventário do app de referência tem colunas "Elemento" e "Usado em (flow)", e
	// ambas levam crase: `Raiz do TouchableOpacity`, `ATLNX-VR`. Ler a linha inteira
	// fazia o gate colher essas células como se fossem testIDs declarados — e o gate
	// irmão depois as acusava de órfãs, inventando dívida a partir da documentação.
	n, g, root := inventarioFixture(t, `<View testID=":abcd-screen" />`)
	spec := "## Test IDs (Maestro)\n\n" +
		"| testID | Elemento | Usado em |\n| -- | -- | -- |\n" +
		"| `abcd-screen` | Raiz do `TouchableOpacity` | `ABCDX-VR` |\n"
	v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID"))
	if v != Pass {
		t.Errorf("só a 1ª célula é o id; as outras são descrição: %v (%s)", v, msg)
	}
	if ids := testIDsDeclarados(spec, "testID"); len(ids) != 1 || ids[0] != "abcd-screen" {
		t.Errorf("inventário deveria ter só o id da 1ª coluna, veio: %v", ids)
	}
}

func TestTestIDDeclared_propDerivadaEhHandle(t *testing.T) {
	// `backTestID=":spending-button-back"` e `confirmTestID: ':x'` (propriedade de
	// objeto): o valor termina num handle real no filho. Ignorá-las fazia o gate
	// declarar inexistente um id que a spec documenta e que os flows usam.
	n, g, root := inventarioFixture(t,
		`<Header backTestID=":abcd-back" /><Alert buttons={[{ confirmTestID: ':abcd-ok' }]} />`)
	spec := "## Test IDs\n\n- `abcd-back`\n- `abcd-ok`\n"
	if v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("prop derivada carrega handle: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_prefixoMontadoNoFilho(t *testing.T) {
	// `testIDPrefix=":otp-input"` vira `otp-input-0`…`otp-input-5` no filho. O literal
	// é a CABEÇA, não um id: registrá-lo como está cobraria da spec um `otp-input` que
	// nunca aparece, e acusaria de inexistentes os `otp-input-N` que 6 flows usam.
	n, g, root := inventarioFixture(t, `<Otp testIDPrefix=":otp-input" />`)
	spec := "## Test IDs\n\n- `otp-input-0`\n- `otp-input-5`\n"
	if v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("prefixo montado cobre os ids concretos: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_curingaCasaNosDoisSentidos(t *testing.T) {
	// Exposto genérico cobre declarado concreto, e vice-versa — a forma genérica tanto
	// nasce no código (template) quanto na spec.
	if !cobre("abcd-item-*", "abcd-item-3") {
		t.Error("curinga exposto deveria cobrir o id concreto")
	}
	if cobre("abcd-item-3", "abcd-item-*") {
		t.Error("id concreto NÃO cobre o curinga (só o inverso)")
	}
	if !cobre("abcd-x", "abcd-x") || cobre("abcd-x", "abcd-y") {
		t.Error("igualdade exata deveria valer, e só ela")
	}
}

func TestTestIDDeclared_ternarioIgnoraACondicao(t *testing.T) {
	// `testID={c2.k === 'push' ? firstToggleTestID : undefined}` — o literal da
	// CONDIÇÃO é estado de domínio, não handle. Colhê-lo registraria `push` como
	// exposto e o gate cobraria da spec um id que não existe.
	n, g, root := inventarioFixture(t,
		`<V testID={k === 'push' ? ':abcd-toggle' : undefined} /><V testID={i === 0 ? ':abcd-first' : undefined} />`)
	spec := "## Test IDs\n\n- `abcd-toggle`\n- `abcd-first`\n"
	if v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("condição não é handle; os dois ramos sim: %v (%s)", v, msg)
	}
}

func TestTestIDDeclared_nomeDoAtributoNaoEhID(t *testing.T) {
	// O átomo genérico (Button, Avatar) RECEBE o handle de fora, e a spec o documenta
	// como `testID` (prop). Declarar isso é dizer "aceito um handle", não "exponho
	// este" — cobrá-lo exigiria do átomo um id fixo, o oposto de ser reusável.
	n, g, root := inventarioFixture(t, `<Touchable testID={testID} />`)
	spec := "## Test IDs\n\n| testID | Elemento |\n| -- | -- |\n| `testID` (prop) | Raiz |\n"
	if v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID")); v != Skip && v != Pass {
		t.Errorf("o nome do atributo não é um id declarado: %v (%s)", v, msg)
	}
}
