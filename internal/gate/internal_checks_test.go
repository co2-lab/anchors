package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

func TestHeaderConforme_binarioNaoCarregaCabecalho(t *testing.T) {
	// O baseline de VR é peça de prova e entrou no mapa como `kind: test` — mas é PNG.
	// Exigir cabeçalho `@anchors` num binário é impossível de cumprir, e barrava todo
	// commit de baseline visual. A identidade dele está no NOME do arquivo, que é o
	// que o `identity-consistent` confronta.
	png := "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00"
	if v, msg := checkHeaderConforms(png, mapx.Node{ID: "X.ABCDX-VR-loaded.png", Kind: mapx.KindTest}); v != Skip {
		t.Errorf("binário não carrega cabeçalho: %v (%s)", v, msg)
	}
	// Texto sem header continua reprovando — a dispensa é só para binário.
	if v, _ := checkHeaderConforms("const x = 1\n", mapx.Node{ID: "x.ts", Kind: mapx.KindCode}); v != Fail {
		t.Errorf("arquivo de texto sem header deve reprovar: %v", v)
	}
}

// "Ao menos uma regra catalogada" é o piso, e sozinho deixa passar o caso mais comum: a
// spec cataloga a primeira regra e escreve as outras em prosa. As outras ficam invisíveis
// para os gates de identidade — que então reportam verde sobre o que não conferiram.
func TestSpecSectionsCobraIrmaSemCodigo(t *testing.T) {
	// Três irmãs sob `## Regras`: duas com código, uma sem. A spec estabeleceu o padrão
	// e uma seção destoa.
	comFuro := `# Unidade

## Regras

### ABCDX-B01 — primeira

Texto.

### ABCDX-B02 — segunda

Texto.

### A terceira regra

Esta não tem código, e é regra igual às irmãs.
`
	v, msg := checkSpecSections(comFuro, mapx.Node{}, "", nil, nil)
	if v != Fail {
		t.Errorf("a irmã sem código deveria reprovar, veio %v", v)
	}
	if !strings.Contains(msg, "A terceira regra") {
		t.Errorf("a mensagem deveria NOMEAR a seção que falta: %s", msg)
	}

	// Sem furo: todas as irmãs catalogam.
	semFuro := `# Unidade

## Regras

### ABCDX-B01 — primeira

Texto.

### ABCDX-B02 — segunda

Texto.
`
	if v, msg := checkSpecSections(semFuro, mapx.Node{}, "", nil, nil); v != Pass {
		t.Errorf("todas as irmãs têm código; não havia o que cobrar: %v — %s", v, msg)
	}
}

// A seção de PROSA não pode ser cobrada: `## Visão Geral` e `## Restrições` não catalogam
// regra, e exigir código delas transformaria o gate num cobrador de formato.
func TestSpecSectionsNaoCobraProsa(t *testing.T) {
	spec := `# Unidade

## Visão Geral

Texto livre, sem código nenhum.

## Regras

### ABCDX-B01 — a regra

Texto.

## Restrições

- Uma restrição em prosa.
- Outra.

## Notas de Implementação

Mais prosa.
`
	if v, msg := checkSpecSections(spec, mapx.Node{}, "", nil, nil); v != Pass {
		t.Errorf("seções de prosa não catalogam regra e não podem ser cobradas: %v — %s", v, msg)
	}
}

// Uma irmã sozinha com código NÃO estabelece padrão. Cobrar as outras a partir de um
// único exemplo inventaria uma regra que a spec não declarou.
func TestSpecSectionsNaoInventaPadraoComUmaSo(t *testing.T) {
	spec := `# Unidade

## Regras

### ABCDX-B01 — a única com código

Texto.

### Uma seção

Texto.

### Outra seção

Texto.
`
	if v, msg := checkSpecSections(spec, mapx.Node{}, "", nil, nil); v != Pass {
		t.Errorf("uma irmã só não estabelece padrão: %v — %s", v, msg)
	}
}

// ── seção VAZIA × seção sem código ───────────────────────────────────────────
//
// São dois estados opostos que a busca por código não distingue sozinha, e confundi-los
// inverte o gate: ele passa a pedir "dê um código a cada uma" para uma tabela que não tem
// nenhuma linha.
//
// O caso real que originou: seis telas de onboarding ESTÁTICAS (zero efeitos no código)
// herdaram o cabeçalho de "Comportamentos Automáticos" do modelo de spec. O gate as
// barrava por não catalogarem regras que não existem.

func comCodigosDe4(t *testing.T) {
	t.Helper()
	orig := config.CodeLengths
	config.CodeLengths = []int{4}
	SetRuleLetters(config.DefaultRuleLetters)
	t.Cleanup(func() {
		config.CodeLengths = orig
		SetRuleLetters(config.DefaultRuleLetters)
	})
}

const specComSecaoVazia = `# Tela

## Rules

### Permissões e Acesso

| Regra | Descrição |
| ---------- | --------- |
| ` + "`PHA1-R01`" + ` | Pública |

### Ações Permitidas

| Regra | Ação | Resultado |
| ---------- | ---- | --------- |
| ` + "`PHA1-A01`" + ` | Toque | Navega |

### Comportamentos Automáticos

| Regra | Gatilho | Ação Automática |
| ---------- | ------- | --------------- |

---
`

// Vazia SEM declaração: o gate não sabe se foi decisão ou esquecimento, e pede que
// alguém diga. Trocar o falso positivo antigo por um falso negativo seria pior — o gate
// passaria a reportar verde sobre seção que ninguém preencheu.
func TestSecaoVaziaSemDeclaracaoPedeQueAlguemDiga(t *testing.T) {
	comCodigosDe4(t)
	d := siblingsWithoutCode(specComSecaoVazia)
	if d == "" {
		t.Fatal("vazia sem declaração tem de ser cobrada — senão o esquecimento passa")
	}
	if !strings.Contains(d, "@no-content") {
		t.Fatalf("o laudo tem de ENSINAR a saída; veio: %s", d)
	}
	// e NÃO pode pedir código para uma tabela sem linhas
	if strings.Contains(d, "Dê um código a cada uma") {
		t.Fatalf("não há 'cada uma' numa seção vazia; veio: %s", d)
	}
}

// Vazia COM declaração: a seção FICA (importa quando é obrigatória por regulação) e o
// gate absolve, porque alguém decidiu e escreveu o porquê.
func TestSecaoVaziaComNoContentEAceita(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"### Comportamentos Automáticos\n\n| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"### Comportamentos Automáticos\n\n@no-content: tela estática — não há efeito, timer nem carga.\n",
		1)
	if d := siblingsWithoutCode(spec); d != "" {
		t.Fatalf("declarada, a seção vazia é aceita. Veio: %s", d)
	}
}

// `@no-content` SEM motivo não conta — dispensa sem porquê é a que ninguém revisa depois.
func TestNoContentExigeMotivo(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"### Comportamentos Automáticos\n\n| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"### Comportamentos Automáticos\n\n@no-content:\n",
		1)
	if d := siblingsWithoutCode(spec); d == "" {
		t.Fatal("`@no-content` sem motivo não pode absolver")
	}
}

// O defeito REAL continua pego: a seção que tem regra escrita em prosa, sem código.
func TestIrmasSemCodigoAindaPegaRegraSemCodigo(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n| — | Abertura | Carrega o perfil |\n",
		1)
	d := siblingsWithoutCode(spec)
	if d == "" {
		t.Fatal("regra SEM código tem de ser acusada — é o defeito que o gate existe para pegar")
	}
	if !strings.Contains(d, "Comportamentos Automáticos") {
		t.Fatalf("o laudo tem de NOMEAR a seção; veio: %s", d)
	}
}

// Prosa também é conteúdo: a seção com texto e sem código segue acusada.
func TestIrmasSemCodigoPegaSecaoComProsa(t *testing.T) {
	comCodigosDe4(t)
	spec := strings.Replace(specComSecaoVazia,
		"### Comportamentos Automáticos\n\n| Regra | Gatilho | Ação Automática |\n| ---------- | ------- | --------------- |\n",
		"### Comportamentos Automáticos\n\nAo abrir, a tela carrega o perfil do usuário.\n",
		1)
	if d := siblingsWithoutCode(spec); d == "" {
		t.Fatal("regra em prosa sem código tem de ser acusada")
	}
}

// O `####` e CONTEUDO do `###` que o precede, nao irmao dele.
//
// Tratando-os como irmaos, a linha de conteudo ia para a ULTIMA secao vista — o `####`
// roubava as linhas do pai, e o `###` aparecia VAZIO. Medido em `SplashScreen.spec.md`:
// `### Caminhos Condicionais` tem tres linhas de tabela dentro de `#### destination`, e o
// gate acusava o pai de vazio com o conteudo logo abaixo.
func TestSubsecaoNaoEsvaziaOPai(t *testing.T) {
	comCodigosDe4(t)
	spec := "# Tela\n\n## Rules\n\n" +
		"### Comportamentos Automáticos\n\n" +
		"| Regra | Gatilho |\n| --- | --- |\n| `SPAX-B01` | Abertura |\n\n" +
		"### Caminhos Condicionais\n\n" +
		"#### `destination` — rota de destino\n\n" +
		"| Data State | Condição |\n| --- | --- |\n| `DS-dest-main` | autenticado |\n\n---\n"

	if d := siblingsWithoutCode(spec); d != "" {
		t.Fatalf("o `###` tem conteúdo no `####` filho; não devia acusar. Veio: %s", d)
	}
}

// O NOME DE UM GATE E' CONTRATO PUBLICO -- ele aparece no `anchors.yaml` de todo
// projeto, na saida do `check` e nas issues que o pipeline abre.
//
// Metade do vocabulario estava em portugues (`regra-implementada`, `cenario-identidade`)
// e parte traduzida ao pe da letra (`header-conforme` -- meia palavra em cada idioma, e
// "conforme" nao diz o que o gate faz). O `board.json` ja foi migrado para ingles pela
// mesma razao: contrato e' superficie, e enquanto metade esta num idioma cada gate novo
// herda a duvida sobre qual convencao seguir.
//
// Esta regua fecha a porta. Nome de gate e' em INGLES, kebab-case.
func TestNomeDeGateEmIngles(t *testing.T) {
	// As palavras que denunciam portugues nos nomes que ja existiram aqui. Nao e' um
	// dicionario -- e' a lista do que ja entrou, para que nao volte.
	emPortugues := []string{
		"cenario", "regra", "conforme", "coerente", "consultado", "preenchido",
		"existe", "implementada", "identidade", "declarada", "alinhado", "dominio",
		"fase", "prova", "trinca", "promovivel", "progresso", "idioma", "carimbado",
		"tipado", "ancorado", "valor", "codigo", "teste", "rastreavel",
	}
	todos := map[string]bool{}
	for n := range internalCheckers {
		todos[n] = true
	}
	for n := range checkersWithRoot {
		todos[n] = true
	}
	for n := range checkersWithGraph {
		todos[n] = true
	}
	if len(todos) == 0 {
		t.Fatal("nenhum gate registrado — o teste passaria vazio")
	}

	for nome := range todos {
		for _, p := range emPortugues {
			// `-p-`, `p-` no comeco ou `-p` no fim: palavra inteira, nao substring.
			if nome == p ||
				strings.HasPrefix(nome, p+"-") ||
				strings.HasSuffix(nome, "-"+p) ||
				strings.Contains(nome, "-"+p+"-") {
				t.Errorf("o gate `%s` tem `%s` no nome — nome de gate e' contrato publico "+
					"e vai em INGLES; ele aparece no `anchors.yaml` de todo projeto e na "+
					"saida do `check`", nome, p)
			}
		}
	}
}

// ── O REGISTRO: como um `check:` declarado vira função ────────────────────────

// O nome declarado no gate roteia para a função registrada sob ele. É o contrato
// inteiro do registro — sem ele, `check:` seria decoração.
func TestRegistroRoteiaONomeDeclarado(t *testing.T) {
	t.Run("INCHN-B01: A declared name routes to the function registered under it", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "vazio.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("   \n  "), 0o644); err != nil {
		t.Fatal(err)
	}
	v, _ := runInternal("non-empty", mapx.Node{ID: alvo}, root, nil, nil)
	if v != Fail {
		t.Errorf("o nome declarado tem de chegar na função dele; veio %v", v)
	}
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("conteúdo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := runInternal("non-empty", mapx.Node{ID: alvo}, root, nil, nil); v != Pass {
		t.Errorf("a mesma rota tem de devolver o veredito oposto para o caso oposto; veio %v", v)
	}
}

// UM NOME QUE NÃO RESOLVE NÃO APROVA. Um checker declarado e inexistente não mediu
// nada — e "não medi" não é "está limpo". Aprovar aqui carimbaria verde sobre uma
// verificação que nunca rodou.
func TestNomeQueNaoResolveNaoAprova(t *testing.T) {
	t.Run("INCHN-B02: A name that does not resolve answers undetermined", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "x.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("conteúdo"), 0o644); err != nil {
		t.Fatal(err)
	}
	v, msg := runInternal("checker-que-ninguem-escreveu", mapx.Node{ID: alvo}, root, nil, nil)
	if v != Pending {
		t.Fatalf("nome não resolvido tem de ser indeterminado, veio %v", v)
	}
	if !strings.Contains(msg, "checker-que-ninguem-escreveu") {
		t.Errorf("o laudo tem de NOMEAR o checker que não roteou: %s", msg)
	}
}

// A ordem das tentativas é o que permite a um checker GANHAR uma dependência sem trocar
// de nome: ele muda de registro, e o `check:` declarado nos projetos continua valendo.
func TestRoteamentoTentaORelacionalPrimeiro(t *testing.T) {
	t.Run("INCHN-B03: The routing tries the relational registry before the simpler ones", func(t *testing.T) {})

	// Um nome plantado nos DOIS registros: o relacional tem de vencer. É o que permite a
	// um checker GANHAR uma dependência sem trocar de nome — ele muda de registro, e o
	// `check:` declarado nos projetos de fora continua valendo.
	const nome = "duplo-para-o-teste"
	marcaPura, marcaRelacional := "respondeu o PURO", "respondeu o RELACIONAL"

	internalCheckers[nome] = func(string, mapx.Node) (Verdict, string) { return Fail, marcaPura }
	checkersWithGraph[nome] = func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string) {
		return Pass, marcaRelacional
	}
	t.Cleanup(func() {
		delete(internalCheckers, nome)
		delete(checkersWithGraph, nome)
	})

	root := t.TempDir()
	alvo := "a.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("texto\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	v, msg := runInternal(nome, mapx.Node{ID: alvo}, root, nil, nil)
	if msg != marcaRelacional || v != Pass {
		t.Errorf("o relacional tem de ser tentado primeiro; respondeu %v — %s", v, msg)
	}

	// E o que existe SÓ no puro continua sendo alcançado: a ordem prioriza, não exclui.
	if vp, _ := runInternal("non-empty", mapx.Node{ID: alvo}, root, nil, nil); vp != Pass {
		t.Errorf("o registro puro continua alcançável: %v", vp)
	}
}

// Um gate GENÉRICO só sabe o que procurar depois de ler a PRÓPRIA configuração, e um
// projeto declara vários dele. Sem o gate chegando ao checker, todas as instâncias
// veriam a mesma configuração — ou nenhuma.
func TestAgregadoComGateRecebeOGateQueOInvocou(t *testing.T) {
	t.Run("INCHN-B07: An aggregate checker that needs the declaring gate receives it", func(t *testing.T) {})

	const nome = "parametrizado-para-o-teste"
	var vistoPeloChecker string
	checkersComGate[nome] = func(g config.Gate, _ string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
		vistoPeloChecker = g.Name
		return Pass, g.Name
	}
	// O MESMO nome também no registro relacional, que NÃO recebe o gate: o roteamento
	// tem de preferir o que recebe — senão a instância se perde.
	checkersWithGraph[nome] = func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string) {
		return Pass, "sem saber quem perguntou"
	}
	t.Cleanup(func() {
		delete(checkersComGate, nome)
		delete(checkersWithGraph, nome)
	})

	root := t.TempDir()
	for _, instancia := range []string{"marker-a", "marker-b"} {
		g := config.Gate{Name: instancia, Check: nome, Scope: config.ScopeProject}
		if _, msg := runInternalAggregate(g, root, nil, nil); msg != instancia {
			t.Errorf("a instância `%s` não chegou ao checker; veio %q", instancia, msg)
		}
		if vistoPeloChecker != instancia {
			t.Errorf("o checker viu %q quando quem perguntou foi %q", vistoPeloChecker, instancia)
		}
	}
}

// O caminho POR NÓ lê o arquivo do alvo, e a leitura que falha é FALHA: um checker de
// conteúdo sem conteúdo não tem o que responder.
func TestCaminhoPorNoReprovaQuandoNaoConsegueLer(t *testing.T) {
	t.Run("INCHN-B04: The per-node path reads the target file, and a failed read is a failure", func(t *testing.T) {})

	v, _ := runInternal("non-empty", mapx.Node{ID: "nao-existe.md"}, t.TempDir(), nil, nil)
	if v != Fail {
		t.Errorf("arquivo ilegível no caminho por nó é falha, veio %v", v)
	}
}

// O caminho AGREGADO não lê arquivo nenhum: o escopo é o CONJUNTO. Tentar ler aqui
// devolveria erro de leitura e o gate reprovaria por um arquivo que nunca existiu.
func TestCaminhoAgregadoNaoTentaLerArquivo(t *testing.T) {
	t.Run("INCHN-B05: The aggregate path reads no file at all", func(t *testing.T) {})

	// Um gate agregado cujo checker é relacional: o nó chega VAZIO, e mesmo assim não há
	// erro de leitura no laudo.
	g := config.Gate{Name: "doc-fresh", Check: "docs-fresh", Scope: config.ScopeProject}
	v, msg := runInternalAggregate(g, t.TempDir(), nil, nil)
	if v == Fail && strings.Contains(msg, "nao-existe") {
		t.Fatalf("o agregado não pode reprovar por leitura de arquivo: %v — %s", v, msg)
	}
	// E o caminho POR NÓ, com o mesmo nó vazio, reprovaria — é a diferença entre os dois.
	if vn, _ := runInternal("docs-fresh", mapx.Node{}, t.TempDir(), nil, nil); vn != Fail {
		t.Errorf("o caminho por nó lê o arquivo e reprova sem ele; veio %v", vn)
	}
}

// O agregado que não resolve também é indeterminado, e o laudo nomeia o check.
func TestAgregadoQueNaoResolveEhIndeterminado(t *testing.T) {
	t.Run("INCHN-B06: An unresolved name in the aggregate path is undetermined too", func(t *testing.T) {})

	g := config.Gate{Name: "x", Check: "agregado-inexistente", Scope: config.ScopeProject}
	v, msg := runInternalAggregate(g, t.TempDir(), nil, nil)
	if v != Pending {
		t.Fatalf("agregado desconhecido tem de ser indeterminado, veio %v", v)
	}
	if !strings.Contains(msg, "agregado-inexistente") {
		t.Errorf("o laudo tem de NOMEAR o check que não roteou: %s", msg)
	}
}

// A GRAMÁTICA DE CÓDIGO é do PROJETO. Um regex adicionado sem registrá-lo aqui fica
// preso às letras canônicas — e um cenário de letra declarada pelo projeto vira
// invisível para aquele gate, que então reporta verde sobre o que não conferiu.
func TestLetrasDeRegraReconfiguramTodosOsRegexes(t *testing.T) {
	t.Run("INCHN-B08: Setting the rule letters reconfigures every dependent pattern together", func(t *testing.T) {})

	t.Cleanup(func() { SetRuleLetters(config.DefaultRuleLetters) })

	const letraDoProjeto = "Z"
	if strings.Contains(config.DefaultRuleLetters, letraDoProjeto) {
		t.Fatalf("o teste precisa de uma letra FORA do vocabulário canônico; %q está nele", letraDoProjeto)
	}
	codigo := "ABCDX-" + letraDoProjeto + "01"

	// Antes de declarar: a letra é invisível para a régua de identidade.
	SetRuleLetters(config.DefaultRuleLetters)
	if anyCodeRE.MatchString(codigo) {
		t.Fatalf("%s não deveria ser reconhecido com o vocabulário canônico", codigo)
	}

	// Declarada, ela passa a valer nos DOIS regexes que dependem do vocabulário.
	SetRuleLetters(config.DefaultRuleLetters + letraDoProjeto)
	if !anyCodeRE.MatchString(codigo) {
		t.Error("a régua de identidade não reconheceu a letra declarada pelo projeto")
	}
	if !featScenarioCodeRE.MatchString("  @" + codigo + " @unit-level") {
		t.Error("a régua de cenário da feature ficou presa às letras canônicas — um regex " +
			"que não é reconfigurado junto reporta verde sobre o que não conferiu")
	}
}

// O arquivo vazio (ou só espaço) é o piso trivial, e ele pega placeholder.
func TestVazioReprovaEConteudoPassa(t *testing.T) {
	t.Run("INCHN-B09: A file that is empty or only whitespace fails the emptiness ruler", func(t *testing.T) {})

	if v, _ := checkNonEmpty("   \n\t  \n", mapx.Node{}); v != Fail {
		t.Errorf("só espaço é vazio, veio %v", v)
	}
	if v, _ := checkNonEmpty("uma linha", mapx.Node{}); v != Pass {
		t.Errorf("conteúdo de verdade tem de passar, veio %v", v)
	}
}

// Sem código, a peça é ÓRFÃ INVISÍVEL — nenhum gate de identidade a enxerga.
func TestIdentidadeCobraOCodigoDeCenario(t *testing.T) {
	t.Run("INCHN-B10: The identity ruler charges the presence of a scenario code", func(t *testing.T) {})

	SetRuleLetters(config.DefaultRuleLetters)
	if v, _ := checkHasScenarioCode("it('LOGIX-A01: faz algo')", mapx.Node{}); v != Pass {
		t.Error("o código presente tem de passar")
	}
	if v, _ := checkHasScenarioCode("it('faz algo')", mapx.Node{}); v != Fail {
		t.Error("sem código a peça é órfã invisível, e tem de reprovar")
	}
}

// Sem BLOCO de cabeçalho nenhum: o arquivo não declara identidade de forma alguma.
func TestHeaderSemBlocoReprova(t *testing.T) {
	t.Run("INCHN-B11: A governed file with no identity block fails the header ruler", func(t *testing.T) {})

	if v, _ := checkHeaderConforms("const x = 1 // nada aqui\n",
		mapx.Node{ID: "x.ts", Kind: mapx.KindCode, Tags: []string{"business-logic"}}); v != Fail {
		t.Error("sem bloco de cabeçalho, a camada regida reprova")
	}
}

// Camada REGIDA exige POSSE (`code:`) ou REFERÊNCIA (`ref:`). Só `layer:` não basta —
// quem tem spec dona ou irmã tem de dizer qual.
func TestHeaderRegidoExigePosseOuReferencia(t *testing.T) {
	t.Run("INCHN-B12: A governed file passes with ownership or with reference, never with layer alone", func(t *testing.T) {})

	regida := mapx.Node{ID: "x.ts", Kind: mapx.KindCode, Tags: []string{"business-logic"}}
	if v, _ := checkHeaderConforms("// @anchors\n//   code: LGNNX\nconst x = 1\n", regida); v != Pass {
		t.Error("posse (code) tem de passar")
	}
	if v, _ := checkHeaderConforms("// @anchors\n//   ref: LGNNX\nconst x = 1\n", regida); v != Pass {
		t.Error("referência (ref) tem de passar")
	}
	if v, _ := checkHeaderConforms("// @anchors\n//   layer: business-logic\nconst x = 1\n", regida); v != Fail {
		t.Error("só layer NÃO basta numa camada regida")
	}
}

// Camada RECONHECIDA não tem spec dona nem irmã a referenciar: a identidade mínima
// honesta dela é a LAYER, e inventar code/ref seria pedir uma identidade falsa.
func TestHeaderReconhecidoPassaComLayer(t *testing.T) {
	t.Run("INCHN-B13: A file of a recognised layer passes with the layer alone", func(t *testing.T) {})

	reconhecida := mapx.Node{ID: "dao.ts", Kind: mapx.KindCode, Tags: []string{"frontend", "presentation"}}
	if v, _ := checkHeaderConforms("// @anchors\n//   layer: presentation\n", reconhecida); v != Pass {
		t.Error("camada reconhecida passa com a layer sozinha")
	}
	// Mas SEM identidade nenhuma ela também reprova — a dispensa é da forma, não do dever.
	if v, _ := checkHeaderConforms("// @anchors\n//   updated_at: x\n", reconhecida); v != Fail {
		t.Error("reconhecida sem layer/code/ref ainda reprova")
	}
}

// Não há sintaxe de comentário num PNG. Cobrar cabeçalho de binário exigiria o
// impossível e barraria todo commit de baseline visual.
func TestHeaderDispensaBinario(t *testing.T) {
	t.Run("INCHN-B14: A binary file steps aside from the header ruler", func(t *testing.T) {})

	png := "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00"
	if v, _ := checkHeaderConforms(png, mapx.Node{ID: "X.ABCDX-VR-loaded.png", Kind: mapx.KindTest}); v != Skip {
		t.Error("binário não carrega cabeçalho")
	}
}

// O ROTEIRO de teste executável sai pelo mesmo motivo por um caminho diferente: o
// formato é do RUNNER, e um bloco nosso no topo é comentário morto para quem o executa.
func TestHeaderDispensaRoteiroExecutavel(t *testing.T) {
	t.Run("INCHN-B15: An executable test script steps aside by a different path", func(t *testing.T) {})

	roteiro := mapx.Node{ID: "flows/ABCDX-A01.yaml", Kind: mapx.KindTest}
	if v, _ := checkHeaderConforms("steps:\n  - open: /login\n", roteiro); v != Skip {
		t.Error("o roteiro do runner tem a identidade no NOME, e não carrega cabeçalho nosso")
	}
	// Mas um teste em linguagem de programação continua cobrado.
	emGo := mapx.Node{ID: "x_test.go", Kind: mapx.KindTest}
	if v, _ := checkHeaderConforms("package x\n", emGo); v != Fail {
		t.Error("teste em código nosso continua cobrado")
	}
}

// Um guide sem PONTOS DE CONFORMIDADE deixa o gate de julgamento por IA recair em
// heurística vaga — e aí ele não mede, ele adivinha.
func TestGuideSemPontosDeConformidadeReprova(t *testing.T) {
	t.Run("INCHN-B16: A guide without compliance points fails", func(t *testing.T) {})

	if v, _ := checkGuideHasChecklist("# Guia\n\nProsa sem seção alguma.\n", mapx.Node{}); v != Fail {
		t.Error("sem a seção, o guide reprova")
	}
	if v, _ := checkGuideHasChecklist("# Guia\n\n## Compliance points\n\nProsa, e nenhum item.\n",
		mapx.Node{}); v != Fail {
		t.Error("a seção sem item nenhum também reprova — ela é o continente, não o conteúdo")
	}
	if v, _ := checkGuideHasChecklist("# Guia\n\n## Compliance points\n\n- CK1 — algo verificável\n",
		mapx.Node{}); v != Pass {
		t.Error("seção com item tem de passar")
	}
}

// NOS DOIS CAMINHOS o nome que não resolve não aprova. Um deles aprovando seria a porta
// lateral: bastaria declarar o escopo certo para carimbar verde sem medir.
func TestNomeNaoResolvidoNuncaAprovaEmNenhumCaminho(t *testing.T) {
	t.Run("INCHN-I01: An unresolved name never approves, on either path", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "x.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("texto"), 0o644); err != nil {
		t.Fatal(err)
	}
	const inexistente = "nao-existe-em-registro-algum"

	if v, _ := runInternal(inexistente, mapx.Node{ID: alvo}, root, nil, nil); v == Pass {
		t.Error("o caminho por nó não pode aprovar um checker que não existe")
	}
	g := config.Gate{Name: "g", Check: inexistente, Scope: config.ScopeProject}
	if v, _ := runInternalAggregate(g, root, nil, nil); v == Pass {
		t.Error("o caminho agregado não pode aprovar um checker que não existe")
	}
}

// Todo nome REGISTRADO é alcançável por exatamente um dos caminhos. Um nome registrado
// em lugar irrecuperável é um gate aceito em silêncio que não mede nada.
func TestTodoNomeRegistradoEhAlcancavel(t *testing.T) {
	t.Run("INCHN-I02: Every registered name is reachable through exactly one routing path", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "x.md"
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("texto\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nomes := map[string]bool{}
	for n := range internalCheckers {
		nomes[n] = true
	}
	for n := range checkersWithRoot {
		nomes[n] = true
	}
	for n := range checkersWithGraph {
		nomes[n] = true
	}
	if len(nomes) == 0 {
		t.Fatal("nenhum checker registrado — o teste passaria vazio")
	}
	for nome := range nomes {
		_, msg := runInternal(nome, mapx.Node{ID: alvo}, root, nil, nil)
		if strings.Contains(msg, i18n.T("gate.checker_not_implemented", nome)) {
			t.Errorf("o checker `%s` está registrado e o roteamento não o alcança — "+
				"um gate aceito em silêncio que não mede nada", nome)
		}
	}
}

// O `anchors init` semeia o guide com o título TRADUZIDO pelo `lang:` do projeto. Uma
// régua que só casasse um idioma faria o projeto nascer reprovando um guide que a
// própria ferramenta acabara de escrever.
func TestPontosDeConformidadeEmQualquerIdiomaDoCatalogo(t *testing.T) {
	t.Run("INCHN-I03: The compliance ruler is recognised in every language of the catalogue", func(t *testing.T) {})

	titulos := i18n.AllTranslations("section.title.compliance_points")
	if len(titulos) < 2 {
		t.Fatalf("o catálogo precisa de mais de um idioma para este confronto: %v", titulos)
	}
	for _, titulo := range titulos {
		guide := "# Guia\n\n## " + titulo + "\n\n- CK1 — algo verificável\n"
		if v, msg := checkGuideHasChecklist(guide, mapx.Node{}); v != Pass {
			t.Errorf("o título %q não foi reconhecido: %v — %s", titulo, v, msg)
		}
	}
}

// O registro não decide o que o projeto roda: quem declara é a Estrutura. Rodar o que
// ninguém pediu cobraria do projeto uma régua que ele nunca adotou.
func TestRegistroNaoDecideOQueOProjetoRoda(t *testing.T) {
	t.Run("INCHN-X01: The registry does not decide which checks a project runs", func(t *testing.T) {})

	root := t.TempDir()
	alvo := "sem-codigo.md"
	// O arquivo reprovaria o `has-code` — mas o projeto declarou só o `non-empty`.
	if err := os.WriteFile(filepath.Join(root, alvo), []byte("prosa sem identidade\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	declarado := config.Gate{Name: "não-vazio", ID: "non-empty", On: []string{"spec"}, Check: "non-empty"}
	res := Run([]config.Gate{declarado}, []mapx.Node{{ID: alvo, Kind: mapx.KindSpec}}, root, nil)

	if len(res) != 1 {
		t.Fatalf("só o gate declarado deveria ter rodado, vieram %d resultados", len(res))
	}
	if res[0].Verdict != Pass {
		t.Errorf("o gate declarado passa; o registro não pode cobrar os outros: %v — %s",
			res[0].Verdict, res[0].Detail)
	}
}

// Estes checkers respondem lendo TEXTO. Se dependessem de binário instalado, a mesma
// pergunta teria respostas diferentes em duas máquinas.
func TestCheckerInternoNaoDependeDeFerramentaInstalada(t *testing.T) {
	t.Run("INCHN-X02: The registry does not invoke external tooling", func(t *testing.T) {})

	t.Setenv("PATH", "")
	if v, _ := checkNonEmpty("conteúdo", mapx.Node{}); v != Pass {
		t.Error("o checker de texto tem de responder sem PATH algum")
	}
	if v, _ := checkHeaderConforms("// @anchors\n//   code: LGNNX\n",
		mapx.Node{ID: "x.ts", Kind: mapx.KindCode}); v != Pass {
		t.Error("a régua de cabeçalho é textual e não pode depender do que está instalado")
	}
}

// A régua aqui é PRESENÇA e FORMA. Julgar se o conteúdo está certo é outra régua, e um
// checker determinístico que tentasse reprovaria por um critério que não sabe medir.
func TestRegistroNaoJulgaSeOTextoEstaCerto(t *testing.T) {
	t.Run("INCHN-X03: The registry does not judge whether the text is good", func(t *testing.T) {})

	// Cabeçalho presente e conforme; o conteúdo, absurdo para qualquer revisor.
	absurdo := "// @anchors\n//   code: LGNNX\n// esta unidade faz exatamente o oposto do que diz\n"
	if v, msg := checkHeaderConforms(absurdo, mapx.Node{ID: "x.ts", Kind: mapx.KindCode}); v != Pass {
		t.Errorf("a régua é presença e forma, não qualidade: %v — %s", v, msg)
	}
}

// A casca vazia: o cabeçalho Gherkin sozinho enche oito linhas sem declarar cenário, e
// `TrimSpace` não vê diferença entre isso e um arquivo com conteúdo. Medido no app de
// referência: 12 features de `services/` exatamente assim, todas aprovadas.
func TestNonEmptyFeatureSemCenario(t *testing.T) {
	casca := "# language: pt\n# @anchors\n#   ref: SGABX\n#   layer: service\n\n@backend @service\nFuncionalidade: auth (service) — Gateway\n"
	v, d := checkNonEmpty(casca, mapx.Node{Kind: mapx.KindFeature})
	if v != Fail {
		t.Fatalf("feature sem cenário deveria reprovar, foi %s (%s)", v, d)
	}
}

func TestNonEmptyFeatureComCenario(t *testing.T) {
	ok := "# language: pt\nFuncionalidade: x\n\n  @comportamento @ABCD-B01\n  Cenário: faz algo\n    Dado que sim\n"
	if v, d := checkNonEmpty(ok, mapx.Node{Kind: mapx.KindFeature}); v != Pass {
		t.Fatalf("feature com cenário deveria passar, foi %s (%s)", v, d)
	}
}

// A regra do cenário vale SÓ para feature: um .ts com conteúdo não declara cenário e
// continua passando, senão o gate reprovaria todo arquivo de código do projeto.
func TestNonEmptyNaoExigeCenarioForaDeFeature(t *testing.T) {
	if v, _ := checkNonEmpty("export const x = 1\n", mapx.Node{Kind: mapx.KindCode}); v != Pass {
		t.Fatalf("código com conteúdo deveria passar, foi %s", v)
	}
}

// O estado de dado leva o NOME junto (`SECU-DS-bio-on`). A regex que lê o TESTE parava
// no `DS-` enquanto a irmã que lê a FEATURE capturava o nome inteiro — as duas pontas do
// mesmo contrato com gramáticas diferentes. Medido no app de referência: 78 dos 124
// códigos "órfãos" do `test-feature-match` eram este truncamento.
func TestAnyCodeRECapturaNomeDoEstadoDeDado(t *testing.T) {
	// Códigos de CINCO caracteres: é o default do engine (`config.CodeLengths`), e este
	// teste não reconfigura o global — o app de referência usa 4 por declarar
	// `code_lengths: [4]` no próprio anchors.yaml.
	casos := map[string]string{
		"it('SECUX-DS-bio-on: ...')":     "SECUX-DS-bio-on",
		"it('TREXX-DS-filter-12m: ...')": "TREXX-DS-filter-12m",
		"it('HOMEX-B01: ...')":           "HOMEX-B01",
		"it('MNPMX-VR: ...')":            "MNPMX-VR",
	}
	for entrada, esperado := range casos {
		got := anyCodeRE.FindString(entrada)
		if got != esperado {
			t.Errorf("%s → %q, esperado %q", entrada, got, esperado)
		}
	}
}
