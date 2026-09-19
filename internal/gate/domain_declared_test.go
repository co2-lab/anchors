package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaDominio(t *testing.T, spec string) (Verdict, string) {
	t.Helper()
	return checkDomainDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
}

const cabecalhoDominio = "## Domínio\n\n| Entrada | Aceita | Fora do domínio | Quem garante |\n| --- | --- | --- | --- |\n"

// A coluna QUEM GARANTE é a razão da seção existir. Sem ela, "fora do domínio" é só
// outra forma de dizer "não é meu problema" — e o problema não fica com ninguém. Foi
// assim que três specs declararam, cada uma corretamente, que não validavam a mesma
// entrada: o dever ficou órfão e a entrada inválida passou.
func TestAnEntryWithNoOwnerIsFailedAndTheVerdictNamesIt(t *testing.T) {
	t.Run("DMDCD-B04: An entry with no owner is failed, and the verdict names it", func(t *testing.T) {})
	t.Run("DMDCD-B07: An entry whose owner is named passes", func(t *testing.T) {})
	semDono := "# Spec\n\n" + cabecalhoDominio +
		"| `chave` | texto não-vazio | `__proto__` |  |\n"
	v, d := rodaDominio(t, semDono)
	if v != Fail {
		t.Fatalf("entrada sem dono deveria reprovar, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "chave") {
		t.Errorf("não nomeou a entrada órfã: %s", d)
	}

	comDono := "# Spec\n\n" + cabecalhoDominio +
		"| `chave` | texto não-vazio | `__proto__` | a interface, antes de chamar |\n"
	if v, d := rodaDominio(t, comDono); v != Pass {
		t.Fatalf("entrada COM dono deveria passar, foi %s (%s)", v, d)
	}
}

// A NÃO-RESPOSTA não conta como dono: é exatamente a frase que cria o órfão.
func TestANonAnswerInTheOwnerColumnIsNotAnOwner(t *testing.T) {
	t.Run("DMDCD-I02: A non-answer in the owner column is not an owner", func(t *testing.T) {})
	for _, naoDono := range []string{"ninguém", "n/a", "-", "—", "não valido", "não é meu", "TODO: decidir"} {
		t.Run(naoDono, func(t *testing.T) {
			spec := "# Spec\n\n" + cabecalhoDominio +
				"| `mês` | `YYYY-MM` | `2026-3` | " + naoDono + " |\n"
			if v, _ := rodaDominio(t, spec); v != Fail {
				t.Fatalf("%q não é dono — deveria reprovar, foi %s", naoDono, v)
			}
		})
	}
}

// ESTE TESTE FIXAVA O DEFEITO, e é por isso que ele sobrevive aqui invertido.
//
// Ele afirmava: "ausência da seção não é falha — o gate cobra quem ABRIU". Com isso
// verde, um gate `blocking: true` passou por 85 specs de um projeto real sem confrontar
// NENHUMA. Ninguém abre uma seção que nada cobra.
//
// A premissa não estava errada: exigir a seção de toda spec vira ritual, e unidade sem
// entrada externa não tem domínio a declarar. Errada estava a SAÍDA — silêncio em vez de
// dispensa declarada. O resto do vocabulário resolve isso com `@no-rule`/`@no-scenario`:
// a razão escrita registra que alguém olhou.
//
// Um teste que trava o comportamento de fuga do gate é pior que teste nenhum: ele faz a
// falha parecer decisão.
func TestASpecWithoutTheDomainSectionIsFailed(t *testing.T) {
	t.Run("DMDCD-B02: A spec without the domain section is failed", func(t *testing.T) {})
	v, msg := rodaDominio(t, "# Spec\n\n## Regras\n\n### AAAAX-B01 — x\n")
	if v != Fail {
		t.Fatalf("spec sem a seção e sem dispensa deveria reprovar, veio %v: %s", v, msg)
	}
}

// Seção aberta e vazia é pior que ausente: AFIRMA que se olhou e não se achou nada.
func TestASectionOpenedAndLeftEmptyIsFailed(t *testing.T) {
	t.Run("DMDCD-B03: A section opened and left empty is failed", func(t *testing.T) {})
	if v, d := rodaDominio(t, "# Spec\n\n"+cabecalhoDominio); v != Fail {
		t.Fatalf("seção só com cabeçalho deveria reprovar, foi %s (%s)", v, d)
	}
	// e a linha do molde não conta como declaração
	molde := "# Spec\n\n" + cabecalhoDominio + "| TODO | TODO | TODO | TODO |\n"
	if v, _ := rodaDominio(t, molde); v != Fail {
		t.Fatal("linha só com TODO não é declaração")
	}
}

// Prosa explicativa não é linha de dados — senão o texto de abertura viraria entrada
// fantasma e o autor aprenderia a não explicar nada.
func TestARowFilledOnlyWithPlaceholdersIsNotADeclaration(t *testing.T) {
	t.Run("DMDCD-B06: A row filled only with placeholders is not a declaration", func(t *testing.T) {})
	spec := "# Spec\n\n## Domínio\n\nEsta unidade recebe o histórico já carregado.\n\n" +
		"| Entrada | Aceita | Fora do domínio | Quem garante |\n| --- | --- | --- | --- |\n" +
		"| `versões` | de UMA chave | lista multi-chave | o chamador (`RDMDX-B03`) |\n"
	if v, d := rodaDominio(t, spec); v != Pass {
		t.Fatalf("prosa não é entrada, foi %s (%s)", v, d)
	}
}

func TestAnArtifactThatIsNotASpecLeavesWithoutAVerdict(t *testing.T) {
	t.Run("DMDCD-B01: An artifact that is not a spec leaves without a verdict", func(t *testing.T) {})
	spec := cabecalhoDominio + "| x | y | z |  |\n"
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindTest, mapx.KindFeature} {
		if v, _ := checkDomainDeclared(spec, mapx.Node{Kind: k}, "", nil, nil); v != Skip {
			t.Errorf("kind %s deveria ser Skip, foi %s", k, v)
		}
	}
}

// O CASO REALX que motivou o gate. A spec `MTVRX` declarava em `## Restrições`:
//
//	"MTVRX-X04 — Não valida conteúdo de chave nem de valor"
//
// e a spec do repositório declarava o mesmo; o modelo não mencionou. Cada uma correta.
// O dever ficou órfão, e a chave `__proto__` sumia do resultado sem erro — dado do
// usuário gravado e invisível, com todos os gates verdes.
//
// Este teste responde: a seção teria forçado o autor a VER? Repare que a não-resposta
// vem com CITAÇÃO ("não valido (MTVRX-X04)") — foi o que quase deixou o gate passar, e é
// como um autor honesto escreveria ao transportar a restrição para a coluna do dono.
func TestANonAnswerInTheOwnerColumnIsNotAnOwnerRealCase(t *testing.T) {
	t.Run("DMDCD-I02: A non-answer in the owner column is not an owner", func(t *testing.T) {})
	comoEstava := `# MTVRX

## Domínio

| Entrada | Aceita | Fora do domínio | Quem garante |
| --- | --- | --- | --- |
| ` + "`chave`" + ` | texto livre | — | não valido (MTVRX-X04) |
`
	v, d := checkDomainDeclared(comoEstava, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
	if v != Fail {
		t.Fatalf("transportar o \"não valido\" para a coluna do dono deveria REPROVAR. Foi %s (%s)", v, d)
	}
	t.Logf("gate diz: %s", d)

	comDono := strings.Replace(comoEstava, "não valido (MTVRX-X04)",
		"a interface de cadastro (KVEDX-V02) rejeita chave reservada", 1)
	if v, _ := checkDomainDeclared(comDono, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil); v != Pass {
		t.Fatal("com dono nomeado deveria passar")
	}
}

// A AUSENCIA DA SECAO era SILENCIO, e silencio nao e' dispensa.
//
// MEDIDO no blue-eyes: 85 specs, ZERO confrontadas -- `domain-declared` e' `blocking:
// true` e nao protegia nada, porque so cobrava quem tinha ABERTO a secao. Ninguem abriu.
//
// A justificativa tinha merito ("exigir a secao de toda spec transformaria instrumento
// em ritual"), mas escolheu a saida errada. O resto do vocabulario dispensa por
// DECLARACAO com razao -- `@no-rule`, `@no-scenario`, `@no-test` --, e essa forma
// registra que alguem OLHOU. O silencio nao distingue "nao tem entrada externa" de
// "ninguem pensou no assunto".
//
// O proprio gate ja fazia essa distincao tres linhas abaixo, sobre a secao vazia: "uma
// secao vazia AFIRMA que se olhou e nao se achou nada a declarar, que e' diferente de
// nao ter olhado". Reconhecia a diferenca e nao a aplicava ao caso principal.
//
// E foi esse silencio que deixou passar o defeito real: o `QueryScope` define o conjunto
// fechado de janelas e recebe entrada de fora -- caso central do gate. Nao abriu a secao,
// o gate calou, e as telas declararam `5m`, `30m` e `1d`, que o contrato nao aceita.
func TestASpecWithoutTheDomainSectionIsFailedAndNamesTheWaiver(t *testing.T) {
	t.Run("DMDCD-B02: A spec without the domain section is failed", func(t *testing.T) {})
	semNada := "# U\n\n## Regras\n\n### UUUUU-B01 — algo\n"
	v, msg := rodaDominio(t, semNada)
	if v == Skip || v == Pass {
		t.Errorf("a spec nao diz NADA sobre dominio e o gate a deixou passar (%v) — "+
			"silencio nao distingue `nao tem entrada externa` de `ninguem pensou`: %s", v, msg)
	}
	if !strings.Contains(msg, "@no-domain") {
		t.Errorf("a mensagem nao ensina a dispensa (`@no-domain`): %q", msg)
	}
}

// A dispensa COM RAZAO vale — e e' o que separa decisao de esquecimento.
func TestAWaiverWithAWrittenReasonSilencesTheGate(t *testing.T) {
	t.Run("DMDCD-B05: A waiver with a written reason silences the gate", func(t *testing.T) {})
	spec := "# U\n\n<!-- @no-domain: recebe so props tipadas do proprio codigo -->\n\n## Regras\n\n### UUUUU-B01 — algo\n"
	if v, msg := rodaDominio(t, spec); v != Skip && v != Pass {
		t.Errorf("`@no-domain` com razao deveria dispensar, veio %v: %s", v, msg)
	}
}

// Marcador NU nao dispensa — mesmo padrao do `@no-rule`. Um marcador sem razao vira um
// jeito silencioso de calar o gate, e some o rastro de que houve decisao.
func TestABareWaiverWithNoReasonDoesNotWaive(t *testing.T) {
	t.Run("DMDCD-I01: A bare waiver, with no reason, does not waive", func(t *testing.T) {})
	spec := "# U\n\n<!-- @no-domain -->\n\n## Regras\n\n### UUUUU-B01 — algo\n"
	if v, _ := rodaDominio(t, spec); v == Skip || v == Pass {
		t.Error("marcador sem razao nao pode dispensar")
	}
}

// A régua aqui é a PRESENÇA da declaração, não o acerto dela. Se o conjunto aceito
// corresponde ao domínio real é julgamento — e julgamento é de outra classe de gate.
// Cobrar isso aqui faria o gate reprovar por um critério que ele não sabe medir.
func TestTheGateDoesNotJudgeWhetherTheDeclaredDomainIsCorrect(t *testing.T) {
	t.Run("DMDCD-X01: The gate does not judge whether the declared domain is correct", func(t *testing.T) {})
	// "qualquer texto" é largo demais para um mês, e o gate passa MESMO ASSIM.
	largoDemais := "# Spec\n\n" + cabecalhoDominio +
		"| `mês` | qualquer texto | — | o chamador, antes de montar a consulta |\n"
	if v, d := rodaDominio(t, largoDemais); v != Pass {
		t.Fatalf("o gate julgou o MÉRITO do domínio declarado, e não deveria: %s (%s)", v, d)
	}
}

// Esta camada lê TEXTO. Cruzar a declaração com a implementação é do gate relacional,
// que tem o mapa — fazer os dois aqui duplicaria a régua em dois lugares que divergiriam.
func TestTheGateDoesNotReadTheCodeToCheckTheValidationExists(t *testing.T) {
	t.Run("DMDCD-X02: The gate does not read the code to check the validation exists", func(t *testing.T) {})
	// A spec declara dono para toda entrada; nenhum código é lido, e o gate passa.
	completa := "# Spec\n\n" + cabecalhoDominio +
		"| `chave` | texto não-vazio | `__proto__` | a interface de cadastro (KVEDX-V02) |\n"
	if v, d := rodaDominio(t, completa); v != Pass {
		t.Fatalf("o gate foi além do texto da spec: %s (%s)", v, d)
	}
}
