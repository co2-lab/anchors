package gate

// Item em aberto dá PENDING, não FAIL: declarar o que a spec ainda não decidiu é o
// comportamento que este gate existe para produzir, e reprová-lo como defeito inverte o
// incentivo — o caminho mais curto para o verde vira apagar a seção.

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func rodaAberto(t *testing.T, content string) (Verdict, string) {
	t.Helper()
	return checkOpenQuestions(content, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
}

// A distinção que decide se o gate é útil ou ritual: quem NÃO abriu a seção não é
// cobrado (specs simples não têm o que declarar); quem ABRIU precisa fechar.
func TestOpenQuestionsSoCobraQuemAbriu(t *testing.T) {
	semSeção := `# Spec — Cálculo de parcelas

## Regras

### PARCX-B01 — divide o total em N parcelas
`
	// Sem a seção: PENDENTE, não Skip. A seção é `Default: true` no catálogo e o
	// `anchors new` já a emite fechada com `nenhuma` — então a ausência é dívida (spec
	// anterior à prática) ou apagamento, e nenhum sinal separa os dois. `Skip` diria "não
	// se aplica", e um gate bloqueante com ✓0 ✗0 ~586 parece vigiar e não vigia.
	if v, msg := rodaAberto(t, semSeção); v != Pending {
		t.Fatalf("spec sem a seção deveria ser Pending (dívida nomeada), foi %s: %s", v, msg)
	}

	comPergunta := semSeção + `
## Decisões em aberto

- Arredondamento da última parcela: sobra vai na primeira ou na última? Contabilidade
  não respondeu.
`
	v, d := rodaAberto(t, comPergunta)
	if v != Pending {
		t.Fatalf("spec com pergunta aberta deveria reprovar, foi %s", v)
	}
	if !strings.Contains(d, "Arredondamento") {
		t.Errorf("a mensagem não mostra QUAL pergunta: %s", d)
	}
}

// O fechamento honesto: afirmar que se olhou e não há dúvida vale, e é diferente de
// omitir a seção. É o opt-out explícito do CONCEPT §5.1.
func TestOpenQuestionsFechamentoHonesto(t *testing.T) {
	for _, fecho := range []string{"nenhuma", "Nenhuma.", "- nenhuma", "none", "N/A", "sem pendências", "—"} {
		t.Run(fecho, func(t *testing.T) {
			spec := "# Spec\n\n## Decisões em aberto\n\n" + fecho + "\n"
			if v, d := rodaAberto(t, spec); v != Pass {
				t.Fatalf("fechamento %q deveria passar, foi %s (%s)", fecho, v, d)
			}
		})
	}
}

// A pergunta respondida vira REGRA; o item fica marcado como resolvido em vez de sumir.
// O rastro tem valor: mostra que a decisão foi tomada, não esquecida.
func TestOpenQuestionsItemResolvidoNaoBloqueia(t *testing.T) {
	spec := `# Spec

## Decisões em aberto

- [x] Arredondamento da última parcela → decidido: sobra na última (virou ` + "`PARCX-R03`" + `)
- ~~Moeda estrangeira~~ → fora de escopo nesta versão
`
	if v, d := rodaAberto(t, spec); v != Pass {
		t.Fatalf("itens resolvidos não deveriam bloquear, foi %s (%s)", v, d)
	}

	misto := spec + "- Fuso horário do vencimento: UTC ou local do usuário?\n"
	v, d := rodaAberto(t, misto)
	if v != Pending {
		t.Fatalf("um item ainda aberto deveria reprovar, foi %s", v)
	}
	if !strings.Contains(d, "1 decisão") {
		t.Errorf("deveria contar só o item ABERTO (1), não os resolvidos: %s", d)
	}
}

// Prosa que explica a seção não é pendência — senão o texto de abertura viraria um item
// fantasma e o autor aprenderia a não escrever nada, que é o oposto do objetivo.
func TestOpenQuestionsProsaNaoEhItem(t *testing.T) {
	spec := `# Spec

## Decisões em aberto

Esta seção lista o que a spec ainda não decide. Quando a resposta vier, ela deve virar
uma regra com código, e o item é marcado como resolvido.

nenhuma
`
	if v, d := rodaAberto(t, spec); v != Pass {
		t.Fatalf("prosa explicativa não é pendência, foi %s (%s)", v, d)
	}
}

// A seção também vale como TABELA — formato comum quando a pergunta tem dono e prazo.
// Cabeçalho e separador não são itens.
func TestOpenQuestionsTabela(t *testing.T) {
	vazia := `# Spec

## Decisões em aberto

| Pergunta | Quem decide | Vira |
| --- | --- | --- |
`
	if v, d := rodaAberto(t, vazia); v != Pass {
		t.Fatalf("tabela só com cabeçalho não tem pendência, foi %s (%s)", v, d)
	}

	comLinha := vazia + "| Fuso do vencimento: UTC ou local? | Produto | `PARCX-R04` |\n"
	v, d := rodaAberto(t, comLinha)
	if v != Pending {
		t.Fatalf("linha de tabela é pendência, foi %s", v)
	}
	if !strings.Contains(d, "Fuso") {
		t.Errorf("não mostrou a pergunta da tabela: %s", d)
	}
}

// A seção termina no próximo cabeçalho: pendência não pode vazar para as seções
// seguintes, nem regras seguintes serem lidas como pendência.
func TestOpenQuestionsRespeitaFronteiraDaSecao(t *testing.T) {
	spec := `# Spec

## Decisões em aberto

nenhuma

## Regras

- ` + "`PARCX-B01`" + ` — divide o total
- ` + "`PARCX-B02`" + ` — arredonda na última
`
	if v, d := rodaAberto(t, spec); v != Pass {
		t.Fatalf("regras da seção seguinte não são pendências, foi %s (%s)", v, d)
	}
}

// Variações de escrita não podem decidir o veredito — reprovar por causa de um acento
// ensinaria o autor a fugir da seção.
func TestOpenQuestionsAceitaVariacoesDoTitulo(t *testing.T) {
	títulos := []string{
		"## Decisões em aberto", "## Decisoes em aberto", "### Decisão em aberto",
		"## Questões em aberto", "## Open Questions", "## Em aberto",
		"## Pendências de decisão",
	}
	for _, tt := range títulos {
		t.Run(tt, func(t *testing.T) {
			spec := "# Spec\n\n" + tt + "\n\n- Pergunta que ninguém respondeu\n"
			if v, d := rodaAberto(t, spec); v != Pending {
				t.Fatalf("título %q não foi reconhecido (veredito %s: %s)", tt, v, d)
			}
		})
	}
}

// O gate é da SPEC: é ela que decide. Cobrar isso de código ou teste seria pedir que
// implemente resolva a ambiguidade — exatamente o chute que se quer evitar.
func TestOpenQuestionsSoSpec(t *testing.T) {
	spec := "# X\n\n## Decisões em aberto\n\n- pergunta\n"
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindTest, mapx.KindFeature} {
		t.Run(string(k), func(t *testing.T) {
			v, _ := checkOpenQuestions(spec, mapx.Node{Kind: k}, "", nil, nil)
			if v != Skip {
				t.Fatalf("kind %s deveria ser Skip, foi %s", k, v)
			}
		})
	}
}
