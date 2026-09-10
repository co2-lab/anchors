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
func TestDomainExigeDono(t *testing.T) {
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
func TestDomainNaoRespostaNaoEhDono(t *testing.T) {
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

// Ausência da seção não é falha: exigir de toda spec vira ritual, e unidade sem entrada
// externa não tem domínio a declarar. O gate cobra quem ABRIU.
func TestDomainAusenciaNaoEhFalha(t *testing.T) {
	if v, _ := rodaDominio(t, "# Spec\n\n## Regras\n\n### AAAAX-B01 — x\n"); v != Skip {
		t.Fatal("spec sem a seção deveria ser Skip")
	}
}

// Seção aberta e vazia é pior que ausente: AFIRMA que se olhou e não se achou nada.
func TestDomainSecaoVaziaReprova(t *testing.T) {
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
func TestDomainProsaNaoEhEntrada(t *testing.T) {
	spec := "# Spec\n\n## Domínio\n\nEsta unidade recebe o histórico já carregado.\n\n" +
		"| Entrada | Aceita | Fora do domínio | Quem garante |\n| --- | --- | --- | --- |\n" +
		"| `versões` | de UMA chave | lista multi-chave | o chamador (`RDMDX-B03`) |\n"
	if v, d := rodaDominio(t, spec); v != Pass {
		t.Fatalf("prosa não é entrada, foi %s (%s)", v, d)
	}
}

func TestDomainSoSpec(t *testing.T) {
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
func TestDominioTeriaPegadoOCasoReal(t *testing.T) {
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
