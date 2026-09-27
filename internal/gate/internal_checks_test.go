package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
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

// The empty shell: the Gherkin header alone fills eight lines without declaring a
// scenario, and `TrimSpace` sees no difference between that and a file with content.
// Measured in the reference app: 12 features under `services/` exactly like this, all
// approved.
func TestNonEmptyFeatureWithoutScenario(t *testing.T) {
	shell := "# language: pt\n# @anchors\n#   ref: SGABX\n#   layer: service\n\n@backend @service\nFuncionalidade: auth (service) — Gateway\n"
	v, d := checkNonEmpty(shell, mapx.Node{Kind: mapx.KindFeature})
	if v != Fail {
		t.Fatalf("a feature with no scenario should fail, got %s (%s)", v, d)
	}
}

func TestNonEmptyFeatureWithScenario(t *testing.T) {
	ok := "# language: pt\nFuncionalidade: x\n\n  @comportamento @ABCD-B01\n  Cenário: faz algo\n    Dado que sim\n"
	if v, d := checkNonEmpty(ok, mapx.Node{Kind: mapx.KindFeature}); v != Pass {
		t.Fatalf("a feature with a scenario should pass, got %s (%s)", v, d)
	}
}

// The scenario rule holds ONLY for features: a .ts with content declares no scenario and
// still passes, otherwise the gate would fail every code file of the project.
func TestNonEmptyDoesNotDemandScenarioOutsideFeatures(t *testing.T) {
	if v, _ := checkNonEmpty("export const x = 1\n", mapx.Node{Kind: mapx.KindCode}); v != Pass {
		t.Fatalf("code with content should pass, got %s", v)
	}
}

// The data state carries the NAME with it (`SECU-DS-bio-on`). The regex that reads the
// TEST stopped at `DS-` while the sibling that reads the FEATURE captured the whole name
// — the two ends of the same contract with different grammars. Measured in the reference
// app: 78 of the 124 "orphan" codes of `test-feature-match` were this truncation.
func TestAnyCodeRECapturesTheDataStateName(t *testing.T) {
	// FIVE-character codes: that is the engine's default (`config.CodeLengths`), and this
	// test does not reconfigure the global — the reference app uses 4 by declaring
	// `code_lengths: [4]` in its own anchors.yaml.
	cases := map[string]string{
		"it('SECUX-DS-bio-on: ...')":     "SECUX-DS-bio-on",
		"it('TREXX-DS-filter-12m: ...')": "TREXX-DS-filter-12m",
		"it('HOMEX-B01: ...')":           "HOMEX-B01",
		"it('MNPMX-VR: ...')":            "MNPMX-VR",
	}
	for input, want := range cases {
		got := anyCodeRE.FindString(input)
		if got != want {
			t.Errorf("%s → %q, want %q", input, got, want)
		}
	}
}

func TestInternalChecks_Errors(t *testing.T) {
	t.Run("INCHN-E01: A test missing from disk does not hide the codes the other tests name", func(t *testing.T) {
		root, g := rootWithTest(t, "const c = \"CREDX-B01\"\nfunc TestX(t *testing.T) {}\n")
		g.Nodes = append([]mapx.Node{{ID: "gone_test.go", Kind: mapx.KindTest}}, g.Nodes...)
		v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil)
		if v != Fail || !strings.Contains(msg, "ingest") || !strings.Contains(msg, "CREDX-B01") {
			t.Fatalf("CREDX-B01 must still count as written despite the missing test, got %v: %s", v, msg)
		}
	})
}

func TestScenarioCoverage_TitlesWhenDeclared(t *testing.T) {
	t.Run("INCHN-B17: With a tests source a scenario is written only when a title cites it", func(t *testing.T) {})
	t.Cleanup(resetProjectTestsCache)
	body := "const c = \"CREDX-B01\"\nfunc TestX(t *testing.T) { t.Run(\"unrelated\", nil) }\n"
	goFamily := &config.Config{Dialect: &config.Dialect{Family: "go"}}
	resetProjectTestsCache()
	root, g := rootWithTest(t, body)
	v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, goFamily)
	if v != Fail || strings.Contains(msg, "ingest") {
		t.Fatalf("with titles declared, a code in a fixture writes no test, got %v: %s", v, msg)
	}
	resetProjectTestsCache()
	root, g = rootWithTest(t, body)
	if v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil); v != Fail || !strings.Contains(msg, "ingest") {
		t.Fatalf("without a declaration the code in the file counts as written, got %v: %s", v, msg)
	}
	resetProjectTestsCache()
	root, g = rootWithTest(t, "func TestX(t *testing.T) { t.Run(\"CREDX-B01: validates the limit\", nil) }\n")
	if _, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, goFamily); !strings.Contains(msg, "ingest") {
		t.Fatalf("a title citing the code writes its test, got %s", msg)
	}
}

func TestScenarioCoverage_FailingSource(t *testing.T) {
	t.Run("INCHN-E02: A failing tests source fails scenario-coverage naming the error", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root, g := rootWithTest(t, "func TestX(t *testing.T) {}\n")
	cfg := &config.Config{Dialect: &config.Dialect{Tests: &config.TestsSource{Script: "echo 'no go toolchain' >&2; exit 1"}}}
	if v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, cfg); v != Fail || !strings.Contains(msg, "no go toolchain") {
		t.Fatalf("a failing source must fail naming its error, got %v: %s", v, msg)
	}
}

func TestSupportFilesAreNotTests(t *testing.T) {
	t.Run("INCHN-B18: A support file is neither run nor counted as naming a scenario", func(t *testing.T) {})
	if v, msg := checkTestsPass("", mapx.Node{ID: "utils/login.yaml", Kind: mapx.KindTest, Support: true}); v != Skip || !strings.Contains(msg, "support") {
		t.Fatalf("tests-pass must skip a support file saying why, got %v (%s)", v, msg)
	}
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root, g := rootWithTest(t, "const c = \"CREDX-B01\"\n")
	for i := range g.Nodes {
		if g.Nodes[i].ID == "credx_test.go" {
			g.Nodes[i].Support = true
		}
	}
	if _, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil); strings.Contains(msg, "ingest") {
		t.Fatalf("a support file must not count as a test naming the code, got %s", msg)
	}
}

func TestScenarioCoverage_FileTheSourceDoesNotDescribe(t *testing.T) {
	t.Run("INCHN-B17: With a tests source a scenario is written only when a title cites it", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root, g := rootWithTest(t, "name: 'CREDX-B01 - validates the limit'\n- launchApp\n")
	ts := &config.Config{Dialect: &config.Dialect{Family: "ts"}}
	if _, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, ts); !strings.Contains(msg, "ingest") {
		t.Fatalf("in a file the source lists no test in, the code in it counts as written, got %s", msg)
	}
}

// The `guide-checklist` recognises the compliance section in ANY language of the catalogue.
//
// The coupling this test locks: `anchors init` seeds HEADER_GUIDE.md with the title
// translated by the project's `lang:`, and the old regex matched only the Portuguese
// form. A `lang: en` project would be born failing the gate over a guide that init itself
// had just written — the framework charging what it does not produce.
//
// Without this test the regression is silent on both sides: whoever translates the guide
// without touching the gate breaks init; whoever narrows the gate back breaks the
// translated projects. Neither shows up in the build.
func TestChecklistHeading_recognisesEveryLanguage(t *testing.T) {
	t.Run("INCHN-I03: The compliance ruler is recognised in every language of the catalogue", func(t *testing.T) {})
	for _, title := range []string{
		"Pontos de conformidade", // pt-BR
		"Compliance points",      // en
		"Puntos de conformidad",  // es
	} {
		if !checklistHeadingRE.MatchString("## " + title + "\n\n- CK1 item\n") {
			t.Errorf("the gate did not recognise the section %q — a project in that language fails over a guide init wrote", title)
		}
	}
	// And it still refuses what is NOT the section: accepting every language must not
	// become accepting any title.
	if checklistHeadingRE.MatchString("## Outra coisa\n\n- CK1 item\n") {
		t.Error("the gate accepted a title that is not the compliance section")
	}
}

// The letter vocabulary belongs to the PROJECT (`rule_types`). A code regex pinned to the
// canonical letters makes EVERY scenario of a declared letter INVISIBLE — and the gate
// reports green over what it did not check, the most dangerous failure mode. It really
// happened with the letter `I` (Invariant): the scenario existed, the gate did not see
// it, and nobody noticed.
func TestRuleLetters_aLetterTheProjectDeclaresIsSeen(t *testing.T) {
	t.Run("INCHN-B08: Setting the rule letters reconfigures every dependent pattern together", func(t *testing.T) {})
	defer SetRuleLetters(config.DefaultRuleLetters) // do not leak into other tests

	feature := "@ABCDX-P01 @nivel-unit\n  Cenário: a política vale sempre\n"

	// before declaring: the letter is not in the vocabulary, so it is not seen — correct.
	SetRuleLetters(config.DefaultRuleLetters)
	if got := parseFeatureScenarios(feature); len(got) != 0 {
		t.Fatalf("an UNdeclared letter should not be recognised, got %+v", got)
	}

	// after declaring: it is seen.
	cfg := &config.Config{RuleTypes: []config.RuleType{
		{Letter: "B", Term: "Behavior"}, {Letter: "P", Term: "Policy"},
	}}
	SetRuleLetters(cfg.RuleLetters())
	got := parseFeatureScenarios(feature)
	if len(got) != 1 || got[0].Code != "ABCDX-P01" {
		t.Fatalf("the scenario of the declared letter is still invisible: %+v", got)
	}
}

// `non-empty` counts SCENARIOS for a feature, and a scenario opens in any language of the
// official Gherkin table — including the synonyms (`Example:` for `Scenario:`).
//
// The first version reused a regex with five keywords nailed in Portuguese and English,
// and it failed valid features on a BLOCKING gate: a Spanish `Escenario:` and an English
// `Rule:` + `Example:` both came out as "feature with no scenario".
func TestNonEmpty_featureScenarioInAnyGherkinLanguage(t *testing.T) {
	t.Run("INCHN-B19: A feature scenario is recognised in any Gherkin language", func(t *testing.T) {})
	n := mapx.Node{Kind: mapx.KindFeature}
	pass := map[string]string{
		"pt":            "Funcionalidade: X\n\n  Cenário: a\n",
		"pt outline":    "Funcionalidade: X\n\n  Esquema do Cenário: a\n",
		"en outline":    "Feature: X\n\n  Scenario Outline: a\n",
		"en example":    "Feature: X\n  Rule: r\n    Example: a\n",
		"es":            "Característica: X\n\n  Escenario: a\n",
		"fr":            "Fonctionnalité: X\n\n  Scénario: a\n",
		"pt unaccented": "Funcionalidade: X\n\n  Cenario: a\n",
	}
	for name, content := range pass {
		if v, msg := checkNonEmpty(content, n); v != Pass {
			t.Errorf("%s: a valid scenario was not recognised (%v: %s)", name, v, msg)
		}
	}

	fail := map[string]string{
		// `Examples:` is the data table of an outline, not a scenario: the `:` right after
		// the keyword keeps `Example` from matching it.
		"examples table only": "Feature: X\n\n  Examples:\n    | a |\n",
		"header only":         "Funcionalidade: X\n# nothing else\n",
	}
	for name, content := range fail {
		if v, _ := checkNonEmpty(content, n); v != Fail {
			t.Errorf("%s: a feature with no scenario passed", name)
		}
	}
}

// The gate charges the requirements the spec DEFINES, not the ones it CITES.
//
// Measured in the reference app: `GoLiveChecklist` defines 6 requirements and the gate
// charged 18 scenarios — 15 of them from other units (`CRPNC-B03`, `MTTLM-B02`,
// `DTSTD-B06`…), cited in the prose while justifying its rules.
//
// None of those scenarios could be proven by a test of this unit: the gate asked for the
// impossible, and the message said the spec was poorly covered. And the side effect is
// worse than the noise — a gate that always fails is a gate one learns to ignore.
func TestScenarioCoverage_doesNotChargeWhatTheSpecOnlyCites(t *testing.T) {
	t.Run("INCHN-B20: Scenario coverage charges what the spec defines, not what it cites", func(t *testing.T) {})
	content := `# GoLiveChecklist

## Regras

### GLCGL-B01 — cada item tem um artefato

É o que o ` + "`PLTFR`" + ` estabelece, e a ` + "`DTSTD-B06`" + ` confirma para
infraestrutura. A ` + "`CRPNC-B03`" + ` usa o mesmo raciocínio na rotação.

### GLCGL-B02 — as dívidas têm estado atual

O ` + "`MTTLM-B02`" + ` nasce desligado, e a ` + "`CRPNC-B06`" + ` o liga.
`
	n := mapx.Node{
		ID:   "packages/infra/GoLiveChecklist.spec.md",
		Code: "GLCGL",
		Signal: &mapx.TestSignal{
			ProvenCodes: []string{"GLCGL-B01", "GLCGL-B02"},
			AtRev:       "abc",
		},
		Rev: "abc",
	}

	v, msg := checkScenarioCoverage(content, n, "", nil, nil)

	if v != Pass {
		t.Errorf("verdict = %v — %s\n  both DEFINED requirements are proven; the "+
			"rest is citation", v, msg)
	}
	for _, cited := range []string{"DTSTD-B06", "CRPNC-B03", "MTTLM-B02", "CRPNC-B06"} {
		if strings.Contains(msg, cited) {
			t.Errorf("the gate charged %q, which this spec only CITES", cited)
		}
	}
}

// And the DEFINED requirement with no proven scenario is still charged — the fix must
// not have switched the gate off.
func TestScenarioCoverage_stillChargesTheDefinedRequirement(t *testing.T) {
	t.Run("INCHN-B20: Scenario coverage charges what the spec defines, not what it cites", func(t *testing.T) {})
	content := "### ABCDX-B01 — a regra\n\nCorpo.\n\n### ABCDX-B02 — outra\n\nCorpo.\n"
	n := mapx.Node{
		Code:   "ABCDX",
		Rev:    "r1",
		Signal: &mapx.TestSignal{ProvenCodes: []string{"ABCDX-B01"}, AtRev: "r1"},
	}

	v, msg := checkScenarioCoverage(content, n, "", nil, nil)

	if v != Fail {
		t.Fatalf("verdict = %v — `ABCDX-B02` has no proven scenario", v)
	}
	if !strings.Contains(msg, "ABCDX-B02") {
		t.Errorf("the message does not name the uncovered requirement: %s", msg)
	}
}

// THE TWO QUESTIONS of `scenario-coverage` — the same ruler as `flag-covered`.
//
// The earlier version asked only "did it pass?", and on a project that never ingested a
// report it answered Pending for everything: "nobody measured", hiding the scenarios
// nobody tested. And the static question alone would be the opposite error — a written
// test may never have run.
//
// Each state has a different fix, which is why the verdict must keep them apart: "no
// test" asks someone to write one; "written and not run" asks someone to run it.

const specWithTwoRequirements = "### CREDX-B01 — validates the limit\n\n### CREDX-B02 — refuses the balance\n"

func specNodeCoverage() mapx.Node {
	return mapx.Node{ID: "credx.spec.md", Kind: mapx.KindSpec, Code: "CREDX"}
}

func rootWithTest(t *testing.T, body string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "credx_test.go"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, &mapx.Graph{Nodes: []mapx.Node{{ID: "credx_test.go", Kind: mapx.KindTest}}}
}

// NO TEST AT ALL: it is reported even without ingested execution. It is the case the
// earlier version hid behind a Pending.
func TestScenarioCoverage_noTestIsReportedEvenWithoutIngestion(t *testing.T) {
	t.Run("INCHN-B21: Scenario coverage tells a missing test apart from a test never run", func(t *testing.T) {})
	root, g := rootWithTest(t, "func TestNothing(t *testing.T) {}\n")

	v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil)
	if v != Fail {
		t.Fatalf("no test names the requirements — expected Fail, got %v", v)
	}
	for _, c := range []string{"CREDX-B01", "CREDX-B02"} {
		if !strings.Contains(msg, c) {
			t.Errorf("the verdict does not name %s: %q", c, msg)
		}
	}
}

// WRITTEN AND NEVER RUN: the verdict has to say which of the two problems it is.
func TestScenarioCoverage_writtenButNotRunSaysWhichOfTheTwo(t *testing.T) {
	t.Run("INCHN-B21: Scenario coverage tells a missing test apart from a test never run", func(t *testing.T) {})
	root, g := rootWithTest(t, "const c = \"CREDX-B01\"\nfunc TestX(t *testing.T) {}\n")

	v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil)
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(msg, "ingest") {
		t.Errorf("the verdict does not tell written-but-not-run apart: %q", msg)
	}
	// B02 still has no test at all, and both states appear in the same verdict.
	if !strings.Contains(msg, "CREDX-B02") {
		t.Errorf("the verdict lost the requirement with no test: %q", msg)
	}
}

// A CODE IN A COMMENT does not count — a citation is a reference, not an implementation.
func TestScenarioCoverage_commentDoesNotCountAsATest(t *testing.T) {
	t.Run("INCHN-B17: With a tests source a scenario is written only when a title cites it", func(t *testing.T) {})
	root, g := rootWithTest(t, "// CREDX-B01 is covered elsewhere\nfunc TestX(t *testing.T) {}\n")

	_, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil)
	if strings.Contains(msg, "ingest") {
		t.Errorf("a code only in a COMMENT passed as a written test: %q", msg)
	}
}

// And what EXECUTION proved leaves the accusation, which is the original behaviour.
func TestScenarioCoverage_provenPasses(t *testing.T) {
	t.Run("INCHN-B21: Scenario coverage tells a missing test apart from a test never run", func(t *testing.T) {})
	root, _ := rootWithTest(t, "func TestNothing(t *testing.T) {}\n")
	n := specNodeCoverage()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CREDX-B01", "CREDX-B02"}}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "credx_test.go", Kind: mapx.KindTest}}}

	if v, msg := checkScenarioCoverage(specWithTwoRequirements, n, root, g, nil); v != Pass {
		t.Errorf("both requirements proven and the verdict was %v: %s", v, msg)
	}
}

// A LAYER THAT DISPENSES `tested-by` has no tests by declaration, and `scenario-coverage`
// honours it as `triad-complete` does. In the reference app every schema-model spec failed "scenario with
// no green test" although the Structure says those tests do not exist. A layer without the
// opt-out is still charged.
func TestScenarioCoverage_honoursTheLayersTestedByOptOut(t *testing.T) {
	t.Run("INCHN-B22: Scenario coverage honours a layer that dispenses tested-by", func(t *testing.T) {})
	root, g := rootWithTest(t, "package credx\n")
	cfg := &config.Config{Layers: map[string]config.Layer{
		"schema-model": {Kind: "spec", OptionalTriadEdges: []string{"covered-by", "tested-by"}},
		"service":      {Kind: "spec"},
	}}
	optedOut := specNodeCoverage()
	optedOut.Tags = []string{"schema-model"}
	if v, msg := checkScenarioCoverage(specWithTwoRequirements, optedOut, root, g, cfg); v != Skip || !strings.Contains(msg, "tested-by") {
		t.Errorf("a layer dispensing tested-by must skip and say why, got %v: %s", v, msg)
	}
	charged := specNodeCoverage()
	charged.Tags = []string{"service"}
	if v, _ := checkScenarioCoverage(specWithTwoRequirements, charged, root, g, cfg); v != Fail {
		t.Errorf("a layer without the opt-out is still charged, got %v", v)
	}
}

// The mutation gate is only worth something if it tells three situations apart; a gate
// that always passes (or always fails) says nothing. Each case below pins one of them.
func TestMutationScore(t *testing.T) {
	t.Run("INCHN-B23: Mutation score passes at the threshold and fails below it naming the survivors", func(t *testing.T) {})
	t.Run("INCHN-B24: A missing or stale mutation signal is pending", func(t *testing.T) {})
	i18n.Set("pt-BR")
	t.Cleanup(func() { i18n.Set(i18n.Default) })
	cases := []struct {
		name     string
		sig      *mapx.TestSignal
		expected Verdict
		contains string
	}{
		{"no signal ingested → Pending, saying what is missing and what is lost",
			nil, Pending, "ingest --mutation"},
		{"score above the threshold → passes",
			&mapx.TestSignal{MutantsKilled: 90, MutantsSurvived: 5, MutationScore: 94.7}, Pass, ""},
		{"too many survivors → fails naming how many",
			&mapx.TestSignal{MutantsKilled: 5, MutantsSurvived: 15, MutationScore: 25}, Fail, "15 mutante(s) sobreviveram"},
		{"exact threshold → passes (the limit does not fail)",
			&mapx.TestSignal{MutantsKilled: 7, MutantsSurvived: 3, MutationScore: 70}, Pass, ""},
		// The tool RAN and ignored everything — a table of constants with `ignoreStatic`.
		// Nothing survived, so it is 100 and the verdict is Pass. Before this the gate said
		// "run the mutation tool" about a file it had already run on: a request running it
		// again would not satisfy, and the kind of noise that teaches people to ignore the gate.
		{"everything ignored → passes, without asking for a new run",
			&mapx.TestSignal{MutantsIgnored: 12, MutationScore: 100}, Pass, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			n := mapx.Node{Kind: mapx.KindCode, Rev: "r1", Signal: c.sig}
			if c.sig != nil {
				c.sig.AtRev = "r1" // fresh signal; the stale one has its own case
			}
			v, d := checkMutationScore("", n)
			if v != c.expected {
				t.Fatalf("verdict = %s, want %s (detail: %s)", v, c.expected, d)
			}
			if c.contains != "" && !strings.Contains(d, c.contains) {
				t.Fatalf("detail %q does not mention %q", d, c.contains)
			}
		})
	}
}

// A mutation signal measured at an EARLIER version of the file proves nothing about the
// current one — the same rule as the other ingested signals, and the most misleading: the
// number looks good.
func TestMutationScoreStale(t *testing.T) {
	t.Run("INCHN-B24: A missing or stale mutation signal is pending", func(t *testing.T) {})
	n := mapx.Node{Kind: mapx.KindCode, Rev: "r2",
		Signal: &mapx.TestSignal{MutantsKilled: 100, MutationScore: 100, AtRev: "r1"}}
	v, d := checkMutationScore("", n)
	if v != Pending {
		t.Fatalf("a perfect score from an old revision should be Pending, was %s", v)
	}
	if !strings.Contains(d, "stale") {
		t.Fatalf("the detail does not explain the staleness: %q", d)
	}
	covered := mapx.Node{Kind: mapx.KindCode, Rev: "r2",
		Signal: &mapx.TestSignal{MutantsKilled: 100, MutationScore: 100, AtRev: "r2", MutationAtRev: "r1"}}
	if v, _ := checkMutationScore("", covered); v != Pending {
		t.Errorf("a coverage ingestion at the new rev does not make an old mutation current, got %s", v)
	}
}

func nodeWithScore(score, low, high float64, survived int) mapx.Node {
	return mapx.Node{
		ID:   "src/regra.ts",
		Kind: mapx.KindCode,
		Signal: &mapx.TestSignal{
			MutantsKilled:   100,
			MutantsSurvived: survived,
			MutationScore:   score,
			MutationLow:     low,
			MutationHigh:    high,
		},
	}
}

// TestMutationRange_belowAcceptableFails — the bottom range is the only one that bars.
func TestMutationRange_belowAcceptableFails(t *testing.T) {
	t.Run("INCHN-B23: Mutation score passes at the threshold and fails below it naming the survivors", func(t *testing.T) {})
	t.Run("INCHN-B25: A score between acceptable and desirable is pending, not failed", func(t *testing.T) {})
	v, detail := checkMutationScore("", nodeWithScore(56, 70, 90, 142))
	if v != Fail {
		t.Fatalf("56%% with a minimum of 70%% must fail; got %v", v)
	}
	if !strings.Contains(detail, "70") {
		t.Errorf("the report must say against which ruler it failed: %q", detail)
	}
}

// TestMutationRange_betweenAcceptableAndDesirableDoesNotBar is the new concept: it passed,
// and it still shows. If this became Fail it would be a threshold of 90 in disguise — and
// the distinction between "must not" and "could be better" would be lost.
func TestMutationRange_betweenAcceptableAndDesirableDoesNotBar(t *testing.T) {
	t.Run("INCHN-B25: A score between acceptable and desirable is pending, not failed", func(t *testing.T) {})
	i18n.Set("pt-BR")
	t.Cleanup(func() { i18n.Set(i18n.Default) })
	v, detail := checkMutationScore("", nodeWithScore(75, 70, 90, 30))
	if v == Fail {
		t.Fatalf("75%% is above the acceptable (70%%) — it must not fail")
	}
	if v != Pending {
		t.Fatalf("the middle range has to SHOW (Pending), not vanish; got %v", v)
	}
	if !strings.Contains(detail, "aceitável") || !strings.Contains(detail, "desejável") {
		t.Errorf("the report must name both ranges: %q", detail)
	}
	if !strings.Contains(detail, "15") {
		t.Errorf("saying HOW MUCH is left is what makes the warning actionable: %q", detail)
	}
}

// TestMutationRange_aboveDesirablePassesClean — no noise for whoever already got there.
func TestMutationRange_aboveDesirablePassesClean(t *testing.T) {
	t.Run("INCHN-B25: A score between acceptable and desirable is pending, not failed", func(t *testing.T) {})
	v, detail := checkMutationScore("", nodeWithScore(92, 70, 90, 5))
	if v != Pass || detail != "" {
		t.Errorf("92%% with a desirable of 90%% passes clean; got %v %q", v, detail)
	}
}

// TestMutationRange_withoutDesirableFallsBackToOneThreshold — no project is forced to
// adopt the concept. With no `high` in the report, the gate behaves as before.
func TestMutationRange_withoutDesirableFallsBackToOneThreshold(t *testing.T) {
	t.Run("INCHN-B25: A score between acceptable and desirable is pending, not failed", func(t *testing.T) {})
	if v, _ := checkMutationScore("", nodeWithScore(75, 70, 0, 30)); v != Pass {
		t.Errorf("with no desirable declared, 75%% above the minimum passes clean; got %v", v)
	}
}

// TestMutationRange_rulerComesFromTheReportNotTheEngine — the project that declares 60 as
// acceptable has 65 approved, even though the engine default is 70. It is what keeps the
// framework from deciding what is good enough quality for everyone.
func TestMutationRange_rulerComesFromTheReportNotTheEngine(t *testing.T) {
	t.Run("INCHN-B26: The mutation thresholds come from the report", func(t *testing.T) {})
	if v, _ := checkMutationScore("", nodeWithScore(65, 60, 0, 40)); v != Pass {
		t.Errorf("with a minimum of 60 declared, 65 passes; got %v", v)
	}
	if v, _ := checkMutationScore("", nodeWithScore(65, 0, 0, 40)); v != Fail {
		t.Errorf("with no ruler in the report, the default 70 holds and 65 fails; got %v", v)
	}
}

// TestMutationRange_invalidDesirableIsIgnored — `high` below `low` is a project
// misconfiguration; the gate must not turn it into an impossible range that always fails.
func TestMutationRange_invalidDesirableIsIgnored(t *testing.T) {
	t.Run("INCHN-B25: A score between acceptable and desirable is pending, not failed", func(t *testing.T) {})
	if v, _ := checkMutationScore("", nodeWithScore(75, 70, 50, 30)); v != Pass {
		t.Errorf("a desirable below the acceptable is incoherent and must be ignored; got %v", v)
	}
}

// ── an old scope does not decide the verdict ───────────────────────────────────

func nodeWithRevScopes(rev string, iso, full mapx.MutationScope, total float64) mapx.Node {
	return mapx.Node{
		ID: "src/regra.ts", Kind: mapx.KindCode, Rev: rev,
		Signal: &mapx.TestSignal{
			MutantsKilled: 100, MutantsSurvived: 50,
			MutationScore: total, MutationLow: 70,
			AtRev:           rev,
			MutationByScope: map[string]mapx.MutationScope{"isolated": iso, "full": full},
		},
	}
}

// TestMutationScope_oldScopeDoesNotDecideTheVerdict — the gate judges by the ISOLATED
// scope. With one stamp for the whole signal, re-ingesting only `full` renewed the stamp
// and the old isolated score rode along as if current: the verdict came from a number
// measured against code that had already changed.
func TestMutationScope_oldScopeDoesNotDecideTheVerdict(t *testing.T) {
	t.Run("INCHN-B28: A scope measured at an older revision decides nothing", func(t *testing.T) {})
	n := nodeWithRevScopes("rev2",
		mapx.MutationScope{Score: 30, Survived: 200, AtRev: "rev1"}, // measured before
		mapx.MutationScope{Score: 95, Survived: 3, AtRev: "rev2"},   // current
		95)
	v, detail := checkMutationScore("", n)
	if v == Fail {
		t.Fatalf("the isolated score of rev1 must not fail rev2; report: %q", detail)
	}
	if strings.Contains(detail, "30") {
		t.Errorf("the old number must not appear in the report: %q", detail)
	}
}

// TestMutationScope_currentScopeStillDecides — the guard must not switch the pair off when
// both are at the same rev; otherwise the coupling finding (low isolated, high full) vanishes.
func TestMutationScope_currentScopeStillDecides(t *testing.T) {
	t.Run("INCHN-B28: A scope measured at an older revision decides nothing", func(t *testing.T) {})
	n := nodeWithRevScopes("rev2",
		mapx.MutationScope{Score: 30, Survived: 200, AtRev: "rev2"},
		mapx.MutationScope{Score: 95, Survived: 3, AtRev: "rev2"},
		95)
	if v, _ := checkMutationScore("", n); v != Fail {
		t.Errorf("an isolated 30%% at the current rev must fail; got %v", v)
	}
}

// TestMutationScope_oldSignalWithoutScopeStampStillCounts — a signal recorded before the
// field existed has no per-scope AtRev. Treating it as old would make the gate ask to
// re-measure everything already in the map, with no basis to claim it is out of date.
func TestMutationScope_oldSignalWithoutScopeStampStillCounts(t *testing.T) {
	t.Run("INCHN-B28: A scope measured at an older revision decides nothing", func(t *testing.T) {})
	n := nodeWithRevScopes("rev2",
		mapx.MutationScope{Score: 30, Survived: 200},
		mapx.MutationScope{Score: 95, Survived: 3},
		95)
	if v, _ := checkMutationScore("", n); v != Fail {
		t.Errorf("with no scope stamp, the pair still holds; got %v", v)
	}
}

// TestMutationScope_reportDoesNotContradictItself — the sentence about the delta was
// concatenated unconditionally, so a LOW delta produced "the two scopes agree … A high
// delta means …" in the same report. A report that contradicts itself is not read: whoever
// reads it stops trusting it.
func TestMutationScope_reportDoesNotContradictItself(t *testing.T) {
	t.Run("INCHN-B27: The verdict follows the isolated scope and the report reads the delta", func(t *testing.T) {})
	n := nodeWithRevScopes("r",
		mapx.MutationScope{Score: 58, Survived: 175, AtRev: "r"},
		mapx.MutationScope{Score: 61, Survived: 170, AtRev: "r"},
		61)
	_, detail := checkMutationScore("", n)
	if strings.Contains(detail, "concordam") && strings.Contains(detail, "Delta alto") {
		t.Errorf("contradictory report: %q", detail)
	}
}

func nodeWithScopes(iso, full mapx.MutationScope) mapx.Node {
	return mapx.Node{
		Kind: mapx.KindCode, ID: "x.ts", Rev: "r1",
		Signal: &mapx.TestSignal{
			MutantsKilled: full.Killed, MutantsSurvived: full.Survived,
			MutationScore: full.Score, AtRev: "r1",
			MutationByScope: map[string]mapx.MutationScope{"isolated": iso, "full": full},
		},
	}
}

// The real case that motivated the change: a UI atom at 8% isolated and 77% full. Looking
// only at the full score it seemed healthy — and 92% of the mutants survive its own tests.
func TestMutationScope_scopesRevealCoupling(t *testing.T) {
	t.Run("INCHN-B27: The verdict follows the isolated scope and the report reads the delta", func(t *testing.T) {})
	i18n.Set("pt-BR")
	t.Cleanup(func() { i18n.Set(i18n.Default) })
	n := nodeWithScopes(
		mapx.MutationScope{Killed: 7, Survived: 81, Score: 8},
		mapx.MutationScope{Killed: 68, Survived: 20, Score: 77})

	v, msg := checkMutationScore("", n)
	if v != Fail {
		t.Fatalf("verdict %v, want Fail — the isolated score is 8%%", v)
	}
	for _, want := range []string{"isolado 8%", "completo 77%", "delta 69p", "dependentes"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message without %q:\n%s", want, msg)
		}
	}
}

// The verdict is about the ISOLATED score: a high full score does not save a unit that
// does not prove itself alone. It is the difference between "someone proves it" and "this
// unit's test proves it".
func TestMutationScope_verdictFollowsIsolatedNotFull(t *testing.T) {
	t.Run("INCHN-B27: The verdict follows the isolated scope and the report reads the delta", func(t *testing.T) {})
	// full at 100%, isolated at 30% → still fails
	n := nodeWithScopes(
		mapx.MutationScope{Killed: 3, Survived: 7, Score: 30},
		mapx.MutationScope{Killed: 10, Survived: 0, Score: 100})
	if v, _ := checkMutationScore("", n); v != Fail {
		t.Errorf("verdict %v, want Fail — a high full score does not make up for a low isolated one", v)
	}
}

// A unit that proves itself alone passes, even if the full score is higher.
func TestMutationScope_isolatedAboveThresholdPasses(t *testing.T) {
	t.Run("INCHN-B27: The verdict follows the isolated scope and the report reads the delta", func(t *testing.T) {})
	n := nodeWithScopes(
		mapx.MutationScope{Killed: 8, Survived: 2, Score: 80},
		mapx.MutationScope{Killed: 9, Survived: 1, Score: 90})
	if v, msg := checkMutationScore("", n); v != Pass {
		t.Errorf("verdict %v (%s), want Pass", v, msg)
	}
}

// A low delta with a low score is another diagnosis: it is not coupling, it is a missing
// assertion — and the message must say so, or the author looks in the wrong place.
func TestMutationScope_lowDeltaPointsAtAssertionNotCoupling(t *testing.T) {
	t.Run("INCHN-B27: The verdict follows the isolated scope and the report reads the delta", func(t *testing.T) {})
	i18n.Set("pt-BR")
	t.Cleanup(func() { i18n.Set(i18n.Default) })
	n := nodeWithScopes(
		mapx.MutationScope{Killed: 2, Survived: 8, Score: 20},
		mapx.MutationScope{Killed: 3, Survived: 7, Score: 25})
	_, msg := checkMutationScore("", n)
	if !strings.Contains(msg, "asserção") {
		t.Errorf("the message does not tell assertion apart from coupling:\n%s", msg)
	}
}

// Backward compatible: whoever ingests without scopes keeps the old ruler over the total.
func TestMutationScope_withoutScopesUsesTheTotal(t *testing.T) {
	t.Run("INCHN-B27: The verdict follows the isolated scope and the report reads the delta", func(t *testing.T) {})
	n := mapx.Node{Kind: mapx.KindCode, ID: "x.ts", Rev: "r1",
		Signal: &mapx.TestSignal{MutantsKilled: 8, MutantsSurvived: 2, MutationScore: 80, AtRev: "r1"}}
	if v, msg := checkMutationScore("", n); v != Pass {
		t.Errorf("verdict %v (%s), want Pass without scopes", v, msg)
	}
}

// The worst message gitmeta's silence produced: with no REPOSITORY, the gate fell into the
// "committed" branch, `LastCommitDate` failed, and the verdict came out as
// "arquivo sem commit no git (novo/untracked)" — a false and SPECIFIC claim about the
// file, which sent the author to investigate exactly where the problem is not.
func TestUpdatedAt_withoutRepoDoesNotBlameTheFile(t *testing.T) {
	t.Run("INCHN-B29: Without a repository updated-at skips instead of blaming the file", func(t *testing.T) {})
	dir := t.TempDir()
	if gitmeta.Check(dir) == gitmeta.Disponível {
		t.Skipf("the temporary directory %s is inside a git repo", dir)
	}
	content := "// @anchors\n// updated_at: 2026-01-01\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	v, d := checkUpdatedAt(content, mapx.Node{ID: "a.go"}, dir)

	if v != Skip {
		t.Fatalf("with no repository there is no verdict to give about the date, got %s (%s)", v, d)
	}
	if strings.Contains(d, "novo/untracked") {
		t.Errorf("blames the FILE for something the REPOSITORY lacks: %s", d)
	}
	if !strings.Contains(d, "reposit") {
		t.Errorf("the message must name the real cause: %s", d)
	}
}

// The counterpart: in a real repo, a new uncommitted file dated today still passes — the
// gate's rule did not change.
func TestUpdatedAt_withRepoStillChecksTheDate(t *testing.T) {
	t.Run("INCHN-B29: Without a repository updated-at skips instead of blaming the file", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("setup: %s", out)
	}
	today := gitmeta.Today()
	content := "// @anchors\n// updated_at: " + today + "\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if v, d := checkUpdatedAt(content, mapx.Node{ID: "a.go"}, dir); v != Pass {
		t.Fatalf("a file being edited dated today should pass, got %s (%s)", v, d)
	}

	// And the wrong date still fails.
	wrong := "// @anchors\n// updated_at: 2020-01-01\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(wrong), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := checkUpdatedAt(wrong, mapx.Node{ID: "a.go"}, dir); v != Fail {
		t.Fatalf("a wrong date in a modified file should fail, got %s", v)
	}
}

// specInEnglish has the catalogued rule (what the gate already charged) and the section
// titles in ENGLISH — the case `enforce_section_language` exists to catch in a pt-BR project.
const specInEnglish = `<!-- @anchors
  code: LGNOI
-->
# Login

## Overview
Entra no app.

## Rules

### LGNOI-B01 — regra
Comportamento.
`

const specInPortuguese = `<!-- @anchors
  code: LGNOI
-->
# Login

## Visão Geral
Entra no app.

## Regras

### LGNOI-B01 — regra
Comportamento.
`

func sectionsGateConfig(enforce *bool) *config.Config {
	return &config.Config{
		Lang:  "pt-BR",
		Gates: []config.Gate{{Name: "spec-complete", Check: "spec-sections", EnforceSectionLanguage: enforce}},
	}
}

// The DEFAULT is to charge: omitting the key must not mean "do not check", or the mixed
// collection is born in silence — the very defect the option exists to make visible.
func TestSpecSections_wrongLanguageFailsByDefault(t *testing.T) {
	t.Run("INCHN-B30: A section title in another language fails unless the gate waives it", func(t *testing.T) {})
	v, msg := checkSpecSections(specInEnglish, mapx.Node{}, "", nil, sectionsGateConfig(nil))
	if v != Fail {
		t.Fatalf("verdict = %v, want Fail: the spec is in English in a pt-BR project", v)
	}
	if !strings.Contains(msg, "Visão Geral") {
		t.Errorf("the message should give the EXPECTED title, so the fix is obvious; got: %s", msg)
	}
}

func TestSpecSections_rightLanguagePasses(t *testing.T) {
	t.Run("INCHN-B30: A section title in another language fails unless the gate waives it", func(t *testing.T) {})
	if v, msg := checkSpecSections(specInPortuguese, mapx.Node{}, "", nil, sectionsGateConfig(nil)); v != Pass {
		t.Fatalf("verdict = %v (%s), want Pass", v, msg)
	}
}

// Switching it off is a legitimate, declared decision: a migration under way, or a
// bilingual project.
func TestSpecSections_enforceFalseDoesNotChargeLanguage(t *testing.T) {
	t.Run("INCHN-B30: A section title in another language fails unless the gate waives it", func(t *testing.T) {})
	off := false
	if v, msg := checkSpecSections(specInEnglish, mapx.Node{}, "", nil, sectionsGateConfig(&off)); v != Pass {
		t.Fatalf("verdict = %v (%s), want Pass with enforce_section_language: false", v, msg)
	}
}

// The project's OWN lexicon is not a wrong language — it is a section the framework does
// not name. Accusing it would impose the engine's vocabulary on the project.
func TestSpecSections_titleOutsideTheCatalogueIsNotWrongLanguage(t *testing.T) {
	t.Run("INCHN-B30: A section title in another language fails unless the gate waives it", func(t *testing.T) {})
	spec := specInPortuguese + "\n## Fora de escopo\nNada.\n\n## Decisões em aberto\nnenhuma\n"
	if v, msg := checkSpecSections(spec, mapx.Node{}, "", nil, sectionsGateConfig(nil)); v != Pass {
		t.Fatalf("verdict = %v (%s), want Pass: the project's own titles are not a wrong language", v, msg)
	}
}

// Without config the gate cannot invent a language: it charges only what it always charged.
func TestSpecSections_withoutConfigDoesNotChargeLanguage(t *testing.T) {
	t.Run("INCHN-B30: A section title in another language fails unless the gate waives it", func(t *testing.T) {})
	if v, msg := checkSpecSections(specInEnglish, mapx.Node{}, "", nil, nil); v != Pass {
		t.Fatalf("verdict = %v (%s), want Pass without config", v, msg)
	}
}

// These tests are the META-GATE of the `anchors new` templates: they ensure the skeleton
// the command emits is BORN CONFORMING — it passes the SAME gate functions `check` runs
// (checkHeaderConforms, checkSpecSections). Without them, the template could drift from
// the ruler and nobody would notice (the template is not in the project's graph).
//
// The strings below are the canonical output of `new` (default) per kind — keeping them in
// sync with cmd/anchors/new_templates.go is the contract. If the gate changes its ruler,
// these tests break and force updating the template along with it.

func newTemplateHeaderNode(kind mapx.Kind) mapx.Node {
	// a generic governed node (not a recognised layer) → the header demands code|ref.
	return mapx.Node{ID: "x/Login." + string(kind), Kind: kind}
}

func TestNewTemplate_specIsBornConforming(t *testing.T) {
	t.Run("INCHN-B31: The skeletons anchors new emits are born conforming", func(t *testing.T) {})
	// mirrors specTemplate (default: title+overview+rules) rendered for "Login"/"LGNOX".
	spec := "<!-- @anchors\n  code: LGNOX\n  updated_at: TODO\n  layer: TODO\n-->\n" +
		"# Login — TODO propósito em uma frase\n\n> **Código**: `LGNOX`\n\n" +
		"## Visão Geral\nTODO: o que a unidade faz e para quem.\n\n" +
		"## Regras\n\n### LGNOX-B01 — TODO regra\nDescreva o comportamento (não a implementação).\n\n"

	if v, msg := checkHeaderConforms(spec, mapx.Node{ID: "x/Login.spec.md", Kind: "spec"}); v != Pass {
		t.Fatalf("the spec from `new` fails header-valid: %s", msg)
	}
	if v, msg := checkSpecSections(spec, mapx.Node{ID: "x/Login.spec.md"}, "", nil, nil); v != Pass {
		t.Fatalf("the spec from `new` fails spec-sections: %s", msg)
	}
}

func TestNewTemplate_featureIsBornConforming(t *testing.T) {
	t.Run("INCHN-B31: The skeletons anchors new emits are born conforming", func(t *testing.T) {})
	feat := "# language: pt\n# @anchors\n#   ref: LGNOX\n#   updated_at: TODO\n#   layer: feature\n" +
		"\n@LGNOX\nFuncionalidade: Login\n\n" +
		"  @LGNOX-B01 @nivel-unit @P2\n  Cenário: TODO\n    Dado TODO\n    Quando TODO\n    Então o efeito LGNOX-B01 se verifica\n\n"

	if v, msg := checkHeaderConforms(feat, newTemplateHeaderNode(mapx.KindFeature)); v != Pass {
		t.Fatalf("the feature from `new` fails header-valid: %s", msg)
	}
	// non-empty: the feature has content beyond the header.
	if strings.TrimSpace(strings.SplitN(feat, "feature\n", 2)[1]) == "" {
		t.Fatal("the feature from `new` is empty after the header")
	}
}

func TestNewTemplate_testIsBornConforming(t *testing.T) {
	t.Run("INCHN-B31: The skeletons anchors new emits are born conforming", func(t *testing.T) {})
	test := "// @anchors\n//   ref: LGNOX\n//   updated_at: TODO\n//   layer: test\n" +
		"\ndescribe('Login', () => {\n  it('[LGNOX-B01] TODO', () => {\n    // TODO\n  })\n})\n"

	if v, msg := checkHeaderConforms(test, newTemplateHeaderNode(mapx.KindTest)); v != Pass {
		t.Fatalf("the test from `new` fails header-valid: %s", msg)
	}
}

func TestMutationScore_UnderLoad(t *testing.T) {
	t.Run("INCHN-B32: A mutation score measured under load is not trusted", func(t *testing.T) {})
	node := func(killed, survived, timedOut int) mapx.Node {
		return mapx.Node{ID: "a.go", Kind: mapx.KindCode, Rev: "r", Signal: &mapx.TestSignal{AtRev: "r",
			MutantsKilled: killed, MutantsSurvived: survived, MutantsTimedOut: timedOut,
			MutationScore: 100 * float64(killed) / float64(killed+survived)}}
	}
	// 74 of 78 killed, 65 of them by the time limit: the reference app's measurement.
	v, msg := checkMutationScoreUnderLoad("", node(74, 4, 65), "", nil, nil)
	if v != Pending || !strings.Contains(msg, "65") || !strings.Contains(msg, "no time limit") {
		t.Fatalf("above the ceiling the score is pending, saying to measure with no time limit, got %v (%s)", v, msg)
	}
	// Exactly at the ceiling (20 of 100) the score still decides.
	if v, _ := checkMutationScoreUnderLoad("", node(95, 5, 19), "", nil, nil); v != Pass {
		t.Fatalf("below the ceiling the score decides, got %v", v)
	}
	if v, _ := checkMutationScoreUnderLoad("", node(95, 5, 20), "", nil, nil); v != Pass {
		t.Fatalf("at the ceiling the score still decides, got %v", v)
	}
	// A declared ceiling of 0.9 lets the same measurement decide.
	cfg := &config.Config{Gates: []config.Gate{{Name: "mutation-score", Check: "mutation-score", TimeoutCeiling: 0.9}}}
	if v, _ := checkMutationScoreUnderLoad("", node(74, 4, 65), "", nil, cfg); v != Pass {
		t.Fatalf("under a declared ceiling of 0.9 the score decides, got %v", v)
	}
	// Below the floor with few timeouts: fails, and says how many timed out.
	v, msg = checkMutationScoreUnderLoad("", node(5, 5, 1), "", nil, nil)
	if v != Fail || !strings.Contains(msg, "1 of the killed") {
		t.Fatalf("a failure says how many were killed by the time limit, got %v (%s)", v, msg)
	}
}
