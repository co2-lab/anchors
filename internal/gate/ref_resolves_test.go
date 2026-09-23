package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaRefResolves(t *testing.T, arquivo, conteudo, specNome, specConteudo string) (Verdict, string) {
	t.Helper()
	dir := t.TempDir()
	if specNome != "" {
		if err := os.WriteFile(filepath.Join(dir, specNome), []byte(specConteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	kind := mapx.KindCode
	if strings.HasSuffix(arquivo, ".feature") {
		kind = mapx.KindFeature
	} else if strings.Contains(arquivo, ".test.") {
		kind = mapx.KindTest
	}
	return checkRefResolves(conteudo, mapx.Node{ID: arquivo, Kind: kind}, dir, nil, nil)
}

// Variant WITH a graph: the cases with no sibling spec need it to say whether the cited
// code exists anywhere in the project.
func runRefResolvesWithGraph(t *testing.T, arquivo, conteudo string, g *mapx.Graph) (Verdict, string) {
	t.Helper()
	dir := t.TempDir()
	kind := mapx.KindCode
	if strings.HasSuffix(arquivo, ".feature") {
		kind = mapx.KindFeature
	} else if strings.Contains(arquivo, ".test.") {
		kind = mapx.KindTest
	}
	return checkRefResolves(conteudo, mapx.Node{ID: arquivo, Kind: kind}, dir, g, nil)
}

// graphWith builds a map holding the given DECLARED identities.
func graphWith(codes ...string) *mapx.Graph {
	g := &mapx.Graph{}
	for _, c := range codes {
		g.Nodes = append(g.Nodes, mapx.Node{
			ID: c + ".spec.md", Kind: mapx.KindSpec, Code: c, CodeDeclarado: true,
		})
	}
	return g
}

// The hole the Skip left: with no sibling spec the gate went silent, and an INVENTED
// `ref:` passed as undetermined. Measured in MIF: 1344 of 3174 refs fell into Skip (42%),
// and a hand-made `ref: KYBDX` — a code that is no spec's `code:` — gave `~1`, not `✗1`.
// An infra file with no spec is legitimate; citing a phantom identity is not.
func TestRefResolvesCodeMissingFromTheProject(t *testing.T) {
	t.Run("RFRSR-B05: A reference to a code that exists nowhere fails", func(t *testing.T) {})
	v, d := runRefResolvesWithGraph(t,
		"kybDocs.ts", "// @anchors\n//   ref: KYBDX\n//   layer: infra\n",
		graphWith("ORAT1", "MORQX"))
	if v != Fail {
		t.Fatalf("a ref to a missing code should fail, got %s (%s)", v, d)
	}
	if !strings.Contains(d, "KYBDX") {
		t.Fatalf("the verdict should name the invented code, got %q", d)
	}
}

// The distinction from `triad-complete` holds: infra with no sibling spec is legitimate,
// as long as the `ref:` points at a REAL identity. Failing here would steal the other
// gate's finding.
func TestRefResolvesNoSiblingSpecButCodeExists(t *testing.T) {
	t.Run("RFRSR-B06: Without sibling spec, an existing code still skips", func(t *testing.T) {})
	v, _ := runRefResolvesWithGraph(t,
		"helper.ts", "// @anchors\n//   ref: MORQX\n//   layer: infra\n",
		graphWith("ORAT1", "MORQX"))
	if v != Skip {
		t.Fatalf("an existing code with no sibling spec should skip, got %s", v)
	}
}

// An INFERRED identity owns nothing: accepting `CodeDeclarado: false` would let through
// the ref that points at an example string in a fixture or in documentation.
func TestRefResolvesInferredIdentityDoesNotCount(t *testing.T) {
	t.Run("RFRSR-B07: An inferred identity does not satisfy the reference", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "fixture.md", Kind: mapx.KindDoc, Code: "FAKEX", CodeDeclarado: false},
	}}
	v, _ := runRefResolvesWithGraph(t, "x.ts", "// @anchors\n//   ref: FAKEX\n", g)
	if v != Fail {
		t.Fatalf("an inferred identity should not satisfy the ref, got %s", v)
	}
}

// Without a graph (a call that does not provide one) the gate cannot assert absence — and
// asserting what was not measured is worse than staying silent.
func TestRefResolvesWithoutGraphDoesNotAssertAbsence(t *testing.T) {
	t.Run("RFRSR-B08: Without a graph, absence is not asserted", func(t *testing.T) {})
	v, _ := runRefResolvesWithGraph(t, "x.ts", "// @anchors\n//   ref: QUALQ\n", nil)
	if v != Skip {
		t.Fatalf("without a graph it should skip, got %s", v)
	}
}

// O caso real: 49 arquivos de modelo com `ref: DTAXX` — a identidade de quando os modelos
// viviam num arquivo só. Depois da desfusão cada um ganhou spec própria, e nenhum `ref:`
// foi propagado. `header-conforme` ficou verde nos 49: ele confere que o campo EXISTE.
func TestRefResolvesRefactoracaoNaoPropagada(t *testing.T) {
	t.Run("RFRSR-B01: A reference that does not match the sibling spec fails", func(t *testing.T) {})
	v, d := rodaRefResolves(t,
		"TaxReceipt.ts", "// @anchors\n//   ref: DTAXX\n//   layer: schema-model\n",
		"TaxReceipt.spec.md", "<!-- @anchors\n  code: TRT1X\n-->\n")
	if v != Fail {
		t.Fatalf("ref divergente deveria reprovar, foi %s (%s)", v, d)
	}
}

// O veredito tem de mostrar OS DOIS lados: o que está escrito e o que a irmã declara.
// Sem isso quem lê abre dois arquivos para descobrir qual é qual.
func TestRefResolvesMensagemMostraOsDoisLados(t *testing.T) {
	t.Run("RFRSR-B02: The failing verdict names both sides of the divergence", func(t *testing.T) {})
	_, d := rodaRefResolves(t,
		"TaxReceipt.ts", "// @anchors\n//   ref: DTAXX\n",
		"TaxReceipt.spec.md", "<!-- @anchors\n  code: TRT1X\n-->\n")
	for _, esperado := range []string{"DTAXX", "TRT1X", "TaxReceipt.spec.md"} {
		if !strings.Contains(d, esperado) {
			t.Errorf("a mensagem não mostra o par (falta %q): %s", esperado, d)
		}
	}
}

func TestRefResolvesQuandoCasa(t *testing.T) {
	t.Run("RFRSR-B03: A reference equal to the sibling spec passes", func(t *testing.T) {})
	casos := map[string][2]string{
		"código":  {"TaxReceipt.ts", "// @anchors\n//   ref: TRT1X\n"},
		"feature": {"TaxReceipt.feature", "# @anchors\n#   ref: TRT1X\n"},
		"teste":   {"TaxReceipt.test.ts", "// @anchors\n//   ref: TRT1X\n"},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			v, d := rodaRefResolves(t, c[0], c[1], "TaxReceipt.spec.md", "<!-- @anchors\n  code: TRT1X\n-->\n")
			if v != Pass {
				t.Fatalf("ref correto deveria passar, foi %s (%s)", v, d)
			}
		})
	}
}

// A spec é DONA (`code:`), não referencia — o gate não se aplica a ela.
func TestRefResolvesNaoSeAplicaASpec(t *testing.T) {
	t.Run("RFRSR-B04: A spec is not confronted", func(t *testing.T) {})
	dir := t.TempDir()
	v, _ := checkRefResolves("<!-- @anchors\n  code: AAAAX\n-->\n",
		mapx.Node{ID: "X.spec.md", Kind: mapx.KindSpec}, dir, nil, nil)
	if v != Skip {
		t.Fatalf("a spec é dona do código — deveria ser Skip, foi %s", v)
	}
	// Um guia ou qualquer outro kind também está fora: só quem REFERENCIA é confrontado.
	if v, _ := checkRefResolves("// ref: AAAAX\n",
		mapx.Node{ID: "X.md", Kind: mapx.KindGuide}, dir, nil, nil); v != Skip {
		t.Fatalf("kind fora do conjunto deveria ser Skip, foi %s", v)
	}
}

// Sem `ref:` declarado, a ausência é do header-conforme — não deste gate.
func TestRefResolvesSemRefEhDoHeaderConforme(t *testing.T) {
	t.Run("RFRSR-B05: An artifact with no reference declared leaves without a verdict", func(t *testing.T) {})
	if v, _ := rodaRefResolves(t, "X.ts", "// nada aqui\n", "X.spec.md", "<!-- @anchors\n  code: AAAAX\n-->\n"); v != Skip {
		t.Fatalf("sem ref deveria ser Skip, foi %s", v)
	}
}

// Cada gate acusa UMA coisa: sem spec irmã, quem cobra é o trinca-completa. Dois gates
// sobre o mesmo defeito viram ruído e o usuário desliga os dois.
func TestRefResolvesSemSpecIrmaEhDoOutroGate(t *testing.T) {
	t.Run("RFRSR-B06: With no sibling spec the gate goes quiet", func(t *testing.T) {})
	if v, d := rodaRefResolves(t, "Solto.ts", "// @anchors\n//   ref: XXXXX\n", "", ""); v != Skip {
		t.Fatalf("sem spec irmã deveria ser Skip, foi %s (%s)", v, d)
	}
}

// Uma spec irmã EXISTE mas não declara `code:`: não há contra o que comparar, e inventar
// uma comparação com string vazia reprovaria todo mundo.
func TestRefResolvesSpecIrmaSemCodigoNaoEhRegua(t *testing.T) {
	t.Run("RFRSR-B07: A sibling spec that declares no identity counts as no sibling", func(t *testing.T) {})
	v, d := rodaRefResolves(t, "Solto.ts", "// @anchors\n//   ref: XXXXX\n",
		"Solto.spec.md", "# Solto\n\nsem cabeçalho de identidade\n")
	if v != Skip {
		t.Fatalf("irmã sem código deveria ser Skip, foi %s (%s)", v, d)
	}
}

// A irmã é achada pela CONVENÇÃO de nome — mesmo tronco, mesmo diretório. Uma spec de
// nome diferente no mesmo diretório não é a irmã de ninguém.
func TestRefResolvesAchaAIrmaPelaConvencao(t *testing.T) {
	t.Run("RFRSR-B08: The sibling is found by the name convention", func(t *testing.T) {})
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	// a irmã certa, no mesmo diretório
	os.WriteFile(filepath.Join(dir, "pkg/Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVOX\n-->\n"), 0o644)
	// uma vizinha de outro tronco, que NÃO pode ser usada como régua
	os.WriteFile(filepath.Join(dir, "pkg/Outro.spec.md"), []byte("<!-- @anchors\n  code: OUTRX\n-->\n"), 0o644)

	if v, d := checkRefResolves("// @anchors\n//   ref: ALVOX\n",
		mapx.Node{ID: "pkg/Alvo.ts", Kind: mapx.KindCode}, dir, nil, nil); v != Pass {
		t.Fatalf("a irmã de mesmo tronco deveria ser a régua, foi %s (%s)", v, d)
	}
	// o arquivo sem irmã própria não herda a vizinha
	if v, d := checkRefResolves("// @anchors\n//   ref: OUTRX\n",
		mapx.Node{ID: "pkg/Terceiro.ts", Kind: mapx.KindCode}, dir, nil, nil); v != Skip {
		t.Fatalf("sem irmã de mesmo tronco deveria ser Skip, foi %s (%s)", v, d)
	}
}

// O sufixo de teste de cada linguagem é removido ANTES do tronco: senão o teste procuraria
// uma spec chamada `X_test.spec.md`, que não existe, e o gate se calaria em todo teste.
func TestRefResolvesTiraOSufixoDeTeste(t *testing.T) {
	t.Run("RFRSR-B09: A test file lands on the same sibling its code does", func(t *testing.T) {})
	for _, arquivo := range []string{
		"Alvo_test.go", "Alvo.test.ts", "Alvo.test.tsx", "Alvo.spec.ts",
		"Alvo_test.py", "Alvo_spec.rb", "Alvo.feature",
	} {
		t.Run(arquivo, func(t *testing.T) {
			dir := t.TempDir()
			os.WriteFile(filepath.Join(dir, "Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVOX\n-->\n"), 0o644)
			marcador := "//"
			if strings.HasSuffix(arquivo, ".feature") {
				marcador = "#"
			}
			v, d := checkRefResolves(marcador+" @anchors\n"+marcador+"   ref: ALVOX\n",
				mapx.Node{ID: arquivo, Kind: mapx.KindTest}, dir, nil, nil)
			if v != Pass {
				t.Fatalf("%s deveria achar Alvo.spec.md, foi %s (%s)", arquivo, v, d)
			}
		})
	}
}

// Uma extensão intermediária (`Alvo.model.ts`) não muda a unidade: o tronco para no
// primeiro ponto, e o arquivo cai na spec do `Alvo`.
func TestRefResolvesTroncoParaNoPrimeiroPonto(t *testing.T) {
	t.Run("RFRSR-B10: An intermediate extension is dropped from the stem", func(t *testing.T) {})
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVOX\n-->\n"), 0o644)
	if v, d := checkRefResolves("// @anchors\n//   ref: ALVOX\n",
		mapx.Node{ID: "Alvo.model.ts", Kind: mapx.KindCode}, dir, nil, nil); v != Pass {
		t.Fatalf("extensão intermediária deveria cair em Alvo.spec.md, foi %s (%s)", v, d)
	}
}

// O cabeçalho é escrito no comentário da LINGUAGEM. Ler só um marcador deixaria famílias
// inteiras de arquivo fora do gate, em silêncio.
func TestRefResolvesLeQualquerMarcadorDeComentario(t *testing.T) {
	t.Run("RFRSR-B11: The reference is read whatever the comment syntax of the language", func(t *testing.T) {})
	for _, cab := range []string{
		"// @anchors\n//   ref: ALVOX\n",
		"# @anchors\n#   ref: ALVOX\n",
		"<!-- @anchors\n  ref: ALVOX\n-->\n",
		" * @anchors\n *   ref: ALVOX\n",
		"@anchors\nref: ALVOX\n",
	} {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVOX\n-->\n"), 0o644)
		if v, d := checkRefResolves(cab, mapx.Node{ID: "Alvo.ts", Kind: mapx.KindCode}, dir, nil, nil); v != Pass {
			t.Errorf("marcador não lido em %q: %s (%s)", cab, v, d)
		}
	}
}

// O comprimento vem da config do PROJETO, lida na chamada. Congelado num `var`, a
// declaração do projeto não teria efeito e o gate se calaria em todo arquivo dele.
func TestRefResolvesComprimentoVemDaEstrutura(t *testing.T) {
	t.Run("RFRSR-B12: The accepted identity length comes from the project's Structure", func(t *testing.T) {})
	original := config.CodeLengths
	t.Cleanup(func() { config.CodeLengths = original })
	config.CodeLengths = []int{4}

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVO\n-->\n"), 0o644)
	if v, d := checkRefResolves("// @anchors\n//   ref: ALVO\n",
		mapx.Node{ID: "Alvo.ts", Kind: mapx.KindCode}, dir, nil, nil); v != Pass {
		t.Fatalf("código de 4 letras declarado pelo projeto deveria ser lido, foi %s (%s)", v, d)
	}
}

// A régua é o arquivo no DISCO. Sem grafo o veredito é o mesmo — e tem de ser: um mapa
// desatualizado é justamente o defeito que este gate existe para pegar.
func TestRefResolvesNaoDependeDoMapa(t *testing.T) {
	t.Run("RFRSR-I01: The ruler is the sibling on disk, never the map", func(t *testing.T) {})
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVOX\n-->\n"), 0o644)
	// o grafo diz outra coisa; o disco manda
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "Alvo.spec.md", Kind: mapx.KindSpec, Code: "DTAXX"}}}
	semGrafo, _ := checkRefResolves("// ref: DTAXX\n", mapx.Node{ID: "Alvo.ts", Kind: mapx.KindCode}, dir, nil, nil)
	comGrafo, _ := checkRefResolves("// ref: DTAXX\n", mapx.Node{ID: "Alvo.ts", Kind: mapx.KindCode}, dir, g, nil)
	if semGrafo != Fail || comGrafo != Fail {
		t.Fatalf("o disco é a régua: sem grafo %s, com grafo %s", semGrafo, comGrafo)
	}
}

// O gate aponta, não conserta. Um gate que reescrevesse o `ref:` passaria na segunda
// execução e o defeito só apareceria em quem clonasse o repositório.
func TestRefResolvesNaoEscreveNoDisco(t *testing.T) {
	t.Run("RFRSR-I02: The gate never repairs what it points at", func(t *testing.T) {})
	dir := t.TempDir()
	conteudo := "// @anchors\n//   ref: DTAXX\n"
	alvo := filepath.Join(dir, "Alvo.ts")
	os.WriteFile(alvo, []byte(conteudo), 0o644)
	os.WriteFile(filepath.Join(dir, "Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVOX\n-->\n"), 0o644)

	checkRefResolves(conteudo, mapx.Node{ID: "Alvo.ts", Kind: mapx.KindCode}, dir, nil, nil)

	depois, err := os.ReadFile(alvo)
	if err != nil {
		t.Fatal(err)
	}
	if string(depois) != conteudo {
		t.Errorf("o gate reescreveu o arquivo — ele aponta, não conserta: %q", string(depois))
	}
}

// Esta é a fronteira com o header-conforme, escrita como restrição: um arquivo COM irmã e
// SEM campo nenhum não é acusado aqui.
func TestRefResolvesNaoCobraAusenciaDoCampo(t *testing.T) {
	t.Run("RFRSR-X01: The gate does not charge the absence of the reference field", func(t *testing.T) {})
	v, d := rodaRefResolves(t, "Alvo.ts", "package alvo\n\nfunc F() {}\n",
		"Alvo.spec.md", "<!-- @anchors\n  code: ALVOX\n-->\n")
	if v != Skip {
		t.Fatalf("ausência do campo é do header-conforme, foi %s (%s)", v, d)
	}
}

// Fronteira com o trinca-completa: a unidade sem spec nenhuma não é acusada aqui.
func TestRefResolvesNaoCobraAusenciaDaSpec(t *testing.T) {
	t.Run("RFRSR-X02: The gate does not charge the absence of the sibling spec", func(t *testing.T) {})
	v, d := rodaRefResolves(t, "Orfao.ts", "// @anchors\n//   ref: ALVOX\n", "", "")
	if v != Skip {
		t.Fatalf("ausência da spec é do trinca-completa, foi %s (%s)", v, d)
	}
}

// Um grafo que CONTRADIZ o disco não muda nada: a resolução não passa pelo mapa.
func TestRefResolvesNaoConsultaOMapa(t *testing.T) {
	t.Run("RFRSR-X03: The gate does not consult the map to resolve the reference", func(t *testing.T) {})
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "Alvo.spec.md"), []byte("<!-- @anchors\n  code: ALVOX\n-->\n"), 0o644)
	// o grafo afirma que a spec do Alvo tem outro código; se o gate olhasse o mapa,
	// este `ref:` reprovaria.
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "Alvo.spec.md", Kind: mapx.KindSpec, Code: "FALSO"}}}
	if v, d := checkRefResolves("// ref: ALVOX\n",
		mapx.Node{ID: "Alvo.ts", Kind: mapx.KindCode}, dir, g, nil); v != Pass {
		t.Fatalf("o veredito tem de seguir o disco, foi %s (%s)", v, d)
	}
}

// A régua é IDENTIDADE, não conteúdo: uma spec que descreve outra coisa passa aqui, e
// julgar isso é de outra classe de gate.
func TestRefResolvesNaoJulgaOConteudoDaSpec(t *testing.T) {
	t.Run("RFRSR-X04: The gate does not judge whether the spec describes the unit well", func(t *testing.T) {})
	v, d := rodaRefResolves(t, "Alvo.ts", "// @anchors\n//   ref: ALVOX\n",
		"Alvo.spec.md", "<!-- @anchors\n  code: ALVOX\n-->\n# Alvo\n\nDescreve um relógio de cozinha.\n")
	if v != Pass {
		t.Fatalf("conteúdo divergente não é assunto deste gate, foi %s (%s)", v, d)
	}
}

// Uma regra com UM cenário por código: este caso continua sendo o comum.
func TestRefResolvesCodigoUnicoContinuaPassando(t *testing.T) {
	if v, _ := rodaRefResolves(t, "A.ts", "// ref: AAAAX\n", "A.spec.md", "<!-- code: AAAAX -->\n"); v != Pass {
		t.Fatalf("queria Pass, foi %s", v)
	}
}

// O ponto INICIAL não separa tronco: um arquivo oculto guarda o nome inteiro. Truncar em
// `i >= 0` deixaria o tronco vazio e faria todo arquivo oculto do diretório cair na mesma
// spec — a de nome só-sufixo.
func TestRefResolvesPontoInicialNaoSeparaTronco(t *testing.T) {
	t.Run("RFRSR-B13: A leading dot is not a stem separator", func(t *testing.T) {})
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".spec.md"), []byte("<!-- @anchors\n  code: DOTFX\n-->\n"), 0o644)
	v, d := checkRefResolves("// @anchors\n//   ref: OUTRX\n",
		mapx.Node{ID: ".hidden.ts", Kind: mapx.KindCode}, dir, nil, nil)
	if v != Skip {
		t.Fatalf("arquivo oculto não tem irmã de nome só-sufixo, foi %s (%s)", v, d)
	}
}
